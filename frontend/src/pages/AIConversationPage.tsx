import { Clock, MessageSquare, Plus, Star } from "lucide-react";
import { useSyncExternalStore } from "react";
import { Link } from "react-router-dom";

import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardDescription, CardTitle } from "@/components/ui/card";
import {
  getCustomScenarios,
  scenarios,
  subscribeToCustomScenarios,
  type Scenario,
  type ScenarioLevel,
} from "@/features/aiConversation/scenarios";
import { cn } from "@/lib/utils";

const stats = [
  {
    label: "Conversations this week",
    value: 4,
    icon: MessageSquare,
    iconColor: "text-primary",
  },
  {
    label: "Avg. score",
    value: 79,
    icon: Star,
    iconColor: "fill-amber-400 text-amber-400",
  },
  {
    label: "Minutes practiced",
    value: 38,
    icon: Clock,
    iconColor: "text-purple-500",
  },
];

const levelBadgeClasses: Record<ScenarioLevel, string> = {
  Beginner: "border-emerald-200 bg-emerald-50 text-emerald-700",
  Intermediate: "border-amber-200 bg-amber-50 text-amber-700",
  Advanced: "border-red-200 bg-red-50 text-red-700",
};

export function AIConversationPage() {
  const myScenarios = useSyncExternalStore(
    subscribeToCustomScenarios,
    getCustomScenarios,
  );

  return (
    <div className="space-y-6">
      <section>
        <h1 className="text-2xl font-bold tracking-normal text-foreground">
          AI Conversation
        </h1>
        <p className="mt-0.5 text-sm text-muted-foreground">
          Practice real English conversations with your AI partner.
        </p>
      </section>

      <section className="flex flex-wrap gap-3">
        {stats.map((stat) => (
          <Badge
            className="gap-2 border-border bg-white px-4 py-2 shadow-sm"
            key={stat.label}
            variant="outline"
          >
            <stat.icon className={cn("size-4", stat.iconColor)} />
            <span className="text-sm font-normal text-muted-foreground">
              {stat.label}
            </span>
            <span className="text-sm font-bold text-foreground">
              {stat.value}
            </span>
          </Badge>
        ))}
      </section>

      <section className="space-y-3">
        <h2 className="text-sm font-semibold text-foreground">
          Choose a scenario
        </h2>

        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {scenarios.map((scenario) => (
            <ScenarioCard key={scenario.slug} scenario={scenario} />
          ))}
        </div>
      </section>

      <section className="space-y-3">
        <div className="flex items-center justify-between gap-4">
          <div>
            <h2 className="text-lg font-bold text-foreground">My scenarios</h2>
            <p className="text-sm text-muted-foreground">
              Scenarios you created
            </p>
          </div>
          <Button asChild className="gap-1.5 rounded-full">
            <Link to="/ai-conversation/create">
              <Plus className="size-4" aria-hidden="true" />
              Create
            </Link>
          </Button>
        </div>

        {myScenarios.length > 0 ? (
          <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
            {myScenarios.map((scenario) => (
              <ScenarioCard key={scenario.slug} scenario={scenario} />
            ))}
          </div>
        ) : (
          <Link
            className="flex w-full flex-col items-center gap-3 rounded-2xl border-2 border-dashed border-border py-12 text-center transition-colors hover:border-primary/40"
            to="/ai-conversation/create"
          >
            <span className="flex size-12 items-center justify-center rounded-2xl bg-muted text-muted-foreground">
              <Plus className="size-5" aria-hidden="true" />
            </span>
            <span>
              <span className="block text-base font-bold text-foreground">
                Create your first scenario
              </span>
              <span className="mt-1 block text-sm text-muted-foreground">
                Add a custom topic or situation to practice
              </span>
            </span>
          </Link>
        )}
      </section>
    </div>
  );
}

function ScenarioCard({ scenario }: { scenario: Scenario }) {
  return (
    <Card
      as={Link}
      className="flex flex-col p-5 text-left transition-colors hover:border-primary/40 hover:shadow-md"
      to={`/ai-conversation/${scenario.slug}`}
    >
      <div className="flex items-start justify-between">
        <div className="flex size-10 items-center justify-center rounded-xl bg-indigo-50">
          <scenario.icon className="size-[18px] text-primary" />
        </div>
        <Badge className={levelBadgeClasses[scenario.level]} variant="outline">
          {scenario.level}
        </Badge>
      </div>

      <CardTitle as="h3" className="mt-4 text-base font-bold text-foreground">
        {scenario.title}
      </CardTitle>
      <CardDescription className="mt-1.5">
        {scenario.description}
      </CardDescription>

      <div className="mt-4 flex items-center gap-1.5 text-xs text-muted-foreground">
        <Clock className="size-3.5" />
        {scenario.duration}
      </div>
    </Card>
  );
}
