import { AlertCircle } from "lucide-react";
import { useParams } from "react-router-dom";

import { APIError } from "@/api/api";
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
  const submitErrorMessage = getSubmitQuizErrorMessage(
    submitQuizMutation.error,
  );

  if (quizQuery.isLoading) {
    return (
      <QuizPageState
        title="Loading quiz"
        description="Preparing this lesson's questions."
      />
    );
  }

  if (errorMessage) {
    return (
      <QuizPageState
        title="Quiz unavailable"
        description={errorMessage}
        tone="error"
      />
    );
  }

  if (!quizQuery.data || quizQuery.data.questions.length === 0) {
    return (
      <QuizPageState
        title="No quiz yet"
        description="This lesson does not have quiz questions yet."
      />
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
      submitError={submitErrorMessage}
    />
  );
}

type QuizPageStateProps = {
  title: string;
  description: string;
  tone?: "default" | "error";
};

function QuizPageState({
  title,
  description,
  tone = "default",
}: QuizPageStateProps) {
  return (
    <div className="min-h-screen bg-[#f8f8ff] px-4 py-10 text-[#10111f] sm:px-6">
      <section className="mx-auto mt-16 max-w-[460px] rounded-[28px] border border-[#e6e6f3] bg-white px-7 py-8 text-center shadow-[0_24px_60px_rgba(26,27,45,0.12)] sm:px-8">
        {tone === "error" ? (
          <AlertCircle
            className="mx-auto size-9 text-destructive"
            aria-hidden="true"
          />
        ) : null}
        <h1 className="mt-4 text-2xl font-bold tracking-normal text-[#10111f]">
          {title}
        </h1>
        <p className="mt-3 text-base font-medium text-[#74768a]">
          {description}
        </p>
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

function getSubmitQuizErrorMessage(error: Error | null) {
  if (!error) {
    return null;
  }

  if (error instanceof APIError) {
    return error.message;
  }

  return "Unable to submit this quiz right now.";
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
