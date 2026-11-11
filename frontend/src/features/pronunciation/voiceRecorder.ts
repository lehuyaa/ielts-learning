type RecorderEvents = {
  onLevel: (level: number) => void
  onElapsed: (seconds: number) => void
}

const errorMessages: Record<string, string> = {
  NotAllowedError: 'Trình duyệt đang chặn micro. Hãy cho phép micro cho trang này rồi thử lại.',
  SecurityError: 'Trình duyệt đang chặn micro. Hãy cho phép micro cho trang này rồi thử lại.',
  NotFoundError: 'Không tìm thấy micro nào.',
  NotReadableError: 'Micro đang được ứng dụng khác dùng.',
}

export function microphoneError(error: unknown): string {
  if (error instanceof DOMException && errorMessages[error.name]) {
    return errorMessages[error.name]
  }
  return 'Không ghi âm được. Hãy thử lại.'
}

export function canRecord() {
  return typeof window !== 'undefined' && 'MediaRecorder' in window && Boolean(navigator.mediaDevices?.getUserMedia)
}

/**
 * Microphone recording with a live level meter and a hard time limit
 * (Task 71). A class keeps the browser objects out of React state.
 */
export class VoiceRecorder {
  private readonly events: RecorderEvents
  private media: MediaRecorder | null = null
  private stream: MediaStream | null = null
  private context: AudioContext | null = null
  private frame = 0
  private timer = 0
  private chunks: Blob[] = []
  private cancelled = false

  constructor(events: RecorderEvents) {
    this.events = events
  }

  /** Records until stop(), cancel() (resolves null) or maxSeconds. */
  async record(maxSeconds: number): Promise<Blob | null> {
    this.release()
    this.cancelled = false
    this.chunks = []
    const stream = await navigator.mediaDevices.getUserMedia({
      audio: { channelCount: 1, echoCancellation: true, noiseSuppression: true },
    })
    this.stream = stream

    const context = new AudioContext()
    const analyser = context.createAnalyser()
    analyser.fftSize = 512
    context.createMediaStreamSource(stream).connect(analyser)
    this.context = context
    const samples = new Uint8Array(analyser.fftSize)
    const tick = () => {
      analyser.getByteTimeDomainData(samples)
      let peak = 0
      for (const sample of samples) peak = Math.max(peak, Math.abs(sample - 128))
      this.events.onLevel(Math.min(1, peak / 64))
      this.frame = requestAnimationFrame(tick)
    }
    tick()

    const media = new MediaRecorder(stream)
    media.ondataavailable = (event) => {
      if (event.data.size > 0) this.chunks.push(event.data)
    }
    const finished = new Promise<Blob | null>((resolve) => {
      media.onstop = () => {
        const recording = this.cancelled ? null : new Blob(this.chunks, { type: media.mimeType })
        this.release()
        resolve(recording)
      }
    })
    this.media = media
    media.start()

    const startedAt = performance.now()
    this.events.onElapsed(0)
    this.timer = window.setInterval(() => {
      const elapsed = (performance.now() - startedAt) / 1000
      this.events.onElapsed(elapsed)
      if (elapsed >= maxSeconds) this.stop()
    }, 200)
    return finished
  }

  stop() {
    if (this.media && this.media.state !== 'inactive') {
      this.media.stop()
    }
  }

  cancel() {
    this.cancelled = true
    if (this.media && this.media.state !== 'inactive') {
      this.media.stop()
    } else {
      this.release()
    }
  }

  private release() {
    cancelAnimationFrame(this.frame)
    window.clearInterval(this.timer)
    this.stream?.getTracks().forEach((track) => track.stop())
    this.stream = null
    void this.context?.close()
    this.context = null
    this.media = null
    this.events.onLevel(0)
  }
}
