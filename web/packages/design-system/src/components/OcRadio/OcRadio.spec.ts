import OcRadio from './OcRadio.vue'
import { PartialComponentProps, mount } from '@ownclouders/web-test-helpers'
import { AvailableSizeType } from '../../helpers'

describe('OcRadio', () => {
  function getWrapper(props: PartialComponentProps<typeof OcRadio> = {}) {
    return mount(OcRadio, {
      props: {
        ...props,
        label: 'Test label'
      }
    })
  }
  const radioElementSelector = "input[type='radio']"
  it('should set provided id to the radio input', () => {
    const wrapper = getWrapper({ id: 'test-radio-input' })
    expect(wrapper.find('input').attributes('id')).toBe('test-radio-input')
    expect(wrapper.find('label').attributes('for')).toBe('test-radio-input')
  })
  describe('size prop', () => {
    it.each(['small', 'medium', 'large'] as AvailableSizeType[])(
      'should set input class according to the provided size',
      (size) => {
        const wrapper = getWrapper({
          size: size
        })
        const radioInput = wrapper.find(radioElementSelector)
        expect(radioInput.attributes('class')).toContain(`oc-radio-${size[0]}`)
      }
    )
  })
  describe('label', () => {
    it('should render the provided label for the input', () => {
      const wrapper = getWrapper()
      const labelEl = wrapper.find('label')
      expect(labelEl.exists()).toBeTruthy()
      expect(labelEl.text()).toBe('Test label')
    })
    describe('hideLabel prop with label', () => {
      it('should have invisible label class if enabled', () => {
        const wrapper = getWrapper({ hideLabel: true })
        const labelEl = wrapper.find('label')
        expect(labelEl.exists()).toBeTruthy()
        expect(labelEl.attributes('class')).toContain('oc-invisible-sr')
      })
      it('should not have invisible label class if disabled', () => {
        const wrapper = getWrapper({ hideLabel: false })
        const labelEl = wrapper.find('label')
        expect(labelEl.exists()).toBeTruthy()
        expect(labelEl.attributes('class')).not.toContain('oc-invisible-sr')
      })
    })
    describe('disabled prop with label', () => {
      it('should not have cursor pointer label class, if disabled', () => {
        const wrapper = getWrapper({ disabled: true })
        const labelEl = wrapper.find('label')
        expect(labelEl.exists()).toBeTruthy()
        expect(labelEl.attributes('class')).not.toContain('oc-cursor-pointer')
      })
      it('should have cursor pointer label class, if not disabled', () => {
        const wrapper = getWrapper({ disabled: false })
        const labelEl = wrapper.find('label')
        expect(labelEl.exists()).toBeTruthy()
        expect(labelEl.attributes('class')).toContain('oc-cursor-pointer')
      })
    })
  })
  describe('radio input', () => {
    it('should emit input event if checked', async () => {
      const wrapper = getWrapper()
      const radioInput = wrapper.find<HTMLInputElement>(radioElementSelector)
      expect(wrapper.emitted('update:modelValue')).toBeFalsy()
      await radioInput.setValue(true)
      expect(radioInput.element.checked).toBeTruthy()
      expect(wrapper.emitted('update:modelValue')).toBeTruthy()
    })
    it('should not emit input event if disabled', async () => {
      const wrapper = getWrapper({ disabled: true })
      const radioInput = wrapper.find<HTMLInputElement>(radioElementSelector)
      expect(wrapper.emitted('update:modelValue')).toBeFalsy()
      await radioInput.trigger('click')
      expect(radioInput.element.checked).toBeFalsy()
      expect(wrapper.emitted('update:modelValue')).toBeFalsy()
    })
  })
  describe('input options', () => {
    const rawData: { availableOptions: string[]; selectedOption: unknown } = {
      availableOptions: ['Water', 'Wine', 'Beer'],
      selectedOption: null
    }
    const Component = {
      template:
        '<div><oc-radio v-for="o in availableOptions" :key="o" v-model="selectedOption" :option="o" :label="o" /></div>',
      data: () => ({ ...rawData })
    }
    const options = {
      data: () => ({ ...rawData }),
      global: {
        stubs: { OcRadio: OcRadio }
      }
    }
    it('should set the provide option as input value', () => {
      const wrapper = mount(Component, options)
      const inputs = wrapper.findAll(radioElementSelector)
      rawData.availableOptions.forEach((option, index) => {
        expect(inputs.at(index).attributes('value')).toBe(option)
      })
    })
    it('should set aria-checked if option equals selected value', async () => {
      const wrapper = mount(Component, options)
      const inputs = wrapper.findAll(radioElementSelector)
      await inputs.at(0).setValue(true)
      expect(inputs.at(0).attributes('aria-checked')).toBe('true')
      expect(inputs.at(1).attributes('aria-checked')).toBe('false')
      expect(inputs.at(2).attributes('aria-checked')).toBe('false')
      await inputs.at(1).setValue(true)
      expect(inputs.at(0).attributes('aria-checked')).toBe('false')
      expect(inputs.at(1).attributes('aria-checked')).toBe('true')
      expect(inputs.at(2).attributes('aria-checked')).toBe('false')
    })
  })
})
