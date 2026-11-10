# Foundation English Spec (CEFR A1 → B2)

Status: draft for Milestones 13-21 (Tasks 49-77). Written for the A1 level; A2-B2 reuse the same model with new content.

This spec is the source of truth for Part 2 of `docs/tasks.md`. Table and endpoint summaries are mirrored in `docs/database-v2.md` §14 and `docs/api-v2.md` §23. Coding rules from `AGENTS.md` still apply (module pattern, React Query, axios only, validation contracts on both sides, Swagger after endpoint changes).

Content is for private, personal study only.

---

## 1. Goals

The learner is a Vietnamese adult who has lost their English foundation. After A1 they can:

- introduce themselves and talk about work, family and daily routine in simple sentences;
- handle 36 everyday situations (office, shopping, travel, health, phone);
- read and pronounce any word from its IPA, including final sounds;
- understand slow, clear speech and short texts (messages, notices, menus, schedules);
- write a 5-7 sentence paragraph with no tense mistakes.

Non-goals for A1: placement test UI, multi-user content authoring, mobile app.

---

## 2. Course Structure

### 2.1 Courses

The platform gets a second course next to the IELTS course.

| Course | `slug` | `kind` | Band levels |
|---|---|---|---|
| IELTS Vocabulary Roadmap | `ielts-vocabulary-roadmap` | `IELTS` | 5.0, 6.0, 7.0, 8.0 |
| Foundation English | `foundation-english` | `CEFR` | A1, A2, B1, B2 |

- `courses.kind` (`IELTS` | `CEFR`) and `band_levels.cefr_level` (`A1`… `C2`, null for IELTS) are added.
- CEFR band levels keep a numeric `band_score` for ordering only: A1 = 1.0, A2 = 2.0, B1 = 3.0, B2 = 4.0. The UI shows `cefr_level`, never the number.
- The roadmap page gets a course switcher. The selected course is remembered per user (`users.active_course_id`, nullable → first published course).
- Migration (Task 49): the temporary "A1 Foundation" band level created by Task 44 moves to the new course by updating `band_levels.course_id`. IDs do not change, so existing progress, flashcards and quiz history are kept.

### 2.2 Tracks inside a level

A CEFR level has two kinds of topics, stored in `topics.track`:

| Track | Purpose | Unlocking |
|---|---|---|
| `CORE` | Week 1…4: one lesson per everyday situation, each lesson = a sequence of skill activities | Sequential inside the level (finish lesson n to unlock n+1); a week's checkpoint must be passed to unlock the next week |
| `VOCABULARY` | Theme word banks (the 14 topics / 92 lessons from Task 44) for flashcards and quick quizzes | Always unlocked, no effect on CORE unlocking |

IELTS topics are `CORE` (default), so the IELTS roadmap behaves exactly as today.

Roadmap rule change: the "previous lesson completed" chain is computed **per course and per track**, not across the whole database.

### 2.3 A1 syllabus

4 weeks × 9 situation lessons + 1 checkpoint per week = 40 CORE lessons.

Every situation lesson contains, in this order (activities marked * are optional):

1. `PRONUNCIATION` — 1-2 sounds or one pronunciation rule
2. `VOCABULARY` — 8-12 words for the situation (linked to existing vocabulary rows)
3. `GRAMMAR` — the week's grammar point (theory card on first use, then exercises)
4. `LISTENING` — model dialogue for the situation
5. `READING`* — a short real-life text for the situation
6. `WRITING`* — 3-7 sentences using the lesson's patterns
7. `SPEAKING` — AI role-play of the situation

