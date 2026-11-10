import { Check, Mic, RotateCcw, X } from "lucide-react";
import { useState, type ReactNode } from "react";

import { Input } from "@/components/ui/input";
import { cn } from "@/lib/utils";
import type {
  AnswerResult,
  Exercise,
  ExerciseAnswer,
  ExerciseType,
} from "@/types/activity";

import { AudioButton } from "./AudioButton";
import { choiceOptions, optionLabel, str, strList } from "./payload";

export type ExerciseRendererProps = {
  exercise: Exercise;
  /** Set once the item has been answered; renderers become read-only. */
  result?: AnswerResult;
  onChange: (answer: ExerciseAnswer | null) => void;
  /** Submit from the keyboard (Enter in text inputs). */
  onSubmit: () => void;
};

type Renderer = (props: ExerciseRendererProps) => ReactNode;

const renderers: Record<ExerciseType, Renderer> = {
  MULTIPLE_CHOICE: ChoiceExercise,
  LISTEN_PICK_IMAGE: ChoiceExercise,
  MINIMAL_PAIR: ChoiceExercise,
  IPA_TO_WORD: ChoiceExercise,
  WORD_TO_IPA: ChoiceExercise,
  ENDING_SOUND: ChoiceExercise,
  TRUE_FALSE: TrueFalseExercise,
  FILL_BLANK: TextExercise,
  TRANSFORM: TextExercise,
  TRANSLATE_VI_EN: TextExercise,
  DICTATION: TextExercise,
  LISTEN_NUMBERS: TextExercise,
  REORDER: ReorderExercise,
  ERROR_CORRECTION: ErrorCorrectionExercise,
  LISTEN_GAP_FILL: GapFillExercise,
  WORD_STRESS: WordStressExercise,
  SPEAK: SpeakExercise,
};

export function ExerciseRenderer(props: ExerciseRendererProps) {
  const Component = renderers[props.exercise.type];
  if (!Component) {
    return (
      <p className="rounded-xl bg-amber-50 p-4 text-sm text-amber-800">
        Dạng bài "{props.exercise.type}" chưa được hỗ trợ.
      </p>
    );
  }
  return <Component {...props} />;
}

// ---------------------------------------------------------------------------

function ExerciseAudio({ exercise }: { exercise: Exercise }) {
  const audioUrl = str(exercise.payload, "audioUrl");
  if (!audioUrl) {
    return null;
  }
  return (
    <div className="mb-5 flex justify-center">
      <AudioButton audioUrl={audioUrl} label="Nghe" />
    </div>
  );
}

function ChoiceExercise({ exercise, result, onChange }: ExerciseRendererProps) {
  const [selected, setSelected] = useState<string | null>(null);
  const options = choiceOptions(exercise);
  const answered = result?.answer as { optionId?: string } | undefined;
  const chosen = answered?.optionId ?? selected;
  const ipa = str(exercise.payload, "ipa");
  const word = str(exercise.payload, "word");
  const hasImages = options.some((option) => option.imageUrl);

  return (
    <div>
      <ExerciseAudio exercise={exercise} />
      {ipa || word ? (
        <p className="mb-5 text-center text-2xl font-bold tracking-normal">
          {word ? <span>{word} </span> : null}
          {ipa ? <span className="font-mono text-primary">{ipa}</span> : null}
        </p>
      ) : null}
      <div className={cn("grid gap-3", hasImages ? "grid-cols-2" : "sm:grid-cols-2")}>
        {options.map((option) => {
          const isChosen = chosen === option.id;
          const isCorrect = result && result.correctAnswer === option.id;
          return (
            <button
              className={cn(
                "flex min-h-14 cursor-pointer flex-col items-center justify-center gap-1 rounded-2xl border-2 px-4 py-3 text-center text-base font-semibold transition-colors disabled:cursor-default",
                !result && isChosen && "border-primary bg-primary/5",
                !result && !isChosen && "border-border bg-white hover:border-primary/40",
                result && isCorrect && "border-emerald-500 bg-emerald-50",
                result && isChosen && !isCorrect && "border-red-400 bg-red-50",
                result && !isChosen && !isCorrect && "border-border bg-white opacity-60",
              )}
              disabled={Boolean(result)}
              key={option.id}
              onClick={() => {
                setSelected(option.id);
                onChange({ optionId: option.id });
              }}
              type="button"
            >
              {option.imageUrl ? (
                <img
                  alt={option.label ?? ""}
                  className="mb-1 h-24 w-full rounded-xl object-cover"
                  src={option.imageUrl}
                />
              ) : null}
              <span>{optionLabel(option)}</span>
              {option.word && option.ipa ? (
                <span className="font-mono text-sm font-normal text-muted-foreground">
                  {option.ipa}
                </span>
              ) : null}
            </button>
          );
        })}
      </div>
    </div>
  );
}

