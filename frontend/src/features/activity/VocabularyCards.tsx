import { ArrowRight } from "lucide-react";
import { useState } from "react";

import { Button } from "@/components/ui/button";
import { AudioButton } from "@/features/exercise/AudioButton";
import { cn } from "@/lib/utils";
import type { VocabularyCard } from "@/types/activity";
import { IpaText } from "@/features/pronunciation/IpaText";

type VocabularyCardsProps = {
  cards: VocabularyCard[];
  onDone: () => void;
  doneLabel: string;
};

/** Word cards for a VOCABULARY activity: tap to flip, listen, then continue. */
export function VocabularyCards({ cards, onDone, doneLabel }: VocabularyCardsProps) {
  const [flipped, setFlipped] = useState<Set<number>>(new Set());
  const allSeen = flipped.size >= cards.length;

  return (
    <div>
      <p className="mb-4 text-center text-sm text-muted-foreground">
        Chạm vào thẻ để xem nghĩa · đã xem {flipped.size}/{cards.length}
      </p>
      <div className="grid gap-3 sm:grid-cols-2">
        {cards.map((card) => {
          const isFlipped = flipped.has(card.id);
          return (
            <div
              className={cn(
                "cursor-pointer rounded-2xl border-2 bg-white p-4 transition-colors",
                isFlipped ? "border-primary/40" : "border-border hover:border-primary/30",
              )}
              key={card.id}
              onClick={() => setFlipped((current) => new Set(current).add(card.id))}
              onKeyDown={(event) => {
                if (event.key === "Enter" || event.key === " ") {
                  setFlipped((current) => new Set(current).add(card.id));
                }
              }}
              role="button"
              tabIndex={0}
            >
              <div className="flex items-start justify-between gap-2">
                <div>
                  <p className="text-xl font-bold tracking-normal">{card.word}</p>
                  <p className="text-sm text-primary"><IpaText ipa={card.ipa} /></p>
                </div>
                <span onClick={(event) => event.stopPropagation()}>
                  <AudioButton
                    audioUrl={card.audioUrl}
                    fallbackText={card.word}
                    label="Nghe"
                    size="sm"
                  />
                </span>
              </div>
              {isFlipped ? (
                <div className="mt-3 border-t pt-3">
                  <p className="font-semibold">{card.meaningVi}</p>
                  {card.exampleSentence ? (
                    <p className="mt-1 text-sm">
                      {card.exampleSentence}
                      <span className="block text-muted-foreground">{card.exampleMeaningVi}</span>
                    </p>
                  ) : null}
                </div>
              ) : (
                <p className="mt-3 text-xs text-muted-foreground">{card.partOfSpeech} · chạm để xem nghĩa</p>
              )}
            </div>
          );
        })}
      </div>
      <div className="mt-6 flex justify-center">
        <Button onClick={onDone} size="lg" variant={allSeen ? "default" : "outline"}>
          {doneLabel}
          <ArrowRight />
        </Button>
      </div>
    </div>
  );
}
