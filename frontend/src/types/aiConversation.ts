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
