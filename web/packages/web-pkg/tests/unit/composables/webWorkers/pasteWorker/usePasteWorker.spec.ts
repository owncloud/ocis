import { getComposableWrapper } from '@ownclouders/web-test-helpers'
import { mock } from 'vitest-mock-extended'
import {
  type WebWorker,
  type WebWorkersStore,
  usePasteWorker,
  useWebWorkersStore
} from '../../../../../src/composables'
import type { TransferData } from '../../../../../src/helpers/resource/conflictHandling'
import type { ClientService, LoadingService } from '../../../../../src/services'

describe('usePasteWorker', () => {
  describe('method "startWorker"', () => {
    it('creates a worker instance and starts loading task', () => {
      getWrapper({
        setup: ({ startWorker }, { webWorkersStore, mocks }) => {
          startWorker([mock<TransferData>()], () => undefined)
          expect(vi.mocked(webWorkersStore.createWorker)).toHaveBeenCalled()
          expect((mocks.$loadingService as LoadingService).addTask).toHaveBeenCalledTimes(1)
        }
      })
    })
  })
})

function getWrapper({
  setup
}: {
  setup: (
    instance: ReturnType<typeof usePasteWorker>,
    {
      webWorkersStore
    }: { webWorkersStore: WebWorkersStore; workerMock: WebWorker; mocks: Record<string, unknown> }
  ) => void
}) {
  const mocks = { $loadingService: mock<LoadingService>(), $clientService: mock<ClientService>() }
  return {
    wrapper: getComposableWrapper(
      () => {
        const instance = usePasteWorker()

        const webWorkersStore = useWebWorkersStore()
        const workerMock = mock<WebWorker>()
        vi.mocked(webWorkersStore.createWorker).mockReturnValue(workerMock)

        setup(instance, { webWorkersStore, workerMock, mocks })
      },
      { mocks, provide: mocks }
    )
  }
}
