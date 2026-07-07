import { AlertCircle } from 'lucide-react'
import { Link } from 'react-router-dom'

import { Button } from '@/components/ui/button'
import { cn } from '@/lib/utils'

type ErrorStateProps = {
  title: string
  description: string
  className?: string
  retryLabel?: string
  onRetry?: () => void
  actionLabel?: string
  actionHref?: string
}

export function ErrorState({
  title,
  description,
  className,
  retryLabel = 'Try again',
  onRetry,
  actionLabel,
  actionHref,
}: ErrorStateProps) {
  return (
    <section
      className={cn(
        'rounded-2xl border border-border bg-white p-6 text-center shadow-sm',
        className,
      )}
    >
      <AlertCircle
        className="mx-auto size-8 text-destructive"
        aria-hidden="true"
      />
      <h2 className="mt-4 text-lg font-bold tracking-normal text-foreground">
        {title}
      </h2>
      <p className="mx-auto mt-2 max-w-md text-sm text-muted-foreground">
        {description}
      </p>
      <div className="mt-6 flex flex-col justify-center gap-3 sm:flex-row">
        {onRetry ? (
          <Button className="rounded-full" onClick={onRetry} type="button">
            {retryLabel}
          </Button>
        ) : null}
        {actionLabel && actionHref ? (
          <Button asChild className="rounded-full" variant="outline">
            <Link to={actionHref}>{actionLabel}</Link>
          </Button>
        ) : null}
      </div>
    </section>
  )
}
