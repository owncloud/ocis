import { getComposableWrapper } from '@ownclouders/web-test-helpers'
import { mock } from 'vitest-mock-extended'
import {
  type WebWorker,
  type WebWorkersStore,
  useDeleteWorker,
  useWebWorkersStore
} from '../../../../../src/composables'
import type { ClientService, LoadingService } from '../../../../../src/services'
import type { Resource, SpaceResource } from '@ownclouders/web-client'

describe('useDeleteWorker', () => {
  describe('method "startWorker"', () => {
    it('creates a worker instance and starts loading task', () => {
      getWrapper({
        setup: ({ startWorker }, { webWorkersStore, mocks }) => {
          startWorker(
            {
              topic: 'fileListDelete',
              space: mock<SpaceResource>(),
              resources: [mock<Resource>()]
            },
            () => undefined
          )
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
    instance: ReturnType<typeof useDeleteWorker>,
    {
      webWorkersStore
    }: { webWorkersStore: WebWorkersStore; workerMock: WebWorker; mocks: Record<string, unknown> }
  ) => void
}) {
  const mocks = { $loadingService: mock<LoadingService>(), $clientService: mock<ClientService>() }
  return {
    wrapper: getComposableWrapper(
      () => {
        const instance = useDeleteWorker()

        const webWorkersStore = useWebWorkersStore()
        const workerMock = mock<WebWorker>()
        vi.mocked(webWorkersStore.createWorker).mockReturnValue(workerMock)

        setup(instance, { webWorkersStore, workerMock, mocks })
      },
      { mocks, provide: mocks }
    )
  }
}
