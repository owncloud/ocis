import { Capabilities } from '@ownclouders/web-client/ocs'
import {
  ApplicationInformation,
  AppMenuItemExtension,
  useExtensionRegistry,
  WebThemeType
} from '@ownclouders/web-pkg'
import { mock } from 'vitest-mock-extended'
import { computed } from 'vue'
import TopBar from '../../../../src/components/Topbar/TopBar.vue'
import defaultTheme from '@ownclouders/web-test-helpers/src/mocks/theme.json'
import {
  defaultComponentMocks,
  defaultPlugins,
  PiniaMockOptions,
  shallowMount
} from '@ownclouders/web-test-helpers'

const defaultOwnCloudTheme = {
  ...defaultTheme.clients.web.defaults,
  ...defaultTheme.clients.web.themes[0]
}

const mockUseEmbedMode = vi.fn().mockReturnValue({ isEnabled: computed(() => false) })

vi.mock('@ownclouders/web-pkg', async (importOriginal) => ({
  ...(await importOriginal<any>()),
  useEmbedMode: vi.fn().mockImplementation(() => mockUseEmbedMode())
}))

describe('Top Bar component', () => {
  describe('applications menu', () => {
    it('Displays applications menu', () => {
      const { wrapper } = getWrapper()
      expect(wrapper.find('applications-menu-stub').exists()).toBeTruthy()
    })
    it('should not display when hideAppSwitcher is "true"', () => {
      const { wrapper } = getWrapper({ options: { hideAppSwitcher: true } })
      expect(wrapper.find('applications-menu-stub').exists()).toBeFalsy()
    })
  })
  describe('notifications bell', () => {
    it('should display in authenticated context if announced via capabilities', () => {
      const { wrapper } = getWrapper({
        capabilities: {
          notifications: { 'ocs-endpoints': ['list', 'get', 'delete'] }
        }
      })
      expect(wrapper.find('notifications-stub').exists()).toBeTruthy()
    })
    it('should not display in an unauthenticated context', () => {
      const { wrapper } = getWrapper({
        userContextReady: false,
        capabilities: {
          notifications: { 'ocs-endpoints': ['list', 'get', 'delete'] }
        }
      })
      expect(wrapper.find('notifications-stub').exists()).toBeFalsy()
    })
    it('should not display if endpoint list is missing', () => {
      const { wrapper } = getWrapper({
        capabilities: { notifications: { 'ocs-endpoints': [] } }
      })
      expect(wrapper.find('notifications-stub').exists()).toBeFalsy()
    })
  })
  it.each(['applications-menu', 'feedback-link', 'notifications', 'user-menu'])(
    'should hide %s when mode is "embed"',
    (componentName) => {
      mockUseEmbedMode.mockReturnValue({
        isEnabled: computed(() => true)
      })

      const { wrapper } = getWrapper()
      expect(wrapper.find(`${componentName}-stub`).exists()).toBeFalsy()
    }
  )
  it.each(['applications-menu', 'feedback-link', 'notifications', 'user-menu'])(
    'should not hide %s when mode is not "embed"',
    (componentName) => {
      mockUseEmbedMode.mockReturnValue({
        isEnabled: computed(() => false)
      })

      const { wrapper } = getWrapper({
        capabilities: {
          notifications: { 'ocs-endpoints': ['list', 'get', 'delete'] }
        }
      })
      expect(wrapper.find(`${componentName}-stub`).exists()).toBeTruthy()
    }
  )
  it.each(['feedback-link', 'notifications', 'user-menu'])(
    'should hide %s when hideAccountMenu is "true"',
    (componentName) => {
      const { wrapper } = getWrapper({ options: { hideAccountMenu: true } })
      expect(wrapper.find(`${componentName}-stub`).exists()).toBeFalsy()
    }
  )
  describe('logo', () => {
    it('links to the internal home route when no href is configured', () => {
      const { wrapper } = getWrapper()
      expect(wrapper.find('.oc-logo-href').attributes('href')).toBeUndefined()
      expect(wrapper.find('router-link-stub').exists()).toBeTruthy()
    })
    it('links to the configured href when set', () => {
      const { wrapper } = getWrapper({ logoHref: 'https://example.com' })
      expect(wrapper.find('router-link-stub').exists()).toBeFalsy()
      expect(wrapper.find('.oc-logo-href').attributes('href')).toBe('https://example.com')
    })
    it('is not rendered when hideLogo is "true", even if href is configured', () => {
      const { wrapper } = getWrapper({
        options: { hideLogo: true },
        logoHref: 'https://example.com'
      })
      expect(wrapper.find('.oc-logo-href').exists()).toBeFalsy()
    })
  })
})

const getWrapper = ({
  capabilities = {},
  userContextReady = true,
  options,
  logoHref
}: {
  capabilities?: Partial<Capabilities['capabilities']>
  userContextReady?: boolean
  options?: PiniaMockOptions['configState']['options']
  logoHref?: string
} = {}) => {
  const mocks = { ...defaultComponentMocks() }

  const plugins = defaultPlugins({
    piniaOptions: {
      authState: { userContextReady },
      capabilityState: { capabilities },
      configState: { options: { disableFeedbackLink: false, ...options } },
      themeState: {
        currentTheme: mock<WebThemeType>({
          ...defaultOwnCloudTheme,
          logo: { ...defaultOwnCloudTheme.logo, href: logoHref }
        })
      }
    }
  })

  const extensionRegistry = useExtensionRegistry()
  vi.mocked(extensionRegistry.requestExtensions).mockReturnValue([mock<AppMenuItemExtension>()])

  return {
    wrapper: shallowMount(TopBar, {
      props: {
        applicationsList: [
          mock<ApplicationInformation>({
            icon: ''
          })
        ]
      },
      global: {
        plugins,
        stubs: { 'router-link': true, 'portal-target': true, notifications: true },
        mocks,
        provide: mocks
      }
    })
  }
}
