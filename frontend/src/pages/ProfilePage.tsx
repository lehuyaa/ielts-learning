import { zodResolver } from '@hookform/resolvers/zod'
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
import { useEffect, useState } from 'react'
import { useForm } from 'react-hook-form'
import { useNavigate } from 'react-router-dom'

import { APIError } from '@/api/api'
import { CardSkeleton } from '@/components/state/CardSkeleton'
import { EmptyState } from '@/components/state/EmptyState'
import { ErrorState } from '@/components/state/ErrorState'
import { Skeleton } from '@/components/ui/skeleton'
import { Button } from '@/components/ui/button'
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { useProfile } from '@/features/profile/hooks/useProfile'
import { useUpdateProfile } from '@/features/profile/hooks/useUpdateProfile'
import {
  normalizeProfileUpdateValues,
  profileUpdateSchema,
  type ProfileUpdateFormValues,
} from '@/features/profile/validation/profileSchemas'
import { cn } from '@/lib/utils'
import type {
  ProfileAchievement,
  ProfileRecentActivity,
  ProfileResponse,
} from '@/types/profile'

const profileTabs = ['Overview', 'Achievements', 'Calendar'] as const
type ProfileTab = (typeof profileTabs)[number]

const weekLabels = ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun']
const emptyHeatmap = Array.from({ length: 84 }, () => 0)

const statIcons = {
  'Words Learned': BookOpen,
  'Current Streak': Flame,
  'Longest Streak': Trophy,
  'Lessons Done': Target,
} as const

type OverviewItem = {
  label: string
  progress: number
  value: string
  color: string
}

