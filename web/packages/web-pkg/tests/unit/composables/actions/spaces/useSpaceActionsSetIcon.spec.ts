import { useSpaceActionsSetIcon } from '../../../../../src/composables/actions/spaces/useSpaceActionsSetIcon'
import { useMessages, useModals } from '../../../../../src/composables/piniaStores'
import {
  defaultComponentMocks,
  RouteLocation,
  getComposableWrapper
} from '@ownclouders/web-test-helpers'
import { unref } from 'vue'
import { SpaceResource } from '@ownclouders/web-client'
import { mock } from 'vitest-mock-extended'

describe('setIcon', () => {
  beforeEach(() => {
    const createElementMock = vi.spyOn(document, 'createElement')
    createElementMock.mockImplementation(() => {
      return {
        insertBefore: vi.fn(),
        toBlob: () => new Blob(),
        getContext: () => ({
          fillText: vi.fn()
        })
      } as unknown as HTMLElement
    })
  })
  describe('isVisible property', () => {
    it('should be false when no resource given', () => {
      getWrapper({
        setup: ({ actions }) => {
          expect(unref(actions)[0].isVisible({ resources: [] })).toBe(false)
        }
      })
    })
    it('should be false when multiple resources are given', () => {
      getWrapper({
        setup: ({ actions }) => {
          expect(
            unref(actions)[0].isVisible({
              resources: [mock<SpaceResource>(), mock<SpaceResource>()]
            })
          ).toBe(false)
        }
      })
    })
    it('should be false when permission is not granted', () => {
      getWrapper({
        setup: ({ actions }) => {
          expect(
            unref(actions)[0].isVisible({
              resources: [mock<SpaceResource>({ canEditImage: () => false })]
            })
          ).toBe(false)
        }
      })
    })
    it('should be true when permission is granted', () => {
      getWrapper({
        setup: ({ actions }) => {
          expect(
            unref(actions)[0].isVisible({
              resources: [mock<SpaceResource>({ canEditImage: () => true })]
            })
          ).toBe(true)
        }
      })
    })
  })
  describe('handler', () => {
    it('should trigger the setIcon modal window with one resource', () => {
      getWrapper({
        setup: async ({ actions }) => {
          const { dispatchModal } = useModals()
          await unref(actions)[0].handler({ resources: [{ id: '1' } as SpaceResource] })

          expect(dispatchModal).toHaveBeenCalledTimes(1)
        }
      })
    })
    it('should not trigger the setIcon modal window with no resource', () => {
      getWrapper({
        setup: async ({ actions }) => {
          const { dispatchModal } = useModals()
          await unref(actions)[0].handler({ resources: [] })

          expect(dispatchModal).toHaveBeenCalledTimes(0)
        }
      })
    })
  })
  describe('method "setIconSpace"', () => {
    it('should show message on success', () => {
      getWrapper({
        setup: async ({ setIconSpace }, { clientService }) => {
          clientService.graphAuthenticated.drives.updateDrive.mockResolvedValue(
            mock<SpaceResource>()
          )
          await setIconSpace(mock<SpaceResource>(), '🐻')

          const { showMessage } = useMessages()
          expect(showMessage).toHaveBeenCalledTimes(1)
        }
      })
    })

    it('should show message on error', () => {
      vi.spyOn(console, 'error').mockImplementation(() => undefined)
      getWrapper({
        setup: async ({ setIconSpace }, { clientService }) => {
          clientService.graphAuthenticated.drives.updateDrive.mockRejectedValue(new Error())
          await setIconSpace(mock<SpaceResource>(), '🐻')

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
    instance: ReturnType<typeof useSpaceActionsSetIcon>,
    {
      clientService
    }: {
      clientService: ReturnType<typeof defaultComponentMocks>['$clientService']
    }
  ) => void
}) {
  const mocks = defaultComponentMocks({
    currentRoute: mock<RouteLocation>({ name: 'files-spaces-generic' })
  })

  return {
    mocks,
    wrapper: getComposableWrapper(
      () => {
        const instance = useSpaceActionsSetIcon()
        setup(instance, { clientService: mocks.$clientService })
      },
      {
        mocks,
        provide: mocks
      }
    )
  }
}
