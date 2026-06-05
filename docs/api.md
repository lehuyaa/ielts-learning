# API Design v2 - Screenshot-Aligned Contracts

This document extends `docs/api.md` after reviewing all screenshots in `docs/screenshots` and comparing them with `docs/database.md` and `docs/architecture.md`.

Use the same backend stack and module pattern from `docs/architecture.md`:

```txt
Go + Gin
GORM
MySQL
JWT auth
Standard API response shape
```

Base URL:

```txt
http://localhost:8080/api/v1
```

---

## 1. Main API Findings

### 1.1 Missing APIs

The screenshots require these APIs beyond `docs/api.md`:

```txt
GET   /api/v1/public/landing
GET   /api/v1/public/pricing

GET   /api/v1/search/vocabularies

GET   /api/v1/dashboard
GET   /api/v1/dashboard/activity

GET   /api/v1/notifications
PATCH /api/v1/notifications/:notificationId/read
PATCH /api/v1/notifications/read-all

GET   /api/v1/challenges/daily
POST  /api/v1/challenges/daily/claim

GET   /api/v1/profile
PATCH /api/v1/profile
GET   /api/v1/profile/overview
GET   /api/v1/profile/achievements
GET   /api/v1/profile/calendar

POST  /api/v1/lessons/:lessonId/vocabularies/:vocabularyId/learned

POST  /api/v1/lessons/:lessonId/quiz/sessions
GET   /api/v1/quiz/sessions/:sessionId
POST  /api/v1/quiz/sessions/:sessionId/answers
POST  /api/v1/quiz/sessions/:sessionId/finish

GET   /api/v1/placement-test
POST  /api/v1/placement-test/attempts
POST  /api/v1/placement-test/attempts/:attemptId/answers
POST  /api/v1/placement-test/attempts/:attemptId/complete

GET   /api/v1/subscription/plans
GET   /api/v1/me/subscription

POST  /api/v1/ai/examples
```

### 1.2 Existing APIs That Need Expanded Response Data

These existing APIs should remain, but their response should include screenshot-driven fields:

```txt
GET /api/v1/auth/me
GET /api/v1/dashboard/summary
GET /api/v1/roadmap
GET /api/v1/lessons/:lessonId
GET /api/v1/lessons/:lessonId/flashcards
GET /api/v1/reviews/due
GET /api/v1/lessons/:lessonId/quiz
POST /api/v1/lessons/:lessonId/quiz/submit
```

### 1.3 APIs That Can Stay Static for MVP

The landing page screenshots show marketing content, testimonials, and pricing. These can be static frontend data in MVP.

If dynamic content is not needed, defer:

```txt
GET /api/v1/public/landing
GET /api/v1/public/pricing
GET /api/v1/subscription/plans
GET /api/v1/me/subscription
```

---

## 2. Standard Response Shape

Keep v1 shape:

### Success

```json
{
  "data": {}
}
```

### Error

```json
{
  "error": {
    "code": "INVALID_INPUT",
    "message": "Invalid request input"
  }
}
```

### Validation Error

