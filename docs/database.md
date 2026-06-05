# Database Design v2 - Screenshot-Aligned Schema

This document extends `docs/database.md` after reviewing all screenshots in `docs/screenshots` and comparing them with `docs/api.md` and `docs/architecture.md`.

Screenshots reviewed:

```txt
landing1.png
landing2.png
landing3.png
landing4.png
landing5.png
landing6.png
dashboard1.png
dashboard2.png
roadmap1.png
roadmap2.png
lesson.png
flashCard1.png
flashCard2.png
flashCard3.png
quiz1.png
quiz2.png
quiz3.png
profile1.png
profile2.png
profile3.png
profile4.png
```

---

## 1. Main Findings

### 1.1 Missing Tables

The v1 schema supports the core learning loop, but the screenshots require these additional tables for a complete product:

```txt
user_xp_events
notifications
daily_challenges
user_daily_challenges
quiz_sessions
quiz_session_answers
placement_tests
placement_questions
placement_options
placement_attempts
placement_attempt_answers
subscription_plans
user_subscriptions
ai_example_requests
```

Optional CMS tables if landing content must be editable from admin:

```txt
landing_testimonials
landing_stats
landing_feature_cards
```

For MVP, landing page content can remain static in frontend, so these CMS tables are optional.

### 1.2 Missing Columns

Important columns missing from v1:

```txt
users.username
users.total_xp
users.level
users.timezone
users.locale
users.starting_band
users.recommended_band
users.placement_completed_at

band_levels.min_score
band_levels.max_score
band_levels.status_label

topics.emoji
topics.color

lessons.band_min
lessons.band_max
lessons.xp_reward
lessons.quiz_time_limit_seconds

vocabularies.target_band
vocabularies.short_definition
vocabularies.example_source

quiz_questions.points
quiz_questions.time_limit_seconds
quiz_questions.topic_id
quiz_options.label

user_vocabulary_progresses.first_learned_at
user_vocabulary_progresses.learned_at
user_vocabulary_progresses.mastered_at

user_lesson_progresses.best_score
user_lesson_progresses.best_xp
user_lesson_progresses.words_learned
user_lesson_progresses.total_words
```

### 1.3 Missing Data Concepts

The screenshots show product concepts not explicitly modeled in v1:

- XP badge in dashboard header and profile level progress.
- Level title such as `Level 18 - Expert`.
- Daily challenge card with goal progress.
- Notifications bell with unread state.
- Immediate-feedback quiz flow with question-by-question answer records.
- Quiz points, timer, progress dots, and per-answer feedback.
- Profile achievement grid with locked and unlocked states.
- Profile calendar heatmap for last 84 days.
- Topic mastery percentages in profile overview.
- Placement test CTA and recommended starting band.
- Pricing tiers and subscriptions.
- AI example generator from landing feature list.

### 1.4 Unnecessary Tables

These v1 tables are not strictly necessary for the screenshots, but can still be kept:

```txt
courses
lesson_vocabularies
```

Reasoning:

- `courses`: the screenshots show one fixed product roadmap. If the app will only ever have one IELTS vocabulary roadmap in MVP, `courses` can be removed and `band_levels` can stand alone. Keep it if future products or multiple courses are expected.
- `lesson_vocabularies`: if every vocabulary item belongs to exactly one lesson, `vocabularies.lesson_id` is enough. Keep the join table if words can appear in multiple lessons or if lesson-specific ordering is needed.

Recommended MVP decision:

```txt
Keep courses.
Keep lesson_vocabularies.
```

They add little complexity and preserve flexibility.

---

## 2. Naming Conventions

Keep v1 conventions:

- Snake case table and column names.
- Plural table names.
- Unsigned integer auto-increment IDs for MVP.
- String constants in Go mapped to `varchar` columns in MySQL.
- `utf8mb4` charset.

---

## 3. Updated Table List

### Core Learning

```txt
users
courses
band_levels
topics
lessons
vocabularies
lesson_vocabularies
quiz_questions
quiz_options
```

