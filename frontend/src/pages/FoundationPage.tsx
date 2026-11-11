import { useQuery } from "@tanstack/react-query";
import {
  ArrowRight,
  AudioLines,
  BookOpen,
  CheckCircle2,
  Circle,
  Ear,
  Flag,
  Flame,
  GraduationCap,
  MessageCircle,
  PenLine,
  Puzzle,
  type LucideIcon,
} from "lucide-react";
import { Link } from "react-router-dom";

import { getApiErrorMessage } from "@/api/api";
import { getFoundationSummary } from "@/api/foundation";
import { ErrorState } from "@/components/state/ErrorState";
import { Button } from "@/components/ui/button";
import { Skeleton } from "@/components/ui/skeleton";
import { scoreTone } from "@/features/pronunciation/scoreColors";
import { cn } from "@/lib/utils";
import type { FoundationGoal, FoundationSkill } from "@/types/foundation";

const skillIcons: Record<FoundationSkill["key"], LucideIcon> = {
  VOCABULARY: BookOpen,
  GRAMMAR: Puzzle,
  SOUNDS: AudioLines,
  LISTENING: Ear,
  READING: BookOpen,
  WRITING: PenLine,
  SPEAKING: MessageCircle,
};

/** A1 dashboard (Task 76): checklist of spec §1, skills, lessons, streak. */
export function FoundationPage() {
  const query = useQuery({ queryKey: ["foundation", "A1"], queryFn: () => getFoundationSummary("A1") });
  const data = query.data;

  return (
    <div className="mx-auto max-w-5xl space-y-6 py-2">
      <header className="flex flex-wrap items-end justify-between gap-4">
        <div>
          <h1 className="flex items-center gap-2 text-2xl font-bold tracking-normal">
            <GraduationCap className="size-6 text-primary" /> {data?.title ?? "A1 Foundation"}
          </h1>
          <p className="mt-1 text-sm text-muted-foreground">Tiến độ của bạn tới trình độ A1: 5 mục tiêu, 7 kỹ năng.</p>
        </div>
        {data?.nextLesson ? (
          <Button asChild>
            <Link to={`/lessons/${data.nextLesson.id}/learn`}>
              Học tiếp: {data.nextLesson.title} <ArrowRight className="size-4" />
            </Link>
          </Button>
        ) : null}
      </header>

      {query.isLoading ? <Skeleton className="h-96 rounded-3xl" /> : null}
      {query.error ? (
        <ErrorState description={getApiErrorMessage(query.error, "Không tải được tiến độ.")} onRetry={() => void query.refetch()} title="Không tải được tiến độ A1" />
      ) : null}

      {data ? (
        <>
          <section className="grid gap-3 sm:grid-cols-3">
            <Stat icon={Flag} label="Bài tình huống" value={`${data.lessons.completed}/${data.lessons.target}`} note={`Hiện có ${data.lessons.total} bài`} />
            <Stat
              icon={CheckCircle2}
              label="Kiểm tra cuối tuần"
              value={`${data.lessons.checkpointsPassed}/${Math.max(data.lessons.checkpointsTotal, 4)}`}
              note="Qua bài kiểm tra để mở tuần sau"
            />
            <Stat icon={Flame} label="Chuỗi ngày học" value={`${data.streak.current} ngày`} note={`Dài nhất ${data.streak.longest} ngày`} />
          </section>

          <section className="rounded-3xl border bg-white p-5 shadow-sm sm:p-6">
            <div className="mb-4 flex flex-wrap items-baseline justify-between gap-2">
              <h2 className="text-lg font-bold">Đạt A1 khi bạn có thể…</h2>
              <span className="text-sm font-semibold text-muted-foreground">
                {data.goalsDone}/{data.checklist.length} mục tiêu
              </span>
            </div>
            <ul className="space-y-4">
              {data.checklist.map((goal) => (
                <GoalRow goal={goal} key={goal.key} />
              ))}
            </ul>
          </section>

          <section>
            <h2 className="mb-3 text-lg font-bold">Kỹ năng</h2>
            <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
              {data.skills.map((skill) => (
                <SkillCard key={skill.key} skill={skill} />
              ))}
            </div>
          </section>
        </>
      ) : null}
    </div>
  );
}

function Stat({ icon: Icon, label, value, note }: { icon: LucideIcon; label: string; value: string; note: string }) {
  return (
    <div className="flex items-center gap-3 rounded-3xl border bg-white p-4 shadow-sm">
      <span className="grid size-11 place-items-center rounded-2xl bg-primary/10 text-primary">
        <Icon className="size-5" />
      </span>
      <div>
        <p className="text-xs text-muted-foreground">{label}</p>
        <p className="text-xl font-bold">{value}</p>
        <p className="text-xs text-muted-foreground">{note}</p>
      </div>
    </div>
  );
}

function GoalRow({ goal }: { goal: FoundationGoal }) {
  return (
    <li className="flex gap-3">
      {goal.done ? <CheckCircle2 className="mt-0.5 size-5 shrink-0 text-emerald-600" /> : <Circle className="mt-0.5 size-5 shrink-0 text-muted-foreground" />}
      <div className="min-w-0 flex-1">
        <p className={cn("font-semibold", goal.done && "text-emerald-700")}>{goal.title}</p>
        <div className="mt-1.5 h-2 overflow-hidden rounded-full bg-muted">
          <div className={cn("h-full rounded-full", goal.done ? "bg-emerald-500" : "bg-primary")} style={{ width: `${goal.progress}%` }} />
        </div>
        <p className="mt-1 text-sm text-muted-foreground">
          {goal.detail}
          {!goal.done ? (
            <>
              {" · "}
              <Link className="font-semibold text-primary hover:underline" to={goal.link}>
                {goal.hint}
              </Link>
            </>
          ) : null}
        </p>
      </div>
    </li>
  );
}

function SkillCard({ skill }: { skill: FoundationSkill }) {
  const Icon = skillIcons[skill.key];
  const coverage = skill.total > 0 ? Math.round((skill.done * 100) / skill.total) : 0;
  return (
    <Link className="block rounded-3xl border bg-white p-4 shadow-sm transition-colors hover:border-primary/40" to={skill.link}>
      <div className="flex items-center justify-between">
        <span className="flex items-center gap-2 font-semibold">
          <Icon className="size-4 text-primary" /> {skill.label}
        </span>
        <span className={cn("text-2xl font-bold", skill.practised ? scoreTone(skill.value) : "text-muted-foreground")}>
          {skill.practised ? skill.value : "—"}
        </span>
      </div>
      <div className="mt-3 h-2 overflow-hidden rounded-full bg-muted" title={`${skill.done}/${skill.total}`}>
        <div className="h-full rounded-full bg-primary" style={{ width: `${coverage}%` }} />
      </div>
      <p className="mt-2 text-xs text-muted-foreground">{skill.detail}</p>
    </Link>
  );
}
