import { useMutation } from '@tanstack/react-query'

import { sendChatMessage } from '@/api/aiConversation'
import type { SendChatMessageInput } from '@/types/aiConversation'

export function useSendChatMessage() {
  return useMutation({
    mutationFn: (payload: SendChatMessageInput) => sendChatMessage(payload),
  })
}
