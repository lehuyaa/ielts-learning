import {
  ArrowRight,
  PartyPopper,
  RotateCcw,
  Sparkles,
} from "lucide-react";
import { useMemo, useState } from "react";
import { Link } from "react-router-dom";

import { APIError } from "@/api/api";
import { ErrorState } from "@/components/state/ErrorState";
import { Skeleton } from "@/components/ui/skeleton";
import { Button } from "@/components/ui/button";
import { ExerciseFeedback } from "@/features/exercise/ExerciseFeedback";
import { ExerciseRenderer } from "@/features/exercise/renderers";
import { cn } from "@/lib/utils";
import type {
  ActivityDetailResponse,
  AnswerResult,
  CompleteActivityResponse,
  ExerciseAnswer,
} from "@/types/activity";

import { GrammarTheoryCard } from "@/features/grammar/GrammarTheoryCard";

import { ActivityContent } from "./ActivityContent";
import { VocabularyCards } from "./VocabularyCards";
import { useActivity } from "./hooks/useActivity";
import { useCompleteActivity } from "./hooks/useCompleteActivity";
import { useSubmitAnswer } from "./hooks/useSubmitAnswer";

type ActivityRunnerProps = {
  activityId: number;
  lessonId: number;
  onNextActivity: (activityId: number) => void;
};

export function ActivityRunner({ activityId, lessonId, onNextActivity }: ActivityRunnerProps) {
  const activityQuery = useActivity(activityId);
  // Kept here (not in the session) so it survives the refetch after finishing.
  const [completion, setCompletion] = useState<CompleteActivityResponse | null>(null);
  // Bumped on retry so the session restarts with clean state.
  const [sessionNumber, setSessionNumber] = useState(0);

  if (completion) {
    return (
      <ActivityResult
        completion={completion}
        onNextActivity={onNextActivity}
        onRetry={() => {
          setCompletion(null);
          setSessionNumber((value) => value + 1);
          void activityQuery.refetch();
        }}
      />
    );
  }

  if (activityQuery.isLoading) {
    return (
      <div className="space-y-4">
        <Skeleton className="h-8 w-2/3" />
        <Skeleton className="h-48 w-full rounded-2xl" />
      </div>
    );
  }

  if (activityQuery.error || !activityQuery.data) {
    return (
      <ErrorState
        description={
          activityQuery.error instanceof APIError
            ? activityQuery.error.message
            : "Không tải được hoạt động này."
        }
        onRetry={() => void activityQuery.refetch()}
        title="Không mở được hoạt động"
      />
    );
  }

  return (
    <ActivitySession
      data={activityQuery.data}
      key={`${activityId}-${sessionNumber}`}
      lessonId={lessonId}
      onCompleted={setCompletion}
    />
  );
}

type Phase = "theory" | "cards" | "exercises";

