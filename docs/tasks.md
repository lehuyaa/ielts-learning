# Task Breakdown - React + Go Gin + MySQL + GORM

## 1. Guiding Rule for Codex / AI Agent

Do not ask the AI to build everything at once.

Bad prompt:

```txt
Build the IELTS app.
```

Good prompt:

```txt
Implement Task 04 only. Follow docs/prd.md, docs/database.md, docs/api.md, docs/architecture.md, and AGENTS.md. Do not modify unrelated files.
```

---

## 2. Milestone 1 - Project Setup

### Task 01 - Create Monorepo Structure

Goal:

Create base project folders.

Structure:

```txt
ielts-vocab-app/
├── frontend/
├── backend/
├── docs/
├── AGENTS.md
├── docker-compose.yml
└── README.md
```

Acceptance criteria:

- Root structure exists
- README explains how to run frontend/backend/database

---

### Task 02 - Setup Frontend React App

Goal:

Create React.js app with Vite and TypeScript.

Requirements:

- Vite
- React
- TypeScript
- Tailwind CSS
- ESLint
- react-router-dom

Acceptance criteria:

- `pnpm dev` runs frontend
- Home page loads
- Tailwind works

---

### Task 03 - Setup shadcn/ui and Base Components

Install and configure:

- shadcn/ui
- lucide-react
- class-variance-authority
- clsx
- tailwind-merge
- react-hook-form
- zod
- @hookform/resolvers

Acceptance criteria:

- Button component works
- Card component works
- Badge component works
- Input component works
- Form component works

---

### Task 04 - Setup Backend Go Gin App

Goal:

Create Go backend project.

Requirements:

- Gin server
- Health check endpoint
- CORS middleware
- Environment config
- Basic logger

Endpoint:

```txt
GET /health
```

Response:

```json
{
  "status": "ok"
}
```

Acceptance criteria:

- `go run ./cmd/api` starts server
- `/health` works
- CORS allows frontend origin

---

### Task 05 - Setup MySQL with Docker Compose

Goal:

Create MySQL local environment.

Requirements:

- docker-compose.yml
- MySQL 8
- Database: ielts_vocab
- Root password: password for local dev

Acceptance criteria:

- `docker compose up -d` starts MySQL
- Backend can connect to MySQL

---

## 3. Milestone 2 - Database and Seed

### Task 06 - Add GORM Models

Goal:

Implement GORM models from `docs/database.md`.

Requirements:

- Create `backend/internal/models`
- Add all core models
- Add enum string constants
- Add JSON fields for synonyms/antonyms/collocations

Acceptance criteria:

- Backend compiles
- Models match database docs

---

### Task 07 - Setup GORM MySQL Connection and AutoMigrate

Goal:

Connect backend to MySQL and migrate schema.

Requirements:

- `backend/internal/database/mysql.go`
- `backend/internal/database/migrate.go`
- AutoMigrate all models

Acceptance criteria:

- Backend connects to MySQL
- Tables are created
- No migration error

---

### Task 08 - Seed Initial Data

Goal:

Create seed data.

Requirements:

- Course
- Band levels
- Topics
- Lessons
- Vocabulary
- LessonVocabulary
- Quiz questions
- Quiz options
- Achievements
- Demo user

Acceptance criteria:

- Seed command works
- At least 50 vocabulary items exist
- Demo user can login later

---

## 4. Milestone 3 - Auth

### Task 09 - Backend Auth APIs

Implement:

```txt
POST /api/v1/auth/register
POST /api/v1/auth/login
GET  /api/v1/auth/me
```

Requirements:

- bcrypt password hashing
- JWT generation
- Auth middleware
- Validation
- Standard response shape

Acceptance criteria:

- User can register
- User can login
- Protected me endpoint works with JWT
- Password hash is never returned

---

### Task 10 - Frontend Auth Pages

Create:

- Register page
- Login page
- Auth context
- ProtectedRoute
- Logout

Acceptance criteria:

- User can register from frontend
- User can login from frontend
- Token is stored
- Protected routes redirect if not logged in
- Logout works

---

## 5. Milestone 4 - Layout and Core UI

### Task 11 - Frontend Layouts

Create:

- Marketing layout
- Auth layout
- Dashboard layout
- Sidebar
- Topbar

Acceptance criteria:

- `/` uses marketing layout
- `/login` and `/register` use auth layout
- `/dashboard` uses dashboard layout
- Responsive navigation works

---

### Task 12 - Landing Page

Create landing page based on design system.

Sections:

- Hero
- Stats
- Features
- Roadmap preview
- Testimonials
- Pricing placeholder
- Footer

Acceptance criteria:

- Looks polished
- CTA buttons work
- Responsive

---

### Task 13 - Dashboard API

Backend endpoint:

```txt
GET /api/v1/dashboard/summary
```

Requirements:

- Protected
- Return user stats
- Return review due count
- Return lesson progress
- Return recent activity placeholder

Acceptance criteria:

- Endpoint returns dashboard data
- Works for demo user
- No hardcoded user ID

---

### Task 14 - Dashboard Page

Frontend page:

```txt
/dashboard
```

Components:

- Welcome card
- Today's progress
- Review due card
- Current streak card
- Target band progress
- Quick actions
- Recent activity
- Vocabulary statistics

Acceptance criteria:

- Fetches real API
- Shows loading state
- Shows empty state where needed
- Responsive

---

## 6. Milestone 5 - Roadmap and Lessons

### Task 15 - Roadmap API

Backend endpoint:

