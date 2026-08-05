import { mock } from 'vitest-mock-extended'
import { Resource, SpaceResource, TrashResource } from '@ownclouders/web-client'
import { defaultComponentMocks, getComposableWrapper } from '@ownclouders/web-test-helpers'
import { useUserStore } from '../../../../../src/composables/piniaStores'
import { useResolveRestorableResources } from '../../../../../src/composables/actions/helpers/useResolveRestorableResources'

describe('resolveRestorableResources', () => {
  it('returns null when the space does not allow trash-restore', () => {
    const space = mock<SpaceResource>({ driveType: 'project', canRestoreFromTrashbin: () => false })
    const deleted = [mock<Resource>({ name: 'file.txt', path: '/file.txt' })]

    getWrapper({
      setup: async ({ resolveRestorableResources }, { clientService }) => {
        const result = await resolveRestorableResources(space, deleted)
        expect(result).toBeNull()
        expect(clientService.webdav.listFiles).not.toHaveBeenCalled()
      }
    })
  })

  it('returns null when listing the trash bin fails', () => {
    const space = mock<SpaceResource>({ driveType: 'personal' })
    const deleted = [mock<Resource>({ name: 'file.txt', path: '/file.txt' })]

    getWrapper({
      setup: async ({ resolveRestorableResources }) => {
        const result = await resolveRestorableResources(space, deleted)
        expect(result).toBeNull()
      },
      listFilesImplementation: () => Promise.reject(new Error('network error'))
    })
  })

  it('returns null when not every deleted resource has a matching trash entry', () => {
    const space = mock<SpaceResource>({ driveType: 'personal' })
    const deleted = [
      mock<Resource>({ name: 'file.txt', path: '/file.txt' }),
      mock<Resource>({ name: 'other.txt', path: '/other.txt' })
    ]
    const trashEntry = mock<TrashResource>({
      name: 'file.txt',
      path: '/file.txt',
      ddate: '2026-01-01T00:00:00Z'
    })

    getWrapper({
      setup: async ({ resolveRestorableResources }) => {
        const result = await resolveRestorableResources(space, deleted)
        expect(result).toBeNull()
      },
      listFilesImplementation: () =>
        Promise.resolve({ resource: mock<Resource>(), children: [trashEntry] })
    })
  })

  it('resolves the matching trash entries for all deleted resources', () => {
    const space = mock<SpaceResource>({ driveType: 'personal' })
    const deleted = [mock<Resource>({ name: 'file.txt', path: '/file.txt' })]
    const trashEntry = mock<TrashResource>({
      name: 'file.txt',
      path: '/file.txt',
      ddate: '2026-01-01T00:00:00Z'
    })

    getWrapper({
      setup: async ({ resolveRestorableResources }) => {
        const result = await resolveRestorableResources(space, deleted)
        expect(result).toEqual([trashEntry])
      },
      listFilesImplementation: () =>
        Promise.resolve({ resource: mock<Resource>(), children: [trashEntry] })
    })
  })

  it('picks the most recently deleted match when there are duplicates', () => {
    const space = mock<SpaceResource>({ driveType: 'personal' })
    const deleted = [mock<Resource>({ name: 'file.txt', path: '/file.txt' })]
    const olderEntry = mock<TrashResource>({
      name: 'file.txt',
      path: '/file.txt',
      ddate: '2025-01-01T00:00:00Z'
    })
    const newerEntry = mock<TrashResource>({
      name: 'file.txt',
      path: '/file.txt',
      ddate: '2026-01-01T00:00:00Z'
    })

    getWrapper({
      setup: async ({ resolveRestorableResources }) => {
        const result = await resolveRestorableResources(space, deleted)
        expect(result).toEqual([newerEntry])
      },
      listFilesImplementation: () =>
        Promise.resolve({ resource: mock<Resource>(), children: [olderEntry, newerEntry] })
    })
  })
})

function getWrapper({
  setup,
  listFilesImplementation,
  user
}: {
  setup: (
    instance: ReturnType<typeof useResolveRestorableResources>,
    { clientService }: { clientService: ReturnType<typeof defaultComponentMocks>['$clientService'] }
  ) => void
  listFilesImplementation?: () => Promise<any>
  user?: any
}) {
  const mocks = defaultComponentMocks()
  mocks.$clientService.webdav.listFiles.mockImplementation(
    listFilesImplementation || (() => Promise.resolve({ resource: mock<Resource>(), children: [] }))
  )

  return {
    mocks,
    wrapper: getComposableWrapper(
      () => {
        if (user) {
          useUserStore().setUser(user)
        }
        const instance = useResolveRestorableResources()
        setup(instance, { clientService: mocks.$clientService })
      },
      {
        mocks,
        provide: mocks
      }
    )
  }
}
