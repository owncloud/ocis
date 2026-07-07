import LayoutContainer from '../../src/LayoutContainer.vue'
import {
  defaultComponentMocks,
  defaultPlugins,
  mount,
  nextTicks
} from '@ownclouders/web-test-helpers'
import { useAppsStore } from '../../src/piniaStores'

const selectors = {
  loadingSpinner: '#app-loading-spinner',
  routerView: '[data-testid="app-store-router-view"]'
}

describe('LayoutContainer', () => {
  it('shows a loading spinner while apps are loading from repositories', () => {
    const { wrapper } = getWrapper({ keepLoadingApps: true })
    expect(wrapper.find(selectors.loadingSpinner).exists()).toBeTruthy()
    expect(wrapper.find(selectors.routerView).exists()).toBeFalsy()
  })
  it('renders the router view when loading apps is done', async () => {
    const { wrapper } = getWrapper({})
    await nextTicks(2)
    expect(wrapper.find(selectors.loadingSpinner).exists()).toBeFalsy()
    expect(wrapper.find(selectors.routerView).exists()).toBeTruthy()
  })
})

function getWrapper({ keepLoadingApps }: { keepLoadingApps?: boolean }) {
  const plugins = defaultPlugins({})

  const { loadApps } = useAppsStore()
  vi.mocked(loadApps).mockReturnValue(
    new Promise((res) => {
      if (!keepLoadingApps) {
        return res()
      }
      return setTimeout(() => res(), 500)
    })
  )

  const mocks = {
    ...defaultComponentMocks(),
    loadApps
  }

  return {
    mocks,
    wrapper: mount(LayoutContainer, {
      global: {
        mocks,
        provide: mocks,
        plugins
      }
    })
  }
}
