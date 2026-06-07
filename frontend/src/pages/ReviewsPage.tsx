import { APIError } from "@/api/api";
import { FlashcardSession } from "@/features/flashcard/FlashcardSession";
import { useDueReviews } from "@/features/flashcard/hooks/useDueReviews";
import { useReviewFlashcard } from "@/features/flashcard/hooks/useReviewFlashcard";
import { mapFlashcardCardToSessionCard } from "@/features/flashcard/mapFlashcard";
import type { FlashcardRating, FlashcardSessionCard } from "@/types/flashcard";

export function ReviewsPage() {
  const dueReviewsQuery = useDueReviews();
  const reviewMutation = useReviewFlashcard();
  const errorMessage = getReviewsErrorMessage(dueReviewsQuery.error);
  const cards =
    dueReviewsQuery.data?.cards.map(mapFlashcardCardToSessionCard) ?? [];

  async function rateCard(
    card: FlashcardSessionCard,
    rating: FlashcardRating,
  ) {
    await reviewMutation.mutateAsync({
      vocabularyId: card.vocabularyId,
      rating,
    });
  }

  if (dueReviewsQuery.isLoading) {
    return (
      <FlashcardSession
        cards={[]}
        completionPrimaryHref="/dashboard"
        completionPrimaryLabel="Back to dashboard"
        emptyActionHref="/dashboard"
        emptyActionLabel="Back to dashboard"
        emptyDescription="Checking your spaced repetition queue."
        emptyTitle="Loading reviews"
        subtitle="Due review queue"
        title="Review Session"
      />
    );
  }

  if (errorMessage) {
    return (
      <FlashcardSession
        cards={[]}
        completionPrimaryHref="/dashboard"
        completionPrimaryLabel="Back to dashboard"
        emptyActionHref="/dashboard"
        emptyActionLabel="Back to dashboard"
        emptyDescription={errorMessage}
        emptyTitle="Reviews unavailable"
        subtitle="Due review queue"
        title="Review Session"
      />
    );
  }

  return (
    <FlashcardSession
      cards={cards}
      completionPrimaryHref="/dashboard"
      completionPrimaryLabel="Back to dashboard"
      emptyActionHref="/dashboard"
      emptyActionLabel="Back to dashboard"
      emptyDescription="You have no due vocabulary cards right now. Your review queue is clear."
      emptyTitle="No reviews due"
      onRateCard={rateCard}
      subtitle="Due review queue"
      title="Review Session"
    />
  );
}

function getReviewsErrorMessage(error: Error | null) {
  if (!error) {
    return null;
  }

  if (error instanceof APIError && error.status === 401) {
    return "Please log in again to load your review queue.";
  }

  return "Unable to load reviews right now.";
}
