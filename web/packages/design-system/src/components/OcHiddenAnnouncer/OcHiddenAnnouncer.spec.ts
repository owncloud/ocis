import OcHiddenAnnouncer from './OcHiddenAnnouncer.vue'
import { mount } from '@ownclouders/web-test-helpers'

describe('OcHiddenAnnouncer', () => {
  function getWrapper(props = {}) {
    return mount(OcHiddenAnnouncer, {
      props: {
        announcement: 'Hidden announcer (please inspect element)',
        ...props
      }
    })
  }
  describe('level', () => {
    it.each(['polite', 'assertive', 'off'])(
      'should set the provided label as aria live',
      (level) => {
        const wrapper = getWrapper({ level: level })
        expect(wrapper.attributes('aria-live')).toBe(level)
      }
    )
  })
  describe('announcement', () => {
    it('should render the provided announcement text', () => {
      const wrapper = getWrapper()
      expect(wrapper.text()).toBe('Hidden announcer (please inspect element)')
    })
  })
})
