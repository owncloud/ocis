import { ClientService, useAuthStore, useConfigStore } from '../../../src/'
import { Language } from 'vue3-gettext'
import { createTestingPinia, writable } from '@ownclouders/web-test-helpers'
import type { OnResponseArgs } from '@ownclouders/web-client'

/**
 * `shouldResponseTriggerMaintenance` is deliberately left unmocked: the decision it encodes
 * — which statuses and which endpoints count — is the thing under test here, and stubbing it
 * would only assert that one function calls another.
 */
vi.mock('@ownclouders/web-client', async (importOriginal) => ({
  ...(await importOriginal<any>()),
  graph: vi.fn(),
  ocs: vi.fn(),
  webdav: vi.fn()
}))

describe('ClientService maintenance mode', () => {
  const language = { current: 'en' }
  const serverUrl = 'someUrl'

  let configStore: ReturnType<typeof useConfigStore>
  let authStore: ReturnType<typeof useAuthStore>
  let service: ClientService
  let onResponse: (args: OnResponseArgs) => void

  beforeEach(() => {
    createTestingPinia({ initialState: { auth: { accessToken: 'token' } } })
    vi.stubGlobal('fetch', vi.fn())

    authStore = useAuthStore()
    configStore = useConfigStore()
    writable(configStore).serverUrl = serverUrl
    configStore.setMaintenanceMode = vi.fn()

    service = new ClientService({
      configStore,
      language: language as Language,
      authStore
    })

    onResponse = service.handleResponse.bind(service)
  })

  it('clears maintenance mode and records the time for a successful response', () => {
    onResponse({
      response: new Response('{}', { status: 200 }),
      status: 200,
      requestUrl: 'some/url'
    })

    expect(configStore.setMaintenanceMode).toHaveBeenCalledWith(false)
    expect(service.lastSuccessfulRequestTime).not.toBeNull()
  })

  it('sets maintenance mode for a 503', () => {
    onResponse({
      response: new Response('{}', { status: 503 }),
      status: 503,
      requestUrl: 'some/url'
    })

    expect(configStore.setMaintenanceMode).toHaveBeenCalledWith(true)
  })

  it('leaves maintenance state untouched for a 404', () => {
    onResponse({
      response: new Response('{}', { status: 404 }),
      status: 404,
      requestUrl: 'some/url'
    })

    expect(configStore.setMaintenanceMode).not.toHaveBeenCalled()
  })

  /**
   * A transport failure has no response to read a status off, so it is reported as 500. As on
   * master — where the axios error interceptor fell back to `error.response?.status || 500` —
   * that is not a maintenance signal, so the flag is left alone rather than cleared.
   */
  it('leaves maintenance state untouched for a transport failure', () => {
    onResponse({ response: null, status: 500, requestUrl: 'some/url' })

    expect(configStore.setMaintenanceMode).not.toHaveBeenCalled()
  })

  /**
   * The allow-list is matched against the relative request url, which is why `onResponse`
   * reports the caller's url rather than `response.url`. The notifications SSE endpoint
   * answers 503 by design and must not raise the banner.
   */
  it('exempts an allow-listed endpoint from a 503', () => {
    onResponse({
      response: new Response('{}', { status: 503 }),
      status: 503,
      requestUrl: 'ocs/v2.php/apps/notifications/api/v1/notifications/sse'
    })

    expect(configStore.setMaintenanceMode).not.toHaveBeenCalled()
  })
})
