import { useMutation, useQueryClient } from '@tanstack/react-query'

import { completeActivity } from '@/api/activity'
import { roadmapQueryKey } from '@/features/roadmap/hooks/useRoadmap'

import { activityQueryKey } from './useActivity'
import { lessonActivitiesQueryKey } from './useLessonActivities'

export function useCompleteActivity(
  activityId: number | string,
  lessonId: number | string,
) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: () => completeActivity(activityId),
    onSuccess: () => {
      void queryClient.invalidateQueries({
        queryKey: lessonActivitiesQueryKey(lessonId),
      })
      void queryClient.invalidateQueries({ queryKey: roadmapQueryKey })
      // The finished attempt is closed; the next visit starts a fresh one.
      void queryClient.invalidateQueries({
        queryKey: activityQueryKey(activityId),
      })
    },
  })
}
