import {
  ArrowLeft,
  Flame,
  Star,
  Target,
  TrendingUp,
  CheckCircle2,
} from "lucide-react";
import type React from "react";
import { useNavigate } from "react-router-dom";

export function RoadmapHeader() {
  const navigate = useNavigate();

  return (
    <>
      <header className="border-b border-[#e6e6f3] bg-white">
        <div className="mx-auto flex min-h-20 max-w-5xl items-center justify-between gap-6 px-4 lg:px-0">
          <div className="flex min-w-0 items-center gap-5">
            <button
              aria-label="Go back"
              className="grid size-9 shrink-0 place-items-center rounded-full text-[#6d7088] transition-colors hover:bg-[#f0f1fb]"
              onClick={() => navigate(-1)}
              type="button"
            >
              <ArrowLeft className="size-5" aria-hidden="true" />
            </button>

            <div className="grid size-10 shrink-0 place-items-center rounded-full bg-primary text-white">
              <Target className="size-6" aria-hidden="true" />
            </div>

            <div className="min-w-0">
              <h1 className="truncate text-xl font-bold tracking-normal text-[#10111f]">
                Vocabulary Roadmap
              </h1>
              <p className="mt-1 text-base font-medium text-[#676982]">
                Band 5.0 → 8.5
              </p>
            </div>
          </div>

          <div className="hidden rounded-full border border-[#ffd75f] bg-[#fffdf5] px-5 py-2 text-base font-bold text-[#b84a00] md:flex md:items-center md:gap-2">
            <Flame className="size-4 text-[#ff6b00]" aria-hidden="true" />
            14 Day Streak
          </div>
        </div>
      </header>

      <section className="mx-auto max-w-5xl px-4 pt-10 lg:px-0">
        <div className="grid gap-5 lg:grid-cols-3">
          <StatCard
            icon={
              <CheckCircle2
                className="size-5 text-success"
                aria-hidden="true"
              />
            }
            label="Topics Completed"
            value="3/19"
          />
          <StatCard
            icon={<Star className="size-5 text-warning" aria-hidden="true" />}
            label="Current Band"
            value="6.5"
          />
          <StatCard
            icon={
              <TrendingUp className="size-5 text-primary" aria-hidden="true" />
            }
            label="Words Mastered"
            value="847"
          />
        </div>
      </section>
    </>
  );
}

type StatCardProps = {
  icon: React.ReactNode;
  label: string;
  value: string;
};

function StatCard({ icon, label, value }: StatCardProps) {
  return (
    <article className="flex items-center gap-4 rounded-2xl border border-[#e3e4f8] bg-white p-5 shadow-sm">
      <div className="grid size-10 shrink-0 place-items-center rounded-full bg-[#ececf6]">
        {icon}
      </div>
      <div>
        <p className="text-2xl font-bold tracking-normal text-[#10111f]">
          {value}
        </p>
        <p className="mt-1 text-base font-medium text-[#676982]">{label}</p>
      </div>
    </article>
  );
}
