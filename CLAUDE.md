# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

IELTS Vocabulary Learning Platform (monorepo): vocabulary roadmap by IELTS band, topic-based lessons, flashcards, quizzes, spaced repetition, progress tracking, streaks, achievements.

Stack: React + TypeScript + Vite + Tailwind + shadcn/ui (frontend) · Go + Gin + GORM + MySQL + JWT (backend).

This repo also has an `AGENTS.md` with detailed, authoritative coding rules (module pattern, validation flow, quiz/spaced-repetition rules, UI fidelity rules, Swagger requirements, etc.) — read it before making non-trivial changes. Do not duplicate its rules from memory; re-read it, since it is the source of truth and may be updated independently of this file.

## Commands

### Frontend (`frontend/`)
```bash
pnpm install
pnpm dev        # http://localhost:5175
pnpm build      # tsc -b && vite build
pnpm lint       # eslint .
pnpm preview
```
No frontend test runner is configured yet.

### Backend (`backend/`)
```bash
go mod tidy
go run ./cmd/api          # http://localhost:8081
go run ./cmd/seed         # run seed data (backend/seeds/seed.go)
go build ./...
go vet ./...
go test ./...              # no _test.go files exist yet in this repo
swag init -g cmd/api/main.go   # regenerate backend/docs Swagger after any endpoint change
```

### Database
```bash
docker compose up -d       # MySQL 8.0, exposed on localhost:3308 -> container 3306
```
Defaults: db `ielts_vocab`, user `root`, password `password` (see `backend/.env.example` / `backend/internal/config/config.go` for env var overrides: `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `JWT_SECRET`, `JWT_ACCESS_TTL_MINUTES`, `APP_ENV`, `APP_PORT`, `FRONTEND_URL`).

Swagger UI is served at `/swagger/*any` once the backend is running.

## Architecture

### Backend: per-module vertical slices
`backend/internal/modules/<name>/` — one folder per feature (auth, dashboard, roadmap, topic, lesson, profile, notification, vocabulary, flashcard, quiz, activity, achievement, xp). Most modules follow a consistent 4-5 file pattern:
- `routes.go` — `RegisterRoutes(router *gin.RouterGroup, db *gorm.DB, jwtManager sharedjwt.Manager)` wires repository → service → handler and mounts routes, applying `middleware.Auth(jwtManager)` for protected groups.
- `handler.go` — parses/validates request, pulls user ID from JWT context (never trusts client-supplied user IDs), calls service, returns response via `internal/shared/response`.
- `service.go` — business logic (spaced repetition, quiz scoring, streaks, achievement unlocks, transactions).
- `repository.go` — GORM queries only.
- `dto.go` (+ `validation.go` in some modules) — request/response shapes and validation.

`cmd/api/main.go` is the composition root: loads config, connects DB, runs `database.AutoMigrate`, builds the JWT manager, registers Gin middleware (`gin.Logger`, `gin.Recovery`, `middleware.CORS`), mounts `/health` and `/swagger/*any`, then calls each module's `RegisterRoutes` under the `/api/v1` group. Not every model in `internal/models/auto_migrate.go` (e.g. placement tests, subscriptions) has a corresponding module/routes yet — check before assuming an API exists.

All GORM models live centrally in `internal/models/` (not per-module); `internal/models/auto_migrate.go` lists every model registered for `AutoMigrate`.

Shared API response envelope (`internal/shared/response`):
```go
type APIResponse struct {
    Data  any       `json:"data,omitempty"`
    Error *APIError `json:"error,omitempty"`
}
```
Handlers use `response.OK`, `response.Created`, `response.Error`, `response.ValidationError` — keep this envelope consistent for every endpoint.

Auth: `middleware.Auth` validates the `Authorization: Bearer <token>` header via the JWT manager and sets `userID` in the Gin context; handlers read it with `middleware.GetUserID(c)`. Never trust a client-submitted user ID.

### Frontend: feature-based, with React Query for server state
- `src/api/` — one Axios-calling module per backend resource (e.g. `auth.ts`, `dashboard.ts`, `roadmap.ts`), all built on the shared instance in `src/api/api.ts`. That file centralizes: base URL from `VITE_API_BASE_URL`, JWT attachment via request interceptor, response unwrapping/error normalization into a typed `APIError`, and a pluggable 401 handler (`setUnauthorizedHandler`, used by the auth context to trigger logout/redirect). Do not use `fetch` or add ad-hoc axios instances.
- `src/features/<name>/` — feature UI + `hooks/` (React Query hooks) + feature-local components; e.g. `features/auth/`, `features/quiz/`, `features/flashcard/`.
- `src/contexts/<name>/` — React Context providers, one folder per context (e.g. `contexts/auth/`, `contexts/toast/`); reserved for auth/theme/UI-preference state only, not server data.
- `src/pages/` — route-level page components, composed from `features/*` and `components/*`.
- `src/components/` — `ui/` (shadcn/ui primitives), `layout/`, `shared/`, `marketing/`, `state/` (loading/empty/error state components).
- `src/app/router.tsx` and `src/app/providers.tsx` — route table and top-level provider tree (React Query client, contexts, etc.).
- `src/types/` — shared TypeScript types (not scattered into feature folders).

Server state is always via TanStack Query (loading/error/caching/refetch come from Query, not manual `useState`/`useEffect`). Forms use React Hook Form + Zod; validation contracts are documented in `docs/validation-contracts.md` and must be mirrored on both frontend (Zod) and backend (DTO validation) — see `AGENTS.md` §5 "Validation".

### Docs worth checking before larger changes
`docs/prd.md`, `docs/database.md` / `docs/database-v2.md`, `docs/api.md` / `docs/api-v2.md`, `docs/architecture.md`, `docs/validation-contracts.md`, `docs/design-system.md`, `docs/user-flow.md`, `docs/seed-data.md`, `docs/tasks.md`, `docs/ui-review-checklist.md`, `docs/design-review-checklist.md`. `docs/screenshots/` is the visual source of truth for existing pages — don't redesign a page that already has an approved screenshot.

`docs/admin-guide.md` documents admin-only features as actually implemented (unlike the spec-style docs above) — currently: promoting a user to `ADMIN` via `go run ./cmd/seed -promote-admin=<email>`, and the Excel vocabulary import flow (`/admin/vocabularies/import*` endpoints, `/admin/vocabularies/import` page). Keep it up to date when admin features change.

### Admin role
`UserRole` has `USER` and `ADMIN` (`internal/models/enums.go`), embedded in the JWT and mirrored into the Gin context by `middleware.Auth`. Admin-only routes are gated with `middleware.RequireAdmin()` chained after `middleware.Auth()`. There is no signup path to `ADMIN` — see `docs/admin-guide.md` §1 to promote a user. Because the role is baked into the JWT, a promoted user must re-login before the new role takes effect.
