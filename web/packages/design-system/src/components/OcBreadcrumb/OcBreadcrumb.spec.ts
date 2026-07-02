import { defaultPlugins, shallowMount } from '@ownclouders/web-test-helpers'
import Breadcrumb from './OcBreadcrumb.vue'

const items = [
  { text: 'First folder', to: { path: 'folder' } },
  { text: 'Subfolder', onClick: () => alert('Breadcrumb clicked!') },
  { text: 'Deep', to: { path: 'folder' } },
  { text: 'Deeper ellipsize in responsive mode' }
]

describe('OcBreadcrumb', () => {
  it('sets correct variation', () => {
    const { wrapper } = getWrapper({ variation: 'lead' })
    expect(wrapper.props().variation).toMatch('lead')
    expect(wrapper.find('.oc-breadcrumb').attributes('class')).toContain('oc-breadcrumb-lead')
    expect(wrapper.html()).toMatchSnapshot()
  })
  it('displays all items', () => {
    const { wrapper } = getWrapper()
    expect(wrapper.findAll('.oc-breadcrumb-list-item:not(.oc-invisible-sr)').length).toBe(
      items.length
    )
    expect(wrapper.html()).toMatchSnapshot()
  })
  it('displays context menu trigger if enabled via property', () => {
    const { wrapper } = getWrapper({ showContextActions: true })
    expect(wrapper.find('#oc-breadcrumb-contextmenu-trigger').exists()).toBe(true)
  })
  it('does not display context menu trigger if not enabled via property', () => {
    const { wrapper } = getWrapper({ showContextActions: false })
    expect(wrapper.find('#oc-breadcrumb-contextmenu-trigger').exists()).toBe(false)
  })
  describe('mobile navigation', () => {
    it.each([
      { items: [], shows: false },
      { items: [items[0]], shows: false },
      { items: [items[0], items[1]], shows: true }
    ])('shows if more than 1 breadcrumb item is given', ({ items, shows }) => {
      const { wrapper } = getWrapper({ items })
      expect(wrapper.find('.oc-breadcrumb-mobile-navigation').exists()).toBe(shows)
    })
  })
  describe('mobile current folder', () => {
    it.each([
      { items: [], shows: false },
      { items: [items[0]], shows: false },
      { items: [items[0], items[1]], shows: true }
    ])('shows if more than 1 breadcrumb item is given', ({ items, shows }) => {
      const { wrapper } = getWrapper({ items })
      expect(wrapper.find('.oc-breadcrumb-mobile-current').exists()).toBe(shows)
    })
  })
})

const getWrapper = (props = {}) => {
  return {
    wrapper: shallowMount(Breadcrumb, {
      props: {
        items,
        ...props
      },
      global: { renderStubDefaultSlot: true, plugins: [...defaultPlugins()] }
    })
  }
}
