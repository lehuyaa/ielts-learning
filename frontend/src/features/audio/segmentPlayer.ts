import { mediaUrl } from '@/lib/media'
import { speakWord, stopPronunciation } from '@/lib/pronunciation'
import type { PassageCard, PassageSegment } from '@/types/passage'

/** Where a segment's sound comes from, or null to use the browser voice. */
export function segmentSource(passage: Pick<PassageCard, 'audioUrl'>, segment: PassageSegment) {
  if (segment.audioUrl) {
    return { url: segment.audioUrl, start: undefined, end: undefined }
  }
  if (passage.audioUrl && segment.start !== undefined && segment.end !== undefined) {
    return { url: passage.audioUrl, start: segment.start, end: segment.end }
  }
  if (passage.audioUrl && segment.start === undefined) {
    return { url: passage.audioUrl, start: undefined, end: undefined }
  }
  return null
}

type Listener = { onSegment: (index: number) => void; onPlaying: (playing: boolean) => void }

/**
 * Plays passage segments one at a time (optionally continuing to the next),
 * from a segment's own file, a start/end slice of the passage audio, or the
 * browser voice. Works with one reused <audio> element so mobile Safari keeps
 * the user-gesture permission.
 */
export class SegmentPlayer {
  private audio: HTMLAudioElement | null = null
  private playId = 0
  private continueAfter = false
  private listener: Listener
  rate = 1

  constructor(listener: Listener) {
    this.listener = listener
  }

  setRate(rate: number) {
    this.rate = rate
    if (this.audio) this.audio.playbackRate = rate
  }

  stop() {
    this.playId++
    this.continueAfter = false
    if (this.audio) {
      this.audio.onended = null
      this.audio.ontimeupdate = null
      this.audio.onerror = null
      this.audio.pause()
    }
    stopPronunciation()
    this.listener.onPlaying(false)
  }

  play(passage: Pick<PassageCard, 'audioUrl'>, segments: PassageSegment[], index: number, continueAfter: boolean) {
    this.stop()
    const segment = segments[index]
    if (!segment) return
    const playId = ++this.playId
    this.continueAfter = continueAfter
    this.listener.onSegment(index)
    this.listener.onPlaying(true)

    const finished = () => {
      if (playId !== this.playId) return
      if (this.continueAfter && index + 1 < segments.length) {
        this.play(passage, segments, index + 1, true)
      } else {
        this.listener.onPlaying(false)
      }
    }
    const speak = () => {
      if (playId === this.playId) speakWord(segment.text, this.rate, finished)
    }

    const source = segmentSource(passage, segment)
    if (!source) {
      speak()
      return
    }

    const audio = this.audio ?? new Audio()
    this.audio = audio
    const url = new URL(mediaUrl(source.url), window.location.href).toString()
    if (audio.src !== url) {
      audio.src = url
    }
    audio.playbackRate = this.rate
    audio.onended = finished
    audio.onerror = speak
    audio.ontimeupdate = null
    const end = source.end
    if (end !== undefined) {
      audio.ontimeupdate = () => {
        if (audio.currentTime >= end) {
          audio.ontimeupdate = null
          audio.pause()
          finished()
        }
      }
    }
    const start = source.start
    if (start !== undefined) {
      try {
        audio.currentTime = start
      } catch {
        // Metadata not loaded yet; seek once it is.
        audio.onloadedmetadata = () => {
          audio.onloadedmetadata = null
          audio.currentTime = start
        }
      }
    } else if (audio.ended || audio.currentTime > 0) {
      audio.currentTime = 0
    }
    audio.play().catch(speak)
  }
}
