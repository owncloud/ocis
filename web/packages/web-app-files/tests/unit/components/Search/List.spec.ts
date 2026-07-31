import { shallowMount } from '@vue/test-utils'
import { merge } from 'lodash-es'
import { ResourceTable } from '@ownclouders/web-pkg'
import List from '../../../../src/components/Search/List.vue'
import { useResourcesViewDefaults } from '../../../../src/composables'
import { useResourcesViewDefaultsMock } from '../../../../tests/mocks/useResourcesViewDefaultsMock'
import { createRouter, createMemoryHistory } from 'vue-router'

import { defaultComponentMocks, defaultPlugins } from '@ownclouders/web-test-helpers'
import { AppBar, ItemFilter, queryItemAsString, useResourcesStore } from '@ownclouders/web-pkg'
import { ref } from 'vue'
import { Resource } from '@ownclouders/web-client'
import { mock } from 'vitest-mock-extended'
import { Capabilities } from '@ownclouders/web-client/ocs'

vi.mock('../../../../src/composables')
vi.mock('@ownclouders/web-pkg', async (importOriginal) => ({
  ...(await importOriginal<any>()),
  queryItemAsString: vi.fn(),
  useAppDefaults: vi.fn(),
  useFileActions: vi.fn(() => ({
    triggerDefaultAction: vi.fn()
  }))
}))

const selectors = {
  noContentMessageStub: 'no-content-message-stub',
  resourceTableStub: 'resource-table-stub',
  tagFilter: '.files-search-filter-tags',
  lastModifiedFilter: '.files-search-filter-last-modified',
  titleOnlyFilter: '.files-search-filter-title-only',
  filter: '.files-search-result-filter'
}

describe('List component', () => {
  it('should render no-content-message if no resources found', () => {
    const { wrapper } = getWrapper()
    expect(wrapper.find(selectors.noContentMessageStub).exists()).toBeTruthy()
  })
  it('should render resource table if resources found', () => {
    const { wrapper } = getWrapper({ resources: [mock<Resource>()] })
    expect(wrapper.find(selectors.resourceTableStub).exists()).toBeTruthy()
  })
  it('resets the initial store file state', () => {
    getWrapper({ resources: [mock<Resource>()] })

    const { clearResourceList } = useResourcesStore()
    expect(clearResourceList).toHaveBeenCalled()
  })
  it('should emit search event on mount', async () => {
    const { wrapper } = getWrapper()
    await (wrapper.vm as any).loadAvailableTagsTask.last
    expect(wrapper.emitted('search').length).toBeGreaterThan(0)
  })
  it('should emit search only if one of the queries changes', async () => {
    const $router = createRouter({
      routes: [{ name: 'files-common-search', path: '/files-common-search', redirect: null }],
      history: createMemoryHistory()
    })

    const { wrapper } = getWrapper({
      mocks: { $router }
    })

    let replacementCounter = 0
    const queries = [{}, {}, {}, {}, {}, {}, {}, {}, { q_tags: 'foo' }, {}, {}, {}, {}]
    for (const query of queries) {
      await $router.replace({
        name: 'files-common-search',
        query: merge(
          {
            q_titleOnly: 'q_titleOnly',
            q_tags: 'q_tags',
            q_lastModified: 'q_lastModified',
            useScope: 'useScope'
          },
          query
        )
      })

      replacementCounter++
    }

    await (wrapper.vm as any).loadAvailableTagsTask.last
    expect(replacementCounter).toBe(queries.length)
    expect(wrapper.emitted('search').length).toBe(4)
  })
  describe('breadcrumbs', () => {
    it('show "Search" when no search term given', () => {
      const { wrapper } = getWrapper()
      const appBar = wrapper.findComponent<typeof AppBar>('app-bar-stub')
      expect(appBar.props('breadcrumbs')[0].text).toEqual('Search')
    })
    it('include the search term if given', () => {
      const searchTerm = 'term'
      const { wrapper } = getWrapper({ searchTerm })
      const appBar = wrapper.findComponent<typeof AppBar>('app-bar-stub')
      expect(appBar.props('breadcrumbs')[0].text).toEqual(`Search results for "${searchTerm}"`)
    })
  })
  describe('filter', () => {
    describe('general', () => {
      it('should not be rendered if no filtering is available', async () => {
        const { wrapper } = getWrapper({ fullTextSearchEnabled: false, availableTags: [] })
        await (wrapper.vm as any).loadAvailableTagsTask.last
        expect(wrapper.find(selectors.filter).exists()).toBeFalsy()
      })
    })
    describe('tags', () => {
      it('should show all available tags', async () => {
        const tag = 'tag1'
        const { wrapper } = getWrapper({ availableTags: [tag] })
        await (wrapper.vm as any).loadAvailableTagsTask.last
        expect(wrapper.find(selectors.tagFilter).exists()).toBeTruthy()
        expect(
          wrapper.findComponent<typeof ItemFilter>(selectors.tagFilter).props('items')
        ).toEqual([{ label: tag, id: tag }])
      })
      it('should set initial filter when tags are given via query param', async () => {
        const searchTerm = 'term'
        const availableTags = ['tag1', 'tag2']
        const { wrapper } = getWrapper({
          availableTags,
          searchTerm,
          tagFilterQuery: availableTags.join('+')
        })
        await (wrapper.vm as any).loadAvailableTagsTask.last
        expect(wrapper.emitted('search')[0][0]).toEqual(
          `(name:"*${searchTerm}*" OR content:"${searchTerm}") AND tag:("${availableTags[0]}" OR "${availableTags[1]}")`
        )
      })
    })

    describe('last modified', () => {
      it('should show available last modified values', async () => {
        const expectation = [
          { label: 'today', id: 'today' },
          { label: 'yesterday', id: 'yesterday' },
          { label: 'this week', id: 'this week' },
          { label: 'last week', id: 'last week' },
          { label: 'last 7 days', id: 'last 7 days' },
          { label: 'this month', id: 'this month' },
          { label: 'last month', id: 'last month' },
          { label: 'last 30 days', id: 'last 30 days' },
          { label: 'this year', id: 'this year' },
          { label: 'last year', id: 'last year' }
        ]
        const lastModifiedValues = {
          keywords: [
            'today',
            'yesterday',
            'this week',
            'last week',
            'last 7 days',
            'this month',
            'last month',
            'last 30 days',
            'this year',
            'last year'
          ]
        }
        const { wrapper } = getWrapper({
          availableLastModifiedValues: lastModifiedValues,
          availableTags: ['tag']
        })
        await (wrapper.vm as any).loadAvailableTagsTask.last

        expect(wrapper.find(selectors.lastModifiedFilter).exists()).toBeTruthy()
        expect(
          wrapper.findComponent<typeof ItemFilter>(selectors.lastModifiedFilter).props('items')
        ).toEqual(expectation)
      })
      it('should set initial filter when last modified is given via query param', async () => {
        const searchTerm = 'Screenshot'
        const lastModifiedFilterQuery = 'today'
        const { wrapper } = getWrapper({
          searchTerm,
          lastModifiedFilterQuery
        })
        await (wrapper.vm as any).loadAvailableTagsTask.last
        expect(wrapper.emitted('search')[0][0]).toEqual(
          `(name:"*${searchTerm}*" OR content:"${searchTerm}") AND mtime:${lastModifiedFilterQuery}`
        )
      })
    })

    describe('titleOnly', () => {
      it('should render filter if enabled via capabilities', () => {
        const { wrapper } = getWrapper({ fullTextSearchEnabled: true })
        expect(wrapper.find(selectors.titleOnlyFilter).exists()).toBeTruthy()
      })
      it('should not render filter if not enabled via capabilities', () => {
        const { wrapper } = getWrapper({ fullTextSearchEnabled: false })
        expect(wrapper.find(selectors.titleOnlyFilter).exists()).toBeFalsy()
      })
      it('should set initial filter when titleOnly is set active via query param', async () => {
        const searchTerm = 'term'
        const { wrapper } = getWrapper({
          searchTerm,
          titleOnlyFilterQuery: 'true',
          fullTextSearchEnabled: true
        })
        await (wrapper.vm as any).loadAvailableTagsTask.last
        expect(wrapper.emitted('search')[0][0]).toEqual(`name:"*${searchTerm}*"`)
      })
    })
  })
})

