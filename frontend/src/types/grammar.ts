import type { Exercise } from '@/types/activity'

export type GrammarExample = { en: string; vi: string; audioUrl?: string }
export type GrammarMistake = { wrong: string; right: string; noteVi: string }

export type GrammarPointCard = {
  id: number
  slug: string
  cefrLevel: string
  title: string
  titleVi: string
  summaryVi: string
  formula: string
  examples: GrammarExample[]
  commonMistakes: GrammarMistake[]
}

export type GrammarProgress = {
  /** Average score (0-100) of the last 20 answers. */
  mastery: number
  correctCount: number
  wrongCount: number
  practised: boolean
  lastPracticedAt: string | null
}

export type GrammarPointItem = {
  id: number
  slug: string
  cefrLevel: string
  title: string
  titleVi: string
  orderIndex: number
  exerciseCount: number
  progress: GrammarProgress
}

export type GrammarListResponse = { items: GrammarPointItem[] }

export type GrammarDetailResponse = {
  grammarPoint: GrammarPointCard
  progress: GrammarProgress
  exerciseCount: number
}

export type PracticeExercise = Exercise & {
  lastResult: 'NEW' | 'WRONG' | 'CORRECT'
}

export type GrammarPracticeResponse = { exercises: PracticeExercise[] }

export type PracticeAnswerResult = {
  exerciseId: number
  attempt: number
  isCorrect: boolean
  scoreRatio: number
  answer: unknown
  correctAnswer: unknown
  explanationVi: string
  answeredAt: string
  progress: GrammarProgress
}
