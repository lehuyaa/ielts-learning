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

export type SendChatMessageInput = {
  message: string
}

export type SendChatMessageResponse = {
  reply: string
}

export type ChatMessageResponse = {
  id: number
  role: ChatRole
  content: string
  createdAt: string
}

export type ListMessagesResponse = {
  items: ChatMessageResponse[]
}
