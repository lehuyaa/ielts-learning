# IELTS Vocabulary Learning Platform

Monorepo for an IELTS vocabulary learning platform focused on band-based lessons, flashcards, quizzes, spaced repetition, progress tracking, streaks, and achievements.

## Stack

```txt
Frontend: React.js + TypeScript + Vite + Tailwind CSS + shadcn/ui
Backend: Go + Gin
Database: MySQL
ORM: GORM
Auth: JWT
```

## Repository Structure

```txt
ielts-learning/
├── frontend/
├── backend/
├── docs/
├── AGENTS.md
├── docker-compose.yml
└── README.md
```

## Current Status

Task 01 creates the base monorepo structure only. The frontend app, backend server, and MySQL service are added in later tasks from `docs/tasks.md`.

## Run Frontend

After Task 02 creates the Vite app:

```bash
cd frontend
pnpm install
pnpm dev
```

The frontend should run at:

```txt
http://localhost:5173
```

## Run Backend

After Task 04 creates the Go Gin app:

```bash
cd backend
go mod tidy
go run ./cmd/api
```

The backend should run at:

```txt
http://localhost:8080
```

## Run Database

After Task 05 adds the MySQL Docker Compose service:

```bash
docker compose up -d
```

Expected local database defaults:

```txt
Host: localhost
Port: 3306
Database: ielts_vocab
User: root
Password: password
```

## Documentation

Read these before implementing tasks:

```txt
AGENTS.md
docs/prd.md
docs/database.md
docs/api.md
docs/architecture.md
docs/user-flow.md
docs/design-system.md
docs/seed-data.md
docs/tasks.md
```

## Development Rule

Implement one task from `docs/tasks.md` at a time. Keep changes small, scoped, and aligned with the documented React + TypeScript + Vite frontend and Go Gin + MySQL + GORM backend.
