import {
  ArrowUpRight,
  BookOpen,
  Brain,
  ChevronRight,
  Flame,
  Play,
  RefreshCw,
  Target,
  Trophy,
  Zap,
} from 'lucide-react'
import { Link } from 'react-router-dom'
import {
  RadialBar,
  RadialBarChart,
  ResponsiveContainer,
} from 'recharts'

import { CardSkeleton } from '@/components/state/CardSkeleton'
import { EmptyState } from '@/components/state/EmptyState'
import { ErrorState } from '@/components/state/ErrorState'
import { ListSkeleton } from '@/components/state/ListSkeleton'
import { Skeleton } from '@/components/ui/skeleton'
import type { DashboardViewModel } from '@/features/dashboard/mapDashboardSummary'
import { cn } from '@/lib/utils'

const quickActions = [
  {
    label: 'Continue Learning',
    icon: Play,
    route: '/lessons/1',
    color: 'border-primary bg-primary text-white hover:bg-primary/90',
    iconColor: 'text-white',
    labelColor: 'text-white',
  },
  {
    label: 'Flashcards',
    icon: BookOpen,
    route: '/lessons/1/flashcards',
    color: 'border-indigo-100 bg-indigo-50 text-primary hover:bg-indigo-100',
    iconColor: 'text-primary',
    labelColor: 'text-primary',
  },
  {
    label: 'Daily Quiz',
    icon: Zap,
    route: '/lessons/1/quiz',
    color:
      'border-emerald-100 bg-emerald-50 text-emerald-700 hover:bg-emerald-100',
    iconColor: 'text-emerald-700',
    labelColor: 'text-emerald-700',
  },
  {
    label: 'Review Due',
    icon: RefreshCw,
    route: '/reviews',
    color: 'border-amber-100 bg-amber-50 text-amber-700 hover:bg-amber-100',
    iconColor: 'text-amber-700',
    labelColor: 'text-amber-700',
  },
] as const

type DashboardOverviewProps = {
  dashboard: DashboardViewModel
  isLoading: boolean
  errorMessage: string | null
  isEmpty: boolean
  onRetry?: () => void
}

export function DashboardOverview({
  dashboard,
  isLoading,
  errorMessage,
  isEmpty,
  onRetry,
}: DashboardOverviewProps) {
  if (isLoading) {
    return <DashboardLoadingSkeleton />
  }

  if (errorMessage) {
    return (
      <ErrorState
        actionHref="/roadmap"
        actionLabel="Go to roadmap"
        description={errorMessage}
        onRetry={onRetry}
        title="Dashboard unavailable"
      />
    )
  }

  if (isEmpty) {
    return (
      <EmptyState
        actionHref="/roadmap"
        actionLabel="Go to roadmap"
        description="We could not find any dashboard summary yet, but you can keep learning from your roadmap."
        title="No dashboard data yet"
      />
    )
  }

  return (
    <div className="space-y-6">
      <section className="flex flex-col justify-between gap-4 sm:flex-row sm:items-center">
        <div>
          <h1 className="text-2xl font-bold tracking-normal text-foreground">
            Welcome back, {dashboard.greetingName}
          </h1>
          <p className="mt-0.5 text-sm text-muted-foreground">
            {dashboard.subtitle}
          </p>
        </div>
        <div className="flex w-fit items-center gap-2 rounded-2xl border border-border bg-white px-4 py-2.5 shadow-sm">
          <Flame className="size-5 fill-orange-500 text-orange-500" />
          <div>
            <div className="text-lg font-bold leading-none text-foreground">
              {dashboard.streakDays}
            </div>
            <div className="text-xs text-muted-foreground">Day streak</div>
          </div>
        </div>
      </section>

      <section className="grid grid-cols-2 gap-4 lg:grid-cols-4">
        {dashboard.progressCards.map((card) => (
          <div
            className="rounded-2xl border border-border bg-white p-4 shadow-sm"
            key={card.label}
          >
            <div
              className={cn(
                'mb-3 flex size-9 items-center justify-center rounded-xl',
                getProgressCardToneClasses(card.tone),
              )}
            >
              {card.label === 'Words Today' ? (
                <BookOpen className="size-[18px]" />
              ) : null}
              {card.label === 'Review Due' ? (
                <RefreshCw className="size-[18px]" />
              ) : null}
              {card.label === 'Current Streak' ? (
                <Flame className="size-[18px]" />
              ) : null}
              {card.label === 'Mastery' ? (
                <Target className="size-[18px]" />
              ) : null}
            </div>
            <div className="text-2xl font-bold leading-tight text-foreground">
              {card.value}
            </div>
            <div className="mt-0.5 text-xs text-muted-foreground">
              {card.label}
            </div>
            <div className="mt-1 text-xs text-emerald-600">
              {card.supportText}
            </div>
          </div>
        ))}
      </section>

      <section className="grid gap-6 lg:grid-cols-3">
        <TargetBandCard dashboard={dashboard} />
        <TodayProgressCard dashboard={dashboard} />
      </section>

      <QuickActions />

      <section className="grid gap-6 lg:grid-cols-2">
        <RecentActivityCard dashboard={dashboard} />
        <DueReviewCard dashboard={dashboard} />
      </section>

      <section className="grid gap-6 lg:grid-cols-[1.6fr_1fr]">
        <XPLevelCard dashboard={dashboard} />
        <VocabularyStatsCard dashboard={dashboard} />
      </section>
    </div>
  )
}

