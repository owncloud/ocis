import { Resource, SpaceResource } from '@ownclouders/web-client'
import {
  defaultComponentMocks,
  getComposableWrapper,
  RouteLocation
} from '../../../../../src/testing'
import { useFileActionsDuplicate } from '../../../../../src/composables/actions'
import { mock } from 'vitest-mock-extended'
import { unref } from 'vue'
import { ListFilesResult } from '@ownclouders/web-client/webdav'
import type { TransferFailure } from '../../../../../src/helpers/resource/conflictHandling/types'
import { usePasteWorker } from '../../../../../src/composables/webWorkers/pasteWorker'
import { useResourcesStore } from '../../../../../src/composables/piniaStores'

vi.mock('../../../../../src/composables/webWorkers/pasteWorker', () => ({
  usePasteWorker: vi.fn().mockReturnValue({
    startWorker: vi.fn()
  })
}))

vi.mock('../../../../../src/helpers/resource/conflictHandling/transfer', () => {
  const ResourceTransfer = vi.fn()
  ResourceTransfer.prototype.getTransferData = vi.fn().mockReturnValue([{}])
  ResourceTransfer.prototype.showResultMessage = vi.fn()

  return {
    ResourceTransfer
  }
})

const mockSpace = mock<SpaceResource>({ id: 'space-1' })

