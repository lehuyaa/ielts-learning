import {
  AlertCircle,
  ArrowLeft,
  BookOpen,
  ChevronDown,
  ChevronLeft,
  ChevronRight,
  CheckCircle2,
  Layers,
  Search,
  SlidersHorizontal,
  Star,
} from "lucide-react";
import type React from "react";
import { useMemo, useState } from "react";
import { Link, useNavigate } from "react-router-dom";

import { APIError } from "@/api/api";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { useVocabularies } from "@/features/vocabulary/hooks/useVocabularies";
import { mapVocabularyListItem } from "@/features/vocabulary/mapVocabulary";
import { cn } from "@/lib/utils";
import type {
  VocabularyDifficulty,
  VocabularyDifficultyLabel,
  VocabularyListItemViewModel,
  VocabularyQueryParams,
  VocabularyStatus,
  VocabularyStatusLabel,
} from "@/types/vocabulary";

const difficulties = ["All", "Beginner", "Intermediate", "Advanced"] as const;
const bands = ["All", "5.0", "6.0", "7.0", "8.0"] as const;
const statuses = ["All", "New", "Learning", "Review", "Mastered"] as const;

type DifficultyFilter = (typeof difficulties)[number];
type BandFilter = (typeof bands)[number];
type StatusFilter = (typeof statuses)[number];

const defaultPage = 1;
const defaultLimit = 20;

export function VocabularyPage() {
  const navigate = useNavigate();
  const [query, setQuery] = useState("");
  const [difficulty, setDifficulty] = useState<DifficultyFilter>("All");
  const [band, setBand] = useState<BandFilter>("All");
  const [status, setStatus] = useState<StatusFilter>("All");
  const [page, setPage] = useState(defaultPage);

  const apiParams = useMemo<VocabularyQueryParams>(
    () => ({
      q: query.trim() || undefined,
      difficulty:
        difficulty === "All" ? undefined : difficultyLabelToApi(difficulty),
      targetBand: band === "All" ? undefined : Number(band),
      status: status === "All" ? undefined : statusLabelToApi(status),
      page,
      limit: defaultLimit,
    }),
    [band, difficulty, page, query, status],
  );

  const vocabularyQuery = useVocabularies(apiParams);
  const items = useMemo(
    () => vocabularyQuery.data?.items.map(mapVocabularyListItem) ?? [],
    [vocabularyQuery.data],
  );
  const pagination = vocabularyQuery.data?.pagination;
  const errorMessage = getVocabularyErrorMessage(vocabularyQuery.error);
  const masteredCount = items.filter(
    (item) => item.status === "Mastered",
  ).length;
  const isEmpty =
    !vocabularyQuery.isLoading && !errorMessage && items.length === 0;

  function resetToFirstPage() {
    setPage(defaultPage);
  }

  return (
    <div className="min-h-screen bg-[#f7f7fc] text-[#10111f]">
      <VocabularyHeader
        masteredCount={masteredCount}
        onBack={() => navigate("/dashboard")}
        page={pagination?.page ?? page}
        totalPages={pagination?.totalPages ?? 0}
        totalWords={pagination?.total ?? 0}
      />

      <main className="mx-auto max-w-5xl px-4 pb-8 pt-6 lg:px-0">
        <section className="mt-6 rounded-2xl border border-[#e6e6f3] bg-white p-4 shadow-sm">
          <div className="grid gap-4 lg:grid-cols-[1fr_auto_auto_auto]">
            <label className="relative block lg:self-end">
              <span className="sr-only">Search vocabulary</span>
              <Search
                className="pointer-events-none absolute left-4 top-1/2 size-5 -translate-y-1/2 text-[#85889c]"
                aria-hidden="true"
              />
              <Input
                className="h-12 rounded-full pl-12 text-base"
                onChange={(event) => {
                  setQuery(event.target.value);
                  resetToFirstPage();
                }}
                placeholder="Search word, topic, or meaning"
                value={query}
              />
            </label>

            <FilterSelect
              label="Difficulty"
              onChange={(value) => {
                setDifficulty(value as DifficultyFilter);
                resetToFirstPage();
              }}
              options={difficulties}
              value={difficulty}
            />
            <FilterSelect
              label="Band"
              onChange={(value) => {
                setBand(value as BandFilter);
                resetToFirstPage();
              }}
              options={bands}
              value={band}
            />
            <FilterSelect
              label="Status"
              onChange={(value) => {
                setStatus(value as StatusFilter);
                resetToFirstPage();
              }}
              options={statuses}
              value={status}
            />
          </div>
        </section>

        {vocabularyQuery.isLoading ? (
          <VocabularyStateMessage
            description="Loading words, progress, and filter results."
            title="Loading vocabulary"
          />
        ) : null}

        {errorMessage ? (
          <VocabularyStateMessage
            description={errorMessage}
            title="Vocabulary unavailable"
            tone="error"
          />
        ) : null}

        {isEmpty ? (
          <VocabularyStateMessage
            description="Try a different search term or clear one of the filters."
            icon="filters"
            title="No matching words"
          />
        ) : null}

        {!vocabularyQuery.isLoading && !errorMessage && items.length > 0 ? (
          <>
            <section className="mt-6 grid gap-4 lg:grid-cols-2">
              {items.map((item) => (
                <VocabularyListCard item={item} key={item.id} />
              ))}
            </section>

            {pagination ? (
              <PaginationControls
                onPageChange={setPage}
                pagination={pagination}
              />
            ) : null}
          </>
        ) : null}
      </main>
    </div>
  );
}

