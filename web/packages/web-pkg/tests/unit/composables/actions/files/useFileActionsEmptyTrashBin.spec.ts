import { useFileActionsEmptyTrashBin } from '../../../../../src/composables/actions'
import { useMessages, useModals } from '../../../../../src/composables/piniaStores'
import { mock } from 'vitest-mock-extended'
import {
  getComposableWrapper,
  defaultComponentMocks,
  RouteLocation
} from '@ownclouders/web-test-helpers'
import { unref } from 'vue'
import { ProjectSpaceResource, TrashResource } from '@ownclouders/web-client'
import { FileActionOptions } from '../../../../../src/composables/actions'

describe('emptyTrashBin', () => {
  describe('isVisible property', () => {
    it('should be false when location is invalid', () => {
      getWrapper({
        invalidLocation: true,
        setup: ({ actions }, { space }) => {
          expect(unref(actions)[0].isVisible({ space, resources: [] })).toBe(false)
        }
      })
    })
    it('should be false in a space trash bin with insufficient permissions', () => {
      getWrapper({
        driveType: 'project',
        setup: ({ actions }, { space }) => {
          expect(
            unref(actions)[0].isVisible({
              space,
              resources: [{ canBeRestored: () => true }] as TrashResource[]
            })
          ).toBe(false)
        }
      })
    })
  })

  describe('empty trashbin action', () => {
    it('should trigger the empty trash bin modal window', () => {
      getWrapper({
        setup: async ({ actions }) => {
          const { dispatchModal } = useModals()
          await unref(actions)[0].handler(mock<FileActionOptions>())

          expect(dispatchModal).toHaveBeenCalledTimes(1)
        }
      })
    })
  })

  describe('method "emptyTrashBin"', () => {
    it('should show message on success', () => {
      getWrapper({
        setup: async ({ emptyTrashBin }, { space }) => {
          await emptyTrashBin({ space })

          const { showMessage } = useMessages()
          expect(showMessage).toHaveBeenCalledTimes(1)
        }
      })
    })

    it('should show message on error', () => {
      vi.spyOn(console, 'error').mockImplementation(() => undefined)

      getWrapper({
        resolveClearTrashBin: false,
        setup: async ({ emptyTrashBin }, { space }) => {
          await emptyTrashBin({ space })

          const { showErrorMessage } = useMessages()
          expect(showErrorMessage).toHaveBeenCalledTimes(1)
        }
      })
    })
  })
})

function getWrapper({
  invalidLocation = false,
  resolveClearTrashBin = true,
  driveType = 'personal',
  setup
}: {
  invalidLocation?: boolean
  resolveClearTrashBin?: boolean
  driveType?: string
  setup: (
    instance: ReturnType<typeof useFileActionsEmptyTrashBin>,
    {
      space
    }: {
      space: ProjectSpaceResource
    }
  ) => void
}) {
  const mocks = {
    ...defaultComponentMocks({
      currentRoute: mock<RouteLocation>({
        name: invalidLocation ? 'files-spaces-generic' : 'files-trash-generic'
      })
    }),
    space: mock<ProjectSpaceResource>({ driveType })
  }

  if (resolveClearTrashBin) {
    mocks.$clientService.webdav.clearTrashBin.mockResolvedValue(undefined)
  } else {
    mocks.$clientService.webdav.clearTrashBin.mockRejectedValue(new Error(''))
  }

  return {
    wrapper: getComposableWrapper(
      () => {
        const instance = useFileActionsEmptyTrashBin()
        setup(instance, { space: mocks.space })
      },
      {
        mocks,
        provide: mocks
      }
    )
  }
}
