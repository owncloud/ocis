import { SpaceResource } from '@ownclouders/web-client'
import SpaceInfo from '../../../../../src/components/SideBar/Spaces/SpaceInfo.vue'
import { defaultPlugins, shallowMount } from '@ownclouders/web-test-helpers'

const spaceMock = {
  type: 'space',
  name: ' space',
  id: '1',
  mdate: 'Wed, 21 Oct 2015 07:28:00 GMT',
  spaceQuota: {
    used: 100
  }
} as unknown as SpaceResource

const selectors = {
  name: '[data-testid="space-info-name"]',
  subtitle: '[data-testid="space-info-subtitle"]'
}

describe('SpaceInfo', () => {
  it('shows space info', () => {
    const { wrapper } = createWrapper(spaceMock)
    expect(wrapper.find(selectors.name).exists()).toBeTruthy()
    expect(wrapper.find(selectors.subtitle).exists()).toBeTruthy()
    expect(wrapper.html()).toMatchSnapshot()
  })
})

function createWrapper(spaceResource: SpaceResource) {
  return {
    wrapper: shallowMount(SpaceInfo, {
      global: {
        plugins: [...defaultPlugins()],
        provide: { resource: spaceResource }
      }
    })
  }
}
