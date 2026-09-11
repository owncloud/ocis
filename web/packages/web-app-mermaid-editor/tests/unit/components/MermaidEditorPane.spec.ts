import { mount } from '@ownclouders/web-test-helpers'
import { nextTick, type Ref } from 'vue'
import { useThemeStore } from '@ownclouders/web-pkg'
import MermaidEditorPane from '../../../src/components/MermaidEditorPane.vue'

// Mock the theme store with a real-shaped, mutable `currentTheme` ref. The editor
// reads `themeStore.currentTheme` directly (no optional chaining), and using a real
// ref lets the dark-mode watch be exercised as a genuine reactive dependency.
vi.mock('@ownclouders/web-pkg', async () => {
  const { ref } = await import('vue')
  const currentTheme = ref({ isDark: false })
  return { useThemeStore: vi.fn(() => ({ currentTheme })) }
})

describe('MermaidEditorPane', () => {
  const currentTheme = useThemeStore().currentTheme as unknown as Ref<{ isDark: boolean }>

  beforeEach(() => {
    currentTheme.value = { isDark: false }
  })

  it('renders a CodeMirror editor for empty content', () => {
    const { wrapper } = getWrapper('')
    expect(wrapper.find('.cm-editor').exists()).toBe(true)
  })

  it('renders a CodeMirror editor for mermaid diagram source', () => {
    const { wrapper } = getWrapper('graph TD\n  A --> B')
    expect(wrapper.find('.cm-editor').exists()).toBe(true)
    expect(wrapper.vm.getView().state.doc.toString()).toContain('A --> B')
  })

  it('emits update:modelValue when the document changes', async () => {
    const { wrapper } = getWrapper('')
    wrapper.vm.getView().dispatch({ changes: { from: 0, insert: 'graph TD' } })
    await nextTick()
    expect(wrapper.emitted('update:modelValue')?.at(-1)).toEqual(['graph TD'])
  })

  it('applies external content changes to the editor', async () => {
    const { wrapper } = getWrapper('graph TD')
    await wrapper.setProps({ modelValue: 'graph LR' })
    expect(wrapper.vm.getView().state.doc.toString()).toBe('graph LR')
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
    wrapper: mount(MermaidEditorPane, {
      props: { modelValue, isReadOnly: false },
      attachTo: document.body
    })
  }
}
