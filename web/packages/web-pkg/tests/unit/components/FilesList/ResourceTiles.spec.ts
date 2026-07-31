import { defaultComponentMocks, defaultPlugins, mount } from '@ownclouders/web-test-helpers'
import ResourceTiles from '../../../../src/components/FilesList/ResourceTiles.vue'
import { sortFields } from '../../../../src/helpers/ui/resourceTiles'
import { Resource, SpaceResource } from '@ownclouders/web-client'
import { mock } from 'vitest-mock-extended'
import { ComponentPublicInstance, computed } from 'vue'
import { extractDomSelector } from '@ownclouders/web-client'
import { useCanBeOpenedWithSecureView } from '../../../../src/composables/resources'
import { displayPositionedDropdown } from '../../../../src/helpers/contextMenuDropdown'
import { OcSelect } from '@ownclouders/design-system/components'
import { ResourceIndicator } from '../../../../src/helpers/statusIndicators'

vi.mock('../../../../src/helpers/contextMenuDropdown')
vi.mock('../../../../src/composables/viewMode', async (importOriginal) => ({
  ...(await importOriginal<any>()),
  useTileSize: vi.fn().mockReturnValue({
    calculateTileSizePixels: vi.fn().mockImplementation((viewSize: number) => 100 * viewSize)
  })
}))

const mockUseEmbedMode = vi.fn().mockReturnValue({ isEnabled: computed(() => false) })
vi.mock('../../../../src/composables/embedMode', () => ({
  useEmbedMode: vi.fn().mockImplementation(() => mockUseEmbedMode())
}))

vi.mock('../../../../src/composables/resources', async (importOriginal) => ({
  ...(await importOriginal<any>()),
  useCanBeOpenedWithSecureView: vi.fn()
}))

vi.mock('../../../../src/composables/actions/files', async (importOriginal) => ({
  ...(await importOriginal<any>()),
  useFileActions: vi.fn().mockReturnValue({
    getDefaultAction: vi.fn().mockReturnValue({ handler: vi.fn() })
  })
}))

const spacesResources = [
  {
    id: '1',
    name: 'Space 1',
    path: '',
    type: 'space',
    isFolder: true,
    indicators: [],
    getDriveAliasAndItem: () => '1'
  },
  {
    id: '2',
    name: 'Space 2',
    path: '',
    type: 'space',
    isFolder: true,
    indicators: [],
    getDriveAliasAndItem: () => '2'
  }
] as unknown as SpaceResource[]

const resources = [
  {
    id: 'forest',
    spaceId: '1',
    driveId: 'forest',
    name: 'forest.jpg',
    path: 'images/nature/forest.jpg',
    extension: 'jpg',
    thumbnail: 'https://cdn.pixabay.com/photo/2015/09/09/16/05/forest-931706_960_720.jpg',
    isFolder: false,
    indicators: [] as ResourceIndicator[],
    type: 'file',
    tags: ['space', 'tag', 'moon'],
    size: '111000234',
    hidden: false,
    syncEnabled: true,
    outgoing: false,
    canRename: vi.fn(),
    getDomSelector: () => extractDomSelector('forest'),
    canDownload: () => true
  } as Resource
]

