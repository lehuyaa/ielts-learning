import { useQuery } from '@tanstack/react-query'

import { getPhoneme, getPhonemes } from '@/api/phoneme'

export const phonemeQueryKey = ['phonemes'] as const

export function usePhonemes() {
  return useQuery({ queryKey: [...phonemeQueryKey, 'list'], queryFn: getPhonemes })
}

export function usePhoneme(symbol: string | undefined) {
  return useQuery({
    queryKey: [...phonemeQueryKey, 'detail', symbol],
    queryFn: () => getPhoneme(symbol as string),
    enabled: Boolean(symbol),
  })
}
