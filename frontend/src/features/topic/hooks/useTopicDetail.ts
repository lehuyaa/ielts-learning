import { useQuery } from '@tanstack/react-query'

import { getTopicDetail } from '@/api/topic'

export function topicDetailQueryKey(topicId: string | number) {
  return ['topic', String(topicId)] as const
}

export function useTopicDetail(topicId: string | number | undefined) {
  return useQuery({
    queryKey: topicDetailQueryKey(topicId ?? ''),
    queryFn: () => getTopicDetail(topicId ?? ''),
    enabled: Boolean(topicId),
    staleTime: 5 * 60 * 1000,
  })
}