describe('ResourceTiles component', () => {
  const originalGetElementById = document.getElementById
  const originalGetComputedStyle = window.getComputedStyle
  beforeEach(() => {
    const mockElement = {
      clientWidth: 800
    } as HTMLElement
    document.getElementById = vi.fn((id) => {
      if (id === 'tiles-view') {
        return mockElement
      }
      return originalGetElementById.call(document, id)
    })
    window.getComputedStyle = vi.fn().mockImplementation(() => {
      return {
        getPropertyValue: (propName: string) => {
          switch (propName) {
            case '--oc-size-tiles-default':
              return '9rem'
            case '--oc-size-tiles-resize-step':
              return '9rem'
            default:
              return originalGetComputedStyle(document.documentElement).getPropertyValue(propName)
          }
        },
        fontSize: '14px'
      }
    })
  })
  it('renders an array of spaces correctly', async () => {
    const { wrapper } = getWrapper({ props: { resources: spacesResources } })
    await wrapper.vm.$nextTick()
    expect(wrapper.html()).toMatchSnapshot()
  })

  it('renders a footer slot', () => {
    const { wrapper } = getWrapper({ slots: { footer: 'Hello, ResourceTiles footer!' } })
    expect(wrapper.html()).toMatchSnapshot()
  })

  describe('file click', () => {
    it('emits fileClick event upon click on tile', async () => {
      const { wrapper } = getWrapper({ props: { resources } })
      await wrapper.find('.oc-tiles-item .oc-resource-name').trigger('click')
      expect(
        wrapper.emitted<{ resources: Resource[] }[]>('fileClick')[0][0].resources[0].name
      ).toMatch('forest.jpg')
    })

    it('does not emit fileClick event upon click on tile when embed mode is enabled', async () => {
      mockUseEmbedMode.mockReturnValue({
        isEnabled: computed(() => true)
      })
      const { wrapper } = getWrapper({ props: { resources } })
      await wrapper.find('.oc-tiles-item .oc-resource-name').trigger('click')
      expect(wrapper.emitted().fileClick).toBeUndefined()
    })

    it('does not emit fileClick event if file can not be opened via secure view', async () => {
      const { wrapper } = getWrapper({
        canBeOpenedWithSecureView: false,
        props: {
          resources: [{ ...resources[0], isFolder: false, canDownload: () => false }]
        }
      })
      await wrapper.find('.oc-tiles-item .oc-resource-name').trigger('click')
      expect(wrapper.emitted().fileClick).toBeUndefined()
    })
  })

  it('emits update:selectedIds event on resource selection and sets the selection', () => {
    const { wrapper } = getWrapper({
      props: {
        resources: spacesResources,
        selectedIds: [spacesResources[0].id]
      }
    })
    ;(wrapper.vm as any).toggleSelection(spacesResources[0])
    expect(
      wrapper.findComponent({ name: 'resource-tile' }).props('isResourceSelected')
    ).toBeTruthy()
    expect(wrapper.emitted('update:selectedIds')).toBeTruthy()
  })

  describe('sorting', () => {
    it('renders the label of the first sort field as default', () => {
      const { wrapper } = getWrapper({ props: { sortFields } })
      expect(wrapper.find('.vs__selected').text()).toEqual(sortFields[0].label)
    })
    it('renders the label of the current sort field as default', () => {
      const sortField = sortFields[2]
      const { wrapper } = getWrapper({
        props: {
          sortFields,
          sortBy: sortField.name,
          sortDir: sortField.sortDir
        }
      })
      expect(wrapper.find('.vs__selected').text()).toEqual(sortField.label)
    })
    it('emits the "sort"-event', () => {
      const { wrapper } = getWrapper({ props: { sortFields }, stubs: { OcSelect: true } })
      wrapper.findComponent<typeof OcSelect>('oc-select-stub').vm.$emit('update:modelValue', 1)
      expect(wrapper.emitted('sort')).toBeTruthy()
    })
  })
  describe('drag and drop', () => {
    it('emits the "update:selectedIds"-event on drag start', async () => {
      const { wrapper } = getWrapper()
      ;(wrapper.vm as any).dragItem = mock<Resource>()
      await wrapper.vm.$nextTick()
      ;(wrapper.vm.$refs.ghostElementRef as ComponentPublicInstance<unknown>).$el = { style: {} }
      ;(wrapper.vm as any).dragStart(mock<Resource>(), {
        dataTransfer: { setDragImage: vi.fn() }
      } as unknown as DragEvent)
      expect(wrapper.emitted('update:selectedIds')).toBeDefined()
    })
    it('emits the "fileDropped"-event on resource drop', () => {
      const { wrapper } = getWrapper()
      ;(wrapper.vm as any).fileDropped(mock<Resource>(), { dataTransfer: {} } as DragEvent)
      expect(wrapper.emitted('fileDropped')).toBeDefined()
    })
  })
  describe('context menu', () => {
    it('triggers the positioned dropdown on click', async () => {
      const spyDisplayPositionedDropdown = vi.mocked(displayPositionedDropdown)
      const { wrapper } = getWrapper({ props: { resources } })
      const btn = wrapper.find('.resource-tiles-btn-action-dropdown')
      await btn.trigger('click')
      expect(spyDisplayPositionedDropdown).toHaveBeenCalled()
    })
    it('does not show for disabled resources', () => {
      const { wrapper } = getWrapper({
        props: { resources: [{ ...resources[0], processing: true }] }
      })
      expect(wrapper.find('.resource-tiles-btn-action-dropdown').exists()).toBeFalsy()
    })
  })
  describe('checkboxes', () => {
    it('update the selected ids on click', async () => {
      const { wrapper } = getWrapper({ props: { resources } })
      const checkbox = wrapper.find('input[type="checkbox"]')
      await checkbox.trigger('click')
      expect(wrapper.emitted('update:selectedIds')).toBeTruthy()
    })
    it('are disabled for disabled resources', () => {
      const { wrapper } = getWrapper({
        props: { resources: [{ ...resources[0], processing: true }] }
      })
      const checkbox = wrapper.find('input[type="checkbox"]')
      expect(Object.keys(checkbox.attributes())).toContain('disabled')
    })
  })

  it.each([
    { viewSize: 1, expected: 'xlarge' },
    { viewSize: 2, expected: 'xlarge' },
    { viewSize: 3, expected: 'xxlarge' },
    { viewSize: 4, expected: 'xxlarge' },
    { viewSize: 5, expected: 'xxxlarge' },
    { viewSize: 6, expected: 'xxxlarge' }
  ])('passes the "viewSize" to the OcTile component', async (data) => {
    const { wrapper } = getWrapper({
      props: { resources: spacesResources, viewSize: data.viewSize }
    })
    await wrapper.vm.$nextTick()
    expect(wrapper.findComponent({ name: 'resource-tile' }).props('resourceIconSize')).toEqual(
      data.expected
    )
  })

  it('should make resource clickable when it is a password protected folder', async () => {
    const { wrapper } = getWrapper({
      props: {
        resources: [
          {
            id: 'protected-folder',
            driveId: 'protected-folder',
            name: 'protected-folder.psec',
            path: '/protected-folder.psec',
            extension: 'psec',
            isFolder: false,
            indicators: [] as ResourceIndicator[],
            type: 'file',
            tags: ['space', 'tag', 'moon'],
            size: '111000234',
            hidden: false,
            syncEnabled: true,
            outgoing: false,
            canRename: () => false,
            getDomSelector: () => extractDomSelector('protected-folder'),
            canDownload: () => false
          }
        ]
      }
    })

    const resourceLink = wrapper.find('.oc-resource-link')
    await resourceLink.trigger('click')

    expect(wrapper.emitted('fileClick')).toBeTruthy()
  })

  it('should not make resource clickable when it is not a password protected folder and does not have enough permissions', () => {
    const { wrapper } = getWrapper({
      props: {
        resources: [
          {
            id: 'forest',
            driveId: 'forest',
            name: 'forest.jpg',
            path: 'images/nature/forest.jpg',
            extension: 'jpg',
            thumbnail: 'https://cdn.pixabay.com/photo/2015/09/09/16/05/forest-931706_960_720.jpg',
            isFolder: false,
            indicators: [] as ResourceIndicator[],
            type: 'file',
            tags: ['space', 'tag', 'moon'],
            size: '111000234',
            hidden: false,
            syncEnabled: true,
            outgoing: false,
            canRename: false,
            getDomSelector: () => extractDomSelector('forest'),
            canDownload: () => false
          }
        ]
      }
    })

    expect(wrapper.find('.oc-resource-link').exists()).toBeFalsy()
  })

  function getWrapper({
    props = {},
    slots = {},
    stubs = {},
    canBeOpenedWithSecureView = true
  } = {}) {
    const mocks = defaultComponentMocks()

    vi.mocked(useCanBeOpenedWithSecureView).mockReturnValue({
      canBeOpenedWithSecureView: () => canBeOpenedWithSecureView
    })

    return {
      wrapper: mount(ResourceTiles, {
        props: {
          lazy: false,
          viewSize: 1,
          ...props
        },
        slots: {
          ...slots
        },
        global: {
          plugins: [...defaultPlugins()],
          mocks: mocks,
          provide: mocks,
          stubs
        }
      })
    }
  }
})
