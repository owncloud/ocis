import OcNotificationMessage from './OcNotificationMessage.vue'
import { defaultPlugins, mount } from '@ownclouders/web-test-helpers'

describe('OcNotificationMessage', () => {
  beforeEach(() => {
    vi.useFakeTimers()
  })

  afterEach(() => {
    vi.useRealTimers()
  })

  describe('title prop', () => {
    it('should set notification message title', () => {
      const wrapper = getWrapper()

      expect(wrapper.find(selectors.messageTitle).text()).toBe('Test passed')
    })
  })

  describe('message prop', () => {
    it('should render message, if message is provided', () => {
      const wrapper = getWrapper({ message: 'This is a test notification message' })
      const messageEl = wrapper.find(selectors.messageContent)

      expect(messageEl.exists()).toBeTruthy()
      expect(messageEl.text()).toBe('This is a test notification message')
    })
    it('should not render message, if message is not provided', () => {
      const wrapper = getWrapper({})
      const messageEl = wrapper.find(selectors.messageContent)

      expect(messageEl.exists()).toBeFalsy()
    })
  })

  describe('status prop', () => {
    it.each(['passive', 'primary', 'success', 'warning', 'danger'])(
      'should set provided status as class for wrapper',
      (status) => {
        const wrapper = getWrapper({ status: status })

        expect(wrapper.attributes('class')).toContain(`oc-notification-message-${status}`)
      }
    )

    it('should set status as icon variation', () => {
      const wrapper = getWrapper({ status: 'primary' })
      const iconElement = wrapper.find('oc-icon-stub')

      expect(iconElement.exists()).toBeTruthy()
      expect(iconElement.attributes('variation')).toBe('primary')
    })

    describe('role and aria live of message content wrapper', () => {
      it("should set role as 'status' and aria-live as 'polite' if status is not danger", () => {
        const wrapper = getWrapper({})
        const messageContentEl = wrapper.find(selectors.messageWrapper)

        expect(messageContentEl.attributes('role')).toBe('status')
        expect(messageContentEl.attributes('aria-live')).toBe('polite')
      })

      it("should set role as 'alert' and aria-live as 'assertive' if status is danger", () => {
        const wrapper = getWrapper({ status: 'danger' })
        const messageContentEl = wrapper.find(selectors.messageWrapper)

        expect(messageContentEl.attributes('role')).toBe('alert')
        expect(messageContentEl.attributes('aria-live')).toBe('assertive')
      })
    })
  })

  describe('errorLogContent prop', () => {
    it('should render OcErrorLogComponent, if errorLogContent is provided', async () => {
      const wrapper = getWrapper({ errorLogContent: 'X-REQUEST-ID: 1234' })
      const errorLogToggleButtonEl = wrapper.find(selectors.errorLogToggleButton)

      expect(errorLogToggleButtonEl.exists()).toBeTruthy()
      await errorLogToggleButtonEl.trigger('click')

      const errorLogEl = wrapper.find(selectors.errorLog)
      expect(errorLogEl.exists()).toBeTruthy()
    })

    it('should toggle aria-expanded on the details button when clicked', async () => {
      const wrapper = getWrapper({ errorLogContent: 'X-REQUEST-ID: 1234' })
      const errorLogToggleButtonEl = wrapper.find(selectors.errorLogToggleButton)

      expect(errorLogToggleButtonEl.attributes('aria-expanded')).toBe('false')
      await errorLogToggleButtonEl.trigger('click')
      expect(errorLogToggleButtonEl.attributes('aria-expanded')).toBe('true')
      await errorLogToggleButtonEl.trigger('click')
      expect(errorLogToggleButtonEl.attributes('aria-expanded')).toBe('false')
    })

    it('should not render OcErrorLogComponent, if errorLogContent is not provided', () => {
      const wrapper = getWrapper()
      const errorLogToggleButtonEl = wrapper.find(selectors.errorLogToggleButton)

      expect(errorLogToggleButtonEl.exists()).toBeFalsy()
    })
  })

  it('should emit close after set timout', () => {
    const wrapper = getWrapper({ timeout: 1 })

    expect(wrapper.emitted('close')).toBeFalsy()
    vi.advanceTimersByTime(1000)
    expect(wrapper.emitted('close')).toBeTruthy()
  })

  const selectors = {
    messageTitle: '.oc-notification-message-title',
    messageContent: '.oc-notification-message-content',
    messageWrapper: '.oc-notification-message div',
    errorLog: '.oc-error-log',
    errorLogToggleButton: '.oc-notification-message-error-log-toggle-button'
  }

  function getWrapper(props = {}) {
    return mount(OcNotificationMessage, {
      props: {
        ...props,
        title: 'Test passed'
      },
      global: {
        stubs: {
          'oc-icon': true
        },
        plugins: defaultPlugins()
      }
    })
  }
})
