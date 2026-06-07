import {
  ArrowUpRight,
  BookOpen,
  Brain,
  CheckCircle2,
  ChevronRight,
  Flame,
  Play,
  RefreshCw,
  Star,
  Target,
  TrendingUp,
  Trophy,
  Zap,
} from "lucide-react";
import { Link } from "react-router-dom";
import {
  Area,
  AreaChart,
  RadialBar,
  RadialBarChart,
  ResponsiveContainer,
  Tooltip,
  XAxis,
} from "recharts";

import { cn } from "@/lib/utils";

const weekData = [
  { day: "Mon", words: 12 },
  { day: "Tue", words: 18 },
  { day: "Wed", words: 8 },
  { day: "Thu", words: 22 },
  { day: "Fri", words: 15 },
  { day: "Sat", words: 28 },
  { day: "Sun", words: 9 },
];

const radialData = [{ name: "progress", value: 65, fill: "#4F46E5" }];

const progressCards = [
  {
    label: "Words Today",
    value: "24",
    icon: BookOpen,
    color: "bg-indigo-50 text-primary",
    change: "+6 vs yesterday",
  },
  {
    label: "Review Due",
    value: "12",
    icon: RefreshCw,
    color: "bg-amber-50 text-amber-600",
    change: "Due now",
  },
  {
    label: "Current Streak",
    value: "14d",
    icon: Flame,
    color: "bg-orange-50 text-orange-500",
    change: "On fire!",
  },
  {
    label: "Accuracy",
    value: "84%",
    icon: Target,
    color: "bg-emerald-50 text-emerald-600",
    change: "+3% this week",
  },
];

const quickActions = [
  {
    label: "Continue Learning",
    icon: Play,
    route: "/lessons/1",
    color: "border-primary bg-primary text-white hover:bg-primary/90",
    iconColor: "text-white",
    labelColor: "text-white",
  },
  {
    label: "Flashcards",
    icon: BookOpen,
    route: "/lessons/1/flashcards",
    color: "border-indigo-100 bg-indigo-50 text-primary hover:bg-indigo-100",
    iconColor: "text-primary",
    labelColor: "text-primary",
  },
  {
    label: "Daily Quiz",
    icon: Zap,
    route: "/lessons/1/quiz",
    color:
      "border-emerald-100 bg-emerald-50 text-emerald-700 hover:bg-emerald-100",
    iconColor: "text-emerald-700",
    labelColor: "text-emerald-700",
  },
  {
    label: "Review Due",
    icon: RefreshCw,
    route: "/reviews",
    color: "border-amber-100 bg-amber-50 text-amber-700 hover:bg-amber-100",
    iconColor: "text-amber-700",
    labelColor: "text-amber-700",
  },
];

const recentActivity = [
  {
    word: "Sustainability",
    topic: "Environment",
    time: "2h ago",
    result: "Easy",
    resultColor: "bg-emerald-50 text-emerald-600",
  },
  {
    word: "Proliferate",
    topic: "Technology",
    time: "3h ago",
    result: "Good",
    resultColor: "bg-blue-50 text-blue-600",
  },
  {
    word: "Deteriorate",
    topic: "Health",
    time: "5h ago",
    result: "Hard",
    resultColor: "bg-orange-50 text-orange-600",
  },
  {
    word: "Endeavour",
    topic: "Education",
    time: "Yesterday",
    result: "Easy",
    resultColor: "bg-emerald-50 text-emerald-600",
  },
  {
    word: "Mitigation",
    topic: "Environment",
    time: "Yesterday",
    result: "Again",
    resultColor: "bg-red-50 text-red-600",
  },
];

const upcomingReviews = [
  { word: "Sophisticated", due: "Now", band: 7.0, urgent: true },
  { word: "Exacerbate", due: "In 2h", band: 7.5, urgent: true },
  { word: "Innate", due: "In 4h", band: 6.5, urgent: false },
  { word: "Pragmatic", due: "Tomorrow", band: 7.0, urgent: false },
];

const achievements = [
  {
    label: "7 Day Streak",
    icon: Flame,
    color: "border-orange-200 bg-orange-50 text-orange-500",
    earned: true,
  },
  {
    label: "100 Words",
    icon: BookOpen,
    color: "border-blue-200 bg-blue-50 text-blue-600",
    earned: true,
  },
  {
    label: "First Quiz",
    icon: Zap,
    color: "border-purple-200 bg-purple-50 text-purple-600",
    earned: true,
  },
  {
    label: "Band 6.0",
    icon: Target,
    color: "border-indigo-200 bg-indigo-50 text-primary",
    earned: true,
  },
  {
    label: "30 Day Streak",
    icon: Trophy,
    color: "border-amber-200 bg-amber-50 text-amber-600",
    earned: false,
  },
  {
    label: "500 Words",
    icon: Star,
    color: "border-emerald-200 bg-emerald-50 text-emerald-600",
    earned: false,
  },
];

