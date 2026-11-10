import { api, unwrapData } from '@/api/api'
import type { APIResponse } from '@/types/api'
import type { ContentImportResult } from '@/types/contentImport'

// Multipart must be declared: the shared instance defaults to JSON.
const multipartConfig = {
  headers: { 'Content-Type': 'multipart/form-data' },
}

async function postBundle(path: string, file: File) {
  const formData = new FormData()
  formData.append('file', file)

  const response = await api.post<APIResponse<ContentImportResult>>(
    path,
    formData,
    multipartConfig,
  )

  return unwrapData<ContentImportResult>(response)
}

export function previewContentImport(file: File) {
  return postBundle('/admin/content/import/preview', file)
}

export function commitContentImport(file: File) {
  return postBundle('/admin/content/import', file)
}
