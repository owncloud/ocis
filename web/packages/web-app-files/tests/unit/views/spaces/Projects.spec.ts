import Projects from '../../../../src/views/spaces/Projects.vue'
import { mock } from 'vitest-mock-extended'
import { h, nextTick, ref } from 'vue'
import {
  queryItemAsString,
  useFileActionsDelete,
  useExtensionRegistry,
  FolderViewExtension,
  AppBar,
  CreateSpace
} from '@ownclouders/web-pkg'

import {
  defaultPlugins,
  mount,
  defaultComponentMocks,
  defaultStubs,
  RouteLocation,
  PiniaMockOptions
} from '@ownclouders/web-test-helpers'
import { AbilityRule, SpaceResource } from '@ownclouders/web-client'
import {
  folderViewsFavoritesExtensionPoint,
  folderViewsFolderExtensionPoint,
  folderViewsProjectSpacesExtensionPoint
} from '../../../../src/extensionPoints'

vi.mock('@ownclouders/web-pkg', async (importOriginal) => ({
  ...(await importOriginal<any>()),
  displayPositionedDropdown: vi.fn(),
  queryItemAsString: vi.fn(),
  appDefaults: vi.fn(),
  useRouteQueryPersisted: vi.fn().mockImplementation(() => ref('resource-table')),
  useFileActions: vi.fn(),
  useFileActionsDelete: vi.fn(() => mock<ReturnType<typeof useFileActionsDelete>>())
}))

const spacesResources = [
  {
    id: '1',
    name: 'Some space',
    driveType: 'project',
    description: 'desc',
    path: '',
    type: 'space',
    isFolder: true,
    disabled: false,
    getDriveAliasAndItem: () => '1'
  },
  {
    id: '2',
    name: 'Some other space',
    driveType: 'project',
    description: 'desc',
    path: '',
    type: 'space',
    isFolder: true,
    disabled: false,
    getDriveAliasAndItem: () => '2'
  },
  {
    id: '3',
    name: 'Some disabled space',
    driveType: 'project',
    description: 'desc',
    path: '',
    type: 'space',
    isFolder: true,
    disabled: true,
    getDriveAliasAndItem: () => '2'
  }
] as unknown as SpaceResource[]

