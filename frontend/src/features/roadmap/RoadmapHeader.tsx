import {
  CheckCircle2,
  ChevronLeft,
  Flame,
  Star,
  Target,
  TrendingUp,
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
      <header className="sticky top-0 z-20 border-b border-border bg-white">
        <div className="mx-auto flex h-16 max-w-5xl items-center gap-4 px-6">
          <button
            aria-label="Go back"
            className="grid size-9 shrink-0 cursor-pointer place-items-center rounded-xl text-muted-foreground transition-colors hover:bg-muted/60"
            onClick={() => navigate("/dashboard")}
            type="button"
          >
            <ChevronLeft className="size-5" aria-hidden="true" />
          </button>

          <div className="grid size-8 shrink-0 place-items-center rounded-lg bg-primary text-white">
            <Target className="size-4" aria-hidden="true" />
          </div>

          <div className="min-w-0">
            <h1 className="truncate text-sm font-semibold tracking-normal text-foreground">
              {title}
            </h1>
            <p className="text-xs text-muted-foreground">{subtitle}</p>
          </div>

          <div className="flex-1" />

          <div className="hidden items-center gap-2 rounded-full border border-amber-200 bg-amber-50 px-3 py-1.5 text-xs font-semibold text-amber-700 md:flex">
            <Flame className="size-3.5 text-orange-500" aria-hidden="true" />
            {currentStreak} Day Streak
          </div>
        </div>
      </header>

      <section className="mx-auto max-w-5xl px-6 pt-8">
        <div className="grid gap-4 sm:grid-cols-3">
          <StatCard
            icon={
              <CheckCircle2
                className="size-4 text-emerald-500"
                aria-hidden="true"
              />
            }
            label="Topics Completed"
            value={`${topicsCompleted}/${totalTopics}`}
          />
          <StatCard
            icon={<Star className="size-4 text-amber-500" aria-hidden="true" />}
            label="Current Band"
            value={currentBand == null ? "-" : currentBand.toFixed(1)}
          />
          <StatCard
            icon={
              <TrendingUp className="size-4 text-primary" aria-hidden="true" />
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
    <article className="flex items-center gap-3 rounded-2xl border border-border bg-white p-4">
      <div className="grid size-8 shrink-0 place-items-center rounded-xl bg-muted">
        {icon}
      </div>
      <div>
        <p className="font-bold tracking-normal text-foreground">{value}</p>
        <p className="text-xs text-muted-foreground">{label}</p>
      </div>
    </article>
  );
}

