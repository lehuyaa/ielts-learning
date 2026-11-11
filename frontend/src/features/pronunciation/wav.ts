const TARGET_RATE = 16000

/**
 * Converts a browser recording (WebM/Opus, MP4/AAC…) to the 16 kHz mono
 * 16-bit PCM WAV that pronunciation scoring expects (spec §7).
 */
export async function toWav16k(recording: Blob, maxSeconds: number): Promise<Blob> {
  const buffer = await recording.arrayBuffer()
  const decoder = new AudioContext()
  let decoded: AudioBuffer
  try {
    decoded = await decoder.decodeAudioData(buffer)
  } finally {
    void decoder.close()
  }
  const seconds = Math.min(decoded.duration, maxSeconds)
  const length = Math.max(1, Math.ceil(seconds * TARGET_RATE))
  const offline = new OfflineAudioContext(1, length, TARGET_RATE)
  const source = offline.createBufferSource()
  source.buffer = decoded
  source.connect(offline.destination)
  source.start()
  const rendered = await offline.startRendering()
  return encodeWav(rendered.getChannelData(0), TARGET_RATE)
}

export function encodeWav(samples: Float32Array, sampleRate: number): Blob {
  const dataBytes = samples.length * 2
  const view = new DataView(new ArrayBuffer(44 + dataBytes))
  const write = (offset: number, text: string) => {
    for (let i = 0; i < text.length; i++) view.setUint8(offset + i, text.charCodeAt(i))
  }
  write(0, 'RIFF')
  view.setUint32(4, 36 + dataBytes, true)
  write(8, 'WAVE')
  write(12, 'fmt ')
  view.setUint32(16, 16, true)
  view.setUint16(20, 1, true)
  view.setUint16(22, 1, true)
  view.setUint32(24, sampleRate, true)
  view.setUint32(28, sampleRate * 2, true)
  view.setUint16(32, 2, true)
  view.setUint16(34, 16, true)
  write(36, 'data')
  view.setUint32(40, dataBytes, true)
  for (let i = 0; i < samples.length; i++) {
    const value = Math.max(-1, Math.min(1, samples[i]))
    view.setInt16(44 + i * 2, value < 0 ? value * 0x8000 : value * 0x7fff, true)
  }
  return new Blob([view], { type: 'audio/wav' })
}
