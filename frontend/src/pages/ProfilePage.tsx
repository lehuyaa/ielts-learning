import { useState } from 'react'
import {
  BookOpen,
  Calendar,
  CheckCircle2,
  ChevronLeft,
  Edit2,
  Flame,
  Star,
  Target,
  Trophy,
} from 'lucide-react'
import { useNavigate } from 'react-router-dom'

import {
  mockProfile,
  profileTabs,
  type ProfileTab,
} from '@/features/profile/mockProfile'
import { cn } from '@/lib/utils'

const statIcons = {
  'Words Learned': BookOpen,
  'Current Streak': Flame,
  'Longest Streak': Trophy,
  'Lessons Done': Target,
} as const

export function ProfilePage() {
  const navigate = useNavigate()
  const [activeTab, setActiveTab] = useState<ProfileTab>('Overview')
  const xpProgress = Math.round(
    (mockProfile.currentXp / mockProfile.nextLevelXp) * 100,
  )

  return (
    <div className="space-y-6">
      <header className="sticky top-0 z-10 -mx-4 border-b border-border bg-background/95 px-4 backdrop-blur supports-[backdrop-filter]:bg-background/80 sm:-mx-6 sm:px-6 lg:-mx-8 lg:px-8">
        <div className="mx-auto flex h-16 max-w-4xl items-center gap-4">
          <button
            className="cursor-pointer rounded-xl p-2 transition-colors hover:bg-muted/60"
            onClick={() => {
              void navigate('/dashboard')
            }}
            type="button"
          >
            <ChevronLeft className="size-5 text-muted-foreground" />
          </button>
          <span className="text-sm font-semibold text-foreground">
            My Profile
          </span>
          <div className="flex-1" />
          <button
            className="flex cursor-pointer items-center gap-2 rounded-xl border border-border px-3 py-1.5 text-sm text-muted-foreground transition-colors hover:bg-muted/50"
            type="button"
          >
            <Edit2 className="size-3.5" />
            Edit
          </button>
        </div>
      </header>

      <div className="mx-auto max-w-4xl space-y-6">
        <section className="overflow-hidden rounded-[28px] border border-border bg-white">
          <div
            className="h-28"
            style={{
              background:
                'linear-gradient(135deg, #4F46E5 0%, #7C3AED 60%, #10B981 100%)',
            }}
          />

          <div className="px-6 pb-6">
            <div className="-mt-10 mb-4 flex items-end justify-between gap-4">
              <div className="flex size-20 items-center justify-center rounded-2xl border-4 border-white bg-indigo-100 text-3xl font-bold text-primary shadow-lg">
                {mockProfile.avatarInitial}
              </div>
              <div className="mb-1 flex items-center gap-2 rounded-full border border-amber-200 bg-amber-50 px-3 py-1.5">
                <Star className="size-3.5 fill-amber-500 text-amber-500" />
                <span className="text-xs font-semibold text-amber-700">
                  Level {mockProfile.level} · {mockProfile.levelTitle}
                </span>
              </div>
            </div>

            <div className="flex flex-col justify-between gap-4 sm:flex-row sm:items-end">
              <div>
                <h1 className="text-xl font-bold text-foreground">
                  {mockProfile.name}
                </h1>
                <div className="text-sm text-muted-foreground">
                  {mockProfile.username} · {mockProfile.memberSince}
                </div>
                <div className="mt-2 flex flex-wrap items-center gap-3">
                  <span className="rounded-full bg-indigo-50 px-2 py-0.5 text-xs font-semibold text-indigo-700">
                    {mockProfile.currentBand}
                  </span>
                  <span className="rounded-full bg-emerald-50 px-2 py-0.5 text-xs font-semibold text-emerald-700">
                    {mockProfile.targetBand}
                  </span>
                </div>
              </div>

              <div className="flex items-center gap-2 rounded-2xl border border-orange-200 bg-orange-50 px-4 py-2.5">
                <Flame className="size-5 fill-orange-500 text-orange-500" />
                <div>
                  <div className="text-lg font-bold leading-none text-foreground">
                    {mockProfile.streakDays}
                  </div>
                  <div className="text-xs text-muted-foreground">Day streak</div>
                </div>
              </div>
            </div>
          </div>
        </section>

        <section className="rounded-2xl border border-border bg-white p-5">
          <div className="mb-2 flex items-center justify-between gap-4">
            <div className="text-sm font-semibold text-foreground">
              Level {mockProfile.level} Progress
            </div>
            <div className="text-xs text-muted-foreground">
              {formatNumber(mockProfile.currentXp)} /{' '}
              {formatNumber(mockProfile.nextLevelXp)} XP
            </div>
          </div>

          <div className="mb-1 h-3 overflow-hidden rounded-full bg-muted">
            <div
              className="h-full rounded-full"
              style={{
                width: `${xpProgress}%`,
                background: 'linear-gradient(90deg, #F59E0B, #EF4444)',
              }}
            />
          </div>

          <div className="text-xs text-muted-foreground">
            {formatNumber(mockProfile.xpRemaining)} XP until Level{' '}
            {mockProfile.level + 1}
          </div>
        </section>

        <section className="grid grid-cols-2 gap-4 md:grid-cols-4">
          {mockProfile.stats.map((stat) => {
            const Icon = statIcons[stat.label as keyof typeof statIcons]

            return (
              <div
                className="rounded-2xl border border-border bg-white p-4"
                key={stat.label}
              >
                <div
                  className={cn(
                    'mb-3 flex size-9 items-center justify-center rounded-xl',
                    getStatToneClasses(stat.tone),
                  )}
                >
                  <Icon className="size-[18px]" />
                </div>
                <div className="text-2xl font-bold text-foreground">
                  {stat.value}
                </div>
                <div className="mt-0.5 text-xs text-muted-foreground">
                  {stat.label}
                </div>
              </div>
            )
          })}
        </section>

        <section className="rounded-2xl border border-border bg-white p-5">
          <div className="mb-4 flex items-center justify-between">
            <div className="text-sm font-semibold text-foreground">
              Vocabulary Mastery
            </div>
            <div className="text-sm font-bold text-primary">
              {mockProfile.masteryPercentage}%
            </div>
          </div>

          <div className="mb-4 h-3 overflow-hidden rounded-full bg-muted">
            <div
              className="h-full rounded-full"
              style={{
                width: `${mockProfile.masteryPercentage}%`,
                background: 'linear-gradient(90deg, #4F46E5, #7C3AED)',
              }}
            />
          </div>

          <div className="grid grid-cols-3 gap-3 text-center text-xs text-muted-foreground">
            {mockProfile.masteryBreakdown.map((item) => (
              <div key={item.label}>
                <div className={cn('text-lg font-bold', item.tone)}>
                  {item.value}
                </div>
                <div>{item.label}</div>
              </div>
            ))}
          </div>
        </section>

        <section className="rounded-xl bg-muted p-1">
          <div className="flex gap-1">
            {profileTabs.map((tab) => (
              <button
                className={cn(
                  'flex-1 cursor-pointer rounded-lg py-2 text-sm font-medium transition-all',
                  activeTab === tab
                    ? 'bg-white text-primary shadow-sm'
                    : 'text-muted-foreground hover:text-foreground',
                )}
                key={tab}
                onClick={() => {
                  setActiveTab(tab)
                }}
                type="button"
              >
                {tab}
              </button>
            ))}
          </div>
        </section>

        {activeTab === 'Overview' ? <OverviewTab /> : null}
        {activeTab === 'Achievements' ? <AchievementsTab /> : null}
        {activeTab === 'Calendar' ? <CalendarTab /> : null}
      </div>
    </div>
  )
}

