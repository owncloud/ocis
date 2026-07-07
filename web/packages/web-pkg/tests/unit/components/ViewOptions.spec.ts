import { ref, unref } from 'vue'
import {
  defaultPlugins,
  defaultComponentMocks,
  mount,
  RouteLocation,
  PartialComponentProps
} from '@ownclouders/web-test-helpers'
import { mock } from 'vitest-mock-extended'
import ViewOptions from '../../../src/components/ViewOptions.vue'
import {
  FolderViewModeConstants,
  useResourcesStore,
  useRouteQuery,
  useRouteQueryPersisted
} from '../../../src/composables'
import { FolderView } from '../../../src'
import { OcPageSize, OcSwitch } from '@ownclouders/design-system/components'

vi.mock('../../../src/composables/router', async (importOriginal) => ({
  ...(await importOriginal<any>()),
  useRouteQueryPersisted: vi.fn(),
  useRouteQuery: vi.fn()
}))

const selectors = {
  pageSizeSelect: '.oc-page-size',
  hiddenFilesSwitch: '[data-testid="files-switch-hidden-files"]',
  flatListSwitch: '[data-testid="files-switch-flat-list"]',
  fileExtensionsSwitch: '[data-testid="files-switch-files-extensions-files"]',
  viewModeSwitchBtns: '.viewmode-switch-buttons',
  tileSizeSlider: '[data-testid="files-tiles-size-slider"]'
}

describe('ViewOptions component', () => {
  describe('pagination', () => {
    it('does not show when disabled', () => {
      const { wrapper } = getWrapper({ props: { hasPagination: false } })
      expect(wrapper.find(selectors.pageSizeSelect).exists()).toBeFalsy()
    })
    it('sets the correct initial files page limit', () => {
      const perPage = '100'
      const { wrapper } = getWrapper({ perPage })
      expect(
        wrapper.findComponent<typeof OcPageSize>(selectors.pageSizeSelect).props().selected
      ).toBe(perPage)
    })
    it('sets the correct files page limit', () => {
      const perPage = '100'
      const newItemsPerPage = '500'
      const { wrapper, mocks } = getWrapper({ perPage }) as any
      wrapper.vm.setItemsPerPage(newItemsPerPage)
      expect(mocks.$router.replace).toHaveBeenCalledWith(
        expect.objectContaining({
          query: expect.objectContaining({ 'items-per-page': newItemsPerPage })
        })
      )
    })
    it('resets the page to 1 if current page is > 1', () => {
      const perPage = '100'
      const newItemsPerPage = '500'
      const { wrapper, mocks } = getWrapper({ perPage, currentPage: '2' }) as any
      wrapper.vm.setItemsPerPage(newItemsPerPage)
      expect(mocks.$router.replace).toHaveBeenCalledWith(
        expect.objectContaining({
          query: expect.objectContaining({ 'items-per-page': newItemsPerPage, page: '1' })
        })
      )
    })
  })
  describe('hidden files toggle', () => {
    it('does not show when disabled', () => {
      const { wrapper } = getWrapper({ props: { hasHiddenFiles: false } })
      expect(wrapper.find(selectors.hiddenFilesSwitch).exists()).toBeFalsy()
    })
    it('toggles the setting to show/hide hidden files', () => {
      const { wrapper } = getWrapper()
      wrapper
        .findComponent<typeof OcSwitch>(selectors.hiddenFilesSwitch)
        .vm.$emit('update:checked', false)

      const { setAreHiddenFilesShown } = useResourcesStore()
      expect(setAreHiddenFilesShown).toHaveBeenCalled()
    })
  })
  describe('flat list toggle', () => {
    it('does not show when disabled', () => {
      const { wrapper } = getWrapper({ props: { shouldShowFlatListToggle: false } })
      expect(wrapper.find(selectors.flatListSwitch).exists()).toBeFalsy()
    })
    it('toggles the setting to show/hide flat list', () => {
      const { wrapper } = getWrapper()
      wrapper
        .findComponent<typeof OcSwitch>(selectors.flatListSwitch)
        .vm.$emit('update:checked', false)

      const { setShouldShowFlatList } = useResourcesStore()
      expect(setShouldShowFlatList).toHaveBeenCalled()
    })
  })
  describe('file extension toggle', () => {
    it('does not show when disabled', () => {
      const { wrapper } = getWrapper({ props: { hasFileExtensions: false } })
      expect(wrapper.find(selectors.fileExtensionsSwitch).exists()).toBeFalsy()
    })
    it('toggles the setting to show/hide file extensions', () => {
      const { wrapper } = getWrapper()
      wrapper
        .findComponent<typeof OcSwitch>(selectors.fileExtensionsSwitch)
        .vm.$emit('update:checked', false)

      const { setAreFileExtensionsShown } = useResourcesStore()
      expect(setAreFileExtensionsShown).toHaveBeenCalled()
    })
  })
  describe('view mode switcher', () => {
    it('does not show initially', () => {
      const { wrapper } = getWrapper()
      expect(wrapper.find(selectors.viewModeSwitchBtns).exists()).toBeFalsy()
    })
    it('shows if more than one viewModes are passed', () => {
      const { wrapper } = getWrapper({
        props: {
          viewModes: [
            mock<FolderView>({ name: '1', label: '' }),
            mock<FolderView>({ name: '2', label: '' })
          ]
        }
      })
      expect(wrapper.find(selectors.viewModeSwitchBtns).exists()).toBeTruthy()
    })
  })
  describe('tile size slider', () => {
    it('does not show initially', () => {
      const { wrapper } = getWrapper()
      expect(wrapper.find(selectors.tileSizeSlider).exists()).toBeFalsy()
    })
    it('shows if the viewModes include "resource-tiles"', () => {
      const { wrapper } = getWrapper({
        props: { viewModes: [mock<FolderView>({ name: FolderViewModeConstants.name.tiles })] }
      })
      expect(wrapper.find(selectors.tileSizeSlider).exists()).toBeTruthy()
    })
    it.each([1, 2, 3, 4, 5, 6])('applies the correct size step', (tileSize) => {
      const { mocks } = getWrapper({
        tileSize: tileSize.toString(),
        props: { viewModes: [mock<FolderView>({ name: FolderViewModeConstants.name.tiles })] }
      })
      expect(unref(mocks.tileSizeQueryMock)).toBe(tileSize.toString())
    })
  })
})

function getWrapper({
  perPage = '100',
  viewMode = FolderViewModeConstants.name.table,
  tileSize = '1',
  props = {},
  currentPage = '1'
}: {
  perPage?: string
  viewMode?: string
  tileSize?: string
  props?: PartialComponentProps<typeof ViewOptions>
  currentPage?: string
} = {}) {
  vi.mocked(useRouteQueryPersisted).mockImplementationOnce(() => ref(perPage))
  vi.mocked(useRouteQueryPersisted).mockImplementationOnce(() => ref(viewMode))
  const tileSizeQueryMock = ref(tileSize)
  vi.mocked(useRouteQueryPersisted).mockImplementationOnce(() => tileSizeQueryMock)
  vi.mocked(useRouteQuery).mockImplementationOnce(() => ref(currentPage))

  const mocks = {
    ...defaultComponentMocks({ currentRoute: mock<RouteLocation>({ path: '/files' }) }),
    tileSizeQueryMock
  }
  return {
    mocks,
    wrapper: mount(ViewOptions, {
      props: {
        perPageStoragePrefix: '',
        ...props
      },
      global: {
        mocks,
        provide: mocks,
        stubs: { OcButton: true, OcPageSize: false, OcSelect: true },
        plugins: [...defaultPlugins()]
      }
    })
  }
}
