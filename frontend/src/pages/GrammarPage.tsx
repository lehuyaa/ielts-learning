import { ArrowRight, Puzzle } from "lucide-react";
import { useState } from "react";
import { Link } from "react-router-dom";

import { APIError } from "@/api/api";
import { EmptyState } from "@/components/state/EmptyState";
import { ErrorState } from "@/components/state/ErrorState";
import { Skeleton } from "@/components/ui/skeleton";
import { MasteryBar } from "@/features/grammar/MasteryBar";
import { useGrammarPoints } from "@/features/grammar/hooks/useGrammar";
import { cn } from "@/lib/utils";
import type { GrammarPointItem } from "@/types/grammar";

type SortMode = "order" | "weakest";

export function GrammarPage() {
  const query = useGrammarPoints();
  const [sort, setSort] = useState<SortMode>("order");
  const items = sortItems(query.data?.items ?? [], sort);

  return (
    <div className="mx-auto max-w-4xl space-y-6 py-2">
      <header className="flex flex-wrap items-end justify-between gap-3">
        <div>
          <h1 className="flex items-center gap-2 text-2xl font-bold tracking-normal">
            <Puzzle className="size-6 text-primary" /> Ngữ pháp
          </h1>
          <p className="mt-1 text-sm text-muted-foreground">
            Mức thành thạo là điểm trung bình của 20 câu gần nhất bạn làm cho mỗi điểm ngữ pháp.
          </p>
        </div>
        <div className="flex rounded-xl border bg-white p-1 text-sm">
          {(["order", "weakest"] as const).map((mode) => (
            <button
              className={cn(
                "cursor-pointer rounded-lg px-3 py-1.5 font-semibold",
                sort === mode ? "bg-primary text-primary-foreground" : "text-muted-foreground hover:bg-muted",
              )}
              key={mode}
              onClick={() => setSort(mode)}
              type="button"
            >
              {mode === "order" ? "Theo lộ trình" : "Yếu nhất trước"}
            </button>
          ))}
        </div>
      </header>

      {query.isLoading ? (
        <div className="grid gap-3 sm:grid-cols-2">
          {Array.from({ length: 4 }).map((_, i) => (
            <Skeleton className="h-32 rounded-2xl" key={i} />
          ))}
        </div>
      ) : null}

      {query.error ? (
        <ErrorState
          description={query.error instanceof APIError ? query.error.message : "Không tải được danh sách ngữ pháp."}
          onRetry={() => void query.refetch()}
          title="Không tải được ngữ pháp"
        />
      ) : null}

      {query.data && items.length === 0 ? (
        <EmptyState
          actionHref="/admin/content/import"
          actionLabel="Import nội dung"
          description="Hãy import file content/a1/a1-grammar.json ở trang Import Lessons."
          title="Chưa có điểm ngữ pháp nào"
        />
      ) : null}

      <div className="grid gap-3 sm:grid-cols-2">
        {items.map((item) => (
          <GrammarCard item={item} key={item.id} />
        ))}
      </div>
    </div>
  );
}

function GrammarCard({ item }: { item: GrammarPointItem }) {
  return (
    <Link
      className="group rounded-2xl border bg-white p-4 shadow-sm transition-colors hover:border-primary/40"
      to={`/grammar/${item.slug}`}
    >
      <div className="flex items-start justify-between gap-2">
        <div>
          <p className="text-xs font-bold text-primary">
            {item.cefrLevel} · G{item.orderIndex}
          </p>
          <h2 className="mt-0.5 font-bold tracking-normal">{item.titleVi || item.title}</h2>
          <p className="text-xs text-muted-foreground">{item.title}</p>
        </div>
        <ArrowRight className="size-4 text-muted-foreground transition-transform group-hover:translate-x-0.5" />
      </div>
      <MasteryBar className="mt-4" mastery={item.progress.mastery} practised={item.progress.practised} />
      <p className="mt-2 text-xs text-muted-foreground">
        {item.exerciseCount > 0
          ? `${item.exerciseCount} câu luyện tập từ các bài bạn đã học`
          : "Học bài có điểm ngữ pháp này để mở phần luyện tập"}
      </p>
    </Link>
  );
}

function sortItems(items: GrammarPointItem[], sort: SortMode) {
  const sorted = [...items];
  if (sort === "weakest") {
    // Practised points by mastery (weakest first), then unpractised ones.
    sorted.sort((a, b) => {
      if (a.progress.practised !== b.progress.practised) return a.progress.practised ? -1 : 1;
      return a.progress.mastery - b.progress.mastery || a.orderIndex - b.orderIndex;
    });
  }
  return sorted;
}
