import { Bell, Search } from 'lucide-react'

import { Button } from '@/components/ui/button'
import { useCurrentUser } from '@/features/auth/hooks/useCurrentUser'

type TopbarProps = {
  title: string
}

export function Topbar({ title }: TopbarProps) {
  const { data: user } = useCurrentUser()
  const initial = user?.name?.trim().charAt(0).toUpperCase() || 'L'

  return (
    <header className="sticky top-0 z-20 border-b border-border bg-background/95 px-4 py-4 backdrop-blur sm:px-6 lg:px-8">
      <div className="flex items-center justify-between gap-6">
        <div className="min-w-0">
          <p className="text-sm font-medium text-muted-foreground">LexPath</p>
          <h1 className="truncate text-2xl font-bold tracking-normal text-foreground">
            {title}
          </h1>
        </div>

        <div className="flex items-center gap-3">
          <div className="hidden h-10 w-64 items-center gap-3 rounded-xl border border-input bg-background px-3 text-sm text-muted-foreground md:flex">
            <Search className="size-4" aria-hidden="true" />
            <span>Search vocabulary</span>
          </div>

          <Button aria-label="Notifications" size="icon" variant="outline">
            <Bell aria-hidden="true" />
          </Button>

          <Button
            aria-label="User menu"
            className="rounded-full font-bold"
            size="icon"
            variant="secondary"
          >
            {initial}
          </Button>
        </div>
      </div>
    </header>
  )
}
