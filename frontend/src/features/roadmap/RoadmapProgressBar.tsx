import { cn } from '@/lib/utils'

type RoadmapProgressBarProps = {
  value: number
  variant?: 'primary' | 'success' | 'muted'
}

export function RoadmapProgressBar({
  value,
  variant = 'primary',
}: RoadmapProgressBarProps) {
  return (
    <div className="h-1.5 overflow-hidden rounded-full bg-muted">
      <div
        className={cn(
          'h-full rounded-full',
          variant === 'primary' && 'bg-gradient-to-r from-primary to-violet-600',
          variant === 'success' && 'bg-success',
          variant === 'muted' && 'bg-muted',
        )}
        style={{ width: `${value}%` }}
      />
    </div>
  )
}
