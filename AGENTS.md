# AGENTS.md — TaladSync

ระบบจัดการแผงตลาดสดรายวัน แก้ pain point: การแย่งจองแผงขาจร (concurrency) + การเช็คชื่อหน้างาน (roll-call) ว่าใครมา/ไม่มาตรงกับระบบ

> **เก็บเงิน 20 บาท:** เจ้าหน้าที่เก็บสดหน้างานเอง — **ไม่ track ในระบบ MVP** (ย้ายไป Phase 2 ถ้าต้องการ financial summary)

## เอกสารอ้างอิง (อ่านก่อนเริ่ม)

- [docs/SYSTEM_DESIGN.md](docs/SYSTEM_DESIGN.md) — ภาพรวม, stack, schema (attendance roll-call MVP)
- [docs/BOOKING_FLOW.md](docs/BOOKING_FLOW.md) — state machine, transaction strategy, concurrent test

---

## Implementation Checkpoint (อัปเดตล่าสุด)

**สถานะรวม:** Backend — **Phase 7 done** (concurrent booking integration test). **Session ถัดไป: Phase 8 (docker-compose)**

### Next session (อ่านก่อนเริ่ม Phase 8)

1. **Commit + push** Phase 7 (ถ้ายังไม่ push)
2. **เริ่ม Phase 8** — `docker-compose.yml` + README polish
3. Optional polish: seed `daily_bookings` สำหรับวันนี้เมื่อข้ามวัน (dev UX)

### Incremental roadmap

| Phase | สถานะ | เป้าหมาย |
|-------|--------|----------|
| 0 | **DONE** | Health API — `cmd/api/main.go` + `GET /health` |
| 1 | **DONE** | Domain models (4 tables) + Postgres connect + migrate |
| 2 | **DONE** | Seed (users 3 roles, stalls, daily_bookings) |
| 3 | **DONE** | Auth mock-login + JWT middleware + `GET /api/auth/me` |
| 4 | **DONE** | Booking — list/book, 409, `date?` body, partial unique index (1 vendor/วัน), unit tests |
| 5 | **DONE** | Vendor leave + staff checklist/attendance/summary + audit (LEAVE/CHECKIN/NOSHOW) |
| 6 | **DONE** | Directory — `GET /api/directory?q=` public, ไม่มี PII |
| 7 | **DONE** | Concurrent test — 50 goroutines, 1 winner (`-tags=integration`) |
| 8 | TODO | DevOps (docker-compose) |

### DONE (มีไฟล์แล้ว)

| Layer | ไฟล์ | หมายเหตุ |
|---|---|---|
| **entrypoint** | `backend/cmd/api/main.go` | Fiber routes — **middleware ต่อ route** (ดูด้านล่าง) |
| **config** | `backend/internal/config/config.go` | env loader (APP_PORT, DB_*, JWT_SECRET) |
| **domain** | `backend/internal/domain/models.go`, `errors.go` | entities, audit constants, inputs, `DirectoryEntry` DTO |
| **database** | `backend/internal/database/database.go`, `seed.go` | Connect, Migrate (+ partial unique index), SeedIfEmpty |
| **auth** | `backend/internal/auth/jwt.go` | TokenService Sign/Parse (HS256) |
| **repository** | `user.go`, `booking.go`, `audit.go` | user; booking CRUD+leave+attendance+SearchDirectory; audit Create |
| **usecase** | `auth.go`, `booking.go`, `vendor.go`, `staff.go`, `directory.go`, `*_test.go`, `booking_concurrent_test.go` | MockLogin; BookStall; Leave; MarkAttendance; Search + unit/integration tests |
| **handler** | `auth.go`, `booking.go`, `vendor.go`, `staff.go`, `directory.go`, `date.go`, `errors.go` | all MVP endpoints + error mapping |
| **middleware** | `auth.go`, `role.go` | JWT Auth, RequireRoles |
| **timeutil** | `backend/internal/pkg/timeutil/market_date.go` | `Today()`, `Now()` Asia/Bangkok |
| **env** | `backend/.env.example` | template สำหรับ dev |
| **deps** | `backend/go.mod`, `go.sum` | Fiber, godotenv, GORM, postgres, jwt/v5, testcontainers-go |

### TODO (Phase 8+)

1. **DevOps** — `docker-compose.yml`, README polish (Phase 8)
2. **Seed daily bookings for today** — dev UX ข้ามวัน (optional)
3. **BOOK audit** — retroactive `audit_logs` ตอนจอง (optional polish)

---

## Route wiring (`cmd/api/main.go`) — สำคัญ

**อย่าใช้ `api.Group("")` หลายครั้งพร้อม `RequireRoles` คนละ role** — Fiber รวม middleware ของ group prefix ว่าง ทำให้ staff route โดน `RequireRoles(WALKIN)` → **403 forbidden**

รูปแบบที่ใช้ตอนนี้:

```go
authMW := middleware.Auth(tokenSvc)

api.Get("/directory", directoryHandler.Search)  // public — ไม่มี authMW

api.Get("/stalls/available", authMW, middleware.RequireRoles(domain.RoleWalkin), ...)
api.Post("/bookings", authMW, middleware.RequireRoles(domain.RoleWalkin), ...)
api.Post("/vendor/leave", authMW, middleware.RequireRoles(domain.RoleRegular), ...)
api.Get("/staff/checklist", authMW, middleware.RequireRoles(domain.RoleStaff), ...)
api.Get("/staff/summary", authMW, middleware.RequireRoles(domain.RoleStaff), ...)
api.Post("/bookings/:id/attendance", authMW, middleware.RequireRoles(domain.RoleStaff), ...)
```

