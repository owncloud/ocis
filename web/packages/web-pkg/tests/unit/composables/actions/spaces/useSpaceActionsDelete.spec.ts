import { useSpaceActionsDelete } from '../../../../../src/composables/actions'
import { useMessages, useModals } from '../../../../../src/composables/piniaStores'
import { SpaceResource } from '@ownclouders/web-client'
import {
  defaultComponentMocks,
  RouteLocation,
  getComposableWrapper
} from '@ownclouders/web-test-helpers'
import { mock } from 'vitest-mock-extended'
import { unref } from 'vue'
import { User } from '@ownclouders/web-client/graph/generated'

describe('delete', () => {
  describe('isVisible property', () => {
    it('should be false when no resource given', () => {
      getWrapper({
        setup: ({ actions }) => {
          expect(unref(actions)[0].isVisible({ resources: [] })).toBe(false)
        }
      })
    })
    it('should be false when the space can not be deleted', () => {
      const spaceMock = mock<SpaceResource>({ driveType: 'project', canBeDeleted: () => false })
      getWrapper({
        setup: ({ actions }) => {
          expect(unref(actions)[0].isVisible({ resources: [spaceMock] })).toBe(false)
        }
      })
    })
    it('should be true when the space can be deleted', () => {
      const spaceMock = mock<SpaceResource>({ driveType: 'project', canBeDeleted: () => true })
      getWrapper({
        setup: ({ actions }) => {
          expect(unref(actions)[0].isVisible({ resources: [spaceMock] })).toBe(true)
        }
      })
    })
  })

  describe('handler', () => {
    it('should trigger the delete modal window', () => {
      getWrapper({
        setup: async ({ actions }) => {
          const { dispatchModal } = useModals()
          await unref(actions)[0].handler({
            resources: [
              mock<SpaceResource>({ id: '1', canBeDeleted: () => true, driveType: 'project' })
            ]
          })

          expect(dispatchModal).toHaveBeenCalledTimes(1)
        }
      })
    })
    it('should not trigger the delete modal window without any resource to delete', () => {
      getWrapper({
        setup: async ({ actions }) => {
          const { dispatchModal } = useModals()
          await unref(actions)[0].handler({
            resources: [
              mock<SpaceResource>({ id: '1', canBeDeleted: () => false, driveType: 'project' })
            ]
          })

          expect(dispatchModal).toHaveBeenCalledTimes(0)
        }
      })
    })
  })

  describe('method "deleteSpace"', () => {
    it('should show message on success', () => {
      getWrapper({
        setup: async ({ deleteSpaces }, { clientService }) => {
          clientService.graphAuthenticated.drives.deleteDrive.mockResolvedValue()

          await deleteSpaces([
            mock<SpaceResource>({ id: '1', canBeDeleted: () => true, driveType: 'project' })
          ])

          const { showMessage } = useMessages()
          expect(showMessage).toHaveBeenCalledTimes(1)
        }
      })
    })

    it('should show message on error', () => {
      vi.spyOn(console, 'error').mockImplementation(() => undefined)
      getWrapper({
        setup: async ({ deleteSpaces }, { clientService }) => {
          clientService.graphAuthenticated.drives.deleteDrive.mockRejectedValue(new Error())
          await deleteSpaces([
            mock<SpaceResource>({ id: '1', canBeDeleted: () => true, driveType: 'project' })
          ])

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
    instance: ReturnType<typeof useSpaceActionsDelete>,
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
        const instance = useSpaceActionsDelete()
        setup(instance, { clientService: mocks.$clientService })
      },
      {
        mocks,
        provide: mocks,
        pluginOptions: {
          piniaOptions: {
            userState: { user: { id: '1', onPremisesSamAccountName: 'alice' } as User }
          }
        }
      }
    )
  }
}
