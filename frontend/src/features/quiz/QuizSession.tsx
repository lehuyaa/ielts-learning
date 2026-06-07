import { ArrowLeft, ArrowRight, Check, X } from "lucide-react";
import { useMemo, useState } from "react";
import { Link } from "react-router-dom";

import { Button } from "@/components/ui/button";
import type {
  CheckQuizAnswerRequest,
  CheckQuizAnswerResponse,
  QuizOption,
  QuizQuestion,
  QuizResultItem,
  QuizSessionViewModel,
  SubmitQuizRequest,
  SubmitQuizResponse,
} from "@/types/quiz";

type QuizSessionProps = {
  quiz: QuizSessionViewModel;
  lessonId: string;
  isSubmitting?: boolean;
  isCheckingAnswer?: boolean;
  submitError?: string | null;
  checkError?: string | null;
  onCheckAnswer: (
    payload: CheckQuizAnswerRequest,
  ) => Promise<CheckQuizAnswerResponse>;
  onSubmit: (payload: SubmitQuizRequest) => Promise<SubmitQuizResponse>;
};

type SelectedAnswers = Record<number, number>;
type CheckedAnswers = Record<number, CheckQuizAnswerResponse>;

export function QuizSession({
  quiz,
  lessonId,
  isSubmitting = false,
  isCheckingAnswer = false,
  submitError,
  checkError,
  onCheckAnswer,
  onSubmit,
}: QuizSessionProps) {
  const [currentQuestionIndex, setCurrentQuestionIndex] = useState(0);
  const [selectedAnswers, setSelectedAnswers] = useState<SelectedAnswers>({});
  const [checkedAnswers, setCheckedAnswers] = useState<CheckedAnswers>({});
  const [result, setResult] = useState<SubmitQuizResponse | null>(null);

  const currentQuestion = quiz.questions[currentQuestionIndex];
  const selectedOptionId = selectedAnswers[currentQuestion.id];
  const checkedAnswer = checkedAnswers[currentQuestion.id];
  const isFinalQuestion = currentQuestionIndex === quiz.questions.length - 1;
  const answeredPoints = useMemo(
    () => calculateEarnedCheckedPoints(checkedAnswers),
    [checkedAnswers],
  );

  async function selectOption(optionId: number) {
    if (result || isSubmitting || selectedAnswers[currentQuestion.id]) {
      return;
    }

    setSelectedAnswers((answers) => ({
      ...answers,
      [currentQuestion.id]: optionId,
    }));

    try {
      const response = await onCheckAnswer({
        questionId: currentQuestion.id,
        optionId,
      });
      setCheckedAnswers((answers) => ({
        ...answers,
        [currentQuestion.id]: response,
      }));
    } catch {
      // The mutation exposes a user-facing error through checkError.
    }
  }

  async function goNext() {
    if (!selectedOptionId) {
      return;
    }

    if (isFinalQuestion) {
      try {
        const response = await onSubmit({
          answers: Object.entries(selectedAnswers).map(
            ([questionId, optionId]) => ({
              questionId: Number(questionId),
              optionId,
            }),
          ),
        });
        setResult(response);
      } catch {
        // The mutation exposes a user-facing error through submitError.
      }
      return;
    }

    setCurrentQuestionIndex((index) => index + 1);
  }

  function retryQuiz() {
    setCurrentQuestionIndex(0);
    setSelectedAnswers({});
    setCheckedAnswers({});
    setResult(null);
  }

  if (result) {
    return (
      <QuizResultShell quiz={quiz}>
        <QuizResult
          lessonId={lessonId}
          onRetry={retryQuiz}
          quiz={quiz}
          result={result}
        />
      </QuizResultShell>
    );
  }

  return (
    <QuizShell
      currentQuestionIndex={currentQuestionIndex}
      earnedPoints={answeredPoints}
      quiz={quiz}
      checkedAnswers={checkedAnswers}
      selectedAnswers={selectedAnswers}
    >
      <section className="mx-auto w-full max-w-[780px] px-4 py-10 sm:px-6 lg:py-12">
        <div className="flex items-center justify-between gap-4">
          <div className="flex flex-wrap items-center gap-3">
            <span className="rounded-full bg-[#efefff] px-3 py-1.5 text-sm font-bold text-[#5147e8]">
              {quiz.topicTitle}
            </span>
            <span className="rounded-full bg-[#ececf4] px-3 py-1.5 text-sm font-bold text-[#74768a]">
              {quiz.band}
            </span>
          </div>
          <TimerBadge seconds={quiz.timeLimitSeconds} />
        </div>

        <div className="mt-10">
          <p className="text-sm font-bold uppercase tracking-normal text-[#676982]">
            Question {currentQuestionIndex + 1}
          </p>
          <h1 className="mt-4 text-2xl font-bold tracking-normal text-[#10111f] sm:text-3xl">
            {currentQuestion.prompt}
          </h1>
        </div>

        <div className="mt-10 grid gap-4">
          {currentQuestion.options.map((option) => (
            <QuizOptionCard
              checkedAnswer={checkedAnswer}
              isSelected={selectedOptionId === option.id}
              isLocked={Boolean(selectedOptionId)}
              key={option.id}
              onSelect={() => selectOption(option.id)}
              option={option}
            />
          ))}
        </div>

        {selectedOptionId ? (
          <QuestionAnsweredMessage
            checkedAnswer={checkedAnswer}
            errorMessage={checkError}
            isChecking={isCheckingAnswer}
          />
        ) : null}
        {submitError ? (
          <p className="mt-4 rounded-2xl border border-[#ffc9c5] bg-[#fff3f2] p-4 text-base font-semibold text-[#c8332d]">
            {submitError}
          </p>
        ) : null}

        <Button
          className="mt-6 h-14 w-full rounded-2xl bg-[#5147e8] text-base font-bold text-white hover:bg-[#453bd4]"
          disabled={!selectedOptionId || isSubmitting || isCheckingAnswer}
          onClick={goNext}
          type="button"
        >
          {isFinalQuestion
            ? isSubmitting
              ? "Submitting..."
              : "Submit Quiz"
            : "Next Question"}
          <ArrowRight className="size-6" aria-hidden="true" />
        </Button>
      </section>
    </QuizShell>
  );
}

