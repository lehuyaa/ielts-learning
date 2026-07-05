export type UserRole = 'USER' | 'ADMIN'

export type User = {
  id: number
  email: string
  name: string
  username: string | null
  avatarUrl: string
  role: UserRole
  targetBand: number
  currentBand: number | null
  startingBand: number | null
  recommendedBand: number | null
  placementCompletedAt: string | null
  totalXp: number
  level: number
  levelTitle: string
  currentStreak: number
  longestStreak: number
  createdAt: string
}
