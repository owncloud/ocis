import UserMenu from '../../../../src/components/Topbar/UserMenu.vue'
import {
  defaultPlugins,
  defaultStubs,
  mount,
  defaultComponentMocks,
  RouteLocation
} from '@ownclouders/web-test-helpers'
import { mock } from 'vitest-mock-extended'
import { SpaceResource } from '@ownclouders/web-client'
import { Instance, Quota } from '@ownclouders/web-client/graph/generated'

const totalQuota = 1000
const basicQuota = 300
const warningQuota = 810
const dangerQuota = 910

const noEmail = ''
const email = 'test@test.de'

const selectors = {
  instanceSwitcher: '[data-testid="instance-switcher"]',
  instanceSwitcherItem: '[data-testid="instance-switcher-item"]',
  instanceSwitcherShowAllButton: '[data-testid="instance-switcher-show-all-button"]'
}

describe('User Menu component', () => {
  describe('when user is not logged in', () => {
    it('renders a navigation with login button', () => {
      const wrapper = getMountedWrapper({}, noEmail, true)
      expect(wrapper.html()).toMatchSnapshot()
    })
  })
  describe('when quota and no email is set', () => {
    it('renders a navigation without email', () => {
      const wrapper = getMountedWrapper({ used: basicQuota, total: totalQuota }, noEmail)
      expect(wrapper.html()).toMatchSnapshot()
    })
  })
  describe('when no quota and email is set', () => {
    it('the user menu does not contain a quota', () => {
      const wrapper = getMountedWrapper(null, email)
      expect(wrapper.html()).toMatchSnapshot()
    })
  })
  describe('when no quota and no email is set', () => {
    it('the user menu does not contain a quota', () => {
      const wrapper = getMountedWrapper(null, noEmail)
      expect(wrapper.html()).toMatchSnapshot()
    })
  })
  describe('when quota is below 80%', () => {
    it('renders a primary quota progress bar', () => {
      const wrapper = getMountedWrapper(
        {
          used: basicQuota,
          total: totalQuota
        },
        email
      )
      expect(wrapper.html()).toMatchSnapshot()
    })
  })
  describe('when quota is above 80% and below 90%', () => {
    it('renders a warning quota progress bar', () => {
      const wrapper = getMountedWrapper(
        {
          used: warningQuota,
          total: totalQuota
        },
        email
      )
      expect(wrapper.html()).toMatchSnapshot()
    })
  })
  describe('when quota is above 90%', () => {
    it('renders a danger quota progress bar', () => {
      const wrapper = getMountedWrapper(
        {
          used: dangerQuota,
          total: totalQuota
        },
        email
      )
      expect(wrapper.html()).toMatchSnapshot()
    })
  })
  describe('when quota is unlimited', () => {
    it('renders no percentag of total and no progress bar', () => {
      const wrapper = getMountedWrapper(
        {
          used: basicQuota,
          total: 0
        },
        email
      )
      expect(wrapper.html()).toMatchSnapshot()
    })
  })
  describe('when quota is not defined', () => {
    it('renders no percentag of total and no progress bar', () => {
      const wrapper = getMountedWrapper(
        {
          used: dangerQuota,
          total: 0
        },
        email
      )
      expect(wrapper.html()).toMatchSnapshot()
    })
  })
  describe('imprint and privacy urls', () => {
    it('should renders imprint and privacy section if urls are defined', () => {
      const wrapper = getMountedWrapper(
        {
          used: dangerQuota,
          total: totalQuota
        },
        email,
        false,
        true
      )
      const element = wrapper.find('.imprint-footer')
      expect(element.exists()).toBeTruthy()
      const output = element.html()
      expect(output).toContain('https://imprint.url')
      expect(output).toContain('https://privacy.url')
    })
  })
  describe('instance switcher', () => {
    it('should render instance switcher if there is at least one instance', () => {
      const wrapper = getMountedWrapper(
        { used: dangerQuota, total: totalQuota },
        email,
        false,
        true,
        [mock<Instance>({ url: 'https://instance1.com', primary: true })]
      )
      expect(wrapper.find(selectors.instanceSwitcher).exists()).toBe(true)
    })
    it('should not render instance switcher if there are no instances', () => {
      const wrapper = getMountedWrapper(
        { used: dangerQuota, total: totalQuota },
        email,
        false,
        true
      )
      expect(wrapper.find(selectors.instanceSwitcher).exists()).toBeFalsy()
    })
    it('should not render more instances than the inline limit', () => {
      const wrapper = getMountedWrapper(
        { used: dangerQuota, total: totalQuota },
        email,
        false,
        true,
        [
          mock<Instance>({ url: 'https://instance1.com', primary: true }),
          mock<Instance>({ url: 'https://instance2.com', primary: false }),
          mock<Instance>({ url: 'https://instance3.com', primary: false }),
          mock<Instance>({ url: 'https://instance4.com', primary: false })
        ]
      )
      expect(wrapper.find(selectors.instanceSwitcher).exists()).toBe(true)
      expect(wrapper.findAll(selectors.instanceSwitcherItem).length).toBe(3)
    })
    it('should render a button to show all instances if there are more instances than the inline limit', () => {
      const wrapper = getMountedWrapper(
        { used: dangerQuota, total: totalQuota },
        email,
        false,
        true,
        [
          mock<Instance>({ url: 'https://instance1.com', primary: true }),
          mock<Instance>({ url: 'https://instance2.com', primary: false }),
          mock<Instance>({ url: 'https://instance3.com', primary: false }),
          mock<Instance>({ url: 'https://instance4.com', primary: false })
        ]
      )
      expect(wrapper.find(selectors.instanceSwitcherShowAllButton).exists()).toBe(true)
    })
    it('should not render a button to show all instances if there are no more instances than the inline limit', () => {
      const wrapper = getMountedWrapper(
        { used: dangerQuota, total: totalQuota },
        email,
        false,
        true,
        [mock<Instance>({ url: 'https://instance1.com', primary: true })]
      )
      expect(wrapper.find(selectors.instanceSwitcherShowAllButton).exists()).toBeFalsy()
    })
  })
})

