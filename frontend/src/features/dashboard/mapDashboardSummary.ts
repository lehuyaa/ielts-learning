import type {
  DashboardRecentActivity,
  DashboardSummaryResponse,
} from '@/types/dashboard'

type DashboardCardTone = 'primary' | 'warning' | 'success' | 'orange'
type ActivityTone = 'success' | 'warning' | 'danger' | 'neutral' | 'primary'

export type DashboardViewModel = {
  greetingName: string
  subtitle: string
  streakDays: number
  progressCards: Array<{
    label: string
    value: string
    tone: DashboardCardTone
    supportText: string
  }>
  targetBand: {
    progressPercentage: number
    currentBandLabel: string
    currentBandCaption: string
    badgeLabel: string
    startLabel: string
    targetLabel: string
  }
  todayProgress: Array<{
    label: string
    value: string
  }>
  xp: {
    level: number
    levelTitle: string
    totalXp: number
    xpToNextLevel: number
    progressPercentage: number
  }
  reviewDue: {
    count: number
    title: string
    description: string
  }
  recentActivity: Array<{
    id: number
    title: string
    subtitle: string
    badgeLabel: string
    badgeTone: ActivityTone
    initial: string
  }>
  vocabularyStats: Array<{
    label: string
    value: string
  }>
}

export const loadingDashboardViewModel: DashboardViewModel = {
  greetingName: 'Learner',
  subtitle: 'Syncing your latest progress and review activity.',
  streakDays: 0,
  progressCards: [
    {
      label: 'Words Today',
      value: '0',
      tone: 'primary',
      supportText: 'Loading today’s activity',
    },
    {
      label: 'Review Due',
      value: '0',
      tone: 'warning',
      supportText: 'Checking due cards',
    },
    {
      label: 'Current Streak',
      value: '0d',
      tone: 'orange',
      supportText: 'Syncing streak',
    },
    {
      label: 'Mastery',
      value: '0%',
      tone: 'success',
      supportText: 'Preparing learning stats',
    },
  ],
  targetBand: {
    progressPercentage: 0,
    currentBandLabel: 'Band --',
    currentBandCaption: 'Current estimated band',
    badgeLabel: '0%',
    startLabel: 'Start',
    targetLabel: 'Target',
  },
  todayProgress: [
    { label: 'Words learned', value: '0' },
    { label: 'Cards reviewed', value: '0' },
    { label: 'Quizzes taken', value: '0' },
    { label: 'Lessons done', value: '0' },
    { label: 'XP earned', value: '0' },
  ],
  xp: {
    level: 1,
    levelTitle: 'Getting Started',
    totalXp: 0,
    xpToNextLevel: 100,
    progressPercentage: 0,
  },
  reviewDue: {
    count: 0,
    title: 'Checking review queue',
    description: 'We are counting cards that are ready to review.',
  },
  recentActivity: [],
  vocabularyStats: [
    { label: 'Total Learned', value: '0' },
    { label: 'Lessons Done', value: '0' },
    { label: 'Mastery', value: '0%' },
    { label: 'Review Due', value: '0' },
  ],
}

