import {
  ResolveConflict,
  ResourceTransfer,
  TransferType,
  resolveFileNameDuplicate
} from '../../../../../src/helpers/resource/conflictHandling'
import { mock, mockDeep, mockReset } from 'vitest-mock-extended'
import { buildSpace, Resource, SpaceResource } from '@ownclouders/web-client'
import { ListFilesResult } from '@ownclouders/web-client/webdav'
import { Drive } from '@ownclouders/web-client/graph/generated'
import { createTestingPinia } from '@ownclouders/web-test-helpers'
import { ClientService } from '../../../../../src/services'
import { computed } from 'vue'

const clientServiceMock = mockDeep<ClientService>()
let resourcesToMove: Resource[]
let sourceSpace: SpaceResource
let targetSpace: SpaceResource
let targetFolder: Resource

describe('resourcesTransfer', () => {
  beforeEach(() => {
    createTestingPinia()
    mockReset(clientServiceMock)
    resourcesToMove = [
      {
        id: 'a',
        name: 'a',
        path: '/a',
        type: 'folder',
        spaceId: '1'
      },
      {
        id: 'b',
        name: 'b',
        path: '/b',
        spaceId: '1'
      }
    ]
    const spaceOptions = {
      id: 'c42c9504-2c19-44fd-87cc-b4fc20ecbb54'
    } as unknown as Drive
    sourceSpace = buildSpace(spaceOptions, {})
    targetSpace = buildSpace(spaceOptions, {})
    targetFolder = {
      id: 'target',
      path: 'target',
      webDavPath: '/target',
      spaceId: '1'
    }
  })
  it.each([
    { name: 'a', extension: '', expectName: 'a (1)' },
    { name: 'a', extension: '', expectName: 'a (2)', existing: [{ name: 'a (1)' }] },
    { name: 'a (1)', extension: '', expectName: 'a (1) (1)' },
    { name: 'b.png', extension: 'png', expectName: 'b (1).png' },
    { name: 'b.png', extension: 'png', expectName: 'b (2).png', existing: [{ name: 'b (1).png' }] }
  ])('should name duplicate file correctly', (dataSet) => {
    const existing = dataSet.existing ? [...resourcesToMove, ...dataSet.existing] : resourcesToMove
    const result = resolveFileNameDuplicate(dataSet.name, dataSet.extension, existing as Resource[])
    expect(result).toEqual(dataSet.expectName)
  })

  it('should prevent recursive paste', async () => {
    const resourcesTransfer = new ResourceTransfer(
      sourceSpace,
      resourcesToMove,
      targetSpace,
      resourcesToMove[0],
      computed(() => mock<Resource>()),
      clientServiceMock,
      vi.fn(),
      vi.fn()
    )
    const result = await resourcesTransfer.getTransferData(TransferType.COPY)
    expect(result.length).toBe(0)
  })

  describe('copyMoveResource without conflicts', () => {
    it.each([TransferType.COPY, TransferType.MOVE])(
      'should copy / move files without renaming them if no conflicts exist',
      async (action: TransferType) => {
        const listFilesResult: ListFilesResult = {
          resource: {} as Resource,
          children: []
        }
        clientServiceMock.webdav.listFiles.mockReturnValueOnce(
          new Promise((resolve) => resolve(listFilesResult))
        )
        const resourcesTransfer = new ResourceTransfer(
          sourceSpace,
          resourcesToMove,
          targetSpace,
          targetFolder,
          computed(() => mock<Resource>()),
          clientServiceMock,
          vi.fn(),
          vi.fn()
        )
        const transferData = await resourcesTransfer.getTransferData(action)

        expect(transferData.length).toBe(resourcesToMove.length)

        for (let i = 0; i < resourcesToMove.length; i++) {
          const input = resourcesToMove[i]
          const output = transferData[i]
          expect(input.name).toBe(output.resource.name)
        }
      }
    )
  })
  it('should show message if conflict exists', async () => {
    const targetFolderItems = [
      {
        id: 'a',
        path: 'target/a',
        webDavPath: '/target/a',
        name: '/target/a',
        spaceId: '1'
      }
    ]
    const resourcesTransfer = new ResourceTransfer(
      sourceSpace,
      resourcesToMove,
      targetSpace,
      resourcesToMove[0],
      computed(() => mock<Resource>()),
      clientServiceMock,
      vi.fn(),
      vi.fn()
    )
    resourcesTransfer.resolveFileExists = vi
      .fn()
      .mockImplementation(() => Promise.resolve({ strategy: 0 } as ResolveConflict))
    await resourcesTransfer.resolveAllConflicts(resourcesToMove, targetFolder, targetFolderItems)

    expect(resourcesTransfer.resolveFileExists).toHaveBeenCalled()
  })
  it('should show error message if trying to overwrite parent', async () => {
    const targetFolderItems = [
      {
        id: 'a',
        path: 'target/a',
        webDavPath: '/target/a',
        name: '/target/a',
        spaceId: '1'
      }
    ]
    const resourcesTransfer = new ResourceTransfer(
      sourceSpace,
      resourcesToMove,
      targetSpace,
      resourcesToMove[0],
      computed(() => mock<Resource>()),
      clientServiceMock,
      vi.fn(),
      vi.fn()
    )
    const namingClash = await resourcesTransfer.isOverwritingParentFolder(
      resourcesToMove[0],
      targetFolder,
      targetFolderItems
    )
    const noNamingClash = await resourcesTransfer.isOverwritingParentFolder(
      resourcesToMove[1],
      targetFolder,
      targetFolderItems
    )

    expect(namingClash).toBeTruthy()
    expect(noNamingClash).toBeFalsy()
  })
})
