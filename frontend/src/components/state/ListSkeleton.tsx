import { cn } from '@/lib/utils'

import { CardSkeleton } from './CardSkeleton'

type ListSkeletonProps = {
  className?: string
  count?: number
  itemClassName?: string
}

export function ListSkeleton({
  className,
  count = 3,
  itemClassName,
}: ListSkeletonProps) {
  return (
    <div className={cn('grid gap-4', className)}>
      {Array.from({ length: count }).map((_, index) => (
        <CardSkeleton className={itemClassName} key={index} />
      ))}
    </div>
  )
}
