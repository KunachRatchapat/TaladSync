# AGENTS.md — TaladSync

ระบบจัดการแผงตลาดสดรายวัน แก้ pain point: การแย่งจองแผงขาจร (concurrency) + การเช็คชื่อหน้างาน (roll-call) ว่าใครมา/ไม่มาตรงกับระบบ

> **เก็บเงิน 20 บาท:** เจ้าหน้าที่เก็บสดหน้างานเอง — **ไม่ track ในระบบ MVP** (ย้ายไป Phase 2 ถ้าต้องการ financial summary)

## เอกสารอ้างอิง (อ่านก่อนเริ่ม)

- [docs/SYSTEM_DESIGN.md](docs/SYSTEM_DESIGN.md) — ภาพรวม, stack, schema (attendance roll-call MVP)
- [docs/BOOKING_FLOW.md](docs/BOOKING_FLOW.md) — state machine, transaction strategy, concurrent test

---

## Implementation Checkpoint (อัปเดตล่าสุด)

**สถานะรวม:** Backend — **Phase 5 done** (vendor leave + staff attendance + audit log). **Session ถัดไป: Phase 6 (directory)**

### Next session (อ่านก่อนเริ่ม Phase 6)

1. **Commit + push** Phase 4–5 changes (ถ้ายังไม่ push)
2. **เริ่ม Phase 6** — `GET /api/directory?q=`

### Incremental roadmap

| Phase | สถานะ | เป้าหมาย |
|-------|--------|----------|
| 0 | **DONE** | Health API — `cmd/api/main.go` + `GET /health` |
| 1 | **DONE** | Domain models (4 tables) + Postgres connect + migrate |
| 2 | **DONE** | Seed (users 3 roles, stalls, daily_bookings) |
| 3 | **DONE** | Auth mock-login + JWT middleware + `GET /api/auth/me` |
| 4 | **DONE** | Booking core — `GET /stalls/available`, `POST /bookings`, 409, partial unique index, usecase unit test |
| 5 | **DONE** | Vendor leave + staff checklist/attendance/summary + audit log |
| 6 | TODO | Directory |
| 6 | TODO | Directory |
| 7 | TODO | Concurrent test (testcontainers) |
| 8 | TODO | DevOps (docker-compose) |

### DONE (มีไฟล์แล้ว)

| Layer | ไฟล์ | หมายเหตุ |
|---|---|---|
| **entrypoint** | `backend/cmd/api/main.go` | Fiber + DB + auth + booking + vendor + staff routes |
| **config** | `backend/internal/config/config.go` | env loader (APP_PORT, DB_*, JWT_SECRET) |
| **domain** | `backend/internal/domain/models.go`, `errors.go` | entities, audit constants, Leave/MarkAttendance inputs |
| **database** | `backend/internal/database/database.go`, `seed.go` | Connect, Migrate (+ partial unique index), SeedIfEmpty |
| **auth** | `backend/internal/auth/jwt.go` | TokenService Sign/Parse (HS256) |
| **repository** | `user.go`, `booking.go`, `audit.go` | user; booking CRUD+leave+attendance; audit Create |
| **usecase** | `auth.go`, `booking.go`, `vendor.go`, `staff.go`, `*_test.go` | MockLogin; BookStall; Leave; MarkAttendance + tests |
| **handler** | `auth.go`, `booking.go`, `vendor.go`, `staff.go`, `date.go`, `errors.go` | all MVP endpoints + error mapping |
| **middleware** | `auth.go`, `role.go` | JWT Auth, RequireRoles |
| **timeutil** | `backend/internal/pkg/timeutil/market_date.go` | `Today()`, `Now()` Asia/Bangkok |
| **env** | `backend/.env.example` | template สำหรับ dev |
| **deps** | `backend/go.mod`, `go.sum` | Fiber, godotenv, GORM, postgres, jwt/v5 |

### TODO (Phase 6+)

