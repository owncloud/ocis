import { ClientService, useAuthStore, useConfigStore } from '../../../src/'
import { Language } from 'vue3-gettext'
import { createTestingPinia, writable } from '@ownclouders/web-test-helpers'
import { shouldResponseTriggerMaintenance } from '@ownclouders/web-client'
import type { OnResponseArgs } from '@ownclouders/web-client'

vi.mock('@ownclouders/web-client', async (importOriginal) => ({
  ...(await importOriginal<any>()),
  graph: vi.fn(),
  ocs: vi.fn(),
  webdav: vi.fn(),
  shouldResponseTriggerMaintenance: vi.fn()
}))

describe('ClientService maintenance mode', () => {
  const language = { current: 'en' }
  const serverUrl = 'someUrl'

  let configStore: ReturnType<typeof useConfigStore>
  let authStore: ReturnType<typeof useAuthStore>
  let onResponse: (args: OnResponseArgs) => void

  beforeEach(() => {
    createTestingPinia({ initialState: { auth: { accessToken: 'token' } } })
    vi.mocked(shouldResponseTriggerMaintenance).mockReset()
    vi.stubGlobal('fetch', vi.fn())

    authStore = useAuthStore()
    configStore = useConfigStore()
    writable(configStore).serverUrl = serverUrl
    configStore.setMaintenanceMode = vi.fn()

    const service = new ClientService({
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
  })

  it('sets maintenance mode when shouldResponseTriggerMaintenance returns true', () => {
    vi.mocked(shouldResponseTriggerMaintenance).mockReturnValue(true)

    onResponse({
      response: new Response('{}', { status: 503 }),
      status: 503,
      requestUrl: 'some/url'
    })

    expect(shouldResponseTriggerMaintenance).toHaveBeenCalledWith(503, 'some/url')
    expect(configStore.setMaintenanceMode).toHaveBeenCalledWith(true)
  })

  it('leaves maintenance state untouched for a 404', () => {
    vi.mocked(shouldResponseTriggerMaintenance).mockReturnValue(false)

    onResponse({
      response: new Response('{}', { status: 404 }),
      status: 404,
      requestUrl: 'some/url'
    })

    expect(shouldResponseTriggerMaintenance).toHaveBeenCalledWith(404, 'some/url')
    expect(configStore.setMaintenanceMode).not.toHaveBeenCalled()
  })

  it('trap 5: treats a transport failure as 503-eligible via status 500', () => {
    vi.mocked(shouldResponseTriggerMaintenance).mockReturnValue(true)

    onResponse({ response: null, status: 500, requestUrl: 'some/url' })

    expect(shouldResponseTriggerMaintenance).toHaveBeenCalledWith(500, 'some/url')
    expect(configStore.setMaintenanceMode).toHaveBeenCalledWith(true)
  })

  it('trap 4: forwards the relative request url to the maintenance check', () => {
    vi.mocked(shouldResponseTriggerMaintenance).mockReturnValue(false)
    const sseUrl = 'ocs/v2.php/apps/notifications/api/v1/notifications/sse'

    onResponse({
      response: new Response('{}', { status: 503 }),
      status: 503,
      requestUrl: sseUrl
    })

    expect(shouldResponseTriggerMaintenance).toHaveBeenCalledWith(503, sseUrl)
  })
})
