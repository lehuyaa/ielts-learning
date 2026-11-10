import {
  ArrowLeft,
  BookOpen,
  CheckCircle2,
  Ear,
  FileText,
  Flag,
  Lock,
  Mic,
  PenLine,
  Puzzle,
  Volume2,
  type LucideIcon,
} from "lucide-react";
import { useState } from "react";
import { Link, useParams, useSearchParams } from "react-router-dom";

import { APIError } from "@/api/api";
import { ErrorState } from "@/components/state/ErrorState";
import { Skeleton } from "@/components/ui/skeleton";
import { ActivityRunner } from "@/features/activity/ActivityRunner";
import { useLessonActivities } from "@/features/activity/hooks/useLessonActivities";
import { cn } from "@/lib/utils";
import type { ActivityItem, ActivityType } from "@/types/activity";

const activityIcons: Record<ActivityType, LucideIcon> = {
  PRONUNCIATION: Volume2,
  VOCABULARY: BookOpen,
  GRAMMAR: Puzzle,
  LISTENING: Ear,
  READING: FileText,
  WRITING: PenLine,
  SPEAKING: Mic,
  CHECKPOINT: Flag,
};

const activityLabels: Record<ActivityType, string> = {
  PRONUNCIATION: "Phát âm",
  VOCABULARY: "Từ vựng",
  GRAMMAR: "Ngữ pháp",
  LISTENING: "Nghe",
  READING: "Đọc",
  WRITING: "Viết",
  SPEAKING: "Nói",
  CHECKPOINT: "Kiểm tra",
};

export function LessonPlayerPage() {
  const { lessonId } = useParams();
  const [searchParams, setSearchParams] = useSearchParams();
  const activitiesQuery = useLessonActivities(lessonId);
  const [selectedId, setSelectedId] = useState<number | null>(() => {
    const fromUrl = Number(searchParams.get("activity"));
    return Number.isFinite(fromUrl) && fromUrl > 0 ? fromUrl : null;
  });

  function select(activityId: number) {
    setSelectedId(activityId);
    setSearchParams({ activity: String(activityId) }, { replace: true });
    window.scrollTo({ top: 0, behavior: "smooth" });
  }

  if (activitiesQuery.isLoading) {
    return (
      <PlayerShell lessonId={lessonId} title="">
        <Skeleton className="mb-6 h-14 w-full rounded-2xl" />
        <Skeleton className="h-64 w-full rounded-3xl" />
      </PlayerShell>
    );
  }

  if (activitiesQuery.error || !activitiesQuery.data) {
    return (
      <PlayerShell lessonId={lessonId} title="">
        <ErrorState
          actionHref="/roadmap"
          actionLabel="Về lộ trình"
          description={
            activitiesQuery.error instanceof APIError
              ? activitiesQuery.error.message
              : "Không tải được bài học."
          }
          onRetry={() => void activitiesQuery.refetch()}
          title="Không mở được bài học"
        />
      </PlayerShell>
    );
  }

  const { lesson, activities, progress } = activitiesQuery.data;
  const selectable = activities.filter((activity) => !activity.isLocked);
  const resumeAt =
    selectable.find((activity) => activity.id === progress.nextActivityId) ??
    selectable[0];

  // Pin the activity being studied: progress refetches after finishing an
  // activity must not swap it out from under the result screen.
  if (selectedId === null && resumeAt) {
    setSelectedId(resumeAt.id);
  }
  const current =
    selectable.find((activity) => activity.id === selectedId) ?? resumeAt;

  return (
    <PlayerShell lessonId={lessonId} title={lesson.title}>
      <ActivityStepper
        activities={activities}
        currentId={current?.id}
        onSelect={select}
        passed={progress.passedRequired}
        total={progress.totalRequired}
      />

      {activities.length === 0 ? (
        <p className="rounded-2xl border border-dashed p-8 text-center text-muted-foreground">
          Bài học này chưa có hoạt động nào.
        </p>
      ) : current ? (
        <ActivityRunner
          activityId={current.id}
          key={current.id}
          lessonId={lesson.id}
          onNextActivity={select}
        />
      ) : null}
    </PlayerShell>
  );
}

function PlayerShell({
  title,
  children,
}: {
  lessonId?: string;
  title: string;
  children: React.ReactNode;
}) {
  return (
    <div className="min-h-screen bg-[#f8f8ff] text-foreground">
      <header className="sticky top-0 z-20 border-b bg-white/90 backdrop-blur">
        <div className="mx-auto flex max-w-3xl items-center gap-3 px-4 py-3">
          <Link
            aria-label="Quay lại"
            className="grid size-10 place-items-center rounded-full hover:bg-muted"
            // Not the lesson page: it redirects skill lessons back here.
            to="/roadmap"
          >
            <ArrowLeft className="size-5" />
          </Link>
          <h1 className="truncate text-base font-bold tracking-normal">{title}</h1>
        </div>
      </header>
      <main className="mx-auto max-w-3xl px-4 py-6 sm:py-8">{children}</main>
    </div>
  );
}

function ActivityStepper({
  activities,
  currentId,
  onSelect,
  passed,
  total,
}: {
  activities: ActivityItem[];
  currentId: number | undefined;
  onSelect: (activityId: number) => void;
  passed: number;
  total: number;
}) {
  return (
    <nav aria-label="Các hoạt động trong bài" className="mb-6">
      <div className="mb-2 flex items-center justify-between text-xs text-muted-foreground">
        <span>Tiến độ bài học</span>
        <span>
          {passed}/{total} hoạt động bắt buộc
        </span>
      </div>
      <ol className="flex gap-2 overflow-x-auto pb-1">
        {activities.map((activity) => {
          const Icon = activity.isLocked ? Lock : activityIcons[activity.type] ?? Flag;
          const isPassed = activity.status === "PASSED";
          const isCurrent = activity.id === currentId;
          return (
            <li className="min-w-[88px] flex-1" key={activity.id}>
              <button
                aria-current={isCurrent ? "step" : undefined}
                className={cn(
                  "flex w-full cursor-pointer flex-col items-center gap-1 rounded-2xl border-2 px-2 py-2.5 text-xs font-semibold transition-colors disabled:cursor-not-allowed",
                  isCurrent && "border-primary bg-primary/5 text-primary",
                  !isCurrent && isPassed && "border-emerald-200 bg-emerald-50 text-emerald-700",
                  !isCurrent && !isPassed && "border-border bg-white text-muted-foreground",
                  activity.isLocked && "opacity-50",
                )}
                disabled={activity.isLocked}
                onClick={() => onSelect(activity.id)}
                title={activity.title}
                type="button"
              >
                <span className="relative">
                  <Icon className="size-5" aria-hidden="true" />
                  {isPassed ? (
                    <CheckCircle2 className="absolute -right-2 -top-2 size-3.5 fill-white text-emerald-600" />
                  ) : null}
                </span>
                {activityLabels[activity.type] ?? activity.type}
                {!activity.isRequired ? (
                  <span className="text-[10px] font-normal">tùy chọn</span>
                ) : null}
              </button>
            </li>
          );
        })}
      </ol>
    </nav>
  );
}
