import { defaultPlugins, mount } from '@ownclouders/web-test-helpers'
import OcFilterChip from './OcFilterChip.vue'

const selectors = {
  filterChipBtn: '.oc-filter-chip-button',
  filterChipLabel: '.oc-filter-chip-label',
  filterChipDrop: '.oc-filter-chip-drop',
  filterChipBtnActive: '.oc-filter-chip-button-selected',
  clearBtn: '.oc-filter-chip-clear'
}

describe('OcFilterChip', () => {
  it('renders the filterLabel when no selected item has been given', () => {
    const filterLabel = 'Users'
    const { wrapper } = getWrapper({ props: { filterLabel } })
    expect(wrapper.find(selectors.filterChipLabel).text()).toEqual(filterLabel)
  })
  it('renders the first selected item name when selected items have been given', () => {
    const selectedItemNames = ['Einstein', 'Marie']
    const { wrapper } = getWrapper({ props: { selectedItemNames } })
    expect(wrapper.find(selectors.filterChipLabel).text()).toEqual(selectedItemNames[0])
  })
  it('emits the "clearFilter"-event when clicking the clear-button', async () => {
    const { wrapper } = getWrapper({ props: { selectedItemNames: ['Einstein', 'Marie'] } })
    await wrapper.find(selectors.clearBtn).trigger('click')
    expect(wrapper.emitted('clearFilter')).toBeTruthy()
  })
  describe('isToggle is true', () => {
    it('does not render a dropdown', () => {
      const { wrapper } = getWrapper({ props: { isToggle: true } })
      expect(wrapper.find(selectors.filterChipDrop).exists()).toBeFalsy()
    })
    it('marks filter as active when isToggleActive is true', () => {
      const { wrapper } = getWrapper({ props: { isToggle: true, isToggleActive: true } })
      expect(wrapper.find(selectors.filterChipBtnActive).exists()).toBeTruthy()
    })
    it('emits the "toggleFilter"-event when clicking the button', async () => {
      const { wrapper } = getWrapper({ props: { isToggle: true } })
      await wrapper.find(selectors.filterChipBtn).trigger('click')
      expect(wrapper.emitted('toggleFilter')).toBeTruthy()
    })
  })
})

const getWrapper = ({ props = {} } = {}) => {
  return {
    wrapper: mount(OcFilterChip, {
      props: { filterLabel: 'Users', selectedItemNames: [], ...props },
      global: { plugins: [...defaultPlugins()] }
    })
  }
}
