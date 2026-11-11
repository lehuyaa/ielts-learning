#!/usr/bin/env bash
# Generates audio for A1 lesson content (dialogues, dictation, sounds) with
# macOS text-to-speech, from every audio_manifest_*.tsv in this folder.
# Run on a Mac from this folder:   ./generate_lesson_audio.sh
# Options: VOICE=Samantha MALE_VOICE=Daniel FEMALE_VOICE=Samantha RATE=150 FORCE=1 ./generate_lesson_audio.sh
# A manifest line is "path<TAB>text" or "path<TAB>text<TAB>voice"; voice is
# "male" / "female" (dialogue speakers) or a macOS voice name. List installed
# voices with:  say -v '?' | grep en_
set -euo pipefail
cd "$(dirname "$0")"

OUT="../../frontend/public/audio/a1"
VOICE="${VOICE:-Samantha}"

# Installed voice names (e.g. "Samantha", "Evan (Enhanced)").
INSTALLED="$(say -v '?' | sed -E 's/ {2,}.*//')"

# First installed voice of the list, else $VOICE.
pick_voice() {
  for name in "$@"; do
    [ -n "$name" ] && grep -qxF "$name" <<<"$INSTALLED" && { echo "$name"; return; }
  done
  echo "$VOICE"
}
MALE_VOICE="$(pick_voice "${MALE_VOICE:-}" "Evan (Enhanced)" "Nathan (Enhanced)" "Tom (Enhanced)" "Aaron" "Alex" "Daniel" "Ralph" "Fred")"
FEMALE_VOICE="$(pick_voice "${FEMALE_VOICE:-}" "Ava (Premium)" "Zoe (Premium)" "Samantha (Enhanced)" "Samantha" "Kathy")"
echo "Voices: default=$VOICE, male=$MALE_VOICE, female=$FEMALE_VOICE"
RATE="${RATE:-150}"
FORCE="${FORCE:-0}"

count=0
for manifest in audio_manifest_*.tsv; do
  while IFS=$'\t' read -r path text voice; do
    [ -z "$path" ] && continue
    file="$OUT/$path.m4a"
    if [ -f "$file" ] && [ "$FORCE" != "1" ]; then
      continue
    fi
    mkdir -p "$(dirname "$file")"
    case "$voice" in
      male) voice="$MALE_VOICE" ;;
      female) voice="$FEMALE_VOICE" ;;
    esac
    if ! say -v "${voice:-$VOICE}" -r "$RATE" -o "$file" --file-format=m4af --data-format=aac "$text" 2>/dev/null; then
      say -v "$VOICE" -r "$RATE" -o "$file" --file-format=m4af --data-format=aac "$text"
    fi
    count=$((count + 1))
  done < "$manifest"
done

echo "Generated $count audio files in $OUT"
