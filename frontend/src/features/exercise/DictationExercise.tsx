import { Snail, Volume2 } from "lucide-react";
import { useState } from "react";

import { segmentSource } from "@/features/audio/segmentPlayer";
import { usePassage } from "@/features/passage/hooks";
import { playPronunciation } from "@/lib/pronunciation";
import { cn } from "@/lib/utils";

import { str } from "./payload";
import type { ExerciseRendererProps } from "./renderers";
import { diffWords } from "./wordDiff";

/**
 * DICTATION: hear one sentence (its own file, or one segment of a passage),
 * at normal or slow speed, and type it. After checking, a word-level diff
 * shows correct, missing and extra words.
 */
export function DictationExercise({ exercise, result, onChange, onSubmit }: ExerciseRendererProps) {
  const [text, setText] = useState("");
  const [plays, setPlays] = useState(0);
  const answered = result?.answer as { text?: string } | undefined;
  const passageId = Number(exercise.payload.passageId) || undefined;
  const segmentIndex = Number(exercise.payload.segmentIndex);
  const passage = usePassage(passageId);
  const wordCount = Number(exercise.payload.wordCount) || 0;

  function play(rate: number) {
    setPlays((value) => value + 1);
    const audioUrl = str(exercise.payload, "audioUrl");
    if (audioUrl || !passage.data) {
      // Without a file the browser voice cannot help (it would need the answer).
      playPronunciation("", audioUrl || null, { rate });
      return;
    }
    const segment = passage.data.segments[segmentIndex];
    if (!segment) return;
    const source = segmentSource(passage.data, segment);
    if (!source || source.start === undefined) {
      playPronunciation(segment.text, source?.url, { rate });
      return;
    }
    const audio = new Audio(source.url);
    audio.playbackRate = rate;
    audio.currentTime = source.start;
    audio.ontimeupdate = () => {
      if (source.end !== undefined && audio.currentTime >= source.end) audio.pause();
    };
    void audio.play();
  }

  const correct = typeof result?.correctAnswer === "string" ? result.correctAnswer : "";

  return (
    <div>
      <div className="mb-4 flex flex-wrap items-center justify-center gap-2">
        <button
          className="inline-flex cursor-pointer items-center gap-2 rounded-full bg-primary/10 px-5 py-3 font-semibold text-primary hover:bg-primary/20"
          onClick={() => play(1)}
          type="button"
        >
          <Volume2 className="size-5" /> Nghe
        </button>
        <button
          className="inline-flex cursor-pointer items-center gap-2 rounded-full border px-4 py-3 font-semibold text-muted-foreground hover:border-primary/40 hover:text-foreground"
          onClick={() => play(0.75)}
          type="button"
        >
          <Snail className="size-5" /> Nghe chậm
        </button>
      </div>
      <p className="mb-3 text-center text-xs text-muted-foreground">
        {wordCount ? `Câu có ${wordCount} từ` : null}
        {plays > 0 ? ` · đã nghe ${plays} lần` : null}
      </p>
      <textarea
        autoCapitalize="sentences"
        autoComplete="off"
        autoFocus
        className={cn(
          "min-h-20 w-full resize-none rounded-2xl border-2 border-input bg-white p-4 text-lg outline-none focus:border-primary disabled:opacity-100",
          result && (result.isCorrect ? "border-emerald-500 bg-emerald-50" : "border-amber-400 bg-amber-50"),
        )}
        disabled={Boolean(result)}
        onChange={(event) => {
          setText(event.target.value);
          onChange(event.target.value.trim() ? { text: event.target.value } : null);
        }}
        onKeyDown={(event) => {
          if (event.key === "Enter") {
            event.preventDefault();
            if (text.trim()) onSubmit();
          }
        }}
        placeholder="Gõ lại câu bạn nghe được…"
        rows={2}
        spellCheck={false}
        value={answered?.text ?? text}
      />
      {result && correct && result.scoreRatio < 1 ? <WordDiff answer={answered?.text ?? ""} correct={correct} /> : null}
    </div>
  );
}

function WordDiff({ answer, correct }: { answer: string; correct: string }) {
  const tokens = diffWords(answer, correct);
  if (tokens.every((token) => token.kind === "ok")) {
    return null;
  }
  return (
    <div className="mt-3 rounded-2xl bg-muted/40 p-3 text-base leading-loose">
      <p className="mb-1 text-xs font-semibold text-muted-foreground">
        So sánh từng từ: <span className="text-emerald-700">đúng</span> ·{" "}
        <span className="text-amber-700">bạn thiếu</span> · <span className="text-red-600 line-through">bạn viết thừa/sai</span>
      </p>
      {tokens.map((token, index) => (
        <span
          className={cn(
            "mr-1.5 inline-block rounded px-1",
            token.kind === "ok" && "text-emerald-700",
            token.kind === "missing" && "bg-amber-100 font-semibold text-amber-800",
            token.kind === "extra" && "text-red-600 line-through",
          )}
          key={index}
        >
          {token.word}
        </span>
      ))}
    </div>
  );
}
