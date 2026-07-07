import { mockDeep } from 'vitest-mock-extended'
import ResourceUpload from '../../../../../src/components/AppBar/Upload/ResourceUpload.vue'
import {
  defaultComponentMocks,
  defaultPlugins,
  defaultStubs,
  mount
} from '@ownclouders/web-test-helpers'
import { UppyService } from '@ownclouders/web-pkg'
import { OcButton } from '@ownclouders/design-system/components'

describe('Resource Upload Component', () => {
  describe('file upload', () => {
    it('should render component', () => {
      const { wrapper } = getWrapper()
      expect(wrapper.exists()).toBeTruthy()
      expect(wrapper.html()).toMatchSnapshot()
    })
  })

  describe('folder upload', () => {
    it('should render component', () => {
      const { wrapper } = getWrapper({ isFolder: true })
      expect(wrapper.exists()).toBeTruthy()
      expect(wrapper.html()).toMatchSnapshot()
    })
  })

  describe('when upload file button is clicked', () => {
    it('should call "triggerUpload"', async () => {
      const { wrapper } = getWrapper()

      const spyTriggerUpload = vi.spyOn(wrapper.vm, 'triggerUpload')
      const uploadButton = wrapper.find('button')
      const fileUploadInput = wrapper.find('#files-file-upload-input')

      ;(fileUploadInput.element as HTMLElement).click = vi.fn()
      await wrapper.vm.$forceUpdate()

      await uploadButton.trigger('click')

      expect(spyTriggerUpload).toHaveBeenCalledTimes(1)
      expect((fileUploadInput.element as HTMLElement).click).toHaveBeenCalledTimes(1)
    })
  })

  it('should be disabled when a remote upload is running', () => {
    const uppyService = mockDeep<UppyService>()
    uppyService.isRemoteUploadInProgress.mockReturnValue(true)
    const { wrapper } = getWrapper({ isFolder: true }, uppyService)
    expect(wrapper.findComponent<typeof OcButton>('button').props('disabled')).toBeTruthy()
  })
})

function getWrapper(props = {}, uppyService = mockDeep<UppyService>()) {
  const mocks = {
    ...defaultComponentMocks(),
    $uppyService: uppyService
  }
  return {
    mocks,
    wrapper: mount(ResourceUpload, {
      props,
      global: {
        mocks,
        stubs: defaultStubs,
        provide: mocks,
        plugins: [...defaultPlugins()]
      }
    })
  }
}
