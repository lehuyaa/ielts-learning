import { CheckCircle2, TrendingUp, Target } from 'lucide-react'
import type React from 'react'

import type { TopicDetailViewModel } from '@/types/topic'

type TopicHeroProps = {
  topic: TopicDetailViewModel
}

export function TopicHero({ topic }: TopicHeroProps) {
  return (
    <>
      <section className="rounded-2xl border border-[#bfdbff] bg-[#eff6ff] px-6 py-6 md:px-7">
        <div className="flex flex-col gap-4 md:flex-row md:items-center">
          <span className="text-5xl leading-none" aria-hidden="true">
            {topic.icon}
          </span>
          <div className="min-w-0">
            <h2 className="text-2xl font-bold tracking-normal text-[#10111f] md:text-3xl">
              {topic.title} Vocabulary
            </h2>
            <p className="mt-3 text-base font-medium text-[#676982]">
              {topic.description}
            </p>
            <div className="mt-5 flex flex-wrap items-center gap-x-5 gap-y-2 text-base font-medium text-[#676982]">
              <span className="inline-flex items-center gap-3 font-bold text-[#10111f]">
                <span className="size-4 rounded-full bg-[#6385ff]" />
                {topic.band}
              </span>
              <span className="hidden h-6 w-px bg-[#dfe3f4] sm:inline-block" />
              <span>{topic.totalLessons} lessons</span>
              <span className="hidden h-6 w-px bg-[#dfe3f4] sm:inline-block" />
              <span className="font-bold text-[#2fa878]">
                {topic.completedLessons} completed
              </span>
            </div>
          </div>
        </div>
      </section>

      <section className="grid gap-5 md:grid-cols-3">
        <TopicStatCard
          icon={<Target className="size-5 text-primary" aria-hidden="true" />}
          iconClassName="bg-[#eeedff]"
          label="Progress"
          value={`${topic.progressPercentage}%`}
        />
        <TopicStatCard
          icon={
            <CheckCircle2 className="size-5 text-success" aria-hidden="true" />
          }
          iconClassName="bg-[#e8fff3]"
          label="Completed"
          value={`${topic.completedLessons}/${topic.totalLessons}`}
        />
        <TopicStatCard
          icon={<TrendingUp className="size-5 text-warning" aria-hidden="true" />}
          iconClassName="bg-[#fff8df]"
          label="Total XP"
          value={topic.totalXP.toLocaleString()}
        />
      </section>
    </>
  )
}

type TopicStatCardProps = {
  icon: React.ReactNode
  iconClassName: string
  label: string
  value: string
}

function TopicStatCard({
  icon,
  iconClassName,
  label,
  value,
}: TopicStatCardProps) {
  return (
    <article className="flex items-center gap-4 rounded-2xl border border-[#e6e6f3] bg-white p-5 shadow-sm">
      <span className={`grid size-10 place-items-center rounded-full ${iconClassName}`}>
        {icon}
      </span>
      <div>
        <p className="text-base font-medium text-[#676982]">{label}</p>
        <p className="mt-1 text-xl font-bold tracking-normal text-[#10111f]">
          {value}
        </p>
      </div>
    </article>
  )
}
