import MembersRoleSection from '../../../../../src/components/Spaces/SideBar/MembersRoleSection.vue'
import { defaultPlugins, shallowMount } from '@ownclouders/web-test-helpers'
import { mock } from 'vitest-mock-extended'
import { SpaceMember } from '@ownclouders/web-client'

describe('MembersRoleSection', () => {
  it('should render all members accordingly', () => {
    const members = [
      mock<SpaceMember>({ grantedTo: { user: { displayName: 'einstein' }, group: undefined } }),
      mock<SpaceMember>({ grantedTo: { group: { displayName: 'physic-lovers' }, user: undefined } })
    ]
    const { wrapper } = getWrapper({ members })
    expect(wrapper.html()).toMatchSnapshot()
  })
})

function getWrapper({ members = [] }: { members?: SpaceMember[] } = {}) {
  return {
    wrapper: shallowMount(MembersRoleSection, {
      props: {
        members
      },
      global: {
        plugins: [...defaultPlugins()]
      }
    })
  }
}
