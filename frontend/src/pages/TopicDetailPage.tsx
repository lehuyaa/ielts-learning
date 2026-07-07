import { useParams } from "react-router-dom";

import { APIError } from "@/api/api";
import { CardSkeleton } from "@/components/state/CardSkeleton";
import { EmptyState } from "@/components/state/EmptyState";
import { ErrorState } from "@/components/state/ErrorState";
import { Skeleton } from "@/components/ui/skeleton";
import { LessonList } from "@/features/topic/LessonList";
import { TopicHeader } from "@/features/topic/TopicHeader";
import { TopicHero } from "@/features/topic/TopicHero";
import { useTopicDetail } from "@/features/topic/hooks/useTopicDetail";
import { getTopicDisplayIcon } from "@/lib/topicIcons";
import type {
  TopicDetailResponse,
  TopicDetailViewModel,
  TopicLessonStatus,
} from "@/types/topic";

const loadingTopic: TopicDetailViewModel = {
  id: "loading",
  title: "Topic",
  icon: "📘",
  band: "Loading",
  description: "Loading topic lessons and progress.",
  progressPercentage: 0,
  completedLessons: 0,
  totalLessons: 0,
  totalXP: 0,
  lessons: [],
};

export function TopicDetailPage() {
  const { topicId } = useParams();
  const topicQuery = useTopicDetail(topicId);
  const topic = topicQuery.data
    ? mapTopicDetailToViewModel(topicQuery.data)
    : loadingTopic;
  const errorMessage = getTopicErrorMessage(topicQuery.error);

  const isEmpty =
    !topicQuery.isLoading &&
    !errorMessage &&
    Boolean(topicQuery.data) &&
    topic.lessons.length === 0;

  return (
    <div className="min-h-screen bg-[#f8f8ff] text-[#10111f]">
      <TopicHeader topic={topic} />

      <main className="mx-auto grid max-w-[1200px] gap-6 px-4 py-8 md:px-6">
        {topicQuery.isLoading ? (
          <TopicLoadingSkeleton />
        ) : null}

        {errorMessage ? (
          <ErrorState
            actionHref="/roadmap"
            actionLabel="Back to roadmap"
            description={errorMessage}
            onRetry={() => {
              void topicQuery.refetch();
            }}
            title="Topic unavailable"
          />
        ) : null}

        {isEmpty ? (
          <>
            <TopicHero topic={topic} />
            <EmptyState
              actionHref="/roadmap"
              actionLabel="Back to roadmap"
              description="This topic exists, but no lessons have been published for it."
              title="No lessons yet"
            />
          </>
        ) : null}

        {!topicQuery.isLoading && !errorMessage && !isEmpty ? (
          <>
            <TopicHero topic={topic} />
            <LessonList topic={topic} />
          </>
        ) : null}
      </main>
    </div>
  );
}

function mapTopicDetailToViewModel(
  response: TopicDetailResponse,
): TopicDetailViewModel {
  return {
    id: String(response.topic.id),
    title: response.topic.title,
    icon: getTopicDisplayIcon(response.topic),
    band: response.bandLevel.title,
    description: response.topic.description,
    progressPercentage: response.summary.progressPercentage,
    completedLessons: response.summary.completedLessons,
    totalLessons: response.summary.totalLessons,
    totalXP: response.summary.totalXP,
    lessons: response.lessons.map((lesson) => ({
      id: String(lesson.id),
      title: lesson.title,
      description: lesson.description,
      wordCount: lesson.wordCount,
      estimatedMinutes: lesson.estimatedMinutes,
      xpReward: lesson.xpReward,
      status: mapLessonStatus(lesson.status),
      progressPercentage: lesson.progressPercentage,
      lockedReason: lesson.lockedReason,
    })),
  };
}

function mapLessonStatus(status: TopicLessonStatus) {
  switch (status) {
    case "COMPLETED":
      return "completed";
    case "IN_PROGRESS":
      return "in-progress";
    case "LOCKED":
      return "locked";
    case "UNLOCKED":
    default:
      return "unlocked";
  }
}

function getTopicErrorMessage(error: Error | null) {
  if (!error) {
    return null;
  }

  if (error instanceof APIError) {
    if (error.status === 404) {
      return "This topic could not be found.";
    }

    return error.message;
  }

  return "Unable to load this topic right now.";
}

function TopicLoadingSkeleton() {
  return (
    <>
      <section className="rounded-[28px] border border-[#e6e6f3] bg-white p-6 shadow-sm">
        <div className="flex flex-col gap-5 lg:flex-row lg:items-start lg:justify-between">
          <div className="space-y-3">
            <div className="flex items-center gap-3">
              <Skeleton className="size-12 rounded-2xl" />
              <Skeleton className="h-7 w-40" />
            </div>
            <Skeleton className="h-4 w-28" />
            <Skeleton className="h-4 w-[28rem] max-w-full" />
            <Skeleton className="h-4 w-[24rem] max-w-full" />
          </div>
          <div className="grid w-full max-w-xs gap-3">
            <Skeleton className="h-20 rounded-2xl" />
            <Skeleton className="h-20 rounded-2xl" />
          </div>
        </div>
      </section>

      <div className="grid gap-4">
        {Array.from({ length: 4 }).map((_, index) => (
          <CardSkeleton className="border-[#e6e6f3]" key={index} lines={3} />
        ))}
      </div>
    </>
  );
}
