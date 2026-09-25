import { HttpError } from '@ownclouders/web-client'

const maxBackoffMs = 10_000
// Cap attempts so a sustained outage surfaces a 503 (maintenance banner) instead of hanging login.
const maxAttempts = 60

export const isTransientError = (e: unknown): e is HttpError =>
  e instanceof HttpError && (e.statusCode === 503 || e.statusCode === 429)

// Retry-After is either delta-seconds or an HTTP-date; ms to wait, or undefined if absent/unparseable.
const parseRetryAfter = (header: string | null | undefined): number | undefined => {
  if (!header) {
    return undefined
  }
  const seconds = Number(header)
  if (Number.isFinite(seconds)) {
    return seconds >= 0 ? seconds * 1000 : undefined
  }
  const date = Date.parse(header)
  return Number.isNaN(date) ? undefined : Math.max(0, date - Date.now())
}

const retryAfterMs = (e: HttpError, fallbackMs: number): number =>
  parseRetryAfter(e.response?.headers?.get('retry-after')) ?? fallbackMs

const defaultSleep = (ms: number) => new Promise((resolve) => setTimeout(resolve, ms))

// Retry a transient 503/429 (IdP failure, #12999) honoring Retry-After with
// capped backoff; any other error (e.g. a genuine 401) propagates.
export const retryOnTransientError = async <T>(
  fn: () => Promise<T>,
  { sleep = defaultSleep, attempts = maxAttempts }: { sleep?: (ms: number) => Promise<unknown>; attempts?: number } = {}
): Promise<T> => {
  for (let attempt = 0; ; attempt++) {
    try {
      return await fn()
    } catch (e) {
      if (!isTransientError(e) || attempt >= attempts - 1) {
        throw e
      }
      const backoff = Math.min(maxBackoffMs, 1000 * 2 ** attempt)
      await sleep(retryAfterMs(e, backoff))
    }
  }
}
