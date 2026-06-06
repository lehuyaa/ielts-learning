import { ArrowLeft, CheckCircle2, RotateCcw, Volume2 } from "lucide-react";
import type React from "react";
import { useRef, useState } from "react";
import { Link, useNavigate } from "react-router-dom";

import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";

import type { FlashcardRating, MockFlashcard } from "./mockFlashcards";

type RatingCounts = Record<FlashcardRating, number>;

type FlashcardSessionProps = {
  cards: MockFlashcard[];
  title: string;
  subtitle: string;
  completionPrimaryHref: string;
  completionPrimaryLabel: string;
  emptyTitle?: string;
  emptyDescription?: string;
};

const emptyRatings: RatingCounts = {
  again: 0,
  hard: 0,
  good: 0,
  easy: 0,
};

const ratingOptions: Array<{
  value: FlashcardRating;
  label: string;
  interval: string;
  className: string;
}> = [
  {
    value: "again",
    label: "Again",
    interval: "<10 min",
    className: "border-[#ffc8c8] bg-[#fff8f8] text-[#dc2626]",
  },
  {
    value: "hard",
    label: "Hard",
    interval: "1 day",
    className: "border-[#ffd89b] bg-[#fffaf0] text-[#ea580c]",
  },
  {
    value: "good",
    label: "Good",
    interval: "3 days",
    className: "border-[#bfd8ff] bg-[#f4f9ff] text-[#2563eb]",
  },
  {
    value: "easy",
    label: "Easy",
    interval: "7 days",
    className: "border-[#bff0d5] bg-[#f2fff8] text-[#15803d]",
  },
];

export function FlashcardSession({
  cards,
  title,
  subtitle,
  completionPrimaryHref,
  completionPrimaryLabel,
  emptyTitle = "No reviews due",
  emptyDescription = "You are all caught up. Come back later for more review cards.",
}: FlashcardSessionProps) {
  const navigate = useNavigate();
  const [currentIndex, setCurrentIndex] = useState(0);
  const [isFlipped, setIsFlipped] = useState(false);
  const [isRatingVisible, setIsRatingVisible] = useState(false);
  const [ratingCounts, setRatingCounts] = useState<RatingCounts>(emptyRatings);
  const [isCompleted, setIsCompleted] = useState(false);
  const ratingTimerRef = useRef<number | null>(null);

  const currentCard = cards[currentIndex];
  const totalCards = cards.length;
  const doneCount = isCompleted ? totalCards : currentIndex;
  const remainingCount = Math.max(totalCards - doneCount, 0);
  const progressPercentage =
    totalCards === 0 ? 0 : Math.round((doneCount / totalCards) * 100);

  function resetSession() {
    clearRatingTimer();
    setCurrentIndex(0);
    setIsFlipped(false);
    setIsRatingVisible(false);
    setRatingCounts(emptyRatings);
    setIsCompleted(false);
  }

  function flipCard() {
    clearRatingTimer();
    setIsRatingVisible(false);

    if (isFlipped) {
      setIsFlipped(false);
      return;
    }

    setIsFlipped(true);

    ratingTimerRef.current = window.setTimeout(() => {
      setIsRatingVisible(true);
      ratingTimerRef.current = null;
    }, 600);
  }

  function clearRatingTimer() {
    if (ratingTimerRef.current === null) {
      return;
    }

    window.clearTimeout(ratingTimerRef.current);
    ratingTimerRef.current = null;
  }

  function rateCard(rating: FlashcardRating) {
    clearRatingTimer();
    setRatingCounts((current) => ({
      ...current,
      [rating]: current[rating] + 1,
    }));

    if (currentIndex >= totalCards - 1) {
      setIsCompleted(true);
      setIsFlipped(false);
      setIsRatingVisible(false);
      return;
    }

    setCurrentIndex((index) => index + 1);
    setIsFlipped(false);
    setIsRatingVisible(false);
  }

  if (totalCards === 0) {
    return (
      <FlashcardShell
        doneCount={0}
        onBack={() => navigate(-1)}
        onReset={resetSession}
        progressPercentage={0}
        remainingCount={0}
      >
        <EmptyFlashcardState
          description={emptyDescription}
          title={emptyTitle}
        />
      </FlashcardShell>
    );
  }

  if (isCompleted) {
    return (
      <FlashcardShell
        doneCount={doneCount}
        onBack={() => navigate(-1)}
        onReset={resetSession}
        progressPercentage={100}
        remainingCount={0}
      >
        <CompletionSummary
          completionPrimaryHref={completionPrimaryHref}
          completionPrimaryLabel={completionPrimaryLabel}
          ratingCounts={ratingCounts}
          resetSession={resetSession}
          totalCards={totalCards}
        />
      </FlashcardShell>
    );
  }

  return (
    <FlashcardShell
      doneCount={doneCount}
      onBack={() => navigate(-1)}
      onReset={resetSession}
      progressPercentage={progressPercentage}
      remainingCount={remainingCount}
    >
      <main className="mx-auto flex min-h-[calc(100vh-72px)] w-full max-w-xl flex-col items-center justify-center px-4 py-8">
        <div className="mb-6 flex items-center justify-center gap-3 text-sm font-semibold text-[#6d7088]">
          <span className="rounded-full bg-[#eff1ff] px-3 py-1.5 text-primary">
            {currentCard.topicTitle}
          </span>
          <span>
            {currentIndex + 1}/{totalCards}
          </span>
          <span className="rounded-full bg-[#ececf3] px-3 py-1 text-sm">
            {currentCard.band}
          </span>
        </div>

        <FlipCard card={currentCard} isFlipped={isFlipped} onFlip={flipCard} />

        {isRatingVisible ? (
          <RatingControls onRate={rateCard} />
        ) : (
          <div className="mt-5 h-10" aria-hidden="true" />
        )}

        <StepDots currentIndex={currentIndex} totalCards={totalCards} />

        <div className="sr-only" aria-live="polite">
          {title}. {subtitle}. Card {currentIndex + 1} of {totalCards}.
        </div>
      </main>
    </FlashcardShell>
  );
}