export function mapDashboardSummaryToViewModel(
  summary: DashboardSummaryResponse,
): DashboardViewModel {
  const currentBand =
    summary.targetBand.currentBand ?? summary.user.currentBand ?? null

  return {
    greetingName: firstName(summary.user.name) || 'Learner',
    subtitle: buildSubtitle(summary),
    streakDays: summary.user.currentStreak,
    progressCards: [
      {
        label: 'Words Today',
        value: String(summary.todayProgress.wordsLearned),
        tone: 'primary',
        supportText: `${summary.todayProgress.wordsReviewed} reviewed today`,
      },
      {
        label: 'Review Due',
        value: String(summary.reviewDue.count),
        tone: 'warning',
        supportText:
          summary.reviewDue.count > 0 ? 'Ready to practice now' : 'All caught up',
      },
      {
        label: 'Current Streak',
        value: `${summary.user.currentStreak}d`,
        tone: 'orange',
        supportText:
          summary.user.currentStreak > 0
            ? 'Keep the streak alive'
            : 'Start today to build a streak',
      },
      {
        label: 'Mastery',
        value: `${summary.learningStats.masteryPercentage}%`,
        tone: 'success',
        supportText: `${summary.learningStats.lessonsCompleted} lessons completed`,
      },
    ],
    targetBand: {
      progressPercentage: clampPercentage(summary.targetBand.progressPercentage),
      currentBandLabel: formatBand(currentBand),
      currentBandCaption: 'Current estimated band',
      badgeLabel: `${clampPercentage(summary.targetBand.progressPercentage)}%`,
      startLabel: formatBand(summary.targetBand.startingBand),
      targetLabel: `${formatBand(summary.targetBand.targetBand)} target`,
    },
    todayProgress: [
      {
        label: 'Words learned',
        value: String(summary.todayProgress.wordsLearned),
      },
      {
        label: 'Cards reviewed',
        value: String(summary.todayProgress.wordsReviewed),
      },
      {
        label: 'Quizzes taken',
        value: String(summary.todayProgress.quizzesTaken),
      },
      {
        label: 'Lessons done',
        value: String(summary.todayProgress.lessonsDone),
      },
      {
        label: 'XP earned',
        value: String(summary.todayProgress.xpEarned),
      },
    ],
    xp: {
      level: summary.xp.level,
      levelTitle: summary.xp.levelTitle,
      totalXp: summary.xp.totalXp,
      xpToNextLevel: summary.xp.xpToNextLevel,
      progressPercentage: clampPercentage(summary.xp.progressPercentage),
    },
    reviewDue: {
      count: summary.reviewDue.count,
      title:
        summary.reviewDue.count > 0
          ? `${summary.reviewDue.count} cards are waiting`
          : 'No reviews due right now',
      description:
        summary.reviewDue.count > 0
          ? 'Open your review queue to keep words fresh and protect your streak.'
          : 'Nice work. New review cards will appear here when their schedule is due.',
    },
    recentActivity: summary.recentActivity.map(mapActivity),
    vocabularyStats: [
      {
        label: 'Total Learned',
        value: String(summary.learningStats.totalWordsLearned),
      },
      {
        label: 'Lessons Done',
        value: String(summary.learningStats.lessonsCompleted),
      },
      {
        label: 'Mastery',
        value: `${summary.learningStats.masteryPercentage}%`,
      },
      {
        label: 'Review Due',
        value: String(summary.reviewDue.count),
      },
    ],
  }
}

function buildSubtitle(summary: DashboardSummaryResponse) {
  const currentBand = summary.targetBand.currentBand ?? summary.user.currentBand
  const currentBandText = currentBand
    ? `Current ${formatBand(currentBand)}`
    : 'Progress is still warming up'

  return `${currentBandText} • Target ${formatBand(summary.user.targetBand)} • Level ${summary.xp.level}`
}

function mapActivity(activity: DashboardRecentActivity) {
  return {
    id: activity.id,
    title: activity.title,
    subtitle: `${activity.description} • ${formatRelativeTime(activity.createdAt)}`,
    badgeLabel:
      activity.xp > 0 ? `+${activity.xp} XP` : formatActivityType(activity.type),
    badgeTone: getActivityTone(activity.type),
    initial: activity.title.charAt(0).toUpperCase() || 'A',
  }
}

function getActivityTone(type: string): ActivityTone {
  switch (type.toUpperCase()) {
    case 'QUIZ_PASSED':
    case 'QUIZ_COMPLETED':
      return 'success'
    case 'FLASHCARD_REVIEW':
      return 'primary'
    case 'LESSON_COMPLETED':
      return 'warning'
    case 'QUIZ_FAILED':
      return 'danger'
    default:
      return 'neutral'
  }
}

function formatBand(value: number | null | undefined) {
  return typeof value === 'number' ? `Band ${value.toFixed(1)}` : 'Band --'
}

function firstName(name: string) {
  return name.trim().split(/\s+/)[0] ?? ''
}

function clampPercentage(value: number) {
  return Math.max(0, Math.min(100, Math.round(value)))
}

function formatActivityType(type: string) {
  return type
    .toLowerCase()
    .split('_')
    .map((part) => part.charAt(0).toUpperCase() + part.slice(1))
    .join(' ')
}

function formatRelativeTime(value: string) {
  const timestamp = new Date(value)
  const diffMs = timestamp.getTime() - Date.now()

  if (Number.isNaN(timestamp.getTime())) {
    return 'Recently'
  }

  const formatter = new Intl.RelativeTimeFormat('en', { numeric: 'auto' })
  const minutes = Math.round(diffMs / (1000 * 60))

  if (Math.abs(minutes) < 60) {
    return formatter.format(minutes, 'minute')
  }

  const hours = Math.round(minutes / 60)
  if (Math.abs(hours) < 24) {
    return formatter.format(hours, 'hour')
  }

  const days = Math.round(hours / 24)
  if (Math.abs(days) < 7) {
    return formatter.format(days, 'day')
  }

  const weeks = Math.round(days / 7)
  return formatter.format(weeks, 'week')
}