export function ProfilePage() {
  const navigate = useNavigate()
  const profileQuery = useProfile()
  const updateProfileMutation = useUpdateProfile()
  const [activeTab, setActiveTab] = useState<ProfileTab>('Overview')
  const [isEditing, setIsEditing] = useState(false)

  const form = useForm<ProfileUpdateFormValues>({
    resolver: zodResolver(profileUpdateSchema),
    defaultValues: {
      name: '',
      username: '',
      targetBand: 7,
      timezone: 'UTC',
      locale: 'en',
    },
  })

  useEffect(() => {
    if (!profileQuery.data) {
      return
    }

    form.reset({
      name: profileQuery.data.user.name,
      username: profileQuery.data.user.username ?? '',
      targetBand: profileQuery.data.user.targetBand,
      timezone: profileQuery.data.user.timezone,
      locale: profileQuery.data.user.locale,
    })
  }, [form, profileQuery.data])

  const errorMessage = getProfileErrorMessage(profileQuery.error)
  const isEmpty =
    !profileQuery.isLoading && !errorMessage && !profileQuery.data
  const profile = profileQuery.data

  const xpProgress = profile
    ? profile.user.nextLevelXp > 0
      ? Math.round(
          (profile.user.currentLevelXp / profile.user.nextLevelXp) * 100,
        )
      : 0
    : 0

  const levelProgressLabel = profile
    ? `${formatNumber(profile.user.totalXp)} / ${formatNumber(
        profile.user.totalXp + profile.user.xpUntilNextLevel,
      )} XP`
    : '0 / 0 XP'

  const profileStats = profile ? buildStatCards(profile) : []
  const overviewItems = profile ? buildOverviewItems(profile) : []
  async function onSubmit(values: ProfileUpdateFormValues) {
    try {
      await updateProfileMutation.mutateAsync(normalizeProfileUpdateValues(values))
      setIsEditing(false)
    } catch (error) {
      if (error instanceof APIError && error.fields) {
        applyFieldErrors(error.fields)
      }
    }
  }

  function applyFieldErrors(fields: Record<string, string>) {
    for (const [field, message] of Object.entries(fields)) {
      if (
        field === 'name' ||
        field === 'username' ||
        field === 'targetBand' ||
        field === 'timezone' ||
        field === 'locale'
      ) {
        form.setError(field, { message })
      }
    }
  }

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
            onClick={() => {
              setIsEditing((current) => !current)
            }}
            type="button"
          >
            <Edit2 className="size-3.5" />
            {isEditing ? 'Close' : 'Edit'}
          </button>
        </div>
      </header>

      <div className="mx-auto max-w-4xl space-y-6">
        {profileQuery.isLoading ? (
          <ProfileLoadingSkeleton />
        ) : null}

        {errorMessage ? (
          <ErrorState
            actionHref="/dashboard"
            actionLabel="Back to dashboard"
            description={errorMessage}
            onRetry={() => {
              void profileQuery.refetch()
            }}
            title="Profile unavailable"
          />
        ) : null}

        {isEmpty ? (
          <EmptyState
            actionHref="/dashboard"
            actionLabel="Back to dashboard"
            description="We could not find profile information for this account yet."
            title="No profile data yet"
          />
        ) : null}

        {profile ? (
          <>
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
                    {profile.user.initials}
                  </div>
                  <div className="mb-1 flex items-center gap-2 rounded-full border border-amber-200 bg-amber-50 px-3 py-1.5">
                    <Star className="size-3.5 fill-amber-500 text-amber-500" />
                    <span className="text-xs font-semibold text-amber-700">
                      Level {profile.user.level} · {profile.user.levelTitle}
                    </span>
                  </div>
                </div>

                <div className="flex flex-col justify-between gap-4 sm:flex-row sm:items-end">
                  <div>
                    <h1 className="text-xl font-bold text-foreground">
                      {profile.user.name}
                    </h1>
                    <div className="text-sm text-muted-foreground">
                      {formatUsername(profile.user.username)} ·{' '}
                      {formatMemberSince(profile.user.memberSince)}
                    </div>
                    <div className="mt-2 flex flex-wrap items-center gap-3">
                      <span className="rounded-full bg-indigo-50 px-2 py-0.5 text-xs font-semibold text-indigo-700">
                        {profile.user.currentBand
                          ? `Band ${profile.user.currentBand.toFixed(1)} Current`
                          : 'Band Pending'}
                      </span>
                      <span className="rounded-full bg-emerald-50 px-2 py-0.5 text-xs font-semibold text-emerald-700">
                        Band {profile.user.targetBand.toFixed(1)} Target
                      </span>
                    </div>
                  </div>

                  <div className="flex items-center gap-2 rounded-2xl border border-orange-200 bg-orange-50 px-4 py-2.5">
                    <Flame className="size-5 fill-orange-500 text-orange-500" />
                    <div>
                      <div className="text-lg font-bold leading-none text-foreground">
                        {profile.user.currentStreak}
                      </div>
                      <div className="text-xs text-muted-foreground">
                        Day streak
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </section>

            {isEditing ? (
              <section className="rounded-2xl border border-border bg-white p-5">
                <div className="mb-4">
                  <h2 className="text-sm font-semibold text-foreground">
                    Edit profile
                  </h2>
                  <p className="mt-1 text-xs text-muted-foreground">
                    Update your public profile and learning preferences.
                  </p>
                </div>

                <Form {...form}>
                  <form onSubmit={form.handleSubmit(onSubmit)}>
                    <div className="space-y-6">
                      <div className="grid gap-6 sm:grid-cols-2">
                        <FormField
                          control={form.control}
                          name="name"
                          render={({ field }) => (
                            <FormItem>
                              <FormLabel>Name</FormLabel>
                              <FormControl>
                                <Input placeholder="Alex Johnson" {...field} />
                              </FormControl>
                              <FormMessage />
                            </FormItem>
                          )}
                        />

                        <FormField
                          control={form.control}
                          name="username"
                          render={({ field }) => (
                            <FormItem>
                              <FormLabel>Username</FormLabel>
                              <FormControl>
                                <Input placeholder="alexj" {...field} />
                              </FormControl>
                              <FormMessage />
                            </FormItem>
                          )}
                        />
                      </div>

                      <div className="grid gap-6 sm:grid-cols-3">
                        <FormField
                          control={form.control}
                          name="targetBand"
                          render={({ field }) => (
                            <FormItem>
                              <FormLabel>Target band</FormLabel>
                              <FormControl>
                                <Input
                                  max="8.5"
                                  min="5"
                                  step="0.5"
                                  type="number"
                                  value={field.value}
                                  onBlur={field.onBlur}
                                  onChange={(event) => {
                                    field.onChange(Number(event.target.value))
                                  }}
                                />
                              </FormControl>
                              <FormMessage />
                            </FormItem>
                          )}
                        />

                        <FormField
                          control={form.control}
                          name="timezone"
                          render={({ field }) => (
                            <FormItem className="sm:col-span-2">
                              <FormLabel>Timezone</FormLabel>
                              <FormControl>
                                <Input placeholder="Asia/Pontianak" {...field} />
                              </FormControl>
                              <FormMessage />
                            </FormItem>
                          )}
                        />
                      </div>

                      <FormField
                        control={form.control}
                        name="locale"
                        render={({ field }) => (
                          <FormItem>
                            <FormLabel>Locale</FormLabel>
                            <FormControl>
                              <Input placeholder="en" {...field} />
                            </FormControl>
                            <FormMessage />
                          </FormItem>
                        )}
                      />
                    </div>

                    <div className="mt-8 flex flex-col gap-3 sm:flex-row">
                      <Button
                        disabled={updateProfileMutation.isPending}
                        type="submit"
                      >
                        {updateProfileMutation.isPending
                          ? 'Saving changes'
                          : 'Save changes'}
                      </Button>
                      <Button
                        disabled={updateProfileMutation.isPending}
                        onClick={() => {
                          setIsEditing(false)
                          form.reset({
                            name: profile.user.name,
                            username: profile.user.username ?? '',
                            targetBand: profile.user.targetBand,
                            timezone: profile.user.timezone,
                            locale: profile.user.locale,
                          })
                        }}
                        type="button"
                        variant="outline"
                      >
                        Cancel
                      </Button>
                    </div>
                  </form>
                </Form>
              </section>
            ) : null}

            <section className="rounded-2xl border border-border bg-white p-5">
              <div className="mb-2 flex items-center justify-between gap-4">
                <div className="text-sm font-semibold text-foreground">
                  Level {profile.user.level} Progress
                </div>
                <div className="text-xs text-muted-foreground">
                  {levelProgressLabel}
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
                {formatNumber(profile.user.xpUntilNextLevel)} XP until Level{' '}
                {profile.user.level + 1}
              </div>
            </section>

            <section className="grid grid-cols-2 gap-4 md:grid-cols-4">
              {profileStats.map((stat) => {
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
                  {profile.stats.vocabularyMasteryPercentage}%
                </div>
              </div>

              <div className="mb-4 h-3 overflow-hidden rounded-full bg-muted">
                <div
                  className="h-full rounded-full"
                  style={{
                    width: `${profile.stats.vocabularyMasteryPercentage}%`,
                    background: 'linear-gradient(90deg, #4F46E5, #7C3AED)',
                  }}
                />
              </div>

              <div className="grid grid-cols-3 gap-3 text-center text-xs text-muted-foreground">
                {[
                  {
                    label: 'Mastered',
                    value: formatNumber(profile.stats.masteredWords),
                    tone: 'text-emerald-600',
                  },
                  {
                    label: 'Learning',
                    value: formatNumber(profile.stats.learningWords),
                    tone: 'text-amber-600',
                  },
                  {
                    label: 'New',
                    value: formatNumber(profile.stats.newWords),
                    tone: 'text-blue-500',
                  },
                ].map((item) => (
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

            {activeTab === 'Overview' ? (
              <OverviewTab
                activitySummary={profile.activitySummary}
                items={overviewItems}
              />
            ) : null}
            {activeTab === 'Achievements' ? (
              <AchievementsTab achievements={profile.achievements} />
            ) : null}
            {activeTab === 'Calendar' ? (
              <CalendarTab
                activeDays={profile.activitySummary.activeDays}
                averageWordsPerDay={profile.activitySummary.averageWordsPerDay}
                wordsLearnedLast84Days={
                  profile.activitySummary.wordsLearnedLast84Days
                }
              />
            ) : null}
          </>
        ) : null}
      </div>
    </div>
  )
}

function OverviewTab({
  items,
  activitySummary,
}: {
  items: OverviewItem[]
  activitySummary: ProfileResponse['activitySummary']
}) {
  return (
    <section className="space-y-4">
      <div className="grid gap-4 md:grid-cols-2">
        {items.map((item) => (
          <div
            className="flex items-center gap-4 rounded-2xl border border-border bg-white p-4"
            key={item.label}
          >
            <div className="w-28 shrink-0 text-sm font-medium text-foreground">
              {item.label}
            </div>
            <div className="h-2 flex-1 overflow-hidden rounded-full bg-muted">
              <div
                className={cn('h-full rounded-full', item.color)}
                style={{ width: `${item.progress}%` }}
              />
            </div>
            <div className="w-12 text-right text-sm font-semibold text-foreground">
              {item.value}
            </div>
          </div>
        ))}
      </div>

      <div className="rounded-2xl border border-border bg-white p-5">
        <div className="mb-4 flex items-center justify-between">
          <div className="text-sm font-semibold text-foreground">
            Recent Activity
          </div>
          <div className="text-xs text-muted-foreground">
            {activitySummary.activeDays} active days
          </div>
        </div>

        {activitySummary.recentActivity.length > 0 ? (
          <div className="space-y-3">
            {activitySummary.recentActivity.map((activity) => (
              <RecentActivityRow activity={activity} key={activity.id} />
            ))}
          </div>
        ) : (
          <p className="text-sm text-muted-foreground">
            Recent profile activity has not been recorded yet.
          </p>
        )}
      </div>
    </section>
  )
}

function AchievementsTab({
  achievements,
}: {
  achievements: ProfileAchievement[]
}) {
  return (
    <section className="grid gap-4 md:grid-cols-3">
      {achievements.length > 0 ? (
        achievements.map((achievement) => (
          <div
            className={cn(
              'rounded-2xl border bg-white p-5',
              achievement.unlocked
                ? 'border-border'
                : 'border-border opacity-50',
            )}
            key={achievement.id}
          >
            <div className="flex items-start gap-3">
              <div className={cn('text-3xl', !achievement.unlocked && 'grayscale')}>
                {getAchievementIcon(achievement.icon)}
              </div>
              <div className="min-w-0 flex-1">
                <div className="flex items-center gap-2">
                  <div className="text-sm font-semibold text-foreground">
                    {achievement.title}
                  </div>
                  {achievement.unlocked ? (
                    <CheckCircle2 className="size-3.5 shrink-0 text-emerald-500" />
                  ) : null}
                </div>
                <div className="mt-0.5 text-xs text-muted-foreground">
                  {achievement.description}
                </div>
                <div className="mt-1 text-xs text-primary">
                  {achievement.unlockedAt
                    ? formatAchievementDate(achievement.unlockedAt)
                    : `${achievement.progressValue} / ${achievement.requirementValue}`}
                </div>
              </div>
            </div>
          </div>
        ))
      ) : (
        <ProfileStateMessage
          title="No achievements yet"
          description="Achievements will appear here once they are unlocked."
        />
      )}
    </section>
  )
}

function CalendarTab({
  wordsLearnedLast84Days,
  averageWordsPerDay,
  activeDays,
}: {
  wordsLearnedLast84Days: number
  averageWordsPerDay: number
  activeDays: number
}) {
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
          {weekLabels.map((label) => (
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
                  <HeatCell count={emptyHeatmap[index] ?? 0} key={`${week}-${day}`} />
                )
              })}
            </div>
          ))}
        </div>
      </div>

      <div className="mt-4 space-y-2 text-xs text-muted-foreground">
        <div className="flex items-center gap-2">
          <Calendar className="size-[13px]" />
          {formatNumber(wordsLearnedLast84Days)} words learned in the last 84
          days · Avg. {averageWordsPerDay.toFixed(1)}/day
        </div>
        <p>
          Detailed day-by-day heatmap data is not available from the current
          profile endpoint yet. Active days recorded: {activeDays}.
        </p>
      </div>
    </section>
  )
}

