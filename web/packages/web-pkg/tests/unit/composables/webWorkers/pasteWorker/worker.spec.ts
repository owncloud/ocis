import { unref } from 'vue'
import { useWebWorker } from '@vueuse/core'
import PasteWorker from '../../../../../src/composables/webWorkers/pasteWorker/worker?worker'
import { TransferType } from '../../../../../src/helpers/resource/conflictHandling'
import { mock } from 'vitest-mock-extended'
import type { WebDAV } from '@ownclouders/web-client/webdav'
import type { Resource } from '@ownclouders/web-client'

let webDavMock: ReturnType<typeof mock<WebDAV>>

vi.mock('@ownclouders/web-client', async (importOriginal) => ({
  ...(await importOriginal<any>()),
  webdav: () => webDavMock
}))

const resourceMock = {
  id: 'resourceId',
  name: 'resourceName'
}

const targetSpaceMock = {
  id: 'targetSpaceId'
}

const sourceSpaceMock = {
  id: 'sourceSpaceId'
}

const targetFolderMock = {
  id: 'targetFolderId',
  path: '/',
  webDavPath: '/'
}

const transferDataMock = {
  resource: resourceMock,
  sourceSpace: sourceSpaceMock,
  targetSpace: targetSpaceMock,
  targetFolder: targetFolderMock,
  path: ''
}

// Real web workers can be slow to spin up and reply under heavy CI parallelism;
// give these tests more headroom than the 5s default to avoid flaky timeouts.
vi.setConfig({ testTimeout: 20000, hookTimeout: 20000 })

