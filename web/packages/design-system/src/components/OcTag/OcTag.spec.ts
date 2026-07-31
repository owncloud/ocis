import { shallowMount } from '@ownclouders/web-test-helpers'
import Tag from './OcTag.vue'

describe('OcTag', () => {
  it('uses correct component when type is specified', () => {
    const wrapper = shallowMount(Tag, {
      props: {
        type: 'button'
      }
    })

    expect(wrapper.element.tagName.toLowerCase()).toMatch('button')
    expect(wrapper.html()).toMatchSnapshot()
  })

  it('emits a click event', async () => {
    const wrapper = shallowMount(Tag, {
      props: {
        type: 'a'
      }
    })

    await wrapper.trigger('click')
    expect(wrapper.emitted().click).toBeTruthy()
  })
})
