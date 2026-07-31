import { defaultPlugins, mount } from '@ownclouders/web-test-helpers'
import AppActions from '../../../src/components/AppActions.vue'
import { App, AppVersion } from '../../../src/types'
import { mock } from 'vitest-mock-extended'

const version1: AppVersion = {
  version: '1.0.0',
  url: 'https://example.com/app-1.0.0.zip'
}
const version2: AppVersion = {
  version: '1.1.0',
  url: 'https://example.com/app-1.1.0.zip'
}
const versions = [version1, version2]
const mostRecentVersion = version2

const selectors = {
  downloadButton: 'button'
}

describe('AppActions', () => {
  it('renders a "Download" button', () => {
    const { wrapper } = getWrapper({})
    expect(wrapper.find(selectors.downloadButton).text()).toBe('Download')
  })
  describe('calling the "download" handler', () => {
    it('uses the most recent version when none is specified', async () => {
      const { wrapper } = getWrapper({})
      await wrapper.find(selectors.downloadButton).trigger('click')
      expect(window.location.href).toBe(mostRecentVersion.url)
    })
    it('uses the version provided via props', async () => {
      const { wrapper } = getWrapper({ version: version1 })
      await wrapper.find(selectors.downloadButton).trigger('click')
      expect(window.location.href).toBe(version1.url)
    })
  })
})

const getWrapper = ({ version }: { version?: AppVersion }) => {
  const app = { ...mock<App>({}), versions, mostRecentVersion }

  return {
    wrapper: mount(AppActions, {
      props: {
        app,
        version
      },
      global: {
        plugins: [...defaultPlugins()]
      }
    })
  }
}
