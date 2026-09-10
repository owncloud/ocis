/**
 * @vitest-environment node
 *
 * Deliberately not happy-dom. Its `Headers` constructor applies a record init with `set`
 * semantics, whereas the spec — and therefore both browsers and undici — applies it with
 * `append`. Every header assertion below is about that merge, so running them against a
 * forgiving implementation would let a real duplicate-header bug pass.
 */
import { FetchClient } from '../../../src/http'
import { HttpError } from '../../../src/errors'

const jsonResponse = (body: unknown, init: ResponseInit = {}) =>
  new Response(JSON.stringify(body), {
    status: 200,
    headers: { 'Content-Type': 'application/json' },
    ...init
  })

describe('FetchClient', () => {
  let fetchMock: ReturnType<typeof vi.fn>

  beforeEach(() => {
    fetchMock = vi.fn()
    vi.stubGlobal('fetch', fetchMock)
  })

  const lastCall = () => fetchMock.mock.calls[0]

  describe('request envelope', () => {
    it('returns data, status, statusText and headers', async () => {
      fetchMock.mockResolvedValue(jsonResponse({ some: 'value' }, { statusText: 'OK' }))

      const result = await new FetchClient().request<{ some: string }>('https://host/foo')

      expect(result.data).toEqual({ some: 'value' })
      expect(result.status).toBe(200)
      expect(result.statusText).toBe('OK')
      expect(result.headers.get('Content-Type')).toBe('application/json')
    })

    it('trap 2: headers also answer to bracket access, as axios headers did', async () => {
      fetchMock.mockResolvedValue(jsonResponse({}, { headers: { 'Lock-Token': '<token>' } }))

      const result = await new FetchClient().request('https://host/foo')

      expect(result.headers['lock-token']).toBe('<token>')
      expect(result.headers['Lock-Token']).toBe('<token>')
    })
  })

  describe('trap 1: throws on non-2xx', () => {
    it.each([400, 404, 500, 503])('throws HttpError for %i', async (status) => {
      fetchMock.mockResolvedValue(new Response('{}', { status }))

      const client = new FetchClient()
      await expect(client.request('https://host/foo')).rejects.toBeInstanceOf(HttpError)
    })

    it('sets statusCode, not status', async () => {
      fetchMock.mockResolvedValue(new Response('{}', { status: 423 }))

      const error: HttpError = await new FetchClient().request('https://host/foo').catch((e) => e)

      expect(error.statusCode).toBe(423)
    })

    it('carries the parsed error body on data', async () => {
      fetchMock.mockResolvedValue(
        new Response(JSON.stringify({ error: { message: 'nope' } }), { status: 400 })
      )

      const error: HttpError = await new FetchClient().request('https://host/foo').catch((e) => e)

      expect(error.data).toEqual({ error: { message: 'nope' } })
    })

    it('does not fail the throw path on a non-JSON error body', async () => {
      fetchMock.mockResolvedValue(new Response('<html>gateway</html>', { status: 502 }))

      const error: HttpError = await new FetchClient().request('https://host/foo').catch((e) => e)

      expect(error.statusCode).toBe(502)
      expect(error.data).toBe('<html>gateway</html>')
    })

    it('keeps the axios error message, which callers may match on', async () => {
      fetchMock.mockResolvedValue(new Response('{}', { status: 404, statusText: 'Not Found' }))

      const error: HttpError = await new FetchClient().request('https://host/foo').catch((e) => e)

      expect(error.message).toBe('Request failed with status code 404')
    })

    it('also exposes the body and headers under their axios names on response', async () => {
      fetchMock.mockResolvedValue(
        new Response(JSON.stringify({ error: 'nope' }), {
          status: 429,
          headers: { 'Retry-After': '30' }
        })
      )

      const error: HttpError = await new FetchClient().request('https://host/foo').catch((e) => e)

      expect((error.response as Response & { data: unknown }).data).toEqual({ error: 'nope' })
      expect(error.response.headers['retry-after']).toBe('30')
      expect(error.response.headers.get('retry-after')).toBe('30')
    })

    it('reads the error body with the responseType the caller asked for', async () => {
      fetchMock.mockResolvedValue(new Response('{"error":"nope"}', { status: 400 }))

      const error: HttpError = await new FetchClient()
        .request('https://host/foo', { responseType: 'blob' })
        .catch((e) => e)

      expect(error.data).toBeInstanceOf(Blob)
    })

    it('returns the envelope instead of throwing when throwOnError is false', async () => {
      fetchMock.mockResolvedValue(new Response('{}', { status: 404 }))

      const result = await new FetchClient().request('https://host/foo', {
        throwOnError: false
      })

      expect(result.status).toBe(404)
    })
  })

  describe('onResponse hook', () => {
    it('fires for a successful response', async () => {
      fetchMock.mockResolvedValue(jsonResponse({}))
      const onResponse = vi.fn()

      await new FetchClient({ onResponse }).request('https://host/foo')

      expect(onResponse).toHaveBeenCalledWith(
        expect.objectContaining({ status: 200, requestUrl: 'https://host/foo' })
      )
    })

    it('fires for a non-2xx response before the throw', async () => {
      fetchMock.mockResolvedValue(new Response('{}', { status: 503 }))
      const onResponse = vi.fn()

      await new FetchClient({ onResponse }).request('https://host/foo').catch(() => undefined)

      expect(onResponse).toHaveBeenCalledWith(expect.objectContaining({ status: 503 }))
    })

    it('trap 4: reports the caller URL, not the resolved response.url', async () => {
      const relative = 'ocs/v2.php/apps/notifications/api/v1/notifications/sse'
      fetchMock.mockResolvedValue(new Response('{}', { status: 503 }))
      const onResponse = vi.fn()

      await new FetchClient({ baseUrl: 'https://host/', onResponse })
        .request(relative)
        .catch(() => undefined)

      expect(onResponse).toHaveBeenCalledWith(expect.objectContaining({ requestUrl: relative }))
    })

    it('trap 5: reports a transport failure as status 500 with a null response, then throws', async () => {
      fetchMock.mockRejectedValue(new TypeError('Failed to fetch'))
      const onResponse = vi.fn()

      const error = await new FetchClient({ onResponse })
        .request('https://host/foo')
        .catch((e) => e)

      expect(onResponse).toHaveBeenCalledWith({
        response: null,
        status: 500,
        requestUrl: 'https://host/foo'
      })
      expect(error).toBeInstanceOf(HttpError)
      expect(error.statusCode).toBe(500)
    })

    it('does not invoke onResponse when the request is aborted', async () => {
      const abortError = new DOMException('aborted', 'AbortError')
      fetchMock.mockRejectedValue(abortError)
      const onResponse = vi.fn()

      await expect(new FetchClient({ onResponse }).request('https://host/foo')).rejects.toBe(
        abortError
      )
      expect(onResponse).not.toHaveBeenCalled()
    })

    /**
     * `AbortController.abort(reason)` rejects the fetch with that reason verbatim, so an
     * aborted request is not always identifiable by the error name. Treating one as a
     * transport failure would report a 500 and could trip maintenance detection.
     */
    it('propagates an abort whose reason is not named AbortError', async () => {
      const controller = new AbortController()
      const reason = new Error('superseded by a newer request')
      fetchMock.mockImplementation(() => {
        controller.abort(reason)
        return Promise.reject(reason)
      })
      const onResponse = vi.fn()

      await expect(
        new FetchClient({ onResponse }).request('https://host/foo', { signal: controller.signal })
      ).rejects.toBe(reason)
      expect(onResponse).not.toHaveBeenCalled()
    })
  })

  describe('url and params', () => {
    it('joins a relative url onto baseUrl without doubling slashes', async () => {
      fetchMock.mockResolvedValue(jsonResponse({}))

      await new FetchClient({ baseUrl: 'https://host/' }).request('/foo')

      expect(lastCall()[0]).toBe('https://host/foo')
    })

    it('leaves an absolute url untouched', async () => {
      fetchMock.mockResolvedValue(jsonResponse({}))

      await new FetchClient({ baseUrl: 'https://host/' }).request('https://other/foo')

      expect(lastCall()[0]).toBe('https://other/foo')
    })

    it('serializes params', async () => {
      fetchMock.mockResolvedValue(jsonResponse({}))

      await new FetchClient().request('https://host/foo', { params: { a: '1', b: 2 } })

      expect(lastCall()[0]).toBe('https://host/foo?a=1&b=2')
    })

    it('appends params onto a url that already has a query', async () => {
      fetchMock.mockResolvedValue(jsonResponse({}))

      await new FetchClient().request('https://host/foo?x=0', { params: { a: '1' } })

      expect(lastCall()[0]).toBe('https://host/foo?x=0&a=1')
    })

    it('leaves the url unchanged for empty or absent params', async () => {
      // a fresh Response per call: a body is single-use
      fetchMock.mockImplementation(() => Promise.resolve(jsonResponse({})))
      const client = new FetchClient()

      await client.request('https://host/foo', { params: {} })
      expect(lastCall()[0]).toBe('https://host/foo')

      fetchMock.mockClear()
      await client.request('https://host/foo')
      expect(lastCall()[0]).toBe('https://host/foo')
    })
  })

  describe('body encoding', () => {
    it('JSON-stringifies a plain object and sets Content-Type', async () => {
      fetchMock.mockResolvedValue(jsonResponse({}))

      await new FetchClient().request('https://host/foo', {
        method: 'POST',
        body: { a: 1 }
      })

      const init = lastCall()[1]
      expect(init.body).toBe('{"a":1}')
      expect((init.headers as Headers).get('Content-Type')).toBe('application/json')
    })

    it('passes a BodyInit through untouched and sets no Content-Type', async () => {
      fetchMock.mockResolvedValue(jsonResponse({}))
      const form = new FormData()

      await new FetchClient().request('https://host/foo', { method: 'POST', body: form })

      const init = lastCall()[1]
      expect(init.body).toBe(form)
      expect((init.headers as Headers).get('Content-Type')).toBeNull()
    })

    /**
     * A `multipart/form-data` header written by hand has no boundary, and fetch would send it
     * verbatim, leaving the server unable to parse the body. Callers such as the logo upload
     * relied on axios replacing the header, so it has to be dropped here.
     */
    it('drops a boundary-less multipart Content-Type for FormData', async () => {
      fetchMock.mockResolvedValue(jsonResponse({}))

      await new FetchClient().request('https://host/foo', {
        method: 'POST',
        body: new FormData(),
        headers: { 'Content-Type': 'multipart/form-data' }
      })

      expect((lastCall()[1].headers as Headers).get('Content-Type')).toBeNull()
    })

    it('keeps a multipart Content-Type that already carries a boundary', async () => {
      fetchMock.mockResolvedValue(jsonResponse({}))
      const contentType = 'multipart/form-data; boundary=--abc'

      await new FetchClient().request('https://host/foo', {
        method: 'POST',
        body: new FormData(),
        headers: { 'Content-Type': contentType }
      })

      expect((lastCall()[1].headers as Headers).get('Content-Type')).toBe(contentType)
    })
  })

  describe('responseType', () => {
    it('reads text', async () => {
      fetchMock.mockResolvedValue(new Response('hello', { status: 200 }))

      const result = await new FetchClient().request('https://host/foo', {
        responseType: 'text'
      })

      expect(result.data).toBe('hello')
    })

    it('reads blob', async () => {
      fetchMock.mockResolvedValue(new Response('hello', { status: 200 }))

      const result = await new FetchClient().request<Blob>('https://host/foo', {
        responseType: 'blob'
      })

      expect(result.data).toBeInstanceOf(Blob)
    })

    it('reads arraybuffer', async () => {
      fetchMock.mockResolvedValue(new Response('hello', { status: 200 }))

      const result = await new FetchClient().request<ArrayBuffer>('https://host/foo', {
        responseType: 'arraybuffer'
      })

      expect(result.data).toBeInstanceOf(ArrayBuffer)
    })

    it('returns undefined for responseType none and for 204', async () => {
      const client = new FetchClient()

      fetchMock.mockResolvedValue(new Response('ignored', { status: 200 }))
      expect(
        (await client.request('https://host/a', { responseType: 'none' })).data
      ).toBeUndefined()

      fetchMock.mockResolvedValue(new Response(null, { status: 204 }))
      expect((await client.request('https://host/b')).data).toBeUndefined()
    })

    it('returns undefined for an empty json body rather than throwing', async () => {
      fetchMock.mockResolvedValue(new Response('', { status: 200 }))

      const result = await new FetchClient().request('https://host/foo')

      expect(result.data).toBeUndefined()
    })

    it('falls back to the raw body when a default-responseType body is not json', async () => {
      const multistatus = '<?xml version="1.0"?><d:multistatus xmlns:d="DAV:" />'
      fetchMock.mockResolvedValue(new Response(multistatus, { status: 207 }))

      const result = await new FetchClient().request('https://host/foo', { method: 'REPORT' })

      expect(result.data).toBe(multistatus)
    })
  })

  describe('headers', () => {
    it('applies staticHeaders, then headers(), then per-request headers', async () => {
      fetchMock.mockResolvedValue(jsonResponse({}))

      await new FetchClient({
        staticHeaders: { 'X-Static': 'a', 'X-Shared': 'static' },
        headers: () => ({ 'X-Dynamic': 'b', 'X-Shared': 'dynamic' })
      }).request('https://host/foo', { headers: { 'X-Shared': 'request' } })

      const headers = lastCall()[1].headers as Headers
      expect(headers.get('X-Static')).toBe('a')
      expect(headers.get('X-Dynamic')).toBe('b')
      expect(headers.get('X-Shared')).toBe('request')
    })

    /**
     * Header names are case-insensitive, so a later layer has to replace an earlier one no
     * matter how either spelled it. The graph bridge lowercases every per-request name, which
     * makes the mismatch against the canonically-cased client headers the normal case rather
     * than an edge one — and a spec-compliant `Headers` would join the two values with a comma,
     * sending `Bearer stale, Bearer fresh` instead of overriding the token.
     */
    it.each([
      ['authorization', 'Authorization'],
      ['Authorization', 'authorization'],
      ['X-REQUEST-ID', 'x-request-id']
    ])('overrides a %s client header with a per-request %s', async (clientName, requestName) => {
      fetchMock.mockResolvedValue(jsonResponse({}))

      await new FetchClient({ headers: () => ({ [clientName]: 'stale' }) }).request(
        'https://host/foo',
        { headers: { [requestName]: 'fresh' } }
      )

      const headers = lastCall()[1].headers as Headers
      expect(headers.get(clientName)).toBe('fresh')
      expect([...headers]).toHaveLength(1)
    })

    it('overrides a staticHeader whose casing differs from the per-request one', async () => {
      fetchMock.mockResolvedValue(jsonResponse({}))

      await new FetchClient({ staticHeaders: { 'X-Requested-With': 'XMLHttpRequest' } }).request(
        'https://host/foo',
        { headers: { 'x-requested-with': 'fetch' } }
      )

      const headers = lastCall()[1].headers as Headers
      expect(headers.get('x-requested-with')).toBe('fetch')
    })

    it('accepts a Headers instance as per-request headers', async () => {
      fetchMock.mockResolvedValue(jsonResponse({}))

      await new FetchClient({ headers: () => ({ Authorization: 'stale' }) }).request(
        'https://host/foo',
        { headers: new Headers({ authorization: 'fresh' }) }
      )

      expect((lastCall()[1].headers as Headers).get('authorization')).toBe('fresh')
    })

    it('evaluates headers() on every request', async () => {
      // a fresh Response per call: a body is single-use
      fetchMock.mockImplementation(() => Promise.resolve(jsonResponse({})))
      const headers = vi.fn().mockReturnValue({})
      const client = new FetchClient({ headers })

      await client.request('https://host/a')
      await client.request('https://host/b')

      expect(headers).toHaveBeenCalledTimes(2)
    })
  })

  it('forwards the abort signal', async () => {
    fetchMock.mockResolvedValue(jsonResponse({}))
    const controller = new AbortController()

    await new FetchClient().request('https://host/foo', { signal: controller.signal })

    expect(lastCall()[1].signal).toBe(controller.signal)
  })
})
