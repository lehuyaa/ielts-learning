import { useMutation } from '@tanstack/react-query'
import { useNavigate } from 'react-router-dom'

import { startSession } from '@/api/aiConversation'

/** Starts (or resumes) a role-play and opens its page. */
export function useStartSession() {
  const navigate = useNavigate()
  return useMutation({
    mutationFn: (input: { scenarioSlug: string; activityId?: number; fresh?: boolean; returnTo?: string }) =>
      startSession({ scenarioSlug: input.scenarioSlug, activityId: input.activityId, fresh: input.fresh }),
    onSuccess: (session, input) => {
      void navigate(`/ai-conversation/sessions/${session.id}`, { state: { returnTo: input.returnTo } })
    },
  })
}