| Week | Grammar points | Situations (lessons) | Sounds |
|---|---|---|---|
| 1 | G1 *to be* + subject pronouns · G2 nouns, plurals, a/an/the, possessives | Greetings · Introducing yourself · Introducing others · Personal information · Family · Hobbies · Describing your job · Finding places in the building · Borrowing things | /iː/–/ɪ/ · /e/–/æ/ · /ʌ/–/ɑː/ · /uː/–/ʊ/ · /ɜːr/–/ə/ · plural -s /s/ /z/ /ɪz/ |
| 2 | G3 present simple · G4 present continuous · G5 *there is/are* + prepositions | Daily work schedule · Asking for leave · Scheduling a meeting · Booking a table · Ordering food · Paying the bill · Asking prices · Choosing a size · Returning items | /eɪ/ /aɪ/ /ɔɪ/ · /aʊ/ /oʊ/ · /p/–/b/ · /t/–/d/ · /k/–/ɡ/ · word stress (2 syllables) |
| 3 | G6 adjectives, comparatives, superlatives · G7 *can / must / should* | Asking for directions · Taking a taxi · Buying tickets · Airport check-in · Hotel check-in/out · Train and bus times · Describing symptoms · Buying medicine · At the bank | /θ/–/ð/ · /s/–/z/ · /ʃ/–/ʒ/ · /tʃ/–/dʒ/ · /f/–/v/ · /l/–/r/ · final consonant clusters |
| 4 | G8 past simple (regular + 20 irregular) · G9 future: *will / be going to* | Phone messages · Invitations · Arranging to meet · Changing plans · Weather small talk · Your weekend · Sending a parcel · Reporting a problem · Emergency help | /m/ /n/ /ŋ/ · /h/ /w/ /j/ · past -ed /t/ /d/ /ɪd/ · sentence stress · linking sounds |

Checkpoint lesson (end of each week): one `CHECKPOINT` activity mixing exercises from the week (grammar, listening, reading) and one speaking role-play without hints. Pass mark 70.

---

## 3. Data Model

All new tables follow `docs/database-v2.md` §2 conventions (snake_case, plural, uint IDs, `utf8mb4`, string enums as varchar). GORM models live in `internal/models/`.

### 3.1 Changes to existing tables

```txt
courses            + kind varchar(20) default 'IELTS'
band_levels        + cefr_level varchar(5) null
topics             + track varchar(20) default 'CORE'
users              + active_course_id uint null
quiz_questions     (unchanged — vocabulary quizzes keep using it)
ai_conversation_scenarios
                   ~ user_id becomes nullable (null = system scenario)
                   + is_system bool default false
                   + cefr_level varchar(5) null
                   + role text, task text            (role the AI plays, learner's task)
                   + target_patterns_json json       (["I work as a ...", "I'm from ..."])
                   + max_turns int default 8
ai_conversation_messages
                   + session_id uint null index
```

### 3.2 Activities and exercises

`lesson_activities` — the steps of a CORE lesson.

```txt
id uint pk
lesson_id uint index not null
slug varchar(150) not null              unique (lesson_id, slug)
type varchar(30) not null               VOCABULARY | PRONUNCIATION | GRAMMAR | LISTENING | READING | WRITING | SPEAKING | CHECKPOINT
title varchar(255) not null
instructions_vi text                    shown above the activity
grammar_point_id uint null              GRAMMAR
passage_id uint null                    LISTENING / READING
scenario_id uint null                   SPEAKING / CHECKPOINT
content_json json null                  type-specific settings, see §4
pass_score int default 70
xp_reward int default 10
is_required bool default true
order_index int default 0
status varchar(20) default 'PUBLISHED'  DRAFT | PUBLISHED
created_at, updated_at, deleted_at
```

`exercises` — scored items inside an activity. Kept separate from `quiz_questions` because answers are not always options (typed text, reordering, audio) and grading needs private answer data.

```txt
id uint pk
activity_id uint index not null
type varchar(40) not null               see §5
prompt text not null                    question shown to the learner
payload_json json not null              PUBLIC data needed to render the item
answer_json json not null               PRIVATE: accepted answers, correct option, grading flags — never sent before answering
explanation_vi text                     shown after answering
skill_tag varchar(80) null              e.g. grammar:present_simple, sound:th, vocab:a1-office
points int default 10
order_index int default 0
created_at, updated_at
```

`user_activity_progresses`

```txt
id uint pk
user_id uint                            unique (user_id, activity_id)
activity_id uint
status varchar(20) default 'IN_PROGRESS' IN_PROGRESS | PASSED | FAILED
score int null                          last score 0-100
best_score int null
attempts int default 0
passed_at datetime null
last_attempt_at datetime null
created_at, updated_at
```

`user_exercise_answers` — one row per answered item (for mistakes and review).

