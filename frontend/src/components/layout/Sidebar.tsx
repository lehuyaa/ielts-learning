import {
  BookOpenText,
  LayoutDashboard,
  LogOut,
  Map,
  RotateCcw,
  UserCircle,
} from 'lucide-react'
import { NavLink, useNavigate } from 'react-router-dom'

import { BrandMark } from '@/components/layout/BrandMark'
import { Button } from '@/components/ui/button'
import { useAuth } from '@/contexts/auth/useAuth'
import { useCurrentUser } from '@/features/auth/hooks/useCurrentUser'
import { cn } from '@/lib/utils'

const navigationItems = [
  {
    label: 'Dashboard',
    to: '/dashboard',
    icon: LayoutDashboard,
  },
  {
    label: 'Roadmap',
    to: '/roadmap',
    icon: Map,
  },
  {
    label: 'Reviews',
    to: '/reviews',
    icon: RotateCcw,
  },
  {
    label: 'Vocabulary',
    to: '/vocabulary',
    icon: BookOpenText,
  },
  {
    label: 'Profile',
    to: '/profile',
    icon: UserCircle,
  },
]

function getNavLinkClass(isActive: boolean) {
  return cn(
    'flex items-center gap-3 rounded-xl px-3 py-2 text-sm font-medium text-muted-foreground transition-colors',
    'hover:bg-accent hover:text-accent-foreground',
    isActive &&
      'bg-primary text-primary-foreground hover:bg-primary hover:text-primary-foreground',
  )
}

export function Sidebar() {
  const { logout } = useAuth()
  const { data: user } = useCurrentUser()
  const navigate = useNavigate()

  function handleLogout() {
    logout()
    navigate('/login', { replace: true })
  }

  return (
    <aside className="hidden min-h-screen border-r border-border bg-card px-4 py-6 lg:flex lg:flex-col">
      <BrandMark className="px-2" />

      <nav className="mt-12 grid gap-2">
        {navigationItems.map((item) => (
          <NavLink
            className={({ isActive }) => getNavLinkClass(isActive)}
            key={item.to}
            to={item.to}
          >
            <item.icon className="size-4" aria-hidden="true" />
            {item.label}
          </NavLink>
        ))}
      </nav>

      <div className="mt-auto grid gap-4 rounded-xl border border-border bg-background p-4">
        <div className="min-w-0">
          <p className="truncate text-sm font-semibold">{user?.name ?? 'Learner'}</p>
          <p className="truncate text-xs text-muted-foreground">{user?.email}</p>
        </div>
        <Button onClick={handleLogout} variant="outline">
          <LogOut aria-hidden="true" />
          Logout
        </Button>
      </div>
    </aside>
  )
}

export function DashboardMobileNav() {
  return (
    <nav className="border-b border-border bg-card px-4 py-3 lg:hidden">
      <div className="flex gap-2 overflow-x-auto">
        {navigationItems.map((item) => (
          <NavLink
            className={({ isActive }) =>
              cn(
                'flex shrink-0 items-center gap-2 rounded-xl px-3 py-2 text-sm font-medium text-muted-foreground',
                'hover:bg-accent hover:text-accent-foreground',
                isActive &&
                  'bg-primary text-primary-foreground hover:bg-primary hover:text-primary-foreground',
              )
            }
            key={item.to}
            to={item.to}
          >
            <item.icon className="size-4" aria-hidden="true" />
            {item.label}
          </NavLink>
        ))}
      </div>
    </nav>
  )
}
