import { ArrowRight, AudioLines, Mic, Volume2 } from "lucide-react";
import { useState, type ReactNode } from "react";
import { Link } from "react-router-dom";

import { APIError } from "@/api/api";
import { EmptyState } from "@/components/state/EmptyState";
import { ErrorState } from "@/components/state/ErrorState";
import { Skeleton } from "@/components/ui/skeleton";
import { IpaText } from "@/features/pronunciation/IpaText";
import { usePhonemes } from "@/features/pronunciation/hooks/usePhonemes";
import { scoreBg, scoreTone } from "@/features/pronunciation/scoreColors";
import { phonemePath } from "@/lib/ipa";
import { playPronunciation } from "@/lib/pronunciation";
import { cn } from "@/lib/utils";
import type { PhonemeItem } from "@/types/phoneme";

/** Voiceless / voiced consonant pairs, shown side by side on the chart. */
const CONSONANT_PAIRS: [string, string][] = [
  ["p", "b"], ["t", "d"], ["k", "ɡ"], ["f", "v"], ["θ", "ð"], ["s", "z"], ["ʃ", "ʒ"], ["tʃ", "dʒ"],
];

export function PronunciationPage() {
  const query = usePhonemes();
  const items = query.data?.items ?? [];
  const bySymbol = new Map(items.map((item) => [item.symbol, item]));
  const [selected, setSelected] = useState<string | null>(null);
  const current = selected ? bySymbol.get(selected) : undefined;

  function pick(item: PhonemeItem) {
    setSelected(item.symbol);
    playPronunciation(item.exampleWord?.word ?? item.symbol, item.audioUrl || item.exampleWord?.audioUrl);
  }

  const paired = new Set(CONSONANT_PAIRS.flat());
  const vowels = items.filter((item) => item.kind === "VOWEL");
  const diphthongs = items.filter((item) => item.kind === "DIPHTHONG");
  const otherConsonants = items.filter((item) => item.kind === "CONSONANT" && !paired.has(item.symbol));
  // Weakest sounds (Task 72): practised enough to trust, lowest accuracy first.
  const weakest = items
    .filter((item) => item.progress.practised && item.progress.attempts >= 3 && item.progress.accuracy < 80)
    .sort((a, b) => a.progress.accuracy - b.progress.accuracy)
    .slice(0, 6);

  return (
    <div className="mx-auto max-w-5xl space-y-6 py-2">
      <header>
        <h1 className="flex items-center gap-2 text-2xl font-bold tracking-normal">
          <AudioLines className="size-6 text-primary" /> Phát âm (IPA)
        </h1>
        <p className="mt-1 text-sm text-muted-foreground">
          40 âm của tiếng Anh Mỹ. Chạm vào một âm để nghe ví dụ, rồi mở thẻ âm để xem khẩu hình và luyện các cặp từ dễ nhầm.
        </p>
      </header>

      {query.isLoading ? <Skeleton className="h-96 rounded-3xl" /> : null}
      {query.error ? (
        <ErrorState
          description={query.error instanceof APIError ? query.error.message : "Không tải được bảng âm."}
          onRetry={() => void query.refetch()}
          title="Không tải được bảng IPA"
        />
      ) : null}
      {query.data && items.length === 0 ? (
        <EmptyState
          actionHref="/admin/content/import"
          actionLabel="Import nội dung"
          description="Hãy import file content/a1/a1-sounds.json ở trang Import Lessons."
          title="Chưa có thẻ âm nào"
        />
      ) : null}

      {weakest.length > 0 ? (
        <section className="rounded-3xl border border-amber-200 bg-amber-50/50 p-4 sm:p-5">
          <h2 className="flex items-center gap-2 font-bold">
            <Mic className="size-4 text-amber-600" /> Âm bạn cần luyện
          </h2>
          <p className="mb-3 text-xs text-muted-foreground">Từ các lần bạn nói được chấm điểm (trung bình ~30 lần gần nhất).</p>
          <div className="grid gap-2 sm:grid-cols-2 lg:grid-cols-3">
            {weakest.map((item) => (
              <Link
                className={cn("flex items-center gap-3 rounded-2xl border bg-white p-3 hover:shadow-sm", scoreBg(item.progress.accuracy))}
                key={item.symbol}
                to={phonemePath(item.symbol)}
              >
                <span className="font-mono text-2xl font-bold">/{item.symbol}/</span>
                <span className="min-w-0 flex-1">
                  <span className={cn("block font-bold", scoreTone(item.progress.accuracy))}>{item.progress.accuracy}/100</span>
                  <span className="text-xs text-muted-foreground">
                    {item.progress.attempts} lần{item.exampleWord ? ` · ${item.exampleWord.word}` : ""}
                  </span>
                </span>
                <span className="inline-flex items-center gap-1 rounded-xl bg-primary px-3 py-1.5 text-sm font-semibold text-primary-foreground">
                  Luyện <ArrowRight className="size-3.5" />
                </span>
              </Link>
            ))}
          </div>
        </section>
      ) : null}

      {items.length > 0 ? (
        <>
          <ChartSection title="Nguyên âm đơn" subtitle="Dấu ː nghĩa là âm dài">
            <div className="grid grid-cols-3 gap-2 sm:grid-cols-4 lg:grid-cols-6">
              {vowels.map((item) => (
                <Cell item={item} key={item.symbol} onPick={pick} selected={selected === item.symbol} />
              ))}
            </div>
          </ChartSection>

          <ChartSection title="Nguyên âm đôi" subtitle="Trượt từ âm này sang âm kia">
            <div className="grid grid-cols-3 gap-2 sm:grid-cols-5">
              {diphthongs.map((item) => (
                <Cell item={item} key={item.symbol} onPick={pick} selected={selected === item.symbol} />
              ))}
            </div>
          </ChartSection>

          <ChartSection title="Phụ âm" subtitle="Cặp vô thanh (không rung cổ) – hữu thanh (rung cổ)">
            <div className="grid grid-cols-2 gap-x-4 gap-y-2 sm:grid-cols-4">
              {CONSONANT_PAIRS.map(([voiceless, voiced]) => (
                <div className="grid grid-cols-2 gap-1.5 rounded-2xl bg-muted/40 p-1.5" key={voiceless}>
                  {[voiceless, voiced].map((symbol) => {
                    const item = bySymbol.get(symbol);
                    return item ? (
                      <Cell compact item={item} key={symbol} onPick={pick} selected={selected === symbol} />
                    ) : null;
                  })}
                </div>
              ))}
            </div>
            <p className="mb-2 mt-4 text-xs font-semibold uppercase text-muted-foreground">Phụ âm khác</p>
            <div className="grid grid-cols-4 gap-2 sm:grid-cols-8">
              {otherConsonants.map((item) => (
                <Cell compact item={item} key={item.symbol} onPick={pick} selected={selected === item.symbol} />
              ))}
            </div>
          </ChartSection>

          {current ? <Preview item={current} /> : null}

          {(query.data?.practiceTopics ?? []).map((topic) => (
            <Link
              className="flex items-center justify-between gap-4 rounded-3xl border bg-gradient-to-r from-primary/10 to-white p-5 transition-colors hover:border-primary/40"
              key={topic.id}
              to={`/topics/${topic.id}`}
            >
              <div>
                <p className="font-bold">{topic.title}</p>
                <p className="text-sm text-muted-foreground">
                  {topic.description} · {topic.lessonCount} bài
                </p>
              </div>
              <ArrowRight className="size-5 shrink-0 text-primary" />
            </Link>
          ))}
        </>
      ) : null}
    </div>
  );
}

