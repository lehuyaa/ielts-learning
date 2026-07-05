import { api, unwrapData } from '@/api/api'
import type { APIResponse } from '@/types/api'
import type { ProfileResponse, UpdateProfileInput, ProfileUser } from '@/types/profile'

export async function getProfile() {
  const response = await api.get<APIResponse<ProfileResponse>>('/me/profile')

  return unwrapData<ProfileResponse>(response)
}

export async function updateProfile(payload: UpdateProfileInput) {
  const response = await api.patch<APIResponse<ProfileUser>>('/me/profile', payload)

  return unwrapData<ProfileUser>(response)
}
