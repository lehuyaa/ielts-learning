import { CheckCircle2, Eye, Volume2 } from "lucide-react";
import { useState } from "react";
import { Link } from "react-router-dom";

import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";
import type { LessonVocabularyItem } from "@/types/lesson";

type VocabularyPreviewListProps = {
  vocabulary: LessonVocabularyItem[];
};

export function VocabularyPreviewList({
  vocabulary,
}: VocabularyPreviewListProps) {
  const [expandedVocabularyId, setExpandedVocabularyId] = useState<
    string | null
  >(null);

  function toggleVocabulary(vocabularyId: string) {
    setExpandedVocabularyId((currentId) =>
      currentId === vocabularyId ? null : vocabularyId,
    );
  }

  return (
    <div className="grid gap-5">
      {vocabulary.map((item) => (
        <VocabularyCard
          expanded={expandedVocabularyId === item.id}
          item={item}
          key={item.id}
          onToggle={() => toggleVocabulary(item.id)}
        />
      ))}
    </div>
  );
}

type VocabularyCardProps = {
  item: LessonVocabularyItem;
  expanded: boolean;
  onToggle: () => void;
};

function VocabularyCard({ item, expanded, onToggle }: VocabularyCardProps) {
  const detailId = `vocabulary-detail-${item.id}`;

  return (
    <article
      className={cn(
        "rounded-2xl border bg-white p-5 shadow-sm",
        expanded ? "border-[#c6c4ff]" : "border-[#e6e6f3]",
      )}
    >
      <div className="flex flex-col gap-5 sm:flex-row sm:items-center">
        <div className="flex min-w-0 flex-1 items-start gap-5">
          <span className="mt-1 grid size-11 shrink-0 place-items-center rounded-full bg-[#45c486] text-white">
            <CheckCircle2 className="size-5" aria-hidden="true" />
          </span>

          <div className="min-w-0">
            <div className="flex flex-wrap items-baseline gap-x-3 gap-y-1">
              <h2 className="text-2xl font-bold tracking-normal text-[#10111f]">
                {item.word}
              </h2>
              <span className="font-mono text-base font-semibold text-[#85889c]">
                {item.ipa}
              </span>
              <span className="text-base font-semibold text-[#85889c]">
                {item.partOfSpeech}
              </span>
            </div>
            <p className="mt-2 line-clamp-1 text-lg font-medium text-[#676982]">
              {item.shortDefinition}
            </p>
          </div>
        </div>

        <div className="flex shrink-0 flex-wrap items-center gap-3 sm:justify-end">
          <DifficultyBadge difficulty={item.difficulty} />
          <span className="rounded-full bg-[#f0efff] px-3 py-1 text-base font-bold text-primary">
            {item.band}
          </span>
          <Button
            aria-label={`Listen to ${item.word}`}
            size="icon"
            variant="ghost"
          >
            <Volume2 className="size-5 text-[#6d7088]" aria-hidden="true" />
          </Button>
          <Button
            aria-controls={detailId}
            aria-expanded={expanded}
            className="rounded-full px-4 text-base cursor-pointer"
            onClick={onToggle}
            size="sm"
            type="button"
            variant="outline"
          >
            <Eye aria-hidden="true" />
            {expanded ? "Hide" : "View"}
          </Button>
        </div>
      </div>

      {expanded ? (
        <div className="mt-6 border-t border-[#e6e6f3] pt-5" id={detailId}>
          <p className="text-sm font-bold uppercase tracking-normal text-[#676982]">
            Definition
          </p>
          <p className="mt-2 text-lg font-medium text-[#232431]">
            {item.definition}
          </p>

          <p className="mt-5 text-sm font-bold uppercase tracking-normal text-[#676982]">
            IELTS Example
          </p>
          <blockquote className="mt-2 rounded-2xl border-l-4 border-primary bg-[#eff1ff] px-5 py-4 text-lg font-medium text-[#30348e]">
            "{item.example}"
          </blockquote>

          <div className="mt-4 flex flex-wrap gap-3">
            <Button
              asChild
              className="rounded-full px-5 text-white hover:text-white"
              size="sm"
              style={{
                color: "white",
              }}
            >
              <Link to={`/vocabulary/${item.id}`}>Full Detail</Link>
            </Button>
            <Button className="rounded-full px-5" size="sm" variant="outline">
              Practice Card
            </Button>
          </div>
        </div>
      ) : null}
    </article>
  );
}

type DifficultyBadgeProps = {
  difficulty: LessonVocabularyItem["difficulty"];
};

function DifficultyBadge({ difficulty }: DifficultyBadgeProps) {
  const styles = {
    Beginner: "bg-[#dcfce7] text-[#138a53]",
    Intermediate: "bg-[#fff3c4] text-[#c46700]",
    Advanced: "bg-[#fde5e7] text-[#d22f38]",
  };

  return (
    <span
      className={cn(
        "rounded-full px-3 py-1 text-base font-bold",
        styles[difficulty],
      )}
    >
      {difficulty}
    </span>
  );
}
