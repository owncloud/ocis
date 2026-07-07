import ApplicationsMenu from '../../../../src/components/Topbar/ApplicationsMenu.vue'
import {
  RouteLocation,
  defaultComponentMocks,
  defaultPlugins,
  shallowMount
} from '@ownclouders/web-test-helpers'
import { mock } from 'vitest-mock-extended'
import { AppMenuItemExtension } from '@ownclouders/web-pkg'

describe('ApplicationsMenu component', () => {
  describe('type', () => {
    it('renders internal button menu items', () => {
      const menuItem = { label: () => '', handler: () => undefined } as AppMenuItemExtension
      const { wrapper } = getWrapper({ menuItems: [menuItem] })
      const menuItemType = wrapper.find('oc-list-stub oc-button-stub').attributes('type')
      expect(menuItemType).toEqual('button')
    })
    it('renders internal link menu items', () => {
      const menuItem = { label: () => '', path: '/files' } as AppMenuItemExtension
      const { wrapper } = getWrapper({ menuItems: [menuItem] })
      const menuItemType = wrapper.find('oc-list-stub oc-button-stub').attributes('type')
      expect(menuItemType).toEqual('router-link')
    })
    it('renders external menu items', () => {
      const menuItem = { label: () => '', url: 'foo.bar' } as AppMenuItemExtension
      const { wrapper } = getWrapper({ menuItems: [menuItem] })
      const menuItemType = wrapper.find('oc-list-stub oc-button-stub').attributes('type')
      expect(menuItemType).toEqual('a')
    })
  })
  it('correctly sorts menu items by priority', () => {
    const menuItems = [
      { label: () => '1', priority: 50 },
      { label: () => '2', priority: 40 }
    ] as AppMenuItemExtension[]

    const { wrapper } = getWrapper({ menuItems })
    const firstElTxt = wrapper.findAll('oc-list-stub oc-button-stub')[0].find('span').text()
    expect(firstElTxt).toEqual(menuItems[1].label())
  })
  describe('active state', () => {
    it('checks the current route path against the menu item path', () => {
      const menuItems = [
        { label: () => '1', path: '/1' },
        { label: () => '2', path: '/2' }
      ] as AppMenuItemExtension[]

      const { wrapper } = getWrapper({ menuItems, path: '/2' })
      const activeElTxt = wrapper.find('.router-link-active span').text()
      expect(activeElTxt).toEqual(menuItems[1].label())
    })
  })
})

function getWrapper({
  menuItems = [],
  path = '/'
}: { menuItems?: AppMenuItemExtension[]; path?: string } = {}) {
  const mocks = {
    ...defaultComponentMocks({
      currentRoute: mock<RouteLocation>({ path })
    })
  }

  return {
    wrapper: shallowMount(ApplicationsMenu, {
      props: {
        menuItems
      },
      global: {
        renderStubDefaultSlot: true,
        mocks,
        provide: mocks,
        plugins: [...defaultPlugins()]
      }
    })
  }
}
