import { zodResolver } from '@hookform/resolvers/zod'
import { UserPlus } from 'lucide-react'
import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { Link, useNavigate } from 'react-router-dom'

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
  normalizeRegisterValues,
  registerSchema,
  type RegisterFormValues,
} from './validation/authSchemas'

export function RegisterForm() {
  const { register } = useAuth()
  const navigate = useNavigate()
  const [error, setError] = useState('')

  const form = useForm<RegisterFormValues>({
    resolver: zodResolver(registerSchema),
    defaultValues: {
      name: '',
      email: '',
      username: '',
      password: '',
      targetBand: 7,
    },
  })

  async function onSubmit(values: RegisterFormValues) {
    setError('')
    try {
      await register(normalizeRegisterValues(values))
      navigate('/dashboard', { replace: true })
    } catch (err) {
      if (err instanceof APIError && err.fields) {
        applyFieldErrors(err.fields)
      }
      setError(err instanceof Error ? err.message : 'Registration failed')
    }
  }

  function applyFieldErrors(fields: Record<string, string>) {
    for (const [field, message] of Object.entries(fields)) {
      if (
        field === 'name' ||
        field === 'username' ||
        field === 'email' ||
        field === 'targetBand' ||
        field === 'password'
      ) {
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
            name="name"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Name</FormLabel>
                <FormControl>
                  <Input
                    autoComplete="name"
                    placeholder="Alex Johnson"
                    {...field}
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={form.control}
            name="email"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Email</FormLabel>
                <FormControl>
                  <Input
                    autoComplete="email"
                    placeholder="alex@example.com"
                    type="email"
                    {...field}
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />

          <div className="grid gap-6 sm:grid-cols-2">
            <FormField
              control={form.control}
              name="username"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Username</FormLabel>
                  <FormControl>
                    <Input
                      autoComplete="username"
                      placeholder="alexj"
                      {...field}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="targetBand"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Target band</FormLabel>
                  <FormControl>
                    <Input
                      min="0"
                      max="9"
                      step="0.5"
                      type="number"
                      value={field.value}
                      onBlur={field.onBlur}
                      onChange={(event) =>
                        field.onChange(Number(event.target.value))
                      }
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
          </div>

          <FormField
            control={form.control}
            name="password"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Password</FormLabel>
                <FormControl>
                  <Input
                    autoComplete="new-password"
                    placeholder="At least 8 characters"
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
          <UserPlus aria-hidden="true" />
          {form.formState.isSubmitting ? 'Creating account' : 'Create account'}
        </Button>

        <p className="mt-6 text-center text-sm text-muted-foreground">
          Already have an account?{' '}
          <Link className="font-semibold text-primary hover:underline" to="/login">
            Sign in
          </Link>
        </p>
      </form>
    </Form>
  )
}
