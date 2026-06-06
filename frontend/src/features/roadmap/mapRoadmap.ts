import type {
  RoadmapBandLevel,
  RoadmapResponse,
  RoadmapStatus,
  RoadmapTopic,
} from "@/types/roadmap";
import type {
  RoadmapBand,
  RoadmapLessonStatus,
  RoadmapViewModel,
} from "@/features/roadmap/types";
import { getTopicDisplayIcon } from "@/lib/topicIcons";

export function mapRoadmapToViewModel(
  roadmap: RoadmapResponse,
): RoadmapViewModel {
  return {
    title: roadmap.course.title || "Vocabulary Roadmap",
    subtitle: `Band ${roadmap.course.bandMin.toFixed(1)} → ${roadmap.course.bandMax.toFixed(1)}`,
    topicsCompleted: roadmap.summary.topicsCompleted,
    totalTopics: roadmap.summary.totalTopics,
    currentBand: roadmap.summary.currentBand,
    wordsMastered: roadmap.summary.wordsMastered,
    currentStreak: roadmap.summary.currentStreak,
    bands: roadmap.bandLevels.map(mapBand),
  };
}

function mapBand(band: RoadmapBandLevel): RoadmapBand {
  return {
    id: String(band.id),
    band: band.bandScore.toFixed(1),
    title: band.title,
    description: band.description,
    status: mapStatus(band.status),
    progress: band.progressPercentage,
    lessonCount: band.topics.reduce(
      (total, topic) => total + topic.totalLessons,
      0,
    ),
    topicCount: band.totalTopics,
    topics: band.topics.map(mapTopic),
  };
}

function mapTopic(topic: RoadmapTopic) {
  return {
    id: topic.id,
    title: topic.title,
    description: "",
    progress: topic.progressPercentage,
    icon: getTopicDisplayIcon(topic),
    completedLessons: topic.lessonsCompleted,
    totalLessons: topic.totalLessons,
    lessons: topic.lessons.map((lesson) => ({
      id: String(lesson.id),
      title: lesson.title,
      duration: `${lesson.estimatedMinutes} min`,
      wordCount: lesson.xpReward,
      status: mapStatus(lesson.status),
    })),
  };
}

function mapStatus(status: RoadmapStatus): RoadmapLessonStatus {
  switch (status) {
    case "COMPLETED":
      return "completed";
    case "IN_PROGRESS":
      return "in-progress";
    case "UNLOCKED":
      return "unlocked";
    case "LOCKED":
    default:
      return "locked";
  }
}
