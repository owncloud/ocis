import { useFileActionsDeleteResources } from '../../../../../src/composables/actions'
import { mock, mockDeep } from 'vitest-mock-extended'
import { FolderResource, Resource, SpaceResource, TrashResource } from '@ownclouders/web-client'
import {
  defaultComponentMocks,
  getComposableWrapper,
  RouteLocation,
  useGetMatchingSpaceMock
} from '../../../../../src/testing'
import { useDeleteWorker } from '../../../../../src/composables/webWorkers/deleteWorker'
import { useRestoreWorker } from '../../../../../src/composables/webWorkers/restoreWorker'
import { useGetMatchingSpace } from '../../../../../src/composables/spaces/useGetMatchingSpace'
import {
  useMessages,
  useResourcesStore,
  useSpacesStore
} from '../../../../../src/composables/piniaStores'
import { MockedFunction } from 'vitest'

vi.mock('../../../../../src/composables/webWorkers/deleteWorker')
vi.mock('../../../../../src/composables/webWorkers/restoreWorker')
vi.mock('../../../../../src/composables/spaces/useGetMatchingSpace')

const currentFolder = {
  id: '1',
  path: '/folder',
  spaceId: '1'
}

const passwordProtectedFolder = mock<Resource>({
  path: '/.PasswordProtectedFolders/projects/Personal/folder/psecFolder',
  storageId: 'personal',
  canBeDeleted: () => true
})

