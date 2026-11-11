import type { AnswerResult, Exercise } from '@/types/activity'

export type PracticeSet = {
  id: number
  skillTag: string
  label: string
  /** AI = generated for you, BANK = picked from lessons. */
  source: 'AI' | 'BANK'
  day: string
  exercises: Exercise[]
  answers: AnswerResult[]
}