function DashboardLoadingSkeleton() {
  return (
    <div className="space-y-6">
      <section className="flex flex-col justify-between gap-4 sm:flex-row sm:items-center">
        <div className="space-y-2">
          <Skeleton className="h-8 w-64" />
          <Skeleton className="h-4 w-72 max-w-full" />
        </div>
        <div className="rounded-2xl border border-border bg-white px-4 py-2.5 shadow-sm">
          <Skeleton className="h-10 w-24" />
        </div>
      </section>

      <section className="grid grid-cols-2 gap-4 lg:grid-cols-4">
        {Array.from({ length: 4 }).map((_, index) => (
          <CardSkeleton className="p-4" key={index} lines={2} showIcon />
        ))}
      </section>

      <section className="grid gap-6 lg:grid-cols-3">
        <CardSkeleton className="min-h-[320px]" lines={5} />
        <ListSkeleton className="lg:col-span-2" count={1} itemClassName="min-h-[320px]" />
      </section>

      <section>
        <Skeleton className="mb-3 h-4 w-28" />
        <div className="grid grid-cols-2 gap-3 md:grid-cols-4">
          {Array.from({ length: 4 }).map((_, index) => (
            <Skeleton className="h-14 rounded-xl" key={index} />
          ))}
        </div>
      </section>

      <section className="grid gap-6 lg:grid-cols-2">
        <ListSkeleton count={2} itemClassName="min-h-[220px]" />
      </section>

      <section className="grid gap-6 lg:grid-cols-[1.6fr_1fr]">
        <ListSkeleton count={2} itemClassName="min-h-[220px]" />
      </section>
    </div>
  )
}

function TargetBandCard({ dashboard }: { dashboard: DashboardViewModel }) {
  const radialData = [
    {
      name: 'progress',
      value: dashboard.targetBand.progressPercentage,
      fill: '#4F46E5',
    },
  ]

  return (
    <div className="rounded-2xl border border-border bg-white p-5 shadow-sm">
      <div className="mb-4 flex items-center justify-between">
        <div>
          <div className="text-sm font-semibold text-foreground">
            Target Band Progress
          </div>
          <div className="text-xs text-muted-foreground">
            {dashboard.targetBand.targetLabel}
          </div>
        </div>
        <div className="rounded-full bg-indigo-50 px-2 py-1 text-xs font-semibold text-primary">
          {dashboard.targetBand.badgeLabel}
        </div>
      </div>

      <div className="flex h-[140px] items-center justify-center">
        <ResponsiveContainer height="100%" width="100%">
          <RadialBarChart
            cx="50%"
            cy="50%"
            data={radialData}
            endAngle={-270}
            innerRadius="65%"
            outerRadius="90%"
            startAngle={90}
          >
            <RadialBar
              background={{ fill: '#ede9fe' }}
              cornerRadius={8}
              dataKey="value"
            />
          </RadialBarChart>
        </ResponsiveContainer>
      </div>

      <div className="-mt-2 text-center">
        <div className="text-2xl font-bold text-primary">
          {dashboard.targetBand.currentBandLabel}
        </div>
        <div className="text-xs text-muted-foreground">
          {dashboard.targetBand.currentBandCaption}
        </div>
      </div>

      <div className="mt-4 h-2 overflow-hidden rounded-full bg-muted">
        <div
          className="h-full rounded-full bg-gradient-to-r from-primary to-violet-600"
          style={{ width: `${dashboard.targetBand.progressPercentage}%` }}
        />
      </div>
      <div className="mt-1 flex justify-between text-xs text-muted-foreground">
        <span>{dashboard.targetBand.startLabel}</span>
        <span>{dashboard.targetBand.targetLabel}</span>
      </div>
    </div>
  )
}

