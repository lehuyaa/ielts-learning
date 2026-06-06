import {
  BookOpen,
  ChevronRight,
  Clock3,
  GraduationCap,
  Star,
  Target,
  Trophy,
} from 'lucide-react'
import { Link } from 'react-router-dom'

import { Button } from '@/components/ui/button'

import type { MockLesson } from './mockLesson'

type LessonSidebarProps = {
  lesson: MockLesson
  lessonId: string
  learnedCount: number
  progressPercentage: number
}

export function LessonSidebar({
  lesson,
  lessonId,
  learnedCount,
  progressPercentage,
}: LessonSidebarProps) {
  return (
    <aside className="grid gap-6 lg:sticky lg:top-8">
      <section className="rounded-2xl border border-[#e6e6f3] bg-white p-6 shadow-sm">
        <h2 className="text-xl font-bold tracking-normal text-[#10111f]">
          Lesson Info
        </h2>
        <div className="mt-6 grid gap-6">
          <LessonInfoRow
            icon={BookOpen}
            label="Words in Lesson"
            value={String(lesson.vocabulary.length)}
          />
          <LessonInfoRow
            icon={Clock3}
            label="Est. Time"
            value={`~${lesson.estimatedMinutes} min`}
          />
          <LessonInfoRow icon={Target} label="Band Range" value={lesson.bandRange} />
          <LessonInfoRow
            icon={Star}
            label="Required Score"
            value={`${lesson.requiredScore}%`}
          />
          <LessonInfoRow
            icon={Trophy}
            label="XP Available"
            value={`${lesson.xpReward} XP`}
          />
        </div>
      </section>

      <section className="rounded-2xl border border-[#e6e6f3] bg-white p-6 shadow-sm">
        <h2 className="text-xl font-bold tracking-normal text-[#10111f]">
          Your Progress
        </h2>
        <div className="mt-6 grid justify-items-center gap-5 text-center">
          <ProgressRing progressPercentage={progressPercentage} />
          <p className="text-base font-medium text-[#676982]">
            {learnedCount} of {lesson.vocabulary.length} words learned
          </p>
          <div className="grid w-full gap-3">
            <Button asChild className="h-12 rounded-full text-base">
              <Link to={`/lessons/${lessonId}/flashcards`}>
                Start Flashcards
                <ChevronRight aria-hidden="true" />
              </Link>
            </Button>
            <Button
              asChild
              className="h-12 rounded-full text-base"
              variant="outline"
            >
              <Link to={`/lessons/${lessonId}/quiz`}>
                <GraduationCap aria-hidden="true" />
                Start Quiz
              </Link>
            </Button>
          </div>
        </div>
      </section>

      <section className="rounded-2xl bg-gradient-to-br from-[#5b53f0] to-[#7a32f0] p-6 text-white shadow-sm">
        <p className="text-base font-bold text-white/75">Lesson Score</p>
        <p className="mt-2 text-5xl font-bold tracking-normal">
          {lesson.lessonScore}
        </p>
        <p className="mt-3 text-base font-medium text-white/75">
          Personal best: {lesson.personalBest}
        </p>
        <div className="mt-6 h-2 overflow-hidden rounded-full bg-white/20">
          <div className="h-full w-4/5 rounded-full bg-white" />
        </div>
      </section>
    </aside>
  )
}

type LessonInfoRowProps = {
  icon: typeof BookOpen
  label: string
  value: string
}

function LessonInfoRow({ icon: Icon, label, value }: LessonInfoRowProps) {
  return (
    <div className="flex items-center gap-4">
      <span className="grid size-10 shrink-0 place-items-center rounded-full bg-[#eeeef8] text-[#7b7e92]">
        <Icon className="size-5" aria-hidden="true" />
      </span>
      <div>
        <p className="text-base font-medium text-[#676982]">{label}</p>
        <p className="mt-1 text-xl font-bold tracking-normal text-[#10111f]">
          {value}
        </p>
      </div>
    </div>
  )
}

type ProgressRingProps = {
  progressPercentage: number
}

function ProgressRing({ progressPercentage }: ProgressRingProps) {
  return (
    <div
      className="grid size-32 place-items-center rounded-full"
      style={{
        background: `conic-gradient(#5147e8 ${progressPercentage}%, #ecebff 0)`,
      }}
    >
      <div className="grid size-24 place-items-center rounded-full bg-white">
        <div>
          <p className="text-3xl font-bold tracking-normal text-[#10111f]">
            {Math.round(progressPercentage)}%
          </p>
          <p className="text-base font-medium text-[#676982]">done</p>
        </div>
      </div>
    </div>
  )
}