describe('deleteResources', () => {
  describe('method "filesList_delete"', () => {
    it('should call the delete action on a resource in the file list', () => {
      const filesToDelete = [{ id: '2', path: '/folder/fileToDelete.txt', spaceId: '1' }]

      getWrapper({
        currentFolder,
        result: filesToDelete,
        setup: ({ filesList_delete }, { router }) => {
          filesList_delete(filesToDelete)

          expect(router.push).toHaveBeenCalledTimes(0)
        }
      })
    })

    it('should call the delete action on the current folder', () => {
      const resourcesToDelete = [currentFolder]
      getWrapper({
        currentFolder,
        setup: ({ filesList_delete }, { router }) => {
          filesList_delete(resourcesToDelete)

          expect(router.push).toHaveBeenCalledTimes(1)
        }
      })
    })

    it('should push resources into delete queue', () => {
      const filesToDelete = [{ id: '2', path: '/folder/fileToDelete.txt', spaceId: '1' }]
      getWrapper({
        currentFolder,
        result: filesToDelete,
        setup: ({ filesList_delete }) => {
          filesList_delete(filesToDelete)
        }
      })

      const { addResourcesIntoDeleteQueue } = useResourcesStore()
      expect(addResourcesIntoDeleteQueue).toHaveBeenCalledWith(['2'])
    })

    it('should delete password protected folders when deleting psec file', () => {
      const filesToDelete = [
        mock<Resource>({
          id: '2',
          path: '/folder/psecFolder.psec',
          storageId: 'personal',
          extension: 'psec',
          name: 'psecFolder.psec',
          spaceId: '1'
        })
      ]
      getWrapper({
        currentFolder,
        getFileInfoResult: passwordProtectedFolder,
        setup: async ({ filesList_delete }, { space }) => {
          await filesList_delete(filesToDelete)

          const { startWorker } = vi.mocked(useDeleteWorker)()
          expect(startWorker).toHaveBeenCalledWith(
            {
              resources: [...filesToDelete, passwordProtectedFolder],
              space: space,
              topic: 'fileListDelete'
            },
            expect.any(Function)
          )
        }
      })
    })

    it('should delete psec file when deleting password protected folder', () => {
      const psecFile = mock<Resource>({
        id: '2',
        path: '/folder/psecFolder.psec',
        spaceId: '1',
        storageId: 'personal',
        extension: 'psec',
        name: 'psecFolder.psec',
        canBeDeleted: () => true
      })

      getWrapper({
        currentFolder,
        getFileInfoResult: psecFile,
        setup: async ({ filesList_delete }, { space }) => {
          const { getSpacesByName } = useSpacesStore()
          ;(getSpacesByName as MockedFunction<typeof getSpacesByName>).mockReturnValue([space])

          await filesList_delete([passwordProtectedFolder])

          const { startWorker } = vi.mocked(useDeleteWorker)()
          expect(startWorker).toHaveBeenCalledWith(
            {
              resources: [passwordProtectedFolder, psecFile],
              space: space,
              topic: 'fileListDelete'
            },
            expect.any(Function)
          )
        }
      })
    })

    it('updates the size of the project space after deleting from its root', async () => {
      const projectSpace = {
        id: 'project-1',
        driveType: 'project',
        size: 104947000,
        spaceQuota: { used: 104947000 }
      } as unknown as SpaceResource
      const deletedFile = mock<Resource>({
        id: '2',
        name: 'big.bin',
        path: '/big.bin',
        storageId: projectSpace.id,
        spaceId: projectSpace.id
      })
      const updatedDrive = mock<SpaceResource>({
        id: projectSpace.id,
        size: 47000,
        spaceQuota: { used: 47000 }
      })

      const { getWorkerCallbackDone } = getWrapper({
        currentFolder: { id: projectSpace.id, path: '/', spaceId: projectSpace.id },
        currentRouteName: 'files-spaces-generic',
        spaceDriveType: 'project',
        result: [deletedFile],
        getDriveResult: updatedDrive,
        stubActions: false,
        setup: ({ filesList_delete }) => {
          useSpacesStore().addSpaces([projectSpace])
          filesList_delete([deletedFile])
        }
      })
      await getWorkerCallbackDone()

      const space = useSpacesStore().spaces.find(({ id }) => id === projectSpace.id)
      expect(space.size).toBe(47000)
      expect(space.spaceQuota).toEqual(updatedDrive.spaceQuota)
    })

    describe('undo action on the success message', () => {
      const deletedFile = mock<Resource>({
        id: '2',
        name: 'fileToDelete.txt',
        path: '/folder/fileToDelete.txt',
        storageId: 'personal',
        spaceId: '1'
      })

      it('is included when the space allows trash-restore and the trash entry resolves', async () => {
        const trashEntry = mock<TrashResource>({
          name: deletedFile.name,
          path: deletedFile.path,
          ddate: '2026-01-01T00:00:00Z'
        })

        const { getWorkerCallbackDone } = getWrapper({
          currentFolder,
          result: [deletedFile],
          listFilesResult: { resource: mock<Resource>(), children: [trashEntry] },
          setup: ({ filesList_delete }) => {
            filesList_delete([deletedFile])
          }
        })
        await getWorkerCallbackDone()

        const { showMessage } = useMessages()
        expect(showMessage).toHaveBeenCalledWith(
          expect.objectContaining({
            timeout: 5,
            actions: [expect.objectContaining({ label: 'Undo' })]
          })
        )
      })

      it('is omitted when the trash entry cannot be resolved', async () => {
        const { getWorkerCallbackDone } = getWrapper({
          currentFolder,
          result: [deletedFile],
          listFilesResult: { resource: mock<Resource>(), children: [] },
          setup: ({ filesList_delete }) => {
            filesList_delete([deletedFile])
          }
        })
        await getWorkerCallbackDone()

        const { showMessage } = useMessages()
        expect(showMessage).toHaveBeenCalledWith(
          expect.not.objectContaining({ actions: expect.anything() })
        )
      })

      it('is omitted when the space does not allow trash-restore', async () => {
        const { getWorkerCallbackDone } = getWrapper({
          currentFolder,
          result: [deletedFile],
          spaceDriveType: 'project',
          canRestoreFromTrashbin: false,
          setup: ({ filesList_delete }) => {
            filesList_delete([deletedFile])
          }
        })
        await getWorkerCallbackDone()

        const { showMessage } = useMessages()
        expect(showMessage).toHaveBeenCalledWith(
          expect.not.objectContaining({ actions: expect.anything() })
        )
      })

      it('restores the resolved trash entry when clicked', async () => {
        const trashEntry = mock<TrashResource>({
          name: deletedFile.name,
          path: deletedFile.path,
          ddate: '2026-01-01T00:00:00Z'
        })

        let spaceUnderTest: SpaceResource
        const { getWorkerCallbackDone } = getWrapper({
          currentFolder,
          result: [deletedFile],
          listFilesResult: { resource: mock<Resource>(), children: [trashEntry] },
          setup: ({ filesList_delete }, { space }) => {
            spaceUnderTest = space
            filesList_delete([deletedFile])
          }
        })
        await getWorkerCallbackDone()

        const { showMessage } = useMessages()
        const call = (showMessage as MockedFunction<typeof showMessage>).mock.calls[0][0]
        await call.actions[0].onClick()

        const { startWorker } = vi.mocked(useRestoreWorker)()
        expect(startWorker).toHaveBeenCalledWith(
          expect.objectContaining({ space: spaceUnderTest, resources: [trashEntry] }),
          expect.any(Function)
        )
      })
    })
  })
})

