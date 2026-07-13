import HelpMenu from '../../../../src/components/Topbar/HelpMenu.vue'
import { defaultPlugins, mount } from '@ownclouders/web-test-helpers'
import { WebThemeType } from '@ownclouders/web-pkg'
import { mock } from 'vitest-mock-extended'

describe('HelpMenu component', () => {
  describe('trigger button', () => {
    it('is rendered when at least one url is set', () => {
      const { wrapper } = getWrapper({ softwareLicenseUrl: 'https://example.com/license' })
      expect(wrapper.find('button').exists()).toBeTruthy()
    })
    it('is not rendered when both urls are unset', () => {
      const { wrapper } = getWrapper()
      expect(wrapper.find('button').exists()).toBeFalsy()
    })
  })
  describe('software license link', () => {
    it('is rendered when softwareLicenseUrl is set', () => {
      const { wrapper } = getWrapper({ softwareLicenseUrl: 'https://example.com/license' })
      const link = wrapper.find('[data-testid="help-menu-software-license-link"]')
      expect(link.exists()).toBeTruthy()
      expect(link.attributes('href')).toEqual('https://example.com/license')
    })
    it('is not rendered when softwareLicenseUrl is unset', () => {
      const { wrapper } = getWrapper({ helpPageUrl: 'https://example.com/help' })
      const link = wrapper.find('[data-testid="help-menu-software-license-link"]')
      expect(link.exists()).toBeFalsy()
    })
  })
  describe('help page link', () => {
    it('is rendered when helpPageUrl is set', () => {
      const { wrapper } = getWrapper({ helpPageUrl: 'https://example.com/help' })
      const link = wrapper.find('[data-testid="help-menu-help-page-link"]')
      expect(link.exists()).toBeTruthy()
      expect(link.attributes('href')).toEqual('https://example.com/help')
    })
    it('is not rendered when helpPageUrl is unset', () => {
      const { wrapper } = getWrapper({ softwareLicenseUrl: 'https://example.com/license' })
      const link = wrapper.find('[data-testid="help-menu-help-page-link"]')
      expect(link.exists()).toBeFalsy()
    })
  })
})

function getWrapper({
  softwareLicenseUrl = undefined,
  helpPageUrl = undefined
}: {
  softwareLicenseUrl?: string
  helpPageUrl?: string
} = {}) {
  return {
    wrapper: mount(HelpMenu, {
      global: {
        plugins: [
          ...defaultPlugins({
            piniaOptions: {
              themeState: {
                currentTheme: mock<WebThemeType>({
                  common: {
                    urls: {
                      softwareLicense: softwareLicenseUrl,
                      helpPage: helpPageUrl
                    }
                  }
                })
              }
            }
          })
        ]
      }
    })
  }
}
