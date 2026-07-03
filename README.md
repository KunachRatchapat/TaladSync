# TaladSync 🏪

A daily fresh-market stall management system that solves two core pain points:

1. **Walk-in stall booking races** (concurrency) — prevents double-booking when several vendors tap "book" at the same time.
2. **On-site roll-call** — lets staff record who showed up vs. who didn't, keeping the system in sync with reality.

> 💡 The 20 THB stall fee is collected in cash on-site by staff — it is **not tracked in the system** (MVP).

---

## 🚧 Project Status

> **Backend MVP complete (Phase 0–8).** All API endpoints, unit tests, concurrent integration test, and Docker Compose are implemented.
> **Phase 8 not yet committed** — see file list in [`AGENTS.md`](AGENTS.md).
> Frontend (Vue 3) is not started — see Product Phase 2 below.

---

## 🛠 Tech Stack

| Layer | Technology |
|-------|------------|
| Backend | Go 1.24 + Fiber v2 + GORM |
| Database | PostgreSQL 17 |
| Auth | JWT (mock login for MVP) → production: LINE LIFF |
| Frontend | Vue 3 + Tailwind + Pinia *(not started)* |

---

## 🏗 Architecture

Modular Monolith + Clean Architecture — dependencies point inward:

```
handler → usecase → domain
              ↑
         repository (implements usecase interfaces)
```

```
taladsync/
├── .env.example             # Docker Compose secrets template (copy → .env)
├── docker-compose.yml       # Postgres 17 + API
└── backend/
    ├── Dockerfile
    ├── .dockerignore
    ├── .env.example         # local go run (copy → backend/.env)
    ├── cmd/api/main.go
    └── internal/ ...
```

**Two env files:** root `.env` for Docker; `backend/.env` for `go run` — keep `JWT_SECRET` in sync if using both.

---

## ⚡ Core Concurrency Model

Stall booking relies on a **conditional atomic UPDATE**, not an application-level lock:

```sql
UPDATE daily_bookings
SET status = 'BOOKED', vendor_id = ?, ...
WHERE stall_id = ? AND market_date = ? AND status = 'AVAILABLE';
```

- `RowsAffected = 1` → booking succeeded
- `RowsAffected = 0` → someone got there first (HTTP 409)
- Double-booking is blocked on two levels: the `WHERE` condition + a `UNIQUE(stall_id, market_date)` constraint

More detail: [`docs/BOOKING_FLOW.md`](docs/BOOKING_FLOW.md)

---

## 🔌 API

| Method | Endpoint | Role | Status |
|--------|----------|------|--------|
| GET | `/health` | public | ✅ |
| POST | `/api/auth/mock-login` | public | ✅ |
| GET | `/api/auth/me` | JWT | ✅ |
| GET | `/api/stalls/available?date=` | WALKIN | ✅ |
| POST | `/api/bookings` | WALKIN | ✅ |
| POST | `/api/vendor/leave` | REGULAR | ✅|
| GET | `/api/staff/checklist?date=` | STAFF | ✅|
| POST | `/api/bookings/:id/attendance` | STAFF | ✅ |
| GET | `/api/staff/summary?date=` | STAFF | ✅ |
| GET | `/api/directory?q=` | public | ✅ |

---

## 🚀 Running

### Option A: Docker Compose (full stack)

Requires [Docker Desktop](https://www.docker.com/products/docker-desktop/). Stop local Postgres/API if ports `5432` / `8080` are already in use.

```powershell
copy .env.example .env
# Edit .env — set POSTGRES_PASSWORD and JWT_SECRET (do not commit .env)

docker compose up --build
# GET http://localhost:8080/health  →  {"status":"ok"}
# First run auto-seeds users, stalls, daily_bookings
```

Stop: `docker compose down`  
Reset DB (e.g. after changing `POSTGRES_PASSWORD`): `docker compose down -v` then `up --build` again

**Troubleshooting:** `password authentication failed for user "postgres"` — volume was initialized with an old password. Run `docker compose down -v` and start fresh, or set `.env` password to match the existing volume.

### Option B: Local dev (go run)

```powershell
# 1. Set up environment
copy backend\.env.example backend\.env
# Edit DB_PASSWORD to match your PostgreSQL

# 2. Create the database (first time only)
createdb -U postgres taladsync

# 3. Run the API (auto-seeds when users table is empty)
cd backend
go run ./cmd/api
# GET http://localhost:8080/health  →  {"status":"ok"}
```

### Testing with Mock Login

```json
POST /api/auth/mock-login
{ "role": "WALKIN" }
```

Available roles: `REGULAR` | `WALKIN` | `STAFF`  
Send the returned token in the header: `Authorization: Bearer <token>`

### Booking flow (WALKIN)

1. Mock-login with `{ "role": "WALKIN" }`
2. `GET /api/stalls/available?date=2026-06-30` — list pool stalls (`P-01`…`P-08` in seed)
3. `POST /api/bookings` with body:

```json
{ "stallId": 13, "shopName": "ร้านทด", "productType": "ผลไม้" }
```

Returns `201` on success, `409` if the stall was taken or vendor already booked today.

### Tests

```powershell
cd backend

# unit tests (no Docker)
go test ./internal/usecase/... -v

# concurrent integration test (requires Docker)
go test -tags=integration ./internal/usecase/... -v -run Concurrent -count=1
```

---

## 📚 Documentation

- [`docs/SYSTEM_DESIGN.md`](docs/SYSTEM_DESIGN.md) — overview, stack, schema
- [`docs/BOOKING_FLOW.md`](docs/BOOKING_FLOW.md) — state machine, transactions, concurrent test
- [`AGENTS.md`](AGENTS.md) — development status + remaining work

---

## 🗺 Roadmap (backend MVP)

| Phase | Status | Goal |
|-------|--------|------|
| 0–2 | ✅ | Health, schema, seed |
| 3 | ✅ | Auth (mock-login, JWT, `/me`) |
| 4 | ✅ | Booking (list + book + 409) |
| 5 | ✅ | Vendor leave + staff attendance |
| 6 | ✅ | Customer directory |
| 7 | ✅ | Concurrent load test (testcontainers) |
| 8 | ✅ | Docker Compose (Postgres + API, secrets via `.env`) |

### Phase 8 files (pending commit)

- `docker-compose.yml`, `backend/Dockerfile`, `backend/.dockerignore`, `.env.example`

## 🗺 Product Phase 2 (later)

- Real LINE LIFF / Push notifications (Flash Drop)
- Day-off quotas, emergency leave, long leave + approval
- Force Absent + automatic release of no-show stalls
- Financial tracking (stall-fee totals for treasury reporting)
- `national_id` encryption (PDPA compliance)
- Vue 3 frontend