function TrueFalseExercise({ exercise, result, onChange }: ExerciseRendererProps) {
  const [selected, setSelected] = useState<boolean | null>(null);
  const answered = result?.answer as { value?: boolean } | undefined;
  const chosen = answered?.value ?? selected;

  return (
    <div>
      <ExerciseAudio exercise={exercise} />
      <p className="mb-5 rounded-2xl bg-muted/50 p-4 text-center text-lg font-semibold">
        {str(exercise.payload, "statement")}
      </p>
      <div className="grid grid-cols-2 gap-3">
        {[true, false].map((value) => {
          const isChosen = chosen === value;
          const isCorrect = result && result.correctAnswer === value;
          return (
            <button
              className={cn(
                "flex h-14 cursor-pointer items-center justify-center gap-2 rounded-2xl border-2 text-base font-semibold disabled:cursor-default",
                !result && isChosen && "border-primary bg-primary/5",
                !result && !isChosen && "border-border bg-white hover:border-primary/40",
                result && isCorrect && "border-emerald-500 bg-emerald-50",
                result && isChosen && !isCorrect && "border-red-400 bg-red-50",
              )}
              disabled={Boolean(result)}
              key={String(value)}
              onClick={() => {
                setSelected(value);
                onChange({ value });
              }}
              type="button"
            >
              {value ? <Check className="size-5" /> : <X className="size-5" />}
              {value ? "Đúng" : "Sai"}
            </button>
          );
        })}
      </div>
    </div>
  );
}

function TextExercise({ exercise, result, onChange, onSubmit }: ExerciseRendererProps) {
  const [text, setText] = useState("");
  const answered = result?.answer as { text?: string } | undefined;
  const sentence = str(exercise.payload, "sentence");
  const sourceVi = str(exercise.payload, "sourceVi");
  const instruction = str(exercise.payload, "instruction");
  const isLong = exercise.type === "DICTATION" || exercise.type === "TRANSLATE_VI_EN" || exercise.type === "TRANSFORM";

  return (
    <div>
      <ExerciseAudio exercise={exercise} />
      {sentence ? (
        <p className="mb-4 rounded-2xl bg-muted/50 p-4 text-center text-lg font-semibold">
          {sentence}
        </p>
      ) : null}
      {sourceVi ? (
        <p className="mb-4 rounded-2xl bg-muted/50 p-4 text-center text-lg font-semibold">
          {sourceVi}
        </p>
      ) : null}
      {instruction ? (
        <p className="mb-3 text-center text-sm font-semibold uppercase text-primary">
          {transformLabel(instruction)}
        </p>
      ) : null}
      <Input
        autoComplete="off"
        autoFocus
        className={cn("h-14 text-lg", isLong && "text-base")}
        disabled={Boolean(result)}
        onChange={(event) => {
          setText(event.target.value);
          onChange(event.target.value.trim() ? { text: event.target.value } : null);
        }}
        onKeyDown={(event) => {
          if (event.key === "Enter" && text.trim()) {
            onSubmit();
          }
        }}
        placeholder={exercise.type === "LISTEN_NUMBERS" ? "Ví dụ: 7:30" : "Nhập câu trả lời…"}
        spellCheck={false}
        value={answered?.text ?? text}
      />
    </div>
  );
}

function transformLabel(instruction: string) {
  const labels: Record<string, string> = {
    NEGATIVE: "Chuyển sang câu phủ định",
    QUESTION: "Chuyển sang câu hỏi",
    PAST: "Chuyển sang thì quá khứ",
    THIRD_PERSON: "Đổi chủ ngữ sang ngôi thứ ba",
  };
  return labels[instruction] ?? instruction;
}

