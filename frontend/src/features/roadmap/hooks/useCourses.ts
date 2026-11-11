import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'

import { getCourses, setActiveCourse } from '@/api/course'
import { useToast } from '@/contexts/toast/useToast'

import { roadmapQueryKey } from './useRoadmap'

export const coursesQueryKey = ['courses'] as const

export function useCourses() {
  return useQuery({
    queryKey: coursesQueryKey,
    queryFn: getCourses,
    staleTime: 5 * 60 * 1000,
  })
}

export function useSetActiveCourse() {
  const queryClient = useQueryClient()
  const { showToast } = useToast()

  return useMutation({
    mutationFn: (courseId: number) => setActiveCourse(courseId),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: coursesQueryKey })
      void queryClient.invalidateQueries({ queryKey: roadmapQueryKey })
    },
    onError: (error) => {
      showToast({
        title: 'Could not switch course',
        description:
          error instanceof Error ? error.message : 'Please try again.',
        tone: 'error',
      })
    },
  })
}
