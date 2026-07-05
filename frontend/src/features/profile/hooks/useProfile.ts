import { useQuery } from '@tanstack/react-query'

import { getProfile } from '@/api/profile'

export const profileQueryKey = ['profile'] as const

export function useProfile() {
  return useQuery({
    queryKey: profileQueryKey,
    queryFn: getProfile,
    staleTime: 5 * 60 * 1000,
  })
}
