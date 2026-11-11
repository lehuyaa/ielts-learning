import { CheckCircle2, MessageCircle, Mic, Target } from "lucide-react";

import { getApiErrorMessage } from "@/api/api";
import { Button } from "@/components/ui/button";
import { AudioPlayer } from "@/features/audio/AudioPlayer";
import { usePassage } from "@/features/passage/hooks";
import type { ScenarioCard, SpeakingSummary } from "@/types/activity";

import { useStartSession } from "./hooks";

type SpeakingTaskProps = {
  activityId: number;
  lessonId: number;
  scenario: ScenarioCard;
  latest: SpeakingSummary | null;
  finishing: boolean;
  onFinish: () => void;
};

/**
 * SPEAKING activity with a role-play (Tasks 68-69): brief + model dialogue →
 * talk with the AI → back here to complete with the latest session's score.
 */
export function SpeakingTask({ activityId, lessonId, scenario, latest, finishing, onFinish }: SpeakingTaskProps) {
  const start = useStartSession();
  const model = usePassage(scenario.passageId ?? undefined);
  const returnTo = `/lessons/${lessonId}/learn?activity=${activityId}`;
  const done = latest?.taskCompleted ?? false;

  function begin(fresh: boolean) {
    start.mutate({ scenarioSlug: scenario.slug, activityId, fresh, returnTo });
  }

  return (
    <div className="space-y-5">
      <section className="rounded-3xl border bg-white p-5 shadow-sm">
        <p className="flex items-center gap-2 text-lg font-bold">
          <MessageCircle className="size-5 text-primary" /> {scenario.title}
        </p>
        <p className="mt-1 text-muted-foreground">{scenario.description}</p>
        <dl className="mt-4 space-y-2 text-sm">
          <div>
            <dt className="inline font-semibold">AI đóng vai: </dt>
            <dd className="inline">{scenario.role}</dd>
          </div>
          <div>
            <dt className="inline font-semibold">
              <Target className="mr-1 inline size-4 text-primary" />
              Nhiệm vụ:{" "}
            </dt>
            <dd className="inline">{scenario.task}</dd>
          </div>
        </dl>
        {scenario.targetPatterns.length > 0 ? (
          <div className="mt-4">
            <p className="mb-2 text-sm font-semibold">Mẫu câu nên dùng</p>
            <div className="flex flex-wrap gap-2">
              {scenario.targetPatterns.map((pattern) => (
                <span className="rounded-full bg-primary/5 px-3 py-1 text-sm font-medium text-primary" key={pattern}>
                  {pattern}
                </span>
              ))}
            </div>
          </div>
        ) : null}
        <p className="mt-4 text-xs text-muted-foreground">
          Khoảng {scenario.maxTurns} lượt nói. Bạn có thể nói qua micro hoặc gõ. Hoàn thành nhiệm vụ để qua bài.
        </p>
      </section>

      {model.data ? (
        <section>
          <p className="mb-2 text-sm font-semibold">Nghe hội thoại mẫu trước khi nói</p>
          <AudioPlayer passage={model.data} passed />
        </section>
      ) : null}

      {latest ? (
        <section className="flex flex-wrap items-center justify-between gap-3 rounded-3xl border bg-white p-4">
          <div>
            <p className="text-sm text-muted-foreground">Lần nói gần nhất</p>
            <p className="text-2xl font-bold">{latest.score ?? 0}/100</p>
          </div>
          <span
            className={
              done
                ? "inline-flex items-center gap-1 rounded-full bg-emerald-50 px-3 py-1 text-sm font-semibold text-emerald-700"
                : "rounded-full bg-amber-50 px-3 py-1 text-sm font-semibold text-amber-700"
            }
          >
            {done ? <CheckCircle2 className="size-4" /> : null}
            {done ? "Đã hoàn thành nhiệm vụ" : "Chưa hoàn thành nhiệm vụ"}
          </span>
        </section>
      ) : null}

      {start.error ? (
        <p className="text-center text-sm text-red-600">{getApiErrorMessage(start.error, "Không bắt đầu được cuộc hội thoại.")}</p>
      ) : null}

      <div className="flex flex-wrap justify-center gap-3">
        {latest ? (
          <>
            <Button disabled={start.isPending} onClick={() => begin(true)} variant={done ? "outline" : "default"}>
              <Mic className="size-4" /> Nói lại
            </Button>
            <Button disabled={finishing} onClick={onFinish} variant={done ? "default" : "outline"}>
              <CheckCircle2 className="size-4" /> Hoàn thành hoạt động
            </Button>
          </>
        ) : (
          <Button disabled={start.isPending} onClick={() => begin(false)} size="lg">
            <Mic className="size-4" /> {start.isPending ? "Đang mở…" : "Bắt đầu nói"}
          </Button>
        )}
      </div>
    </div>
  );
}