function RecentActivityRow({ activity }: { activity: ProfileRecentActivity }) {
  return (
    <div className="flex items-center gap-3">
      <div className="flex size-8 items-center justify-center rounded-xl bg-muted text-xs font-bold text-foreground">
        {activity.title.charAt(0).toUpperCase()}
      </div>
      <div className="min-w-0 flex-1">
        <div className="truncate text-sm font-medium text-foreground">
          {activity.title}
        </div>
        <div className="text-xs text-muted-foreground">
          {activity.description} · {formatRelativeTime(activity.createdAt)}
        </div>
      </div>
      <div className="rounded-full bg-emerald-50 px-2 py-0.5 text-xs font-medium text-emerald-600">
        +{activity.xp} XP
      </div>
    </div>
  )
}

function ProfileStateMessage({
  title,
  description,
  tone = 'default',
}: {
  title: string
  description: string
  tone?: 'default' | 'error'
}) {
  return (
    <section className="rounded-2xl border border-border bg-white p-6 text-center shadow-sm">
      <h2
        className={cn(
          'text-lg font-bold tracking-normal',
          tone === 'error' ? 'text-destructive' : 'text-foreground',
        )}
      >
        {title}
      </h2>
      <p className="mt-2 text-sm text-muted-foreground">{description}</p>
    </section>
  )
}

