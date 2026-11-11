import { useQueryClient } from "@tanstack/react-query";
import { ArrowLeft, CheckCircle2, Eye, PartyPopper, Volume2, XCircle } from "lucide-react";
import { useState } from "react";
import { Link } from "react-router-dom";

import { getApiErrorMessage } from "@/api/api";
import { ErrorState } from "@/components/state/ErrorState";
import { Button } from "@/components/ui/button";
import { Skeleton } from "@/components/ui/skeleton";
import { mistakesQueryKey, useMistakeReview, useReviewMistake } from "@/features/mistakes/hooks";
import { phonemePath } from "@/lib/ipa";
import { playPronunciation } from "@/lib/pronunciation";
import { cn } from "@/lib/utils";
import type { MistakeItem, MistakeRating } from "@/types/mistakes";

const ratings: { rating: MistakeRating; label: string; className: string }[] = [
  { rating: "AGAIN", label: "Chưa nhớ", className: "border-red-200 text-red-700 hover:bg-red-50" },
  { rating: "HARD", label: "Khó", className: "border-amber-200 text-amber-700 hover:bg-amber-50" },
  { rating: "GOOD", label: "Nhớ rồi", className: "border-emerald-200 text-emerald-700 hover:bg-emerald-50" },
  { rating: "EASY", label: "Dễ", className: "border-primary/30 text-primary hover:bg-primary/5" },
];

