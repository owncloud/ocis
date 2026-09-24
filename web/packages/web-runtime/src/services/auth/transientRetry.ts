import { HttpError } from '@ownclouders/web-client'

const maxBackoffMs = 10_000

const isTransientStatus = (e: unknown): e is HttpError =>
  e instanceof HttpError && (e.statusCode === 503 || e.statusCode === 429)

const retryAfterMs = (e: HttpError, fallbackMs: number): number => {
  const header = e.response?.headers?.get?.('retry-after')
  const seconds = header ? Number.parseInt(header, 10) : NaN
  return Number.isFinite(seconds) && seconds >= 0 ? seconds * 1000 : fallbackMs
}

const defaultSleep = (ms: number) => new Promise((resolve) => setTimeout(resolve, ms))

// A 503/429 from the proxy means the IdP userinfo call failed transiently
// (OCISDEV-1411), not that the session is invalid. Wait it out — honoring
// Retry-After with capped exponential backoff — so a transient backend outage
// surfaces as the maintenance banner and recovers on its own, instead of being
// mistaken for an auth failure and logging the user out. Only a non-transient
// error (e.g. a genuine 401) propagates. A truly persistent outage keeps
// waiting; the session ends only when the token itself expires (401).
export const retryOnTransientError = async <T>(
  fn: () => Promise<T>,
  { sleep = defaultSleep }: { sleep?: (ms: number) => Promise<unknown> } = {}
): Promise<T> => {
  for (let attempt = 0; ; attempt++) {
    try {
      return await fn()
    } catch (e) {
      if (!isTransientStatus(e)) {
        throw e
      }
      const backoff = Math.min(maxBackoffMs, 1000 * 2 ** attempt)
      await sleep(retryAfterMs(e, backoff))
    }
  }
}
