import { useFileActionsDeleteResources } from '../../../../../src/composables/actions'
import { mock, mockDeep } from 'vitest-mock-extended'
import { FolderResource, Resource, SpaceResource } from '@ownclouders/web-client'
import {
  defaultComponentMocks,
  getComposableWrapper,
  useGetMatchingSpaceMock
} from '@ownclouders/web-test-helpers'
import { useDeleteWorker } from '../../../../../src/composables/webWorkers/deleteWorker'
import { useGetMatchingSpace } from '../../../../../src/composables/spaces/useGetMatchingSpace'
import { useResourcesStore, useSpacesStore } from '../../../../../src/composables/piniaStores'
import { MockedFunction } from 'vitest'

vi.mock('../../../../../src/composables/webWorkers/deleteWorker')
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
  })
})

function getWrapper({
  currentFolder,
  setup,
  result = [],
  getFileInfoResult
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
}) {
  const mocks = {
    ...defaultComponentMocks(),
    space: mockDeep<SpaceResource>({ id: 'personal' })
  }
  mocks.$clientService.webdav.deleteFile.mockResolvedValue(undefined)
  mocks.$clientService.webdav.getFileInfo.mockResolvedValue(getFileInfoResult)

  vi.mocked(useDeleteWorker).mockReturnValue({
    startWorker: vi.fn().mockImplementation((_, callback) => {
      callback({ successful: result, failed: [] })
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
    wrapper: getComposableWrapper(
      () => {
        const instance = useFileActionsDeleteResources()
        setup(instance, { space: mocks.space, router: mocks.$router })
      },
      {
        mocks,
        provide: mocks,
        pluginOptions: { piniaOptions: { resourcesStore: { currentFolder } } }
      }
    )
  }
}