1. **Directory** — `GET /api/directory?q=` (Phase 6)
2. **Concurrent test** — 50 goroutines, assert 1 winner (Phase 7)
3. **DevOps** — `docker-compose.yml`, README polish (Phase 8)
4. **Seed daily bookings for today** — dev UX ข้ามวัน (optional polish)

---

## Tech Stack

- **Backend:** Go 1.24 + Fiber v2 + GORM + PostgreSQL 17 (`localhost:5432`)
- **Frontend (ยังไม่เริ่ม):** Vue 3 + Tailwind + Pinia
- **Auth:** JWT (MVP mock login), production = LINE LIFF
- **Project root:** `D:\taladsync`

## Architecture: Modular Monolith + Clean Architecture

**Dependency rule:** ลูกศรชี้เข้าด้านใน `handler → usecase → domain` ; `repository` implement interface ของ `usecase`

```
cmd/api/main.go          ← DONE (health + auth + booking)
internal/
  config/                ← DONE
  domain/                ← DONE
  database/              ← DONE (connect + migrate + seed)
  pkg/timeutil/          ← DONE
  auth/                  ← DONE (JWT)
  repository/            ← DONE (user, booking)
  usecase/               ← DONE (auth, booking + unit test)
  handler/               ← DONE (auth, booking)
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

## API (MVP — implemented vs TODO)

| Method | Endpoint | Role | สถานะ | คำอธิบาย |
|--------|----------|------|--------|----------|
| GET | `/health` | public | **DONE** | health check |
| POST | `/api/auth/mock-login` | public | **DONE** | `{role}` → `{token, user}` |
| GET | `/api/auth/me` | any JWT | **DONE** | profile จาก token |
| GET | `/api/stalls/available?date=` | WALKIN | **DONE** | แผงว่างใน pool |
| POST | `/api/bookings` | WALKIN | **DONE** | `{stallId, shopName, productType, date?}` → 409 ถ้าเต็ม |
| POST | `/api/vendor/leave` | REGULAR | **DONE** | `{stallId, date?}` แจ้งหยุด |
| GET | `/api/staff/checklist?date=` | STAFF | **DONE** | roll-call ทุกแผงที่มีคน |
| POST | `/api/bookings/:id/attendance` | STAFF | **DONE** | `{present: true/false}` |
| GET | `/api/staff/summary?date=` | STAFF | **DONE** | `{present, noShow, pending, total}` |
| GET | `/api/directory?q=` | public | TODO | ค้นหาร้าน (ไม่มี PII) |

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

## Booking (ทดสอบ — WALKIN)

```
1. POST /api/auth/mock-login  { "role": "WALKIN" }
2. GET  /api/stalls/available?date=YYYY-MM-DD   (Bearer token)
3. POST /api/bookings  { "stallId": 13, "shopName": "ร้านทด", "productType": "ผลไม้", "date": "YYYY-MM-DD" }
```

Unit test: `go test ./internal/usecase/... -v`

## Phase 5 smoke test

```
# REGULAR leave
POST /api/auth/mock-login  { "role": "REGULAR" }
POST /api/vendor/leave     { "stallId": 1, "date": "YYYY-MM-DD" }

# STAFF roll-call
POST /api/auth/mock-login  { "role": "STAFF" }
GET  /api/staff/checklist?date=YYYY-MM-DD
POST /api/bookings/:id/attendance  { "present": true }
GET  /api/staff/summary?date=YYYY-MM-DD
```

## Phase 2 (ยังไม่ทำ)

- LINE LIFF / Push จริง (Flash Drop)
- โควตาวันหยุด, แจ้งหยุดฉุกเฉิน 2 ครั้ง/เดือน, ลาพักยาว+อนุมัติ
- Force Absent + ปล่อยแผง no-show อัตโนมัติ
- แจ้งมาช้า (late arrival lock)
- **Financial tracking** (เก็บเงิน 20 บ. ในระบบ + สรุปยอดส่งคลัง) — ตัดออกจาก MVP แล้ว
- national_id encryption เต็มรูปแบบ (PDPA)
- Frontend Vue 3 (vendor grab UI, staff checklist, customer map)
