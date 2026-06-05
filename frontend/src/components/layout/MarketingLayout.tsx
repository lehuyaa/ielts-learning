import { Link, Outlet } from 'react-router-dom'

import { BrandMark } from '@/components/layout/BrandMark'
import { Button } from '@/components/ui/button'

export function MarketingLayout() {
  return (
    <main className="min-h-screen bg-background text-foreground">
      <section className="mx-auto flex min-h-screen w-full max-w-6xl flex-col px-4 py-8 sm:px-6 lg:px-8">
        <nav className="flex items-center justify-between gap-6">
          <BrandMark />

          <Button asChild variant="outline">
            <Link to="/register">Start Learning</Link>
          </Button>
        </nav>

        <Outlet />
      </section>
    </main>
  )
}