Use this response shape for request validation failures:

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Validation failed",
    "fields": {
      "email": "Please enter a valid email address",
      "password": "Password must be at least 8 characters"
    }
  }
}
```

Notes:

- `fields` maps request field names to user-facing validation messages.
- Validation rules must match `docs/validation-contracts.md`.
- Frontend validation is for user experience; backend validation is the final authority.

### Pagination

```json
{
  "data": {
    "items": [],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 100,
      "totalPages": 5
    }
  }
}
```

---

## 3. Error Codes

Keep v1 error codes and add:

```txt
VALIDATION_ERROR
QUIZ_SESSION_NOT_FOUND
QUIZ_SESSION_COMPLETED
QUESTION_ALREADY_ANSWERED
QUESTION_NOT_IN_SESSION
OPTION_NOT_IN_QUESTION
PLACEMENT_ATTEMPT_NOT_FOUND
DAILY_CHALLENGE_NOT_FOUND
DAILY_CHALLENGE_NOT_COMPLETED
DAILY_CHALLENGE_ALREADY_CLAIMED
NOTIFICATION_NOT_FOUND
SUBSCRIPTION_REQUIRED
AI_GENERATION_FAILED
RATE_LIMITED
```

---

## 4. Public APIs

### GET /api/v1/public/landing

Optional for MVP. Returns public landing content if not hardcoded in frontend.

Response:

```json
{
  "data": {
    "stats": [
      { "label": "IELTS Words", "value": "5,000+" },
      { "label": "Lessons", "value": "100+" },
      { "label": "Topics", "value": "50+" },
      { "label": "Band Range", "value": "5.0 -> 8.5" }
    ],
    "features": [
      {
        "title": "Smart Flashcards",
        "description": "Interactive flashcards with IPA pronunciation, examples, and synonyms.",
        "icon": "book-open"
      }
    ],
    "testimonials": [
      {
        "name": "Priya S.",
        "initials": "PS",
        "achievedBand": 7.5,
        "rating": 5,
        "quote": "I went from Band 6.0 to 7.5 in 3 months."
      }
    ]
  }
}
```

### GET /api/v1/public/pricing

Optional for MVP.

Response:

```json
{
  "data": {
    "plans": [
      {
        "code": "FREE",
        "name": "Free",
        "priceCents": 0,
        "currency": "USD",
        "billingInterval": "FOREVER",
        "description": "Get started with essential vocabulary tools.",
        "features": [
          "200 IELTS Words",
          "5 Topics",
          "Basic Flashcards",
          "Daily Quiz (5 questions)",
          "Progress Tracking"
        ],
        "isPopular": false
      },
      {
        "code": "PRO",
        "name": "Pro",
        "priceCents": 900,
        "currency": "USD",
        "billingInterval": "MONTH",
        "description": "Full access to all vocabulary and learning tools.",
        "features": [
          "5,000+ IELTS Words",
          "All 50+ Topics",
          "Unlimited Flashcards",
          "Unlimited Quizzes",
          "Spaced Repetition",
          "AI Example Generator"
        ],
        "isPopular": true
      }
    ]
  }
}
```

---

## 5. Auth APIs

### GET /api/v1/auth/me

Protected.

Expanded response:

```json
{
  "data": {
    "id": 1,
    "email": "alex@example.com",
    "name": "Alex Johnson",
    "username": "alexj",
    "avatarUrl": null,
    "targetBand": 7.0,
    "currentBand": 6.5,
    "startingBand": 5.0,
    "recommendedBand": 6.0,
    "totalXp": 2840,
    "level": 18,
    "levelTitle": "Expert",
    "currentStreak": 14,
    "longestStreak": 21,
    "createdAt": "2024-10-01T08:00:00+07:00"
  }
}
```

Notes:

- Never expose `passwordHash`.
- `levelTitle` can be computed from `level`.

---

## 6. Search APIs

### GET /api/v1/search/vocabularies

Protected.

Used by the global dashboard search input.

Query params:

```txt
q
limit
```

Response:

```json
{
  "data": {
    "items": [
      {
        "id": 1,
        "word": "sustainable",
        "ipa": "/səˈsteɪnəbl/",
        "partOfSpeech": "adjective",
        "meaningEn": "Able to be maintained without depleting natural resources",
        "topic": {
          "id": 3,
          "title": "Environment"
        },
        "targetBand": 7.0,
        "difficulty": "INTERMEDIATE",
        "status": "REVIEW"
      }
    ]
  }
}
```

---

## 7. Dashboard APIs

### GET /api/v1/dashboard

Protected.

Replaces or expands `GET /api/v1/dashboard/summary`.

Response:

```json
{
  "data": {
    "user": {
      "name": "Alex",
      "avatarUrl": null,
      "initials": "A",
      "targetBand": 7.0,
      "currentBand": 6.5,
      "totalXp": 2840,
      "level": 18,
      "levelTitle": "Expert",
      "currentStreak": 14
    },
    "notifications": {
      "unreadCount": 1
    },
    "stats": {
      "wordsLearnedToday": 24,
      "wordsLearnedDelta": 6,
      "reviewDue": 12,
      "currentStreak": 14,
      "accuracyPercentage": 84,
      "accuracyDelta": 3,
      "totalWordsLearned": 847,
      "lessonsCompleted": 32,
      "masteryPercentage": 68
    },
    "targetBandProgress": {
      "startingBand": 5.0,
      "currentBand": 6.5,
      "targetBand": 7.0,
      "progressPercentage": 65
    },
    "weeklyWords": {
      "changePercentage": 18,
      "items": [
        { "date": "2026-06-01", "label": "Mon", "wordsLearned": 18 },
        { "date": "2026-06-02", "label": "Tue", "wordsLearned": 24 }
      ]
    },
    "quickActions": {
      "continueLearning": {
        "type": "LESSON",
        "lessonId": 10,
        "label": "Continue Learning"
      },
      "flashcardsUrl": "/reviews",
      "dailyQuizUrl": "/quiz",
      "reviewDueUrl": "/reviews"
    },
    "dailyChallenge": {
      "title": "Daily Challenge",
      "description": "Learn 10 words before midnight!",
      "metric": "WORDS_LEARNED",
      "progressValue": 6,
      "targetValue": 10,
      "xpReward": 50,
      "completed": false,
      "claimed": false
    },
    "recentActivity": [
      {
        "type": "VOCABULARY_REVIEWED",
        "title": "Sustainability",
        "subtitle": "Environment",
        "rating": "EASY",
        "createdAt": "2026-06-05T08:00:00+07:00"
      }
    ],
    "dueForReview": [
      {
        "id": 1,
        "word": "Sophisticated",
        "initial": "S",
        "targetBand": 7.0,
        "dueLabel": "Now",
        "nextReviewAt": "2026-06-05T08:00:00+07:00"
      }
    ],
    "recentAchievements": [
      {
        "id": 1,
        "code": "SEVEN_DAY_STREAK",
        "title": "7 Day Streak",
        "icon": "flame",
        "unlocked": true,
        "unlockedAt": "2026-06-01T08:00:00+07:00"
      }
    ]
  }
}
```

### GET /api/v1/dashboard/activity

Protected.

Returns a fuller activity feed when the user clicks "View all".

Query params:

```txt
page
limit
```

Response:

```json
{
  "data": {
    "items": [
      {
        "type": "VOCABULARY_REVIEWED",
        "title": "Proliferate",
        "subtitle": "Technology",
        "rating": "GOOD",
        "createdAt": "2026-06-05T05:00:00+07:00"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 80,
      "totalPages": 4
    }
  }
}
```

---

## 8. Notifications APIs

### GET /api/v1/notifications

Protected.

Query params:

```txt
unreadOnly
page
limit
```

Response:

```json
{
  "data": {
    "items": [
      {
        "id": 1,
        "type": "ACHIEVEMENT",
        "title": "Achievement unlocked",
        "body": "You earned 7 Day Streak.",
        "actionUrl": "/profile?tab=achievements",
        "readAt": null,
        "createdAt": "2026-06-05T08:00:00+07:00"
      }
    ],
    "unreadCount": 1,
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 1,
      "totalPages": 1
    }
  }
}
```

### PATCH /api/v1/notifications/:notificationId/read

Protected.

Response:

```json
{
  "data": {
    "id": 1,
    "readAt": "2026-06-05T08:10:00+07:00"
  }
}
```

### PATCH /api/v1/notifications/read-all

Protected.

Response:

```json
{
  "data": {
    "updatedCount": 4
  }
}
```

---

## 9. Daily Challenge APIs

### GET /api/v1/challenges/daily

Protected.

Response:

```json
{
  "data": {
    "date": "2026-06-05",
    "title": "Daily Challenge",
    "description": "Learn 10 words before midnight!",
    "metric": "WORDS_LEARNED",
    "progressValue": 6,
    "targetValue": 10,
    "xpReward": 50,
    "completed": false,
    "claimed": false
  }
}
```

### POST /api/v1/challenges/daily/claim

Protected.

Claims XP after challenge completion.

Response:

```json
{
  "data": {
    "claimed": true,
    "xpAwarded": 50,
    "totalXp": 2890,
    "level": 18
  }
}
```

Errors:

```txt
DAILY_CHALLENGE_NOT_COMPLETED
DAILY_CHALLENGE_ALREADY_CLAIMED
```

---

## 10. Roadmap APIs

### GET /api/v1/roadmap

Protected.

Expanded response:

```json
{
  "data": {
    "summary": {
      "bandRangeLabel": "Band 5.0 -> 8.5",
      "topicsCompleted": 3,
      "totalTopics": 19,
      "currentBand": 6.5,
      "wordsMastered": 847,
      "currentStreak": 14
    },
    "course": {
      "id": 1,
      "title": "IELTS Vocabulary Roadmap",
      "bandMin": 5.0,
      "bandMax": 8.5
    },
    "bandLevels": [
      {
        "id": 1,
        "bandScore": 5.0,
        "title": "Band 5.0",
        "status": "COMPLETED",
        "statusLabel": "Complete",
        "isLocked": false,
        "lockedReason": null,
        "topicsCompleted": 3,
        "totalTopics": 3,
        "progressPercentage": 100,
        "topics": [
          {
            "id": 1,
            "title": "Education",
            "slug": "education",
            "emoji": "graduation-cap",
            "progressPercentage": 100,
            "lessonsCompleted": 8,
            "totalLessons": 8,
            "status": "COMPLETED",
            "isLocked": false
          }
        ]
      },
      {
        "id": 3,
        "bandScore": 7.0,
        "title": "Band 7.0",
        "status": "LOCKED",
        "statusLabel": "Locked",
        "isLocked": true,
        "lockedReason": "Complete previous band",
        "topicsCompleted": 0,
        "totalTopics": 0,
        "progressPercentage": 0,
        "topics": []
      }
    ],
    "masterMilestone": {
      "title": "Band 8.5+ Master",
      "description": "Complete all bands to unlock",
      "isLocked": true
    }
  }
}
```

---

## 11. Lesson APIs

### GET /api/v1/lessons/:lessonId

Protected.

Expanded response:

```json
{
  "data": {
    "id": 1,
    "title": "Education Vocabulary",
    "description": "Core education vocabulary for IELTS.",
    "requiredScore": 80,
    "estimatedMinutes": 25,
    "bandMin": 5.5,
    "bandMax": 7.0,
    "xpReward": 350,
    "progress": {
      "status": "COMPLETED",
      "wordsLearned": 7,
      "totalWords": 7,
      "progressPercentage": 100,
      "score": 920,
      "bestScore": 920,
      "personalBest": 850,
      "completedAt": "2026-06-05T08:00:00+07:00"
    },
    "topic": {
      "id": 1,
      "title": "Education",
      "slug": "education"
    },
    "vocabularies": [
      {
        "id": 1,
        "word": "Accommodate",
        "ipa": "/əˈkɒmədeɪt/",
        "audioUrl": "/audio/accommodate.mp3",
        "partOfSpeech": "verb",
        "meaningEn": "To provide space or housing for; to adapt or adjust to...",
        "shortDefinition": "To provide space or housing for; to adapt or adjust to...",
        "difficulty": "INTERMEDIATE",
        "targetBand": 6.0,
        "status": "MASTERED",
        "learned": true,
        "learnedAt": "2026-06-05T08:00:00+07:00"
      }
    ]
  }
}
```

### POST /api/v1/lessons/:lessonId/vocabularies/:vocabularyId/learned

Protected.

Marks a vocabulary item as learned from the lesson list or flashcard flow.

Response:

```json
{
  "data": {
    "vocabularyId": 1,
    "status": "LEARNING",
    "learned": true,
    "learnedAt": "2026-06-05T08:00:00+07:00",
    "lessonProgress": {
      "wordsLearned": 7,
      "totalWords": 7,
      "progressPercentage": 100
    },
    "dailyChallenge": {
      "progressValue": 6,
      "targetValue": 10,
      "completed": false
    }
  }
}
```

---

## 12. Vocabulary APIs

### GET /api/v1/vocabularies

Protected.

Keep v1 query params and add:

```txt
targetBand
learned
mastered
```

Each item should include:

```json
{
  "id": 1,
  "word": "sustainable",
  "ipa": "/səˈsteɪnəbl/",
  "partOfSpeech": "adjective",
  "meaningEn": "Able to be maintained without depleting natural resources",
  "meaningVi": "bền vững",
  "difficulty": "INTERMEDIATE",
  "targetBand": 7.0,
  "status": "REVIEW",
  "topic": {
    "id": 3,
    "title": "Environment"
  }
}
```

---

## 13. Flashcard APIs

### GET /api/v1/lessons/:lessonId/flashcards

Protected.

Expanded response:

```json
{
  "data": {
    "lessonId": 1,
    "topic": {
      "id": 3,
      "title": "Environment"
    },
    "bandLabel": "Band 7",
    "progress": {
      "done": 0,
      "remaining": 5,
      "total": 5
    },
    "items": [
      {
        "id": 1,
        "word": "Sustainable",
        "ipa": "/səˈsteɪnəbl/",
        "audioUrl": "/audio/sustainable.mp3",
        "partOfSpeech": "adjective",
        "meaningEn": "Able to be maintained at a certain rate or level; avoiding depletion of natural resources.",
        "meaningVi": "bền vững",
        "exampleSentence": "Governments must develop sustainable policies to address climate change effectively.",
        "synonyms": ["Renewable", "Viable", "Eco-friendly", "Long-term"],
        "targetBand": 7.0,
        "status": "NEW"
      }
    ]
  }
}
```

### GET /api/v1/reviews/due

Protected.

Expanded response:

```json
{
  "data": {
    "items": [
      {
        "id": 1,
        "word": "Sophisticated",
        "initial": "S",
        "ipa": "/səˈfɪstɪkeɪtɪd/",
        "meaningEn": "Having a refined knowledge of the world.",
        "meaningVi": "tinh tế",
        "exampleSentence": "The essay uses sophisticated vocabulary.",
        "targetBand": 7.0,
        "dueLabel": "Now",
        "nextReviewAt": "2026-06-05T08:00:00+07:00"
      }
    ],
    "count": 12
  }
}
```

### POST /api/v1/flashcards/review

Protected.

Keep v1 request and expand response:

```json
{
  "data": {
    "vocabularyId": 1,
    "status": "REVIEW",
    "rating": "GOOD",
    "reviewCount": 2,
    "nextReviewAt": "2026-06-08T08:00:00+07:00",
    "nextReviewLabel": "3 days",
    "xpAwarded": 5,
    "totalXp": 2845,
    "dailyChallenge": {
      "progressValue": 7,
      "targetValue": 10,
      "completed": false
    }
  }
}
```

---

## 14. Quiz APIs

The v1 `GET /lessons/:lessonId/quiz` and `POST /lessons/:lessonId/quiz/submit` work for submit-at-end quizzes. The screenshots show immediate feedback after each answer, so v2 adds quiz sessions.

### GET /api/v1/lessons/:lessonId/quiz

Protected.

Still useful for previewing a quiz before starting a session. Do not expose correct answers.

Expanded response:

```json
{
  "data": {
    "lessonId": 1,
    "topic": {
      "id": 3,
      "title": "Environment"
    },
    "bandLabel": "Band 7",
    "totalQuestions": 5,
    "timeLimitSeconds": 150,
    "pointsPerQuestion": 20,
    "questions": [
      {
        "id": 1,
        "type": "MEANING_CHOICE",
        "question": "What does \"sustainable\" mean?",
        "points": 20,
        "timeLimitSeconds": 30,
        "options": [
          {
            "id": 1,
            "label": "A",
            "content": "Able to be maintained without depleting natural resources"
          }
        ]
      }
    ]
  }
}
```

### POST /api/v1/lessons/:lessonId/quiz/sessions

Protected.

Starts a live quiz session.

Request:

```json
{
  "questionCount": 5
}
```

Response:

```json
{
  "data": {
    "sessionId": 101,
    "lessonId": 1,
    "status": "IN_PROGRESS",
    "currentQuestionIndex": 0,
    "totalQuestions": 5,
    "points": 0,
    "startedAt": "2026-06-05T08:00:00+07:00",
    "question": {
      "id": 1,
      "type": "MEANING_CHOICE",
      "topic": {
        "id": 3,
        "title": "Environment"
      },
      "bandLabel": "Band 7",
      "questionNumber": 1,
      "question": "What does \"sustainable\" mean?",
      "points": 20,
      "timeLimitSeconds": 30,
      "options": [
        {
          "id": 1,
          "label": "A",
          "content": "Able to be maintained without depleting natural resources"
        },
        {
          "id": 2,
          "label": "B",
          "content": "Relating to the study of living organisms"
        }
      ]
    }
  }
}
```

### GET /api/v1/quiz/sessions/:sessionId

Protected.

Returns current session state.

Response:

```json
{
  "data": {
    "sessionId": 101,
    "lessonId": 1,
    "status": "IN_PROGRESS",
    "currentQuestionIndex": 2,
    "totalQuestions": 5,
    "points": 40,
    "progress": [
      { "questionId": 1, "status": "CORRECT" },
      { "questionId": 2, "status": "CORRECT" },
      { "questionId": 3, "status": "CURRENT" },
      { "questionId": 4, "status": "PENDING" },
      { "questionId": 5, "status": "PENDING" }
    ],
    "question": {
      "id": 3,
      "type": "USAGE_CHOICE",
      "topic": {
        "id": 2,
        "title": "Technology"
      },
      "bandLabel": "Band 7.5",
      "questionNumber": 3,
      "question": "Which sentence uses \"proliferate\" correctly?",
      "points": 20,
      "timeLimitSeconds": 30,
      "options": [
        {
          "id": 9,
          "label": "A",
          "content": "Scientists proliferate the new vaccine to ensure safety."
        }
      ]
    }
  }
}
```

### POST /api/v1/quiz/sessions/:sessionId/answers

Protected.

Submits one answer and returns immediate feedback. This is the first time correct-answer data may be returned for that question.

Request:

```json
{
  "questionId": 3,
  "optionId": 12,
  "timeSpentSeconds": 8
}
```

Response:

```json
{
  "data": {
    "sessionId": 101,
    "questionId": 3,
    "selectedOptionId": 12,
    "correctOptionId": 11,
    "isCorrect": false,
    "pointsAwarded": 0,
    "points": 40,
    "feedback": {
      "title": "Incorrect - The answer is option C",
      "message": "The correct answer is: \"Online misinformation has proliferated rapidly since 2016.\""
    },
    "options": [
      {
        "id": 11,
        "label": "C",
        "content": "Online misinformation has proliferated rapidly since 2016.",
        "isCorrect": true,
        "isSelected": false
      },
      {
        "id": 12,
        "label": "D",
        "content": "Doctors proliferate with their medical opinions on the matter.",
        "isCorrect": false,
        "isSelected": true
      }
    ],
    "progress": [
      { "questionId": 1, "status": "CORRECT" },
      { "questionId": 2, "status": "CORRECT" },
      { "questionId": 3, "status": "INCORRECT" },
      { "questionId": 4, "status": "PENDING" },
      { "questionId": 5, "status": "PENDING" }
    ],
    "nextQuestionId": 4
  }
}
```

Errors:

```txt
QUIZ_SESSION_COMPLETED
QUESTION_ALREADY_ANSWERED
QUESTION_NOT_IN_SESSION
OPTION_NOT_IN_QUESTION
```

### POST /api/v1/quiz/sessions/:sessionId/finish

Protected.

Finishes the session and writes a `user_quiz_attempts` row.

Response:

```json
{
  "data": {
    "sessionId": 101,
    "lessonId": 1,
    "status": "COMPLETED",
    "score": 80,
    "points": 80,
    "totalQuestions": 5,
    "correctAnswers": 4,
    "wrongAnswers": 1,
    "passed": true,
    "requiredScore": 80,
    "xpEarned": 120,
    "lessonProgress": {
      "status": "COMPLETED",
      "bestScore": 920,
      "bestXp": 350,
      "completedAt": "2026-06-05T08:05:00+07:00"
    },
    "unlockedAchievements": [
      {
        "id": 3,
        "code": "FIRST_QUIZ",
        "title": "First Quiz",
        "icon": "zap"
      }
    ],
    "unlockedNextLesson": {
      "id": 2,
      "title": "Education Vocabulary 2"
    }
  }
}
```

### POST /api/v1/lessons/:lessonId/quiz/submit

Protected.

Keep for classic submit-at-end quiz support. Internally, this can create a completed `quiz_session` or write directly to `user_quiz_attempts`.

Important:

- Do not expose correct answers before submission.
- Mark lesson completed only when score is at least `requiredScore`.
- Do not downgrade a completed lesson on a worse later attempt.

---

## 15. Profile APIs

### GET /api/v1/profile

Protected.

Returns the profile header and summary cards.

Response:

```json
{
  "data": {
    "user": {
      "id": 1,
      "name": "Alex Johnson",
      "username": "alexj",
      "email": "alex@example.com",
      "avatarUrl": null,
      "initials": "A",
      "memberSince": "2024-10-01T08:00:00+07:00",
      "currentBand": 6.5,
      "targetBand": 7.0,
      "totalXp": 2840,
      "level": 18,
      "levelTitle": "Expert",
      "currentLevelXp": 2840,
      "nextLevelXp": 3500,
      "xpUntilNextLevel": 660,
      "currentStreak": 14,
      "longestStreak": 21
    },
    "stats": {
      "wordsLearned": 847,
      "lessonsDone": 32,
      "masteredWords": 574,
      "learningWords": 183,
      "newWords": 90,
      "vocabularyMasteryPercentage": 68
    }
  }
}
```

### PATCH /api/v1/profile

Protected.

Request:

```json
{
  "name": "Alex Johnson",
  "username": "alexj",
  "avatarUrl": "https://example.com/avatar.png",
  "targetBand": 7.5,
  "timezone": "Asia/Pontianak",
  "locale": "en"
}
```

Response:

```json
{
  "data": {
    "id": 1,
    "name": "Alex Johnson",
    "username": "alexj",
    "avatarUrl": "https://example.com/avatar.png",
    "targetBand": 7.5,
    "timezone": "Asia/Pontianak",
    "locale": "en"
  }
}
```

Validation:

```txt
Profile update validation must be defined in docs/validation-contracts.md before implementation.
Backend DTO validation and frontend Zod validation must match that contract.
```

### GET /api/v1/profile/overview

Protected.

Returns topic mastery rows for the Overview tab.

Response:

```json
{
  "data": {
    "topics": [
      {
        "id": 1,
        "title": "Education",
        "progressPercentage": 85,
        "masteredWords": 120,
        "totalWords": 140
      },
      {
        "id": 3,
        "title": "Environment",
        "progressPercentage": 60,
        "masteredWords": 84,
        "totalWords": 140
      }
    ]
  }
}
```

### GET /api/v1/profile/achievements

Protected.

Returns unlocked and locked achievements for the Achievements tab.

Response:

```json
{
  "data": {
    "items": [
      {
        "id": 1,
        "code": "SEVEN_DAY_STREAK",
        "title": "7 Day Streak",
        "description": "Studied 7 days in a row",
        "icon": "flame",
        "unlocked": true,
        "unlockedAt": "2026-12-01T08:00:00+07:00",
        "progressValue": 7,
        "requirementValue": 7
      },
      {
        "id": 6,
        "code": "THIRTY_DAY_STREAK",
        "title": "30 Day Streak",
        "description": "Study for 30 consecutive days",
        "icon": "trophy",
        "unlocked": false,
        "unlockedAt": null,
        "progressValue": 14,
        "requirementValue": 30
      }
    ]
  }
}
```

### GET /api/v1/profile/calendar

Protected.

Query params:

```txt
days default 84
```

Response:

```json
{
  "data": {
    "range": {
      "from": "2026-03-14",
      "to": "2026-06-05",
      "days": 84
    },
    "summary": {
      "wordsLearned": 847,
      "averageWordsPerDay": 10.1
    },
    "items": [
      {
        "date": "2026-06-05",
        "wordsLearned": 24,
        "wordsReviewed": 12,
        "quizzesTaken": 1,
        "lessonsDone": 1,
        "xpEarned": 120,
        "intensity": 5
      }
    ]
  }
}
```

Notes:

- `intensity` should be 0-5 for heatmap color levels.
- Dates should respect `users.timezone`.

---

## 16. Placement Test APIs

### GET /api/v1/placement-test

Public or protected. MVP can require auth if easier.

Do not expose correct answers.

Response:

```json
{
  "data": {
    "id": 1,
    "title": "IELTS Vocabulary Placement Test",
    "description": "Estimate your vocabulary band.",
    "questionCount": 20,
    "questions": [
      {
        "id": 1,
        "question": "What does \"sustainable\" mean?",
        "bandScore": 7.0,
        "options": [
          {
            "id": 1,
            "label": "A",
            "content": "Able to continue without depleting resources"
          }
        ]
      }
    ]
  }
}
```

### POST /api/v1/placement-test/attempts

Public or protected.

Request:

```json
{
  "placementTestId": 1
}
```

Response:

```json
{
  "data": {
    "attemptId": 1,
    "placementTestId": 1,
    "startedAt": "2026-06-05T08:00:00+07:00"
  }
}
```

### POST /api/v1/placement-test/attempts/:attemptId/answers

Public or protected.

Request:

```json
{
  "questionId": 1,
  "optionId": 1
}
```

Response:

```json
{
  "data": {
    "questionId": 1,
    "answered": true
  }
}
```

### POST /api/v1/placement-test/attempts/:attemptId/complete

Public or protected.

Response:

```json
{
  "data": {
    "attemptId": 1,
    "score": 72,
    "totalQuestions": 20,
    "correctAnswers": 14,
    "estimatedBand": 6.0,
    "recommendedBand": 6.0,
    "message": "Start with Band 6.0 Core IELTS Vocabulary."
  }
}
```

Notes:

- If protected, update `users.currentBand`, `users.recommendedBand`, and `users.placementCompletedAt`.
- If public, return result and attach after registration.

---

## 17. Subscription APIs

### GET /api/v1/subscription/plans

Protected or public.

Same response shape as `GET /api/v1/public/pricing`.

### GET /api/v1/me/subscription

Protected.

Response:

```json
{
  "data": {
    "plan": {
      "code": "FREE",
      "name": "Free"
    },
    "status": "ACTIVE",
    "startedAt": "2026-06-05T08:00:00+07:00",
    "trialEndsAt": null,
    "currentPeriodEndsAt": null
  }
}
```

MVP note:

- Pricing can stay static and unauthenticated until billing exists.

---

## 18. AI Example APIs

### POST /api/v1/ai/examples

Protected.

Generates an IELTS example for a learned vocabulary word. Defer unless AI generation is in scope.

Request:

```json
{
  "vocabularyId": 1,
  "topicId": 3,
  "taskType": "WRITING_TASK_2",
  "prompt": "Generate an IELTS Writing Task 2 example sentence."
}
```

Response:

```json
{
  "data": {
    "id": 1,
    "vocabularyId": 1,
    "word": "sustainable",
    "generatedText": "Governments should prioritize sustainable transport policies to reduce urban pollution.",
    "createdAt": "2026-06-05T08:00:00+07:00"
  }
}
```

Errors:

```txt
SUBSCRIPTION_REQUIRED
AI_GENERATION_FAILED
RATE_LIMITED
```

---

## 19. Backend Module Updates

Add or expand modules:

```txt
backend/internal/modules/public/
backend/internal/modules/search/
backend/internal/modules/notification/
backend/internal/modules/challenge/
backend/internal/modules/profile/
backend/internal/modules/placement/
backend/internal/modules/subscription/
backend/internal/modules/ai/
```

Expand existing modules:

```txt
backend/internal/modules/dashboard/
backend/internal/modules/roadmap/
backend/internal/modules/lesson/
backend/internal/modules/flashcard/
backend/internal/modules/quiz/
```

Each module should keep the existing pattern:

```txt
handler.go
service.go
repository.go
dto.go
routes.go
```

---

## 20. Protected Route List

Protected APIs:

```txt
/api/v1/auth/me
/api/v1/dashboard/*
/api/v1/search/*
/api/v1/roadmap
/api/v1/lessons/*
/api/v1/reviews/*
/api/v1/flashcards/*
/api/v1/vocabularies/*
/api/v1/quiz/*
/api/v1/profile/*
/api/v1/notifications/*
/api/v1/challenges/*
/api/v1/me/*
/api/v1/ai/*
```

Public or optional:

```txt
/api/v1/public/*
/api/v1/auth/register
/api/v1/auth/login
/api/v1/placement-test
```

MVP may make placement protected to simplify attaching results to users.

---

## 21. Implementation Priority

1. Expand `GET /dashboard` or `GET /dashboard/summary` to match dashboard screenshots.
2. Expand `GET /roadmap` with summary, lock state, and topic progress.
3. Expand `GET /lessons/:lessonId` with lesson score, XP, band range, and learned word state.
4. Expand flashcard APIs with progress, band labels, XP, and challenge updates.
5. Add quiz session APIs for immediate feedback.
6. Add profile APIs for header, overview, achievements, and calendar.
7. Add notifications and daily challenge APIs.
8. Add placement APIs if the placement CTA is implemented.
9. Add subscription/pricing APIs only if pricing is dynamic.
10. Add AI example API only if AI generation is in scope.

---

## 22. Security Rules

- Never trust `userId` in request body.
- Always use authenticated user ID from JWT context.
- Never expose password hash.
- Never expose quiz or placement correct answers before answer/submit.
- Ensure quiz answer option belongs to the submitted question.
- Ensure question belongs to the active quiz session.
- Ensure lesson, flashcard, quiz, profile, challenge, and notification reads are scoped to the current user.
- Rate limit AI generation and auth endpoints.
