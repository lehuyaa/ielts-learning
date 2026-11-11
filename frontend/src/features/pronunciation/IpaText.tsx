import { Link } from "react-router-dom";

import { phonemePath, splitIpa } from "@/lib/ipa";
import { cn } from "@/lib/utils";

type IpaTextProps = {
  ipa: string | null | undefined;
  className?: string;
  /** Plain text when false (e.g. inside an exercise). */
  linked?: boolean;
};

/**
 * An IPA transcription where every sound links to its sound card
 * (/pronunciation/:symbol). Falls back to plain text when the value is not a
 * real transcription.
 */
export function IpaText({ ipa, className, linked = true }: IpaTextProps) {
  if (!ipa) {
    return null;
  }
  const segments = linked ? splitIpa(ipa) : null;
  if (!segments) {
    return <span className={cn("font-mono", className)}>{ipa}</span>;
  }

  return (
    <span className={cn("font-mono", className)}>
      {segments.map((segment, index) =>
        segment.symbol ? (
          <Link
            className="rounded-sm underline decoration-current/25 decoration-dotted underline-offset-4 hover:bg-primary/10 hover:decoration-current"
            key={index}
            onClick={(event) => event.stopPropagation()}
            title={`Âm /${segment.symbol}/ — xem cách phát âm`}
            to={phonemePath(segment.symbol)}
          >
            {segment.text}
          </Link>
        ) : (
          <span key={index}>{segment.text}</span>
        ),
      )}
    </span>
  );
}
