import { Inbox } from 'lucide-react'
import type React from 'react'
import { Link } from 'react-router-dom'

import { Button } from '@/components/ui/button'
import { cn } from '@/lib/utils'

type EmptyStateProps = {
  title: string
  description: string
  className?: string
  icon?: React.ReactNode
  actionLabel?: string
  actionHref?: string
  onAction?: () => void
}

export function EmptyState({
  title,
  description,
  className,
  icon,
  actionLabel,
  actionHref,
  onAction,
}: EmptyStateProps) {
  return (
    <section
      className={cn(
        'rounded-2xl border border-border bg-white p-6 text-center shadow-sm',
        className,
      )}
    >
      <div className="mx-auto flex size-10 items-center justify-center rounded-full bg-muted text-muted-foreground">
        {icon ?? <Inbox className="size-5" aria-hidden="true" />}
      </div>
      <h2 className="mt-4 text-lg font-bold tracking-normal text-foreground">
        {title}
      </h2>
      <p className="mx-auto mt-2 max-w-md text-sm text-muted-foreground">
        {description}
      </p>
      {actionLabel ? (
        actionHref ? (
          <Button asChild className="mt-6 rounded-full">
            <Link to={actionHref}>{actionLabel}</Link>
          </Button>
        ) : (
          <Button className="mt-6 rounded-full" onClick={onAction} type="button">
            {actionLabel}
          </Button>
        )
      ) : null}
    </section>
  )
}
