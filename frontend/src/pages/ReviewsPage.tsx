import { FlashcardSession } from "@/features/flashcard/FlashcardSession";
import { reviewFlashcards } from "@/features/flashcard/mockFlashcards";

export function ReviewsPage() {
  return (
    <FlashcardSession
      cards={reviewFlashcards}
      completionPrimaryHref="/dashboard"
      completionPrimaryLabel="Back to dashboard"
      emptyDescription="You have no due vocabulary cards right now. Your review queue is clear."
      emptyTitle="No reviews due"
      subtitle="Due review queue"
      title="Review Session"
    />
  );
}
