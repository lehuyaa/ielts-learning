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

Frontend UI task rule:

- Every frontend UI task must be checked against `docs/ui-review-checklist.md` before completion.
- Every frontend UI task must be checked against `docs/design-review-checklist.md` before completion.
- Every frontend UI task must check `docs/screenshots/README.md` when it exists.
- For every frontend UI task, analyze the reference screenshot in `docs/screenshots/` before coding when one exists.
- The screenshot is the source of truth for layout, section order, component hierarchy, spacing hierarchy, typography hierarchy, colors, and interaction placement.
- Do not redesign a page when a screenshot exists.
- The final task summary should mention that the UI review checklist was considered for frontend UI work.
- The final task summary should mention screenshot/design review and report visual differences before marking frontend UI work complete.

Validation contract rule:

- Every task that implements a user-facing form must check `docs/validation-contracts.md` before implementation.
- If the form does not have a validation contract, add one before implementing the form.
- If the form submits to the backend, implement both frontend Zod validation and backend request DTO validation from the same contract.
- Never implement validation on only one side for backend-submitted forms.
- Any validation change must update `docs/validation-contracts.md`, the frontend Zod schema, and backend request DTO validation.

Frontend API consumption rule:

- Every frontend task that consumes backend APIs must use React Query.
- Do not fetch API data inside `useEffect`.
- Do not manage loading/error state manually with `useState` + `useEffect` when React Query can handle it.
- API modules must stay in `frontend/src/api`.
- React Query hooks should live in the relevant feature folder, for example `frontend/src/features/roadmap/hooks/useRoadmap.ts`.
- Pages and components should consume React Query hooks instead of calling API modules directly.

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
- axios
- @tanstack/react-query

Acceptance criteria:

- `pnpm dev` runs frontend
- Home page loads
- Tailwind works
- Shared frontend folders exist:
  - `src/api`
  - `src/contexts`
  - `src/types`

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
- Backend request DTO validation must match `docs/validation-contracts.md`
- Return field-level validation errors for validation failures
- Standard response shape

Acceptance criteria:

- User can register
- User can login
- Protected me endpoint works with JWT
- Password hash is never returned
- Register/Login backend validation matches the auth validation contracts

---

### Task 10 - Frontend Auth Pages

Create:

- Register page
- Login page
- Auth context in `frontend/src/contexts/auth`
- ProtectedRoute
- Logout
- Shared axios instance in `frontend/src/api/api.ts`
- Auth API module in `frontend/src/api/auth.ts`
- Shared auth types in `frontend/src/types/auth.ts`
- Check `docs/validation-contracts.md`
- Frontend Zod schemas must match the Register/Login validation contracts

Requirements:

- Analyze the reference screenshot before coding if one exists.
- The screenshot is the source of truth.
- Do not redesign the page.

Acceptance criteria:

- User can register from frontend
- User can login from frontend
- Token is stored
- Protected routes redirect if not logged in
- Logout works
- Frontend API calls use axios only
- No frontend code uses `fetch`
- No shared types are placed in feature folders
- Register/Login frontend validation matches the auth validation contracts
- Layout closely matches screenshot if one exists.
- Typography hierarchy matches screenshot if one exists.
- Spacing hierarchy matches screenshot if one exists.
- Responsive behavior remains correct.

---

## 5. Milestone 4 - Layout and Core UI

### Task 11 - Frontend Layouts

Create:

- Marketing layout
- Auth layout
- Dashboard layout
- Sidebar
- Topbar

Requirements:

- Analyze the reference screenshot before coding if one exists.
- The screenshot is the source of truth.
- Do not redesign the page.

Acceptance criteria:

- `/` uses marketing layout
- `/login` and `/register` use auth layout
- `/dashboard` uses dashboard layout
- Responsive navigation works
- Layout closely matches screenshot if one exists.
- Typography hierarchy matches screenshot if one exists.
- Spacing hierarchy matches screenshot if one exists.
- Responsive behavior remains correct.

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

Requirements:

- Analyze the reference screenshot before coding if one exists.
- The screenshot is the source of truth.
- Do not redesign the page.

