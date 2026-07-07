import { Skeleton } from '@/components/ui/skeleton'
import { cn } from '@/lib/utils'

type CardSkeletonProps = {
  className?: string
  lines?: number
  showIcon?: boolean
}

export function CardSkeleton({
  className,
  lines = 3,
  showIcon = false,
}: CardSkeletonProps) {
  return (
    <div className={cn('rounded-2xl border border-border bg-white p-5 shadow-sm', className)}>
      {showIcon ? <Skeleton className="mb-4 size-10 rounded-xl" /> : null}
      <Skeleton className="h-5 w-32" />
      <div className="mt-3 space-y-2">
        {Array.from({ length: lines }).map((_, index) => (
          <Skeleton
            className={cn('h-4', index === lines - 1 ? 'w-2/3' : 'w-full')}
            key={index}
          />
        ))}
      </div>
    </div>
  )
}
