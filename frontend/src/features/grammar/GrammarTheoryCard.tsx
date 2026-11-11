import { ArrowRight, Check, X } from "lucide-react";

import { AudioButton } from "@/features/exercise/AudioButton";
import type { GrammarPointCard } from "@/types/grammar";

type GrammarTheoryCardProps = {
  point: GrammarPointCard;
};

/** Theory card: Vietnamese explanation, formula, examples, common mistakes. */
export function GrammarTheoryCard({ point }: GrammarTheoryCardProps) {
  return (
    <article className="space-y-5 rounded-3xl border bg-white p-5 shadow-sm sm:p-7">
      <header>
        <p className="text-xs font-bold uppercase tracking-wide text-primary">
          Ngữ pháp {point.cefrLevel}
        </p>
        <h2 className="mt-1 text-2xl font-bold tracking-normal">{point.titleVi || point.title}</h2>
        {point.titleVi ? <p className="text-sm text-muted-foreground">{point.title}</p> : null}
      </header>

      <p className="text-base leading-relaxed">{point.summaryVi}</p>

      {point.formula ? (
        <section>
          <h3 className="mb-2 text-sm font-bold text-muted-foreground">Công thức</h3>
          <pre className="overflow-x-auto whitespace-pre-wrap rounded-2xl bg-primary/5 p-4 font-mono text-sm leading-relaxed text-foreground">
            {point.formula}
          </pre>
        </section>
      ) : null}

      {point.examples.length > 0 ? (
        <section>
          <h3 className="mb-2 text-sm font-bold text-muted-foreground">Ví dụ</h3>
          <ul className="space-y-2">
            {point.examples.map((example) => (
              <li className="flex items-start gap-3 rounded-2xl border p-3" key={example.en}>
                <AudioButton
                  audioUrl={example.audioUrl}
                  fallbackText={example.en}
                  label="Nghe"
                  size="sm"
                />
                <div>
                  <p className="font-semibold">{example.en}</p>
                  <p className="text-sm text-muted-foreground">{example.vi}</p>
                </div>
              </li>
            ))}
          </ul>
        </section>
      ) : null}

      {point.commonMistakes.length > 0 ? (
        <section>
          <h3 className="mb-2 text-sm font-bold text-muted-foreground">Lỗi người Việt hay mắc</h3>
          <ul className="space-y-2">
            {point.commonMistakes.map((mistake) => (
              <li className="rounded-2xl bg-muted/40 p-3" key={mistake.wrong}>
                <p className="flex flex-wrap items-center gap-2 text-sm">
                  <span className="inline-flex items-center gap-1 text-red-600 line-through decoration-red-300">
                    <X className="size-4 no-underline" aria-hidden="true" />
                    {mistake.wrong}
                  </span>
                  <ArrowRight className="size-4 text-muted-foreground" aria-hidden="true" />
                  <span className="inline-flex items-center gap-1 font-semibold text-emerald-700">
                    <Check className="size-4" aria-hidden="true" />
                    {mistake.right}
                  </span>
                </p>
                {mistake.noteVi ? <p className="mt-1 text-xs text-muted-foreground">{mistake.noteVi}</p> : null}
              </li>
            ))}
          </ul>
        </section>
      ) : null}
    </article>
  );
}
