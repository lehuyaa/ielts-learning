import { Trophy } from "lucide-react";

import { APIError } from "@/api/api";
import { CardSkeleton } from "@/components/state/CardSkeleton";
import { EmptyState } from "@/components/state/EmptyState";
import { ErrorState } from "@/components/state/ErrorState";
import { Skeleton } from "@/components/ui/skeleton";
import { RoadmapBandSection } from "@/features/roadmap/RoadmapBandSection";
import { RoadmapHeader } from "@/features/roadmap/RoadmapHeader";
import { useRoadmap } from "@/features/roadmap/hooks/useRoadmap";
import { mapRoadmapToViewModel } from "@/features/roadmap/mapRoadmap";
import type { RoadmapViewModel } from "@/features/roadmap/types";

const loadingRoadmap: RoadmapViewModel = {
  title: "Vocabulary Roadmap",
  subtitle: "Band 5.0 → 8.5",
  topicsCompleted: 0,
  totalTopics: 0,
  currentBand: null,
  wordsMastered: 0,
  currentStreak: 0,
  bands: [],
};

export function RoadmapPage() {
  const roadmapQuery = useRoadmap();
  const roadmap = roadmapQuery.data
    ? mapRoadmapToViewModel(roadmapQuery.data)
    : null;
  const viewModel = roadmap ?? loadingRoadmap;
  const errorMessage = getRoadmapErrorMessage(roadmapQuery.error);
  const isEmpty =
    !roadmapQuery.isLoading && !errorMessage && viewModel.bands.length === 0;

  return (
    <div className="min-h-screen bg-background text-foreground">
      <RoadmapHeader
        title={viewModel.title}
        subtitle={viewModel.subtitle}
        topicsCompleted={viewModel.topicsCompleted}
        totalTopics={viewModel.totalTopics}
        currentBand={viewModel.currentBand}
        wordsMastered={viewModel.wordsMastered}
        currentStreak={viewModel.currentStreak}
      />

      <div className="relative mx-auto max-w-5xl px-6 py-8">
        <div
          className="absolute bottom-28 left-1/2 top-8 hidden w-0.5 -translate-x-1/2 bg-border md:block"
          aria-hidden="true"
        />

        {roadmapQuery.isLoading ? (
          <RoadmapLoadingSkeleton />
        ) : null}

        {errorMessage ? (
          <ErrorState
            actionHref="/dashboard"
            actionLabel="Go to dashboard"
            description={errorMessage}
            onRetry={() => {
              void roadmapQuery.refetch();
            }}
            title="Roadmap unavailable"
          />
        ) : null}

        {isEmpty ? (
          <EmptyState
            actionHref="/dashboard"
            actionLabel="Go to dashboard"
            description="Seed or publish a course to show band levels, topics, and lessons here."
            title="No roadmap data yet"
          />
        ) : null}

        {!roadmapQuery.isLoading && !errorMessage
          ? viewModel.bands.map((band, index) => (
              <RoadmapBandSection
                band={band}
                isLast={index === viewModel.bands.length - 1}
                key={band.id}
              />
            ))
          : null}

        {!roadmapQuery.isLoading && !errorMessage && !isEmpty ? (
          <RoadmapMasterCard />
        ) : null}
      </div>
    </div>
  );
}

function getRoadmapErrorMessage(error: Error | null) {
  if (!error) {
    return null;
  }

  return error instanceof APIError
    ? error.message
    : "Unable to load the roadmap right now.";
}

function RoadmapLoadingSkeleton() {
  return (
    <section className="relative space-y-8 py-2">
      {Array.from({ length: 3 }).map((_, index) => (
        <div className="relative z-10" key={index}>
          <div className="mx-auto mb-6 flex w-full max-w-[320px] justify-center">
            <div className="w-full rounded-2xl border border-[#e3e4f8] bg-white p-4 shadow-sm">
              <div className="flex items-center gap-3">
                <Skeleton className="size-10 rounded-full" />
                <div className="flex-1 space-y-2">
                  <Skeleton className="h-5 w-28" />
                  <Skeleton className="h-4 w-20" />
                </div>
                <Skeleton className="h-6 w-20 rounded-full" />
              </div>
            </div>
          </div>

          <div className="grid gap-4 md:grid-cols-2">
            <CardSkeleton className="border-[#e3e4f8]" lines={3} showIcon />
            <CardSkeleton className="border-[#e3e4f8]" lines={3} showIcon />
          </div>
        </div>
      ))}
    </section>
  );
}

function RoadmapMasterCard() {
  return (
    <section className="relative z-10 flex justify-center py-8">
      <div className="flex flex-col items-center gap-2 rounded-2xl border border-amber-200 bg-amber-50 px-8 py-5 text-center opacity-50">
        <Trophy className="size-7 text-amber-500" aria-hidden="true" />
        <h2 className="text-base font-bold tracking-normal text-foreground">
          Band 8.5+ Master
        </h2>
        <p className="text-xs text-muted-foreground">
          Complete all bands to unlock
        </p>
      </div>
    </section>
  );
}
