export type FoundationSkill = {
  key: 'VOCABULARY' | 'GRAMMAR' | 'SOUNDS' | 'LISTENING' | 'READING' | 'WRITING' | 'SPEAKING'
  label: string
  value: number
  practised: boolean
  done: number
  total: number
  detail: string
  link: string
}

export type FoundationGoal = {
  key: string
  title: string
  done: boolean
  progress: number
  detail: string
  hint: string
  link: string
}

export type FoundationSummary = {
  level: string
  title: string
  lessons: { completed: number; total: number; target: number; checkpointsPassed: number; checkpointsTotal: number }
  streak: { current: number; longest: number }
  skills: FoundationSkill[]
  checklist: FoundationGoal[]
  goalsDone: number
  nextLesson: { id: number; title: string; topicTitle: string } | null
}
