import { useMutation } from "@tanstack/react-query";
import { ChevronLeft, Download, FileSearch, TriangleAlert } from "lucide-react";
import { useState } from "react";
import { useNavigate } from "react-router-dom";

import { getApiErrorMessage } from "@/api/api";
import { extractContent } from "@/api/contentImport";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { downloadText } from "@/lib/download";

/**
 * Task 65: upload textbook pages (images or one PDF) → the AI extracts theory
 * and exercises into a content bundle draft → review it in Import Lessons.
 */
export function AdminContentExtractPage() {
  const navigate = useNavigate();
  const [files, setFiles] = useState<File[]>([]);
  const [form, setForm] = useState({ topicTitle: "", topicSlug: "", notes: "", publish: false });
  const extract = useMutation({ mutationFn: extractContent });
  const bundleText = extract.data ? JSON.stringify(extract.data.bundle, null, 2) : "";
  const fileName = `${form.topicSlug || "extracted"}.json`;

  return (
    <div className="space-y-6">
      <header className="sticky top-0 z-10 -mx-4 border-b border-border bg-background/95 px-4 backdrop-blur sm:-mx-6 sm:px-6 lg:-mx-8 lg:px-8">
        <div className="mx-auto flex h-16 max-w-4xl items-center gap-4">
          <button className="cursor-pointer rounded-xl p-2 hover:bg-muted/60" onClick={() => void navigate("/dashboard")} type="button">
            <ChevronLeft className="size-5 text-muted-foreground" />
          </button>
          <span className="text-sm font-semibold">Extract Lessons from PDF / Images</span>
        </div>
      </header>

      <div className="mx-auto max-w-4xl space-y-6">
        <section className="rounded-2xl border bg-white p-5">
          <h2 className="text-sm font-semibold">Upload textbook pages</h2>
          <p className="mt-1 text-xs text-muted-foreground">
            Up to 10 page photos/scans, or one PDF (max 20MB). The AI reads the grammar explanations and exercises and writes a
            lesson bundle (one lesson in an always-open topic). Nothing is saved until you import it; activities stay DRAFT
            (hidden from learners) unless you tick "publish".
          </p>
          <input
            accept="image/*,.pdf,application/pdf"
            className="mt-4 block w-full cursor-pointer text-sm text-muted-foreground file:mr-3 file:cursor-pointer file:rounded-lg file:border-0 file:bg-muted file:px-3 file:py-2 file:text-sm file:font-medium file:text-foreground"
            multiple
            onChange={(event) => setFiles(Array.from(event.target.files ?? []))}
            type="file"
          />
          {files.length > 0 ? (
            <p className="mt-2 text-xs text-muted-foreground">{files.map((file) => file.name).join(", ")}</p>
          ) : null}
          <div className="mt-4 grid gap-3 sm:grid-cols-2">
            <Input
              onChange={(event) => setForm({ ...form, topicTitle: event.target.value })}
              placeholder="Topic title, e.g. Essential Grammar · Unit 5"
              value={form.topicTitle}
            />
            <Input
              onChange={(event) => setForm({ ...form, topicSlug: event.target.value })}
              placeholder="Topic slug, e.g. egu-unit-5"
              value={form.topicSlug}
            />
            <Input
              className="sm:col-span-2"
              onChange={(event) => setForm({ ...form, notes: event.target.value })}
              placeholder="Notes for the AI (optional), e.g. only exercises 5.1-5.3"
              value={form.notes}
            />
          </div>
          <label className="mt-3 flex items-center gap-2 text-sm">
            <input
              checked={form.publish}
              onChange={(event) => setForm({ ...form, publish: event.target.checked })}
              type="checkbox"
            />
            Publish activities right away (skip DRAFT)
          </label>
          <Button
            className="mt-4"
            disabled={files.length === 0 || extract.isPending}
            onClick={() => extract.mutate({ files, ...form })}
            type="button"
          >
            <FileSearch className="size-4" />
            {extract.isPending ? "Reading pages… (up to a few minutes)" : "Extract"}
          </Button>
          {extract.error ? (
            <p className="mt-3 text-sm text-destructive">{getApiErrorMessage(extract.error, "Extraction failed.")}</p>
          ) : null}
        </section>

        {extract.data ? (
          <section className="space-y-4 rounded-2xl border bg-white p-5">
            <h2 className="text-sm font-semibold">Result: {extract.data.summary}</h2>
            {extract.data.warnings.map((warning) => (
              <p className="flex gap-2 text-sm text-amber-700" key={warning}>
                <TriangleAlert className="mt-0.5 size-4 shrink-0" />
                {warning}
              </p>
            ))}
            <p className="text-xs text-muted-foreground">
              Check every answer before importing — the AI can misread a page. The import preview self-checks each exercise.
            </p>
            <div className="flex flex-wrap gap-3">
              <Button
                onClick={() => void navigate("/admin/content/import", { state: { bundleText, fileName } })}
                type="button"
              >
                Review in Import Lessons
              </Button>
              <Button onClick={() => downloadText(fileName, bundleText)} type="button" variant="outline">
                <Download className="size-4" /> Download JSON
              </Button>
            </div>
            <details>
              <summary className="cursor-pointer text-sm font-semibold text-primary">Show JSON</summary>
              <pre className="mt-2 max-h-96 overflow-auto rounded-xl bg-muted/50 p-3 text-xs">{bundleText}</pre>
            </details>
          </section>
        ) : null}
      </div>
    </div>
  );
}
