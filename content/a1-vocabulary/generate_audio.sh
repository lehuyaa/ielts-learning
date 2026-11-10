#!/usr/bin/env bash
# Generates pronunciation audio for every A1 word with macOS text-to-speech.
# Run on a Mac from this folder:  ./generate_audio.sh
# Optional: VOICE=Alex ./generate_audio.sh   (list voices: say -v '?' | grep en_US)
set -euo pipefail
cd "$(dirname "$0")"

OUT="../../frontend/public/audio/vocab"
VOICE="${VOICE:-Samantha}"
RATE="${RATE:-150}"
mkdir -p "$OUT"

count=0
while IFS=$'\t' read -r slug word; do
  [ -z "$slug" ] && continue
  if [ -f "$OUT/$slug.m4a" ]; then
    continue
  fi
  say -v "$VOICE" -r "$RATE" -o "$OUT/$slug.m4a" --file-format=m4af --data-format=aac "$word"
  count=$((count + 1))
done < audio_manifest.tsv

echo "Generated $count new audio files in $OUT"
