import DateFilter from '../../../../src/components/Filters/DateFilter.vue'
import {
  PartialComponentProps,
  defaultComponentMocks,
  defaultPlugins,
  mount
} from '@ownclouders/web-test-helpers'
import { queryItemAsString } from '../../../../src/composables/appDefaults'
import { DateTime } from 'luxon'

vi.mock('../../../../src/composables/appDefaults')

const filterItems = [
  { id: '1', name: 'today' },
  { id: '2', name: 'yesterday' }
]

const selectors = {
  filterListItem: '.date-filter-list-item',
  filterChipLabel: '.oc-filter-chip-label',
  activeFilterListItemSpan: '.date-filter-list-item-active .oc-text-truncate span',
  customDateRangeBtn: '[data-testid="custom-date-range"]',
  customDateRangePanel: '.date-filter-range-panel',
  customDateRangeApplyBtn: '.date-filter-apply-btn button',
  customDateRangeBackBtn: '.date-filter-range-panel-back'
}

describe('DateFilter', () => {
  it('renders all items', () => {
    const { wrapper } = getWrapper()
    expect(wrapper.html()).toMatchSnapshot()
  })
  it('can use a custom attribute as display name', () => {
    const filterItems = [
      { id: '1', displayName: 'today' },
      { id: '2', displayName: 'yesterday' }
    ]
    const { wrapper } = getWrapper({
      props: { displayNameAttribute: 'displayName', items: filterItems }
    })
    expect(wrapper.html()).toMatchSnapshot()
  })

  describe('route query', () => {
    it('sets the selected item as route query param', async () => {
      const { wrapper, mocks } = getWrapper()
      const item = wrapper.findAll(selectors.filterListItem).at(0)
      expect(mocks.$router.push).not.toHaveBeenCalled()
      await item.trigger('click')
      expect(mocks.$router.push).toHaveBeenCalledWith(
        expect.objectContaining({
          query: expect.objectContaining({ [(wrapper.vm as any).queryParam]: '1' })
        })
      )
    })
    it('sets the selected items initially when given via query param', () => {
      const { wrapper } = getWrapper({ initialQuery: '1' })
      expect((wrapper.vm as any).selectedItem).toEqual(filterItems[0])
    })
  })

  describe('custom date range', () => {
    it('shows the option and panel for picking a custom date range', () => {
      const { wrapper } = getWrapper()
      expect(wrapper.find(selectors.customDateRangeBtn).exists()).toBeTruthy()
      expect(wrapper.find(selectors.customDateRangePanel).exists()).toBeTruthy()
    })
    it('shows the custom date range as filter chip label when selected', async () => {
      const fromDate = DateTime.now().minus({ days: 1 })
      const toDate = DateTime.now()
      const { wrapper } = getWrapper({
        initialQuery: `range:${fromDate.toMillis()} - ${toDate.toMillis()}`
      })
      await wrapper.vm.$nextTick()
      const filterChipLabel = wrapper.find(selectors.filterChipLabel)
      expect(filterChipLabel.text()).toContain(
        `${fromDate.setLocale('en').toLocaleString()} - ${toDate.setLocale('en').toLocaleString()}`
      )
    })
    it('correctly marks the "Custom date range"-option as selected', async () => {
      const fromDate = DateTime.now().minus({ days: 1 })
      const toDate = DateTime.now()
      const { wrapper } = getWrapper({
        initialQuery: `range:${fromDate.toMillis()} - ${toDate.toMillis()}`
      })
      await wrapper.vm.$nextTick()
      expect(wrapper.find(selectors.activeFilterListItemSpan).text()).toEqual('Custom date range')
    })
    it('back button should close the custom date range panel', async () => {
      const { wrapper } = getWrapper()
      ;(wrapper.vm as any).dateRangeClicked = true
      await wrapper.vm.$nextTick()
      const backBtn = wrapper.find(selectors.customDateRangeBackBtn)
      await backBtn.trigger('click')
      expect((wrapper.vm as any).dateRangeClicked).toBeFalsy()
    })
    describe('apply button', () => {
      it('is not clickable without dates entered', async () => {
        console.error = vi.fn()
        const { wrapper } = getWrapper()
        const applyBtn = wrapper.find(selectors.customDateRangeApplyBtn)
        await applyBtn.trigger('click')
        expect(wrapper.emitted('selectionChange')).not.toBeDefined()
      })
      it('is not clickable when from date is after to date', async () => {
        const { wrapper } = getWrapper()
        ;(wrapper.vm as any).fromDate = DateTime.now().plus({ days: 1 })
        ;(wrapper.vm as any).toDate = DateTime.now()
        await wrapper.vm.$nextTick()
        const applyBtn = wrapper.find(selectors.customDateRangeApplyBtn)
        await applyBtn.trigger('click')
        expect(wrapper.emitted('selectionChange')).not.toBeDefined()
      })
      it('emits a selection change on click with today entered as date', async () => {
        const { wrapper } = getWrapper()
        ;(wrapper.vm as any).fromDate = DateTime.now()
        ;(wrapper.vm as any).toDate = DateTime.now()
        await wrapper.vm.$nextTick()
        const applyBtn = wrapper.find(selectors.customDateRangeApplyBtn)
        await applyBtn.trigger('click')
        expect(wrapper.emitted('selectionChange')).toBeDefined()
      })
    })
  })
})

function getWrapper({
  props = {},
  initialQuery = ''
}: { props?: PartialComponentProps<typeof DateFilter>; initialQuery?: string } = {}) {
  vi.mocked(queryItemAsString).mockImplementation(() => initialQuery)
  const mocks = defaultComponentMocks()

  return {
    mocks,
    wrapper: mount(DateFilter, {
      props: {
        filterLabel: 'Users',
        filterName: 'users',
        items: filterItems,
        ...props
      },
      slots: {
        item(data) {
          return props.displayNameAttribute ? data.item[props.displayNameAttribute] : data.item.name
        }
      },
      global: {
        plugins: [...defaultPlugins()],
        mocks,
        provide: mocks,
        stubs: { OcCheckbox: true }
      }
    })
  }
}
