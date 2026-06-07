import { AlertCircle } from "lucide-react";
import { useParams } from "react-router-dom";

import { APIError } from "@/api/api";
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
          <TopicStateMessage
            title="Loading topic"
            description="Preparing the lesson list and your current progress."
          />
        ) : null}

        {errorMessage ? (
          <TopicStateMessage
            title="Topic unavailable"
            description={errorMessage}
            tone="error"
          />
        ) : null}

        {isEmpty ? (
          <>
            <TopicHero topic={topic} />
            <TopicStateMessage
              title="No lessons yet"
              description="This topic exists, but no lessons have been published for it."
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

type TopicStateMessageProps = {
  title: string;
  description: string;
  tone?: "default" | "error";
};

function TopicStateMessage({
  title,
  description,
  tone = "default",
}: TopicStateMessageProps) {
  return (
    <section className="rounded-2xl border border-[#e6e6f3] bg-white p-6 text-center shadow-sm">
      {tone === "error" ? (
        <AlertCircle
          className="mx-auto size-8 text-destructive"
          aria-hidden="true"
        />
      ) : null}
      <h2 className="text-xl font-bold tracking-normal text-[#10111f]">
        {title}
      </h2>
      <p className="mt-3 text-base font-medium text-[#676982]">{description}</p>
    </section>
  );
}
