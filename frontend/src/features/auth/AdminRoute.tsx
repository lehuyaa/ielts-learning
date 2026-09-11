import { Navigate, Outlet } from 'react-router-dom'

import { useCurrentUser } from '@/features/auth/hooks/useCurrentUser'

export function AdminRoute() {
  const currentUserQuery = useCurrentUser()

  if (currentUserQuery.isLoading) {
    return (
      <main className="grid min-h-screen place-items-center bg-background px-4 text-foreground">
        <div className="rounded-xl border border-border bg-card px-5 py-4 text-sm font-medium text-muted-foreground shadow-sm">
          Loading
        </div>
      </main>
    )
  }

  if (currentUserQuery.data?.role !== 'ADMIN') {
    return <Navigate replace to="/dashboard" />
  }

  return <Outlet />
}
