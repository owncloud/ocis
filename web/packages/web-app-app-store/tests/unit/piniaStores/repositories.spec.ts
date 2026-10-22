import { setActivePinia } from 'pinia'
import { createTestingPinia } from '@ownclouders/web-test-helpers'
import { useRepositoriesStore } from '../../../src/piniaStores/repositories'
import { AppStoreRepository } from '../../../src/types'

describe('useRepositoriesStore', () => {
  beforeEach(() => {
    setActivePinia(createTestingPinia({ stubActions: false }))
  })

  describe('method "setRepositories"', () => {
    it('updates the repositories state', () => {
      const store = useRepositoriesStore()
      expect(store.repositories).toEqual([])

      const repositories: AppStoreRepository[] = [
        { name: 'repo1', url: 'https://example.test/repo1.json' },
        { name: 'repo2', url: 'https://example.test/repo2.json' }
      ]
      store.setRepositories(repositories)

      expect(store.repositories).toEqual(repositories)
    })
  })
})