const vocabularyStats = [
  { label: "Learned", value: "847" },
  { label: "Mastered", value: "312" },
  { label: "In Review", value: "128" },
  { label: "New", value: "76" },
];

export function DashboardOverview() {
  return (
    <div className="space-y-6">
      <section className="flex flex-col justify-between gap-4 sm:flex-row sm:items-center">
        <div>
          <h1 className="text-2xl font-bold tracking-normal text-foreground">
            Welcome back, Alex
          </h1>
          <p className="mt-0.5 text-sm text-muted-foreground">
            You are on a roll! Keep building your streak.
          </p>
        </div>
        <div className="flex w-fit items-center gap-2 rounded-2xl border border-border bg-white px-4 py-2.5 shadow-sm">
          <Flame className="size-5 fill-orange-500 text-orange-500" />
          <div>
            <div className="text-lg font-bold leading-none text-foreground">
              14
            </div>
            <div className="text-xs text-muted-foreground">Day streak</div>
          </div>
        </div>
      </section>

      <section className="grid grid-cols-2 gap-4 lg:grid-cols-4">
        {progressCards.map((card) => (
          <div
            className="rounded-2xl border border-border bg-white p-4 shadow-sm"
            key={card.label}
          >
            <div
              className={cn(
                "mb-3 flex size-9 items-center justify-center rounded-xl",
                card.color,
              )}
            >
              <card.icon className="size-[18px]" />
            </div>
            <div className="text-2xl font-bold leading-tight text-foreground">
              {card.value}
            </div>
            <div className="mt-0.5 text-xs text-muted-foreground">
              {card.label}
            </div>
            <div className="mt-1 text-xs text-emerald-600">{card.change}</div>
          </div>
        ))}
      </section>

      <section className="grid gap-6 lg:grid-cols-3">
        <TargetBandCard />
        <WeeklyWordsCard />
      </section>

      <QuickActions />

      <section className="grid gap-6 lg:grid-cols-2">
        <RecentActivityCard />
        <DueReviewCard />
      </section>

      <section className="grid gap-6 lg:grid-cols-[1.6fr_1fr]">
        <AchievementsCard />
        <VocabularyStatsCard />
      </section>
    </div>
  );
}

function TargetBandCard() {
  return (
    <div className="rounded-2xl border border-border bg-white p-5 shadow-sm">
      <div className="mb-4 flex items-center justify-between">
        <div>
          <div className="text-sm font-semibold text-foreground">
            Target Band Progress
          </div>
          <div className="text-xs text-muted-foreground">Band 7.0 Goal</div>
        </div>
        <div className="rounded-full bg-indigo-50 px-2 py-1 text-xs font-semibold text-primary">
          65%
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
              background={{ fill: "#ede9fe" }}
              cornerRadius={8}
              dataKey="value"
            />
          </RadialBarChart>
        </ResponsiveContainer>
      </div>

      <div className="-mt-2 text-center">
        <div className="text-2xl font-bold text-primary">Band 6.5</div>
        <div className="text-xs text-muted-foreground">
          Current estimated band
        </div>
      </div>

      <div className="mt-4 h-2 overflow-hidden rounded-full bg-muted">
        <div
          className="h-full rounded-full bg-gradient-to-r from-primary to-violet-600"
          style={{ width: "65%" }}
        />
      </div>
      <div className="mt-1 flex justify-between text-xs text-muted-foreground">
        <span>Band 5.0</span>
        <span>Band 7.0 target</span>
      </div>
    </div>
  );
}

