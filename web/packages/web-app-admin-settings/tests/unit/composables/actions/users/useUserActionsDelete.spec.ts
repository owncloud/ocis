import { useUserActionsDelete } from '../../../../../src/composables/actions/users/useUserActionsDelete'
import DeleteUserModal from '../../../../../src/components/Users/DeleteUserModal.vue'
import { mock } from 'vitest-mock-extended'
import { Mock } from 'vitest'
import { unref } from 'vue'
import { User } from '@ownclouders/web-client/graph/generated'
import { useCapabilityStore, useModals } from '@ownclouders/web-pkg'
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
    it('should not delete the current user when included in the selection', () => {
      getWrapper({
        currentUserId: 'self',
        setup: async ({ deleteUsers }, { clientService }) => {
          const currentUser = mock<User>({ id: 'self' })
          const otherUser = mock<User>({ id: 'other' })
          await deleteUsers([currentUser, otherUser])
          expect(clientService.graphAuthenticated.users.deleteUser).toHaveBeenCalledWith(
            otherUser.id
          )
          expect(clientService.graphAuthenticated.users.deleteUser).not.toHaveBeenCalledWith(
            currentUser.id
          )
        }
      })
    })
    it('should do nothing when the current user is the only selected user', () => {
      getWrapper({
        currentUserId: 'self',
        setup: async ({ deleteUsers }, { clientService }) => {
          const currentUser = mock<User>({ id: 'self' })
          await deleteUsers([currentUser])
          expect(clientService.graphAuthenticated.users.deleteUser).not.toHaveBeenCalled()
        }
      })
    })
  })
  describe('method "handler"', () => {
    it('dispatches the delete modal with the selected users', () => {
      getWrapper({
        setup: ({ actions }) => {
          const { dispatchModal } = useModals()
          const resources = [mock<User>({ id: 'self' }), mock<User>({ id: 'other' })]
          unref(actions)[0].handler({ resources })
          expect(dispatchModal).toHaveBeenCalledWith(
            expect.objectContaining({
              variation: 'danger',
              customComponent: DeleteUserModal,
              customComponentAttrs: expect.any(Function)
            })
          )
          const attrs = (dispatchModal as unknown as Mock).mock.calls[0][0].customComponentAttrs()
          expect(attrs.users).toEqual(resources)
        }
      })
    })
  })
})

function getWrapper({
  setup,
  currentUserId = '0'
}: {
  setup: (
    instance: ReturnType<typeof useUserActionsDelete>,
    {
      clientService
    }: {
      clientService: ReturnType<typeof defaultComponentMocks>['$clientService']
    }
  ) => void
  currentUserId?: string
}) {
  const mocks = defaultComponentMocks()
  return {
    wrapper: getComposableWrapper(
      () => {
        const instance = useUserActionsDelete()
        setup(instance, { clientService: mocks.$clientService })
      },
      {
        mocks,
        provide: mocks,
        pluginOptions: {
          piniaOptions: {
            userState: { user: { id: currentUserId } as User }
          }
        }
      }
    )
  }
}
