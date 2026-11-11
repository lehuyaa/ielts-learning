import { api, unwrapData } from '@/api/api'
import type { APIResponse } from '@/types/api'
import type { AssessInput, PronunciationResult, PronunciationStatus } from '@/types/pronunciation'

export async function getPronunciationStatus() {
  const response = await api.get<APIResponse<PronunciationStatus>>('/pronunciation/status')
  return unwrapData<PronunciationStatus>(response)
}

export async function assessPronunciation(input: AssessInput) {
  const formData = new FormData()
  formData.append('audio', input.audio, 'recording.wav')
  formData.append('referenceText', input.referenceText)
  formData.append('source', input.source)
  if (input.exerciseId) {
    formData.append('exerciseId', String(input.exerciseId))
  }
  const response = await api.post<APIResponse<PronunciationResult>>('/pronunciation/assess', formData, {
    headers: { 'Content-Type': 'multipart/form-data' },
    timeout: 45_000,
  })
  return unwrapData<PronunciationResult>(response)
}
