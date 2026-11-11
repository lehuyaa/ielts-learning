import { ArrowRight, BookmarkPlus, Check, NotebookPen, RotateCcw, Sparkles, Undo2 } from "lucide-react";
import { useState } from "react";
import { Link } from "react-router-dom";

import { getApiErrorMessage } from "@/api/api";
import { EmptyState } from "@/components/state/EmptyState";
import { ErrorState } from "@/components/state/ErrorState";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Skeleton } from "@/components/ui/skeleton";
import { useMistakes, useStartPractice, useUpdateMistake } from "@/features/mistakes/hooks";
import { cn } from "@/lib/utils";
import type { MistakeGroup, MistakeItem, MistakeSource } from "@/types/mistakes";

const sourceLabel: Record<MistakeSource, string> = {
  EXERCISE: "Bài tập",
  WRITING: "Viết",
  SPEAKING: "Nói với AI",
  PRONUNCIATION: "Phát âm",
};

const FIXED_STREAK = 3;

/** Mistake log (Task 73): grouped by skill, add to review, practise a weakness. */
export function MistakesPage() {
  const [filter, setFilter] = useState<"open" | "fixed">("open");
  const query = useMistakes(filter);
  const summary = query.data?.summary;

  return (
    <div className="mx-auto max-w-4xl space-y-5 py-2">
      <header className="flex flex-wrap items-end justify-between gap-4">
        <div>
          <h1 className="flex items-center gap-2 text-2xl font-bold tracking-normal">
            <NotebookPen className="size-6 text-primary" /> Sổ lỗi sai
          </h1>
          <p className="mt-1 text-sm text-muted-foreground">
            Lỗi từ bài tập, bài viết, hội thoại và phát âm, gom theo điểm yếu. Thêm lỗi vào ôn tập; đúng 3 lần liên tiếp là sửa xong.
          </p>
        </div>
        <Button asChild disabled={!summary?.due}>
          <Link aria-disabled={!summary?.due} className={cn(!summary?.due && "pointer-events-none opacity-50")} to="/mistakes/review">
            <RotateCcw className="size-4" /> Ôn lỗi{summary?.due ? ` (${summary.due})` : ""}
          </Link>
        </Button>
      </header>

      {summary ? (
        <div className="flex flex-wrap gap-2">
          <Badge className="bg-white px-3 py-1.5" variant="outline">
            Đang mở <span className="ml-1 font-bold">{summary.open}</span>
          </Badge>
          <Badge className="bg-white px-3 py-1.5" variant="outline">
            Đang ôn <span className="ml-1 font-bold">{summary.reviewing}</span>
          </Badge>
          <Badge className="bg-white px-3 py-1.5" variant="outline">
            Đã sửa <span className="ml-1 font-bold text-emerald-600">{summary.fixed}</span>
          </Badge>
        </div>
      ) : null}

      <div className="flex gap-1 rounded-full border bg-white p-1 text-sm font-semibold w-fit" role="tablist">
        {(["open", "fixed"] as const).map((value) => (
          <button
            aria-selected={filter === value}
            className={cn("cursor-pointer rounded-full px-4 py-1.5", filter === value ? "bg-primary text-primary-foreground" : "text-muted-foreground hover:text-foreground")}
            key={value}
            onClick={() => setFilter(value)}
            role="tab"
            type="button"
          >
            {value === "open" ? "Cần sửa" : "Đã sửa"}
          </button>
        ))}
      </div>

      {query.isLoading ? <Skeleton className="h-64 rounded-3xl" /> : null}
      {query.error ? (
        <ErrorState
          description={getApiErrorMessage(query.error, "Không tải được sổ lỗi.")}
          onRetry={() => void query.refetch()}
          title="Không tải được sổ lỗi"
        />
      ) : null}
      {query.data && query.data.groups.length === 0 ? (
        <EmptyState
          description={
            filter === "open"
              ? "Khi bạn làm sai bài tập, bài viết, hội thoại hay phát âm, lỗi sẽ được ghi ở đây."
              : "Chưa có lỗi nào được sửa xong. Thêm lỗi vào ôn tập để bắt đầu."
          }
          title={filter === "open" ? "Chưa có lỗi nào" : "Chưa sửa xong lỗi nào"}
        />
      ) : null}

      {query.data?.groups.map((group) => <GroupCard group={group} key={group.skillTag} />)}
    </div>
  );
}

