import { ClientService, useAuthStore, useConfigStore } from '../../../src/'
import { Language } from 'vue3-gettext'
import { createTestingPinia, writable } from '@ownclouders/web-test-helpers'
import { AxiosError, AxiosResponse } from 'axios'
import { shouldResponseTriggerMaintenance } from '@ownclouders/web-client'
import { mock } from 'vitest-mock-extended'

vi.mock('@ownclouders/web-client', async (importOriginal) => ({
  ...(await importOriginal<any>()),
  graph: vi.fn(),
  ocs: vi.fn(),
  webdav: vi.fn(),
  shouldResponseTriggerMaintenance: vi.fn()
}))

let responseSuccessInterceptorFn: (response: AxiosResponse) => AxiosResponse
let responseErrorInterceptorFn: (error: AxiosError) => Promise<AxiosError>

vi.mock('axios', () => {
  return {
    default: {
      create: vi.fn().mockReturnValue({
        interceptors: {
          response: {
            use: vi.fn().mockImplementation((successFn, errorFn) => {
              responseSuccessInterceptorFn = successFn
              responseErrorInterceptorFn = errorFn
            })
          },
          request: { use: vi.fn() }
        }
      }),
      CancelToken: { source: vi.fn() }
    }
  }
})

describe('ClientService maintenance mode', () => {
  const language = { current: 'en' }
  const serverUrl = 'someUrl'

  let configStore: ReturnType<typeof useConfigStore>
  let authStore: ReturnType<typeof useAuthStore>

  beforeEach(() => {
    createTestingPinia({ initialState: { auth: { accessToken: 'token' } } })

    vi.mocked(shouldResponseTriggerMaintenance).mockReset()

    authStore = useAuthStore()
    configStore = useConfigStore()
    writable(configStore).serverUrl = serverUrl
    configStore.setMaintenanceMode = vi.fn()

    new ClientService({
      configStore,
      language: language as Language,
      authStore
    })
  })

  describe('handling axios responses', () => {
    it('should turn off maintenance mode for successful responses', () => {
      const response = mock<AxiosResponse>({
        status: 200,
        data: { some: 'data' }
      })

      responseSuccessInterceptorFn(response)
      expect(configStore.setMaintenanceMode).toHaveBeenCalledWith(false)
    })

    it('should not turn off maintenance mode for 503 responses', () => {
      const response = mock<AxiosResponse>({
        status: 503,
        data: { error: 'Service Unavailable' }
      })

      responseSuccessInterceptorFn(response)
      expect(configStore.setMaintenanceMode).not.toHaveBeenCalledWith(false)
    })
  })

  describe('handling axios errors', () => {
    it('should turn on maintenance mode when shouldResponseTriggerMaintenance returns true', () => {
      vi.mocked(shouldResponseTriggerMaintenance).mockReturnValue(true)

      const error = mock<AxiosError>({
        response: { status: 503 },
        config: { url: 'some/url' }
      })

      expect(responseErrorInterceptorFn(error)).rejects.toEqual(error)
      expect(shouldResponseTriggerMaintenance).toHaveBeenCalledWith(503, 'some/url')
      expect(configStore.setMaintenanceMode).toHaveBeenCalledWith(true)
    })

    it('should not turn on maintenance mode when shouldResponseTriggerMaintenance returns false', () => {
      vi.mocked(shouldResponseTriggerMaintenance).mockReturnValue(false)

      const error = mock<AxiosError>({
        response: { status: 404 },
        config: { url: 'some/url' }
      })

      expect(responseErrorInterceptorFn(error)).rejects.toEqual(error)
      expect(shouldResponseTriggerMaintenance).toHaveBeenCalledWith(404, 'some/url')
      expect(configStore.setMaintenanceMode).not.toHaveBeenCalled()
    })
  })
})
