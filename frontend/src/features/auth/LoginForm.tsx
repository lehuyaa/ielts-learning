import { zodResolver } from '@hookform/resolvers/zod'
import { LogIn } from 'lucide-react'
import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { Link, useLocation, useNavigate } from 'react-router-dom'

import { Button } from '@/components/ui/button'
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { APIError } from '@/api/api'

import { useAuth } from '@/contexts/auth/useAuth'
import {
  loginSchema,
  normalizeLoginValues,
  type LoginFormValues,
} from './validation/authSchemas'

type LocationState = {
  from?: {
    pathname?: string
  }
}

export function LoginForm() {
  const { login } = useAuth()
  const navigate = useNavigate()
  const location = useLocation()
  const [error, setError] = useState('')

  const form = useForm<LoginFormValues>({
    resolver: zodResolver(loginSchema),
    defaultValues: {
      email: '',
      password: '',
    },
  })

  async function onSubmit(values: LoginFormValues) {
    setError('')
    try {
      await login(normalizeLoginValues(values))
      const state = location.state as LocationState | null
      navigate(state?.from?.pathname ?? '/dashboard', { replace: true })
    } catch (err) {
      if (err instanceof APIError && err.fields) {
        applyFieldErrors(err.fields)
      }
      setError(err instanceof Error ? err.message : 'Login failed')
    }
  }

  function applyFieldErrors(fields: Record<string, string>) {
    for (const [field, message] of Object.entries(fields)) {
      if (field === 'email' || field === 'password') {
        form.setError(field, { message })
      }
    }
  }

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)}>
        {error ? (
          <div className="mb-6 rounded-xl border border-destructive/30 bg-destructive/10 px-4 py-3 text-sm font-medium text-destructive">
            {error}
          </div>
        ) : null}

        <div className="space-y-6">
          <FormField
            control={form.control}
            name="email"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Email</FormLabel>
                <FormControl>
                  <Input
                    autoComplete="email"
                    placeholder="demo@example.com"
                    type="email"
                    {...field}
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={form.control}
            name="password"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Password</FormLabel>
                <FormControl>
                  <Input
                    autoComplete="current-password"
                    placeholder="password"
                    type="password"
                    {...field}
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
        </div>

        <Button className="mt-8 w-full" disabled={form.formState.isSubmitting} type="submit">
          <LogIn aria-hidden="true" />
          {form.formState.isSubmitting ? 'Signing in' : 'Sign in'}
        </Button>

        <p className="mt-6 text-center text-sm text-muted-foreground">
          New here?{' '}
          <Link className="font-semibold text-primary hover:underline" to="/register">
            Create an account
          </Link>
        </p>
      </form>
    </Form>
  )
}
