import { DateTime } from 'luxon'
import ResourceTable from '../../../../src/components/FilesList/ResourceTable.vue'
import {
  extractDomSelector,
  IncomingShareResource,
  OutgoingShareResource,
  Resource,
  ShareTypes,
  SpaceResource
} from '@ownclouders/web-client'
import { defaultPlugins, mount, PartialComponentProps } from '@ownclouders/web-test-helpers'
import { CapabilityStore, useResourcesStore } from '../../../../src/composables/piniaStores'
import { useCanBeOpenedWithSecureView } from '../../../../src/composables/resources'
import { displayPositionedDropdown } from '../../../../src/helpers/contextMenuDropdown'
import { eventBus } from '../../../../src/services/eventBus'
import { SideBarEventTopics } from '../../../../src/composables/sideBar'
import { mock } from 'vitest-mock-extended'
import { computed } from 'vue'
import { Identity } from '@ownclouders/web-client/graph/generated'
import { describe } from 'vitest'
import { useFileActionsRename } from '../../../../src/composables/actions/files'
import { useSpaceActionsRename } from '../../../../src/composables/actions'
import { FileAction, SpaceAction } from '../../../../src/composables/actions/types'
import { ResourceIndicator } from '../../../../src/helpers/statusIndicators'

const mockUseEmbedMode = vi.fn().mockReturnValue({
  isLocationPicker: computed(() => false),
  isFilePicker: computed(() => false),
  isEnabled: computed(() => false)
})

vi.mock('../../../../src/helpers/contextMenuDropdown')
vi.mock('../../../../src/composables/embedMode', () => ({
  useEmbedMode: vi.fn().mockImplementation(() => mockUseEmbedMode())
}))

vi.mock('../../../../src/composables/resources', async (importOriginal) => ({
  ...(await importOriginal<any>()),
  useCanBeOpenedWithSecureView: vi.fn()
}))

vi.mock('../../../../src/composables/actions/files', async (importOriginal) => ({
  ...(await importOriginal<any>()),
  useFileActionsRename: vi.fn(),
  useFileActions: vi.fn().mockReturnValue({
    getDefaultAction: vi.fn().mockReturnValue({ handler: vi.fn() })
  })
}))

vi.mock('../../../../src/composables/actions', async (importOriginal) => ({
  ...(await importOriginal<any>()),
  useSpaceActionsRename: vi.fn()
}))

const router = {
  push: vi.fn(),
  afterEach: vi.fn(),
  currentRoute: {
    name: 'some-route-name',
    query: {},
    params: {
      driveAliasAndItem: ''
    }
  },
  resolve: (r: { name: string }) => {
    return { href: r.name }
  }
}

const getCurrentDate = () => {
  return DateTime.fromJSDate(new Date()).minus({ days: 1 }).toFormat('EEE, dd MMM yyyy HH:mm:ss')
}

const fields = ['name', 'size', 'mdate', 'sdate', 'ddate', 'actions', 'sharedBy', 'sharedWith']

const sharedWith = [
  {
    id: 'bob',
    displayName: 'Bob',
    shareType: ShareTypes.user.value
  },
  {
    id: 'marie',
    displayName: 'Marie',
    shareType: ShareTypes.user.value
  },
  {
    id: 'john',
    displayName: 'John Richards Emperor of long names',
    shareType: ShareTypes.user.value
  }
] as Array<{ shareType: number } & Identity>

const owner = {
  id: 'bob',
  displayName: 'Bob'
} as Resource['owner']

const sharedBy = [
  {
    id: 'bob',
    displayName: 'Bob'
  }
] as Identity[]

const indicators = [
  {
    id: 'files-sharing',
    label: 'Shared with other people',
    icon: 'group',
    handler: (resource) => alert(`Resource: ${resource.name}`)
  },
  {
    id: 'file-link',
    label: 'Shared via link',
    icon: 'link'
  }
] as ResourceIndicator[]

