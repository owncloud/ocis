import { defaultPlugins, mount, shallowMount } from '@ownclouders/web-test-helpers'
import { FocusTrap } from 'focus-trap-vue'
import Modal from './OcModal.vue'
import OcButton from './../OcButton/OcButton.vue'

const defaultProps = {
  title: 'Example title',
  message: 'Example message'
}

const inputProps = {
  title: 'Create new folder',
  hasInput: true,
  inputValue: 'New folder',
  inputLabel: 'Folder name'
}

describe('OcModal', () => {
  it('displays correct variation', () => {
    const wrapper = shallowMount(Modal, {
      props: {
        ...defaultProps,
        variation: 'danger'
      },
      global: {
        components: {
          FocusTrap
        },

        renderStubDefaultSlot: true,
        plugins: [...defaultPlugins()]
      }
    })

    expect(wrapper.findAll('.oc-modal-danger').length).toBe(1)
  })

  it('hides icon if not specified', () => {
    const wrapper = shallowMount(Modal, {
      global: {
        components: {
          FocusTrap
        },
        renderStubDefaultSlot: true,
        plugins: [...defaultPlugins()]
      },
      props: {
        ...defaultProps
      }
    })

    expect(wrapper.findAll('.oc-icon').length).toBe(0)
    expect(wrapper.html()).toMatchSnapshot()
  })

  it('overrides props message with slot', () => {
    const wrapper = shallowMount(Modal, {
      global: {
        components: {
          FocusTrap
        },
        renderStubDefaultSlot: true,
        plugins: [...defaultPlugins()]
      },
      props: {
        ...defaultProps
      },
      slots: {
        content: '<p>Slot message</p>'
      }
    })

    expect(wrapper.find('.oc-modal-body-message > p').text()).toMatch('Slot message')
    expect(wrapper.html()).toMatchSnapshot()
  })

  it('matches snapshot', () => {
    const wrapper = shallowMount(Modal, {
      global: {
        components: {
          FocusTrap
        },
        renderStubDefaultSlot: true,
        plugins: [...defaultPlugins()]
      },
      props: {
        ...defaultProps,
        icon: 'info'
      }
    })

    expect(wrapper.html()).toMatchSnapshot()
  })

  it('displays input', () => {
    const wrapper = shallowMount(Modal, {
      global: {
        components: {
          FocusTrap
        },
        renderStubDefaultSlot: true,
        plugins: [...defaultPlugins()]
      },
      props: inputProps
    })

    expect(wrapper.findAll('.oc-modal-body-input').length).toBe(1)
    expect(wrapper.html()).toMatchSnapshot()
  })

  it('displays loading state', async () => {
    const waitForSpinnerToShow = async () => {
      await wrapper.vm.$nextTick()
      return new Promise((resolve) => setTimeout(resolve, 1000))
    }

    const wrapper = mount(Modal, {
      global: {
        components: {
          FocusTrap
        },
        renderStubDefaultSlot: true,
        plugins: [...defaultPlugins()],
        stubs: {
          'focus-trap': true
        }
      },
      props: {
        ...defaultProps,
        isLoading: true
      }
    })

    const cancelButton = wrapper.find('.oc-modal-body-actions-cancel')
    const confirmButton = wrapper.find('.oc-modal-body-actions-confirm')

    expect(cancelButton.attributes('disabled')).toBeDefined()
    expect(confirmButton.attributes('disabled')).toBeDefined()
    expect(
      wrapper.findComponent<typeof OcButton>('.oc-modal-body-actions-confirm').props('showSpinner')
    ).toBeFalsy()
    expect(
      wrapper.findComponent<typeof OcButton>('.oc-modal-body-actions-confirm').props('appearance')
    ).toEqual('filled')

    await waitForSpinnerToShow()

    expect(
      wrapper.findComponent<typeof OcButton>('.oc-modal-body-actions-confirm').props('showSpinner')
    ).toBeTruthy()
    expect(
      wrapper.findComponent<typeof OcButton>('.oc-modal-body-actions-confirm').props('appearance')
    ).toEqual('outline')
    expect(wrapper.html()).toMatchSnapshot()
  })
})