function ReorderExercise({ exercise, result, onChange }: ExerciseRendererProps) {
  const tokens = strList(exercise.payload, "tokens");
  const [picked, setPicked] = useState<number[]>([]);
  const answered = result?.answer as { tokens?: string[] } | undefined;
  const sentence = answered?.tokens ?? picked.map((index) => tokens[index]);

  function update(next: number[]) {
    setPicked(next);
    onChange(next.length === tokens.length ? { tokens: next.map((index) => tokens[index]) } : null);
  }

  return (
    <div>
      <div className="mb-4 flex min-h-16 flex-wrap items-center gap-2 rounded-2xl border-2 border-dashed border-border bg-white p-3">
        {sentence.length === 0 ? (
          <span className="text-sm text-muted-foreground">Chạm vào các từ bên dưới theo đúng thứ tự</span>
        ) : null}
        {sentence.map((token, position) => (
          <button
            className="cursor-pointer rounded-xl bg-primary px-3 py-2 text-base font-semibold text-primary-foreground disabled:cursor-default"
            disabled={Boolean(result)}
            key={`${token}-${position}`}
            onClick={() => update(picked.filter((_, i) => i !== position))}
            type="button"
          >
            {token}
          </button>
        ))}
      </div>
      {!result ? (
        <div className="flex flex-wrap justify-center gap-2">
          {tokens.map((token, index) => (
            <button
              className="cursor-pointer rounded-xl border-2 border-border bg-white px-3 py-2 text-base font-semibold hover:border-primary/40 disabled:invisible"
              disabled={picked.includes(index)}
              key={`${token}-${index}`}
              onClick={() => update([...picked, index])}
              type="button"
            >
              {token}
            </button>
          ))}
          {picked.length > 0 ? (
            <button
              aria-label="Làm lại"
              className="cursor-pointer rounded-xl px-3 py-2 text-muted-foreground hover:bg-muted"
              onClick={() => update([])}
              type="button"
            >
              <RotateCcw className="size-4" />
            </button>
          ) : null}
        </div>
      ) : null}
    </div>
  );
}

function ErrorCorrectionExercise({ exercise, result, onChange, onSubmit }: ExerciseRendererProps) {
  const tokens = strList(exercise.payload, "tokens");
  const [index, setIndex] = useState<number | null>(null);
  const [text, setText] = useState("");
  const answered = result?.answer as { index?: number; text?: string } | undefined;
  const chosen = answered?.index ?? index;

  function update(nextIndex: number | null, nextText: string) {
    setIndex(nextIndex);
    setText(nextText);
    onChange(nextIndex !== null && nextText.trim() ? { index: nextIndex, text: nextText } : null);
  }

  return (
    <div>
      <p className="mb-2 text-center text-sm text-muted-foreground">1. Chạm vào từ sai</p>
      <div className="mb-5 flex flex-wrap justify-center gap-2">
        {tokens.map((token, i) => (
          <button
            className={cn(
              "cursor-pointer rounded-xl border-2 px-3 py-2 text-lg font-semibold disabled:cursor-default",
              chosen === i ? "border-red-400 bg-red-50 text-red-700" : "border-transparent bg-muted/50 hover:border-primary/40",
            )}
            disabled={Boolean(result)}
            key={`${token}-${i}`}
            onClick={() => update(i, text)}
            type="button"
          >
            {token}
          </button>
        ))}
      </div>
      <p className="mb-2 text-center text-sm text-muted-foreground">2. Viết từ đúng</p>
      <Input
        autoComplete="off"
        className="h-12 text-lg"
        disabled={Boolean(result) || chosen === null}
        onChange={(event) => update(index, event.target.value)}
        onKeyDown={(event) => {
          if (event.key === "Enter" && index !== null && text.trim()) {
            onSubmit();
          }
        }}
        placeholder={chosen === null ? "Chọn từ sai trước" : "Từ đúng…"}
        spellCheck={false}
        value={answered?.text ?? text}
      />
    </div>
  );
}

const blankPattern = /\{\{(\w+)\}\}/g;

