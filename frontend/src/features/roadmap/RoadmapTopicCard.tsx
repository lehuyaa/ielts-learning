import { ArrowRight, CheckCircle2 } from "lucide-react";
import { Link } from "react-router-dom";

import { RoadmapProgressBar } from "@/features/roadmap/RoadmapProgressBar";
import type { RoadmapTopic } from "@/features/roadmap/types";

type RoadmapTopicCardProps = {
  topic: RoadmapTopic;
};

export function RoadmapTopicCard({ topic }: RoadmapTopicCardProps) {
  const isComplete = topic.progress === 100;

  return (
    <Link aria-label={`Open ${topic.title}`} to={`/topics/${topic.id}`}>
      <article className="rounded-2xl border border-[#e3e4f8] bg-white p-5 shadow-sm">
        <div className="flex items-center justify-between gap-4">
          <div className="flex min-w-0 items-center gap-3">
            <span className="text-2xl" aria-hidden="true">
              {topic.icon}
            </span>
            <h3 className="truncate text-xl font-bold tracking-normal text-[#10111f]">
              {topic.title}
            </h3>
          </div>

          {isComplete ? (
            <CheckCircle2
              className="size-5 shrink-0 text-success"
              aria-hidden="true"
            />
          ) : (
            <ArrowRight
              className="size-5 shrink-0 text-[#6d7088]"
              aria-hidden="true"
            />
          )}
        </div>

        <div className="mt-6">
          <RoadmapProgressBar
            value={topic.progress}
            variant={isComplete ? "success" : "primary"}
          />
        </div>

        <div className="mt-4 flex items-center justify-between gap-4 text-lg font-medium text-[#676982]">
          <span>
            {topic.completedLessons}/{topic.totalLessons} lessons
          </span>
          <span className={isComplete ? "text-[#009f73]" : undefined}>
            {topic.progress}%
          </span>
        </div>
      </article>
    </Link>
  );
}
