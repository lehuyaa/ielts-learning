#!/usr/bin/env bash
# Generates audio for A1 lesson content (dialogues, dictation, sounds) with
# macOS text-to-speech, from every audio_manifest_*.tsv in this folder.
# Run on a Mac from this folder:   ./generate_lesson_audio.sh
# Options: VOICE=Alex RATE=150 FORCE=1 ./generate_lesson_audio.sh
# A manifest line is "path<TAB>text" or "path<TAB>text<TAB>voice" (dialogue
# speakers); a voice that is not installed falls back to VOICE.
set -euo pipefail
cd "$(dirname "$0")"

OUT="../../frontend/public/audio/a1"
VOICE="${VOICE:-Samantha}"
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
    if ! say -v "${voice:-$VOICE}" -r "$RATE" -o "$file" --file-format=m4af --data-format=aac "$text" 2>/dev/null; then
      say -v "$VOICE" -r "$RATE" -o "$file" --file-format=m4af --data-format=aac "$text"
    fi
    count=$((count + 1))
  done < "$manifest"
done

echo "Generated $count audio files in $OUT"
