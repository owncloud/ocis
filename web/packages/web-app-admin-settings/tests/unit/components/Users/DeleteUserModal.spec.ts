import DeleteUserModal from '../../../../src/components/Users/DeleteUserModal.vue'
import { defaultComponentMocks, defaultPlugins, shallowMount } from '@ownclouders/web-test-helpers'
import { mock } from 'vitest-mock-extended'
import { User } from '@ownclouders/web-client/graph/generated'
import { Modal } from '@ownclouders/web-pkg'

describe('DeleteUserModal', () => {
  it('shows a hint when the current user is part of the selection', () => {
    // the default mocked user store uses id '1'
    const { wrapper } = getWrapper([mock<User>({ id: '1' }), mock<User>({ id: '2' })])
    expect(wrapper.find('.delete-user-own-account-hint').exists()).toBe(true)
  })
  it('does not show a hint when the current user is not part of the selection', () => {
    const { wrapper } = getWrapper([mock<User>({ id: '2' })])
    expect(wrapper.find('.delete-user-own-account-hint').exists()).toBe(false)
  })
})

function getWrapper(users = [mock<User>()]) {
  const mocks = defaultComponentMocks()

  return {
    mocks,
    wrapper: shallowMount(DeleteUserModal, {
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
