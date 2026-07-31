import { computed, ref } from 'vue'
import { mock, mockDeep } from 'vitest-mock-extended'
import { AbilityRule, Resource, SpaceResource } from '@ownclouders/web-client'
import GenericSpace from '../../../../src/views/spaces/GenericSpace.vue'
import { useResourcesViewDefaults } from '../../../../src/composables/resourcesViewDefaults'
import { useResourcesViewDefaultsMock } from '../../../../tests/mocks/useResourcesViewDefaultsMock'
import {
  defaultPlugins,
  mount,
  defaultComponentMocks,
  defaultStubs,
  RouteLocation,
  ComponentProps,
  PartialComponentProps,
  PiniaMockOptions
} from '@ownclouders/web-test-helpers'
import {
  AppBar,
  FolderViewExtension,
  useBreadcrumbsFromPath,
  useExtensionRegistry
} from '@ownclouders/web-pkg'
import { useBreadcrumbsFromPathMock } from '../../../mocks/useBreadcrumbsFromPathMock'
import { h } from 'vue'
import { BreadcrumbItem } from '@ownclouders/design-system/helpers'
import {
  folderViewsFavoritesExtensionPoint,
  folderViewsFolderExtensionPoint,
  folderViewsProjectSpacesExtensionPoint
} from '../../../../src/extensionPoints'

const mockCreateFolder = vi.fn()
const mockUseEmbedMode = vi.fn().mockReturnValue({ isEnabled: computed(() => false) })

vi.mock('../../../../src/composables/resourcesViewDefaults')
vi.mock('../../../../src/composables/keyboardActions')
vi.mock('@ownclouders/web-pkg', async (importOriginal) => ({
  ...(await importOriginal<any>()),
  useBreadcrumbsFromPath: vi.fn(),
  useFileActionsCreateNewFolder: () => ({
    actions: [{ handler: mockCreateFolder }]
  }),
  useEmbedMode: vi.fn().mockImplementation(() => mockUseEmbedMode()),
  useFileActions: vi.fn(() => ({})),
  useOpenWithDefaultApp: vi.fn(() => ({}))
}))

const selectors = Object.freeze({
  btnCreateFolder: '[data-testid="btn-new-folder"]',
  actionsCreateAndUpload: '[data-testid="actions-create-and-upload"]'
})

