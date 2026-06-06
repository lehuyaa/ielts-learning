import { CheckCircle2, Clock3, Lock, PlayCircle } from 'lucide-react'

import type {
  RoadmapLesson,
  RoadmapLessonStatus,
} from '@/features/roadmap/mockRoadmap'
import { cn } from '@/lib/utils'

type RoadmapLessonNodeProps = {
  lesson: RoadmapLesson
}

const statusConfig: Record<
  RoadmapLessonStatus,
  {
    label: string
    icon: typeof CheckCircle2
    className: string
    iconClassName: string
  }
> = {
  completed: {
    label: 'Completed',
    icon: CheckCircle2,
    className: 'border-success/40 bg-emerald-50 text-slate-950',
    iconClassName: 'text-success',
  },
  'in-progress': {
    label: 'In progress',
    icon: Clock3,
    className: 'border-warning/50 bg-amber-50 text-slate-950',
    iconClassName: 'text-warning',
  },
  unlocked: {
    label: 'Unlocked',
    icon: PlayCircle,
    className: 'border-primary/40 bg-indigo-50 text-slate-950',
    iconClassName: 'text-primary',
  },
  locked: {
    label: 'Locked',
    icon: Lock,
    className: 'border-border bg-muted text-muted-foreground opacity-75',
    iconClassName: 'text-muted-foreground',
  },
}

export function RoadmapLessonNode({ lesson }: RoadmapLessonNodeProps) {
  const status = statusConfig[lesson.status]
  const StatusIcon = status.icon

  return (
    <div
      className={cn(
        'rounded-xl border p-4 transition-colors',
        status.className,
      )}
    >
      <div className="flex items-start gap-3">
        <StatusIcon
          className={cn('mt-0.5 size-5 shrink-0', status.iconClassName)}
          aria-hidden="true"
        />
        <div className="min-w-0">
          <div className="flex flex-wrap items-center gap-2">
            <h4 className="font-semibold tracking-normal">{lesson.title}</h4>
            <span className="rounded-full bg-background px-2 py-1 text-xs font-semibold text-muted-foreground">
              {status.label}
            </span>
          </div>
          <p className="mt-2 text-sm text-muted-foreground">
            {lesson.wordCount} words · {lesson.duration}
          </p>
        </div>
      </div>
    </div>
  )
}
