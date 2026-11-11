import { api, unwrapData } from '@/api/api'
import type { AnswerResult, SubmitAnswerInput } from '@/types/activity'
import type { APIResponse } from '@/types/api'
import type { PracticeSet } from '@/types/practice'

export async function generatePractice(skillTag: string) {
  const response = await api.post<APIResponse<PracticeSet>>('/practice/generate', { skillTag }, { timeout: 90_000 })
  return unwrapData<PracticeSet>(response)
}

export async function getPracticeSet(setId: number) {
  const response = await api.get<APIResponse<PracticeSet>>(`/practice/sets/${setId}`)
  return unwrapData<PracticeSet>(response)
}

export async function answerPractice(setId: number, input: SubmitAnswerInput) {
  const response = await api.post<APIResponse<AnswerResult>>(`/practice/sets/${setId}/answers`, input)
  return unwrapData<AnswerResult>(response)
}
