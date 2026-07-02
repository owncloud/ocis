import { useFileActionsRestore } from '../../../../../src/composables/actions'
import { mock } from 'vitest-mock-extended'
import {
  defaultComponentMocks,
  getComposableWrapper,
  RouteLocation
} from '@ownclouders/web-test-helpers'
import { useMessages, useResourcesStore } from '../../../../../src/composables/piniaStores'
import { unref } from 'vue'
import { HttpError, Resource, TrashResource } from '@ownclouders/web-client'
import { ProjectSpaceResource, SpaceResource } from '@ownclouders/web-client'
import { useRestoreWorker } from '../../../../../src/composables/webWorkers/restoreWorker'

vi.mock('../../../../../src/composables/webWorkers/restoreWorker')

describe('restore', () => {
  describe('isVisible property', () => {
    it('should be false when no resource is given', () => {
      getWrapper({
        setup: ({ actions }, { space }) => {
          expect(unref(actions)[0].isVisible({ space, resources: [] })).toBe(false)
        }
      })
    })
    it('should be true when permission is sufficient', () => {
      getWrapper({
        setup: ({ actions }, { space }) => {
          expect(
            unref(actions)[0].isVisible({
              space,
              resources: [{ canBeRestored: () => true, ddate: '2020-01-01' }] as TrashResource[]
            })
          ).toBe(true)
        }
      })
    })
    it('should be false when permission is not sufficient', () => {
      getWrapper({
        setup: ({ actions }, { space }) => {
          expect(
            unref(actions)[0].isVisible({
              space,
              resources: [{ canBeRestored: () => false }] as TrashResource[]
            })
          ).toBe(false)
        }
      })
    })
    it('should be false when location is invalid', () => {
      getWrapper({
        invalidLocation: true,
        setup: ({ actions }, { space }) => {
          expect(unref(actions)[0].isVisible({ space, resources: [{}] as TrashResource[] })).toBe(
            false
          )
        }
      })
    })
    it('should be false in a space trash bin with insufficient permissions', () => {
      getWrapper({
        driveType: 'project',
        setup: ({ actions }, { space }) => {
          expect(
            unref(actions)[0].isVisible({
              space,
              resources: [{ canBeRestored: () => true }] as TrashResource[]
            })
          ).toBe(false)
        }
      })
    })
  })

  describe('method "restoreResources"', () => {
    it('should show message on success', () => {
      const resourcesToRestore = [{ id: '1', path: '/1' }] as TrashResource[]

      getWrapper({
        setup: ({ restoreResources }, { space }) => {
          restoreResources(space, resourcesToRestore, [])

          const { showMessage } = useMessages()
          expect(showMessage).toHaveBeenCalledTimes(1)

          const { removeResources, resetSelection } = useResourcesStore()
          expect(removeResources).toHaveBeenCalledTimes(1)
          expect(resetSelection).toHaveBeenCalledTimes(1)
        },
        restoreResult: { successful: resourcesToRestore, failed: [] }
      })
    })

    it('should show message on error', () => {
      vi.spyOn(console, 'error').mockImplementation(() => undefined)
      const resourcesToRestore = [{ id: '1', path: '/1' }] as TrashResource[]

      getWrapper({
        setup: ({ restoreResources }, { space }) => {
          restoreResources(space, resourcesToRestore, [])

          const { showErrorMessage } = useMessages()
          expect(showErrorMessage).toHaveBeenCalledTimes(1)

          const { removeResources } = useResourcesStore()
          expect(removeResources).toHaveBeenCalledTimes(0)
        },
        restoreResult: {
          successful: [],
          failed: [{ resource: resourcesToRestore[0], error: new HttpError('', undefined) }]
        }
      })
    })
  })

  it('should request parent folder on collecting restore conflicts', () => {
    getWrapper({
      setup: async ({ collectConflicts }, { space, clientService }) => {
        const resource = { id: '1', path: '1', name: '1' } as TrashResource
        await collectConflicts(space, [resource])

        expect(clientService.webdav.listFiles).toHaveBeenCalledWith(expect.anything(), {
          path: '.'
        })
      }
    })
  })

  it('should find conflict within resources', () => {
    getWrapper({
      setup: async ({ collectConflicts }, { space }) => {
        const resourceOne = { id: '1', path: '1', name: '1' } as TrashResource
        const resourceTwo = { id: '2', path: '1', name: '1' } as TrashResource
        const { conflicts } = await collectConflicts(space, [resourceOne, resourceTwo])

        expect(conflicts).toContain(resourceTwo)
      }
    })
  })

  it('should add files without conflict to resolved resources', () => {
    getWrapper({
      setup: async ({ collectConflicts }, { space }) => {
        const resource = { id: '1', path: '1', name: '1' } as TrashResource
        const { resolvedResources } = await collectConflicts(space, [resource])

        expect(resolvedResources).toContain(resource)
      }
    })
  })
})

function getWrapper({
  invalidLocation = false,
  driveType = 'personal',
  restoreResult = { successful: [], failed: [] },
  setup
}: {
  invalidLocation?: boolean
  driveType?: string
  restoreResult?: {
    successful: TrashResource[]
    failed: { resource: TrashResource; error: HttpError }[]
  }
  setup: (
    instance: ReturnType<typeof useFileActionsRestore>,
    {
      space
    }: {
      space: SpaceResource
      router: ReturnType<typeof defaultComponentMocks>['$router']
      clientService: ReturnType<typeof defaultComponentMocks>['$clientService']
    }
  ) => void
}) {
  vi.mocked(useRestoreWorker).mockReturnValue({
    startWorker: vi.fn().mockImplementation((_, callback) => {
      callback(restoreResult)
    })
  })

  const mocks = {
    ...defaultComponentMocks({
      currentRoute: mock<RouteLocation>({
        name: invalidLocation ? 'files-spaces-generic' : 'files-trash-generic'
      })
    }),
    space: mock<ProjectSpaceResource>({ driveType })
  }
  mocks.$clientService.webdav.listFiles.mockImplementation(() => {
    return Promise.resolve({ resource: mock<Resource>(), children: [] })
  })
  mocks.$clientService.graphAuthenticated.drives.getDrive.mockResolvedValue(mock<SpaceResource>())

  return {
    mocks,
    wrapper: getComposableWrapper(
      () => {
        const instance = useFileActionsRestore()
        setup(instance, {
          space: mocks.space,
          router: mocks.$router,
          clientService: mocks.$clientService
        })
      },
      {
        mocks,
        provide: mocks
      }
    )
  }
}
