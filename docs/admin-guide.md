# Admin Guide

Operational guide for admin-only features. Unlike the other files in `docs/`,
this describes what is actually implemented today, not a pre-implementation
spec — keep it in sync with the code when admin features change.

## 1. Promoting a user to ADMIN

There is no UI or API to grant the `ADMIN` role — `POST /api/v1/auth/register`
always creates users with `role = USER` (`backend/internal/modules/auth/service.go`).
To promote an existing user, use the `-promote-admin` flag on the seed binary:

```bash
cd backend
docker compose up -d   # from repo root, if MySQL is not already running
go run ./cmd/seed -promote-admin=someone@example.com
```

This updates `users.role` to `ADMIN` for that email and exits without running
the normal seed. It fails with `no user found with email "..."` if the email
doesn't exist yet — register the account first.

Implementation: `backend/cmd/seed/main.go` (`promoteAdmin` function).

**Important:** `role` is embedded in the JWT at login time
(`backend/internal/shared/jwt/jwt.go`), so a user must **log out and log back
in** after being promoted — an already-issued token keeps the old role until
it expires (`JWT_ACCESS_TTL_MINUTES`, default 1440 minutes).

Admin-only routes are guarded by `middleware.RequireAdmin()`, chained after
`middleware.Auth()`. To add a new admin-only endpoint, register its route
group the same way:

```go
adminGroup := router.Group("/admin/<resource>", middleware.Auth(jwtManager), middleware.RequireAdmin())
```

## 2. Importing vocabulary from Excel

Admins can bulk-create/update vocabulary and assign each word to a specific
lesson (and therefore topic) via an `.xlsx` file.

### Where

- UI: `/admin/vocabularies/import` (linked from the sidebar for admins only,
  `frontend/src/components/layout/Sidebar.tsx`). Page:
  `frontend/src/pages/AdminVocabularyImportPage.tsx`.
- API: `backend/internal/modules/vocabulary/` — `import.go` (parsing +
  row validation), `service.go` (`Import`), `repository.go`
  (`CommitImportRows` and the upsert helpers).

### Endpoints (admin only, `multipart/form-data`, field name `file`, max 5MB)

| Method | Path | Effect |
| --- | --- | --- |
| POST | `/api/v1/admin/vocabularies/import/preview` | Parses and validates the file. **Does not write to the database.** Use this to show the admin what will happen. |
| POST | `/api/v1/admin/vocabularies/import` | Same validation, then commits every valid row in one transaction. |

Both return `ImportResultResponse`: a summary (`totalRows`/`validRows`/`invalidRows`)
and a per-row result (`valid`, `errors`, and — for the commit endpoint —
whether the vocabulary/lesson link was `CREATE`d or `UPDATE`d).

### Excel format

First row = header (case-insensitive, spaces ignored). One data row = one word
assigned to one lesson.

Required columns:

| Column | Notes |
| --- | --- |
| `word` | |
| `meaningVi` | |
| `topicSlug` | Must match an existing `topics.slug`. |
| `lessonSlug` | Must match an existing `lessons.slug` **under that topic** (lesson slugs are only unique per topic, not globally). |

Optional columns: `slug` (auto-generated from `word` if omitted), `ipa`,
`audioUrl` (pronunciation audio: a path served by the frontend such as
`/audio/vocab/a1-apple.m4a`, or an absolute URL; max 500 characters),
`partOfSpeech`, `meaningEn`, `shortDefinition`, `exampleSentence`,
`exampleMeaningVi`, `exampleSource`, `synonyms`, `antonyms`, `collocations`
(comma-separated lists), `difficulty` (`BEGINNER`/`INTERMEDIATE`/`ADVANCED`,
default `INTERMEDIATE`), `targetBand` (number), `orderIndex` (int, default 0),
`isRequired` (`true`/`false`, default `true`).

To assign the same word to more than one lesson, repeat the row with a
different `topicSlug`/`lessonSlug`.

### Upsert behavior

- **Vocabulary**: matched by `slug`. If it exists, all fields are overwritten
  with the file's values (this doubles as a bulk-edit mechanism). If not, it's
  created.
- **Lesson link** (`LessonVocabulary`): matched by `(lessonId, vocabularyId)`.
  If it exists, `orderIndex`/`isRequired` are updated. If not, it's created.
  Existing links to *other* lessons are never touched or removed.

### Validation failure reasons a row can get

- Missing required column value.
- Invalid `difficulty`, `targetBand`, `orderIndex`, or `isRequired` format.
- `topicSlug` not found.
- `lessonSlug` not found under that `topicSlug`.
- Same `(lessonSlug, word)` pair repeated elsewhere in the same file.

Invalid rows are skipped on commit; valid rows in the same file still import.

## 3. Foundation English course (A1 vocabulary library)

`go run ./cmd/seed -only-foundation` creates the **Foundation English** course
(`slug = foundation-english`, `kind = CEFR`) with level A1 and its vocabulary
library: 14 topics / 92 lessons with `track = VOCABULARY` (always unlocked).
Topics and lessons are defined in `backend/seeds/data/foundation_a1.json`, which
is generated together with the matching Excel import file, so slugs line up.

```bash
cd backend
go run ./cmd/seed -only-foundation   # 1. create course, A1 level, topics, lessons
# 2. import content/a1-vocabulary/a1_vocabulary_import.xlsx at /admin/vocabularies/import
go run ./cmd/seed -only-foundation   # 3. build each lesson's quiz from the imported words
```

`-only-foundation` is idempotent and never touches demo users or progress.
Databases seeded before Task 49 had "A1 Foundation" inside the IELTS course
(band 3.0); the command moves that band level to the Foundation English course
and keeps its IDs, so earlier A1 progress is preserved.

