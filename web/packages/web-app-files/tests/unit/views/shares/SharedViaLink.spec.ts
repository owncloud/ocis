import SharedViaLink from '../../../../src/views/shares/SharedViaLink.vue'
import { useResourcesViewDefaults } from '../../../../src/composables'
import { useResourcesViewDefaultsMock } from '../../../../tests/mocks/useResourcesViewDefaultsMock'
import { ref } from 'vue'
import { AppBar, ResourceTable } from '@ownclouders/web-pkg'
import { mock, mockDeep } from 'vitest-mock-extended'
import { OutgoingShareResource } from '@ownclouders/web-client'
import {
  defaultPlugins,
  mount,
  defaultComponentMocks,
  defaultStubs,
  RouteLocation
} from '@ownclouders/web-test-helpers'

vi.mock('../../../../src/composables')
vi.mock('@ownclouders/web-pkg', async (importOriginal) => ({
  ...(await importOriginal<any>()),
  useFileActions: vi.fn(() => ({
    triggerDefaultAction: vi.fn()
  }))
}))

describe('SharedViaLink view', () => {
  it('appBar always present', () => {
    const { wrapper } = getMountedWrapper()
    expect(wrapper.find('app-bar-stub').exists()).toBeTruthy()
  })
  it('sideBar always present', () => {
    const { wrapper } = getMountedWrapper()
    expect(wrapper.find('file-side-bar-stub').exists()).toBeTruthy()
  })
  describe('different files view states', () => {
    it('shows the loading spinner during loading', () => {
      const { wrapper } = getMountedWrapper({ loading: true })
      expect(wrapper.find('oc-spinner-stub').exists()).toBeTruthy()
    })
    it('shows the no-content-message after loading', () => {
      const { wrapper } = getMountedWrapper()
      expect(wrapper.find('oc-spinner-stub').exists()).toBeFalsy()
      expect(wrapper.find('.no-content-message').exists()).toBeTruthy()
    })
    it('shows the files table when files are available', () => {
      const mockedFiles: OutgoingShareResource[] = [
        mockDeep<OutgoingShareResource>({
          id: '1',
          fileId: 'f1',
          name: 'file1',
          sdate: '2024-01-01'
        }),
        mockDeep<OutgoingShareResource>({
          id: '2',
          fileId: 'f2',
          name: 'file2',
          sdate: '2024-01-02'
        })
      ]

      const { wrapper } = getMountedWrapper({ files: mockedFiles })
      expect(wrapper.find('.no-content-message').exists()).toBeFalsy()
      expect(wrapper.find('resource-table-stub').exists()).toBeTruthy()
      expect(
        wrapper.findComponent<typeof ResourceTable>('resource-table-stub').props().resources.length
      ).toEqual(mockedFiles.length)
    })
  })
})

function getMountedWrapper({
  mocks = {},
  files = [],
  loading = false
}: { mocks?: Record<string, unknown>; files?: OutgoingShareResource[]; loading?: boolean } = {}) {
  vi.mocked(useResourcesViewDefaults).mockImplementation(() =>
    useResourcesViewDefaultsMock({
      paginatedResources: ref(files),
      areResourcesLoading: ref(loading)
    })
  )
  const defaultMocks = {
    ...defaultComponentMocks({
      currentRoute: mock<RouteLocation>({ name: 'files-shares-via-link' })
    }),
    ...(mocks && mocks)
  }

  return {
    mocks: defaultMocks,
    wrapper: mount(SharedViaLink, {
      global: {
        components: {
          AppBar,
          ResourceTable
        },
        plugins: [...defaultPlugins()],
        mocks: defaultMocks,
        provide: defaultMocks,
        stubs: defaultStubs
      }
    })
  }
}
