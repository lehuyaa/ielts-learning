import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { ArrowLeft, Eye, EyeOff, Flag, Lightbulb, Mic, RotateCcw, Send, Target, Volume2 } from "lucide-react";
import { useEffect, useRef, useState } from "react";
import { Link, useLocation, useParams } from "react-router-dom";

import { getApiErrorMessage } from "@/api/api";
import { finishSession, getSession, getSessionHint, sendSessionMessage } from "@/api/aiConversation";
import { ErrorState } from "@/components/state/ErrorState";
import { Button } from "@/components/ui/button";
import { useToast } from "@/contexts/toast/useToast";
import { scenariosQueryKey } from "@/features/aiConversation/hooks/useScenarios";
import { useSpeechRecognition } from "@/features/aiConversation/hooks/useSpeechRecognition";
import { SessionReport } from "@/features/speaking/SessionReport";
import { useStartSession } from "@/features/speaking/hooks";
import { cn } from "@/lib/utils";
import type { SessionMessage, SpeakingHint, SpeakingSession } from "@/types/aiConversation";

/** Beginners hear the AI slower (Task 67). */
function speechRate(level: string) {
  return level === "A1" || level === "A2" ? 0.85 : 1;
}

function speak(text: string, rate: number) {
  if (!("speechSynthesis" in window)) return;
  window.speechSynthesis.cancel();
  const utterance = new SpeechSynthesisUtterance(text);
  utterance.lang = "en-US";
  utterance.rate = rate;
  window.speechSynthesis.speak(utterance);
}

function stopSpeaking() {
  if ("speechSynthesis" in window) window.speechSynthesis.cancel();
}

/** Role-play session (Tasks 67-70): chat with the AI in a scenario, hints, finish → report. */
export function SpeakingSessionPage() {
  const sessionId = Number(useParams().sessionId);
  const location = useLocation();
  const returnTo = (location.state as { returnTo?: string } | null)?.returnTo;
  const query = useQuery({ queryKey: ["speaking-session", sessionId], queryFn: () => getSession(sessionId), enabled: sessionId > 0 });

  if (query.isLoading) {
    return <div className="grid min-h-[60vh] place-items-center text-sm text-muted-foreground">Đang mở cuộc hội thoại…</div>;
  }
  if (query.error || !query.data) {
    return (
      <ErrorState
        description={getApiErrorMessage(query.error, "Không mở được cuộc hội thoại này.")}
        onRetry={() => void query.refetch()}
        title="Không mở được cuộc hội thoại"
      />
    );
  }
  return <SessionView initial={query.data} key={query.data.id} returnTo={returnTo} />;
}

