import { BookOpen, Search, SlidersHorizontal } from 'lucide-react'
import { useMemo, useState } from 'react'
import { Link } from 'react-router-dom'

import { Badge } from '@/components/ui/badge'
import { Input } from '@/components/ui/input'
import {
  mockVocabulary,
  type MockVocabulary,
  type VocabularyDifficulty,
  type VocabularyStatus,
} from '@/features/vocabulary/mockVocabulary'
import { cn } from '@/lib/utils'

const difficulties = ['All', 'Beginner', 'Intermediate', 'Advanced'] as const
const bands = ['All', '5.5', '6.5', '7'] as const
const statuses = ['All', 'New', 'Learning', 'Review', 'Mastered'] as const

type DifficultyFilter = (typeof difficulties)[number]
type BandFilter = (typeof bands)[number]
type StatusFilter = (typeof statuses)[number]

export function VocabularyPage() {
  const [query, setQuery] = useState('')
  const [difficulty, setDifficulty] = useState<DifficultyFilter>('All')
  const [band, setBand] = useState<BandFilter>('All')
  const [status, setStatus] = useState<StatusFilter>('All')

  const filteredVocabulary = useMemo(() => {
    const normalizedQuery = query.trim().toLowerCase()

    return mockVocabulary.filter((item) => {
      const matchesQuery =
        normalizedQuery.length === 0 ||
        item.word.toLowerCase().includes(normalizedQuery) ||
        item.topic.toLowerCase().includes(normalizedQuery) ||
        item.shortDefinition.toLowerCase().includes(normalizedQuery)
      const matchesDifficulty =
        difficulty === 'All' || item.difficulty === difficulty
      const matchesBand = band === 'All' || item.bandScore.toFixed(1) === band
      const matchesStatus = status === 'All' || item.status === status

      return matchesQuery && matchesDifficulty && matchesBand && matchesStatus
    })
  }, [band, difficulty, query, status])

  return (
    <div className="min-h-screen bg-[#f8f8ff] px-4 py-8 text-[#10111f] md:px-8">
      <main className="mx-auto max-w-[1280px]">
        <section className="rounded-2xl border border-[#e6e6f3] bg-white p-6 shadow-sm md:p-8">
          <div className="flex flex-col gap-6 lg:flex-row lg:items-end lg:justify-between">
            <div>
              <p className="text-sm font-bold uppercase tracking-normal text-primary">
                Vocabulary
              </p>
              <h1 className="mt-2 text-3xl font-bold tracking-normal text-[#10111f] md:text-4xl">
                Vocabulary Library
              </h1>
              <p className="mt-3 max-w-2xl text-lg font-medium text-[#676982]">
                Browse IELTS words, review your progress, and open full word
                detail pages for dictionary-style study.
              </p>
            </div>

            <div className="flex items-center gap-3 rounded-2xl bg-[#f0efff] px-4 py-3 text-primary">
              <BookOpen className="size-5" aria-hidden="true" />
              <span className="text-base font-bold">
                {mockVocabulary.length} mock words
              </span>
            </div>
          </div>
        </section>

        <section className="mt-8 rounded-2xl border border-[#e6e6f3] bg-white p-5 shadow-sm">
          <div className="grid gap-4 lg:grid-cols-[1fr_auto_auto_auto]">
            <label className="relative block">
              <span className="sr-only">Search vocabulary</span>
              <Search
                className="pointer-events-none absolute left-4 top-1/2 size-5 -translate-y-1/2 text-[#85889c]"
                aria-hidden="true"
              />
              <Input
                className="h-12 rounded-full pl-12 text-base"
                onChange={(event) => setQuery(event.target.value)}
                placeholder="Search word, topic, or meaning"
                value={query}
              />
            </label>

            <FilterSelect
              label="Difficulty"
              onChange={(value) => setDifficulty(value as DifficultyFilter)}
              options={difficulties}
              value={difficulty}
            />
            <FilterSelect
              label="Band"
              onChange={(value) => setBand(value as BandFilter)}
              options={bands}
              value={band}
            />
            <FilterSelect
              label="Status"
              onChange={(value) => setStatus(value as StatusFilter)}
              options={statuses}
              value={status}
            />
          </div>
        </section>

        {filteredVocabulary.length > 0 ? (
          <section className="mt-8 grid gap-5 lg:grid-cols-2">
            {filteredVocabulary.map((item) => (
              <VocabularyListCard item={item} key={item.id} />
            ))}
          </section>
        ) : (
          <section className="mt-8 rounded-2xl border border-[#e6e6f3] bg-white p-8 text-center shadow-sm">
            <SlidersHorizontal
              className="mx-auto size-9 text-[#85889c]"
              aria-hidden="true"
            />
            <h2 className="mt-4 text-2xl font-bold tracking-normal">
              No matching words
            </h2>
            <p className="mx-auto mt-3 max-w-md text-base font-medium text-[#676982]">
              Try a different search term or clear one of the mock filters.
            </p>
          </section>
        )}
      </main>
    </div>
  )
}

