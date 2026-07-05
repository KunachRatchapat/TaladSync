# AGENTS.md — TaladSync

ระบบจัดการแผงตลาดสดรายวัน แก้ pain point: การแย่งจองแผงขาจร (concurrency) + การเช็คชื่อหน้างาน (roll-call) ว่าใครมา/ไม่มาตรงกับระบบ

> **เก็บเงิน 20 บาท:** เจ้าหน้าที่เก็บสดหน้างานเอง — **ไม่ track ในระบบ MVP** (ย้ายไป Phase 2 ถ้าต้องการ financial summary)

## เอกสารอ้างอิง (อ่านก่อนเริ่ม)

- [docs/SYSTEM_DESIGN.md](docs/SYSTEM_DESIGN.md) — ภาพรวม, stack, schema (attendance roll-call MVP)
- [docs/BOOKING_FLOW.md](docs/BOOKING_FLOW.md) — state machine, transaction strategy, concurrent test

---

## Implementation Checkpoint (อัปเดตล่าสุด)

**สถานะรวม:** Backend MVP **Phase 0–8 DONE** (commit `ff85bba`). **Session ถัดไป: Frontend F1 (WALKIN)**

### Next session (อ่านก่อนเริ่ม Frontend)

1. อ่าน section **Frontend roadmap** ด้านล่าง — **ยังไม่ implement** จนกว่าจะคุย scope กับ user อีกครั้ง
2. Backend เตรียม: เพิ่ม **CORS** ใน `main.go` (หรือ Vite proxy ตอน dev)
3. รัน API ก่อนเทส UI: `docker compose up` หรือ `go run ./cmd/api`
4. ใช้ `?date=` / body `date` ตรง `market_date` ใน DB (ดู log seed)

### Phase 8 — DONE (commit แล้ว)

| ไฟล์ | งาน |
|------|-----|
| [`docker-compose.yml`](docker-compose.yml) | `postgres:17-alpine` + `api`, secrets จาก `.env` |
| [`backend/Dockerfile`](backend/Dockerfile) | multi-stage Go 1.24 → alpine, `TZ=Asia/Bangkok` |
| [`backend/.dockerignore`](backend/.dockerignore) | กัน `.env`, binaries |
| [`.env.example`](.env.example) | template secrets (commit ได้) — copy → `.env` |

**Secrets:** compose ใช้ `${POSTGRES_PASSWORD}`, `${JWT_SECRET}` จาก `.env` (root) — **ไม่ hardcode ใน repo**

**สองไฟล์ env:**

| ไฟล์ | ใช้เมื่อ |
|------|---------|
| `.env` (root) | `docker compose up` |
| `backend/.env` | `go run ./cmd/api` |

**Pitfalls Docker:** เปลี่ยน `POSTGRES_PASSWORD` หลัง volume สร้างแล้ว → `docker compose down -v` แล้ว up ใหม่; port 5432/8080 ชน local ได้

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
| 8 | **DONE** | DevOps — compose + Dockerfile + `.env.example` |

### Frontend roadmap (ถัดไป — **ยังไม่เริ่มโค้ด**)

**ตัดสินใจแล้ว:** เริ่มจาก **WALKIN (ขาจรจองแผง)** ก่อน role อื่น

| Phase | สถานะ | เป้าหมาย |
|-------|--------|----------|
| F1 | **TODO** | WALKIN — mock-login, รายการแผงว่าง, ฟอร์มจอง, แสดง 409 |
| F2 | TODO | STAFF — checklist roll-call + summary |
| F3 | TODO | REGULAR — แจ้งหยุด |
| F4 | TODO | Public — directory ค้นหาร้าน (+ แผนที่ posX/posY ถ้าต้องการ) |

**Stack (F1):** Vue 3 + **Vite** + **TypeScript** + **Tailwind CSS** + **Pinia** + Vue Router

**โฟลเดอร์ที่วางแผน:**

```
frontend/                    ← ยังไม่มี — สร้างตอน implement F1
├── src/
│   ├── api/                 # fetch client, types, error mapping
│   ├── stores/auth.ts       # Pinia: token, user, mockLogin(WALKIN)
│   ├── views/WalkinView.vue # list available + book
│   ├── components/
│   └── router/
├── .env.example             # VITE_API_URL=http://localhost:8080
└── vite.config.ts           # dev proxy → :8080 (ทางเลือกแทน CORS)
```

**F1 user flow:**

```
1. POST /api/auth/mock-login  { "role": "WALKIN" }  → เก็บ token (Pinia + localStorage)
2. GET  /api/stalls/available?date=YYYY-MM-DD       → แสดงรายการแผง pool
3. POST /api/bookings  { stallId, shopName, productType, date? }  → 201 หรือ 409
```

**F1 scope (MVP UI):**

- Mobile-first, Tailwind
- เลือก/แสดง `market_date` (default วันนี้ — ต้องตรง seed ใน DB)
- Loading + error states ทุก API call
- ข้อความ 409 เป็นภาษาไทย (แผงเต็ม / จองวันนี้แล้ว)
- **Out of scope F1:** แผนที่, STAFF/REGULAR, LINE LIFF, role picker หลาย role

**Backend ที่ต้องทำก่อน/พร้อมกับ F1:**

- [ ] **CORS** — `fiber/middleware/cors` ใน `main.go` (ตอนนี้ยังไม่มี; browser จาก `:5173` จะ block)
- [ ] (ทางเลือก) Vite `server.proxy` `/api` → `http://localhost:8080` ช่วง dev

