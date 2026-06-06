export type RoadmapLessonStatus =
  | 'locked'
  | 'unlocked'
  | 'in-progress'
  | 'completed'

export type RoadmapLesson = {
  id: string
  title: string
  duration: string
  wordCount: number
  status: RoadmapLessonStatus
}

export type RoadmapTopic = {
  id: string
  title: string
  description: string
  progress: number
  icon: string
  completedLessons: number
  totalLessons: number
  lessons: RoadmapLesson[]
}

export type RoadmapBand = {
  id: string
  band: string
  title: string
  description: string
  status: RoadmapLessonStatus
  progress: number
  lessonCount: number
  topicCount: number
  topics: RoadmapTopic[]
}

export type RoadmapViewModel = {
  title: string
  subtitle: string
  topicsCompleted: number
  totalTopics: number
  currentBand: number | null
  wordsMastered: number
  currentStreak: number
  bands: RoadmapBand[]
}
