import { ArrowLeft, ArrowRight, Check, X } from "lucide-react";
import { useMemo, useState } from "react";
import { Link } from "react-router-dom";

import { Button } from "@/components/ui/button";

import type { MockQuiz, MockQuizOption, MockQuizQuestion } from "./mockQuiz";

type QuizSessionProps = {
  quiz: MockQuiz;
  lessonId: string;
};

type SelectedAnswers = Record<string, string>;

export function QuizSession({ quiz, lessonId }: QuizSessionProps) {
  const [currentQuestionIndex, setCurrentQuestionIndex] = useState(0);
  const [selectedAnswers, setSelectedAnswers] = useState<SelectedAnswers>({});
  const [isSubmitted, setIsSubmitted] = useState(false);

  const currentQuestion = quiz.questions[currentQuestionIndex];
  const selectedOptionId = selectedAnswers[currentQuestion.id];
  const isFinalQuestion = currentQuestionIndex === quiz.questions.length - 1;
  const earnedPoints = useMemo(
    () => calculateEarnedPoints(quiz.questions, selectedAnswers),
    [quiz.questions, selectedAnswers],
  );
  const result = useMemo(
    () => calculateResult(quiz, selectedAnswers),
    [quiz, selectedAnswers],
  );

  function selectOption(optionId: string) {
    if (isSubmitted || selectedAnswers[currentQuestion.id]) {
      return;
    }

    setSelectedAnswers((answers) => ({
      ...answers,
      [currentQuestion.id]: optionId,
    }));
  }

  function goNext() {
    if (!selectedOptionId) {
      return;
    }

    if (isFinalQuestion) {
      setIsSubmitted(true);
      return;
    }

    setCurrentQuestionIndex((index) => index + 1);
  }

  function retryQuiz() {
    setCurrentQuestionIndex(0);
    setSelectedAnswers({});
    setIsSubmitted(false);
  }

  if (isSubmitted) {
    return (
      <QuizResultShell quiz={quiz}>
        <QuizResult
          lessonId={lessonId}
          onRetry={retryQuiz}
          quiz={quiz}
          result={result}
          selectedAnswers={selectedAnswers}
        />
      </QuizResultShell>
    );
  }

  return (
    <QuizShell
      currentQuestionIndex={currentQuestionIndex}
      earnedPoints={earnedPoints}
      quiz={quiz}
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
              isCorrect={option.id === currentQuestion.correctOptionId}
              isSelected={selectedOptionId === option.id}
              isLocked={Boolean(selectedOptionId)}
              key={option.id}
              onSelect={() => selectOption(option.id)}
              option={option}
              showFeedback={Boolean(selectedOptionId)}
            />
          ))}
        </div>

        {selectedOptionId ? (
          <QuestionFeedback
            question={currentQuestion}
            selectedOptionId={selectedOptionId}
          />
        ) : null}

        <Button
          className="mt-6 h-14 w-full rounded-2xl bg-[#5147e8] text-base font-bold text-white hover:bg-[#453bd4]"
          disabled={!selectedOptionId}
          onClick={goNext}
          type="button"
        >
          {isFinalQuestion ? "Submit Quiz" : "Next Question"}
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
  quiz: MockQuiz;
  selectedAnswers: SelectedAnswers;
};

