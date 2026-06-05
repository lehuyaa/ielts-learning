import { BookOpen } from 'lucide-react'
import { Link, Navigate } from 'react-router-dom'

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
    <main className="grid min-h-screen bg-background px-4 py-8 text-foreground">
      <div className="mx-auto grid w-full max-w-xl content-center gap-8">
        <Link className="mx-auto flex items-center gap-3 font-semibold" to="/">
          <span className="grid size-10 place-items-center rounded-2xl bg-primary text-primary-foreground">
            <BookOpen className="size-5" aria-hidden="true" />
          </span>
          <span>LexPath</span>
        </Link>

        <Card className="p-6 shadow-xl shadow-slate-200/70 md:p-8">
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
      </div>
    </main>
  )
}