type QuizShellProps = {
  children: React.ReactNode;
  currentQuestionIndex: number;
  earnedPoints: number;
  quiz: QuizSessionViewModel;
  checkedAnswers: CheckedAnswers;
  selectedAnswers: SelectedAnswers;
};

function QuizShell({
  children,
  currentQuestionIndex,
  earnedPoints,
  quiz,
  checkedAnswers,
  selectedAnswers,
}: QuizShellProps) {
  const progressPercentage =
    (Math.min(currentQuestionIndex, quiz.questions.length) /
      quiz.questions.length) *
    100;

  return (
    <div className="min-h-screen bg-[#f8f8ff] text-[#10111f]">
      <header className="border-b border-[#e6e6f3] bg-white">
        <div className="mx-auto flex max-w-[780px] items-center gap-5 px-4 py-4 sm:px-6">
          <Button
            asChild
            className="size-9 shrink-0 rounded-full text-[#676982]"
            size="icon"
            variant="ghost"
          >
            <Link aria-label="Back to lesson" to={`/lessons/${quiz.lessonId}`}>
              <ArrowLeft className="size-5" aria-hidden="true" />
            </Link>
          </Button>

          <div className="min-w-0 flex-1">
            <div className="flex items-center justify-between gap-4">
              <p className="text-sm font-bold text-[#676982]">
                Question{" "}
                {Math.min(currentQuestionIndex + 1, quiz.questions.length)} of{" "}
                {quiz.questions.length}
              </p>
              <p className="text-sm font-bold text-[#5147e8]">
                {earnedPoints} pts
              </p>
            </div>
            <div className="mt-2 h-2 overflow-hidden rounded-full bg-[#ececf4]">
              <div
                className="h-full rounded-full bg-[#5147e8] transition-all duration-300"
                style={{ width: `${progressPercentage}%` }}
              />
            </div>
          </div>

          <div
            className="hidden items-center gap-2 sm:flex"
            aria-label="Question result progress"
          >
            {quiz.questions.map((question, index) => (
              <span
                className={`size-2.5 rounded-full ${questionDotClass(question, selectedAnswers, checkedAnswers, index === currentQuestionIndex)}`}
                key={question.id}
              />
            ))}
          </div>
        </div>
      </header>

      {children}
    </div>
  );
}

