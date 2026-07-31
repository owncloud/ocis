import { defaultComponentMocks, getComposableWrapper } from '@ownclouders/web-test-helpers'
import { CapabilityStore, useSearch } from '../../../../src/composables'
import { SearchResource, SpaceResource } from '@ownclouders/web-client'

describe('useSearch', () => {
  describe('method "buildSearchTerm"', () => {
    it('appends vault:true when isVault is true', () => {
      const wrapper = createWrapper()
      const result = wrapper.vm.buildSearchTerm({ term: 'test', isVault: true })
      expect(result).toContain('vault:true')
    })
    it('does not append vault:true when isVault is false', () => {
      const wrapper = createWrapper()
      const result = wrapper.vm.buildSearchTerm({ term: 'test', isVault: false })
      expect(result).not.toContain('vault:true')
    })
    it('does not append vault:true when isVault is undefined', () => {
      const wrapper = createWrapper()
      const result = wrapper.vm.buildSearchTerm({ term: 'test' })
      expect(result).not.toContain('vault:true')
    })
    it('combines vault:true with other query parts', () => {
      const wrapper = createWrapper()
      const result = wrapper.vm.buildSearchTerm({
        term: 'test',
        tags: 'lorem',
        isVault: true
      })
      expect(result).toContain('vault:true')
      expect(result).toContain('tag:("lorem")')
      expect(result).toContain('name:"*test*"')
    })
    it('places scope: at the end of the query', () => {
      const wrapper = createWrapper()
      const result = wrapper.vm.buildSearchTerm({
        term: 'test',
        scope: 'lorem',
        useScope: true,
        isVault: true
      })
      expect(result).toMatch(/scope:lorem$/)
    })
  })
  describe('method "search"', () => {
    it('can search', async () => {
      const files = [
        { id: 'foo', name: 'foo' },
        { id: 'bar', name: 'bar' },
        { id: 'baz', name: 'baz' }
      ] as SearchResource[]

      const wrapper = createWrapper({ resources: files })

      const noTermResult = await wrapper.vm.search('')
      expect(noTermResult).toEqual({ totalResults: null, values: [] })

      const withTermResult = await wrapper.vm.search('foo')
      expect(withTermResult.values.map((r) => r.data)).toMatchObject(files)
    })
    it('properly returns space resources', async () => {
      const files = [{ id: 'foo', name: 'foo', parentFolderId: '2' }] as SearchResource[]

      const wrapper = createWrapper({ resources: files })

      const withTerm = await wrapper.vm.search('foo')
      expect(withTerm.values.map((r) => r.data)[0].id).toEqual('2')
    })
  })
})

const createWrapper = ({ resources = [] }: { resources?: SearchResource[] } = {}) => {
  const spaces = [
    {
      id: '1',
      fileId: '1',
      driveType: 'personal',
      getDriveAliasAndItem: () => 'personal/admin'
    },
    {
      id: '2',
      driveType: 'project',
      name: 'New space',
      getDriveAliasAndItem: vi.fn()
    }
  ] as unknown as SpaceResource[]

  const mocks = defaultComponentMocks({})
  const capabilities = {
    spaces: { projects: true }
  } satisfies Partial<CapabilityStore['capabilities']>

  mocks.$clientService.webdav.search.mockResolvedValue({
    resources,
    totalResults: resources.length
  })

  return getComposableWrapper(
    () => {
      const { search, buildSearchTerm } = useSearch()

      return {
        search,
        buildSearchTerm
      }
    },
    {
      mocks,
      provide: mocks,
      pluginOptions: {
        piniaOptions: { spacesState: { spaces }, capabilityState: { capabilities } }
      }
    }
  )
}
