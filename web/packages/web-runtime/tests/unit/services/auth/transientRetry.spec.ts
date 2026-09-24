import { HttpError } from '@ownclouders/web-client'
import { retryOnTransientError } from '../../../../src/services/auth/transientRetry'

const httpError = (statusCode: number, retryAfter?: string) =>
  new HttpError(
    `status ${statusCode}`,
    { headers: new Headers(retryAfter ? { 'retry-after': retryAfter } : {}) } as Response,
    statusCode
  )

describe('retryOnTransientError', () => {
  // inject a no-op sleep so the retry loop doesn't wait in tests
  const noSleep = { sleep: vi.fn().mockResolvedValue(undefined) }

  it('returns the result without retrying on success', async () => {
    const fn = vi.fn().mockResolvedValue('ok')
    await expect(retryOnTransientError(fn, noSleep)).resolves.toBe('ok')
    expect(fn).toHaveBeenCalledTimes(1)
  })

  it('waits out a transient 503 and resolves once the backend recovers', async () => {
    const fn = vi
      .fn()
      .mockRejectedValueOnce(httpError(503))
      .mockRejectedValueOnce(httpError(503))
      .mockRejectedValueOnce(httpError(503))
      .mockResolvedValue('ok')
    await expect(retryOnTransientError(fn, noSleep)).resolves.toBe('ok')
    expect(fn).toHaveBeenCalledTimes(4)
  })

  it('retries on a 429', async () => {
    const fn = vi.fn().mockRejectedValueOnce(httpError(429)).mockResolvedValue('ok')
    await expect(retryOnTransientError(fn, noSleep)).resolves.toBe('ok')
    expect(fn).toHaveBeenCalledTimes(2)
  })

  it('does not retry a 401 and rethrows immediately', async () => {
    const err = httpError(401)
    const fn = vi.fn().mockRejectedValue(err)
    await expect(retryOnTransientError(fn, noSleep)).rejects.toBe(err)
    expect(fn).toHaveBeenCalledTimes(1)
  })

  it('rethrows a non-HttpError immediately', async () => {
    const err = new Error('boom')
    const fn = vi.fn().mockRejectedValue(err)
    await expect(retryOnTransientError(fn, noSleep)).rejects.toBe(err)
    expect(fn).toHaveBeenCalledTimes(1)
  })

  it('honors the Retry-After header for the wait duration', async () => {
    const sleep = vi.fn().mockResolvedValue(undefined)
    const fn = vi.fn().mockRejectedValueOnce(httpError(503, '2')).mockResolvedValue('ok')
    await expect(retryOnTransientError(fn, { sleep })).resolves.toBe('ok')
    expect(sleep).toHaveBeenCalledWith(2000)
  })

  it('falls back to capped backoff when Retry-After is absent', async () => {
    const sleep = vi.fn().mockResolvedValue(undefined)
    const fn = vi
      .fn()
      .mockRejectedValueOnce(httpError(503))
      .mockRejectedValueOnce(httpError(503))
      .mockResolvedValue('ok')
    await expect(retryOnTransientError(fn, { sleep })).resolves.toBe('ok')
    // attempt 0 -> 1000ms, attempt 1 -> 2000ms
    expect(sleep).toHaveBeenNthCalledWith(1, 1000)
    expect(sleep).toHaveBeenNthCalledWith(2, 2000)
  })
})
