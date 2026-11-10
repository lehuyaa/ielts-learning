import type { Exercise } from "@/types/activity";

// Typed accessors for exercise payloads (docs/foundation-spec.md §5).

export type ChoiceOption = {
  id: string;
  text?: string;
  word?: string;
  ipa?: string;
  imageUrl?: string;
  label?: string;
};

export function str(payload: Exercise["payload"], key: string): string {
  const value = payload[key];
  return typeof value === "string" ? value : "";
}

export function strList(payload: Exercise["payload"], key: string): string[] {
  const value = payload[key];
  return Array.isArray(value)
    ? value.filter((item): item is string => typeof item === "string")
    : [];
}

/** Options for choice-style exercises; ENDING_SOUND uses plain strings. */
export function choiceOptions(exercise: Exercise): ChoiceOption[] {
  const value = exercise.payload.options;
  if (!Array.isArray(value)) {
    return [];
  }
  return value.map((option) =>
    typeof option === "string"
      ? { id: option, text: option }
      : (option as ChoiceOption),
  );
}

export function optionLabel(option: ChoiceOption): string {
  return option.text ?? option.word ?? option.ipa ?? option.label ?? option.id;
}

/** Human-readable correct answer for the feedback panel. */
export function formatCorrectAnswer(exercise: Exercise, correct: unknown): string {
  if (correct === null || correct === undefined) {
    return "";
  }

  switch (exercise.type) {
    case "TRUE_FALSE":
      return correct ? "Đúng (True)" : "Sai (False)";
    case "WORD_STRESS": {
      const syllables = strList(exercise.payload, "syllables");
      const index = Number(correct);
      return syllables
        .map((syllable, i) => (i === index ? syllable.toUpperCase() : syllable))
        .join("-");
    }
    case "FILL_BLANK": {
      // "She ___ (work) in a bank." -> "She works in a bank."
      const sentence = str(exercise.payload, "sentence");
      return /_{2,}/.test(sentence)
        ? sentence.replace(/_{2,}(\s*\([^)]*\))?/, String(correct))
        : String(correct);
    }
    case "ERROR_CORRECTION": {
      const value = correct as { wrongIndex: number; text: string };
      const tokens = strList(exercise.payload, "tokens");
      return `"${tokens[value.wrongIndex] ?? ""}" → "${value.text}"`;
    }
    case "LISTEN_GAP_FILL":
      return Object.entries(correct as Record<string, string>)
        .sort(([a], [b]) => Number(a) - Number(b))
        .map(([blank, value]) => `(${blank}) ${value}`)
        .join(", ");
    case "MULTIPLE_CHOICE":
    case "LISTEN_PICK_IMAGE":
    case "MINIMAL_PAIR":
    case "IPA_TO_WORD":
    case "WORD_TO_IPA":
    case "ENDING_SOUND": {
      const option = choiceOptions(exercise).find((item) => item.id === correct);
      return option ? optionLabel(option) : String(correct);
    }
    default:
      return String(correct);
  }
}
