import { useQuery } from '@tanstack/react-query'

import { getActivity } from '@/api/activity'

export function activityQueryKey(activityId: number | string) {
  return ['activity', String(activityId)] as const
}

export function useActivity(activityId: number | string | undefined) {
  return useQuery({
    queryKey: activityQueryKey(activityId ?? ''),
    queryFn: () => getActivity(activityId as number | string),
    enabled: activityId !== undefined && activityId !== '',
  })
}
