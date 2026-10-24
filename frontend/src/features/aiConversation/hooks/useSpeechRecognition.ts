import { useEffect, useRef, useState } from 'react'

// TypeScript's DOM lib ships the speech event types but not the recognizer.
type SpeechRecognitionInstance = {
  lang: string
  continuous: boolean
  interimResults: boolean
  onresult: ((event: SpeechRecognitionEvent) => void) | null
  onerror: ((event: SpeechRecognitionErrorEvent) => void) | null
  onend: (() => void) | null
  start: () => void
  abort: () => void
}

type SpeechRecognitionConstructor = new () => SpeechRecognitionInstance

type WindowWithSpeechRecognition = Window & {
  SpeechRecognition?: SpeechRecognitionConstructor
  webkitSpeechRecognition?: SpeechRecognitionConstructor
}

const errorMessages: Partial<Record<SpeechRecognitionErrorCode, string>> = {
  'not-allowed':
    'Microphone access is blocked. Allow it in your browser settings.',
  'service-not-allowed':
    'Microphone access is blocked. Allow it in your browser settings.',
  'audio-capture': 'No microphone was found.',
  network: 'Speech recognition needs an internet connection.',
}

export function useSpeechRecognition(onError: (message: string) => void) {
  const recognitionRef = useRef<SpeechRecognitionInstance | null>(null)
  const [isListening, setIsListening] = useState(false)

  useEffect(
    () => () => {
      const recognition = recognitionRef.current
      recognitionRef.current = null
      recognition?.abort()
    },
    [],
  )

  function start(onTranscript: (transcript: string) => void) {
    const { SpeechRecognition, webkitSpeechRecognition } =
      window as WindowWithSpeechRecognition
    const Recognition = SpeechRecognition ?? webkitSpeechRecognition
    if (!Recognition) {
      onError(
        'Voice input is not supported in this browser. Try Chrome or Edge.',
      )
      return
    }

    const recognition = new Recognition()
    recognition.lang = 'en-US'
    recognition.continuous = true
    recognition.interimResults = true

    // Final phrases are kept here because each auto-restart below begins a
    // new session whose results start again from index 0.
    const finalPhrases: string[] = []
    recognition.onresult = (event) => {
      const interimPhrases: string[] = []
      for (let i = event.resultIndex; i < event.results.length; i++) {
        const result = event.results[i]
        const phrase = result[0].transcript.trim()
        if (result.isFinal) {
          finalPhrases.push(phrase)
        } else {
          interimPhrases.push(phrase)
        }
      }
      onTranscript(
        [...finalPhrases, ...interimPhrases].filter(Boolean).join(' '),
      )
    }

    recognition.onerror = (event) => {
      // Silence is not a failure; onend simply restarts listening.
      if (event.error === 'no-speech' || event.error === 'aborted') {
        return
      }
      recognitionRef.current = null
      setIsListening(false)
      onError(
        errorMessages[event.error] ?? 'Voice input stopped unexpectedly.',
      )
    }

    // Browsers end a session after a pause; keep listening until stop().
    recognition.onend = () => {
      if (recognitionRef.current === recognition) {
        recognition.start()
      }
    }

    recognitionRef.current = recognition
    recognition.start()
    setIsListening(true)
  }

  function stop() {
    const recognition = recognitionRef.current
    recognitionRef.current = null
    if (recognition) {
      // Drop in-flight results so the input is not refilled after sending.
      recognition.onresult = null
      recognition.abort()
    }
    setIsListening(false)
  }

  return { isListening, start, stop }
}
