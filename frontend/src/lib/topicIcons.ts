type TopicIconInput = {
  slug: string
  icon?: string
  emoji?: string
}

const topicIconBySlug: Record<string, string> = {
  education: '🎓',
  health: '🏥',
  environment: '🌿',
  technology: '💻',
  business: '💼',
  society: '👥',
  media: '📺',
  science: '🔬',
  government: '🏛️',
  'global-issues': '🌍',
  economics: '📈',
  economy: '📈',
  law: '⚖️',
  culture: '🎭',
  innovation: '💡',
  psychology: '🧠',
  travel: '✈️',
  work: '💼',
}

const topicIconByName: Record<string, string> = {
  'graduation-cap': '🎓',
  'heart-pulse': '🏥',
  leaf: '🌿',
  cpu: '💻',
  'briefcase-business': '💼',
  briefcase: '💼',
  'users-round': '👥',
  newspaper: '📺',
  microscope: '🔬',
  landmark: '🏛️',
  'globe-2': '🌍',
  'chart-no-axes-combined': '📈',
  scale: '⚖️',
  palette: '🎭',
  lightbulb: '💡',
  brain: '🧠',
  plane: '✈️',
}

export function getTopicDisplayIcon({ slug, icon, emoji }: TopicIconInput) {
  if (emoji && !/^[A-Z]{2}$/.test(emoji)) {
    return emoji
  }

  if (icon && topicIconByName[icon]) {
    return topicIconByName[icon]
  }

  return topicIconBySlug[slug] ?? '📚'
}
