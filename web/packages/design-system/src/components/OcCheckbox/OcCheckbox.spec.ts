import OcCheckbox from './OcCheckbox.vue'
import { PartialComponentProps, mount } from '@ownclouders/web-test-helpers'

describe('OcCheckbox', () => {
  function getWrapperWithProps(props: PartialComponentProps<typeof OcCheckbox>) {
    return mount(OcCheckbox, {
      props: {
        label: 'test label',
        ...props
      }
    })
  }

  const checkboxSelector = "input[type='checkbox']"

  describe('input id', () => {
    it('should set the provided id for the input', () => {
      const wrapper = getWrapperWithProps({ id: 'test-id' })
      const checkBoxElement = wrapper.find(checkboxSelector)
      expect(checkBoxElement.exists()).toBeTruthy()
      expect(checkBoxElement.attributes('id')).toBe('test-id')
    })
  })
  describe('input label', () => {
    it('should set the provided label for the input', () => {
      const wrapper = getWrapperWithProps({ id: 'test-id' })
      const checkBoxLabelElement = wrapper.find('label')
      expect(checkBoxLabelElement.exists()).toBeTruthy()
      expect(checkBoxLabelElement.attributes('for')).toBe('test-id')
      expect(checkBoxLabelElement.text()).toBe('test label')
    })
    it("should hide label if 'labelHidden' prop is enabled", () => {
      const wrapper = getWrapperWithProps({ labelHidden: true })
      const checkBoxLabelElement = wrapper.find('label')
      expect(checkBoxLabelElement.exists()).toBeFalsy()
      const checkboxElement = wrapper.find<HTMLInputElement>(checkboxSelector)
      expect(checkboxElement.attributes('aria-label')).toContain('test label')
    })
    it('should have cursor pointer property if not disabled', () => {
      const wrapper = getWrapperWithProps({ disabled: false })
      const checkBoxLabelElement = wrapper.find('label')
      expect(checkBoxLabelElement.exists()).toBeTruthy()
      expect(checkBoxLabelElement.attributes('class')).toContain('oc-cursor-pointer')
    })
  })
  describe('input size', () => {
    type Item = {
      size: 'small' | 'medium' | 'large'
      class: string
    }
    it.each([
      { size: 'small', class: 'oc-checkbox-s' },
      { size: 'medium', class: 'oc-checkbox-m' },
      { size: 'large', class: 'oc-checkbox-l' }
    ])('valid size options', (item: Item) => {
      const wrapper = getWrapperWithProps({ size: item.size })
      const checkboxElement = wrapper.find<HTMLInputElement>(checkboxSelector)
      expect(checkboxElement.exists()).toBeTruthy()
      expect(checkboxElement.attributes('class')).toContain(item.class)
    })
  })
  describe('set checked', () => {
    it('should set check on input change', async () => {
      const wrapper = await getWrapperWithProps({})
      const checkbox = wrapper.find<HTMLInputElement>(checkboxSelector)
      expect(checkbox.element.checked).toBeFalsy()
      await checkbox.setValue(true)
      expect(wrapper.emitted('update:modelValue')).toBeTruthy()
      // the model is owned by the parent, the native state follows it
      await wrapper.setProps({ modelValue: true })
      expect(checkbox.element.checked).toBeTruthy()
    })
    it('restores the native checked state if the model does not follow the input', async () => {
      const wrapper = getWrapperWithProps({ modelValue: true })
      const checkbox = wrapper.find<HTMLInputElement>(checkboxSelector)
      expect(checkbox.element.checked).toBeTruthy()
      // the browser toggles the native state on click, the parent keeps the model as it is
      await checkbox.setValue(false)
      expect(wrapper.emitted('update:modelValue')).toEqual([[false]])
      await wrapper.vm.$nextTick()
      expect(checkbox.element.checked).toBeTruthy()
    })
  })
})
