import { AlertCircle, Trophy } from "lucide-react";

import { APIError } from "@/api/api";
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
          <RoadmapStateMessage
            title="Loading roadmap"
            description="Preparing your band levels, topics, and lesson progress."
          />
        ) : null}

        {errorMessage ? (
          <RoadmapStateMessage
            title="Roadmap unavailable"
            description={errorMessage}
            tone="error"
          />
        ) : null}

        {isEmpty ? (
          <RoadmapStateMessage
            title="No roadmap data yet"
            description="Seed or publish a course to show band levels, topics, and lessons here."
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

type RoadmapStateMessageProps = {
  title: string;
  description: string;
  tone?: "default" | "error";
};

function RoadmapStateMessage({
  title,
  description,
  tone = "default",
}: RoadmapStateMessageProps) {
  return (
    <section className="relative py-10">
      <div className="mx-auto max-w-sm rounded-2xl border border-[#e3e4f8] bg-white p-6 text-center">
        {tone === "error" ? (
          <AlertCircle
            className="mx-auto size-8 text-destructive"
            aria-hidden="true"
          />
        ) : null}
        <h2 className="text-xl font-bold tracking-normal text-[#676982]">
          {title}
        </h2>
        <p className="mt-3 text-base font-medium text-[#a6a8bb]">
          {description}
        </p>
      </div>
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
