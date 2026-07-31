import ContextActionMenu from '../../../../src/components/ContextActions/ContextActionMenu.vue'
import { Action } from '../../../../src/composables/actions'
import { defaultPlugins, shallowMount } from '@ownclouders/web-test-helpers'

describe('ContextActionMenu component', () => {
  it('renders the menu with actions', () => {
    const menuSections = [
      { name: 'action 1', items: [] as Action[] },
      { name: 'action 2', items: [] as Action[] }
    ]
    const { wrapper } = getShallowWrapper(menuSections)
    expect(wrapper.html()).toMatchSnapshot()
    expect(wrapper.find('.oc-files-context-actions').exists()).toBeTruthy()
    expect(wrapper.findAll('.oc-files-context-actions').length).toEqual(menuSections.length)
  })
})

function getShallowWrapper(menuSections: { name: string; items: Action[] }[]) {
  return {
    wrapper: shallowMount(ContextActionMenu, {
      props: {
        menuSections,
        actionOptions: { resources: [] }
      },
      global: {
        plugins: [...defaultPlugins()]
      }
    })
  }
}
