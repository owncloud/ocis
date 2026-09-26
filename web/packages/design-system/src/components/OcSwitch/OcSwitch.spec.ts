import Switch from './OcSwitch.vue'
import { shallowMount } from '../../testing'

const defaultProps = {
  label: 'Test label'
}

describe('OcSwitch', () => {
  it('can be toggled', async () => {
    const wrapper = shallowMount(Switch, {
      props: defaultProps
    })

    await wrapper.find('[data-testid="oc-switch-btn"]').trigger('click')

    expect(wrapper.emitted('update:checked')[0][0]).toEqual(true)

    await wrapper.find('[data-testid="oc-switch-btn"]').trigger('click')

    expect(wrapper.emitted('update:checked')[0][0]).toEqual(true)
  })
})
