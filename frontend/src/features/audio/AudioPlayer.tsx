import { Eye, Mic, Pause, Play, RotateCcw, SkipBack, SkipForward } from "lucide-react";
import { useEffect, useMemo, useState, type ReactNode } from "react";

import { TappableText } from "@/features/passage/TappableText";
import { PronunciationRecorder } from "@/features/pronunciation/PronunciationRecorder";
import { ScoredText } from "@/features/pronunciation/ScoredText";
import { scoreTone } from "@/features/pronunciation/scoreColors";
import { usePronunciationStatus } from "@/features/pronunciation/hooks/usePronunciation";
import { canRecord } from "@/features/pronunciation/voiceRecorder";
import { cn } from "@/lib/utils";
import type { PassageCard, TranscriptMode } from "@/types/passage";
import type { PronunciationResult } from "@/types/pronunciation";

import { SegmentPlayer } from "./segmentPlayer";

type AudioPlayerProps = {
  passage: PassageCard;
  speeds?: number[];
  transcript?: TranscriptMode;
  /** The learner has answered every exercise (AFTER_ANSWER shows the text). */
  answered?: boolean;
  /** The activity is passed (every mode shows the text). */
  passed?: boolean;
};

/**
 * Listening player: play the whole passage or one sentence at a time, at
 * 0.75x or 1x, with an optional transcript that highlights the current
 * sentence. A segment plays its own audio file, or its start/end slice of the
 * passage audio; without audio it is read by the browser voice.
 */