```txt
id uint pk
user_id uint index
exercise_id uint index
activity_attempt int                    attempt number from user_activity_progresses.attempts
answer_json json                        what the learner submitted
is_correct bool
score_ratio decimal(4,3)                1.000 = fully correct; partial credit for dictation / AI grading
answered_at datetime
```

### 3.3 Skill content

`grammar_points`

```txt
id uint pk
slug varchar(100) unique                g1-to-be, g3-present-simple
cefr_level varchar(5)
title varchar(255)                      "Present simple"
title_vi varchar(255)                   "Thì hiện tại đơn"
summary_vi text                         when to use it, in Vietnamese
formula_md text                         "S + V(s/es) + O"; markdown table allowed
examples_json json                      [{ "en": "...", "vi": "...", "audioUrl": "..." }]
common_mistakes_json json               [{ "wrong": "He work here.", "right": "He works here.", "noteVi": "..." }]
order_index int
created_at, updated_at
```

`user_grammar_progresses` — (user_id, grammar_point_id) unique; `mastery` 0-100, `correct_count`, `wrong_count`, `last_practiced_at`. Mastery = rolling score of the last 20 answers tagged `grammar:<slug>` (exercises, writing and speaking mistakes).

`phonemes`

```txt
id uint pk
symbol varchar(10) unique               "iː", "θ", "eɪ"
kind varchar(20)                        VOWEL | DIPHTHONG | CONSONANT
voiced bool null                        consonants only
example_words_json json                 [{ "word": "see", "ipa": "/siː/", "audioUrl": "..." }]
minimal_pairs_json json                 [{ "a": "ship", "b": "sheep" }]
mouth_vi text                           lips, tongue, jaw — in Vietnamese
mistakes_vi text                        typical Vietnamese-learner mistakes
audio_url varchar(500)
image_url varchar(500) null
order_index int
```

Phoneme set: American English as transcribed in our vocabulary data (Oxford NAmE style): 11 vowels (iː ɪ e æ ɑː ɔː ʊ uː ʌ ɜːr ə), 5 diphthongs (eɪ aɪ ɔɪ aʊ oʊ), 24 consonants. Total 40 cards; r-coloured vowels (ɪr, er, ʊr, ɑːr, ɔːr) are taught as combinations, not separate cards.

`user_phoneme_progresses` — (user_id, phoneme_id) unique; `accuracy` 0-100 (rolling average from pronunciation scoring), `attempts`, `last_practiced_at`.

`passages` — text and/or audio used by LISTENING and READING.

```txt
id uint pk
slug varchar(150) unique
cefr_level varchar(5)
genre varchar(30)                       DIALOGUE | MESSAGE | EMAIL | NOTICE | MENU | SCHEDULE | ARTICLE | PROFILE | MONOLOGUE
title varchar(255)
body_text text                          full transcript / reading text
audio_url varchar(500) null
segments_json json null                 [{ "index": 0, "speaker": "A", "start": 0.0, "end": 2.4, "text": "Hi, I'm Lan." }]
glossary_json json null                 [{ "word": "colleague", "vocabularyId": 12 }] — tap-to-translate overrides
source varchar(255)                     book / unit / page, or "generated"
word_count int
created_at, updated_at
```

`segments_json` is what makes sentence replay, dictation and shadowing possible. It comes from Whisper (Task 64) for real recordings, or from the TTS step for generated dialogues.

### 3.4 Writing, speaking, pronunciation, mistakes

`writing_submissions`

```txt
id uint pk
user_id uint index
activity_id uint index
text text
status varchar(20)                      GRADED | REWRITE_REQUIRED | COMPLETED
score int                               0-100
report_json json                        AI grading result, see §6.3
rewrite_text text null
completed_at datetime null
created_at, updated_at
```

`ai_conversation_sessions`

```txt
id uint pk
user_id uint index
scenario_id uint index
activity_id uint null                   set when started from a lesson
status varchar(20)                      IN_PROGRESS | FINISHED | ABANDONED
turn_count int default 0
hint_count int default 0
score int null
report_json json null                   see §6.2
started_at datetime, finished_at datetime null
created_at, updated_at
```

`pronunciation_attempts`

