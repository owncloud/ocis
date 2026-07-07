import { mount } from '@ownclouders/web-test-helpers'
import { nextTick, type Ref } from 'vue'
import { useThemeStore } from '@ownclouders/web-pkg'
import HtmlEditorPane from '../../../src/components/HtmlEditorPane.vue'

// Mock the theme store with a real-shaped, mutable `currentTheme` ref. The editor
// reads `themeStore.currentTheme` directly (no optional chaining), and using a real
// ref lets the dark-mode watch be exercised as a genuine reactive dependency.
vi.mock('@ownclouders/web-pkg', async () => {
  const { ref } = await import('vue')
  const currentTheme = ref({ isDark: false })
  return { useThemeStore: vi.fn(() => ({ currentTheme })) }
})

describe('HtmlEditorPane', () => {
  // The mocked store hands out the same `currentTheme` ref the component reads. The
  // real store type unwraps it (the component reads it ref-free via `unref`), so the
  // cast lets the test flip dark mode on the very ref the watch tracks.
  const currentTheme = useThemeStore().currentTheme as unknown as Ref<{ isDark: boolean }>

  beforeEach(() => {
    currentTheme.value = { isDark: false }
  })

  it('renders a CodeMirror editor for empty content', () => {
    const { wrapper } = getWrapper('')
    expect(wrapper.find('.cm-editor').exists()).toBe(true)
  })

  it('renders a CodeMirror editor for a valid HTML document', () => {
    const { wrapper } = getWrapper('<!doctype html><html><body><h1>hi</h1></body></html>')
    expect(wrapper.find('.cm-editor').exists()).toBe(true)
    expect(wrapper.vm.getView().state.doc.toString()).toContain('<h1>hi</h1>')
  })

  it('emits update:modelValue when the document changes', async () => {
    const { wrapper } = getWrapper('')
    wrapper.vm.getView().dispatch({ changes: { from: 0, insert: '<p>x</p>' } })
    await nextTick()
    expect(wrapper.emitted('update:modelValue')?.at(-1)).toEqual(['<p>x</p>'])
  })

  it('applies external content changes to the editor', async () => {
    const { wrapper } = getWrapper('<p>one</p>')
    await wrapper.setProps({ modelValue: '<p>two</p>' })
    expect(wrapper.vm.getView().state.doc.toString()).toBe('<p>two</p>')
  })

  it('reconfigures the editor when the theme dark mode changes', async () => {
    const { wrapper } = getWrapper('')
    const view = wrapper.vm.getView()
    const dispatchSpy = vi.spyOn(view, 'dispatch')
    currentTheme.value = { isDark: true }
    await nextTick()
    // the dark-mode watch fired and pushed a reconfigure effect into the editor
    expect(dispatchSpy).toHaveBeenCalledWith(
      expect.objectContaining({ effects: expect.anything() })
    )
  })
})

function getWrapper(modelValue: string) {
  return {
    wrapper: mount(HtmlEditorPane, {
      props: { modelValue, isReadOnly: false },
      attachTo: document.body
    })
  }
}
