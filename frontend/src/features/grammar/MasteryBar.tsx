import { cn } from "@/lib/utils";

function masteryTone(mastery: number, practised: boolean) {
  if (!practised) return { bar: "bg-muted", text: "text-muted-foreground", label: "Chưa luyện" };
  if (mastery >= 80) return { bar: "bg-emerald-500", text: "text-emerald-700", label: "Vững" };
  if (mastery >= 50) return { bar: "bg-amber-400", text: "text-amber-700", label: "Đang tiến bộ" };
  return { bar: "bg-red-400", text: "text-red-700", label: "Cần ôn" };
}

export function MasteryBar({ mastery, practised, className }: { mastery: number; practised: boolean; className?: string }) {
  const tone = masteryTone(mastery, practised);
  return (
    <div className={className}>
      <div className="mb-1 flex items-center justify-between text-xs">
        <span className={cn("font-semibold", tone.text)}>{tone.label}</span>
        <span className="text-muted-foreground">{practised ? `${mastery}%` : "—"}</span>
      </div>
      <div className="h-2 overflow-hidden rounded-full bg-muted">
        <div className={cn("h-full rounded-full transition-all", tone.bar)} style={{ width: `${practised ? Math.max(mastery, 3) : 0}%` }} />
      </div>
    </div>
  );
}