export function AudioPlayer({ passage, speeds = [0.75, 1], transcript = "AFTER_ANSWER", answered = false, passed = false }: AudioPlayerProps) {
  const segments = useMemo(
    () => (passage.segments.length > 0 ? passage.segments : [{ index: 0, text: passage.bodyText }]),
    [passage],
  );
  const [current, setCurrent] = useState(0);
  const [playing, setPlaying] = useState(false);
  const [speed, setSpeed] = useState(speeds.includes(1) ? 1 : speeds[0] ?? 1);
  const [unblurred, setUnblurred] = useState(false);
  const [player] = useState(() => new SegmentPlayer({ onSegment: setCurrent, onPlaying: setPlaying }));
  const pronunciation = usePronunciationStatus();
  const [shadowing, setShadowing] = useState(false);
  const [shadowScores, setShadowScores] = useState<Record<number, PronunciationResult>>({});

  useEffect(() => () => player.stop(), [player]);

  function playSegment(index: number, continueAfter: boolean) {
    player.play(passage, segments, index, continueAfter);
  }

  function stop() {
    player.stop();
  }

  function changeSpeed(value: number) {
    setSpeed(value);
    player.setRate(value);
  }

  const showText = passed || (transcript === "AFTER_ANSWER" && answered) || (transcript === "BLURRED" && unblurred);
  const blurred = transcript === "BLURRED" && !showText;
  // Shadowing (Task 72): play a sentence, say it, see each word scored.
  const canShadow = Boolean(pronunciation.data?.enabled) && canRecord() && showText;
  const shadowOn = shadowing && canShadow;

  return (
    <div className="space-y-4">
      <div className="flex flex-wrap items-center gap-2">
        <IconButton disabled={current === 0} label="Câu trước" onClick={() => playSegment(Math.max(current - 1, 0), false)}>
          <SkipBack className="size-4" />
        </IconButton>
        <button
          aria-label={playing ? "Tạm dừng" : "Nghe cả bài"}
          className="inline-flex cursor-pointer items-center gap-2 rounded-full bg-primary px-5 py-2.5 font-semibold text-primary-foreground hover:bg-primary/90"
          onClick={() => (playing ? stop() : playSegment(current >= segments.length - 1 && !playing ? 0 : current, true))}
          type="button"
        >
          {playing ? <Pause className="size-4" /> : <Play className="size-4" />}
          {playing ? "Tạm dừng" : "Nghe cả bài"}
        </button>
        <IconButton disabled={current >= segments.length - 1} label="Câu sau" onClick={() => playSegment(Math.min(current + 1, segments.length - 1), false)}>
          <SkipForward className="size-4" />
        </IconButton>
        <button
          className="inline-flex cursor-pointer items-center gap-1.5 rounded-full border px-3 py-2 text-sm font-semibold hover:border-primary/40"
          onClick={() => playSegment(current, false)}
          type="button"
        >
          <RotateCcw className="size-4" /> Nghe lại câu {current + 1}
        </button>
        {canShadow ? (
          <button
            aria-pressed={shadowOn}
            className={cn(
              "inline-flex cursor-pointer items-center gap-1.5 rounded-full border px-3 py-2 text-sm font-semibold",
              shadowOn ? "border-primary bg-primary/10 text-primary" : "hover:border-primary/40",
            )}
            onClick={() => setShadowing(!shadowing)}
            title="Nghe từng câu rồi nói theo, AI chấm phát âm từng từ"
            type="button"
          >
            <Mic className="size-4" /> Nói theo
          </button>
        ) : null}
        <div className="ml-auto flex rounded-full border p-0.5 text-sm" role="group" aria-label="Tốc độ">
          {speeds.map((value) => (
            <button
              aria-pressed={speed === value}
              className={cn(
                "cursor-pointer rounded-full px-3 py-1 font-semibold",
                speed === value ? "bg-primary text-primary-foreground" : "text-muted-foreground hover:text-foreground",
              )}
              key={value}
              onClick={() => changeSpeed(value)}
              type="button"
            >
              {value}x
            </button>
          ))}
        </div>
      </div>

      <div className="flex gap-1" aria-hidden="true">
        {segments.map((segment, index) => (
          <button
            className={cn(
              "h-1.5 flex-1 cursor-pointer rounded-full",
              index === current ? "bg-primary" : index < current ? "bg-primary/40" : "bg-muted",
            )}
            key={segment.index}
            onClick={() => playSegment(index, false)}
            tabIndex={-1}
            type="button"
          />
        ))}
      </div>

      {transcript === "HIDDEN" && !showText ? null : (
        <div className="relative">
          {!showText && !blurred ? (
            <p className="text-xs text-muted-foreground">Lời thoại sẽ hiện sau khi bạn trả lời hết các câu hỏi.</p>
          ) : (
            <ol className={cn("space-y-1.5", blurred && "pointer-events-none select-none blur-sm")}>
              {segments.map((segment, index) => (
                <li
                  className={cn(
                    "flex gap-2 rounded-xl px-2 py-1 transition-colors",
                    index === current && "bg-primary/10",
                  )}
                  key={segment.index}
                >
                  <button
                    aria-label={`Nghe câu ${index + 1}`}
                    className="mt-0.5 shrink-0 cursor-pointer text-muted-foreground hover:text-primary"
                    onClick={() => playSegment(index, false)}
                    type="button"
                  >
                    <Play className="size-3.5" />
                  </button>
                  <div className="min-w-0 flex-1">
                    <span>
                      {segment.speaker ? <span className="mr-1.5 font-bold">{segment.speaker}:</span> : null}
                      {showText ? <TappableText passageId={passage.id} text={segment.text} /> : segment.text}
                    </span>
                    {shadowOn && shadowScores[index] ? (
                      <div className="mt-1 flex items-start gap-2 text-sm">
                        <span className={cn("shrink-0 pt-1 font-bold", scoreTone(shadowScores[index].accuracy))}>
                          {shadowScores[index].accuracy}
                        </span>
                        <ScoredText className="[&>p]:text-base" key={shadowScores[index].attemptId} words={shadowScores[index].words} />
                      </div>
                    ) : null}
                    {shadowOn && index === current ? (
                      <div className="mt-2">
                        <PronunciationRecorder
                          compact
                          hideResult
                          key={segment.index}
                          onScored={(result) => setShadowScores((scores) => ({ ...scores, [index]: result }))}
                          referenceText={segment.text}
                          source="SHADOWING"
                        />
                      </div>
                    ) : null}
                  </div>
                </li>
              ))}
            </ol>
          )}
          {blurred ? (
            <button
              className="absolute inset-0 m-auto inline-flex h-10 w-fit cursor-pointer items-center gap-2 rounded-full border bg-white px-4 text-sm font-semibold shadow-sm hover:border-primary/40"
              onClick={() => setUnblurred(true)}
              type="button"
            >
              <Eye className="size-4" /> Hiện lời thoại
            </button>
          ) : null}
        </div>
      )}
    </div>
  );
}

function IconButton({ label, disabled, onClick, children }: { label: string; disabled?: boolean; onClick: () => void; children: ReactNode }) {
  return (
    <button
      aria-label={label}
      className="inline-flex size-10 cursor-pointer items-center justify-center rounded-full border hover:border-primary/40 disabled:cursor-not-allowed disabled:opacity-40"
      disabled={disabled}
      onClick={onClick}
      title={label}
      type="button"
    >
      {children}
    </button>
  );
}