type VocabularyHeaderProps = {
  totalWords: number;
  page: number;
  totalPages: number;
  masteredCount: number;
  onBack: () => void;
};

function VocabularyHeader({
  totalWords,
  page,
  totalPages,
  masteredCount,
  onBack,
}: VocabularyHeaderProps) {
  return (
    <>
      <header className="border-b border-[#e6e6f3] bg-white">
        <div className="mx-auto flex min-h-[72px] max-w-5xl items-center justify-between gap-5 px-4 lg:px-0">
          <div className="flex min-w-0 items-center gap-4">
            <button
              aria-label="Go back"
              className="grid size-9 shrink-0 cursor-pointer place-items-center rounded-full text-[#6d7088] transition-colors hover:bg-[#f0f1fb]"
              onClick={onBack}
              type="button"
            >
              <ArrowLeft className="size-5" aria-hidden="true" />
            </button>

            <div className="grid size-9 shrink-0 place-items-center rounded-full bg-primary text-white">
              <BookOpen className="size-5" aria-hidden="true" />
            </div>

            <div className="min-w-0">
              <h1 className="truncate text-xl font-bold tracking-normal text-[#10111f]">
                Vocabulary Library
              </h1>
              <p className="mt-1 text-base font-medium text-[#676982]">
                Search, filter, and review IELTS vocabulary
              </p>
            </div>
          </div>

          <div className="hidden rounded-full border border-[#dde3ff] bg-[#f7f7ff] px-4 py-2 text-sm font-bold text-primary md:flex md:items-center md:gap-2">
            <Layers className="size-4" aria-hidden="true" />
            Word Bank
          </div>
        </div>
      </header>

      <section className="mx-auto max-w-5xl px-4 pt-8 lg:px-0">
        <div className="grid gap-4 lg:grid-cols-3">
          <HeaderStatCard
            icon={
              <BookOpen className="size-5 text-primary" aria-hidden="true" />
            }
            label="Total Words"
            value={String(totalWords)}
          />
          <HeaderStatCard
            icon={<Star className="size-5 text-warning" aria-hidden="true" />}
            label="Current Page"
            value={totalPages > 0 ? `${page}/${totalPages}` : "-"}
          />
          <HeaderStatCard
            icon={
              <CheckCircle2
                className="size-5 text-success"
                aria-hidden="true"
              />
            }
            label="Mastered Here"
            value={String(masteredCount)}
          />
        </div>
      </section>
    </>
  );
}