const resourcesWithAllFields = [
  {
    id: 'forest',
    spaceId: '1',
    driveId: 'forest',
    name: 'forest.jpg',
    path: 'images/nature/forest.jpg',
    extension: 'jpg',
    thumbnail: 'https://cdn.pixabay.com/photo/2015/09/09/16/05/forest-931706_960_720.jpg',
    indicators,
    isFolder: false,
    type: 'file',
    tags: ['space', 'tag', 'moon'],
    size: '111000234',
    mdate: getCurrentDate(),
    sdate: getCurrentDate(),
    ddate: getCurrentDate(),
    owner,
    sharedBy,
    sharedWith,
    hidden: false,
    syncEnabled: true,
    outgoing: false,
    shareRoles: [],
    sharePermissions: [],
    shareTypes: [],
    canRename: vi.fn(),
    getDomSelector: () => extractDomSelector('forest'),
    canDownload: () => true
  },
  {
    id: 'notes',
    spaceId: '1',
    driveId: 'notes',
    name: 'notes.txt',
    path: '/Documents/notes.txt',
    extension: 'txt',
    indicators,
    isFolder: false,
    type: 'file',
    size: 'big',
    tags: ['space', 'tag'],
    mdate: getCurrentDate(),
    sdate: getCurrentDate(),
    owner,
    sharedBy,
    sharedWith,
    hidden: false,
    syncEnabled: true,
    outgoing: false,
    shareRoles: [],
    sharePermissions: [],
    shareTypes: [],
    canRename: vi.fn(),
    getDomSelector: () => extractDomSelector('notes'),
    canDownload: () => true
  },
  {
    id: 'documents',
    spaceId: '1',
    driveId: 'documents',
    name: 'Documents',
    path: '/Documents',
    isFolder: true,
    type: 'folder',
    size: '-1',
    tags: [],
    mdate: getCurrentDate(),
    sdate: getCurrentDate(),
    owner,
    sharedBy,
    sharedWith,
    hidden: false,
    syncEnabled: true,
    outgoing: false,
    shareRoles: [],
    sharePermissions: [],
    shareTypes: [],
    canRename: vi.fn(),
    getDomSelector: () => extractDomSelector('documents'),
    canDownload: () => true,
    canListVersions: () => true
  },
  {
    id: 'another-one==',
    driveId: 'another-one==',
    name: 'Another one',
    path: '/Another one',
    isFolder: true,
    type: 'folder',
    size: '237895',
    mdate: getCurrentDate(),
    sdate: getCurrentDate(),
    owner,
    sharedBy,
    sharedWith,
    hidden: false,
    syncEnabled: true,
    outgoing: false,
    shareRoles: [],
    sharePermissions: [],
    shareTypes: [],
    tags: [],
    canRename: vi.fn(),
    getDomSelector: () => extractDomSelector('another-one=='),
    canDownload: () => true,
    canListVersions: () => true
  },
  {
    id: 'in-delete-queue==',
    driveId: 'another-one==',
    name: 'In delete queue',
    path: '/In delete queue',
    isFolder: true,
    type: 'folder',
    size: '237895',
    mdate: getCurrentDate(),
    sdate: getCurrentDate(),
    owner,
    sharedBy,
    sharedWith,
    hidden: false,
    syncEnabled: true,
    outgoing: false,
    shareRoles: [],
    sharePermissions: [],
    shareTypes: [],
    tags: [],
    canRename: vi.fn(),
    getDomSelector: () => extractDomSelector('in-delete-queue=='),
    canDownload: () => true,
    canListVersions: () => true
  }
] as IncomingShareResource[]

const processingResourcesWithAllFields = [
  {
    id: 'rainforest',
    driveId: 'rainforest',
    name: 'rainforest.jpg',
    path: 'images/nature/rainforest.jpg',
    extension: 'jpg',
    thumbnail: 'https://cdn.pixabay.com/photo/2015/09/09/16/05/forest-931706_960_720.jpg',
    indicators,
    isFolder: false,
    type: 'file',
    tags: ['space', 'tag', 'moon'],
    size: '111000234',
    mdate: getCurrentDate(),
    sdate: getCurrentDate(),
    ddate: getCurrentDate(),
    owner,
    sharedBy,
    sharedWith,
    hidden: false,
    syncEnabled: true,
    outgoing: false,
    shareRoles: [],
    sharePermissions: [],
    shareTypes: [],
    canRename: vi.fn(),
    getDomSelector: () => extractDomSelector('forest'),
    canDownload: () => true,
    canListVersions: () => true,
    processing: true,
    spaceId: '1'
  },
  {
    id: 'personalnotes',
    driveId: 'personalnotes',
    name: 'personalnotes.txt',
    path: '/Documents/personalnotes.txt',
    extension: 'txt',
    isFolder: false,
    type: 'file',
    size: 'big',
    tags: ['space', 'tag'],
    mdate: getCurrentDate(),
    sdate: getCurrentDate(),
    owner,
    sharedBy,
    sharedWith,
    hidden: false,
    syncEnabled: true,
    outgoing: false,
    shareRoles: [],
    sharePermissions: [],
    shareTypes: [],
    canRename: vi.fn(),
    getDomSelector: () => extractDomSelector('notes'),
    canDownload: () => true,
    canListVersions: () => true,
    processing: true,
    spaceId: '1'
  }
] as IncomingShareResource[]