type FlashcardShellProps = {
  children: React.ReactNode;
  doneCount: number;
  remainingCount: number;
  progressPercentage: number;
  onBack: () => void;
  onReset: () => void;
};

function FlashcardShell({
  children,
  doneCount,
  remainingCount,
  progressPercentage,
  onBack,
  onReset,
}: FlashcardShellProps) {
  return (
    <div className="min-h-screen bg-[#f8f8ff] text-[#10111f]">
      <header className="border-b border-[#e6e6f3] bg-white">
        <div className="mx-auto grid min-h-[60px] max-w-3xl grid-cols-[36px_1fr_36px] items-center gap-4 px-4">
          <button
            aria-label="Go back"
            className="grid size-9 place-items-center rounded-full text-[#6d7088] transition-colors hover:bg-[#f0f1fb]"
            onClick={onBack}
            type="button"
          >
            <ArrowLeft className="size-4" aria-hidden="true" />
          </button>

          <div>
            <div className="h-2 overflow-hidden rounded-full bg-[#e9e9f3]">
              <div
                className="h-full rounded-full bg-gradient-to-r from-primary to-[#55c79a]"
                style={{ width: `${progressPercentage}%` }}
              />
            </div>
            <div className="mt-1.5 flex items-center justify-between text-sm font-medium text-[#676982]">
              <span>{doneCount} done</span>
              <span>{remainingCount} remaining</span>
            </div>
          </div>

          <button
            aria-label="Restart session"
            className="grid size-9 place-items-center rounded-full text-[#6d7088] transition-colors hover:bg-[#f0f1fb]"
            onClick={onReset}
            type="button"
          >
            <RotateCcw className="size-4" aria-hidden="true" />
          </button>
        </div>
      </header>

      {children}
    </div>
  );
}

function FlipCard({
  card,
  isFlipped,
  onFlip,
}: {
  card: MockFlashcard;
  isFlipped: boolean;
  onFlip: () => void;
}) {
  return (
    <div
      className={cn(
        "w-full max-w-xl [perspective:1400px]",
        isFlipped ? "min-h-[300px]" : "h-[300px]",
      )}
    >
      <div
        className={cn(
          "relative grid transition-transform duration-[600ms] [transform-style:preserve-3d] motion-reduce:transition-none",
          isFlipped ? "[transform:rotateY(180deg)]" : "h-full",
        )}
      >
        <div
          className={cn(
            "col-start-1 row-start-1 [backface-visibility:hidden]",
            isFlipped ? "absolute inset-0" : "h-full",
          )}
        >
          <FlashcardFront card={card} onFlip={onFlip} />
        </div>
        <div
          className={cn(
            "col-start-1 row-start-1 [backface-visibility:hidden] [transform:rotateY(180deg)]",
            isFlipped ? "relative" : "absolute inset-0",
          )}
        >
          <FlashcardBack card={card} onFlip={onFlip} />
        </div>
      </div>
    </div>
  );
}

