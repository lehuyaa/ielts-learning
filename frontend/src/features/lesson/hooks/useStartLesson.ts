import { useMutation, useQueryClient } from '@tanstack/react-query'

import { startLesson } from '@/api/lesson'

import { lessonDetailQueryKey } from './useLessonDetail'

type UseStartLessonOptions = {
  lessonId: string | number | undefined
  topicId?: string | number
}

export function useStartLesson({ lessonId, topicId }: UseStartLessonOptions) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: () => startLesson(lessonId ?? ''),
    onSuccess: () => {
      if (lessonId) {
        void queryClient.invalidateQueries({
          queryKey: lessonDetailQueryKey(lessonId),
        })
      }

      if (topicId) {
        void queryClient.invalidateQueries({
          queryKey: ['topic', String(topicId)],
        })
      }

      void queryClient.invalidateQueries({ queryKey: ['roadmap'] })
    },
  })
}
