import { shallowMount } from '@ownclouders/web-test-helpers'
import Progress from './OcProgress.vue'

describe('OcProgress', () => {
  it('sets correct classes', () => {
    const wrapper = shallowMount(Progress, {
      props: {
        value: 3,
        max: 10,
        variation: 'warning',
        size: 'small'
      }
    })

    expect(wrapper.attributes('class')).toContain('oc-progress-small')
    expect(wrapper.attributes('class')).toContain('oc-progress-warning')
    expect(wrapper.html()).toMatchSnapshot()
  })
})
