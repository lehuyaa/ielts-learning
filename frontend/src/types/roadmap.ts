export type RoadmapStatus =
  | 'LOCKED'
  | 'UNLOCKED'
  | 'IN_PROGRESS'
  | 'COMPLETED'

export type RoadmapCourse = {
  id: number
  title: string
  slug: string
  bandMin: number
  bandMax: number
  totalWords: number
  totalLessons: number
  totalTopics: number
}

export type RoadmapSummary = {
  topicsCompleted: number
  totalTopics: number
  currentBand: number | null
  wordsMastered: number
  currentStreak: number
}

export type RoadmapLesson = {
  id: number
  title: string
  slug: string
  status: RoadmapStatus
  requiredScore: number
  estimatedMinutes: number
  xpReward: number
  score: number | null
  bestScore: number | null
}

export type RoadmapTopic = {
  id: number
  title: string
  slug: string
  icon?: string
  emoji: string
  color: string
  status: RoadmapStatus
  lessonsCompleted: number
  totalLessons: number
  progressPercentage: number
  lessons: RoadmapLesson[]
}

export type RoadmapBandLevel = {
  id: number
  bandScore: number
  title: string
  description: string
  status: RoadmapStatus
  progressPercentage: number
  topicsCompleted: number
  totalTopics: number
  topics: RoadmapTopic[]
}

export type RoadmapResponse = {
  course: RoadmapCourse
  summary: RoadmapSummary
  bandLevels: RoadmapBandLevel[]
}
