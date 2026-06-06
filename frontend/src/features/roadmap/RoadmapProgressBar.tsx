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
    <div className="h-2 overflow-hidden rounded-full bg-[#e8e8f2]">
      <div
        className={cn(
          'h-full rounded-full',
          variant === 'primary' && 'bg-primary',
          variant === 'success' && 'bg-success',
          variant === 'muted' && 'bg-[#e8e8f2]',
        )}
        style={{ width: `${value}%` }}
      />
    </div>
  )
}
