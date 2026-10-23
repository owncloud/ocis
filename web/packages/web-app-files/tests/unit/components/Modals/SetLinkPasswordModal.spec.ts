import { mock } from 'vitest-mock-extended'
import SetLinkPasswordModal from '../../../../src/components/Modals/SetLinkPasswordModal.vue'
import { defaultComponentMocks, defaultPlugins, shallowMount } from '@ownclouders/web-test-helpers'
import { Modal, useMessages, useSharesStore } from '@ownclouders/web-pkg'
import { Resource, LinkShare, SpaceResource } from '@ownclouders/web-client'

describe('SetLinkPasswordModal', () => {
  it('should render a text input field for the password', () => {
    const { wrapper } = getWrapper()

    expect(wrapper.find('oc-text-input-stub').exists()).toBeTruthy()
  })
  describe('confirm disabled state', () => {
    it('emits update:confirmDisabled with true after mount', () => {
      const { wrapper } = getWrapper()

      expect(wrapper.emitted('update:confirmDisabled')).toContainEqual([true])
    })
    it('re-emits update:confirmDisabled with true when the field is cleared', async () => {
      const { wrapper } = getWrapper()
      const input = wrapper.findComponent({ name: 'OcTextInput' })

      await input.vm.$emit('update:modelValue', 'abc')
      await input.vm.$emit('update:modelValue', '')

      const emissions = wrapper.emitted('update:confirmDisabled')
      expect(emissions[emissions.length - 1]).toEqual([true])
    })
    it('does not emit update:confirmDisabled when a non-empty value is entered', async () => {
      const { wrapper } = getWrapper()
      const input = wrapper.findComponent({ name: 'OcTextInput' })
      const emissionCountAfterMount = wrapper.emitted('update:confirmDisabled').length

      await input.vm.$emit('update:modelValue', 'abc')

      expect(wrapper.emitted('update:confirmDisabled').length).toBe(emissionCountAfterMount)
    })
  })
  describe('pressing Enter in the password field', () => {
    it('does not emit confirm when the field is untouched/empty', async () => {
      const { wrapper } = getWrapper()
      const input = wrapper.findComponent({ name: 'OcTextInput' })

      await input.trigger('keydown.enter')

      expect(wrapper.emitted('confirm')).toBeFalsy()
    })
    it('emits confirm after entering a value', async () => {
      const { wrapper } = getWrapper()
      const input = wrapper.findComponent({ name: 'OcTextInput' })

      await input.vm.$emit('update:modelValue', 'abc')
      await input.trigger('keydown.enter')

      expect(wrapper.emitted('confirm')).toBeTruthy()
    })
  })
  describe('method "onConfirm"', () => {
    it('updates the link', async () => {
      const { wrapper } = getWrapper()
      await wrapper.vm.onConfirm()

      const sharesStore = useSharesStore()
      expect(sharesStore.updateLink).toHaveBeenCalled()
      const { showMessage } = useMessages()
      expect(showMessage).toHaveBeenCalled()
    })
    it('shows an error message on error', async () => {
      const { wrapper } = getWrapper()
      const sharesStore = useSharesStore()
      vi.mocked(sharesStore.updateLink).mockRejectedValue(new Error(''))
      await wrapper.vm.onConfirm()

      const { showErrorMessage } = useMessages()
      expect(showErrorMessage).toHaveBeenCalled()
    })
  })
})

function getWrapper() {
  const mocks = { ...defaultComponentMocks() }

  return {
    mocks,
    wrapper: shallowMount(SetLinkPasswordModal, {
      props: {
        modal: mock<Modal>(),
        link: mock<LinkShare>(),
        space: mock<SpaceResource>(),
        resource: mock<Resource>()
      },
      global: {
        plugins: [...defaultPlugins()],
        mocks,
        provide: mocks
      }
    })
  }
}
