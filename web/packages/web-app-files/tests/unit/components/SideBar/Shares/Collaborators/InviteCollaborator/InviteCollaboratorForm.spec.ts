import { mock } from 'vitest-mock-extended'
import InviteCollaboratorForm from '../../../../../../../src/components/SideBar/Shares/Collaborators/InviteCollaborator/InviteCollaboratorForm.vue'
import {
  defaultComponentMocks,
  defaultPlugins,
  RouteLocation,
  shallowMount,
  VueWrapper
} from '@ownclouders/web-test-helpers'
import { Resource, SpaceResource } from '@ownclouders/web-client'
import { useSharesStore } from '@ownclouders/web-pkg'
import { CollaboratorAutoCompleteItem, CollaboratorShare, ShareRole } from '@ownclouders/web-client'
import { Group, User } from '@ownclouders/web-client/graph/generated'
import { OcButton } from '@ownclouders/design-system/components'
import RoleDropdown from '../../../../../../../src/components/SideBar/Shares/Collaborators/RoleDropdown.vue'
import { ShareRoleType } from '../../../../../../../src/components/SideBar/Shares/Collaborators/InviteCollaborator/InviteCollaboratorForm.vue'

vi.mock('lodash-es', () => ({ debounce: (fn: any) => fn() }))

const folderMock = {
  id: '1',
  type: 'folder',
  isFolder: true,
  mdate: 'Wed, 21 Oct 2015 07:28:00 GMT',
  size: '740',
  isMounted: vi.fn(() => true),
  name: 'lorem.txt',
  privateLink: 'some-link',
  canShare: vi.fn(() => true),
  path: '/documents',
  canDeny: () => false,
  spaceId: '1'
} as Resource

const spaceMock = {
  id: '1',
  type: 'space'
}

