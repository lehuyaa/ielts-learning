import { CheckCircle2, ChevronLeft, TriangleAlert, Upload } from 'lucide-react'
import { useRef, useState } from 'react'
import { useNavigate } from 'react-router-dom'

import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import {
  useVocabularyImportCommit,
  useVocabularyImportPreview,
} from '@/features/admin/hooks/useVocabularyImport'
import { cn } from '@/lib/utils'
import type { ImportResultResponse } from '@/types/vocabulary'

export function AdminVocabularyImportPage() {
  const navigate = useNavigate()
  const fileInputRef = useRef<HTMLInputElement>(null)
  const [selectedFile, setSelectedFile] = useState<File | null>(null)
  const [preview, setPreview] = useState<ImportResultResponse | null>(null)
  const [result, setResult] = useState<ImportResultResponse | null>(null)

  const previewMutation = useVocabularyImportPreview()
  const commitMutation = useVocabularyImportCommit()

  function handleFileChange(event: React.ChangeEvent<HTMLInputElement>) {
    const file = event.target.files?.[0] ?? null
    setSelectedFile(file)
    setPreview(null)
    setResult(null)
  }

  async function handlePreview() {
    if (!selectedFile) {
      return
    }

    setResult(null)
    const data = await previewMutation.mutateAsync(selectedFile)
    setPreview(data)
  }

  async function handleConfirmImport() {
    if (!selectedFile) {
      return
    }

    const data = await commitMutation.mutateAsync(selectedFile)
    setResult(data)
    setPreview(null)
  }

  function handleReset() {
    setSelectedFile(null)
    setPreview(null)
    setResult(null)
    if (fileInputRef.current) {
      fileInputRef.current.value = ''
    }
  }

  const summary = result?.summary ?? preview?.summary
  const rows = result?.rows ?? preview?.rows ?? []
  const canConfirm = Boolean(preview) && preview!.summary.validRows > 0

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
          <span className="text-sm font-semibold text-foreground">
            Import Vocabulary
          </span>
        </div>
      </header>

      <div className="mx-auto max-w-4xl space-y-6">
        <section className="rounded-2xl border border-border bg-white p-5">
          <div className="mb-4">
            <h2 className="text-sm font-semibold text-foreground">
              Upload Excel file
            </h2>
            <p className="mt-1 text-xs text-muted-foreground">
              Required columns: <code>word</code>, <code>meaningVi</code>,{' '}
              <code>topicSlug</code>, <code>lessonSlug</code>. Optional:{' '}
              <code>ipa</code>, <code>partOfSpeech</code>, <code>meaningEn</code>,{' '}
              <code>shortDefinition</code>, <code>exampleSentence</code>,{' '}
              <code>exampleMeaningVi</code>, <code>exampleSource</code>,{' '}
              <code>synonyms</code>, <code>antonyms</code>,{' '}
              <code>collocations</code>, <code>difficulty</code>,{' '}
              <code>targetBand</code>, <code>orderIndex</code>,{' '}
              <code>isRequired</code>. Each row assigns the word to the given
              lesson within the given topic.
            </p>
          </div>

          <div className="flex flex-col gap-3 sm:flex-row sm:items-center">
            <input
              accept=".xlsx,.xls"
              className="block w-full cursor-pointer text-sm text-muted-foreground file:mr-3 file:cursor-pointer file:rounded-lg file:border-0 file:bg-muted file:px-3 file:py-2 file:text-sm file:font-medium file:text-foreground hover:file:bg-muted/70"
              onChange={handleFileChange}
              ref={fileInputRef}
              type="file"
            />
          </div>

          <div className="mt-4 flex flex-wrap gap-3">
            <Button
              disabled={!selectedFile || previewMutation.isPending}
              onClick={() => {
                void handlePreview()
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
                  void handleConfirmImport()
                }}
                type="button"
                variant="success"
              >
                {commitMutation.isPending
                  ? 'Importing'
                  : `Confirm import (${preview!.summary.validRows} words)`}
              </Button>
            ) : null}

            {selectedFile || preview || result ? (
              <Button onClick={handleReset} type="button" variant="outline">
                Start over
              </Button>
            ) : null}
          </div>
        </section>

        {result ? (
          <section className="rounded-2xl border border-emerald-200 bg-emerald-50 p-5">
            <div className="flex items-center gap-2 text-emerald-800">
              <CheckCircle2 className="size-4" />
              <h2 className="text-sm font-semibold">Import complete</h2>
            </div>
            <div className="mt-3 grid grid-cols-2 gap-3 text-sm text-emerald-800 sm:grid-cols-4">
              <div>
                <div className="text-lg font-bold">
                  {result.vocabulariesCreated}
                </div>
                <div className="text-xs">Words created</div>
              </div>
              <div>
                <div className="text-lg font-bold">
                  {result.vocabulariesUpdated}
                </div>
                <div className="text-xs">Words updated</div>
              </div>
              <div>
                <div className="text-lg font-bold">
                  {result.lessonLinksCreated}
                </div>
                <div className="text-xs">Lesson links created</div>
              </div>
              <div>
                <div className="text-lg font-bold">
                  {result.lessonLinksUpdated}
                </div>
                <div className="text-xs">Lesson links updated</div>
              </div>
            </div>
          </section>
        ) : null}

        {summary ? (
          <section className="rounded-2xl border border-border bg-white p-5">
            <div className="mb-4 flex flex-wrap items-center justify-between gap-3">
              <h2 className="text-sm font-semibold text-foreground">
                {result ? 'Import result' : 'Preview result'}
              </h2>
              <div className="flex gap-2 text-xs">
                <Badge variant="outline">{summary.totalRows} rows</Badge>
                <Badge variant="success">{summary.validRows} valid</Badge>
                {summary.invalidRows > 0 ? (
                  <Badge variant="destructive">
                    {summary.invalidRows} invalid
                  </Badge>
                ) : null}
              </div>
            </div>

            <div className="overflow-x-auto">
              <table className="w-full min-w-[640px] text-left text-sm">
                <thead>
                  <tr className="border-b border-border text-xs text-muted-foreground">
                    <th className="py-2 pr-3 font-medium">Row</th>
                    <th className="py-2 pr-3 font-medium">Word</th>
                    <th className="py-2 pr-3 font-medium">Topic / Lesson</th>
                    <th className="py-2 pr-3 font-medium">Status</th>
                  </tr>
                </thead>
                <tbody>
                  {rows.map((row) => (
                    <tr className="border-b border-border/60" key={row.row}>
                      <td className="py-2 pr-3 text-muted-foreground">
                        {row.row}
                      </td>
                      <td className="py-2 pr-3 font-medium text-foreground">
                        {row.word || '—'}
                      </td>
                      <td className="py-2 pr-3 text-muted-foreground">
                        {row.topicSlug} / {row.lessonSlug}
                      </td>
                      <td className="py-2 pr-3">
                        {row.valid ? (
                          <span className="inline-flex items-center gap-1 text-emerald-600">
                            <CheckCircle2 className="size-3.5" />
                            {row.vocabularyAction || row.linkAction
                              ? [row.vocabularyAction, row.linkAction]
                                  .filter(Boolean)
                                  .join(' / ')
                              : 'Valid'}
                          </span>
                        ) : (
                          <span
                            className={cn(
                              'inline-flex items-start gap-1 text-destructive',
                            )}
                          >
                            <TriangleAlert className="mt-0.5 size-3.5 shrink-0" />
                            <span>{row.errors?.join('; ')}</span>
                          </span>
                        )}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </section>
        ) : null}
      </div>
    </div>
  )
}
