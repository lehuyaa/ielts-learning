import { ArrowLeft, ChevronLeft, ChevronRight, Ear, GraduationCap, Mic, TriangleAlert, Volume2 } from "lucide-react";
import { useState, type ReactNode } from "react";
import { Link, useParams } from "react-router-dom";

import { APIError } from "@/api/api";
import { ErrorState } from "@/components/state/ErrorState";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Skeleton } from "@/components/ui/skeleton";
import { IpaText } from "@/features/pronunciation/IpaText";
import { PairDrill } from "@/features/pronunciation/PairDrill";
import { PronunciationRecorder } from "@/features/pronunciation/PronunciationRecorder";
import { scoreTone } from "@/features/pronunciation/scoreColors";
import { usePronunciationStatus } from "@/features/pronunciation/hooks/usePronunciation";
import { cn } from "@/lib/utils";
import { usePhoneme } from "@/features/pronunciation/hooks/usePhonemes";
import { phonemePath } from "@/lib/ipa";
import { playPronunciation } from "@/lib/pronunciation";
import type { MinimalPair, PhonemeCard } from "@/types/phoneme";

const kindLabel: Record<PhonemeCard["kind"], string> = {
  VOWEL: "Nguyên âm",
  DIPHTHONG: "Nguyên âm đôi",
  CONSONANT: "Phụ âm",
};

export function PronunciationDetailPage() {
  const { symbol = "" } = useParams();
  const query = usePhoneme(symbol);
  const [drill, setDrill] = useState(0);
  const pronunciation = usePronunciationStatus();

  const data = query.data;
  const phoneme = data?.phoneme;

  return (
    <div className="mx-auto max-w-3xl space-y-5 py-2" key={symbol}>
      <div className="flex items-center justify-between gap-3">
        <Link className="inline-flex items-center gap-1 text-sm font-semibold text-muted-foreground hover:text-foreground" to="/pronunciation">
          <ArrowLeft className="size-4" /> Bảng IPA
        </Link>
        {data ? (
          <div className="flex gap-1">
            <NeighbourLink direction="previous" symbol={data.previous} />
            <NeighbourLink direction="next" symbol={data.next} />
          </div>
        ) : null}
      </div>

      {query.isLoading ? <Skeleton className="h-96 rounded-3xl" /> : null}
      {query.error ? (
        <ErrorState
          description={query.error instanceof APIError ? query.error.message : "Không tải được thẻ âm."}
          onRetry={() => void query.refetch()}
          title="Không tải được thẻ âm"
        />
      ) : null}

      {phoneme && data ? (
        <>
          <section className="flex flex-wrap items-center gap-5 rounded-3xl border bg-white p-6 shadow-sm">
            <button
              aria-label="Nghe ví dụ"
              className="flex size-24 shrink-0 cursor-pointer flex-col items-center justify-center rounded-3xl bg-primary/10 font-mono text-4xl font-bold text-primary hover:bg-primary/20"
              onClick={() =>
                playPronunciation(phoneme.exampleWords.map((word) => word.word).join(", "), phoneme.audioUrl)
              }
              type="button"
            >
              /{phoneme.symbol}/
              <Volume2 className="mt-1 size-4" />
            </button>
            <div className="min-w-0 flex-1">
              <h1 className="text-2xl font-bold tracking-normal">Âm {phoneme.nameVi}</h1>
              <div className="mt-2 flex flex-wrap gap-2">
                <Badge variant="outline">{kindLabel[phoneme.kind]}</Badge>
                {phoneme.voiced !== null ? (
                  <Badge variant="outline">{phoneme.voiced ? "Hữu thanh — rung cổ" : "Vô thanh — không rung cổ"}</Badge>
                ) : null}
              </div>
            </div>
          </section>

          <Section icon={<GraduationCap className="size-4" />} title="Khẩu hình">
            <p className="leading-relaxed">{phoneme.mouthVi}</p>
          </Section>

          <Section icon={<Volume2 className="size-4" />} title="Từ ví dụ">
            <div className="grid gap-2 sm:grid-cols-3">
              {phoneme.exampleWords.map((word) => (
                <button
                  className="flex cursor-pointer items-center justify-between gap-2 rounded-2xl border px-4 py-3 text-left hover:border-primary/40"
                  key={word.word}
                  onClick={() => playPronunciation(word.word, word.audioUrl)}
                  type="button"
                >
                  <span>
                    <span className="block font-semibold">{word.word}</span>
                    <IpaText className="text-sm text-muted-foreground" ipa={word.ipa} linked={false} />
                  </span>
                  <Volume2 className="size-4 shrink-0 text-primary" />
                </button>
              ))}
            </div>
          </Section>

          {pronunciation.data?.enabled && phoneme.exampleWords.length > 0 ? (
            <Section icon={<Mic className="size-4" />} title="Luyện nói">
              <p className="mb-3 text-sm text-muted-foreground">
                {data.progress.practised ? (
                  <>
                    Độ chính xác âm /{phoneme.symbol}/ của bạn:{" "}
                    <span className={cn("font-bold", scoreTone(data.progress.accuracy))}>{data.progress.accuracy}/100</span> ({data.progress.attempts} lần).
                  </>
                ) : (
                  <>Nói từng từ, AI chấm từng âm. Kết quả cập nhật độ chính xác của âm này.</>
                )}
              </p>
              <ul className="divide-y">
                {phoneme.exampleWords.map((word) => (
                  <li className="flex flex-wrap items-start gap-3 py-3" key={word.word}>
                    <button
                      className="flex w-32 shrink-0 cursor-pointer items-center gap-2 text-left hover:text-primary"
                      onClick={() => playPronunciation(word.word, word.audioUrl)}
                      type="button"
                    >
                      <Volume2 className="size-4 text-primary" />
                      <span>
                        <span className="block font-semibold">{word.word}</span>
                        <IpaText className="text-xs text-muted-foreground" ipa={word.ipa} linked={false} />
                      </span>
                    </button>
                    <div className="min-w-0 flex-1">
                      <PronunciationRecorder compact referenceText={word.word} source="WORD" />
                    </div>
                  </li>
                ))}
              </ul>
            </Section>
          ) : null}

          {phoneme.minimalPairs.length > 0 ? (
            <Section
              action={
                drill === 0 ? (
                  <Button onClick={() => setDrill(1)} size="sm">
                    <Ear className="size-4" /> Luyện nghe
                  </Button>
                ) : null
              }
              icon={<Ear className="size-4" />}
              title={`Cặp từ dễ nhầm (${phoneme.minimalPairs.length})`}
            >
              {drill > 0 ? (
                <div className="mb-4">
                  <PairDrill key={drill} onClose={() => setDrill(0)} pairs={phoneme.minimalPairs} symbol={phoneme.symbol} />
                </div>
              ) : null}
              <PairList pairs={phoneme.minimalPairs} symbol={phoneme.symbol} />
            </Section>
          ) : null}

          {phoneme.mistakesVi ? (
            <Section icon={<TriangleAlert className="size-4" />} title="Lỗi người Việt hay mắc">
              <p className="leading-relaxed">{phoneme.mistakesVi}</p>
            </Section>
          ) : null}

          {data.lessons.length > 0 ? (
            <Section icon={<GraduationCap className="size-4" />} title="Bài học có âm này">
              <ul className="divide-y">
                {data.lessons.map((lesson) => (
                  <li key={lesson.activityId}>
                    <Link
                      className="flex items-center justify-between gap-3 py-2.5 hover:text-primary"
                      to={`/lessons/${lesson.lessonId}/learn?activity=${lesson.activityId}`}
                    >
                      <span>
                        <span className="block font-semibold">{lesson.activityTitle}</span>
                        <span className="text-xs text-muted-foreground">
                          {lesson.topicTitle} · {lesson.lessonTitle}
                        </span>
                      </span>
                      <ChevronRight className="size-4 shrink-0" />
                    </Link>
                  </li>
                ))}
              </ul>
            </Section>
          ) : null}
        </>
      ) : null}
    </div>
  );
}

