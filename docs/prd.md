# PRD - IELTS Vocabulary Learning Platform

## 1. Product Overview

### Product Name

IELTS Vocabulary Trainer

### Product Vision

Build a modern IELTS vocabulary learning platform that helps learners improve vocabulary systematically through band-based lessons, flashcards, quizzes, spaced repetition, and progress tracking.

Version 1 focuses only on vocabulary learning. Later versions can expand into IELTS Writing, Speaking, Reading, and Listening.

### Chosen Tech Stack

```txt
Frontend: React.js + TypeScript + Vite
Backend: Go + Gin
Database: MySQL
ORM: GORM
Auth: JWT
```

### Product Positioning

The product should feel like:

- Duolingo for motivation and gamification
- Anki for spaced repetition
- Notion/Linear for clean and premium SaaS UI
- IELTS-specific vocabulary roadmap based on target band score

### Core Value Proposition

Users do not just browse IELTS vocabulary lists. They follow a structured learning path, review words at the right time, test themselves, and track progress toward their target band.

---

## 2. Target Users

### Primary Users

IELTS learners targeting Band 0 - 7.0.

### Secondary Users

- Students preparing for university admission
- Working professionals preparing for migration or work abroad
- English learners who want academic vocabulary

### User Problems

1. They do not know which IELTS words to learn first.
2. They forget vocabulary quickly after learning.
3. They lack structure and motivation.
4. They do not know how to use vocabulary in IELTS Writing and Speaking.
5. Existing vocabulary lists feel boring and disconnected.

---

## 3. MVP Scope

Version 1 focuses on:

- Authentication
- Vocabulary roadmap by IELTS band
- Topic-based lessons
- Flashcard learning
- Quiz practice
- Spaced repetition review
- User progress tracking
- Streak and basic achievements

---

## 4. MVP Features

### 4.1 Landing Page

Purpose:

Convert visitors into registered users.

Sections:

- Hero
- Product benefits
- Feature cards
- IELTS band roadmap preview
- Testimonials
- Pricing placeholder
- Footer

Primary CTA:

- Start Learning

Secondary CTA:

- Take Placement Test

---

### 4.2 Authentication

Users can:

- Register with email/password
- Login
- Logout
- Access protected pages with JWT

MVP auth:

- Email/password
- JWT access token
- Store token in frontend securely enough for MVP
- Backend middleware protects private APIs

Recommended frontend storage for MVP:

- httpOnly cookie is preferred if implemented
- localStorage is acceptable for quick MVP but less secure

---

### 4.3 Dashboard

Dashboard shows:

- Welcome message
- Current streak
- Words learned today
- Review due
- Target band progress
- Continue learning button
- Quick actions
- Recent activity
- Vocabulary statistics

---

### 4.4 Placement Test

Optional but recommended for V1.

Purpose:

Estimate user's vocabulary level and recommend starting band.

Flow:

1. User answers 20-30 vocabulary questions.
2. Backend calculates score.
3. User receives estimated level.
4. User is assigned recommended roadmap.

For first MVP, this can be simplified:

- Static test questions
- Multiple choice only
- Local or backend scoring

---

### 4.5 Vocabulary Roadmap

Users can see band levels:

- Band 0 - 4.0
- Band 4.0
- Band 5.0
- Band 6.0
- Band 7.0
- Band 8.0

Each band contains topics:

- Education
- Technology
- Environment
- Health
- Crime
- Government
- Culture
- Work
- Travel
- Society

Each topic contains lessons.

Roadmap behavior:

- Some lessons are unlocked by default.
- Completed lessons unlock next lessons.
- Locked lessons show lock icon.
- Progress is displayed visually.

---

### 4.6 Lesson Detail

Each lesson contains 10-20 vocabulary items.

Each vocabulary item includes:

- Word
- IPA
- Part of speech
- Vietnamese meaning
- English meaning
- Example sentence
- Synonyms
- Collocations
- Difficulty
- IELTS usage note

Lesson actions:

- Start Flashcards
- Start Quiz
- Mark as completed after passing quiz

---

### 4.7 Flashcard Learning

Flashcard flow:

1. Show word.
2. User clicks show meaning.
3. Show meaning, IPA, example, synonyms.
4. User rates memory:

- Again
- Hard
- Good
- Easy

Result:

- User progress is updated.
- Next review date is calculated.

---

### 4.8 Quiz

Quiz types for V1:

- Choose correct meaning
- Choose correct word for sentence
- Match word with meaning

Simplify MVP:

- Multiple choice only

Quiz completion:

- Show score
- Show correct/wrong answers
- If score >= required score, mark lesson as completed

Default passing score:

```txt
80%
```

---

### 4.9 Vocabulary Detail

Vocabulary detail page includes:

- Word
- IPA
- Audio pronunciation placeholder
- Meanings
- Examples
- Synonyms
- Antonyms
- Collocations
- IELTS Writing usage
- Related words
- User progress status

---

### 4.10 Spaced Repetition

Each vocabulary has review status per user.

Rating rules for MVP:

- Again: review today
- Hard: review tomorrow
- Good: review in 3 days
- Easy: review in 7 days

Later:

- Implement FSRS or SM-2 algorithm

---

### 4.11 User Profile

Profile page shows:

- Name
- Email
- Target band
- Total words learned
- Current streak
- Longest streak
- Total lessons completed
- Vocabulary mastery %
- Achievements
- Learning heatmap placeholder

---

## 5. Out of Scope for V1

Do not build these in V1 unless MVP is complete:

- Full IELTS Writing checker
- Speaking practice with voice
- Listening tests
- Reading passages
- Real payment integration
- Admin dashboard
- User-generated vocabulary lists
- Community/forum
- Native mobile app
- Complex AI features

---

## 6. Success Metrics

### Activation

- User completes first lesson.
- User completes first quiz.
- User reviews at least 10 words.

### Engagement

- Daily active users
- Average review completion
- Lessons completed per user
- Streak length

### Learning

- Words mastered
- Quiz accuracy
- Review retention rate

---

## 7. MVP Acceptance Criteria

MVP is acceptable when:

1. User can register/login.
2. User can see dashboard.
3. User can browse roadmap.
4. User can open a lesson.
5. User can study flashcards.
6. User can complete quiz.
7. User progress is saved in MySQL.
8. Review due words appear correctly.
9. UI is responsive on desktop and mobile.
10. Seed data provides at least 3 bands, 5 topics, 10 lessons, and 100 vocabulary items.

---

## 8. Product Principles

1. Learning should feel structured, not random.
2. Every screen should guide the user to the next action.
3. Vocabulary should be IELTS-specific.
4. Progress should be visible and motivating.
5. V1 should be simple but extensible.
6. Avoid over-engineering early.
7. Prefer clear data models over clever UI tricks.
8. Keep frontend and backend separated cleanly.
