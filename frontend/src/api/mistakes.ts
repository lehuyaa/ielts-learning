import { api, unwrapData } from '@/api/api'
import type { APIResponse } from '@/types/api'
import type {
  MistakeItem,
  MistakeList,
  MistakeRating,
  MistakeReviewQueue,
  MistakeReviewResult,
  MistakeStatus,
} from '@/types/mistakes'

export async function getMistakes(status: 'open' | 'fixed' | 'all') {
  const response = await api.get<APIResponse<MistakeList>>('/mistakes', { params: { status } })
  return unwrapData<MistakeList>(response)
}

export async function updateMistake(id: number, status: MistakeStatus) {
  const response = await api.patch<APIResponse<MistakeItem>>(`/mistakes/${id}`, { status })
  return unwrapData<MistakeItem>(response)
}

export async function getMistakeReview() {
  const response = await api.get<APIResponse<MistakeReviewQueue>>('/mistakes/review', { params: { limit: 30 } })
  return unwrapData<MistakeReviewQueue>(response)
}

export async function reviewMistake(id: number, rating: MistakeRating) {
  const response = await api.post<APIResponse<MistakeReviewResult>>(`/mistakes/${id}/review`, { rating })
  return unwrapData<MistakeReviewResult>(response)
}
