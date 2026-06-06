import { useQuery } from '@tanstack/react-query'

import { getLessonDetail } from '@/api/lesson'

export function lessonDetailQueryKey(lessonId: string | number) {
  return ['lesson', String(lessonId)] as const
}

export function useLessonDetail(lessonId: string | number | undefined) {
  return useQuery({
    queryKey: lessonDetailQueryKey(lessonId ?? ''),
    queryFn: () => getLessonDetail(lessonId ?? ''),
    enabled: Boolean(lessonId),
    staleTime: 5 * 60 * 1000,
  })
}
