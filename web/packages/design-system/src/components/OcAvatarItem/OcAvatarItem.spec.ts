import OcAvatarItem from './OcAvatarItem.vue'
import { mount } from '@ownclouders/web-test-helpers'
describe('OcAvatarItem', () => {
  function getWrapperWithProps(props = {}) {
    return mount(OcAvatarItem, {
      props: {
        ...props,
        name: 'test link'
      }
    })
  }
  describe("when prop 'name' is provided", () => {
    it('should set data test user attribute on wrapper', () => {
      const wrapper = getWrapperWithProps()
      expect(wrapper.attributes('data-test-item-name')).toBe('test link')
    })
  })
  describe('accessibleLabel', () => {
    it('should not be set when value is empty string', () => {
      const wrapper = getWrapperWithProps({
        accessibleLabel: ''
      })
      expect(wrapper.attributes('aria-label')).toBeFalsy()
      expect(wrapper.attributes('role')).toBeFalsy()
      expect(wrapper.attributes('aria-hidden')).toBe('true')
      expect(wrapper.attributes('focusable')).toBe('false')
    })
    it('should be set when value is not empty string', () => {
      const wrapper = getWrapperWithProps({
        accessibleLabel: 'test label'
      })
      expect(wrapper.attributes('aria-label')).toBe('test label')
      expect(wrapper.attributes('role')).toBe('img')
      expect(wrapper.attributes('aria-hidden')).toBeFalsy()
      expect(wrapper.attributes('focusable')).toBeFalsy()
    })
  })
})
