# A1 Vocabulary Content

Source of truth for the A1 Foundation vocabulary (Milestone 11, Tasks D1-D4).

- `src/*.txt` — authored word lists. `## topic-slug | Title | Vietnamese description | icon | color`
  starts a topic, `### Lesson title` starts a lesson, and each word line is
  `word|partOfSpeech|meaningVi|simple English definition|example|example in Vietnamese`.
- `build_a1.py` — builds the import file and the seed structure.
- `a1_vocabulary_import.xlsx` — upload at `/admin/vocabularies/import`.
- `generate_audio.sh` + `audio_manifest.tsv` — pronunciation audio (see below).

## Audio

Every row's `audioUrl` points to `/audio/vocab/<slug>.m4a`. Generate the files on a Mac
(offline, built-in text-to-speech, American voice):

```bash
cd content/a1-vocabulary
./generate_audio.sh              # default voice: Samantha
VOICE=Alex ./generate_audio.sh   # another en_US voice
```

Files land in `frontend/public/audio/vocab/` and are served by Vite. Until a file exists,
the speaker button falls back to the browser's speech synthesis, so importing first is fine.
To replace a voice later, delete the `.m4a` files and run the script again.

Word list based on the CEFR-J Wordlist (A1) plus workplace words for adult learners.
IPA (American English) comes from the CMU Pronouncing Dictionary.

## Rebuild after editing `src/*.txt`

```bash
cd content/a1-vocabulary
curl -L -o /tmp/cmudict.dict https://raw.githubusercontent.com/cmusphinx/cmudict/master/cmudict.dict
pip install openpyxl
python3 build_a1.py src /tmp/cmudict.dict out
cp out/foundation_a1.json ../../backend/seeds/data/foundation_a1.json
cp out/a1_vocabulary_import.xlsx out/audio_manifest.tsv .
```

Lessons longer than 15 words are split automatically ("Months & Seasons (1)", "(2)").
A word that appears twice with the same part of speech keeps only its first occurrence.

Then run `go run ./cmd/seed -only-foundation`, import the Excel file, and run the
seed command again to rebuild the lesson quizzes.
