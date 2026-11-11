import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { useNavigate } from 'react-router-dom'

import { getMistakeReview, getMistakes, reviewMistake, updateMistake } from '@/api/mistakes'
import { generatePractice } from '@/api/practice'
import type { MistakeRating, MistakeStatus } from '@/types/mistakes'

export const mistakesQueryKey = ['mistakes'] as const

export function useMistakes(status: 'open' | 'fixed' | 'all') {
  return useQuery({ queryKey: [...mistakesQueryKey, 'list', status], queryFn: () => getMistakes(status) })
}

export function useUpdateMistake() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ id, status }: { id: number; status: MistakeStatus }) => updateMistake(id, status),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: mistakesQueryKey })
    },
  })
}

export function useMistakeReview() {
  return useQuery({ queryKey: [...mistakesQueryKey, 'review'], queryFn: getMistakeReview, staleTime: Infinity })
}

export function useReviewMistake() {
  return useMutation({
    mutationFn: ({ id, rating }: { id: number; rating: MistakeRating }) => reviewMistake(id, rating),
  })
}

/** Builds (or reopens) today's practice for a skill and opens it. */
export function useStartPractice() {
  const navigate = useNavigate()
  return useMutation({
    mutationFn: (input: { skillTag: string; returnTo?: string }) => generatePractice(input.skillTag),
    onSuccess: (set, input) => {
      void navigate(`/practice/${set.id}`, { state: { returnTo: input.returnTo } })
    },
  })
}
