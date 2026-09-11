import { LogOut } from 'lucide-react'
import { Navigate, Outlet, useLocation } from 'react-router-dom'

import { APIError } from '@/api/api'
import { Button } from '@/components/ui/button'
import { useAuth } from '@/contexts/auth/useAuth'
import { useCurrentUser } from '@/features/auth/hooks/useCurrentUser'

export function ProtectedRoute() {
  const { isAuthenticated, logout } = useAuth()
  const currentUserQuery = useCurrentUser()
  const location = useLocation()

  if (!isAuthenticated) {
    return <Navigate replace state={{ from: location }} to="/login" />
  }

  if (currentUserQuery.isLoading) {
    return (
      <main className="grid min-h-screen place-items-center bg-background px-4 text-foreground">
        <div className="rounded-xl border border-border bg-card px-5 py-4 text-sm font-medium text-muted-foreground shadow-sm">
          Loading
        </div>
      </main>
    )
  }

  if (
    currentUserQuery.isError &&
    currentUserQuery.error instanceof APIError &&
    currentUserQuery.error.status === 401
  ) {
    return <Navigate replace state={{ from: location }} to="/login" />
  }

  if (currentUserQuery.isError) {
    return (
      <main className="grid min-h-screen place-items-center bg-background px-4 text-foreground">
        <div className="rounded-xl border border-border bg-card px-5 py-4 text-center text-sm shadow-sm">
          <p className="font-medium text-foreground">
            Unable to verify your session right now.
          </p>
          <p className="mt-1 text-muted-foreground">
            Please try again in a moment.
          </p>
          <Button
            className="mt-4 gap-1.5 rounded-full"
            onClick={logout}
            type="button"
          >
            <LogOut className="size-4" aria-hidden="true" />
            Log out
          </Button>
        </div>
      </main>
    )
  }

  return <Outlet />
}
