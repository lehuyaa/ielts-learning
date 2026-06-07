import type { VocabularyProgressStatus } from "@/types/progress";

export type FlashcardRating = "AGAIN" | "HARD" | "GOOD" | "EASY";

export type FlashcardDifficulty = "BEGINNER" | "INTERMEDIATE" | "ADVANCED";

export type FlashcardCard = {
  id: number;
  vocabularyId: number;
  word: string;
  slug: string;
  ipa: string;
  audioUrl: string;
  partOfSpeech: string;
  meaningVi: string;
  meaningEn: string;
  exampleSentence: string;
  synonyms: string[];
  collocations: string[];
  ieltsUsage: string;
  difficulty: FlashcardDifficulty;
  band: number | null;
  targetBand: number | null;
  topicTitle: string;
  status: VocabularyProgressStatus;
  reviewCount: number;
  nextReviewAt: string | null;
};

export type FlashcardLesson = {
  id: number;
  title: string;
  topicId: number;
  topicTitle: string;
  bandLabel: string;
};

export type FlashcardSessionProgress = {
  done: number;
  remaining: number;
  total: number;
};

export type LessonFlashcardsResponse = {
  lesson: FlashcardLesson;
  progress: FlashcardSessionProgress;
  cards: FlashcardCard[];
};

export type DueReviewsQueryParams = {
  limit?: number;
  topicId?: number;
};

export type DueReviewsResponse = {
  cards: FlashcardCard[];
  count: number;
};

export type ReviewFlashcardRequest = {
  vocabularyId: number;
  lessonId?: number;
  rating: FlashcardRating;
};

export type ReviewFlashcardResponse = {
  vocabularyId: number;
  status: VocabularyProgressStatus;
  reviewCount: number;
  correctCount: number;
  wrongCount: number;
  lastRating: FlashcardRating;
  lastReviewedAt: string | null;
  nextReviewAt: string | null;
  masteryScore: number;
  xpAwarded: number;
  totalXp: number;
};

export type FlashcardStatusLabel = "New" | "Learning" | "Review" | "Mastered";

export type FlashcardDifficultyLabel =
  | "Beginner"
  | "Intermediate"
  | "Advanced";

export type FlashcardSessionCard = {
  id: string;
  vocabularyId: number;
  word: string;
  ipa: string;
  partOfSpeech: string;
  meaningEn: string;
  meaningVi: string;
  exampleSentence: string;
  synonyms: string[];
  collocations: string[];
  ieltsUsage: string;
  difficulty: FlashcardDifficultyLabel;
  band: string;
  topicTitle: string;
  status: FlashcardStatusLabel;
};
