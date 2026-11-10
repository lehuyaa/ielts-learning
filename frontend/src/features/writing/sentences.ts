/** Counts sentences the same way as the backend (ends with . ! ? and has a letter). */
export function countSentences(text: string) {
  return text
    .trim()
    .split(/[.!?]+(?:\s+|$)/)
    .filter((part) => /\p{L}/u.test(part)).length
}