```txt
GET /api/v1/roadmap
```

Requirements:

- Return course
- Return band levels
- Return topics
- Return progress
- Return locked/unlocked states

Acceptance criteria:

- Endpoint returns nested roadmap
- Uses authenticated user progress
- No N+1 query issue if possible

---

### Task 16 - Roadmap Page

Frontend page:

```txt
/roadmap
```

Requirements:

- Show band levels
- Show topics
- Show lessons
- Locked/unlocked/completed states
- Progress percentage

Acceptance criteria:

- User can click unlocked lesson
- Locked lessons visually disabled
- Completed lessons show check mark

---

### Task 17 - Lesson Detail API

Backend endpoint:

```txt
GET  /api/v1/lessons/:lessonId
POST /api/v1/lessons/:lessonId/start
```

Requirements:

- Return lesson detail
- Return vocabulary list
- Return user progress
- Start endpoint creates/updates progress

Acceptance criteria:

- Lesson loads from database
- Start endpoint updates status
- Locked lessons are handled

---

### Task 18 - Lesson Detail Page

Frontend page:

```txt
/lessons/:lessonId
```

Requirements:

- Lesson header
- Vocabulary list
- Progress
- Start flashcards button
- Start quiz button
- Right sidebar

Acceptance criteria:

- Loads lesson from API
- Shows vocabulary items
- Buttons route correctly

---

## 7. Milestone 6 - Flashcards and Review

### Task 19 - Flashcard APIs

Backend endpoints:

```txt
GET  /api/v1/lessons/:lessonId/flashcards
GET  /api/v1/reviews/due
POST /api/v1/flashcards/review
```

Requirements:

- Return flashcards for lesson
- Return due review cards
- Save flashcard rating
- Calculate nextReviewAt

Acceptance criteria:

- UserVocabularyProgress is updated
- Review due endpoint works
- Spaced repetition logic follows docs/database.md

---

### Task 20 - Flashcard Page

Frontend pages:

```txt
/lessons/:lessonId/flashcards
/reviews
```

Requirements:

- Show one card at a time
- Flip/show meaning
- Rating buttons: Again, Hard, Good, Easy
- Update progress after rating
- Show session progress

Acceptance criteria:

- User can finish a flashcard session
- API is called after each rating
- UI is focused and responsive

---

## 8. Milestone 7 - Quiz

### Task 21 - Quiz APIs

Backend endpoints:

```txt
GET  /api/v1/lessons/:lessonId/quiz
POST /api/v1/lessons/:lessonId/quiz/submit
```

Requirements:

- Do not expose `isCorrect` before submit
- Score quiz on backend
- Save quiz attempt
- Complete lesson if passed
- Unlock next lesson if passed

Acceptance criteria:

- Correct answers hidden before submit
- Score is calculated correctly
- Lesson completion works

---

### Task 22 - Quiz Page

Frontend page:

```txt
/lessons/:lessonId/quiz
```

Requirements:

- Multiple choice questions
- Progress indicator
- Submit answers
- Show result

Acceptance criteria:

- User can answer and submit quiz
- Result page shows score
- Passed/failed state works

---

## 9. Milestone 8 - Vocabulary and Profile

### Task 23 - Vocabulary APIs

Backend endpoints:

```txt
GET /api/v1/vocabularies
GET /api/v1/vocabularies/:vocabularyId
```

Requirements:

- Search
- Filter by difficulty
- Filter by status
- Pagination
- Detail page data

Acceptance criteria:

- Search works
- Filters work
- Pagination works
- Detail endpoint works

---

### Task 24 - Vocabulary Pages

Frontend pages:

```txt
/vocabulary
/vocabulary/:vocabularyId
```

Requirements:

- Vocabulary list
- Search
- Filters
- Detail tabs
- Progress display

Acceptance criteria:

- Pages fetch real API
- Empty states work
- Responsive

---

### Task 25 - Profile APIs

Backend endpoints:

```txt
GET   /api/v1/me/profile
PATCH /api/v1/me/profile
```

Requirements:

- Return user profile
- Return stats
- Return achievements
- Update name and target band

Acceptance criteria:

- Profile loads
- User can update profile
- Validation works

---

### Task 26 - Profile Page

Frontend page:

```txt
/profile
```

Requirements:

- Profile card
- Stats
- Achievements
- Learning heatmap placeholder
- Edit profile form

Acceptance criteria:

- User can view profile
- User can update name and target band

---

## 10. Milestone 9 - Polish

### Task 27 - Streak and Daily Activity

Requirements:

- Track daily activity
- Update streak
- Show daily activity stats

Acceptance criteria:

- Studying today updates activity
- Current streak updates correctly
- Dashboard shows streak

---

### Task 28 - Achievements

Requirements:

- Unlock achievements based on rules
- Show achievements on profile

Initial achievements:

- First Lesson
- 7 Day Streak
- 100 Words Learned
- Education Master

Acceptance criteria:

- Achievement unlocks once
- Profile displays unlocked achievements

---

### Task 29 - Loading, Empty, and Error States

Add:

- Skeletons
- Empty states
- Error states
- Toasts

Acceptance criteria:

- No page feels broken while loading
- Empty data has clear CTA
- API errors show useful messages

---

### Task 30 - Final QA

Checklist:

- Register
- Login
- Dashboard
- Roadmap
- Lesson
- Flashcards
- Quiz
- Review due
- Vocabulary search
- Profile
- Logout

Acceptance criteria:

- Main happy path works end-to-end
