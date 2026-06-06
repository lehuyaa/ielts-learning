export type TopicLessonStatus =
  | 'LOCKED'
  | 'UNLOCKED'
  | 'IN_PROGRESS'
  | 'COMPLETED'

export type Topic = {
  id: number
  title: string
  slug: string
  icon: string
  emoji: string
  color: string
  description: string
}

export type TopicBandLevel = {
  id: number
  bandScore: number
  title: string
}

export type TopicSummary = {
  progressPercentage: number
  completedLessons: number
  totalLessons: number
  totalXP: number
}

export type TopicLesson = {
  id: number
  title: string
  slug: string
  description: string
  wordCount: number
  estimatedMinutes: number
  xpReward: number
  status: TopicLessonStatus
  progressPercentage: number
  lockedReason: string | null
}

export type TopicDetailResponse = {
  topic: Topic
  bandLevel: TopicBandLevel
  summary: TopicSummary
  lessons: TopicLesson[]
}

export type TopicLessonViewStatus =
  | 'completed'
  | 'in-progress'
  | 'unlocked'
  | 'locked'

export type TopicLessonViewModel = {
  id: string
  title: string
  description: string
  wordCount: number
  estimatedMinutes: number
  xpReward: number
  status: TopicLessonViewStatus
  progressPercentage: number
  lockedReason: string | null
}

export type TopicDetailViewModel = {
  id: string
  title: string
  icon: string
  band: string
  description: string
  progressPercentage: number
  completedLessons: number
  totalLessons: number
  totalXP: number
  lessons: TopicLessonViewModel[]
}
