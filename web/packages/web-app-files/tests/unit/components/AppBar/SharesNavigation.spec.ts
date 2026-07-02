import SharesNavigation from '../../../../src/components/AppBar/SharesNavigation.vue'
import { locationSharesWithMe } from '@ownclouders/web-pkg'
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
})

function getWrapper({ currentRouteName = locationSharesWithMe.name } = {}) {
  const mocks = defaultComponentMocks({
    currentRoute: mock<RouteLocation>({ name: currentRouteName })
  })
  mocks.$router.getRoutes.mockImplementation(() => routes)
  return {
    mocks,
    wrapper: shallowMount(SharesNavigation, {
      global: {
        stubs: defaultStubs,
        renderStubDefaultSlot: true,
        mocks,
        provide: mocks,
        plugins: [...defaultPlugins()]
      }
    })
  }
}
