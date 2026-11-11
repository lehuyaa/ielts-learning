import { useMutation, useQuery } from "@tanstack/react-query";
import { ArrowLeft, ArrowRight, Sparkles } from "lucide-react";
import { useState } from "react";
import { Link, useLocation, useParams } from "react-router-dom";

import { getApiErrorMessage } from "@/api/api";
import { answerPractice, getPracticeSet } from "@/api/practice";
import { ErrorState } from "@/components/state/ErrorState";
import { Button } from "@/components/ui/button";
import { Skeleton } from "@/components/ui/skeleton";
import { ExerciseFeedback } from "@/features/exercise/ExerciseFeedback";
import { ExerciseRenderer } from "@/features/exercise/renderers";
import { cn } from "@/lib/utils";
import type { AnswerResult, ExerciseAnswer } from "@/types/activity";
import type { PracticeSet } from "@/types/practice";

/** Extra practice for one weakness (Task 74), played like a lesson activity. */
export function PracticePage() {
  const setId = Number(useParams().setId);
  const location = useLocation();
  const returnTo = (location.state as { returnTo?: string } | null)?.returnTo ?? "/mistakes";
  const query = useQuery({ queryKey: ["practice", setId], queryFn: () => getPracticeSet(setId), enabled: setId > 0, staleTime: Infinity });

  return (
    <div className="mx-auto max-w-3xl space-y-5 py-2">
      <Link className="inline-flex items-center gap-1 text-sm font-semibold text-muted-foreground hover:text-foreground" to={returnTo}>
        <ArrowLeft className="size-4" /> Quay lại
      </Link>
      {query.isLoading ? <Skeleton className="h-96 rounded-3xl" /> : null}
      {query.error ? (
        <ErrorState description={getApiErrorMessage(query.error, "Không mở được bài luyện.")} onRetry={() => void query.refetch()} title="Không mở được bài luyện" />
      ) : null}
      {query.data ? <PracticeRunner key={query.data.id} returnTo={returnTo} set={query.data} /> : null}
    </div>
  );
}

function PracticeRunner({ set, returnTo }: { set: PracticeSet; returnTo: string }) {
  const exercises = set.exercises;
  const [results, setResults] = useState(() => new Map(set.answers.map((answer) => [answer.exerciseId, answer])));
  const firstOpen = exercises.findIndex((exercise) => !results.has(exercise.id));
  const [index, setIndex] = useState(firstOpen === -1 ? exercises.length : firstOpen);
  const [draft, setDraft] = useState<ExerciseAnswer | null>(null);
  const answer = useMutation({
    mutationFn: (input: { exerciseId: number; answer: ExerciseAnswer }) => answerPractice(set.id, input),
    onSuccess: (result: AnswerResult) => setResults((current) => new Map(current).set(result.exerciseId, result)),
  });

  const exercise = exercises[index];
  const result = exercise ? results.get(exercise.id) : undefined;
  const correct = [...results.values()].filter((item) => item.isCorrect).length;

  function check() {
    if (!exercise || !draft || result || answer.isPending) return;
    answer.mutate({ exerciseId: exercise.id, answer: draft });
  }

  function next() {
    setDraft(null);
    setIndex(index + 1);
  }

  return (
    <div>
      <header className="mb-5">
        <p className="flex items-center gap-2 text-xs font-semibold uppercase text-muted-foreground">
          {set.source === "AI" ? (
            <>
              <Sparkles className="size-3.5 text-violet-600" /> AI tạo riêng cho bạn hôm nay
            </>
          ) : (
            "Chọn từ các bài học"
          )}
        </p>
        <h1 className="mt-1 text-2xl font-bold tracking-normal">Luyện: {set.label}</h1>
      </header>

      <div aria-hidden="true" className="mb-4 flex gap-1.5">
        {exercises.map((item, i) => {
          const itemResult = results.get(item.id);
          return (
            <span
              className={cn(
                "h-2 flex-1 rounded-full",
                itemResult?.isCorrect && "bg-emerald-500",
                itemResult && !itemResult.isCorrect && (itemResult.scoreRatio > 0 ? "bg-amber-400" : "bg-red-400"),
                !itemResult && (i === index ? "bg-primary" : "bg-muted"),
              )}
              key={item.id}
            />
          );
        })}
      </div>

      {!exercise ? (
        <div className="rounded-3xl border bg-white p-8 text-center shadow-sm">
          <p className="text-5xl font-bold tracking-normal">
            {correct}/{exercises.length}
          </p>
          <p className="mt-1 text-muted-foreground">câu đúng</p>
          <p className="mx-auto mt-3 max-w-md text-sm text-muted-foreground">
            Câu sai được ghi vào sổ lỗi. Ngày mai bạn có thể tạo một bài luyện mới cho điểm yếu này.
          </p>
          <Button asChild className="mt-6">
            <Link to={returnTo}>Xong</Link>
          </Button>
        </div>
      ) : (
        <>
          <section className="rounded-3xl border bg-white p-5 shadow-sm sm:p-7">
            <p className="mb-5 text-center text-lg font-semibold">{exercise.prompt}</p>
            <ExerciseRenderer exercise={exercise} key={exercise.id} onChange={setDraft} onSubmit={check} result={result} />
          </section>
          {result ? <ExerciseFeedback exercise={exercise} result={result} /> : null}
          {answer.error ? <p className="mt-3 text-center text-sm text-red-600">{getApiErrorMessage(answer.error, "Không gửi được câu trả lời.")}</p> : null}
          <div className="mt-6 flex items-center justify-between gap-3">
            <span className="text-sm text-muted-foreground">
              Câu {index + 1}/{exercises.length}
            </span>
            {result ? (
              <Button onClick={next}>
                {index >= exercises.length - 1 ? "Xem kết quả" : "Tiếp tục"} <ArrowRight className="size-4" />
              </Button>
            ) : (
              <Button disabled={!draft || answer.isPending} onClick={check}>
                Kiểm tra
              </Button>
            )}
          </div>
        </>
      )}
    </div>
  );
}
