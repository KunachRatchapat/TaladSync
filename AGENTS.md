# AGENTS.md — TaladSync

ระบบจัดการแผงตลาดสดรายวัน แก้ pain point: การแย่งจองแผงขาจร (concurrency) + การเช็คชื่อหน้างาน (roll-call) ว่าใครมา/ไม่มาตรงกับระบบ

> **เก็บเงิน 20 บาท:** เจ้าหน้าที่เก็บสดหน้างานเอง — **ไม่ track ในระบบ MVP** (ย้ายไป Phase 2 ถ้าต้องการ financial summary)

## เอกสารอ้างอิง (อ่านก่อนเริ่ม)

- [docs/SYSTEM_DESIGN.md](docs/SYSTEM_DESIGN.md) — ภาพรวม, stack, schema (attendance roll-call MVP)
- [docs/BOOKING_FLOW.md](docs/BOOKING_FLOW.md) — state machine, transaction strategy, concurrent test

---

## Implementation Checkpoint (อัปเดตล่าสุด)

**สถานะรวม:** Backend — **Phase 3 done** (auth: mock-login, JWT middleware, `/api/auth/me`). ถัดไป Phase 4 booking core

### Incremental roadmap

| Phase | สถานะ | เป้าหมาย |
|-------|--------|----------|
| 0 | **DONE** | Health API — `cmd/api/main.go` + `GET /health` |
| 1 | **DONE** | Domain models (4 tables) + Postgres connect + migrate |
| 2 | **DONE** | Seed (users 3 roles, stalls, daily_bookings) |
| 3 | **DONE** | Auth mock-login + JWT middleware + `GET /api/auth/me` |
| 4 | TODO | Booking core (list + book + 409) |
| 5 | TODO | Vendor leave + staff attendance |
| 6 | TODO | Directory |
| 7 | TODO | Concurrent test |
| 8 | TODO | DevOps (docker-compose) |

### DONE (มีไฟล์แล้ว)

| Layer | ไฟล์ | หมายเหตุ |
|---|---|---|
| **entrypoint** | `backend/cmd/api/main.go` | Fiber + DB + auth routes |
| **config** | `backend/internal/config/config.go` | env loader (APP_PORT, DB_*, JWT_SECRET) |
| **domain** | `backend/internal/domain/models.go`, `errors.go` | User, Stall, DailyBooking, AuditLog + business errors |
| **database** | `backend/internal/database/database.go`, `seed.go` | Connect, Migrate, SeedIfEmpty |
| **auth** | `backend/internal/auth/jwt.go` | TokenService Sign/Parse (HS256) |
| **repository** | `backend/internal/repository/user.go` | FindByRole, FindByID |
| **usecase** | `backend/internal/usecase/auth.go` | MockLogin, GetMe |
| **handler** | `backend/internal/handler/auth.go`, `errors.go` | mock-login, /me, error mapping |
| **middleware** | `backend/internal/middleware/auth.go`, `role.go` | JWT Auth, RequireRoles |
| **timeutil** | `backend/internal/pkg/timeutil/market_date.go` | `Today()` Asia/Bangkok |
| **env** | `backend/.env.example` | template สำหรับ dev |
| **deps** | `backend/go.mod`, `go.sum` | Fiber, godotenv, GORM, postgres, jwt/v5 |

### TODO (Phase 4 — ทำต่อจากตรงนี้)

1. **Booking core** — repository, usecase, handler, concurrent test
2. **Vendor + Staff + Directory** — handlers ที่เหลือ
3. **DevOps** — `docker-compose.yml`, README dev setup

---

## Tech Stack

- **Backend:** Go 1.24 + Fiber v2 + GORM + PostgreSQL 17 (`localhost:5432`)
- **Frontend (ยังไม่เริ่ม):** Vue 3 + Tailwind + Pinia
- **Auth:** JWT (MVP mock login), production = LINE LIFF
- **Project root:** `D:\taladsync`

## Architecture: Modular Monolith + Clean Architecture

**Dependency rule:** ลูกศรชี้เข้าด้านใน `handler → usecase → domain` ; `repository` implement interface ของ `usecase`

```
cmd/api/main.go          ← DONE (health + auth routes)
internal/
  config/                ← DONE
  domain/                ← DONE
  database/              ← DONE (connect + migrate + seed)
  pkg/timeutil/          ← DONE
  auth/                  ← DONE (JWT)
  repository/            ← DONE (user); TODO booking Phase 4
  usecase/               ← DONE (auth); TODO booking Phase 4
  handler/               ← DONE (auth); TODO booking Phase 4
  middleware/            ← DONE (Auth, RequireRoles)
```

**หมายเหตุ:** domain model มี gorm tag (pragmatic); testability มาจาก usecase พึ่ง interface ไม่ใช่ GORM ตรง ๆ

## Core Concurrency Rule (สำคัญที่สุด)

