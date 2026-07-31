import ContextMenuQuickAction from '../../../../src/components/ContextActions/ContextMenuQuickAction.vue'
import { defaultPlugins, mount } from '@ownclouders/web-test-helpers'

describe('ContextMenuQuickAction component', () => {
  it('renders component', () => {
    const { wrapper } = getWrapper()
    expect(wrapper.html()).toMatchSnapshot()
  })
  it('triggers the "quickActionClicked"-event on click', async () => {
    const { wrapper } = getWrapper()
    await wrapper.find('button').trigger('click')
    expect(wrapper.emitted('quickActionClicked')).toBeTruthy()
  })
})

function getWrapper({ item = { id: '1' } } = {}) {
  return {
    wrapper: mount(ContextMenuQuickAction, {
      props: {
        item
      },
      global: {
        plugins: [...defaultPlugins()]
      }
    })
  }
}
