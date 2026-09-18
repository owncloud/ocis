import { flushPromises } from '@vue/test-utils'
import { nextTick } from 'vue'
import { mount } from '@ownclouders/web-test-helpers-core'
import { defaultPlugins } from '../../../../src/testing'
import TextEditor from '../../../../src/components/TextEditor/TextEditor.vue'

// Mounts the real md-editor-v3 rather than the stub used by TextEditor.spec.ts: the `footers`
// prop picks defFooters slot children by index, so anything that shifts that index makes the
// footer silently render the wrong node. Only the real Footer component exercises that lookup.

class MockFileReader {
  onload: (() => void) | null = null
  onerror: (() => void) | null = null
  result: string | null = null
  readAsDataURL(file: File) {
    this.result = `data:${file.type};base64,${'A'.repeat(file.size)}`
    queueMicrotask(() => this.onload?.())
  }
}

describe('TextEditor - footer slot', () => {
  it('renders the slot content in the footer', async () => {
    vi.stubGlobal('FileReader', MockFileReader)
    const wrapper = mount(TextEditor, {
      props: {
        currentContent: '',
        markdownMode: true,
        isReadOnly: false,
        applicationConfig: { maxImageSize: 200, maxDocumentImageSize: 500 }
      },
      global: { plugins: [...defaultPlugins()] }
    })
    await flushPromises()
    expect(wrapper.find('.md-editor-footer .footer-links').exists()).toBe(true)

    const md = wrapper.findComponent({ name: 'MdEditorV3' })
    md.vm.$emit(
      'on-upload-img',
      [new File([new Uint8Array(100)], 'a.png', { type: 'image/png' })],
      () => {}
    )
    await nextTick()
    expect(wrapper.find('.md-editor-footer .footer-image-progress').exists()).toBe(true)

    await flushPromises()
    await nextTick()
    expect(wrapper.find('.md-editor-footer .footer-image-progress').exists()).toBe(false)
    expect(wrapper.find('.md-editor-footer .footer-links').exists()).toBe(true)
  })
})
