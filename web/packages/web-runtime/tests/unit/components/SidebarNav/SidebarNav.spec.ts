import SidebarNav from '../../../../src/components/SidebarNav/SidebarNav.vue'
import sidebarNavItemFixtures from '../../../__fixtures__/sidebarNavItems'
import { defaultPlugins, mount } from '@ownclouders/web-test-helpers'

vi.mock('uuid', () => ({
  v4: () => {
    return '00000000-0000-0000-0000-000000000000'
  }
}))

const slots = {
  bottom: '<span class="footer">Footer</span>'
}

describe('OcSidebarNav', () => {
  it('displays a bottom slot if given', () => {
    const { wrapper } = getWrapper()
    expect(wrapper.findAll('.footer').length).toBe(1)
  })
  it('renders navItems into a list', () => {
    const { wrapper } = getWrapper()
    expect(wrapper.html()).toMatchSnapshot()
  })
  it('expands the navbar in open state', () => {
    const { wrapper } = getWrapper({ closed: false })
    expect(wrapper.find('.oc-app-navigation-expanded').exists).toBeTruthy()
  })
  it('collapses the navbar in closed state', () => {
    const { wrapper } = getWrapper({ closed: true })
    expect(wrapper.find('.oc-app-navigation-collapsed').exists).toBeTruthy()
  })
  it('emits "update:nav-bar-closed" upon button click', async () => {
    const { wrapper } = getWrapper()
    await wrapper.find('.toggle-sidebar-button').trigger('click')
    expect(wrapper.emitted('update:nav-bar-closed').length).toBeGreaterThan(0)
  })
  it('initially sets the highlighter to the active nav item', async () => {
    const { wrapper } = getWrapper()
    await wrapper.vm.$nextTick()
    expect(wrapper.vm.highlighterAttrs).toEqual({
      style: {
        transform: 'translateY(0px)',
        'transition-duration': '0.2s'
      }
    })
  })
})

function getWrapper({ closed = false } = {}) {
  return {
    wrapper: mount(SidebarNav, {
      slots,
      props: {
        navItems: sidebarNavItemFixtures,
        closed
      },
      global: {
        renderStubDefaultSlot: true,
        plugins: [...defaultPlugins()],
        stubs: { SidebarNavItem: true }
      }
    })
  }
}
