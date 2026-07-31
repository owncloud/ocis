import { mock } from 'vitest-mock-extended'
import { computed, unref } from 'vue'
import {
  defaultComponentMocks,
  RouteLocation,
  getComposableWrapper
} from '@ownclouders/web-test-helpers'
import {
  useFileActionsOpenWithApp,
  useIsFilesAppActive,
  useModals
} from '../../../../../src/composables'
import { Resource, SpaceResource } from '@ownclouders/web-client'
import { ApplicationInformation } from '../../../../../src'

vi.mock('../../../../../src/composables/actions/helpers/useIsFilesAppActive')

const spaceMock = mock<SpaceResource>({
  id: '1'
})
describe('openWithApp', () => {
  describe('computed property "actions"', () => {
    describe('method "isVisible"', () => {
      it.each([
        {
          isFilesAppActive: false,
          expectedStatus: true
        },
        {
          isFilesAppActive: true,
          expectedStatus: false
        }
      ])('should be set correctly', ({ isFilesAppActive, expectedStatus }) => {
        getWrapper({
          isFilesAppActive,
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
            await unref(actions)[0].handler({
              resources: [mock<Resource>({ spaceId: spaceMock.id, path: '/' })],
              space: mock<SpaceResource>()
            })
            expect(dispatchModal).toHaveBeenCalled()
          }
        })
      })
    })
  })
})

function getWrapper({
  setup,
  isFilesAppActive = false
}: {
  setup: (instance: ReturnType<typeof useFileActionsOpenWithApp>) => void
  isFilesAppActive?: boolean
}) {
  vi.mocked(useIsFilesAppActive).mockReturnValueOnce(computed(() => isFilesAppActive))

  const mocks = {
    ...defaultComponentMocks({
      currentRoute: mock<RouteLocation>({ name: 'text-editor' })
    })
  }

  return {
    wrapper: getComposableWrapper(
      () => {
        const instance = useFileActionsOpenWithApp({ appId: 'text-editor' })
        setup(instance)
      },
      {
        mocks,
        provide: mocks,
        pluginOptions: {
          piniaOptions: {
            spacesState: { spaces: [spaceMock] },
            appsState: {
              apps: { 'text-editor': mock<ApplicationInformation>({ name: 'text-editor' }) }
            }
          }
        }
      }
    )
  }
}
