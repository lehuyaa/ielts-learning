# A1 Lesson Content (Foundation English)

Skill lessons for the A1 level, imported at `/admin/content/import`
(format: `docs/foundation-spec.md` §8.2).

| File | What it is |
|---|---|
| `build_grammar.py` → `a1-grammar.json` | The 9 A1 grammar theory cards (G1-G9). **Import first** |
| `build_sounds.py` → `a1-sounds.json` | The 40 sound cards (IPA chart) + the always-open "Sound Lab" topic: 178 minimal pairs, ending -s / -ed, 2-syllable stress. **Import first** |
| `audio_manifest_sounds.tsv` | Words and minimal-pair audio for the sound cards and Sound Lab |
| `lesson_kit.py` | Shared helpers for the week builders: exercise / activity helpers, role-plays, writing tasks, checkpoint pool, output files |
| `build_week1.py` … `build_week4.py` | Authoring source for each week (edit these, not the JSON) |
| `build_scenarios.py` → `a1-scenarios.json` | The 40 AI role-play scenarios: 9 per week (Task D10) + one per weekly checkpoint (Task 75): role, task, target patterns, max turns, opening line. **Import before the week bundles** |
| `writing_prompts.py` | The 20 A1 writing prompts with frames and model answers (Task D9); weeks add them with `writing_activity(key)` |
| `week-1.json` | Generated bundle: 18 passages (9 dialogues, 9 reading texts), 1 topic, 9 situation lessons + review, 69 activities (5 optional writing tasks, 9 optional AI role-plays), 229 exercises |
| `week-2.json` | Generated bundle: 18 passages, 1 topic, 9 situation lessons + review, 69 activities (5 optional writing tasks, 9 optional AI role-plays), 297 exercises |
| `week-3.json` | Generated bundle: 18 passages, 1 topic, 9 situation lessons + review, 69 activities, 307 exercises |
| `week-4.json` | Generated bundle: 18 passages, 1 topic, 9 situation lessons + review, 69 activities, 301 exercises |
| `audio_manifest_week-N.tsv` | Sentences/words to synthesise (`path<TAB>text[<TAB>voice]`) |
| `generate_lesson_audio.sh` | Creates the `.m4a` files with macOS text-to-speech |

## Workflow

```bash
cd content/a1
python3 build_grammar.py
python3 build_scenarios.py
python3 build_sounds.py --cmu path/to/cmudict.dict                  # CMU dict: github.com/cmusphinx/cmudict
python3 build_week1.py --vocab ../a1-vocabulary/out/a1_rows.json   # --vocab optional: checks word slugs
python3 build_week2.py --vocab ../a1-vocabulary/out/a1_rows.json
python3 build_week3.py --vocab ../a1-vocabulary/out/a1_rows.json
python3 build_week4.py --vocab ../a1-vocabulary/out/a1_rows.json
./generate_lesson_audio.sh                                          # → frontend/public/audio/a1/{w1,w2,w3,w4,sounds,words}/*.m4a
```

Then upload `a1-grammar.json`, `a1-sounds.json` and `a1-scenarios.json`, then `week-1.json` … `week-4.json`, and `a1-scenarios.json` once more (it links each scenario to its week's dialogue) at **Admin → Import Lessons**
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
| R | Week 1 review | | Checkpoint: 15 random items from an 86-item pool + a role-play without hints (70 to pass; unlocks Week 2) |

Reading and speaking activities are optional. Each lesson has a read-aloud
activity (slug `read-aloud`, SPEAK items scored by Azure when configured, otherwise a self-check) and an AI role-play "Nói chuyện với AI" (slug
`speaking`, `scenarioSlug` → `a1-scenarios.json`) that passes when the learner
completes the scenario task.

## Week 2 · Work and Money

| # | Lesson | Sounds | Grammar |
|---|---|---|---|
| 1 | Daily work schedule | /eɪ/ – /aɪ/ | Present simple: I work / she works, always / usually |
| 2 | Asking for leave | /aʊ/ – /oʊ/ | Present simple: don't / doesn't |
| 3 | Scheduling a meeting | /p/ – /b/ | Present simple: Do / Does questions; on Monday, at two |
| 4 | Booking a table | /t/ – /d/ | There is / there are, Is there…? |
| 5 | Ordering food | /k/ – /ɡ/ | Present simple: like / want / have |
| 6 | Paying the bill | /ɔɪ/ vs /aɪ/ | Wh- questions: How much does…? Where do…? |
| 7 | Asking prices | fifteen vs fifty (stress) | Prepositions of place: on, under, next to, between |
| 8 | Choosing a size | 2-syllable word stress | Present continuous: I'm looking for… |
| 9 | Returning items | Week review | Present simple or present continuous? |
| R | Week 2 review | | Checkpoint: 15 random items + a role-play without hints (70 to pass; unlocks Week 3) |

## Week 3 · Travel and Health

| # | Lesson | Sounds | Grammar |
|---|---|---|---|
| 1 | Asking for directions | /θ/ – /ð/ | Can you tell me…? + directions (imperatives) |
| 2 | Taking a taxi | /s/ – /z/ | Adjectives: position, very / too |
| 3 | Buying tickets | /ʃ/ – /ʒ/ | Comparatives: cheaper than, more expensive than |
| 4 | Airport check-in | /tʃ/ – /dʒ/ | must / mustn't |
| 5 | Hotel check-in and check-out | /f/ – /v/ | Can I…? / Can you…? |
| 6 | Train and bus times | /l/ – /r/ | Superlatives: the fastest, the most expensive |
| 7 | Describing symptoms | Final consonant clusters | should / shouldn't |
| 8 | Buying medicine | /s/ – /ʃ/ | can / must / should review |
| 9 | At the bank | Week review | Comparatives + can / must review |
| R | Week 3 review | | Checkpoint: 15 random items + a role-play without hints (70 to pass; unlocks Week 4) |

## Week 4 · Plans and Problems

| # | Lesson | Sounds | Grammar |
|---|---|---|---|
| 1 | Phone messages | /m/ – /n/ at the end of words | will: I'll tell him, I'll call you back |
| 2 | Invitations | /w/ vs /v/ | be going to (plans) |
| 3 | Arranging to meet | /h/ and /j/ | will or going to? Let's… |
| 4 | Changing plans | Past -ed: /t/ /d/ /ɪd/ | won't, not going to, can't |
| 5 | Weather small talk | Sentence stress | was / were, will (predictions) |
| 6 | Your weekend | Linking sounds | Past simple: regular and irregular |
| 7 | Sending a parcel | /n/ vs /ŋ/, -ing | Past simple questions: Did…? |
| 8 | Reporting a problem | Stress on can't / didn't | Past simple negatives: didn't |
| 9 | Emergency help | Week review | Past simple + will review |
| R | Week 4 review | | Checkpoint: 15 random items + a role-play without hints (70 to pass; completes A1) |

## Sound cards and Sound Lab

`build_sounds.py` holds the 40 cards (Vietnamese mouth notes and typical
mistakes), the minimal pairs and the ending/stress word sets. IPA comes from the
CMU dictionary converted to our Oxford NAmE style; every minimal pair is checked
to differ in exactly the two sounds it teaches (the script stops otherwise), and
`IPA_OVERRIDES` fixes words CMU gets wrong for this style. Minimal-pair audio
uses neutral file names (`sounds/mp-…`) so the answer is not visible in the URL.