function WeeklyWordsCard() {
  return (
    <div className="rounded-2xl border border-border bg-white p-5 shadow-sm lg:col-span-2">
      <div className="mb-4 flex items-center justify-between gap-4">
        <div className="text-sm font-semibold text-foreground">
          Words Learned This Week
        </div>
        <div className="flex items-center gap-1 text-xs text-emerald-600">
          <TrendingUp className="size-[13px]" />
          +18% vs last week
        </div>
      </div>

      <ResponsiveContainer height={140} width="100%">
        <AreaChart data={weekData}>
          <defs>
            <linearGradient id="wordGrad" x1="0" x2="0" y1="0" y2="1">
              <stop offset="5%" stopColor="#4F46E5" stopOpacity={0.15} />
              <stop offset="95%" stopColor="#4F46E5" stopOpacity={0} />
            </linearGradient>
          </defs>
          <XAxis
            axisLine={false}
            dataKey="day"
            tick={{ fontSize: 11, fill: "#6b6b80" }}
            tickLine={false}
          />
          <Tooltip
            contentStyle={{
              border: "1px solid #e5e7eb",
              borderRadius: 12,
              fontSize: 12,
            }}
          />
          <Area
            dataKey="words"
            dot={{ fill: "#4F46E5", r: 3 }}
            fill="url(#wordGrad)"
            stroke="#4F46E5"
            strokeWidth={2.5}
            type="monotone"
          />
        </AreaChart>
      </ResponsiveContainer>
    </div>
  );
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
              "flex cursor-pointer items-center gap-2.5 rounded-xl border px-4 py-3 text-sm font-medium transition-all hover:-translate-y-0.5",
              action.color,
            )}
            key={action.label}
            to={action.route}
          >
            <action.icon className={cn("size-[18px]", action.iconColor)} />
            <span className={action.labelColor}>{action.label}</span>
          </Link>
        ))}
      </div>
    </section>
  );
}

function RecentActivityCard() {
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
      <div className="space-y-3">
        {recentActivity.map((activity) => (
          <div className="flex items-center gap-3" key={activity.word}>
            <div className="flex size-8 items-center justify-center rounded-xl bg-muted text-xs font-bold text-foreground">
              {activity.word[0]}
            </div>
            <div className="min-w-0 flex-1">
              <div className="truncate text-sm font-medium text-foreground">
                {activity.word}
              </div>
              <div className="text-xs text-muted-foreground">
                {activity.topic} - {activity.time}
              </div>
            </div>
            <div
              className={cn(
                "rounded-full px-2 py-0.5 text-xs font-medium",
                activity.resultColor,
              )}
            >
              {activity.result}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}

function DueReviewCard() {
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
      <div className="space-y-3">
        {upcomingReviews.map((review) => (
          <div
            className={cn(
              "flex items-center gap-3 rounded-xl p-3",
              review.urgent ? "border border-red-100 bg-red-50" : "bg-muted/30",
            )}
            key={review.word}
          >
            <div className="flex size-8 items-center justify-center rounded-xl border border-border bg-white text-xs font-bold text-foreground">
              {review.word[0]}
            </div>
            <div className="min-w-0 flex-1">
              <div className="truncate text-sm font-medium text-foreground">
                {review.word}
              </div>
              <div className="text-xs text-muted-foreground">
                Band {review.band}+
              </div>
            </div>
            <div
              className={cn(
                "text-xs font-medium",
                review.urgent ? "text-red-600" : "text-muted-foreground",
              )}
            >
              {review.due}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}

function AchievementsCard() {
  return (
    <div className="rounded-2xl border border-border bg-white p-5 shadow-sm">
      <div className="mb-4 flex items-center justify-between">
        <div className="text-sm font-semibold text-foreground">
          Recent Achievements
        </div>
        <Link
          className="cursor-pointer text-xs text-primary hover:underline"
          to="/profile"
        >
          View all
        </Link>
      </div>
      <div className="flex gap-4 overflow-x-auto pb-1">
        {achievements.map((achievement) => (
          <div
            className={cn(
              "flex w-20 shrink-0 flex-col items-center gap-1.5 rounded-2xl border p-3 text-center",
              achievement.color,
              !achievement.earned && "grayscale opacity-40",
            )}
            key={achievement.label}
          >
            <achievement.icon className="size-6" />
            <div className="text-xs font-medium leading-tight text-foreground">
              {achievement.label}
            </div>
            {achievement.earned ? (
              <CheckCircle2 className="size-3 text-emerald-500" />
            ) : null}
          </div>
        ))}
      </div>
    </div>
  );
}

function VocabularyStatsCard() {
  return (
    <div className="rounded-2xl border border-border bg-white p-5 shadow-sm">
      <div className="mb-4 flex items-center justify-between">
        <div className="text-sm font-semibold text-foreground">
          Vocabulary Statistics
        </div>
        <Brain className="size-4 text-primary" />
      </div>
      <div className="grid grid-cols-2 gap-3">
        {vocabularyStats.map((stat) => (
          <div className="rounded-xl bg-muted/40 p-3" key={stat.label}>
            <div className="text-xl font-bold text-foreground">
              {stat.value}
            </div>
            <div className="text-xs text-muted-foreground">{stat.label}</div>
          </div>
        ))}
      </div>
    </div>
  );
}
