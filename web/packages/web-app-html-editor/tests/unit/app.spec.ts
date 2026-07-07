import { defaultPlugins, mount, PartialComponentProps } from '@ownclouders/web-test-helpers'
import { nextTick } from 'vue'
import App from '../../src/App.vue'
import HtmlToolbar from '../../src/components/HtmlToolbar.vue'
import HtmlEditorPane from '../../src/components/HtmlEditorPane.vue'
import HtmlPreviewPane from '../../src/components/HtmlPreviewPane.vue'
import { PREVIEW_SIZE_LIMIT } from '../../src/helpers/preview'

describe('HTML editor app', () => {
  it('renders the toolbar, the editor and the preview', () => {
    const { wrapper } = getWrapper()
    expect(wrapper.findComponent(HtmlToolbar).exists()).toBe(true)
    expect(wrapper.findComponent(HtmlEditorPane).exists()).toBe(true)
    expect(wrapper.findComponent(HtmlPreviewPane).exists()).toBe(true)
  })

  it('defaults to split view', () => {
    const { wrapper } = getWrapper()
    expect(wrapper.find('.html-editor-body').classes()).toContain('html-editor-body-split')
  })

  it('re-emits editor changes as update:currentContent', async () => {
    const { wrapper } = getWrapper()
    wrapper.findComponent(HtmlEditorPane).vm.$emit('update:modelValue', '<p>hi</p>')
    await nextTick()
    expect(wrapper.emitted('update:currentContent')?.[0]).toEqual(['<p>hi</p>'])
  })

  it('switches the view mode from the toolbar', async () => {
    const { wrapper } = getWrapper()
    wrapper.findComponent(HtmlToolbar).vm.$emit('changeMode', 'preview')
    await nextTick()
    expect(wrapper.find('.html-editor-body').classes()).toContain('html-editor-body-preview-only')
  })

  it('wraps the preview content with a strict iframe CSP', () => {
    const { wrapper } = getWrapper({ currentContent: '<h1>hi</h1>' })
    const content = wrapper.findComponent(HtmlPreviewPane).props('content') as string
    expect(content).toContain('Content-Security-Policy')
    expect(content).toContain("default-src 'none'")
    expect(content).toContain('<h1>hi</h1>')
  })

  it('feeds debounced content to the preview', async () => {
    const { wrapper } = getWrapper({ currentContent: '<h1>start</h1>' })
    expect(wrapper.findComponent(HtmlPreviewPane).props('content')).toContain('<h1>start</h1>')

    // Use real timers (no fake-timer manipulation) so this file can never leak
    // timer state into other test projects in the shared run.
    await wrapper.setProps({ currentContent: '<h1>changed</h1>' })
    await nextTick()
    // debounced: not updated immediately after the change
    expect(wrapper.findComponent(HtmlPreviewPane).props('content')).toContain('<h1>start</h1>')
    // wait out the 250ms preview debounce
    await new Promise((resolve) => setTimeout(resolve, 350))
    expect(wrapper.findComponent(HtmlPreviewPane).props('content')).toContain('<h1>changed</h1>')
  })

  it('pauses the live preview for large files until the user opts in', async () => {
    const big = '<p>'.repeat(Math.ceil((PREVIEW_SIZE_LIMIT + 100) / 3))
    const { wrapper } = getWrapper({ currentContent: big })
    // preview is paused: pane not rendered, opt-in button shown
    expect(wrapper.findComponent(HtmlPreviewPane).exists()).toBe(false)
    const renderButton = wrapper.find('.html-editor-preview-render')
    expect(renderButton.exists()).toBe(true)

    await renderButton.trigger('click')
    await nextTick()
    expect(wrapper.findComponent(HtmlPreviewPane).exists()).toBe(true)
    // rendered synchronously on opt-in, not after the 250ms debounce
    expect(wrapper.findComponent(HtmlPreviewPane).props('content')).toContain(
      'Content-Security-Policy'
    )
  })

  it('re-pauses after opt-in when the content changes (large-file guard re-arms)', async () => {
    const big = '<p>'.repeat(Math.ceil((PREVIEW_SIZE_LIMIT + 100) / 3))
    const { wrapper } = getWrapper({ currentContent: big })
    await wrapper.find('.html-editor-preview-render').trigger('click')
    await nextTick()
    expect(wrapper.findComponent(HtmlPreviewPane).exists()).toBe(true)

    // a later change (e.g. an external conflict-reload) must not silently
    // auto-render another large document — it has to re-pause and re-prompt
    await wrapper.setProps({ currentContent: big + '<p>more</p>' })
    await nextTick()
    expect(wrapper.findComponent(HtmlPreviewPane).exists()).toBe(false)
    expect(wrapper.find('.html-editor-preview-render').exists()).toBe(true)
  })
})

function getWrapper(props: PartialComponentProps<typeof App> = {}) {
  return {
    wrapper: mount(App, {
      props: {
        applicationConfig: {},
        currentContent: '',
        isReadOnly: false,
        resource: undefined,
        ...props
      },
      global: {
        plugins: [...defaultPlugins()],
        stubs: {
          HtmlEditorPane: true,
          HtmlPreviewPane: true,
          HtmlToolbar: true
        }
      }
    })
  }
}
