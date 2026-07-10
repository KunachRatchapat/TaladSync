# Database — TaladSync

PostgreSQL 17 · ORM = GORM (AutoMigrate สร้างตารางจาก Go struct) · เวลา `Asia/Bangkok` เสมอ
Schema อยู่ที่ [`backend/internal/domain/models.go`](../backend/internal/domain/models.go) · migration ที่ [`database/database.go`](../backend/internal/database/database.go)

---

## 4 ตาราง

### 1. `users` — คน 3 บทบาท
| คอลัมน์ | หมายเหตุ |
|--------|---------|
| `role` | `REGULAR` (แม่ค้าประจำ, เจ้าของแผง) · `WALKIN` (ขาจร) · `STAFF` (เจ้าหน้าที่) |
| `national_id` | **PII** → `json:"-"` ไม่ส่งออก API (production ต้อง encrypt) |
| `line_user_id` | เผื่อ production ต่อ LINE LIFF |

### 2. `stalls` — แผง (ข้อมูล**ถาวร**)
| คอลัมน์ | หมายเหตุ |
|--------|---------|
| `code` | `A-01`..`A-12` (Zone A) · `P-01`..`P-08` (Pool) |
| `owner_vendor_id` | **เจ้าของถาวร** (nullable) — Zone A มีเจ้าของ (แม่ค้าประจำ), Pool ไม่มี (ว่างให้ขาจร) |
| `pos_x`, `pos_y` | ตำแหน่งบนแผนที่ (เผื่ออนาคต — ยังไม่ใช้) |

### 3. `daily_bookings` — ⭐ หัวใจของระบบ
**1 แถว = 1 แผง ต่อ 1 วัน** — สถานะ "รายวัน" อยู่ที่นี่ ไม่ใช่ที่ `stalls` (เพราะตลาด reset ทุกวัน)

| คอลัมน์ | หมายเหตุ |
|--------|---------|
| `stall_id` + `market_date` | คู่นี้ **UNIQUE** (1 แผง 1 แถว/วัน) |
| `status` | `EXPECTED` / `AVAILABLE` / `BOOKED` / `ABSENT` (ดู state machine ล่าง) |
| `vendor_id` | **ใครใช้แผงนี้วันนี้** (อาจต่างจาก owner! · nullable) |
| `shop_name`, `product_type` | ชื่อร้าน/สินค้าของวันนั้น |
| `attendance` | `PENDING` / `PRESENT` / `NOSHOW` — staff เช็คชื่อ (lifecycle แยกจาก status) |
| `checked_by_staff_id`, `checked_at`, `booked_at` | ร่องรอยเวลา |

### 4. `audit_logs` — บันทึกการกระทำ
`action` = `LEAVE` / `CHECKIN` / `NOSHOW` · เก็บว่า **ใคร** (`actor_user_id`) ทำอะไรกับ booking ไหน

---

## 🔑 จุดที่มักงงที่สุด: `owner_vendor_id` vs `vendor_id`

| | อยู่ที่ | ความหมาย | เปลี่ยนไหม |
|--|-------|---------|-----------|
| `stalls.owner_vendor_id` | ตาราง stalls | **เจ้าของถาวร** | ไม่เปลี่ยน |
| `daily_bookings.vendor_id` | ตาราง daily_bookings | **คนใช้แผ่งวันนี้** | เปลี่ยน/เป็น null ได้ |

**ตัวอย่าง A-01 ของ "แม่ค้าประจำ01" (owner_vendor_id = 1):**
| เหตุการณ์ | status | vendor_id |
|-----------|--------|-----------|
| ปกติ (จะมาขาย) | `EXPECTED` | 1 (เจ้าของ) |
| แจ้งหยุด | `AVAILABLE` | **null** (เคลียร์ทิ้ง) |
| ขาจรจองต่อ | `BOOKED` | 7 (ขาจร) |

➡️ นี่คือเหตุผลที่ API "แผงของฉัน" (`GET /api/vendor/mine`) ต้อง query ด้วย **`owner_vendor_id`** ไม่ใช่ `vendor_id` — ไม่งั้นพอแจ้งหยุดแล้วแผงจะหายจากรายการ (เพราะ vendor_id ถูกเคลียร์)

---

## State machine (สถานะแผงรายวัน)

```
status (daily_bookings.status):
  EXPECTED  ──แม่ค้าแจ้งหยุด──▶ AVAILABLE ──ขาจรจอง (conditional UPDATE)──▶ BOOKED

attendance (แยก lifecycle, staff เป็นคนเปลี่ยน):
  PENDING ──เช็คว่ามา──▶ PRESENT
  PENDING ──เช็คว่าไม่มา──▶ NOSHOW
```

---

## ⭐⭐ หัวใจ: กัน double-booking (concurrency)

**ปัญหา:** ขาจรหลายคนกดจองแผงเดียว "พร้อมกัน" → ต้องมีคนจองได้แค่ **1 คน**

**วิธีแก้ — conditional atomic UPDATE (ไม่ล็อกใน application):**
```sql
UPDATE daily_bookings
SET status = 'BOOKED', vendor_id = ?, shop_name = ?, ...
WHERE stall_id = ? AND market_date = ? AND status = 'AVAILABLE';
```
- `RowsAffected = 1` → **คนนี้ชนะ** (201)
- `RowsAffected = 0` → แผงถูกจองไปแล้ว → **409 Conflict**

DB รับประกันว่า `UPDATE` เป็น atomic ต่อ row → คำขอที่มาพร้อมกัน จะมีแค่ **คนเดียว**ที่เจอ `status='AVAILABLE'` ตอนอัปเดต ที่เหลือเจอ 0 rows

**กัน 2 ชั้น (defense in depth):**
1. `UNIQUE(stall_id, market_date)` — 1 แผง 1 แถว/วัน (กัน double-book แผงเดียวกัน)
2. **partial unique index** — `(vendor_id, market_date) WHERE status='BOOKED'` — กัน 1 ขาจร จองเกิน 1 แผง/วัน แม้ยิงหลาย request พร้อมกัน

**พิสูจน์:** integration test 50 goroutines จองแผงเดียว → assert เหลือผู้ชนะ 1 คน
([`booking_concurrent_test.go`](../backend/internal/usecase/booking_concurrent_test.go))
```powershell
cd backend
go test -tags=integration ./internal/usecase/... -run Concurrent -count=1 -v
```

**ทำไมไม่ทำ I/O นอก (เช่น LINE push) ใน transaction ที่ถือ row lock?** → กัน lock ค้างนาน (แยกไปทำหลัง commit)

---

## Migration & Seed
- `Migrate()` — AutoMigrate 4 ตาราง + `CREATE UNIQUE INDEX ... WHERE status='BOOKED'` (partial index สร้างเองเพราะ GORM tag ทำ partial ไม่ได้)
- `SeedIfEmpty()` — ใส่ข้อมูล dev เมื่อ `users` ว่าง: **12 REGULAR (คนละ 1 แผง A-01..A-12)** + 1 WALKIN + 1 STAFF + pool P-01..P-08
- `SeedDailyBookingsForToday()` — สร้าง `daily_bookings` ของ "วันนี้" ทุกครั้งที่ API start (แผงที่มีเจ้าของ = EXPECTED, pool = AVAILABLE)