type QuizResultShellProps = {
  children: React.ReactNode;
  quiz: QuizSessionViewModel;
};

function QuizResultShell({ children, quiz }: QuizResultShellProps) {
  return (
    <div className="min-h-screen bg-[#f8f8ff] text-[#10111f]">
      <header className="border-b border-[#e6e6f3] bg-white">
        <div className="mx-auto flex max-w-[780px] items-center gap-5 px-4 py-4 sm:px-6">
          <Button
            asChild
            className="size-9 shrink-0 rounded-full text-[#676982]"
            size="icon"
            variant="ghost"
          >
            <Link aria-label="Back to lesson" to={`/lessons/${quiz.lessonId}`}>
              <ArrowLeft className="size-5" aria-hidden="true" />
            </Link>
          </Button>
          <h1 className="text-lg font-bold tracking-normal text-[#10111f]">
            Quiz Results
          </h1>
        </div>
      </header>

      {children}
    </div>
  );
}

function questionDotClass(
  question: QuizQuestion,
  selectedAnswers: SelectedAnswers,
  checkedAnswers: CheckedAnswers,
  isCurrent: boolean,
) {
  const checkedAnswer = checkedAnswers[question.id];
  if (checkedAnswer) {
    return checkedAnswer.isCorrect ? "bg-[#49c389]" : "bg-[#ef6a67]";
  }

  const selectedOptionId = selectedAnswers[question.id];
  if (!selectedOptionId) {
    return isCurrent ? "bg-[#c9c8f8]" : "bg-[#e2e3ee]";
  }

  return "bg-[#5147e8]";
}

type QuizOptionCardProps = {
  option: QuizOption;
  checkedAnswer: CheckQuizAnswerResponse | undefined;
  isSelected: boolean;
  isLocked: boolean;
  onSelect: () => void;
};

function QuizOptionCard({
  option,
  checkedAnswer,
  isSelected,
  isLocked,
  onSelect,
}: QuizOptionCardProps) {
  const isCorrectOption = checkedAnswer?.correctOptionId === option.id;
  const isWrongSelection =
    checkedAnswer &&
    checkedAnswer.selectedOptionId === option.id &&
    !checkedAnswer.isCorrect;

  return (
    <button
      className={`flex min-h-[72px] w-full items-center gap-4 rounded-2xl border-2 px-5 py-3.5 text-left transition ${
        isCorrectOption
          ? "border-[#49c389] bg-[#effcf6] text-[#236b55]"
          : isWrongSelection
            ? "border-[#fb6a68] bg-[#fff3f2] text-[#a93633]"
            : isSelected
          ? "border-[#5147e8] bg-[#f0efff] text-[#28204d]"
          : "border-[#e6e6f3] bg-white text-[#2a2b38] hover:border-[#cfcef8]"
      } ${isLocked ? "cursor-not-allowed" : "cursor-pointer"}`}
      aria-disabled={isLocked}
      onClick={onSelect}
      type="button"
    >
      <span
        className={`grid size-11 shrink-0 place-items-center rounded-full text-base font-bold ${
          isCorrectOption
            ? "bg-[#49c389] text-white"
            : isWrongSelection
              ? "bg-[#ef4444] text-white"
              : isSelected
            ? "bg-[#5147e8] text-white"
            : "bg-[#eeeef8] text-[#74768a]"
        }`}
      >
        {isCorrectOption ? (
          <Check className="size-5" aria-hidden="true" />
        ) : isWrongSelection ? (
          <X className="size-5" aria-hidden="true" />
        ) : (
          option.label
        )}
      </span>
      <span className="text-base font-bold leading-6">{option.text}</span>
    </button>
  );
}

