import { useQuery } from '@tanstack/react-query'

import { getVocabularies } from '@/api/vocabulary'
import type { VocabularyQueryParams } from '@/types/vocabulary'

export function vocabulariesQueryKey(params: VocabularyQueryParams) {
  return ['vocabularies', params] as const
}

export function useVocabularies(params: VocabularyQueryParams) {
  return useQuery({
    queryKey: vocabulariesQueryKey(params),
    queryFn: () => getVocabularies(params),
    staleTime: 5 * 60 * 1000,
  })
}
