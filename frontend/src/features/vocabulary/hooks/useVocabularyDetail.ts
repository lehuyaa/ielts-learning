import { useQuery } from '@tanstack/react-query'

import { getVocabularyDetail } from '@/api/vocabulary'

export function vocabularyDetailQueryKey(vocabularyId: string | number) {
  return ['vocabulary', String(vocabularyId)] as const
}

export function useVocabularyDetail(vocabularyId: string | number | undefined) {
  return useQuery({
    queryKey: vocabularyDetailQueryKey(vocabularyId ?? ''),
    queryFn: () => getVocabularyDetail(vocabularyId ?? ''),
    enabled: Boolean(vocabularyId),
    staleTime: 5 * 60 * 1000,
  })
}
