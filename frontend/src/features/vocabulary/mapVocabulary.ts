import type {
  VocabularyDetailResponse,
  VocabularyDetailViewModel,
  VocabularyDifficulty,
  VocabularyDifficultyLabel,
  VocabularyListItem,
  VocabularyListItemViewModel,
  VocabularyStatus,
  VocabularyStatusLabel,
} from '@/types/vocabulary'

export function mapVocabularyListItem(
  item: VocabularyListItem,
): VocabularyListItemViewModel {
  return {
    id: String(item.id),
    word: item.word,
    ipa: item.ipa,
    partOfSpeech: item.partOfSpeech,
    topic: item.topic?.title ?? 'General',
    band: formatBand(item.targetBand),
    bandScore: item.targetBand,
    difficulty: mapDifficulty(item.difficulty),
    status: mapStatus(item.progress.status),
    frequency: frequencyFromBand(item.targetBand),
    masteryScore: item.progress.masteryScore,
    shortDefinition:
      item.shortDefinition || item.meaningEn || item.meaningVi || 'No meaning yet',
  }
}

export function mapVocabularyDetail(
  response: VocabularyDetailResponse,
): VocabularyDetailViewModel {
  const masteryScore = response.masteryScore ?? response.userProgress.masteryScore
  const examples = response.exampleSentences.length
    ? response.exampleSentences.map((sentence) => ({
        sentence,
        note:
          response.ieltsUsage ||
          'Useful for IELTS vocabulary review and writing practice.',
      }))
    : [
        {
          sentence:
            response.primaryMeaning ||
            response.meaningEn ||
            'Example content has not been added yet.',
          note: 'Example content will appear here when available.',
        },
      ]

  return {
    id: String(response.id),
    word: response.word,
    slug: response.slug,
    ipa: response.ipa,
    partOfSpeech: response.partOfSpeech,
    topic: response.topic?.title ?? 'General',
    band: formatBand(response.targetBand ?? response.bandLevel?.bandScore ?? null),
    bandScore: response.targetBand ?? response.bandLevel?.bandScore ?? null,
    difficulty: mapDifficulty(response.difficulty),
    status: mapStatus(response.userProgress.status),
    frequency: response.frequency || frequencyFromBand(response.targetBand),
    frequencyScore: ratingToScore(response.rating),
    masteryScore,
    reviewCount: response.userProgress.reviewCount,
    lastReviewedAt: formatRelativeDate(response.userProgress.lastReviewedAt),
    primaryMeaning:
      response.primaryMeaning ||
      response.meaningEn ||
      response.meaningVi ||
      'No primary meaning has been added yet.',
    secondaryMeaning:
      response.secondaryMeaning ||
      response.meaningVi ||
      'No secondary meaning has been added yet.',
    meaningVi: response.meaningVi,
    meaningEn: response.meaningEn,
    shortDefinition:
      response.primaryMeaning ||
      response.meaningEn ||
      response.meaningVi ||
      'No meaning yet',
    examples,
    synonyms: response.synonyms,
    antonyms: response.antonyms,
    collocations: response.collocations,
    ieltsUsage:
      response.ieltsUsage ||
      'IELTS usage guidance has not been added for this word yet.',
    relatedForms: response.relatedForms,
  }
}

export function mapDifficulty(
  difficulty: VocabularyDifficulty,
): VocabularyDifficultyLabel {
  switch (difficulty) {
    case 'BEGINNER':
      return 'Beginner'
    case 'ADVANCED':
      return 'Advanced'
    case 'INTERMEDIATE':
    default:
      return 'Intermediate'
  }
}

export function mapStatus(status: VocabularyStatus): VocabularyStatusLabel {
  switch (status) {
    case 'LEARNING':
      return 'Learning'
    case 'REVIEW':
      return 'Review'
    case 'MASTERED':
      return 'Mastered'
    case 'NEW':
    default:
      return 'New'
  }
}

export function formatBand(band: number | null | undefined) {
  if (band === null || band === undefined) {
    return 'Band -'
  }

  return Number.isInteger(band) ? `Band ${band.toFixed(1)}` : `Band ${band}`
}

function frequencyFromBand(band: number | null | undefined) {
  if (!band) {
    return 'Medium'
  }
  if (band >= 7) {
    return 'Very High'
  }
  if (band >= 6) {
    return 'High'
  }
  return 'Medium'
}

function ratingToScore(rating: number) {
  if (rating <= 0) {
    return 20
  }

  return Math.min(100, Math.max(20, rating * 20))
}

function formatRelativeDate(value: string | null) {
  if (!value) {
    return 'not reviewed'
  }

  const date = new Date(value)
  if (Number.isNaN(date.getTime())) {
    return 'recently'
  }

  const diffMs = Date.now() - date.getTime()
  const diffDays = Math.floor(diffMs / 86_400_000)

  if (diffDays <= 0) {
    return 'today'
  }
  if (diffDays === 1) {
    return 'yesterday'
  }
  if (diffDays < 30) {
    return `${diffDays} days ago`
  }

  return date.toLocaleDateString()
}
