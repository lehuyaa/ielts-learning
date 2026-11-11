export type PhonemeKind = 'VOWEL' | 'DIPHTHONG' | 'CONSONANT'

export type PhonemeWord = { word: string; ipa: string; audioUrl?: string }

export type MinimalPair = {
  /** The word with this sound. */
  a: string
  aIpa: string
  aAudioUrl?: string
  /** The word with the contrasting sound (contrast is "" when the sound is dropped). */
  b: string
  bIpa: string
  bAudioUrl?: string
  contrast: string
}

export type PhonemeProgress = {
  accuracy: number
  attempts: number
  practised: boolean
  lastPracticedAt: string | null
}

export type PhonemeItem = {
  id: number
  symbol: string
  kind: PhonemeKind
  voiced: boolean | null
  nameVi: string
  audioUrl: string
  exampleWord: PhonemeWord | null
  orderIndex: number
  progress: PhonemeProgress
}

export type PhonemeCard = {
  id: number
  symbol: string
  kind: PhonemeKind
  voiced: boolean | null
  nameVi: string
  exampleWords: PhonemeWord[]
  minimalPairs: MinimalPair[]
  mouthVi: string
  mistakesVi: string
  audioUrl: string
  imageUrl: string | null
  orderIndex: number
}

export type PhonemeLessonLink = {
  lessonId: number
  lessonTitle: string
  topicTitle: string
  activityId: number
  activityTitle: string
}

export type PracticeTopic = { id: number; title: string; description: string; lessonCount: number }

export type PhonemeListResponse = { items: PhonemeItem[]; practiceTopics: PracticeTopic[] }

export type PhonemeDetailResponse = {
  phoneme: PhonemeCard
  progress: PhonemeProgress
  lessons: PhonemeLessonLink[]
  previous: string | null
  next: string | null
}
