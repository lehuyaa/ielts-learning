import { useMutation } from "@tanstack/react-query";
import { CheckCircle2, PenLine, RotateCcw, Send, Sparkles } from "lucide-react";
import { useRef, useState } from "react";

import { getApiErrorMessage } from "@/api/api";
import { rewriteWriting, submitWriting } from "@/api/writing";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";
import type { WritingContent, WritingSubmission } from "@/types/writing";

import { MistakeText } from "./MistakeText";
import { ScoreBars } from "./ScoreBars";
import { countSentences } from "./sentences";

type WritingTaskProps = {
  activityId: number;
  content: WritingContent;
  initial: WritingSubmission | null;
  finishing: boolean;
  onFinish: () => void;
};

/**
 * WRITING activity (Task 66): prompt with sentence frames → AI report with
 * mistakes marked in the text → rewrite → finish.
 */
export function WritingTask({ activityId, content, initial, finishing, onFinish }: WritingTaskProps) {
  const [submission, setSubmission] = useState<WritingSubmission | null>(initial);
  const [writingNew, setWritingNew] = useState(false);

  const showEditor = !submission || writingNew;
  return (
    <div className="space-y-5">
      <PromptCard content={content} />
      {showEditor ? (
        <Editor
          activityId={activityId}
          content={content}
          onSubmitted={(result) => {
            setSubmission(result);
            setWritingNew(false);
          }}
        />
      ) : null}
      {submission && !writingNew ? (
        <Result
          finishing={finishing}
          onFinish={onFinish}
          onUpdated={setSubmission}
          onWriteNew={() => setWritingNew(true)}
          submission={submission}
        />
      ) : null}
    </div>
  );
}

function PromptCard({ content }: { content: WritingContent }) {
  return (
    <section className="rounded-3xl border bg-white p-5 shadow-sm">
      <p className="flex items-center gap-2 text-xs font-bold uppercase text-primary">
        <PenLine className="size-4" /> Đề bài
      </p>
      <p className="mt-2 text-lg font-semibold">{content.promptVi}</p>
      {content.promptEn ? <p className="mt-1 text-sm text-muted-foreground">{content.promptEn}</p> : null}
      <p className="mt-2 text-xs text-muted-foreground">
        Viết {content.minSentences ?? 3}–{content.maxSentences ?? 7} câu.
      </p>
    </section>
  );
}

function Editor({
  activityId,
  content,
  onSubmitted,
}: {
  activityId: number;
  content: WritingContent;
  onSubmitted: (submission: WritingSubmission) => void;
}) {
  const [text, setText] = useState("");
  const textareaRef = useRef<HTMLTextAreaElement>(null);
  const submit = useMutation({ mutationFn: submitWriting, onSuccess: onSubmitted });
  const sentences = countSentences(text);
  const min = content.minSentences ?? 3;
  const max = content.maxSentences ?? 7;

  /** Inserts a frame at the cursor ("..." becomes the selection to replace). */
  function insertFrame(frame: string) {
    const area = textareaRef.current;
    const start = area?.selectionStart ?? text.length;
    const before = text.slice(0, start);
    const spacer = before && !/\s$/.test(before) ? " " : "";
    const next = before + spacer + frame + text.slice(area?.selectionEnd ?? start);
    setText(next);
    requestAnimationFrame(() => {
      if (!area) return;
      const dots = next.indexOf("...", before.length);
      area.focus();
      if (dots >= 0) area.setSelectionRange(dots, dots + 3);
    });
  }

  return (
    <section className="rounded-3xl border bg-white p-5 shadow-sm">
      {content.frame?.length ? (
        <div className="mb-3">
          <p className="mb-2 text-xs font-semibold text-muted-foreground">Mẫu câu gợi ý (chạm để chèn, rồi sửa phần "...")</p>
          <div className="flex flex-wrap gap-2">
            {content.frame.map((frame) => (
              <button
                className="cursor-pointer rounded-full border border-primary/30 bg-primary/5 px-3 py-1 text-sm font-medium text-primary hover:bg-primary/10"
                key={frame}
                onClick={() => insertFrame(frame)}
                type="button"
              >
                {frame}
              </button>
            ))}
          </div>
        </div>
      ) : null}
      <textarea
        className="min-h-44 w-full resize-y rounded-2xl border-2 border-input bg-white p-4 text-lg leading-relaxed outline-none focus:border-primary"
        onChange={(event) => setText(event.target.value)}
        placeholder="Viết bài của bạn ở đây…"
        ref={textareaRef}
        spellCheck={false}
        value={text}
      />
      <div className="mt-3 flex flex-wrap items-center justify-between gap-3">
        <span className={cn("text-sm font-semibold", sentences < min ? "text-muted-foreground" : sentences > max ? "text-amber-600" : "text-emerald-600")}>
          {sentences} / {min}–{max} câu
        </span>
        <Button disabled={sentences < min || submit.isPending} onClick={() => submit.mutate({ activityId, text })} size="lg">
          <Send className="size-4" /> {submit.isPending ? "AI đang chấm…" : "Nộp bài"}
        </Button>
      </div>
      {submit.error ? <p className="mt-2 text-sm text-red-600">{getApiErrorMessage(submit.error, "Không nộp được bài.")}</p> : null}
    </section>
  );
}

