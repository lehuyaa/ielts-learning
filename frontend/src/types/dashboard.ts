export type DashboardSummaryResponse = {
  user: DashboardUserSummary
  todayProgress: DashboardTodayProgress
  learningStats: DashboardLearningStats
  targetBand: DashboardTargetBand
  xp: DashboardXP
  reviewDue: DashboardReviewDue
  recentActivity: DashboardRecentActivity[]
}

export type DashboardUserSummary = {
  id: number
  name: string
  email: string
  avatarUrl: string
  targetBand: number
  currentBand: number | null
  currentStreak: number
}

export type DashboardTodayProgress = {
  date: string
  wordsLearned: number
  wordsReviewed: number
  quizzesTaken: number
  lessonsDone: number
  xpEarned: number
}

export type DashboardLearningStats = {
  totalWordsLearned: number
  lessonsCompleted: number
  masteryPercentage: number
}

export type DashboardTargetBand = {
  startingBand: number | null
  currentBand: number | null
  targetBand: number
  progressPercentage: number
}

export type DashboardXP = {
  totalXp: number
  level: number
  levelTitle: string
  xpToNextLevel: number
  progressPercentage: number
}

export type DashboardReviewDue = {
  count: number
}

export type DashboardRecentActivity = {
  id: number
  type: string
  title: string
  description: string
  xp: number
  createdAt: string
}
