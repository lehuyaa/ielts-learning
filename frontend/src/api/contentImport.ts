import { api, unwrapData } from '@/api/api'
import type { APIResponse } from '@/types/api'
import type { ContentImportResult } from '@/types/contentImport'
import type { PassageCard } from '@/types/passage'

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

export type PassageDraft = {
  slug: string
  cefrLevel: string
  genre: string
  title: string
  bodyText: string
  audioUrl: string
  segments: { index: number; speaker?: string; text: string; start?: number; end?: number; audioUrl?: string }[]
  glossary: unknown[] | null
  source: string
  duration?: number
}

/** Whisper import (Task 64): audio → passage draft. Can take a minute. */
export async function transcribeAudio(input: { file: File; title: string; slug: string; genre: string; source: string }) {
  const formData = new FormData()
  formData.append('file', input.file)
  formData.append('title', input.title)
  formData.append('slug', input.slug)
  formData.append('genre', input.genre)
  formData.append('source', input.source)
  const response = await api.post<APIResponse<PassageDraft>>('/admin/passages/transcribe', formData, {
    ...multipartConfig,
    timeout: 5 * 60_000,
  })
  return unwrapData<PassageDraft>(response)
}

export async function savePassage(draft: Omit<PassageDraft, 'duration'>) {
  const response = await api.post<APIResponse<PassageCard>>('/admin/passages', draft)
  return unwrapData<PassageCard>(response)
}

export type ExtractResult = { bundle: unknown; summary: string; warnings: string[] }

/** PDF / image extraction (Task 65): pages → content bundle draft. */
export async function extractContent(input: {
  files: File[]
  topicTitle: string
  topicSlug: string
  notes: string
  publish: boolean
}) {
  const formData = new FormData()
  for (const file of input.files) formData.append('files', file)
  formData.append('topicTitle', input.topicTitle)
  formData.append('topicSlug', input.topicSlug)
  formData.append('notes', input.notes)
  formData.append('publish', String(input.publish))
  const response = await api.post<APIResponse<ExtractResult>>('/admin/content/extract', formData, {
    ...multipartConfig,
    timeout: 5 * 60_000,
  })
  return unwrapData<ExtractResult>(response)
}
