import { api, unwrapData } from '@/api/api'
import type { APIResponse } from '@/types/api'
import type { LessonDetailResponse, StartLessonResponse } from '@/types/lesson'

export async function getLessonDetail(lessonId: string | number) {
  const response = await api.get<APIResponse<LessonDetailResponse>>(
    `/lessons/${lessonId}`,
  )

  return unwrapData<LessonDetailResponse>(response)
}

export async function startLesson(lessonId: string | number) {
  const response = await api.post<APIResponse<StartLessonResponse>>(
    `/lessons/${lessonId}/start`,
  )

  return unwrapData<StartLessonResponse>(response)
}
