export type PassageGenre =
  | 'DIALOGUE'
  | 'MESSAGE'
  | 'EMAIL'
  | 'NOTICE'
  | 'MENU'
  | 'SCHEDULE'
  | 'ARTICLE'
  | 'PROFILE'
  | 'MONOLOGUE'

export type PassageSegment = {
  index: number
  speaker?: string
  text: string
  /** Times inside the passage audio (Whisper imports). */
  start?: number
  end?: number
  /** Own audio file (generated dialogues). */
  audioUrl?: string
}

export type GlossaryEntry = {
  word: string
  vocabularySlug?: string
  meaningVi?: string
  ipa?: string
}

export type PassageCard = {
  id: number
  slug: string
  cefrLevel: string
  genre: PassageGenre
  title: string
  bodyText: string
  audioUrl: string | null
  segments: PassageSegment[]
  glossary: GlossaryEntry[]
  source: string
  wordCount: number
}

export type WordLookup = {
  word: string
  found: boolean
  note: string
  ipa: string
  vocabulary: {
    id: number
    word: string
    slug: string
    ipa: string
    partOfSpeech: string
    meaningVi: string
    audioUrl: string
    inFlashcards: boolean
  } | null
}

export type AddFlashcardResult = { vocabularyId: number; added: boolean; status: string }

export type TranscriptMode = 'HIDDEN' | 'BLURRED' | 'AFTER_ANSWER'