describe('GenericSpace view', () => {
  it('appBar always present', () => {
    const { wrapper } = getMountedWrapper()
    expect(wrapper.find('app-bar-stub').exists()).toBeTruthy()
  })
  it('sideBar always present', () => {
    const { wrapper } = getMountedWrapper()
    expect(wrapper.find('file-side-bar-stub').exists()).toBeTruthy()
  })
  describe('space header', () => {
    it('does not render the space header in the personal space', () => {
      const { wrapper } = getMountedWrapper()
      expect(wrapper.find('space-header-stub').exists()).toBeFalsy()
    })
    it('does not render the space header in a nested project space', () => {
      const { wrapper } = getMountedWrapper({
        props: {
          item: '/someFolder',
          space: mock<SpaceResource>({ driveType: 'project' })
        }
      })
      expect(wrapper.find('space-header-stub').exists()).toBeFalsy()
    })
    it('renders the space header on a space frontpage', () => {
      const { wrapper } = getMountedWrapper({
        props: {
          space: mock<SpaceResource>({ driveType: 'project' })
        }
      })
      expect(wrapper.find('space-header-stub').exists()).toBeTruthy()
    })
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
      const { wrapper } = getMountedWrapper({ files: [mock<Resource>()] })
      expect(wrapper.find('.no-content-message').exists()).toBeFalsy()
      expect(wrapper.find('.resource-table').exists()).toBeTruthy()
    })
  })
  describe('breadcrumbs', () => {
    it.each([
      { driveType: 'personal', expectedItems: 1 },
      { driveType: 'project', expectedItems: 2 },
      { driveType: 'share', expectedItems: 3 }
    ])('include root item(s)', ({ driveType, expectedItems }) => {
      const space = mock<SpaceResource>({
        id: '1',
        getDriveAliasAndItem: vi.fn(),
        driveType,
        isOwner: () => driveType === 'personal'
      })
      const { wrapper } = getMountedWrapper({ files: [mockDeep<Resource>()], props: { space } })
      expect(wrapper.findComponent<typeof AppBar>('app-bar-stub').props().breadcrumbs.length).toBe(
        expectedItems
      )
    })
    describe('personal space with vault access', () => {
      it('shows "Drive" as root breadcrumb when scope is not vault', () => {
        const space = mock<SpaceResource>({
          id: '1',
          getDriveAliasAndItem: vi.fn(),
          driveType: 'personal',
          name: 'Personal space',
          isOwner: () => true
        })
        const { wrapper } = getMountedWrapper({
          files: [mockDeep<Resource>()],
          props: { space },
          abilities: [{ action: 'read-all', subject: 'Vault' }],
          capabilityState: { capabilities: { vault: { enabled: true } } }
        })
        const breadcrumbs = wrapper.findComponent<typeof AppBar>('app-bar-stub').props().breadcrumbs
        expect(breadcrumbs.length).toBe(2)
        expect(breadcrumbs[0].text).toBe('Drive')
        expect(breadcrumbs[1].text).toBe('Personal space')
      })
      it('shows "Vault" as root breadcrumb when scope is vault', () => {
        const space = mock<SpaceResource>({
          id: '1',
          getDriveAliasAndItem: vi.fn(),
          driveType: 'personal',
          name: 'Personal space',
          isOwner: () => true
        })
        const { wrapper } = getMountedWrapper({
          files: [mockDeep<Resource>()],
          props: { space },
          abilities: [{ action: 'read-all', subject: 'Vault' }],
          capabilityState: { capabilities: { vault: { enabled: true } } },
          currentRoute: { name: 'files-spaces-generic', path: '/', params: { scope: 'vault' } }
        })
        const breadcrumbs = wrapper.findComponent<typeof AppBar>('app-bar-stub').props().breadcrumbs
        expect(breadcrumbs.length).toBe(2)
        expect(breadcrumbs[0].text).toBe('Vault')
        expect(breadcrumbs[1].text).toBe('Personal space')
      })
      it('shows only space name when user cannot access vault', () => {
        const space = mock<SpaceResource>({
          id: '1',
          getDriveAliasAndItem: vi.fn(),
          driveType: 'personal',
          name: 'Personal space',
          isOwner: () => true
        })
        const { wrapper } = getMountedWrapper({
          files: [mockDeep<Resource>()],
          props: { space }
        })
        const breadcrumbs = wrapper.findComponent<typeof AppBar>('app-bar-stub').props().breadcrumbs
        expect(breadcrumbs.length).toBe(1)
        expect(breadcrumbs[0].text).toBe('Personal space')
      })
    })
    describe('project space breadcrumbs', () => {
      it('shows "Spaces" when user cannot access vault', () => {
        const space = mock<SpaceResource>({
          id: '1',
          getDriveAliasAndItem: vi.fn(),
          driveType: 'project'
        })
        const { wrapper } = getMountedWrapper({
          files: [mockDeep<Resource>()],
          props: { space }
        })
        const breadcrumbs = wrapper.findComponent<typeof AppBar>('app-bar-stub').props().breadcrumbs
        expect(breadcrumbs[0].text).toBe('Spaces')
      })
      it('shows "Drive" when user can access vault and scope is not vault', () => {
        const space = mock<SpaceResource>({
          id: '1',
          getDriveAliasAndItem: vi.fn(),
          driveType: 'project'
        })
        const { wrapper } = getMountedWrapper({
          files: [mockDeep<Resource>()],
          props: { space },
          abilities: [{ action: 'read-all', subject: 'Vault' }],
          capabilityState: { capabilities: { vault: { enabled: true } } }
        })
        const breadcrumbs = wrapper.findComponent<typeof AppBar>('app-bar-stub').props().breadcrumbs
        expect(breadcrumbs[0].text).toBe('Drive')
      })
      it('shows "Vault" when user can access vault and scope is vault', () => {
        const space = mock<SpaceResource>({
          id: '1',
          getDriveAliasAndItem: vi.fn(),
          driveType: 'project'
        })
        const { wrapper } = getMountedWrapper({
          files: [mockDeep<Resource>()],
          props: { space },
          abilities: [{ action: 'read-all', subject: 'Vault' }],
          capabilityState: { capabilities: { vault: { enabled: true } } },
          currentRoute: { name: 'files-spaces-generic', path: '/', params: { scope: 'vault' } }
        })
        const breadcrumbs = wrapper.findComponent<typeof AppBar>('app-bar-stub').props().breadcrumbs
        expect(breadcrumbs[0].text).toBe('Vault')
      })
    })
    it('include the root item and the current folder', () => {
      const folderName = 'someFolder'
      const { wrapper } = getMountedWrapper({
        files: [mockDeep<Resource>()],
        props: { item: `/${folderName}` },
        breadcrumbsFromPath: [{ text: folderName }]
      })
      expect(wrapper.findComponent<typeof AppBar>('app-bar-stub').props().breadcrumbs.length).toBe(
        2
      )
      expect(
        wrapper.findComponent<typeof AppBar>('app-bar-stub').props().breadcrumbs[1].text
      ).toEqual(folderName)
    })
    it('omit the "page"-query of the current route', () => {
      const currentRoute = { name: 'files-spaces-generic', path: '/', query: { page: '2' } }
      const { wrapper } = getMountedWrapper({
        files: [mockDeep<Resource>()],
        props: { item: 'someFolder' },
        currentRoute
      })
      const breadCrumbItem = wrapper.findComponent<typeof AppBar>('app-bar-stub').props()
        .breadcrumbs[0]
      expect((breadCrumbItem.to as RouteLocation).query.page).toBeUndefined()
    })
  })
  describe('loader task', () => {
    it('re-loads the resources on item change', async () => {
      const { wrapper, mocks } = getMountedWrapper()
      await (wrapper.vm as any).loadResourcesTask.last
      expect(mocks.refreshFileListHeaderPosition).toHaveBeenCalledTimes(1)
      await wrapper.setProps({ item: 'newItem' })
      await (wrapper.vm as any).loadResourcesTask.last
      expect(mocks.refreshFileListHeaderPosition).toHaveBeenCalledTimes(2)
    })
    it('re-loads the resources on space change', async () => {
      const { wrapper, mocks } = getMountedWrapper()
      await (wrapper.vm as any).loadResourcesTask.last
      expect(mocks.refreshFileListHeaderPosition).toHaveBeenCalledTimes(1)
      await wrapper.setProps({ space: mockDeep<SpaceResource>() })
      await (wrapper.vm as any).loadResourcesTask.last
      expect(mocks.refreshFileListHeaderPosition).toHaveBeenCalledTimes(2)
    })
  })
  describe('empty folder upload hint', () => {
    it('renders if the user can upload to the current folder', () => {
      const { wrapper } = getMountedWrapper({
        currentFolder: mock<Resource>({ canUpload: () => true })
      })
      expect(wrapper.find('.file-empty-upload-hint').exists()).toBeTruthy()
    })
    it('does not render if the user can not upload to the current folder', () => {
      const { wrapper } = getMountedWrapper({
        currentFolder: mock<Resource>({ canUpload: () => false })
      })
      expect(wrapper.find('.file-empty-upload-hint').exists()).toBeFalsy()
    })
  })
  describe('whitespace context menu', () => {
    it('shows whitespace context menu on right click in whitespace', async () => {
      const { wrapper } = getMountedWrapper()
      await (wrapper.vm as any).loadResourcesTask.last
      await wrapper.find('#files-view').trigger('contextmenu')
      await wrapper.vm.$nextTick()
      expect((wrapper.vm as any).whitespaceContextMenu).toBeDefined()
    })
  })
  describe('for a single file', () => {
    describe('on EOS for single shared resources', () => {
      it('renders the ResourceDetails component if no currentFolder id is present', () => {
        const { wrapper } = getMountedWrapper({
          currentFolder: mock<Resource>({ fileId: '' }),
          files: [mock<Resource>({ isFolder: false })],
          runningOnEos: true
        })
        expect(wrapper.find('resource-details-stub').exists()).toBeTruthy()
      })
      it('renders the ResourceDetails component if currentFolder path matches single shared resource path', () => {
        const path = 'foo'
        const { wrapper } = getMountedWrapper({
          currentFolder: {
            ...mock<Resource>(),
            path
          },
          files: [{ ...mock<Resource>(), path }],
          runningOnEos: true
        })
        expect(wrapper.find('resource-details-stub').exists()).toBeTruthy()
      })
    })
    describe('on public links', () => {
      it('renders the ResourceDetails component', () => {
        const { wrapper } = getMountedWrapper({
          currentFolder: {
            ...mock<Resource>()
          },
          files: [{ ...mock<Resource>(), isFolder: false }],
          space: mock<SpaceResource>({
            id: '1',
            getDriveAliasAndItem: vi.fn(),
            name: 'Personal space',
            driveType: 'public'
          })
        })
        expect(wrapper.find('resource-details-stub').exists()).toBeTruthy()
      })
    })
  })
  describe('create and upload actions', () => {
    const AppBarStub = { template: '<div><slot name="actions" /></div>' }

    it('should not render create folder button when not in embed mode', () => {
      const { wrapper } = getMountedWrapper({
        stubs: { 'app-bar': AppBarStub, CreateAndUpload: true }
      })

      expect(wrapper.find(selectors.btnCreateFolder).exists()).toBe(false)
    })

    it('should render create and upload actions when not in embed mode', () => {
      const { wrapper } = getMountedWrapper({
        stubs: { 'app-bar': AppBarStub, CreateAndUpload: true }
      })

      expect(wrapper.find(selectors.actionsCreateAndUpload).exists()).toBe(true)
    })
  })
})

