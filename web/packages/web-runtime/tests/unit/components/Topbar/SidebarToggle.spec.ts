import SidebarToggle from '../../../../src/components/Topbar/SideBarToggle.vue'
import { eventBus } from '@ownclouders/web-pkg/src/services'
import { defaultPlugins, mount, defaultComponentMocks } from '@ownclouders/web-test-helpers'

const selectors = {
  toggleSidebarBtn: '#files-toggle-sidebar'
}

describe('SidebarToggle component', () => {
  it.each([true, false])(
    'should show the "Toggle sidebar"-button with sidebar opened and closed',
    (isSideBarOpen) => {
      const { wrapper } = getWrapper({ isSideBarOpen })
      expect(wrapper.find(selectors.toggleSidebarBtn).exists()).toBeTruthy()
      expect(wrapper.html()).toMatchSnapshot()
    }
  )
  it('publishes the toggle-event to the sidebar on click', async () => {
    const { wrapper } = getWrapper()
    const eventSpy = vi.spyOn(eventBus, 'publish')
    await wrapper.find(selectors.toggleSidebarBtn).trigger('click')
    expect(eventSpy).toHaveBeenCalled()
  })
})

function getWrapper({ isSideBarOpen = false } = {}) {
  const mocks = defaultComponentMocks()
  return {
    mocks,
    wrapper: mount(SidebarToggle, {
      props: { isSideBarOpen },
      global: {
        mocks,
        plugins: [...defaultPlugins()]
      }
    })
  }
}
