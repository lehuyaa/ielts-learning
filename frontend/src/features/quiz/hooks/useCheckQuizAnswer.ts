import { useMutation } from "@tanstack/react-query";

import { checkQuizAnswer } from "@/api/quiz";
import type { CheckQuizAnswerRequest } from "@/types/quiz";

type UseCheckQuizAnswerOptions = {
  lessonId: string | number | undefined;
};

export function useCheckQuizAnswer({ lessonId }: UseCheckQuizAnswerOptions) {
  return useMutation({
    mutationFn: (payload: CheckQuizAnswerRequest) =>
      checkQuizAnswer(lessonId ?? "", payload),
  });
}
