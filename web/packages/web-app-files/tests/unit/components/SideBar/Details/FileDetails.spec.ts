import FileDetails from '../../../../../src/components/SideBar/Details/FileDetails.vue'
import { Resource, ShareResource, ShareTypes } from '@ownclouders/web-client'
import {
  mount,
  defaultComponentMocks,
  defaultPlugins,
  RouteLocation
} from '@ownclouders/web-test-helpers'
import { mock, mockDeep } from 'vitest-mock-extended'
import { SpaceResource } from '@ownclouders/web-client'
import { AncestorMetaData } from '@ownclouders/web-pkg/'
import { User } from '@ownclouders/web-client/graph/generated'

const getResourceMock = ({
  type = 'file',
  mimeType = 'image/jpeg',
  tags = [],
  thumbnail = null,
  shareTypes = [],
  path = '/somePath/someResource',
  locked = false,
  canEditTags = true,
  sharedBy = undefined
} = {}) =>
  mock<ShareResource>({
    id: '1',
    type,
    isFolder: type === 'folder',
    mimeType,
    owner: {
      id: 'marie',
      displayName: 'Marie'
    },
    sharedBy,
    mdate: 'Wed, 21 Oct 2015 07:28:00 GMT',
    tags,
    size: '740',
    path,
    thumbnail,
    shareTypes,
    locked,
    canEditTags: vi.fn(() => canEditTags),
    ...(sharedBy && { sharedWith: [] })
  })

const selectors = {
  ownerDisplayName: '[data-testid="ownerDisplayName"]',
  preview: '[data-testid="preview"]',
  resourceIcon: '.details-icon',
  lockedBy: '[data-testid="locked-by"]',
  sharedBy: '[data-testid="shared-by"]',
  sharedVia: '[data-testid="shared-via"]',
  sharingInfo: '[data-testid="sharingInfo"]',
  sizeInfo: '[data-testid="sizeInfo"]',
  tags: '[data-testid="tags"]',
  timestamp: '[data-testid="timestamp"]',
  versionsInfo: '[data-testid="versionsInfo"]'
}

