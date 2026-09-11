import { useQuery } from '@tanstack/react-query'

import { listScenarios } from '@/api/aiConversation'

export const scenariosQueryKey = ['ai-conversation', 'scenarios'] as const

export function useScenarios() {
  return useQuery({
    queryKey: scenariosQueryKey,
    queryFn: listScenarios,
  })
}
