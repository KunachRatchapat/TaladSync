# CLAUDE.md — TaladSync

Daily fresh-market stall manager. Solves two pain points: **walk-in booking races** (concurrency) and **on-site roll-call** (attendance sync).

> The primary dev log lives in **[AGENTS.md](AGENTS.md)** — read it for phase status, next steps, and full file map. This file is the quick-reference; keep it short.

## Stack & layout

- **Backend:** Go 1.24 + Fiber v2 + GORM + PostgreSQL 17 — `backend/`, Clean Architecture (`handler → usecase → domain`, `repository` implements usecase interfaces).
- **Frontend:** Vue 3 + Vite + TypeScript + Tailwind + Pinia + Vue Router — `frontend/`. Four role pages: `/` WALKIN, `/staff` STAFF, `/regular` REGULAR, `/directory` public.
- **Status:** Backend Phase 0–8 done; Frontend F1–F4 done, merged to `main`. **Active work: UI redesign on branch `bookingvendor`** (designer-led, one component at a time) — see the UI Redesign section in [AGENTS.md](AGENTS.md) for what's done and what's blocked on designs.

## Common commands

```powershell
# Backend (local)
cd backend; go run ./cmd/api            # auto-seeds when users table empty
cd backend; go run ./cmd/api -reseed    # wipe + reseed dev data, then exit
cd backend; go test ./internal/usecase/... -v                                   # unit
cd backend; go test -tags=integration ./internal/usecase/... -run Concurrent    # needs Docker

# Frontend
cd frontend; npm run dev                # http://localhost:5173
cd frontend; npm run type-check         # vue-tsc --noEmit

# Full backend stack
docker compose up --build               # Postgres 17 + API on :8080
```

## Critical gotchas

- **Route wiring:** do not use multiple `api.Group("")` with different `RequireRoles` — Fiber merges middleware; staff routes can get WALKIN-only 403. Wire `RequireRoles` per route (see `cmd/api/main.go`).
- **Booking concurrency:** `TryBookStall` uses conditional `UPDATE ... WHERE status='AVAILABLE'`; success = `RowsAffected == 1`, not app-level mutex.
- **Two env files:** root `.env` for Docker Compose; `backend/.env` for `go run` — keep `JWT_SECRET` in sync if using both.
- **Vendor mine:** query by `stalls.owner_vendor_id`, not `daily_bookings.vendor_id` (cleared on leave).
- **Two headers coexist during redesign:** `PageHeader.vue` (new design, WalkinView only) vs `AppHeader.vue` (still used by staff/regular/directory) — do not delete `AppHeader` until every page is redesigned.
- **Zones:** seed now creates **โซน A** (owned by regulars) + **โซน B** (unowned, walk-in pool). `/api/stalls/available` only returns AVAILABLE rows, so โซน A appears only after a regular calls leave.
- **Workflow:** user reviews each change before committing — make the edit, then stop and report; do not commit or push unless asked.
- Match surrounding style. Backend comments and AGENTS.md are partly Thai — mirror the file you're editing.
- After finishing meaningful work, update the checkpoint section in `AGENTS.md`.

## Docs map

| File | Use when |
|------|----------|
| [AGENTS.md](AGENTS.md) | phase status, file map, conventions |
| [docs/SYSTEM_DESIGN.md](docs/SYSTEM_DESIGN.md) | architecture overview, stack, schema |
| [docs/DATABASE.md](docs/DATABASE.md) | schema + concurrency (source of truth for MVP) |
| [docs/DEMO.md](docs/DEMO.md) | run, demo script, `-reseed` |
