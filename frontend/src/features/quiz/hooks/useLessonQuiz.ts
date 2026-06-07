import { useQuery } from "@tanstack/react-query";

import { getLessonQuiz } from "@/api/quiz";

export function lessonQuizQueryKey(lessonId: string | number) {
  return ["lessonQuiz", String(lessonId)] as const;
}

export function useLessonQuiz(lessonId: string | number | undefined) {
  return useQuery({
    queryKey: lessonQuizQueryKey(lessonId ?? ""),
    queryFn: () => getLessonQuiz(lessonId ?? ""),
    enabled: Boolean(lessonId),
    staleTime: 5 * 60 * 1000,
  });
}