describe('ResourceTable', () => {
  it('displays all known fields of the resources', () => {
    const { wrapper } = getMountedWrapper()
    for (const field of fields) {
      expect(wrapper.findAll('.oc-table-header-cell-' + field).length).toEqual(1)
      expect(wrapper.findAll('.oc-table-data-cell-' + field).length).toEqual(
        resourcesWithAllFields.length
      )
    }
  })

  it('accepts resourceDomId closure', () => {
    const { wrapper } = getMountedWrapper({
      props: {
        resourceDomSelector: (resource: Resource) => ['custom', resource.getDomSelector()].join('-')
      }
    })
    resourcesWithAllFields.slice(0, -1).forEach((resource) => {
      ;['.oc-tbody-tr', '#resource-table-select', '#context-menu-drop'].forEach((baseSelector) => {
        expect(
          wrapper.find([baseSelector, 'custom', resource.getDomSelector()].join('-')).exists()
        ).toBeTruthy()
      })
    })
  })

  it('formats the resource size to a human readable format', () => {
    const { wrapper } = getMountedWrapper()
    expect(wrapper.find('.oc-tbody-tr-forest .oc-table-data-cell-size').text()).toEqual('111 MB')
    expect(wrapper.find('.oc-tbody-tr-documents .oc-table-data-cell-size').text()).toEqual('--')
    expect(wrapper.find('.oc-tbody-tr-notes .oc-table-data-cell-size').text()).toEqual('?')
  })

  describe('resource selection', () => {
    it('adds resources to selection model via checkboxes', async () => {
      const { wrapper } = getMountedWrapper()
      await wrapper.find('.resource-table-select-all .oc-checkbox').trigger('click')
      await wrapper.find('.oc-tbody-tr-documents .oc-checkbox').trigger('click')
      expect(wrapper.emitted('update:selectedIds').length).toBe(2)
    })
    it('does not add resources that are disabled to selection model via checkboxes', async () => {
      const { wrapper } = getMountedWrapper({ addProcessingResources: true })
      await wrapper.find('.resource-table-select-all .oc-checkbox').trigger('click')
      expect(wrapper.emitted('update:selectedIds')[0][0]).not.toContain(
        processingResourcesWithAllFields[0].id
      )
      expect(wrapper.emitted('update:selectedIds')[0][0]).not.toContain(
        processingResourcesWithAllFields[1].id
      )
    })
    it('should replace checkbox with spinner when resource is in delete queue', () => {
      const { wrapper } = getMountedWrapper()
      expect(wrapper.find('.oc-tbody-tr-in-delete-queue .oc-checkbox').exists()).toBe(false)
      expect(wrapper.find('.oc-tbody-tr-in-delete-queue .oc-spinner').exists()).toBe(true)
    })

    describe('all rows already selected', () => {
      it('de-selects all resources via the select-all checkbox', async () => {
        const { wrapper } = getMountedWrapper({
          props: {
            selectedIds: resourcesWithAllFields.map((resource) => resource.id),
            resources: resourcesWithAllFields.slice(0, -1)
          }
        })

        await wrapper.find('.resource-table-select-all .oc-checkbox').trigger('click')
        expect(wrapper.emitted<string>('update:selectedIds')[0][0].length).toBe(0)
      })
    })

    describe('embed mode file types', () => {
      it('does not set folders disabled ', async () => {
        mockUseEmbedMode.mockReturnValue({
          isEnabled: computed(() => true),
          fileTypes: computed(() => ['txt'])
        })

        const { wrapper } = getMountedWrapper()
        const trFolder = await wrapper.find('.oc-tbody-tr-documents')
        expect(trFolder.attributes('class').includes('oc-table-disabled')).toBeFalsy()
      })
      it('sets files disabled if does not occur in it', async () => {
        mockUseEmbedMode.mockReturnValue({
          isEnabled: computed(() => true),
          fileTypes: computed(() => ['txt'])
        })

        const { wrapper } = getMountedWrapper()
        const trTxt = await wrapper.find('.oc-tbody-tr-notes')
        const trJpg = await wrapper.find('.oc-tbody-tr-forest')
        expect(trTxt.attributes('class').includes('oc-table-disabled')).toBeFalsy()
        expect(trJpg.attributes('class').includes('oc-table-disabled')).toBeTruthy()
      })
    })

    describe('embed mode location target', () => {
      it('should not hide checkboxes when embed mode does not have location as target', () => {
        mockUseEmbedMode.mockReturnValue({
          isLocationPicker: computed(() => false)
        })

        const { wrapper } = getMountedWrapper()

        expect(wrapper.find('.resource-table-select-all').exists()).toBe(true)
        expect(wrapper.find('.resource-table-select-all .oc-checkbox').exists()).toBe(true)
      })
      it('should not hide checkboxes when embed mode does not have file as target', () => {
        mockUseEmbedMode.mockReturnValue({
          isFilePicker: computed(() => false)
        })

        const { wrapper } = getMountedWrapper()

        expect(wrapper.find('.resource-table-select-all').exists()).toBe(true)
        expect(wrapper.find('.resource-table-select-all .oc-checkbox').exists()).toBe(true)
      })

      it('should hide checkboxes when embed mode has location as target', () => {
        mockUseEmbedMode.mockReturnValue({
          isLocationPicker: computed(() => true)
        })

        const { wrapper } = getMountedWrapper()

        expect(wrapper.find('.resource-table-select-all').exists()).toBe(false)
        expect(wrapper.find('.resource-table-select-all .oc-checkbox').exists()).toBe(false)
      })
      it('should hide checkboxes when embed mode has file as target', () => {
        mockUseEmbedMode.mockReturnValue({
          isFilePicker: computed(() => true)
        })

        const { wrapper } = getMountedWrapper()

        expect(wrapper.find('.resource-table-select-all').exists()).toBe(false)
        expect(wrapper.find('.resource-table-select-all .oc-checkbox').exists()).toBe(false)
      })
    })
  })

  describe('resource activation', () => {
    it('emits fileClick upon clicking on a resource name', async () => {
      const { wrapper } = getMountedWrapper()
      const tr = await wrapper.find('.oc-tbody-tr-forest .oc-resource-name')
      await tr.trigger('click')

      expect(
        wrapper.emitted<{ resources: Resource[] }[]>('fileClick')[0][0].resources[0].name
      ).toMatch('forest.jpg')
    })

    it('does not emit fileClick upon clicking on a disabled resource name', async () => {
      const { wrapper } = getMountedWrapper({ addProcessingResources: true })
      const tr = await wrapper.find('.oc-tbody-tr-rainforest .oc-resource-name')
      await tr.trigger('click')

      expect(wrapper.emitted().fileClick).toBeUndefined()
    })

    it('does not emit fileClick upon clicking on a resource when embed mode is enabled', async () => {
      mockUseEmbedMode.mockReturnValue({
        isEnabled: computed(() => true)
      })
      const { wrapper } = getMountedWrapper()
      const tr = await wrapper.find('.oc-tbody-tr-forest .oc-resource-name')
      await tr.trigger('click')
      expect(wrapper.emitted().fileClick).toBeUndefined()
    })

    it('does not emit fileClick event if file can not be opened via secure view', async () => {
      const { wrapper } = getMountedWrapper({
        canBeOpenedWithSecureView: false,
        resources: [
          {
            ...resourcesWithAllFields[0],
            isFolder: false,
            canDownload: () => false
          }
        ]
      })
      const tr = await wrapper.find('.oc-tbody-tr-forest .oc-resource-name')
      await tr.trigger('click')
      expect(wrapper.emitted().fileClick).toBeUndefined()
    })

    it('should emit "fileClick" when psec file is clicked', async () => {
      const resource = mock<Resource>({
        id: 'psec-file',
        name: 'psec-file.psec',
        path: '/psec-file.psec',
        remoteItemPath: '/psec-file.psec',
        canDownload: () => false,
        isFolder: false,
        getDomSelector: () => extractDomSelector('psec-file')
      })

      const { wrapper } = getMountedWrapper({
        canBeOpenedWithSecureView: false,
        resources: [resource]
      })
      await wrapper.find('.oc-tbody-tr-psec-file .oc-resource-name').trigger('click')

      expect(
        wrapper.emitted<{ resources: Resource[] }[]>('fileClick')[0][0].resources[0].name
      ).toMatch('psec-file.psec')
    })
  })

  describe('resource details', () => {
    it('emits select event when clicking on the row', async () => {
      const { wrapper } = getMountedWrapper()
      const setLastSelectedIdSpy = vi.spyOn(useResourcesStore(), 'setLastSelectedId')
      const tableRow = await wrapper.find('.oc-tbody-tr-forest .oc-table-data-cell-size')
      await tableRow.trigger('click')
      expect(wrapper.emitted('update:selectedIds')).toBeTruthy()
      expect(setLastSelectedIdSpy).toHaveBeenCalledWith('forest')
    })

    it('does not emit select event when clicking on the row of a disabled resource', async () => {
      const { wrapper } = getMountedWrapper({ addProcessingResources: true })
      const tableRow = await wrapper.find('.oc-tbody-tr-rainforest .oc-table-data-cell-size')
      await tableRow.trigger('click')

      expect(wrapper.emitted('update:selectedIds')).toBeUndefined()
    })

    it('should disable resource when it is in delete queue', () => {
      const { wrapper } = getMountedWrapper()
      expect(wrapper.find('.oc-tbody-tr-in-delete-queue.oc-table-disabled').exists()).toBe(true)
    })

    it('should not emit select event when clicking on the row of an already selected resource', async () => {
      const { wrapper } = getMountedWrapper({ props: { selectedIds: ['forest'] } })
      const tableRow = wrapper.find('.oc-tbody-tr-forest .oc-table-data-cell-size')
      await tableRow.trigger('click')

      expect(wrapper.emitted('update:selectedIds')).toBeUndefined()
    })
  })

  describe('context menu', () => {
    it('emits select event on contextmenu click', async () => {
      const spyDisplayPositionedDropdown = vi.mocked(displayPositionedDropdown)
      const { wrapper } = getMountedWrapper()
      await wrapper.find('.oc-tbody-tr').trigger('contextmenu')
      expect(wrapper.emitted('update:selectedIds').length).toBe(1)
      expect(spyDisplayPositionedDropdown).toHaveBeenCalledTimes(1)
    })

    it('does not emit select event on contextmenu click of disabled resource', async () => {
      const spyDisplayPositionedDropdown = vi.mocked(displayPositionedDropdown)
      const { wrapper } = getMountedWrapper({ addProcessingResources: true })
      await wrapper.find('.oc-tbody-tr-rainforest').trigger('contextmenu')
      expect(wrapper.emitted('update:selectedIds')).toBeUndefined()
      expect(spyDisplayPositionedDropdown).not.toHaveBeenCalled()
    })

    it('emits select event on clicking the three-dot icon in table row', async () => {
      const spyDisplayPositionedDropdown = vi.mocked(displayPositionedDropdown)
      const { wrapper } = getMountedWrapper()
      await wrapper
        .find('.oc-table-data-cell-actions .resource-table-btn-action-dropdown')
        .trigger('click')
      expect(wrapper.emitted('update:selectedIds').length).toBe(1)
      expect(spyDisplayPositionedDropdown).toHaveBeenCalledTimes(1)
    })

    it('does not show the three-dot icon in table row of a disabled resource', () => {
      const { wrapper } = getMountedWrapper({ addProcessingResources: true })
      expect(
        wrapper
          .find(
            '.oc-tbody-tr-rainforest .oc-table-data-cell-actions .resource-table-btn-action-dropdown'
          )
          .exists()
      ).toBeFalsy()
    })

    it('removes invalid chars from item ids for usage in html template', async () => {
      const { wrapper } = getMountedWrapper()
      const contextMenuTriggers = await wrapper.findAll('.resource-table-btn-action-dropdown')
      for (let i = 0; i < contextMenuTriggers.length; i++) {
        const id = contextMenuTriggers.at(i).attributes().id
        expect(id).not.toBeUndefined()
        expect(id).toEqual(expect.not.stringContaining('='))
      }
    })
  })

  describe('hover effect', () => {
    it('is disabled by default', () => {
      const { wrapper } = getMountedWrapper({ props: { hover: false } })
      expect(wrapper.find('table').attributes('class')).not.toContain('oc-table-hover')
    })

    it('can be enabled', () => {
      const { wrapper } = getMountedWrapper({ props: { hover: true } })
      expect(wrapper.find('table').attributes('class')).toContain('oc-table-hover')
    })
  })

  describe('tags', () => {
    describe('inline', () => {
      it.each([
        { tags: [], tagCount: 0 },
        { tags: ['1'], tagCount: 1 },
        { tags: ['1', '2'], tagCount: 2 },
        { tags: ['1', '2', '3'], tagCount: 2 },
        { tags: ['1', '2', '3', '4'], tagCount: 2 }
      ])('render 2 tags max', (data) => {
        const { tags, tagCount } = data
        const resource = mock<Resource>({ id: '1', tags })
        const { wrapper } = getMountedWrapper({ props: { resources: [resource] } })
        const resourceRow = wrapper.find(`[data-item-id="${resource.id}"]`)
        expect(resourceRow.findAll('.resource-table-tag').length).toBe(tagCount)
      })
      it('render router link if user is authenticated', () => {
        const resource = mock<Resource>({ id: '1', tags: ['1'] })
        const { wrapper } = getMountedWrapper({ props: { resources: [resource] } })
        const resourceRow = wrapper.find(`[data-item-id="${resource.id}"]`)
        expect(resourceRow.find('.resource-table-tag-wrapper').element.tagName).toEqual(
          'ROUTER-LINK-STUB'
        )
      })
      it('do not render router link if user is not authenticated', () => {
        const resource = mock<Resource>({ id: '1', tags: ['1'] })
        const { wrapper } = getMountedWrapper({
          props: { resources: [resource] },
          userContextReady: false
        })
        const resourceRow = wrapper.find(`[data-item-id="${resource.id}"]`)
        expect(resourceRow.find('.resource-table-tag-wrapper').element.tagName).toEqual('SPAN')
      })
    })
    describe('"more"-button', () => {
      it.each([
        { tags: [], renderButton: false },
        { tags: ['1'], renderButton: false },
        { tags: ['1', '2'], renderButton: false },
        { tags: ['1', '2', '3'], renderButton: true },
        { tags: ['1', '2', '3', '4'], renderButton: true }
      ])('does only render when the resource has 3 tags or more', (data) => {
        const { tags, renderButton } = data
        const resource = mock<Resource>({ id: '1', tags })
        const { wrapper } = getMountedWrapper({ props: { resources: [resource] } })
        const resourceRow = wrapper.find(`[data-item-id="${resource.id}"]`)
        expect(resourceRow.find('.resource-table-tag-more').exists()).toBe(renderButton)
      })
      it('opens sidebar on click', async () => {
        const spyBus = vi.spyOn(eventBus, 'publish')
        const resource = mock<Resource>({ id: '1', tags: ['1', '2', '3'] })
        const { wrapper } = getMountedWrapper({ props: { resources: [resource] } })
        const resourceRow = wrapper.find(`[data-item-id="${resource.id}"]`)
        await resourceRow.find('.resource-table-tag-more').trigger('click')
        expect(spyBus).toHaveBeenCalledWith(SideBarEventTopics.open)
      })
    })
  })
  describe('"shared with" field', () => {
    it('only displays authenticated shares', () => {
      const resource = mock<OutgoingShareResource>({ id: '1' })
      resource.sharedWith = [
        {
          id: 'bob',
          displayName: 'Bob',
          shareType: ShareTypes.user.value
        },
        { displayName: 'Link', shareType: ShareTypes.link.value }
      ]

      const { wrapper } = getMountedWrapper({ resources: [resource] })

      expect(wrapper.find('.resource-table-shared-with').exists()).toBeTruthy()
      expect(wrapper.findAll('.resource-table-shared-with .oc-avatar').length).toBe(1)
    })
  })
  describe('rename action', () => {
    it('shows if available', () => {
      const { wrapper } = getMountedWrapper()
      expect(wrapper.find('.resource-table-edit-name').exists()).toBeTruthy()
    })
    it('does not show if not available', () => {
      const { wrapper } = getMountedWrapper({ hasRenameAction: false })
      expect(wrapper.find('.resource-table-edit-name').exists()).toBeFalsy()
    })
  })

  describe('getRenameButtonAriaLabel', () => {
    it('returns "Rename space" for space resources', () => {
      const spaceResource = mock<SpaceResource>({
        id: 'space-1',
        name: 'My Space',
        path: '/My Space',
        type: 'space',
        driveType: 'project',
        getDomSelector: () => 'space-1',
        canDownload: () => true
      })
      const { wrapper } = getMountedWrapper({ resources: [spaceResource] })
      const editButton = wrapper.find('.resource-table-edit-name')
      expect(editButton.attributes('aria-label')).toBe('Rename space')
    })

    it('returns "Rename folder" for folder resources', () => {
      const folderResource = mock<Resource>({
        id: 'folder-1',
        name: 'My Folder',
        path: '/My Folder',
        isFolder: true,
        type: 'folder',
        getDomSelector: () => 'folder-1',
        canDownload: () => true
      })
      const { wrapper } = getMountedWrapper({ resources: [folderResource] })
      const editButton = wrapper.find('.resource-table-edit-name')
      expect(editButton.attributes('aria-label')).toBe('Rename folder')
    })

    it('returns "Rename file" for file resources', () => {
      const fileResource = mock<Resource>({
        id: 'file-1',
        name: 'document.txt',
        path: '/document.txt',
        isFolder: false,
        type: 'file',
        getDomSelector: () => 'file-1',
        canDownload: () => true
      })
      const { wrapper } = getMountedWrapper({ resources: [fileResource] })
      const editButton = wrapper.find('.resource-table-edit-name')
      expect(editButton.attributes('aria-label')).toBe('Rename file')
    })
  })
})

