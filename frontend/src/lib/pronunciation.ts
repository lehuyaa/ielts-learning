import { mediaUrl } from '@/lib/media'

let currentAudio: HTMLAudioElement | null = null

type PlayOptions = {
  /** Playback speed (1 = normal). */
  rate?: number
  /** Called when playback (or the speech fallback) ends. */
  onEnd?: () => void
}

/**
 * Plays a word's or sentence's recorded audio. Falls back to the browser's
 * speech synthesis when there is no audio file or it cannot be played.
 */
export function playPronunciation(word: string, audioUrl?: string | null, options: PlayOptions = {}) {
  stopPronunciation()
  const rate = options.rate ?? 1

  if (!audioUrl) {
    speakWord(word, rate, options.onEnd)
    return
  }

  const audio = new Audio(mediaUrl(audioUrl))
  audio.playbackRate = rate
  currentAudio = audio
  audio.onended = () => {
    if (currentAudio === audio) {
      currentAudio = null
      options.onEnd?.()
    }
  }
  audio.play().catch(() => {
    if (currentAudio === audio) {
      speakWord(word, rate, options.onEnd)
    }
  })
}

export function stopPronunciation() {
  if (currentAudio) {
    currentAudio.onended = null
    currentAudio.pause()
  }
  currentAudio = null
  if ('speechSynthesis' in window) {
    window.speechSynthesis.cancel()
  }
}

/** Reads text aloud with the browser voice (no audio file). */
export function speakWord(word: string, rate = 1, onEnd?: () => void) {
  if (!('speechSynthesis' in window) || !word) {
    onEnd?.()
    return
  }

  const utterance = new SpeechSynthesisUtterance(word)
  utterance.lang = 'en-US'
  // Slightly slower than normal so beginners can hear every sound.
  utterance.rate = 0.85 * rate
  if (onEnd) {
    utterance.onend = () => onEnd()
  }
  window.speechSynthesis.speak(utterance)
}
