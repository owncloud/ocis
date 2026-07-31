import { useUserActionsEditQuota } from '../../../../../src/composables/actions/users/useUserActionsEditQuota'
import {
  defaultComponentMocks,
  getComposableWrapper,
  writable
} from '@ownclouders/web-test-helpers'
import { unref } from 'vue'
import { useCapabilityStore, useModals } from '@ownclouders/web-pkg'
import { User } from '@ownclouders/web-client/graph/generated'

describe('useUserActionsEditQuota', () => {
  describe('isVisible property', () => {
    it('should be false when not resource given', () => {
      getWrapper({
        setup: ({ actions }) => {
          expect(unref(actions)[0].isVisible({ resources: [] })).toBe(false)
        }
      })
    })
    it('should be true when the current user has the "set-space-quota"-permission', () => {
      const userMock = {
        id: '1',
        drive: {
          name: 'some-drive',
          quota: {}
        }
      } as User
      getWrapper({
        canEditSpaceQuota: true,
        setup: ({ actions }) => {
          expect(unref(actions)[0].isVisible({ resources: [userMock] })).toBe(true)
        }
      })
    })
    it('should be false when the current user does not have the "set-space-quota"-permission', () => {
      const userMock = {
        id: '1',
        drive: {
          name: 'some-drive',
          quota: {}
        }
      } as User
      getWrapper({
        canEditSpaceQuota: false,
        setup: ({ actions }) => {
          expect(unref(actions)[0].isVisible({ resources: [userMock] })).toBe(false)
        }
      })
    })
    it('should false if included in capability readOnlyUserAttributes list', () => {
      getWrapper({
        setup: ({ actions }) => {
          const userMock = {
            id: '1',
            drive: {
              name: 'some-drive',
              quota: {}
            }
          } as User

          const capabilityStore = useCapabilityStore()
          writable(capabilityStore).graphUsersReadOnlyAttributes = ['drive.quota']

          expect(unref(actions)[0].isVisible({ resources: [userMock] })).toEqual(false)
        }
      })
    })
  })
  describe('handler', () => {
    it('should create a modal', () => {
      getWrapper({
        setup: async ({ actions }) => {
          const { dispatchModal } = useModals()
          await unref(actions)[0].handler({ resources: [] })
          expect(dispatchModal).toHaveBeenCalled()
        }
      })
    })
  })
})

function getWrapper({
  canEditSpaceQuota = false,
  setup
}: {
  canEditSpaceQuota?: boolean
  setup: (instance: ReturnType<typeof useUserActionsEditQuota>) => void
}) {
  const mocks = defaultComponentMocks()

  return {
    wrapper: getComposableWrapper(
      () => {
        const instance = useUserActionsEditQuota()
        setup(instance)
      },
      {
        mocks,
        pluginOptions: {
          abilities: canEditSpaceQuota ? [{ action: 'set-quota-all', subject: 'Drive' }] : []
        }
      }
    )
  }
}
