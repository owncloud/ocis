import { defaultPlugins, shallowMount } from '@ownclouders/web-test-helpers'
import ResourceStatusIndicators from '../../../../src/components/FilesList/ResourceStatusIndicators.vue'
import { Resource, SpaceResource } from '@ownclouders/web-client'
import { mock } from 'vitest-mock-extended'
import { getIndicators, ResourceIndicator } from '../../../../src/helpers/statusIndicators'
import { OcStatusIndicators } from '@ownclouders/design-system/components'

vi.mock('../../../../src/helpers/statusIndicators', async (importOriginal) => ({
  ...(await importOriginal<any>()),
  getIndicators: vi.fn(() => [])
}))

describe('ResourceStatusIndicators component', () => {
  it('renders indicators from getIndicators', () => {
    const space = mock<SpaceResource>({ driveType: 'project' })
    const resource = mock<Resource>({ id: 'resource' })

    const indicators = [
      {
        id: 'some-id',
        type: 'some-type',
        category: 'system',
        label: 'label',
        accessibleDescription: 'accessible description',
        icon: 'download',
        fillType: 'fill'
      }
    ] satisfies ResourceIndicator[]

    vi.mocked(getIndicators).mockReturnValue(indicators)

    const wrapper = getWrapper({ space, resource })
    expect(wrapper.findComponent(OcStatusIndicators).props('indicators')).toEqual(indicators)
  })
  function getWrapper(props: {
    space: SpaceResource
    resource: Resource
    filter?: (resource: ResourceIndicator) => boolean
  }) {
    return shallowMount(ResourceStatusIndicators, {
      props,
      global: {
        plugins: [...defaultPlugins()]
      }
    })
  }
})