type QuestionAnsweredMessageProps = {
  checkedAnswer: CheckQuizAnswerResponse | undefined;
  errorMessage?: string | null;
  isChecking: boolean;
};

function QuestionAnsweredMessage({
  checkedAnswer,
  errorMessage,
  isChecking,
}: QuestionAnsweredMessageProps) {
  if (errorMessage) {
    return (
      <section className="mt-5 rounded-2xl border border-[#ffc9c5] bg-[#fff3f2] p-4">
        <div className="flex items-center gap-3 text-base font-bold text-[#c8332d]">
          <X className="size-5" aria-hidden="true" />
          Could not check answer
        </div>
        <p className="mt-2 text-base font-medium text-[#676982]">
          {errorMessage}
        </p>
      </section>
    );
  }

  if (isChecking || !checkedAnswer) {
    return (
      <section className="mt-5 rounded-2xl border border-[#e6e6f3] bg-white p-4">
        <div className="flex items-center gap-3 text-base font-bold text-[#5147e8]">
          <Check className="size-5" aria-hidden="true" />
          Checking answer...
        </div>
      </section>
    );
  }

  return (
    <section
      className={`mt-5 rounded-2xl border p-4 ${
        checkedAnswer.isCorrect
          ? "border-[#a8efd1] bg-[#effcf6]"
          : "border-[#ffc9c5] bg-[#fff3f2]"
      }`}
    >
      <div
        className={`flex items-center gap-3 text-base font-bold ${
          checkedAnswer.isCorrect ? "text-[#168653]" : "text-[#c8332d]"
        }`}
      >
        {checkedAnswer.isCorrect ? (
          <Check className="size-5" aria-hidden="true" />
        ) : (
          <X className="size-5" aria-hidden="true" />
        )}
        {checkedAnswer.isCorrect
          ? `Correct! +${checkedAnswer.earnedPoints} points`
          : "Incorrect - review the correct answer above"}
      </div>
      <p className="mt-2 text-base font-medium text-[#676982]">
        {checkedAnswer.explanation}
      </p>
    </section>
  );
}

type TimerBadgeProps = {
  seconds: number;
};

function TimerBadge({ seconds }: TimerBadgeProps) {
  return (
    <div className="grid size-12 place-items-center rounded-full border-4 border-[#5cc493] bg-white text-sm font-bold text-[#49a875]">
      <span>{seconds}</span>
    </div>
  );
}

type ScoreRingProps = {
  scorePercentage: number;
};

function ScoreRing({ scorePercentage }: ScoreRingProps) {
  return (
    <div
      className="mx-auto mt-7 grid size-36 place-items-center rounded-full"
      style={{
        background: `conic-gradient(#5147e8 ${scorePercentage}%, #ece8ff 0)`,
      }}
    >
      <div className="grid size-28 place-items-center rounded-full bg-white">
        <div>
          <p className="text-3xl font-bold tracking-normal text-[#5147e8]">
            {scorePercentage}%
          </p>
          <p className="mt-1 text-base font-medium text-[#74768a]">Score</p>
        </div>
      </div>
    </div>
  );
}

type ResultOutcomeDotProps = {
  isCorrect: boolean;
};

