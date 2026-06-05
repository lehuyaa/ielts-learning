import { RouterProvider } from 'react-router-dom'

import { AuthProvider } from '@/contexts/auth/AuthContext.tsx'

import { router } from './router'

export function AppProviders() {
  return (
    <AuthProvider>
      <RouterProvider router={router} />
    </AuthProvider>
  )
}
