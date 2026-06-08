import { ArrowRight, CheckCircle2, Lock } from "lucide-react";
import { Link } from "react-router-dom";

import { RoadmapProgressBar } from "@/features/roadmap/RoadmapProgressBar";
import type { RoadmapTopic } from "@/features/roadmap/types";
import { cn } from "@/lib/utils";

type RoadmapTopicCardProps = {
  topic: RoadmapTopic;
  isLocked: boolean;
};

export function RoadmapTopicCard({ topic, isLocked }: RoadmapTopicCardProps) {
  const isComplete = topic.progress === 100;
  const content = (
    <article
      className={cn(
        "relative w-full rounded-2xl border p-4 text-left transition-all",
        isLocked
          ? "cursor-not-allowed border-border/50 bg-muted/30"
          : "cursor-pointer border-border bg-white hover:-translate-y-0.5 hover:shadow-md",
      )}
    >
      {isLocked ? (
        <div className="absolute inset-0 z-10 flex items-center justify-center rounded-2xl bg-white/60">
          <Lock className="size-5 text-muted-foreground" aria-hidden="true" />
        </div>
      ) : null}

      <div className={cn("transition-all", isLocked && "opacity-40")}>
        <div className="mb-3 flex items-center justify-between gap-3">
          <div className="flex min-w-0 items-center gap-2">
            <span className="text-xl" aria-hidden="true">
              {topic.icon}
            </span>
            <h3 className="truncate text-sm font-semibold tracking-normal text-foreground">
              {topic.title}
            </h3>
          </div>

          {isComplete ? (
            <CheckCircle2
              className="size-4 shrink-0 text-emerald-500"
              aria-hidden="true"
            />
          ) : (
            <ArrowRight
              className="size-3.5 shrink-0 text-muted-foreground"
              aria-hidden="true"
            />
          )}
        </div>

        <RoadmapProgressBar
          value={topic.progress}
          variant={isComplete ? "success" : "primary"}
        />

        <div className="mt-2 flex items-center justify-between gap-3 text-xs text-muted-foreground">
          <span>
            {topic.completedLessons}/{topic.totalLessons} lessons
          </span>
          <span className={isComplete ? "font-medium text-emerald-600" : undefined}>
            {topic.progress}%
          </span>
        </div>
      </div>
    </article>
  );

  if (isLocked) {
    return content;
  }

  return (
    <Link aria-label={`Open ${topic.title}`} to={`/topics/${topic.id}`}>
      {content}
    </Link>
  );
}

