import { cn } from "@/lib/utils";
import type { WritingReport } from "@/types/writing";

const labels: [keyof WritingReport["scores"], string][] = [
  ["taskAchievement", "Đúng yêu cầu đề"],
  ["grammar", "Ngữ pháp"],
  ["vocabulary", "Từ vựng"],
  ["coherence", "Mạch lạc"],
];

/** The four 1-10 writing scores as bars. */
export function ScoreBars({ scores, className }: { scores: WritingReport["scores"]; className?: string }) {
  return (
    <div className={cn("space-y-1.5", className)}>
      {labels.map(([key, label]) => (
        <div className="flex items-center gap-3 text-sm" key={key}>
          <span className="w-32 shrink-0 text-muted-foreground">{label}</span>
          <div className="h-2 flex-1 rounded-full bg-muted">
            <div
              className={cn("h-2 rounded-full", scores[key] >= 8 ? "bg-emerald-500" : scores[key] >= 6 ? "bg-amber-400" : "bg-red-400")}
              style={{ width: `${scores[key] * 10}%` }}
            />
          </div>
          <span className="w-8 text-right font-semibold">{scores[key]}</span>
        </div>
      ))}
    </div>
  );
}
