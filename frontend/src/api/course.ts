import { api, unwrapData } from '@/api/api'
import type { APIResponse } from '@/types/api'
import type {
  CourseListResponse,
  SetActiveCourseResponse,
} from '@/types/course'

export async function getCourses() {
  const response = await api.get<APIResponse<CourseListResponse>>('/courses')

  return unwrapData<CourseListResponse>(response)
}

export async function setActiveCourse(courseId: number) {
  const response = await api.patch<APIResponse<SetActiveCourseResponse>>(
    '/me/active-course',
    { courseId },
  )

  return unwrapData<SetActiveCourseResponse>(response)
}
