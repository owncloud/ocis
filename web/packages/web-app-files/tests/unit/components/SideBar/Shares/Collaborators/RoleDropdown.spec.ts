import RoleDropdown from '../../../../../../src/components/SideBar/Shares/Collaborators/RoleDropdown.vue'
import { ShareRole } from '@ownclouders/web-client'
import { defaultPlugins, mount, shallowMount } from '@ownclouders/web-test-helpers'
import { mock } from 'vitest-mock-extended'
import { Resource } from '@ownclouders/web-client'
import { User } from '@ownclouders/web-client/graph/generated'

const selectors = {
  recipientRoleBtn: '.files-recipient-role-select-btn',
  roleButton: '.files-recipient-role-drop-btn',
  filesPermissionActionsList: '.files-permission-actions-list'
}

describe('RoleDropdown', () => {
  it('renders a button with invite text if no existing role given', () => {
    const { wrapper } = getWrapper({ mountType: shallowMount })
    expect(wrapper.find(selectors.recipientRoleBtn).exists()).toBeTruthy()
    expect(wrapper.find(`${selectors.recipientRoleBtn} span`).text()).toEqual('Can view')
  })
  it('renders a button with existing role if given', () => {
    const { wrapper } = getWrapper({
      mountType: shallowMount,
      existingShareRole: mock<ShareRole>({ displayName: 'Can edit' })
    })
    expect(wrapper.find(selectors.recipientRoleBtn).exists()).toBeTruthy()
    expect(wrapper.find(`${selectors.recipientRoleBtn} span`).text()).toEqual('Can edit')
    expect(wrapper.find(selectors.filesPermissionActionsList).exists()).toBeFalsy()
  })
  it('lists permission actions if a role is unknown', () => {
    const { wrapper } = getWrapper({
      mountType: shallowMount,
      existingSharePermissions: ['read', 'update']
    })
    expect(wrapper.find(selectors.recipientRoleBtn).exists()).toBeTruthy()
    expect(wrapper.find(selectors.filesPermissionActionsList).exists()).toBeTruthy()
  })
  it('does not render a button if only one role is available', () => {
    const { wrapper } = getWrapper({
      mountType: shallowMount,
      availableInternalShareRoles: [mock<ShareRole>({ displayName: 'Can view', description: '' })]
    })
    expect(wrapper.find(selectors.recipientRoleBtn).exists()).toBeFalsy()
  })
  it('emits "optionChange"-event on role click', async () => {
    const { wrapper } = getWrapper()
    ;(wrapper.vm.$refs.rolesDrop as any).tippy = { hide: vi.fn() }
    await wrapper.find(selectors.roleButton).trigger('click')
    expect(wrapper.emitted('optionChange')).toBeTruthy()
  })
  it('renders a button for each available role', () => {
    const { wrapper } = getWrapper({ mountType: shallowMount })
    expect(wrapper.findAll(selectors.roleButton).length).toBe(2)
  })
  it('uses available external share roles if "isExternal" is given', () => {
    const externalShareRole2 = mock<ShareRole>({ id: 'external1', displayName: '' })
    const externalShareRole1 = mock<ShareRole>({ id: 'external2', displayName: '' })
    const { wrapper } = getWrapper({
      mountType: shallowMount,
      isExternal: true,
      availableExternalShareRoles: [externalShareRole1, externalShareRole2]
    })

    expect(
      wrapper.find(`oc-button-stub#files-recipient-role-drop-btn-${externalShareRole1.id}`).exists()
    ).toBeTruthy()
  })
})

function getWrapper({
  mountType = mount,
  existingShareRole = null,
  existingSharePermissions = null,
  isExternal = false,
  availableInternalShareRoles = [
    mock<ShareRole>({ displayName: 'Can view', description: '' }),
    mock<ShareRole>({ displayName: 'Can edit', description: '' })
  ],
  availableExternalShareRoles = []
}: {
  mountType?: typeof mount
  existingShareRole?: ShareRole
  existingSharePermissions?: string[]
  isExternal?: boolean
  availableInternalShareRoles?: ShareRole[]
  availableExternalShareRoles?: ShareRole[]
} = {}) {
  return {
    wrapper: mountType(RoleDropdown, {
      props: {
        existingShareRole,
        existingSharePermissions: existingSharePermissions ?? [],
        isExternal
      },
      global: {
        plugins: [
          ...defaultPlugins({
            piniaOptions: { userState: { user: { onPremisesSamAccountName: 'name' } as User } }
          })
        ],
        renderStubDefaultSlot: true,
        provide: {
          resource: mock<Resource>(),
          availableInternalShareRoles,
          availableExternalShareRoles
        }
      }
    })
  }
}
