import { defaultComponentMocks, defaultPlugins, shallowMount } from '@ownclouders/web-test-helpers'
import FolderViewModal from '../../../src/components/FolderViewModal.vue'
import { Modal } from '@ownclouders/web-pkg'
import { mock } from 'vitest-mock-extended'

const SERVER_URL = 'https://example.org/'

const SELECTORS = Object.freeze({
  iframe: '#iframe-folder-view'
})

describe('FolderViewModal', () => {
  it('should set iframe src', () => {
    const { wrapper } = getWrapper()
    const iframe = wrapper.find(SELECTORS.iframe)

    expect(iframe.attributes('src')).toEqual(
      'https://example.org/public-link?hide-logo=true&hide-app-switcher=true&hide-account-menu=true&hide-navigation=true&lang=en'
    )
  })

  it.each(['javascript:alert(1)', 'data:text/html,<script>alert(1)</script>', 'blob:fake'])(
    'should throw when publicLink has a non-http(s) scheme: %s',
    (invalidUrl) => {
      expect(() => getWrapper({ publicLink: invalidUrl })).toThrow('Invalid URL scheme for iframe')
    }
  )

  it('should throw when publicLink points to a different server', () => {
    expect(() => getWrapper({ publicLink: 'https://other.example.com/public-link' })).toThrow(
      'URL does not belong to this server'
    )
  })
})

function getWrapper({
  publicLink = 'https://example.org/public-link',
  serverUrl = SERVER_URL
} = {}) {
  const mocks = defaultComponentMocks()

  return {
    mocks,
    wrapper: shallowMount(FolderViewModal, {
      props: {
        modal: mock<Modal>(),
        publicLink,
        serverUrl
      },
      global: {
        plugins: defaultPlugins(),
        mocks,
        provide: mocks
      }
    })
  }
}
