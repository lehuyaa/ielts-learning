import { api, unwrapData } from '@/api/api'
import type { APIResponse } from '@/types/api'
import type { FoundationSummary } from '@/types/foundation'

export async function getFoundationSummary(level = 'A1') {
  const response = await api.get<APIResponse<FoundationSummary>>('/foundation/summary', { params: { level } })
  return unwrapData<FoundationSummary>(response)
}
