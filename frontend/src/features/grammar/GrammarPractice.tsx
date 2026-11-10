import { ArrowRight, CheckCircle2, RotateCcw, XCircle } from "lucide-react";
import { useState } from "react";

import { APIError } from "@/api/api";
import { Button } from "@/components/ui/button";
import { formatCorrectAnswer } from "@/features/exercise/payload";
import { ExerciseRenderer } from "@/features/exercise/renderers";
import { cn } from "@/lib/utils";
import type { AnswerResult, ExerciseAnswer } from "@/types/activity";
import type { GrammarProgress, PracticeExercise } from "@/types/grammar";

import { useAnswerGrammarPractice } from "./hooks/useGrammar";

type GrammarPracticeProps = {
  slug: string;
  exercises: PracticeExercise[];
  onFinish: () => void;
  onRestart: () => void;
};

/** A short practice set for one grammar point (answers update mastery). */
export function GrammarPractice({ slug, exercises, onFinish, onRestart }: GrammarPracticeProps) {
  const [index, setIndex] = useState(0);
  const [draft, setDraft] = useState<ExerciseAnswer | null>(null);
  const [results, setResults] = useState<Map<number, AnswerResult>>(new Map());
  const [progress, setProgress] = useState<GrammarProgress | null>(null);
  const [done, setDone] = useState(false);
  const answer = useAnswerGrammarPractice(slug);

  const exercise = exercises[index];
  const result = exercise ? results.get(exercise.id) : undefined;
  const correct = [...results.values()].filter((item) => item.isCorrect).length;

  function check() {
    if (!exercise || !draft || result || answer.isPending) return;
    answer.mutate(
      { exerciseId: exercise.id, answer: draft },
      {
        onSuccess: (response) => {
          setResults((current) =>
            new Map(current).set(response.exerciseId, response as unknown as AnswerResult),
          );
          setProgress(response.progress);
        },
      },
    );
  }

  function next() {
    setDraft(null);
    if (index >= exercises.length - 1) {
      setDone(true);
      return;
    }
    setIndex(index + 1);
  }

  if (done) {
    return (
      <div className="rounded-3xl border bg-white p-8 text-center shadow-sm">
        <p className="text-5xl font-bold tracking-normal">
          {correct}/{exercises.length}
        </p>
        <p className="mt-1 text-muted-foreground">câu đúng</p>
        {progress ? (
          <p className="mt-4 font-semibold">
            Mức thành thạo hiện tại: <span className="text-primary">{progress.mastery}%</span>
          </p>
        ) : null}
        <div className="mt-6 flex flex-wrap justify-center gap-3">
          <Button onClick={onRestart} variant="outline">
            <RotateCcw /> Luyện tiếp
          </Button>
          <Button onClick={onFinish}>Xong</Button>
        </div>
      </div>
    );
  }

  if (!exercise) return null;

  return (
    <div>
      <div className="mb-4 flex gap-1.5" aria-hidden="true">
        {exercises.map((item, i) => {
          const itemResult = results.get(item.id);
          return (
            <span
              className={cn(
                "h-2 flex-1 rounded-full",
                itemResult?.isCorrect && "bg-emerald-500",
                itemResult && !itemResult.isCorrect && "bg-red-400",
                !itemResult && (i === index ? "bg-primary" : "bg-muted"),
              )}
              key={item.id}
            />
          );
        })}
      </div>

      <section className="rounded-3xl border bg-white p-5 shadow-sm sm:p-7">
        {exercise.lastResult === "WRONG" ? (
          <p className="mb-2 text-center text-xs font-semibold text-red-600">Lần trước bạn làm sai câu này</p>
        ) : null}
        <p className="mb-5 text-center text-lg font-semibold">{exercise.prompt}</p>
        <ExerciseRenderer exercise={exercise} key={exercise.id} onChange={setDraft} onSubmit={check} result={result} />
      </section>

      {result ? (
        <div
          className={cn(
            "mt-4 rounded-2xl border p-4",
            result.isCorrect ? "border-emerald-200 bg-emerald-50" : "border-red-200 bg-red-50",
          )}
          role="status"
        >
          <p className="flex items-center gap-2 font-bold">
            {result.isCorrect ? <CheckCircle2 className="size-5 text-emerald-600" /> : <XCircle className="size-5 text-red-600" />}
            {result.isCorrect ? "Chính xác!" : "Chưa đúng"}
          </p>
          {!result.isCorrect ? (
            <p className="mt-1">
              Đáp án: <span className="font-semibold">{formatCorrectAnswer(exercise, result.correctAnswer)}</span>
            </p>
          ) : null}
          {result.explanationVi ? <p className="mt-1 text-sm text-muted-foreground">{result.explanationVi}</p> : null}
        </div>
      ) : null}

      {answer.error ? (
        <p className="mt-3 text-center text-sm text-red-600">
          {answer.error instanceof APIError ? answer.error.message : "Không gửi được câu trả lời."}
        </p>
      ) : null}

      <div className="mt-6 flex items-center justify-between gap-3">
        <span className="text-sm text-muted-foreground">
          Câu {index + 1}/{exercises.length}
        </span>
        {result ? (
          <Button onClick={next} size="lg">
            {index >= exercises.length - 1 ? "Xem kết quả" : "Tiếp tục"} <ArrowRight />
          </Button>
        ) : (
          <Button disabled={!draft || answer.isPending} onClick={check} size="lg">
            Kiểm tra
          </Button>
        )}
      </div>
    </div>
  );
}
