import ListInfo from '../../../../src/components/FilesList/ListInfo.vue'
import { defaultComponentMocks, defaultPlugins, shallowMount } from '@ownclouders/web-test-helpers'
import { mock } from 'vitest-mock-extended'
import { Resource } from '@ownclouders/web-client'

describe('ListInfo', () => {
  it('it renders content summary correctly', () => {
    const wrapper = getWrapper()
    expect(wrapper.html()).toMatchSnapshot()
  })
})

function getWrapper() {
  const resources = [
    mock<Resource>({
      isFolder: false,
      size: 10000,
      type: 'file',
      name: 'file1'
    }),
    mock<Resource>({
      isFolder: false,
      size: 5,
      type: 'file',
      name: 'file2'
    }),
    mock<Resource>({
      isFolder: true,
      size: 40000,
      type: 'folder',
      name: 'folder1'
    }),
    mock<Resource>({
      isFolder: true,
      size: 0,
      type: 'folder',
      name: 'folder2'
    })
  ]
  const mocks = {
    ...defaultComponentMocks()
  }

  return shallowMount(ListInfo, {
    global: {
      mocks,
      plugins: [...defaultPlugins({ piniaOptions: { resourcesStore: { resources } } })],
      provide: mocks
    }
  })
}
