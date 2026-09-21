import { defineComponent, nextTick } from 'vue'
import { defaultPlugins, mount } from '../../../../src/testing'
import { flushPromises } from '@vue/test-utils'
import { config } from 'md-editor-v3'
import { useMessages } from '../../../../src/composables'
import TextEditor from '../../../../src/components/TextEditor/TextEditor.vue'
import { AppConfigObject } from '../../../../src/apps'

vi.mock('md-editor-v3', () => ({
  config: vi.fn(),
  MdEditor: defineComponent({
    name: 'MdEditor',
    emits: ['on-upload-img', 'on-change'],
    template: '<div><slot name="defFooters" /></div>'
  }),
  MdPreview: defineComponent({
    name: 'MdPreview',
    template: '<div />'
  })
}))

// FileReader returns a URI where the base64 section is exactly file.size 'A' characters.
// URI length = file.size + len('data:<type>;base64,'), making lengths predictable.
class MockFileReader {
  onload: (() => void) | null = null
  onerror: ((e: ErrorEvent) => void) | null = null
  result: string | null = null

  readAsDataURL(file: File) {
    this.result = `data:${file.type};base64,${'A'.repeat(file.size)}`
    queueMicrotask(() => this.onload?.())
  }
}

beforeEach(() => {
  vi.stubGlobal('FileReader', MockFileReader)
})

const makeFile = (name: string, size: number, type = 'image/png') =>
  new File([new Uint8Array(size)], name, { type })

const contentWithInlinedBytes = (totalEncodedBytes: number) =>
  `![x](data:image/png;base64,${'A'.repeat(totalEncodedBytes)})`

function getWrapper(
  overrides: {
    currentContent?: string
    markdownMode?: boolean
    applicationConfig?: AppConfigObject
    resource?: any
    isReadOnly?: boolean
  } = {}
) {
  return mount(TextEditor, {
    props: {
      currentContent: '',
      markdownMode: true,
      isReadOnly: false,
      applicationConfig: { maxImageSize: 200, maxDocumentImageSize: 500 },
      ...overrides
    },
    global: { plugins: [...defaultPlugins()] }
  })
}

async function triggerUpload(
  wrapper: ReturnType<typeof getWrapper>,
  files: File[]
): Promise<ReturnType<typeof vi.fn>> {
  const callback = vi.fn()
  const mdEditor = wrapper.findComponent({ name: 'MdEditor' })
  await mdEditor.vm.$emit('on-upload-img', files, callback)
  await flushPromises()
  return callback
}

