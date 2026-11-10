import { api, unwrapData } from '@/api/api'
import type {
  ActivityDetailResponse,
  AnswerResult,
  LessonActivitiesResponse,
  SubmitAnswerInput,
} from '@/types/activity'
import type { APIResponse } from '@/types/api'

export async function getLessonActivities(lessonId: number | string) {
  const response = await api.get<APIResponse<LessonActivitiesResponse>>(
    `/lessons/${lessonId}/activities`,
  )

  return unwrapData<LessonActivitiesResponse>(response)
}

export async function getActivity(activityId: number | string) {
  const response = await api.get<APIResponse<ActivityDetailResponse>>(
    `/activities/${activityId}`,
  )

  return unwrapData<ActivityDetailResponse>(response)
}

export async function submitActivityAnswer(
  activityId: number | string,
  input: SubmitAnswerInput,
) {
  const response = await api.post<APIResponse<AnswerResult>>(
    `/activities/${activityId}/answers`,
    input,
  )

  return unwrapData<AnswerResult>(response)
}