function GapFillExercise({ exercise, result, onChange }: ExerciseRendererProps) {
  const text = str(exercise.payload, "text");
  const blanks = Array.from(text.matchAll(blankPattern), (match) => match[1]);
  const [values, setValues] = useState<Record<string, string>>({});
  const answered = result?.answer as { blanks?: Record<string, string> } | undefined;
  const shown = answered?.blanks ?? values;
  const parts = text.split(blankPattern);

  function update(blank: string, value: string) {
    const next = { ...values, [blank]: value };
    setValues(next);
    onChange(blanks.every((key) => next[key]?.trim()) ? { blanks: next } : null);
  }

  return (
    <div>
      <ExerciseAudio exercise={exercise} />
      <p className="rounded-2xl bg-muted/50 p-4 text-center text-lg leading-loose">
        {parts.map((part, i) =>
          i % 2 === 1 ? (
            <input
              aria-label={`Chỗ trống ${part}`}
              className="mx-1 w-28 rounded-lg border-2 border-border bg-white px-2 py-1 text-center font-semibold outline-none focus:border-primary disabled:bg-muted"
              disabled={Boolean(result)}
              key={i}
              onChange={(event) => update(part, event.target.value)}
              spellCheck={false}
              value={shown[part] ?? ""}
            />
          ) : (
            <span key={i}>{part}</span>
          ),
        )}
      </p>
    </div>
  );
}

function WordStressExercise({ exercise, result, onChange }: ExerciseRendererProps) {
  const syllables = strList(exercise.payload, "syllables");
  const [selected, setSelected] = useState<number | null>(null);
  const answered = result?.answer as { stressIndex?: number } | undefined;
  const chosen = answered?.stressIndex ?? selected;
  const word = str(exercise.payload, "word");

  return (
    <div>
      <div className="mb-5 flex justify-center">
        <AudioButton audioUrl={str(exercise.payload, "audioUrl")} fallbackText={word} label={`Nghe "${word}"`} />
      </div>
      <div className="flex justify-center gap-2">
        {syllables.map((syllable, i) => {
          const isCorrect = result && result.correctAnswer === i;
          return (
            <button
              className={cn(
                "min-w-20 cursor-pointer rounded-2xl border-2 px-4 py-4 text-xl font-bold disabled:cursor-default",
                !result && chosen === i && "border-primary bg-primary/5",
                !result && chosen !== i && "border-border bg-white hover:border-primary/40",
                result && isCorrect && "border-emerald-500 bg-emerald-50",
                result && chosen === i && !isCorrect && "border-red-400 bg-red-50",
              )}
              disabled={Boolean(result)}
              key={`${syllable}-${i}`}
              onClick={() => {
                setSelected(i);
                onChange({ stressIndex: i });
              }}
              type="button"
            >
              {chosen === i ? syllable.toUpperCase() : syllable}
            </button>
          );
        })}
      </div>
    </div>
  );
}

function SpeakExercise({ exercise, result, onChange }: ExerciseRendererProps) {
  const text = str(exercise.payload, "text");
  const [done, setDone] = useState(false);

  return (
    <div className="text-center">
      <p className="mb-2 text-2xl font-bold tracking-normal">{text}</p>
      {str(exercise.payload, "ipa") ? (
        <p className="mb-4 font-mono text-primary">{str(exercise.payload, "ipa")}</p>
      ) : null}
      <div className="mb-5 flex justify-center">
        <AudioButton audioUrl={str(exercise.payload, "audioUrl")} fallbackText={text} label="Nghe mẫu" />
      </div>
      <button
        className={cn(
          "inline-flex cursor-pointer items-center gap-2 rounded-2xl border-2 px-5 py-3 font-semibold disabled:cursor-default",
          done || result ? "border-emerald-500 bg-emerald-50 text-emerald-700" : "border-border bg-white hover:border-primary/40",
        )}
        disabled={Boolean(result)}
        onClick={() => {
          setDone(true);
          onChange({ selfCheck: true });
        }}
        type="button"
      >
        <Mic className="size-5" />
        Tôi đã đọc to câu này
      </button>
      <p className="mt-3 text-xs text-muted-foreground">
        Chấm phát âm tự động sẽ có ở bản sau. Hiện tại bạn tự đánh giá.
      </p>
    </div>
  );
}
