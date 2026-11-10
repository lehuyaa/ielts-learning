export type WritingContent = {
  promptVi: string
  promptEn?: string
  frame?: string[]
  minSentences?: number
  maxSentences?: number
  targetGrammar?: string[]
  targetVocabulary?: string[]
}

export type WritingMistake = {
  /** UTF-16 offsets in the graded text; -1 when not found. */
  start: number
  end: number
  wrong: string
  correct: string
  skillTag: string
  explanationVi: string
}

export type WritingReport = {
  scores: { taskAchievement: number; grammar: number; vocabulary: number; coherence: number }
  overall: number
  correctedText: string
  mistakes: WritingMistake[]
  usedTargets: { grammar: string[]; vocabulary: string[] }
  feedbackVi: string
  rewriteRequired: boolean
}

export type WritingSubmission = {
  id: number
  activityId: number
  text: string
  status: 'GRADED' | 'REWRITE_REQUIRED' | 'COMPLETED'
  score: number
  passed: boolean
  passScore: number
  report: WritingReport
  rewriteText: string | null
  rewriteCount: number
  rewriteReport: WritingReport | null
  modelAnswer: string
  completedAt: string | null
  createdAt: string
}