`/api/auth/me` ใช้ group แยก `api.Group("/auth", authMW)` — ไม่มี RequireRoles

## Tech Stack

- **Backend:** Go 1.24 + Fiber v2 + GORM + PostgreSQL 17 (`localhost:5432`)
- **Frontend (ยังไม่เริ่ม):** Vue 3 + Tailwind + Pinia
- **Auth:** JWT (MVP mock login), production = LINE LIFF
- **Project root:** `D:\taladsync`

## Architecture: Modular Monolith + Clean Architecture

**Dependency rule:** ลูกศรชี้เข้าด้านใน `handler → usecase → domain` ; `repository` implement interface ของ `usecase`

```
cmd/api/main.go          ← DONE (health + auth + booking + vendor + staff + directory)
internal/
  config/                ← DONE
  domain/                ← DONE (models, errors, inputs, DirectoryEntry)
  database/              ← DONE (connect + migrate + seed + idx_vendor_one_booked_per_day)
  pkg/timeutil/          ← DONE (Today, Now, MarketDate)
  auth/                  ← DONE (JWT)
  repository/            ← DONE (user, booking, audit)
  usecase/               ← DONE (auth, booking, vendor, staff, directory + tests)
  handler/               ← DONE (auth, booking, vendor, staff, directory, date, errors)
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
- 1 ขาจร 1 แผง/วัน: partial unique index `idx_vendor_one_booked_per_day` บน `(vendor_id, market_date) WHERE status='BOOKED'`

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
- Audit: `CHECKIN`, `NOSHOW`, `LEAVE` → `audit_logs` (MVP; `BOOK` ยังไม่เขียน)
- `parseMarketDate` ใน `handler/date.go` — query/body `date` ว่าง = วันนี้ Bangkok; ใช้ร่วม list/book/leave/checklist

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
| GET | `/api/directory?q=` | public | **DONE** | ค้นหาร้าน (ไม่มี PII) — `?date=` optional |

## การรัน (dev)

```powershell
copy backend\.env.example backend\.env

# สร้าง database (ครั้งแรก)
createdb -U postgres taladsync

# รัน API (seed อัตโนมัติถ้า users ว่าง)
cd backend
go run ./cmd/api
# GET http://localhost:8080/health → {"status":"ok"}
```

**หมายเหตุ dev:** seed สร้าง `daily_bookings` แค่วันที่ seed ครั้งแรก — query `?date=` / body `date` ต้องตรงวันใน DB (หรือ reset DB / รอ seed-today polish)

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

ใช้ `YYYY-MM-DD` ตรงกับ `market_date` ใน DB (ดูด้วย `SELECT DISTINCT market_date FROM daily_bookings`)

```
# REGULAR แจ้งหยุด (stallId = แผง A-01..A-12, ไม่ใช่ daily_bookings.id)
POST /api/auth/mock-login  { "role": "REGULAR" }
POST /api/vendor/leave     { "stallId": 1, "date": "YYYY-MM-DD" }

# WALKIN ดู pool หลัง leave
POST /api/auth/mock-login  { "role": "WALKIN" }
GET  /api/stalls/available?date=YYYY-MM-DD

# STAFF roll-call (:id = daily_bookings.id จาก checklist ไม่ใช่ stallId)
POST /api/auth/mock-login  { "role": "STAFF" }
GET  /api/staff/checklist?date=YYYY-MM-DD
POST /api/bookings/:id/attendance  { "present": true }
GET  /api/staff/summary?date=YYYY-MM-DD
```

## Phase 6 smoke test

ใช้ `YYYY-MM-DD` ตรงกับ `market_date` ใน DB (ดูด้วย `SELECT DISTINCT market_date FROM daily_bookings`) — **ไม่ต้อง JWT**

```
GET /api/directory?q=ของสด&date=YYYY-MM-DD
GET /api/directory?date=YYYY-MM-DD          # q ว่าง = list ทุกร้านวันนั้น
GET /api/directory?q=ร้านA-01&date=YYYY-MM-DD
```

Response: `{"results":[{"stallCode","zone","shopName","productType","posX","posY"}]}`

Unit test: `go test ./internal/usecase/... -v -run Search_`

## Phase 7 concurrent test (integration)

พิสูจน์ conditional UPDATE กัน double-booking — **ต้อง Docker Desktop รันอยู่**

```powershell
cd backend

# unit tests ปกติ (ไม่ต้อง Docker)
go test ./internal/usecase/... -v

# killer demo — 50 goroutines จองแผงเดียว, assert 1 winner
go test -tags=integration ./internal/usecase/... -v -run Concurrent -count=1
```

ถ้า Docker ไม่พร้อม test จะ **skip** (ไม่ fail) พร้อมข้อความให้เปิด Docker Desktop

## Phase 2 (ยังไม่ทำ)

- LINE LIFF / Push จริง (Flash Drop)
- โควตาวันหยุด, แจ้งหยุดฉุกเฉิน 2 ครั้ง/เดือน, ลาพักยาว+อนุมัติ
- Force Absent + ปล่อยแผง no-show อัตโนมัติ
- แจ้งมาช้า (late arrival lock)
- **Financial tracking** (เก็บเงิน 20 บ. ในระบบ + สรุปยอดส่งคลัง) — ตัดออกจาก MVP แล้ว
- national_id encryption เต็มรูปแบบ (PDPA)
- Frontend Vue 3 (vendor grab UI, staff checklist, customer map)