```txt
id uint pk
user_id uint index
reference_text varchar(500)
source varchar(30)                      WORD | SHADOWING | READ_ALOUD | EXERCISE
exercise_id uint null
accuracy int, fluency int, completeness int, pronunciation int   0-100 from Azure
words_json json                         per-word and per-phoneme scores (IPA)
created_at
```

Audio is scored and discarded; it is not stored.

`learner_mistakes` — one row per distinct mistake, counted when it repeats.

```txt
id uint pk
user_id uint index
source varchar(20)                      EXERCISE | WRITING | SPEAKING | PRONUNCIATION
skill_tag varchar(80) index             grammar:past_simple, sound:th, vocab:word-choice, spelling
wrong_text varchar(500)
correct_text varchar(500)
explanation_vi text
occurrences int default 1
status varchar(20) default 'OPEN'       OPEN | REVIEWING | FIXED
vocabulary_id uint null                 flashcard created for this mistake (Task 73)
last_seen_at datetime
created_at, updated_at
unique (user_id, skill_tag, correct_text)
```

---

## 4. Activity Types

`content_json` per type. Unknown keys are ignored; missing optional keys use the defaults shown.

| Type | `content_json` | Completed when |
|---|---|---|
| `VOCABULARY` | `{ "vocabularySlugs": ["a1-hello", …], "quizSize": 5 }` | learner flips every card once and scores ≥ pass_score on the mini quiz (built from these words, MEANING_CHOICE) |
| `PRONUNCIATION` | `{ "phonemes": ["iː","ɪ"], "rule": null }` or `{ "phonemes": [], "rule": "ending-s" }` | exercises ≥ pass_score |
| `GRAMMAR` | `{ "showTheory": true }` | exercises ≥ pass_score; theory card shown first when `showTheory` |
| `LISTENING` | `{ "speeds": [0.75, 1.0], "transcript": "AFTER_ANSWER" }` (`HIDDEN` / `BLURRED` / `AFTER_ANSWER`) | exercises ≥ pass_score |
| `READING` | `{ "tapToTranslate": true, "readAloud": false }` | exercises ≥ pass_score |
| `WRITING` | see §4.1 | AI score ≥ pass_score **and** rewrite submitted |
| `SPEAKING` | `{ "hintsAllowed": true }` | session finished with `taskCompleted = true` and score ≥ pass_score |
| `CHECKPOINT` | `{ "fromActivityTypes": ["GRAMMAR","LISTENING","READING"], "itemCount": 15, "speaking": true }` | exercises ≥ pass_score and speaking passed without hints |

Activity score = round(100 × Σ(points × score_ratio) / Σ points). Retries are unlimited; `best_score` is kept.

Lesson completion (CORE lessons that have activities): all `is_required` activities are `PASSED`. This sets `user_lesson_progresses.status = COMPLETED` and `score = average of required activity best scores`, and awards the lesson XP once (`LESSON_COMPLETED`). Lessons without activities (IELTS and VOCABULARY track) keep the existing quiz-pass rule.

### 4.1 Writing prompt

```json
{
  "promptVi": "Giới thiệu bản thân: tên, công việc, nơi sống và một sở thích.",
  "promptEn": "Introduce yourself.",
  "frame": ["My name is ...", "I am a ...", "I live in ...", "I like ..."],
  "minSentences": 4,
  "maxSentences": 7,
  "targetGrammar": ["g1-to-be", "g3-present-simple"],
  "targetVocabulary": ["a1-name", "a1-job", "a1-live"],
  "modelAnswer": "My name is Huy. I am a developer. I live in Hanoi. I like football."
}
```

---

## 5. Exercise Types

`payload_json` is public; `answer_json` is private (never in GET responses). All typed answers are normalised before comparison (§5.2).

