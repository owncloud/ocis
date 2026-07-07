import OcNotifications from './OcNotifications.vue'
import { shallowMount } from '@ownclouders/web-test-helpers'

describe('OcNotifications', () => {
  function getWrapper(options = {}) {
    return shallowMount(OcNotifications, options)
  }
  describe('position prop', () => {
    it.each(['top-left', 'top-center', 'top-right'])(
      'should set provided position as class for wrapper',
      (position) => {
        const wrapper = getWrapper({
          props: { position: position }
        })
        expect(wrapper.attributes('class')).toContain(`oc-notification-${position}`)
      }
    )
  })

  it('should render provided slot html', () => {
    const wrapper = getWrapper({
      slots: {
        default:
          "<oc-notification-message title='test notification title' message='Testing is good.'/>"
      },
      global: {
        renderStubDefaultSlot: true,
        stubs: {
          OcNotificationMessage: true
        }
      }
    })
    const slotEl = wrapper.find('oc-notification-message-stub')
    expect(slotEl.exists()).toBeTruthy()
    expect(slotEl.attributes('title')).toBe('test notification title')
    expect(slotEl.attributes('message')).toBe('Testing is good.')
  })
})
