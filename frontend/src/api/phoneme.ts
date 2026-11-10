import { api, unwrapData } from '@/api/api'
import type { APIResponse } from '@/types/api'
import type { PhonemeDetailResponse, PhonemeListResponse } from '@/types/phoneme'

export async function getPhonemes() {
  const response = await api.get<APIResponse<PhonemeListResponse>>('/phonemes')
  return unwrapData<PhonemeListResponse>(response)
}

export async function getPhoneme(symbol: string) {
  const response = await api.get<APIResponse<PhonemeDetailResponse>>(
    `/phonemes/${encodeURIComponent(symbol)}`,
  )
  return unwrapData<PhonemeDetailResponse>(response)
}
