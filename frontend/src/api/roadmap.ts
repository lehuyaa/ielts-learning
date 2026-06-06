import { api, unwrapData } from '@/api/api'
import type { APIResponse } from '@/types/api'
import type { RoadmapResponse } from '@/types/roadmap'

export async function getRoadmap() {
  const response = await api.get<APIResponse<RoadmapResponse>>('/roadmap')

  return unwrapData<RoadmapResponse>(response)
}