function getMountedWrapper({
  mocks = {},
  props = {},
  files = [],
  loading = false,
  currentRoute = { name: 'files-spaces-generic', path: '/' },
  currentFolder = mock<Resource>(),
  runningOnEos = false,
  space = mock<SpaceResource>({
    id: '1',
    getDriveAliasAndItem: vi.fn(),
    name: 'Personal space',
    driveType: ''
  }),
  breadcrumbsFromPath = [],
  stubs = {},
  abilities = [],
  capabilityState = {}
}: {
  mocks?: Record<string, unknown>
  props?: PartialComponentProps<typeof GenericSpace>
  files?: Resource[]
  loading?: boolean
  currentRoute?: { name?: string; path?: string; params?: Record<string, string> }
  currentFolder?: Resource
  runningOnEos?: boolean
  space?: SpaceResource
  breadcrumbsFromPath?: BreadcrumbItem[]
  stubs?: any
  abilities?: AbilityRule[]
  capabilityState?: PiniaMockOptions['capabilityState']
} = {}) {
  const plugins = defaultPlugins({
    abilities,
    piniaOptions: {
      configState: { options: { runningOnEos } },
      resourcesStore: { currentFolder },
      capabilityState
    }
  })

  const resourcesViewDetailsMock = useResourcesViewDefaultsMock({
    paginatedResources: ref(files),
    areResourcesLoading: ref(loading)
  })
  vi.mocked(useResourcesViewDefaults).mockImplementation(() => resourcesViewDetailsMock)
  vi.mocked(useBreadcrumbsFromPath).mockImplementation(() =>
    useBreadcrumbsFromPathMock({ breadcrumbsFromPath: vi.fn(() => breadcrumbsFromPath) })
  )

  const extensions = [
    {
      id: 'com.github.owncloud.web.files.folder-view.resource-table',
      type: 'folderView',
      extensionPointIds: [
        folderViewsFolderExtensionPoint.id,
        folderViewsProjectSpacesExtensionPoint.id,
        folderViewsFavoritesExtensionPoint.id
      ],
      folderView: {
        name: 'resource-table',
        label: 'Switch to default view',
        icon: {
          name: 'menu-line',
          fillType: 'none'
        },
        component: h('div', { class: 'resource-table' })
      }
    }
  ] satisfies FolderViewExtension[]
  const { requestExtensions } = useExtensionRegistry()
  vi.mocked(requestExtensions).mockReturnValue(extensions)

  const defaultMocks = {
    ...defaultComponentMocks({ currentRoute: mock<RouteLocation>(currentRoute) }),
    ...(mocks && mocks)
  }

  const propsData: ComponentProps<typeof GenericSpace> = {
    space,
    item: '/',
    itemId: undefined,
    ...props
  }

  return {
    mocks: { ...defaultMocks, ...resourcesViewDetailsMock },
    wrapper: mount(GenericSpace, {
      props: propsData,
      global: {
        components: {
          AppBar
        },
        plugins,
        mocks: defaultMocks,
        provide: defaultMocks,
        stubs: { ...defaultStubs, 'resource-details': true, portal: true, ...stubs }
      }
    })
  }
}
