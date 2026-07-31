import FileInfo from '../../../../../src/components/SideBar/Files/FileInfo.vue'
import {
  defaultComponentMocks,
  defaultPlugins,
  shallowMount,
  RouteLocation
} from '@ownclouders/web-test-helpers'
import { mock } from 'vitest-mock-extended'
import { Resource } from '@ownclouders/web-client'

const selectors = {
  name: '[data-testid="files-info-name"]'
}

describe('FileInfo', () => {
  it('shows file info', () => {
    const { wrapper } = createWrapper()
    expect(wrapper.find(selectors.name).exists()).toBeTruthy()
  })
})

function createWrapper() {
  const file = mock<Resource>({
    name: 'someFolder',
    webDavPath: '',
    type: 'folder',
    extension: ''
  })

  return {
    wrapper: shallowMount(FileInfo, {
      global: {
        plugins: [...defaultPlugins()],
        provide: {
          resource: file
        },
        mocks: {
          ...defaultComponentMocks({ currentRoute: mock<RouteLocation>({ path: '/files' }) })
        }
      }
    })
  }
}
