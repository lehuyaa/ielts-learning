import { api, unwrapData } from '@/api/api'
import type { APIResponse } from '@/types/api'
import type {
  AuthResponse,
  LoginInput,
  RegisterInput,
} from '@/types/auth'
import type { User } from '@/types/user'

export async function login(input: LoginInput) {
  const response = await api.post<APIResponse<AuthResponse>>(
    '/auth/login',
    input,
  )

  return unwrapData<AuthResponse>(response)
}

export async function register(input: RegisterInput) {
  const response = await api.post<APIResponse<AuthResponse>>(
    '/auth/register',
    input,
  )

  return unwrapData<AuthResponse>(response)
}

export async function getMe() {
  const response = await api.get<APIResponse<User>>('/auth/me')

  return unwrapData<User>(response)
}
