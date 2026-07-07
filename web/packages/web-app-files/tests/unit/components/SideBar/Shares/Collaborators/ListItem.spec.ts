import ListItem from '../../../../../../src/components/SideBar/Shares/Collaborators/ListItem.vue'
import {
  CollaboratorShare,
  GraphSharePermission,
  ShareRole,
  ShareTypes
} from '@ownclouders/web-client'
import {
  defaultPlugins,
  mount,
  defaultStubs,
  defaultComponentMocks,
  flushPromises
} from '@ownclouders/web-test-helpers'
import { useMessages, useSharesStore } from '@ownclouders/web-pkg'
import EditDropdown from '../../../../../../src/components/SideBar/Shares/Collaborators/EditDropdown.vue'
import RoleDropdown from '../../../../../../src/components/SideBar/Shares/Collaborators/RoleDropdown.vue'
import { mock } from 'vitest-mock-extended'
import { Mock } from 'vitest'
import { Resource, SpaceResource } from '@ownclouders/web-client'
import { RouteLocationNamedRaw } from 'vue-router'

const selectors = {
  userAvatarImage: 'avatar-image-stub.files-collaborators-collaborator-indicator',
  notUserAvatar: 'oc-avatar-item-stub.files-collaborators-collaborator-indicator',
  collaboratorName: '.files-collaborators-collaborator-name',
  collaboratorRole: '.files-collaborators-collaborator-role',
  shareInheritanceIndicators: '.files-collaborators-collaborator-shared-via',
  expirationDateIcon: '[data-testid="recipient-info-expiration-date"]',
  externalContextHelper: '.files-collaborators-collaborator-name-wrapper .oc-contextual-helper'
}

const getShareMock = ({
  sharedWith,
  shareType,
  expirationDateTime
}: Partial<CollaboratorShare> = {}): CollaboratorShare => ({
  id: '1',
  sharedWith: sharedWith || {
    id: 'ZWluc3RlaW5AaHR0cHM6Ly93d3cubG9yZW0uY29t@www.lorem.com',
    displayName: 'einstein'
  },
  sharedBy: { id: '2', displayName: 'marie' },
  permissions: [],
  shareType: shareType || ShareTypes.user.value,
  role: mock<ShareRole>({ description: '', displayName: '' }),
  resourceId: '1',
  indirect: false,
  expirationDateTime: expirationDateTime || '',
  createdDateTime: '2024-01-01'
})