function OverviewTab() {
  return (
    <section className="grid gap-4 md:grid-cols-2">
      {mockProfile.topicProgress.map((topic) => (
        <div
          className="flex items-center gap-4 rounded-2xl border border-border bg-white p-4"
          key={topic.topic}
        >
          <div className="w-24 shrink-0 text-sm font-medium text-foreground">
            {topic.topic}
          </div>
          <div className="h-2 flex-1 overflow-hidden rounded-full bg-muted">
            <div
              className={cn('h-full rounded-full', topic.color)}
              style={{ width: `${topic.progress}%` }}
            />
          </div>
          <div className="w-10 text-right text-sm font-semibold text-foreground">
            {topic.progress}%
          </div>
        </div>
      ))}
    </section>
  )
}

function AchievementsTab() {
  return (
    <section className="grid gap-4 md:grid-cols-3">
      {mockProfile.achievements.map((achievement) => (
        <div
          className={cn(
            'rounded-2xl border bg-white p-5',
            achievement.earned ? 'border-border' : 'border-border opacity-50',
          )}
          key={achievement.label}
        >
          <div className="flex items-start gap-3">
            <div className={cn('text-3xl', !achievement.earned && 'grayscale')}>
              {achievement.icon}
            </div>
            <div className="min-w-0 flex-1">
              <div className="flex items-center gap-2">
                <div className="text-sm font-semibold text-foreground">
                  {achievement.label}
                </div>
                {achievement.earned ? (
                  <CheckCircle2 className="size-3.5 shrink-0 text-emerald-500" />
                ) : null}
              </div>
              <div className="mt-0.5 text-xs text-muted-foreground">
                {achievement.description}
              </div>
              {achievement.earned ? (
                <div className="mt-1 text-xs text-primary">
                  {achievement.date}
                </div>
              ) : null}
            </div>
          </div>
        </div>
      ))}
    </section>
  )
}

