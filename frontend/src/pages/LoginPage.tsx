import { Navigate } from 'react-router-dom'

import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { LoginForm } from '@/features/auth/LoginForm'
import { useAuth } from '@/contexts/auth/useAuth'

export function LoginPage() {
  const { isAuthenticated } = useAuth()

  if (isAuthenticated) {
    return <Navigate replace to="/dashboard" />
  }

  return (
    <Card className="mx-auto w-full max-w-md p-6 shadow-xl shadow-slate-200/70 md:p-8">
      <CardHeader className="space-y-3 p-0">
        <CardTitle className="text-2xl">Welcome back</CardTitle>
        <CardDescription>
          Sign in to continue your IELTS vocabulary roadmap.
        </CardDescription>
      </CardHeader>
      <CardContent className="mt-12 p-0">
        <LoginForm />
      </CardContent>
    </Card>
  )
}