function ChartSection({ title, subtitle, children }: { title: string; subtitle: string; children: ReactNode }) {
  return (
    <section className="rounded-3xl border bg-white p-4 sm:p-5">
      <div className="mb-3 flex flex-wrap items-baseline justify-between gap-2">
        <h2 className="font-bold">{title}</h2>
        <p className="text-xs text-muted-foreground">{subtitle}</p>
      </div>
      {children}
    </section>
  );
}

function Cell({
  item,
  selected,
  onPick,
  compact = false,
}: {
  item: PhonemeItem;
  selected: boolean;
  onPick: (item: PhonemeItem) => void;
  compact?: boolean;
}) {
  return (
    <button
      aria-label={`Âm /${item.symbol}/ như trong ${item.exampleWord?.word ?? ""}`}
      className={cn(
        "flex cursor-pointer flex-col items-center justify-center rounded-xl border-2 bg-white transition-colors",
        compact ? "h-16 px-1" : "h-20 px-2",
        selected ? "border-primary bg-primary/5" : "border-transparent shadow-sm hover:border-primary/40",
      )}
      onClick={() => onPick(item)}
      type="button"
    >
      <span className={cn("font-mono font-bold", compact ? "text-xl" : "text-2xl")}>{item.symbol}</span>
      {item.exampleWord ? (
        <span className="truncate text-xs text-muted-foreground">{item.exampleWord.word}</span>
      ) : null}
    </button>
  );
}

function Preview({ item }: { item: PhonemeItem }) {
  return (
    <div className="sticky bottom-4 z-10 flex flex-wrap items-center gap-4 rounded-3xl border border-primary/30 bg-white p-4 shadow-lg">
      <span className="font-mono text-4xl font-bold text-primary">/{item.symbol}/</span>
      <div className="min-w-0 flex-1">
        <p className="font-semibold">{item.nameVi}</p>
        {item.exampleWord ? (
          <p className="text-sm text-muted-foreground">
            như trong <span className="font-semibold text-foreground">{item.exampleWord.word}</span>{" "}
            <IpaText ipa={item.exampleWord.ipa} />
          </p>
        ) : null}
      </div>
      <button
        aria-label="Nghe lại"
        className="cursor-pointer rounded-full bg-primary/10 p-3 text-primary hover:bg-primary/20"
        onClick={() => playPronunciation(item.exampleWord?.word ?? item.symbol, item.audioUrl || item.exampleWord?.audioUrl)}
        type="button"
      >
        <Volume2 className="size-5" />
      </button>
      <Link
        className="inline-flex items-center gap-1 rounded-xl bg-primary px-4 py-2 text-sm font-semibold text-primary-foreground hover:bg-primary/90"
        to={phonemePath(item.symbol)}
      >
        Xem thẻ âm <ArrowRight className="size-4" />
      </Link>
    </div>
  );
}
