import { z } from 'zod'

const usernameRegex = /^[a-z0-9_]+$/
const localeRegex = /^[a-z]{2}(-[a-z]{2})?$/
const targetBandValues = [5, 5.5, 6, 6.5, 7, 7.5, 8, 8.5] as const

function isValidTimeZone(value: string) {
  try {
    Intl.DateTimeFormat(undefined, { timeZone: value })
    return true
  } catch {
    return false
  }
}

export const profileUpdateSchema = z.object({
  name: z
    .string()
    .trim()
    .min(2, 'Name must be at least 2 characters.')
    .max(80, 'Name must be at most 80 characters.'),
  username: z
    .string()
    .trim()
    .toLowerCase()
    .min(3, 'Username must be at least 3 characters.')
    .max(30, 'Username must be at most 30 characters.')
    .regex(
      usernameRegex,
      'Username can only use lowercase letters, numbers, and underscore.',
    ),
  targetBand: z
    .number()
    .refine(
      (value) =>
        targetBandValues.includes(
          value as (typeof targetBandValues)[number],
        ),
      'Choose a valid IELTS target band.',
    ),
  timezone: z
    .string()
    .trim()
    .min(1, 'Timezone is required.')
    .max(80, 'Timezone must be at most 80 characters.')
    .refine(isValidTimeZone, 'Timezone must be a valid IANA timezone.'),
  locale: z
    .string()
    .trim()
    .toLowerCase()
    .min(2, 'Locale must be at least 2 characters.')
    .max(20, 'Locale must be at most 20 characters.')
    .regex(
      localeRegex,
      'Locale must use language or language-region format.',
    ),
})

export type ProfileUpdateFormValues = z.infer<typeof profileUpdateSchema>

export function normalizeProfileUpdateValues(values: ProfileUpdateFormValues) {
  return {
    name: values.name.trim(),
    username: values.username.trim().toLowerCase(),
    targetBand: values.targetBand,
    timezone: values.timezone.trim(),
    locale: values.locale.trim().toLowerCase(),
  }
}