| Type | Skill | payload_json | answer_json | Grading |
|---|---|---|---|---|
| `MULTIPLE_CHOICE` | any | `{ "options": [{ "id": "a", "text": "…" }], "audioUrl"?, "imageUrl"? }` | `{ "correct": "a" }` | exact |
| `TRUE_FALSE` | reading, listening | `{ "statement": "…" }` | `{ "value": true }` | exact |
| `FILL_BLANK` | grammar | `{ "sentence": "She ___ (work) in a bank.", "hint"? }` | `{ "accepted": ["works"] }` | normalised match |
| `REORDER` | grammar | `{ "tokens": ["name","is","My","Huy"] }` (shuffled) | `{ "accepted": ["My name is Huy."] }` | normalised match |
| `ERROR_CORRECTION` | grammar | `{ "tokens": ["He","work","here","."] }` | `{ "wrongIndex": 1, "accepted": ["works"] }` | both parts right = 1, index only = 0.5 |
| `TRANSFORM` | grammar | `{ "sentence": "I like coffee.", "instruction": "NEGATIVE" }` (`NEGATIVE` / `QUESTION` / `PAST` / `THIRD_PERSON`) | `{ "accepted": ["I don't like coffee.", "I do not like coffee."], "aiFallback": true }` | match, else AI (§6.4) |
| `TRANSLATE_VI_EN` | grammar, vocab | `{ "sourceVi": "Tôi làm việc ở Hà Nội." }` | `{ "reference": ["I work in Hanoi."], "aiGraded": true }` | match, else AI |
| `DICTATION` | listening | `{ "audioUrl", "passageId"?, "segmentIndex"?, "wordCount" }` | `{ "accepted": ["I usually take the bus."] }` | word-level diff → score_ratio |
| `LISTEN_GAP_FILL` | listening | `{ "audioUrl", "text": "I start work at {{1}} o'clock." }` | `{ "blanks": { "1": ["nine","9"] } }` | per blank |
| `LISTEN_PICK_IMAGE` | listening | `{ "audioUrl", "options": [{ "id", "imageUrl", "label" }] }` | `{ "correct": "b" }` | exact |
| `LISTEN_NUMBERS` | listening | `{ "audioUrl", "format": "TIME" }` (`NUMBER` / `TIME` / `PRICE` / `PHONE` / `DATE`) | `{ "accepted": ["7:30"] }` | format-aware normalisation |
| `MINIMAL_PAIR` | pronunciation | `{ "audioUrl", "options": [{ "id": "a", "word": "ship", "ipa": "/ʃɪp/" }, …] }` | `{ "correct": "b" }` | exact |
| `IPA_TO_WORD` | pronunciation | `{ "ipa": "/wɜːrk/", "options": [{ "id", "word" }] }` | `{ "correct": "a" }` | exact |
| `WORD_TO_IPA` | pronunciation | `{ "word": "work", "audioUrl", "options": [{ "id", "ipa" }] }` | `{ "correct": "c" }` | exact |
| `ENDING_SOUND` | pronunciation | `{ "word": "watched", "audioUrl", "options": ["/t/","/d/","/ɪd/"] }` | `{ "correct": "/t/" }` | exact |
| `WORD_STRESS` | pronunciation | `{ "word": "computer", "syllables": ["com","pu","ter"], "audioUrl" }` | `{ "stressIndex": 1 }` | exact |
| `SPEAK` | pronunciation, shadowing | `{ "text": "I'm from Vietnam.", "audioUrl", "ipa"? }` | `{ "minAccuracy": 70 }` | Azure accuracy ≥ minAccuracy (Task 71); before Task 71: self-check, always 1 |

### 5.1 Answer submission

```json
{ "exerciseId": 31, "answer": { "text": "works" } }
{ "exerciseId": 32, "answer": { "optionId": "b" } }
{ "exerciseId": 33, "answer": { "tokens": ["My","name","is","Huy"] } }
{ "exerciseId": 34, "answer": { "index": 1, "text": "works" } }
{ "exerciseId": 35, "answer": { "blanks": { "1": "nine" } } }
{ "exerciseId": 36, "answer": { "stressIndex": 1 } }
```

### 5.2 Normalisation

Applied to typed answers and accepted answers before comparing:

1. trim, collapse whitespace, lower-case;
2. curly quotes/apostrophes → straight; strip final `. ! ?`;
3. expand contractions both ways (`I'm` ≡ `I am`, `don't` ≡ `do not`, `can't` ≡ `cannot`, `it's` ≡ `it is`, …);
4. numbers: `nine` ≡ `9` for `LISTEN_NUMBERS` and gap-fills; times `7.30` ≡ `7:30` ≡ `seven thirty`.

Dictation score_ratio = 1 − (word edit distance / reference word count), floored at 0; ≥ 0.9 counts as correct.

