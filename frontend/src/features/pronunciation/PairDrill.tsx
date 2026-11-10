import { CheckCircle2, RotateCcw, Volume2, XCircle } from "lucide-react";
import { useState } from "react";

import { Button } from "@/components/ui/button";
import { playPronunciation } from "@/lib/pronunciation";
import { cn } from "@/lib/utils";
import type { MinimalPair } from "@/types/phoneme";

type Round = { pair: MinimalPair; answer: "a" | "b"; order: ("a" | "b")[] };

const ROUNDS = 10;

function makeRounds(pairs: MinimalPair[]): Round[] {
  const rounds: Round[] = [];
  for (let i = 0; i < Math.min(ROUNDS, Math.max(pairs.length, 6)); i++) {
    const pair = pairs[Math.floor(Math.random() * pairs.length)];
    rounds.push({
      pair,
      answer: Math.random() < 0.5 ? "a" : "b",
      order: Math.random() < 0.5 ? ["a", "b"] : ["b", "a"],
    });
  }
  return rounds;
}

function wordOf(pair: MinimalPair, side: "a" | "b") {
  return side === "a"
    ? { word: pair.a, ipa: pair.aIpa, audioUrl: pair.aAudioUrl }
    : { word: pair.b, ipa: pair.bIpa, audioUrl: pair.bAudioUrl };
}

/**
 * Quick listening drill on a sound card: hear one word of a minimal pair and
 * pick it. Practice only — nothing is saved.
 */
export function PairDrill({ pairs, symbol, onClose }: { pairs: MinimalPair[]; symbol: string; onClose: () => void }) {
  const [rounds, setRounds] = useState(() => makeRounds(pairs));
  const [index, setIndex] = useState(0);
  const [picked, setPicked] = useState<"a" | "b" | null>(null);
  const [correct, setCorrect] = useState(0);

  const round = rounds[index];
  const done = index >= rounds.length;

  function play(target: Round = round) {
    const heard = wordOf(target.pair, target.answer);
    playPronunciation(heard.word, heard.audioUrl);
  }

  function choose(side: "a" | "b") {
    if (picked) return;
    setPicked(side);
    if (side === round.answer) setCorrect((value) => value + 1);
  }

  function next() {
    const nextIndex = index + 1;
    setPicked(null);
    setIndex(nextIndex);
    if (nextIndex < rounds.length) play(rounds[nextIndex]);
  }

  function restart() {
    setRounds(makeRounds(pairs));
    setIndex(0);
    setPicked(null);
    setCorrect(0);
  }

  if (done) {
    return (
      <div className="rounded-3xl border bg-white p-6 text-center">
        <p className="text-4xl font-bold tracking-normal">
          {correct}/{rounds.length}
        </p>
        <p className="mt-1 text-sm text-muted-foreground">câu nghe đúng với âm /{symbol}/</p>
        <div className="mt-4 flex justify-center gap-3">
          <Button onClick={restart} variant="outline">
            <RotateCcw className="size-4" /> Luyện lại
          </Button>
          <Button onClick={onClose}>Xong</Button>
        </div>
      </div>
    );
  }

  return (
    <div className="rounded-3xl border bg-white p-5">
      <div className="mb-4 flex items-center justify-between text-sm text-muted-foreground">
        <span>
          Câu {index + 1}/{rounds.length} · đúng {correct}
        </span>
        <button className="cursor-pointer hover:text-foreground" onClick={onClose} type="button">
          Đóng
        </button>
      </div>
      <div className="mb-5 flex justify-center">
        <button
          className="inline-flex cursor-pointer items-center gap-2 rounded-full bg-primary/10 px-5 py-3 font-semibold text-primary hover:bg-primary/20"
          onClick={() => play()}
          type="button"
        >
          <Volume2 className="size-5" /> {index === 0 && !picked ? "Bấm để nghe" : "Nghe lại"}
        </button>
      </div>
      <div className="grid grid-cols-2 gap-3">
        {round.order.map((side) => {
          const word = wordOf(round.pair, side);
          const isAnswer = side === round.answer;
          return (
            <button
              className={cn(
                "flex cursor-pointer flex-col items-center rounded-2xl border-2 px-3 py-4 text-lg font-semibold disabled:cursor-default",
                !picked && "border-border hover:border-primary/40",
                picked && isAnswer && "border-emerald-500 bg-emerald-50",
                picked === side && !isAnswer && "border-red-400 bg-red-50",
                picked && picked !== side && !isAnswer && "border-border opacity-60",
              )}
              disabled={Boolean(picked)}
              key={side}
              onClick={() => choose(side)}
              type="button"
            >
              {word.word}
              <span className="font-mono text-sm font-normal text-muted-foreground">{word.ipa}</span>
            </button>
          );
        })}
      </div>
      {picked ? (
        <div className="mt-4 flex items-center justify-between gap-3">
          <p className="flex items-center gap-2 font-semibold">
            {picked === round.answer ? (
              <CheckCircle2 className="size-5 text-emerald-600" />
            ) : (
              <XCircle className="size-5 text-red-600" />
            )}
            {picked === round.answer ? "Chính xác!" : `Bạn vừa nghe "${wordOf(round.pair, round.answer).word}"`}
          </p>
          <Button onClick={next}>Tiếp</Button>
        </div>
      ) : null}
    </div>
  );
}
