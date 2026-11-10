/**
 * The 40 sound cards (docs/foundation-spec.md §3.3), as written in our
 * vocabulary data (American English, Oxford NAmE style).
 */
export const PHONEME_SYMBOLS = [
  'iː', 'ɪ', 'e', 'æ', 'ɑː', 'ɔː', 'ʊ', 'uː', 'ʌ', 'ɜːr', 'ə',
  'eɪ', 'aɪ', 'ɔɪ', 'aʊ', 'oʊ',
  'p', 'b', 't', 'd', 'k', 'ɡ', 'f', 'v', 'θ', 'ð', 's', 'z', 'ʃ', 'ʒ', 'h', 'tʃ', 'dʒ',
  'm', 'n', 'ŋ', 'l', 'r', 'w', 'j',
] as const

/** Other spellings found in dictionaries, mapped to their sound card. */
const ALIASES: Record<string, string> = {
  i: 'iː', // happy /ˈhæpi/
  u: 'uː', // July /dʒuˈlaɪ/
  ɑ: 'ɑː',
  ɒ: 'ɑː', // British "hot"
  ɔ: 'ɔː',
  ɜː: 'ɜːr', // British "work"
  ɜ: 'ɜːr',
  ɝ: 'ɜːr',
  ɚ: 'ə',
  əʊ: 'oʊ', // British "go"
  g: 'ɡ',
  ʧ: 'tʃ',
  ʤ: 'dʒ',
}

const MATCHES = [...PHONEME_SYMBOLS, ...Object.keys(ALIASES)].sort((a, b) => b.length - a.length)

/** Characters that are part of a transcription but not a sound. */
const MARKS = new Set(['/', 'ˈ', 'ˌ', '.', ' ', '(', ')', '[', ']', 'ː', '-', ','])

export type IpaSegment = { text: string; symbol?: string }

/**
 * Splits an IPA transcription into segments, linking each sound to its card
 * (longest match first, so "tʃ" wins over "t"). Returns null when the text is
 * not a real transcription (e.g. a word typed in the IPA field).
 */
export function splitIpa(ipa: string): IpaSegment[] | null {
  const segments: IpaSegment[] = []
  let rest = ipa
  while (rest.length > 0) {
    const match = MATCHES.find((candidate) => rest.startsWith(candidate))
    if (match) {
      segments.push({ text: match, symbol: ALIASES[match] ?? match })
      rest = rest.slice(match.length)
      continue
    }
    const char = rest[0]
    if (!MARKS.has(char)) {
      return null
    }
    const last = segments.at(-1)
    if (last && !last.symbol) {
      last.text += char
    } else {
      segments.push({ text: char })
    }
    rest = rest.slice(1)
  }
  return segments
}

/** URL of a sound card page. */
export function phonemePath(symbol: string) {
  return `/pronunciation/${encodeURIComponent(symbol)}`
}
