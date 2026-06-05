# User Flow - IELTS Vocabulary Platform

## 1. Main User Journey

```txt
Visitor
→ Landing Page
→ Register / Login
→ Optional Placement Test
→ Dashboard
→ Roadmap
→ Lesson
→ Flashcards
→ Quiz
→ Lesson Completed
→ Review Due
→ Continue Learning
```

---

## 2. Guest Flow

### 2.1 Visit Landing Page

User sees:

- Product value proposition
- IELTS vocabulary roadmap preview
- Flashcards preview
- Quiz preview
- CTA buttons

Actions:

- Start Learning
- Take Placement Test
- Login

### 2.2 Start Learning

If not logged in:

```txt
Start Learning
→ Register Page
→ Create Account
→ Dashboard
```

### 2.3 Take Placement Test

```txt
Take Placement Test
→ Answer 20-30 questions
→ View estimated band
→ Register to save result
→ Dashboard
```

For MVP, placement test can require login first if easier.

---

## 3. Authentication Flow

### Register

```txt
Register Page
→ Enter name, email, password
→ React validates form
→ POST /api/v1/auth/register
→ Backend creates user
→ Backend returns JWT
→ Frontend saves JWT
→ Redirect to Dashboard
```

### Login

```txt
Login Page
→ Enter email/password
→ POST /api/v1/auth/login
→ Backend validates credentials
→ Backend returns JWT
→ Frontend saves JWT
→ Redirect to Dashboard
```

### Logout

```txt
User menu
→ Remove JWT
→ Redirect to Landing Page
```

---

## 4. Dashboard Flow

Dashboard is the user's home base.

User sees:

- Today's learning goal
- Review due
- Current streak
- Continue learning card
- Target band progress
- Recent activity

Primary actions:

```txt
Continue Learning → Next recommended lesson
Review Due → Review flashcards
Start Quiz → Quiz practice
Roadmap → Full roadmap
```

Recommendation logic for Continue Learning:

1. If review due count > 0, suggest review.
2. Else if current lesson in progress, continue current lesson.
3. Else suggest next unlocked lesson.

---

## 5. Roadmap Flow

```txt
Dashboard
→ Roadmap
→ Select Band
→ Select Topic
→ Select Lesson
→ Lesson Detail
```

Roadmap states:

### Locked Lesson

User sees:

```txt
This lesson is locked. Complete the previous lesson first.
```

### Unlocked Lesson

User can open lesson.

### In Progress Lesson

User can continue.

### Completed Lesson

User sees completion mark and score.

---

## 6. Lesson Flow

```txt
Lesson Detail
→ View vocabulary list
→ Start Flashcards
→ Learn all cards
→ Start Quiz
→ Submit Quiz
→ Pass or Retry
```

### Lesson Detail Page

User sees:

- Lesson title
- Topic
- Estimated time
- Required score
- Vocabulary list
- Current progress
- Start/Continue buttons

### Start Lesson

Frontend calls:

```txt
POST /api/v1/lessons/:lessonId/start
```

Backend:

- Creates or updates UserLessonProgress
- Status becomes IN_PROGRESS

---

## 7. Flashcard Flow

```txt
Flashcard Page
→ GET /api/v1/lessons/:lessonId/flashcards
→ Show word
→ User clicks Show Meaning
→ User rates memory
→ POST /api/v1/flashcards/review
→ Next card
→ Finish session
```

Card front:

- Word
- IPA optional
- Part of speech

Card back:

- Meaning VI
- Meaning EN
- Example
- Synonyms
- Collocations

Rating actions:

- Again
- Hard
- Good
- Easy

After rating:

- Backend updates UserVocabularyProgress
- Backend calculates nextReviewAt
- Frontend updates session progress

When all cards are done:

```txt
Session Complete
→ Show summary
→ Start Quiz
```

---

## 8. Quiz Flow

```txt
Quiz Page
→ GET /api/v1/lessons/:lessonId/quiz
→ Answer questions
→ POST /api/v1/lessons/:lessonId/quiz/submit
→ See result
→ Pass or Retry
```

### During Quiz

User sees:

- Current question number
- Question text
- Options
- Progress
- Timer optional

### Submit Quiz

Backend:

- Calculates score
- Saves attempt
- Shows correct answers only after submit
- Updates lesson progress if passed

### If Passed

```txt
Score >= requiredScore
→ Mark lesson completed
→ Unlock next lesson
→ Show success screen
```

### If Failed

```txt
Score < requiredScore
→ Keep lesson in progress
→ Suggest review flashcards
→ Retry quiz
```

---

## 9. Review Due Flow

```txt
Dashboard
→ Review Due
→ GET /api/v1/reviews/due
→ Flashcard Review Session
→ Rate cards
→ Update next review dates
→ Session summary
```

Review due list includes:

```txt
next_review_at <= now
```

If no review due:

Show empty state:

```txt
Great job! No words due for review today.
```

CTA:

```txt
Go to Roadmap
```

---

## 10. Vocabulary Search Flow

```txt
Dashboard
→ Vocabulary
→ Search word
→ GET /api/v1/vocabularies?q=...
→ Open detail page
```

Vocabulary list supports:

- Search
- Filter by topic
- Filter by band
- Filter by difficulty
- Filter by status

Vocabulary detail lets user:

- View details
- Save word
- Start quick review

---

## 11. Profile Flow

```txt
Dashboard
→ Profile
→ GET /api/v1/me/profile
→ View stats
→ Edit profile
→ PATCH /api/v1/me/profile
→ Save
```

User can edit:

- Name
- Target band

User can view:

- Total words learned
- Streak
- Achievements
- Lesson count
- Mastery percentage

---

## 12. Main Happy Path

```txt
User registers
→ Dashboard
→ Roadmap
→ Opens Band 6.0 Education Lesson 1
→ Starts flashcards
→ Reviews 10 words
→ Takes quiz
→ Scores 90%
→ Lesson completed
→ Next lesson unlocked
→ Dashboard shows updated progress
```

---

## 13. Main Edge Cases

### User tries locked lesson

Show:

```txt
This lesson is locked. Complete the previous lesson first.
```

### User submits quiz with missing answers

Show validation:

```txt
Please answer all questions before submitting.
```

### User has no progress record for vocabulary

Backend creates progress record automatically on first review.

### User revisits completed lesson

Allow:

- Review vocabulary
- Retake quiz
- Practice flashcards

Do not remove completed status if later quiz score is lower.