**F1 implementation todos (เมื่อเริ่ม session):**

1. Scaffold `frontend/` (Vite + Vue + TS + Tailwind + Pinia + Router)
2. CORS หรือ Vite proxy
3. `api/client` + types + auth store
4. `WalkinView` — list + book form + 409 handling
5. README/AGENTS.md อัปเดต + commit

### Backend optional polish (ไม่ block F1)

1. **Seed daily bookings for today** — dev UX ข้ามวัน
2. **BOOK audit** — retroactive `audit_logs` ตอนจอง

### DONE — Backend (มีไฟล์แล้ว)

| Layer | ไฟล์ | หมายเหตุ |
|---|---|---|
| **entrypoint** | `backend/cmd/api/main.go` | Fiber routes — **middleware ต่อ route** (ดูด้านล่าง) |
| **devops** | `docker-compose.yml`, `backend/Dockerfile`, `backend/.dockerignore`, `.env.example` | Postgres 17 + API; secrets ใน `.env` |
| **config** | `backend/internal/config/config.go` | env loader (APP_PORT, DB_*, JWT_SECRET) |
| **domain** | `backend/internal/domain/models.go`, `errors.go` | entities, audit constants, inputs, `DirectoryEntry` DTO |
| **database** | `backend/internal/database/database.go`, `seed.go` | Connect, Migrate (+ partial unique index), SeedIfEmpty |
| **auth** | `backend/internal/auth/jwt.go` | TokenService Sign/Parse (HS256) |
| **repository** | `user.go`, `booking.go`, `audit.go` | user; booking CRUD+leave+attendance+SearchDirectory; audit Create |
| **usecase** | `auth.go`, `booking.go`, `vendor.go`, `staff.go`, `directory.go`, `*_test.go`, `booking_concurrent_test.go` | MockLogin; BookStall; Leave; MarkAttendance; Search + unit/integration tests |
| **handler** | `auth.go`, `booking.go`, `vendor.go`, `staff.go`, `directory.go`, `date.go`, `errors.go` | all MVP endpoints + error mapping |
| **middleware** | `auth.go`, `role.go` | JWT Auth, RequireRoles |
| **timeutil** | `backend/internal/pkg/timeutil/market_date.go` | `Today()`, `Now()` Asia/Bangkok |
| **env** | `backend/.env.example`, `.env.example` (root) | local dev vs docker compose |
| **deps** | `backend/go.mod`, `go.sum` | Fiber, godotenv, GORM, postgres, jwt/v5, testcontainers-go |

### TODO

- **Frontend F1** — ดู Frontend roadmap ด้านบน (ยังไม่ implement)
- Backend polish: seed-today, BOOK audit
- **Product Phase 2** — LINE LIFF, financial, PDPA (ดู section ท้ายไฟล์)

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

- **Backend:** Go 1.24 + Fiber v2 + GORM + PostgreSQL 17
- **Frontend (F1 วางแผนแล้ว):** Vue 3 + Vite + TypeScript + Tailwind + Pinia + Vue Router — โฟลเดอร์ `frontend/` ยังไม่มี
- **Auth:** JWT mock login (MVP) → production LINE LIFF
- **Project root:** `D:\taladsync`

## Architecture: Modular Monolith + Clean Architecture

**Dependency rule:** ลูกศรชี้เข้าด้านใน `handler → usecase → domain` ; `repository` implement interface ของ `usecase`

```
backend/                     ← DONE (Phase 0–8)
frontend/                    ← TODO F1 (WALKIN) — ยังไม่สร้าง
```

**Backend** (`cmd/api/main.go` …) — DONE: health, auth, booking, vendor, staff, directory

**Frontend** — ยังไม่มี; เริ่มที่ WALKIN ตาม Frontend roadmap


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

### Docker Compose (full stack)

```powershell
copy .env.example .env
# แก้ POSTGRES_PASSWORD, JWT_SECRET ใน .env (ไม่ commit .env)

docker compose up --build
# GET http://localhost:8080/health → {"status":"ok"}
```

หยุด local Postgres/API ก่อนถ้า port `5432` / `8080` ชน — reset volume: `docker compose down -v`

### Local (go run)

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

## Phase 8 Docker Compose

```powershell
# จาก repo root
copy .env.example .env
# แก้ POSTGRES_PASSWORD, JWT_SECRET

docker compose up --build
# GET http://localhost:8080/health → {"status":"ok"}
```

```powershell
docker compose down          # หยุด
docker compose down -v       # หยุด + ลบ volume (reset DB / แก้รหัสผ่าน PG)
docker compose up --build -d # รัน background
docker compose logs -f api   # ดู log
```

**Verify แล้ว:** build ผ่าน, seed อัตโนมัติ, health OK

## Phase 2 Product (ยังไม่ทำ — หลัง Frontend MVP)

- LINE LIFF / Push จริง (Flash Drop)
- โควตาวันหยุด, แจ้งหยุดฉุกเฉิน 2 ครั้ง/เดือน, ลาพักยาว+อนุมัติ
- Force Absent + ปล่อยแผง no-show อัตโนมัติ
- แจ้งมาช้า (late arrival lock)
- **Financial tracking** (เก็บเงิน 20 บ. ในระบบ + สรุปยอดส่งคลัง) — ตัดออกจาก MVP แล้ว
- national_id encryption เต็มรูปแบบ (PDPA)
