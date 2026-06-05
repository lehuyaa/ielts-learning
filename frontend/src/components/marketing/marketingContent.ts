import {
  BarChart3,
  Brain,
  CheckCircle2,
  ClipboardList,
  Layers3,
  Map,
  RotateCcw,
  Sparkles,
  Target,
  Trophy,
} from 'lucide-react'

export const heroBands = ['5.0', '6.0', '7.0', '8.0']

export const landingStats = [
  {
    value: '2000+',
    label: 'Vocabulary Words',
  },
  {
    value: '50+',
    label: 'Lessons',
  },
  {
    value: '8',
    label: 'IELTS Band Levels',
  },
  {
    value: '100+',
    label: 'Quiz Questions',
  },
]

export const landingFeatures = [
  {
    title: 'Structured Learning Roadmap',
    description:
      'Move through band-based lessons that show what to learn next and why it matters.',
    icon: Map,
  },
  {
    title: 'Flashcards & Spaced Repetition',
    description:
      'Review words at the right time with simple ratings that keep practice focused.',
    icon: RotateCcw,
  },
  {
    title: 'IELTS Topic-Based Vocabulary',
    description:
      'Learn academic words by common IELTS themes like education, technology, and society.',
    icon: Layers3,
  },
  {
    title: 'Progress Tracking',
    description:
      'See target band progress, streaks, review due counts, and vocabulary growth.',
    icon: BarChart3,
  },
  {
    title: 'Quiz & Review System',
    description:
      'Check understanding with quick quizzes before moving ahead in the roadmap.',
    icon: ClipboardList,
  },
  {
    title: 'Personalized Learning Journey',
    description:
      'Start from your target band and build a repeatable plan around your goal.',
    icon: Target,
  },
]

export const roadmapPreview = [
  {
    band: 'Band 5.0',
    title: 'Core IELTS Vocabulary',
    lessons: '12 lessons',
    icon: Brain,
  },
  {
    band: 'Band 6.0',
    title: 'Topic Fluency',
    lessons: '14 lessons',
    icon: CheckCircle2,
  },
  {
    band: 'Band 7.0',
    title: 'Academic Precision',
    lessons: '16 lessons',
    icon: Sparkles,
  },
  {
    band: 'Band 8.0',
    title: 'Advanced Expression',
    lessons: '8 lessons',
    icon: Trophy,
  },
]

export const testimonials = [
  {
    name: 'Minh Anh',
    targetBand: 'Target Band 7.0',
    feedback:
      'The roadmap made IELTS vocabulary feel organized instead of overwhelming.',
  },
  {
    name: 'Sarah Nguyen',
    targetBand: 'Target Band 6.5',
    feedback:
      'Flashcard reviews helped me remember words I used to forget after one day.',
  },
  {
    name: 'Duc Tran',
    targetBand: 'Target Band 8.0',
    feedback:
      'The topic lessons gave me better phrases for writing and speaking practice.',
  },
]

export const pricingPlans = [
  {
    name: 'Free Plan',
    price: '$0',
    description: 'Start building your IELTS vocabulary roadmap today.',
    features: ['Band roadmap preview', 'Starter lessons', 'Basic flashcards'],
    cta: 'Start Learning Free',
    to: '/register',
    highlighted: true,
  },
  {
    name: 'Premium Plan',
    price: 'Coming Soon',
    description: 'Advanced practice, deeper analytics, and expanded lessons.',
    features: ['Full roadmap access', 'Advanced reviews', 'Progress insights'],
    cta: 'Coming Soon',
    to: '/register',
    highlighted: false,
  },
]