function ActivitySession({
  data,
  lessonId,
  onCompleted,
}: {
  data: ActivityDetailResponse;
  lessonId: number;
  onCompleted: (completion: CompleteActivityResponse) => void;
}) {
  const { activity, exercises, vocabularies, grammarPoint, phonemes } = data;
  const showTheory = Boolean(grammarPoint) && activity.content?.showTheory !== false;
  const answersById = useMemo(
    () => new Map(data.answers.map((answer) => [answer.exerciseId, answer])),
    [data.answers],
  );

  // Resume at the first unanswered exercise of an unfinished attempt.
  const firstUnanswered = exercises.findIndex((exercise) => !answersById.has(exercise.id));
  const resuming = data.answers.length > 0;
  const [phase, setPhase] = useState<Phase>(() => {
    if (resuming) return "exercises";
    // Theory is shown before the first attempt; later it is one click away.
    if (showTheory && activity.attempts === 0) return "theory";
    return vocabularies.length > 0 ? "cards" : "exercises";
  });
  const [index, setIndex] = useState(firstUnanswered === -1 ? Math.max(exercises.length - 1, 0) : firstUnanswered);
  const [draft, setDraft] = useState<ExerciseAnswer | null>(null);
  const [results, setResults] = useState<Map<number, AnswerResult>>(answersById);

  const submitAnswer = useSubmitAnswer(activity.id);
  const complete = useCompleteActivity(activity.id, lessonId);

  const exercise = exercises[index];
  const result = exercise ? results.get(exercise.id) : undefined;
  const isLast = index >= exercises.length - 1;
  const answeredCount = results.size;

  function check() {
    if (!exercise || !draft || submitAnswer.isPending || result) {
      return;
    }
    submitAnswer.mutate(
      { exerciseId: exercise.id, answer: draft },
      {
        onSuccess: (answer) => {
          setResults((current) => new Map(current).set(answer.exerciseId, answer));
        },
      },
    );
  }

  function next() {
    setDraft(null);
    if (!isLast) {
      setIndex(index + 1);
      return;
    }
    finish();
  }

  function finish() {
    complete.mutate(undefined, { onSuccess: onCompleted });
  }

  return (
    <div>
      <header className="mb-6">
        <h2 className="text-2xl font-bold tracking-normal">{activity.title}</h2>
        {activity.instructionsVi ? (
          <p className="mt-1 text-muted-foreground">{activity.instructionsVi}</p>
        ) : null}
      </header>

      {phase === "theory" && grammarPoint ? (
        <div>
          <GrammarTheoryCard point={grammarPoint} />
          <div className="mt-6 flex justify-center">
            <Button onClick={() => setPhase("exercises")} size="lg">
              Đã hiểu, làm bài tập <ArrowRight />
            </Button>
          </div>
        </div>
      ) : null}

      {phase === "exercises" && grammarPoint ? (
        <button
          className="mb-4 cursor-pointer text-sm font-semibold text-primary hover:underline"
          onClick={() => setPhase("theory")}
          type="button"
        >
          Xem lại lý thuyết: {grammarPoint.titleVi || grammarPoint.title}
        </button>
      ) : null}

      {phase === "cards" ? (
        <VocabularyCards
          cards={vocabularies}
          doneLabel={exercises.length > 0 ? "Làm bài kiểm tra" : "Hoàn thành"}
          onDone={() => (exercises.length > 0 ? setPhase("exercises") : finish())}
        />
      ) : null}

      {phase === "exercises" && exercises.length === 0 ? (
        <div className="rounded-2xl border border-dashed p-8 text-center">
          <p className="mb-4 text-muted-foreground">Hoạt động này chưa có bài tập.</p>
          <Button disabled={complete.isPending} onClick={finish}>
            Đánh dấu hoàn thành
          </Button>
        </div>
      ) : null}

      {phase === "exercises" && exercise ? (
        <div>
          <ActivityContent activity={activity} phonemes={phonemes} />
          <ProgressDots count={exercises.length} current={index} results={results} exercises={exercises} />

          <section className="rounded-3xl border bg-white p-5 shadow-sm sm:p-7">
            <p className="mb-5 text-center text-lg font-semibold">{exercise.prompt}</p>
            <ExerciseRenderer
              exercise={exercise}
              key={exercise.id}
              onChange={setDraft}
              onSubmit={check}
              result={result}
            />
          </section>

          {result ? <ExerciseFeedback exercise={exercise} result={result} /> : null}

          {submitAnswer.error ? (
            <p className="mt-3 text-center text-sm text-red-600">
              {submitAnswer.error instanceof APIError ? submitAnswer.error.message : "Không gửi được câu trả lời."}
            </p>
          ) : null}
          {complete.error ? (
            <p className="mt-3 text-center text-sm text-red-600">
              {complete.error instanceof APIError ? complete.error.message : "Không hoàn thành được hoạt động."}
            </p>
          ) : null}

          <div className="mt-6 flex items-center justify-between gap-3">
            <span className="text-sm text-muted-foreground">
              Câu {index + 1}/{exercises.length} · đã trả lời {answeredCount}
            </span>
            {result ? (
              <Button disabled={complete.isPending} onClick={next} size="lg">
                {isLast ? "Xem kết quả" : "Tiếp tục"}
                <ArrowRight />
              </Button>
            ) : (
              <Button disabled={!draft || submitAnswer.isPending} onClick={check} size="lg">
                {submitAnswer.isPending ? "Đang chấm…" : "Kiểm tra"}
              </Button>
            )}
          </div>
        </div>
      ) : null}
    </div>
  );
}

