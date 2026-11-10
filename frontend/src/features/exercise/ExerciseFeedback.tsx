import { CheckCircle2, Sparkles, XCircle } from "lucide-react";

import { cn } from "@/lib/utils";
import type { AnswerResult, Exercise } from "@/types/activity";

import { formatCorrectAnswer } from "./payload";

const aiStatusNote: Record<NonNullable<AnswerResult["aiStatus"]>, string> = {
  LIMIT_REACHED: "Hôm nay đã hết lượt AI chấm, câu này được so với đáp án mẫu.",
  UNAVAILABLE: "AI chấm tạm thời không dùng được, câu này được so với đáp án mẫu.",
};

/** Result panel shown under an answered exercise (lesson player and practice). */
export function ExerciseFeedback({ exercise, result }: { exercise: Exercise; result: AnswerResult }) {
  const partial = !result.isCorrect && result.scoreRatio > 0;
  const correctText = formatCorrectAnswer(exercise, result.correctAnswer);
  const ai = result.gradedBy === "AI" ? result.aiFeedback : undefined;
  const showSuggested = ai && !result.isCorrect && ai.suggested;
  // For AI-accepted answers the reference is still useful as another way to say it.
  const showReference = correctText && (!result.isCorrect || (ai && !sameText(ai.suggested, correctText)));

  return (
    <div
      className={cn(
        "mt-4 rounded-2xl border p-4",
        result.isCorrect ? "border-emerald-200 bg-emerald-50" : partial ? "border-amber-200 bg-amber-50" : "border-red-200 bg-red-50",
      )}
      role="status"
    >
      <p className="flex flex-wrap items-center gap-2 font-bold">
        {result.isCorrect ? (
          <CheckCircle2 className="size-5 text-emerald-600" />
        ) : (
          <XCircle className={cn("size-5", partial ? "text-amber-600" : "text-red-600")} />
        )}
        {result.isCorrect ? "Chính xác!" : partial ? `Gần đúng (${Math.round(result.scoreRatio * 100)}%)` : "Chưa đúng"}
        {ai ? (
          <span className="inline-flex items-center gap-1 rounded-full bg-white/80 px-2 py-0.5 text-xs font-semibold text-violet-700">
            <Sparkles className="size-3" />
            AI chấm
          </span>
        ) : null}
      </p>

      {ai?.explanationVi ? <p className="mt-1">{ai.explanationVi}</p> : null}
      {showSuggested ? (
        <p className="mt-1">
          Sửa lại: <span className="font-semibold">{ai.suggested}</span>
        </p>
      ) : null}
      {showReference ? (
        <p className="mt-1">
          {result.isCorrect ? "Cách nói khác: " : ai ? "Đáp án mẫu: " : "Đáp án: "}
          <span className="font-semibold">{correctText}</span>
        </p>
      ) : null}
      {result.explanationVi ? <p className="mt-1 text-sm text-muted-foreground">{result.explanationVi}</p> : null}
      {result.aiStatus && !result.isCorrect ? (
        <p className="mt-2 text-xs text-muted-foreground">{aiStatusNote[result.aiStatus]}</p>
      ) : null}
    </div>
  );
}

function sameText(a: string, b: string) {
  const clean = (value: string) => value.toLowerCase().replace(/[.!?\s]+$/g, "").replace(/\s+/g, " ").trim();
  return clean(a) === clean(b);
}
