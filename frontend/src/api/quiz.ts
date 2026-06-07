import { api, unwrapData } from "@/api/api";
import type { APIResponse } from "@/types/api";
import type {
  CheckQuizAnswerRequest,
  CheckQuizAnswerResponse,
  LessonQuizResponse,
  SubmitQuizRequest,
  SubmitQuizResponse,
} from "@/types/quiz";

export async function getLessonQuiz(lessonId: string | number) {
  const response = await api.get<APIResponse<LessonQuizResponse>>(
    `/lessons/${lessonId}/quiz`,
  );

  return unwrapData<LessonQuizResponse>(response);
}

export async function submitLessonQuiz(
  lessonId: string | number,
  payload: SubmitQuizRequest,
) {
  const response = await api.post<APIResponse<SubmitQuizResponse>>(
    `/lessons/${lessonId}/quiz/submit`,
    payload,
  );

  return unwrapData<SubmitQuizResponse>(response);
}

export async function checkQuizAnswer(
  lessonId: string | number,
  payload: CheckQuizAnswerRequest,
) {
  const response = await api.post<APIResponse<CheckQuizAnswerResponse>>(
    `/lessons/${lessonId}/quiz/check-answer`,
    payload,
  );

  return unwrapData<CheckQuizAnswerResponse>(response);
}
