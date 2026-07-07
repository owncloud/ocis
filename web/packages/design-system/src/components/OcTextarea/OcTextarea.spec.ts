import { shallowMount } from '@ownclouders/web-test-helpers'
import OcTextarea from './OcTextarea.vue'

const defaultProps = {
  label: 'label'
}

describe('OcTextarea', () => {
  function getShallowWrapper(props = {}) {
    return shallowMount(OcTextarea, {
      props: {
        ...defaultProps,
        ...props
      }
    })
  }

  const selectors = {
    textareaMessage: '.oc-textarea-message span',
    textArea: '.oc-textarea'
  }
  describe('id prop', () => {
    const wrapper = getShallowWrapper({
      id: 'test-textarea-id',
      descriptionMessage: 'hello'
    })
    it('should set provided id to the textarea', () => {
      expect(wrapper.find('textarea').attributes().id).toBe('test-textarea-id')
    })
    it('should set label target for provided id', () => {
      expect(wrapper.find('label').attributes().for).toBe('test-textarea-id')
    })
    it('should set message id according to provided id', () => {
      expect(wrapper.find(selectors.textareaMessage).attributes().id).toBe(
        'test-textarea-id-message'
      )
    })
  })
  describe('label prop', () => {
    it('should set provided label to the textarea', () => {
      const wrapper = getShallowWrapper()
      expect(wrapper.find('label').text()).toBe('label')
    })
  })
  describe('when a description message is provided', () => {
    const wrapper = getShallowWrapper({ descriptionMessage: 'You should pass.' })
    it('should add the description class to the textarea message', () => {
      expect(wrapper.find(selectors.textareaMessage).attributes().class).toContain(
        'oc-textarea-description'
      )
    })
    it('should show the description message as the input message text', () => {
      expect(wrapper.find(selectors.textareaMessage).text()).toBe('You should pass.')
    })
  })
  describe('when a warning message is provided', () => {
    const wrapper = getShallowWrapper({ warningMessage: 'You may pass.' })
    it('should add the warning class to the textarea', () => {
      expect(wrapper.find('textarea').attributes().class).toContain('oc-textarea-warning')
    })
    it('should add the warning class to the textarea message', () => {
      expect(wrapper.find(selectors.textareaMessage).attributes().class).toContain(
        'oc-textarea-warning'
      )
    })
    it('should show the warning message as the textarea message text', () => {
      expect(wrapper.find(selectors.textareaMessage).text()).toBe('You may pass.')
    })
  })
  describe('when an error message is provided', () => {
    const wrapper = getShallowWrapper({ errorMessage: 'You shall not pass.' })
    it('should add the error class to the textarea', () => {
      expect(wrapper.find('textarea').attributes().class).toContain('oc-textarea-danger')
    })
    it('should add the error class to the textarea message', () => {
      expect(wrapper.find(selectors.textareaMessage).attributes().class).toContain(
        'oc-textarea-danger'
      )
    })
    it('should show the error message as the textarea message text', () => {
      expect(wrapper.find(selectors.textareaMessage).text()).toBe('You shall not pass.')
    })
    it('should set the input aria-invalid attribute to true', () => {
      expect(wrapper.find('textarea').attributes('aria-invalid')).toBe('true')
    })
  })
  describe('message priority', () => {
    it('should give error message top priority', () => {
      const wrapper = getShallowWrapper({
        errorMessage: 'You shall not pass.',
        warningMessage: 'You may pass.',
        descriptionMessage: 'Your should pass.'
      })
      const messageEl = wrapper.find('.oc-textarea-message span')
      expect(messageEl.attributes().class).toBe(
        'oc-textarea-description oc-textarea-warning oc-textarea-danger'
      )
      expect(messageEl.text()).toBe('You shall not pass.')
    })
    it('should give warning message priority over description message', () => {
      const wrapper = getShallowWrapper({
        warningMessage: 'You may pass.',
        descriptionMessage: 'Your should pass.'
      })
      const messageEl = wrapper.find(selectors.textareaMessage)
      expect(messageEl.attributes().class).toBe('oc-textarea-description oc-textarea-warning')
      expect(messageEl.text()).toBe('You may pass.')
    })
  })
  describe('input events', () => {
    it('should emit an input event on typing', async () => {
      const wrapper = getShallowWrapper()
      expect(wrapper.emitted('update:modelValue')).toBeFalsy()
      await wrapper.find('textarea').setValue('a')
      expect(wrapper.emitted('update:modelValue')).toBeTruthy()
      expect(wrapper.emitted('update:modelValue')[0][0]).toBe('a')
    })
  })
  describe('change events', () => {
    it('should emit a change event if submitOnEnter is true', async () => {
      const wrapper = getShallowWrapper({ submitOnEnter: true })
      expect(wrapper.emitted().change).toBeFalsy()
      await wrapper.find('textarea').trigger('keydown.enter')
      expect(wrapper.emitted().change).toBeTruthy()
    })
    it("shouldn't emit a change event if submitOnEnter is false", async () => {
      const wrapper = getShallowWrapper({ submitOnEnter: false })
      expect(wrapper.emitted().change).toBeFalsy()
      await wrapper.find('textarea').trigger('keydown.enter')
      expect(wrapper.emitted().change).toBeFalsy()
    })
  })
})
