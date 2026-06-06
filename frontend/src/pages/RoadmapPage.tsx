import { RoadmapBandSection } from "@/features/roadmap/RoadmapBandSection";
import { RoadmapHeader } from "@/features/roadmap/RoadmapHeader";
import { roadmapBands } from "@/features/roadmap/mockRoadmap";
import { Trophy } from "lucide-react";

export function RoadmapPage() {
  return (
    <div className="min-h-screen bg-[#f7f7fc] text-[#10111f]">
      <RoadmapHeader />

      <div className="relative mx-auto max-w-5xl px-4 pb-10 pt-6 lg:px-0">
        <div
          className="absolute bottom-10 left-1/2 top-6 hidden w-0.5 -translate-x-1/2 bg-[#e6e6ff] lg:block"
          aria-hidden="true"
        />

        {roadmapBands.map((band) => (
          <RoadmapBandSection band={band} key={band.id} />
        ))}

        <section className="relative py-10">
          <div className="mx-auto max-w-sm rounded-2xl border border-[#f5d675] bg-[#fffdf7] p-6 text-center">
            <Trophy
              className="mx-auto size-9 text-[#ffc15a]"
              aria-hidden="true"
            />
            <h2 className="mt-5 text-xl font-bold tracking-normal text-[#676982]">
              Band 8.5+ Master
            </h2>
            <p className="mt-5 text-base font-medium text-[#a6a8bb]">
              Complete all bands to unlock
            </p>
          </div>
        </section>
      </div>
    </div>
  );
}