describe('paste worker', () => {
  let worker: ReturnType<typeof useWebWorker>

  let resolveTest: (value: boolean) => unknown
  let workerPromise: Promise<unknown>

  beforeEach(() => {
    worker = useWebWorker(PasteWorker as unknown as string, { type: 'module' })
    webDavMock = mock<WebDAV>()

    workerPromise = new Promise((resolve) => {
      resolveTest = resolve
    })
  })

  afterEach(() => {
    worker.terminate()

    workerPromise = undefined
    resolveTest = undefined
    webDavMock = undefined
  })

  it('calls webdav copy operation for copy actions', async () => {
    webDavMock.copyFiles.mockResolvedValue(undefined)

    unref(worker.worker).onmessage = (e: MessageEvent) => {
      const { successful } = JSON.parse(e.data)
      expect(successful.length).toBe(1)
      expect(webDavMock.copyFiles).toHaveBeenCalledTimes(1)

      resolveTest(true)
    }

    worker.post(
      JSON.stringify({
        topic: 'startProcess',
        data: {
          baseUrl: 'https://example.com',
          transferData: [{ ...transferDataMock, transferType: TransferType.COPY }]
        }
      })
    )

    await workerPromise
  })

  it('calls webdav move operation for move actions', async () => {
    webDavMock.moveFiles.mockResolvedValue(undefined)

    unref(worker.worker).onmessage = (e: MessageEvent) => {
      const { successful } = JSON.parse(e.data)
      expect(successful.length).toBe(1)
      expect(webDavMock.moveFiles).toHaveBeenCalledTimes(1)

      resolveTest(true)
    }

    worker.post(
      JSON.stringify({
        topic: 'startProcess',
        data: {
          baseUrl: 'https://example.com',
          transferData: [{ ...transferDataMock, transferType: TransferType.MOVE }]
        }
      })
    )

    await workerPromise
  })

  it('returns failed files', async () => {
    vi.spyOn(console, 'error').mockImplementation(() => undefined)
    webDavMock.copyFiles.mockRejectedValue({ response: {} })

    unref(worker.worker).onmessage = (e: MessageEvent) => {
      const { failed } = JSON.parse(e.data)
      expect(failed.length).toBe(1)

      resolveTest(true)
    }

    worker.post(
      JSON.stringify({
        topic: 'startProcess',
        data: {
          baseUrl: 'https://example.com',
          transferData: [{ ...transferDataMock, transferType: TransferType.COPY }]
        }
      })
    )

    await workerPromise
  })

  it('returns a copy that got no HTTP response as failed', async () => {
    vi.spyOn(console, 'error').mockImplementation(() => undefined)
    webDavMock.copyFiles.mockRejectedValue(new TypeError('Failed to fetch'))

    unref(worker.worker).onmessage = (e: MessageEvent) => {
      const { successful, failed } = JSON.parse(e.data)
      expect(successful.length).toBe(0)
      expect(failed).toEqual([{ resourceName: resourceMock.name, message: 'Failed to fetch' }])

      resolveTest(true)
    }

    worker.post(
      JSON.stringify({
        topic: 'startProcess',
        data: {
          baseUrl: 'https://example.com',
          transferData: [{ ...transferDataMock, transferType: TransferType.DUPLICATE }]
        }
      })
    )

    await workerPromise
  })

  describe('folder copy', () => {
    const folderMock = { id: 'folderId', name: 'folder', path: '/folder', isFolder: true }

    // same space for source and target, like a duplication
    const folderTransferData = (overrides: Record<string, unknown> = {}) => ({
      ...transferDataMock,
      targetSpace: sourceSpaceMock,
      resource: folderMock,
      path: '/folder (1)',
      overwrite: false,
      transferType: TransferType.DUPLICATE,
      ...overrides
    })

    // the source folder contains "a.txt" and the subfolder "sub" with "b.txt"
    const folderTree: Record<string, Partial<Resource>[]> = {
      '/folder': [
        { name: 'a.txt', path: '/folder/a.txt', isFolder: false },
        { name: 'sub', path: '/folder/sub', isFolder: true }
      ],
      '/folder/sub': [{ name: 'b.txt', path: '/folder/sub/b.txt', isFolder: false }]
    }

    const mockFolderTree = (tree = folderTree) => {
      webDavMock.createFolder.mockResolvedValue(undefined)
      webDavMock.copyFiles.mockResolvedValue(undefined)
      webDavMock.listFiles.mockImplementation((_space, { path }) =>
        Promise.resolve({
          resource: {} as Resource,
          children: (tree[path] || []) as Resource[]
        })
      )
    }

    const runWorker = (transferData: unknown[]) =>
      new Promise<{ successful: Resource[]; failed: Record<string, unknown>[] }>((resolve) => {
        unref(worker.worker).onmessage = (e: MessageEvent) => resolve(JSON.parse(e.data))
        worker.post(
          JSON.stringify({
            topic: 'startProcess',
            data: { baseUrl: 'https://example.com', transferData }
          })
        )
      })

    const copiedPaths = () =>
      webDavMock.copyFiles.mock.calls.map(([, source, , target]) => [source.path, target.path])
    const createdPaths = () => webDavMock.createFolder.mock.calls.map(([, { path }]) => path)

    it.each([TransferType.DUPLICATE, TransferType.COPY])(
      'copies a folder file by file',
      async (transferType) => {
        mockFolderTree()

        const { successful, failed } = await runWorker([folderTransferData({ transferType })])

        expect(failed).toEqual([])
        expect(successful).toEqual([expect.objectContaining({ name: 'folder' })])
        expect(createdPaths().sort()).toEqual(['/folder (1)', '/folder (1)/sub'])
        expect(copiedPaths().sort()).toEqual([
          ['/folder/a.txt', '/folder (1)/a.txt'],
          ['/folder/sub/b.txt', '/folder (1)/sub/b.txt']
        ])
      }
    )

    it.each([
      {
        failure: 'a file copy',
        mockFailure: () =>
          webDavMock.copyFiles.mockImplementation((_s, source) =>
            source.path === '/folder/a.txt'
              ? Promise.reject(new TypeError('Failed to fetch'))
              : Promise.resolve(undefined)
          )
      },
      {
        failure: 'creating a subfolder',
        mockFailure: () =>
          webDavMock.createFolder.mockImplementation((_s, { path }) =>
            path === '/folder (1)/sub'
              ? Promise.reject(new TypeError('Failed to fetch'))
              : Promise.resolve(undefined)
          )
      },
      {
        failure: 'listing a subfolder',
        mockFailure: () =>
          webDavMock.listFiles.mockImplementation((_s, { path }) =>
            path === '/folder/sub'
              ? Promise.reject(new TypeError('Failed to fetch'))
              : Promise.resolve({
                  resource: {} as Resource,
                  children: folderTree[path] as Resource[]
                })
          )
      }
    ])('returns a folder as incomplete when $failure fails', async ({ mockFailure }) => {
      vi.spyOn(console, 'error').mockImplementation(() => undefined)
      mockFolderTree()
      mockFailure()

      const { successful, failed } = await runWorker([folderTransferData()])

      expect(successful).toEqual([])
      expect(failed).toEqual([
        { resourceName: 'folder', incomplete: true, message: 'Failed to fetch' }
      ])
    })

    it('keeps copying the rest of a folder after one file failed', async () => {
      vi.spyOn(console, 'error').mockImplementation(() => undefined)
      mockFolderTree()
      webDavMock.copyFiles.mockImplementation((_s, source) =>
        source.path === '/folder/a.txt'
          ? Promise.reject(new TypeError('Failed to fetch'))
          : Promise.resolve(undefined)
      )

      await runWorker([folderTransferData()])

      expect(copiedPaths()).toContainEqual(['/folder/sub/b.txt', '/folder (1)/sub/b.txt'])
    })

    it('returns a folder as failed when the folder itself cannot be created', async () => {
      vi.spyOn(console, 'error').mockImplementation(() => undefined)
      mockFolderTree()
      webDavMock.createFolder.mockRejectedValue(new TypeError('Failed to fetch'))

      const { successful, failed } = await runWorker([folderTransferData()])

      expect(successful).toEqual([])
      expect(failed).toEqual([{ resourceName: 'folder', message: 'Failed to fetch' }])
      expect(webDavMock.createFolder).toHaveBeenCalledTimes(1)
      expect(webDavMock.copyFiles).not.toHaveBeenCalled()
    })

    it('returns a folder as failed and creates nothing when it cannot be listed', async () => {
      vi.spyOn(console, 'error').mockImplementation(() => undefined)
      mockFolderTree()
      webDavMock.listFiles.mockRejectedValue(new TypeError('Failed to fetch'))

      const { successful, failed } = await runWorker([folderTransferData()])

      expect(successful).toEqual([])
      expect(failed).toEqual([{ resourceName: 'folder', message: 'Failed to fetch' }])
      expect(webDavMock.createFolder).not.toHaveBeenCalled()
      expect(webDavMock.copyFiles).not.toHaveBeenCalled()
    })

    it('lists the whole folder before creating anything', async () => {
      mockFolderTree()

      await runWorker([folderTransferData()])

      const lastListing = Math.max(...webDavMock.listFiles.mock.invocationCallOrder)
      const firstCreation = Math.min(...webDavMock.createFolder.mock.invocationCallOrder)
      expect(webDavMock.listFiles).toHaveBeenCalledTimes(2)
      expect(lastListing).toBeLessThan(firstCreation)
    })

    it('replaces a folder with one request', async () => {
      mockFolderTree()

      await runWorker([folderTransferData({ transferType: TransferType.COPY, overwrite: true })])

      expect(webDavMock.createFolder).not.toHaveBeenCalled()
      expect(copiedPaths()).toEqual([['/folder', '/folder (1)']])
    })

    it('leaves copying a folder into one of its subfolders to the server', async () => {
      mockFolderTree()

      await runWorker([
        folderTransferData({ transferType: TransferType.COPY, path: '/folder/sub/folder' })
      ])

      expect(webDavMock.listFiles).not.toHaveBeenCalled()
      expect(webDavMock.createFolder).not.toHaveBeenCalled()
      expect(copiedPaths()).toEqual([['/folder', '/folder/sub/folder']])
    })

    it('moves a folder with one request', async () => {
      mockFolderTree()
      webDavMock.moveFiles.mockResolvedValue(undefined)

      await runWorker([folderTransferData({ transferType: TransferType.MOVE })])

      expect(webDavMock.moveFiles).toHaveBeenCalledTimes(1)
      expect(webDavMock.createFolder).not.toHaveBeenCalled()
      expect(webDavMock.copyFiles).not.toHaveBeenCalled()
    })

    it('copies more folders than requests run at a time, never more than 4 at once', async () => {
      mockFolderTree()
      let running = 0
      let mostRunning = 0
      const request = async () => {
        running++
        mostRunning = Math.max(mostRunning, running)
        await new Promise((resolve) => setTimeout(resolve, 5))
        running--
      }
      webDavMock.createFolder.mockImplementation(async () => {
        await request()
        return undefined
      })
      webDavMock.copyFiles.mockImplementation(async () => {
        await request()
        return undefined
      })
      webDavMock.listFiles.mockImplementation(async (_space, { path }) => {
        await request()
        return { resource: {} as Resource, children: folderTree[path] as Resource[] }
      })

      const folders = ['f1', 'f2', 'f3', 'f4', 'f5'].map((name) =>
        folderTransferData({
          resource: { ...folderMock, id: name, name, path: '/folder' },
          path: `/${name} (1)`
        })
      )
      const { successful, failed } = await runWorker(folders)

      expect(failed).toEqual([])
      expect(successful.length).toBe(5)
      expect(mostRunning).toBeLessThanOrEqual(4)
    })
  })
})