function SessionView({ initial, returnTo }: { initial: SpeakingSession; returnTo?: string }) {
  const queryClient = useQueryClient();
  const { showToast } = useToast();
  const [session, setSession] = useState(initial);
  const [messages, setMessages] = useState<SessionMessage[]>(initial.messages);
  const [input, setInput] = useState("");
  const [showAiText, setShowAiText] = useState(false);
  const [hint, setHint] = useState<SpeakingHint | null>(null);
  const [showBrief, setShowBrief] = useState(initial.turnCount === 0);
  const scrollRef = useRef<HTMLDivElement>(null);
  const rate = speechRate(session.scenario.cefrLevel);
  const finished = session.status !== "IN_PROGRESS";
  const turnsLeft = session.maxTurns - session.turnCount;
  const backTo = returnTo ?? (session.activityId ? undefined : "/ai-conversation");
  const restart = useStartSession();

  const speech = useSpeechRecognition((message) => showToast({ tone: "error", title: "Không dùng được micro", description: message }));

  useEffect(() => stopSpeaking, []);
  useEffect(() => {
    scrollRef.current?.scrollTo({ top: scrollRef.current.scrollHeight, behavior: "smooth" });
  }, [messages, hint]);

  const send = useMutation({
    mutationFn: (text: string) => sendSessionMessage(session.id, text),
    onMutate: (text) => {
      setHint(null);
      setMessages((current) => [...current, { id: -Date.now(), role: "user", content: text, createdAt: new Date().toISOString() }]);
      setInput("");
    },
    onSuccess: (turn) => {
      setMessages((current) => [...current, turn.reply]);
      setSession((current) => ({ ...current, turnCount: turn.turnCount }));
      speak(turn.reply.content, rate);
    },
    onError: (error, text) => {
      setMessages((current) => current.filter((message) => !(message.id < 0 && message.content === text)));
      setInput((current) => current || text);
      showToast({ tone: "error", title: "Chưa gửi được", description: getApiErrorMessage(error, "Hãy thử lại.") });
    },
  });

  const askHint = useMutation({
    mutationFn: () => getSessionHint(session.id),
    onSuccess: (result) => {
      setHint(result);
      setSession((current) => ({ ...current, hintCount: result.hintCount }));
    },
    onError: (error) => showToast({ tone: "error", title: "Không lấy được gợi ý", description: getApiErrorMessage(error, "Hãy thử lại.") }),
  });

  const finish = useMutation({
    mutationFn: () => finishSession(session.id),
    onSuccess: (result) => {
      stopSpeaking();
      setSession(result);
      void queryClient.invalidateQueries({ queryKey: ["activity"] });
      void queryClient.invalidateQueries({ queryKey: scenariosQueryKey });
    },
    onError: (error) => showToast({ tone: "error", title: "Chưa kết thúc được", description: getApiErrorMessage(error, "Hãy thử lại.") }),
  });

  function submit() {
    speech.stop();
    const text = input.trim();
    if (!text || send.isPending || turnsLeft <= 0) return;
    send.mutate(text);
  }

  function toggleMic() {
    if (speech.isListening) {
      speech.stop();
      return;
    }
    stopSpeaking();
    const typed = input.trim();
    speech.start((transcript) => setInput([typed, transcript].filter(Boolean).join(" ")));
  }

  const brief = session.scenario;

  return (
    <div className="mx-auto flex max-w-3xl flex-col gap-4 py-2">
      <header className="flex items-center justify-between gap-3">
        {backTo ? (
          <Link className="inline-flex items-center gap-1 text-sm font-semibold text-muted-foreground hover:text-foreground" to={backTo}>
            <ArrowLeft className="size-4" /> {returnTo ? "Bài học" : "Tình huống"}
          </Link>
        ) : (
          <span />
        )}
        <div className="text-center">
          <p className="font-bold">{brief.title}</p>
          <p className="text-xs text-muted-foreground">
            {finished ? "Đã kết thúc" : `Lượt ${Math.min(session.turnCount + 1, session.maxTurns)}/${session.maxTurns}`} · {brief.cefrLevel}
          </p>
        </div>
        <button
          className="inline-flex cursor-pointer items-center gap-1 text-sm font-semibold text-muted-foreground hover:text-foreground"
          onClick={() => setShowAiText((value) => !value)}
          type="button"
        >
          {showAiText ? <EyeOff className="size-4" /> : <Eye className="size-4" />}
          {showAiText ? "Ẩn chữ" : "Hiện chữ"}
        </button>
      </header>

      <section className="rounded-3xl border bg-white p-4">
        <button className="flex w-full cursor-pointer items-center justify-between text-left" onClick={() => setShowBrief((value) => !value)} type="button">
          <span className="flex items-center gap-2 font-bold">
            <Target className="size-4 text-primary" /> Nhiệm vụ của bạn
          </span>
          <span className="text-xs text-muted-foreground">{showBrief ? "Thu gọn" : "Xem"}</span>
        </button>
        {showBrief ? (
          <div className="mt-3 space-y-2 text-sm">
            <p className="text-muted-foreground">{brief.description}</p>
            <p>
              <span className="font-semibold">AI đóng vai:</span> {brief.role}
            </p>
            <p>
              <span className="font-semibold">Bạn cần:</span> {brief.task}
            </p>
            {brief.targetPatterns.length > 0 ? (
              <div className="flex flex-wrap gap-2 pt-1">
                {brief.targetPatterns.map((pattern) => (
                  <button
                    className="cursor-pointer rounded-full border border-primary/30 bg-primary/5 px-3 py-1 text-sm font-medium text-primary hover:bg-primary/10 disabled:cursor-default"
                    disabled={finished}
                    key={pattern}
                    onClick={() => setInput((current) => (current ? `${current} ${pattern}` : pattern))}
                    type="button"
                  >
                    {pattern}
                  </button>
                ))}
              </div>
            ) : null}
          </div>
        ) : null}
      </section>

      <div className="max-h-[55vh] min-h-64 space-y-3 overflow-y-auto rounded-3xl border bg-muted/20 p-4" ref={scrollRef}>
        {messages.map((message, index) => (
          <Bubble
            blurred={!showAiText && message.role === "assistant" && index > 0 && !finished}
            key={message.id}
            message={message}
            onPlay={() => speak(message.content, rate)}
          />
        ))}
        {send.isPending ? (
          <div className="flex justify-start">
            <div className="rounded-2xl rounded-bl-md border bg-white px-4 py-3 text-sm text-muted-foreground">Đang trả lời…</div>
          </div>
        ) : null}
        {hint ? (
          <div className="rounded-2xl border border-amber-200 bg-amber-50 p-3 text-sm">
            <p className="flex items-center gap-2 font-semibold text-amber-800">
              <Lightbulb className="size-4" /> {hint.hintVi}
            </p>
            <div className="mt-2 flex flex-wrap items-center gap-2">
              <span className="font-medium">“{hint.sayThis}”</span>
              <button className="cursor-pointer rounded-lg p-1 text-amber-800 hover:bg-amber-100" onClick={() => speak(hint.sayThis, rate)} type="button" aria-label="Nghe câu gợi ý">
                <Volume2 className="size-4" />
              </button>
              <button className="cursor-pointer rounded-lg bg-white px-2 py-1 text-xs font-semibold text-amber-800 hover:bg-amber-100" onClick={() => setInput(hint.sayThis)} type="button">
                Dùng câu này
              </button>
            </div>
          </div>
        ) : null}
      </div>

      {finished && session.report ? (
        <>
          <SessionReport report={session.report} />
          <div className="flex flex-wrap justify-center gap-3">
            <Button
              disabled={restart.isPending}
              onClick={() => restart.mutate({ scenarioSlug: brief.slug, activityId: session.activityId ?? undefined, fresh: true, returnTo })}
              variant="outline"
            >
              <RotateCcw className="size-4" /> Nói lại
            </Button>
            {returnTo ? (
              <Button asChild>
                <Link to={returnTo}>Quay lại bài học</Link>
              </Button>
            ) : (
              <Button asChild>
                <Link to="/ai-conversation">Chọn tình huống khác</Link>
              </Button>
            )}
          </div>
        </>
      ) : null}

      {!finished ? (
        <section className="rounded-3xl border bg-white p-4">
          {turnsLeft <= 0 ? (
            <p className="mb-3 rounded-2xl bg-emerald-50 p-3 text-center text-sm font-semibold text-emerald-700">
              Bạn đã nói đủ {session.maxTurns} lượt. Bấm "Kết thúc & xem báo cáo".
            </p>
          ) : null}
          <form
            className="flex items-center gap-2"
            onSubmit={(event) => {
              event.preventDefault();
              submit();
            }}
          >
            <button
              aria-label={speech.isListening ? "Dừng ghi âm" : "Nói"}
              className={cn(
                "grid size-12 shrink-0 cursor-pointer place-items-center rounded-full text-white",
                speech.isListening ? "animate-pulse bg-red-500" : "bg-primary hover:bg-primary/90",
              )}
              disabled={turnsLeft <= 0}
              onClick={toggleMic}
              type="button"
            >
              <Mic className="size-5" />
            </button>
            <input
              className="h-12 min-w-0 flex-1 rounded-full border bg-muted/30 px-4 outline-none focus:border-primary"
              disabled={turnsLeft <= 0}
              onChange={(event) => setInput(event.target.value)}
              placeholder={speech.isListening ? "Đang nghe…" : "Nói hoặc gõ câu trả lời…"}
              value={input}
            />
            <button
              aria-label="Gửi"
              className="grid size-12 shrink-0 cursor-pointer place-items-center rounded-full bg-primary text-white disabled:opacity-40"
              disabled={!input.trim() || send.isPending || turnsLeft <= 0}
              type="submit"
            >
              <Send className="size-5" />
            </button>
          </form>
          <div className="mt-3 flex flex-wrap items-center justify-between gap-2">
            {session.hintsAllowed ? (
              <Button disabled={askHint.isPending || send.isPending} onClick={() => askHint.mutate()} size="sm" variant="outline">
                <Lightbulb className="size-4" /> {askHint.isPending ? "Đang nghĩ…" : "Gợi ý"}
                {session.hintCount > 0 ? <span className="text-xs text-muted-foreground">({session.hintCount})</span> : null}
              </Button>
            ) : (
              <span className="text-xs text-muted-foreground">Bài kiểm tra: không dùng gợi ý.</span>
            )}
            <Button
              disabled={session.turnCount === 0 || finish.isPending || send.isPending}
              onClick={() => finish.mutate()}
              size="sm"
              variant={turnsLeft <= 0 ? "default" : "outline"}
            >
              <Flag className="size-4" /> {finish.isPending ? "AI đang chấm…" : "Kết thúc & xem báo cáo"}
            </Button>
          </div>
        </section>
      ) : null}
    </div>
  );
}

function Bubble({ message, blurred, onPlay }: { message: SessionMessage; blurred: boolean; onPlay: () => void }) {
  const mine = message.role === "user";
  return (
    <div className={cn("flex", mine ? "justify-end" : "justify-start")}>
      <div
        className={cn(
          "flex max-w-[80%] items-start gap-2 rounded-2xl px-4 py-2.5",
          mine ? "rounded-br-md bg-primary text-primary-foreground" : "rounded-bl-md border bg-white",
        )}
      >
        <span className={cn("transition-[filter]", blurred && "select-none blur-sm")}>{message.content}</span>
        {!mine ? (
          <button aria-label="Nghe lại" className="mt-0.5 shrink-0 cursor-pointer text-primary hover:text-primary/70" onClick={onPlay} type="button">
            <Volume2 className="size-4" />
          </button>
        ) : null}
      </div>
    </div>
  );
}
