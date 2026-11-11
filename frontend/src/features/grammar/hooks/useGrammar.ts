import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'

import {
  answerGrammarPractice,
  getGrammarPoint,
  getGrammarPoints,
  getGrammarPractice,
} from '@/api/grammar'
import type { ExerciseAnswer } from '@/types/activity'

export const grammarQueryKey = ['grammar-points'] as const

export function useGrammarPoints(level?: string) {
  return useQuery({
    queryKey: [...grammarQueryKey, 'list', level ?? 'all'],
    queryFn: () => getGrammarPoints(level),
  })
}

export function useGrammarPoint(slug: string | undefined) {
  return useQuery({
    queryKey: [...grammarQueryKey, 'detail', slug],
    queryFn: () => getGrammarPoint(slug as string),
    enabled: Boolean(slug),
  })
}

/** A practice set is fetched on demand so every "Practice" click is fresh. */
export function useStartGrammarPractice(slug: string) {
  return useMutation({ mutationFn: () => getGrammarPractice(slug) })
}

export function useAnswerGrammarPractice(slug: string) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (input: { exerciseId: number; answer: ExerciseAnswer }) =>
      answerGrammarPractice(slug, input),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: grammarQueryKey })
    },
  })
}
