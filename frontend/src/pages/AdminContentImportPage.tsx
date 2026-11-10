import {
  CheckCircle2,
  ChevronLeft,
  ChevronRight,
  TriangleAlert,
  Upload,
} from 'lucide-react'
import { useMemo, useRef, useState } from 'react'
import { useNavigate } from 'react-router-dom'

import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import {
  useContentImportCommit,
  useContentImportPreview,
} from '@/features/admin/hooks/useContentImport'
import { cn } from '@/lib/utils'
import type {
  ContentImportItem,
  ContentImportResult,
} from '@/types/contentImport'

/** Lessons with their activities and exercises, grouped for display. */
type LessonGroup = {
  lesson: ContentImportItem
  activities: { activity: ContentImportItem; exercises: ContentImportItem[] }[]
}

export function AdminContentImportPage() {
  const navigate = useNavigate()
  const fileInputRef = useRef<HTMLInputElement>(null)
  const [selectedFile, setSelectedFile] = useState<File | null>(null)
  const [report, setReport] = useState<ContentImportResult | null>(null)

  const previewMutation = useContentImportPreview()
  const commitMutation = useContentImportCommit()

  function reset() {
    setSelectedFile(null)
    setReport(null)
    if (fileInputRef.current) {
      fileInputRef.current.value = ''
    }
  }

  const groups = useMemo(() => groupItems(report?.items ?? []), [report])
  const topics = report?.items.filter((item) => item.kind === 'TOPIC') ?? []
  const validCount = report ? report.summary.created + report.summary.updated : 0
  const canConfirm =
    report && !report.committed && report.bundleErrors.length === 0 && validCount > 0

  return (
    <div className="space-y-6">
      <header className="sticky top-0 z-10 -mx-4 border-b border-border bg-background/95 px-4 backdrop-blur supports-[backdrop-filter]:bg-background/80 sm:-mx-6 sm:px-6 lg:-mx-8 lg:px-8">
        <div className="mx-auto flex h-16 max-w-4xl items-center gap-4">
          <button
            className="cursor-pointer rounded-xl p-2 transition-colors hover:bg-muted/60"
            onClick={() => {
              void navigate('/dashboard')
            }}
            type="button"
          >
            <ChevronLeft className="size-5 text-muted-foreground" />
          </button>
          <span className="text-sm font-semibold text-foreground">Import Lesson Content</span>
        </div>
      </header>

      <div className="mx-auto max-w-4xl space-y-6">
        <section className="rounded-2xl border border-border bg-white p-5">
          <h2 className="text-sm font-semibold text-foreground">Upload content bundle (.json)</h2>
          <p className="mt-1 text-xs text-muted-foreground">
            A bundle contains topics → lessons → activities → exercises for one course level
            (format: <code>docs/foundation-spec.md</code> §8.2, example:{' '}
            <code>content/a1/week-1.json</code>). Items are matched by slug, so importing the same
            file again updates it. Preview never writes to the database.
          </p>

          <input
            accept=".json,application/json"
            className="mt-4 block w-full cursor-pointer text-sm text-muted-foreground file:mr-3 file:cursor-pointer file:rounded-lg file:border-0 file:bg-muted file:px-3 file:py-2 file:text-sm file:font-medium file:text-foreground hover:file:bg-muted/70"
            onChange={(event) => {
              setSelectedFile(event.target.files?.[0] ?? null)
              setReport(null)
            }}
            ref={fileInputRef}
            type="file"
          />

          <div className="mt-4 flex flex-wrap gap-3">
            <Button
              disabled={!selectedFile || previewMutation.isPending}
              onClick={() => {
                if (selectedFile) {
                  previewMutation.mutate(selectedFile, { onSuccess: setReport })
                }
              }}
              type="button"
            >
              <Upload className="size-4" />
              {previewMutation.isPending ? 'Checking file' : 'Preview import'}
            </Button>
            {canConfirm ? (
              <Button
                disabled={commitMutation.isPending}
                onClick={() => {
                  if (selectedFile) {
                    commitMutation.mutate(selectedFile, { onSuccess: setReport })
                  }
                }}
                type="button"
                variant="success"
              >
                {commitMutation.isPending ? 'Importing' : `Confirm import (${validCount} items)`}
              </Button>
            ) : null}
            {selectedFile || report ? (
              <Button onClick={reset} type="button" variant="outline">
                Start over
              </Button>
            ) : null}
          </div>
        </section>

        {report ? (
          <section
            className={cn(
              'rounded-2xl border p-5',
              report.committed ? 'border-emerald-200 bg-emerald-50' : 'border-border bg-white',
            )}
          >
            <div className="mb-3 flex flex-wrap items-center justify-between gap-3">
              <h2 className="flex items-center gap-2 text-sm font-semibold text-foreground">
                {report.committed ? <CheckCircle2 className="size-4 text-emerald-600" /> : null}
                {report.committed ? 'Import complete' : 'Preview result'}
              </h2>
              <div className="flex flex-wrap gap-2 text-xs">
                <Badge variant="outline">{report.summary.lessons} lessons</Badge>
                <Badge variant="outline">{report.summary.activities} activities</Badge>
                <Badge variant="outline">{report.summary.exercises} exercises</Badge>
                <Badge variant="success">
                  {report.summary.created} new · {report.summary.updated} updated
                </Badge>
                {report.summary.invalid > 0 ? (
                  <Badge variant="destructive">{report.summary.invalid} invalid</Badge>
                ) : null}
              </div>
            </div>

            {report.bundleErrors.map((error) => (
              <p className="mb-2 flex gap-2 text-sm text-destructive" key={error}>
                <TriangleAlert className="mt-0.5 size-4 shrink-0" />
                {error}
              </p>
            ))}
            {report.warnings.map((warning) => (
              <p className="mb-2 text-sm text-amber-700" key={warning}>
                {warning}
              </p>
            ))}

            {topics.map((topic) => (
              <p className="mb-3 text-sm" key={topic.path}>
                <ItemStatus item={topic} /> <span className="font-semibold">{topic.title}</span>{' '}
                <span className="text-muted-foreground">({topic.path})</span>
              </p>
            ))}

            <div className="space-y-2">
              {groups.map((group) => (
                <LessonRow group={group} key={group.lesson.path} />
              ))}
            </div>
          </section>
        ) : null}
      </div>
    </div>
  )
}

