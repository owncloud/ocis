import Datepicker from './OcDatepicker.vue'
import { ComponentProps, defaultPlugins, shallowMount } from '@ownclouders/web-test-helpers'
import { DateTime } from 'luxon'
import { nextTick } from 'vue'

describe('OcDatePicker', () => {
  it('renders', () => {
    const wrapper = getWrapper({ label: 'Datepicker label' })
    expect(wrapper.html()).toMatchSnapshot()
  })
  it('sets the initial date correctly', async () => {
    const wrapper = getWrapper({ label: 'Datepicker label', currentDate: DateTime.now() })
    await nextTick()
    const inputEl = wrapper.find('.oc-text-input').element as HTMLInputElement
    expect(inputEl.value).toEqual(DateTime.now().toISODate())
  })
  it('sets the minimum date correctly', async () => {
    const wrapper = getWrapper({ label: 'Datepicker label', minDate: DateTime.now() })
    await nextTick()
    const inputEl = wrapper.find('.oc-text-input')
    expect(inputEl.attributes('min')).toEqual(DateTime.now().toISODate())
  })
  it('emits event on date change', async () => {
    const wrapper = getWrapper({ label: 'Datepicker label' })
    const inputEl = wrapper.find('.oc-text-input')
    await inputEl.setValue(DateTime.now().toISODate())
    expect(wrapper.emitted('dateChanged')).toBeTruthy()
  })
})

function getWrapper(props: ComponentProps<typeof Datepicker>) {
  return shallowMount(Datepicker, {
    props,
    global: {
      plugins: [...defaultPlugins()],
      stubs: {
        OcTextInput: false
      }
    }
  })
}
