import { BookmarkCheck, BookmarkPlus, CheckCircle2, Lightbulb, Sparkles, Star, XCircle } from "lucide-react";
import { useState } from "react";

import { addWordToFlashcards, lookupWord } from "@/api/passage";
import { cn } from "@/lib/utils";
import type { SpeakingReport } from "@/types/aiConversation";

const labels: [keyof SpeakingReport["scores"], string][] = [
  ["taskAchievement", "Hoàn thành nhiệm vụ"],
  ["grammar", "Ngữ pháp"],
  ["vocabulary", "Từ vựng"],
  ["fluency", "Trôi chảy"],
];

/** Feedback report of a finished role-play (Task 69). */
export function SessionReport({ report, passScore }: { report: SpeakingReport; passScore?: number }) {
  const passed = report.taskCompleted && (passScore === undefined || report.overall >= passScore);
  return (
    <div className="space-y-4">
      <section className="rounded-3xl border bg-white p-5 shadow-sm">
        <div className="flex flex-wrap items-center gap-5">
          <div className={cn("text-center", passed ? "text-emerald-600" : "text-amber-600")}>
            <p className="text-5xl font-bold tracking-normal">{report.overall}</p>
            <p className="text-xs font-semibold">điểm{passScore !== undefined ? ` · cần ${passScore}` : ""}</p>
          </div>
          <div className="min-w-56 flex-1 space-y-1.5">
            {labels.map(([key, label]) => (
              <div className="flex items-center gap-3 text-sm" key={key}>
                <span className="w-36 shrink-0 text-muted-foreground">{label}</span>
                <div className="h-2 flex-1 rounded-full bg-muted">
                  <div
                    className={cn("h-2 rounded-full", report.scores[key] >= 8 ? "bg-emerald-500" : report.scores[key] >= 6 ? "bg-amber-400" : "bg-red-400")}
                    style={{ width: `${report.scores[key] * 10}%` }}
                  />
                </div>
                <span className="w-6 text-right font-semibold">{report.scores[key]}</span>
              </div>
            ))}
          </div>
        </div>
        <div className="mt-4 flex flex-wrap gap-2 text-sm">
          <span className={cn("inline-flex items-center gap-1.5 rounded-full px-3 py-1 font-semibold", report.taskCompleted ? "bg-emerald-50 text-emerald-700" : "bg-amber-50 text-amber-700")}>
            {report.taskCompleted ? <CheckCircle2 className="size-4" /> : <XCircle className="size-4" />}
            {report.taskCompleted ? "Đã hoàn thành nhiệm vụ" : "Chưa hoàn thành nhiệm vụ"}
          </span>
          {report.xpAwarded > 0 ? (
            <span className="inline-flex items-center gap-1.5 rounded-full bg-amber-50 px-3 py-1 font-semibold text-amber-700">
              <Star className="size-4 fill-amber-400 text-amber-400" /> +{report.xpAwarded} XP
            </span>
          ) : null}
        </div>
        {report.tipVi ? (
          <p className="mt-4 flex gap-2 rounded-2xl bg-violet-50 p-3 text-sm">
            <Lightbulb className="mt-0.5 size-4 shrink-0 text-violet-600" />
            {report.tipVi}
          </p>
        ) : null}
      </section>

      {report.mistakes.length > 0 ? (
        <section className="rounded-3xl border bg-white p-5">
          <h3 className="mb-3 font-bold">Câu cần sửa ({report.mistakes.length})</h3>
          <ul className="space-y-3">
            {report.mistakes.map((mistake, index) => (
              <li className="rounded-2xl bg-muted/30 p-3" key={index}>
                <p className="text-red-600 line-through">{mistake.said}</p>
                <p className="font-semibold text-emerald-700">{mistake.correct}</p>
                {mistake.explanationVi ? <p className="mt-1 text-sm text-muted-foreground">{mistake.explanationVi}</p> : null}
              </li>
            ))}
          </ul>
          <p className="mt-3 text-xs text-muted-foreground">Các lỗi này đã được lưu vào sổ lỗi sai để ôn lại.</p>
        </section>
      ) : null}

      <div className="grid gap-4 sm:grid-cols-2">
        {report.goodPhrases.length > 0 ? (
          <section className="rounded-3xl border bg-white p-5">
            <h3 className="mb-3 flex items-center gap-2 font-bold">
              <Sparkles className="size-4 text-emerald-600" /> Câu bạn nói tốt
            </h3>
            <ul className="space-y-1.5 text-sm">
              {report.goodPhrases.map((phrase) => (
                <li className="rounded-xl bg-emerald-50 px-3 py-2" key={phrase}>
                  {phrase}
                </li>
              ))}
            </ul>
          </section>
        ) : null}
        {report.newWords.length > 0 ? (
          <section className="rounded-3xl border bg-white p-5">
            <h3 className="mb-3 font-bold">Từ mới</h3>
            <ul className="space-y-1.5">
              {report.newWords.map((word) => (
                <NewWord key={word} word={word} />
              ))}
            </ul>
          </section>
        ) : null}
      </div>
    </div>
  );
}

type WordState = "idle" | "adding" | "added" | "missing";

function NewWord({ word }: { word: string }) {
  const [state, setState] = useState<WordState>("idle");

  async function add() {
    setState("adding");
    try {
      const found = await lookupWord(word);
      if (!found.vocabulary) {
        setState("missing");
        return;
      }
      await addWordToFlashcards(found.vocabulary.id);
      setState("added");
    } catch {
      setState("idle");
    }
  }

  return (
    <li className="flex items-center justify-between gap-2 rounded-xl bg-muted/30 px-3 py-1.5">
      <span className="font-semibold">{word}</span>
      <button
        className="inline-flex cursor-pointer items-center gap-1 rounded-lg px-2 py-1 text-xs font-semibold text-primary hover:bg-primary/10 disabled:cursor-default disabled:text-muted-foreground disabled:hover:bg-transparent"
        disabled={state !== "idle"}
        onClick={() => void add()}
        type="button"
      >
        {state === "added" ? <BookmarkCheck className="size-4" /> : <BookmarkPlus className="size-4" />}
        {state === "added" ? "Đã thêm" : state === "missing" ? "Chưa có trong từ điển" : state === "adding" ? "Đang thêm…" : "Thêm flashcard"}
      </button>
    </li>
  );
}
