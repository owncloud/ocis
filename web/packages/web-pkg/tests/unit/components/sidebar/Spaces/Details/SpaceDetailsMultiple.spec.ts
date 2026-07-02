import { SpaceResource } from '@ownclouders/web-client'
import SpaceDetailsMultiple from '../../../../../../src/components/SideBar/Spaces/Details/SpaceDetailsMultiple.vue'
import { defaultPlugins, shallowMount } from '@ownclouders/web-test-helpers'

const spaceMock = {
  type: 'space',
  name: ' space',
  id: '1',
  mdate: 'Wed, 21 Oct 2015 07:28:00 GMT',
  spaceQuota: {
    used: 100,
    total: 1000,
    remaining: 900
  }
} as unknown as SpaceResource

describe('Multiple Details SideBar Panel', () => {
  it('displays the details side panel', () => {
    const { wrapper } = createWrapper(spaceMock)
    expect(wrapper.html()).toMatchSnapshot()
  })
})

function createWrapper(spaceResource: SpaceResource) {
  return {
    wrapper: shallowMount(SpaceDetailsMultiple, {
      global: {
        plugins: [...defaultPlugins()]
      },
      props: {
        selectedSpaces: [spaceResource]
      }
    })
  }
}
