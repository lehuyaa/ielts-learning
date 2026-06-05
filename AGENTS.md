# AGENTS.md - Coding Rules for AI Coding Agents

This file defines how Codex, Claude Code, Cursor, or other AI coding agents should work in this repository.

## 1. Project Summary

This is an IELTS Vocabulary Learning Platform.

Version 1 focuses on:

- Vocabulary roadmap by IELTS band
- Topic-based lessons
- Flashcards
- Quizzes
- Spaced repetition
- User progress tracking
- Streaks and achievements

Chosen stack:

```txt
Frontend: React.js + TypeScript + Vite
Backend: Go + Gin
Database: MySQL
ORM: GORM
Auth: JWT
Styling: Tailwind CSS + shadcn/ui
```

Read these files before coding:

```txt
docs/prd.md
docs/database.md
docs/api.md
docs/architecture.md
docs/user-flow.md
docs/design-system.md
docs/seed-data.md
docs/tasks.md
```

---

## 2. Tech Stack Rules

Frontend must use:

- React.js
- TypeScript
- Vite
- Tailwind CSS
- shadcn/ui
- react-router-dom
- React Hook Form
- Zod
- Lucide React

Backend must use:

- Go
- Gin
- GORM
- MySQL
- JWT
- bcrypt
- validator

Do not use:

- Next.js
- Prisma
- PostgreSQL
- Redux
- Large UI libraries outside shadcn/ui

---

## 3. Repository Structure

Use:

```txt
frontend/
backend/
docs/
AGENTS.md
docker-compose.yml
README.md
```

Frontend structure:

```txt
frontend/src/
├── app/
├── pages/
├── components/
├── features/
├── lib/
├── hooks/
└── types/
```

Backend structure:

```txt
backend/
├── cmd/api/
├── internal/config/
├── internal/database/
├── internal/middleware/
├── internal/models/
├── internal/modules/
├── internal/shared/
└── seeds/
```

---

## 4. Coding Principles

1. Keep code simple.
2. Prefer readable code over clever code.
3. Use TypeScript strictly.
4. Do not use `any` casually.
5. Keep React components small and focused.
6. Keep Go handlers thin.
7. Put Go business logic in services.
8. Put GORM queries in repositories.
9. Validate all external input.
10. Never trust client-submitted user ID.
11. Always use authenticated user ID from JWT context.
12. Never expose password hash.
13. Never expose quiz correct answers before submission.

---

## 5. Frontend Rules

### API Client

Use one API client:

```txt
frontend/src/lib/api.ts
```

It should:

- Use `VITE_API_BASE_URL`
- Attach JWT token
- Parse standard backend response
- Throw meaningful errors

### Auth

Use:

```txt
frontend/src/features/auth/
```

Must include:

- Auth context
- Login form
- Register form
- ProtectedRoute
- Logout

### Routing

Use:

```txt
react-router-dom
```

Protected routes must redirect to:

```txt
/login
```

### UI

Follow:

```txt
docs/design-system.md
```

All pages should have:

- Loading state
- Empty state if relevant
- Error state if relevant
- Responsive design

---

## 6. Backend Rules

### Module Pattern

Each backend module should have:

```txt
handler.go
service.go
repository.go
dto.go
routes.go
```

### Handler

Handler should:

- Parse request
- Validate input
- Get user ID from context
- Call service
- Return response

Handler should not:

- Contain complex business logic
- Contain long GORM queries

### Service

Service should contain:

- Business rules
- Transactions
- Spaced repetition
- Quiz scoring
- Lesson completion
- Streak update
- Achievement unlock checks

### Repository

Repository should contain:

- GORM queries
- Create/update/find methods

---

## 7. API Rules

Use response shape:

```go
type APIResponse struct {
	Data  any       `json:"data,omitempty"`
	Error *APIError `json:"error,omitempty"`
}
```

Rules:

1. Protected APIs must check authenticated user.
2. Never trust `userId` from client.
3. Use user ID from JWT context.
4. Validate request body.
5. Use GORM transactions for multi-step mutations.
6. Return clear error codes.

---

## 8. Database Rules

Use GORM models from:

```txt
docs/database.md
```

Rules:

1. Use unsigned integer auto-increment IDs for MVP.
2. Add indexes for common query fields.
3. Use unique constraints for progress tables.
4. Do not delete user progress accidentally.
5. Use soft delete only where useful.
6. Seed data should be deterministic.
7. MySQL charset should be utf8mb4.

---

## 9. Spaced Repetition Rules

For MVP:

```txt
AGAIN -> review today / now
HARD  -> review tomorrow
GOOD  -> review in 3 days
EASY  -> review in 7 days
```

Update:

- status
- reviewCount
- correctCount / wrongCount
- lastRating
- lastReviewedAt
- nextReviewAt

Keep spaced repetition logic in:

```txt
backend/internal/modules/flashcard/service.go
```

or:

```txt
backend/internal/modules/progress/service.go
```

---

## 10. Quiz Rules

Important:

- Do not send `isCorrect` to frontend before submission.
- Score quiz on backend.
- Save quiz attempt.
- Mark lesson completed only if score >= requiredScore.
- Keep completed lesson completed even if later attempt is worse.
- Unlock next lesson after passing quiz.

---

## 11. Auth Rules

Protected frontend routes:

```txt
/dashboard
/roadmap
/lessons/*
/reviews
/vocabulary/*
/profile
```

Protected backend APIs:

```txt
/api/v1/dashboard/*
/api/v1/roadmap
/api/v1/lessons/*
/api/v1/reviews/*
/api/v1/flashcards/*
/api/v1/vocabularies/*
/api/v1/me/*
```

If unauthenticated:

- Frontend redirects to `/login`
- Backend returns `401 UNAUTHORIZED`

---

## 12. Task Execution Rules

When given a task:

1. Read relevant docs.
2. Implement only the requested task.
3. Do not modify unrelated files.
4. Explain what changed.
5. Mention any assumptions.
6. Mention how to test.
7. Keep changes small and reviewable.

Bad behavior:

- Rewriting unrelated architecture
- Adding unrequested features
- Changing database schema without reason
- Ignoring existing docs
- Creating duplicate components

---

## 13. Quality Checklist

Before finishing a task, check:

- Frontend TypeScript has no errors.
- Backend compiles.
- API response shape is consistent.
- UI is responsive.
- Empty states exist.
- Loading states exist where needed.
- API validates input.
- Auth is respected.
- No `any` used casually.
- No correct quiz answers exposed.
- No hardcoded user ID.
- Password hash is never returned.

---

## 14. Prompt Template for Codex

Use this style when asking the AI to code:

```txt
Implement Task [number] from docs/tasks.md.

Context:
- Read AGENTS.md first.
- Follow docs/prd.md.
- Follow docs/database.md.
- Follow docs/api.md.
- Follow docs/architecture.md.
- Follow docs/design-system.md.

Stack:
- React.js + TypeScript + Vite frontend
- Go Gin backend
- MySQL database
- GORM ORM

Constraints:
- Implement only this task.
- Do not modify unrelated files.
- Keep handlers thin and services clean.
- Add loading and empty states if relevant.

After implementation:
- Summarize changed files.
- Explain how to test.
- Mention any assumptions.
```