### User Progress

```txt
user_vocabulary_progresses
user_lesson_progresses
user_quiz_attempts
daily_activities
```

### Quiz Runtime

```txt
quiz_sessions
quiz_session_answers
```

### Gamification

```txt
achievements
user_achievements
user_xp_events
daily_challenges
user_daily_challenges
notifications
```

### Placement

```txt
placement_tests
placement_questions
placement_options
placement_attempts
placement_attempt_answers
```

### Subscription and Product

```txt
subscription_plans
user_subscriptions
ai_example_requests
```

Optional editable landing content:

```txt
landing_testimonials
landing_stats
landing_feature_cards
```

---

## 4. Enums

```go
type UserRole string
const (
	UserRoleUser  UserRole = "USER"
	UserRoleAdmin UserRole = "ADMIN"
)

type LessonStatus string
const (
	LessonStatusLocked     LessonStatus = "LOCKED"
	LessonStatusUnlocked   LessonStatus = "UNLOCKED"
	LessonStatusInProgress LessonStatus = "IN_PROGRESS"
	LessonStatusCompleted  LessonStatus = "COMPLETED"
)

type VocabularyStatus string
const (
	VocabularyStatusNew      VocabularyStatus = "NEW"
	VocabularyStatusLearning VocabularyStatus = "LEARNING"
	VocabularyStatusReview   VocabularyStatus = "REVIEW"
	VocabularyStatusMastered VocabularyStatus = "MASTERED"
)

type DifficultyLevel string
const (
	DifficultyBeginner     DifficultyLevel = "BEGINNER"
	DifficultyIntermediate DifficultyLevel = "INTERMEDIATE"
	DifficultyAdvanced     DifficultyLevel = "ADVANCED"
)

type QuizQuestionType string
const (
	QuizQuestionMeaningChoice  QuizQuestionType = "MEANING_CHOICE"
	QuizQuestionWordChoice     QuizQuestionType = "WORD_CHOICE"
	QuizQuestionUsageChoice    QuizQuestionType = "USAGE_CHOICE"
	QuizQuestionSentenceChoice QuizQuestionType = "SENTENCE_CHOICE"
)

type FlashcardRating string
const (
	FlashcardRatingAgain FlashcardRating = "AGAIN"
	FlashcardRatingHard  FlashcardRating = "HARD"
	FlashcardRatingGood  FlashcardRating = "GOOD"
	FlashcardRatingEasy  FlashcardRating = "EASY"
)

type QuizSessionStatus string
const (
	QuizSessionInProgress QuizSessionStatus = "IN_PROGRESS"
	QuizSessionCompleted  QuizSessionStatus = "COMPLETED"
	QuizSessionAbandoned  QuizSessionStatus = "ABANDONED"
)

type NotificationType string
const (
	NotificationAchievement NotificationType = "ACHIEVEMENT"
	NotificationReviewDue   NotificationType = "REVIEW_DUE"
	NotificationStreak      NotificationType = "STREAK"
	NotificationSystem      NotificationType = "SYSTEM"
)

type SubscriptionStatus string
const (
	SubscriptionActive   SubscriptionStatus = "ACTIVE"
	SubscriptionTrialing SubscriptionStatus = "TRIALING"
	SubscriptionPastDue  SubscriptionStatus = "PAST_DUE"
	SubscriptionCanceled SubscriptionStatus = "CANCELED"
)
```

---

## 5. Updated Core Tables

### 5.1 users

Adds profile identity, XP, level, placement state, and localization.

```txt
id uint pk
email varchar(255) unique not null
name varchar(255)
username varchar(80) unique null
password_hash varchar(255) not null
avatar_url varchar(500)
role varchar(20) default USER
target_band decimal(3,1) default 7.0
current_band decimal(3,1) null
starting_band decimal(3,1) null
recommended_band decimal(3,1) null
placement_completed_at datetime null
total_xp int default 0
level int default 1
current_streak int default 0
longest_streak int default 0
last_active_at datetime null
timezone varchar(80) default 'UTC'
locale varchar(20) default 'en'
created_at datetime
updated_at datetime
deleted_at datetime null
```

