import { PartialComponentProps, defaultPlugins, shallowMount } from '@ownclouders/web-test-helpers'
import OcButton from './OcButton.vue'

describe('OcButton', () => {
  it('should display slot html', () => {
    const wrapper = getWrapperWithTestSlot()
    const slot = wrapper.find('p')
    expect(slot).toBeTruthy()
    expect(slot.attributes('class')).toBe('text')
    expect(slot.text()).toBe('Test button')
  })

  describe('click event', () => {
    it('should emit click event when click is triggered', async () => {
      const wrapper = getWrapperWithProps({})
      await wrapper.trigger('click')
      expect(wrapper.emitted('click')).toBeTruthy()
    })
    it.each`
      type
      ${'a'}
      ${'router-link'}
    `('should not emit click event when type is $type', async ({ type }) => {
      const wrapper = getWrapperWithProps({ type: type })
      await wrapper.trigger('click')
      expect(wrapper.emitted('click')).toBeFalsy()
    })
  })
  describe('when oc button is disabled', () => {
    let wrapper: ReturnType<typeof getWrapperWithProps>
    beforeEach(() => {
      wrapper = getWrapperWithProps({ disabled: true })
    })
    it('should have disabled attribute set to disabled', () => {
      // disabled: true => '' disabled: false => undefined ¯\_(ツ)_/¯
      expect(wrapper.attributes('disabled')).toBe('')
    })
    it('should not emit click event', async () => {
      await wrapper.trigger('click')
      expect(wrapper.emitted('click')).toBeFalsy()
    })
  })
  describe('different types of button', () => {
    it.each`
      type             | expectLink | expectButton | expectRouterLink
      ${'a'}           | ${true}    | ${false}     | ${false}
      ${'button'}      | ${false}   | ${true}      | ${false}
      ${'router-link'} | ${false}   | ${false}     | ${true}
    `('can behave as a $type', ({ type, expectLink, expectButton, expectRouterLink }) => {
      const wrapper = getWrapperWithProps({ type: type })
      expect(wrapper.find('a').exists()).toBe(expectLink)
      expect(wrapper.find('button').exists()).toBe(expectButton)
      expect(wrapper.find('router-link-stub').exists()).toBe(expectRouterLink)
    })
  })
  describe('different sizes of button', () => {
    it.each`
      size         | expectedClass
      ${'small'}   | ${'oc-button-s'}
      ${'medium'}  | ${'oc-button-m'}
      ${'large'}   | ${'oc-button-l'}
      ${'x-small'} | ${'oc-button-undefined'}
    `(
      'when size prop is set as $size class $expectedClass should be assigned',
      ({ size, expectedClass }) => {
        const wrapper = getWrapperWithProps({
          size: size
        })
        expect(wrapper.attributes('class')).toContain(expectedClass)
      }
    )
  })
  describe('default prop values', () => {
    it.each`
      name                 | expected
      ${'size'}            | ${'oc-button-m'}
      ${'variation'}       | ${'oc-button-passive'}
      ${'justify content'} | ${'oc-button-justify-content-center'}
      ${'gap size'}        | ${'oc-button-gap-m'}
      ${'appearance'}      | ${'oc-button-passive-outline'}
    `('should have attribute "$name" as "$expected"', ({ expected }) => {
      const wrapper = getWrapperWithProps({})
      expect(wrapper.attributes('class')).toContain(expected)
    })
  })
  describe('oc button appearance', () => {
    // appearance prop is combined with variation prop
    describe('when appearance is "filled"', () => {
      it('should not have extra appearance class', () => {
        const wrapper = getWrapperWithProps({
          appearance: 'filled'
        })
        expect(wrapper.attributes('class')).toContain('oc-button-passive')
        expect(wrapper.attributes('class')).not.toContain('oc-button-passive-raw')
        expect(wrapper.attributes('class')).not.toContain('oc-button-passive-outline')
      })
    })
    describe('when oc button is initialized with variation and appearance', () => {
      it.each`
        variation    | appearance   | expectedClass
        ${'success'} | ${'raw'}     | ${'oc-button-success oc-button-success-raw'}
        ${'success'} | ${'outline'} | ${'oc-button-success oc-button-success-outline'}
        ${'primary'} | ${'raw'}     | ${'oc-button-primary oc-button-primary-raw'}
        ${'primary'} | ${'outline'} | ${'oc-button-primary-outline'}
      `('should have extra appearance class', ({ variation, appearance, expectedClass }) => {
        const wrapper = getWrapperWithProps({
          appearance: appearance,
          variation: variation
        })
        expect(wrapper.attributes('class')).toContain(expectedClass)
      })
    })
  })
})

function getWrapperWithProps(props: PartialComponentProps<typeof OcButton>) {
  return shallowMount(OcButton, { props, global: { plugins: [...defaultPlugins()] } })
}
function getWrapperWithTestSlot() {
  const testSlots = { default: '<p class="text">Test button</p>' }
  return shallowMount(OcButton, { slots: testSlots, global: { plugins: [...defaultPlugins()] } })
}