Acceptance criteria:

- Looks polished
- CTA buttons work
- Responsive
- Layout closely matches screenshot if one exists.
- Typography hierarchy matches screenshot if one exists.
- Spacing hierarchy matches screenshot if one exists.
- Responsive behavior remains correct.

---

### Task 13 - Roadmap Page UI

Frontend page:

```txt
/roadmap
```

Goal:

Create the main learning roadmap page using mock data first.

Requirements:

- Analyze the reference screenshot before coding if one exists.
- The screenshot is the source of truth.
- Do not redesign the page.
- Use the existing Dashboard layout.
- Show band levels:
  - Band 5.0
  - Band 6.0
  - Band 7.0
  - Band 8.0
- Show topic cards under each band.
- Show topic locked, unlocked, in-progress, and completed states.
- Show topic progress percentage.
- Topic cards route to `/topics/:topicId`.
- Use polished UI based on `docs/design-system.md`.
- Check `docs/ui-review-checklist.md`.
- Use mock data only in this task.
- Do not implement backend API in this task.
- Do not implement Topic Detail, Lesson Detail, quiz, or flashcard functionality.

Acceptance criteria:

- User can open `/roadmap`.
- Roadmap is visually clear and responsive.
- Locked topics are visually disabled.
- Completed topics show a clear completed state.
- Mock data structure is easy to replace with real API data later.
- Layout closely matches screenshot if one exists.
- Typography hierarchy matches screenshot if one exists.
- Spacing hierarchy matches screenshot if one exists.
- Responsive behavior remains correct.

---

### Task 14 - Roadmap API

Backend endpoint:

```txt
GET /api/v1/roadmap
```

Goal:

Return real roadmap data from MySQL.

Requirements:

- Protected endpoint.
- Use authenticated user ID from JWT context.
- Return course.
- Return band levels.
- Return topics.
- Return topic progress.
- Return locked/unlocked/completed states.
- Return progress percentages.
- Avoid hardcoded user ID.
- Avoid N+1 query problems where reasonable.
- Use standard API response shape.

Acceptance criteria:

- Endpoint returns nested roadmap data.
- Works for demo user.
- Uses authenticated user topic progress.
- Does not expose unrelated user data.
- Frontend can consume the response directly in Task 15.

---

### Task 15 - Connect Roadmap Page to API

Frontend page:

```txt
/roadmap
```

Goal:

Replace roadmap mock data with real backend data.

Requirements:

- Analyze the reference screenshot before coding if one exists.
- The screenshot is the source of truth.
- Do not redesign the page.
- Add roadmap API function in `frontend/src/api/roadmap.ts`.
- Add shared roadmap types in `frontend/src/types/roadmap.ts`.
- Use axios through the shared API client.
- Do not use fetch.
- Show loading state.
- Show error state.
- Show empty state if no roadmap data exists.
- Keep existing UI from Task 13.
- Topic cards navigate to `/topics/:topicId`.
- Check `docs/ui-review-checklist.md`.

Acceptance criteria:

- Roadmap page fetches real API data.
- Loading and error states work.
- The page remains responsive.
- No mock roadmap data remains in production page logic.
- Layout closely matches screenshot if one exists.
- Typography hierarchy matches screenshot if one exists.
- Spacing hierarchy matches screenshot if one exists.
- Responsive behavior remains correct.

---

## 6. Milestone 5 - Lesson and Vocabulary

### Task 16 - Lesson Detail Page UI

Frontend page:

```txt
/lessons/:lessonId
```

Goal:

Create the lesson detail page using mock data first.

Requirements:

- Analyze the reference screenshot before coding if one exists.
- The screenshot is the source of truth.
- Do not redesign the page.
- Use Dashboard layout.
- Show lesson title.
- Show topic/band information.
- Show lesson description.
- Show estimated minutes.
- Show required score.
- Show XP reward if available.
- Show vocabulary preview list.
- Show progress card.
- Show CTA buttons:
  - Start Flashcards
  - Start Quiz
