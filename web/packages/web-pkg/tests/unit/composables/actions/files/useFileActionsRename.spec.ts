import { useFileActionsRename } from '../../../../../src/composables/actions'
import { useResolvePasswordProtectedFolderCounterpart } from '../../../../../src/composables/actions/helpers/useResolvePasswordProtectedFolderCounterpart'
import {
  useConfigStore,
  useMessages,
  useModals,
  useResourcesStore
} from '../../../../../src/composables/piniaStores'
import { mock, mockDeep } from 'vitest-mock-extended'
import { Resource, SpaceResource } from '@ownclouders/web-client'
import { defaultComponentMocks, getComposableWrapper } from '@ownclouders/web-test-helpers'
import { unref } from 'vue'

vi.mock(
  '../../../../../src/composables/actions/helpers/useResolvePasswordProtectedFolderCounterpart'
)

const currentFolder = {
  id: '1',
  path: '/folder',
  spaceId: '1'
}

describe('rename', () => {
  describe('computed property "actions"', () => {
    describe('isVisible property of returned element', () => {
      it.each([
        { resources: [{ canRename: () => true }] as Resource[], expectedStatus: true },
        { resources: [{ canRename: () => false }] as Resource[], expectedStatus: false },
        {
          resources: [{ canRename: () => true }, { canRename: () => true }] as Resource[],
          expectedStatus: false
        },
        {
          resources: [{ canRename: () => true, locked: true }] as Resource[],
          expectedStatus: false
        }
      ])('should be set correctly', (inputData) => {
        getWrapper({
          setup: ({ actions }, { space }) => {
            const resources = inputData.resources
            expect(unref(actions)[0].isVisible({ space, resources })).toBe(inputData.expectedStatus)
          }
        })
      })
    })
  })

  describe('rename action handler', () => {
    it('should trigger the rename modal window', () => {
      getWrapper({
        setup: async ({ actions }, { space }) => {
          const { dispatchModal } = useModals()
          const resources = [currentFolder]
          await unref(actions)[0].handler({ space, resources })
          expect(dispatchModal).toHaveBeenCalledTimes(1)
        }
      })
    })
  })

  describe('method "getNameErrorMsg"', () => {
    it('should not show an error if new name not taken', () => {
      getWrapper({
        setup: ({ getNameErrorMsg }) => {
          const resourcesStore = useResourcesStore()
          resourcesStore.resources = [{ name: 'file1', path: '/file1' }] as Resource[]
          const message = getNameErrorMsg(
            { name: 'currentName', path: '/currentName' } as Resource,
            'newName'
          )
          expect(message).toEqual(null)
        }
      })
    })

    it('should not show an error if new name already exists but in different folder', () => {
      getWrapper({
        setup: ({ getNameErrorMsg }) => {
          const resourcesStore = useResourcesStore()
          resourcesStore.resources = [{ name: 'file1', path: '/file1' }] as Resource[]

          const message = getNameErrorMsg(
            mock<Resource>({ name: 'currentName', path: '/favorites/currentName' }),
            'file1'
          )
          expect(message).toEqual(null)
        }
      })
    })

    it.each([
      { currentName: 'currentName', newName: '', message: 'The name cannot be empty' },
      { currentName: 'currentName', newName: 'new/name', message: 'The name cannot contain "/"' },
      { currentName: 'currentName', newName: '.', message: 'The name cannot be equal to "."' },
      { currentName: 'currentName', newName: '..', message: 'The name cannot be equal to ".."' },
      {
        currentName: 'currentName',
        newName: 'newname ',
        message: 'The name cannot end with whitespace'
      },
      {
        currentName: 'currentName',
        newName: 'file1',
        message: 'The name "file1" is already taken'
      },
      {
        currentName: 'currentName',
        newName: 'newname',
        parentResources: [{ name: 'newname', path: '/newname' } as Resource],
        message: 'The name "newname" is already taken'
      }
    ])('should detect name errors and display error messages accordingly', (inputData) => {
      getWrapper({
        setup: ({ getNameErrorMsg }) => {
          const resourcesStore = useResourcesStore()
          resourcesStore.resources = [{ name: 'file1', path: '/file1' }] as Resource[]

          const message = getNameErrorMsg(
            mock<Resource>({ name: inputData.currentName, path: `/${inputData.currentName}` }),
            inputData.newName,
            inputData.parentResources
          )
          expect(message).toEqual(inputData.message)
        }
      })
    })
  })

  describe('method "renameResource"', () => {
    it('should call the rename action on a resource in the file list', () => {
      getWrapper({
        setup: async ({ renameResource }, { space }) => {
          const resource = {
            id: '2',
            path: '/folder',
            webDavPath: '/files/admin/folder',
            spaceId: '1'
          }
          await renameResource(space, resource, 'new name')

          const { upsertResource } = useResourcesStore()
          expect(upsertResource).toHaveBeenCalledTimes(1)
        }
      })
    })

    it('should call the rename action on the current folder', () => {
      getWrapper({
        setup: async ({ renameResource }, { space }) => {
          await renameResource(space, currentFolder, 'new name')

          const { upsertResource } = useResourcesStore()
          expect(upsertResource).toHaveBeenCalledTimes(1)
        }
      })
    })

    it('should handle errors properly', () => {
      vi.spyOn(console, 'error').mockImplementation(() => undefined)

      getWrapper({
        setup: async ({ renameResource }, { space, clientService }) => {
          clientService.webdav.moveFiles.mockRejectedValueOnce(new Error())

          await renameResource(space, currentFolder, 'new name')
          const { showErrorMessage } = useMessages()
          expect(showErrorMessage).toHaveBeenCalledTimes(1)
        }
      })
    })

    describe('password protected folder sync', () => {
      const passwordProtectedFolder = {
        id: 'ppf',
        isFolder: true,
        name: 'test',
        path: '/.PasswordProtectedFolders/projects/Personal/test',
        spaceId: '1'
      } as Resource

      it('renames the .psec counterpart when the real folder is renamed', () => {
        const psecSpace = mockDeep<SpaceResource>({ webDavPath: 'irrelevant' })
        const psecFile = mock<Resource>({ id: 'psec', name: 'test.psec', path: '/test.psec' })
        const getPsecFile = vi.fn().mockResolvedValue({ psecFile, space: psecSpace })

        getWrapper({
          getPsecFile,
          setup: async ({ renameResource }, { space, clientService }) => {
            await renameResource(space, passwordProtectedFolder, 'something')

            expect(getPsecFile).toHaveBeenCalledWith(passwordProtectedFolder)
            // both the folder itself and the .psec counterpart get moved
            expect(clientService.webdav.moveFiles).toHaveBeenCalledWith(
              psecSpace,
              psecFile,
              psecSpace,
              { path: '/something.psec' }
            )
            // the renamed folder + the renamed .psec file both get upserted
            const { upsertResource } = useResourcesStore()
            expect(upsertResource).toHaveBeenCalledTimes(2)
          }
        })
      })

      it('resolves the .psec counterpart before moving the folder (uses the old name)', () => {
        const psecSpace = mockDeep<SpaceResource>({ webDavPath: 'irrelevant' })
        const psecFile = mock<Resource>({ id: 'psec', name: 'test.psec', path: '/test.psec' })
        const getPsecFile = vi.fn().mockResolvedValue({ psecFile, space: psecSpace })

        getWrapper({
          getPsecFile,
          setup: async ({ renameResource }, { space, clientService }) => {
            const callOrder: string[] = []
            getPsecFile.mockImplementation(() => {
              callOrder.push('getPsecFile')
              return Promise.resolve({ psecFile, space: psecSpace })
            })
            clientService.webdav.moveFiles.mockImplementation(() => {
              callOrder.push('moveFiles')
              return Promise.resolve(undefined)
            })

            await renameResource(space, passwordProtectedFolder, 'something')

            // getPsecFile must run before any move so it can derive the path from the old folder name
            expect(callOrder[0]).toBe('getPsecFile')
          }
        })
      })

      it('does not resolve a .psec counterpart for a regular (non-ppf) rename', () => {
        const getPsecFile = vi.fn().mockResolvedValue(null)

        getWrapper({
          getPsecFile,
          setup: async ({ renameResource }, { space }) => {
            const resource = {
              id: '2',
              isFolder: true,
              name: 'folder',
              path: '/folder',
              spaceId: '1'
            } as Resource
            await renameResource(space, resource, 'new name')

            expect(getPsecFile).not.toHaveBeenCalled()
          }
        })
      })

      it('does not resolve a .psec counterpart when the name is unchanged', () => {
        const getPsecFile = vi.fn().mockResolvedValue(null)

        getWrapper({
          getPsecFile,
          setup: async ({ renameResource }, { space }) => {
            await renameResource(space, passwordProtectedFolder, passwordProtectedFolder.name)

            expect(getPsecFile).not.toHaveBeenCalled()
          }
        })
      })

      it('renames the folder normally when no .psec counterpart is found', () => {
        const getPsecFile = vi.fn().mockResolvedValue(null)

        getWrapper({
          getPsecFile,
          setup: async ({ renameResource }, { space, clientService }) => {
            await renameResource(space, passwordProtectedFolder, 'something')

            expect(getPsecFile).toHaveBeenCalledWith(passwordProtectedFolder)
            // only the folder itself is moved, no counterpart move
            expect(clientService.webdav.moveFiles).toHaveBeenCalledTimes(1)
            const { upsertResource } = useResourcesStore()
            expect(upsertResource).toHaveBeenCalledTimes(1)
          }
        })
      })
    })

    // When this app instance runs framed inside the password-protected-folder view modal, a
    // rename of the shared folder (the public link root) is posted to the parent window so it
    // can keep the owner's `.psec` pointer file in sync.
    describe('password protected folder view modal notification', () => {
      const publicSpace = mockDeep<SpaceResource>({
        driveType: 'public',
        webDavPath: 'irrelevant',
        fileId: 'root-file-id'
      })
      const rootFolder = mock<Resource>({
        fileId: 'root-file-id',
        isFolder: true,
        name: 'shared',
        path: '/shared'
      })

      let parentSpy: ReturnType<typeof vi.spyOn>
      let originalParent: Window

      beforeEach(() => {
        originalParent = window.parent
        // jsdom sets window.parent === window; simulate being framed
        Object.defineProperty(window, 'parent', {
          value: { postMessage: vi.fn() },
          configurable: true
        })
        parentSpy = vi.spyOn(window.parent, 'postMessage')
      })

      afterEach(() => {
        Object.defineProperty(window, 'parent', { value: originalParent, configurable: true })
      })

      it('notifies the parent window when the shared folder is renamed', () => {
        getWrapper({
          passwordProtectedFolderView: true,
          setup: async ({ renameResource }) => {
            await renameResource(publicSpace, rootFolder, 'renamed')

            expect(parentSpy).toHaveBeenCalledWith(
              {
                name: 'owncloud-password-protected-folder:renamed',
                data: { newName: 'renamed' }
              },
              window.location.origin
            )
          }
        })
      })

      it('does not notify for a child resource inside the shared folder', () => {
        const childResource = mock<Resource>({
          fileId: 'child-file-id',
          isFolder: true,
          name: 'child',
          path: '/child'
        })

        getWrapper({
          passwordProtectedFolderView: true,
          setup: async ({ renameResource }) => {
            await renameResource(publicSpace, childResource, 'renamed')

            expect(parentSpy).not.toHaveBeenCalled()
          }
        })
      })

      it('does not notify when not running inside the folder view modal', () => {
        getWrapper({
          passwordProtectedFolderView: false,
          setup: async ({ renameResource }) => {
            await renameResource(publicSpace, rootFolder, 'renamed')

            expect(parentSpy).not.toHaveBeenCalled()
          }
        })
      })

      it('does not notify for a non-public space', () => {
        const personalSpace = mockDeep<SpaceResource>({
          driveType: 'personal',
          webDavPath: 'irrelevant',
          fileId: 'root-file-id'
        })

        getWrapper({
          passwordProtectedFolderView: true,
          setup: async ({ renameResource }) => {
            await renameResource(personalSpace, rootFolder, 'renamed')

            expect(parentSpy).not.toHaveBeenCalled()
          }
        })
      })
    })
  })
})

