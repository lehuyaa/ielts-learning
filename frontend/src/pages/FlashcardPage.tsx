import { useParams } from "react-router-dom";

import { APIError } from "@/api/api";
import { FlashcardSession } from "@/features/flashcard/FlashcardSession";
import { useLessonFlashcards } from "@/features/flashcard/hooks/useLessonFlashcards";
import { useReviewFlashcard } from "@/features/flashcard/hooks/useReviewFlashcard";
import { mapFlashcardCardToSessionCard } from "@/features/flashcard/mapFlashcard";
import type { FlashcardRating, FlashcardSessionCard } from "@/types/flashcard";

export function FlashcardPage() {
  const { lessonId } = useParams();
  const flashcardsQuery = useLessonFlashcards(lessonId);
  const reviewMutation = useReviewFlashcard({
    lessonId,
    topicId: flashcardsQuery.data?.lesson.topicId,
  });
  const errorMessage = getFlashcardErrorMessage(flashcardsQuery.error);

  const cards =
    flashcardsQuery.data?.cards.map(mapFlashcardCardToSessionCard) ?? [];
  const lessonTitle = flashcardsQuery.data?.lesson.title ?? "Lesson";
  const completionHref = lessonId ? `/lessons/${lessonId}` : "/roadmap";

  async function rateCard(card: FlashcardSessionCard, rating: FlashcardRating) {
    await reviewMutation.mutateAsync({
      vocabularyId: card.vocabularyId,
      lessonId: parseOptionalNumber(lessonId),
      rating,
    });
  }

  if (flashcardsQuery.isLoading) {
    return (
      <FlashcardSession
        cards={[]}
        completionPrimaryHref={completionHref}
        completionPrimaryLabel="Back to lesson"
        emptyActionHref={completionHref}
        emptyActionLabel="Back to lesson"
        emptyDescription="Preparing this lesson's vocabulary cards."
        emptyTitle="Loading flashcards"
        subtitle="Lesson flashcards"
        title="Flashcard Learning"
      />
    );
  }

  if (errorMessage) {
    return (
      <FlashcardSession
        cards={[]}
        completionPrimaryHref={completionHref}
        completionPrimaryLabel="Back to lesson"
        emptyActionHref={completionHref}
        emptyActionLabel="Back to lesson"
        emptyDescription={errorMessage}
        emptyTitle="Flashcards unavailable"
        subtitle="Lesson flashcards"
        title="Flashcard Learning"
      />
    );
  }

  return (
    <FlashcardSession
      cards={cards}
      completionPrimaryHref={completionHref}
      completionPrimaryLabel="Back to lesson"
      emptyActionHref={completionHref}
      emptyActionLabel="Back to lesson"
      emptyDescription="This lesson does not have flashcards yet."
      emptyTitle="No flashcards yet"
      onRateCard={rateCard}
      subtitle={lessonTitle}
      title="Flashcard Learning"
    />
  );
}

function parseOptionalNumber(value: string | undefined) {
  if (!value) {
    return undefined;
  }

  const parsed = Number(value);
  return Number.isFinite(parsed) ? parsed : undefined;
}

function getFlashcardErrorMessage(error: Error | null) {
  if (!error) {
    return null;
  }

  if (error instanceof APIError) {
    if (error.status === 404) {
      return "This lesson could not be found.";
    }

    if (error.status === 403) {
      return "Complete previous lessons to unlock these flashcards.";
    }
  }

  return "Unable to load flashcards right now.";
}
