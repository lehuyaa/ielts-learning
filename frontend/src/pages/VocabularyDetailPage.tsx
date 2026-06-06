import {
  AlertCircle,
  ArrowLeft,
  ArrowRight,
  Bookmark,
  BookOpen,
  CheckCircle2,
  Share2,
  Star,
  Volume2,
} from 'lucide-react'
import type React from 'react'
import { useState } from 'react'
import { Link, useNavigate, useParams } from 'react-router-dom'

import { APIError } from '@/api/api'
import { Button } from '@/components/ui/button'
import { useVocabularyDetail } from '@/features/vocabulary/hooks/useVocabularyDetail'
import { mapVocabularyDetail } from '@/features/vocabulary/mapVocabulary'
import { cn } from '@/lib/utils'
import type { VocabularyDetailViewModel } from '@/types/vocabulary'

const tabs = [
  'Meaning',
  'Examples',
  'Synonyms',
  'Collocations',
  'IELTS Usage',
] as const

type VocabularyTab = (typeof tabs)[number]

export function VocabularyDetailPage() {
  const { vocabularyId } = useParams()
  const navigate = useNavigate()
  const vocabularyQuery = useVocabularyDetail(vocabularyId)
  const vocabulary = vocabularyQuery.data
    ? mapVocabularyDetail(vocabularyQuery.data)
    : null
  const [activeTab, setActiveTab] = useState<VocabularyTab>('Meaning')
  const errorMessage = getVocabularyDetailErrorMessage(vocabularyQuery.error)

  if (vocabularyQuery.isLoading) {
    return (
      <VocabularyDetailShell navigateBack={() => navigate(-1)}>
        <VocabularyDetailState
          description="Loading word metadata, examples, and your progress."
          title="Loading word"
        />
      </VocabularyDetailShell>
    )
  }

  if (errorMessage) {
    return (
      <VocabularyDetailShell navigateBack={() => navigate(-1)}>
        <VocabularyDetailState
          description={errorMessage}
          title="Word unavailable"
          tone="error"
        />
      </VocabularyDetailShell>
    )
  }

  if (!vocabulary) {
    return (
      <VocabularyDetailShell navigateBack={() => navigate(-1)}>
        <VocabularyNotFound />
      </VocabularyDetailShell>
    )
  }

  return (
    <VocabularyDetailShell navigateBack={() => navigate(-1)}>
      <WordHero vocabulary={vocabulary} />

      <section className="mt-8 grid gap-8 lg:grid-cols-[1fr_360px]">
        <div className="min-w-0">
          <VocabularyTabs activeTab={activeTab} onChange={setActiveTab} />
          <TabContent activeTab={activeTab} vocabulary={vocabulary} />
        </div>

        <VocabularySidebar vocabulary={vocabulary} />
      </section>
    </VocabularyDetailShell>
  )
}

type VocabularyDetailShellProps = {
  children: React.ReactNode
  navigateBack: () => void
}

function VocabularyDetailShell({
  children,
  navigateBack,
}: VocabularyDetailShellProps) {
  return (
    <div className="min-h-screen bg-[#f8f8ff] text-[#10111f]">
      <header className="border-b border-[#e6e6f3] bg-white">
        <div className="mx-auto flex min-h-[88px] max-w-[1280px] items-center justify-between gap-4 px-4 md:px-8">
          <div className="flex items-center gap-6">
            <button
              aria-label="Go back"
              className="grid size-10 place-items-center rounded-full text-[#6d7088] transition-colors hover:bg-[#f0f1fb]"
              onClick={navigateBack}
              type="button"
            >
              <ArrowLeft className="size-5" aria-hidden="true" />
            </button>
            <h1 className="text-xl font-bold tracking-normal">Word Detail</h1>
          </div>

          <div className="flex items-center gap-3 text-[#6d7088]">
            <Button aria-label="Save word" size="icon" variant="ghost">
              <Bookmark className="size-5" aria-hidden="true" />
            </Button>
            <Button aria-label="Share word" size="icon" variant="ghost">
              <Share2 className="size-5" aria-hidden="true" />
            </Button>
          </div>
        </div>
      </header>

      <main className="mx-auto max-w-[1280px] px-4 py-8 md:px-8">
        {children}
      </main>
    </div>
  )
}

