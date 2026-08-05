import { defineComponent, nextTick } from 'vue'
import { defaultPlugins, mount } from '@ownclouders/web-test-helpers'
import SideBar from '../../../../src/components/SideBar/SideBar.vue'
import { SideBarPanel } from '../../../../src/components/SideBar/types'

const panelComponent = defineComponent({ template: '<div>Panel content</div>' })
const rootPanel = {
  name: 'details',
  icon: 'information',
  title: () => 'Details',
  isVisible: () => true,
  isRoot: () => true,
  component: panelComponent
} as SideBarPanel<any, any, any>
const versionsPanel = {
  name: 'versions',
  icon: 'history',
  title: () => 'Versions',
  isVisible: () => true,
  component: panelComponent
} as SideBarPanel<any, any, any>

const waitForAnimationFrame = () =>
  new Promise<void>((resolve) => requestAnimationFrame(() => resolve()))

describe('SideBar', () => {
  it('moves focus into the active panel and restores it when closed', async () => {
    const trigger = document.createElement('button')
    document.body.append(trigger)
    trigger.focus()

    const wrapper = mount(SideBar, {
      attachTo: document.body,
      global: { plugins: defaultPlugins() },
      props: {
        isOpen: true,
        loading: false,
        availablePanels: [rootPanel],
        panelContext: {}
      }
    })

    await nextTick()
    await waitForAnimationFrame()
    expect(document.activeElement).toBe(
      wrapper.find('#sidebar-panel-details .header__title').element
    )

    wrapper.unmount()
    await nextTick()
    expect(document.activeElement).toBe(trigger)

    trigger.remove()
  })

  it('moves focus into the active panel when loading finishes', async () => {
    const wrapper = mount(SideBar, {
      attachTo: document.body,
      global: { plugins: defaultPlugins() },
      props: {
        isOpen: true,
        loading: true,
        availablePanels: [rootPanel],
        panelContext: {}
      }
    })

    await wrapper.setProps({ loading: false })
    await nextTick()
    await waitForAnimationFrame()
    expect(document.activeElement).toBe(
      wrapper.find('#sidebar-panel-details .header__title').element
    )

    wrapper.unmount()
  })

  it('does not move external focus into the sidebar after a background reload', async () => {
    const wrapper = mount(SideBar, {
      attachTo: document.body,
      global: { plugins: defaultPlugins() },
      props: {
        isOpen: true,
        loading: false,
        availablePanels: [rootPanel],
        panelContext: {}
      }
    })

    await nextTick()
    await waitForAnimationFrame()

    const filterOption = document.createElement('button')
    document.body.append(filterOption)
    filterOption.focus()

    await wrapper.setProps({ loading: true })
    await wrapper.setProps({ loading: false })
    await nextTick()
    await waitForAnimationFrame()

    expect(document.activeElement).toBe(filterOption)

    wrapper.unmount()
    filterOption.remove()
  })

  it('moves focus to the back button when the versions panel opens', async () => {
    const wrapper = mount(SideBar, {
      attachTo: document.body,
      global: { plugins: defaultPlugins() },
      props: {
        isOpen: true,
        loading: false,
        availablePanels: [rootPanel, versionsPanel],
        panelContext: {},
        onSelectPanel: (panel: string) => wrapper.setProps({ activePanel: panel })
      }
    })

    const versionsSelect = wrapper.find('#sidebar-panel-versions-select')
    ;(versionsSelect.element as HTMLElement).focus()
    expect(document.activeElement).toBe(versionsSelect.element)

    await versionsSelect.trigger('click')
    await nextTick()
    expect(document.activeElement).toBe(wrapper.find('.sidebar-focus-guard').element)

    await waitForAnimationFrame()

    expect(document.activeElement).toBe(
      wrapper.find('#sidebar-panel-versions .header__back').element
    )

    await wrapper.setProps({ loading: true })
    expect(document.activeElement).toBe(wrapper.find('.sidebar-focus-guard').element)

    wrapper.unmount()
  })
})
