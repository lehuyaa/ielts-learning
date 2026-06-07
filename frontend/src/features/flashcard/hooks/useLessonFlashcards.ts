import { useQuery } from "@tanstack/react-query";

import { getLessonFlashcards } from "@/api/flashcard";

export function lessonFlashcardsQueryKey(lessonId: string | number) {
  return ["lessonFlashcards", String(lessonId)] as const;
}

export function useLessonFlashcards(lessonId: string | number | undefined) {
  return useQuery({
    queryKey: lessonFlashcardsQueryKey(lessonId ?? ""),
    queryFn: () => getLessonFlashcards(lessonId ?? ""),
    enabled: Boolean(lessonId),
    staleTime: 5 * 60 * 1000,
  });
}
