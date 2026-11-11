import type { WritingSubmission } from '@/types/writing'
import type { PassageCard } from '@/types/passage'
import type { PhonemeCard } from '@/types/phoneme'
import type { GrammarPointCard } from '@/types/grammar'

// Lesson activities and exercises (docs/foundation-spec.md §4-§5).

export type ActivityType =
  | 'VOCABULARY'
  | 'PRONUNCIATION'
  | 'GRAMMAR'
  | 'LISTENING'
  | 'READING'
  | 'WRITING'
  | 'SPEAKING'
  | 'CHECKPOINT'

export type ActivityStatus = 'NOT_STARTED' | 'IN_PROGRESS' | 'PASSED' | 'FAILED'

export type ExerciseType =
  | 'MULTIPLE_CHOICE'
  | 'TRUE_FALSE'
  | 'FILL_BLANK'
  | 'REORDER'
  | 'ERROR_CORRECTION'
  | 'TRANSFORM'
  | 'TRANSLATE_VI_EN'
  | 'DICTATION'
  | 'LISTEN_GAP_FILL'
  | 'LISTEN_PICK_IMAGE'
  | 'LISTEN_NUMBERS'
  | 'MINIMAL_PAIR'
  | 'IPA_TO_WORD'
  | 'WORD_TO_IPA'
  | 'ENDING_SOUND'
  | 'WORD_STRESS'
  | 'SPEAK'

export type ActivityItem = {
  id: number
  slug: string
  type: ActivityType
  title: string
  instructionsVi: string
  orderIndex: number
  isRequired: boolean
  passScore: number
  xpReward: number
  exerciseCount: number
  status: ActivityStatus
  score: number | null
  bestScore: number | null
  attempts: number
  isLocked: boolean
}

export type LessonActivitiesResponse = {
  lesson: {
    id: number
    title: string
    slug: string
    status: 'LOCKED' | 'UNLOCKED' | 'IN_PROGRESS' | 'COMPLETED'
  }
  progress: {
    passedRequired: number
    totalRequired: number
    nextActivityId: number | null
  }
  activities: ActivityItem[]
}

export type ActivityDetail = ActivityItem & {
  lessonId: number
  content: Record<string, unknown> | null
  grammarPointId: number | null
  passageId: number | null
  scenarioId: number | null
  currentAttempt: number
}

/** Public exercise data — the answer key never reaches the client. */
export type Exercise = {
  id: number
  type: ExerciseType
  prompt: string
  payload: Record<string, unknown>
  points: number
  orderIndex: number
}

/** Answer shapes per exercise type (spec §5.1). */
export type ExerciseAnswer =
  | { optionId: string }
  | { value: boolean }
  | { text: string }
  | { tokens: string[] }
  | { index: number; text: string }
  | { blanks: Record<string, string> }
  | { stressIndex: number }
  | { selfCheck: boolean }

export type AnswerResult = {
  exerciseId: number
  attempt: number
  isCorrect: boolean
  scoreRatio: number
  answer: ExerciseAnswer
  correctAnswer: unknown
  explanationVi: string
  answeredAt: string
  /** MATCH or AI; absent for answers loaded from history. */
  gradedBy?: 'MATCH' | 'AI'
  aiFeedback?: AIFeedback
  /** Set when AI grading was allowed but skipped. */
  aiStatus?: 'LIMIT_REACHED' | 'UNAVAILABLE'
}

export type AIFeedback = {
  suggested: string
  explanationVi: string
}

export type VocabularyCard = {
  id: number
  word: string
  slug: string
  ipa: string
  audioUrl: string
  partOfSpeech: string
  meaningVi: string
  meaningEn: string
  exampleSentence: string
  exampleMeaningVi: string
}

export type ActivityDetailResponse = {
  activity: ActivityDetail
  exercises: Exercise[]
  vocabularies: VocabularyCard[]
  grammarPoint: GrammarPointCard | null
  /** Sound cards of a PRONUNCIATION activity (content.phonemes). */
  phonemes?: PhonemeCard[]
  /** Text/recording of a LISTENING or READING activity. */
  passage?: PassageCard | null
  /** Latest submission of a WRITING activity. */
  writing?: WritingSubmission | null
  /** Role-play brief of a SPEAKING activity with a scenario. */
  scenario?: ScenarioCard | null
  /** Latest finished role-play of that activity. */
  speaking?: SpeakingSummary | null
  answers: AnswerResult[]
}

export type ScenarioCard = {
  id: number
  slug: string
  title: string
  description: string
  role: string
  task: string
  targetPatterns: string[]
  maxTurns: number
  passageId: number | null
}

export type SpeakingSummary = {
  sessionId: number
  score: number | null
  taskCompleted: boolean
  finishedAt: string | null
}

export type CompleteActivityResponse = {
  activity: ActivityItem
  score: number
  passed: boolean
  correctCount: number
  answeredCount: number
  totalCount: number
  xpAwarded: number
  nextActivityId: number | null
  lesson: {
    id: number
    status: 'LOCKED' | 'UNLOCKED' | 'IN_PROGRESS' | 'COMPLETED'
    passedRequired: number
    totalRequired: number
    justCompleted: boolean
    score: number | null
    xpAwarded: number
  }
}

export type SubmitAnswerInput = {
  exerciseId: number
  answer: ExerciseAnswer
}