const getMountedWrapper = (
  quota: Quota,
  userEmail: string,
  noUser = false,
  areThemeUrlsSet = false,
  instances = []
) => {
  const mocks = {
    ...defaultComponentMocks({
      currentRoute: mock<RouteLocation>({ path: '/files', fullPath: '/files' })
    })
  }

  return mount(UserMenu, {
    global: {
      provide: mocks,
      renderStubDefaultSlot: true,
      plugins: [
        ...defaultPlugins({
          piniaOptions: {
            themeState: {
              currentTheme: {
                isDark: false,
                name: 'ownCloud',
                common: {
                  name: 'ownCloud',
                  slogan: 'ownCloud',
                  logo: 'https://logo.url.theme',
                  shareRoles: {},
                  urls: {
                    privacy: areThemeUrlsSet ? 'https://privacy.url.theme' : '',
                    imprint: areThemeUrlsSet ? 'https://imprint.url.theme' : '',
                    accessDeniedHelp: areThemeUrlsSet ? 'https://access-denied-help.url.theme' : ''
                  }
                }
              }
            },
            userState: {
              user: noUser
                ? null
                : {
                    id: '1',
                    onPremisesSamAccountName: 'einstein',
                    displayName: 'Albert Einstein',
                    mail: userEmail || '',
                    instances
                  }
            },
            spacesState: {
              spaces: [
                mock<SpaceResource>({
                  spaceQuota: quota,
                  isOwner: () => true,
                  driveType: 'personal'
                })
              ]
            }
          }
        })
      ],
      stubs: {
        ...defaultStubs,
        'oc-button': true,
        'oc-drop': true,
        'oc-list': true,
        'avatar-image': true,
        'oc-icon': true,
        'oc-progress': true
      },
      mocks
    }
  })
}
