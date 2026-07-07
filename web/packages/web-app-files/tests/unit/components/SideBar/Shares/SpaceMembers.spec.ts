import SpaceMembers from '../../../../../src/components/SideBar/Shares/SpaceMembers.vue'
import {
  ShareTypes,
  ShareRole,
  CollaboratorShare,
  GraphSharePermission
} from '@ownclouders/web-client'
import { mock } from 'vitest-mock-extended'
import { ProjectSpaceResource, SpaceResource } from '@ownclouders/web-client'
import {
  defaultPlugins,
  mount,
  shallowMount,
  defaultComponentMocks,
  RouteLocation,
  VueWrapper
} from '@ownclouders/web-test-helpers'
import { User } from '@ownclouders/web-client/graph/generated'
import { useCanShare, useModals, useSpacesStore } from '@ownclouders/web-pkg'
import ListItem from '../../../../../src/components/SideBar/Shares/Collaborators/ListItem.vue'

vi.mock('@ownclouders/web-pkg', async (importOriginal) => ({
  ...(await importOriginal<any>()),
  useCanShare: vi.fn()
}))

const memberMocks = [
  {
    id: '1',
    shareType: ShareTypes.user.value,
    sharedWith: {
      id: 'alice',
      displayName: 'alice'
    },
    role: mock<ShareRole>(),
    permissions: [GraphSharePermission.updatePermissions],
    resourceId: '1',
    indirect: false,
    sharedBy: { id: 'admin', displayName: 'admin' }
  },
  {
    id: '2',
    shareType: ShareTypes.user.value,
    sharedWith: {
      onPremisesSamAccountName: 'Einstein',
      displayName: 'einstein'
    },
    role: mock<ShareRole>(),
    permissions: [],
    resourceId: '1',
    indirect: false,
    sharedBy: { id: 'admin', displayName: 'admin' }
  },
  {
    id: '3',
    shareType: ShareTypes.user.value,
    sharedWith: {
      onPremisesSamAccountName: 'Marie',
      displayName: 'marie'
    },
    role: mock<ShareRole>(),
    permissions: [],
    resourceId: '1',
    indirect: false,
    sharedBy: { id: 'admin', displayName: 'admin' }
  }
] as CollaboratorShare[]

describe('SpaceMembers', () => {
  describe('invite collaborator form', () => {
    it('renders the form when the current user can share', () => {
      const wrapper = getWrapper({ canShare: true })
      expect(wrapper.find('invite-collaborator-form-stub').exists()).toBeTruthy()
    })
    it('does not render the form when the current user can not share', () => {
      const wrapper = getWrapper({ canShare: false })
      expect(wrapper.find('invite-collaborator-form-stub').exists()).toBeFalsy()
    })
  })

  describe('existing members', () => {
    it('can edit when current user can share', () => {
      const wrapper = getWrapper({ canShare: true })
      expect(
        wrapper.findAllComponents<any>('collaborator-list-item-stub').at(1).props().modifiable
      ).toEqual(true)
    })
    it('can not edit when current user can not share', () => {
      const wrapper = getWrapper({ canShare: false })
      expect(
        wrapper.findAllComponents<any>('collaborator-list-item-stub').at(1).props().modifiable
      ).toEqual(false)
    })
    it('can not edit current user when they are the only space manager', () => {
      const wrapper = getWrapper({ spaceMembers: [memberMocks[0]], canShare: true })
      expect(
        wrapper.findAllComponents<any>('collaborator-list-item-stub').at(0).props().modifiable
      ).toEqual(false)
    })
  })

  describe('deleting members', () => {
    it('reacts on delete events by collaborator list items', async () => {
      const user = mock<User>({ id: 'admin' })
      const wrapper = getWrapper({ user })
      wrapper.findComponent<typeof ListItem>('collaborator-list-item-stub').vm.$emit('onDelete')
      await wrapper.vm.$nextTick()

      const { dispatchModal } = useModals()
      expect(dispatchModal).toHaveBeenCalledTimes(1)
    })
  })

  describe('filter', () => {
    it('toggles the filter on click', async () => {
      const space = mock<ProjectSpaceResource>()
      const wrapper = getWrapper({ mountType: mount, space })
      expect(wrapper.vm.isFilterOpen).toBeFalsy()
      await wrapper.find('.open-filter-btn').trigger('click')
      expect(wrapper.vm.isFilterOpen).toBeTruthy()
    })
  })
})

function getWrapper({
  mountType = shallowMount,
  space = mock<SpaceResource>(),
  spaceMembers = memberMocks,
  user = mock<User>(),
  currentRouteName = 'files-spaces-generic',
  canShare = true
} = {}) {
  vi.mocked(useCanShare).mockReturnValue({ canShare: () => canShare })

  const mocks = defaultComponentMocks({
    currentRoute: mock<RouteLocation>({ name: currentRouteName })
  })

  const plugins = defaultPlugins({
    piniaOptions: {
      userState: { user },
      configState: {
        options: { contextHelpers: true }
      }
    }
  })

  const spacesStore = useSpacesStore()
  vi.mocked(spacesStore).getSpaceMembers.mockReturnValue(spaceMembers)

  return mountType(SpaceMembers, {
    global: {
      plugins,
      mocks,
      provide: {
        ...mocks,
        space,
        resource: space
      },
      stubs: {
        OcButton: false,
        'oc-icon': true,
        'oc-spinner': true,
        'avatar-image': true,
        'role-dropdown': true,
        'edit-dropdown': true,
        'invite-collaborator-form': true,
        'collaborator-list-item': true
      }
    }
  }) as VueWrapper<any, any>
}
