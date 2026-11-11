export type PronunciationSource = 'WORD' | 'SHADOWING' | 'READ_ALOUD' | 'EXERCISE'

export type PronunciationStatus = {
  enabled: boolean
  maxSeconds: number
  dailyLimit: number
}

export type PhonemeScore = {
  /** Phoneme as Azure wrote it. */
  azure: string
  /** The sound in our IPA notation. */
  display: string
  /** Sound cards it counts for (link to /pronunciation/:symbol). */
  symbols: string[]
  accuracy: number
}

export type WordScore = {
  word: string
  accuracy: number
  errorType: 'None' | 'Mispronunciation' | 'Omission' | 'Insertion' | string
  phonemes: PhonemeScore[]
}

export type PronunciationResult = {
  attemptId: number
  referenceText: string
  recognizedText: string
  accuracy: number
  fluency: number
  completeness: number
  pronunciation: number
  words: WordScore[]
  minAccuracy?: number
  passed?: boolean
}

export type AssessInput = {
  audio: Blob
  referenceText: string
  exerciseId?: number
  source: PronunciationSource
}
