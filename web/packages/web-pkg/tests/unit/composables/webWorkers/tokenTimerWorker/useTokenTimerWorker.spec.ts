import { unref } from 'vue'
import { getComposableWrapper } from '@ownclouders/web-test-helpers'
import { mock } from 'vitest-mock-extended'
import {
  AuthServiceInterface,
  WebWorker,
  WebWorkersStore,
  useTokenTimerWorker,
  useWebWorkersStore
} from '../../../../../src/composables'

describe('useTokenTimerWorker', () => {
  describe('method "startWorker"', () => {
    it('creates a worker instance', () => {
      getWrapper({
        setup: ({ startWorker }, { webWorkersStore }) => {
          startWorker()
          expect(vi.mocked(webWorkersStore.createWorker)).toHaveBeenCalled()
        }
      })
    })
  })

  describe('method "setTokenTimer"', () => {
    it('posts an event to the worker when started', () => {
      getWrapper({
        setup: ({ startWorker, setTokenTimer }, { workerMock }) => {
          startWorker()
          const expiryOptions = { expiry: 10, expiryThreshold: 1 }
          setTokenTimer(expiryOptions)

          expect(unref(workerMock).post).toHaveBeenCalledWith(
            JSON.stringify({ topic: 'set', ...expiryOptions })
          )
        }
      })
    })
    it('does not post an event to the worker when not started', () => {
      const consoleSpy = vi.fn()
      vi.spyOn(console, 'error').mockImplementation(consoleSpy)

      getWrapper({
        setup: ({ setTokenTimer }) => {
          const expiryOptions = { expiry: 10, expiryThreshold: 1 }
          setTokenTimer(expiryOptions)

          expect(consoleSpy).toHaveBeenCalled()
        }
      })
    })
  })

  describe('method "resetTokenTimer"', () => {
    it('posts an event to the worker when started', () => {
      getWrapper({
        setup: ({ startWorker, resetTokenTimer }, { workerMock }) => {
          startWorker()
          resetTokenTimer()

          expect(unref(workerMock).post).toHaveBeenCalledWith(JSON.stringify({ topic: 'reset' }))
        }
      })
    })
    it('does not post an event to the worker when not started', () => {
      const consoleSpy = vi.fn()
      vi.spyOn(console, 'error').mockImplementation(consoleSpy)

      getWrapper({
        setup: ({ resetTokenTimer }) => {
          resetTokenTimer()

          expect(consoleSpy).toHaveBeenCalled()
        }
      })
    })
  })
})

function getWrapper({
  setup
}: {
  setup: (
    instance: ReturnType<typeof useTokenTimerWorker>,
    { webWorkersStore }: { webWorkersStore: WebWorkersStore; workerMock: WebWorker }
  ) => void
}) {
  return {
    wrapper: getComposableWrapper(() => {
      const instance = useTokenTimerWorker({ authService: mock<AuthServiceInterface>() })

      const webWorkersStore = useWebWorkersStore()

      const workerMock = mock<WebWorker>()
      vi.mocked(webWorkersStore.createWorker).mockReturnValue(workerMock)

      setup(instance, { webWorkersStore, workerMock })
    })
  }
}
