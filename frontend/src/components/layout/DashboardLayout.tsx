import { Outlet, useLocation } from 'react-router-dom'

import { DashboardMobileNav, Sidebar } from '@/components/layout/Sidebar'
import { Topbar } from '@/components/layout/Topbar'

const pageTitles: Record<string, string> = {
  '/dashboard': 'Dashboard',
  '/roadmap': 'Roadmap',
  '/reviews': 'Reviews',
  '/vocabulary': 'Vocabulary',
  '/profile': 'Profile',
}

export function DashboardLayout() {
  const { pathname } = useLocation()
  const title = pageTitles[pathname] ?? 'Dashboard'
  const isRoadmapPage = pathname === '/roadmap'

  if (isRoadmapPage) {
    return <Outlet />
  }

  return (
    <div className="min-h-screen bg-background text-foreground lg:grid lg:grid-cols-[280px_1fr]">
      <Sidebar />

      <div className="min-w-0">
        <Topbar title={title} />
        <DashboardMobileNav />

        <main className="px-4 py-8 sm:px-6 lg:px-8">
          <div className="mx-auto max-w-6xl">
            <Outlet />
          </div>
        </main>
      </div>
    </div>
  )
}
