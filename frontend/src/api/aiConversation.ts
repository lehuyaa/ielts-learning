import { api, unwrapData } from '@/api/api'
import type { APIResponse } from '@/types/api'
import type {
  CreateScenarioInput,
  ListMessagesResponse,
  ListScenariosResponse,
  ScenarioResponse,
  SendChatMessageInput,
  SendChatMessageResponse,
  SessionTurn,
  SpeakingHint,
  SpeakingSession,
} from '@/types/aiConversation'

export async function createScenario(payload: CreateScenarioInput) {
  const response = await api.post<APIResponse<ScenarioResponse>>(
    '/ai-conversations/scenarios',
    payload,
  )

  return unwrapData<ScenarioResponse>(response)
}

export async function listScenarios() {
  const response = await api.get<APIResponse<ListScenariosResponse>>(
    '/ai-conversations/scenarios',
  )

  return unwrapData<ListScenariosResponse>(response)
}

export async function listChatMessages(scenarioSlug: string) {
  const response = await api.get<APIResponse<ListMessagesResponse>>(
    `/ai-conversations/scenarios/${scenarioSlug}/messages`,
  )

  return unwrapData<ListMessagesResponse>(response)
}

export async function sendChatMessage(
  scenarioSlug: string,
  payload: SendChatMessageInput,
) {
  const response = await api.post<APIResponse<SendChatMessageResponse>>(
    `/ai-conversations/scenarios/${scenarioSlug}/messages`,
    payload,
  )

  return unwrapData<SendChatMessageResponse>(response)
}

// AI calls can take a few seconds (reports longer).
const slow = { timeout: 90_000 }

export async function startSession(input: { scenarioSlug: string; activityId?: number; fresh?: boolean }) {
  const response = await api.post<APIResponse<SpeakingSession>>('/ai-conversations/sessions', input)
  return unwrapData<SpeakingSession>(response)
}

export async function getSession(sessionId: number) {
  const response = await api.get<APIResponse<SpeakingSession>>(`/ai-conversations/sessions/${sessionId}`)
  return unwrapData<SpeakingSession>(response)
}

export async function sendSessionMessage(sessionId: number, message: string) {
  const response = await api.post<APIResponse<SessionTurn>>(
    `/ai-conversations/sessions/${sessionId}/messages`,
    { message },
    slow,
  )
  return unwrapData<SessionTurn>(response)
}

export async function getSessionHint(sessionId: number) {
  const response = await api.post<APIResponse<SpeakingHint>>(`/ai-conversations/sessions/${sessionId}/hint`, null, slow)
  return unwrapData<SpeakingHint>(response)
}

export async function finishSession(sessionId: number) {
  const response = await api.post<APIResponse<SpeakingSession>>(`/ai-conversations/sessions/${sessionId}/finish`, null, slow)
  return unwrapData<SpeakingSession>(response)
}
