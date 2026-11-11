export type MistakeStatus = 'OPEN' | 'REVIEWING' | 'FIXED'
export type MistakeSource = 'EXERCISE' | 'WRITING' | 'SPEAKING' | 'PRONUNCIATION'
export type SkillKind = 'GRAMMAR' | 'SOUND' | 'VOCABULARY' | 'LISTENING' | 'READING' | 'SPEAKING' | 'OTHER'

export type MistakeItem = {
  id: number
  source: MistakeSource
  skillTag: string
  label: string
  wrongText: string
  correctText: string
  explanationVi: string
  occurrences: number
  status: MistakeStatus
  reviewStreak: number
  nextReviewAt: string | null
  fixedAt: string | null
  lastSeenAt: string
}

export type MistakeGroup = {
  skillTag: string
  label: string
  kind: SkillKind
  practiceable: boolean
  count: number
  occurrences: number
  lastSeenAt: string
  items: MistakeItem[]
}

export type MistakeList = {
  summary: { open: number; reviewing: number; fixed: number; due: number }
  groups: MistakeGroup[]
}

export type MistakeReviewQueue = { cards: MistakeItem[]; count: number }

export type MistakeRating = 'AGAIN' | 'HARD' | 'GOOD' | 'EASY'

export type MistakeReviewResult = { item: MistakeItem; fixed: boolean }
