import OcGrid from './OcGrid.vue'
import { mount } from '@ownclouders/web-test-helpers'

describe('OcGrid', () => {
  function getWrapper(props = {}) {
    return mount(OcGrid, {
      props: props
    })
  }
  describe('gutter', () => {
    it.each(['small', 'medium', 'large', 'collapse'])(
      'should set provided gutter value',
      (gutter) => {
        const wrapper = getWrapper({
          gutter: gutter
        })
        expect(wrapper.attributes('class')).toBe('oc-grid-' + gutter)
      }
    )
  })
  describe('when flex prop is true', () => {
    it('should set grid flex class', () => {
      const wrapper = getWrapper({ flex: true })
      expect(wrapper.attributes('class')).toContain('oc-flex')
      expect(wrapper.attributes('class')).toContain('oc-flex-middle')
    })
  })
})
