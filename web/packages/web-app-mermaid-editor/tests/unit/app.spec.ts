import { defaultPlugins, mount, PartialComponentProps } from '@ownclouders/web-test-helpers'
import { nextTick } from 'vue'
import App from '../../src/App.vue'
import MermaidToolbar from '../../src/components/MermaidToolbar.vue'
import MermaidEditorPane from '../../src/components/MermaidEditorPane.vue'
import MermaidPreviewPane from '../../src/components/MermaidPreviewPane.vue'
import { PREVIEW_SIZE_LIMIT } from '../../src/helpers/preview'

describe('Mermaid editor app', () => {
  it('renders the toolbar, the editor and the preview', () => {
    const { wrapper } = getWrapper()
    expect(wrapper.findComponent(MermaidToolbar).exists()).toBe(true)
    expect(wrapper.findComponent(MermaidEditorPane).exists()).toBe(true)
    expect(wrapper.findComponent(MermaidPreviewPane).exists()).toBe(true)
  })

  it('defaults to split view', () => {
    const { wrapper } = getWrapper()
    expect(wrapper.find('.mermaid-editor-body').classes()).toContain('mermaid-editor-body-split')
  })

  it('re-emits editor changes as update:currentContent', async () => {
    const { wrapper } = getWrapper()
    wrapper.findComponent(MermaidEditorPane).vm.$emit('update:modelValue', 'graph TD')
    await nextTick()
    expect(wrapper.emitted('update:currentContent')?.[0]).toEqual(['graph TD'])
  })

  it('switches the view mode from the toolbar', async () => {
    const { wrapper } = getWrapper()
    wrapper.findComponent(MermaidToolbar).vm.$emit('changeMode', 'preview')
    await nextTick()
    expect(wrapper.find('.mermaid-editor-body').classes()).toContain(
      'mermaid-editor-body-preview-only'
    )
  })

  it('passes the raw diagram source to the preview, unwrapped', () => {
    const { wrapper } = getWrapper({ currentContent: 'graph TD\n  A --> B' })
    const content = wrapper.findComponent(MermaidPreviewPane).props('content') as string
    expect(content).toBe('graph TD\n  A --> B')
  })

  it('feeds debounced content to the preview', async () => {
    const { wrapper } = getWrapper({ currentContent: 'graph TD' })
    expect(wrapper.findComponent(MermaidPreviewPane).props('content')).toBe('graph TD')

    // Use real timers (no fake-timer manipulation) so this file can never leak
    // timer state into other test projects in the shared run.
    await wrapper.setProps({ currentContent: 'graph LR' })
    await nextTick()
    // debounced: not updated immediately after the change
    expect(wrapper.findComponent(MermaidPreviewPane).props('content')).toBe('graph TD')
    // wait out the 250ms preview debounce
    await new Promise((resolve) => setTimeout(resolve, 350))
    expect(wrapper.findComponent(MermaidPreviewPane).props('content')).toBe('graph LR')
  })

  it('pauses the live preview for large files until the user opts in', async () => {
    const big = 'A'.repeat(PREVIEW_SIZE_LIMIT + 100)
    const { wrapper } = getWrapper({ currentContent: big })
    // preview is paused: pane not rendered, opt-in button shown
    expect(wrapper.findComponent(MermaidPreviewPane).exists()).toBe(false)
    const renderButton = wrapper.find('.mermaid-editor-preview-render')
    expect(renderButton.exists()).toBe(true)

    await renderButton.trigger('click')
    await nextTick()
    expect(wrapper.findComponent(MermaidPreviewPane).exists()).toBe(true)
    // rendered synchronously on opt-in, not after the 250ms debounce
    expect(wrapper.findComponent(MermaidPreviewPane).props('content')).toBe(big)
  })

  it('re-pauses after opt-in when the content changes (large-file guard re-arms)', async () => {
    const big = 'A'.repeat(PREVIEW_SIZE_LIMIT + 100)
    const { wrapper } = getWrapper({ currentContent: big })
    await wrapper.find('.mermaid-editor-preview-render').trigger('click')
    await nextTick()
    expect(wrapper.findComponent(MermaidPreviewPane).exists()).toBe(true)

    // a later change (e.g. an external conflict-reload) must not silently
    // auto-render another large document — it has to re-pause and re-prompt
    await wrapper.setProps({ currentContent: big + 'more' })
    await nextTick()
    expect(wrapper.findComponent(MermaidPreviewPane).exists()).toBe(false)
    expect(wrapper.find('.mermaid-editor-preview-render').exists()).toBe(true)
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
          MermaidEditorPane: true,
          MermaidPreviewPane: true,
          MermaidToolbar: true
        }
      }
    })
  }
}
