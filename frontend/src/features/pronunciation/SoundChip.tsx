import { Volume2 } from "lucide-react";
import { Link } from "react-router-dom";

import { phonemePath } from "@/lib/ipa";
import { playPronunciation } from "@/lib/pronunciation";
import type { PhonemeCard } from "@/types/phoneme";

import { IpaText } from "./IpaText";

/** Compact sound card shown at the top of a PRONUNCIATION activity. */
export function SoundChip({ phoneme }: { phoneme: PhonemeCard }) {
  const examples = phoneme.exampleWords.slice(0, 3);
  return (
    <div className="flex min-w-0 flex-1 basis-60 items-start gap-3 rounded-2xl border bg-white p-3">
      <button
        aria-label={`Nghe ví dụ âm /${phoneme.symbol}/`}
        className="flex size-14 shrink-0 cursor-pointer flex-col items-center justify-center rounded-xl bg-primary/10 font-mono text-xl font-bold text-primary hover:bg-primary/20"
        onClick={() => playPronunciation(examples.map((word) => word.word).join(", "), phoneme.audioUrl)}
        type="button"
      >
        /{phoneme.symbol}/
        <Volume2 className="size-3" />
      </button>
      <div className="min-w-0 text-sm">
        <p className="font-semibold">{phoneme.nameVi}</p>
        <p className="mt-0.5 flex flex-wrap gap-x-2 text-muted-foreground">
          {examples.map((word) => (
            <button
              className="cursor-pointer hover:text-primary"
              key={word.word}
              onClick={() => playPronunciation(word.word, word.audioUrl)}
              type="button"
            >
              {word.word} <IpaText ipa={word.ipa} linked={false} className="text-xs" />
            </button>
          ))}
        </p>
        <Link className="mt-1 inline-block text-xs font-semibold text-primary hover:underline" to={phonemePath(phoneme.symbol)}>
          Xem khẩu hình →
        </Link>
      </div>
    </div>
  );
}
