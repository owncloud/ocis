import { shallowMount } from '@ownclouders/web-test-helpers-core'

import Table from './OcTableSimple.vue'

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
