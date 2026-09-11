import { api, unwrapData } from '@/api/api'
import type { APIResponse } from '@/types/api'
import type {
  CreateScenarioInput,
  ListScenariosResponse,
  ScenarioResponse,
  SendChatMessageInput,
  SendChatMessageResponse,
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

export async function sendChatMessage(payload: SendChatMessageInput) {
  const response = await api.post<APIResponse<SendChatMessageResponse>>(
    '/ai-conversations/chat',
    payload,
  )

  return unwrapData<SendChatMessageResponse>(response)
}
