import { FocusTrap } from 'focus-trap-vue'
import OcDrop from '../OcDrop/OcDrop.vue'
import OcInfoDrop from './OcInfoDrop.vue'
import { PartialComponentProps, defaultPlugins, shallowMount } from '@ownclouders/web-test-helpers'

describe('OcInfoDrop', () => {
  function getWrapperWithProps(props: PartialComponentProps<typeof OcInfoDrop>) {
    return shallowMount(OcInfoDrop, {
      props: {
        ...props,
        title: props.title || 'test-title'
      },
      global: {
        components: {
          FocusTrap
        },
        plugins: [...defaultPlugins()],
        renderStubDefaultSlot: true,
        stubs: {
          OcDrop: true
        }
      }
    })
  }
  describe('should use props correctly', () => {
    it('should set title prop', () => {
      const wrapper = getWrapperWithProps({ title: 'test-my-title' })
      expect(wrapper.find('.info-title').text()).toBe('test-my-title')
    })
    it('should set text prop', () => {
      const wrapper = getWrapperWithProps({ text: 'test-my-text' })
      expect(wrapper.find('.info-text').text()).toBe('test-my-text')
    })
    it('should set list prop', () => {
      const listValues = [
        { text: 'a-list-value' },
        { text: 'b-list-value' },
        { text: 'c-list-value' }
      ]
      const wrapper = getWrapperWithProps({ list: listValues })
      const result = wrapper.find('.info-list').text()
      listValues.forEach((value) => {
        expect(result).toContain(value.text)
      })
    })
    it('should set a readMore link', () => {
      const wrapper = getWrapperWithProps({ readMoreLink: 'owncloud.design' })
      const attributes = wrapper.find('.info-more-link').attributes()
      expect(attributes['href']).toBe('owncloud.design')
      expect(attributes['target']).toBe('_blank')
    })
    it('should set end-text prop', () => {
      const wrapper = getWrapperWithProps({ endText: 'test-my-text' })
      expect(wrapper.find('.info-text-end').text()).toBe('test-my-text')
    })
    describe('focus trap', () => {
      it('is active if the drop is open', async () => {
        const wrapper = getWrapperWithProps({ title: 'title' })
        wrapper.findComponent<typeof OcDrop>('oc-drop-stub').vm.$emit('show-drop')
        await wrapper.vm.$nextTick()
        const focusTrap = wrapper.findComponent<typeof FocusTrap>('focus-trap-stub')
        expect(focusTrap.props('active')).toBeTruthy()
      })
      it('is not active if the drop is closed', () => {
        const wrapper = getWrapperWithProps({ title: 'title' })
        const focusTrap = wrapper.findComponent<typeof FocusTrap>('focus-trap-stub')
        expect(focusTrap.props('active')).toBeFalsy()
      })
    })
  })
})
