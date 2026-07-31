import FeedbackLink from '../../../../src/components/Topbar/FeedbackLink.vue'
import { defaultPlugins, mount } from '@ownclouders/web-test-helpers'

describe('FeedbackLink component', () => {
  describe('properties', () => {
    it('allows to overwrite the link href', async () => {
      const { wrapper } = getWrapper()
      const url = 'https://some-link.tld/'
      await wrapper.setProps({ href: url })
      expect(wrapper.find('a').attributes().href).toEqual(url)
    })
    it('allows to overwrite the link ariaLabel', async () => {
      const { wrapper } = getWrapper()
      const ariaLabel = 'some aria label'
      await wrapper.setProps({ ariaLabel })
      expect(wrapper.find('a').attributes()['aria-label']).toEqual(ariaLabel)
    })
    it('allows to overwrite the link description', async () => {
      const { wrapper } = getWrapper()
      const description = 'some lengthy description'
      await wrapper.setProps({ description })
      expect(wrapper.find('#oc-feedback-link-description').text()).toEqual(description)
    })
  })
})

const getWrapper = () => {
  return {
    wrapper: mount(FeedbackLink, {
      global: {
        plugins: [...defaultPlugins()]
      }
    })
  }
}