Indexes:

```txt
unique email
unique username
index current_band
index target_band
```

Notes:

- `total_xp` and `level` support the dashboard/profile XP display.
- `username` supports profile handles such as `@alexj`.
- `created_at` supports "Member since Oct 2024".

### 5.2 courses

Keep v1 table.

Recommended added columns:

```txt
band_min decimal(3,1) default 5.0
band_max decimal(3,1) default 8.5
total_words int default 0
total_lessons int default 0
total_topics int default 0
```

These support landing and roadmap summary numbers.

### 5.3 band_levels

```txt
id uint pk
course_id uint index not null
band_score decimal(3,1) not null
min_score decimal(3,1) null
max_score decimal(3,1) null
title varchar(255) not null
description text
status_label varchar(50) null
order_index int default 0
created_at datetime
updated_at datetime
deleted_at datetime null
```

Notes:

- `min_score` and `max_score` support labels like `Band 5.0 -> 8.5`.
- Lock/completion state should be computed per user from lesson/topic progress, not stored here globally.

### 5.4 topics

```txt
id uint pk
band_level_id uint index not null
title varchar(255) not null
slug varchar(255) not null
description text
icon varchar(100)
emoji varchar(20)
color varchar(30)
order_index int default 0
created_at datetime
updated_at datetime
deleted_at datetime null
```

Indexes:

```txt
unique band_level_id + slug
```

### 5.5 lessons

```txt
id uint pk
topic_id uint index not null
title varchar(255) not null
slug varchar(255) not null
description text
required_score int default 80
estimated_minutes int default 10
band_min decimal(3,1) null
band_max decimal(3,1) null
xp_reward int default 0
quiz_time_limit_seconds int null
order_index int default 0
is_published bool default true
created_at datetime
updated_at datetime
deleted_at datetime null
```

Indexes:

```txt
unique topic_id + slug
index topic_id + order_index
```

Notes:

- `band_min`, `band_max`, and `xp_reward` support the lesson sidebar.
- `quiz_time_limit_seconds` supports timed quizzes.

### 5.6 vocabularies

```txt
id uint pk
word varchar(255) index not null
slug varchar(255) unique not null
ipa varchar(255)
audio_url varchar(500)
part_of_speech varchar(100)
meaning_vi text not null
meaning_en text
short_definition text
example_sentence text
example_meaning_vi text
example_source varchar(255)
synonyms_json json
antonyms_json json
collocations_json json
ielts_usage text
difficulty varchar(30) default INTERMEDIATE
target_band decimal(3,1) null
created_at datetime
updated_at datetime
deleted_at datetime null
```

Indexes:

```txt
index word
index difficulty
index target_band
```

Notes:

- `short_definition` supports truncated lesson list rows.
- `target_band` supports badges like `Band 7` or `Band 7.5`.
- `difficulty` should use screenshot labels: `BEGINNER`, `INTERMEDIATE`, `ADVANCED`.

### 5.7 lesson_vocabularies

Keep v1 table.

Recommended added columns:

```txt
is_required bool default true
```

This allows optional bonus words later without changing the lesson count.

### 5.8 quiz_questions

```txt
id uint pk
lesson_id uint index not null
topic_id uint index null
vocabulary_id uint index null
type varchar(50) not null
question text not null
explanation text
points int default 20
time_limit_seconds int null
order_index int default 0
created_at datetime
updated_at datetime
```

Notes:

- `points` supports `+20 points` and quiz score display.
- `time_limit_seconds` supports per-question timers.
- Correct answers remain hidden from frontend until answer/submit time.

### 5.9 quiz_options

```txt
id uint pk
question_id uint index not null
label varchar(5)
content text not null
is_correct bool default false
order_index int default 0
created_at datetime
updated_at datetime
```

Notes:

