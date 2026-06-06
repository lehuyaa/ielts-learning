import { api, unwrapData } from '@/api/api'
import type { APIResponse } from '@/types/api'
import type {
  VocabularyDetailResponse,
  VocabularyListResponse,
  VocabularyQueryParams,
} from '@/types/vocabulary'

export async function getVocabularies(params: VocabularyQueryParams) {
  const response = await api.get<APIResponse<VocabularyListResponse>>(
    '/vocabularies',
    {
      params,
    },
  )

  return unwrapData<VocabularyListResponse>(response)
}

export async function getVocabularyDetail(vocabularyId: string | number) {
  const response = await api.get<APIResponse<VocabularyDetailResponse>>(
    `/vocabularies/${vocabularyId}`,
  )

  return unwrapData<VocabularyDetailResponse>(response)
}
