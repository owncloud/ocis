import { mock } from 'vitest-mock-extended'
import { ref, unref } from 'vue'
import {
  defaultComponentMocks,
  RouteLocation,
  getComposableWrapper
} from '@ownclouders/web-test-helpers'
import { useFileActionsCreateNewShortcut, useModals } from '../../../../../src/composables'
import { Resource, SpaceResource } from '@ownclouders/web-client'

describe('createNewShortcut', () => {
  describe('computed property "actions"', () => {
    describe('method "isVisible"', () => {
      it.each([
        {
          currentFolderCanCreate: true,
          expectedStatus: true
        },
        {
          currentFolderCanCreate: false,
          expectedStatus: false
        }
      ])('should be set correctly', ({ currentFolderCanCreate, expectedStatus }) => {
        getWrapper({
          currentFolder: mock<Resource>({ canCreate: () => currentFolderCanCreate }),
          setup: ({ actions }) => {
            expect(unref(actions)[0].isVisible()).toBe(expectedStatus)
          }
        })
      })
    })
    describe('method "handler"', () => {
      it('creates a modal', () => {
        getWrapper({
          setup: async ({ actions }) => {
            const { dispatchModal } = useModals()
            await unref(actions)[0].handler()
            expect(dispatchModal).toHaveBeenCalled()
          }
        })
      })
    })
  })
})

function getWrapper({
  setup,
  currentFolder = mock<Resource>()
}: {
  setup: (instance: ReturnType<typeof useFileActionsCreateNewShortcut>) => void
  currentFolder?: Resource
}) {
  const mocks = {
    ...defaultComponentMocks({
      currentRoute: mock<RouteLocation>({ name: 'files-spaces-generic' })
    })
  }

  return {
    wrapper: getComposableWrapper(
      () => {
        const instance = useFileActionsCreateNewShortcut({ space: ref(mock<SpaceResource>()) })
        setup(instance)
      },
      {
        mocks,
        provide: mocks,
        pluginOptions: { piniaOptions: { resourcesStore: { currentFolder } } }
      }
    )
  }
}
