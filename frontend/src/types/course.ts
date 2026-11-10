export type CourseKind = 'IELTS' | 'CEFR'

export type CourseItem = {
  id: number
  title: string
  slug: string
  description: string
  kind: CourseKind
  bandMin: number
  bandMax: number
  totalWords: number
  totalLessons: number
  totalTopics: number
  completedLessons: number
  isActive: boolean
}

export type CourseListResponse = {
  activeCourseId: number
  items: CourseItem[]
}

export type SetActiveCourseResponse = {
  activeCourseId: number
}
