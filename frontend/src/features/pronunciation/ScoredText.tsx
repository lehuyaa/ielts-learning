import { useState } from "react";
import { Link } from "react-router-dom";

import { phonemePath } from "@/lib/ipa";
import { cn } from "@/lib/utils";
import type { PronunciationResult, WordScore } from "@/types/pronunciation";

import { scoreBg, scoreTone } from "./scoreColors";

/** The reference text with each word coloured; tap a word to see its sounds. */
export function ScoredText({ words, className }: { words: WordScore[]; className?: string }) {
  const [open, setOpen] = useState<number | null>(() => {
    const worst = words.reduce((best, word, index) => (word.errorType !== "Omission" && (best < 0 || word.accuracy < words[best].accuracy) ? index : best), -1);
    return worst >= 0 && words[worst].accuracy < 80 ? worst : null;
  });
  const selected = open !== null ? words[open] : undefined;

  return (
    <div className={className}>
      <p className="flex flex-wrap gap-x-1.5 gap-y-1 text-lg">
        {words.map((word, index) => (
          <button
            className={cn(
              "cursor-pointer rounded-md px-1 font-semibold underline-offset-4 hover:bg-muted",
              word.errorType === "Omission" ? "text-muted-foreground line-through" : scoreTone(word.accuracy),
              open === index && "bg-muted underline",
            )}
            key={`${word.word}-${index}`}
            onClick={() => setOpen(open === index ? null : index)}
            title={word.errorType === "Omission" ? "Bạn bỏ sót từ này" : `${word.accuracy}/100`}
            type="button"
          >
            {word.word}
          </button>
        ))}
      </p>
      {selected ? <WordDetail word={selected} /> : null}
    </div>
  );
}

function WordDetail({ word }: { word: WordScore }) {
  if (word.errorType === "Omission") {
    return <p className="mt-2 text-sm text-muted-foreground">Bạn bỏ sót từ “{word.word}”.</p>;
  }
  const weak = word.phonemes.filter((phoneme) => phoneme.accuracy < 60 && phoneme.symbols.length > 0);
  return (
    <div className="mt-2 rounded-2xl border bg-white p-3">
      <div className="flex flex-wrap items-baseline gap-2">
        <span className="font-semibold">{word.word}</span>
        <span className="font-mono text-lg">
          /
          {word.phonemes.map((phoneme, index) => (
            <span className={cn("font-semibold", scoreTone(phoneme.accuracy))} key={index} title={`${phoneme.accuracy}/100`}>
              {phoneme.display}
            </span>
          ))}
          /
        </span>
        <span className={cn("ml-auto text-sm font-bold", scoreTone(word.accuracy))}>{word.accuracy}/100</span>
      </div>
      {weak.length > 0 ? (
        <p className="mt-2 flex flex-wrap items-center gap-1.5 text-sm text-muted-foreground">
          Âm cần luyện:
          {weak.map((phoneme, index) => (
            <Link
              className="rounded-full border border-red-200 bg-red-50 px-2 py-0.5 font-mono font-semibold text-red-700 hover:bg-red-100"
              key={index}
              to={phonemePath(phoneme.symbols[0])}
            >
              /{phoneme.symbols[0]}/ {phoneme.accuracy}
            </Link>
          ))}
        </p>
      ) : null}
    </div>
  );
}

/** Overall scores + coloured text. */
export function PronunciationScore({ result, compact = false }: { result: PronunciationResult; compact?: boolean }) {
  const scores = [
    { label: "Chính xác", value: result.accuracy },
    { label: "Trôi chảy", value: result.fluency },
    { label: "Đầy đủ", value: result.completeness },
  ];
  return (
    <div className={cn("space-y-3 text-left", !compact && "rounded-2xl border bg-muted/20 p-4")}>
      <div className="flex flex-wrap items-center gap-2">
        {scores.map((score) => (
          <span className={cn("rounded-full border px-3 py-1 text-sm", scoreBg(score.value))} key={score.label}>
            {score.label} <span className={cn("font-bold", scoreTone(score.value))}>{score.value}</span>
          </span>
        ))}
        {result.minAccuracy !== undefined ? (
          <span className="text-xs text-muted-foreground">Cần độ chính xác ≥ {result.minAccuracy}</span>
        ) : null}
      </div>
      <ScoredText key={result.attemptId} words={result.words} />
    </div>
  );
}
