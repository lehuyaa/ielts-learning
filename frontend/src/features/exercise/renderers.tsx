import { ArrowDown, Check, CornerDownLeft, Mic, RotateCcw, Snail, X } from "lucide-react";
import { useState, type ReactNode } from "react";

import { Input } from "@/components/ui/input";
import { playPronunciation } from "@/lib/pronunciation";
import { cn } from "@/lib/utils";
import type {
  AnswerResult,
  Exercise,
  ExerciseAnswer,
  ExerciseType,
} from "@/types/activity";

import { AudioButton } from "./AudioButton";
import { DictationExercise } from "./DictationExercise";
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
  FILL_BLANK: FillBlankExercise,
  TRANSFORM: WritingExercise,
  TRANSLATE_VI_EN: WritingExercise,
  DICTATION: DictationExercise,
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
    <div className="mb-5 flex flex-wrap justify-center gap-2">
      <AudioButton audioUrl={audioUrl} label="Nghe" />
      <button
        className="inline-flex cursor-pointer items-center gap-2 rounded-full border px-4 py-3 text-sm font-semibold text-muted-foreground hover:border-primary/40 hover:text-foreground"
        onClick={() => playPronunciation("", audioUrl, { rate: 0.75 })}
        type="button"
      >
        <Snail className="size-4" /> Nghe chậm
      </button>
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
              <span className={cn(isIpaOnly(option) && "font-mono text-xl")}>{optionLabel(option)}</span>
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

