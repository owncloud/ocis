import { mock } from 'vitest-mock-extended'
import { useRequest } from '../../../../src/composables/authContext/useRequest'
import { HttpClient } from '../../../../src/http'
import { ClientService } from '../../../../src/services'
import { AuthStore } from '../../../../src/composables/piniaStores'

const basicAuth = (password: string) =>
  'Basic ' + Buffer.from(`public:${password}`).toString('base64')

describe('useRequest', () => {
  let fetchMock: ReturnType<typeof vi.fn>

  beforeEach(() => {
    fetchMock = vi.fn().mockResolvedValue(new Response('{}', { status: 200 }))
    vi.stubGlobal('fetch', fetchMock)
  })

  afterEach(() => {
    vi.unstubAllGlobals()
  })

  // Uses a real HttpClient that injects the user's Bearer token as a dynamic header, like
  // the authenticated client of the ClientService does, to verify what is sent over the wire.
  const getSentHeaders = async (authStore: Partial<AuthStore>) => {
    const httpAuthenticated = new HttpClient({
      baseUrl: 'https://ocis.test',
      headers: () => ({ Authorization: 'Bearer user-token' })
    })
    const clientService = mock<ClientService>({ httpAuthenticated })

    const { makeRequest } = useRequest({
      clientService,
      authStore: authStore as AuthStore
    })
    await makeRequest('GET', '/some/url')

    return fetchMock.mock.calls[0][1].headers as Headers
  }

  it('authenticates with the public link password instead of the user token', async () => {
    const headers = await getSentHeaders({
      accessToken: 'user-token',
      publicLinkContextReady: true,
      publicLinkToken: 'link-token',
      publicLinkPassword: 'secret'
    })

    expect(headers.get('Authorization')).toBe(basicAuth('secret'))
    expect(headers.get('public-token')).toBe('link-token')
  })

  it('authenticates with the user token outside of a public link context', async () => {
    const headers = await getSentHeaders({
      accessToken: 'user-token',
      publicLinkContextReady: false
    })

    expect(headers.get('Authorization')).toBe('Bearer user-token')
    expect(headers.has('public-token')).toBe(false)
  })
})
