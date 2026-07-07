import { useUserActionsDelete } from '../../../../../src/composables/actions/users/useUserActionsDelete'
import { mock } from 'vitest-mock-extended'
import { unref } from 'vue'
import { User } from '@ownclouders/web-client/graph/generated'
import { useCapabilityStore } from '@ownclouders/web-pkg'
import {
  defaultComponentMocks,
  getComposableWrapper,
  writable
} from '@ownclouders/web-test-helpers'
import { useUserSettingsStore } from '../../../../../src/composables/stores/userSettings'

describe('useUserActionsDelete', () => {
  describe('method "isVisible"', () => {
    it.each([
      { resources: [], disabledViaCapability: false, isVisible: false },
      { resources: [mock<User>()], disabledViaCapability: false, isVisible: true },
      { resources: [mock<User>(), mock<User>()], disabledViaCapability: false, isVisible: true },
      { resources: [mock<User>(), mock<User>()], disabledViaCapability: true, isVisible: false }
    ])(
      'should only return true if 1 or more users are selected and not disabled via capability',
      ({ resources, disabledViaCapability, isVisible }) => {
        getWrapper({
          setup: ({ actions }) => {
            const capabilityStore = useCapabilityStore()
            writable(capabilityStore).graphUsersDeleteDisabled = !!disabledViaCapability
            expect(unref(actions)[0].isVisible({ resources })).toEqual(isVisible)
          }
        })
      }
    )
  })
  describe('method "deleteUsers"', () => {
    it('should successfully delete all given users and reload the users list', () => {
      getWrapper({
        setup: async ({ deleteUsers }, { clientService }) => {
          const user = mock<User>({ id: '1' })
          await deleteUsers([user])
          expect(clientService.graphAuthenticated.users.deleteUser).toHaveBeenCalledWith(user.id)
          const { removeUsers } = useUserSettingsStore()
          expect(removeUsers).toHaveBeenCalled()
        }
      })
    })
    it('should handle errors', () => {
      vi.spyOn(console, 'error').mockImplementation(() => undefined)
      getWrapper({
        setup: async ({ deleteUsers }, { clientService }) => {
          clientService.graphAuthenticated.users.deleteUser.mockRejectedValue({})
          const user = mock<User>({ id: '1' })
          await deleteUsers([user])
          expect(clientService.graphAuthenticated.users.deleteUser).toHaveBeenCalledWith(user.id)
          const { removeUsers } = useUserSettingsStore()
          expect(removeUsers).toHaveBeenCalled()
        }
      })
    })
  })
})

function getWrapper({
  setup
}: {
  setup: (
    instance: ReturnType<typeof useUserActionsDelete>,
    {
      clientService
    }: {
      clientService: ReturnType<typeof defaultComponentMocks>['$clientService']
    }
  ) => void
}) {
  const mocks = defaultComponentMocks()
  return {
    wrapper: getComposableWrapper(
      () => {
        const instance = useUserActionsDelete()
        setup(instance, { clientService: mocks.$clientService })
      },
      { mocks, provide: mocks }
    )
  }
}
