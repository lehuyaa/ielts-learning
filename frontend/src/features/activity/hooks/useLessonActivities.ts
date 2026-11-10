import { useQuery } from '@tanstack/react-query'

import { getLessonActivities } from '@/api/activity'

export function lessonActivitiesQueryKey(lessonId: number | string) {
  return ['lesson-activities', String(lessonId)] as const
}

export function useLessonActivities(lessonId: number | string | undefined) {
  return useQuery({
    queryKey: lessonActivitiesQueryKey(lessonId ?? ''),
    queryFn: () => getLessonActivities(lessonId as number | string),
    enabled: lessonId !== undefined && lessonId !== '',
  })
}
