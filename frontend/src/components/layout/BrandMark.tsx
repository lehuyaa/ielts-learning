import { BookOpen } from 'lucide-react'
import { Link } from 'react-router-dom'

import { cn } from '@/lib/utils'

type BrandMarkProps = {
  className?: string
}

export function BrandMark({ className }: BrandMarkProps) {
  return (
    <Link className={cn('flex items-center gap-3 font-semibold', className)} to="/">
      <span className="grid size-10 place-items-center rounded-2xl bg-primary text-primary-foreground">
        <BookOpen className="size-5" aria-hidden="true" />
      </span>
      <span>LexPath</span>
    </Link>
  )
}