type HeaderStatCardProps = {
  icon: React.ReactNode;
  label: string;
  value: string;
};

function HeaderStatCard({ icon, label, value }: HeaderStatCardProps) {
  return (
    <article className="flex items-center gap-4 rounded-2xl border border-[#e3e4f8] bg-white p-4 shadow-sm">
      <div className="grid size-9 shrink-0 place-items-center rounded-full bg-[#ececf6]">
        {icon}
      </div>
      <div>
        <p className="text-xl font-bold tracking-normal text-[#10111f]">
          {value}
        </p>
        <p className="mt-1 text-sm font-medium text-[#676982]">{label}</p>
      </div>
    </article>
  );
}

type FilterSelectProps = {
  label: string;
  options: readonly string[];
  value: string;
  onChange: (value: string) => void;
};

function FilterSelect({ label, options, value, onChange }: FilterSelectProps) {
  return (
    <label className="grid gap-2">
      <span className="text-sm font-bold text-[#676982]">{label}</span>
      <span className="relative block">
        <select
          className="h-11 min-w-32 cursor-pointer appearance-none rounded-full border border-[#e6e6f3] bg-white py-0 pl-4 pr-9 text-sm font-semibold text-[#10111f] shadow-sm outline-none transition-colors focus:border-primary disabled:cursor-not-allowed disabled:opacity-50"
          onChange={(event) => onChange(event.target.value)}
          value={value}
        >
          {options.map((option) => (
            <option key={option} value={option}>
              {option}
            </option>
          ))}
        </select>
        <ChevronDown
          className="pointer-events-none absolute right-3 top-1/2 size-4 -translate-y-1/2 text-[#10111f]"
          aria-hidden="true"
        />
      </span>
    </label>
  );
}

type VocabularyListCardProps = {
  item: VocabularyListItemViewModel;
};

function VocabularyListCard({ item }: VocabularyListCardProps) {
  return (
    <Link aria-label={`Open ${item.word}`} to={`/vocabulary/${item.id}`}>
      <article className="h-full cursor-pointer rounded-2xl border border-[#e6e6f3] bg-white p-5 shadow-sm transition-colors hover:border-[#c6c4ff]">
        <div className="flex flex-col gap-4 sm:flex-row sm:items-start sm:justify-between">
          <div className="min-w-0">
            <div className="flex flex-wrap items-center gap-3">
              <h2 className="text-xl font-bold tracking-normal text-[#10111f]">
                {item.word}
              </h2>
              <DifficultyBadge difficulty={item.difficulty} />
            </div>
            <p className="mt-2 font-mono text-base font-semibold text-[#85889c]">
              {item.ipa} · {item.partOfSpeech}
            </p>
            <p className="mt-3 line-clamp-2 text-base font-medium text-[#676982]">
              {item.shortDefinition}
            </p>
          </div>

          <StatusBadge status={item.status} />
        </div>

        <div className="mt-5 flex flex-wrap items-center gap-3 text-sm font-bold">
          <Badge className="bg-[#f0efff] text-primary">{item.band}</Badge>
          <Badge className="bg-[#eff6ff] text-[#2563eb]">{item.topic}</Badge>
          <span className="text-[#676982]">{item.frequency} frequency</span>
        </div>

        <div className="mt-5 h-2 overflow-hidden rounded-full bg-[#ebeaff]">
          <div
            className="h-full rounded-full bg-primary"
            style={{ width: `${item.masteryScore}%` }}
          />
        </div>
        <p className="mt-2 text-sm font-bold text-[#676982]">
          {item.masteryScore}% mastery
        </p>
      </article>
    </Link>
  );
}

