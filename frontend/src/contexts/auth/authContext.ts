import { createContext } from 'react'

import type { LoginInput, RegisterInput } from '@/types/auth'

export type AuthContextValue = {
  accessToken: string | null
  isAuthenticated: boolean
  login: (input: LoginInput) => Promise<void>
  register: (input: RegisterInput) => Promise<void>
  logout: () => void
}

export const AuthContext = createContext<AuthContextValue | null>(null)
