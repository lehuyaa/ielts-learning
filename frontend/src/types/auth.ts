import type { User } from '@/types/user'

export type AuthResponse = {
  accessToken: string
  tokenType: 'Bearer'
  user: User
}

export type LoginInput = {
  email: string
  password: string
}

export type RegisterInput = {
  name: string
  email: string
  username?: string
  password: string
  targetBand: number
}
