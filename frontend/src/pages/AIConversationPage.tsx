import { Clock, MessageSquare, Mic, Plus, Star, Target } from "lucide-react";
import { useMemo } from "react";
import { Link } from "react-router-dom";

import { getApiErrorMessage } from "@/api/api";
import { CardSkeleton } from "@/components/state/CardSkeleton";
import { ErrorState } from "@/components/state/ErrorState";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { useScenarios } from "@/features/aiConversation/hooks/useScenarios";
import { useStartSession } from "@/features/speaking/hooks";
import { cn } from "@/lib/utils";
import type { ScenarioResponse } from "@/types/aiConversation";

const levelBadge: Record<string, string> = {
  A1: "border-emerald-200 bg-emerald-50 text-emerald-700",
  A2: "border-teal-200 bg-teal-50 text-teal-700",
  B1: "border-amber-200 bg-amber-50 text-amber-700",
  B2: "border-red-200 bg-red-50 text-red-700",
};

const customLevel: Record<string, string> = { BEGINNER: "A2", INTERMEDIATE: "B1", ADVANCED: "B2" };

function levelOf(scenario: ScenarioResponse) {
  return scenario.cefrLevel || customLevel[scenario.level] || "A2";
}

/** Role-play library (Task 68): system scenarios by week, then the learner's own. */
export function AIConversationPage() {
  const scenariosQuery = useScenarios();
  const items = useMemo(() => scenariosQuery.data?.items ?? [], [scenariosQuery.data]);

  const groups = useMemo(() => {
    const byGroup = new Map<string, ScenarioResponse[]>();
    for (const scenario of items) {
      const key = scenario.isSystem ? scenario.group || "Tình huống mẫu" : "Tình huống của bạn";
      byGroup.set(key, [...(byGroup.get(key) ?? []), scenario]);
    }
    return [...byGroup.entries()];
  }, [items]);

  const stats = useMemo(() => {
    const played = items.filter((scenario) => scenario.sessions > 0);
    const sessions = played.reduce((sum, scenario) => sum + scenario.sessions, 0);
    const best = played.map((scenario) => scenario.bestScore ?? 0);
    const average = best.length > 0 ? Math.round(best.reduce((a, b) => a + b, 0) / best.length) : null;
    return [
      { label: "Lượt hội thoại", value: sessions, icon: MessageSquare, color: "text-primary" },
      { label: "Tình huống đã luyện", value: `${played.length}/${items.length}`, icon: Target, color: "text-purple-500" },
      { label: "Điểm cao nhất TB", value: average ?? "—", icon: Star, color: "fill-amber-400 text-amber-400" },
    ];
  }, [items]);

  return (
    <div className="space-y-6">
      <section className="flex flex-wrap items-end justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold tracking-normal text-foreground">Nói chuyện với AI</h1>
          <p className="mt-0.5 text-sm text-muted-foreground">
            Chọn một tình huống, nhận vai và hoàn thành nhiệm vụ bằng tiếng Anh. Kết thúc để nhận báo cáo và góp ý.
          </p>
        </div>
        <Button asChild className="gap-1.5 rounded-full" variant="outline">
          <Link to="/ai-conversation/create">
            <Plus className="size-4" aria-hidden="true" />
            Tự tạo tình huống
          </Link>
        </Button>
      </section>

      <section className="flex flex-wrap gap-3">
        {stats.map((stat) => (
          <Badge className="gap-2 border-border bg-white px-4 py-2 shadow-sm" key={stat.label} variant="outline">
            <stat.icon className={cn("size-4", stat.color)} />
            <span className="text-sm font-normal text-muted-foreground">{stat.label}</span>
            <span className="text-sm font-bold text-foreground">{stat.value}</span>
          </Badge>
        ))}
      </section>

      {scenariosQuery.isLoading ? (
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {Array.from({ length: 6 }).map((_, index) => (
            <CardSkeleton key={index} lines={2} showIcon />
          ))}
        </div>
      ) : scenariosQuery.error ? (
        <ErrorState
          description={getApiErrorMessage(scenariosQuery.error, "Không tải được danh sách tình huống.")}
          onRetry={() => void scenariosQuery.refetch()}
          title="Không tải được tình huống"
        />
      ) : groups.length > 0 ? (
        groups.map(([group, scenarios]) => (
          <section className="space-y-3" key={group}>
            <h2 className="text-lg font-bold text-foreground">{group}</h2>
            <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
              {scenarios.map((scenario) => (
                <ScenarioTile key={scenario.slug} scenario={scenario} />
              ))}
            </div>
          </section>
        ))
      ) : (
        <Link
          className="flex w-full flex-col items-center gap-3 rounded-2xl border-2 border-dashed border-border py-12 text-center transition-colors hover:border-primary/40"
          to="/ai-conversation/create"
        >
          <Plus className="size-5 text-muted-foreground" aria-hidden="true" />
          <span className="text-base font-bold text-foreground">Tạo tình huống đầu tiên</span>
          <span className="text-sm text-muted-foreground">Chưa có tình huống mẫu — admin cần import a1-scenarios.json.</span>
        </Link>
      )}
    </div>
  );
}

function ScenarioTile({ scenario }: { scenario: ScenarioResponse }) {
  const start = useStartSession();
  const level = levelOf(scenario);
  return (
    <button
      className="flex cursor-pointer flex-col rounded-2xl border bg-white p-5 text-left shadow-sm transition-colors hover:border-primary/40 hover:shadow-md disabled:cursor-wait disabled:opacity-60"
      disabled={start.isPending}
      onClick={() => start.mutate({ scenarioSlug: scenario.slug, returnTo: "/ai-conversation" })}
      type="button"
    >
      <div className="flex w-full items-start justify-between">
        <span className="flex size-10 items-center justify-center rounded-xl bg-indigo-50">
          <Mic className="size-[18px] text-primary" />
        </span>
        <Badge className={levelBadge[level] ?? levelBadge.A2} variant="outline">
          {level}
        </Badge>
      </div>
      <h3 className="mt-4 text-base font-bold text-foreground">{scenario.title}</h3>
      <p className="mt-1.5 line-clamp-2 text-sm text-muted-foreground">{scenario.task || scenario.description}</p>
      <div className="mt-4 flex w-full items-center justify-between text-xs text-muted-foreground">
        <span className="flex items-center gap-1.5">
          <Clock className="size-3.5" />
          {scenario.duration || `${scenario.maxTurns} lượt`}
        </span>
        {scenario.sessions > 0 ? (
          <span className="font-semibold text-foreground">
            Cao nhất {scenario.bestScore ?? 0} · {scenario.sessions} lần
          </span>
        ) : (
          <span>Chưa luyện</span>
        )}
      </div>
      {start.error ? <p className="mt-2 text-xs text-red-600">{getApiErrorMessage(start.error, "Không mở được.")}</p> : null}
    </button>
  );
}
