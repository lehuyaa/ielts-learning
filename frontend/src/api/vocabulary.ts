import { api, unwrapData } from '@/api/api'
import type { APIResponse } from '@/types/api'
import type {
  ImportResultResponse,
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

export async function previewVocabularyImport(file: File) {
  const formData = new FormData()
  formData.append('file', file)

  const response = await api.post<APIResponse<ImportResultResponse>>(
    '/admin/vocabularies/import/preview',
    formData,
  )

  return unwrapData<ImportResultResponse>(response)
}

export async function commitVocabularyImport(file: File) {
  const formData = new FormData()
  formData.append('file', file)

  const response = await api.post<APIResponse<ImportResultResponse>>(
    '/admin/vocabularies/import',
    formData,
  )

  return unwrapData<ImportResultResponse>(response)
}