function ProfileLoadingSkeleton() {
  return (
    <>
      <section className="overflow-hidden rounded-[28px] border border-border bg-white">
        <div className="h-28 bg-muted/70" />
        <div className="px-6 pb-6">
          <div className="-mt-10 mb-4 flex items-end justify-between gap-4">
            <Skeleton className="size-20 rounded-2xl border-4 border-white" />
            <Skeleton className="h-9 w-32 rounded-full" />
          </div>

          <div className="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
            <div className="space-y-3">
              <Skeleton className="h-7 w-40" />
              <Skeleton className="h-4 w-48" />
              <div className="flex gap-3">
                <Skeleton className="h-6 w-28 rounded-full" />
                <Skeleton className="h-6 w-28 rounded-full" />
              </div>
            </div>
            <Skeleton className="h-16 w-28 rounded-2xl" />
          </div>
        </div>
      </section>

      <CardSkeleton lines={3} />

      <section className="grid grid-cols-2 gap-4 md:grid-cols-4">
        {Array.from({ length: 4 }).map((_, index) => (
          <CardSkeleton className="p-4" key={index} lines={2} showIcon />
        ))}
      </section>

      <CardSkeleton lines={4} />
      <CardSkeleton lines={2} />
      <CardSkeleton lines={6} />
    </>
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

function buildStatCards(profile: ProfileResponse) {
  return [
    {
      label: 'Words Learned',
      value: formatNumber(profile.stats.wordsLearned),
      tone: 'primary',
    },
    {
      label: 'Current Streak',
      value: `${profile.user.currentStreak}d`,
      tone: 'orange',
    },
    {
      label: 'Longest Streak',
      value: `${profile.user.longestStreak}d`,
      tone: 'amber',
    },
    {
      label: 'Lessons Done',
      value: formatNumber(profile.stats.lessonsDone),
      tone: 'emerald',
    },
  ] as const
}

function buildOverviewItems(profile: ProfileResponse): OverviewItem[] {
  const unlockedCount = profile.achievements.filter(
    (achievement) => achievement.unlocked,
  ).length
  const targetBandProgress = profile.user.currentBand
    ? Math.min(
        100,
        Math.round((profile.user.currentBand / profile.user.targetBand) * 100),
      )
    : 0
  const streakProgress =
    profile.user.longestStreak > 0
      ? Math.min(
          100,
          Math.round(
            (profile.user.currentStreak / profile.user.longestStreak) * 100,
          ),
        )
      : 0
  const achievementProgress =
    profile.achievements.length > 0
      ? Math.round((unlockedCount / profile.achievements.length) * 100)
      : 0
  const activityProgress = Math.round(
    (profile.activitySummary.activeDays / 84) * 100,
  )

  return [
    {
      label: 'Target Band',
      progress: targetBandProgress,
      value: `${targetBandProgress}%`,
      color: 'bg-indigo-500',
    },
    {
      label: 'Mastery',
      progress: profile.stats.vocabularyMasteryPercentage,
      value: `${profile.stats.vocabularyMasteryPercentage}%`,
      color: 'bg-emerald-500',
    },
    {
      label: 'Streak Consistency',
      progress: streakProgress,
      value: `${streakProgress}%`,
      color: 'bg-orange-500',
    },
    {
      label: 'Achievements',
      progress: achievementProgress,
      value: `${achievementProgress}%`,
      color: 'bg-amber-500',
    },
    {
      label: 'Active Days',
      progress: activityProgress,
      value: `${profile.activitySummary.activeDays}`,
      color: 'bg-purple-500',
    },
    {
      label: 'Recent XP',
      progress: Math.min(100, unlockedCount * 10),
      value: `${profile.activitySummary.recentActivity.length}`,
      color: 'bg-blue-500',
    },
  ] as const
}

function getProfileErrorMessage(error: Error | null) {
  if (!error) {
    return null
  }

  if (error instanceof APIError) {
    return error.message
  }

  return 'Unable to load your profile right now.'
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

function getAchievementIcon(icon: string) {
  switch (icon) {
    case 'flame':
      return '🔥'
    case 'layers':
      return '📚'
    case 'circle-check':
      return '⚡'
    case 'trophy':
      return '🏆'
    case 'book-open-check':
      return '🎓'
    default:
      return '🌟'
  }
}

function formatMemberSince(value: string) {
  const date = new Date(value)
  return `Member since ${date.toLocaleDateString('en-US', {
    month: 'short',
    year: 'numeric',
  })}`
}

function formatAchievementDate(value: string) {
  const date = new Date(value)
  return date.toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
  })
}

function formatRelativeTime(value: string) {
  const date = new Date(value)
  const diffMs = date.getTime() - Date.now()

  if (Number.isNaN(date.getTime())) {
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
  return formatter.format(days, 'day')
}

function formatNumber(value: number) {
  return new Intl.NumberFormat('en-US').format(value)
}

function formatUsername(username: string | null) {
  return username ? `@${username}` : '@learner'
}