function FlashcardFront({
  card,
  onFlip,
}: {
  card: MockFlashcard;
  onFlip: () => void;
}) {
  return (
    <div
      className="flex h-full w-full flex-col items-center cursor-pointer justify-center rounded-[24px] bg-gradient-to-br from-[#6258f6] to-[#8318e8] px-6 py-10 text-center text-white shadow-[0_22px_60px_rgba(79,70,229,0.18)] md:px-8"
      onClick={onFlip}
      onKeyDown={(event) => handleFlipKeyDown(event, onFlip)}
      role="button"
      tabIndex={0}
    >
      <p className="font-mono text-base font-semibold text-white/75">
        {card.ipa}
      </p>
      <h1 className="mt-4 text-4xl font-bold leading-tight tracking-normal md:text-5xl">
        {card.word}
      </h1>
      <p className="mt-4 text-base font-medium text-white/85">
        {card.partOfSpeech}
      </p>
      <p className="mt-8 inline-flex items-center gap-2 text-sm font-medium text-white/80">
        <Volume2 className="size-3.5" aria-hidden="true" />
        Tap to reveal meaning
      </p>
    </div>
  );
}

function FlashcardBack({
  card,
  onFlip,
}: {
  card: MockFlashcard;
  onFlip: () => void;
}) {
  return (
    <div
      onClick={onFlip}
      onKeyDown={(event) => handleFlipKeyDown(event, onFlip)}
      className="min-h-[300px] w-full cursor-pointer rounded-[24px] border border-[#e6e6f3] bg-white p-5 shadow-[0_22px_60px_rgba(16,17,31,0.10)] md:p-7"
      role="button"
      tabIndex={0}
    >
      <div className="text-center">
        <h1 className="text-2xl font-bold tracking-normal text-[#10111f] md:text-3xl">
          {card.word}
        </h1>
        <p className="mt-2 font-mono text-base font-semibold text-[#85889c]">
          {card.ipa}
        </p>
      </div>

      <section className="mt-6 text-left">
        <p className="text-sm font-bold uppercase tracking-normal text-[#676982]">
          Meaning
        </p>
        <p className="mt-2 text-lg font-medium leading-relaxed text-[#10111f]">
          {card.meaningEn}
        </p>
        <p className="mt-1.5 text-sm font-medium text-[#676982]">
          {card.meaningVi}
        </p>
      </section>

      <section className="mt-5 text-left">
        <p className="text-sm font-bold uppercase tracking-normal text-[#676982]">
          Example
        </p>
        <blockquote className="mt-2 rounded-2xl border-l-4 border-primary bg-[#eff1ff] px-4 py-3 text-base font-medium text-[#30348e]">
          "{card.exampleSentence}"
        </blockquote>
      </section>

      <TokenSection label="Synonyms" tokens={card.synonyms} />
    </div>
  );
}

function handleFlipKeyDown(
  event: React.KeyboardEvent<HTMLDivElement>,
  onFlip: () => void,
) {
  if (event.key !== "Enter" && event.key !== " ") {
    return;
  }

  event.preventDefault();
  onFlip();
}

function TokenSection({ label, tokens }: { label: string; tokens: string[] }) {
  if (tokens.length === 0) {
    return null;
  }

  return (
    <section className="mt-5">
      <p className="text-sm text-left font-bold uppercase tracking-normal text-[#676982]">
        {label}
      </p>
      <div className="mt-2 flex flex-wrap gap-2">
        {tokens.map((token) => (
          <span
            className="rounded-full bg-[#ececf3] px-3 py-1.5 text-sm font-bold text-[#676982]"
            key={token}
          >
            {token}
          </span>
        ))}
      </div>
    </section>
  );
}