- `label` supports A/B/C/D display and stable result feedback.
- `is_correct` must never be returned before a user answers or submits.

---

## 6. Updated Progress Tables

### 6.1 user_vocabulary_progresses

```txt
id uint pk
user_id uint uniqueIndex:user_vocab not null
vocabulary_id uint uniqueIndex:user_vocab not null
status varchar(30) default NEW
ease_factor decimal(5,2) default 2.5
interval_days int default 0
review_count int default 0
correct_count int default 0
wrong_count int default 0
last_rating varchar(20) null
first_learned_at datetime null
learned_at datetime null
mastered_at datetime null
last_reviewed_at datetime null
next_review_at datetime index null
created_at datetime
updated_at datetime
```

Notes:

- `first_learned_at` supports "recent activity".
- `learned_at` supports "words learned today" and lesson learned counts.
- `mastered_at` supports roadmap/profile "words mastered".

### 6.2 user_lesson_progresses

```txt
id uint pk
user_id uint uniqueIndex:user_lesson not null
lesson_id uint uniqueIndex:user_lesson not null
status varchar(30) default UNLOCKED
score int null
best_score int null
best_xp int default 0
words_learned int default 0
total_words int default 0
completed_at datetime null
started_at datetime null
last_studied_at datetime null
created_at datetime
updated_at datetime
```

Notes:

- Keep `score` as latest passed score if desired.
- `best_score` and `best_xp` support lesson score and personal best UI.

### 6.3 user_quiz_attempts

Keep v1 table, but add:

```txt
quiz_session_id uint index null
points int default 0
xp_earned int default 0
started_at datetime null
finished_at datetime null
```

Notes:

- Use `user_quiz_attempts` for final historical attempts.
- Use `quiz_sessions` and `quiz_session_answers` for live question-by-question state.

### 6.4 daily_activities

Keep v1 table, but add:

```txt
accuracy_percent decimal(5,2) null
active_minutes int default 0
challenge_progress int default 0
challenge_completed bool default false
```

Notes:

- Supports dashboard weekly chart, profile heatmap, and daily challenge progress.

---

## 7. New Quiz Runtime Tables

### 7.1 quiz_sessions

Stores in-progress quiz state for immediate feedback screens.

```txt
id uint pk
user_id uint index not null
lesson_id uint index not null
status varchar(30) default IN_PROGRESS
current_question_index int default 0
total_questions int default 0
points int default 0
correct_answers int default 0
wrong_answers int default 0
started_at datetime not null
finished_at datetime null
expires_at datetime null
created_at datetime
updated_at datetime
```

Indexes:

```txt
index user_id + status
index lesson_id
```

### 7.2 quiz_session_answers

Stores each answered question.

```txt
id uint pk
quiz_session_id uint uniqueIndex:session_question not null
question_id uint uniqueIndex:session_question not null
selected_option_id uint null
correct_option_id uint null
is_correct bool not null
points_awarded int default 0
answered_at datetime not null
time_spent_seconds int null
created_at datetime
updated_at datetime
```

Notes:

- `correct_option_id` can be stored after answering because the user has already received feedback.
- Never expose `correct_option_id` for unanswered questions.

---

## 8. New Gamification Tables

### 8.1 user_xp_events

Append-only XP ledger.

```txt
id uint pk
user_id uint index not null
source_type varchar(50) not null
source_id uint null
xp int not null
description varchar(255)
created_at datetime
```

Examples for `source_type`:

```txt
FLASHCARD_REVIEW
QUIZ_CORRECT
LESSON_COMPLETED
DAILY_CHALLENGE
ACHIEVEMENT
PLACEMENT_TEST
```

Notes:

- `users.total_xp` is the cached sum.
- This table supports audits and level recalculation.

### 8.2 daily_challenges

```txt
id uint pk
code varchar(100) unique not null
title varchar(255) not null
description text
metric varchar(50) not null
target_value int not null
xp_reward int default 0
is_active bool default true
created_at datetime
updated_at datetime
```