function TodayProgressCard({ dashboard }: { dashboard: DashboardViewModel }) {
  return (
    <div className="rounded-2xl border border-border bg-white p-5 shadow-sm lg:col-span-2">
      <div className="mb-4 flex items-center justify-between gap-4">
        <div className="text-sm font-semibold text-foreground">
          Today&apos;s Progress
        </div>
        <div className="flex items-center gap-1 text-xs text-emerald-600">
          <Zap className="size-[13px]" />
          Live from /dashboard/summary
        </div>
      </div>

      <div className="grid gap-3 sm:grid-cols-2 xl:grid-cols-5">
        {dashboard.todayProgress.map((item) => (
          <div className="rounded-xl bg-muted/40 p-4" key={item.label}>
            <div className="text-xl font-bold text-foreground">{item.value}</div>
            <div className="mt-1 text-xs text-muted-foreground">
              {item.label}
            </div>
          </div>
        ))}
      </div>
    </div>
  )
}

function QuickActions() {
  return (
    <section>
      <div className="mb-3 text-sm font-semibold text-foreground">
        Quick Actions
      </div>
      <div className="grid grid-cols-2 gap-3 md:grid-cols-4">
        {quickActions.map((action) => (
          <Link
            className={cn(
              'flex cursor-pointer items-center gap-2.5 rounded-xl border px-4 py-3 text-sm font-medium transition-all hover:-translate-y-0.5',
              action.color,
            )}
            key={action.label}
            to={action.route}
          >
            <action.icon className={cn('size-[18px]', action.iconColor)} />
            <span className={action.labelColor}>{action.label}</span>
          </Link>
        ))}
      </div>
    </section>
  )
}

function RecentActivityCard({ dashboard }: { dashboard: DashboardViewModel }) {
  return (
    <div className="rounded-2xl border border-border bg-white p-5 shadow-sm">
      <div className="mb-4 flex items-center justify-between">
        <div className="text-sm font-semibold text-foreground">
          Recent Activity
        </div>
        <Link
          className="flex cursor-pointer items-center gap-1 text-xs text-primary hover:underline"
          to="/vocabulary"
        >
          View all
          <ChevronRight className="size-3" />
        </Link>
      </div>

      {dashboard.recentActivity.length > 0 ? (
        <div className="space-y-3">
          {dashboard.recentActivity.map((activity) => (
            <div className="flex items-center gap-3" key={activity.id}>
              <div className="flex size-8 items-center justify-center rounded-xl bg-muted text-xs font-bold text-foreground">
                {activity.initial}
              </div>
              <div className="min-w-0 flex-1">
                <div className="truncate text-sm font-medium text-foreground">
                  {activity.title}
                </div>
                <div className="text-xs text-muted-foreground">
                  {activity.subtitle}
                </div>
              </div>
              <div
                className={cn(
                  'rounded-full px-2 py-0.5 text-xs font-medium',
                  getActivityBadgeClasses(activity.badgeTone),
                )}
              >
                {activity.badgeLabel}
              </div>
            </div>
          ))}
        </div>
      ) : (
        <div className="rounded-xl bg-muted/30 p-4 text-sm text-muted-foreground">
          Your latest lesson completions, reviews, and quiz results will appear
          here once activity is recorded.
        </div>
      )}
    </div>
  )
}