function ProgressDots({
  count,
  current,
  results,
  exercises,
}: {
  count: number;
  current: number;
  results: Map<number, AnswerResult>;
  exercises: ActivityDetailResponse["exercises"];
}) {
  return (
    <div className="mb-4 flex gap-1.5" aria-hidden="true">
      {Array.from({ length: count }).map((_, i) => {
        const answer = results.get(exercises[i].id);
        return (
          <span
            className={cn(
              "h-2 flex-1 rounded-full",
              answer?.isCorrect && "bg-emerald-500",
              answer && !answer.isCorrect && (answer.scoreRatio > 0 ? "bg-amber-400" : "bg-red-400"),
              !answer && (i === current ? "bg-primary" : "bg-muted"),
            )}
            key={i}
          />
        );
      })}
    </div>
  );
}

function ActivityResult({
  completion,
  onNextActivity,
  onRetry,
}: {
  completion: CompleteActivityResponse;
  onNextActivity: (activityId: number) => void;
  onRetry: () => void;
}) {
  const { lesson } = completion;

  if (lesson.justCompleted) {
    return (
      <div className="rounded-3xl border bg-white p-8 text-center shadow-sm">
        <PartyPopper className="mx-auto mb-3 size-12 text-amber-500" />
        <h2 className="text-2xl font-bold tracking-normal">Hoàn thành bài học!</h2>
        <p className="mt-2 text-muted-foreground">
          Điểm bài học: <span className="font-semibold text-foreground">{lesson.score ?? completion.score}</span>
          {" · "}+{completion.xpAwarded + lesson.xpAwarded} XP
        </p>
        <div className="mt-6 flex flex-wrap justify-center gap-3">
          {completion.nextActivityId ? (
            <Button onClick={() => onNextActivity(completion.nextActivityId as number)} variant="outline">
              Làm hoạt động tùy chọn
            </Button>
          ) : null}
          <Link
            className="inline-flex h-11 items-center gap-2 rounded-xl bg-primary px-4 text-sm font-semibold text-primary-foreground hover:bg-primary/90"
            to="/roadmap"
          >
            Về lộ trình <ArrowRight className="size-4" />
          </Link>
        </div>
      </div>
    );
  }

  return (
    <div className="rounded-3xl border bg-white p-8 text-center shadow-sm">
      {completion.passed ? (
        <Sparkles className="mx-auto mb-3 size-12 text-emerald-500" />
      ) : (
        <RotateCcw className="mx-auto mb-3 size-12 text-amber-500" />
      )}
      <p className="text-5xl font-bold tracking-normal">{completion.score}</p>
      <p className="mt-1 text-sm text-muted-foreground">
        Đúng {completion.correctCount}/{completion.totalCount} câu · cần {completion.activity.passScore} điểm để qua
      </p>
      <h2 className="mt-4 text-xl font-bold tracking-normal">
        {completion.passed ? "Đạt rồi!" : "Chưa đạt, thử lại nhé"}
      </h2>
      {completion.xpAwarded > 0 ? <p className="mt-1 font-semibold text-primary">+{completion.xpAwarded} XP</p> : null}
      <p className="mt-1 text-sm text-muted-foreground">
        Bài học: {lesson.passedRequired}/{lesson.totalRequired} hoạt động bắt buộc đã đạt
      </p>
      <div className="mt-6 flex flex-wrap justify-center gap-3">
        <Button onClick={onRetry} variant="outline">
          <RotateCcw /> Làm lại
        </Button>
        {completion.nextActivityId && completion.nextActivityId !== completion.activity.id ? (
          <Button onClick={() => onNextActivity(completion.nextActivityId as number)}>
            Hoạt động tiếp theo <ArrowRight />
          </Button>
        ) : null}
      </div>
    </div>
  );
}
