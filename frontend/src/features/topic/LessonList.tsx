import {
  ArrowRight,
  BookOpen,
  CheckCircle2,
  Clock3,
  Lock,
  Star,
  Zap,
} from "lucide-react";
import { Link } from "react-router-dom";

import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";
import type { TopicDetailViewModel, TopicLessonViewModel } from "@/types/topic";

type LessonListProps = {
  topic: TopicDetailViewModel;
};

export function LessonList({ topic }: LessonListProps) {
  const availableLessons = topic.lessons.filter(
    (lesson) => lesson.status !== "locked",
  ).length;

  return (
    <section>
      <div className="mb-6 flex items-center justify-between gap-4">
        <h2 className="text-xl font-bold tracking-normal text-[#10111f]">
          Lessons
        </h2>
        <p className="text-base font-medium text-[#676982]">
          {topic.completedLessons} of {topic.totalLessons} completed
        </p>
      </div>

      <div className="grid gap-4">
        {topic.lessons.map((lesson) => (
          <LessonCard key={lesson.id} lesson={lesson} />
        ))}
      </div>

      <section className="mt-8 rounded-2xl bg-linear-to-br from-[#5d55f1] to-[#7c2ef0] p-6 text-white shadow-sm md:flex md:items-center md:justify-between">
        <div>
          <h2 className="text-xl font-bold tracking-normal">
            Ready to continue learning?
          </h2>
          <p className="mt-2 text-base font-medium text-white/85">
            {availableLessons} lessons available
          </p>
        </div>
        <Button
          asChild
          className="mt-5 h-12 rounded-full bg-white px-6 text-base font-bold text-primary hover:bg-white/95 md:mt-0"
          variant="secondary"
        >
          <Link to={`/lessons/${getContinueLessonId(topic.lessons)}`}>
            Continue Learning
            <Zap className="size-5" aria-hidden="true" />
          </Link>
        </Button>
      </section>
    </section>
  );
}

type LessonCardProps = {
  lesson: TopicLessonViewModel;
};

function LessonCard({ lesson }: LessonCardProps) {
  const isLocked = lesson.status === "locked";
  const content = (
    <article
      className={cn(
        "relative overflow-hidden rounded-2xl border border-[#e6e6f3] bg-white p-6 shadow-sm",
        "transition-colors",
        !isLocked && "cursor-pointer hover:border-[#cbc9ff]",
        isLocked && "cursor-not-allowed bg-white/55 text-[#cfd0d8]",
      )}
    >
      <div
        className={cn(
          "flex items-start justify-between gap-6",
          isLocked && "opacity-35",
        )}
      >
        <div className="min-w-0">
          <div className="flex flex-wrap items-center gap-3">
            <h3 className="text-xl font-bold tracking-normal text-[#10111f]">
              {lesson.title}
            </h3>
            {lesson.status === "completed" ? (
              <CheckCircle2
                className="size-5 text-success"
                aria-hidden="true"
              />
            ) : null}
            {lesson.status === "in-progress" ? (
              <span className="rounded-full bg-[#eceaff] px-3 py-1 text-sm font-bold text-primary">
                In Progress
              </span>
            ) : null}
          </div>

          <p className="mt-2 text-base font-medium text-[#676982]">
            {lesson.description}
          </p>

          {lesson.status === "in-progress" ? (
            <div className="mt-5">
              <div className="h-2 overflow-hidden rounded-full bg-[#ebeaff]">
                <div
                  className="h-full rounded-full bg-primary"
                  style={{ width: `${lesson.progressPercentage}%` }}
                />
              </div>
              <p className="mt-2 text-base font-bold text-[#676982]">
                {lesson.progressPercentage}% complete
              </p>
            </div>
          ) : null}

          <div className="mt-4 flex flex-wrap items-center gap-x-4 gap-y-2 text-sm font-bold text-[#77798e]">
            <LessonMeta icon={BookOpen} label={`${lesson.wordCount} words`} />
            <LessonMeta
              icon={Clock3}
              label={`~${lesson.estimatedMinutes} min`}
            />
            <LessonMeta
              className="text-[#c46600]"
              icon={Star}
              label={`${lesson.xpReward} XP`}
            />
          </div>
        </div>

        <ArrowRight
          className="mt-2 size-5 shrink-0 text-[#77798e]"
          aria-hidden="true"
        />
      </div>

      {isLocked ? (
        <div className="absolute inset-0 grid place-items-center bg-white/55">
          <div className="text-center">
            <Lock
              className="mx-auto size-6 text-[#6f7184]"
              aria-hidden="true"
            />
            <p className="mt-2 text-sm font-bold text-[#6f7184]">
              {lesson.lockedReason}
            </p>
          </div>
        </div>
      ) : null}
    </article>
  );

  if (isLocked) {
    return content;
  }

  return (
    <Link aria-label={`Open ${lesson.title}`} to={`/lessons/${lesson.id}`}>
      {content}
    </Link>
  );
}

type LessonMetaProps = {
  icon: typeof BookOpen;
  label: string;
  className?: string;
};

function LessonMeta({ icon: Icon, label, className }: LessonMetaProps) {
  return (
    <span className={cn("inline-flex items-center gap-2", className)}>
      <Icon className="size-4" aria-hidden="true" />
      {label}
    </span>
  );
}

function getContinueLessonId(lessons: TopicLessonViewModel[]) {
  return (
    lessons.find((lesson) => lesson.status === "in-progress") ??
    lessons.find((lesson) => lesson.status === "unlocked") ??
    lessons[0]
  ).id;
}
