import { useState } from "react";

import { cn } from "@/lib/utils";

import { WordSheet } from "./WordSheet";

const TOKEN = /([A-Za-z]+(?:['’][A-Za-z]+)*)/;

type TappableTextProps = {
  text: string;
  /** Its glossary is checked first. */
  passageId?: number;
  className?: string;
};

/**
 * Text whose words can be tapped to see their meaning, IPA and audio and to
 * add them to flashcards (Task 63). Line breaks are kept.
 */
export function TappableText({ text, passageId, className }: TappableTextProps) {
  const [selected, setSelected] = useState<{ word: string; key: string } | null>(null);
  const parts = text.split(TOKEN);

  return (
    <>
      <span className={cn("whitespace-pre-line", className)}>
        {parts.map((part, index) => {
          if (index % 2 === 0) {
            return <span key={index}>{part}</span>;
          }
          const key = `${index}`;
          return (
            <button
              className={cn(
                "cursor-pointer rounded-sm decoration-primary/40 decoration-dotted underline-offset-4 hover:bg-primary/10 hover:underline",
                selected?.key === key && "bg-primary/15",
              )}
              key={index}
              onClick={(event) => {
                event.stopPropagation();
                setSelected({ word: part, key });
              }}
              type="button"
            >
              {part}
            </button>
          );
        })}
      </span>
      {selected ? <WordSheet onClose={() => setSelected(null)} passageId={passageId} word={selected.word} /> : null}
    </>
  );
}
