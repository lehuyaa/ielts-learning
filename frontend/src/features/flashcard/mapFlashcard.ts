import type {
  FlashcardCard,
  FlashcardDifficulty,
  FlashcardDifficultyLabel,
  FlashcardSessionCard,
} from "@/types/flashcard";
import type { VocabularyProgressStatus } from "@/types/progress";

export function mapFlashcardCardToSessionCard(
  card: FlashcardCard,
): FlashcardSessionCard {
  return {
    id: String(card.id),
    vocabularyId: card.vocabularyId,
    word: card.word,
    ipa: card.ipa,
    partOfSpeech: card.partOfSpeech,
    meaningEn: card.meaningEn,
    meaningVi: card.meaningVi,
    exampleSentence: card.exampleSentence,
    synonyms: card.synonyms,
    collocations: card.collocations,
    ieltsUsage: card.ieltsUsage,
    difficulty: mapDifficulty(card.difficulty),
    band: formatBandLabel(card.band ?? card.targetBand),
    topicTitle: card.topicTitle,
    status: mapStatus(card.status),
  };
}

function mapDifficulty(
  difficulty: FlashcardDifficulty,
): FlashcardDifficultyLabel {
  switch (difficulty) {
    case "BEGINNER":
      return "Beginner";
    case "ADVANCED":
      return "Advanced";
    case "INTERMEDIATE":
    default:
      return "Intermediate";
  }
}

function mapStatus(status: VocabularyProgressStatus) {
  switch (status) {
    case "LEARNING":
      return "Learning";
    case "REVIEW":
      return "Review";
    case "MASTERED":
      return "Mastered";
    case "NEW":
    default:
      return "New";
  }
}

function formatBandLabel(band: number | null) {
  if (band === null) {
    return "Band";
  }

  return `Band ${Number.isInteger(band) ? band.toFixed(0) : band.toFixed(1)}`;
}
