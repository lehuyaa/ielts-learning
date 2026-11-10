# A1 Lesson Content (Foundation English)

Skill lessons for the A1 level, imported at `/admin/content/import`
(format: `docs/foundation-spec.md` §8.2).

| File | What it is |
|---|---|
| `build_week1.py` | Authoring source for Week 1 (edit this, not the JSON) |
| `week-1.json` | Generated bundle: 1 topic, 9 situation lessons + review, 49 activities, 184 exercises |
| `audio_manifest_week-1.tsv` | Sentences/words to synthesise (`path<TAB>text`) |
| `generate_lesson_audio.sh` | Creates the `.m4a` files with macOS text-to-speech |

## Workflow

```bash
cd content/a1
python3 build_week1.py --vocab ../a1-vocabulary/out/a1_rows.json   # --vocab optional: checks word slugs
./generate_lesson_audio.sh                                          # → frontend/public/audio/a1/w1/*.m4a
```

Then upload `week-1.json` at **Admin → Import Lessons**: Preview, then Confirm import.
Re-importing updates lessons in place (matched by slug; exercises by position),
so learner progress is kept. Every exercise is self-checked on import: the
expected answer must grade as correct, otherwise the item is skipped.

Prerequisites: `go run ./cmd/seed -only-foundation` and the A1 vocabulary import
(VOCABULARY activities reference words such as `a1-hello`).

## Week 1 · Meeting People

| # | Lesson | Sounds | Grammar |
|---|---|---|---|
| 1 | Greetings | /iː/ – /ɪ/ | to be: am / is / are |
| 2 | Introducing yourself | /e/ – /æ/ | I'm…, live in, from |
| 3 | Introducing others | /ʌ/ – /ɑː/ | This is…, He/She is… |
| 4 | Asking for personal information | /uː/ – /ʊ/ | Questions and negatives with to be |
| 5 | Talking about your family | /ɜːr/ – /ə/ | Plurals, a / an |
| 6 | Talking about hobbies | -s endings /s/ /z/ /ɪz/ | my / your / his / her |
| 7 | Describing your job | Word stress (jobs) | a / an with jobs |
| 8 | Finding places in the building | thirteen vs thirty | There is / there are + prepositions |
| 9 | Borrowing things at work | Review | this / that / these / those, Can I…? |
| R | Week 1 review | | Checkpoint (70 to pass) |

Reading and speaking activities are optional; speaking is a read-aloud
self-check until AI speaking and pronunciation scoring arrive (Tasks 67-72).
