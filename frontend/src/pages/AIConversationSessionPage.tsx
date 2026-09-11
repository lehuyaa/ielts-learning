import { ArrowLeft, HelpCircle, Mic, Send } from 'lucide-react'
import { useEffect, useRef, useState } from 'react'
import { Link, Navigate, useNavigate, useParams } from 'react-router-dom'

import { getApiErrorMessage } from '@/api/api'
import { useScenarios } from '@/features/aiConversation/hooks/useScenarios'
import { useSendChatMessage } from '@/features/aiConversation/hooks/useSendChatMessage'
import { mapScenarioResponseToScenario } from '@/features/aiConversation/scenarios'
import { cn } from '@/lib/utils'

type Message = {
  id: number
  sender: 'ai' | 'user'
  text: string
  time: string
}

const MAX_HISTORY_MESSAGES = 20

function formatElapsed(totalSeconds: number) {
  const minutes = Math.floor(totalSeconds / 60)
  const seconds = totalSeconds % 60
  return `${minutes}:${String(seconds).padStart(2, '0')}`
}

export function AIConversationSessionPage() {
  const { scenarioSlug } = useParams()
  const navigate = useNavigate()
  const scenariosQuery = useScenarios()

  const scenarioResponse = scenariosQuery.data?.items.find(
    (item) => item.slug === scenarioSlug,
  )
  const scenario = scenarioResponse
    ? mapScenarioResponseToScenario(scenarioResponse)
    : undefined
  const isResolvingScenario = !scenario && scenariosQuery.isLoading

  const [elapsedSeconds, setElapsedSeconds] = useState(0)
  const [messages, setMessages] = useState<Message[]>([])
  const [inputValue, setInputValue] = useState('')
  const [isRecording, setIsRecording] = useState(false)
  const [isAiTyping, setIsAiTyping] = useState(false)

  const nextMessageId = useRef(1)
  const scrollRef = useRef<HTMLDivElement>(null)
  const sendChatMessageMutation = useSendChatMessage()

  useEffect(() => {
    const interval = setInterval(() => {
      setElapsedSeconds((seconds) => seconds + 1)
    }, 1000)
    return () => clearInterval(interval)
  }, [])

  useEffect(() => {
    if (!scenario) {
      return
    }

    setMessages([
      {
        id: nextMessageId.current++,
        sender: 'ai',
        text: scenario.openingLine,
        time: '0:00',
      },
    ])
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [scenario?.slug])

  useEffect(() => {
    scrollRef.current?.scrollTo({ top: scrollRef.current.scrollHeight })
  }, [messages, isAiTyping])

  if (!scenario) {
    if (isResolvingScenario) {
      return (
        <div className="grid min-h-screen place-items-center bg-background text-sm text-muted-foreground">
          Loading scenario…
        </div>
      )
    }

    return <Navigate replace to="/ai-conversation" />
  }

  function sendMessage() {
    const text = inputValue.trim()
    if (!text || !scenario || sendChatMessageMutation.isPending) {
      return
    }

    const history = messages.slice(-MAX_HISTORY_MESSAGES).map((message) => ({
      role: message.sender === 'user' ? ('user' as const) : ('assistant' as const),
      content: message.text,
    }))

    setMessages((current) => [
      ...current,
      {
        id: nextMessageId.current++,
        sender: 'user',
        text,
        time: formatElapsed(elapsedSeconds),
      },
    ])
    setInputValue('')
    setIsAiTyping(true)

    sendChatMessageMutation.mutate(
      {
        message: text,
        systemPrompt: scenario.situationContext,
        history,
      },
      {
        onSuccess: (result) => {
          setMessages((current) => [
            ...current,
            {
              id: nextMessageId.current++,
              sender: 'ai',
              text: result.reply,
              time: formatElapsed(elapsedSeconds),
            },
          ])
          setIsAiTyping(false)
        },
        onError: (error) => {
          setMessages((current) => [
            ...current,
            {
              id: nextMessageId.current++,
              sender: 'ai',
              text: getApiErrorMessage(
                error,
                "Sorry, I couldn't respond right now. Please try again.",
              ),
              time: formatElapsed(elapsedSeconds),
            },
          ])
          setIsAiTyping(false)
        },
      },
    )
  }

  return (
    <div className="flex min-h-screen flex-col bg-background">
      <header className="flex items-center justify-between border-b border-border bg-white px-4 py-3">
        <Link
          className="flex items-center gap-1.5 text-sm text-muted-foreground transition-colors hover:text-foreground"
          to="/ai-conversation"
        >
          <ArrowLeft className="size-4" aria-hidden="true" />
          Back
        </Link>

        <div className="text-center">
          <div className="text-sm font-bold text-foreground">
            {scenario.title}
          </div>
          <div className="text-xs text-muted-foreground">
            {formatElapsed(elapsedSeconds)}
          </div>
        </div>

        <button
          className="text-sm font-medium text-red-500 transition-colors hover:text-red-600"
          onClick={() => navigate('/ai-conversation')}
          type="button"
        >
          End
        </button>
      </header>

      <div className="border-b border-border bg-white px-4 py-6 text-center">
        <div className="mx-auto flex size-14 items-center justify-center rounded-full bg-indigo-50">
          <scenario.icon className="size-6 text-primary" aria-hidden="true" />
        </div>
        <div className="mt-2 text-sm font-bold text-foreground">AI Tutor</div>
        <div className="text-xs text-muted-foreground">
          {isRecording ? 'Listening…' : isAiTyping ? 'Typing…' : 'Ready'}
        </div>
      </div>

      <div
        className="flex-1 space-y-4 overflow-y-auto px-4 py-4"
        ref={scrollRef}
      >
        {messages.map((message) => (
          <div
            className={cn(
              'flex',
              message.sender === 'user' ? 'justify-end' : 'justify-start',
            )}
            key={message.id}
          >
            <div
              className={cn(
                'max-w-[75%] rounded-2xl px-4 py-2.5 text-sm',
                message.sender === 'user'
                  ? 'rounded-br-md bg-primary text-white'
                  : 'rounded-bl-md border border-border bg-white text-foreground shadow-sm',
              )}
            >
              {message.text}
              <div
                className={cn(
                  'mt-1 text-[11px]',
                  message.sender === 'user'
                    ? 'text-white/70'
                    : 'text-muted-foreground',
                )}
              >
                {message.time}
              </div>
            </div>
          </div>
        ))}

        {isAiTyping ? (
          <div className="flex justify-start">
            <div className="rounded-2xl rounded-bl-md border border-border bg-white px-4 py-2.5 text-sm text-muted-foreground shadow-sm">
              <span className="inline-flex gap-1">
                <span className="size-1.5 animate-bounce rounded-full bg-muted-foreground/50 [animation-delay:-0.3s]" />
                <span className="size-1.5 animate-bounce rounded-full bg-muted-foreground/50 [animation-delay:-0.15s]" />
                <span className="size-1.5 animate-bounce rounded-full bg-muted-foreground/50" />
              </span>
            </div>
          </div>
        ) : null}
      </div>

      <div className="relative border-t border-border bg-white px-4 pb-4 pt-6">
        <button
          aria-label="Help"
          className="absolute -top-4 right-4 grid size-8 place-items-center rounded-full bg-foreground text-white shadow-md"
          type="button"
        >
          <HelpCircle className="size-4" aria-hidden="true" />
        </button>

        <div className="flex flex-col items-center gap-2 pb-4">
          <button
            aria-label={isRecording ? 'Stop recording' : 'Start recording'}
            className={cn(
              'grid size-14 place-items-center rounded-full text-white transition-colors',
              isRecording ? 'bg-red-500' : 'bg-primary hover:bg-primary/90',
            )}
            onClick={() => setIsRecording((current) => !current)}
            type="button"
          >
            <Mic className="size-6" aria-hidden="true" />
          </button>

          <div className="flex items-center gap-1">
            {Array.from({ length: 7 }).map((_, index) => (
              <span
                className={cn(
                  'h-1.5 w-1.5 rounded-full bg-primary/30',
                  isRecording && 'animate-pulse',
                )}
                key={index}
                style={{ animationDelay: `${index * 100}ms` }}
              />
            ))}
          </div>

          <div className="text-xs text-muted-foreground">
            {isRecording ? 'Listening…' : 'Tap to speak'}
          </div>
        </div>

        <form
          className="flex items-center gap-2"
          onSubmit={(event) => {
            event.preventDefault()
            sendMessage()
          }}
        >
          <input
            className="h-11 flex-1 rounded-full border border-border bg-muted/30 px-4 text-sm text-foreground outline-none transition-colors focus:border-primary"
            onChange={(event) => setInputValue(event.target.value)}
            placeholder="Or type your message..."
            value={inputValue}
          />
          <button
            aria-label="Send message"
            className="grid size-11 shrink-0 place-items-center rounded-full bg-primary text-white transition-opacity disabled:opacity-40"
            disabled={!inputValue.trim() || sendChatMessageMutation.isPending}
            type="submit"
          >
            <Send className="size-4" aria-hidden="true" />
          </button>
        </form>
      </div>
    </div>
  )
}