function ResultOutcomeDot({ isCorrect }: ResultOutcomeDotProps) {
  return (
    <span
      className={`grid size-9 place-items-center rounded-full text-white ${
        isCorrect ? "bg-[#49c389]" : "bg-[#ef6a67]"
      }`}
    >
      {isCorrect ? (
        <Check className="size-5" aria-hidden="true" />
      ) : (
        <X className="size-5" aria-hidden="true" />
      )}
    </span>
  );
}

type QuizResultProps = {
  result: SubmitQuizResponse;
  quiz: QuizSessionViewModel;
  lessonId: string;
  onRetry: () => void;
};

function QuizResult({
  result,
  quiz,
  lessonId,
  onRetry,
}: QuizResultProps) {
  const incorrectCount = result.totalQuestions - result.correctCount;
  const resultByQuestion = mapResultByQuestion(result.results);

  return (
    <main className="mx-auto grid min-h-[calc(100vh-80px)] w-full max-w-[460px] place-items-center px-4 py-10 sm:px-6">
      <section className="w-full rounded-[28px] border border-[#e6e6f3] bg-white px-7 py-8 text-center shadow-[0_24px_60px_rgba(26,27,45,0.12)] sm:px-8">
        <div className="text-5xl leading-none" aria-hidden="true">
          ⭐
        </div>
        <h1 className="mt-6 text-2xl font-bold tracking-normal text-[#10111f]">
          Quiz Complete!
        </h1>
        <p className="mt-2 text-base font-medium text-[#74768a]">
          Here's how you did
        </p>
        <p className="mt-2 text-sm font-bold text-[#74768a]">
          Required score: {result.requiredScore}% -{" "}
          {result.passed ? "Passed" : "Needs review"}
        </p>

        <ScoreRing scorePercentage={result.score} />

        <div className="mt-8 grid grid-cols-2 gap-4">
          <ResultStat
            label="Correct"
            tone="success"
            value={String(result.correctCount)}
          />
          <ResultStat
            label="Incorrect"
            tone="error"
            value={String(incorrectCount)}
          />
        </div>

        <div className="mt-7 flex flex-wrap justify-center gap-3">
          {quiz.questions.map((question) => (
            <ResultOutcomeDot
              isCorrect={resultByQuestion[question.id]?.isCorrect ?? false}
              key={question.id}
            />
          ))}
        </div>

        <div className="mt-8 grid gap-4 sm:grid-cols-2">
          <Button
            className="h-14 rounded-3xl border-[#e6e6f3] text-base font-bold text-[#2a2b38]"
            onClick={onRetry}
            type="button"
            variant="outline"
          >
            Try Again
          </Button>
          <Button
            asChild
            style={{ color: "white" }}
            className="h-14 rounded-3xl bg-[#5147e8] text-base font-bold text-white hover:bg-[#453bd4]"
          >
            <Link to={`/lessons/${lessonId}`}>Done</Link>
          </Button>
        </div>
      </section>
    </main>
  );
}

type ResultStatProps = {
  label: string;
  value: string;
  tone: "success" | "error";
};

function ResultStat({ label, value, tone }: ResultStatProps) {
  return (
    <div
      className={`rounded-2xl border p-5 ${
        tone === "success"
          ? "border-[#c7f5df] bg-[#effcf6] text-[#168653]"
          : "border-[#ffd5d2] bg-[#fff3f2] text-[#ef4444]"
      }`}
    >
      <p className="text-2xl font-bold tracking-normal">{value}</p>
      <p className="mt-2 text-sm font-medium">{label}</p>
    </div>
  );
}

function calculateEarnedCheckedPoints(checkedAnswers: CheckedAnswers) {
  return Object.values(checkedAnswers).reduce(
    (total, answer) => total + answer.earnedPoints,
    0,
  );
}

function mapResultByQuestion(results: QuizResultItem[]) {
  return results.reduce<Record<number, QuizResultItem>>((lookup, result) => {
    lookup[result.questionId] = result;
    return lookup;
  }, {});
}
