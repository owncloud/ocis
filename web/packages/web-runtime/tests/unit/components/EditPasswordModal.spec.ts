import { Modal } from '@ownclouders/web-pkg'
import EditPasswordModal from '../../../src/components/EditPasswordModal.vue'
import { defaultPlugins, shallowMount } from '@ownclouders/web-test-helpers'
import { mock } from 'vitest-mock-extended'

describe('EditPasswordModal', () => {
  describe('computed method "confirmButtonDisabled"', () => {
    it('should be true if any data set is invalid', () => {
      const { wrapper } = getWrapper()
      wrapper.vm.currentPassword = ''
      expect(wrapper.vm.confirmButtonDisabled).toBeTruthy()
    })
    it('should be false if no data set is invalid', () => {
      const { wrapper } = getWrapper()
      wrapper.vm.currentPassword = 'password'
      wrapper.vm.newPassword = 'newpassword'
      wrapper.vm.newPasswordConfirm = 'newpassword'
      expect(wrapper.vm.confirmButtonDisabled).toBeFalsy()
    })
  })

  describe('method "validatePasswordConfirm"', () => {
    it('should be true if passwords are identical', () => {
      const { wrapper } = getWrapper()
      wrapper.vm.newPassword = 'newpassword'
      wrapper.vm.newPasswordConfirm = 'newpassword'
      expect(wrapper.vm.validatePasswordConfirm).toBeTruthy()
    })
    it('should be false if passwords are not identical', () => {
      const { wrapper } = getWrapper()
      wrapper.vm.newPassword = 'newpassword'
      wrapper.vm.newPasswordConfirm = 'anothernewpassword'
      expect(wrapper.vm.validatePasswordConfirm).toBeTruthy()
    })
  })
})

function getWrapper() {
  return {
    wrapper: shallowMount(EditPasswordModal, {
      props: {
        modal: mock<Modal>()
      },
      global: {
        plugins: [...defaultPlugins()]
      }
    })
  }
}