describe('Projects view', () => {
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
      const { wrapper } = getMountedWrapper()
      expect(wrapper.find('oc-spinner-stub').exists()).toBeTruthy()
    })
    it('shows the no-content-message after loading', async () => {
      const { wrapper } = getMountedWrapper()
      await (wrapper.vm as any).loadResourcesTask.last
      expect(wrapper.find('oc-spinner-stub').exists()).toBeFalsy()
      expect(wrapper.find('.no-content-message').exists()).toBeTruthy()
    })
    it('lists all available project spaces', async () => {
      const spaces = spacesResources
      const { wrapper } = getMountedWrapper({ spaces })
      await (wrapper.vm as any).loadResourcesTask.last
      expect(wrapper.html()).toMatchSnapshot()
      expect(wrapper.find('.no-content-message').exists()).toBeFalsy()
      expect(wrapper.find('.spaces-list').exists()).toBeTruthy()
    })
    it('shows only filtered spaces if filter applied', async () => {
      const { wrapper } = getMountedWrapper({ spaces: spacesResources })
      ;(wrapper.vm as any).filterTerm = 'Some other space'
      await nextTick()
      expect((wrapper.vm as any).items).toEqual([spacesResources[1]])
    })
    it('shows only enabled spaces if includeDisabled filter is not applied', async () => {
      const { wrapper } = getMountedWrapper({ spaces: spacesResources })
      await nextTick()
      expect((wrapper.vm as any).items.length).toEqual(2)
    })
    it('shows all spaces if includeDisabled filter is applied', async () => {
      const { wrapper } = getMountedWrapper({ spaces: spacesResources, includeDisabled: true })
      await nextTick()
      expect((wrapper.vm as any).items.length).toEqual(3)
    })
  })
  it('should display the "Create Space"-button when permission given', () => {
    const { wrapper } = getMountedWrapper({
      abilities: [{ action: 'create-all', subject: 'Drive' }],
      stubAppBar: false
    })
    expect(wrapper.find('create-space-stub').exists()).toBeTruthy()
  })
  describe('breadcrumbs', () => {
    it('shows "Spaces" when user cannot access vault', async () => {
      const { wrapper } = getMountedWrapper({ spaces: spacesResources })
      await (wrapper.vm as any).loadResourcesTask.last
      expect(wrapper.findComponent<typeof AppBar>('app-bar-stub').props().breadcrumbs[0].text).toBe(
        'Spaces'
      )
    })
    it('shows "Drive" when user can access vault and scope is not vault', async () => {
      const { wrapper } = getMountedWrapper({
        spaces: spacesResources,
        abilities: [{ action: 'read-all', subject: 'Vault' }],
        store: { capabilityState: { capabilities: { vault: { enabled: true } } } }
      })
      await (wrapper.vm as any).loadResourcesTask.last
      expect(wrapper.findComponent<typeof AppBar>('app-bar-stub').props().breadcrumbs[0].text).toBe(
        'Drive'
      )
    })
    it('shows "Vault" when user can access vault and scope is vault', async () => {
      const { wrapper } = getMountedWrapper({
        spaces: spacesResources,
        abilities: [{ action: 'read-all', subject: 'Vault' }],
        store: { capabilityState: { capabilities: { vault: { enabled: true } } } },
        currentRoute: mock<RouteLocation>({
          name: 'files-spaces-projects',
          params: { scope: 'vault' }
        })
      })
      await (wrapper.vm as any).loadResourcesTask.last
      expect(wrapper.findComponent<typeof AppBar>('app-bar-stub').props().breadcrumbs[0].text).toBe(
        'Vault'
      )
    })
  })
  it('should not pass selected resource as space to sidebar when driveType is not "project"', () => {
    const resource = mock<SpaceResource>({ id: 'selected-resource', driveType: 'personal' })
    const { wrapper } = getMountedWrapper({
      store: { resourcesStore: { resources: [resource], selectedIds: ['selected-resource'] } }
    })

    expect((wrapper.vm as any).selectedSpace).toStrictEqual(null)
  })
  it('should pass selected resource as space to sidebar when driveType is "project"', () => {
    const resource = mock<SpaceResource>({ id: 'selected-resource', driveType: 'project' })
    const { wrapper } = getMountedWrapper({
      store: { resourcesStore: { resources: [resource], selectedIds: ['selected-resource'] } }
    })

    expect((wrapper.vm as any).selectedSpace.id).toStrictEqual('selected-resource')
  })
})

function getMountedWrapper({
  mocks = {},
  spaces = [],
  abilities = [],
  stubAppBar = true,
  includeDisabled = false,
  store = {},
  currentRoute = mock<RouteLocation>({ name: 'files-spaces-projects' })
}: {
  mocks?: Record<string, unknown>
  spaces?: SpaceResource[]
  abilities?: AbilityRule[]
  stubAppBar?: boolean
  includeDisabled?: boolean
  store?: PiniaMockOptions
  currentRoute?: RouteLocation
} = {}) {
  const plugins = defaultPlugins({ abilities, piniaOptions: { spacesState: { spaces }, ...store } })

  vi.mocked(queryItemAsString).mockImplementation(() => includeDisabled.toString())

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
    ...defaultComponentMocks({ currentRoute }),
    ...(mocks && mocks)
  }

  return {
    mocks: defaultMocks,
    wrapper: mount(Projects, {
      global: {
        components: {
          AppBar,
          CreateSpace
        },
        plugins,
        mocks: defaultMocks,
        provide: defaultMocks,
        stubs: {
          ...defaultStubs,
          'space-context-actions': true,
          'app-bar': stubAppBar,
          CreateSpace: true
        }
      }
    })
  }
}