Examples for `metric`:

```txt
WORDS_LEARNED
WORDS_REVIEWED
QUIZZES_TAKEN
LESSONS_COMPLETED
```

### 8.3 user_daily_challenges

```txt
id uint pk
user_id uint uniqueIndex:user_challenge_day not null
daily_challenge_id uint uniqueIndex:user_challenge_day not null
date date uniqueIndex:user_challenge_day not null
progress_value int default 0
target_value int not null
completed_at datetime null
claimed_at datetime null
created_at datetime
updated_at datetime
```

### 8.4 notifications

```txt
id uint pk
user_id uint index not null
type varchar(50) not null
title varchar(255) not null
body text
action_url varchar(500)
read_at datetime null
created_at datetime
```

Indexes:

```txt
index user_id + read_at
index user_id + created_at
```

### 8.5 achievements

Keep v1 table, but add:

```txt
category varchar(50)
requirement_type varchar(50)
requirement_value int
sort_order int default 0
is_active bool default true
```

Examples:

```txt
STREAK_DAYS 7
WORDS_LEARNED 100
FIRST_QUIZ 1
BAND_REACHED 60
TOPIC_COMPLETED 1
```

### 8.6 user_achievements

Keep v1 table, but add:

```txt
progress_value int default 0
is_seen bool default false
```

Notes:

- Locked achievements can be returned by joining all `achievements` with optional `user_achievements`.

---

## 9. New Placement Tables

### 9.1 placement_tests

```txt
id uint pk
title varchar(255) not null
description text
question_count int default 20
is_active bool default true
created_at datetime
updated_at datetime
```

### 9.2 placement_questions

```txt
id uint pk
placement_test_id uint index not null
vocabulary_id uint index null
question text not null
band_score decimal(3,1) null
difficulty varchar(30)
order_index int default 0
created_at datetime
updated_at datetime
```

### 9.3 placement_options

```txt
id uint pk
placement_question_id uint index not null
label varchar(5)
content text not null
is_correct bool default false
order_index int default 0
created_at datetime
updated_at datetime
```

### 9.4 placement_attempts

```txt
id uint pk
user_id uint index null
placement_test_id uint index not null
score int default 0
total_questions int default 0
correct_answers int default 0
estimated_band decimal(3,1) null
recommended_band decimal(3,1) null
started_at datetime not null
completed_at datetime null
created_at datetime
updated_at datetime
```

Notes:

- `user_id` can be null if guests can take the placement test before registration.
- If a guest later registers, attach the attempt to the new user.

### 9.5 placement_attempt_answers

```txt
id uint pk
placement_attempt_id uint uniqueIndex:placement_attempt_question not null
placement_question_id uint uniqueIndex:placement_attempt_question not null
selected_option_id uint null
is_correct bool not null
answered_at datetime not null
created_at datetime
updated_at datetime
```

---

## 10. Subscription and Product Tables

### 10.1 subscription_plans

Supports pricing cards in landing screenshot.

```txt
id uint pk
code varchar(50) unique not null
name varchar(100) not null
description text
price_cents int default 0
currency varchar(10) default 'USD'
billing_interval varchar(20) default 'MONTH'
max_words int null
max_topics int null
max_users int default 1
features_json json
is_popular bool default false
is_active bool default true
sort_order int default 0
created_at datetime
updated_at datetime
```

MVP note:

- If pricing is only a static placeholder, this table can wait.

### 10.2 user_subscriptions

```txt
id uint pk
user_id uint index not null
subscription_plan_id uint index not null
status varchar(30) default ACTIVE
started_at datetime not null
trial_ends_at datetime null
current_period_ends_at datetime null
canceled_at datetime null
created_at datetime
updated_at datetime
```

### 10.3 ai_example_requests

Supports the "AI Example Generator" feature card if implemented in-app.

```txt
id uint pk
user_id uint index not null
vocabulary_id uint index null
topic_id uint index null
prompt text not null
generated_text text
status varchar(30) default PENDING
error_message text
created_at datetime
updated_at datetime
```

