import {
  Briefcase,
  Coffee,
  MessageSquare,
  Plane,
  Sparkles,
  Users,
  Utensils,
  type LucideIcon,
} from 'lucide-react'

export type ScenarioLevel = 'Beginner' | 'Intermediate' | 'Advanced'

export type Scenario = {
  slug: string
  title: string
  description: string
  duration: string
  level: ScenarioLevel
  icon: LucideIcon
  openingLine: string
  aiReplies: string[]
}

export const scenarios: Scenario[] = [
  {
    slug: 'free-conversation',
    title: 'Free Conversation',
    description:
      'Open-ended chat on any topic. Build confidence and fluency naturally.',
    duration: '5-15 min',
    level: 'Beginner',
    icon: MessageSquare,
    openingLine: "Hi! I'm ready to chat about anything. What's on your mind today?",
    aiReplies: [
      "That's interesting! Can you tell me more about that?",
      'I see. How did that make you feel?',
      "Nice! What happened next?",
      "Got it. What's your opinion on that?",
    ],
  },
  {
    slug: 'daily-conversation',
    title: 'Daily Conversation',
    description:
      'Talk about everyday life — hobbies, routines, and current events.',
    duration: '5-10 min',
    level: 'Beginner',
    icon: Coffee,
    openingLine: 'Hey there! How has your day been so far?',
    aiReplies: [
      'That sounds productive! What kind of project was it? Did you face any challenges?',
      "That's great to hear. Do you usually do that every day?",
      'Interesting! What do you like to do in your free time?',
      'Nice routine. What are you planning for later today?',
    ],
  },
  {
    slug: 'travel',
    title: 'Travel',
    description:
      'Navigate airports, hotels, and tourist attractions in English.',
    duration: '8-12 min',
    level: 'Intermediate',
    icon: Plane,
    openingLine:
      "Welcome to the airport check-in desk. Can I see your passport and ticket, please?",
    aiReplies: [
      'Great, everything looks in order. Do you have any luggage to check in?',
      'Would you prefer a window or an aisle seat?',
      'Your gate is B12, boarding starts in 40 minutes. Anything else I can help with?',
      'Enjoy your trip! Let me know if you need directions to the lounge.',
    ],
  },
  {
    slug: 'restaurant',
    title: 'Restaurant',
    description:
      'Order food, make reservations, and handle dining situations.',
    duration: '5-10 min',
    level: 'Intermediate',
    icon: Utensils,
    openingLine: 'Good evening! Welcome in. Do you have a reservation with us tonight?',
    aiReplies: [
      'Perfect, right this way. Can I start you off with something to drink?',
      'Great choice. Would you like that with a side salad or fries?',
      'No problem, I will let the kitchen know about the allergy.',
      'Your order will be ready in about fifteen minutes.',
    ],
  },
  {
    slug: 'job-interview',
    title: 'Job Interview',
    description:
      'Practice common interview questions and professional vocabulary.',
    duration: '10-20 min',
    level: 'Advanced',
    icon: Briefcase,
    openingLine: "Thanks for coming in today. Could you start by telling me a bit about yourself?",
    aiReplies: [
      "That's a strong background. What would you say is your greatest strength?",
      'Can you walk me through a time you solved a difficult problem at work?',
      'Why are you interested in this position specifically?',
      'Do you have any questions for us before we wrap up?',
    ],
  },
  {
    slug: 'work-meetings',
    title: 'Work & Meetings',
    description: 'Lead meetings, present ideas, and collaborate with colleagues.',
    duration: '10-15 min',
    level: 'Advanced',
    icon: Users,
    openingLine: "Let's get started. Can you give the team a quick update on your progress?",
    aiReplies: [
      "That's good progress. What's blocking you from finishing this sprint?",
      'Can you share that data with the team after this call?',
      'Let’s circle back on that in the next meeting. Anything else to add?',
      'Sounds good. Who should own the next action item?',
    ],
  },
]

export function getScenarioBySlug(slug: string | undefined) {
  return getAllScenarios().find((scenario) => scenario.slug === slug)
}

// In-memory store for user-created scenarios. Not persisted to a backend or
// localStorage yet — it lives only for the current browser session/tab.
let customScenarios: Scenario[] = []
const listeners = new Set<() => void>()

function notifyListeners() {
  listeners.forEach((listener) => listener())
}

export function getCustomScenarios() {
  return customScenarios
}

export function getAllScenarios() {
  return [...scenarios, ...customScenarios]
}

export function subscribeToCustomScenarios(listener: () => void) {
  listeners.add(listener)
  return () => listeners.delete(listener)
}

export type NewScenarioInput = {
  title: string
  description: string
  duration: string
  level: ScenarioLevel
  situationContext?: string
}

export function addCustomScenario(input: NewScenarioInput): Scenario {
  const slug = `custom-${slugify(input.title)}-${Date.now().toString(36)}`
  const context = input.situationContext?.trim()

  const scenario: Scenario = {
    slug,
    title: input.title,
    description: input.description,
    duration: input.duration || '5-10 min',
    level: input.level,
    icon: Sparkles,
    openingLine: context
      ? `Let's begin. ${context}`
      : `Let's practice "${input.title}". ${input.description}`,
    aiReplies: [
      "That's a great point. Can you tell me more?",
      'I see. What would you say or do next in this situation?',
      'Good. How would you respond if that did not go as planned?',
      "Let's keep going — what happens after that?",
    ],
  }

  customScenarios = [...customScenarios, scenario]
  notifyListeners()
  return scenario
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
