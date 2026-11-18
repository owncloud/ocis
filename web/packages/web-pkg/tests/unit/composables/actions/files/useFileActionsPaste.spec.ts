// In tests, importing the paste action first runs into an import cycle (helpers index -> transfer ->
// composables index -> paste action), which leaves `ResourceTransfer` undefined in the action.
// Loading the transfer helpers first avoids that.
import { ResourceTransfer } from '../../../../../src/helpers/resource/conflictHandling/transfer'
import { Resource, SpaceResource } from '@ownclouders/web-client'
import {
  defaultComponentMocks,
  getComposableWrapper,
  RouteLocation
} from '../../../../../src/testing'
import { useFileActionsPaste } from '../../../../../src/composables/actions'
import { mock } from 'vitest-mock-extended'
import { unref } from 'vue'
import { ListFilesResult } from '@ownclouders/web-client/webdav'
import { usePasteWorker } from '../../../../../src/composables/webWorkers/pasteWorker'
import { useResourcesStore } from '../../../../../src/composables/piniaStores'
import { ClipboardActions } from '../../../../../src/helpers/clipboardActions'
import {
  TransferType,
  type TransferData,
  type TransferFailure
} from '../../../../../src/helpers/resource/conflictHandling/types'

vi.mock('../../../../../src/composables/webWorkers/pasteWorker', () => ({
  usePasteWorker: vi.fn().mockReturnValue({
    startWorker: vi.fn()
  })
}))

const mockSpace = mock<SpaceResource>({ id: 'space-1', driveType: 'personal' })

vi.mock('../../../../../src/composables/spaces/useGetMatchingSpace', () => ({
  useGetMatchingSpace: () => ({ getMatchingSpace: () => mockSpace })
}))

// what the paste worker reports back, see `usePasteWorker`
const runWorkerCallback = async (result: { successful: Resource[]; failed: TransferFailure[] }) => {
  const { startWorker } = usePasteWorker()
  expect(startWorker).toHaveBeenCalled()
  const [, callback] = vi.mocked(startWorker).mock.calls.at(-1)
  await callback(result)
}

describe('paste', () => {
  describe('handler', () => {
    afterEach(() => {
      vi.restoreAllMocks()
    })

    it('shows an incomplete folder copy in the files list right away', async () => {
      await getWrapper({
        setup: async ({ actions }, { $clientService }) => {
          const copy = mock<Resource>({ path: '/target/folder' })
          $clientService.webdav.listFiles.mockResolvedValue(mock<ListFilesResult>({ children: [] }))
          $clientService.webdav.getFileInfo.mockResolvedValue(copy)

          await unref(actions)[0].handler({ space: mockSpace, resources: [] })
          await runWorkerCallback({
            successful: [],
            failed: [
              {
                resourceName: 'folder',
                incomplete: {
                  path: '/target/folder',
                  copiedItems: 1,
                  totalItems: 2,
                  notCopied: ['folder/b.txt']
                },
                error: new Error()
              }
            ]
          })

          expect($clientService.webdav.getFileInfo).toHaveBeenCalledWith(mockSpace, {
            path: '/target/folder'
          })
          const { upsertResources } = useResourcesStore()
          expect(upsertResources).toHaveBeenCalledWith([copy])
        }
      }).done
    })

    it('reports a cut that became a copy as a copy', async () => {
      // cutting into another space asks to copy instead, the transfer data then holds a copy
      vi.spyOn(ResourceTransfer.prototype, 'getTransferData').mockResolvedValue([
        mock<TransferData>({ transferType: TransferType.COPY })
      ])
      const showResultMessage = vi.spyOn(ResourceTransfer.prototype, 'showResultMessage')

      await getWrapper({
        clipboardAction: ClipboardActions.Cut,
        setup: async ({ actions }) => {
          await unref(actions)[0].handler({ space: mockSpace, resources: [] })
          await runWorkerCallback({ successful: [], failed: [] })

          expect(showResultMessage).toHaveBeenCalledWith([], [], TransferType.COPY)
        }
      }).done
    })

    it('leaves the files list alone when nothing was created', async () => {
      await getWrapper({
        setup: async ({ actions }, { $clientService }) => {
          $clientService.webdav.listFiles.mockResolvedValue(mock<ListFilesResult>({ children: [] }))

          await unref(actions)[0].handler({ space: mockSpace, resources: [] })
          await runWorkerCallback({
            successful: [],
            failed: [{ resourceName: 'folder', error: new Error() }]
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
  clipboardAction = ClipboardActions.Copy
}: {
  clipboardAction?: string
  setup: (
    instance: ReturnType<typeof useFileActionsPaste>,
    mocks: ReturnType<typeof defaultComponentMocks>
  ) => void | Promise<void>
}) {
  const mocks = defaultComponentMocks({
    currentRoute: mock<RouteLocation>({ name: 'files-spaces-generic' })
  })
  const currentFolder = mock<Resource>({ id: 'target-id', path: '/target', type: 'folder' })
  const clipboardResource = mock<Resource>({
    id: 'folder-id',
    name: 'folder',
    extension: '',
    path: '/folder',
    isFolder: true,
    parentFolderId: 'other-folder-id'
  })

  // the composable wrapper runs `setup` on mount without waiting for it, `done` lets a test wait
  let done: void | Promise<void>
  return {
    mocks,
    wrapper: getComposableWrapper(
      () => {
        const instance = useFileActionsPaste()
        done = setup(instance, mocks)
      },
      {
        mocks,
        provide: mocks,
        pluginOptions: {
          piniaOptions: {
            clipboardState: { action: clipboardAction, resources: [clipboardResource] },
            resourcesStore: { currentFolder },
            spacesState: { spaces: [mockSpace] }
          }
        }
      }
    ),
    done
  }
}
