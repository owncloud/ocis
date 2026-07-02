import { defaultPlugins, mount, shallowMount } from '@ownclouders/web-test-helpers'
import OcSearchBar from './OcSearchBar.vue'

describe('OcSearchBar', () => {
  const selectors = {
    searchButton: '.oc-search-button',
    searchInput: '.oc-search-input',
    searchButtonWrapper: '.oc-search-button-wrapper',
    searchClearButton: '.oc-search-clear'
  }
  describe('search grid role', () => {
    it('should be undefined if filter is enabled', () => {
      const wrapper = getWrapper({ isFilter: true })
      expect(wrapper.attributes('role')).toBe(undefined)
    })
    it('should be search if filter is disabled', () => {
      const wrapper = getWrapper()
      expect(wrapper.attributes('role')).toBe('search')
    })
  })
  describe('small prop', () => {
    describe('when small is true', () => {
      const wrapper = getWrapper({ small: true })
      it('should set small search button', () => {
        expect(wrapper.find(selectors.searchButton).attributes('size')).toBe('small')
      })
      it('should set small search grid', () => {
        expect(wrapper.attributes('class')).toContain('oc-search-small')
      })
      it('should set spinner size as xsmall', () => {
        const spinnerStub = wrapper.find('oc-spinner-stub')
        expect(spinnerStub.attributes('size')).toBe('xsmall')
      })
    })
    it('should set medium search grid and search button if not enabled', () => {
      const wrapper = getWrapper({ small: false })
      expect(wrapper.attributes('class')).not.toContain('oc-search-small')
      expect(wrapper.find(selectors.searchButton).attributes('size')).toBe('medium')
    })
  })
  describe('icon prop', () => {
    describe('when icon prop is not false', () => {
      const wrapper = getWrapper({ icon: 'mdi-icon' })
      it('should render icon', () => {
        const iconStub = wrapper.find('oc-icon-stub[name="mdi-icon"]')
        expect(iconStub.exists()).toBeTruthy()
        expect(iconStub.attributes('name')).toBe('mdi-icon')
      })
    })
    it('should not render icon if false', () => {
      const wrapper = getWrapper({ icon: false })
      const iconStub = wrapper.find('.oc-search-icon')
      expect(iconStub.exists()).toBeFalsy()
    })
  })
  describe('loading prop', () => {
    describe('when loading', () => {
      const wrapper = getWrapper({ icon: 'mdi-icon', loading: true })
      it('should show spinner', () => {
        const spinnerStub = wrapper.find('oc-spinner-stub')
        expect(spinnerStub.exists()).toBeTruthy()
        expect(spinnerStub.attributes('style')).not.toBe('display: none;')
      })
      it('should not show icon if loading', () => {
        const iconStub = wrapper.find('oc-icon-stub')
        expect(iconStub.exists()).toBeTruthy()
        expect(iconStub.attributes('style')).toBe('display: none;')
      })
      it('should set search input as disabled', () => {
        const searchInput = wrapper.find(selectors.searchInput)
        expect(searchInput.attributes('disabled')).toBe('')
      })
      it('should set search button as disabled', () => {
        const wrapper = getWrapper({ icon: 'mdi-icon', loading: true, value: 'kiran' })
        const searchInput = wrapper.find(selectors.searchButton)
        expect(searchInput.attributes('disabled')).toBe('true')
      })
    })
    describe('when not loading', () => {
      const wrapper = getWrapper({ icon: 'mdi-icon', loading: false })
      it('should not show spinner', () => {
        const spinnerStub = wrapper.find('oc-spinner-stub')
        expect(spinnerStub.exists()).toBeTruthy()
        expect(spinnerStub.attributes('style')).toBe('display: none;')
      })
      it('should enable search input', () => {
        const searchInput = wrapper.find(selectors.searchInput)
        expect(searchInput.attributes('disabled')).toBe(undefined)
      })
      it('should enable search button', async () => {
        await wrapper.find(selectors.searchInput).setValue('a') // search query should also be not null
        expect(wrapper.find(selectors.searchButton).attributes('disabled')).toBe('false')
      })
    })
  })
  describe('button label prop', () => {
    it('should set the provided button label on search button', () => {
      const wrapper = getWrapper({ buttonLabel: 'Search Elastic' })
      expect(wrapper.find(selectors.searchButton).text()).toBe('Search Elastic')
    })
  })
  describe('button hidden prop', () => {
    it('should add invisible class to search button if enabled', () => {
      const wrapper = getWrapper({ buttonHidden: true })
      const searchButtonWrapper = wrapper.find(selectors.searchButtonWrapper)
      expect(searchButtonWrapper.attributes('class')).toContain('oc-invisible-sr')
    })
    it('should add button class to input if disabled', () => {
      const wrapper = getWrapper({ buttonHidden: false })
      const searchInput = wrapper.find(selectors.searchInput)
      expect(searchInput.attributes('class')).toContain('oc-search-input-button')
    })
  })
  describe('aria-label for the loading spinner', () => {
    it('should add provided loading accessible label', () => {
      const wrapper = getWrapper({ loadingAccessibleLabel: 'Spinner is spinning' })
      const spinnerEl = wrapper.find('oc-spinner-stub')
      expect(spinnerEl.attributes('arialabel')).toBe('Spinner is spinning')
    })
    it('should add loading accessible label if not provided', () => {
      const wrapper = getWrapper()
      const spinnerEl = wrapper.find('oc-spinner-stub')
      expect(spinnerEl.attributes('arialabel')).toBe('Loading results')
    })
  })
  describe('search input', () => {
    it('should set provided label as input aria label', () => {
      const wrapper = getWrapper()
      const searchInput = wrapper.find(selectors.searchInput)
      expect(searchInput.attributes('aria-label')).toBe('Test search label')
    })
    it('should set provided placeholder as input placeholder', () => {
      const wrapper = getWrapper({ placeholder: 'Start typing..' })
      const searchInput = wrapper.find(selectors.searchInput)
      expect(searchInput.attributes('placeholder')).toBe('Start typing..')
    })
    it('should emit input event on typing', async () => {
      const wrapper = getWrapper()
      const searchInput = wrapper.find(selectors.searchInput)
      expect(wrapper.emitted('input')).toBeFalsy()
      await searchInput.setValue('abc')
      expect(wrapper.emitted('input')).toBeTruthy()
    })
  })
  describe('type ahead prop', () => {
    it('the search event is triggered on each entered character if enabled', async () => {
      const wrapper = getWrapper({ typeAhead: true })
      const searchInput = wrapper.find(selectors.searchInput)
      expect(wrapper.emitted('input')).toBeFalsy()
      expect(wrapper.emitted('search')).toBeFalsy()
      await searchInput.setValue('a')
      expect(wrapper.emitted('input')).toBeTruthy()
      expect(wrapper.emitted('search')).toBeTruthy()
    })
    it('the search event is not triggered on each entered character if disabled', async () => {
      const wrapper = getWrapper({ typeAhead: false })
      const searchInput = wrapper.find(selectors.searchInput)
      expect(wrapper.emitted('input')).toBeFalsy()
      expect(wrapper.emitted('search')).toBeFalsy()
      await searchInput.setValue('a')
      expect(wrapper.emitted('input')).toBeTruthy()
      expect(wrapper.emitted('search')).toBeFalsy()
    })
  })
  describe('when search button is clicked', () => {
    it('should emit search event if search query is not null', async () => {
      const wrapper = getMountedWrapper()
      const searchInput = wrapper.find(selectors.searchInput)
      await searchInput.setValue('a')
      expect(wrapper.emitted('search')).toBeFalsy()
      const searchButton = wrapper.find(selectors.searchButton)
      await searchButton.trigger('click')
      expect(wrapper.emitted('search')).toBeTruthy()
    })
    it('should not emit search event if search query is null', async () => {
      const wrapper = getMountedWrapper()
      expect(wrapper.emitted('search')).toBeFalsy()
      const searchButton = wrapper.find(selectors.searchButton)
      await searchButton.trigger('click')
      expect(wrapper.emitted('search')).toBeFalsy()
    })
  })
})

function getWrapper(props = {}) {
  return shallowMount(OcSearchBar, {
    props: {
      ...props,
      label: 'Test search label'
    },
    global: {
      renderStubDefaultSlot: true,
      plugins: [...defaultPlugins()]
    }
  })
}
function getMountedWrapper() {
  return mount(OcSearchBar, {
    props: { label: 'abc' },
    global: {
      plugins: [...defaultPlugins()]
    }
  })
}
