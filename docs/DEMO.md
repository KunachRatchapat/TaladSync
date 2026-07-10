# DEMO — คู่มือรัน + สคริปต์นำเสนอ

## 0. ต้องมีอะไรบ้าง
- PostgreSQL 17 รันอยู่ (Windows service เปิดเองตอนบูต) · database ชื่อ `taladsync`
- Go 1.24 · Node.js (npm)
- (เฉพาะ concurrent test) Docker Desktop เปิดอยู่

---

## 1. รัน (2 terminal)

```powershell
# Terminal 1 — Backend
cd D:\taladsync\backend
go run ./cmd/api
#   → http://localhost:8080/health  ควรได้ {"status":"ok"}
#   → seed อัตโนมัติถ้า DB ว่าง

# Terminal 2 — Frontend
cd D:\taladsync\frontend
npm run dev
#   → เปิด http://localhost:5173
```

> ต้องรัน **ทั้งสอง** พร้อมกัน · frontend เรียก API ที่ `http://localhost:8080`

---

## 2. รีเซ็ตข้อมูลให้สะอาด (ก่อนเดโม)

ถ้าเคยเทสจนแผงถูกจอง/แจ้งหยุดไปแล้ว อยากได้ข้อมูลสดใหม่:

```powershell
# หยุด backend ก่อน (Ctrl+C ใน Terminal 1) แล้ว:
$env:PGPASSWORD = "<DB_PASSWORD จาก backend/.env>"
& "C:\Program Files\PostgreSQL\17\bin\psql" -h localhost -U postgres -d taladsync `
  -c "TRUNCATE audit_logs, daily_bookings, stalls, users RESTART IDENTITY CASCADE;"

# แล้วรัน backend ใหม่ → seed สดอัตโนมัติ
go run ./cmd/api
```

หรือถ้าใช้ Docker: `docker compose down -v; docker compose up --build`

---

## 3. สคริปต์เดโม (เดินเรื่องตาม 4 หน้า)

ใช้ **ปุ่มสลับหน้ามุมล่าง (DevNav)** สลับได้ไว — ทุกหน้าอยู่ที่ `localhost:5173`

| # | หน้า | เดโมอะไร | จุดพูด |
|---|------|---------|--------|
| 1 | **ขาจร** `/` | เลือกแผงว่าง → กรอกชื่อร้าน → จอง → ขึ้นการ์ด "จองแล้ว" | จองแล้วจองซ้ำไม่ได้ (1 แผง/วัน), error ภาษาไทย |
| 2 | **เจ้าหน้าที่** `/staff` | เช็คชื่อ มา/ไม่มา → progress bar เดินหน้า → "เช็คครบแล้ว" | optimistic UI (กดปุ๊บเปลี่ยนทันที), summary นับ client-side |
| 3 | **แม่ค้าประจำ** `/regular` | เห็นแผงตัวเอง (A-01) → แจ้งหยุด → modal ยืนยัน → ปล่อยเข้าพูล | แจ้งหยุดแล้วขาจรจองต่อได้ (ลองสลับไปหน้าขาจรดู A-01 โผล่เป็นแผงว่าง) |
| 4 | **ค้นหาร้าน** `/directory` | พิมพ์ค้นหา (debounce) → กรองผล · public ไม่ต้อง login | ไม่มี PII, search-as-you-type |

**เส้นเรื่องที่ลื่น:** ขาจรจอง (1) → staff เช็คชื่อ (2) → แม่ค้าแจ้งหยุด (3) → ขาจรกลับมาเห็นแผงว่างเพิ่ม (1) → ลูกค้าค้นหาร้าน (4)

---

## 4. ⭐ Killer demo — พิสูจน์กัน race (ถ้ามีเวลา + Docker เปิด)

```powershell
cd D:\taladsync\backend
go test -tags=integration ./internal/usecase/... -run Concurrent -count=1 -v
```
50 goroutines จองแผงเดียวพร้อมกัน → assert เหลือผู้ชนะ **1 คน** (ที่เหลือ 409)
👉 อธิบายว่าใช้ **conditional atomic UPDATE** + partial unique index (ดู [DATABASE.md](DATABASE.md))

---

## 5. แก้ปัญหาเฉพาะหน้า
| อาการ | แก้ |
|------|-----|
| frontend เปิดไม่ได้ | ยังไม่ได้ `npm run dev` |
| จอง/หน้าเพี้ยน 404 | ยังไม่ได้ restart backend หลังแก้โค้ด → รัน `go run ./cmd/api` ใหม่ |
| port 8080/5173 ชน | มี process เก่าค้าง → ปิดก่อน |
| จองไม่ได้ ขึ้น "จองครบแล้ว" | ขาจรคนนี้จองไปแล้ววันนี้ → รีเซ็ต DB (ข้อ 2) |