function DifficultyBadge({
  difficulty,
}: {
  difficulty: VocabularyDifficultyLabel;
}) {
  const styles = {
    Beginner: "bg-[#dcfce7] text-[#138a53]",
    Intermediate: "bg-[#fff3c4] text-[#c46700]",
    Advanced: "bg-[#fde5e7] text-[#d22f38]",
  };

  return (
    <span
      className={cn(
        "rounded-full px-3 py-1 text-sm font-bold",
        styles[difficulty],
      )}
    >
      {difficulty}
    </span>
  );
}

function StatusBadge({ status }: { status: VocabularyStatusLabel }) {
  const styles = {
    New: "bg-[#f3f4f6] text-[#4b5563]",
    Learning: "bg-[#eff6ff] text-[#2563eb]",
    Review: "bg-[#fff7ed] text-[#c46600]",
    Mastered: "bg-[#e8fff3] text-[#138a53]",
  };

  return (
    <span
      className={cn("rounded-full px-3 py-1 text-sm font-bold", styles[status])}
    >
      {status}
    </span>
  );
}

type PaginationControlsProps = {
  pagination: {
    page: number;
    totalPages: number;
    total: number;
  };
  onPageChange: (page: number) => void;
};

function PaginationControls({
  pagination,
  onPageChange,
}: PaginationControlsProps) {
  if (pagination.totalPages <= 1) {
    return null;
  }

  return (
    <section className="mt-6 flex flex-col gap-4 rounded-2xl border border-[#e6e6f3] bg-white p-4 shadow-sm sm:flex-row sm:items-center sm:justify-between">
      <p className="text-sm font-bold text-[#676982]">
        Page {pagination.page} of {pagination.totalPages} · {pagination.total}{" "}
        words
      </p>
      <div className="flex gap-3">
        <Button
          className="rounded-full"
          disabled={pagination.page <= 1}
          onClick={() => onPageChange(pagination.page - 1)}
          type="button"
          variant="outline"
        >
          <ChevronLeft aria-hidden="true" />
          Previous
        </Button>
        <Button
          className="rounded-full"
          disabled={pagination.page >= pagination.totalPages}
          onClick={() => onPageChange(pagination.page + 1)}
          type="button"
        >
          Next
          <ChevronRight aria-hidden="true" />
        </Button>
      </div>
    </section>
  );
}

type VocabularyStateMessageProps = {
  title: string;
  description: string;
  tone?: "default" | "error";
  icon?: "filters";
};

function VocabularyStateMessage({
  title,
  description,
  tone = "default",
  icon,
}: VocabularyStateMessageProps) {
  return (
    <section className="mt-6 rounded-2xl border border-[#e6e6f3] bg-white p-6 text-center shadow-sm">
      {tone === "error" ? (
        <AlertCircle
          className="mx-auto size-9 text-destructive"
          aria-hidden="true"
        />
      ) : null}
      {icon === "filters" ? (
        <SlidersHorizontal
          className="mx-auto size-9 text-[#85889c]"
          aria-hidden="true"
        />
      ) : null}
      <h2 className="mt-4 text-xl font-bold tracking-normal">{title}</h2>
      <p className="mx-auto mt-3 max-w-md text-base font-medium text-[#676982]">
        {description}
      </p>
    </section>
  );
}

function difficultyLabelToApi(
  difficulty: Exclude<DifficultyFilter, "All">,
): VocabularyDifficulty {
  switch (difficulty) {
    case "Beginner":
      return "BEGINNER";
    case "Advanced":
      return "ADVANCED";
    case "Intermediate":
    default:
      return "INTERMEDIATE";
  }
}

function statusLabelToApi(
  status: Exclude<StatusFilter, "All">,
): VocabularyStatus {
  switch (status) {
    case "Learning":
      return "LEARNING";
    case "Review":
      return "REVIEW";
    case "Mastered":
      return "MASTERED";
    case "New":
    default:
      return "NEW";
  }
}

function getVocabularyErrorMessage(error: Error | null) {
  if (!error) {
    return null;
  }

  if (error instanceof APIError) {
    return error.message;
  }

  return "Unable to load vocabulary right now.";
}
