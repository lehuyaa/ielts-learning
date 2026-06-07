export type VocabularyProgressStatus =
  | "NEW"
  | "LEARNING"
  | "REVIEW"
  | "MASTERED";

export type VocabularyProgress = {
  status: VocabularyProgressStatus;
  reviewCount: number;
  correctCount?: number;
  wrongCount?: number;
  nextReviewAt: string | null;
  lastReviewedAt?: string | null;
  learnedAt?: string | null;
  masteryScore?: number;
};
