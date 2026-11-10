import { useMutation } from "@tanstack/react-query";
import { ChevronLeft, Download, Mic, Save, Trash2 } from "lucide-react";
import { useRef, useState } from "react";
import { useNavigate } from "react-router-dom";

import { getApiErrorMessage } from "@/api/api";
import { savePassage, transcribeAudio, type PassageDraft } from "@/api/contentImport";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { AudioPlayer } from "@/features/audio/AudioPlayer";
import { downloadText } from "@/lib/download";
import type { PassageCard, PassageGenre } from "@/types/passage";

const GENRES: PassageGenre[] = ["DIALOGUE", "MONOLOGUE", "MESSAGE", "EMAIL", "NOTICE", "MENU", "SCHEDULE", "ARTICLE", "PROFILE"];

/** Task 64: upload a recording → Whisper transcript with sentence times → edit → save as a passage. */
export function AdminPassageImportPage() {
  const navigate = useNavigate();
  const fileInputRef = useRef<HTMLInputElement>(null);
  const [file, setFile] = useState<File | null>(null);
  const [form, setForm] = useState({ title: "", slug: "", genre: "DIALOGUE", source: "" });
  const [draft, setDraft] = useState<PassageDraft | null>(null);

  const transcribe = useMutation({ mutationFn: transcribeAudio, onSuccess: setDraft });
  const save = useMutation({ mutationFn: savePassage });

  function update(changes: Partial<PassageDraft>) {
    setDraft((current) => (current ? { ...current, ...changes } : current));
    save.reset();
  }

  function updateSegment(index: number, changes: Partial<PassageDraft["segments"][number]>) {
    if (!draft) return;
    update({ segments: draft.segments.map((segment, i) => (i === index ? { ...segment, ...changes } : segment)) });
  }

  function removeSegment(index: number) {
    if (!draft) return;
    update({ segments: draft.segments.filter((_, i) => i !== index).map((segment, i) => ({ ...segment, index: i })) });
  }

  /** The passage as it will be saved: body text rebuilt from the segments. */
  function finalDraft(): Omit<PassageDraft, "duration"> | null {
    if (!draft) return null;
    const segments = draft.segments
      .filter((segment) => segment.text.trim())
      .map((segment, index) => ({ ...segment, index, speaker: segment.speaker?.trim() || undefined, text: segment.text.trim() }));
    const bodyText = segments.map((segment) => (segment.speaker ? `${segment.speaker}: ${segment.text}` : segment.text)).join("\n");
    return {
      slug: draft.slug,
      cefrLevel: draft.cefrLevel,
      genre: draft.genre,
      title: draft.title,
      audioUrl: draft.audioUrl,
      source: draft.source,
      segments,
      bodyText,
      glossary: draft.glossary ?? [],
    };
  }

  const preview: PassageCard | null = draft
    ? {
        id: 0,
        slug: draft.slug,
        cefrLevel: draft.cefrLevel,
        genre: draft.genre as PassageGenre,
        title: draft.title,
        bodyText: draft.bodyText,
        audioUrl: draft.audioUrl,
        segments: draft.segments,
        glossary: [],
        source: draft.source,
        wordCount: 0,
      }
    : null;

  return (
    <div className="space-y-6">
      <header className="sticky top-0 z-10 -mx-4 border-b border-border bg-background/95 px-4 backdrop-blur sm:-mx-6 sm:px-6 lg:-mx-8 lg:px-8">
        <div className="mx-auto flex h-16 max-w-4xl items-center gap-4">
          <button className="cursor-pointer rounded-xl p-2 hover:bg-muted/60" onClick={() => void navigate("/dashboard")} type="button">
            <ChevronLeft className="size-5 text-muted-foreground" />
          </button>
          <span className="text-sm font-semibold">Import Audio (Whisper)</span>
        </div>
      </header>

      <div className="mx-auto max-w-4xl space-y-6">
        <section className="rounded-2xl border bg-white p-5">
          <h2 className="text-sm font-semibold">1. Upload a recording</h2>
          <p className="mt-1 text-xs text-muted-foreground">
            mp3 / m4a / wav / webm, max 25MB. Whisper writes the transcript and splits it into sentences with times, so learners
            can replay sentence by sentence and use them for dictation. The file is stored under <code>/media/passages</code> on the
            server.
          </p>
          <div className="mt-4 grid gap-3 sm:grid-cols-2">
            <input
              accept="audio/*,.mp3,.m4a,.wav,.webm,.ogg"
              className="block w-full cursor-pointer text-sm text-muted-foreground file:mr-3 file:cursor-pointer file:rounded-lg file:border-0 file:bg-muted file:px-3 file:py-2 file:text-sm file:font-medium file:text-foreground sm:col-span-2"
              onChange={(event) => setFile(event.target.files?.[0] ?? null)}
              ref={fileInputRef}
              type="file"
            />
            <Input onChange={(event) => setForm({ ...form, title: event.target.value })} placeholder="Title, e.g. At the cafe" value={form.title} />
            <Input onChange={(event) => setForm({ ...form, slug: event.target.value })} placeholder="Slug (optional, from the title)" value={form.slug} />
            <select
              className="h-10 rounded-md border border-input bg-white px-3 text-sm"
              onChange={(event) => setForm({ ...form, genre: event.target.value })}
              value={form.genre}
            >
              {GENRES.map((genre) => (
                <option key={genre}>{genre}</option>
              ))}
            </select>
            <Input onChange={(event) => setForm({ ...form, source: event.target.value })} placeholder="Source, e.g. English File Elementary 2.3" value={form.source} />
          </div>
          <Button
            className="mt-4"
            disabled={!file || transcribe.isPending}
            onClick={() => file && transcribe.mutate({ file, ...form })}
            type="button"
          >
            <Mic className="size-4" />
            {transcribe.isPending ? "Transcribing… (up to a minute)" : "Transcribe"}
          </Button>
          {transcribe.error ? (
            <p className="mt-3 text-sm text-destructive">{getApiErrorMessage(transcribe.error, "Transcription failed.")}</p>
          ) : null}
        </section>

        {draft && preview ? (
          <section className="space-y-4 rounded-2xl border bg-white p-5">
            <h2 className="text-sm font-semibold">2. Check the transcript</h2>
            <div className="grid gap-3 sm:grid-cols-3">
              <Input onChange={(event) => update({ title: event.target.value })} value={draft.title} />
              <Input onChange={(event) => update({ slug: event.target.value })} value={draft.slug} />
              <select
                className="h-10 rounded-md border border-input bg-white px-3 text-sm"
                onChange={(event) => update({ genre: event.target.value })}
                value={draft.genre}
              >
                {GENRES.map((genre) => (
                  <option key={genre}>{genre}</option>
                ))}
              </select>
            </div>

            <div className="rounded-2xl bg-muted/30 p-4">
              <AudioPlayer key={draft.segments.length} passage={preview} passed transcript="AFTER_ANSWER" />
            </div>

            <p className="text-xs text-muted-foreground">
              Fix words Whisper misheard, add the speaker (A, B…) for dialogues, delete noise. Times stay as transcribed.
            </p>
            <ol className="space-y-2">
              {draft.segments.map((segment, index) => (
                <li className="flex items-center gap-2" key={index}>
                  <span className="w-24 shrink-0 font-mono text-xs text-muted-foreground">
                    {segment.start?.toFixed(1)}–{segment.end?.toFixed(1)}s
                  </span>
                  <Input
                    className="w-16 shrink-0"
                    onChange={(event) => updateSegment(index, { speaker: event.target.value })}
                    placeholder="A"
                    value={segment.speaker ?? ""}
                  />
                  <Input onChange={(event) => updateSegment(index, { text: event.target.value })} value={segment.text} />
                  <button
                    aria-label="Delete sentence"
                    className="cursor-pointer rounded-lg p-2 text-muted-foreground hover:bg-red-50 hover:text-red-600"
                    onClick={() => removeSegment(index)}
                    type="button"
                  >
                    <Trash2 className="size-4" />
                  </button>
                </li>
              ))}
            </ol>

            <div className="flex flex-wrap gap-3">
              <Button
                disabled={save.isPending}
                onClick={() => {
                  const passage = finalDraft();
                  if (passage) save.mutate(passage);
                }}
                type="button"
                variant="success"
              >
                <Save className="size-4" /> {save.isPending ? "Saving" : "Save passage"}
              </Button>
              <Button
                onClick={() => {
                  const passage = finalDraft();
                  if (passage) downloadText(`${passage.slug}.json`, JSON.stringify({ passages: [passage] }, null, 2));
                }}
                type="button"
                variant="outline"
              >
                <Download className="size-4" /> Download JSON
              </Button>
              <Button
                onClick={() => {
                  setDraft(null);
                  setFile(null);
                  save.reset();
                  if (fileInputRef.current) fileInputRef.current.value = "";
                }}
                type="button"
                variant="outline"
              >
                Start over
              </Button>
            </div>
            {save.data ? (
              <p className="rounded-xl bg-emerald-50 p-3 text-sm text-emerald-800">
                Saved passage #{save.data.id}. Use it in a bundle with <code>"passageSlug": "{save.data.slug}"</code> on a LISTENING
                activity, and dictation items with <code>{`"passageSlug": "${save.data.slug}", "segmentIndex": 0`}</code>.
              </p>
            ) : null}
            {save.error ? <p className="text-sm text-destructive">{getApiErrorMessage(save.error, "Could not save.")}</p> : null}
          </section>
        ) : null}
      </div>
    </div>
  );
}