function getMountedWrapper({
  props = {},
  userContextReady = true,
  addProcessingResources = false,
  canBeOpenedWithSecureView = true,
  hasRenameAction = true,
  resources = resourcesWithAllFields
}: {
  props?: PartialComponentProps<typeof ResourceTable>
  userContextReady?: boolean
  addProcessingResources?: boolean
  canBeOpenedWithSecureView?: boolean
  hasRenameAction?: boolean
  resources?: Resource[]
} = {}) {
  const capabilities = {
    files: { tags: true }
  } satisfies Partial<CapabilityStore['capabilities']>

  vi.mocked(useCanBeOpenedWithSecureView).mockReturnValue({
    canBeOpenedWithSecureView: () => canBeOpenedWithSecureView
  })

  const useFileActionsRenameMock = mock<ReturnType<typeof useFileActionsRename>>()
  useFileActionsRenameMock.actions = computed(() => [
    mock<FileAction>({ isVisible: () => hasRenameAction })
  ])
  vi.mocked(useFileActionsRename).mockReturnValue(useFileActionsRenameMock)

  const useSpaceActionsRenameMock = mock<ReturnType<typeof useSpaceActionsRename>>()
  useSpaceActionsRenameMock.actions = computed(() => [
    mock<SpaceAction>({ isVisible: () => hasRenameAction })
  ])
  vi.mocked(useSpaceActionsRename).mockReturnValue(useSpaceActionsRenameMock)

  return {
    wrapper: mount(ResourceTable, {
      props: {
        resources: [
          ...resources,
          ...(addProcessingResources ? processingResourcesWithAllFields : [])
        ],
        selection: [],
        hover: false,
        lazy: false,
        space: mock<SpaceResource>({
          getDriveAliasAndItem: vi.fn()
        }),
        ...props
      },
      global: {
        renderStubDefaultSlot: true,
        plugins: [
          ...defaultPlugins({
            piniaOptions: {
              authState: { userContextReady },
              capabilityState: { capabilities },
              resourcesStore: { deleteQueue: ['in-delete-queue=='] }
            }
          })
        ],
        stubs: {
          OcButton: false,
          'router-link': true
        },
        mocks: {
          $route: router.currentRoute,
          $router: router
        },
        provide: {
          $router: router,
          $route: router.currentRoute
        }
      }
    })
  }
}
