import { useParams } from "react-router-dom";

import { FlashcardSession } from "@/features/flashcard/FlashcardSession";
import { lessonFlashcards } from "@/features/flashcard/mockFlashcards";

export function FlashcardPage() {
  const { lessonId = "1" } = useParams();

  return (
    <FlashcardSession
      cards={lessonFlashcards}
      completionPrimaryHref={`/lessons/${lessonId}`}
      completionPrimaryLabel="Back to lesson"
      subtitle={`Lesson ${lessonId}`}
      title="Flashcard Learning"
    />
  );
}