function normalize(text: string) {
  return text.toLowerCase().replace(/[’']/g, "'").replace(/[^a-z0-9'/ːəɪʊɛæʌɑɔθðʃʒŋ]+/g, " ").trim();
}

/** Mistake review cards (Task 73): wrong → try the correction → rate. */
export function MistakeReviewPage() {
  const queryClient = useQueryClient();
  const query = useMistakeReview();
  const review = useReviewMistake();
  const [index, setIndex] = useState(0);
  const [fixedCount, setFixedCount] = useState(0);
  const cards = query.data?.cards ?? [];
  const card = cards[index];

  function rate(rating: MistakeRating) {
    if (!card) return;
    review.mutate(
      { id: card.id, rating },
      {
        onSuccess: (result) => {
          if (result.fixed) setFixedCount((count) => count + 1);
          setIndex((current) => current + 1);
          void queryClient.invalidateQueries({ queryKey: [...mistakesQueryKey, "list"] });
        },
      },
    );
  }

  return (
    <div className="mx-auto max-w-2xl space-y-5 py-2">
      <div className="flex items-center justify-between">
        <Link className="inline-flex items-center gap-1 text-sm font-semibold text-muted-foreground hover:text-foreground" to="/mistakes">
          <ArrowLeft className="size-4" /> Sổ lỗi
        </Link>
        {cards.length > 0 && card ? (
          <span className="text-sm text-muted-foreground">
            {index + 1}/{cards.length}
          </span>
        ) : null}
      </div>

      {query.isLoading ? <Skeleton className="h-80 rounded-3xl" /> : null}
      {query.error ? (
        <ErrorState description={getApiErrorMessage(query.error, "Không tải được thẻ ôn.")} onRetry={() => void query.refetch()} title="Không tải được thẻ ôn" />
      ) : null}

      {query.data && !card ? (
        <section className="rounded-3xl border bg-white p-8 text-center shadow-sm">
          <PartyPopper className="mx-auto size-10 text-primary" />
          <h1 className="mt-3 text-xl font-bold">{cards.length > 0 ? "Xong lượt ôn!" : "Không có lỗi nào cần ôn lúc này"}</h1>
          <p className="mt-1 text-muted-foreground">
            {cards.length > 0
              ? `Bạn đã ôn ${cards.length} lỗi${fixedCount > 0 ? `, sửa xong ${fixedCount} lỗi` : ""}.`
              : "Thêm lỗi vào ôn tập ở Sổ lỗi sai."}
          </p>
          <Button asChild className="mt-5">
            <Link to="/mistakes">Về sổ lỗi</Link>
          </Button>
        </section>
      ) : null}

      {card ? <ReviewCard busy={review.isPending} card={card} key={card.id} onRate={rate} /> : null}
      {review.error ? <p className="text-center text-sm text-red-600">{getApiErrorMessage(review.error, "Không lưu được.")}</p> : null}
    </div>
  );
}

function ReviewCard({ card, busy, onRate }: { card: MistakeItem; busy: boolean; onRate: (rating: MistakeRating) => void }) {
  const [attempt, setAttempt] = useState("");
  const [revealed, setRevealed] = useState(false);
  const typed = attempt.trim() !== "";
  const matched = typed && normalize(attempt) === normalize(card.correctText);
  // Sound mistakes are a word to say again, not a sentence to correct.
  const sound = card.source === "PRONUNCIATION" && card.skillTag.startsWith("sound:");
  const symbol = card.skillTag.replace(/^sound:/, "");

  if (sound) {
    return (
      <section className="rounded-3xl border bg-white p-6 shadow-sm">
        <p className="text-xs font-semibold uppercase text-muted-foreground">{card.label}</p>
        <p className="mt-3 text-sm text-muted-foreground">Nghe rồi nói to từ này, chú ý âm {card.correctText}.</p>
        <button
          className="mt-2 inline-flex cursor-pointer items-center gap-2 text-2xl font-semibold hover:text-primary"
          onClick={() => playPronunciation(card.wrongText)}
          type="button"
        >
          {card.wrongText} <Volume2 className="size-5 text-primary" />
        </button>
        {card.explanationVi ? <p className="mt-3 text-sm text-muted-foreground">{card.explanationVi}</p> : null}
        <Link className="mt-2 inline-block text-sm font-semibold text-primary hover:underline" to={phonemePath(symbol)}>
          Mở thẻ âm /{symbol}/ để luyện nói →
        </Link>
        <RatingButtons busy={busy} onRate={onRate} streak={card.reviewStreak} />
      </section>
    );
  }

  return (
    <section className="rounded-3xl border bg-white p-6 shadow-sm">
      <p className="text-xs font-semibold uppercase text-muted-foreground">{card.label}</p>
      <p className="mt-3 text-sm text-muted-foreground">Câu này bạn từng làm sai. Câu đúng là gì?</p>
      <p className="mt-1 text-xl font-semibold text-red-600">{card.wrongText}</p>

      <form
        className="mt-4"
        onSubmit={(event) => {
          event.preventDefault();
          setRevealed(true);
        }}
      >
        <input
          className="h-12 w-full rounded-2xl border bg-muted/30 px-4 outline-none focus:border-primary disabled:opacity-70"
          disabled={revealed}
          onChange={(event) => setAttempt(event.target.value)}
          placeholder="Gõ câu đúng (không bắt buộc), rồi bấm Xem đáp án"
          value={attempt}
        />
        {!revealed ? (
          <Button className="mt-3 w-full" type="submit">
            <Eye className="size-4" /> Xem đáp án
          </Button>
        ) : null}
      </form>

      {revealed ? (
        <div className="mt-4 space-y-3">
          <div className={cn("rounded-2xl border p-4", typed ? (matched ? "border-emerald-200 bg-emerald-50" : "border-amber-200 bg-amber-50") : "bg-muted/30")}>
            {typed ? (
              <p className="mb-1 flex items-center gap-1.5 text-sm font-semibold">
                {matched ? <CheckCircle2 className="size-4 text-emerald-600" /> : <XCircle className="size-4 text-amber-600" />}
                {matched ? "Chính xác!" : "Chưa giống đáp án"}
              </p>
            ) : null}
            <p className="text-lg font-semibold text-emerald-700">{card.correctText}</p>
            {card.explanationVi ? <p className="mt-1 text-sm text-muted-foreground">{card.explanationVi}</p> : null}
          </div>
          <RatingButtons busy={busy} onRate={onRate} streak={card.reviewStreak} />
        </div>
      ) : null}
    </section>
  );
}

function RatingButtons({ busy, streak, onRate }: { busy: boolean; streak: number; onRate: (rating: MistakeRating) => void }) {
  return (
    <div className="mt-4 space-y-2">
      <p className="text-center text-xs text-muted-foreground">Bạn nhớ lỗi này tới đâu? (đúng {streak}/3 lần liên tiếp)</p>
      <div className="grid grid-cols-2 gap-2 sm:grid-cols-4">
        {ratings.map((option) => (
          <button
            className={cn("cursor-pointer rounded-2xl border bg-white px-3 py-2.5 font-semibold disabled:opacity-50", option.className)}
            disabled={busy}
            key={option.rating}
            onClick={() => onRate(option.rating)}
            type="button"
          >
            {option.label}
          </button>
        ))}
      </div>
    </div>
  );
}
