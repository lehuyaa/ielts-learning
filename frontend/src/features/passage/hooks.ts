import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'

import { addWordToFlashcards, getPassage, lookupWord } from '@/api/passage'

export function usePassage(passageId: number | undefined) {
  return useQuery({
    queryKey: ['passages', passageId],
    queryFn: () => getPassage(passageId as number),
    enabled: Boolean(passageId),
    staleTime: 5 * 60_000,
  })
}

export function useWordLookup(word: string | null, passageId?: number) {
  return useQuery({
    queryKey: ['word-lookup', word?.toLowerCase(), passageId ?? 0],
    queryFn: () => lookupWord(word as string, passageId),
    enabled: Boolean(word),
    staleTime: 10 * 60_000,
  })
}

export function useAddWordToFlashcards() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: addWordToFlashcards,
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: ['word-lookup'] })
    },
  })
}
