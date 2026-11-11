import { api, unwrapData } from '@/api/api'
import type { APIResponse } from '@/types/api'
import type { AddFlashcardResult, PassageCard, WordLookup } from '@/types/passage'

export async function getPassage(passageId: number) {
  const response = await api.get<APIResponse<PassageCard>>(`/passages/${passageId}`)
  return unwrapData<PassageCard>(response)
}

export async function lookupWord(word: string, passageId?: number) {
  const response = await api.get<APIResponse<WordLookup>>('/words/lookup', {
    params: passageId ? { word, passageId } : { word },
  })
  return unwrapData<WordLookup>(response)
}

export async function addWordToFlashcards(vocabularyId: number) {
  const response = await api.post<APIResponse<AddFlashcardResult>>(`/words/${vocabularyId}/flashcard`)
  return unwrapData<AddFlashcardResult>(response)
}
