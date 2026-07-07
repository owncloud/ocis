import { defaultPlugins, mount, PartialComponentProps } from '@ownclouders/web-test-helpers'
import OcSelect from './OcSelect.vue'

const selectors = {
  ocSelect: '.oc-select',
  selectedOptions: '.vs__selected-options .vs__selected',
  deselectBtn: '.vs__selected-options .vs__deselect',
  deselectLockIcon: '.vs__deselect-lock',
  clearBtn: '.vs__clear',
  searchInput: '.vs__search',
  ocSpinner: '.oc-spinner',
  warningMessage: '.oc-text-input-warning',
  errorMessage: '.oc-text-input-danger',
  descriptionMessage: '.oc-text-input-description'
}

describe('OcSelect', () => {
  it('passes the options to the vue-select component', () => {
    const options = [{ label: 'label1' }, { label: 'label2' }]
    const wrapper = getWrapper({ options })
    expect(
      // options is just passed through, so it does not exist as prop on OcSelect, hence we need the any cast
      wrapper.findComponent<typeof OcSelect>(selectors.ocSelect).props('options' as any)
    ).toEqual(options)
  })
  it('shows ocSpinner component when loading', () => {
    const wrapper = getWrapper({ loading: true })
    expect(wrapper.find(selectors.ocSpinner).exists()).toBeTruthy()
  })
  it('triggers the "search:input"-event on search input', async () => {
    const wrapper = getWrapper()
    await wrapper.find(selectors.searchInput).trigger('input')
    expect(wrapper.emitted('search:input')).toBeDefined()
  })
  describe('clear button', () => {
    it('is hidden by default', () => {
      const options = [{ label: 'label1' }, { label: 'label2' }]
      const wrapper = getWrapper({ options, modelValue: options[0] })
      expect(wrapper.find(selectors.clearBtn).attributes('style')).toEqual('display: none;')
    })
    it('is visible if "clearable" is set to true', () => {
      const options = [{ label: 'label1' }, { label: 'label2' }]
      const wrapper = getWrapper({ options, modelValue: options[0], clearable: true })
      expect(wrapper.find(selectors.clearBtn).attributes('style')).toBeUndefined()
    })
  })
  describe('selected option', () => {
    it('displays', () => {
      const options = [{ label: 'label1' }, { label: 'label2' }]
      const wrapper = getWrapper({ options, modelValue: options[0] })
      expect(wrapper.findAll(selectors.selectedOptions).length).toBe(1)
      expect(wrapper.findAll(selectors.selectedOptions).at(0).text()).toEqual(options[0].label)
    })
    it('displays with a custom label property', () => {
      const options = [{ customLabel: 'label1' }, { customLabel: 'label2' }]
      const wrapper = getWrapper({ options, modelValue: options[0], optionLabel: 'customLabel' })
      expect(wrapper.findAll(selectors.selectedOptions).at(0).text()).toEqual(
        options[0].customLabel
      )
    })
    it('displays with a custom label function', () => {
      const options = [{ customLabel: 'label1' }, { customLabel: 'label2' }]
      const wrapper = getWrapper<(typeof options)[0]>({
        options,
        modelValue: options[0],
        getOptionLabel: (o) => o.customLabel
      })
      expect(wrapper.findAll(selectors.selectedOptions).at(0).text()).toEqual(
        options[0].customLabel
      )
    })
    it('can be cleared if multi-select is allowed', () => {
      const options = [{ label: 'label1' }, { label: 'label2' }]
      const wrapper = getWrapper({ options, modelValue: options[0], multiple: true })
      expect(wrapper.find(selectors.deselectBtn).exists()).toBeTruthy()
      expect(wrapper.find(selectors.deselectLockIcon).exists()).toBeFalsy()
    })
    it('can not be cleared if readonly', () => {
      const options = [{ label: 'label1', readonly: true }, { label: 'label2' }]
      const wrapper = getWrapper({ options, modelValue: options[0], multiple: true })
      expect(wrapper.find(selectors.deselectBtn).exists()).toBeFalsy()
      expect(wrapper.find(selectors.deselectLockIcon).exists()).toBeTruthy()
    })
  })
  describe('message', () => {
    it('displays a warning message', () => {
      const wrapper = getWrapper({ warningMessage: 'foo' })
      expect(wrapper.find(selectors.warningMessage).exists()).toBeTruthy()
    })
    it('displays an error message', () => {
      const wrapper = getWrapper({ errorMessage: 'foo' })
      expect(wrapper.find(selectors.errorMessage).exists()).toBeTruthy()
    })
    it('displays a description message', () => {
      const wrapper = getWrapper({ descriptionMessage: 'foo' })
      expect(wrapper.find(selectors.descriptionMessage).exists()).toBeTruthy()
    })
  })
})

function getWrapper<T>(
  props: Partial<
    Omit<PartialComponentProps<typeof OcSelect>, 'getOptionLabel'> & {
      options: T[]
      getOptionLabel: (o: T) => string
      modelValue: T
    }
  > = {}
) {
  return mount(OcSelect, {
    props: {
      label: 'Select label',
      ...props
    },
    global: {
      plugins: [...defaultPlugins()]
    }
  })
}