describe('InviteCollaboratorForm', () => {
  describe('renders correctly', () => {
    it.todo('renders a select field for share receivers')
    it.todo('renders an inviteDescriptionMessage')
    it.todo('renders a role selector component')
    it.todo('renders an expiration datepicker component')
    it.todo('renders an hidden-announcer')
  })
  describe('share button', () => {
    it('is present', () => {
      const { wrapper } = getWrapper()
      expect(wrapper.find('#new-collaborators-form-create-button').exists()).toBeTruthy()
    })
    it('is disabled when no collaborators selected', () => {
      const { wrapper } = getWrapper()
      const btn = wrapper.findComponent<typeof OcButton>('#new-collaborators-form-create-button')
      expect(btn.props('disabled')).toBeTruthy()
    })
    it('is disabled when currently saving', async () => {
      const { wrapper } = getWrapper()
      wrapper.vm.selectedCollaborators = [mock<CollaboratorAutoCompleteItem>()]
      wrapper.vm.saving = true
      await wrapper.vm.$nextTick()
      const btn = wrapper.findComponent<typeof OcButton>('#new-collaborators-form-create-button')
      expect(btn.props('disabled')).toBeTruthy()
    })
    it('is enabled when collaborators selected and not saving', async () => {
      const { wrapper } = getWrapper()
      wrapper.vm.selectedCollaborators = [mock<CollaboratorAutoCompleteItem>()]
      wrapper.vm.saving = false
      await wrapper.vm.$nextTick()
      const btn = wrapper.findComponent<typeof OcButton>('#new-collaborators-form-create-button')
      expect(btn.props('disabled')).toBeFalsy()
    })
  })
  describe('fetching recipients', () => {
    it('fetches recipients upon mount', async () => {
      const { wrapper, mocks } = getWrapper()
      await wrapper.vm.fetchRecipientsTask.last

      expect(mocks.$clientService.graphAuthenticated.users.listUsers).toHaveBeenCalledTimes(1)
      expect(mocks.$clientService.graphAuthenticated.groups.listGroups).toHaveBeenCalledTimes(1)
    })
    it('fetches users and groups returned from the server', async () => {
      const { wrapper } = getWrapper({ users: [{ id: '2' } as User], groups: [{ id: '3' }] })
      await wrapper.vm.fetchRecipientsTask.last

      expect(wrapper.vm.autocompleteResults.length).toBe(2)
    })
    it('filters out the current user', async () => {
      const { wrapper } = getWrapper({ users: [{ id: '1' } as User], groups: [{ id: '3' }] })
      await wrapper.vm.fetchRecipientsTask.last

      expect(wrapper.vm.autocompleteResults.length).toBe(1)
    })
    it('filters out selected users', async () => {
      const { wrapper } = getWrapper({ users: [{ id: '2' } as User], groups: [{ id: '3' }] })
      wrapper.vm.selectedCollaborators = [mock<CollaboratorAutoCompleteItem>({ id: '2' })]
      await wrapper.vm.fetchRecipientsTask.last

      expect(wrapper.vm.autocompleteResults.length).toBe(1)
    })
    it('filters out existing direct shares', async () => {
      const { wrapper } = getWrapper({
        users: [{ id: '2' } as User],
        groups: [{ id: '3' }],
        existingCollaborators: [
          mock<CollaboratorShare>({ sharedWith: { id: '2' }, indirect: false })
        ]
      })

      await wrapper.vm.fetchRecipientsTask.last

      expect(wrapper.vm.autocompleteResults.length).toBe(1)
    })
  })
  describe('share action', () => {
    it('clicking the invite-sharees button calls the "share"-action', async () => {
      const { wrapper } = getWrapper()
      const shareSpy = vi.spyOn(wrapper.vm, 'share')
      const shareBtn = wrapper.find('#new-collaborators-form-create-button')
      expect(shareBtn.exists()).toBeTruthy()

      await wrapper.vm.$nextTick()
      await shareBtn.trigger('click')
      expect(shareSpy).toHaveBeenCalledTimes(0)
    })
    it.each([
      { storageId: undefined, resource: mock<Resource>(folderMock) },
      { storageId: undefined, resource: mock<SpaceResource>(spaceMock) },
      { storageId: '1', resource: mock<Resource>(folderMock) }
    ] as const)('calls the "addShare" action', async (dataSet) => {
      const { wrapper } = getWrapper({
        storageId: dataSet.storageId,
        resource: dataSet.resource
      })

      wrapper.vm.selectedCollaborators = [mock<CollaboratorAutoCompleteItem>()]

      const { addShare } = useSharesStore()
      vi.mocked(addShare).mockResolvedValue(undefined)

      await wrapper.vm.$nextTick()
      await wrapper.vm.share()

      expect(addShare).toHaveBeenCalled()
    })
    it.todo('resets focus upon selecting an invitee')
  })
  describe('share role type filter', () => {
    it.each([
      { externalRoles: [], available: false },
      { externalRoles: [mock<ShareRole>()], available: true }
    ])(
      'is present depending on the available external share roles',
      ({ externalRoles, available }) => {
        const { wrapper } = getWrapper({ externalShareRoles: externalRoles })
        expect(wrapper.find('.invite-form-share-role-type').exists()).toBe(available)
      }
    )
    it('correctly passes the external prop to the role dropdown component', async () => {
      const externalRoles = [mock<ShareRole>()]
      const { wrapper } = getWrapper({ externalShareRoles: externalRoles })
      wrapper.vm.currentShareRoleType = mock<ShareRoleType>({ id: '2' })
      await wrapper.vm.$nextTick()

      const roleDropdown = wrapper.findComponent<typeof RoleDropdown>('role-dropdown-stub')
      expect(roleDropdown.props('isExternal')).toBeTruthy()
    })
    it('is not present for project space resources', () => {
      const space = mock<SpaceResource>({ driveType: 'project' })
      const externalRoles = [mock<ShareRole>()]
      const { wrapper } = getWrapper({ externalShareRoles: externalRoles, resource: space })

      expect(wrapper.find('.invite-form-share-role-type').exists()).toBeFalsy()
    })
  })
})

function getWrapper({
  storageId = 'fake-storage-id',
  resource = mock<Resource>(folderMock),
  users = [],
  groups = [],
  existingCollaborators = [],
  externalShareRoles = [],
  user = mock<User>({ id: '1' })
}: {
  storageId?: string
  resource?: Resource
  users?: User[]
  groups?: Group[]
  existingCollaborators?: CollaboratorShare[]
  externalShareRoles?: ShareRole[]
  user?: User
} = {}) {
  const mocks = defaultComponentMocks({
    currentRoute: mock<RouteLocation>({ params: { storageId } })
  })

  mocks.$clientService.graphAuthenticated.users.listUsers.mockResolvedValue(users)
  mocks.$clientService.graphAuthenticated.groups.listGroups.mockResolvedValue(groups)

  const capabilities = { files_sharing: { federation: { incoming: true, outgoing: true } } }

  return {
    mocks,
    wrapper: shallowMount(InviteCollaboratorForm, {
      global: {
        plugins: [
          ...defaultPlugins({
            piniaOptions: {
              userState: { user },
              capabilityState: { capabilities },
              configState: { options: { concurrentRequests: { shares: { create: 1 } } } },
              sharesState: {
                collaboratorShares: existingCollaborators
              }
            }
          })
        ],
        provide: {
          ...mocks,
          resource,
          availableExternalShareRoles: externalShareRoles,
          availableInternalShareRoles: [mock<ShareRole>()]
        },
        mocks,
        stubs: { OcSelect: false, VueSelect: false }
      }
    }) as VueWrapper<any, any>
  }
}
