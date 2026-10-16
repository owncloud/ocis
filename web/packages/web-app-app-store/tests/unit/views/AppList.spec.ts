import { RouteLocationNormalizedLoaded } from 'vue-router'
import { setActivePinia } from 'pinia'
import {
  createTestingPinia,
  defaultComponentMocks,
  defaultPlugins,
  flushPromises,
  shallowMount,
  VueWrapper
} from '@ownclouders/web-test-helpers'
import AppList from '../../../src/views/AppList.vue'
import { useAppsStore } from '../../../src/piniaStores'
import { App } from '../../../src/types'

const selectors = {
  noContentMessage: 'no-content-message-stub',
  appTile: 'app-tile-stub',
  filterInput: '#apps-filter'
}

const makeApp = (id: string, name: string): App => ({
  id,
  name,
  subtitle: `${name} subtitle`,
  license: 'MIT',
  versions: [{ version: '1.0.0', url: 'https://example.test/app.zip' }],
  authors: [{ name: 'ownCloud' }],
  tags: ['productivity'],
  screenshots: [],
  resources: [],
  repository: { name: 'repo', url: 'https://example.test/repo.json' },
  mostRecentVersion: { version: '1.0.0', url: 'https://example.test/app.zip' }
})

describe('AppList', () => {
  it('shows the no-content message when no apps match', () => {
    const { wrapper } = getWrapper({ apps: [] })
    expect(wrapper.find(selectors.noContentMessage).exists()).toBe(true)
    expect(wrapper.findAll(selectors.appTile)).toHaveLength(0)
  })

  it('renders a tile for each app', () => {
    const { wrapper } = getWrapper({ apps: [makeApp('alpha', 'Alpha'), makeApp('bravo', 'Bravo')] })
    expect(wrapper.findAll(selectors.appTile)).toHaveLength(2)
    expect(wrapper.find(selectors.noContentMessage).exists()).toBe(false)
  })

  it('reflects the search term into the filter route query on input', async () => {
    const { wrapper, mocks } = getWrapper({ apps: [makeApp('alpha', 'Alpha')] })
    await (wrapper.findComponent(selectors.filterInput) as VueWrapper).vm.$emit(
      'update:modelValue',
      'foo'
    )
    expect(mocks.$router.replace).toHaveBeenCalledWith({ query: { filter: 'foo' } })
  })

  it('narrows the rendered tiles to the ones matching the active filter query', async () => {
    const { wrapper } = getWrapper({
      apps: [makeApp('alpha', 'Alpha'), makeApp('bravo', 'Bravo')],
      filter: 'Alpha'
    })
    await flushPromises()
    const filteredApps = (wrapper.vm as unknown as { filteredApps: App[] }).filteredApps
    expect(filteredApps.map((app) => app.name)).toEqual(['Alpha'])
    expect(wrapper.findAll(selectors.appTile)).toHaveLength(1)
  })
})

function getWrapper({ apps, filter = '' }: { apps: App[]; filter?: string }) {
  const pinia = createTestingPinia({ stubActions: false })
  setActivePinia(pinia)
  const appsStore = useAppsStore()
  appsStore.apps = apps
  const plugins = [...defaultPlugins({ pinia: false }), pinia]

  const mocks = {
    ...defaultComponentMocks({
      currentRoute: {
        name: 'app-store',
        path: '/',
        params: {},
        query: { ...(filter && { filter }) },
        meta: {}
      } as unknown as RouteLocationNormalizedLoaded
    })
  }

  return {
    mocks,
    wrapper: shallowMount(AppList, {
      global: {
        plugins,
        mocks,
        provide: mocks,
        // render oc-list's default slot so the app tiles inside become visible stubs
        stubs: { 'oc-list': { template: '<ul><slot /></ul>' } }
      }
    })
  }
}