type FilterSelectProps = {
  label: string
  options: readonly string[]
  value: string
  onChange: (value: string) => void
}

function FilterSelect({ label, options, value, onChange }: FilterSelectProps) {
  return (
    <label className="grid gap-2">
      <span className="text-sm font-bold text-[#676982]">{label}</span>
      <select
        className="h-12 min-w-36 rounded-full border border-[#e6e6f3] bg-white px-4 text-base font-semibold text-[#10111f] shadow-sm outline-none transition-colors focus:border-primary"
        onChange={(event) => onChange(event.target.value)}
        value={value}
      >
        {options.map((option) => (
          <option key={option} value={option}>
            {option}
          </option>
        ))}
      </select>
    </label>
  )
}

type VocabularyListCardProps = {
  item: MockVocabulary
}

function VocabularyListCard({ item }: VocabularyListCardProps) {
  return (
    <Link aria-label={`Open ${item.word}`} to={`/vocabulary/${item.id}`}>
      <article className="h-full rounded-2xl border border-[#e6e6f3] bg-white p-6 shadow-sm transition-colors hover:border-[#c6c4ff]">
        <div className="flex flex-col gap-5 sm:flex-row sm:items-start sm:justify-between">
          <div className="min-w-0">
            <div className="flex flex-wrap items-center gap-3">
              <h2 className="text-2xl font-bold tracking-normal text-[#10111f]">
                {item.word}
              </h2>
              <DifficultyBadge difficulty={item.difficulty} />
            </div>
            <p className="mt-2 font-mono text-base font-semibold text-[#85889c]">
              {item.ipa} · {item.partOfSpeech}
            </p>
            <p className="mt-3 line-clamp-2 text-lg font-medium text-[#676982]">
              {item.shortDefinition}
            </p>
          </div>

          <StatusBadge status={item.status} />
        </div>

        <div className="mt-6 flex flex-wrap items-center gap-3 text-base font-bold">
          <Badge className="bg-[#f0efff] text-primary">{item.band}</Badge>
          <Badge className="bg-[#eff6ff] text-[#2563eb]">{item.topic}</Badge>
          <span className="text-[#676982]">{item.frequency} frequency</span>
        </div>

        <div className="mt-5 h-2 overflow-hidden rounded-full bg-[#ebeaff]">
          <div
            className="h-full rounded-full bg-primary"
            style={{ width: `${item.masteryScore}%` }}
          />
        </div>
        <p className="mt-2 text-base font-bold text-[#676982]">
          {item.masteryScore}% mastery
        </p>
      </article>
    </Link>
  )
}

function DifficultyBadge({ difficulty }: { difficulty: VocabularyDifficulty }) {
  const styles = {
    Beginner: 'bg-[#dcfce7] text-[#138a53]',
    Intermediate: 'bg-[#fff3c4] text-[#c46700]',
    Advanced: 'bg-[#fde5e7] text-[#d22f38]',
  }

  return (
    <span className={cn('rounded-full px-3 py-1 text-base font-bold', styles[difficulty])}>
      {difficulty}
    </span>
  )
}

function StatusBadge({ status }: { status: VocabularyStatus }) {
  const styles = {
    New: 'bg-[#f3f4f6] text-[#4b5563]',
    Learning: 'bg-[#eff6ff] text-[#2563eb]',
    Review: 'bg-[#fff7ed] text-[#c46600]',
    Mastered: 'bg-[#e8fff3] text-[#138a53]',
  }

  return (
    <span className={cn('rounded-full px-3 py-1 text-base font-bold', styles[status])}>
      {status}
    </span>
  )
}
