import axios, { type AxiosResponse } from 'axios'

import type { APIErrorBody, APIResponse } from '@/types/api'

const API_BASE_URL =
  import.meta.env.VITE_API_BASE_URL ?? 'http://localhost:8080/api/v1'

export class APIError extends Error {
  code: string
  status: number
  fields?: Record<string, string>

  constructor(
    message: string,
    code: string,
    status: number,
    fields?: Record<string, string>,
  ) {
    super(message)
    this.name = 'APIError'
    this.code = code
    this.status = status
    this.fields = fields
  }
}

export function getAccessToken() {
  return localStorage.getItem('accessToken')
}

export function setAccessToken(token: string) {
  localStorage.setItem('accessToken', token)
}

export function clearAccessToken() {
  localStorage.removeItem('accessToken')
}

export const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
})

api.interceptors.request.use((config) => {
  const token = getAccessToken()

  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }

  return config
})

api.interceptors.response.use(
  (response) => response,
  (error: unknown) => {
    if (axios.isAxiosError<APIResponse<unknown>>(error)) {
      const status = error.response?.status ?? 0
      const apiError = error.response?.data?.error

      if (status === 401) {
        clearAccessToken()
      }

      return Promise.reject(
        new APIError(
          apiError?.message ?? error.message ?? 'Something went wrong',
          apiError?.code ?? 'REQUEST_FAILED',
          status,
          apiError?.fields,
        ),
      )
    }

    return Promise.reject(
      new APIError('Something went wrong', 'REQUEST_FAILED', 0),
    )
  },
)

export function unwrapData<T>(response: AxiosResponse<APIResponse<T>>): T {
  const body = response.data

  if (body.error) {
    throw new APIError(
      body.error.message,
      body.error.code,
      response.status,
      body.error.fields,
    )
  }

  return body.data as T
}

export type { APIErrorBody, APIResponse }
