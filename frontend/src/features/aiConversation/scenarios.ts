import { Sparkles, type LucideIcon } from 'lucide-react'

import type {
  ScenarioLevelApi,
  ScenarioResponse,
} from '@/types/aiConversation'

export type ScenarioLevel = 'Beginner' | 'Intermediate' | 'Advanced'

export type Scenario = {
  slug: string
  title: string
  description: string
  duration: string
  level: ScenarioLevel
  icon: LucideIcon
  openingLine: string
  situationContext: string
}

const apiLevelToLocal: Record<ScenarioLevelApi, ScenarioLevel> = {
  BEGINNER: 'Beginner',
  INTERMEDIATE: 'Intermediate',
  ADVANCED: 'Advanced',
}

export function getScenarioBySlug(slug: string | undefined) {
  return scenarios.find((scenario) => scenario.slug === slug)
}

// In-memory store for the authenticated user's scenarios (the default
// starter set plus any they created). Hydrated from the backend — see
// replaceScenariosFromApi — each time the AI Conversation list page loads.
let scenarios: Scenario[] = []
const listeners = new Set<() => void>()

function notifyListeners() {
  listeners.forEach((listener) => listener())
}

export function getScenarios() {
  return scenarios
}

export function subscribeToScenarios(listener: () => void) {
  listeners.add(listener)
  return () => listeners.delete(listener)
}

export type NewScenarioInput = {
  title: string
  description: string
  duration: string
  level: ScenarioLevel
  situationContext?: string
  /** Use a server-assigned slug (e.g. after persisting via the backend API) instead of generating one locally. */
  slug?: string
}

function buildScenario(input: {
  slug: string
  title: string
  description: string
  duration: string
  level: ScenarioLevel
  situationContext?: string
}): Scenario {
  const context = input.situationContext?.trim() ?? ''

  return {
    slug: input.slug,
    title: input.title,
    description: input.description,
    duration: input.duration || '5-10 min',
    level: input.level,
    icon: Sparkles,
    openingLine: context
      ? `Let's begin. ${context}`
      : `Let's practice "${input.title}". ${input.description}`,
    situationContext: context,
  }
}

export function addScenario(input: NewScenarioInput): Scenario {
  const slug =
    input.slug ?? `custom-${slugify(input.title)}-${Date.now().toString(36)}`
  const scenario = buildScenario({ ...input, slug })

  scenarios = [...scenarios, scenario]
  notifyListeners()
  return scenario
}

/**
 * Maps one backend scenario record to the local chat shape. The icon and
 * opening line are generated locally since the backend does not store them.
 */
export function mapScenarioResponseToScenario(item: ScenarioResponse): Scenario {
  return buildScenario({
    slug: item.slug,
    title: item.title,
    description: item.description,
    duration: item.duration,
    level: apiLevelToLocal[item.level],
    situationContext: item.situationContext,
  })
}

/**
 * Replaces the in-memory scenario list with the authenticated user's
 * scenarios fetched from the backend (source of truth).
 */
export function replaceScenariosFromApi(items: ScenarioResponse[]) {
  scenarios = items.map(mapScenarioResponseToScenario)
  notifyListeners()
}

function slugify(value: string) {
  return (
    value
      .toLowerCase()
      .trim()
      .replace(/[^a-z0-9]+/g, '-')
      .replace(/^-+|-+$/g, '') || 'scenario'
  )
}
