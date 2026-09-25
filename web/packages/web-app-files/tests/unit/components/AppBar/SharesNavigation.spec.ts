import SharesNavigation from '../../../../src/components/AppBar/SharesNavigation.vue'
import { CapabilityStore, locationSharesWithMe } from '@ownclouders/web-pkg'
import { mock } from 'vitest-mock-extended'
import { RouteRecordNormalized } from 'vue-router'
import {
  defaultPlugins,
  defaultStubs,
  shallowMount,
  defaultComponentMocks,
  RouteLocation
} from '@ownclouders/web-test-helpers'

const routes = [
  mock<RouteRecordNormalized>({
    path: '/files/shares/with-me/',
    name: 'files-shares-with-me'
  }),
  mock<RouteRecordNormalized>({
    path: '/files/shares/with-others/',
    name: 'files-shares-with-others'
  }),
  mock<RouteRecordNormalized>({
    path: '/files/shares/via-link/',
    name: 'files-shares-via-link'
  })
]

describe('SharesNavigation component', () => {
  it('renders a shares navigation for both mobile and a desktop viewports', () => {
    const { wrapper } = getWrapper()
    expect(wrapper.html()).toMatchSnapshot()
  })

  describe('when public sharing is disabled', () => {
    it('does not render the "Shared via link" entry on desktop or mobile', () => {
      const { wrapper } = getWrapper({ sharingPublicEnabled: false })
      expect(wrapper.findAll('[to="/files/shares/with-me/"]').length).toBe(2)
      expect(wrapper.findAll('[to="/files/shares/with-others/"]').length).toBe(2)
      expect(wrapper.find('[to="/files/shares/via-link/"]').exists()).toBeFalsy()
      expect(wrapper.html()).not.toContain('Shared via link')
    })

    it('does not throw and shows a sensible mobile toggle label when mounted on the via-link route', () => {
      expect(() =>
        getWrapper({
          currentRouteName: 'files-shares-via-link',
          sharingPublicEnabled: false
        })
      ).not.toThrow()

      const { wrapper } = getWrapper({
        currentRouteName: 'files-shares-via-link',
        sharingPublicEnabled: false
      })
      expect(wrapper.find('#shares_navigation_mobile').text()).toContain('Shared with me')
    })
  })

  describe('when public sharing is enabled', () => {
    it('renders all three entries on desktop and mobile', () => {
      const { wrapper } = getWrapper({ sharingPublicEnabled: true })
      expect(wrapper.findAll('[to="/files/shares/with-me/"]').length).toBe(2)
      expect(wrapper.findAll('[to="/files/shares/with-others/"]').length).toBe(2)
      expect(wrapper.findAll('[to="/files/shares/via-link/"]').length).toBe(2)
    })
  })
})

function getWrapper({
  currentRouteName = locationSharesWithMe.name,
  sharingPublicEnabled = true
} = {}) {
  const mocks = defaultComponentMocks({
    currentRoute: mock<RouteLocation>({ name: currentRouteName })
  })
  mocks.$router.getRoutes.mockImplementation(() => routes)
  const capabilities = {
    files_sharing: { public: { enabled: sharingPublicEnabled } }
  } satisfies Partial<CapabilityStore['capabilities']>
  return {
    mocks,
    wrapper: shallowMount(SharesNavigation, {
      global: {
        stubs: defaultStubs,
        renderStubDefaultSlot: true,
        mocks,
        provide: mocks,
        plugins: [
          ...defaultPlugins({
            piniaOptions: { capabilityState: { capabilities } }
          })
        ]
      }
    })
  }
}