describe('Collaborator ListItem component', () => {
  describe('displays the correct image/icon according to the shareType', () => {
    describe('user share type', () => {
      it('should display a users avatar', () => {
        const { wrapper } = createWrapper({
          share: getShareMock({ shareType: ShareTypes.user.value })
        })
        expect(wrapper.find(selectors.userAvatarImage).exists()).toBeTruthy()
        expect(wrapper.find(selectors.notUserAvatar).exists()).toBeFalsy()
      })
      it('sets user info on the avatar', () => {
        const share = getShareMock()
        const { wrapper } = createWrapper({ share })
        expect(wrapper.find(selectors.userAvatarImage).attributes('userid')).toEqual(
          share.sharedWith.id
        )
        expect(wrapper.find(selectors.userAvatarImage).attributes('user-name')).toEqual(
          share.sharedWith.displayName
        )
      })
    })
    describe('non-user share types', () => {
      it.each(ShareTypes.all.filter((shareType) => shareType !== ShareTypes.user))(
        'should display an oc-avatar-item for any non-user share types',
        (shareType) => {
          const { wrapper } = createWrapper({ share: getShareMock({ shareType: shareType.value }) })
          expect(wrapper.find(selectors.userAvatarImage).exists()).toBeFalsy()
          expect(wrapper.find(selectors.notUserAvatar).exists()).toBeTruthy()
          expect(wrapper.find(selectors.notUserAvatar).attributes().name).toEqual(shareType.key)
        }
      )
      it('should display an oc-avatar-item for space group shares', () => {
        const { wrapper } = createWrapper({
          share: getShareMock({
            shareType: ShareTypes.group.value,
            sharedWith: { id: '1', displayName: 'someGroup' }
          })
        })
        expect(wrapper.find(selectors.userAvatarImage).exists()).toBeFalsy()
        expect(wrapper.find(selectors.notUserAvatar).exists()).toBeTruthy()
      })
    })
  })
  describe('share info', () => {
    it('shows the collaborator display name', () => {
      const share = getShareMock()
      const { wrapper } = createWrapper({ share })
      expect(wrapper.find(selectors.collaboratorName).text()).toEqual(share.sharedWith.displayName)
    })
    it('shows the share expiration date if given', () => {
      const { wrapper } = createWrapper({
        share: getShareMock({ expirationDateTime: '2000-01-01' })
      })
      expect(wrapper.find(selectors.expirationDateIcon).exists()).toBeTruthy()
    })
  })
  describe('modifiable property', () => {
    it('shows interactive elements when modifiable', () => {
      const { wrapper } = createWrapper({ modifiable: true })
      expect(wrapper.find(selectors.collaboratorRole).exists()).toBeTruthy()
    })
    it('hides interactive elements when not modifiable', () => {
      const { wrapper } = createWrapper({ modifiable: false })
      expect(wrapper.find(selectors.collaboratorRole).exists()).toBeFalsy()
    })
  })
  describe('share inheritance indicators', () => {
    it('show when sharedParentRoute is given', () => {
      const { wrapper } = createWrapper({
        sharedParentRoute: { params: { driveAliasAndItem: '/folder' } }
      })
      expect(wrapper.find(selectors.shareInheritanceIndicators).exists()).toBeTruthy()
      expect(wrapper.html()).toMatchSnapshot()
    })
    it('do not show when sharedParentRoute is not given', () => {
      const { wrapper } = createWrapper()
      expect(wrapper.find(selectors.shareInheritanceIndicators).exists()).toBeFalsy()
    })
  })
  describe('remove share', () => {
    it('emits the "removeShare" event', () => {
      const { wrapper } = createWrapper()
      wrapper.findComponent<typeof EditDropdown>('edit-dropdown-stub').vm.$emit('removeShare')
      expect(wrapper.emitted().onDelete).toBeTruthy()
    })
  })
  describe('change share role', () => {
    it('calls "changeShare" for regular resources', () => {
      const { wrapper } = createWrapper()
      wrapper.findComponent<typeof RoleDropdown>('role-dropdown-stub').vm.$emit('optionChange', {
        permissions: [GraphSharePermission.readBasic]
      })
      const sharesStore = useSharesStore()
      expect(sharesStore.updateShare).toHaveBeenCalled()
    })
    it('shows a message on error', async () => {
      const resource = mock<SpaceResource>({ driveType: 'project' })
      vi.spyOn(console, 'error').mockImplementation(() => undefined)
      const { wrapper } = createWrapper()
      const sharesStore = useSharesStore()
      ;(sharesStore.updateShare as Mock).mockRejectedValueOnce(new Error())
      wrapper.findComponent<typeof RoleDropdown>('role-dropdown-stub').vm.$emit('optionChange', {
        share: getShareMock({ shareType: ShareTypes.user.value }),
        resource
      })

      await flushPromises()

      const messagesStore = useMessages()
      expect(messagesStore.showErrorMessage).toHaveBeenCalled()
    })
  })
  describe('change expiration date', () => {
    it('calls "changeShare" for regular resources', () => {
      const { wrapper } = createWrapper()
      wrapper
        .findComponent<typeof EditDropdown>('edit-dropdown-stub')
        .vm.$emit('expirationDateChanged', {
          shareExpirationChanged: new Date()
        })
      const sharesStore = useSharesStore()
      expect(sharesStore.updateShare).toHaveBeenCalled()
    })
    it('shows a message on error', async () => {
      vi.spyOn(console, 'error').mockImplementation(() => undefined)
      const { wrapper } = createWrapper()
      const sharesStore = useSharesStore()
      ;(sharesStore.updateShare as Mock).mockRejectedValueOnce(new Error())
      wrapper
        .findComponent<typeof EditDropdown>('edit-dropdown-stub')
        .vm.$emit('expirationDateChanged', {
          shareExpirationChanged: new Date()
        })

      await flushPromises()
      const messagesStore = useMessages()
      expect(messagesStore.showErrorMessage).toHaveBeenCalled()
    })
  })
  describe('external user shares', () => {
    it('correctly identifies external shares', () => {
      const share = getShareMock({ shareType: ShareTypes.remote.value })
      const { wrapper } = createWrapper({ share })
      const roleDropDown = wrapper.findComponent<typeof RoleDropdown>('role-dropdown-stub')

      expect(roleDropDown.props('isExternal')).toBeTruthy()
      expect(wrapper.find(selectors.externalContextHelper).exists()).toBeTruthy()
    })

    it('should show domain name below external user name', () => {
      const share = getShareMock({ shareType: ShareTypes.remote.value })
      const { wrapper } = createWrapper({ share })
      const el = wrapper.find('[data-testid="external-share-domain"]')

      expect(el.text()).toEqual('www.lorem.com')
    })
  })
})

function createWrapper({
  share = getShareMock(),
  modifiable = true,
  sharedParentRoute = null,
  resource = mock<Resource>()
}: {
  share?: CollaboratorShare
  modifiable?: boolean
  sharedParentRoute?: RouteLocationNamedRaw
  resource?: Resource
} = {}) {
  const mocks = defaultComponentMocks()
  mocks.$clientService.graphAuthenticated.drives.getDrive.mockResolvedValue(undefined)

  return {
    wrapper: mount(ListItem, {
      props: {
        share,
        modifiable,
        sharedParentRoute
      },
      global: {
        plugins: [...defaultPlugins()],
        mocks,
        provide: { ...mocks, resource },
        renderStubDefaultSlot: true,
        stubs: {
          ...defaultStubs,
          'oc-icon': true,
          'avatar-image': true,
          'router-link': true,
          'oc-info-drop': true,
          'oc-table-simple': true,
          'oc-tr': true,
          'oc-td': true,
          'oc-tag': true,
          'oc-pagination': true,
          'oc-avatar-item': true,
          'role-dropdown': true,
          'edit-dropdown': true,
          translate: false
        }
      }
    })
  }
}
