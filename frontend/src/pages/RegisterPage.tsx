import { Navigate } from 'react-router-dom'

import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { RegisterForm } from '@/features/auth/RegisterForm'
import { useAuth } from '@/contexts/auth/useAuth'

export function RegisterPage() {
  const { isAuthenticated } = useAuth()

  if (isAuthenticated) {
    return <Navigate replace to="/dashboard" />
  }

  return (
    <Card className="w-full p-6 shadow-xl shadow-slate-200/70 md:p-8">
      <CardHeader className="space-y-3 p-0">
        <CardTitle className="text-2xl">Create your account</CardTitle>
        <CardDescription>
          Set your target band and start learning with structured lessons.
        </CardDescription>
      </CardHeader>
      <CardContent className="mt-12 p-0">
        <RegisterForm />
      </CardContent>
    </Card>
  )
}
