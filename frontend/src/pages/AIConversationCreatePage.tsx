import { ArrowLeft, HelpCircle, MessageSquare, Pencil } from 'lucide-react'
import { useState } from 'react'
import { useNavigate } from 'react-router-dom'

import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import {
  addCustomScenario,
  type ScenarioLevel,
} from '@/features/aiConversation/scenarios'
import { cn } from '@/lib/utils'

const levels: ScenarioLevel[] = ['Beginner', 'Intermediate', 'Advanced']

const levelSelectedClasses: Record<ScenarioLevel, string> = {
  Beginner: 'border-emerald-300 bg-emerald-50 text-emerald-700',
  Intermediate: 'border-amber-300 bg-amber-50 text-amber-700',
  Advanced: 'border-red-300 bg-red-50 text-red-700',
}

export function AIConversationCreatePage() {
  const navigate = useNavigate()
  const [title, setTitle] = useState('')
  const [description, setDescription] = useState('')
  const [level, setLevel] = useState<ScenarioLevel>('Beginner')
  const [situationContext, setSituationContext] = useState('')

  const canSubmit = title.trim().length > 0 && description.trim().length > 0

  return (
    <div className="min-h-screen bg-background">
      <div className="mx-auto max-w-2xl px-4 py-8">
        <div className="mb-8 flex items-start gap-4">
          <button
            aria-label="Go back"
            className="grid size-10 shrink-0 place-items-center rounded-full bg-muted text-foreground transition-colors hover:bg-muted/70"
            onClick={() => navigate('/ai-conversation')}
            type="button"
          >
            <ArrowLeft className="size-5" aria-hidden="true" />
          </button>
          <div>
            <h1 className="text-2xl font-bold tracking-normal text-foreground">
              Create Scenario
            </h1>
            <p className="mt-0.5 text-sm text-muted-foreground">
              Design your own conversation practice.
            </p>
          </div>
        </div>

        <form
          className="space-y-6"
          onSubmit={(event) => {
            event.preventDefault()
            if (!canSubmit) {
              return
            }

            const scenario = addCustomScenario({
              title: title.trim(),
              description: description.trim(),
              duration: '5-10 min',
              level,
              situationContext,
            })
            navigate(`/ai-conversation/${scenario.slug}`)
          }}
        >
          <div className="relative w-fit">
            <div className="flex size-16 items-center justify-center rounded-2xl bg-indigo-50">
              <MessageSquare className="size-7 text-primary" aria-hidden="true" />
            </div>
            <span className="absolute -bottom-1 -right-1 grid size-6 place-items-center rounded-full bg-primary text-white ring-2 ring-background">
              <Pencil className="size-3" aria-hidden="true" />
            </span>
          </div>

          <div className="space-y-1.5">
            <label className="text-sm font-semibold text-foreground" htmlFor="create-scenario-title">
              Scenario title <span className="text-red-500">*</span>
            </label>
            <Input
              id="create-scenario-title"
              onChange={(event) => setTitle(event.target.value)}
              placeholder="e.g. Doctor's Appointment"
              value={title}
            />
          </div>

          <div className="space-y-1.5">
            <label
              className="text-sm font-semibold text-foreground"
              htmlFor="create-scenario-description"
            >
              Short description <span className="text-red-500">*</span>
            </label>
            <Input
              id="create-scenario-description"
              onChange={(event) => setDescription(event.target.value)}
              placeholder="e.g. Practice explaining symptoms and asking the doctor questions."
              value={description}
            />
          </div>

          <div className="space-y-2">
            <span className="text-sm font-semibold text-foreground">
              Difficulty
            </span>
            <div className="grid grid-cols-3 gap-3">
              {levels.map((option) => (
                <button
                  className={cn(
                    'rounded-xl border px-3 py-3 text-sm font-medium transition-colors',
                    level === option
                      ? levelSelectedClasses[option]
                      : 'border-border text-muted-foreground hover:bg-muted/50',
                  )}
                  key={option}
                  onClick={() => setLevel(option)}
                  type="button"
                >
                  {option}
                </button>
              ))}
            </div>
          </div>

          <div className="space-y-1.5">
            <div className="flex items-center justify-between">
              <label
                className="text-sm font-semibold text-foreground"
                htmlFor="create-scenario-context"
              >
                Situation context
              </label>
              <span className="text-xs text-muted-foreground">Optional</span>
            </div>
            <textarea
              className="w-full resize-none rounded-xl border border-border bg-white px-4 py-3 text-sm text-foreground outline-none transition-colors placeholder:text-muted-foreground focus:border-primary"
              id="create-scenario-context"
              onChange={(event) => setSituationContext(event.target.value)}
              placeholder={
                'Describe the setting in detail so the AI plays the right role.\n\nExample: "You are a patient at a busy clinic. The AI plays the role of a doctor. Start by asking about symptoms."'
              }
              rows={4}
              value={situationContext}
            />
            <p className="text-xs text-muted-foreground">
              The AI will use this context to guide the conversation.
            </p>
          </div>

          <Button className="w-full rounded-xl py-6 text-base" disabled={!canSubmit} type="submit">
            Create Scenario
          </Button>
        </form>
      </div>

      <button
        aria-label="Help"
        className="fixed bottom-6 right-6 grid size-10 place-items-center rounded-full bg-foreground text-white shadow-md"
        type="button"
      >
        <HelpCircle className="size-4" aria-hidden="true" />
      </button>
    </div>
  )
}
