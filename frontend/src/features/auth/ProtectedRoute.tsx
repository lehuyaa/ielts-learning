import { useEffect } from 'react'
import { Navigate, Outlet, useLocation } from 'react-router-dom'

import { useAuth } from '@/contexts/auth/useAuth'
import { useCurrentUser } from '@/features/auth/hooks/useCurrentUser'

export function ProtectedRoute() {
  const { isAuthenticated, logout } = useAuth()
  const currentUserQuery = useCurrentUser()
  const location = useLocation()

  useEffect(() => {
    if (isAuthenticated && currentUserQuery.isError) {
      logout()
    }
  }, [currentUserQuery.isError, isAuthenticated, logout])

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

  if (currentUserQuery.isError) {
    return <Navigate replace state={{ from: location }} to="/login" />
  }

  return <Outlet />
}
