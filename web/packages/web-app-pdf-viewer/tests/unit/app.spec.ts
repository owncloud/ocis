import App from '../../src/App.vue'
import { defaultPlugins, shallowMount } from '@ownclouders/web-test-helpers'

const selectors = {
  object: 'object',
  openButton: 'oc-button-stub'
}

afterEach(() => {
  vi.restoreAllMocks()
})

describe('PDF Viewer App', () => {
  describe('on non-iOS', () => {
    it('renders the object element with the PDF URL', () => {
      const { wrapper } = getWrapper()
      expect(wrapper.find(selectors.object).exists()).toBe(true)
      expect(wrapper.find(selectors.object).attributes('data')).toBe('blob:test-url')
    })

    it('does not render the Open PDF button', () => {
      const { wrapper } = getWrapper()
      expect(wrapper.find(selectors.openButton).exists()).toBe(false)
    })
  })

  describe('on iOS', () => {
    beforeEach(() => {
      vi.spyOn(navigator, 'userAgent', 'get').mockReturnValue(
        'Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/605.1.15'
      )
    })

    it('does not render the object element', () => {
      const { wrapper } = getWrapper()
      expect(wrapper.find(selectors.object).exists()).toBe(false)
    })

    it('renders the Open PDF button linking to the PDF URL', () => {
      const { wrapper } = getWrapper()
      expect(wrapper.find(selectors.openButton).attributes('href')).toBe('blob:test-url')
    })

    it('opens the PDF in a new tab', () => {
      const { wrapper } = getWrapper()
      expect(wrapper.find(selectors.openButton).attributes('target')).toBe('_blank')
    })
  })

  describe('on iPadOS', () => {
    beforeEach(() => {
      vi.spyOn(navigator, 'userAgent', 'get').mockReturnValue(
        'Mozilla/5.0 (iPad; CPU OS 17_0 like Mac OS X) AppleWebKit/605.1.15'
      )
    })

    it('does not render the object element', () => {
      const { wrapper } = getWrapper()
      expect(wrapper.find(selectors.object).exists()).toBe(false)
    })

    it('renders the Open PDF button linking to the PDF URL', () => {
      const { wrapper } = getWrapper()
      expect(wrapper.find(selectors.openButton).attributes('href')).toBe('blob:test-url')
    })
  })

  describe('on iPadOS in desktop mode', () => {
    beforeEach(() => {
      vi.spyOn(navigator, 'userAgent', 'get').mockReturnValue(
        'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15'
      )
    })

    it('renders the object element', () => {
      const { wrapper } = getWrapper()
      expect(wrapper.find(selectors.object).exists()).toBe(true)
    })

    it('does not render the Open PDF button', () => {
      const { wrapper } = getWrapper()
      expect(wrapper.find(selectors.openButton).exists()).toBe(false)
    })
  })
})

function getWrapper() {
  return {
    wrapper: shallowMount(App, {
      props: { url: 'blob:test-url' },
      global: { plugins: [...defaultPlugins()] }
    })
  }
}
