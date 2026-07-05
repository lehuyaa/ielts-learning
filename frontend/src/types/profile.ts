export type ProfileResponse = {
  user: ProfileUser
  stats: ProfileStats
  achievements: ProfileAchievement[]
  activitySummary: ProfileActivitySummary
}

export type ProfileUser = {
  id: number
  name: string
  username: string | null
  email: string
  avatarUrl: string
  initials: string
  memberSince: string
  currentBand: number | null
  targetBand: number
  totalXp: number
  level: number
  levelTitle: string
  currentLevelXp: number
  nextLevelXp: number
  xpUntilNextLevel: number
  currentStreak: number
  longestStreak: number
  timezone: string
  locale: string
}

export type ProfileStats = {
  wordsLearned: number
  lessonsDone: number
  masteredWords: number
  learningWords: number
  newWords: number
  vocabularyMasteryPercentage: number
}

export type ProfileAchievement = {
  id: number
  code: string
  title: string
  description: string
  icon: string
  unlocked: boolean
  unlockedAt: string | null
  progressValue: number
  requirementValue: number
}

export type ProfileActivitySummary = {
  activeDays: number
  wordsLearnedLast84Days: number
  averageWordsPerDay: number
  recentActivity: ProfileRecentActivity[]
}

export type ProfileRecentActivity = {
  id: number
  type: string
  title: string
  description: string
  xp: number
  createdAt: string
}

export type UpdateProfileInput = {
  name?: string
  username?: string
  targetBand?: number
  timezone?: string
  locale?: string
}
