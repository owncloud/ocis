import DriveRedirect from '../../../../src/views/spaces/DriveRedirect.vue'
import { mock } from 'vitest-mock-extended'
import {
  defaultPlugins,
  mount,
  defaultComponentMocks,
  defaultStubs,
  RouteLocation
} from '@ownclouders/web-test-helpers'

const selectors = Object.freeze({
  spaceNotFound: '#files-space-not-found'
})

describe('DriveRedirect view', () => {
  it('redirects to "projects" route if no personal space exist', () => {
    const { mocks } = getMountedWrapper()
    expect(mocks.$router.replace).toHaveBeenCalledWith({
      name: 'files-spaces-projects'
    })
  })

  it('should show a message if the space is not found', () => {
    const { mocks, wrapper } = getMountedWrapper({ props: { driveAliasAndItem: 'missing-space' } })

    expect(wrapper.find(selectors.spaceNotFound).exists()).toBe(true)
    expect(mocks.$router.replace).not.toHaveBeenCalledWith({
      name: 'files-spaces-projects'
    })
  })

  it('should redirect to personal space if the alias is personal drive', () => {
    const { mocks } = getMountedWrapper({ props: { driveAliasAndItem: 'personal' } })
    expect(mocks.$router.replace).toHaveBeenCalledWith({
      name: 'files-spaces-projects'
    })
  })

  it('should redirect to personal space if the alias is the fake personal drive alias', () => {
    const { mocks } = getMountedWrapper({ props: { driveAliasAndItem: 'personal/home' } })
    expect(mocks.$router.replace).toHaveBeenCalledWith({
      name: 'files-spaces-projects'
    })
  })
})

function getMountedWrapper({ currentRouteName = 'files-spaces-generic', props = {} } = {}) {
  const mocks = {
    ...defaultComponentMocks({ currentRoute: mock<RouteLocation>({ name: currentRouteName }) })
  }

  return {
    mocks,
    wrapper: mount(DriveRedirect, {
      props,
      global: {
        plugins: defaultPlugins(),
        stubs: defaultStubs,
        mocks,
        provide: mocks
      }
    })
  }
}
