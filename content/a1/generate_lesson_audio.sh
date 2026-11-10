#!/usr/bin/env bash
# Generates audio for A1 lesson content (dialogues, dictation, sounds) with
# macOS text-to-speech, from every audio_manifest_*.tsv in this folder.
# Run on a Mac from this folder:   ./generate_lesson_audio.sh
# Options: VOICE=Alex RATE=150 FORCE=1 ./generate_lesson_audio.sh
set -euo pipefail
cd "$(dirname "$0")"

OUT="../../frontend/public/audio/a1"
VOICE="${VOICE:-Samantha}"
RATE="${RATE:-150}"
FORCE="${FORCE:-0}"

count=0
for manifest in audio_manifest_*.tsv; do
  while IFS=$'\t' read -r path text; do
    [ -z "$path" ] && continue
    file="$OUT/$path.m4a"
    if [ -f "$file" ] && [ "$FORCE" != "1" ]; then
      continue
    fi
    mkdir -p "$(dirname "$file")"
    say -v "$VOICE" -r "$RATE" -o "$file" --file-format=m4af --data-format=aac "$text"
    count=$((count + 1))
  done < "$manifest"
done

echo "Generated $count audio files in $OUT"