describe('duplicate', () => {
  describe('visibility', () => {
    it.each([
      {
        scenario: 'should not be visible when location is files-space-projects',
        location: 'files-space-projects',
        currentFolder: mock<Resource>({ canCreate: () => true }),
        resources: [mock<Resource>({ canDownload: () => true })],
        expectedStatus: false
      },
      {
        scenario:
          'should not be visible when location is files-public-link and create permission is not granted',
        location: 'files-public-link',
        currentFolder: mock<Resource>({ canCreate: () => false }),
        resources: [mock<Resource>({ canDownload: () => true })],
        expectedStatus: false
      },
      {
        scenario:
          'should be visible when location is files-public-link and create permission is granted',
        location: 'files-public-link',
        currentFolder: mock<Resource>({ canCreate: () => true }),
        resources: [mock<Resource>({ canDownload: () => true })],
        expectedStatus: true
      },
      {
        scenario:
          'should not be visible when location is files-common-search and all resources are project spaces',
        location: 'files-common-search',
        currentFolder: mock<Resource>({ canCreate: () => true }),
        resources: [mock<SpaceResource>({ driveType: 'project' })],
        expectedStatus: false
      },
      {
        scenario:
          'should be visible when location is files-common-search and at least one resource is not a project space',
        location: 'files-common-search',
        currentFolder: mock<Resource>({ canCreate: () => true }),
        resources: [mock<Resource>({ canDownload: () => true })],
        expectedStatus: true
      },
      {
        scenario: 'should not be visible when download permission is not granted',
        location: 'files-spaces-generic',
        currentFolder: mock<Resource>({ canCreate: () => true }),
        resources: [mock<Resource>({ canDownload: () => false })],
        expectedStatus: false
      },
      {
        scenario: 'should be visible when download permission is granted',
        location: 'files-spaces-generic',
        currentFolder: mock<Resource>({ canCreate: () => true }),
        resources: [mock<Resource>({ canDownload: () => true })],
        expectedStatus: true
      }
    ])('$scenario', ({ location, currentFolder, resources, expectedStatus }) => {
      getWrapper({
        routeName: location,
        currentFolder,
        setup: ({ actions }) => {
          expect(unref(actions)[0].isVisible({ resources, space: mock<SpaceResource>() })).toBe(
            expectedStatus
          )
        }
      })
    })
  })

  describe('handler', () => {
    it('should start paste worker', async () => {
      await getWrapper({
        routeName: 'files-spaces-generic',
        currentFolder: mock<Resource>({ path: '/', canCreate: () => true }),
        setup: async ({ actions }, { $clientService }) => {
          const resource = mock<Resource>({ storageId: mockSpace.id, path: '/' })

          $clientService.webdav.listFiles.mockResolvedValue(
            mock<ListFilesResult>({
              children: [resource]
            })
          )

          await unref(actions)[0].handler({
            space: mockSpace,
            resources: [resource]
          })

          const { startWorker } = usePasteWorker()
          expect(startWorker).toHaveBeenCalled()
        }
      })
    })

    // what the paste worker reports back, see `usePasteWorker`
    const runWorkerCallback = async (result: {
      successful: Resource[]
      failed: TransferFailure[]
    }) => {
      const { startWorker } = usePasteWorker()
      const [, callback] = vi.mocked(startWorker).mock.calls.at(-1)
      await callback(result)
    }

    it('shows an incomplete folder copy in the files list right away', async () => {
      await getWrapper({
        routeName: 'files-spaces-generic',
        currentFolder: mock<Resource>({ path: '/', canCreate: () => true }),
        setup: async ({ actions }, { $clientService }) => {
          const resource = mock<Resource>({ storageId: mockSpace.id, path: '/folder' })
          const copy = mock<Resource>({ path: '/folder (1)' })
          $clientService.webdav.getFileInfo.mockResolvedValue(copy)

          await unref(actions)[0].handler({ space: mockSpace, resources: [resource] })
          await runWorkerCallback({
            successful: [],
            failed: [
              {
                resourceName: 'folder (1)',
                incomplete: {
                  path: '/folder (1)',
                  copiedItems: 1,
                  totalItems: 2,
                  notCopied: ['folder/b.txt']
                },
                error: new Error()
              }
            ]
          })

          expect($clientService.webdav.getFileInfo).toHaveBeenCalledWith(mockSpace, {
            path: '/folder (1)'
          })
          const { upsertResources } = useResourcesStore()
          expect(upsertResources).toHaveBeenCalledWith([copy])
        }
      }).done
    })

    it('leaves the files list alone when the user opened another folder meanwhile', async () => {
      await getWrapper({
        routeName: 'files-spaces-generic',
        currentFolder: mock<Resource>({ id: 'folder-1', path: '/', canCreate: () => true }),
        setup: async ({ actions }, { $clientService }) => {
          const resource = mock<Resource>({ storageId: mockSpace.id, path: '/folder' })

          await unref(actions)[0].handler({ space: mockSpace, resources: [resource] })
          // a big copy takes a while, the user opens another folder before it is done
          const resourcesStore = useResourcesStore()
          resourcesStore.currentFolder = mock<Resource>({ id: 'folder-2', path: '/other' })
          await runWorkerCallback({
            successful: [mock<Resource>({ path: '/folder (1)' })],
            failed: []
          })

          expect($clientService.webdav.getFileInfo).not.toHaveBeenCalled()
          expect(resourcesStore.upsertResources).not.toHaveBeenCalled()
        }
      }).done
    })

    it('leaves the files list alone when nothing was created', async () => {
      await getWrapper({
        routeName: 'files-spaces-generic',
        currentFolder: mock<Resource>({ path: '/', canCreate: () => true }),
        setup: async ({ actions }, { $clientService }) => {
          const resource = mock<Resource>({ storageId: mockSpace.id, path: '/folder' })

          await unref(actions)[0].handler({ space: mockSpace, resources: [resource] })
          await runWorkerCallback({
            successful: [],
            failed: [{ resourceName: 'folder (1)', error: new Error() }]
          })

          expect($clientService.webdav.getFileInfo).not.toHaveBeenCalled()
          const { upsertResources } = useResourcesStore()
          expect(upsertResources).not.toHaveBeenCalled()
        }
      }).done
    })
  })
})

function getWrapper({
  setup,
  routeName,
  currentFolder
}: {
  setup: (
    instance: ReturnType<typeof useFileActionsDuplicate>,
    mocks: ReturnType<typeof defaultComponentMocks>
  ) => void | Promise<void>
  routeName: string
  currentFolder: Resource
}) {
  const mocks = defaultComponentMocks({ currentRoute: mock<RouteLocation>({ name: routeName }) })

  // the composable wrapper runs `setup` on mount without waiting for it, `done` lets a test wait
  let done: void | Promise<void>
  return {
    mocks,
    wrapper: getComposableWrapper(
      () => {
        const instance = useFileActionsDuplicate()
        done = setup(instance, mocks)
      },
      {
        mocks,
        provide: mocks,
        pluginOptions: {
          piniaOptions: {
            resourcesStore: { currentFolder },
            spacesState: { spaces: [mockSpace] }
          }
        }
      }
    ),
    done
  }
}
