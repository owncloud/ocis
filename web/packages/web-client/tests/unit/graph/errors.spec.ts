import { graph } from '../../../src/graph'
import { FetchClient } from '../../../src/http'
import { HttpError } from '../../../src/errors'

describe('graph error propagation', () => {
  let fetchMock: ReturnType<typeof vi.fn>

  beforeEach(() => {
    fetchMock = vi.fn()
    vi.stubGlobal('fetch', fetchMock)
  })

  /**
   * The generated runtime wraps anything its `fetchApi` throws in a `FetchError`. Callers
   * such as the "banned password" hint read `statusCode` and `data` off the rejection, so
   * the `HttpError` from the core has to survive the round trip.
   */
  it('rejects with the HttpError raised for a non-2xx response', async () => {
    const body = { error: { message: 'password is commonly used' } }
    fetchMock.mockResolvedValue(
      new Response(JSON.stringify(body), {
        status: 400,
        headers: { 'Content-Type': 'application/json' }
      })
    )

    const permissions = graph('https://host', new FetchClient()).permissions
    const request = permissions.setPermissionPassword('space-id', 'item-id', 'link-id', {
      password: 'ownCloud-1'
    })

    await expect(request).rejects.toBeInstanceOf(HttpError)
    await expect(request).rejects.toMatchObject({ statusCode: 400, data: body })
  })

  it('keeps an abort an AbortError', async () => {
    const controller = new AbortController()
    fetchMock.mockImplementation(
      (_url: string, init: RequestInit) =>
        new Promise((_resolve, reject) => {
          // the generated client awaits before it calls fetch, so the abort may already be in
          if (init.signal.aborted) {
            reject(init.signal.reason)
            return
          }
          init.signal.addEventListener('abort', () => reject(init.signal.reason))
        })
    )

    const request = graph('https://host', new FetchClient()).drives.listMyDrives(
      {},
      {},
      { signal: controller.signal }
    )
    controller.abort(new DOMException('gone', 'AbortError'))

    await expect(request).rejects.toMatchObject({ name: 'AbortError' })
  })
})
