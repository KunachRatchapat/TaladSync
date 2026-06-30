# TaladSync 🏪

A daily fresh-market stall management system that solves two core pain points:

1. **Walk-in stall booking races** (concurrency) — prevents double-booking when several vendors tap "book" at the same time.
2. **On-site roll-call** — lets staff record who showed up vs. who didn't, keeping the system in sync with reality.

> 💡 The 20 THB stall fee is collected in cash on-site by staff — it is **not tracked in the system** (MVP).

---

## 🚧 Project Status

> **Phase 4 done.** Auth (mock-login, JWT, `/api/auth/me`) + booking (`GET /stalls/available`, `POST /bookings`, HTTP 409). Usecase unit tests for `BookStall`.
> Vendor leave, staff roll-call, directory, concurrent test, and Docker — not yet implemented.
> See the incremental roadmap in [`AGENTS.md`](AGENTS.md).

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
backend/
├── .env.example
├── cmd/api/main.go          # entrypoint
└── internal/
    ├── auth/                # JWT Sign/Parse
    ├── config/              # env loading
    ├── domain/              # models + business errors
    ├── database/            # connect + migrate + seed
    ├── handler/             # auth, booking HTTP
    ├── middleware/          # JWT auth, RequireRoles
    ├── repository/          # user, booking (GORM)
    ├── usecase/             # auth, booking + unit tests
    └── pkg/timeutil/        # market_date (Asia/Bangkok)
```

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
| POST | `/api/vendor/leave` | REGULAR | — |
| GET | `/api/staff/checklist?date=` | STAFF | — |
| POST | `/api/bookings/:id/attendance` | STAFF | — |
| GET | `/api/staff/summary?date=` | STAFF | — |
| GET | `/api/directory?q=` | public | — |

---

## 🚀 Running (dev)

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
go test ./internal/usecase/... -v -run BookStall
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
| 5 | — | Vendor leave + staff attendance |
| 6 | — | Customer directory |
| 7 | — | Concurrent load test (testcontainers) |
| 8 | — | Docker Compose |

## 🗺 Product Phase 2 (later)

- Real LINE LIFF / Push notifications (Flash Drop)
- Day-off quotas, emergency leave, long leave + approval
- Force Absent + automatic release of no-show stalls
- Financial tracking (stall-fee totals for treasury reporting)
- `national_id` encryption (PDPA compliance)
- Vue 3 frontend
