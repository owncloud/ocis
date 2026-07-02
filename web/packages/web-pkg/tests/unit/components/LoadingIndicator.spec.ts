import LoadingIndicator from '../../../src/components/LoadingIndicator.vue'
import { defaultPlugins, shallowMount } from '@ownclouders/web-test-helpers'
import { mock } from 'vitest-mock-extended'
import { LoadingService } from '../../../src/services'
import { OcProgress } from '@ownclouders/design-system/components'

const selectors = {
  loadingIndicator: '#oc-loading-indicator',
  progressStub: 'oc-progress-stub'
}

describe('LoadingIndicator', () => {
  it('should not render when not loading', () => {
    const { wrapper } = getWrapper()
    expect(wrapper.find(selectors.loadingIndicator).exists()).toBeFalsy()
  })
  it('should render when loading', () => {
    const { wrapper } = getWrapper({ isLoading: true })
    expect(wrapper.find(selectors.loadingIndicator).exists()).toBeTruthy()
  })
  describe('indeterminate', () => {
    it('progress bar should be in indeterminate when no progress given', () => {
      const { wrapper } = getWrapper({ isLoading: true })
      expect(
        wrapper.findComponent<typeof OcProgress>(selectors.progressStub).props('indeterminate')
      ).toBeTruthy()
    })
    it('progress bar should not be in indeterminate when progress given', () => {
      const { wrapper } = getWrapper({ isLoading: true, currentProgress: 50 })
      expect(
        wrapper.findComponent<typeof OcProgress>(selectors.progressStub).props('indeterminate')
      ).toBeFalsy()
    })
  })
})

function getWrapper({
  isLoading = false,
  currentProgress = null
}: { isLoading?: boolean; currentProgress?: number } = {}) {
  const mocks = { $loadingService: mock<LoadingService>({ isLoading, currentProgress }) }
  return {
    wrapper: shallowMount(LoadingIndicator, {
      global: {
        plugins: [...defaultPlugins()],
        mocks,
        provide: mocks
      }
    })
  }
}
