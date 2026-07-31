import { useMessages } from '@ownclouders/web-pkg'
import { useGeneralActionsResetLogo } from '../../../../../src/composables/actions/general/useGeneralActionsResetLogo'
import { mock } from 'vitest-mock-extended'
import { unref } from 'vue'
import {
  defaultComponentMocks,
  RouteLocation,
  mockAxiosResolve,
  getComposableWrapper
} from '@ownclouders/web-test-helpers'

describe('resetLogo', () => {
  beforeEach(() => {
    vi.useFakeTimers()
  })

  afterEach(() => {
    vi.useRealTimers()
  })

  describe('handler', () => {
    it('should show message on request success', () => {
      getWrapper({
        setup: async ({ actions }, { clientService, router }) => {
          clientService.httpAuthenticated.delete.mockResolvedValue(mockAxiosResolve())
          await unref(actions)[0].handler()
          vi.runAllTimers()
          expect(router.go).toHaveBeenCalledTimes(1)
          const { showMessage } = useMessages()
          expect(showMessage).toHaveBeenCalledTimes(1)
        }
      })
    })

    it('should show message on request error', () => {
      vi.spyOn(console, 'error').mockImplementation(() => undefined)
      getWrapper({
        setup: async ({ actions }, { clientService, router }) => {
          clientService.httpAuthenticated.delete.mockRejectedValue(new Error(''))
          await unref(actions)[0].handler()
          vi.runAllTimers()
          expect(router.go).toHaveBeenCalledTimes(0)
          const { showErrorMessage } = useMessages()
          expect(showErrorMessage).toHaveBeenCalledTimes(1)
        }
      })
    })
  })
})

function getWrapper({
  setup
}: {
  setup: (
    instance: ReturnType<typeof useGeneralActionsResetLogo>,
    {
      clientService,
      router
    }: {
      clientService: ReturnType<typeof defaultComponentMocks>['$clientService']
      router: ReturnType<typeof defaultComponentMocks>['$router']
    }
  ) => void
}) {
  const mocks = defaultComponentMocks({
    currentRoute: mock<RouteLocation>({ name: 'admin-settings-general' })
  })
  return {
    wrapper: getComposableWrapper(
      () => {
        const instance = useGeneralActionsResetLogo()
        setup(instance, {
          clientService: mocks.$clientService,
          router: mocks.$router
        })
      },
      {
        mocks,
        provide: mocks
      }
    )
  }
}
