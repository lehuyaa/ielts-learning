export type VocabularyDifficulty = 'BEGINNER' | 'INTERMEDIATE' | 'ADVANCED'

export type VocabularyStatus = 'NEW' | 'LEARNING' | 'REVIEW' | 'MASTERED'

export type VocabularyQueryParams = {
  q?: string
  difficulty?: VocabularyDifficulty
  targetBand?: number
  status?: VocabularyStatus
  page?: number
  limit?: number
}

export type VocabularyTopic = {
  id: number
  title: string
  slug: string
  icon: string
  emoji: string
  color: string
}

export type VocabularyBandLevel = {
  id: number
  bandScore: number
  title: string
}

export type VocabularyProgress = {
  status: VocabularyStatus
  reviewCount: number
  nextReviewAt: string | null
  masteryScore: number
}

export type VocabularyDetailProgress = VocabularyProgress & {
  correctCount: number
  wrongCount: number
  lastReviewedAt: string | null
  learnedAt: string | null
}

export type VocabularyListItem = {
  id: number
  word: string
  slug: string
  ipa: string
  partOfSpeech: string
  meaningVi: string
  meaningEn: string
  shortDefinition: string
  difficulty: VocabularyDifficulty
  targetBand: number | null
  topic: VocabularyTopic | null
  progress: VocabularyProgress
  status: VocabularyStatus
  masteryScore: number
}

export type VocabularyPagination = {
  page: number
  limit: number
  total: number
  totalPages: number
}

export type VocabularyListResponse = {
  items: VocabularyListItem[]
  pagination: VocabularyPagination
}

export type RelatedForm = {
  word: string
  partOfSpeech: string
}

export type VocabularyDetailResponse = {
  id: number
  word: string
  slug: string
  ipa: string
  audioUrl: string
  partOfSpeech: string
  meaningVi: string
  meaningEn: string
  primaryMeaning: string
  secondaryMeaning: string
  exampleSentences: string[]
  synonyms: string[]
  antonyms: string[]
  collocations: string[]
  ieltsUsage: string
  relatedForms: RelatedForm[]
  difficulty: VocabularyDifficulty
  targetBand: number | null
  frequency: string
  rating: number
  topic: VocabularyTopic | null
  bandLevel: VocabularyBandLevel | null
  progress: VocabularyDetailProgress
  userProgress: VocabularyDetailProgress
  masteryScore: number
}

export type VocabularyDifficultyLabel =
  | 'Beginner'
  | 'Intermediate'
  | 'Advanced'

export type VocabularyStatusLabel = 'New' | 'Learning' | 'Review' | 'Mastered'

export type VocabularyExampleViewModel = {
  sentence: string
  note: string
}

export type VocabularyListItemViewModel = {
  id: string
  word: string
  ipa: string
  partOfSpeech: string
  topic: string
  band: string
  bandScore: number | null
  difficulty: VocabularyDifficultyLabel
  status: VocabularyStatusLabel
  frequency: string
  masteryScore: number
  shortDefinition: string
}

export type VocabularyDetailViewModel = VocabularyListItemViewModel & {
  slug: string
  frequencyScore: number
  reviewCount: number
  lastReviewedAt: string
  primaryMeaning: string
  secondaryMeaning: string
  meaningVi: string
  meaningEn: string
  examples: VocabularyExampleViewModel[]
  synonyms: string[]
  antonyms: string[]
  collocations: string[]
  ieltsUsage: string
  relatedForms: RelatedForm[]
}
