import { SDKSearch } from '../../../src/search'
import { RouteLocation, Router } from 'vue-router'
import { mock } from 'vitest-mock-extended'
import { computed } from 'vue'
import { createTestingPinia, writable } from '@ownclouders/web-test-helpers'
import { ConfigStore, useCapabilityStore } from '@ownclouders/web-pkg'

const getStore = (reports: string[] = []) => {
  createTestingPinia({
    initialState: { capabilities: { capabilities: { dav: { reports } } } }
  })
  return useCapabilityStore()
}

describe('SDKProvider', () => {
  it('is only available if announced via capabilities', () => {
    const search = new SDKSearch(getStore(), mock<Router>(), vi.fn(), mock<ConfigStore>())
    expect(search.available).toBe(false)
  })

  describe('SDKProvider previewSearch', () => {
    it('is not available on certain routes', () => {
      ;[
        { route: 'foo', available: true },
        { route: 'search-provider-list' },
        { route: 'bar', available: true }
      ].forEach((v) => {
        const router = mock<Router>()
        writable(router).currentRoute = computed(() => mock<RouteLocation>({ name: v.route }))

        const search = new SDKSearch(
          getStore(['search-files']),
          router,
          vi.fn(),
          mock<ConfigStore>()
        )

        expect(!!search.previewSearch.available).toBe(!!v.available)
      })
    })
  })
})