---

## 6. AI Usage

All AI calls go through one backend client (`internal/shared/ai`, extracted from `aiconversation/openai_client.go`) with:

- model from config (`OPENAI_MODEL`), timeout 30 s, one retry on 5xx;
- `response_format: { type: "json_schema", strict: true }` for every call that returns data;
- Vietnamese explanations, English corrections;
- the API key never leaves the backend.

### 6.1 Speaking role-play system prompt (per CEFR level)

```txt
You are {role}. Situation: {situation}.
The learner is a Vietnamese adult at CEFR {level}. Their task: {task}.
Useful patterns for this lesson: {targetPatterns}.

Rules for {level}:
A1: max 10 words per sentence, only A1 words, one question per turn, speak slowly and clearly.
A2: max 15 words per sentence, A1-A2 words, one or two questions per turn.
B1: natural sentences, everyday idioms allowed, ask follow-up "why/how" questions.
B2: natural speed and vocabulary, challenge the learner's opinions politely.

Always:
- Stay in role. Never switch to Vietnamese unless asked for a hint.
- If the learner makes a mistake, do not stop the conversation; repeat their idea correctly in your reply (recast).
- If the learner is silent, off-topic or writes Vietnamese, give one short English sentence they can copy.
- End naturally after {maxTurns} learner turns or when the task is complete, with a closing line.
```

### 6.2 Speaking feedback report (session finish)

```json
{
  "taskCompleted": true,
  "scores": { "grammar": 6, "vocabulary": 7, "fluency": 5, "taskAchievement": 8 },
  "overall": 65,
  "mistakes": [
    { "said": "I go to office yesterday", "correct": "I went to the office yesterday",
      "skillTag": "grammar:g8-past-simple", "explanationVi": "Hành động đã xảy ra trong quá khứ nên dùng 'went'." }
  ],
  "goodPhrases": ["Could I have the bill, please?"],
  "newWords": ["receipt"],
  "tipVi": "Bạn hay quên chia động từ ở quá khứ."
}
```

`overall` = 10 × average of the four scores. Every mistake is upserted into `learner_mistakes`.

### 6.3 Writing grading

```json
{
  "scores": { "taskAchievement": 8, "grammar": 6, "vocabulary": 7, "coherence": 7 },
  "overall": 70,
  "correctedText": "My name is Huy. I am a developer. I live in Hanoi.",
  "mistakes": [
    { "start": 18, "end": 20, "wrong": "is", "correct": "am",
      "skillTag": "grammar:g1-to-be", "explanationVi": "Chủ ngữ 'I' đi với 'am'." }
  ],
  "usedTargets": { "grammar": ["g1-to-be"], "vocabulary": ["a1-live"] },
  "feedbackVi": "Bài viết rõ ràng. Chú ý động từ to be.",
  "rewriteRequired": true
}
```

`start`/`end` are character offsets into the learner's text, used to draw the diff. `rewriteRequired` is true whenever there is at least one mistake.

### 6.4 Short-answer grading (TRANSFORM, TRANSLATE_VI_EN)

Only called when normalised matching fails.

```json
{ "isCorrect": false, "scoreRatio": 0.5, "suggested": "I don't like coffee.", "explanationVi": "Thiếu trợ động từ 'do'." }
```

### 6.5 Hint (speaking)

Input: last 6 turns + task. Output: `{ "hintVi": "Hãy nói tên và công việc của bạn.", "sayThis": "My name is Huy and I'm a developer." }`. Each hint increments `hint_count`; checkpoints disallow hints.

### 6.6 Content generation (admin only, Task 74 and content pipeline)

Input: activity type, CEFR level, grammar point / situation / phonemes, allowed vocabulary list, count. Output: a JSON array of exercises in the §5 format, always saved as `DRAFT`. A second call validates: one correct answer per item, no words above the level, sentences ≤ 12 words at A1.

### 6.7 Cost guard

- Speaking: send at most the last 20 messages plus the system prompt.
- Per-user daily limits from config: 200 AI chat turns, 30 writing gradings, 300 short-answer gradings. Exceeding returns `AI_LIMIT_REACHED`.

---

## 7. Pronunciation Scoring

