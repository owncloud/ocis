/**
 * @vitest-environment node
 *
 * Deliberately not happy-dom, whose `Headers` constructor applies a record init with `set`
 * semantics where the spec — and therefore browsers and undici — applies it with `append`.
 * These assertions are about a header merge, so a forgiving implementation would hide a
 * duplicate-header bug rather than fail on it.
 */
import { graph } from '../../../src/graph'
import { FetchClient } from '../../../src/http'

describe('graph request headers', () => {
  let fetchMock: ReturnType<typeof vi.fn>

  beforeEach(() => {
    fetchMock = vi.fn().mockResolvedValue(
      new Response(JSON.stringify({ value: [] }), {
        status: 200,
        headers: { 'Content-Type': 'application/json' }
      })
    )
    vi.stubGlobal('fetch', fetchMock)
  })

  const sentHeaders = () => new Headers(fetchMock.mock.calls[0][1].headers)

  it('sends the client-wide headers', async () => {
    const client = new FetchClient({
      staticHeaders: { 'X-Requested-With': 'XMLHttpRequest' },
      headers: () => ({ Authorization: 'Bearer token' })
    })

    await graph('https://host', client).tags.listTags()

    expect(sentHeaders().get('authorization')).toBe('Bearer token')
    expect(sentHeaders().get('x-requested-with')).toBe('XMLHttpRequest')
  })

  /**
   * A per-request override has to replace the client-wide header it names, not be appended
   * next to it. Both spellings surviving would put `Bearer stale, Bearer fresh` on the wire,
   * which no server accepts as a token.
   */
  it('lets a per-request header override a differently-cased client-wide one', async () => {
    const client = new FetchClient({ headers: () => ({ Authorization: 'Bearer stale' }) })

    await graph('https://host', client).tags.listTags({
      headers: { authorization: 'Bearer fresh' }
    })

    expect(sentHeaders().get('authorization')).toBe('Bearer fresh')
  })

  /**
   * The generated client writes `Content-Type` itself for a request with a body. An override
   * has to replace that entry rather than land beside it as a second, lowercase one.
   */
  it('lets a per-request header override a differently-cased generated one', async () => {
    await graph('https://host', new FetchClient()).tags.assignTags(
      { resourceId: 'storage$space!node', tags: ['a'] },
      { headers: { 'content-type': 'application/json; charset=utf-8' } }
    )

    expect(sentHeaders().get('content-type')).toBe('application/json; charset=utf-8')
    // the body still goes out as JSON: the generated runtime decides that before the override
    expect(fetchMock.mock.calls[0][1].body).toBe('{"resourceId":"storage$space!node","tags":["a"]}')
  })

  it('keeps a client-wide header a request does not name', async () => {
    const client = new FetchClient({ headers: () => ({ Authorization: 'Bearer token' }) })

    await graph('https://host', client).tags.assignTags(
      { resourceId: 'storage$space!node', tags: ['a'] },
      { headers: { Purge: 'T' } }
    )

    const headers = sentHeaders()
    expect(headers.get('purge')).toBe('T')
    expect(headers.get('authorization')).toBe('Bearer token')
    // set by the generated client for the request body, untouched by the override
    expect(headers.get('content-type')).toBe('application/json')
  })
})