function getWrapper({
  currentFolder,
  setup,
  result = [],
  getFileInfoResult,
  listFilesResult = { resource: mock<Resource>(), children: [] },
  spaceDriveType = 'personal',
  canRestoreFromTrashbin = true,
  currentRouteName = undefined,
  getDriveResult = undefined,
  stubActions = true
}: {
  currentFolder: FolderResource
  setup: (
    instance: ReturnType<typeof useFileActionsDeleteResources>,
    {
      space,
      router
    }: {
      space: SpaceResource
      router: ReturnType<typeof defaultComponentMocks>['$router']
    }
  ) => void
  result?: Resource[]
  getFileInfoResult?: Resource
  listFilesResult?: { resource: Resource; children: Resource[] }
  spaceDriveType?: string
  canRestoreFromTrashbin?: boolean
  currentRouteName?: string
  getDriveResult?: SpaceResource
  stubActions?: boolean
}) {
  const mocks = {
    ...defaultComponentMocks({
      ...(currentRouteName && { currentRoute: mock<RouteLocation>({ name: currentRouteName }) })
    }),
    space: mockDeep<SpaceResource>({
      id: 'personal',
      driveType: spaceDriveType,
      canRestoreFromTrashbin: () => canRestoreFromTrashbin
    })
  }
  mocks.$clientService.webdav.deleteFile.mockResolvedValue(undefined)
  mocks.$clientService.graphAuthenticated.drives.getDrive.mockResolvedValue(getDriveResult)
  mocks.$clientService.webdav.getFileInfo.mockResolvedValue(getFileInfoResult)
  mocks.$clientService.webdav.listFiles.mockImplementation((_, __, opts) => {
    if (opts?.isTrash) {
      return Promise.resolve(listFilesResult)
    }
    return Promise.resolve({ resource: mock<Resource>(), children: [] })
  })

  let workerCallbackDone: Promise<unknown> = Promise.resolve()
  vi.mocked(useDeleteWorker).mockReturnValue({
    startWorker: vi.fn().mockImplementation((_, callback) => {
      workerCallbackDone = Promise.resolve(callback({ successful: result, failed: [] }))
    })
  })

  vi.mocked(useRestoreWorker).mockReturnValue({
    startWorker: vi.fn().mockImplementation((_, callback) => {
      callback({ successful: [], failed: [] })
    })
  })

  vi.mocked(useGetMatchingSpace).mockImplementation(() =>
    useGetMatchingSpaceMock({
      getInternalSpace: () => mocks.space,
      getMatchingSpace: () => mocks.space
    })
  )

  return {
    mocks,
    getWorkerCallbackDone: () => workerCallbackDone,
    wrapper: getComposableWrapper(
      () => {
        const instance = useFileActionsDeleteResources()
        setup(instance, { space: mocks.space, router: mocks.$router })
      },
      {
        mocks,
        provide: mocks,
        pluginOptions: { piniaOptions: { stubActions, resourcesStore: { currentFolder } } }
      }
    )
  }
}
