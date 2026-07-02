import { unref } from 'vue'
import { useWebWorker } from '@vueuse/core'
import TokenWorker from '../../../../../src/composables/webWorkers/tokenTimerWorker/worker?worker'

// Real web workers can be slow to spin up and reply under heavy CI parallelism;
// give these tests more headroom than the 5s default to avoid flaky timeouts.
vi.setConfig({ testTimeout: 20000, hookTimeout: 20000 })

describe('token timer worker', () => {
  let worker: ReturnType<typeof useWebWorker>

  beforeEach(() => {
    worker = useWebWorker(TokenWorker as unknown as string, { type: 'module' })
  })

  afterEach(() => {
    worker.terminate()
  })

  it('resets the timer with "reset" topic', async () => {
    let workerResolve: (val: boolean) => unknown
    const promise = new Promise((resolve) => {
      workerResolve = resolve
    })

    const clearTimeoutSpy = vi.spyOn(window, 'clearTimeout')
    clearTimeoutSpy.mockImplementation(() => {
      workerResolve(true)
    })

    worker.post(JSON.stringify({ topic: 'reset' }))
    await promise

    expect(clearTimeoutSpy).toHaveBeenCalled()
  })

  it('should post a message with "set" topic', async () => {
    let workerResolve: (val: boolean) => unknown
    const promise = new Promise((resolve) => {
      workerResolve = resolve
    })

    const messageSpy = vi.fn()
    unref(worker.worker).onmessage = () => {
      workerResolve(true)
      messageSpy()
    }

    worker.post(JSON.stringify({ topic: 'set', expiry: 1, expiryThreshold: 1 }))
    await promise

    expect(messageSpy).toHaveBeenCalled()
  })
})
