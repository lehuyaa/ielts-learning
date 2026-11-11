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
  isSystem: boolean
  cefrLevel: string
  role: string
  task: string
  targetPatterns: string[]
  maxTurns: number
  openingLine: string
  group: string
  passageId: number | null
  /** Finished role-plays of this scenario and their best score. */
  sessions: number
  bestScore: number | null
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

export type ScenarioBrief = {
  id: number
  slug: string
  title: string
  description: string
  cefrLevel: string
  role: string
  task: string
  targetPatterns: string[]
  maxTurns: number
  openingLine: string
  isSystem: boolean
  passageId: number | null
  group: string
}

export type SessionMessage = { id: number; role: ChatRole; content: string; createdAt: string }

export type SpeakingReport = {
  taskCompleted: boolean
  scores: { grammar: number; vocabulary: number; fluency: number; taskAchievement: number }
  overall: number
  mistakes: { said: string; correct: string; skillTag: string; explanationVi: string }[]
  goodPhrases: string[]
  newWords: string[]
  tipVi: string
  xpAwarded: number
}

export type SpeakingSession = {
  id: number
  status: 'IN_PROGRESS' | 'FINISHED' | 'ABANDONED'
  activityId: number | null
  turnCount: number
  maxTurns: number
  hintCount: number
  hintsAllowed: boolean
  score: number | null
  taskCompleted: boolean
  report: SpeakingReport | null
  scenario: ScenarioBrief
  messages: SessionMessage[]
  startedAt: string
  finishedAt: string | null
}

export type SessionTurn = {
  reply: SessionMessage
  turnCount: number
  maxTurns: number
  shouldFinish: boolean
}

export type SpeakingHint = { hintVi: string; sayThis: string; hintCount: number }
