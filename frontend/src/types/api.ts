export type APIErrorBody = {
  code: string
  message: string
  fields?: Record<string, string>
}

export type APIResponse<T> = {
  data?: T
  error?: APIErrorBody
}
