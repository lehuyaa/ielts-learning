import { useMutation } from '@tanstack/react-query'

import { sendChatMessage } from '@/api/aiConversation'

type SendChatMessageVariables = {
  scenarioSlug: string
  message: string
}

export function useSendChatMessage() {
  return useMutation({
    mutationFn: ({ scenarioSlug, message }: SendChatMessageVariables) =>
      sendChatMessage(scenarioSlug, { message }),
  })
}