describe('TextEditor — image upload', () => {
  it('counts images already in the document against the budget', async () => {
    const intoEmpty = await triggerUpload(getWrapper(), [makeFile('a.png', 100)])
    const [accepted] = intoEmpty.mock.calls[0]
    expect(accepted).toHaveLength(1)
    expect(accepted[0].url).toContain('data:image/png;base64,')

    // The same file must be rejected once the document already holds enough to leave no room.
    const intoFilled = await triggerUpload(
      getWrapper({ currentContent: contentWithInlinedBytes(400) }),
      [makeFile('a.png', 100)]
    )
    expect(intoFilled).toHaveBeenCalledWith([])
  })

  it('rejects a file that exceeds the per-image size cap', async () => {
    const wrapper = getWrapper()
    const callback = await triggerUpload(wrapper, [makeFile('big.png', 250)])

    expect(callback).toHaveBeenCalledWith([])
    const { showErrorMessage } = useMessages()
    expect(showErrorMessage).toHaveBeenCalledTimes(1)
    expect(vi.mocked(showErrorMessage).mock.calls[0][0].desc).toContain('big.png')
    expect(vi.mocked(showErrorMessage).mock.calls[0][0].desc).toContain('200')
  })

  it('reports one message per reason rather than one per rejected file', async () => {
    await triggerUpload(getWrapper(), [
      makeFile('a.png', 250),
      makeFile('b.png', 250),
      makeFile('c.png', 250),
      makeFile('d.png', 250)
    ])

    const { showErrorMessage } = useMessages()
    expect(showErrorMessage).toHaveBeenCalledTimes(1)
    const { desc } = vi.mocked(showErrorMessage).mock.calls[0][0]
    // The file list is truncated so a large drop cannot produce an unbounded message.
    expect(desc).toContain('a.png, b.png, c.png')
    expect(desc).not.toContain('d.png')
    expect(desc).toContain('1')
  })

  it('rejects a file that would push the document over budget and reports remaining space', async () => {
    const existingBytes = 400
    const wrapper = getWrapper({ currentContent: contentWithInlinedBytes(existingBytes) })
    const callback = await triggerUpload(wrapper, [makeFile('over.png', 150)])

    expect(callback).toHaveBeenCalledWith([])
    // 400 bytes of base64 plus the 22-byte data URI prefix leaves 78 of the 500-byte budget.
    const { showErrorMessage } = useMessages()
    expect(vi.mocked(showErrorMessage).mock.calls[0][0].desc).toContain('78')
  })

  it('accepts files that fit and rejects those that would overflow in the same drop', async () => {
    const wrapper = getWrapper({
      applicationConfig: { maxImageSize: 200, maxDocumentImageSize: 300 }
    })
    const callback = await triggerUpload(wrapper, [
      makeFile('f1.png', 100),
      makeFile('f2.png', 100),
      makeFile('f3.png', 100)
    ])

    const [accepted] = callback.mock.calls[0]
    expect(accepted).toHaveLength(2)
    expect(useMessages().showErrorMessage).toHaveBeenCalledTimes(1)
  })

  it('keeps the document budget hard across uploads the parent has not echoed back yet', async () => {
    const wrapper = getWrapper({
      applicationConfig: { maxImageSize: 200, maxDocumentImageSize: 200 }
    })

    // 100 bytes plus the 22-byte data URI prefix fits; a second one of the same size does not.
    const first = await triggerUpload(wrapper, [makeFile('f1.png', 100)])
    expect(first.mock.calls[0][0]).toHaveLength(1)

    // currentContent is still '' here - the parent has not propagated the insertion.
    const second = await triggerUpload(wrapper, [makeFile('f2.png', 100)])
    expect(second).toHaveBeenCalledWith([])
    expect(useMessages().showErrorMessage).toHaveBeenCalledTimes(1)
  })

  it('shows a busy cue while the picked images are being read', async () => {
    const wrapper = getWrapper()
    const mdEditor = wrapper.findComponent({ name: 'MdEditor' })

    mdEditor.vm.$emit('on-upload-img', [makeFile('a.png', 100)], vi.fn())
    await nextTick()
    expect(wrapper.find('.footer-image-progress').exists()).toBe(true)

    await flushPromises()
    await nextTick()
    expect(wrapper.find('.footer-image-progress').exists()).toBe(false)
  })

  it('data URIs survive the DOMPurify sanitize call used by the editor', () => {
    const wrapper = getWrapper()
    const sanitize = wrapper.findComponent({ name: 'MdEditor' }).vm.$attrs.sanitize as (
      html: string
    ) => string
    const html = '<div><img src="data:image/png;base64,ABC123==" alt="test"></div>'
    expect(sanitize(html)).toContain('data:image/png;base64,ABC123==')
  })

  it('does not process uploads when isMarkdown is false', async () => {
    const wrapper = getWrapper({
      markdownMode: false,
      applicationConfig: { showPreviewOnlyMd: true, maxImageSize: 200, maxDocumentImageSize: 500 }
    })
    const callback = await triggerUpload(wrapper, [makeFile('img.png', 50)])

    expect(callback).toHaveBeenCalledWith([])
    expect(useMessages().showErrorMessage).not.toHaveBeenCalled()
  })

  it('silently skips non-image files without a validation message', async () => {
    const wrapper = getWrapper()
    const callback = await triggerUpload(wrapper, [makeFile('doc.txt', 10, 'text/plain')])

    expect(callback).toHaveBeenCalledWith([])
    expect(useMessages().showErrorMessage).not.toHaveBeenCalled()
  })

  // md-editor-v3 injects cropperjs from unpkg.com on mount unless an instance is
  // registered, so dropping this registration silently adds an external request.
  it('registers a cropper instance so md-editor-v3 skips the unpkg CDN injection', () => {
    getWrapper()

    const [{ editorExtensions }] = vi.mocked(config).mock.calls.at(-1)
    expect(editorExtensions.cropper.instance).toBeDefined()
  })
})
