import ItemFilterToggle from '../../../src/components/ItemFilterToggle.vue'
import { defaultComponentMocks, defaultPlugins, mount } from '@ownclouders/web-test-helpers'
import { queryItemAsString } from '../../../src/composables/appDefaults'

vi.mock('../../../src/composables/appDefaults', () => ({
  appDefaults: vi.fn(),
  queryItemAsString: vi.fn()
}))

const selectors = {
  labelSpan: '.oc-filter-chip-label',
  filterBtn: '.oc-filter-chip-button'
}

describe('ItemFilterToggle', () => {
  it('renders the toggle filter including its label', () => {
    const filterLabel = 'Toggle'
    const { wrapper } = getWrapper({ props: { filterLabel } })
    expect(wrapper.find(selectors.labelSpan).text()).toEqual(filterLabel)
  })
  it('emits the "toggleFilter"-event on click', async () => {
    const { wrapper } = getWrapper()
    await wrapper.find(selectors.filterBtn).trigger('click')
    expect(wrapper.emitted('toggleFilter').length).toBeGreaterThan(0)
  })
  describe('route query', () => {
    it('sets the active state as query param', async () => {
      const { wrapper, mocks } = getWrapper()
      expect(mocks.$router.push).not.toHaveBeenCalled()
      await wrapper.find(selectors.filterBtn).trigger('click')
      expect(mocks.$router.push).toHaveBeenCalledWith(
        expect.objectContaining({
          query: expect.objectContaining({ [(wrapper.vm as any).queryParam]: 'true' })
        })
      )
    })
    it('sets the active state initially when given via query param', () => {
      const { wrapper } = getWrapper({ initialQuery: 'true' })
      expect((wrapper.vm as any).filterActive).toEqual(true)
    })
  })
})

function getWrapper({ props = {}, initialQuery = '' } = {}) {
  vi.mocked(queryItemAsString).mockImplementation(() => initialQuery)
  const mocks = defaultComponentMocks()
  return {
    mocks,
    wrapper: mount(ItemFilterToggle, {
      props: {
        filterLabel: 'Toggle',
        filterName: 'toggle',
        ...props
      },
      global: {
        plugins: [...defaultPlugins()],
        mocks,
        provide: mocks
      }
    })
  }
}
