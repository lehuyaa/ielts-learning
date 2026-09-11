import { useMutation, useQueryClient } from '@tanstack/react-query'

import { APIError, getApiErrorMessage } from '@/api/api'
import { commitVocabularyImport, previewVocabularyImport } from '@/api/vocabulary'
import { useToast } from '@/contexts/toast/useToast'

export function useVocabularyImportPreview() {
  const { showToast } = useToast()

  return useMutation({
    mutationFn: (file: File) => previewVocabularyImport(file),
    onError: (error) => {
      showToast({
        title: 'Could not read the file',
        description: getApiErrorMessage(
          error,
          'Please check the file format and try again.',
        ),
        tone: 'error',
      })
    },
  })
}

export function useVocabularyImportCommit() {
  const queryClient = useQueryClient()
  const { showToast } = useToast()

  return useMutation({
    mutationFn: (file: File) => commitVocabularyImport(file),
    onSuccess: (result) => {
      showToast({
        title: 'Import complete',
        description: `${result.vocabulariesCreated} created, ${result.vocabulariesUpdated} updated.`,
        tone: 'success',
      })
      void queryClient.invalidateQueries({ queryKey: ['vocabularies'] })
    },
    onError: (error) => {
      showToast({
        title: 'Import failed',
        description:
          error instanceof APIError
            ? error.message
            : 'Please try importing the file again.',
        tone: 'error',
      })
    },
  })
}
