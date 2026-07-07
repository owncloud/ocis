import { shallowMount } from '@ownclouders/web-test-helpers'
import Cell from './OcTableCell.vue'

describe('OcTableCell', () => {
  it('Uses correct element', () => {
    const wrapper = shallowMount(Cell, {
      props: {
        type: 'th',
        alignH: 'right',
        alignV: 'bottom',
        width: 'shrink'
      },
      slots: {
        default: 'Hello world!'
      }
    })

    expect(wrapper.element.tagName).toBe('TH')
    expect(wrapper.attributes('class')).toContain('oc-table-cell-align-right')
    expect(wrapper.attributes('class')).toContain('oc-table-cell-align-bottom')
    expect(wrapper.attributes('class')).toContain('oc-table-cell-width-shrink')
    expect(wrapper.html()).toMatchSnapshot()
  })
})
