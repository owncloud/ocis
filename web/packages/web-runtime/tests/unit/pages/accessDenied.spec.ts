import accessDenied from '../../../src/pages/accessDenied.vue'
import { defaultComponentMocks, defaultPlugins, mount } from '@ownclouders/web-test-helpers'
import { mock } from 'vitest-mock-extended'
import { RouteLocationNormalizedLoaded } from 'vue-router'

const selectors = {
  logInAgainButton: '#exitAnchor',
  cardTitle: '.oc-login-card-title'
}

describe('access denied page', () => {
  it('renders component', () => {
    const { wrapper } = getWrapper()
    expect(wrapper.html()).toMatchSnapshot()
  })
  describe('"Log in again" button', () => {
    it('navigates to "loginUrl" if set in config', () => {
      const loginUrl = 'https://myidp.int/login'
      const { wrapper } = getWrapper({ loginUrl })

      const logInAgainButton = wrapper.find(selectors.logInAgainButton)
      const loginAgainUrl = new URL(logInAgainButton.attributes().href)
      loginAgainUrl.search = ''

      expect(logInAgainButton.exists()).toBeTruthy()
      expect(loginAgainUrl.toString()).toEqual(loginUrl)
    })
  })
  describe('reason=forbidden (still logged in, lacking permission)', () => {
    it('does not claim the user was logged out, and sends them back into the app', () => {
      const { wrapper } = getWrapper({ reason: 'forbidden' })

      expect(wrapper.find(selectors.cardTitle).text()).not.toContain('logged')
      expect(wrapper.html()).not.toContain('log out')

      const goBackButton = wrapper.find(selectors.logInAgainButton)
      const destination = Object.values(goBackButton.attributes()).join(' ') + goBackButton.html()
      expect(destination.toLowerCase()).not.toContain('login')
    })
  })
})

function getWrapper({ loginUrl = '', reason = '' } = {}) {
  const currentRoute = mock<RouteLocationNormalizedLoaded>({
    query: { ...(reason && { reason }) }
  })
  const mocks = {
    ...defaultComponentMocks({ currentRoute })
  }

  return {
    mocks,
    wrapper: mount(accessDenied, {
      global: {
        plugins: [...defaultPlugins({ piniaOptions: { configState: { options: { loginUrl } } } })],
        mocks,
        provide: mocks
      }
    })
  }
}
