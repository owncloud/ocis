import { ref } from 'vue'
import SpaceHeader from '../../../../src/components/Spaces/SpaceHeader.vue'
import { DriveItem } from '@ownclouders/web-client/graph/generated'
import { SpaceResource, Resource, buildSpaceImageResource } from '@ownclouders/web-client'

import {
  defaultPlugins,
  mount,
  defaultComponentMocks,
  flushPromises
} from '@ownclouders/web-test-helpers'
import { mock } from 'vitest-mock-extended'

vi.mock('@ownclouders/web-pkg', async (importOriginal) => ({
  ...(await importOriginal<any>()),
  useFileActions: vi.fn().mockReturnValue({
    getDefaultAction: vi.fn().mockReturnValue({ handler: vi.fn() })
  }),
  useLoadPreview: vi.fn().mockReturnValue({
    loadPreview: vi.fn(() => 'blob:image')
  })
}))

vi.mock('@ownclouders/web-client', async (importOriginal) => ({
  ...(await importOriginal<any>()),
  buildSpaceImageResource: vi.fn()
}))

const getSpaceMock = (spaceImageData: DriveItem = undefined, options?: Partial<SpaceResource>) =>
  mock<SpaceResource>({
    id: '1',
    name: '',
    description: '',
    spaceReadmeData: undefined,
    spaceImageData,
    ...options
  })

describe('SpaceHeader', () => {
  it('should add the "squashed"-class when the sidebar is opened', () => {
    const wrapper = getWrapper({ space: getSpaceMock(), isSideBarOpen: true })
    expect(wrapper.find('.space-header-squashed').exists()).toBeTruthy()
    expect(wrapper.html()).toMatchSnapshot()
  })
  describe('space image', () => {
    it('should show the default image if no other image is set', () => {
      const wrapper = getWrapper({ space: getSpaceMock() })
      expect(wrapper.find('.space-header-image-default').exists()).toBeTruthy()
      expect(wrapper.html()).toMatchSnapshot()
    })
    it('should show the set image', async () => {
      const wrapper = getWrapper({ space: getSpaceMock({ webDavUrl: '/' }) })
      await wrapper.vm.$nextTick()
      expect(wrapper.find('.space-header-image-default').exists()).toBeFalsy()
      expect(wrapper.find('.space-header-image img').exists()).toBeTruthy()
      expect(wrapper.html()).toMatchSnapshot()
    })
    it('should take full width in mobile view', () => {
      const wrapper = getWrapper({
        space: getSpaceMock({ webDavUrl: '/' }),
        isMobileWidth: true
      })
      expect(wrapper.find('.space-header').attributes('class')).not.toContain('oc-flex')
      expect(wrapper.find('.space-header-image').attributes('class')).toContain(
        'space-header-image-expanded'
      )
    })
  })
  describe('space description', () => {
    it('should show the description', async () => {
      const wrapper = getWrapper({
        space: getSpaceMock(undefined, {
          spaceReadmeData: {
            name: 'lorem'
          }
        })
      })
      await flushPromises()
      expect(wrapper.find('.markdown-container').exists()).toBeTruthy()
      expect(wrapper.html()).toMatchSnapshot()
    })
  })
})

function getWrapper({ space = {} as SpaceResource, isSideBarOpen = false, isMobileWidth = false }) {
  const mocks = defaultComponentMocks()
  mocks.$previewService.loadPreview.mockResolvedValue('blob:image')
  mocks.$clientService.webdav.getFileContents.mockResolvedValue({ body: 'lorem body' })
  mocks.$clientService.webdav.getFileInfo.mockResolvedValue({
    id: '1',
    path: 'lorem/path',
    spaceId: '1'
  })
  vi.mocked(buildSpaceImageResource).mockReturnValue(mock<Resource>())

  return mount(SpaceHeader, {
    props: {
      space,
      isSideBarOpen
    },
    global: {
      mocks,
      plugins: [...defaultPlugins()],
      provide: { ...mocks, isMobileWidth: ref(isMobileWidth) },
      stubs: {
        'quota-modal': true,
        'space-context-actions': true
      }
    }
  })
}
