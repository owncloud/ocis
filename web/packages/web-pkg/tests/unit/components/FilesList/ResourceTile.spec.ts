import { defaultPlugins, shallowMount } from '@ownclouders/web-test-helpers'
import { RouteLocationRaw } from 'vue-router'
import { SpaceResource, Resource } from '@ownclouders/web-client'
import ResourceTile from '../../../../src/components/FilesList/ResourceTile.vue'

type TypeIconSize = 'large' | 'xlarge' | 'xxlarge' | 'xxxlarge'
interface Props {
  resource: SpaceResource | Resource
  resourceRoute?: RouteLocationRaw | null
  isResourceSelected?: boolean
  isResourceClickable?: boolean
  isResourceDisabled?: boolean
  isExtensionDisplayed?: boolean
  resourceIconSize?: TypeIconSize
  lazy?: boolean
}
const getSpaceMock = (props = {}) => ({
  id: 'lorem-id',
  name: 'Space 1',
  path: '',
  type: 'space',
  isFolder: true,
  disabled: false,
  spaceId: '1',
  getDriveAliasAndItem: () => '1',
  ...props
})

describe('OcTile component', () => {
  it('renders default space correctly', () => {
    const wrapper = getWrapper({ resource: getSpaceMock() })
    expect(wrapper.html()).toMatchSnapshot()
  })
  it('renders disabled space correctly', () => {
    const wrapper = getWrapper({
      resource: getSpaceMock({ disabled: true }),
      isResourceDisabled: true
    })
    expect(wrapper.html()).toMatchSnapshot()
  })
  it('renders selected resource correctly', () => {
    const wrapper = getWrapper({ resource: getSpaceMock(), isResourceSelected: true })
    expect(wrapper.find('.oc-tile-card-selected').exists()).toBeTruthy()
  })
  it('should emit click event when resource is clicked', () => {
    const wrapper = getWrapper({ resource: getSpaceMock(), isResourceSelected: true })
    const resourceLink = wrapper.find('resource-link-stub')
    resourceLink.trigger('click')
    expect(wrapper.emitted()).toHaveProperty('click')
  })
  it('should load lazily and show shimmering tile cards', () => {
    const wrapper = getWrapper({ resource: getSpaceMock(), isResourceSelected: false, lazy: true })
    expect(wrapper.find('.oc-tile-card-lazy-shimmer').exists()).toBeTruthy()
  })
  it('should show locked resource', () => {
    const wrapper = getWrapper({
      resource: getSpaceMock({ locked: true }),
      isResourceSelected: true
    })
    const element = wrapper.find('.oc-tile-card-preview')
    expect(element.attributes('aria-label')).toEqual('This item is locked')
  })
  it.each(['xlarge, xxlarge, xxxlarge'])(
    'renders resource icon size correctly',
    (resourceIconSize) => {
      const wrapper = getWrapper({
        resource: getSpaceMock(),
        resourceIconSize: resourceIconSize as TypeIconSize
      })
      expect(wrapper.find('resource-icon-stub').attributes().size).toEqual(resourceIconSize)
    }
  )

  function getWrapper(props: Props) {
    return shallowMount(ResourceTile, {
      props,
      global: { plugins: [...defaultPlugins()], renderStubDefaultSlot: true }
    })
  }
})
