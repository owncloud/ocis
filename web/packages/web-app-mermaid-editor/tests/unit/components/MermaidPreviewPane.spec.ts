import { defaultPlugins, mount } from '@ownclouders/web-test-helpers'
import { nextTick, type Ref } from 'vue'
import mermaid from 'mermaid'
import { useThemeStore } from '@ownclouders/web-pkg'
import MermaidPreviewPane from '../../../src/components/MermaidPreviewPane.vue'

vi.mock('mermaid', () => ({
  default: {
    initialize: vi.fn(),
    parse: vi.fn(),
    render: vi.fn()
  }
}))

vi.mock('@ownclouders/web-pkg', async () => {
  const { ref } = await import('vue')
  const currentTheme = ref({ isDark: false })
  return { useThemeStore: vi.fn(() => ({ currentTheme })) }
})

describe('MermaidPreviewPane', () => {
  const currentTheme = useThemeStore().currentTheme as unknown as Ref<{ isDark: boolean }>

  beforeEach(() => {
    currentTheme.value = { isDark: false }
    vi.mocked(mermaid.parse).mockReset()
    vi.mocked(mermaid.render).mockReset()
    vi.mocked(mermaid.initialize).mockReset()
  })

  it('renders nothing for empty content, without calling mermaid', async () => {
    const { wrapper } = getWrapper('')
    await flush()
    expect(mermaid.parse).not.toHaveBeenCalled()
    expect(wrapper.find('.mermaid-preview-pane-diagram').html()).not.toContain('<svg')
    expect(wrapper.find('.mermaid-preview-pane-error').exists()).toBe(false)
  })

  it('renders the sanitized svg for a diagram that parses and renders', async () => {
    vi.mocked(mermaid.parse).mockResolvedValue({ diagramType: 'flowchart' } as never)
    vi.mocked(mermaid.render).mockResolvedValue({ svg: '<svg><text>A--&gt;B</text></svg>' } as never)

    const { wrapper } = getWrapper('graph TD\n  A --> B')
    await flush()

    expect(wrapper.find('.mermaid-preview-pane-error').exists()).toBe(false)
    expect(wrapper.find('.mermaid-preview-pane-diagram').html()).toContain('A--&gt;B')
  })

  it('shows an error and never calls render when the syntax is invalid', async () => {
    vi.mocked(mermaid.parse).mockResolvedValue(false as never)

    const { wrapper } = getWrapper('not a mermaid diagram {{{')
    await flush()

    expect(mermaid.render).not.toHaveBeenCalled()
    expect(wrapper.find('.mermaid-preview-pane-error').exists()).toBe(true)
    expect(wrapper.find('.mermaid-preview-pane-error__text').text()).toContain(
      'could not be rendered'
    )
  })

  it('shows the render error message when render() rejects after a valid parse', async () => {
    vi.mocked(mermaid.parse).mockResolvedValue({ diagramType: 'flowchart' } as never)
    vi.mocked(mermaid.render).mockRejectedValue(new Error('boom'))

    const { wrapper } = getWrapper('graph TD')
    await flush()

    expect(wrapper.find('.mermaid-preview-pane-error__detail').text()).toBe('boom')
  })

  it('re-renders with the dark theme when the theme store flips to dark', async () => {
    vi.mocked(mermaid.parse).mockResolvedValue({ diagramType: 'flowchart' } as never)
    vi.mocked(mermaid.render).mockResolvedValue({ svg: '<svg></svg>' } as never)

    getWrapper('graph TD')
    await flush()
    expect(mermaid.initialize).toHaveBeenLastCalledWith(
      expect.objectContaining({ theme: 'default' })
    )

    currentTheme.value = { isDark: true }
    await flush()
    expect(mermaid.initialize).toHaveBeenLastCalledWith(expect.objectContaining({ theme: 'dark' }))
  })

  it('drops a stale render so it cannot clobber a fresher one', async () => {
    let resolveFirst: (value: { svg: string }) => void
    const firstRender = new Promise<{ svg: string }>((resolve) => {
      resolveFirst = resolve
    })
    vi.mocked(mermaid.parse).mockResolvedValue({ diagramType: 'flowchart' } as never)
    vi.mocked(mermaid.render)
      .mockReturnValueOnce(firstRender as never)
      .mockResolvedValueOnce({ svg: '<svg><text>second</text></svg>' } as never)

    const { wrapper } = getWrapper('graph TD')
    await wrapper.setProps({ content: 'graph LR' })
    await flush()

    // the second (newer) render has already resolved and painted its result
    expect(wrapper.find('.mermaid-preview-pane-diagram').html()).toContain('second')

    // the first (stale) render now resolves — it must not overwrite the newer result
    resolveFirst!({ svg: '<svg><text>first</text></svg>' })
    await flush()
    expect(wrapper.find('.mermaid-preview-pane-diagram').html()).toContain('second')
    expect(wrapper.find('.mermaid-preview-pane-diagram').html()).not.toContain('first')
  })
})

async function flush() {
  await nextTick()
  await nextTick()
  await nextTick()
}

function getWrapper(content: string) {
  return {
    wrapper: mount(MermaidPreviewPane, {
      props: { content },
      global: { plugins: [...defaultPlugins()] }
    })
  }
}
