import { AudioButton } from "@/features/exercise/AudioButton";
import { SoundChip } from "@/features/pronunciation/SoundChip";
import type { ActivityDetail } from "@/types/activity";
import type { PhonemeCard } from "@/types/phoneme";

type TranscriptLine = { speaker?: string; text: string };

/**
 * Shared material shown above an activity's exercises: the dialogue audio of a
 * LISTENING activity (transcript only after it has been passed) and the text
 * of a READING activity, and the sound cards of a PRONUNCIATION activity.
 */
export function ActivityContent({ activity, phonemes = [] }: { activity: ActivityDetail; phonemes?: PhonemeCard[] }) {
  const content = activity.content ?? {};
  const audioUrl = typeof content.audioUrl === "string" ? content.audioUrl : "";
  const text = typeof content.text === "string" ? content.text : "";
  const transcript = Array.isArray(content.transcript)
    ? (content.transcript as TranscriptLine[])
    : [];
  const canSeeTranscript = activity.bestScore !== null && activity.bestScore >= activity.passScore;

  if (phonemes.length > 0) {
    return (
      <section className="mb-6 flex flex-wrap gap-3">
        {phonemes.map((phoneme) => (
          <SoundChip key={phoneme.symbol} phoneme={phoneme} />
        ))}
      </section>
    );
  }
  if (!audioUrl && !text) {
    return null;
  }

  return (
    <section className="mb-6 rounded-3xl border bg-white p-5 shadow-sm">
      {audioUrl ? (
        <div className="flex flex-wrap items-center gap-3">
          <AudioButton audioUrl={audioUrl} label="Nghe hội thoại" />
          <span className="text-xs text-muted-foreground">Bạn có thể nghe lại bất cứ lúc nào.</span>
        </div>
      ) : null}

      {audioUrl && transcript.length > 0 ? (
        canSeeTranscript ? (
          <details className="mt-4">
            <summary className="cursor-pointer text-sm font-semibold text-primary">Xem lời thoại</summary>
            <ul className="mt-2 space-y-1 text-sm">
              {transcript.map((line, index) => (
                <li key={index}>
                  {line.speaker ? <span className="mr-2 font-bold">{line.speaker}:</span> : null}
                  {line.text}
                </li>
              ))}
            </ul>
          </details>
        ) : (
          <p className="mt-3 text-xs text-muted-foreground">Lời thoại sẽ hiện sau khi bạn đạt hoạt động này.</p>
        )
      ) : null}

      {text ? <p className="whitespace-pre-line text-lg leading-relaxed">{text}</p> : null}
    </section>
  );
}
