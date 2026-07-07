import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { RouterProvider } from 'react-router-dom'

import { APIError } from '@/api/api'
import { AuthProvider } from '@/contexts/auth/AuthContext.tsx'
import { ToastProvider } from '@/contexts/toast/ToastProvider'

import { router } from './router'

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: (failureCount, error) => {
        if (error instanceof APIError) {
          if (error.status === 401 || error.status === 403 || error.status === 404) {
            return false
          }
        }

        return failureCount < 1
      },
    },
    mutations: {
      retry: false,
    },
  },
})

export function AppProviders() {
  return (
    <QueryClientProvider client={queryClient}>
      <ToastProvider>
        <AuthProvider>
          <RouterProvider router={router} />
        </AuthProvider>
      </ToastProvider>
    </QueryClientProvider>
  )
}
