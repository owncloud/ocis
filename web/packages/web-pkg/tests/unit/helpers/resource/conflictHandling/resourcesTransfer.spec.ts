import {
  ResolveConflict,
  ResourceTransfer,
  TransferType,
  resolveFileNameDuplicate
} from '../../../../../src/helpers/resource/conflictHandling'
import { mock, mockDeep, mockReset } from 'vitest-mock-extended'
import { buildSpace, createHttpError, Resource, SpaceResource } from '@ownclouders/web-client'
import { ListFilesResult } from '@ownclouders/web-client/webdav'
import { Drive } from '@ownclouders/web-client/graph/generated'
import { createTestingPinia } from '@ownclouders/web-test-helpers'
import { ClientService } from '../../../../../src/services'
import { useMessages } from '../../../../../src/composables'
import { computed } from 'vue'

/**
 * Unlike a plain spy, this interpolates the message, so that assertions can be made on what
 * the user actually reads.
 */
const interpolate = (msgid: string, params: Record<string, unknown> = {}) =>
  Object.entries(params).reduce(
    (message, [key, value]) => message.replace(`%{${key}}`, String(value)),
    msgid
  )

/** Picks the form English asks for, so assertions read like the English messages do. */
const interpolatePlural = (
  singular: string,
  plural: string,
  count: number,
  params: Record<string, unknown> = {}
) => interpolate(count === 1 ? singular : plural, params)

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
  describe('showResultMessage', () => {
    const buildTransfer = (space: SpaceResource) =>
      new ResourceTransfer(
        sourceSpace,
        resourcesToMove,
        space,
        targetFolder,
        computed(() => mock<Resource>()),
        clientServiceMock,
        interpolate,
        interpolatePlural
      )

    const projectSpace = () =>
      mock<SpaceResource>({ id: '1', name: 'Group Space', driveType: 'project' })

    const quotaError = () =>
      createHttpError({ message: 'Insufficient Storage', statusCode: 507, xReqId: '1' })

    it.each([
      { transferType: TransferType.DUPLICATE, title: 'Failed to duplicate "a"' },
      { transferType: TransferType.COPY, title: 'Failed to copy "a"' },
      { transferType: TransferType.MOVE, title: 'Failed to move "a"' }
    ])('names the action that failed for a single resource', ({ transferType, title }) => {
      const showErrorMessage = vi.spyOn(useMessages(), 'showErrorMessage')

      buildTransfer(projectSpace()).showResultMessage(
        [{ resourceName: 'a', error: new Error() }],
        [],
        transferType
      )

      expect(showErrorMessage).toHaveBeenCalledWith(expect.objectContaining({ title }))
    })

    it.each([
      { transferType: TransferType.DUPLICATE, title: 'Failed to duplicate 2 resources' },
      { transferType: TransferType.COPY, title: 'Failed to copy 2 resources' },
      { transferType: TransferType.MOVE, title: 'Failed to move 2 resources' }
    ])('names the action that failed for several resources', ({ transferType, title }) => {
      const showErrorMessage = vi.spyOn(useMessages(), 'showErrorMessage')

      buildTransfer(projectSpace()).showResultMessage(
        [
          { resourceName: 'a', error: new Error() },
          { resourceName: 'b', error: new Error() }
        ],
        [],
        transferType
      )

      expect(showErrorMessage).toHaveBeenCalledWith(expect.objectContaining({ title }))
    })

    it.each([
      {
        transferType: TransferType.DUPLICATE,
        desc: 'The file cannot be duplicated because there is not enough storage left in "Group Space".'
      },
      {
        transferType: TransferType.COPY,
        desc: 'The file cannot be copied because there is not enough storage left in "Group Space".'
      },
      {
        transferType: TransferType.MOVE,
        desc: 'The file cannot be moved because there is not enough storage left in "Group Space".'
      }
    ])('says why the action failed', ({ transferType, desc }) => {
      const showErrorMessage = vi.spyOn(useMessages(), 'showErrorMessage')

      buildTransfer(projectSpace()).showResultMessage(
        [{ resourceName: 'a', error: quotaError() }],
        [],
        transferType
      )

      expect(showErrorMessage).toHaveBeenCalledWith(expect.objectContaining({ desc }))
    })

    it('refers to a personal space as "Personal" rather than by its name', () => {
      const showErrorMessage = vi.spyOn(useMessages(), 'showErrorMessage')

      buildTransfer(
        mock<SpaceResource>({ id: '1', name: 'Admin', driveType: 'personal' })
      ).showResultMessage([{ resourceName: 'a', error: quotaError() }], [], TransferType.DUPLICATE)

      expect(showErrorMessage).toHaveBeenCalledWith(
        expect.objectContaining({
          desc: 'The file cannot be duplicated because there is not enough storage left in "Personal".'
        })
      )
    })

    it('leaves the description out when the failure was not about quota', () => {
      const showErrorMessage = vi.spyOn(useMessages(), 'showErrorMessage')

      buildTransfer(projectSpace()).showResultMessage(
        [
          {
            resourceName: 'a',
            error: createHttpError({ message: 'nope', statusCode: 403, xReqId: '1' })
          }
        ],
        [],
        TransferType.DUPLICATE
      )

      expect(showErrorMessage).toHaveBeenCalledWith(
        expect.not.objectContaining({ desc: expect.anything() })
      )
    })
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
