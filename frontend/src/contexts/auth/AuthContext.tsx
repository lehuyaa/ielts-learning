import { useCallback, useMemo, useState, type ReactNode } from 'react'
import { useMutation, useQueryClient } from '@tanstack/react-query'

import {
  login as loginRequest,
  register as registerRequest,
} from '@/api/auth'
import {
  clearAccessToken,
  getAccessToken,
  setAccessToken as persistAccessToken,
} from '@/api/api'
import { currentUserQueryKey } from '@/features/auth/hooks/useCurrentUser'
import type { LoginInput, RegisterInput } from '@/types/auth'
import type { User } from '@/types/user'

import { AuthContext } from './authContext.ts'

export function AuthProvider({ children }: { children: ReactNode }) {
  const queryClient = useQueryClient()
  const [accessToken, setAccessToken] = useState(() => getAccessToken())

  const saveSession = useCallback(
    (token: string, user: User) => {
      persistAccessToken(token)
      setAccessToken(token)
      queryClient.clear()
      queryClient.setQueryData(currentUserQueryKey, user)
    },
    [queryClient],
  )

  const loginMutation = useMutation({
    mutationFn: loginRequest,
    onSuccess: (response) => {
      saveSession(response.accessToken, response.user)
    },
  })

  const registerMutation = useMutation({
    mutationFn: registerRequest,
    onSuccess: (response) => {
      saveSession(response.accessToken, response.user)
    },
  })

  const login = useCallback(
    async (input: LoginInput) => {
      await loginMutation.mutateAsync(input)
    },
    [loginMutation],
  )

  const register = useCallback(
    async (input: RegisterInput) => {
      await registerMutation.mutateAsync(input)
    },
    [registerMutation],
  )

  const logout = useCallback(() => {
    clearAccessToken()
    setAccessToken(null)
    queryClient.clear()
  }, [queryClient])

  const value = useMemo(
    () => ({
      accessToken,
      isAuthenticated: Boolean(accessToken),
      login,
      register,
      logout,
    }),
    [accessToken, login, logout, register],
  )

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}
