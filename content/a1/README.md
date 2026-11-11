# A1 Lesson Content (Foundation English)

Skill lessons for the A1 level, imported at `/admin/content/import`
(format: `docs/foundation-spec.md` §8.2).

| File | What it is |
|---|---|
| `build_grammar.py` → `a1-grammar.json` | The 9 A1 grammar theory cards (G1-G9). **Import first** |
| `build_sounds.py` → `a1-sounds.json` | The 40 sound cards (IPA chart) + the always-open "Sound Lab" topic: 178 minimal pairs, ending -s / -ed, 2-syllable stress. **Import first** |
| `audio_manifest_sounds.tsv` | Words and minimal-pair audio for the sound cards and Sound Lab |
| `build_week1.py` | Authoring source for Week 1 (edit this, not the JSON) |
| `build_scenarios.py` → `a1-scenarios.json` | The 36 AI role-play scenarios (Task D10), 9 per week: role, task, target patterns, max turns, opening line. **Import before the week bundles** |
| `writing_prompts.py` | The 20 A1 writing prompts with frames and model answers (Task D9); weeks add them with `writing_activity(key)` |
| `week-1.json` | Generated bundle: 18 passages (9 dialogues, 9 reading texts), 1 topic, 9 situation lessons + review, 69 activities (5 optional writing tasks, 9 optional AI role-plays), 229 exercises |
| `audio_manifest_week-1.tsv` | Sentences/words to synthesise (`path<TAB>text[<TAB>voice]`) |
| `generate_lesson_audio.sh` | Creates the `.m4a` files with macOS text-to-speech |

## Workflow

```bash
cd content/a1
python3 build_grammar.py
python3 build_scenarios.py
python3 build_sounds.py --cmu path/to/cmudict.dict                  # CMU dict: github.com/cmusphinx/cmudict
python3 build_week1.py --vocab ../a1-vocabulary/out/a1_rows.json   # --vocab optional: checks word slugs
./generate_lesson_audio.sh                                          # → frontend/public/audio/a1/{w1,sounds,words}/*.m4a
```

Then upload `a1-grammar.json`, `a1-sounds.json` and `a1-scenarios.json`, then `week-1.json` at **Admin → Import Lessons**
(Preview, then Confirm import). GRAMMAR activities link to their theory card
through `content.grammarPointSlug`; exercises practise a card through the skill
tag `grammar:<slug>`, which also drives the learner's mastery on `/grammar`.
LISTENING / READING activities point to a passage with `passageSlug`; each
dialogue line is its own audio file (one passage segment, voiced per speaker),
which gives the player sentence replay, 0.75x speed and dictation items.
PRONUNCIATION activities list their sounds in `content.phonemes` (symbols as on
the sound cards); the player shows those cards above the exercises and
`/pronunciation/:symbol` lists the lessons that teach each sound.
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

Reading and speaking activities are optional. Each lesson has a read-aloud
self-check (slug `read-aloud`) and an AI role-play "Nói chuyện với AI" (slug
`speaking`, `scenarioSlug` → `a1-scenarios.json`) that passes when the learner
completes the scenario task.

## Sound cards and Sound Lab

`build_sounds.py` holds the 40 cards (Vietnamese mouth notes and typical
mistakes), the minimal pairs and the ending/stress word sets. IPA comes from the
CMU dictionary converted to our Oxford NAmE style; every minimal pair is checked
to differ in exactly the two sounds it teaches (the script stops otherwise), and
`IPA_OVERRIDES` fixes words CMU gets wrong for this style. Minimal-pair audio
uses neutral file names (`sounds/mp-…`) so the answer is not visible in the URL.
