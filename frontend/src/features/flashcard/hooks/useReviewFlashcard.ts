import { useMutation, useQueryClient } from "@tanstack/react-query";

import { reviewFlashcard } from "@/api/flashcard";
import { lessonDetailQueryKey } from "@/features/lesson/hooks/useLessonDetail";
import { roadmapQueryKey } from "@/features/roadmap/hooks/useRoadmap";
import { topicDetailQueryKey } from "@/features/topic/hooks/useTopicDetail";

import { lessonFlashcardsQueryKey } from "./useLessonFlashcards";

type UseReviewFlashcardOptions = {
  lessonId?: string | number;
  topicId?: string | number;
};

export function useReviewFlashcard(options: UseReviewFlashcardOptions = {}) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: reviewFlashcard,
    onSuccess: (result) => {
      void queryClient.invalidateQueries({
        queryKey: ["dueReviews"],
      });

      if (options.lessonId) {
        void queryClient.invalidateQueries({
          queryKey: lessonFlashcardsQueryKey(options.lessonId),
        });
        void queryClient.invalidateQueries({
          queryKey: lessonDetailQueryKey(options.lessonId),
        });
      }

      if (options.topicId) {
        void queryClient.invalidateQueries({
          queryKey: topicDetailQueryKey(options.topicId),
        });
      }

      void queryClient.invalidateQueries({
        queryKey: roadmapQueryKey,
      });
      void queryClient.invalidateQueries({
        queryKey: ["vocabulary", String(result.vocabularyId)],
      });
      void queryClient.invalidateQueries({
        queryKey: ["vocabularies"],
      });
    },
  });
}
