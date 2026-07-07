import { useCallback, useEffect, useMemo, useState, type ReactNode } from 'react'
import { useMutation, useQueryClient } from '@tanstack/react-query'

import {
  login as loginRequest,
  register as registerRequest,
} from '@/api/auth'
import {
  APIError,
  clearAccessToken,
  getAccessToken,
  setUnauthorizedHandler,
  setAccessToken as persistAccessToken,
} from '@/api/api'
import { router } from '@/app/router'
import { currentUserQueryKey } from '@/features/auth/hooks/useCurrentUser'
import type { LoginInput, RegisterInput } from '@/types/auth'
import type { User } from '@/types/user'

import { AuthContext } from './authContext.ts'

export function AuthProvider({ children }: { children: ReactNode }) {
  const queryClient = useQueryClient()
  const [accessToken, setAccessToken] = useState(() => getAccessToken())

  const clearSession = useCallback(() => {
    clearAccessToken()
    setAccessToken(null)
    queryClient.clear()
  }, [queryClient])

  const saveSession = useCallback(
    (token: string, user: User) => {
      persistAccessToken(token)
      setAccessToken(token)
      queryClient.clear()
      queryClient.setQueryData(currentUserQueryKey, user)
    },
    [queryClient],
  )

  useEffect(() => {
    setUnauthorizedHandler((error: APIError) => {
      const pathname = window.location.pathname
      const isAuthPage =
        pathname === '/login' || pathname === '/register'

      clearSession()

      if (!isAuthPage) {
        void router.navigate('/login', {
          replace: true,
          state: {
            sessionExpired: true,
            reason: error.code,
          },
        })
      }
    })

    return () => {
      setUnauthorizedHandler(null)
    }
  }, [clearSession])

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
    clearSession()
  }, [clearSession])

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
