import { useMutation, useQueryClient } from '@tanstack/react-query'

import { updateProfile } from '@/api/profile'
import { currentUserQueryKey } from '@/features/auth/hooks/useCurrentUser'
import type { UpdateProfileInput } from '@/types/profile'

import { profileQueryKey } from './useProfile'

export function useUpdateProfile() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (payload: UpdateProfileInput) => updateProfile(payload),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: profileQueryKey })
      void queryClient.invalidateQueries({ queryKey: currentUserQueryKey })
    },
  })
}
