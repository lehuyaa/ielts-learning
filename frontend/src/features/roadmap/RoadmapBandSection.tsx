import { CheckCircle2, Lock } from "lucide-react";

import { RoadmapTopicCard } from "@/features/roadmap/RoadmapTopicCard";
import type {
  RoadmapBand,
  RoadmapLessonStatus,
} from "@/features/roadmap/types";
import { cn } from "@/lib/utils";

type RoadmapBandSectionProps = {
  band: RoadmapBand;
  isLast: boolean;
};

export function RoadmapBandSection({ band, isLast }: RoadmapBandSectionProps) {
  const isLocked = band.status === "locked";

  return (
    <section className="relative z-10 py-2">
      <BandNode
        band={band.band}
        status={bandStatusLabel(band.status)}
        statusVariant={band.status}
        topicCount={band.topicCount}
      />

      <div className="mx-auto mb-2 mt-4 grid max-w-2xl gap-3 md:grid-cols-2">
        {band.topics.map((topic) => (
          <RoadmapTopicCard isLocked={isLocked} key={topic.id} topic={topic} />
        ))}
      </div>

      {!isLast ? <ConnectorDots isActive={!isLocked} /> : null}
    </section>
  );
}

function bandStatusLabel(status: RoadmapLessonStatus) {
  switch (status) {
    case "completed":
      return "Complete";
    case "in-progress":
    case "unlocked":
      return "In Progress";
    case "locked":
    default:
      return "Locked";
  }
}

type BandNodeProps = {
  band: string;
  status: string;
  statusVariant: RoadmapLessonStatus;
  topicCount: number;
};

function BandNode({ band, status, statusVariant, topicCount }: BandNodeProps) {
  const isLocked = statusVariant === "locked";
  const isComplete = statusVariant === "completed";
  const color = getBandColor(band);

  return (
    <div className="flex justify-center">
      <div
        className={cn(
          "flex items-center gap-3 rounded-2xl border px-6 py-3",
          isLocked
            ? "border-border bg-muted/30"
            : `${color.border} ${color.bg} transition-shadow hover:shadow-md`,
        )}
      >
        {isLocked ? (
          <Lock className="size-5 text-muted-foreground" aria-hidden="true" />
        ) : isComplete ? (
          <CheckCircle2 className="size-5 text-emerald-500" aria-hidden="true" />
        ) : (
          <span
            className={cn("size-5 rounded-full bg-gradient-to-br", color.dot)}
            aria-hidden="true"
          />
        )}

        <div className="text-left">
          <h2
            className={cn(
              "text-sm font-bold tracking-normal",
              isLocked ? "text-muted-foreground" : "text-foreground",
            )}
          >
            Band {band}
          </h2>
          <p className="text-xs text-muted-foreground">
            {isLocked ? "Complete previous band" : `${topicCount} topics`}
          </p>
        </div>

        <span
          className={cn(
            "rounded-full px-2 py-0.5 text-xs font-semibold",
            isLocked ? "bg-muted text-muted-foreground" : color.badge,
          )}
        >
          {status}
        </span>
      </div>
    </div>
  );
}

function ConnectorDots({ isActive }: { isActive: boolean }) {
  return (
    <div className="flex justify-center">
      <div className="flex flex-col items-center gap-1 py-2">
        {[0, 1, 2].map((dot) => (
          <span
            className={cn(
              "size-1.5 rounded-full",
              isActive ? "bg-primary/40" : "bg-muted",
            )}
            key={dot}
          />
        ))}
      </div>
    </div>
  );
}

function getBandColor(band: string) {
  const numericBand = Number.parseFloat(band);

  if (numericBand < 6) {
    return {
      bg: "bg-slate-50",
      border: "border-slate-200",
      dot: "from-slate-400 to-slate-500",
      badge: "bg-slate-100 text-slate-600",
    };
  }

  if (numericBand < 7) {
    return {
      bg: "bg-blue-50",
      border: "border-blue-200",
      dot: "from-blue-400 to-indigo-500",
      badge: "bg-blue-100 text-blue-700",
    };
  }

  if (numericBand < 8) {
    return {
      bg: "bg-indigo-50",
      border: "border-indigo-200",
      dot: "from-indigo-500 to-violet-600",
      badge: "bg-indigo-100 text-indigo-700",
    };
  }

  return {
    bg: "bg-purple-50",
    border: "border-purple-200",
    dot: "from-violet-500 to-purple-700",
    badge: "bg-purple-100 text-purple-700",
  };
}

