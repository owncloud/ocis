import { unref } from 'vue'
import { useWebWorker } from '@vueuse/core'
import DeleteWorker from '../../../../../src/composables/webWorkers/deleteWorker/worker?worker'
import { mock } from 'vitest-mock-extended'
import { type WebDAV } from '@ownclouders/web-client/webdav'

let webDavMock: ReturnType<typeof mock<WebDAV>>

vi.mock('@ownclouders/web-client', async (importOriginal) => ({
  ...(await importOriginal<any>()),
  webdav: () => webDavMock
}))

const resourceMock = {
  id: 'resourceId',
  name: 'resourceName',
  path: '/'
}

const spaceMock = {
  id: 'sourceSpaceId'
}

const dataMock = {
  resources: [resourceMock],
  space: spaceMock,
  concurrentRequests: 4,
  baseUrl: 'https://example.com'
}

// Real web workers can be slow to spin up and reply under heavy CI parallelism;
// give these tests more headroom than the 5s default to avoid flaky timeouts.
vi.setConfig({ testTimeout: 20000, hookTimeout: 20000 })

describe('delete worker', () => {
  let worker: ReturnType<typeof useWebWorker>

  let resolveTest: (value: boolean) => unknown
  let workerPromise: Promise<unknown>

  beforeEach(() => {
    worker = useWebWorker(DeleteWorker as unknown as string, { type: 'module' })
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

  it('calls webdav delete file operation for file list delete', async () => {
    webDavMock.deleteFile.mockResolvedValue({ status: 204, result: undefined, body: undefined })

    unref(worker.worker).onmessage = (e: MessageEvent) => {
      const { successful } = JSON.parse(e.data)
      expect(successful.length).toBe(1)
      expect(webDavMock.deleteFile).toHaveBeenCalledTimes(1)

      resolveTest(true)
    }

    worker.post(
      JSON.stringify({
        topic: 'fileListDelete',
        data: dataMock
      })
    )

    await workerPromise
  })

  it('calls webdav clear trashbin operation for trash bin delete', async () => {
    webDavMock.clearTrashBin.mockResolvedValue({ status: 204, result: undefined, body: undefined })

    unref(worker.worker).onmessage = (e: MessageEvent) => {
      const { successful } = JSON.parse(e.data)
      expect(successful.length).toBe(1)
      expect(webDavMock.clearTrashBin).toHaveBeenCalledTimes(1)

      resolveTest(true)
    }

    worker.post(
      JSON.stringify({
        topic: 'trashBinDelete',
        data: dataMock
      })
    )

    await workerPromise
  })

  it('retries on 423 status code', async () => {
    webDavMock.deleteFile.mockResolvedValueOnce({ status: 423, result: undefined, body: undefined })
    webDavMock.deleteFile.mockResolvedValueOnce({ status: 204, result: undefined, body: undefined })

    unref(worker.worker).onmessage = () => {
      expect(webDavMock.deleteFile).toHaveBeenCalledTimes(2)

      resolveTest(true)
    }

    worker.post(
      JSON.stringify({
        topic: 'fileListDelete',
        data: dataMock
      })
    )

    await workerPromise
  })

  it('returns failed files', async () => {
    vi.spyOn(console, 'error').mockImplementation(() => undefined)
    webDavMock.deleteFile.mockRejectedValue({ response: {} })

    unref(worker.worker).onmessage = (e: MessageEvent) => {
      const { failed } = JSON.parse(e.data)
      expect(failed.length).toBe(1)

      resolveTest(true)
    }

    worker.post(
      JSON.stringify({
        topic: 'fileListDelete',
        data: dataMock
      })
    )

    await workerPromise
  })
})
