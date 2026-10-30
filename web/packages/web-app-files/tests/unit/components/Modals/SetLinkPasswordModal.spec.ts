import { mock } from 'vitest-mock-extended'
import SetLinkPasswordModal from '../../../../src/components/Modals/SetLinkPasswordModal.vue'
import { defaultComponentMocks, defaultPlugins, shallowMount } from '@ownclouders/web-test-helpers'
import {
  Modal,
  PasswordPolicyService,
  useMessages,
  usePasswordPolicyService,
  useSharesStore
} from '@ownclouders/web-pkg'
import { Resource, LinkShare, SpaceResource } from '@ownclouders/web-client'
import { PasswordPolicy } from '@ownclouders/design-system/helpers'

vi.mock('@ownclouders/web-pkg', async (importOriginal) => ({
  ...(await importOriginal<any>()),
  usePasswordPolicyService: vi.fn()
}))

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
    it('leaves the confirm state to the password policy check while typing', async () => {
      const { wrapper } = getWrapper()
      const input = wrapper.findComponent({ name: 'OcTextInput' })
      const emissionCountAfterMount = wrapper.emitted('update:confirmDisabled').length

      await input.vm.$emit('update:modelValue', 'abc')
      await input.vm.$emit('update:modelValue', '')

      expect(wrapper.emitted('update:confirmDisabled').length).toBe(emissionCountAfterMount)
    })
  })
  describe('pressing Enter in the password field', () => {
    it('does not emit confirm when the field is untouched/empty', async () => {
      const { wrapper } = getWrapper({ passwordPolicyFulfilled: false })
      const input = wrapper.findComponent({ name: 'OcTextInput' })

      await input.trigger('keydown.enter')

      expect(wrapper.emitted('confirm')).toBeFalsy()
    })
    it('does not emit confirm when the value does not fulfill the password policy', async () => {
      const { wrapper } = getWrapper({ passwordPolicyFulfilled: false })
      const input = wrapper.findComponent({ name: 'OcTextInput' })

      await input.vm.$emit('update:modelValue', ' ')
      await input.trigger('keydown.enter')

      expect(wrapper.emitted('confirm')).toBeFalsy()
    })
    it('emits confirm when the value fulfills the password policy', async () => {
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

function getWrapper({
  passwordPolicyFulfilled = true
}: { passwordPolicyFulfilled?: boolean } = {}) {
  vi.mocked(usePasswordPolicyService).mockReturnValue(
    mock<PasswordPolicyService>({
      getPolicy: () => mock<PasswordPolicy>({ check: () => passwordPolicyFulfilled })
    })
  )
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
