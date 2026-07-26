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

  describe('actions prop', () => {
    it('should not render actions if none are provided', () => {
      const wrapper = getWrapper()

      expect(wrapper.find(selectors.actionButton).exists()).toBeFalsy()
    })

    it('should render an action button for each action', () => {
      const actions = [
        { label: 'Undo', ariaLabel: 'Undo delete', onClick: vi.fn() },
        { label: 'Retry', onClick: vi.fn() }
      ]
      const wrapper = getWrapper({ actions })
      const buttons = wrapper.findAll(selectors.actionButton)

      expect(buttons.length).toBe(2)
      expect(buttons[0].text()).toBe('Undo')
      expect(buttons[0].attributes('aria-label')).toBe('Undo delete')
      expect(buttons[1].attributes('aria-label')).toBe('Retry')
    })

    it('should call onClick when an action button is clicked', async () => {
      const onClick = vi.fn()
      const wrapper = getWrapper({ actions: [{ label: 'Undo', onClick }] })

      await wrapper.find(selectors.actionButton).trigger('click')

      expect(onClick).toHaveBeenCalledTimes(1)
    })

    it('should close the notification when an action button is clicked, so it cannot be clicked twice', async () => {
      const onClick = vi.fn()
      const wrapper = getWrapper({ actions: [{ label: 'Undo', onClick }] })

      await wrapper.find(selectors.actionButton).trigger('click')

      expect(wrapper.emitted('close')).toBeTruthy()
      expect(wrapper.emitted('close').length).toBe(1)
    })

    it('should not move focus into the notification on mount', () => {
      const wrapper = getWrapper(
        { actions: [{ label: 'Undo', onClick: vi.fn() }] },
        { attachTo: document.body }
      )

      expect(document.activeElement).toBe(document.body)
      wrapper.unmount()
    })
  })

  describe('auto-dismiss pause/resume', () => {
    it('should pause the dismiss timer on hover and resume on mouse leave', async () => {
      const wrapper = getWrapper({ timeout: 10 })
      const interactiveEl = wrapper.find(selectors.interactiveWrapper)

      await interactiveEl.trigger('mouseenter')
      vi.advanceTimersByTime(10000)
      expect(wrapper.emitted('close')).toBeFalsy()

      await interactiveEl.trigger('mouseleave')
      vi.advanceTimersByTime(10000)
      expect(wrapper.emitted('close')).toBeTruthy()
    })

    it('should pause the dismiss timer on focus and resume on blur', async () => {
      const wrapper = getWrapper({ timeout: 10 })
      const interactiveEl = wrapper.find(selectors.interactiveWrapper)

      await interactiveEl.trigger('focusin')
      vi.advanceTimersByTime(10000)
      expect(wrapper.emitted('close')).toBeFalsy()

      await interactiveEl.trigger('focusout')
      vi.advanceTimersByTime(10000)
      expect(wrapper.emitted('close')).toBeTruthy()
    })

    it('should still auto-dismiss after hover and focus overlap (e.g. clicking a button)', async () => {
      const wrapper = getWrapper({ timeout: 10 })
      const interactiveEl = wrapper.find(selectors.interactiveWrapper)

      // hover and focus both engage (as happens when clicking an action button by mouse)
      await interactiveEl.trigger('mouseenter')
      await interactiveEl.trigger('focusin')
      // then only the mouse leaves, focus remains (e.g. keyboard user tabs through)
      await interactiveEl.trigger('mouseleave')
      vi.advanceTimersByTime(10000)
      expect(wrapper.emitted('close')).toBeFalsy()

      await interactiveEl.trigger('focusout')
      vi.advanceTimersByTime(10000)
      expect(wrapper.emitted('close')).toBeTruthy()
    })
  })

  const selectors = {
    messageTitle: '.oc-notification-message-title',
    messageContent: '.oc-notification-message-content',
    messageWrapper: '.oc-notification-message div',
    interactiveWrapper: '.oc-notification-message-interactive',
    errorLog: '.oc-error-log',
    errorLogToggleButton: '.oc-notification-message-error-log-toggle-button',
    actionButton: '.oc-notification-message-action-button'
  }

  function getWrapper(props = {}, mountOptions = {}) {
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
      },
      ...mountOptions
    })
  }
})
