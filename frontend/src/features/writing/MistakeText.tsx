import { useState } from "react";

import { cn } from "@/lib/utils";
import type { WritingMistake } from "@/types/writing";

type Piece = { text: string; mistake?: number };

/** Splits the text into plain pieces and marked mistakes (by offsets, in order). */
function pieces(text: string, mistakes: WritingMistake[]): Piece[] {
  const located = mistakes
    .map((mistake, index) => ({ ...mistake, index }))
    .filter((mistake) => mistake.start >= 0 && mistake.end <= text.length && mistake.end > mistake.start)
    .sort((a, b) => a.start - b.start);
  const out: Piece[] = [];
  let cursor = 0;
  for (const mistake of located) {
    if (mistake.start < cursor) continue; // overlapping marks: keep the first
    if (mistake.start > cursor) out.push({ text: text.slice(cursor, mistake.start) });
    out.push({ text: text.slice(mistake.start, mistake.end), mistake: mistake.index });
    cursor = mistake.end;
  }
  if (cursor < text.length) out.push({ text: text.slice(cursor) });
  return out;
}

/**
 * The learner's text with each mistake struck through in red, the fix in
 * green after it, and the Vietnamese explanation of the tapped mistake.
 */
export function MistakeText({ text, mistakes }: { text: string; mistakes: WritingMistake[] }) {
  const [selected, setSelected] = useState<number | null>(mistakes.length > 0 ? 0 : null);
  const unplaced = mistakes.map((mistake, index) => ({ mistake, index })).filter(({ mistake }) => mistake.start < 0);
  const current = selected !== null ? mistakes[selected] : null;

  return (
    <div>
      <p className="whitespace-pre-line rounded-2xl bg-muted/30 p-4 text-lg leading-loose">
        {pieces(text, mistakes).map((piece, index) =>
          piece.mistake === undefined ? (
            <span key={index}>{piece.text}</span>
          ) : (
            <button
              className={cn(
                "cursor-pointer rounded-md px-0.5 text-left",
                selected === piece.mistake ? "bg-red-100 ring-2 ring-red-300" : "hover:bg-red-50",
              )}
              key={index}
              onClick={() => setSelected(piece.mistake ?? null)}
              type="button"
            >
              <span className="text-red-600 line-through decoration-2">{piece.text}</span>{" "}
              <span className="font-semibold text-emerald-700">{mistakes[piece.mistake].correct}</span>
            </button>
          ),
        )}
      </p>

      {unplaced.length > 0 ? (
        <ul className="mt-2 space-y-1 text-sm">
          {unplaced.map(({ mistake, index }) => (
            <li key={index}>
              <button className="cursor-pointer text-left hover:underline" onClick={() => setSelected(index)} type="button">
                <span className="text-red-600 line-through">{mistake.wrong}</span> →{" "}
                <span className="font-semibold text-emerald-700">{mistake.correct}</span>
              </button>
            </li>
          ))}
        </ul>
      ) : null}

      {current ? (
        <div className="mt-3 rounded-2xl border border-red-200 bg-white p-3 text-sm" role="status">
          <p>
            <span className="text-red-600 line-through">{current.wrong}</span> →{" "}
            <span className="font-semibold text-emerald-700">{current.correct}</span>
            {current.skillTag ? (
              <span className="ml-2 rounded bg-muted px-1.5 py-0.5 text-xs text-muted-foreground">{current.skillTag}</span>
            ) : null}
          </p>
          {current.explanationVi ? <p className="mt-1">{current.explanationVi}</p> : null}
        </div>
      ) : null}
    </div>
  );
}