function Result({
  submission,
  finishing,
  onFinish,
  onUpdated,
  onWriteNew,
}: {
  submission: WritingSubmission;
  finishing: boolean;
  onFinish: () => void;
  onUpdated: (submission: WritingSubmission) => void;
  onWriteNew: () => void;
}) {
  const { report } = submission;
  const completed = submission.status === "COMPLETED";

  return (
    <>
      <section className="rounded-3xl border bg-white p-5 shadow-sm">
        <div className="flex flex-wrap items-center gap-5">
          <div className={cn("text-center", submission.passed ? "text-emerald-600" : "text-amber-600")}>
            <p className="text-5xl font-bold tracking-normal">{report.overall}</p>
            <p className="text-xs font-semibold">điểm · cần {submission.passScore}</p>
          </div>
          <ScoreBars className="min-w-56 flex-1" scores={report.scores} />
        </div>
        {report.feedbackVi ? (
          <p className="mt-4 flex gap-2 rounded-2xl bg-violet-50 p-3 text-sm">
            <Sparkles className="mt-0.5 size-4 shrink-0 text-violet-600" />
            {report.feedbackVi}
          </p>
        ) : null}

        <h3 className="mb-2 mt-5 font-bold">
          Bài của bạn {report.mistakes.length > 0 ? `· ${report.mistakes.length} lỗi (chạm vào lỗi để xem giải thích)` : "· không có lỗi"}
        </h3>
        <MistakeText mistakes={report.mistakes} text={submission.text} />

        {submission.modelAnswer ? (
          <details className="mt-4">
            <summary className="cursor-pointer text-sm font-semibold text-primary">Xem bài mẫu</summary>
            <p className="mt-2 rounded-2xl bg-muted/40 p-3 leading-relaxed">{submission.modelAnswer}</p>
          </details>
        ) : null}
      </section>

      {submission.status === "REWRITE_REQUIRED" ? (
        <Rewrite onUpdated={onUpdated} submission={submission} />
      ) : null}

      {completed ? (
        <section className="rounded-3xl border border-emerald-200 bg-emerald-50 p-5">
          <p className="flex items-center gap-2 font-bold text-emerald-700">
            <CheckCircle2 className="size-5" />
            {submission.rewriteCount > 0 ? "Bài sửa đã đúng hết!" : "Bài viết không có lỗi!"}
          </p>
          {!submission.passed ? (
            <p className="mt-1 text-sm text-amber-700">
              Điểm {report.overall} chưa đạt {submission.passScore}. Bạn có thể hoàn thành để lưu kết quả, hoặc viết một bài mới đầy đủ hơn.
            </p>
          ) : null}
          <div className="mt-4 flex flex-wrap gap-3">
            <Button disabled={finishing} onClick={onFinish} size="lg">
              Hoàn thành hoạt động
            </Button>
            <Button onClick={onWriteNew} size="lg" variant="outline">
              <RotateCcw className="size-4" /> Viết bài mới
            </Button>
          </div>
        </section>
      ) : null}
    </>
  );
}

function Rewrite({ submission, onUpdated }: { submission: WritingSubmission; onUpdated: (submission: WritingSubmission) => void }) {
  const [text, setText] = useState(submission.rewriteText ?? submission.text);
  const rewrite = useMutation({
    mutationFn: (value: string) => rewriteWriting(submission.id, value),
    onSuccess: onUpdated,
  });
  const remaining = submission.rewriteReport?.mistakes ?? [];

  return (
    <section className="rounded-3xl border border-amber-200 bg-amber-50/60 p-5">
      <h3 className="font-bold">Sửa lại bài</h3>
      <p className="mt-1 text-sm text-muted-foreground">
        Sửa các lỗi được đánh dấu rồi gửi lại. Hoạt động hoàn thành khi bài sửa không còn lỗi.
      </p>
      {remaining.length > 0 ? (
        <div className="mt-3 rounded-2xl bg-white p-3">
          <p className="mb-2 text-sm font-semibold text-red-600">Bài sửa lần {submission.rewriteCount} còn {remaining.length} lỗi:</p>
          <MistakeText mistakes={remaining} text={submission.rewriteText ?? ""} />
        </div>
      ) : null}
      <textarea
        className="mt-3 min-h-40 w-full resize-y rounded-2xl border-2 border-input bg-white p-4 text-lg leading-relaxed outline-none focus:border-primary"
        onChange={(event) => setText(event.target.value)}
        spellCheck={false}
        value={text}
      />
      <div className="mt-3 flex justify-end">
        <Button disabled={!text.trim() || rewrite.isPending} onClick={() => rewrite.mutate(text)} size="lg">
          <Send className="size-4" /> {rewrite.isPending ? "Đang kiểm tra…" : "Gửi bài sửa"}
        </Button>
      </div>
      {rewrite.error ? <p className="mt-2 text-sm text-red-600">{getApiErrorMessage(rewrite.error, "Không gửi được bài sửa.")}</p> : null}
    </section>
  );
}