function Section({
  title,
  icon,
  action,
  children,
}: {
  title: string;
  icon: ReactNode;
  action?: ReactNode;
  children: ReactNode;
}) {
  return (
    <section className="rounded-3xl border bg-white p-5">
      <div className="mb-3 flex items-center justify-between gap-3">
        <h2 className="flex items-center gap-2 font-bold">
          <span className="text-primary">{icon}</span>
          {title}
        </h2>
        {action}
      </div>
      {children}
    </section>
  );
}

function NeighbourLink({ symbol, direction }: { symbol: string | null; direction: "previous" | "next" }) {
  if (!symbol) {
    return null;
  }
  return (
    <Link
      className="inline-flex items-center gap-1 rounded-xl border bg-white px-3 py-1.5 font-mono text-sm font-semibold hover:border-primary/40"
      to={phonemePath(symbol)}
    >
      {direction === "previous" ? <ChevronLeft className="size-4" /> : null}/{symbol}/
      {direction === "next" ? <ChevronRight className="size-4" /> : null}
    </Link>
  );
}

/** Pairs grouped by the contrasting sound, each word playable. */
function PairList({ pairs, symbol }: { pairs: MinimalPair[]; symbol: string }) {
  const groups = new Map<string, MinimalPair[]>();
  for (const pair of pairs) {
    groups.set(pair.contrast, [...(groups.get(pair.contrast) ?? []), pair]);
  }

  return (
    <div className="space-y-4">
      {[...groups.entries()].map(([contrast, items]) => (
        <div key={contrast}>
          <p className="mb-2 text-sm font-semibold text-muted-foreground">
            /{symbol}/ và{" "}
            {contrast ? (
              <Link className="font-mono text-primary hover:underline" to={phonemePath(contrast)}>
                /{contrast}/
              </Link>
            ) : (
              "không có âm này"
            )}
          </p>
          <div className="grid gap-2 sm:grid-cols-2">
            {items.map((pair) => (
              <div className="flex items-center gap-2 rounded-2xl bg-muted/40 p-2" key={`${pair.a}-${pair.b}`}>
                <WordButton audioUrl={pair.aAudioUrl} ipa={pair.aIpa} word={pair.a} />
                <span className="text-xs text-muted-foreground">vs</span>
                <WordButton audioUrl={pair.bAudioUrl} ipa={pair.bIpa} word={pair.b} />
              </div>
            ))}
          </div>
        </div>
      ))}
    </div>
  );
}

function WordButton({ word, ipa, audioUrl }: { word: string; ipa: string; audioUrl?: string }) {
  return (
    <button
      className="flex flex-1 cursor-pointer flex-col items-center rounded-xl bg-white px-2 py-1.5 hover:ring-2 hover:ring-primary/30"
      onClick={() => playPronunciation(word, audioUrl)}
      type="button"
    >
      <span className="font-semibold">{word}</span>
      <span className="font-mono text-xs text-muted-foreground">{ipa}</span>
    </button>
  );
}