function WordHero({ vocabulary }: { vocabulary: VocabularyDetailViewModel }) {
  return (
    <section className="rounded-[28px] bg-gradient-to-br from-[#6258f6] to-[#8318e8] p-6 text-white shadow-sm md:p-8">
      <div className="flex flex-col gap-8 md:flex-row md:items-start md:justify-between">
        <div>
          <div className="flex flex-wrap gap-3">
            <span className="rounded-full bg-white/20 px-4 py-2 text-base font-bold">
              {vocabulary.topic}
            </span>
            <span className="rounded-full bg-white/20 px-4 py-2 text-base font-bold">
              {vocabulary.band}
            </span>
          </div>
          <h2 className="mt-5 text-5xl font-bold tracking-normal md:text-6xl">
            {vocabulary.word}
          </h2>
          <p className="mt-4 font-mono text-xl font-semibold text-white/75">
            {vocabulary.ipa}
            <span className="ml-4 font-sans text-white/80">
              {vocabulary.partOfSpeech}
            </span>
          </p>
          <p className="mt-6 flex items-center gap-3 text-lg font-medium text-white/80">
            <CheckCircle2 className="size-5" aria-hidden="true" />
            {vocabulary.frequency} frequency in IELTS Writing Task 2
          </p>
        </div>

        <div className="grid justify-items-start gap-4 md:justify-items-end">
          <Button className="rounded-full bg-white/20 px-6 text-base text-white hover:bg-white/25">
            <Volume2 aria-hidden="true" />
            Listen
          </Button>
          <div
            aria-label={`${vocabulary.frequencyScore}% frequency rating`}
            className="flex gap-1 text-[#ffd235]"
          >
            {Array.from({ length: 5 }).map((_, index) => (
              <Star
                className={cn(
                  'size-5',
                  index < Math.round(vocabulary.frequencyScore / 20)
                    ? 'fill-current'
                    : 'text-white/35',
                )}
                aria-hidden="true"
                key={index}
              />
            ))}
          </div>
        </div>
      </div>
    </section>
  )
}

type VocabularyTabsProps = {
  activeTab: VocabularyTab
  onChange: (tab: VocabularyTab) => void
}

function VocabularyTabs({ activeTab, onChange }: VocabularyTabsProps) {
  return (
    <div className="flex overflow-x-auto rounded-[24px] bg-[#e9e9f3] p-1">
      {tabs.map((tab) => (
        <button
          className={cn(
            'min-w-fit rounded-[20px] px-6 py-3 text-lg font-bold transition-colors',
            activeTab === tab
              ? 'bg-white text-primary shadow-sm'
              : 'text-[#676982] hover:text-[#10111f]',
          )}
          key={tab}
          onClick={() => onChange(tab)}
          type="button"
        >
          {tab}
        </button>
      ))}
    </div>
  )
}

type TabContentProps = {
  activeTab: VocabularyTab
  vocabulary: VocabularyDetailViewModel
}

function TabContent({ activeTab, vocabulary }: TabContentProps) {
  if (activeTab === 'Examples') {
    return (
      <div className="mt-6 grid gap-5">
        {vocabulary.examples.map((example) => (
          <ContentCard eyebrow="IELTS Example" key={example.sentence}>
            <blockquote className="text-xl font-medium text-[#30348e]">
              "{example.sentence}"
            </blockquote>
            <p className="mt-3 text-base font-medium text-[#676982]">
              {example.note}
            </p>
          </ContentCard>
        ))}
      </div>
    )
  }

  if (activeTab === 'Synonyms') {
    return (
      <div className="mt-6 grid gap-5">
        <TokenCard eyebrow="Synonyms" tokens={vocabulary.synonyms} />
        <TokenCard eyebrow="Antonyms" tokens={vocabulary.antonyms} />
      </div>
    )
  }

  if (activeTab === 'Collocations') {
    return (
      <div className="mt-6">
        <TokenCard eyebrow="Collocations" tokens={vocabulary.collocations} />
      </div>
    )
  }

  if (activeTab === 'IELTS Usage') {
    return (
      <div className="mt-6">
        <ContentCard eyebrow="IELTS Usage">
          <p className="text-xl font-medium leading-relaxed">
            {vocabulary.ieltsUsage}
          </p>
        </ContentCard>
      </div>
    )
  }

  return (
    <div className="mt-6 grid gap-5">
      <ContentCard eyebrow="Primary Meaning">
        <p className="text-xl font-medium leading-relaxed">
          {vocabulary.primaryMeaning}
        </p>
      </ContentCard>

      <ContentCard eyebrow="Secondary Meaning">
        <p className="text-xl font-medium leading-relaxed">
          {vocabulary.secondaryMeaning}
        </p>
      </ContentCard>

      <ContentCard eyebrow="Related Forms">
        <div className="flex flex-wrap gap-3">
          {vocabulary.relatedForms.length > 0 ? (
            vocabulary.relatedForms.map((form) => (
              <span
                className="rounded-full border border-[#dde3ff] bg-[#f0efff] px-4 py-3 text-base font-bold text-primary"
                key={`${form.word}-${form.partOfSpeech}`}
              >
                {form.word}
                <span className="ml-2 font-medium text-[#676982]">
                  {form.partOfSpeech}
                </span>
              </span>
            ))
          ) : (
            <p className="text-lg font-medium text-[#676982]">
              No related forms have been added yet.
            </p>
          )}
        </div>
      </ContentCard>
    </div>
  )
}

type ContentCardProps = {
  eyebrow: string
  children: React.ReactNode
}

function ContentCard({ eyebrow, children }: ContentCardProps) {
  return (
    <section className="rounded-2xl border border-[#e6e6f3] bg-white p-6 shadow-sm">
      <p className="text-base font-bold uppercase tracking-normal text-[#676982]">
        {eyebrow}
      </p>
      <div className="mt-4 text-[#10111f]">{children}</div>
    </section>
  )
}