function DueReviewCard({ dashboard }: { dashboard: DashboardViewModel }) {
  return (
    <div className="rounded-2xl border border-border bg-white p-5 shadow-sm">
      <div className="mb-4 flex items-center justify-between">
        <div className="text-sm font-semibold text-foreground">
          Due for Review
        </div>
        <Link
          className="flex cursor-pointer items-center gap-1 text-xs text-primary hover:underline"
          to="/reviews"
        >
          Start review
          <ArrowUpRight className="size-3" />
        </Link>
      </div>

      <div className="rounded-2xl border border-border bg-muted/20 p-4">
        <div className="text-3xl font-bold text-foreground">
          {dashboard.reviewDue.count}
        </div>
        <div className="mt-1 text-sm font-semibold text-foreground">
          {dashboard.reviewDue.title}
        </div>
        <p className="mt-2 text-sm text-muted-foreground">
          {dashboard.reviewDue.description}
        </p>
      </div>
    </div>
  )
}

function XPLevelCard({ dashboard }: { dashboard: DashboardViewModel }) {
  return (
    <div className="rounded-2xl border border-border bg-white p-5 shadow-sm">
      <div className="mb-4 flex items-center justify-between">
        <div>
          <div className="text-sm font-semibold text-foreground">
            XP &amp; Level Summary
          </div>
          <div className="text-xs text-muted-foreground">
            {dashboard.xp.levelTitle}
          </div>
        </div>
        <div className="flex size-10 items-center justify-center rounded-xl bg-amber-50 text-amber-600">
          <Trophy className="size-5" />
        </div>
      </div>

      <div className="grid gap-4 md:grid-cols-3">
        <div className="rounded-xl bg-muted/40 p-4">
          <div className="text-xs text-muted-foreground">Current Level</div>
          <div className="mt-1 text-2xl font-bold text-foreground">
            {dashboard.xp.level}
          </div>
        </div>
        <div className="rounded-xl bg-muted/40 p-4">
          <div className="text-xs text-muted-foreground">Total XP</div>
          <div className="mt-1 text-2xl font-bold text-foreground">
            {dashboard.xp.totalXp}
          </div>
        </div>
        <div className="rounded-xl bg-muted/40 p-4">
          <div className="text-xs text-muted-foreground">Next Level In</div>
          <div className="mt-1 text-2xl font-bold text-foreground">
            {dashboard.xp.xpToNextLevel} XP
          </div>
        </div>
      </div>

      <div className="mt-4 h-2 overflow-hidden rounded-full bg-muted">
        <div
          className="h-full rounded-full bg-gradient-to-r from-amber-500 to-orange-500"
          style={{ width: `${dashboard.xp.progressPercentage}%` }}
        />
      </div>
      <div className="mt-2 flex items-center justify-between text-xs text-muted-foreground">
        <span>Level progress</span>
        <span>{dashboard.xp.progressPercentage}%</span>
      </div>
    </div>
  )
}

function VocabularyStatsCard({ dashboard }: { dashboard: DashboardViewModel }) {
  return (
    <div className="rounded-2xl border border-border bg-white p-5 shadow-sm">
      <div className="mb-4 flex items-center justify-between">
        <div className="text-sm font-semibold text-foreground">
          Vocabulary Statistics
        </div>
        <Brain className="size-4 text-primary" />
      </div>
      <div className="grid grid-cols-2 gap-3">
        {dashboard.vocabularyStats.map((stat) => (
          <div className="rounded-xl bg-muted/40 p-3" key={stat.label}>
            <div className="text-xl font-bold text-foreground">
              {stat.value}
            </div>
            <div className="text-xs text-muted-foreground">{stat.label}</div>
          </div>
        ))}
      </div>
    </div>
  )
}

function getProgressCardToneClasses(
  tone: DashboardViewModel['progressCards'][number]['tone'],
) {
  switch (tone) {
    case 'warning':
      return 'bg-amber-50 text-amber-600'
    case 'success':
      return 'bg-emerald-50 text-emerald-600'
    case 'orange':
      return 'bg-orange-50 text-orange-500'
    case 'primary':
    default:
      return 'bg-indigo-50 text-primary'
  }
}

function getActivityBadgeClasses(
  tone: DashboardViewModel['recentActivity'][number]['badgeTone'],
) {
  switch (tone) {
    case 'success':
      return 'bg-emerald-50 text-emerald-600'
    case 'warning':
      return 'bg-amber-50 text-amber-600'
    case 'danger':
      return 'bg-red-50 text-red-600'
    case 'primary':
      return 'bg-blue-50 text-blue-600'
    case 'neutral':
    default:
      return 'bg-muted text-muted-foreground'
  }
}
