export type ContentImportKind = 'GRAMMAR_POINT' | 'TOPIC' | 'LESSON' | 'ACTIVITY' | 'EXERCISE'
export type ContentImportAction = 'CREATE' | 'UPDATE' | 'SKIP'

export type ContentImportItem = {
  kind: ContentImportKind
  path: string
  title: string
  action: ContentImportAction
  status?: string
  errors: string[]
}

export type ContentImportResult = {
  committed: boolean
  summary: {
    topics: number
    lessons: number
    activities: number
    exercises: number
    created: number
    updated: number
    invalid: number
  }
  bundleErrors: string[]
  warnings: string[]
  items: ContentImportItem[]
}
