import SidebarNavItem from '../../../../src/components/SidebarNav/SidebarNavItem.vue'
import sidebarNavItemFixtures from '../../../__fixtures__/sidebarNavItems'
import { defaultPlugins, mount } from '@ownclouders/web-test-helpers'

const exampleNavItem = sidebarNavItemFixtures[0]

const propsData = {
  name: exampleNavItem.name,
  active: false,
  target: exampleNavItem.route.path,
  icon: exampleNavItem.icon,
  index: '5',
  id: '123'
}

describe('OcSidebarNav', () => {
  it('renders navItem without toolTip if expanded', () => {
    const { wrapper } = getWrapper(false)
    expect(wrapper.html()).toMatchSnapshot()
  })

  it('renders navItem with toolTip if collapsed', () => {
    const { wrapper } = getWrapper(true)
    expect(wrapper.html()).toMatchSnapshot()
  })
})

function getWrapper(collapsed: boolean) {
  return {
    wrapper: mount(SidebarNavItem, {
      props: {
        ...propsData,
        collapsed
      },
      global: {
        plugins: [...defaultPlugins()],
        stubs: { 'router-link': true }
      }
    })
  }
}
