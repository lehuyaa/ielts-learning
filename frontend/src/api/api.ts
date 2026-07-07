import axios, { AxiosError, type AxiosResponse } from 'axios'

import type { APIErrorBody, APIResponse } from '@/types/api'

const API_BASE_URL =
  import.meta.env.VITE_API_BASE_URL ?? 'http://localhost:8080/api/v1'

export class APIError extends Error {
  code: string
  status: number
  fields?: Record<string, string>
  originalError?: unknown
  isNetworkError: boolean

  constructor(
    message: string,
    code: string,
    status: number,
    fields?: Record<string, string>,
    options?: {
      originalError?: unknown
      isNetworkError?: boolean
    },
  ) {
    super(message)
    this.name = 'APIError'
    this.code = code
    this.status = status
    this.fields = fields
    this.originalError = options?.originalError
    this.isNetworkError = options?.isNetworkError ?? false
  }
}

type UnauthorizedHandler = (error: APIError) => void

let unauthorizedHandler: UnauthorizedHandler | null = null
let lastUnauthorizedHandledAt = 0
const UNAUTHORIZED_COOLDOWN_MS = 1000

export function getAccessToken() {
  return localStorage.getItem('accessToken')
}

export function setAccessToken(token: string) {
  localStorage.setItem('accessToken', token)
}

export function clearAccessToken() {
  localStorage.removeItem('accessToken')
}

export function setUnauthorizedHandler(handler: UnauthorizedHandler | null) {
  unauthorizedHandler = handler
}

export function getApiErrorMessage(
  error: unknown,
  fallback = 'Something went wrong. Please try again.',
) {
  if (error instanceof APIError) {
    return error.message
  }

  if (error instanceof Error && error.message.trim()) {
    return error.message
  }

  return fallback
}

export function isAPIError(error: unknown): error is APIError {
  return error instanceof APIError
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
      const normalizedError = normalizeAxiosError(error)

      if (normalizedError.status === 401) {
        const now = Date.now()
        if (
          unauthorizedHandler &&
          now - lastUnauthorizedHandledAt > UNAUTHORIZED_COOLDOWN_MS
        ) {
          lastUnauthorizedHandledAt = now
          unauthorizedHandler(normalizedError)
        }
      }

      return Promise.reject(normalizedError)
    }

    return Promise.reject(
      new APIError(
        'Something went wrong. Please try again.',
        'REQUEST_FAILED',
        0,
        undefined,
        { originalError: error },
      ),
    )
  },
)

export function unwrapData<T>(response: AxiosResponse<APIResponse<T>>): T {
  const body = response.data

  if (body.error) {
    throw new APIError(
      sanitizeMessage(body.error.message, response.status),
      body.error.code,
      response.status,
      body.error.fields,
    )
  }

  return body.data as T
}

export type { APIErrorBody, APIResponse }

function normalizeAxiosError(error: AxiosError<APIResponse<unknown>>) {
  const status = error.response?.status ?? 0
  const apiError = error.response?.data?.error
  const isNetworkError = !error.response

  return new APIError(
    resolveErrorMessage(error, apiError, status),
    apiError?.code ?? (isNetworkError ? 'NETWORK_ERROR' : 'REQUEST_FAILED'),
    status,
    apiError?.fields,
    {
      originalError: error,
      isNetworkError,
    },
  )
}

function resolveErrorMessage(
  error: AxiosError<APIResponse<unknown>>,
  apiError: APIErrorBody | undefined,
  status: number,
) {
  if (!error.response) {
    return 'Unable to connect. Please check your internet connection.'
  }

  return sanitizeMessage(apiError?.message, status)
}

function sanitizeMessage(message: string | undefined, status: number) {
  const trimmedMessage = message?.trim()
  if (trimmedMessage) {
    return trimmedMessage
  }

  if (status === 401) {
    return 'Your session has expired. Please sign in again.'
  }

  return 'Something went wrong. Please try again.'
}