- Show right sidebar summary.
- Use mock data only in this task.
- Check `docs/ui-review-checklist.md`.

Acceptance criteria:

- User can open `/lessons/:lessonId`.
- Page looks polished and responsive.
- CTA buttons route to expected pages:
  - `/lessons/:lessonId/flashcards`
  - `/lessons/:lessonId/quiz`
- No backend integration required yet.
- Layout closely matches screenshot if one exists.
- Typography hierarchy matches screenshot if one exists.
- Spacing hierarchy matches screenshot if one exists.
- Responsive behavior remains correct.

---

### Task 17 - Topic Detail / Lesson List UI

Frontend page:

```txt
/topics/:topicId
```

Goal:

Create the Topic Detail page using mock data first.

Requirements:

- Analyze the reference screenshot before coding if one exists.
- The screenshot is the source of truth.
- Do not redesign the page.
- Use Dashboard layout.
- Show topic title, icon, band, and description.
- Show topic progress.
- Show statistics:
  - Progress percentage
  - Completed lessons
  - Total XP
- Show lesson list.
- Show locked lessons.
- Show unlocked lessons.
- Show in-progress lessons.
- Show completed lessons.
- Lesson cards route to `/lessons/:lessonId`.
- Use mock data only in this task.
- Do not integrate backend API.
- Do not modify Lesson Detail UI except for route compatibility if needed.
- Check `docs/ui-review-checklist.md`.

Acceptance criteria:

- User can open `/topics/:topicId`.
- Page closely matches the Topic Detail screenshot.
- Lesson cards route to Lesson Detail.
- Locked lessons are visually disabled and show unlock reason when available.
- In-progress lessons show progress percentage.
- Page is responsive.
- No backend integration required yet.

---

### Task 18 - Topic Detail API

Backend endpoint:

```txt
GET /api/v1/topics/:topicId
```

Goal:

Return Topic Detail data from MySQL.

Requirements:

- Protected endpoint.
- Use authenticated user ID from JWT context.
- Return topic metadata.
- Return parent band level.
- Return topic summary:
  - progressPercentage
  - completedLessons
  - totalLessons
  - totalXP
- Return lessons for the topic.
- Return lesson status:
  - LOCKED
  - UNLOCKED
  - IN_PROGRESS
  - COMPLETED
- Return lesson card fields:
  - id
  - title
  - description
  - wordCount
  - estimatedMinutes
  - xpReward
  - status
  - progressPercentage
  - lockedReason
- Query progress only for authenticated user.
- Avoid hardcoded user ID.
- Avoid N+1 query problems where reasonable.
- Use standard API response shape.

Acceptance criteria:

- Endpoint returns Topic Detail data.
- Works for demo user.
- Uses authenticated user progress.
- Does not expose unrelated user data.
- Missing topic returns NOT_FOUND.
- Frontend can consume the response directly in Task 19.

---

### Task 19 - Connect Topic Detail API

Frontend page:

```txt
/topics/:topicId
```

Goal:

Replace Topic Detail mock data with real backend data.

Requirements:

- Analyze the reference screenshot before coding if one exists.
- The screenshot is the source of truth.
- Do not redesign the page.
- Add topic API function in `frontend/src/api/topic.ts`.
- Add shared topic types in `frontend/src/types/topic.ts`.
- Add React Query hook in `frontend/src/features/topic/hooks/useTopicDetail.ts`.
- Use axios through shared API client.
- Do not use fetch.
- Show loading state.
- Show error state.
- Show not found/empty state.
- Keep existing UI from Task 17.
- Check `docs/ui-review-checklist.md`.

Acceptance criteria:

- Topic Detail page fetches real API data.
- Loading and error states work.
- Lesson cards route correctly to `/lessons/:lessonId`.
- No mock topic data remains in production page logic.
- Layout closely matches screenshot if one exists.
- Typography hierarchy matches screenshot if one exists.
- Spacing hierarchy matches screenshot if one exists.
- Responsive behavior remains correct.

---

### Task 20 - Lesson Detail API

Backend endpoints:

```txt
GET  /api/v1/lessons/:lessonId
POST /api/v1/lessons/:lessonId/start
```

Goal:

Return lesson detail and allow starting a lesson.

Requirements:

- Protected endpoints.
- Use authenticated user ID from JWT context.
- `GET /lessons/:lessonId` returns:
  - lesson detail
  - topic
  - band level
  - vocabulary list
  - user lesson progress
  - user vocabulary progress where available
- `POST /lessons/:lessonId/start`:
  - creates or updates user_lesson_progresses
  - sets status to IN_PROGRESS when appropriate
  - sets started_at if not already set
  - sets last_studied_at
- Handle locked lessons if lock logic already exists.
- Use standard API response shape.

Acceptance criteria:

- Lesson loads from database.
- Start endpoint updates progress.
- No hardcoded user ID.
- Missing lesson returns NOT_FOUND.

---

### Task 21 - Connect Lesson Detail Page to API

Frontend page:

```txt
/lessons/:lessonId
```

Goal:

Replace mock lesson data with real API data.

Requirements:

- Analyze the reference screenshot before coding if one exists.
- The screenshot is the source of truth.
- Do not redesign the page.
- Add lesson API functions in `frontend/src/api/lesson.ts`.
- Add shared lesson types in `frontend/src/types/lesson.ts`.
- Use axios through shared API client.
- Do not use fetch.
- Show loading state.
- Show error state.
- Show not found/empty state.
- Start button should call `POST /lessons/:lessonId/start` if needed.
- Keep existing UI from Task 16.
- Check `docs/ui-review-checklist.md`.

Acceptance criteria:

- Lesson detail page fetches real API data.
- Start lesson action works.
- CTA buttons route correctly.
- No mock lesson data remains in production page logic.
- Layout closely matches screenshot if one exists.
- Typography hierarchy matches screenshot if one exists.
- Spacing hierarchy matches screenshot if one exists.
- Responsive behavior remains correct.

---

### Task 22 - Vocabulary List and Detail UI

Frontend pages:

```txt
/vocabulary
/vocabulary/:vocabularyId
```

Goal:

Create vocabulary browsing UI using mock data first.

Requirements:

- Analyze the reference screenshot before coding if one exists.
- The screenshot is the source of truth.
- Do not redesign the page.

Vocabulary list page:

- Search input.
- Difficulty filter.
- Band filter if useful.
- Vocabulary cards or table.
- Difficulty badges.
- Progress/status badges.
- Empty state.
- Responsive layout.

Vocabulary detail page:

- Word header.
- IPA.
- Part of speech.
- Meaning VI.
- Meaning EN.
- Example sentence.
- Synonyms.
- Antonyms.
- Collocations.
- IELTS usage.
- User progress status placeholder.

Use mock data only in this task.

Acceptance criteria:

- `/vocabulary` page works with mock data.
- `/vocabulary/:vocabularyId` page works with mock data.
- UI is responsive.
- Check `docs/ui-review-checklist.md`.
- Layout closely matches screenshot if one exists.
- Typography hierarchy matches screenshot if one exists.
- Spacing hierarchy matches screenshot if one exists.
- Responsive behavior remains correct.

---

### Task 23 - Vocabulary APIs

Backend endpoints:

```txt
GET /api/v1/vocabularies
GET /api/v1/vocabularies/:vocabularyId
```

Goal:

Return vocabulary list and vocabulary detail from database.

Requirements:

- Protected endpoints.
- Use authenticated user ID from JWT context.
- Support query params:
  - q
  - difficulty
  - targetBand
  - status
  - page
  - limit
- Return pagination metadata.
- Include user progress status where available.
- Detail endpoint returns full vocabulary data.
- Never expose unrelated user data.

Acceptance criteria:

- Search works.
- Filters work.
- Pagination works.
- Detail endpoint works.
- Works for demo user.

---

### Task 24 - Connect Vocabulary Pages to API

Frontend pages:

```txt
/vocabulary
/vocabulary/:vocabularyId
```

Goal:

Replace mock vocabulary data with real backend data.

Requirements:

- Analyze the reference screenshot before coding if one exists.
- The screenshot is the source of truth.
- Do not redesign the page.
- Add vocabulary API functions in `frontend/src/api/vocabulary.ts`.
- Add shared vocabulary types in `frontend/src/types/vocabulary.ts`.
- Use axios through shared API client.
- Do not use fetch.
- Search and filters update API query params.
- Show loading, error, and empty states.
- Check `docs/ui-review-checklist.md`.

Acceptance criteria:

- Vocabulary pages load real API data.
- Search and filters work.
- Pagination works if implemented.
- No mock vocabulary data remains in production page logic.
- Layout closely matches screenshot if one exists.
- Typography hierarchy matches screenshot if one exists.
- Spacing hierarchy matches screenshot if one exists.
- Responsive behavior remains correct.

---

## 7. Milestone 6 - Flashcards and Review

### Task 25 - Flashcard Learning UI

Frontend pages:

```txt
/lessons/:lessonId/flashcards
/reviews
```

Goal:

Create flashcard learning and review UI using mock data first.

Requirements:

- Analyze the reference screenshot before coding if one exists.
- The screenshot is the source of truth.
- Do not redesign the page.
- Show one card at a time.
- Front side:
  - word
  - IPA
  - part of speech if available
- Back side:
  - meaning
  - example
  - synonyms/collocations if available
- Show flip/show meaning behavior.
- Show session progress bar.
- Rating buttons:
  - Again
  - Hard
  - Good
  - Easy
- Show completion summary.
- Use mock data only in this task.
- Check `docs/ui-review-checklist.md`.

Acceptance criteria:

- User can complete a mock flashcard session.
- UI is focused and responsive.
- No backend integration required yet.
- Layout closely matches screenshot if one exists.
- Typography hierarchy matches screenshot if one exists.
- Spacing hierarchy matches screenshot if one exists.
- Responsive behavior remains correct.

---

### Task 26 - Flashcard and Review APIs

Backend endpoints:

```txt
GET  /api/v1/lessons/:lessonId/flashcards
GET  /api/v1/reviews/due
POST /api/v1/flashcards/review
```

Goal:

Return flashcards and save review progress.

Requirements:

- Protected endpoints.
- Use authenticated user ID from JWT context.
- `GET /lessons/:lessonId/flashcards` returns lesson vocabulary.
- `GET /reviews/due` returns vocabulary where next_review_at <= now.
- `POST /flashcards/review` saves rating:
  - AGAIN
  - HARD
  - GOOD
  - EASY
- Update:
  - status
  - review_count
  - correct_count
  - wrong_count
  - last_rating
  - last_reviewed_at
  - next_review_at
  - learned_at/mastered_at when appropriate
- Update daily activity if implemented.
- Add XP event if XP system already exists.
- Use GORM transaction where appropriate.

Acceptance criteria:

- Review due endpoint works.
- Flashcard review updates user_vocabulary_progresses.
- Spaced repetition follows `docs/database.md`.
- No hardcoded user ID.

---

### Task 27 - Connect Flashcard Pages to API

Frontend pages:

```txt
/lessons/:lessonId/flashcards
/reviews
```

Goal:

Replace mock flashcard data with real backend data.

Requirements:

- Analyze the reference screenshot before coding if one exists.
- The screenshot is the source of truth.
- Do not redesign the page.
- Add flashcard API functions in `frontend/src/api/flashcard.ts`.
- Add shared flashcard/progress types in:
  - `frontend/src/types/flashcard.ts`
  - `frontend/src/types/progress.ts`
- Use axios through shared API client.
- Do not use fetch.
- Submit rating after each card.
- Show loading, error, and empty states.
- If no reviews are due, show a positive empty state.
- Check `docs/ui-review-checklist.md`.

Acceptance criteria:

- Lesson flashcards load from API.
- Review due cards load from API.
- Ratings are saved.
- Session summary works with real data.
- Layout closely matches screenshot if one exists.
- Typography hierarchy matches screenshot if one exists.
- Spacing hierarchy matches screenshot if one exists.
- Responsive behavior remains correct.

