import { api, unwrapData } from '@/api/api'
import type { APIResponse } from '@/types/api'
import type {
  CreateScenarioInput,
  ListMessagesResponse,
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
