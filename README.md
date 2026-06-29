# TaladSync 🏪

A daily fresh-market stall management system that solves two core pain points:

1. **Walk-in stall booking races** (concurrency) — prevents double-booking when several vendors tap "book" at the same time.
2. **On-site roll-call** — lets staff record who showed up vs. who didn't, keeping the system in sync with reality.

> 💡 The 20 THB stall fee is collected in cash on-site by staff — it is **not tracked in the system** (MVP).

---

## 🚧 Project Status

> **Fresh restart — Phase 2 done.** Dev seed data loads automatically when the database is empty (3 users, 20 stalls, today's daily_bookings).
> Business API endpoints are not implemented yet.
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
├── cmd/api/main.go   # entrypoint (health + DB migrate)
└── internal/
    ├── config/       # env loading
    ├── domain/       # models + business errors
    ├── database/     # connect + migrate + seed
    └── pkg/timeutil/ # market_date (Asia/Bangkok)
    # usecase, repository, handler, … — TODO Phase 3+
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

## 🔌 API (MVP)

| Method | Endpoint | Role | Description |
|--------|----------|------|-------------|
| POST | `/api/auth/mock-login` | public | `{role}` → `{token, user}` |
| GET  | `/api/stalls/available?date=` | WALKIN | available stalls in the pool |
| POST | `/api/bookings` | WALKIN | book a stall (409 if taken) |
| POST | `/api/vendor/leave` | REGULAR | report a day off |
| GET  | `/api/staff/checklist?date=` | STAFF | roll-call of every occupied stall |
| POST | `/api/bookings/:id/attendance` | STAFF | mark present/absent (idempotent) |
| GET  | `/api/staff/summary?date=` | STAFF | counts of present/no-show/pending |
| GET  | `/api/directory?q=` | public | search shops (no PII) |
| GET  | `/health` | public | health check |

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

---

## 📚 Documentation

- [`docs/SYSTEM_DESIGN.md`](docs/SYSTEM_DESIGN.md) — overview, stack, schema
- [`docs/BOOKING_FLOW.md`](docs/BOOKING_FLOW.md) — state machine, transactions, concurrent test
- [`AGENTS.md`](AGENTS.md) — development status + remaining work

---

## 🗺 Roadmap (Phase 2)

- Real LINE LIFF / Push notifications (Flash Drop)
- Day-off quotas, emergency leave, long leave + approval
- Force Absent + automatic release of no-show stalls
- Financial tracking (stall-fee totals for treasury reporting)
- `national_id` encryption (PDPA compliance)
- Vue 3 frontend
