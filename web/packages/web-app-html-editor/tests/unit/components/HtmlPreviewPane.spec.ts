import { defaultPlugins, mount } from '@ownclouders/web-test-helpers'
import HtmlPreviewPane from '../../../src/components/HtmlPreviewPane.vue'

// These assertions pin the security-critical contract of the preview. The whole
// isolation model is "the iframe attribute string is exactly right", so a future
// refactor that loosens it must break a test rather than ship silently.
const FORBIDDEN_SANDBOX_TOKENS = [
  'allow-same-origin',
  'allow-top-navigation',
  'allow-top-navigation-by-user-activation',
  'allow-popups',
  'allow-popups-to-escape-sandbox',
  'allow-forms',
  'allow-modals',
  'allow-downloads',
  'allow-pointer-lock'
]

describe('HtmlPreviewPane', () => {
  it('renders the content as the iframe srcdoc (verbatim)', () => {
    const { wrapper } = getWrapper('<p>hello</p>')
    expect(wrapper.find('iframe').attributes('srcdoc')).toBe('<p>hello</p>')
  })

  it('injects content via srcdoc, never via a src URL', () => {
    const { wrapper } = getWrapper('<p>x</p>')
    const iframe = wrapper.find('iframe')
    expect(iframe.attributes('srcdoc')).toBeDefined()
    expect(iframe.attributes('src')).toBeUndefined()
  })

  it('uses the minimal sandbox (only allow-scripts)', () => {
    const { wrapper } = getWrapper('<p>x</p>')
    expect(wrapper.find('iframe').attributes('sandbox')).toBe('allow-scripts')
  })

  it('never grants any origin-escaping or interaction sandbox token', () => {
    const { wrapper } = getWrapper('<p>x</p>')
    const sandbox = wrapper.find('iframe').attributes('sandbox') ?? ''
    for (const token of FORBIDDEN_SANDBOX_TOKENS) {
      expect(sandbox).not.toContain(token)
    }
  })

  it('does not leak the referrer', () => {
    const { wrapper } = getWrapper('<p>x</p>')
    expect(wrapper.find('iframe').attributes('referrerpolicy')).toBe('no-referrer')
  })
})

function getWrapper(content: string) {
  return {
    wrapper: mount(HtmlPreviewPane, {
      props: { content },
      global: { plugins: [...defaultPlugins()] }
    })
  }
}
