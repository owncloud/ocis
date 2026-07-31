import { defaultPlugins, mount } from '@ownclouders/web-test-helpers'
import HtmlToolbar from '../../../src/components/HtmlToolbar.vue'

describe('HtmlToolbar', () => {
  it('renders the three view-mode buttons', () => {
    const { wrapper } = getWrapper('split')
    expect(wrapper.findAllComponents({ name: 'OcButton' })).toHaveLength(3)
    expect(wrapper.find('.html-editor-viewmode-editor').exists()).toBe(true)
    expect(wrapper.find('.html-editor-viewmode-split').exists()).toBe(true)
    expect(wrapper.find('.html-editor-viewmode-preview').exists()).toBe(true)
  })

  it('marks the active view mode with the filled appearance', () => {
    const { wrapper } = getWrapper('preview')
    const buttons = wrapper.findAllComponents({ name: 'OcButton' })
    expect(buttons[0].props('appearance')).toBe('outline') // editor
    expect(buttons[1].props('appearance')).toBe('outline') // split
    expect(buttons[2].props('appearance')).toBe('filled') // preview
  })

  it('emits changeMode with the chosen mode', async () => {
    const { wrapper } = getWrapper('split')
    await wrapper.find('.html-editor-viewmode-preview').trigger('click')
    expect(wrapper.emitted('changeMode')?.[0]).toEqual(['preview'])
  })
})

function getWrapper(viewMode: 'editor' | 'split' | 'preview') {
  return {
    wrapper: mount(HtmlToolbar, {
      props: { viewMode },
      global: { plugins: [...defaultPlugins()] }
    })
  }
}
