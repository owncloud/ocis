import { getComposableWrapper, useGetMatchingSpaceMock } from '@ownclouders/web-test-helpers'
import { mock } from 'vitest-mock-extended'
import {
  getPermissionsForSpaceMember,
  GraphSharePermission,
  IncomingShareResource,
  PersonalSpaceResource,
  PublicSpaceResource,
  Resource,
  ShareSpaceResource,
  SpaceResource,
  TrashResource
} from '@ownclouders/web-client'
import { useCanListShares } from '../../../../src/composables/shares'
import { useCapabilityStore } from '../../../../src/composables/piniaStores'
import { useGetMatchingSpace } from '../../../../src/composables/spaces/useGetMatchingSpace'
import { Identity } from '@ownclouders/web-client/graph/generated'

vi.mock('../../../../src/composables/spaces/useGetMatchingSpace')
vi.mock('@ownclouders/web-client', async (importOriginal) => ({
  ...(await importOriginal<any>()),
  getPermissionsForSpaceMember: vi.fn()
}))

describe('useCanListShares', () => {
  describe('canListShares', () => {
    it('returns true with sharing enabled and sufficient permissions', () => {
      getWrapper({
        setup: ({ canListShares }) => {
          const space = mock<SpaceResource>()
          const resource = mock<Resource>()

          const capabilityStore = useCapabilityStore()
          vi.mocked(capabilityStore).sharingApiEnabled = true

          const canList = canListShares({ space, resource })
          expect(canList).toBeTruthy()
        }
      })
    })
    it('returns false when sharing not enabled', () => {
      getWrapper({
        setup: ({ canListShares }) => {
          const space = mock<SpaceResource>()
          const resource = mock<Resource>()

          const capabilityStore = useCapabilityStore()
          vi.mocked(capabilityStore).sharingApiEnabled = false

          const canList = canListShares({ space, resource })
          expect(canList).toBeFalsy()
        }
      })
    })
    it('returns false in public spaces', () => {
      getWrapper({
        setup: ({ canListShares }) => {
          const space = mock<PublicSpaceResource>({ driveType: 'public' })
          const resource = mock<Resource>()

          const capabilityStore = useCapabilityStore()
          vi.mocked(capabilityStore).sharingApiEnabled = true

          const canList = canListShares({ space, resource })
          expect(canList).toBeFalsy()
        }
      })
    })
    it('returns false for personal space root resources', () => {
      getWrapper({
        setup: ({ canListShares }) => {
          const space = mock<PersonalSpaceResource>()
          const resource = mock<Resource>()

          const capabilityStore = useCapabilityStore()
          vi.mocked(capabilityStore).sharingApiEnabled = true

          const canList = canListShares({ space, resource })
          expect(canList).toBeFalsy()
        },
        isPersonalSpaceRoot: true
      })
    })
    it('returns false for trash resources', () => {
      getWrapper({
        setup: ({ canListShares }) => {
          const space = mock<SpaceResource>()
          const resource = mock<TrashResource>({ ddate: '2021-01-01T00:00:00Z' })

          const capabilityStore = useCapabilityStore()
          vi.mocked(capabilityStore).sharingApiEnabled = true

          const canList = canListShares({ space, resource })
          expect(canList).toBeFalsy()
        }
      })
    })
    describe('incoming share resources', () => {
      it('returns true with sufficient permissions', () => {
        getWrapper({
          setup: ({ canListShares }) => {
            const space = mock<SpaceResource>()
            const resource = mock<IncomingShareResource>({
              sharedWith: [mock<Identity>()],
              sharePermissions: [GraphSharePermission.readPermissions],
              outgoing: false
            })

            const capabilityStore = useCapabilityStore()
            vi.mocked(capabilityStore).sharingApiEnabled = true

            const canList = canListShares({ space, resource })
            expect(canList).toBeTruthy()
          }
        })
      })
      it('returns false with insufficient permissions', () => {
        getWrapper({
          setup: ({ canListShares }) => {
            const space = mock<SpaceResource>()
            const resource = mock<IncomingShareResource>({
              sharedWith: [mock<Identity>()],
              sharePermissions: [],
              outgoing: false
            })

            const capabilityStore = useCapabilityStore()
            vi.mocked(capabilityStore).sharingApiEnabled = true

            const canList = canListShares({ space, resource })
            expect(canList).toBeFalsy()
          }
        })
      })
    })
    describe('share spaces', () => {
      it('returns true with sufficient permissions', () => {
        getWrapper({
          setup: ({ canListShares }) => {
            const space = mock<ShareSpaceResource>({ driveType: 'share' })
            const resource = mock<Resource>()

            const capabilityStore = useCapabilityStore()
            vi.mocked(capabilityStore).sharingApiEnabled = true

            const canList = canListShares({ space, resource })
            expect(canList).toBeTruthy()
          },
          shareSpacePermissions: [GraphSharePermission.readPermissions]
        })
      })
      it('returns false with insufficient permissions', () => {
        getWrapper({
          setup: ({ canListShares }) => {
            const space = mock<ShareSpaceResource>({ driveType: 'share' })
            const resource = mock<Resource>()

            const capabilityStore = useCapabilityStore()
            vi.mocked(capabilityStore).sharingApiEnabled = true

            const canList = canListShares({ space, resource })
            expect(canList).toBeFalsy()
          },
          shareSpacePermissions: []
        })
      })
    })
  })
})

function getWrapper({
  setup,
  isPersonalSpaceRoot = false,
  shareSpacePermissions = []
}: {
  setup: (instance: ReturnType<typeof useCanListShares>) => void
  isPersonalSpaceRoot?: boolean
  shareSpacePermissions?: GraphSharePermission[]
}) {
  vi.mocked(useGetMatchingSpace).mockImplementation(() =>
    useGetMatchingSpaceMock({
      isPersonalSpaceRoot: () => isPersonalSpaceRoot
    })
  )

  vi.mocked(getPermissionsForSpaceMember).mockReturnValue(shareSpacePermissions)

  return {
    wrapper: getComposableWrapper(() => {
      const instance = useCanListShares()
      setup(instance)
    })
  }
}
