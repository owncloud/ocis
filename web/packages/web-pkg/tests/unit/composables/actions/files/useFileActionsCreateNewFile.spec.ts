import { mock } from 'vitest-mock-extended'
import { ref, unref } from 'vue'
import { useFileActionsCreateNewFile } from '../../../../../src/composables/actions'
import { useModals } from '../../../../../src/composables/piniaStores'
import { SpaceResource } from '@ownclouders/web-client'
import { Resource } from '@ownclouders/web-client'
import { FileActionOptions } from '../../../../../src/composables/actions'
import { useFileActions } from '../../../../../src/composables/actions/files/useFileActions'
import {
  RouteLocation,
  defaultComponentMocks,
  getComposableWrapper
} from '@ownclouders/web-test-helpers'
import { ApplicationFileExtension } from '../../../../../types'
import { useResourcesStore } from '../../../../../src/composables/piniaStores'

vi.mock('../../../../../src/composables/actions/files/useFileActions', async (importOriginal) => ({
  ...(await importOriginal<any>()),
  useFileActions: vi.fn(() => mock<ReturnType<typeof useFileActions>>())
}))

describe('useFileActionsCreateNewFile', () => {
  describe('checkFileName', () => {
    it.each([
      { input: '', output: 'File name cannot be empty' },
      { input: '/', output: 'File name cannot contain "/"' },
      { input: '.', output: 'File name cannot be equal to "."' },
      { input: '..', output: 'File name cannot be equal to ".."' },
      { input: 'myfile.txt', output: null }
    ])('should validate file name %s', (data) => {
      const space = mock<SpaceResource>({ id: '1' })
      getWrapper({
        space,
        setup: ({ getNameErrorMsg }) => {
          const result = getNameErrorMsg(data.input)
          expect(result).toBe(data.output)
        }
      })
    })
  })

  describe('openFile', () => {
    it('upserts the resource before opening', () => {
      const space = mock<SpaceResource>({ id: '1' })
      getWrapper({
        space,
        setup: ({ openFile }) => {
          openFile(mock<Resource>(), null)

          const { upsertResource } = useResourcesStore()
          expect(upsertResource).toHaveBeenCalled()
        }
      })
    })
  })
  describe('createNewFileModal', () => {
    it('should show modal', () => {
      const space = mock<SpaceResource>({ id: '1' })
      getWrapper({
        space,
        setup: async ({ actions }) => {
          const { dispatchModal } = useModals()
          const fileActionOptions: FileActionOptions = { space, resources: [] } as FileActionOptions
          await unref(actions)[0].handler(fileActionOptions)

          expect(dispatchModal).toHaveBeenCalled()
        }
      })
    })
  })

  describe('customHandler', () => {
    it('should use customHandler when it is provided', () => {
      const action = mock<ApplicationFileExtension>({ app: 'link-opener', customHandler: vi.fn() })

      getWrapper({
        action,
        setup: ({ actions }) => {
          unref(actions).at(0).handler()
          expect(action.customHandler).toHaveBeenCalled()
        }
      })
    })
  })

  describe('actions', () => {
    it('should check if the new file menu is visible when the isVisible is a function', () => {
      const isVisible = vi.fn().mockReturnValue(false)
      const action = mock<ApplicationFileExtension>({
        app: 'text-editor',
        extension: '.txt',
        newFileMenu: { menuTitle: vi.fn(), isVisible },
        customHandler: null
      })

      getWrapper({
        action,
        setup(instance) {
          expect(unref(instance.actions).at(0).isVisible()).toBe(false)
          expect(isVisible).toHaveBeenCalled()
        }
      })
    })

    it('should not check if the new file menu is visible when the isVisible is not a function', () => {
      const action = mock<ApplicationFileExtension>({
        app: 'text-editor',
        extension: '.txt',
        newFileMenu: { menuTitle: vi.fn(), isVisible: undefined },
        customHandler: null
      })

      getWrapper({
        action,
        currentFolder: mock<Resource>({ canUpload: () => true }),
        setup(instance) {
          expect(unref(instance.actions).at(0).isVisible()).toBe(true)
        }
      })
    })
  })
})

function getWrapper({
  resolveCreateFile = true,
  space = undefined,
  setup,
  action = mock<ApplicationFileExtension>({
    app: 'text-editor',
    extension: '.txt',
    newFileMenu: { menuTitle: vi.fn() },
    customHandler: null
  }),
  currentFolder = mock<Resource>({ id: '1', path: '/' })
}: {
  resolveCreateFile?: boolean
  space?: SpaceResource
  setup: (instance: ReturnType<typeof useFileActionsCreateNewFile>) => void
  action?: ApplicationFileExtension
  currentFolder?: Resource
}) {
  const mocks = {
    ...defaultComponentMocks({
      currentRoute: mock<RouteLocation>({ name: 'files-spaces-generic' })
    }),
    space
  }
  mocks.$clientService.webdav.putFileContents.mockImplementation(() => {
    if (resolveCreateFile) {
      return Promise.resolve({
        id: '1',
        type: 'folder',
        path: '/',
        isReceivedShare: vi.fn(),
        spaceId: '1'
      } as Resource)
    }
    return Promise.reject('error')
  })

  return {
    wrapper: getComposableWrapper(
      () => {
        const instance = useFileActionsCreateNewFile({ space: ref(space) })
        setup(instance)
      },
      {
        provide: mocks,
        mocks,
        pluginOptions: {
          piniaOptions: {
            appsState: {
              fileExtensions: [action]
            },
            resourcesStore: { currentFolder }
          }
        }
      }
    )
  }
}
