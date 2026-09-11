export type ScenarioLevelApi = 'BEGINNER' | 'INTERMEDIATE' | 'ADVANCED'

export type CreateScenarioInput = {
  title: string
  description: string
  situationContext?: string
  level: ScenarioLevelApi
  duration?: string
}

export type ScenarioResponse = {
  id: number
  slug: string
  title: string
  description: string
  situationContext: string
  level: ScenarioLevelApi
  duration: string
  createdAt: string
}

export type ListScenariosResponse = {
  items: ScenarioResponse[]
}

export type ChatRole = 'user' | 'assistant'

export type ChatHistoryItem = {
  role: ChatRole
  content: string
}

export type SendChatMessageInput = {
  message: string
  systemPrompt?: string
  history?: ChatHistoryItem[]
}

export type SendChatMessageResponse = {
  reply: string
}
