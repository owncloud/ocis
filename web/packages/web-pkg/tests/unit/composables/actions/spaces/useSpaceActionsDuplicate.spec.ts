import { useSpaceActionsDuplicate } from '../../../../../src/composables/actions/spaces'
import { AbilityRule, SpaceResource } from '@ownclouders/web-client'
import { mock } from 'vitest-mock-extended'
import {
  defaultComponentMocks,
  RouteLocation,
  getComposableWrapper
} from '@ownclouders/web-test-helpers'
import { unref } from 'vue'
import { ListFilesResult } from '@ownclouders/web-client/webdav'
import {
  useMessages,
  useResourcesStore,
  useSpacesStore
} from '../../../../../src/composables/piniaStores'

const spaces = [
  mock<SpaceResource>({
    name: 'Moon',
    description: 'To the moon',
    type: 'project',
    spaceImageData: null,
    spaceReadmeData: null,
    spaceQuota: { total: Math.pow(10, 9) }
  })
]
describe('restore', () => {
  describe('isVisible property', () => {
    it('should be false when no resource given', () => {
      getWrapper({
        setup: ({ actions }) => {
          expect(unref(actions)[0].isVisible({ resources: [] })).toBe(false)
        }
      })
    })
    it('should be false when the space is disabled', () => {
      getWrapper({
        setup: ({ actions }) => {
          expect(
            unref(actions)[0].isVisible({
              resources: [
                mock<SpaceResource>({
                  disabled: true,
                  driveType: 'project'
                })
              ]
            })
          ).toBe(false)
        }
      })
    })
    it('should be false when the space is no project space', () => {
      getWrapper({
        setup: ({ actions }) => {
          expect(
            unref(actions)[0].isVisible({
              resources: [
                mock<SpaceResource>({
                  disabled: false,
                  driveType: 'personal'
                })
              ]
            })
          ).toBe(false)
        }
      })
    })
    it('should be false when the current user can not create spaces', () => {
      getWrapper({
        abilities: [],
        setup: ({ actions }) => {
          expect(
            unref(actions)[0].isVisible({
              resources: [mock<SpaceResource>({ disabled: false, driveType: 'project' })]
            })
          ).toBe(false)
        }
      })
    })
    it('should be true when the current user can create spaces', () => {
      getWrapper({
        setup: ({ actions }) => {
          expect(
            unref(actions)[0].isVisible({
              resources: [
                mock<SpaceResource>({ name: 'Moon', disabled: false, driveType: 'project' }),
                mock<SpaceResource>({ name: 'Sun', disabled: false, driveType: 'project' })
              ]
            })
          ).toBe(true)
        }
      })
    })
  })
  describe('method "duplicateSpace"', () => {
    it('should show error message on error', () => {
      vi.spyOn(console, 'error').mockImplementation(() => undefined)
      getWrapper({
        setup: async ({ duplicateSpace }, { clientService }) => {
          clientService.graphAuthenticated.drives.createDrive.mockRejectedValue(new Error())
          await duplicateSpace(spaces[0])
          const { showErrorMessage } = useMessages()
          expect(showErrorMessage).toHaveBeenCalledTimes(1)
        }
      })
    })
    it('should show message on success', () => {
      getWrapper({
        setup: async ({ duplicateSpace }, { clientService }) => {
          clientService.graphAuthenticated.drives.createDrive.mockResolvedValue(
            mock<SpaceResource>({
              id: '1',
              name: 'Moon (1)'
            })
          )
          clientService.webdav.listFiles.mockResolvedValue({ children: [] } as ListFilesResult)
          await duplicateSpace(spaces[0])
          expect(clientService.graphAuthenticated.drives.createDrive).toHaveBeenCalledWith(
            {
              description: 'To the moon',
              name: 'Moon (1)',
              quota: {
                total: Math.pow(10, 9)
              }
            },
            expect.anything()
          )
          const spacesStore = useSpacesStore()
          expect(spacesStore.upsertSpace).toHaveBeenCalled()
          const { showMessage } = useMessages()
          expect(showMessage).toHaveBeenCalled()
        }
      })
    })
    it('should upsert a space as resource on the projects page', () => {
      getWrapper({
        currentRouteName: 'files-spaces-projects',
        setup: async ({ duplicateSpace }, { clientService }) => {
          clientService.graphAuthenticated.drives.createDrive.mockResolvedValue(
            mock<SpaceResource>({
              id: '1',
              name: 'Moon (1)'
            })
          )
          clientService.webdav.listFiles.mockResolvedValue({ children: [] } as ListFilesResult)
          await duplicateSpace(spaces[0])

          const { upsertResource } = useResourcesStore()
          expect(upsertResource).toHaveBeenCalled()
        }
      })
    })
  })
})

function getWrapper({
  setup,
  abilities = [{ action: 'create-all', subject: 'Drive' }],
  currentRouteName = 'files-spaces-generic'
}: {
  setup: (
    instance: ReturnType<typeof useSpaceActionsDuplicate>,
    {
      clientService
    }: {
      clientService: ReturnType<typeof defaultComponentMocks>['$clientService']
    }
  ) => void
  abilities?: AbilityRule[]
  currentRouteName?: string
}) {
  const mocks = defaultComponentMocks({
    currentRoute: mock<RouteLocation>({ name: currentRouteName })
  })
  return {
    mocks,
    wrapper: getComposableWrapper(
      () => {
        const instance = useSpaceActionsDuplicate()
        setup(instance, { clientService: mocks.$clientService })
      },
      {
        mocks,
        provide: mocks,
        pluginOptions: { abilities, piniaOptions: { spacesState: { spaces } } }
      }
    )
  }
}