function QuizShell({
  children,
  currentQuestionIndex,
  earnedPoints,
  quiz,
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
                className={`size-2.5 rounded-full ${questionDotClass(question, selectedAnswers, index === currentQuestionIndex)}`}
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
  quiz: MockQuiz;
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
  question: MockQuizQuestion,
  selectedAnswers: SelectedAnswers,
  isCurrent: boolean,
) {
  const selectedOptionId = selectedAnswers[question.id];
  if (!selectedOptionId) {
    return isCurrent ? "bg-[#c9c8f8]" : "bg-[#e2e3ee]";
  }

  if (selectedOptionId === question.correctOptionId) {
    return "bg-[#49c389]";
  }

  return "bg-[#ef6a67]";
}

type QuizOptionCardProps = {
  option: MockQuizOption;
  isSelected: boolean;
  isCorrect: boolean;
  isLocked: boolean;
  showFeedback: boolean;
  onSelect: () => void;
};

function QuizOptionCard({
  option,
  isSelected,
  isCorrect,
  isLocked,
  showFeedback,
  onSelect,
}: QuizOptionCardProps) {
  const isWrongSelection = showFeedback && isSelected && !isCorrect;
  const shouldShowCorrect = showFeedback && isCorrect;

  return (
    <button
      className={`flex min-h-[72px] w-full items-center gap-4 rounded-2xl border-2 px-5 py-3.5 text-left transition ${
        shouldShowCorrect
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
          shouldShowCorrect
            ? "bg-[#49c389] text-white"
            : isWrongSelection
              ? "bg-[#ef4444] text-white"
              : "bg-[#eeeef8] text-[#74768a]"
        }`}
      >
        {shouldShowCorrect ? (
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

type QuestionFeedbackProps = {
  question: MockQuizQuestion;
  selectedOptionId: string;
};

function QuestionFeedback({
  question,
  selectedOptionId,
}: QuestionFeedbackProps) {
  const isCorrect = selectedOptionId === question.correctOptionId;

  return (
    <section
      className={`mt-5 rounded-2xl border p-4 ${
        isCorrect
          ? "border-[#a8efd1] bg-[#effcf6]"
          : "border-[#ffc9c5] bg-[#fff3f2]"
      }`}
    >
      <div
        className={`flex items-center gap-3 text-base font-bold ${
          isCorrect ? "text-[#168653]" : "text-[#c8332d]"
        }`}
      >
        {isCorrect ? (
          <Check className="size-5" aria-hidden="true" />
        ) : (
          <X className="size-5" aria-hidden="true" />
        )}
        {isCorrect
          ? `Correct! +${question.points} points`
          : `Incorrect - The answer is option ${correctOptionLabel(question)}`}
      </div>
      <p className="mt-2 text-base font-medium text-[#676982]">
        {question.explanation}
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
  result: QuizResult;
  quiz: MockQuiz;
  selectedAnswers: SelectedAnswers;
  lessonId: string;
  onRetry: () => void;
};

function QuizResult({
  result,
  quiz,
  selectedAnswers,
  lessonId,
  onRetry,
}: QuizResultProps) {
  const incorrectCount = result.totalQuestions - result.correctCount;

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

        <ScoreRing scorePercentage={result.scorePercentage} />

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
              isCorrect={
                selectedAnswers[question.id] === question.correctOptionId
              }
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

type QuizResult = {
  scorePercentage: number;
  correctCount: number;
  totalQuestions: number;
  earnedPoints: number;
  requiredScore: number;
  passed: boolean;
};

function calculateResult(
  quiz: MockQuiz,
  selectedAnswers: SelectedAnswers,
): QuizResult {
  const correctCount = quiz.questions.filter(
    (question) => selectedAnswers[question.id] === question.correctOptionId,
  ).length;
  const scorePercentage = Math.round(
    (correctCount / quiz.questions.length) * 100,
  );

  return {
    scorePercentage,
    correctCount,
    totalQuestions: quiz.questions.length,
    earnedPoints: calculateEarnedPoints(quiz.questions, selectedAnswers),
    requiredScore: quiz.requiredScore,
    passed: scorePercentage >= quiz.requiredScore,
  };
}

function calculateEarnedPoints(
  questions: MockQuizQuestion[],
  selectedAnswers: SelectedAnswers,
) {
  return questions.reduce((total, question) => {
    if (selectedAnswers[question.id] === question.correctOptionId) {
      return total + question.points;
    }

    return total;
  }, 0);
}

function correctOptionLabel(question: MockQuizQuestion) {
  return (
    question.options.find((option) => option.id === question.correctOptionId)
      ?.label ?? ""
  );
}