describe('Details SideBar Panel', () => {
  describe('preview', () => {
    describe('shows preview area', () => {
      it('while trying to load a preview', () => {
        const resource = getResourceMock()
        const { wrapper } = createWrapper({ resource })
        expect(wrapper.find(selectors.preview).exists()).toBeTruthy()
        expect(wrapper.find(selectors.resourceIcon).exists()).toBeFalsy()
      })
      it('for allowed mime types', () => {
        const resource = getResourceMock()
        const { wrapper } = createWrapper({ resource })
        expect(wrapper.find(selectors.preview).exists()).toBeTruthy()
        expect(wrapper.find(selectors.resourceIcon).exists()).toBeFalsy()
      })
    })
  })
  describe('status indicators', () => {
    it('show if given on non-public page', () => {
      const resource = getResourceMock({ shareTypes: [ShareTypes.user.value] })
      const { wrapper } = createWrapper({ resource })
      expect(wrapper.find(selectors.sharingInfo).exists()).toBeTruthy()
    })
    it('do not show on a public page', () => {
      const resource = getResourceMock({ shareTypes: [ShareTypes.user.value] })
      const { wrapper } = createWrapper({ resource, isPublicLinkContext: true })
      expect(wrapper.find(selectors.sharingInfo).exists()).toBeFalsy()
    })
  })
  describe('timestamp', () => {
    it('shows if given', () => {
      const resource = getResourceMock()
      const { wrapper } = createWrapper({ resource })
      expect(wrapper.find(selectors.timestamp).exists()).toBeTruthy()
    })
  })
  describe('locked by', () => {
    it('shows if the resource is locked', () => {
      const resource = getResourceMock({ locked: true })
      const { wrapper } = createWrapper({ resource })
      expect(wrapper.find(selectors.lockedBy).exists()).toBeTruthy()
    })
  })
  describe('shared via', () => {
    it('shows if the resource has an indirect share', () => {
      const resource = getResourceMock()
      const ancestorMetaData = {
        '/somePath': { path: '/somePath', shareTypes: [ShareTypes.user.value] }
      } as unknown as AncestorMetaData
      const { wrapper } = createWrapper({ resource, ancestorMetaData })
      expect(wrapper.find(selectors.sharedVia).exists()).toBeTruthy()
    })
  })
  describe('shared by', () => {
    it('shows if the resource is a share from another user', () => {
      const resource = getResourceMock({
        shareTypes: [ShareTypes.user.value],
        sharedBy: [{ id: '1', displayName: 'Marie' }]
      })
      const { wrapper } = createWrapper({
        resource,
        user: { onPremisesSamAccountName: 'einstein' } as User
      })
      expect(wrapper.find(selectors.sharedBy).exists()).toBeTruthy()
    })
  })
  describe('owner display name', () => {
    it('shows if given', () => {
      const resource = getResourceMock()
      const { wrapper } = createWrapper({ resource })
      expect(wrapper.find(selectors.ownerDisplayName).exists()).toBeTruthy()
    })
  })
  describe('size', () => {
    it('shows if given', () => {
      const resource = getResourceMock()
      const { wrapper } = createWrapper({ resource })
      expect(wrapper.find(selectors.sizeInfo).exists()).toBeTruthy()
    })
  })
  describe('versions', () => {
    it('show if given for files on a private page', async () => {
      const resource = getResourceMock()
      const { wrapper } = createWrapper({ resource, versions: ['1'] })
      await wrapper.vm.$nextTick()
      await wrapper.vm.$nextTick()
      await wrapper.vm.$nextTick()
      expect(wrapper.find(selectors.versionsInfo).exists()).toBeTruthy()
    })
    it('do not show for folders on a private page', () => {
      const resource = getResourceMock({ type: 'folder' })
      const { wrapper } = createWrapper({ resource, versions: ['1'] })
      expect(wrapper.find(selectors.versionsInfo).exists()).toBeFalsy()
    })
    it('do not show on public pages', () => {
      const resource = getResourceMock()
      const { wrapper } = createWrapper({ resource, versions: ['1'], isPublicLinkContext: true })
      expect(wrapper.find(selectors.versionsInfo).exists()).toBeFalsy()
    })
  })

  describe('tags', () => {
    it('shows when enabled via capabilities', () => {
      const resource = getResourceMock()
      const { wrapper } = createWrapper({ resource })
      expect(wrapper.find(selectors.tags).exists()).toBeTruthy()
    })
    it('does not show when disabled via capabilities', () => {
      const resource = getResourceMock()
      const { wrapper } = createWrapper({ resource, tagsEnabled: false })
      expect(wrapper.find(selectors.tags).exists()).toBeFalsy()
    })
    it('does not show for root folders', () => {
      const resource = getResourceMock({ path: '/' })
      const { wrapper } = createWrapper({ resource })
      expect(wrapper.find(selectors.tags).exists()).toBeTruthy()
    })
    it('shows as disabled when permission not set', () => {
      const resource = getResourceMock({ canEditTags: false })
      const { wrapper } = createWrapper({ resource })
      expect(wrapper.find(selectors.tags).find('.vs--disabled ').exists()).toBeTruthy()
    })
    it('should use router-link on private page', async () => {
      const resource = getResourceMock({ tags: ['moon', 'mars'] })
      const { wrapper } = createWrapper({ resource })
      await wrapper.vm.$nextTick()
      expect(wrapper.find(selectors.tags).find('router-link-stub').exists()).toBeTruthy()
    })
    it('should not use router-link on public page', async () => {
      const resource = getResourceMock({ tags: ['moon', 'mars'] })
      const { wrapper } = createWrapper({ resource, isPublicLinkContext: true })
      await wrapper.vm.$nextTick()
      expect(wrapper.find(selectors.tags).find('router-link-stub').exists()).toBeFalsy()
    })
  })
})

function createWrapper({
  resource = null,
  isPublicLinkContext = false,
  ancestorMetaData = {},
  user = { onPremisesSamAccountName: 'marie' } as User,
  versions = [],
  tagsEnabled = true
}: {
  resource?: Resource
  isPublicLinkContext?: boolean
  ancestorMetaData?: AncestorMetaData
  user?: User
  versions?: string[]
  tagsEnabled?: boolean
} = {}) {
  const currentRouteName = isPublicLinkContext ? 'files-public-link' : 'files-spaces-generic'
  const mocks = defaultComponentMocks({
    currentRoute: mock<RouteLocation>({ name: currentRouteName })
  })
  const capabilities = { files: { tags: tagsEnabled } }
  return {
    wrapper: mount(FileDetails, {
      global: {
        stubs: { 'router-link': true, 'resource-icon': true },
        provide: {
          ...mocks,
          versions,
          resource,
          space: mockDeep<SpaceResource>({ driveType: 'personal', isOwner: () => true })
        },
        plugins: [
          ...defaultPlugins({
            piniaOptions: {
              userState: { user },
              authState: { publicLinkContextReady: isPublicLinkContext },
              capabilityState: { capabilities },
              resourcesStore: { ancestorMetaData, currentFolder: mock<Resource>() }
            }
          })
        ],
        mocks
      }
    })
  }
}
