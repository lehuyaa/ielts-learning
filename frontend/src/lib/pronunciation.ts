let currentAudio: HTMLAudioElement | null = null

/**
 * Plays a word's recorded pronunciation. Falls back to the browser's speech
 * synthesis when the word has no audio file or the file cannot be played.
 */
export function playPronunciation(word: string, audioUrl?: string | null) {
  stopPronunciation()

  if (!audioUrl) {
    speakWord(word)
    return
  }

  const audio = new Audio(audioUrl)
  currentAudio = audio
  audio.play().catch(() => {
    if (currentAudio === audio) {
      speakWord(word)
    }
  })
}

export function stopPronunciation() {
  currentAudio?.pause()
  currentAudio = null
  if ('speechSynthesis' in window) {
    window.speechSynthesis.cancel()
  }
}

function speakWord(word: string) {
  if (!('speechSynthesis' in window)) {
    return
  }

  const utterance = new SpeechSynthesisUtterance(word)
  utterance.lang = 'en-US'
  // Slightly slower than normal so beginners can hear every sound.
  utterance.rate = 0.85
  window.speechSynthesis.speak(utterance)
}