function LessonRow({ group }: { group: LessonGroup }) {
  const problems = [group.lesson, ...group.activities.flatMap((a) => [a.activity, ...a.exercises])].filter(
    (item) => item.errors.length > 0,
  ).length
  // Lessons with problems start expanded so errors are visible.
  const [open, setOpen] = useState(problems > 0)
  const exerciseCount = group.activities.reduce((total, a) => total + a.exercises.length, 0)

  return (
    <div className="rounded-xl border border-border bg-white">
      <button
        className="flex w-full cursor-pointer items-center gap-2 px-3 py-2.5 text-left text-sm"
        onClick={() => setOpen(!open)}
        type="button"
      >
        <ChevronRight className={cn('size-4 shrink-0 transition-transform', open && 'rotate-90')} />
        <ItemStatus item={group.lesson} />
        <span className="font-semibold">{group.lesson.title}</span>
        <span className="ml-auto shrink-0 text-xs text-muted-foreground">
          {group.activities.length} activities · {exerciseCount} exercises
          {problems > 0 ? <span className="ml-2 text-destructive">{problems} problems</span> : null}
        </span>
      </button>
      {open ? (
        <div className="space-y-2 border-t px-3 py-2">
          {group.lesson.errors.length > 0 ? <Errors item={group.lesson} /> : null}
          {group.activities.map(({ activity, exercises }) => (
            <div key={activity.path}>
              <p className="text-sm">
                <ItemStatus item={activity} /> <span className="font-medium">{activity.title}</span>{' '}
                {activity.status === 'DRAFT' ? <Badge variant="outline">draft</Badge> : null}
              </p>
              <Errors item={activity} />
              <ul className="ml-6 mt-1 space-y-0.5">
                {exercises.map((exercise) => (
                  <li className="text-xs text-muted-foreground" key={exercise.path}>
                    <ItemStatus item={exercise} /> {exercise.title}
                    <Errors item={exercise} />
                  </li>
                ))}
              </ul>
            </div>
          ))}
        </div>
      ) : null}
    </div>
  )
}

function ItemStatus({ item }: { item: ContentImportItem }) {
  const styles: Record<ContentImportItem['action'], string> = {
    CREATE: 'bg-emerald-100 text-emerald-700',
    UPDATE: 'bg-blue-100 text-blue-700',
    SKIP: 'bg-red-100 text-red-700',
  }
  return (
    <span className={cn('rounded px-1.5 py-0.5 text-[10px] font-bold', styles[item.action])}>
      {item.action}
    </span>
  )
}

function Errors({ item }: { item: ContentImportItem }) {
  if (item.errors.length === 0) {
    return null
  }
  return (
    <span className="mt-0.5 flex items-start gap-1 text-xs text-destructive">
      <TriangleAlert className="mt-0.5 size-3 shrink-0" />
      {item.errors.join('; ')}
    </span>
  )
}

function groupItems(items: ContentImportItem[]): LessonGroup[] {
  const groups: LessonGroup[] = []
  for (const item of items) {
    if (item.kind === 'LESSON') {
      groups.push({ lesson: item, activities: [] })
    } else if (item.kind === 'ACTIVITY') {
      groups.at(-1)?.activities.push({ activity: item, exercises: [] })
    } else if (item.kind === 'EXERCISE') {
      groups.at(-1)?.activities.at(-1)?.exercises.push(item)
    }
  }
  return groups
}
