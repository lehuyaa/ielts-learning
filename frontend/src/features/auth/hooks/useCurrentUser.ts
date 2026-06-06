import { useQuery } from '@tanstack/react-query'

import { getMe } from '@/api/auth'
import { useAuth } from '@/contexts/auth/useAuth'

export const currentUserQueryKey = ['auth', 'currentUser'] as const

export function useCurrentUser() {
  const { accessToken } = useAuth()

  return useQuery({
    queryKey: currentUserQueryKey,
    queryFn: getMe,
    enabled: Boolean(accessToken),
    staleTime: 5 * 60 * 1000,
    retry: false,
  })
}
