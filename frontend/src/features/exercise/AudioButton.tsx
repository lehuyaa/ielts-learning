import { Volume2 } from "lucide-react";

import { playPronunciation } from "@/lib/pronunciation";
import { cn } from "@/lib/utils";

type AudioButtonProps = {
  audioUrl?: string;
  /** Spoken by the browser when there is no audio file. */
  fallbackText?: string;
  label?: string;
  size?: "sm" | "lg";
  className?: string;
};

export function AudioButton({
  audioUrl,
  fallbackText,
  label = "Nghe",
  size = "lg",
  className,
}: AudioButtonProps) {
  const disabled = !audioUrl && !fallbackText;

  return (
    <button
      aria-label={label}
      className={cn(
        "inline-flex cursor-pointer items-center gap-2 rounded-full bg-primary/10 font-semibold text-primary transition-colors hover:bg-primary/20 disabled:cursor-not-allowed disabled:opacity-40",
        size === "lg" ? "px-5 py-3 text-base" : "px-3 py-1.5 text-sm",
        className,
      )}
      disabled={disabled}
      onClick={() => playPronunciation(fallbackText ?? "", audioUrl)}
      title={disabled ? "Chưa có audio" : label}
      type="button"
    >
      <Volume2 className={size === "lg" ? "size-5" : "size-4"} aria-hidden="true" />
      {label}
    </button>
  );
}