- Provider: Azure Speech Pronunciation Assessment, `en-US`, `PhonemeAlphabet = IPA`, granularity `Phoneme`, scripted mode with the reference text.
- The browser records with `MediaRecorder` (WebM/Opus or WAV 16 kHz mono); the backend converts if needed and calls Azure. The key stays on the backend (`AZURE_SPEECH_KEY`, `AZURE_SPEECH_REGION`).
- Response stored in `pronunciation_attempts.words_json`; every phoneme score updates `user_phoneme_progresses` (rolling average of the last 30 occurrences).
- UI colours: ≥ 80 green, 60-79 amber, < 60 red, per word and per phoneme in the IPA.
- Phonemes with accuracy < 60 after ≥ 5 attempts become `learner_mistakes` (`sound:<symbol>`).

---

## 8. Content Pipeline

### 8.1 Source files

Content lives in the repo as JSON, one folder per level:

```txt
content/
  a1-vocabulary/            (Task 44 / D1-D4, Excel import)
  a1/
    grammar-points.json
    phonemes.json
    passages/*.json
    scenarios.json
    lessons/week-1.json … week-4.json
    audio/                   (generated, git-ignored, copied to frontend/public/audio/a1/)
```

### 8.2 Bundle format

`POST /admin/content/import/preview` and `POST /admin/content/import` accept one JSON bundle (same preview/commit pattern as the vocabulary import). Everything is upserted by slug.

```json
{
  "version": 1,
  "courseSlug": "foundation-english",
  "cefrLevel": "A1",
  "grammarPoints": [ { "slug": "g1-to-be", "titleVi": "Động từ to be", "...": "..." } ],
  "phonemes": [ { "symbol": "iː", "kind": "VOWEL", "...": "..." } ],
  "passages": [ { "slug": "a1-w1-greetings-dialogue", "genre": "DIALOGUE", "...": "..." } ],
  "scenarios": [ { "slug": "a1-first-day", "role": "a friendly new colleague", "task": "...", "...": "..." } ],
  "topics": [
    {
      "slug": "a1-week-1", "title": "Week 1 - Meeting People", "track": "CORE",
      "lessons": [
        {
          "slug": "a1-w1-l1-greetings", "title": "Greetings",
          "activities": [
            { "type": "PRONUNCIATION", "slug": "sounds", "title": "/iː/ and /ɪ/",
              "content": { "phonemes": ["iː","ɪ"] },
              "exercises": [ { "type": "MINIMAL_PAIR", "prompt": "Bạn nghe thấy từ nào?",
                               "payload": { "...": "..." }, "answer": { "correct": "b" },
                               "explanationVi": "...", "skillTag": "sound:iː" } ] },
            { "type": "GRAMMAR", "slug": "grammar", "grammarPointSlug": "g1-to-be", "exercises": [] },
            { "type": "SPEAKING", "slug": "speaking", "scenarioSlug": "a1-first-day" }
          ]
        }
      ]
    }
  ]
}
```

Validation (preview reports per item): unknown references (grammarPointSlug, passageSlug, scenarioSlug, vocabularySlugs), duplicate slugs, exercise payload/answer shape per type, MULTIPLE_CHOICE correct id present in options, at least one accepted answer.

### 8.3 Tools

| Tool | Task | What it does |
|---|---|---|
| `cmd/contentgen` | 74 + content | AI-generates draft exercises (§6.6) into a bundle file for review |
| `cmd/tts` | D7 | Generates dialogue/sentence audio (OpenAI TTS or macOS `say`) and fills `segments_json` timings |
| Audio import (admin page) | 64 | Upload mp3 → Whisper (`verbose_json`, segment timestamps) → passage draft |
| PDF / image import (admin page) | 65 | Upload pages → vision model extracts exercises into a bundle draft |

### 8.4 Sources (personal study)

Grammar: *Essential Grammar in Use* (units for G1-G9). Listening/reading: *English File Elementary* audio + transcripts, BBC Learning English, VOA Learning English. Pronunciation: *Ship or Sheep?*, *English Pronunciation in Use Elementary*. Sentences: Tatoeba. Everything imported is stored with `source` so it can be traced.

---

## 9. API Summary

All endpoints are under `/api/v1`, protected unless noted, using the standard response envelope. Admin endpoints use `middleware.RequireAdmin()`.