function GroupCard({ group }: { group: MistakeGroup }) {
  const practice = useStartPractice();
  const [expanded, setExpanded] = useState(false);
  const visible = expanded ? group.items : group.items.slice(0, 3);
  const open = group.items.some((item) => item.status !== "FIXED");

  return (
    <section className="rounded-3xl border bg-white p-5 shadow-sm">
      <div className="flex flex-wrap items-start justify-between gap-3">
        <div>
          <h2 className="font-bold">{group.label}</h2>
          <p className="text-xs text-muted-foreground">
            {group.count} lỗi · gặp {group.occurrences} lần
          </p>
        </div>
        {group.practiceable && open ? (
          <Button
            disabled={practice.isPending}
            onClick={() => practice.mutate({ skillTag: group.skillTag, returnTo: "/mistakes" })}
            size="sm"
            variant="outline"
          >
            <Sparkles className="size-4" /> {practice.isPending ? "Đang tạo bài…" : "Luyện 10 câu"}
          </Button>
        ) : null}
      </div>
      {practice.error ? (
        <p className="mt-2 text-sm text-red-600">{getApiErrorMessage(practice.error, "Không tạo được bài luyện.")}</p>
      ) : null}
      <ul className="mt-3 divide-y">
        {visible.map((item) => (
          <MistakeRow item={item} key={item.id} />
        ))}
      </ul>
      {group.items.length > 3 ? (
        <button className="mt-2 cursor-pointer text-sm font-semibold text-primary hover:underline" onClick={() => setExpanded(!expanded)} type="button">
          {expanded ? "Thu gọn" : `Xem thêm ${group.items.length - 3} lỗi`}
        </button>
      ) : null}
    </section>
  );
}

function MistakeRow({ item }: { item: MistakeItem }) {
  const update = useUpdateMistake();
  const busy = update.isPending;
  return (
    <li className="flex flex-wrap items-start gap-3 py-3">
      <div className="min-w-0 flex-1 space-y-1">
        <p className="flex flex-wrap items-baseline gap-x-2">
          <span className="text-red-600 line-through decoration-red-300">{item.wrongText}</span>
          <ArrowRight className="size-3.5 shrink-0 self-center text-muted-foreground" />
          <span className="font-semibold text-emerald-700">{item.correctText}</span>
        </p>
        {item.explanationVi ? <p className="text-sm text-muted-foreground">{item.explanationVi}</p> : null}
        <p className="text-xs text-muted-foreground">
          {sourceLabel[item.source]} · {item.occurrences} lần
          {item.status === "REVIEWING" ? ` · đang ôn ${item.reviewStreak}/${FIXED_STREAK}` : ""}
        </p>
      </div>
      <div className="flex shrink-0 gap-2">
        {item.status === "OPEN" ? (
          <Button disabled={busy} onClick={() => update.mutate({ id: item.id, status: "REVIEWING" })} size="sm" variant="outline">
            <BookmarkPlus className="size-4" /> Thêm vào ôn tập
          </Button>
        ) : null}
        {item.status === "REVIEWING" ? (
          <Button disabled={busy} onClick={() => update.mutate({ id: item.id, status: "OPEN" })} size="sm" variant="ghost">
            <Undo2 className="size-4" /> Bỏ ôn
          </Button>
        ) : null}
        {item.status !== "FIXED" ? (
          <Button disabled={busy} onClick={() => update.mutate({ id: item.id, status: "FIXED" })} size="sm" variant="ghost">
            <Check className="size-4" /> Đã nhớ
          </Button>
        ) : (
          <Button disabled={busy} onClick={() => update.mutate({ id: item.id, status: "OPEN" })} size="sm" variant="ghost">
            <Undo2 className="size-4" /> Mở lại
          </Button>
        )}
      </div>
    </li>
  );
}
