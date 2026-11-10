import { api, unwrapData } from '@/api/api'
import type { ExerciseAnswer } from '@/types/activity'
import type { APIResponse } from '@/types/api'
import type {
  GrammarDetailResponse,
  GrammarListResponse,
  GrammarPracticeResponse,
  PracticeAnswerResult,
} from '@/types/grammar'

export async function getGrammarPoints(level?: string) {
  const response = await api.get<APIResponse<GrammarListResponse>>('/grammar-points', {
    params: level ? { level } : undefined,
  })
  return unwrapData<GrammarListResponse>(response)
}

export async function getGrammarPoint(slug: string) {
  const response = await api.get<APIResponse<GrammarDetailResponse>>(`/grammar-points/${slug}`)
  return unwrapData<GrammarDetailResponse>(response)
}

export async function getGrammarPractice(slug: string) {
  const response = await api.get<APIResponse<GrammarPracticeResponse>>(
    `/grammar-points/${slug}/practice`,
  )
  return unwrapData<GrammarPracticeResponse>(response)
}

export async function answerGrammarPractice(
  slug: string,
  input: { exerciseId: number; answer: ExerciseAnswer },
) {
  const response = await api.post<APIResponse<PracticeAnswerResult>>(
    `/grammar-points/${slug}/practice/answers`,
    input,
  )
  return unwrapData<PracticeAnswerResult>(response)
}
