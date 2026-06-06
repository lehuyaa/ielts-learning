import { api, unwrapData } from '@/api/api'
import type { APIResponse } from '@/types/api'
import type { TopicDetailResponse } from '@/types/topic'

export async function getTopicDetail(topicId: string | number) {
  const response = await api.get<APIResponse<TopicDetailResponse>>(
    `/topics/${topicId}`,
  )

  return unwrapData<TopicDetailResponse>(response)
}
