import { Skeleton } from '@/components/ui/skeleton'
import { cn } from '@/lib/utils'

type PageSkeletonProps = {
  className?: string
  showHeader?: boolean
  lines?: number
}

export function PageSkeleton({
  className,
  showHeader = true,
  lines = 3,
}: PageSkeletonProps) {
  return (
    <section className={cn('rounded-2xl border border-border bg-white p-6 shadow-sm', className)}>
      {showHeader ? (
        <div className="mb-5 space-y-3">
          <Skeleton className="h-6 w-40" />
          <Skeleton className="h-4 w-72 max-w-full" />
        </div>
      ) : null}
      <div className="space-y-3">
        {Array.from({ length: lines }).map((_, index) => (
          <Skeleton
            className={cn('h-4', index === lines - 1 ? 'w-2/3' : 'w-full')}
            key={index}
          />
        ))}
      </div>
    </section>
  )
}
