import {
  ArrowLeft,
  BookOpen,
  Clock3,
  GraduationCap,
  Trophy,
} from "lucide-react";
import { useNavigate } from "react-router-dom";

import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";
import type { LessonDetailViewModel } from "@/types/lesson";

type LessonHeaderProps = {
  lesson: LessonDetailViewModel;
  learnedCount: number;
  progressPercentage: number;
  isStartDisabled: boolean;
  isStarting: boolean;
  onStartFlashcards: () => void;
};

export function LessonHeader({
  lesson,
  learnedCount,
  progressPercentage,
  isStartDisabled,
  isStarting,
  onStartFlashcards,
}: LessonHeaderProps) {
  const navigate = useNavigate();

  return (
    <header className="border-b border-[#e6e6f3] bg-white">
      <div className="mx-auto grid min-h-[72px] max-w-[1320px] gap-4 px-4 py-3 md:grid-cols-[240px_1fr_auto] md:items-center md:px-6">
        <div className="flex min-w-0 items-center gap-4">
          <button
            aria-label="Go back"
            className="grid size-9 shrink-0 cursor-pointer place-items-center rounded-full text-[#6d7088] transition-colors hover:bg-[#f0f1fb]"
            onClick={() => navigate(-1)}
            type="button"
          >
            <ArrowLeft className="size-5" aria-hidden="true" />
          </button>
          <div className="min-w-0">
            <h1 className="truncate text-xl font-bold tracking-normal text-[#10111f]">
              {lesson.title}
            </h1>
            <p className="mt-1 text-base font-medium text-[#676982]">
              Band {lesson.bandRange} - {lesson.vocabulary.length} words
            </p>
          </div>
        </div>

        <div className="min-w-0">
          <div className="h-3 overflow-hidden rounded-full bg-[#ebeaff]">
            <div
              className="h-full rounded-full bg-primary"
              style={{ width: `${progressPercentage}%` }}
            />
          </div>
          <p className="mt-2 text-sm font-medium text-[#676982]">
            {learnedCount}/{lesson.vocabulary.length} words learned
          </p>
        </div>

        <div className="flex flex-wrap items-center gap-3 md:justify-end">
          <LessonHeaderMeta icon={BookOpen} label={lesson.topic} />
          <LessonHeaderMeta
            icon={Clock3}
            label={`~${lesson.estimatedMinutes} min`}
          />
          <LessonHeaderMeta icon={Trophy} label={`${lesson.xpReward} XP`} />
          <Button
            className="h-10 rounded-full px-5 text-sm"
            disabled={isStartDisabled}
            onClick={onStartFlashcards}
            type="button"
          >
            <GraduationCap aria-hidden="true" />
            {isStarting ? "Starting..." : "Start Flashcards"}
          </Button>
        </div>
      </div>
    </header>
  );
}

type LessonHeaderMetaProps = {
  icon: typeof BookOpen;
  label: string;
};

function LessonHeaderMeta({ icon: Icon, label }: LessonHeaderMetaProps) {
  return (
    <span
      className={cn(
        "hidden items-center gap-2 rounded-full border border-[#ecebff] bg-[#f8f8ff]",
        "px-3 py-1.5 text-sm font-bold text-[#676982] xl:inline-flex",
      )}
    >
      <Icon className="size-4" aria-hidden="true" />
      {label}
    </span>
  );
}
