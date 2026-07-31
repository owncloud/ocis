import { mock } from 'vitest-mock-extended'
import { unref } from 'vue'
import { useFileActionsDownloadArchive } from '../../../../../src/composables/actions'
import { Resource, SpaceResource } from '@ownclouders/web-client'
import {
  defaultComponentMocks,
  RouteLocation,
  getComposableWrapper
} from '@ownclouders/web-test-helpers'
import { useArchiverService } from '../../../../../src/composables'
import { ArchiverService } from '../../../../../src'

vi.mock('../../../../../src/composables/archiverService/useArchiverService')

describe('downloadArchive', () => {
  describe('search context', () => {
    describe('computed property "actions"', () => {
      describe('handler', () => {
        it.each([
          {
            resources: [
              { id: '1', fileId: '1', canDownload: () => true },
              { id: '2', fileId: '2', canDownload: () => true }
            ] as Resource[],
            downloadableResourceIds: ['1', '2']
          },
          {
            resources: [
              { id: '1', fileId: '1', canDownload: () => true },
              { id: '2', fileId: '2', canDownload: () => true },
              { id: '3', fileId: '3', canDownload: () => true },
              { id: '4', fileId: '4', canDownload: () => false },
              { id: '5', fileId: '5', canDownload: () => true, driveType: 'project' }
            ] as Resource[],
            downloadableResourceIds: ['1', '2', '3']
          }
        ])('should filter non downloadable resources', ({ resources, downloadableResourceIds }) => {
          const triggerDownloadMock = vi.fn().mockResolvedValue(true)
          getWrapper({
            searchLocation: true,
            triggerDownloadMock,
            setup: () => {
              const { actions } = useFileActionsDownloadArchive()

              unref(actions)[0].handler({ space: null, resources })

              expect(triggerDownloadMock).toHaveBeenCalledWith({ fileIds: downloadableResourceIds })
            }
          })
        })
      })
    })
  })
})

function getWrapper({
  searchLocation = false,
  triggerDownloadMock = vi.fn() as (...args: unknown[]) => unknown,
  setup = () => undefined
} = {}) {
  const routeName = searchLocation ? 'files-common-search' : 'files-spaces-generic'

  vi.mocked(useArchiverService).mockImplementation(() => {
    return {
      triggerDownload: triggerDownloadMock,
      fileIdsSupported: true
    } as ArchiverService
  })

  const mocks = {
    ...defaultComponentMocks({ currentRoute: mock<RouteLocation>({ name: routeName }) }),
    space: {
      driveType: 'personal'
    } as unknown as SpaceResource
  }

  return {
    wrapper: getComposableWrapper(setup, {
      mocks,
      provide: mocks
    })
  }
}
