import UsersList from '../../../../src/components/Users/UsersList.vue'
import {
  defaultComponentMocks,
  defaultPlugins,
  mount,
  shallowMount
} from '@ownclouders/web-test-helpers'
import { displayPositionedDropdown, eventBus, queryItemAsString } from '@ownclouders/web-pkg'
import { SideBarEventTopics } from '@ownclouders/web-pkg'
import { useUserSettingsStore } from '../../../../src/composables/stores/userSettings'
import { User } from '@ownclouders/web-client/graph/generated'

const getUserMocks = () => [{ id: '1', displayName: 'jan' }] as User[]
vi.mock('@ownclouders/web-pkg', async (importOriginal) => ({
  ...(await importOriginal<any>()),
  displayPositionedDropdown: vi.fn(),
  queryItemAsString: vi.fn()
}))

describe('UsersList', () => {
  describe('computed method "allUsersSelected"', () => {
    it('should be true if all users are selected', () => {
      const { wrapper } = getWrapper({
        users: getUserMocks(),
        selectedUsers: getUserMocks()
      })
      expect((wrapper.vm as any).allUsersSelected).toBeTruthy()
    })
    it('should be false if not every user is selected', () => {
      const { wrapper } = getWrapper({
        users: getUserMocks(),
        selectedUsers: []
      })
      expect((wrapper.vm as any).allUsersSelected).toBeFalsy()
    })
  })

  describe('method "orderBy"', () => {
    it('should return an ascending ordered list while desc is set to false', () => {
      const { wrapper } = getWrapper()

      expect(
        (wrapper.vm as any).orderBy(
          [{ displayName: 'user' }, { displayName: 'admin' }] as User[],
          'displayName',
          false
        )
      ).toEqual([{ displayName: 'admin' }, { displayName: 'user' }])
    })
    it('should return an descending ordered list based on role while desc is set to true', () => {
      const { wrapper } = getWrapper()

      expect(
        (wrapper.vm as any).orderBy(
          [{ displayName: 'admin' }, { displayName: 'user' }] as User[],
          'displayName',
          true
        )
      ).toEqual([{ displayName: 'user' }, { displayName: 'admin' }])
    })

    it('should return ascending ordered list based on role while desc is set to false', () => {
      const { wrapper } = getWrapper()

      expect(
        (wrapper.vm as any).orderBy(
          [
            { appRoleAssignments: [{ appRoleId: '1' }] },
            { appRoleAssignments: [{ appRoleId: '2' }] }
          ] as User[],
          'role',
          false
        )
      ).toEqual([
        { appRoleAssignments: [{ appRoleId: '1' }] },
        { appRoleAssignments: [{ appRoleId: '2' }] }
      ])
    })
    it('should return an role based descending ordered list while desc is set to true', () => {
      const { wrapper } = getWrapper()

      expect(
        (wrapper.vm as any).orderBy(
          [
            { appRoleAssignments: [{ appRoleId: '1' }] },
            { appRoleAssignments: [{ appRoleId: '2' }] }
          ] as User[],
          'role',
          true
        )
      ).toEqual([
        { appRoleAssignments: [{ appRoleId: '2' }] },
        { appRoleAssignments: [{ appRoleId: '1' }] }
      ])
    })
  })
  it('should show the context menu on right click', async () => {
    const users = getUserMocks()
    const spyDisplayPositionedDropdown = vi.mocked(displayPositionedDropdown)
    const { wrapper } = getWrapper({ mountType: mount, users })
    await wrapper.find(`[data-item-id="${users[0].id}"]`).trigger('contextmenu')
    expect(spyDisplayPositionedDropdown).toHaveBeenCalledTimes(1)
  })
  it('should show the context menu on context menu button click', async () => {
    const users = getUserMocks()
    const spyDisplayPositionedDropdown = vi.mocked(displayPositionedDropdown)
    const { wrapper } = getWrapper({ mountType: mount, users })
    await wrapper.find('.users-table-btn-action-dropdown').trigger('click')
    expect(spyDisplayPositionedDropdown).toHaveBeenCalledTimes(1)
  })
  it('should show the user details on details button click', async () => {
    const users = getUserMocks()
    const eventBusSpy = vi.spyOn(eventBus, 'publish')
    const { wrapper } = getWrapper({ mountType: mount, users })
    await wrapper.find('.users-table-btn-details').trigger('click')
    expect(eventBusSpy).toHaveBeenCalledWith(SideBarEventTopics.open)
  })
  it('should show the user edit panel on edit button click', async () => {
    const users = getUserMocks()
    const eventBusSpy = vi.spyOn(eventBus, 'publish')
    const { wrapper } = getWrapper({ mountType: mount, users })
    await wrapper.find('.users-table-btn-edit').trigger('click')
    expect(eventBusSpy).toHaveBeenCalledWith(SideBarEventTopics.openWithPanel, 'EditPanel')
  })
  describe('toggle selection', () => {
    describe('selectUsers method', () => {
      it('selects all users', () => {
        const users = getUserMocks()
        const { wrapper } = getWrapper({ mountType: shallowMount, users })
        ;(wrapper.vm as any).selectUsers(users)
        const { setSelectedUsers } = useUserSettingsStore()
        expect(setSelectedUsers).toHaveBeenCalledWith(users)
      })
    })
    describe('selectUsers method', () => {
      it('selects a user', () => {
        const users = getUserMocks()
        const { wrapper } = getWrapper({ mountType: shallowMount, users: [users[0]] })
        ;(wrapper.vm as any).selectUser(users[0])
        const { addSelectedUser } = useUserSettingsStore()
        expect(addSelectedUser).toHaveBeenCalledWith(users[0])
      })
      it('de-selects a selected user', () => {
        const users = getUserMocks()
        const { wrapper } = getWrapper({
          mountType: shallowMount,
          users: [users[0]],
          selectedUsers: [users[0]]
        })
        ;(wrapper.vm as any).selectUser(users[0])
        const { setSelectedUsers } = useUserSettingsStore()
        expect(setSelectedUsers).toHaveBeenCalledWith([])
      })
    })
    describe('unselectAllUsers method', () => {
      it('de-selects all selected users', () => {
        const users = getUserMocks()
        const { wrapper } = getWrapper({
          mountType: shallowMount,
          users: [users[0]],
          selectedUsers: [users[0]]
        })
        ;(wrapper.vm as any).unselectAllUsers()
        const { setSelectedUsers } = useUserSettingsStore()
        expect(setSelectedUsers).toHaveBeenCalledWith([])
      })
    })
  })
})

function getWrapper({
  mountType = shallowMount,
  users = [],
  selectedUsers = []
}: { mountType?: typeof mount; users?: User[]; selectedUsers?: User[] } = {}) {
  vi.mocked(queryItemAsString).mockImplementationOnce(() => '1')
  vi.mocked(queryItemAsString).mockImplementationOnce(() => '100')
  const mocks = defaultComponentMocks()
  return {
    wrapper: mountType(UsersList, {
      props: {
        roles: [
          {
            displayName: 'Admin',
            id: '1'
          },
          {
            displayName: 'Guest',
            id: '2'
          },
          {
            displayName: 'Space Admin',
            id: '3'
          },
          {
            displayName: 'User',
            id: '4'
          }
        ],
        headerPosition: 0
      },
      global: {
        plugins: [
          ...defaultPlugins({
            piniaOptions: {
              userSettingsStore: { users, selectedUsers }
            }
          })
        ],
        mocks,
        provide: mocks,
        stubs: {
          OcCheckbox: true
        }
      }
    })
  }
}
