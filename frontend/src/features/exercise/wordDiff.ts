export type DiffToken = { word: string; kind: 'ok' | 'missing' | 'extra' }

function normalize(word: string) {
  return word
    .toLowerCase()
    .replace(/[’‘]/g, "'")
    .replace(/^[^a-z0-9']+|[^a-z0-9']+$/g, '')
}

/**
 * Word-level diff of a dictation answer against the correct sentence (longest
 * common subsequence): correct words, words the learner missed, and words
 * they wrote that are not in the sentence.
 */
export function diffWords(answer: string, correct: string): DiffToken[] {
  const typed = answer.split(/\s+/).filter(Boolean)
  const expected = correct.split(/\s+/).filter(Boolean)
  const a = typed.map(normalize)
  const b = expected.map(normalize)

  const table: number[][] = Array.from({ length: a.length + 1 }, () => new Array<number>(b.length + 1).fill(0))
  for (let i = a.length - 1; i >= 0; i--) {
    for (let j = b.length - 1; j >= 0; j--) {
      table[i][j] = a[i] === b[j] ? table[i + 1][j + 1] + 1 : Math.max(table[i + 1][j], table[i][j + 1])
    }
  }

  const out: DiffToken[] = []
  let i = 0
  let j = 0
  while (i < a.length && j < b.length) {
    if (a[i] === b[j]) {
      out.push({ word: expected[j], kind: 'ok' })
      i++
      j++
    } else if (table[i + 1][j] >= table[i][j + 1]) {
      out.push({ word: typed[i], kind: 'extra' })
      i++
    } else {
      out.push({ word: expected[j], kind: 'missing' })
      j++
    }
  }
  for (; i < a.length; i++) out.push({ word: typed[i], kind: 'extra' })
  for (; j < b.length; j++) out.push({ word: expected[j], kind: 'missing' })
  return out
}
