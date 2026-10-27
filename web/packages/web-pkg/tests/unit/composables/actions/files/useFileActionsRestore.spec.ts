import { useFileActionsRestore } from '../../../../../src/composables/actions'
import { mock } from 'vitest-mock-extended'
import {
  defaultComponentMocks,
  getComposableWrapper,
  RouteLocation
} from '../../../../../src/testing'
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

    it('should upsert restored resources into the current view when not in the trash bin', async () => {
      const resourcesToRestore = [{ id: '1', path: '/1', name: 'file1' }] as TrashResource[]
      const restoredResource = mock<Resource>({ id: '1', path: '/1', name: 'file1' })

      const { getWorkerCallbackDone } = getWrapper({
        invalidLocation: true,
        setup: ({ restoreResources }, { space }) => {
          restoreResources(space, resourcesToRestore, [])
        },
        restoreResult: { successful: resourcesToRestore, failed: [] },
        getFileInfoResult: restoredResource
      })
      await getWorkerCallbackDone()

      const { removeResources, upsertResources } = useResourcesStore()
      expect(removeResources).toHaveBeenCalledTimes(0)
      expect(upsertResources).toHaveBeenCalledWith([restoredResource])
    })

    it('should upsert the resolvable resources even when one lookup fails', async () => {
      const resourcesToRestore = [
        { id: '1', path: '/1', name: 'file1' },
        { id: '2', path: '/2', name: 'file2' }
      ] as TrashResource[]
      const restoredResource = mock<Resource>({ id: '2', path: '/2', name: 'file2' })

      const { getWorkerCallbackDone } = getWrapper({
        invalidLocation: true,
        setup: ({ restoreResources }, { space }) => {
          restoreResources(space, resourcesToRestore, [])
        },
        restoreResult: { successful: resourcesToRestore, failed: [] },
        getFileInfoImplementation: (_, { path }) =>
          path === '/2'
            ? Promise.resolve(restoredResource)
            : Promise.reject(new Error('network error'))
      })
      await getWorkerCallbackDone()

      const { upsertResources } = useResourcesStore()
      expect(upsertResources).toHaveBeenCalledWith([restoredResource])
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

  it('should resolve a conflict within resources with a unique name', async () => {
    const resourceOne = { id: '1', path: '/file.txt', name: 'file.txt' } as TrashResource
    const resourceTwo = { id: '2', path: '/file.txt', name: 'file.txt' } as TrashResource
    let result: ReturnType<ReturnType<typeof useFileActionsRestore>['collectConflicts']>
    getWrapper({
      setup: ({ collectConflicts }, { space }) => {
        result = collectConflicts(space, [resourceOne, resourceTwo])
      }
    })
    const { conflicts, resolvedResources } = await result

    expect(conflicts).toHaveLength(0)
    expect(resolvedResources.map(({ id, name, path }) => ({ id, name, path }))).toEqual([
      { id: '1', name: 'file.txt', path: '/file.txt' },
      { id: '2', name: 'file (1).txt', path: '/file (1).txt' }
    ])
  })

  it('should find conflict with existing resources in the parent folder', async () => {
    const resource = { id: '1', path: '/file.txt', name: 'file.txt' } as TrashResource
    let result: ReturnType<ReturnType<typeof useFileActionsRestore>['collectConflicts']>
    getWrapper({
      existingResources: [mock<Resource>({ name: 'file.txt' })],
      setup: ({ collectConflicts }, { space }) => {
        result = collectConflicts(space, [resource])
      }
    })
    const { conflicts, resolvedResources } = await result

    expect(conflicts).toEqual([resource])
    expect(resolvedResources).toHaveLength(0)
  })

  it('should not pick a name of an existing resource when resolving a conflict within resources', async () => {
    const resourceOne = { id: '1', path: '/file.txt', name: 'file.txt' } as TrashResource
    const resourceTwo = { id: '2', path: '/file.txt', name: 'file.txt' } as TrashResource
    let result: ReturnType<ReturnType<typeof useFileActionsRestore>['collectConflicts']>
    getWrapper({
      existingResources: [mock<Resource>({ name: 'file (1).txt' })],
      setup: ({ collectConflicts }, { space }) => {
        result = collectConflicts(space, [resourceOne, resourceTwo])
      }
    })
    const { resolvedResources } = await result

    expect(resolvedResources.map(({ path }) => path)).toEqual(['/file.txt', '/file (2).txt'])
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
  getFileInfoResult,
  getFileInfoImplementation,
  existingResources = [],
  setup
}: {
  existingResources?: Resource[]
  invalidLocation?: boolean
  driveType?: string
  restoreResult?: {
    successful: TrashResource[]
    failed: { resource: TrashResource; error: HttpError }[]
  }
  getFileInfoResult?: Resource
  getFileInfoImplementation?: (...args: any[]) => Promise<Resource>
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
  let workerCallbackDone: Promise<unknown> = Promise.resolve()
  vi.mocked(useRestoreWorker).mockReturnValue({
    startWorker: vi.fn().mockImplementation((_, callback) => {
      workerCallbackDone = Promise.resolve(callback(restoreResult))
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
    return Promise.resolve({ resource: mock<Resource>(), children: existingResources })
  })
  if (getFileInfoImplementation) {
    mocks.$clientService.webdav.getFileInfo.mockImplementation(getFileInfoImplementation)
  } else {
    mocks.$clientService.webdav.getFileInfo.mockResolvedValue(getFileInfoResult)
  }
  mocks.$clientService.graphAuthenticated.drives.getDrive.mockResolvedValue(mock<SpaceResource>())

  return {
    mocks,
    getWorkerCallbackDone: () => workerCallbackDone,
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
