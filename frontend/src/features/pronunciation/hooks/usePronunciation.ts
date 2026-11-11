import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'

import { assessPronunciation, getPronunciationStatus } from '@/api/pronunciation'

import { phonemeQueryKey } from './usePhonemes'

export function usePronunciationStatus() {
  return useQuery({
    queryKey: ['pronunciation', 'status'],
    queryFn: getPronunciationStatus,
    staleTime: 10 * 60_000,
  })
}

export function useAssessPronunciation() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: assessPronunciation,
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: phonemeQueryKey })
    },
  })
}
