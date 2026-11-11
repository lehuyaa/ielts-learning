import { useMutation, useQueryClient } from '@tanstack/react-query'

import { submitActivityAnswer } from '@/api/activity'
import type { ActivityDetailResponse, SubmitAnswerInput } from '@/types/activity'

import { activityQueryKey } from './useActivity'

export function useSubmitAnswer(activityId: number | string) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (input: SubmitAnswerInput) =>
      submitActivityAnswer(activityId, input),
    onSuccess: (result) => {
      // Keep the cached activity in sync so resuming shows this answer.
      queryClient.setQueryData<ActivityDetailResponse>(
        activityQueryKey(activityId),
        (current) => {
          if (!current) {
            return current
          }
          const answers = current.answers.filter(
            (answer) => answer.exerciseId !== result.exerciseId,
          )
          return {
            ...current,
            activity: { ...current.activity, currentAttempt: result.attempt },
            answers: [...answers, result],
          }
        },
      )
    },
  })
}
