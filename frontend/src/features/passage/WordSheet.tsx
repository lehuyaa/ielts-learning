import { BookmarkCheck, BookmarkPlus, Loader2, Volume2, X } from "lucide-react";
import { useEffect } from "react";
import { createPortal } from "react-dom";

import { Button } from "@/components/ui/button";
import { IpaText } from "@/features/pronunciation/IpaText";
import { playPronunciation } from "@/lib/pronunciation";

import { useAddWordToFlashcards, useWordLookup } from "./hooks";

/** Bottom sheet with a tapped word: meaning, IPA, audio, add to flashcards. */
export function WordSheet({ word, passageId, onClose }: { word: string; passageId?: number; onClose: () => void }) {
  const lookup = useWordLookup(word, passageId);
  const add = useAddWordToFlashcards();
  const data = lookup.data;
  const vocabulary = data?.vocabulary ?? null;
  const ipa = data?.ipa || vocabulary?.ipa || "";
  const inFlashcards = vocabulary?.inFlashcards || add.data?.vocabularyId === vocabulary?.id;

  useEffect(() => {
    const onKey = (event: KeyboardEvent) => {
      if (event.key === "Escape") onClose();
    };
    window.addEventListener("keydown", onKey);
    return () => window.removeEventListener("keydown", onKey);
  }, [onClose]);

  return createPortal(
    <div
      aria-label={`Nghĩa của từ ${word}`}
      className="fixed inset-x-0 bottom-0 z-50 mx-auto w-full max-w-lg p-3 sm:bottom-4"
      role="dialog"
    >
      <div className="rounded-3xl border bg-white p-5 shadow-2xl">
        <div className="flex items-start justify-between gap-3">
          <div className="min-w-0">
            <p className="text-2xl font-bold tracking-normal">{vocabulary?.word ?? word}</p>
            {ipa ? <IpaText className="text-sm text-muted-foreground" ipa={ipa} /> : null}
            {vocabulary?.partOfSpeech ? (
              <span className="ml-2 text-xs font-semibold text-muted-foreground">{vocabulary.partOfSpeech}</span>
            ) : null}
          </div>
          <div className="flex shrink-0 gap-1">
            <button
              aria-label="Nghe"
              className="cursor-pointer rounded-full bg-primary/10 p-2.5 text-primary hover:bg-primary/20"
              onClick={() => playPronunciation(vocabulary?.word ?? word, vocabulary?.audioUrl)}
              type="button"
            >
              <Volume2 className="size-5" />
            </button>
            <button aria-label="Đóng" className="cursor-pointer rounded-full p-2.5 hover:bg-muted" onClick={onClose} type="button">
              <X className="size-5" />
            </button>
          </div>
        </div>

        <div className="mt-3 min-h-10">
          {lookup.isLoading ? (
            <p className="flex items-center gap-2 text-sm text-muted-foreground">
              <Loader2 className="size-4 animate-spin" /> Đang tra từ…
            </p>
          ) : null}
          {data?.note ? <p className="font-semibold">{data.note}</p> : null}
          {vocabulary?.meaningVi && vocabulary.meaningVi !== data?.note ? (
            <p className={data?.note ? "mt-1 text-sm text-muted-foreground" : "font-semibold"}>{vocabulary.meaningVi}</p>
          ) : null}
          {data && !data.found ? (
            <p className="text-sm text-muted-foreground">Chưa có từ này trong từ điển của bạn.</p>
          ) : null}
          {lookup.error ? <p className="text-sm text-red-600">Không tra được từ này.</p> : null}
        </div>

        {vocabulary ? (
          <Button
            className="mt-4 w-full"
            disabled={inFlashcards || add.isPending}
            onClick={() => add.mutate(vocabulary.id)}
            variant={inFlashcards ? "outline" : "default"}
          >
            {inFlashcards ? <BookmarkCheck className="size-4" /> : <BookmarkPlus className="size-4" />}
            {inFlashcards ? "Đã có trong flashcard" : "Thêm vào flashcard"}
          </Button>
        ) : null}
      </div>
    </div>,
    document.body,
  );
}