จองแผงใช้ **conditional atomic UPDATE** ไม่ใช่ lock ใน app:

```sql
UPDATE daily_bookings
SET status = 'BOOKED', vendor_id = ?, shop_name = ?, product_type = ?, booked_at = now()
WHERE stall_id = ? AND market_date = ? AND status = 'AVAILABLE';
```

- `RowsAffected = 1` → สำเร็จ
- `RowsAffected = 0` → `ErrStallNotAvailable` (HTTP 409)
- **ห้าม** ทำ I/O ภายนอก (LINE push) ใน transaction ที่ถือ row lock
- Double-booking กันสองชั้น: conditional `WHERE` + `UNIQUE(stall_id, market_date)`

## State Machine

```
Stall status (daily_bookings.status):
  EXPECTED  → AVAILABLE  (แม่ค้าแจ้งหยุด / cron)
  AVAILABLE → BOOKED     (ขาจรจอง — conditional UPDATE)
  EXPECTED  → ABSENT     (Force Absent — Phase 2)

Attendance (daily_bookings.attendance) — แยก lifecycle:
  PENDING → PRESENT  (staff เช็คว่ามา)
  PENDING → NOSHOW   (staff เช็คว่าไม่มา — บันทึก+audit เท่านั้น ไม่เปลี่ยน status)
```

## Schema สำคัญ (`DailyBooking`)

```go
Status           BookingStatus    // EXPECTED | AVAILABLE | BOOKED | ABSENT
Attendance       AttendanceStatus // PENDING | PRESENT | NOSHOW
CheckedByStaffID *uint
CheckedAt        *time.Time
// ลบแล้ว: PaymentStatus, PaidByStaffID, PaidAt
```

## Conventions

- เวลาใช้ `Asia/Bangkok` เสมอ; `market_date` เป็น type `date` (`internal/pkg/timeutil`)
- Business error อยู่ใน `internal/domain/errors.go`; **handler** เป็นที่เดียวที่ map → HTTP status
- `national_id` = PII: ห้าม return ออก public API (`json:"-"`); production ต้อง encrypt
- Idempotent: การเช็คชื่อ (`MarkAttendance`) และการจอง (middleware `Idempotency-Key`) ต้องกดซ้ำได้ไม่พัง
- Audit: `CHECKIN`, `NOSHOW`, `BOOK`, `LEAVE` → `audit_logs`

## API (MVP — scope ปัจจุบัน)

| Method | Endpoint | Role | คำอธิบาย |
|--------|----------|------|----------|
| POST | `/api/auth/mock-login` | public | `{role}` → `{token, user}` |
| GET | `/api/stalls/available?date=` | WALKIN | แผงว่างใน pool |
| POST | `/api/bookings` | WALKIN | `{stallId, shopName, productType}` → 409 ถ้าเต็ม |
| POST | `/api/vendor/leave` | REGULAR | `{stallId, date?}` แจ้งหยุด |
| GET | `/api/staff/checklist?date=` | STAFF | roll-call ทุกแผงที่มีคน (EXPECTED+BOOKED) |
| POST | `/api/bookings/:id/attendance` | STAFF | `{present: true/false}` เช็คมา/ไม่มา (idempotent) |
| GET | `/api/staff/summary?date=` | STAFF | `{present, noShow, pending, total}` |
| GET | `/api/directory?q=` | public | ค้นหาร้าน (ไม่มี PII) |
| GET | `/health` | public | health check |

## การรัน (dev)

```powershell

copy backend\.env.example backend\.env


# 2. สร้าง database (ครั้งแรก)
createdb -U postgres taladsync

# 3. รัน API (seed อัตโนมัติถ้า users ว่าง)
cd backend
go run ./cmd/api
# GET http://localhost:8080/health  →  {"status":"ok"}
# log: "seed complete: 3 users, 20 stalls, 20 daily_bookings ..."
```

## Mock Login (ทดสอบ)

```json
POST /api/auth/mock-login
{ "role": "WALKIN" }
```

Role ที่ใช้ได้: `REGULAR` | `WALKIN` | `STAFF`  
Header: `Authorization: Bearer <token>`

## Phase 2 (ยังไม่ทำ)

- LINE LIFF / Push จริง (Flash Drop)
- โควตาวันหยุด, แจ้งหยุดฉุกเฉิน 2 ครั้ง/เดือน, ลาพักยาว+อนุมัติ
- Force Absent + ปล่อยแผง no-show อัตโนมัติ
- แจ้งมาช้า (late arrival lock)
- **Financial tracking** (เก็บเงิน 20 บ. ในระบบ + สรุปยอดส่งคลัง) — ตัดออกจาก MVP แล้ว
- national_id encryption เต็มรูปแบบ (PDPA)
- Frontend Vue 3 (vendor grab UI, staff checklist, customer map)
