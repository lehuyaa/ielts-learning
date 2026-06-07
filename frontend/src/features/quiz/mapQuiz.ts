import type { LessonQuizResponse, QuizSessionViewModel } from "@/types/quiz";

export function mapLessonQuizToSessionViewModel(
  response: LessonQuizResponse,
): QuizSessionViewModel {
  return {
    lessonId: response.lesson.id,
    lessonTitle: response.lesson.title,
    topicTitle: "Lesson Quiz",
    band: `Required ${response.lesson.requiredScore}%`,
    requiredScore: response.lesson.requiredScore,
    timeLimitSeconds: response.lesson.timeLimitSeconds ?? 30,
    questions: response.questions,
  };
}