function getWrapper({
  availableTags = [],
  resources = [],
  searchTerm = '',
  tagFilterQuery = null,
  titleOnlyFilterQuery = null,
  fullTextSearchEnabled = true,
  availableLastModifiedValues = {},
  lastModifiedFilterQuery = null,
  mocks = {}
}: {
  availableTags?: string[]
  resources?: Resource[]
  searchTerm?: string
  tagFilterQuery?: string
  titleOnlyFilterQuery?: string
  fullTextSearchEnabled?: boolean
  availableLastModifiedValues?: Record<string, string[]>
  lastModifiedFilterQuery?: string
  mocks?: Record<string, unknown>
} = {}) {
  vi.mocked(queryItemAsString).mockImplementationOnce(() => searchTerm)
  vi.mocked(queryItemAsString).mockImplementationOnce(() => titleOnlyFilterQuery)
  vi.mocked(queryItemAsString).mockImplementationOnce(() => tagFilterQuery)
  vi.mocked(queryItemAsString).mockImplementationOnce(() => lastModifiedFilterQuery)

  const resourcesViewDetailsMock = useResourcesViewDefaultsMock({
    paginatedResources: ref(resources)
  })
  vi.mocked(useResourcesViewDefaults).mockImplementation(() => resourcesViewDetailsMock)

  const localMocks = {
    ...defaultComponentMocks(),
    ...mocks
  }
  localMocks.$clientService.graphAuthenticated.tags.listTags.mockResolvedValue(availableTags)

  const capabilities = {
    files: { tags: true },
    search: {
      property: {
        mtime: availableLastModifiedValues,
        content: { enabled: fullTextSearchEnabled },
        tags: { enabled: true }
      }
    }
  } satisfies Partial<Capabilities['capabilities']>

  return {
    mocks: localMocks,
    wrapper: shallowMount(List, {
      global: {
        components: {
          ResourceTable,
          AppBar
        },
        mocks: localMocks,
        provide: localMocks,
        stubs: {
          FilesViewWrapper: false
        },
        plugins: [...defaultPlugins({ piniaOptions: { capabilityState: { capabilities } } })]
      }
    })
  }
}