/** WORD_TO_IPA and ENDING_SOUND options are transcriptions. */
function isIpaOnly(option: { word?: string; text?: string; ipa?: string }) {
  return (!option.word && Boolean(option.ipa)) || /^\/.*\/$/.test(option.text ?? "");
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
  const isLong = exercise.type === "DICTATION";

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
        className={cn("h-14 text-lg disabled:opacity-100", isLong && "text-base", resultTone(result))}
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

/** Border/background for an answered input. */
function resultTone(result?: AnswerResult) {
  if (!result) return "";
  if (result.isCorrect) return "border-emerald-500 bg-emerald-50 text-emerald-800";
  if (result.scoreRatio > 0) return "border-amber-400 bg-amber-50 text-amber-800";
  return "border-red-400 bg-red-50 text-red-700";
}

const fillBlankPattern = /_{2,}/;

/** FILL_BLANK: the input sits inside the sentence where the blank is. */
function FillBlankExercise(props: ExerciseRendererProps) {
  const { exercise, result, onChange, onSubmit } = props;
  const sentence = str(exercise.payload, "sentence");
  const hint = str(exercise.payload, "hint");
  const [text, setText] = useState("");
  const answered = result?.answer as { text?: string } | undefined;

  if (!fillBlankPattern.test(sentence)) {
    return <TextExercise {...props} />;
  }
  const [before, after] = splitOnce(sentence, fillBlankPattern);
  const value = answered?.text ?? text;

  return (
    <div>
      <ExerciseAudio exercise={exercise} />
      <p className="rounded-2xl bg-muted/50 px-4 py-6 text-center text-xl leading-[2.6rem]">
        <span>{before}</span>
        <input
          aria-label="Chỗ trống"
          autoCapitalize="off"
          autoComplete="off"
          autoFocus
          className={cn(
            "mx-1 rounded-lg border-2 border-primary/40 bg-white px-2 py-0.5 text-center font-semibold outline-none focus:border-primary disabled:opacity-100",
            resultTone(result),
          )}
          disabled={Boolean(result)}
          onChange={(event) => {
            setText(event.target.value);
            onChange(event.target.value.trim() ? { text: event.target.value } : null);
          }}
          onKeyDown={(event) => {
            if (event.key === "Enter" && text.trim()) onSubmit();
          }}
          spellCheck={false}
          style={{ width: `${Math.max(5, value.length + 2)}ch` }}
          value={value}
        />
        <span>{after}</span>
      </p>
      {hint ? <p className="mt-3 text-center text-sm text-muted-foreground">Gợi ý: {hint}</p> : null}
    </div>
  );
}

function splitOnce(value: string, pattern: RegExp): [string, string] {
  const match = pattern.exec(value);
  if (!match) return [value, ""];
  return [value.slice(0, match.index), value.slice(match.index + match[0].length)];
}

/** TRANSFORM and TRANSLATE_VI_EN: a source card, the task, and a multi-line answer. */
function WritingExercise({ exercise, result, onChange, onSubmit }: ExerciseRendererProps) {
  const [text, setText] = useState("");
  const answered = result?.answer as { text?: string } | undefined;
  const isTranslate = exercise.type === "TRANSLATE_VI_EN";
  const source = str(exercise.payload, isTranslate ? "sourceVi" : "sentence");
  const instruction = str(exercise.payload, "instruction");

  return (
    <div>
      <ExerciseAudio exercise={exercise} />
      <div className="rounded-2xl bg-muted/50 p-4 text-center">
        <p className="mb-1 text-xs font-semibold uppercase text-muted-foreground">
          {isTranslate ? "Tiếng Việt" : "Câu gốc"}
        </p>
        <p className="text-lg font-semibold">{source}</p>
      </div>
      <div className="my-3 flex items-center justify-center gap-2 text-sm font-semibold text-primary">
        <ArrowDown className="size-4" />
        {isTranslate ? "Viết câu tiếng Anh" : transformLabel(instruction)}
      </div>
      <textarea
        autoCapitalize="sentences"
        autoComplete="off"
        autoFocus
        className={cn(
          "min-h-24 w-full resize-none rounded-2xl border-2 border-input bg-white p-4 text-lg outline-none focus:border-primary disabled:opacity-100",
          resultTone(result),
        )}
        disabled={Boolean(result)}
        onChange={(event) => {
          setText(event.target.value);
          onChange(event.target.value.trim() ? { text: event.target.value } : null);
        }}
        onKeyDown={(event) => {
          // Enter checks the answer; Shift+Enter is not needed for one sentence.
          if (event.key === "Enter") {
            event.preventDefault();
            if (text.trim()) onSubmit();
          }
        }}
        placeholder={isTranslate ? "Ví dụ: I work in Hanoi." : "Viết lại cả câu…"}
        rows={2}
        spellCheck={false}
        value={answered?.text ?? text}
      />
      {!result ? (
        <p className="mt-2 flex items-center justify-end gap-1 text-xs text-muted-foreground">
          <CornerDownLeft className="size-3" />
          Enter để kiểm tra · cách viết khác đáp án vẫn có thể được AI chấm đúng
        </p>
      ) : null}
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
  const complete = picked.length === tokens.length;

  function update(next: number[]) {
    setPicked(next);
    onChange(next.length === tokens.length ? { tokens: next.map((index) => tokens[index]) } : null);
  }

  return (
    <div>
      <div
        className={cn(
          "mb-2 flex min-h-16 flex-wrap items-center gap-2 rounded-2xl border-2 border-dashed bg-white p-3",
          result ? cn("border-solid", resultTone(result)) : complete ? "border-primary" : "border-border",
        )}
      >
        {sentence.length === 0 ? (
          <span className="text-sm text-muted-foreground">Chạm vào các từ bên dưới theo đúng thứ tự</span>
        ) : null}
        {sentence.map((token, position) => (
          <button
            aria-label={`Bỏ "${token}"`}
            className={cn(
              "cursor-pointer rounded-xl px-3 py-2 text-base font-semibold disabled:cursor-default",
              result ? "bg-white text-foreground shadow-sm" : "bg-primary text-primary-foreground hover:bg-primary/85",
            )}
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
        <p className="mb-4 text-center text-xs text-muted-foreground">
          {picked.length}/{tokens.length} từ · chạm vào từ đã chọn để bỏ ra
        </p>
      ) : null}
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

/** ERROR_CORRECTION: tap the wrong word, then type its fix right in the sentence. */
function ErrorCorrectionExercise({ exercise, result, onChange, onSubmit }: ExerciseRendererProps) {
  const tokens = strList(exercise.payload, "tokens");
  const [index, setIndex] = useState<number | null>(null);
  const [text, setText] = useState("");
  const answered = result?.answer as { index?: number; text?: string } | undefined;
  const chosen = answered?.index ?? index;
  const fix = answered?.text ?? text;
  const correct = result?.correctAnswer as { wrongIndex?: number; text?: string } | undefined;

  function update(nextIndex: number | null, nextText: string) {
    setIndex(nextIndex);
    setText(nextText);
    onChange(nextIndex !== null && nextText.trim() ? { index: nextIndex, text: nextText } : null);
  }

  return (
    <div>
      <p className="mb-3 text-center text-sm text-muted-foreground">
        {chosen === null ? "1. Chạm vào từ bị sai" : result ? "Câu của bạn" : "2. Viết từ đúng thay cho từ đã chọn"}
      </p>
      <div className="mb-4 flex flex-wrap items-end justify-center gap-2">
        {tokens.map((token, i) => {
          const isChosen = chosen === i;
          const isAnswer = result && correct?.wrongIndex === i;
          return (
            <span className="flex flex-col items-center gap-1" key={`${token}-${i}`}>
              <button
                className={cn(
                  "cursor-pointer rounded-xl border-2 px-3 py-2 text-lg font-semibold disabled:cursor-default",
                  isChosen && !result && "border-red-400 bg-red-50 text-red-700 line-through",
                  !isChosen && !result && "border-transparent bg-muted/50 hover:border-primary/40",
                  result && isChosen && "line-through",
                  result && isChosen && !isAnswer && "border-red-400 bg-red-50 text-red-700",
                  result && isAnswer && "border-emerald-500 bg-emerald-50 text-emerald-800",
                  result && !isChosen && !isAnswer && "border-transparent bg-muted/40 opacity-70",
                )}
                disabled={Boolean(result)}
                onClick={() => update(i, i === index ? text : "")}
                type="button"
              >
                {token}
              </button>
              {isChosen ? (
                <span className={cn("rounded-md px-2 text-base font-semibold", result ? resultTone(result) : "text-primary")}>
                  {fix || "…"}
                </span>
              ) : null}
            </span>
          );
        })}
      </div>
      {!result && chosen !== null ? (
        <Input
          aria-label={`Từ đúng thay cho "${tokens[chosen] ?? ""}"`}
          autoComplete="off"
          autoFocus
          className="h-12 text-lg"
          key={chosen}
          onChange={(event) => update(index, event.target.value)}
          onKeyDown={(event) => {
            if (event.key === "Enter" && index !== null && text.trim()) {
              onSubmit();
            }
          }}
          placeholder={`"${tokens[chosen] ?? ""}" phải sửa thành…`}
          spellCheck={false}
          value={text}
        />
      ) : null}
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
