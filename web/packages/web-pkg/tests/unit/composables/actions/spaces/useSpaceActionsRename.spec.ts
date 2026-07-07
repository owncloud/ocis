import { useSpaceActionsRename } from '../../../../../src/composables/actions/spaces'
import { useMessages, useModals } from '../../../../../src/composables/piniaStores'
import { mock } from 'vitest-mock-extended'
import {
  defaultComponentMocks,
  RouteLocation,
  getComposableWrapper
} from '@ownclouders/web-test-helpers'
import { unref } from 'vue'
import { SpaceResource } from '@ownclouders/web-client'

describe('rename', () => {
  describe('handler', () => {
    it('should trigger the rename modal window', () => {
      getWrapper({
        setup: async ({ actions }) => {
          const { dispatchModal } = useModals()
          await unref(actions)[0].handler({
            resources: [{ id: '1', name: 'renamed space' } as SpaceResource]
          })

          expect(dispatchModal).toHaveBeenCalledTimes(1)
        }
      })
    })
    it('should not trigger the rename modal window without any resource', () => {
      getWrapper({
        setup: async ({ actions }) => {
          const { dispatchModal } = useModals()
          await unref(actions)[0].handler({ resources: [] })

          expect(dispatchModal).toHaveBeenCalledTimes(0)
        }
      })
    })
  })
  describe('method "renameSpace"', () => {
    it('should show message on success', () => {
      getWrapper({
        setup: async ({ renameSpace }, { clientService }) => {
          clientService.graphAuthenticated.drives.updateDrive.mockResolvedValue(
            mock<SpaceResource>()
          )
          await renameSpace(mock<SpaceResource>({ id: '1' }), 'renamed space')

          const { showMessage } = useMessages()
          expect(showMessage).toHaveBeenCalledTimes(1)
        }
      })
    })

    it('should show message on error', () => {
      vi.spyOn(console, 'error').mockImplementation(() => undefined)
      getWrapper({
        setup: async ({ renameSpace }, { clientService }) => {
          clientService.graphAuthenticated.drives.updateDrive.mockRejectedValue(new Error())
          await renameSpace(mock<SpaceResource>({ id: '1' }), 'renamed space')

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
    instance: ReturnType<typeof useSpaceActionsRename>,
    {
      clientService
    }: {
      clientService: ReturnType<typeof defaultComponentMocks>['$clientService']
    }
  ) => void
}) {
  const mocks = defaultComponentMocks({
    currentRoute: mock<RouteLocation>({ name: 'files-spaces-projects' })
  })
  return {
    mocks,
    wrapper: getComposableWrapper(
      () => {
        const instance = useSpaceActionsRename()
        setup(instance, { clientService: mocks.$clientService })
      },
      {
        mocks,
        provide: mocks
      }
    )
  }
}
