import { api, unwrapData } from '@/api/api'
import type { APIResponse } from '@/types/api'
import type { WritingSubmission } from '@/types/writing'

// AI grading can take a while.
const slow = { timeout: 90_000 }

export async function submitWriting(input: { activityId: number; text: string }) {
  const response = await api.post<APIResponse<WritingSubmission>>('/writing/submissions', input, slow)
  return unwrapData<WritingSubmission>(response)
}

export async function rewriteWriting(submissionId: number, text: string) {
  const response = await api.post<APIResponse<WritingSubmission>>(
    `/writing/submissions/${submissionId}/rewrite`,
    { text },
    slow,
  )
  return unwrapData<WritingSubmission>(response)
}
