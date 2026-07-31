import StatusIndicators from './OcStatusIndicators.vue'
import { defaultPlugins, mount } from '@ownclouders/web-test-helpers'

const fileResource = {
  name: 'forest.jpg',
  path: 'nature/forest.jpg',
  thumbnail: 'https://cdn.pixabay.com/photo/2015/09/09/16/05/forest-931706_960_720.jpg',
  indicators: [] as unknown[],
  type: 'file',
  isFolder: false,
  extension: 'jpg'
}
const indicator = {
  id: 'testid',
  label: 'testlabel',
  type: 'testtype',
  icon: 'icon',
  handler: vi.fn()
}
describe('OcStatusIndicators', () => {
  it('does call indicator handler on click', () => {
    const spyHandler = vi.spyOn(indicator, 'handler')
    const wrapper = mount(StatusIndicators, {
      props: {
        resource: fileResource,
        indicators: [indicator],
        target: 'test'
      },
      global: {
        plugins: [...defaultPlugins()]
      }
    })
    wrapper.find('.oc-status-indicators-indicator').trigger('click')
    expect(spyHandler).toHaveBeenCalled()
  })
  it('does create indicator with id', () => {
    const wrapper = mount(StatusIndicators, {
      props: {
        resource: fileResource,
        indicators: [indicator],
        target: 'test'
      },
      global: {
        plugins: [...defaultPlugins()]
      }
    })
    expect(wrapper.find(`#${indicator.id}`).exists()).toBeTruthy()
  })
  it('does not render a button if disableHandler is set', () => {
    const wrapper = mount(StatusIndicators, {
      props: {
        resource: fileResource,
        indicators: [indicator],
        disableHandler: true
      },
      global: {
        plugins: [...defaultPlugins()]
      }
    })
    expect(wrapper.find('button').exists()).toBeFalsy()
  })
})
