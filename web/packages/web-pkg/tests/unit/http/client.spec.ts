import { HttpClient } from '../../../src/http'
import { z } from 'zod'
import { mock } from 'vitest-mock-extended'

const schema = z.object({
  someProperty: z.string()
})

type Schema = z.infer<typeof schema>

describe('HttpClient', () => {
  let fetchMock: ReturnType<typeof vi.fn>

  beforeEach(() => {
    fetchMock = vi
      .fn()
      .mockImplementation(() =>
        Promise.resolve(
          new Response('{}', { status: 200, headers: { 'Content-Type': 'application/json' } })
        )
      )
    vi.stubGlobal('fetch', fetchMock)
  })

  test('types', async () => {
    const fn = vi.fn().mockReturnValue({ data: {} })

    const client = mock<HttpClient>({
      delete: fn,
      get: fn,
      head: fn,
      options: fn,
      patch: fn,
      post: fn,
      put: fn,
      request: fn
    })

    // delete
    {
      const { data: untypedResponse } = await client.delete('/foo')
      untypedResponse.property

      const { data: typedResponse } = await client.delete<Schema>('/foo')
      typedResponse.someProperty

      const { data: schemaResponse } = await client.delete('/foo', null, { schema })
      schemaResponse.someProperty
    }

    // get
    {
      const { data: untypedResponse } = await client.get('/foo')
      // @ts-expect-error untypedResponse is unknown, so we cannot access a property without casting
      untypedResponse.property

      const { data: typedResponse } = await client.get<Schema>('/foo')
      typedResponse.someProperty

      const { data: schemaResponse } = await client.get('/foo', { schema })
      schemaResponse.someProperty
    }

    // head
    {
      const { data: untypedResponse } = await client.head('/foo')
      untypedResponse.property

      const { data: typedResponse } = await client.head<Schema>('/foo')
      typedResponse.someProperty

      const { data: schemaResponse } = await client.head('/foo', { schema })
      schemaResponse.someProperty
    }

    // options
    {
      const { data: untypedResponse } = await client.options('/foo')
      untypedResponse.property

      const { data: typedResponse } = await client.options<Schema>('/foo')
      typedResponse.someProperty

      const { data: schemaResponse } = await client.options('/foo', { schema })
      schemaResponse.someProperty
    }

    // patch
    {
      const { data: untypedResponse } = await client.patch('/foo')
      untypedResponse.property

      const { data: typedResponse } = await client.patch<Schema>('/foo')
      typedResponse.someProperty

      const { data: schemaResponse } = await client.patch('/foo', null, { schema })
      schemaResponse.someProperty
    }

    // post
    {
      const { data: untypedResponse } = await client.post('/foo')
      untypedResponse.property

      const { data: typedResponse } = await client.post<Schema>('/foo')
      typedResponse.someProperty

      const { data: schemaResponse } = await client.post('/foo', null, { schema })
      schemaResponse.someProperty
    }

    // put
    {
      const { data: untypedResponse } = await client.put('/foo')
      untypedResponse.property

      const { data: typedResponse } = await client.put<Schema>('/foo')
      typedResponse.someProperty

      const { data: schemaResponse } = await client.put('/foo', null, { schema })
      schemaResponse.someProperty
    }

    // request
    {
      const { data: untypedResponse } = await client.request({ url: '/foo' })
      untypedResponse.property

      const { data: typedResponse } = await client.request<Schema>({ url: '/foo' })
      typedResponse.someProperty

      const { data: schemaResponse } = await client.request({ url: '/foo', schema })
      schemaResponse.someProperty
    }
    expect(true).toBe(true)
  })
  test.each([
    ['delete', 'DELETE'],
    ['get', 'GET'],
    ['head', 'HEAD'],
    ['options', 'OPTIONS'],
    ['patch', 'PATCH'],
    ['post', 'POST'],
    ['put', 'PUT']
  ] as const)('%s issues a %s request', async (method, verb) => {
    const client = new HttpClient()
    await client[method]('https://host/url')

    expect(fetchMock).toHaveBeenCalledTimes(1)
    expect(fetchMock.mock.calls[0][1].method).toBe(verb)
  })

  test('request takes url and method from the config', async () => {
    await new HttpClient().request({ url: 'https://host/url', method: 'GET' })

    expect(fetchMock).toHaveBeenCalledTimes(1)
    expect(fetchMock.mock.calls[0][0]).toBe('https://host/url')
    expect(fetchMock.mock.calls[0][1].method).toBe('GET')
  })

  test('request sends the config data as the request body', async () => {
    await new HttpClient().request({
      url: 'https://host/url',
      method: 'REPORT',
      data: '<?xml version="1.0"?><search />'
    })

    expect(fetchMock.mock.calls[0][1].body).toBe('<?xml version="1.0"?><search />')
  })

  test.each(['patch', 'post', 'put'] as const)(
    '%s sends its positional data as the request body',
    async (method) => {
      await new HttpClient()[method]('https://host/url', 'payload')

      expect(fetchMock.mock.calls[0][1].body).toBe('payload')
    }
  )

  test('applies a zod schema to the response data', async () => {
    fetchMock.mockImplementation(() =>
      Promise.resolve(new Response(JSON.stringify({ someProperty: 'value' }), { status: 200 }))
    )

    const { data } = await new HttpClient().get('https://host/url', { schema })

    expect(data.someProperty).toBe('value')
  })

  test('resolves a non-json body as text, the way a WebDAV REPORT answers', async () => {
    const multistatus = '<?xml version="1.0"?><d:multistatus xmlns:d="DAV:" />'
    fetchMock.mockImplementation(() => Promise.resolve(new Response(multistatus, { status: 207 })))

    const { data } = await new HttpClient().request({ url: 'https://host/dav', method: 'REPORT' })

    expect(data).toBe(multistatus)
  })

  test('exposes response headers both by bracket access and through get', async () => {
    fetchMock.mockImplementation(() =>
      Promise.resolve(new Response('{}', { status: 200, headers: { 'Lock-Token': '<token>' } }))
    )

    const { headers } = await new HttpClient().request({ url: 'https://host/dav', method: 'LOCK' })

    expect(headers['lock-token']).toBe('<token>')
    expect(headers.get('lock-token')).toBe('<token>')
  })

  test('baseUrl is applied to relative urls', async () => {
    await new HttpClient({ baseUrl: 'https://host/' }).get('some/path')

    expect(fetchMock.mock.calls[0][0]).toBe('https://host/some/path')
  })

  describe('cancel', () => {
    /** leaves the request pending until its signal aborts, as a real fetch would */
    const neverResolvingFetch = () =>
      fetchMock.mockImplementation(
        (_url: string, init: RequestInit) =>
          new Promise((_resolve, reject) => {
            init.signal.addEventListener('abort', () => reject(init.signal.reason))
          })
      )

    test('aborts a request that is in flight', async () => {
      neverResolvingFetch()
      const client = new HttpClient()

      const request = client.get('https://host/url')
      client.cancel('gone')

      await expect(request).rejects.toMatchObject({ name: 'AbortError', message: 'gone' })
    })

    test('leaves a per-request signal able to abort on its own', async () => {
      neverResolvingFetch()
      const controller = new AbortController()

      const request = new HttpClient().get('https://host/url', { signal: controller.signal })
      controller.abort(new DOMException('mine', 'AbortError'))

      await expect(request).rejects.toMatchObject({ name: 'AbortError', message: 'mine' })
    })

    test('passes the aborted signal on to requests issued after the cancel', async () => {
      const client = new HttpClient()
      client.cancel()

      await client.get('https://host/url').catch(() => undefined)

      expect(fetchMock.mock.calls[0][1].signal.aborted).toBe(true)
    })

    test('removes its abort listeners once a request has settled', async () => {
      const client = new HttpClient()
      const controller = new AbortController()
      const removeEventListener = vi.spyOn(controller.signal, 'removeEventListener')

      await client.get('https://host/a', { signal: controller.signal })
      await client.get('https://host/b', { signal: controller.signal })

      expect(removeEventListener).toHaveBeenCalledTimes(2)
    })
  })
})
