import { mock } from 'vitest-mock-extended'
import { ResourcePreview, SearchResultValue } from '../../../../src/components'
import { SpaceResource } from '@ownclouders/web-client'
import { useGetMatchingSpace } from '../../../../src/composables/spaces/useGetMatchingSpace'
import {
  defaultComponentMocks,
  defaultPlugins,
  mount,
  useGetMatchingSpaceMock
} from '@ownclouders/web-test-helpers'
import { useFileActions } from '../../../../src/composables/actions'
import { CapabilityStore } from '../../../../src/composables/piniaStores'
import ResourceListItem from '../../../../src/components/FilesList/ResourceListItem.vue'

vi.mock('../../../../src/composables/spaces/useGetMatchingSpace', () => ({
  useGetMatchingSpace: vi.fn()
}))

vi.mock('../../../../src/composables/actions', () => ({
  useFileActions: vi.fn()
}))

describe('Preview component', () => {
  const driveAliasAndItem = '1'
  vi.mocked(useGetMatchingSpace).mockImplementation(() => useGetMatchingSpaceMock())
  it('should render preview component', async () => {
    const { wrapper } = getWrapper({
      space: mock<SpaceResource>({
        id: '1',
        driveType: 'project',
        name: 'New space',
        getDriveAliasAndItem: () => driveAliasAndItem
      })
    })
    ;(wrapper.vm as any).previewData = 'blob:image'
    await wrapper.vm.$nextTick()
    expect(wrapper.html()).toMatchSnapshot()
  })
  it('should render resource component without file extension when areFileExtensionsShown is set to false', () => {
    const { wrapper } = getWrapper({
      areFileExtensionsShown: false,
      space: mock<SpaceResource>({
        id: '1',
        driveType: 'project',
        name: 'New space',
        getDriveAliasAndItem: () => driveAliasAndItem
      })
    })

    expect(wrapper.findComponent<any>(ResourceListItem).props('isExtensionDisplayed')).toBe(false)
  })
})

function getWrapper({
  space = null,
  searchResult = mock<SearchResultValue>({
    id: '1',
    data: {
      storageId: '1',
      name: 'lorem.txt',
      path: '/',
      remoteItemPath: ''
    }
  }),
  areFileExtensionsShown = true
}: {
  space?: SpaceResource
  searchResult?: SearchResultValue
  areFileExtensionsShown?: boolean
} = {}) {
  vi.mocked(useGetMatchingSpace).mockImplementation(() =>
    useGetMatchingSpaceMock({
      isResourceAccessible() {
        return true
      },
      getMatchingSpace() {
        return space
      }
    })
  )
  vi.mocked(useFileActions).mockReturnValue(mock<ReturnType<typeof useFileActions>>())

  const mocks = defaultComponentMocks()
  const capabilities = {
    spaces: { projects: true }
  } satisfies Partial<CapabilityStore['capabilities']>

  return {
    wrapper: mount(ResourcePreview, {
      props: {
        searchResult
      },
      global: {
        provide: mocks,
        renderStubDefaultSlot: true,
        mocks,
        plugins: [
          ...defaultPlugins({
            piniaOptions: {
              capabilityState: { capabilities },
              resourcesStore: { areFileExtensionsShown }
            }
          })
        ]
      }
    })
  }
}
