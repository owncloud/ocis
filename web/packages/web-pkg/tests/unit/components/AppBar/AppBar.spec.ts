import { ref } from 'vue'
import AppBar from '../../../../src/components/AppBar/AppBar.vue'
import { mock } from 'vitest-mock-extended'
import { Resource, SpaceResource } from '@ownclouders/web-client'
import {
  defaultComponentMocks,
  defaultPlugins,
  shallowMount,
  RouteLocation,
  PartialComponentProps
} from '@ownclouders/web-test-helpers'
import { ArchiverService } from '../../../../src/services'
import { FolderView } from '../../../../src/ui/types'
import { useExtensionRegistry, ViewOptions } from '../../../../src'
import { OcBreadcrumb } from '@ownclouders/design-system/components'

const selectors = {
  ocBreadcrumbStub: 'oc-breadcrumb-stub',
  batchActionsStub: 'batch-actions-stub',
  viewOptionsStub: 'view-options-stub',
  sidebarToggleStub: 'sidebar-toggle-stub',
  mobileNavPortal: 'portal-target[name="app.runtime.mobile.nav"]'
}

const selectedFiles = [mock<Resource>(), mock<Resource>()]
const actionSlot = "<button class='action-slot'>Click</button>"
const contentSlot = "<div class='content-slot'>Foo</div>"

const breadcrumbItems = [
  { text: 'Example1', to: '/' },
  { text: 'Example2', to: '/foo' }
]
const breadCrumbItemWithContextActionAllowed = {
  text: 'Example Special',
  to: '/bar',
  allowContextActions: true
}

describe('AppBar component', () => {
  describe('renders', () => {
    it('by default no breadcrumbs, no bulkactions, no sharesnavigation but viewoptions and sidebartoggle', () => {
      const { wrapper } = getShallowWrapper()
      expect(wrapper.html()).toMatchSnapshot()
    })
    describe('breadcrumbs', () => {
      it('if given, by default without breadcrumbsContextActionsItems', () => {
        const { wrapper } = getShallowWrapper([], {}, { breadcrumbs: breadcrumbItems })
        expect(wrapper.find(selectors.ocBreadcrumbStub).exists()).toBeTruthy()
        expect(
          wrapper.findComponent<typeof OcBreadcrumb>(selectors.ocBreadcrumbStub).props('items')
        ).toEqual(breadcrumbItems)
      })
      it('if given, with breadcrumbsContextActionsItems if allowed on last breadcrumb item', () => {
        const { wrapper } = getShallowWrapper(
          [],
          {},
          { breadcrumbs: [...breadcrumbItems, breadCrumbItemWithContextActionAllowed] }
        )
        expect(wrapper.find(selectors.ocBreadcrumbStub).exists()).toBeTruthy()
        expect(
          wrapper.findComponent<typeof OcBreadcrumb>(selectors.ocBreadcrumbStub).props('items')
        ).toEqual([...breadcrumbItems, breadCrumbItemWithContextActionAllowed])
      })
      it('not if no breadcrumb items given', () => {
        const { wrapper } = getShallowWrapper([], {}, { breadcrumbs: [] })
        expect(wrapper.find(selectors.ocBreadcrumbStub).exists()).toBeFalsy()
      })
      it('not if one breadcrumb item is given in mobile view', () => {
        const { wrapper } = getShallowWrapper(
          [],
          {},
          { breadcrumbs: [breadcrumbItems[0]] },
          mock<RouteLocation>({ name: '' }),
          true
        )
        expect(wrapper.find(selectors.ocBreadcrumbStub).exists()).toBeFalsy()
      })
    })
    describe('bulkActions', () => {
      it('if enabled', () => {
        const { wrapper } = getShallowWrapper(selectedFiles, {}, { hasBulkActions: true })
        expect(wrapper.find(selectors.batchActionsStub).exists()).toBeTruthy()
      })
      it('if 1 file selected on trash routes', () => {
        const { wrapper } = getShallowWrapper(
          [selectedFiles[0]],
          {},
          { hasBulkActions: true },
          mock<RouteLocation>({
            name: 'files-trash-generic',
            path: '/files/trash/personal/admin'
          })
        )
        expect(wrapper.find(selectors.batchActionsStub).exists()).toBeTruthy()
      })
    })
    describe('mobile navigation portal', () => {
      it.each([
        { items: [], shows: true },
        { items: [breadcrumbItems[0]], shows: true },
        { items: [breadcrumbItems[0], breadcrumbItems[1]], shows: false }
      ])('if less than 2 breadcrumb items given', ({ items, shows }) => {
        const { wrapper } = getShallowWrapper([], {}, { breadcrumbs: items })
        expect(wrapper.find(selectors.mobileNavPortal).exists()).toBe(shows)
      })
    })
    describe('viewoptions', () => {
      it('show if options are available', () => {
        const { wrapper } = getShallowWrapper([], {}, { hasViewOptions: true })
        expect(wrapper.find(selectors.viewOptionsStub).exists()).toBeTruthy()
      })
      it('hide if options are not available', () => {
        const { wrapper } = getShallowWrapper([], {}, { hasViewOptions: false })
        expect(wrapper.find(selectors.viewOptionsStub).exists()).toBeFalsy()
      })
      it('passes viewModes array to ViewOptions', () => {
        const viewModes = [mock<FolderView>()]
        const { wrapper } = getShallowWrapper([], {}, { hasViewOptions: true, viewModes })
        expect(
          wrapper.findComponent<typeof ViewOptions>(selectors.viewOptionsStub).props('viewModes')
        ).toEqual(viewModes)
      })
    })
    it('if given, with content in the actions slot', () => {
      const { wrapper } = getShallowWrapper([], { actions: actionSlot })
      expect(wrapper.html()).toMatchSnapshot()
    })
    it('if given, with content in the content slot', () => {
      const { wrapper } = getShallowWrapper([], { content: contentSlot })
      expect(wrapper.html()).toMatchSnapshot()
    })
  })
})

function getShallowWrapper(
  selected: Resource[] = [],
  slots = {},
  props: PartialComponentProps<typeof AppBar> = {
    breadcrumbs: [],
    viewModes: [],
    hasBulkActions: false,
    hasViewOptions: true
  },
  currentRoute = mock<RouteLocation>({
    name: 'files-spaces-generic',
    path: '/files/spaces/personal/admin'
  }),
  isMobileWidth = false
) {
  const plugins = defaultPlugins({
    piniaOptions: {
      resourcesStore: { resources: selected, selectedIds: selected.map(({ id }) => id) }
    }
  })

  const { requestExtensions } = useExtensionRegistry()
  vi.mocked(requestExtensions).mockReturnValue([])

  const mocks = {
    ...defaultComponentMocks({
      currentRoute
    }),
    $archiverService: mock<ArchiverService>()
  }
  mocks.$route.meta.title = 'ExampleTitle'

  return {
    wrapper: shallowMount(AppBar, {
      props: { ...props, space: mock<SpaceResource>() },
      slots,
      global: {
        plugins,
        provide: { ...mocks, isMobileWidth: ref(isMobileWidth) },
        mocks
      }
    })
  }
}
