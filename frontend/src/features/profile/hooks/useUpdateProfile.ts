import { useMutation, useQueryClient } from '@tanstack/react-query'

import { APIError } from '@/api/api'
import { updateProfile } from '@/api/profile'
import { useToast } from '@/contexts/toast/useToast'
import { currentUserQueryKey } from '@/features/auth/hooks/useCurrentUser'
import type { UpdateProfileInput } from '@/types/profile'

import { profileQueryKey } from './useProfile'

export function useUpdateProfile() {
  const queryClient = useQueryClient()
  const { showToast } = useToast()

  return useMutation({
    mutationFn: (payload: UpdateProfileInput) => updateProfile(payload),
    onSuccess: () => {
      showToast({
        title: 'Profile updated',
        description: 'Your profile changes were saved successfully.',
        tone: 'success',
      })
      void queryClient.invalidateQueries({ queryKey: profileQueryKey })
      void queryClient.invalidateQueries({ queryKey: currentUserQueryKey })
    },
    onError: (error) => {
      if (error instanceof APIError && error.fields) {
        return
      }

      showToast({
        title: 'Could not update profile',
        description:
          error instanceof Error
            ? error.message
            : 'Please try updating your profile again.',
        tone: 'error',
      })
    },
  })
}
