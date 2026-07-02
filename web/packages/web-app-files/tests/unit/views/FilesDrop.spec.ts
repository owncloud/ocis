import FilesDrop from '../../../src/views/FilesDrop.vue'
import {
  defaultPlugins,
  mount,
  defaultComponentMocks,
  defaultStubs,
  RouteLocation
} from '@ownclouders/web-test-helpers'
import { mock, mockDeep } from 'vitest-mock-extended'
import { ClientService } from '@ownclouders/web-pkg'
import { ListFilesResult } from '@ownclouders/web-client/webdav'

describe('FilesDrop view', () => {
  describe('different files view states', () => {
    it('shows the loading spinner during loading', () => {
      const { wrapper } = getMountedWrapper()
      expect(wrapper.find('#app-loading-spinner').exists()).toBeTruthy()
    })
    it('shows the "resource-upload"-component after loading', async () => {
      const { wrapper } = getMountedWrapper()
      wrapper.vm.loading = false
      await wrapper.vm.$nextTick()
      expect(wrapper.find('#app-loading-spinner').exists()).toBeFalsy()
      expect(wrapper.find('resource-upload-stub').exists()).toBeTruthy()
    })
  })
})

function getMountedWrapper() {
  const $clientService = mockDeep<ClientService>()
  $clientService.webdav.listFiles.mockResolvedValue(mock<ListFilesResult>())
  const defaultMocks = {
    ...defaultComponentMocks({
      currentRoute: mock<RouteLocation>({ name: 'files-common-favorites' })
    }),
    $clientService: $clientService
  }

  return {
    mocks: defaultMocks,
    wrapper: mount(FilesDrop, {
      global: {
        plugins: [...defaultPlugins()],
        mocks: defaultMocks,
        provide: defaultMocks,
        stubs: defaultStubs
      }
    })
  }
}
