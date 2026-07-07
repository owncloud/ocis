import { getComposableWrapper } from '@ownclouders/web-test-helpers'
import { mock } from 'vitest-mock-extended'
import {
  type WebWorker,
  type WebWorkersStore,
  useRestoreWorker,
  useWebWorkersStore
} from '../../../../../src/composables'
import type { ClientService, LoadingService } from '../../../../../src/services'
import type { Resource, SpaceResource } from '@ownclouders/web-client'

describe('useRestoreWorker', () => {
  describe('method "startWorker"', () => {
    it('creates a worker instance and starts loading task', () => {
      getWrapper({
        setup: ({ startWorker }, { webWorkersStore, mocks }) => {
          startWorker(
            { space: mock<SpaceResource>(), resources: [mock<Resource>()], missingFolderPaths: [] },
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
    instance: ReturnType<typeof useRestoreWorker>,
    {
      webWorkersStore
    }: { webWorkersStore: WebWorkersStore; workerMock: WebWorker; mocks: Record<string, unknown> }
  ) => void
}) {
  const mocks = { $loadingService: mock<LoadingService>(), $clientService: mock<ClientService>() }
  return {
    wrapper: getComposableWrapper(
      () => {
        const instance = useRestoreWorker()

        const webWorkersStore = useWebWorkersStore()
        const workerMock = mock<WebWorker>()
        vi.mocked(webWorkersStore.createWorker).mockReturnValue(workerMock)

        setup(instance, { webWorkersStore, workerMock, mocks })
      },
      { mocks, provide: mocks }
    )
  }
}
