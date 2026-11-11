import { Loader2, Mic, RotateCcw, Square, X } from "lucide-react";
import { useEffect, useState } from "react";

import { getApiErrorMessage } from "@/api/api";
import { cn } from "@/lib/utils";
import type { PronunciationResult, PronunciationSource } from "@/types/pronunciation";

import { PronunciationScore } from "./ScoredText";
import { useAssessPronunciation, usePronunciationStatus } from "./hooks/usePronunciation";
import { VoiceRecorder, canRecord, microphoneError } from "./voiceRecorder";
import { toWav16k } from "./wav";

type PronunciationRecorderProps = {
  referenceText: string;
  source: PronunciationSource;
  exerciseId?: number;
  onScored?: (result: PronunciationResult) => void;
  disabled?: boolean;
  compact?: boolean;
  /** Hide the result here (the parent shows it). */
  hideResult?: boolean;
};

type Phase = "idle" | "recording" | "scoring";

/**
 * Record → Azure score → coloured result (Tasks 71-72). Max length comes
 * from /pronunciation/status (15 s).
 */
export function PronunciationRecorder({ referenceText, source, exerciseId, onScored, disabled, compact, hideResult }: PronunciationRecorderProps) {
  const status = usePronunciationStatus();
  const assess = useAssessPronunciation();
  const maxSeconds = status.data?.maxSeconds ?? 15;
  const [phase, setPhase] = useState<Phase>("idle");
  const [level, setLevel] = useState(0);
  const [elapsed, setElapsed] = useState(0);
  const [error, setError] = useState<string | null>(null);
  const [result, setResult] = useState<PronunciationResult | null>(null);
  const [recorder] = useState(() => new VoiceRecorder({ onLevel: setLevel, onElapsed: setElapsed }));

  useEffect(() => () => recorder.cancel(), [recorder]);

  async function record() {
    setError(null);
    setPhase("recording");
    let recording: Blob | null;
    try {
      recording = await recorder.record(maxSeconds);
    } catch (err) {
      setError(microphoneError(err));
      setPhase("idle");
      return;
    }
    if (!recording || recording.size === 0) {
      setPhase("idle");
      return;
    }
    setPhase("scoring");
    try {
      const audio = await toWav16k(recording, maxSeconds);
      const scored = await assess.mutateAsync({ audio, referenceText, exerciseId, source });
      setResult(scored);
      onScored?.(scored);
    } catch (err) {
      setError(getApiErrorMessage(err, "Không chấm được bản ghi. Hãy thử lại."));
    } finally {
      setPhase("idle");
    }
  }

  if (!canRecord()) {
    return <p className="text-sm text-muted-foreground">Trình duyệt này không ghi âm được. Hãy dùng Chrome, Edge hoặc Safari mới.</p>;
  }

  const recording = phase === "recording";
  return (
    <div className={cn("space-y-3", compact ? "" : "text-center")}>
      <div className={cn("flex flex-wrap items-center gap-3", !compact && "justify-center")}>
        {recording ? (
          <>
            <button
              aria-label="Dừng ghi âm"
              className={cn(
                "grid shrink-0 cursor-pointer place-items-center rounded-full bg-red-500 text-white hover:bg-red-600",
                compact ? "size-9" : "size-14",
              )}
              onClick={() => recorder.stop()}
              type="button"
            >
              <Square className={compact ? "size-3.5" : "size-5"} fill="currentColor" />
            </button>
            <LevelMeter level={level} />
            <span className="w-14 font-mono text-sm tabular-nums text-muted-foreground">
              {Math.floor(elapsed)}s/{maxSeconds}s
            </span>
            <button
              aria-label="Huỷ"
              className="cursor-pointer rounded-full p-1 text-muted-foreground hover:text-foreground"
              onClick={() => recorder.cancel()}
              type="button"
            >
              <X className="size-4" />
            </button>
          </>
        ) : (
          <button
            className={cn(
              "inline-flex cursor-pointer items-center gap-2 rounded-full font-semibold disabled:cursor-not-allowed disabled:opacity-50",
              compact ? "border px-3 py-1.5 text-sm hover:border-primary/40" : "bg-primary px-5 py-3 text-primary-foreground hover:bg-primary/90",
            )}
            disabled={disabled || phase === "scoring"}
            onClick={() => void record()}
            type="button"
          >
            {phase === "scoring" ? (
              <Loader2 className="size-4 animate-spin" />
            ) : result ? (
              <RotateCcw className="size-4" />
            ) : (
              <Mic className="size-4" />
            )}
            {phase === "scoring" ? "Đang chấm…" : result ? "Ghi lại" : compact ? "Nói" : "Bấm để nói"}
          </button>
        )}
      </div>
      {error ? <p className="text-sm text-red-600">{error}</p> : null}
      {result && !hideResult && phase !== "recording" ? <PronunciationScore compact={compact} result={result} /> : null}
    </div>
  );
}

function LevelMeter({ level }: { level: number }) {
  const bars = 12;
  return (
    <div aria-hidden="true" className="flex h-8 items-center gap-0.5">
      {Array.from({ length: bars }, (_, index) => {
        const active = level * bars > index;
        return (
          <span
            className={cn("w-1.5 rounded-full transition-all", active ? "bg-primary" : "bg-muted")}
            key={index}
            style={{ height: `${30 + ((index * 37) % 70)}%` }}
          />
        );
      })}
    </div>
  );
}
