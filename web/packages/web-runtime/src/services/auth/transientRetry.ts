import { HttpError } from '@ownclouders/web-client'

const maxBackoffMs = 10_000

export const isTransientError = (e: unknown): e is HttpError =>
  e instanceof HttpError && (e.statusCode === 503 || e.statusCode === 429)

const retryAfterMs = (e: HttpError, fallbackMs: number): number => {
  const header = e.response?.headers?.get('retry-after')
  const seconds = header ? Number.parseInt(header, 10) : NaN
  return Number.isFinite(seconds) && seconds >= 0 ? seconds * 1000 : fallbackMs
}

const defaultSleep = (ms: number) => new Promise((resolve) => setTimeout(resolve, ms))

// Retry a transient 503/429 (IdP failure, #12999) honoring Retry-After with
// capped backoff; any other error (e.g. a genuine 401) propagates.
export const retryOnTransientError = async <T>(
  fn: () => Promise<T>,
  { sleep = defaultSleep }: { sleep?: (ms: number) => Promise<unknown> } = {}
): Promise<T> => {
  for (let attempt = 0; ; attempt++) {
    try {
      return await fn()
    } catch (e) {
      if (!isTransientError(e)) {
        throw e
      }
      const backoff = Math.min(maxBackoffMs, 1000 * 2 ** attempt)
      await sleep(retryAfterMs(e, backoff))
    }
  }
}