function TokenCard({ eyebrow, tokens }: { eyebrow: string; tokens: string[] }) {
  if (tokens.length === 0) {
    return (
      <ContentCard eyebrow={eyebrow}>
        <p className="text-lg font-medium text-[#676982]">
          No {eyebrow.toLowerCase()} have been added yet.
        </p>
      </ContentCard>
    )
  }

  return (
    <ContentCard eyebrow={eyebrow}>
      <div className="flex flex-wrap gap-3">
        {tokens.map((token) => (
          <span
            className="rounded-full bg-[#f0efff] px-4 py-3 text-base font-bold text-primary"
            key={token}
          >
            {token}
          </span>
        ))}
      </div>
    </ContentCard>
  )
}

function VocabularySidebar({
  vocabulary,
}: {
  vocabulary: VocabularyDetailViewModel
}) {
  return (
    <aside className="grid gap-6 lg:sticky lg:top-8">
      <section className="rounded-2xl border border-[#e6e6f3] bg-white p-6 shadow-sm">
        <h2 className="text-xl font-bold tracking-normal">Quick Actions</h2>
        <div className="mt-6 grid gap-3">
          <Button className="h-14 justify-between rounded-2xl bg-[#eff1ff] px-5 text-lg font-bold text-primary hover:bg-[#e6e8ff]">
            <span className="inline-flex items-center gap-3">
              <BookOpen aria-hidden="true" />
              Add to Flashcards
            </span>
            <ArrowRight aria-hidden="true" />
          </Button>
          <Button className="h-14 justify-between rounded-2xl bg-[#e8fff3] px-5 text-lg font-bold text-[#208b5a] hover:bg-[#dcfce7]">
            <span className="inline-flex items-center gap-3">
              <CheckCircle2 aria-hidden="true" />
              Practice in Quiz
            </span>
            <ArrowRight aria-hidden="true" />
          </Button>
        </div>
      </section>

      <section className="rounded-2xl border border-[#e6e6f3] bg-white p-6 shadow-sm">
        <h2 className="text-xl font-bold tracking-normal">Word Stats</h2>
        <div className="mt-5 grid gap-4">
          <StatRow label="IELTS Band" value={vocabulary.band} />
          <StatRow label="Frequency" value={vocabulary.frequency} />
          <StatRow label="Topic" value={vocabulary.topic} />
          <StatRow label="Your Status" value={vocabulary.status} />
        </div>
      </section>

      <section className="rounded-2xl bg-gradient-to-br from-[#45c486] to-[#2e9f83] p-6 text-white shadow-sm">
        <p className="text-base font-bold text-white/75">Mastery Score</p>
        <p className="mt-3 text-5xl font-bold tracking-normal">
          {vocabulary.masteryScore}%
        </p>
        <p className="mt-3 text-base font-medium text-white/80">
          Reviewed {vocabulary.reviewCount} times · Last: {vocabulary.lastReviewedAt}
        </p>
      </section>
    </aside>
  )
}

function StatRow({ label, value }: { label: string; value: string }) {
  return (
    <div className="flex items-center justify-between gap-4 text-lg">
      <span className="font-medium text-[#676982]">{label}</span>
      <span className="font-bold text-[#10111f]">{value}</span>
    </div>
  )
}

function VocabularyNotFound() {
  return (
    <section className="mx-auto w-full max-w-md rounded-2xl border border-[#e6e6f3] bg-white p-8 text-center shadow-sm">
      <h1 className="text-3xl font-bold tracking-normal">Word not found</h1>
      <p className="mt-3 text-base font-medium text-[#676982]">
        This vocabulary item does not exist yet.
      </p>
      <Button asChild className="mt-8 rounded-full">
        <Link to="/vocabulary">Back to vocabulary</Link>
      </Button>
    </section>
  )
}

type VocabularyDetailStateProps = {
  title: string
  description: string
  tone?: 'default' | 'error'
}

function VocabularyDetailState({
  title,
  description,
  tone = 'default',
}: VocabularyDetailStateProps) {
  return (
    <section className="rounded-2xl border border-[#e6e6f3] bg-white p-8 text-center shadow-sm">
      {tone === 'error' ? (
        <AlertCircle
          className="mx-auto size-9 text-destructive"
          aria-hidden="true"
        />
      ) : null}
      <h2 className="mt-4 text-2xl font-bold tracking-normal">{title}</h2>
      <p className="mx-auto mt-3 max-w-md text-base font-medium text-[#676982]">
        {description}
      </p>
      {tone === 'error' ? (
        <Button asChild className="mt-8 rounded-full">
          <Link to="/vocabulary">Back to vocabulary</Link>
        </Button>
      ) : null}
    </section>
  )
}

function getVocabularyDetailErrorMessage(error: Error | null) {
  if (!error) {
    return null
  }

  if (error instanceof APIError) {
    if (error.status === 404) {
      return 'This vocabulary item could not be found.'
    }

    return error.message
  }

  return 'Unable to load this word right now.'
}
