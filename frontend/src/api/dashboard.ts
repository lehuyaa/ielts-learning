import { api, unwrapData } from '@/api/api'
import type { APIResponse } from '@/types/api'
import type { DashboardSummaryResponse } from '@/types/dashboard'

export async function getDashboardSummary() {
  const response = await api.get<APIResponse<DashboardSummaryResponse>>(
    '/dashboard/summary',
  )

  return unwrapData<DashboardSummaryResponse>(response)
}
