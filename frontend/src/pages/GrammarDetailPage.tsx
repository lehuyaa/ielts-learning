import { ArrowLeft, Dumbbell } from "lucide-react";
import { useState } from "react";
import { Link, useParams } from "react-router-dom";

import { APIError } from "@/api/api";
import { ErrorState } from "@/components/state/ErrorState";
import { Button } from "@/components/ui/button";
import { Skeleton } from "@/components/ui/skeleton";
import { GrammarPractice } from "@/features/grammar/GrammarPractice";
import { GrammarTheoryCard } from "@/features/grammar/GrammarTheoryCard";
import { MasteryBar } from "@/features/grammar/MasteryBar";
import { useGrammarPoint, useStartGrammarPractice } from "@/features/grammar/hooks/useGrammar";
import type { PracticeExercise } from "@/types/grammar";

export function GrammarDetailPage() {
  const { slug = "" } = useParams();
  const query = useGrammarPoint(slug);
  const startPractice = useStartGrammarPractice(slug);
  const [practice, setPractice] = useState<{ round: number; exercises: PracticeExercise[] } | null>(null);

  function start() {
    startPractice.mutate(undefined, {
      onSuccess: (data) =>
        setPractice((current) => ({ round: (current?.round ?? 0) + 1, exercises: data.exercises })),
    });
  }

  return (
    <div className="mx-auto max-w-3xl space-y-6 py-2">
      <Link className="inline-flex items-center gap-1 text-sm font-semibold text-muted-foreground hover:text-foreground" to="/grammar">
        <ArrowLeft className="size-4" /> Ngữ pháp
      </Link>

      {query.isLoading ? <Skeleton className="h-96 rounded-3xl" /> : null}
      {query.error ? (
        <ErrorState
          description={query.error instanceof APIError ? query.error.message : "Không tải được điểm ngữ pháp."}
          onRetry={() => void query.refetch()}
          title="Không tải được ngữ pháp"
        />
      ) : null}

      {query.data ? (
        <>
          <section className="flex flex-wrap items-center gap-4 rounded-3xl border bg-white p-5 shadow-sm">
            <MasteryBar
              className="min-w-48 flex-1"
              mastery={query.data.progress.mastery}
              practised={query.data.progress.practised}
            />
            <div className="text-xs text-muted-foreground">
              {query.data.progress.correctCount} đúng · {query.data.progress.wrongCount} sai
            </div>
            {!practice ? (
              <Button disabled={query.data.exerciseCount === 0 || startPractice.isPending} onClick={start}>
                <Dumbbell /> Luyện tập ({Math.min(query.data.exerciseCount, 10)} câu)
              </Button>
            ) : null}
            {query.data.exerciseCount === 0 ? (
              <p className="w-full text-xs text-muted-foreground">
                Học các bài có điểm ngữ pháp này trong Roadmap để mở phần luyện tập.
              </p>
            ) : null}
          </section>

          {practice ? (
            practice.exercises.length > 0 ? (
              <GrammarPractice
                exercises={practice.exercises}
                key={practice.round}
                onFinish={() => setPractice(null)}
                onRestart={start}
                slug={slug}
              />
            ) : (
              <p className="rounded-2xl border border-dashed p-6 text-center text-muted-foreground">
                Chưa có câu luyện tập nào.
              </p>
            )
          ) : (
            <GrammarTheoryCard point={query.data.grammarPoint} />
          )}
        </>
      ) : null}
    </div>
  );
}
