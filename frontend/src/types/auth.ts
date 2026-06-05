export type User = {
  id: number
  email: string
  name: string
  username: string | null
  avatarUrl: string
  role: 'USER' | 'ADMIN'
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

export type AuthResponse = {
  accessToken: string
  tokenType: 'Bearer'
  user: User
}

export type LoginInput = {
  email: string
  password: string
}

export type RegisterInput = {
  name: string
  email: string
  username?: string
  password: string
  targetBand: number
}