function getWrapper({
  setup,
  getPsecFile = vi.fn().mockResolvedValue(null),
  passwordProtectedFolderView = false
}: {
  setup: (
    instance: ReturnType<typeof useFileActionsRename>,
    {
      space,
      clientService
    }: {
      space: SpaceResource
      clientService: ReturnType<typeof defaultComponentMocks>['$clientService']
    }
  ) => void
  getPsecFile?: ReturnType<typeof vi.fn>
  passwordProtectedFolderView?: boolean
}) {
  vi.mocked(useResolvePasswordProtectedFolderCounterpart).mockReturnValue({
    getPasswordProtectedFolder: vi.fn().mockResolvedValue(null),
    getPsecFile
  } as unknown as ReturnType<typeof useResolvePasswordProtectedFolderCounterpart>)

  const mocks = {
    ...defaultComponentMocks(),
    space: mockDeep<SpaceResource>({
      webDavPath: 'irrelevant'
    })
  }

  return {
    mocks,
    wrapper: getComposableWrapper(
      () => {
        const configStore = useConfigStore()
        configStore.options = { ...configStore.options, passwordProtectedFolderView }
        const instance = useFileActionsRename()
        setup(instance, { space: mocks.space, clientService: mocks.$clientService })
      },
      {
        mocks,
        provide: mocks
      }
    )
  }
}
