export type LessonStatus =
  | 'LOCKED'
  | 'UNLOCKED'
  | 'IN_PROGRESS'
  | 'COMPLETED'

export type LessonVocabularyStatus =
  | 'NEW'
  | 'LEARNING'
  | 'REVIEW'
  | 'MASTERED'

export type LessonDifficulty = 'BEGINNER' | 'INTERMEDIATE' | 'ADVANCED'

export type LessonDetail = {
  id: number
  title: string
  slug: string
  description: string
  requiredScore: number
  estimatedMinutes: number
  xpReward: number
  wordCount: number
  orderIndex: number
  status: LessonStatus
  progressPercentage: number
  lockedReason: string | null
}

export type LessonTopic = {
  id: number
  title: string
  slug: string
  icon: string
  emoji: string
  color: string
}

export type LessonBandLevel = {
  id: number
  bandScore: number
  title: string
}

export type LessonProgress = {
  status: LessonStatus
  score: number | null
  bestScore: number | null
  bestXp: number
  wordsLearned: number
  totalWords: number
  progressPercentage: number
  startedAt: string | null
  completedAt: string | null
  lastStudiedAt: string | null
}

export type LessonVocabulary = {
  id: number
  word: string
  slug: string
  ipa: string
  audioUrl: string
  partOfSpeech: string
  meaningVi: string
  meaningEn: string
  shortDefinition: string
  exampleSentence: string
  difficulty: LessonDifficulty
  targetBand: number | null
  status: LessonVocabularyStatus
  reviewCount: number
  correctCount: number
  wrongCount: number
  learned: boolean
  learnedAt: string | null
  lastReviewedAt: string | null
  nextReviewAt: string | null
}

export type LessonDetailResponse = {
  lesson: LessonDetail
  topic: LessonTopic
  bandLevel: LessonBandLevel
  progress: LessonProgress
  vocabularies: LessonVocabulary[]
}

export type StartLessonResponse = {
  lessonId: number
  status: LessonStatus
  startedAt: string | null
  lastStudiedAt: string | null
}

export type LessonVocabularyItem = {
  id: string
  word: string
  ipa: string
  partOfSpeech: string
  shortDefinition: string
  definition: string
  example: string
  difficulty: 'Beginner' | 'Intermediate' | 'Advanced'
  band: string
  learned: boolean
}

export type LessonDetailViewModel = {
  id: string
  title: string
  topic: string
  topicId: number
  bandRange: string
  description: string
  estimatedMinutes: number
  requiredScore: number
  xpReward: number
  lessonScore: number
  personalBest: number
  status: LessonStatus
  lockedReason: string | null
  progressPercentage: number
  vocabulary: LessonVocabularyItem[]
}
