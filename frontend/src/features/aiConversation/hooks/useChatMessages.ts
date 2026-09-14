import { useQuery } from '@tanstack/react-query'

import { listChatMessages } from '@/api/aiConversation'

export function chatMessagesQueryKey(scenarioSlug: string | undefined) {
  return ['ai-conversation', 'messages', scenarioSlug] as const
}

export function useChatMessages(scenarioSlug: string | undefined) {
  return useQuery({
    queryKey: chatMessagesQueryKey(scenarioSlug),
    queryFn: () => listChatMessages(scenarioSlug as string),
    enabled: Boolean(scenarioSlug),
  })
}
