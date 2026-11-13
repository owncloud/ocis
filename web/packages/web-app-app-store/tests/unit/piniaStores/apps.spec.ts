import { setActivePinia } from 'pinia'
import { createTestingPinia } from '@ownclouders/web-test-helpers'
import { useAppsStore } from '../../../src/piniaStores/apps'
import { useRepositoriesStore } from '../../../src/piniaStores/repositories'
import { AppStoreRepository } from '../../../src/types'

const repo1: AppStoreRepository = { name: 'repo1', url: 'https://example.test/repo1.json' }
const repo2: AppStoreRepository = { name: 'repo2', url: 'https://example.test/repo2.json' }

const makeRawApp = (id: string, name: string) => ({
  id,
  name,
  subtitle: `${name} subtitle`,
  license: 'MIT',
  versions: [{ version: '1.0.0', url: `https://example.test/${id}.zip` }],
  authors: [{ name: 'ownCloud' }],
  tags: ['productivity']
})

const stubFetch = (responsesByUrl: Record<string, unknown | 'reject'>) => {
  vi.stubGlobal(
    'fetch',
    vi.fn((url: string) => {
      const entry = responsesByUrl[url]
      if (entry === 'reject') {
        return Promise.reject(new Error('network error'))
      }
      return Promise.resolve({ json: () => Promise.resolve(entry) })
    })
  )
}

describe('useAppsStore', () => {
  beforeEach(() => {
    setActivePinia(createTestingPinia({ stubActions: false }))
  })

  afterEach(() => {
    vi.unstubAllGlobals()
  })

  describe('method "loadApps"', () => {
    it('loads and maps apps from the configured repositories', async () => {
      useRepositoriesStore().setRepositories([repo1])
      stubFetch({ [repo1.url]: { apps: [makeRawApp('alpha', 'Alpha')] } })

      const store = useAppsStore()
      await store.loadApps()

      expect(store.apps).toHaveLength(1)
      expect(store.apps[0].id).toBe('alpha')
      expect(store.apps[0].repository).toEqual(repo1)
      expect(store.apps[0].mostRecentVersion).toEqual({
        version: '1.0.0',
        url: 'https://example.test/alpha.zip'
      })
    })

    it('catches and logs a per-repository failure and returns no apps for that repository', async () => {
      const errorSpy = vi.spyOn(console, 'error').mockImplementation(() => undefined)
      useRepositoriesStore().setRepositories([repo1, repo2])
      stubFetch({
        [repo1.url]: { apps: [makeRawApp('alpha', 'Alpha')] },
        [repo2.url]: 'reject'
      })

      const store = useAppsStore()
      await store.loadApps()

      expect(errorSpy).toHaveBeenCalled()
      expect(store.apps).toHaveLength(1)
      expect(store.apps[0].id).toBe('alpha')
    })

    it('sorts the apps of a repository by name (case-insensitive)', async () => {
      useRepositoriesStore().setRepositories([repo1])
      stubFetch({
        [repo1.url]: {
          apps: [makeRawApp('c', 'Charlie'), makeRawApp('a', 'alpha'), makeRawApp('b', 'Bravo')]
        }
      })

      const store = useAppsStore()
      await store.loadApps()

      expect(store.apps.map((app) => app.name)).toEqual(['alpha', 'Bravo', 'Charlie'])
    })
  })

  describe('method "getById"', () => {
    it('returns the matching app or undefined', async () => {
      useRepositoriesStore().setRepositories([repo1])
      stubFetch({ [repo1.url]: { apps: [makeRawApp('alpha', 'Alpha')] } })

      const store = useAppsStore()
      await store.loadApps()

      expect(store.getById('alpha')?.name).toBe('Alpha')
      expect(store.getById('does-not-exist')).toBeUndefined()
    })
  })
})