---

## 8. Milestone 7 - Quiz

### Task 28 - Quiz Page UI

Frontend page:

```txt
/lessons/:lessonId/quiz
```

Goal:

Create quiz UI using mock data first.

Requirements:

- Analyze the reference screenshot before coding if one exists.
- The screenshot is the source of truth.
- Do not redesign the page.
- Multiple choice question card.
- Question progress indicator.
- Options A/B/C/D.
- Selected answer state.
- Next button.
- Submit button.
- Score/result screen.
- Passed/failed state.
- Timer placeholder if designed.
- Use mock data only in this task.
- Check `docs/ui-review-checklist.md`.

Acceptance criteria:

- User can complete a mock quiz.
- Result screen shows score.
- UI is responsive.
- No backend integration required yet.
- Layout closely matches screenshot if one exists.
- Typography hierarchy matches screenshot if one exists.
- Spacing hierarchy matches screenshot if one exists.
- Responsive behavior remains correct.

---

### Task 29 - Quiz APIs

Backend endpoints:

```txt
GET  /api/v1/lessons/:lessonId/quiz
POST /api/v1/lessons/:lessonId/quiz/submit
```

Goal:

Return quiz questions and grade quiz submission.

Requirements:

- Protected endpoints.
- Use authenticated user ID from JWT context.
- GET endpoint:
  - returns questions and options
  - must not expose `is_correct`
- Submit endpoint:
  - accepts selected answers
  - calculates score on backend
  - returns result details
  - saves user_quiz_attempts
  - completes lesson if score >= required_score
  - updates user_lesson_progresses
  - unlocks next lesson if implemented
  - updates daily activity if implemented
  - creates XP events if XP system already exists
- Use GORM transaction for submit flow.

Acceptance criteria:

- Correct answers are hidden before submit.
- Score is calculated correctly.
- Quiz attempt is saved.
- Lesson completion works.
- No hardcoded user ID.

---

### Task 30 - Connect Quiz Page to API

Frontend page:

```txt
/lessons/:lessonId/quiz
```

Goal:

Replace mock quiz data with real backend data.

Requirements:

- Analyze the reference screenshot before coding if one exists.
- The screenshot is the source of truth.
- Do not redesign the page.
- Add quiz API functions in `frontend/src/api/quiz.ts`.
- Add shared quiz types in `frontend/src/types/quiz.ts`.
- Use axios through shared API client.
- Do not use fetch.
- Submit answers to backend.
- Show loading, error, and empty states.
- Show result from backend response.
- Check `docs/ui-review-checklist.md`.

Acceptance criteria:

- Quiz page loads real questions.
- Correct answers are not visible before submit.
- Submit returns real score.
- Passed/failed state works.
- Lesson completion can be verified after passing.
- Layout closely matches screenshot if one exists.
- Typography hierarchy matches screenshot if one exists.
- Spacing hierarchy matches screenshot if one exists.
- Responsive behavior remains correct.

---

## 9. Milestone 8 - Dashboard and Profile

### Task 31 - Dashboard Page UI

Frontend page:

```txt
/dashboard
```

Goal:

Create dashboard UI using mock data first.

Components:

- Welcome card.
- Continue learning card.
- Today's progress card.
- Review due card.
- Current streak card.
- Target band progress.
- XP/level summary.
- Quick actions.
- Recent activity.
- Vocabulary statistics.

Requirements:

- Analyze the reference screenshot before coding if one exists.
- The screenshot is the source of truth.
- Do not redesign the page.
- Use Dashboard layout.
- Use mock data only in this task.
- Check `docs/ui-review-checklist.md`.

Acceptance criteria:

- Dashboard looks polished and responsive.
- User has clear next action.
- No backend integration required yet.
- Layout closely matches screenshot if one exists.
- Typography hierarchy matches screenshot if one exists.
- Spacing hierarchy matches screenshot if one exists.
- Responsive behavior remains correct.

---

### Task 32 - Dashboard API

Backend endpoint:

```txt
GET /api/v1/dashboard/summary
```

