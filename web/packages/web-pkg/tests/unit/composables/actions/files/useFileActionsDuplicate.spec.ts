import { Resource, SpaceResource } from '@ownclouders/web-client'
import {
  defaultComponentMocks,
  getComposableWrapper,
  RouteLocation
} from '@ownclouders/web-test-helpers'
import { useFileActionsDuplicate } from '../../../../../src/composables/actions'
import { mock } from 'vitest-mock-extended'
import { unref } from 'vue'
import { ListFilesResult } from '@ownclouders/web-client/webdav'
import { usePasteWorker } from '../../../../../src/composables/webWorkers/pasteWorker'

vi.mock('../../../../../src/composables/webWorkers/pasteWorker', () => ({
  usePasteWorker: vi.fn().mockReturnValue({
    startWorker: vi.fn()
  })
}))

vi.mock('../../../../../src/helpers/resource/conflictHandling/transfer', () => {
  const ResourceTransfer = vi.fn()
  ResourceTransfer.prototype.getTransferData = vi.fn().mockReturnValue([{}])

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
  ) => void
  routeName: string
  currentFolder: Resource
}) {
  const mocks = defaultComponentMocks({ currentRoute: mock<RouteLocation>({ name: routeName }) })

  return {
    mocks,
    wrapper: getComposableWrapper(
      () => {
        const instance = useFileActionsDuplicate()
        setup(instance, mocks)
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
    )
  }
}