function RatingControls({
  onRate,
}: {
  onRate: (rating: FlashcardRating) => void;
}) {
  return (
    <section className="mt-6 w-full">
      <p className="text-center text-sm font-semibold text-[#676982]">
        How well did you know this word?
      </p>
      <div className="mt-4 grid gap-3 sm:grid-cols-4">
        {ratingOptions.map((option) => (
          <button
            className={cn(
              "rounded-2xl border px-4 py-3 text-center shadow-sm transition-transform hover:-translate-y-0.5",
              option.className,
            )}
            key={option.value}
            onClick={(event) => {
              event.stopPropagation();
              onRate(option.value);
            }}
            type="button"
          >
            <span className="block text-lg font-bold">{option.label}</span>
            <span className="mt-1 block text-sm font-medium">
              {option.interval}
            </span>
          </button>
        ))}
      </div>
    </section>
  );
}

function StepDots({
  currentIndex,
  totalCards,
}: {
  currentIndex: number;
  totalCards: number;
}) {
  return (
    <div className="mt-6 flex justify-center gap-2">
      {Array.from({ length: totalCards }).map((_, index) => (
        <span
          className={cn(
            "h-1.5 rounded-full transition-all",
            index === currentIndex ? "w-7 bg-primary" : "w-3.5 bg-[#ececf3]",
          )}
          key={index}
        />
      ))}
    </div>
  );
}

function CompletionSummary({
  totalCards,
  ratingCounts,
  resetSession,
  completionPrimaryHref,
  completionPrimaryLabel,
}: {
  totalCards: number;
  ratingCounts: RatingCounts;
  resetSession: () => void;
  completionPrimaryHref: string;
  completionPrimaryLabel: string;
}) {
  return (
    <main className="mx-auto grid min-h-[calc(100vh-72px)] max-w-xl place-items-center px-4 py-8">
      <section className="w-full rounded-[24px] border border-[#e6e6f3] bg-white p-6 text-center shadow-sm md:p-8">
        <div className="mx-auto grid size-12 place-items-center rounded-full bg-[#e8fff3] text-success">
          <CheckCircle2 className="size-7" aria-hidden="true" />
        </div>
        <h1 className="mt-5 text-3xl font-bold tracking-normal">
          Session Complete
        </h1>
        <p className="mx-auto mt-3 max-w-md text-base font-medium text-[#676982]">
          You reviewed {totalCards} IELTS vocabulary cards. Nice focused
          practice.
        </p>

        <div className="mt-6 grid gap-3 sm:grid-cols-4">
          {ratingOptions.map((option) => (
            <div
              className={cn("rounded-2xl border p-3", option.className)}
              key={option.value}
            >
              <p className="text-xl font-bold">{ratingCounts[option.value]}</p>
              <p className="mt-1 text-sm font-semibold">{option.label}</p>
            </div>
          ))}
        </div>

        <div className="mt-6 flex flex-col justify-center gap-3 sm:flex-row">
          <Button asChild className="h-10 rounded-full px-5 text-sm">
            <Link to={completionPrimaryHref}>{completionPrimaryLabel}</Link>
          </Button>
          <Button
            className="h-10 rounded-full px-5 text-sm"
            onClick={resetSession}
            type="button"
            variant="outline"
          >
            Review Again
          </Button>
        </div>
      </section>
    </main>
  );
}

function EmptyFlashcardState({
  title,
  description,
}: {
  title: string;
  description: string;
}) {
  return (
    <main className="mx-auto grid min-h-[calc(100vh-72px)] max-w-lg place-items-center px-4 py-8">
      <section className="w-full rounded-[24px] border border-[#e6e6f3] bg-white p-6 text-center shadow-sm md:p-8">
        <div className="mx-auto grid size-12 place-items-center rounded-full bg-[#e8fff3] text-success">
          <CheckCircle2 className="size-7" aria-hidden="true" />
        </div>
        <h1 className="mt-5 text-2xl font-bold tracking-normal">{title}</h1>
        <p className="mt-3 text-base font-medium text-[#676982]">
          {description}
        </p>
        <Button asChild className="mt-6 h-10 rounded-full px-5 text-sm">
          <Link to="/roadmap">Back to roadmap</Link>
        </Button>
      </section>
    </main>
  );
}
