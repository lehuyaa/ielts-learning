import { API_BASE_URL } from '@/api/api'

/**
 * Files uploaded through the API (e.g. Whisper-imported recordings) are served
 * by the backend under /media; everything else is a frontend public path.
 */
export function mediaUrl(url: string) {
  if (url.startsWith('/media/')) {
    return new URL(url, API_BASE_URL).toString()
  }
  return url
}
