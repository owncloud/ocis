import MembersPanel from '../../../../../src/components/Groups/SideBar/MembersPanel.vue'
import { defaultPlugins, shallowMount } from '@ownclouders/web-test-helpers'
import { mock } from 'vitest-mock-extended'
import { Group } from '@ownclouders/web-client/graph/generated'
import MembersRoleSection from '../../../../../src/components/Groups/SideBar/MembersRoleSection.vue'

const groupMock = mock<Group>({
  id: '1',
  groupTypes: [],
  members: [{ displayName: 'Albert Einstein' }]
})
const selectors = {
  membersRolePanelStub: 'members-role-section-stub'
}

describe('MembersPanel', () => {
  it('should render all members accordingly to their role assignments', () => {
    const { wrapper } = getWrapper()
    expect(wrapper.html()).toMatchSnapshot()
  })
  it('should filter members accordingly to the entered search term', async () => {
    const { wrapper } = getWrapper()
    ;(wrapper.vm as any).filterTerm = 'ein'
    await wrapper.vm.$nextTick
    expect(wrapper.findAll(selectors.membersRolePanelStub).length).toBe(1)
    expect(
      wrapper.findComponent<typeof MembersRoleSection>(selectors.membersRolePanelStub).props()
        .groupMembers[0].displayName
    ).toEqual('Albert Einstein')
  })
  it('should display an empty result if no matching members found', async () => {
    const { wrapper } = getWrapper()
    ;(wrapper.vm as any).filterTerm = 'no-match'
    await wrapper.vm.$nextTick
    expect(wrapper.findAll(selectors.membersRolePanelStub).length).toBe(0)
    expect(wrapper.html()).toMatchSnapshot()
  })
})

function getWrapper({ group = groupMock } = {}) {
  return {
    wrapper: shallowMount(MembersPanel, {
      global: {
        plugins: [...defaultPlugins()],
        provide: { group: group }
      }
    })
  }
}
