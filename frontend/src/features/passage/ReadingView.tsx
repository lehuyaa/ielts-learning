import { BookOpenText } from "lucide-react";

import { AudioButton } from "@/features/exercise/AudioButton";
import type { PassageCard, PassageGenre } from "@/types/passage";

import { TappableText } from "./TappableText";

const genreLabel: Record<PassageGenre, string> = {
  DIALOGUE: "Hội thoại",
  MESSAGE: "Tin nhắn",
  EMAIL: "Email",
  NOTICE: "Thông báo",
  MENU: "Thực đơn",
  SCHEDULE: "Lịch trình",
  ARTICLE: "Bài viết",
  PROFILE: "Hồ sơ",
  MONOLOGUE: "Bài nói",
};

/** A READING passage with tap-to-translate words. */
export function ReadingView({ passage, tapToTranslate = true, readAloud = false }: { passage: PassageCard; tapToTranslate?: boolean; readAloud?: boolean }) {
  return (
    <div>
      <div className="mb-3 flex flex-wrap items-center gap-2 text-xs font-semibold text-muted-foreground">
        <BookOpenText className="size-4 text-primary" />
        <span className="uppercase">{genreLabel[passage.genre] ?? passage.genre}</span>
        <span>· {passage.wordCount} từ</span>
        {tapToTranslate ? <span className="font-normal">· Chạm vào từ để xem nghĩa</span> : null}
        {readAloud || passage.audioUrl ? (
          <AudioButton audioUrl={passage.audioUrl ?? undefined} className="ml-auto" fallbackText={passage.bodyText} label="Nghe đọc" size="sm" />
        ) : null}
      </div>
      <div className="rounded-2xl bg-amber-50/60 p-4 text-lg leading-relaxed sm:p-5">
        {tapToTranslate ? (
          <TappableText passageId={passage.id} text={passage.bodyText} />
        ) : (
          <p className="whitespace-pre-line">{passage.bodyText}</p>
        )}
      </div>
    </div>
  );
}
