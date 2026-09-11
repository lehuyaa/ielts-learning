import { useMutation } from '@tanstack/react-query'

import { createScenario } from '@/api/aiConversation'
import type { CreateScenarioInput } from '@/types/aiConversation'

export function useCreateScenario() {
  return useMutation({
    mutationFn: (payload: CreateScenarioInput) => createScenario(payload),
  })
}
