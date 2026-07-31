import { useSpaceActionsDisable } from '../../../../../src/composables/actions/spaces'
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

describe('disable', () => {
  describe('isVisible property', () => {
    it('should be false when no resource given', () => {
      getWrapper({
        setup: ({ actions }) => {
          expect(unref(actions)[0].isVisible({ resources: [] })).toBe(false)
        }
      })
    })
    it('should be true when the space can be disabled', () => {
      const spaceMock = mock<SpaceResource>({ driveType: 'project', canDisable: () => true })
      getWrapper({
        setup: ({ actions }) => {
          expect(unref(actions)[0].isVisible({ resources: [spaceMock] })).toBe(true)
        }
      })
    })
    it('should be false when the space can not be disabled', () => {
      const spaceMock = mock<SpaceResource>({ driveType: 'project', canDisable: () => false })
      getWrapper({
        setup: ({ actions }) => {
          expect(unref(actions)[0].isVisible({ resources: [spaceMock] })).toBe(false)
        }
      })
    })
  })

  describe('handler', () => {
    it('should trigger the disable modal window', () => {
      getWrapper({
        setup: async ({ actions }) => {
          const { dispatchModal } = useModals()
          await unref(actions)[0].handler({
            resources: [
              mock<SpaceResource>({ id: '1', canDisable: () => true, driveType: 'project' })
            ]
          })

          expect(dispatchModal).toHaveBeenCalledTimes(1)
        }
      })
    })
    it('should not trigger the disable modal window without any resource', () => {
      getWrapper({
        setup: async ({ actions }) => {
          const { dispatchModal } = useModals()
          await unref(actions)[0].handler({
            resources: [
              mock<SpaceResource>({ id: '1', canDisable: () => false, driveType: 'project' })
            ]
          })

          expect(dispatchModal).toHaveBeenCalledTimes(0)
        }
      })
    })
  })

  describe('method "disableSpace"', () => {
    it('should show message on success', () => {
      getWrapper({
        setup: async ({ disableSpaces }, { clientService }) => {
          clientService.graphAuthenticated.drives.disableDrive.mockResolvedValue()
          await disableSpaces([
            mock<SpaceResource>({ id: '1', canDisable: () => true, driveType: 'project' })
          ])

          const { showMessage } = useMessages()
          expect(showMessage).toHaveBeenCalledTimes(1)
        }
      })
    })

    it('should show message on error', () => {
      vi.spyOn(console, 'error').mockImplementation(() => undefined)
      getWrapper({
        setup: async ({ disableSpaces }, { clientService }) => {
          clientService.graphAuthenticated.drives.disableDrive.mockRejectedValue(new Error())
          await disableSpaces([
            mock<SpaceResource>({ id: '1', canDisable: () => true, driveType: 'project' })
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
    instance: ReturnType<typeof useSpaceActionsDisable>,
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
        const instance = useSpaceActionsDisable()
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
