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

type RoadmapHeaderProps = {
  title: string;
  subtitle: string;
  topicsCompleted: number;
  totalTopics: number;
  currentBand: number | null;
  wordsMastered: number;
  currentStreak: number;
};

export function RoadmapHeader({
  title,
  subtitle,
  topicsCompleted,
  totalTopics,
  currentBand,
  wordsMastered,
  currentStreak,
}: RoadmapHeaderProps) {
  const navigate = useNavigate();

  return (
    <>
      <header className="border-b border-[#e6e6f3] bg-white">
        <div className="mx-auto flex min-h-18 max-w-5xl items-center justify-between gap-5 px-4 lg:px-0">
          <div className="flex min-w-0 items-center gap-4">
            <button
              aria-label="Go back"
              className="grid size-9 shrink-0 cursor-pointer place-items-center rounded-full text-[#6d7088] transition-colors hover:bg-[#f0f1fb]"
              onClick={() => navigate("/dashboard")}
              type="button"
            >
              <ArrowLeft className="size-5" aria-hidden="true" />
            </button>

            <div className="grid size-9 shrink-0 place-items-center rounded-full bg-primary text-white">
              <Target className="size-5" aria-hidden="true" />
            </div>

            <div className="min-w-0">
              <h1 className="truncate text-xl font-bold tracking-normal text-[#10111f]">
                {title}
              </h1>
              <p className="mt-1 text-base font-medium text-[#676982]">
                {subtitle}
              </p>
            </div>
          </div>

          <div className="hidden rounded-full border border-[#ffd75f] bg-[#fffdf5] px-4 py-2 text-sm font-bold text-[#b84a00] md:flex md:items-center md:gap-2">
            <Flame className="size-4 text-[#ff6b00]" aria-hidden="true" />
            {currentStreak} Day Streak
          </div>
        </div>
      </header>

      <section className="mx-auto max-w-5xl px-4 pt-8 lg:px-0">
        <div className="grid gap-4 lg:grid-cols-3">
          <StatCard
            icon={
              <CheckCircle2
                className="size-5 text-success"
                aria-hidden="true"
              />
            }
            label="Topics Completed"
            value={`${topicsCompleted}/${totalTopics}`}
          />
          <StatCard
            icon={<Star className="size-5 text-warning" aria-hidden="true" />}
            label="Current Band"
            value={currentBand == null ? "-" : currentBand.toFixed(1)}
          />
          <StatCard
            icon={
              <TrendingUp className="size-5 text-primary" aria-hidden="true" />
            }
            label="Words Mastered"
            value={String(wordsMastered)}
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
