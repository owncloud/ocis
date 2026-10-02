import Count from './OcAvatarCount.vue'
import { shallowMount } from '../../testing'

describe('OcAvatarCount', () => {
  it('dynamically calculates font size', () => {
    const wrapper = shallowMount(Count, {
      props: {
        size: 100,
        count: 2
      }
    })

    expect((wrapper.element as HTMLElement).style.fontSize).toMatch('40px')
    expect(wrapper.html()).toMatchSnapshot()
  })
})
