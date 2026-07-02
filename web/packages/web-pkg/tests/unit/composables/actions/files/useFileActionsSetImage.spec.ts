import { useFileActionsSetImage } from '../../../../../src'
import { useMessages } from '../../../../../src/composables/piniaStores'
import { Resource, SpaceResource } from '@ownclouders/web-client'
import { mock } from 'vitest-mock-extended'
import {
  defaultComponentMocks,
  RouteLocation,
  getComposableWrapper
} from '@ownclouders/web-test-helpers'
import { unref } from 'vue'
import { User } from '@ownclouders/web-client/graph/generated'
import { useSpaceHelpers } from '../../../../../src/composables/spaces/useSpaceHelpers'

vi.mock('../../../../../src/composables/spaces/useSpaceHelpers', () => ({
  useSpaceHelpers: vi.fn()
}))

describe('setImage', () => {
  describe('isVisible property', () => {
    it('should be false when no resource given', () => {
      const space = mock<SpaceResource>({ canEditImage: () => true })

      getWrapper({
        setup: ({ actions }) => {
          expect(unref(actions)[0].isVisible({ space, resources: [] as Resource[] })).toBe(false)
        }
      })
    })
    it('should be false when mimeType is not image', () => {
      const space = mock<SpaceResource>({ canEditImage: () => true })
      getWrapper({
        setup: ({ actions }) => {
          expect(
            unref(actions)[0].isVisible({
              space,
              resources: [{ id: '1', mimeType: 'text/plain' }] as Resource[]
            })
          ).toBe(false)
        },
        isMimetypeSupported: false
      })
    })
    it('should be true when the mimeType is image', () => {
      const space = mock<SpaceResource>({ canEditImage: () => true })
      getWrapper({
        setup: ({ actions }) => {
          expect(
            unref(actions)[0].isVisible({
              space,
              resources: [{ id: '1', mimeType: 'image/png' }] as Resource[]
            })
          ).toBe(true)
        }
      })
    })
    it('should be false when canEditImage is false', () => {
      const space = mock<SpaceResource>({ canEditImage: () => false })
      getWrapper({
        setup: ({ actions }) => {
          expect(
            unref(actions)[0].isVisible({
              space,
              resources: [{ id: '1', mimeType: 'image/png' }] as Resource[]
            })
          ).toBe(false)
        }
      })
    })
  })

  describe('handler', () => {
    it('should show message on success', () => {
      const space = mock<SpaceResource>({ id: '1' })

      getWrapper({
        setup: async ({ actions }, { clientService }) => {
          clientService.graphAuthenticated.drives.updateDrive.mockResolvedValue(space)
          await unref(actions)[0].handler({
            space,
            resources: [
              {
                webDavPath: '/spaces/1fe58d8b-aa69-4c22-baf7-97dd57479f22/subfolder/image.png',
                name: 'image.png'
              }
            ] as Resource[]
          })
          const { showMessage } = useMessages()
          expect(showMessage).toHaveBeenCalledTimes(1)
        }
      })
    })

    it('should show message on error', () => {
      vi.spyOn(console, 'error').mockImplementation(() => undefined)
      const space = mock<SpaceResource>({ id: '1' })
      getWrapper({
        setup: async ({ actions }) => {
          await unref(actions)[0].handler({
            space,
            resources: [
              {
                webDavPath: '/spaces/1fe58d8b-aa69-4c22-baf7-97dd57479f22/subfolder/image.png',
                name: 'image.png'
              }
            ] as Resource[]
          })
          const { showErrorMessage } = useMessages()
          expect(showErrorMessage).toHaveBeenCalledTimes(1)
        }
      })
    })
  })
})

function getWrapper({
  setup,
  isMimetypeSupported = true
}: {
  setup: (
    instance: ReturnType<typeof useFileActionsSetImage>,
    options: {
      clientService: ReturnType<typeof defaultComponentMocks>['$clientService']
    }
  ) => void
  isMimetypeSupported?: boolean
}) {
  vi.mocked(useSpaceHelpers).mockReturnValue({
    checkSpaceNameModalInput: vi.fn(),
    getDefaultMetaFolder: () => new Promise(() => mock<Resource>())
  })

  const mocks = {
    ...defaultComponentMocks({
      currentRoute: mock<RouteLocation>({ name: 'files-spaces-generic' })
    })
  }
  mocks.$previewService.isMimetypeSupported.mockReturnValue(isMimetypeSupported)
  mocks.$clientService.webdav.getFileInfo.mockResolvedValue(mock<Resource>())

  return {
    wrapper: getComposableWrapper(
      () => {
        const instance = useFileActionsSetImage()
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