| Method | Path | Task | Purpose |
|---|---|---|---|
| GET | `/courses` | 49 | Published courses with progress summary |
| PATCH | `/profile/active-course` | 49 | `{ "courseId": 2 }` |
| GET | `/roadmap?courseSlug=` | 49 | Existing roadmap, scoped to a course, with `cefrLevel` and `track` |
| GET | `/lessons/:lessonId/activities` | 50 | Activities + user progress (no exercise answers) |
| GET | `/activities/:activityId` | 51 | Activity detail, exercises with `payload` only, linked grammar point / passage / scenario |
| POST | `/activities/:activityId/answers` | 51 | Check one exercise (§5.1) → `{ isCorrect, scoreRatio, correctAnswer, explanationVi }` |
| POST | `/activities/:activityId/complete` | 52 | Close the attempt → score, status, XP, lesson completion |
| GET | `/grammar-points/:slug` | 54 | Theory card + user mastery |
| GET | `/grammar-points` | 56 | List with mastery, weakest first option |
| GET | `/phonemes` · `/phonemes/:symbol` | 57 | IPA chart and sound cards with user accuracy |
| GET | `/passages/:passageId` | 60 | Text, audio, segments, glossary |
| POST | `/writing/submissions` | 66 | `{ activityId, text }` → AI report |
| POST | `/writing/submissions/:id/rewrite` | 66 | `{ text }` → completes when no new mistakes |
| POST | `/ai-conversations/sessions` | 69 | `{ scenarioSlug, activityId? }` start a session |
| POST | `/ai-conversations/sessions/:id/messages` | 69 | Chat turn (replaces scenario-level thread for sessions) |
| POST | `/ai-conversations/sessions/:id/hint` | 70 | Vietnamese hint |
| POST | `/ai-conversations/sessions/:id/finish` | 69 | Feedback report, XP, mistakes |
| POST | `/pronunciation/assess` | 71 | multipart `audio` + `referenceText` (+ `exerciseId`) → scores |
| GET | `/mistakes` · PATCH `/mistakes/:id` | 73 | Mistake log, mark fixed |
| POST | `/practice/generate` | 74 | AI extra practice for a skill tag (cached per user/day) |
| GET | `/foundation/summary?level=A1` | 76 | Per-skill progress and A1 checklist |
| POST | `/admin/content/import/preview` · `/admin/content/import` | 53 | Content bundle import |
| POST | `/admin/passages/transcribe` | 64 | Whisper import |
| POST | `/admin/content/extract` | 65 | PDF/image extraction |

Error codes added: `ACTIVITY_LOCKED`, `EXERCISE_NOT_IN_ACTIVITY`, `INVALID_ANSWER_SHAPE`, `AI_UNAVAILABLE`, `AI_LIMIT_REACHED`, `SESSION_FINISHED`, `AUDIO_INVALID`.

---

## 10. Frontend

| Route | Page | Task |
|---|---|---|
| `/roadmap` | course switcher, CEFR band labels, CORE weeks + VOCABULARY library | 49 |
| `/lessons/:lessonId/learn` | **Lesson Player**: stepper of activities, resumes at first unfinished one | 52 |
| `/grammar` · `/grammar/:slug` | grammar list with mastery, theory card | 54, 56 |
| `/pronunciation` · `/pronunciation/:symbol` | IPA chart, sound card | 57 |
| `/mistakes` | mistake log and review | 73 |
| `/foundation` | A1 dashboard | 76 |

Components (in `src/features/`): `activity/` (player shell, progress, result), `exercise/` (one renderer per exercise type, registered in a map keyed by type), `audio/` (player with speed, sentence replay, blurred transcript; reuses `lib/pronunciation.ts`), `writing/` (editor, diff view), `speaking/` (existing session page upgraded), `pronunciation/` (recorder, coloured IPA).

UI language: instructions and feedback in Vietnamese, content in English.

---

## 11. Open Questions

- Azure free tier (5 audio hours/month) is enough for one learner; confirm before Task 71.
- Whether CHECKPOINT failures should re-lock the following week or only block it (current rule: only block).
- A2-B2 syllabi are not yet defined; they reuse §2-§9 unchanged.
