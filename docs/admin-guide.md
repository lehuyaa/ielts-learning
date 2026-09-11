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
