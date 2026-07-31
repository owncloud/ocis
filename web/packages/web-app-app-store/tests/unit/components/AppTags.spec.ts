import AppTags from '../../../src/components/AppTags.vue'
import { mount } from '@ownclouders/web-test-helpers'
import { mock } from 'vitest-mock-extended'
import { App } from '../../../src/types'

const tags: string[] = ['someTag', 'anotherTag', 'wololo-tag']

const selectors = {
  button: '[data-testid="tag-button"]',
  markElement: '.mark-element'
}

describe('AppTags.vue', () => {
  it('renders one button per tag', () => {
    const { wrapper } = getWrapper()
    expect(wrapper.findAll(selectors.button)).toHaveLength(tags.length)
  })
  it('shows the tag text as button text', () => {
    const { wrapper } = getWrapper()
    const buttons = wrapper.findAll(selectors.button)
    for (let i = 0; i < buttons.length; i++) {
      expect(buttons[i].text()).toBe(tags[i])
    }
  })
  it('emits click event on tag click', () => {
    const { wrapper } = getWrapper()
    wrapper.find(selectors.button).trigger('click')
    expect(wrapper.emitted('click')).toBeTruthy()
  })
  it('applies mark-element css class to tag text for highlighting', () => {
    const { wrapper } = getWrapper()
    wrapper.findAll(selectors.button).forEach((button) => {
      expect(button.find(selectors.markElement).exists()).toBeTruthy()
    })
  })
})

const getWrapper = () => {
  const app = { ...mock<App>({}), tags }

  return {
    wrapper: mount(AppTags, {
      props: { app }
    })
  }
}
