import { useQuery } from '@tanstack/react-query'

import { getRoadmap } from '@/api/roadmap'

export const roadmapQueryKey = ['roadmap'] as const

export function useRoadmap() {
  return useQuery({
    queryKey: roadmapQueryKey,
    queryFn: getRoadmap,
    staleTime: 5 * 60 * 1000,
  })
}
