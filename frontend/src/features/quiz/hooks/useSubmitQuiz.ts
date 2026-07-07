import { useMutation, useQueryClient } from "@tanstack/react-query";

import { submitLessonQuiz } from "@/api/quiz";
import { useToast } from "@/contexts/toast/useToast";
import { lessonDetailQueryKey } from "@/features/lesson/hooks/useLessonDetail";
import { roadmapQueryKey } from "@/features/roadmap/hooks/useRoadmap";
import { topicDetailQueryKey } from "@/features/topic/hooks/useTopicDetail";
import type { SubmitQuizRequest } from "@/types/quiz";

import { lessonQuizQueryKey } from "./useLessonQuiz";

type UseSubmitQuizOptions = {
  lessonId: string | number | undefined;
  topicId?: string | number;
};

export function useSubmitQuiz({ lessonId, topicId }: UseSubmitQuizOptions) {
  const queryClient = useQueryClient();
  const { showToast } = useToast();

  return useMutation({
    mutationFn: (payload: SubmitQuizRequest) =>
      submitLessonQuiz(lessonId ?? "", payload),
    onSuccess: () => {
      if (lessonId) {
        void queryClient.invalidateQueries({
          queryKey: lessonQuizQueryKey(lessonId),
        });
        void queryClient.invalidateQueries({
          queryKey: lessonDetailQueryKey(lessonId),
        });
      }

      if (topicId) {
        void queryClient.invalidateQueries({
          queryKey: topicDetailQueryKey(topicId),
        });
      }

      void queryClient.invalidateQueries({
        queryKey: roadmapQueryKey,
      });
      void queryClient.invalidateQueries({
        queryKey: ["dashboard"],
      });
    },
    onError: (error) => {
      showToast({
        title: "Quiz submission failed",
        description:
          error instanceof Error
            ? error.message
            : "Please try submitting the quiz again.",
        tone: "error",
      });
    },
  });
}
