import { useMutation, useQueryClient } from '@tanstack/react-query'

import { getApiErrorMessage } from '@/api/api'
import { commitContentImport, previewContentImport } from '@/api/contentImport'
import { useToast } from '@/contexts/toast/useToast'
import { roadmapQueryKey } from '@/features/roadmap/hooks/useRoadmap'

export function useContentImportPreview() {
  const { showToast } = useToast()

  return useMutation({
    mutationFn: (file: File) => previewContentImport(file),
    onError: (error) => {
      showToast({
        title: 'Could not read the file',
        description: getApiErrorMessage(error, 'Please check the JSON file and try again.'),
        tone: 'error',
      })
    },
  })
}

export function useContentImportCommit() {
  const queryClient = useQueryClient()
  const { showToast } = useToast()

  return useMutation({
    mutationFn: (file: File) => commitContentImport(file),
    onSuccess: (result) => {
      showToast({
        title: 'Import complete',
        description: `${result.summary.created} created, ${result.summary.updated} updated, ${result.summary.invalid} skipped.`,
        tone: 'success',
      })
      void queryClient.invalidateQueries({ queryKey: roadmapQueryKey })
      void queryClient.invalidateQueries({ queryKey: ['courses'] })
      void queryClient.invalidateQueries({ queryKey: ['lesson-activities'] })
      void queryClient.invalidateQueries({ queryKey: ['activity'] })
    },
    onError: (error) => {
      showToast({
        title: 'Import failed',
        description: getApiErrorMessage(error, 'Nothing was imported. Please try again.'),
        tone: 'error',
      })
    },
  })
}
