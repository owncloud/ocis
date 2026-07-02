import { useMessages } from '@ownclouders/web-pkg'
import { useGeneralActionsUploadLogo } from '../../../../../src/composables/actions/general/useGeneralActionsUploadLogo'
import { mock } from 'vitest-mock-extended'
import { VNodeRef } from 'vue'
import {
  defaultComponentMocks,
  RouteLocation,
  mockAxiosResolve,
  mockAxiosReject,
  getComposableWrapper
} from '@ownclouders/web-test-helpers'

describe('uploadImage', () => {
  beforeEach(() => {
    vi.useFakeTimers()
  })

  afterEach(() => {
    vi.useRealTimers()
  })

  describe('method "uploadImage"', () => {
    it('should show message on request success', () => {
      getWrapper({
        setup: async ({ uploadImage }, { clientService, router }) => {
          clientService.httpAuthenticated.post.mockResolvedValue(mockAxiosResolve())
          await uploadImage({
            currentTarget: {
              files: [{ name: 'image.png', type: 'image/png' }]
            }
          } as unknown as InputEvent)
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
        setup: async ({ uploadImage }, { clientService, router }) => {
          clientService.httpAuthenticated.post.mockRejectedValue(() => mockAxiosReject())
          await uploadImage({
            currentTarget: {
              files: [{ name: 'image.png', type: 'image/png' }]
            }
          } as unknown as InputEvent)
          vi.runAllTimers()
          expect(router.go).toHaveBeenCalledTimes(0)
          const { showErrorMessage } = useMessages()
          expect(showErrorMessage).toHaveBeenCalledTimes(1)
        }
      })
    })

    it('should show message on invalid mimeType', () => {
      vi.spyOn(console, 'error').mockImplementation(() => undefined)
      getWrapper({
        setup: async ({ uploadImage }, { clientService }) => {
          await uploadImage({
            currentTarget: {
              files: [{ name: 'text.txt', type: 'text/plain' }]
            }
          } as unknown as InputEvent)
          expect(clientService.httpAuthenticated.post).toHaveBeenCalledTimes(0)
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
    instance: ReturnType<typeof useGeneralActionsUploadLogo>,
    {
      imageInput,
      clientService,
      router
    }: {
      imageInput: VNodeRef
      clientService: ReturnType<typeof defaultComponentMocks>['$clientService']
      router: ReturnType<typeof defaultComponentMocks>['$router']
    }
  ) => void
}) {
  const mocks = defaultComponentMocks({
    currentRoute: mock<RouteLocation>({ name: 'admin-settings-general' })
  })
  return {
    mocks,
    wrapper: getComposableWrapper(
      () => {
        const imageInput = mock<VNodeRef>()
        const instance = useGeneralActionsUploadLogo({ imageInput })
        setup(instance, {
          imageInput,
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