function CalendarTab() {
  return (
    <section className="rounded-2xl border border-border bg-white p-5">
      <div className="mb-4 flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <div className="text-sm font-semibold text-foreground">
          Learning Activity
        </div>
        <div className="flex items-center gap-1.5 text-xs text-muted-foreground">
          Less
          {[0.1, 0.3, 0.5, 0.75, 1].map((opacity) => (
            <div
              className="h-3 w-3 rounded-sm"
              key={opacity}
              style={{ backgroundColor: `rgba(79,70,229,${opacity})` }}
            />
          ))}
          More
        </div>
      </div>

      <div className="mb-1 ml-7 flex gap-1.5">
        {Array.from({ length: 12 }, (_, index) => (
          <div
            className="text-xs text-muted-foreground"
            key={index}
            style={{ width: `${7 * 12 + 6 * 1.5}px` }}
          >
            {index % 3 === 0 ? `Week ${index + 1}` : ''}
          </div>
        ))}
      </div>

      <div className="flex gap-1.5">
        <div className="flex flex-col gap-1.5">
          {mockProfile.weekLabels.map((label) => (
            <div
              className="flex h-3 items-center text-xs text-muted-foreground"
              key={label}
            >
              {label[0]}
            </div>
          ))}
        </div>

        <div className="flex gap-1.5 overflow-x-auto">
          {Array.from({ length: 12 }, (_, week) => (
            <div className="flex flex-col gap-1.5" key={week}>
              {Array.from({ length: 7 }, (_, day) => {
                const index = week * 7 + day
                return (
                  <HeatCell
                    count={mockProfile.heatmap[index] ?? 0}
                    key={`${week}-${day}`}
                  />
                )
              })}
            </div>
          ))}
        </div>
      </div>

      <div className="mt-4 flex items-center gap-2 text-xs text-muted-foreground">
        <Calendar className="size-[13px]" />
        847 words learned in the last 84 days · Avg. 10.1/day
      </div>
    </section>
  )
}

function HeatCell({ count }: { count: number }) {
  const opacity =
    count === 0 ? 0 : count < 5 ? 0.2 : count < 12 ? 0.45 : count < 22 ? 0.7 : 1

  return (
    <div
      className="h-3 w-3 rounded-sm"
      style={{
        backgroundColor:
          count === 0 ? '#ede9fe' : `rgba(79, 70, 229, ${opacity})`,
      }}
      title={count > 0 ? `${count} words` : 'No activity'}
    />
  )
}

function getStatToneClasses(tone: string) {
  switch (tone) {
    case 'orange':
      return 'bg-orange-50 text-orange-500'
    case 'amber':
      return 'bg-amber-50 text-amber-500'
    case 'emerald':
      return 'bg-emerald-50 text-emerald-600'
    case 'primary':
    default:
      return 'bg-indigo-50 text-primary'
  }
}

function formatNumber(value: number) {
  return new Intl.NumberFormat('en-US').format(value)
}