Goal:

Return dashboard summary from real user data.

Requirements:

- Protected endpoint.
- Use authenticated user ID from JWT context.
- Return:
  - user summary
  - words learned today
  - review due count
  - current streak
  - total words learned
  - lessons completed
  - mastery percentage
  - target band progress
  - XP/level
  - recent activity placeholder or real data if available
- No hardcoded user ID.

Acceptance criteria:

- Endpoint returns dashboard data.
- Works for demo user.
- Response is directly usable by frontend.

---

### Task 33 - Connect Dashboard Page to API

Frontend page:

```txt
/dashboard
```

Goal:

Replace dashboard mock data with real backend data.

Requirements:

- Analyze the reference screenshot before coding if one exists.
- The screenshot is the source of truth.
- Do not redesign the page.
- Add dashboard API functions in `frontend/src/api/dashboard.ts`.
- Add shared dashboard types in `frontend/src/types/dashboard.ts`.
- Use axios through shared API client.
- Do not use fetch.
- Show loading, error, and empty states.
- Keep existing UI from Task 31.
- Check `docs/ui-review-checklist.md`.

Acceptance criteria:

- Dashboard fetches real data.
- Review due card reflects backend data.
- Progress numbers reflect backend data.
- No mock dashboard data remains in production page logic.
- Layout closely matches screenshot if one exists.
- Typography hierarchy matches screenshot if one exists.
- Spacing hierarchy matches screenshot if one exists.
- Responsive behavior remains correct.

---

### Task 34 - Profile Page UI

Frontend page:

```txt
/profile
```

Goal:

Create profile page UI using mock data first.

Requirements:

- Analyze the reference screenshot before coding if one exists.
- The screenshot is the source of truth.
- Do not redesign the page.
- Profile card.
- Avatar.
- Name/email/username.
- Target band.
- Current level.
- Total XP.
- Current streak.
- Longest streak.
- Lessons completed.
- Words mastered.
- Achievements grid.
- Learning heatmap placeholder.
- Edit profile form if designed.
- If edit form is implemented, check `docs/validation-contracts.md` first.
- Check `docs/ui-review-checklist.md`.

Acceptance criteria:

- Profile page looks polished and responsive.
- Edit profile form follows validation contract if present.
- No backend integration required yet.
- Layout closely matches screenshot if one exists.
- Typography hierarchy matches screenshot if one exists.
- Spacing hierarchy matches screenshot if one exists.
- Responsive behavior remains correct.

---

### Task 35 - Profile APIs

Backend endpoints:

```txt
GET   /api/v1/me/profile
PATCH /api/v1/me/profile
```

Goal:

Return and update user profile.

Requirements:

- Protected endpoints.
- Use authenticated user ID from JWT context.
- Return:
  - user profile
  - stats
  - achievements
  - activity summary if available
- For update:
  - Check `docs/validation-contracts.md`.
  - If profile update contract does not exist, add it before implementing.
  - Backend request DTO validation must match the profile validation contract.
  - Allow updating name, username, target_band, timezone, locale if supported.
- Return field-level validation errors for validation failures.

Acceptance criteria:

- Profile loads real data.
- Profile update works.
- Validation matches contract.
- No hardcoded user ID.

---

### Task 36 - Connect Profile Page to API

Frontend page:

```txt
/profile
```

Goal:

Replace profile mock data with real backend data.

Requirements:

- Analyze the reference screenshot before coding if one exists.
- The screenshot is the source of truth.
- Do not redesign the page.
- Add profile API functions in `frontend/src/api/profile.ts`.
- Add shared profile/user types in:
  - `frontend/src/types/profile.ts`
  - `frontend/src/types/user.ts`
- Use axios through shared API client.
- Do not use fetch.
- If edit form exists:
  - use Zod schema matching `docs/validation-contracts.md`
  - show field-level errors
- Show loading, error, and empty states.
- Check `docs/ui-review-checklist.md`.

Acceptance criteria:

- Profile page loads real API data.
- Profile update works if implemented.
- Validation matches contract.
- No mock profile data remains in production page logic.
- Layout closely matches screenshot if one exists.
- Typography hierarchy matches screenshot if one exists.
- Spacing hierarchy matches screenshot if one exists.
- Responsive behavior remains correct.

---

## 10. Milestone 9 - Gamification and Activity

### Task 37 - XP System

Goal:

Implement XP event tracking.

Requirements:

- Use `user_xp_events`.
- Award XP for:
  - flashcard review if defined
  - quiz correct answer if defined
  - lesson completed
  - achievement unlocked
- Update `users.total_xp`.
- Recalculate `users.level`.
- Use simple MVP formula:

```txt
level = floor(total_xp / 200) + 1
```

Acceptance criteria:

- XP events are recorded.
- User total XP updates.
- Level updates correctly.
- No duplicate XP for the same one-time event.

---

### Task 38 - Achievement System

Goal:

Unlock achievements based on user actions.

Requirements:

- Use achievements and user_achievements.
- Initial achievements:
  - First Lesson
  - 7 Day Streak
  - 100 Words Learned
  - Education Master
- Achievement unlock should be idempotent.
- Achievement unlock may create notification and XP event if supported.

Acceptance criteria:

- Achievement unlocks once.
- Profile can display unlocked achievements.
- No duplicate user_achievements records.

---

### Task 39 - Daily Activity and Streak

Goal:

Track daily learning activity and streak.

Requirements:

- Use daily_activities.
- Use user timezone when calculating activity date.
- Track:
  - words learned
  - words reviewed
  - quizzes taken
  - lessons done
  - active minutes if available
  - XP earned
- Update:
  - users.current_streak
  - users.longest_streak
  - users.last_active_at

Acceptance criteria:

- Studying today creates/updates daily activity.
- Current streak updates correctly.
- Longest streak updates correctly.
- Dashboard can read streak and daily stats.

---

### Task 40 - Notifications

Goal:

Implement basic notifications.

Requirements:

- Use notifications table.
- Create notifications for:
  - achievement unlocked
  - review due reminder if implemented
  - streak milestone if implemented
- Add endpoints if needed:
  - GET /api/v1/notifications
  - PATCH /api/v1/notifications/:id/read
- Frontend topbar notification button can show unread count if API exists.

Acceptance criteria:

- Notifications can be created.
- User can fetch notifications.
- User can mark notification as read.
- No unrelated user notifications are exposed.

---

## 11. Milestone 10 - Polish and QA

### Task 41 - Global Frontend Error Handling

Requirements:

- Improve axios interceptor behavior.
- Handle 401 globally.
- Show user-friendly errors.
- Avoid duplicate error handling.
- Keep auth logout behavior safe.

Acceptance criteria:

- Expired token redirects or logs out consistently.
- API errors show useful messages.
- No app crash on API failure.

---

### Task 42 - Loading, Empty, and Error States

Add or refine:

- Skeletons.
- Empty states.
- Error states.
- Toasts.

Requirements:

- Analyze the reference screenshot before coding if one exists.
- The screenshot is the source of truth.
- Do not redesign the page.

Acceptance criteria:

- No page feels broken while loading.
- Empty data has clear CTA.
- API errors show useful messages.
- All main pages follow `docs/ui-review-checklist.md`.
- Layout closely matches screenshot if one exists.
- Typography hierarchy matches screenshot if one exists.
- Spacing hierarchy matches screenshot if one exists.
- Responsive behavior remains correct.

---

### Task 43 - Final QA

Checklist:

- Register.
- Login.
- Dashboard.
- Roadmap.
- Topic detail.
- Lesson detail.
- Vocabulary list.
- Vocabulary detail.
- Flashcards.
- Review due.
- Quiz.
- Profile.
- Logout.

Acceptance criteria:

- Main happy path works end-to-end.
- Frontend uses axios only.
- No frontend code uses fetch.
- Frontend API data is fetched through React Query.
- Validation contracts are respected.
- Protected routes work.
- No hardcoded user ID.
- Backend does not expose password hash.
- Quiz correct answers are not exposed before submit.