Quiz questions are Vietnamese-meaning multiple choice ("X nghĩa là gì?"),
up to 5 per lesson. Implementation: `backend/seeds/foundation.go`.

Learners switch course on the roadmap (`GET /courses`, `PATCH /me/active-course`).

Pronunciation audio files go in `frontend/public/audio/vocab/<vocabulary-slug>.m4a`
(generated by `content/a1-vocabulary/generate_audio.sh`). Speaker buttons fall
back to the browser's speech synthesis when a file is missing.

## 4. Importing lesson content (JSON bundles)

Skill lessons (topics → lessons → activities → exercises) are imported from a
JSON bundle — format in `docs/foundation-spec.md` §8.2, example
`content/a1/week-1.json` (generated from `content/a1/build_week1.py`).

- UI: `/admin/content/import` (sidebar "Import Lessons", admins only). Page:
  `frontend/src/pages/AdminContentImportPage.tsx`.
- API: `POST /api/v1/admin/content/import/preview` (no writes) and
  `POST /api/v1/admin/content/import` (one transaction), multipart field
  `file`, max 5MB. Module: `backend/internal/modules/contentimport/`.

Rules:

- The bundle names a course and level (`courseSlug`, `cefrLevel`); the level must
  exist (`go run ./cmd/seed -only-foundation`).
- Topics, lessons and activities are upserted by slug; exercises by their
  position inside the activity (extra old exercises are deleted). Re-importing
  a file therefore updates it in place and keeps learner progress.
- Unknown JSON fields reject the whole file (catches typos).
- Every item is validated; each exercise is also self-checked by grading its
  expected answer. Invalid items — and their children — are skipped and listed
  in the report; valid items still import.
- Activities can be `DRAFT` (hidden from learners) or `PUBLISHED` (default).
- `VOCABULARY` activities list word slugs in `content.vocabularySlugs`; those
  words must already be imported.

Passages (listening and reading):

- A bundle's `passages` hold texts and recordings (upserted by slug);
  LISTENING / READING activities use them through `passageSlug`, and a
  DICTATION item can play one sentence with `"passageSlug": "...", "segmentIndex": 2`.
- **Import Audio** (`/admin/passages/import`): upload an mp3/m4a (max 25MB),
  Whisper returns the transcript split into timed sentences; fix words, add
  speakers, then *Save passage* (or download JSON for the repo). Recordings are
  stored in `MEDIA_DIR` (default `backend/storage/media`, git-ignored) and
  served at `/media`. Needs `OPENAI_API_KEY`.
- **Extract from PDF** (`/admin/content/extract`): upload up to 10 page photos
  or one PDF; the AI writes a lesson bundle (activities DRAFT by default) and
  opens it in Import Lessons for review. Check every answer before importing.

Sound cards (IPA):

- Import `content/a1/a1-sounds.json` (built by `content/a1/build_sounds.py`)
  before the week bundles: 40 phonemes (upserted by symbol) plus the always-open
  "Sound Lab" practice topic. PRONUNCIATION activities name their sounds in
  `content.phonemes`; an unknown symbol gives a warning.
- Learners see the chart at `/pronunciation` and a sound card at
  `/pronunciation/:symbol` (API `GET /api/v1/phonemes`, module
  `backend/internal/modules/phoneme/`).

Grammar points:

- Import `content/a1/a1-grammar.json` (built by `content/a1/build_grammar.py`)
  **before** the week bundles. Grammar points are upserted by slug.
- A GRAMMAR activity links to a point with `grammarPointSlug`; a missing slug
  gives a warning and the activity imports without theory.
- Learners see theory and mastery at `/grammar` and `/grammar/:slug`
  (API `GET /api/v1/grammar-points`, module `backend/internal/modules/grammar/`).

Pronunciation scoring (Azure Speech):

- Set `AZURE_SPEECH_KEY` and `AZURE_SPEECH_REGION` in `backend/.env` (Azure portal →
  Speech service → Keys and Endpoint) and restart the API. `GET /api/v1/pronunciation/status`
  then returns `enabled: true`; without a key SPEAK exercises stay a self-check.
- `PRONUNCIATION_DAILY_LIMIT` (default 300 recordings per learner per day) protects the
  free tier (F0: 5 audio hours per month).
- SPEAK exercises can set `answer.minAccuracy` (1-100, default 70).

Role-play scenarios:

- Import `content/a1/a1-scenarios.json` (built by `content/a1/build_scenarios.py`)
  **after** `a1-grammar.json` and **before** the week bundles: 40 system scenarios
  (36 lessons + 4 weekly checkpoints)
  (upserted by slug), listed for every learner at `/ai-conversation`. A scenario's
  `passageSlug` must already exist (Week 1 dialogues come with `week-1.json`), so
  import `a1-scenarios.json` again after the week bundles (`week-1.json` … `week-4.json`) to attach them.
- Weekly checkpoints: a CHECKPOINT activity with `content.itemCount` draws that many
  items per attempt from its exercises, and with `scenarioSlug` ends with a role-play
  without hints. Passing the week's last lesson unlocks the next week's first lesson.
- A SPEAKING activity links to a scenario with `scenarioSlug`; it completes with the
  learner's latest finished session (passes only when the task was completed).
- AI chat quota per learner per day: `AI_CHAT_DAILY_LIMIT` (default 200).

Lesson audio: `content/a1/generate_lesson_audio.sh` (macOS) writes
`frontend/public/audio/a1/**.m4a` from `audio_manifest_*.tsv`.
