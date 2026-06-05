# IELTS Vocabulary Platform Docs

This document package is designed for vibe coding an IELTS Vocabulary Learning Platform with this stack:

```txt
Frontend: React.js + TypeScript + Vite + Tailwind CSS + shadcn/ui
Backend: Go + Gin
Database: MySQL
ORM: GORM
Auth: JWT
```

Recommended monorepo structure:

```txt
ielts-vocab-app/
├── frontend/
├── backend/
├── docs/
├── AGENTS.md
└── README.md
```

## Documents

```txt
docs/
├── prd.md
├── database.md
├── api.md
├── architecture.md
├── user-flow.md
├── design-system.md
├── seed-data.md
├── tasks.md

AGENTS.md
```

## Suggested Usage With Codex / AI Coding Agent

Copy the `docs/` folder and `AGENTS.md` into your project root.

Then ask Codex one task at a time:

```txt
Implement Task 04 from docs/tasks.md.

Read AGENTS.md first.
Follow docs/database.md, docs/api.md, and docs/architecture.md.
Use React.js frontend + Go Gin backend + MySQL + GORM.
Implement only this task.
Do not modify unrelated files.
After implementation, summarize changed files and how to test.
```

## MVP Build Order

1. Project setup
2. Backend setup
3. MySQL + GORM models
4. Seed data
5. Authentication
6. Frontend layout
7. Landing page
8. Dashboard
9. Roadmap
10. Lesson detail
11. Flashcards
12. Quiz
13. Progress tracking
14. Profile
