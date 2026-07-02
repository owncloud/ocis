import FilesViewWrapper from '../../../src/components/FilesViewWrapper.vue'
import { defaultComponentMocks, defaultPlugins, RouteLocation } from '@ownclouders/web-test-helpers'
import { mock } from 'vitest-mock-extended'
import { shallowMount } from '@vue/test-utils'
import { computed } from 'vue'

const selectors = {
  embedActionsStub: 'embed-actions-stub'
}

const mockUseEmbedMode = vi.fn().mockReturnValue({
  isEnabled: computed(() => false)
})

vi.mock('@ownclouders/web-pkg', async (importOriginal) => ({
  ...(await importOriginal<any>()),
  useRoute: vi.fn().mockReturnValue({ query: {} }),
  useEmbedMode: vi.fn().mockImplementation(() => mockUseEmbedMode())
}))

describe('FilesViewWrapper', () => {
  describe('embed actions', () => {
    it('renders when embed mode is enabled', () => {
      mockUseEmbedMode.mockReturnValue({
        isEnabled: computed(() => true)
      })
      const { wrapper } = getWrapper()
      expect(wrapper.findComponent(selectors.embedActionsStub).exists()).toBeTruthy()
    })
    it('does not render when embed mode is disabled', () => {
      mockUseEmbedMode.mockReturnValue({
        isEnabled: computed(() => false)
      })
      const { wrapper } = getWrapper()
      expect(wrapper.findComponent(selectors.embedActionsStub).exists()).toBeFalsy()
    })
  })
})

function getWrapper() {
  const mocks = defaultComponentMocks({ currentRoute: mock<RouteLocation>({ path: '/files' }) })
  return {
    wrapper: shallowMount(FilesViewWrapper, {
      global: {
        mocks,
        provide: mocks,
        plugins: [...defaultPlugins({})]
      }
    })
  }
}
