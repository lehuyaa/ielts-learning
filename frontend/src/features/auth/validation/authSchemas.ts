import { z } from 'zod'

const usernameRegex = /^[a-z0-9_]+$/
const letterRegex = /[A-Za-z]/
const numberRegex = /[0-9]/
const targetBandValues = [5, 5.5, 6, 6.5, 7, 7.5, 8, 8.5] as const

export const loginSchema = z.object({
  email: z
    .string()
    .trim()
    .toLowerCase()
    .min(1, 'Email is required.')
    .max(255, 'Email must be at most 255 characters.')
    .email('Please enter a valid email address.'),
  password: z
    .string()
    .min(8, 'Password must be at least 8 characters.')
    .max(72, 'Password must be at most 72 characters.'),
})

export const registerSchema = z.object({
  name: z
    .string()
    .trim()
    .min(2, 'Name must be at least 2 characters.')
    .max(80, 'Name must be at most 80 characters.'),
  username: z
    .string()
    .trim()
    .toLowerCase()
    .refine(
      (value) => value === '' || value.length >= 3,
      'Username must be at least 3 characters.',
    )
    .refine(
      (value) => value === '' || value.length <= 30,
      'Username must be at most 30 characters.',
    )
    .refine(
      (value) => value === '' || usernameRegex.test(value),
      'Username can only use lowercase letters, numbers, and underscore.',
    ),
  email: z
    .string()
    .trim()
    .toLowerCase()
    .min(1, 'Email is required.')
    .max(255, 'Email must be at most 255 characters.')
    .email('Please enter a valid email address.'),
  targetBand: z
    .number()
    .refine(
      (value) => targetBandValues.includes(value as (typeof targetBandValues)[number]),
      'Choose a valid IELTS target band.',
    ),
  password: z
    .string()
    .min(8, 'Password must be at least 8 characters.')
    .max(72, 'Password must be at most 72 characters.')
    .regex(letterRegex, 'Password must contain at least one letter.')
    .regex(numberRegex, 'Password must contain at least one number.'),
})

export type LoginFormValues = z.infer<typeof loginSchema>
export type RegisterFormValues = z.infer<typeof registerSchema>

export function normalizeLoginValues(values: LoginFormValues) {
  return {
    email: values.email.trim().toLowerCase(),
    password: values.password,
  }
}

export function normalizeRegisterValues(values: RegisterFormValues) {
  const username = values.username.trim().toLowerCase()

  return {
    name: values.name.trim(),
    email: values.email.trim().toLowerCase(),
    username: username || undefined,
    targetBand: values.targetBand,
    password: values.password,
  }
}