MVP note:

- This can be deferred until AI generation is in scope.

---

## 11. Optional Landing CMS Tables

Use only if admins must edit landing content without deployments.

### 11.1 landing_testimonials

```txt
id uint pk
name varchar(255) not null
initials varchar(10)
avatar_url varchar(500)
achieved_band decimal(3,1) null
rating int default 5
quote text not null
is_published bool default true
sort_order int default 0
created_at datetime
updated_at datetime
```

### 11.2 landing_stats

```txt
id uint pk
label varchar(100) not null
value varchar(100) not null
sort_order int default 0
created_at datetime
updated_at datetime
```

### 11.3 landing_feature_cards

```txt
id uint pk
title varchar(255) not null
description text
icon varchar(100)
is_published bool default true
sort_order int default 0
created_at datetime
updated_at datetime
```

---

## 12. Computed Values

Do not store these unless performance becomes a problem:

```txt
dashboard accuracy percentage
weekly chart percentages
topic completion percentage
band completion percentage
words mastered count
review due count
lesson total word count
profile mastery percentage
achievement locked/unlocked display state
calendar heatmap intensity
```

Compute from:

```txt
user_vocabulary_progresses
user_lesson_progresses
user_quiz_attempts
quiz_session_answers
daily_activities
user_achievements
```

Cache only when necessary.

---

## 13. Updated Data Rules

### Progress

- Never trust `user_id` from client requests.
- Use authenticated user ID from JWT context.
- One user can have only one `user_vocabulary_progresses` row per vocabulary.
- One user can have only one `user_lesson_progresses` row per lesson.
- Keep completed lessons completed even if later attempts are worse.
- Keep best score and best XP separately from latest quiz attempt.

### Quiz

- Do not expose `quiz_options.is_correct` before answering.
- Use `quiz_sessions` for live quiz progress.
- Store each answer in `quiz_session_answers`.
- Create or update `user_quiz_attempts` when the session finishes.
- Award XP through `user_xp_events`.

### Flashcards

MVP spaced repetition stays:

```txt
AGAIN -> review today / now
HARD  -> review tomorrow
GOOD  -> review in 3 days
EASY  -> review in 7 days
```

Update:

```txt
status
review_count
correct_count
wrong_count
last_rating
last_reviewed_at
next_review_at
learned_at
mastered_at
daily_activities
user_xp_events
user_daily_challenges
```

### Streak

- Update streak once per local calendar day using `users.timezone`.
- `daily_activities.date` should be based on the user's timezone, not server UTC date.

### XP and Level

- Add XP by inserting `user_xp_events`.
- Update `users.total_xp`.
- Recalculate `users.level`.
- Recommended simple MVP formula:

```txt
level = floor(total_xp / 200) + 1
```

If the product needs custom thresholds, add a `levels` table later.

---

## 14. Migration Notes from v1

1. Keep all v1 tables.
2. Add new columns with safe defaults or nullable values.
3. Add new tables without deleting data.
4. Backfill `users.total_xp` from `daily_activities.xp_earned` if available.
5. Backfill `users.level` from `users.total_xp`.
6. Backfill `user_lesson_progresses.total_words` from `lesson_vocabularies`.
7. Backfill `vocabularies.target_band` from parent lesson/band where possible.
8. Keep `user_quiz_attempts` as historical completed attempts.

---

## 15. MVP Implementation Priority

Implement in this order:

1. Add missing columns to `users`, `lessons`, `vocabularies`, `quiz_questions`, `quiz_options`, and progress tables.
2. Add `user_xp_events`, `daily_challenges`, `user_daily_challenges`, and `notifications`.
3. Add `quiz_sessions` and `quiz_session_answers`.
4. Add placement test tables if the placement CTA is implemented.
5. Add subscription tables only when pricing is not static.
6. Add AI example request table only when AI generation is implemented.
7. Add landing CMS tables only when admin editing is required.
