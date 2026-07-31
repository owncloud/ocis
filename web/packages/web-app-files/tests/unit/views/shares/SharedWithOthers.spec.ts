import SharedWithOthers from '../../../../src/views/shares/SharedWithOthers.vue'
import { useResourcesViewDefaults } from '../../../../src/composables'
import { useResourcesViewDefaultsMock } from '../../../../tests/mocks/useResourcesViewDefaultsMock'
import { ref } from 'vue'
import { defaultStubs, RouteLocation } from '@ownclouders/web-test-helpers'
import { mock, mockDeep } from 'vitest-mock-extended'
import { IncomingShareResource } from '@ownclouders/web-client'
import { defaultPlugins, mount, defaultComponentMocks } from '@ownclouders/web-test-helpers'
import { ShareTypes } from '@ownclouders/web-client'
import { useSortMock } from '../../../mocks/useSortMock'
import { ResourceTable, AppBar } from '@ownclouders/web-pkg'

vi.mock('../../../../src/composables')
vi.mock('@ownclouders/web-pkg', async (importOriginal) => ({
  ...(await importOriginal<any>()),
  useSort: vi.fn().mockImplementation(() => useSortMock()),
  queryItemAsString: vi.fn(),
  useRouteQuery: vi.fn(),
  useFileActions: vi.fn(() => ({
    triggerDefaultAction: vi.fn()
  }))
}))

describe('SharedWithOthers view', () => {
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
      const mockedFiles = [
        mockDeep<IncomingShareResource>({ shareTypes: [ShareTypes.user.value] }),
        mockDeep<IncomingShareResource>({ shareTypes: [ShareTypes.user.value] })
      ]
      const { wrapper } = getMountedWrapper({ files: mockedFiles })
      expect(wrapper.find('.no-content-message').exists()).toBeFalsy()
      expect(wrapper.find('resource-table-stub').exists()).toBeTruthy()
      expect(
        wrapper.findComponent<typeof ResourceTable>('resource-table-stub').props().resources.length
      ).toEqual(mockedFiles.length)
    })
  })
  describe('filter', () => {
    describe('share type', () => {
      it('shows filter if multiple share types are present', () => {
        const { wrapper } = getMountedWrapper({
          files: [
            mock<IncomingShareResource>({ shareTypes: [ShareTypes.user.value] }),
            mock<IncomingShareResource>({ shareTypes: [ShareTypes.group.value] })
          ]
        })
        expect(wrapper.find('.share-type-filter').exists()).toBeTruthy()
      })
      it('does not show filter if only one share type is present', () => {
        const { wrapper } = getMountedWrapper({
          files: [mock<IncomingShareResource>({ shareTypes: [ShareTypes.user.value] })]
        })
        expect(wrapper.find('.share-type-filter').exists()).toBeFalsy()
      })
    })
  })
})

function getMountedWrapper({
  mocks = {},
  files = [],
  loading = false
}: {
  mocks?: Record<string, unknown>
  files?: IncomingShareResource[]
  loading?: boolean
} = {}) {
  vi.mocked(useResourcesViewDefaults).mockImplementation(() =>
    useResourcesViewDefaultsMock({
      paginatedResources: ref(files),
      areResourcesLoading: ref(loading)
    })
  )
  const defaultMocks = {
    ...defaultComponentMocks({
      currentRoute: mock<RouteLocation>({ name: 'files-shares-with-others' })
    }),
    ...(mocks && mocks)
  }

  return {
    mocks: defaultMocks,
    wrapper: mount(SharedWithOthers, {
      global: {
        components: {
          AppBar,
          ResourceTable
        },
        plugins: [...defaultPlugins()],
        mocks: defaultMocks,
        provide: defaultMocks,
        stubs: { ...defaultStubs, ItemFilter: true }
      }
    })
  }
}
