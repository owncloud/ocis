import LoginModal from '../../../../src/components/Users/LoginModal.vue'
import { defaultComponentMocks, defaultPlugins, shallowMount } from '@ownclouders/web-test-helpers'
import { mock } from 'vitest-mock-extended'
import { User } from '@ownclouders/web-client/graph/generated'
import { Modal, useMessages } from '@ownclouders/web-pkg'
import { OcSelect } from '@ownclouders/design-system/components'
import { useUserSettingsStore } from '../../../../src/composables/stores/userSettings'

describe('LoginModal', () => {
  it('renders the input including two options', () => {
    const { wrapper } = getWrapper()
    expect(wrapper.html()).toMatchSnapshot()
  })
  it('shows a warning when the current user is being selected', () => {
    const { wrapper } = getWrapper([mock<User>({ id: '1' })])
    expect(
      wrapper.findComponent<typeof OcSelect>('oc-select-stub').props('warningMessage')
    ).toBeDefined()
  })
  describe('method "onConfirm"', () => {
    it('updates the login for all given users', async () => {
      const users = [mock<User>(), mock<User>()]
      const { wrapper, mocks } = getWrapper(users)
      mocks.$clientService.graphAuthenticated.users.editUser.mockResolvedValue(
        mock<User>({ id: 'e3515ffb-d264-4dfc-8506-6c239f6673b5' })
      )
      mocks.$clientService.graphAuthenticated.users.getUser.mockResolvedValue(
        mock<User>({ id: 'e3515ffb-d264-4dfc-8506-6c239f6673b5' })
      )

      await wrapper.vm.onConfirm()
      const { showMessage } = useMessages()
      expect(showMessage).toHaveBeenCalled()
      const { upsertUser } = useUserSettingsStore()
      expect(upsertUser).toHaveBeenCalledTimes(users.length)
      expect(mocks.$clientService.graphAuthenticated.users.editUser).toHaveBeenCalledTimes(
        users.length
      )
    })
    it('omits the currently logged in user', async () => {
      const users = [mock<User>({ id: '1' }), mock<User>()]
      const { wrapper, mocks } = getWrapper(users)
      mocks.$clientService.graphAuthenticated.users.editUser.mockResolvedValue(
        mock<User>({ id: 'e3515ffb-d264-4dfc-8506-6c239f6673b5' })
      )
      mocks.$clientService.graphAuthenticated.users.getUser.mockResolvedValue(
        mock<User>({ id: 'e3515ffb-d264-4dfc-8506-6c239f6673b5' })
      )

      await wrapper.vm.onConfirm()
      expect(mocks.$clientService.graphAuthenticated.users.editUser).toHaveBeenCalledTimes(1)
    })
    it('should show message on error', async () => {
      vi.spyOn(console, 'error').mockImplementation(() => undefined)

      const users = [mock<User>(), mock<User>()]
      const { wrapper, mocks } = getWrapper(users)
      mocks.$clientService.graphAuthenticated.users.editUser.mockRejectedValue(new Error(''))
      mocks.$clientService.graphAuthenticated.users.getUser.mockRejectedValue(new Error(''))

      await wrapper.vm.onConfirm()
      const { showErrorMessage } = useMessages()
      expect(showErrorMessage).toHaveBeenCalled()
    })
  })
})

function getWrapper(users = [mock<User>()]) {
  const mocks = defaultComponentMocks()

  return {
    mocks,
    wrapper: shallowMount(LoginModal, {
      props: {
        modal: mock<Modal>(),
        users
      },
      global: {
        provide: mocks,
        plugins: [...defaultPlugins()]
      }
    })
  }
}
