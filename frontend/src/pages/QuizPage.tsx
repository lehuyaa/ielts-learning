import { useParams } from "react-router-dom";

import { APIError } from "@/api/api";
import { EmptyState } from "@/components/state/EmptyState";
import { ErrorState } from "@/components/state/ErrorState";
import { CardSkeleton } from "@/components/state/CardSkeleton";
import { Skeleton } from "@/components/ui/skeleton";
import { useCheckQuizAnswer } from "@/features/quiz/hooks/useCheckQuizAnswer";
import { useLessonQuiz } from "@/features/quiz/hooks/useLessonQuiz";
import { useSubmitQuiz } from "@/features/quiz/hooks/useSubmitQuiz";
import { mapLessonQuizToSessionViewModel } from "@/features/quiz/mapQuiz";
import { QuizSession } from "@/features/quiz/QuizSession";

export function QuizPage() {
  const { lessonId } = useParams();
  const quizQuery = useLessonQuiz(lessonId);
  const checkAnswerMutation = useCheckQuizAnswer({ lessonId });
  const submitQuizMutation = useSubmitQuiz({ lessonId });
  const errorMessage = getQuizErrorMessage(quizQuery.error);

  if (quizQuery.isLoading) {
    return <QuizLoadingSkeleton />;
  }

  if (errorMessage) {
    return (
      <div className="min-h-screen bg-[#f8f8ff] px-4 py-10 text-[#10111f] sm:px-6">
        <ErrorState
          actionHref={lessonId ? `/lessons/${lessonId}` : "/roadmap"}
          actionLabel="Back to lesson"
          description={errorMessage}
          onRetry={() => {
            void quizQuery.refetch();
          }}
          title="Quiz unavailable"
          className="mx-auto mt-16 max-w-[460px]"
        />
      </div>
    );
  }

  if (!quizQuery.data || quizQuery.data.questions.length === 0) {
    return (
      <div className="min-h-screen bg-[#f8f8ff] px-4 py-10 text-[#10111f] sm:px-6">
        <EmptyState
          actionHref={lessonId ? `/lessons/${lessonId}` : "/roadmap"}
          actionLabel="Back to lesson"
          description="This lesson does not have quiz questions yet."
          title="No quiz yet"
          className="mx-auto mt-16 max-w-[460px]"
        />
      </div>
    );
  }

  const quiz = mapLessonQuizToSessionViewModel(quizQuery.data);

  return (
    <QuizSession
      checkError={getCheckAnswerErrorMessage(checkAnswerMutation.error)}
      isSubmitting={submitQuizMutation.isPending}
      isCheckingAnswer={checkAnswerMutation.isPending}
      lessonId={lessonId ?? String(quiz.lessonId)}
      onCheckAnswer={(payload) => checkAnswerMutation.mutateAsync(payload)}
      onSubmit={(payload) => submitQuizMutation.mutateAsync(payload)}
      quiz={quiz}
    />
  );
}

function QuizLoadingSkeleton() {
  return (
    <div className="min-h-screen bg-[#f8f8ff] px-4 py-10 text-[#10111f] sm:px-6">
      <section className="mx-auto mt-8 w-full max-w-[780px]">
        <div className="rounded-[28px] border border-[#e6e6f3] bg-white px-7 py-8 shadow-[0_24px_60px_rgba(26,27,45,0.12)] sm:px-8">
          <div className="flex items-center justify-between gap-4">
            <div className="flex flex-wrap items-center gap-3">
              <Skeleton className="h-8 w-28 rounded-full" />
              <Skeleton className="h-8 w-24 rounded-full" />
            </div>
            <Skeleton className="size-12 rounded-full" />
          </div>

          <div className="mt-10">
            <Skeleton className="h-4 w-24" />
            <Skeleton className="mt-4 h-8 w-full max-w-[34rem]" />
          </div>

          <div className="mt-10 grid gap-4">
            {Array.from({ length: 4 }).map((_, index) => (
              <CardSkeleton className="border-[#e6e6f3] p-5" key={index} lines={1} />
            ))}
          </div>

          <Skeleton className="mt-6 h-14 rounded-2xl" />
        </div>
      </section>
    </div>
  );
}

function getQuizErrorMessage(error: Error | null) {
  if (!error) {
    return null;
  }

  if (error instanceof APIError) {
    if (error.status === 404) {
      return "This quiz could not be found.";
    }

    if (error.status === 403) {
      return "Complete previous lessons to unlock this quiz.";
    }

    return error.message;
  }

  return "Unable to load this quiz right now.";
}

function getCheckAnswerErrorMessage(error: Error | null) {
  if (!error) {
    return null;
  }

  if (error instanceof APIError) {
    return error.message;
  }

  return "Unable to check this answer right now.";
}
