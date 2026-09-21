import Table from './OcTableSimple.vue'
import { shallowMount } from '../../testing'

describe('OcTableSimple', () => {
  it('adds hover', () => {
    const wrapper = shallowMount(Table, {
      props: {
        hover: true
      }
    })

    expect(wrapper.attributes('class')).toContain('oc-table-simple-hover')
  })
})
