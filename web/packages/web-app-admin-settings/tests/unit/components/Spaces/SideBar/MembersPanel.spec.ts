import MembersPanel from '../../../../../src/components/Spaces/SideBar/MembersPanel.vue'
import { defaultPlugins, shallowMount } from '@ownclouders/web-test-helpers'
import { mock } from 'vitest-mock-extended'
import { ShareRole, SpaceResource } from '@ownclouders/web-client'
import MembersRoleSection from '../../../../../src/components/Spaces/SideBar/MembersRoleSection.vue'

const graphRoles = {
  '1': mock<ShareRole>({ id: '1', displayName: 'Managers', rolePermissions: [] }),
  '2': mock<ShareRole>({ id: '2', displayName: 'Editors', rolePermissions: [] }),
  '3': mock<ShareRole>({ id: '3', displayName: 'Viewers', rolePermissions: [] })
}

const spaceMock = {
  members: {
    '1': { roleId: '1', grantedTo: { user: { displayName: 'admin' } } },
    '2': { roleId: '2', grantedTo: { user: { displayName: 'marie' } } },
    '3': { roleId: '3', grantedTo: { user: { displayName: 'einstein' } } }
  }
} as undefined as SpaceResource

const selectors = {
  membersRolePanelStub: 'members-role-section-stub',
  spaceMembersCustom: '.space-members-custom'
}

describe('MembersPanel', () => {
  it('should render all members accordingly to their role assignments', () => {
    const { wrapper } = getWrapper()
    expect(wrapper.html()).toMatchSnapshot()
  })
  it('should filter members accordingly to the entered search term', async () => {
    const userToFilterFor = spaceMock.members['3']
    const { wrapper } = getWrapper()
    ;(wrapper.vm as any).filterTerm = 'ein'
    await wrapper.vm.$nextTick()
    expect(wrapper.findAll(selectors.membersRolePanelStub).length).toBe(1)
    expect(
      wrapper.findComponent<typeof MembersRoleSection>(selectors.membersRolePanelStub).props()
        .members[0].grantedTo.user.displayName
    ).toEqual(userToFilterFor.grantedTo.user.displayName)
  })
  it('should display an empty result if no matching members found', async () => {
    const { wrapper } = getWrapper()
    ;(wrapper.vm as any).filterTerm = 'no-match'
    await wrapper.vm.$nextTick()
    expect(wrapper.findAll(selectors.membersRolePanelStub).length).toBe(0)
    expect(wrapper.html()).toMatchSnapshot()
  })
  it('should display members without role under the custom section', () => {
    const spaceResource = {
      members: {
        '1': { grantedTo: { user: { displayName: 'admin' } } }
      }
    } as undefined as SpaceResource
    const { wrapper } = getWrapper({ spaceResource })
    expect(wrapper.find(selectors.spaceMembersCustom).exists()).toBeTruthy()
  })
})

function getWrapper({ spaceResource = spaceMock } = {}) {
  return {
    wrapper: shallowMount(MembersPanel, {
      global: {
        plugins: [...defaultPlugins({ piniaOptions: { sharesState: { graphRoles } } })],
        provide: { resource: spaceResource }
      }
    })
  }
}
