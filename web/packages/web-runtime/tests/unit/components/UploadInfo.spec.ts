import UploadInfo from '../../../src/components/UploadInfo.vue'
import { defaultPlugins, shallowMount, defaultComponentMocks } from '@ownclouders/web-test-helpers'
import { ResourceListItem, OcUppyFile } from '@ownclouders/web-pkg'
import { nextTick } from 'vue'
import { HttpError } from '@ownclouders/web-client'

const selectors = {
  overlay: '#upload-info',
  title: '.upload-info-title p',
  progress: '.upload-info-progress',
  success: '.upload-info-success',
  error: '.upload-info-danger',
  message: '.upload-info-message',
  info: {
    items: '.upload-info-items',
    item: '.upload-info-items li'
  }
}

describe('UploadInfo component', () => {
  it('should render the component in a hidden state per default', () => {
    const { wrapper } = getShallowWrapper()
    const overlay = wrapper.find(selectors.overlay)
    expect(overlay.exists()).toBeFalsy()
  })
  it('should show the component', async () => {
    const { wrapper } = getShallowWrapper()
    wrapper.vm.showInfo = true
    await nextTick()
    const overlay = wrapper.find(selectors.overlay)
    expect(overlay.exists()).toBeTruthy()
  })
  describe('title', () => {
    it('should show that an upload is in progress', async () => {
      const { wrapper } = getShallowWrapper()
      wrapper.vm.showInfo = true
      wrapper.vm.inPreparation = false
      wrapper.vm.filesInProgressCount = 1
      wrapper.vm.runningUploads = 1
      await nextTick()

      const uploadTitle = wrapper.find(selectors.title).text()
      expect(uploadTitle).toBe('1 item uploading...')
    })
    it('should show that an upload was successful', async () => {
      const { wrapper } = getShallowWrapper()
      wrapper.vm.showInfo = true
      wrapper.vm.successful = ['1']
      await nextTick()

      const uploadTitle = wrapper.find(selectors.title).text()
      expect(uploadTitle).toBe('Upload complete')
    })
    it('should show that an upload failed', async () => {
      const { wrapper } = getShallowWrapper()
      wrapper.vm.showInfo = true
      ;((wrapper.vm.errors = { '1': new HttpError('', undefined) }), await nextTick())

      const uploadTitle = wrapper.find(selectors.title).text()
      expect(uploadTitle).toBe('Upload failed')
    })
    it('should show that an upload was cancelled', async () => {
      const { wrapper } = getShallowWrapper()
      wrapper.vm.showInfo = true
      wrapper.vm.uploadsCancelled = true
      await nextTick()

      const uploadTitle = wrapper.find(selectors.title).text()
      expect(uploadTitle).toBe('Upload cancelled')
    })
    it('should show that an upload is preparing', async () => {
      const { wrapper } = getShallowWrapper()
      wrapper.vm.showInfo = true
      wrapper.vm.runningUploads = 1
      wrapper.vm.inPreparation = true
      await nextTick()

      const uploadTitle = wrapper.find(selectors.title).text()
      expect(uploadTitle).toBe('Preparing upload...')
    })
    it('should show that an upload is being finalized', async () => {
      const { wrapper } = getShallowWrapper()
      wrapper.vm.showInfo = true
      wrapper.vm.filesInProgressCount = 1
      wrapper.vm.runningUploads = 1
      wrapper.vm.inFinalization = true
      await nextTick()

      const uploadTitle = wrapper.find(selectors.title).text()
      expect(uploadTitle).toBe('Finalizing upload...')
    })
  })
  describe('progress bar', () => {
    it('should show the progress bar when an upload is in progress', async () => {
      const { wrapper } = getShallowWrapper()
      wrapper.vm.showInfo = true
      wrapper.vm.filesInProgressCount = 1
      wrapper.vm.runningUploads = 1
      await nextTick()

      const progressBar = wrapper.find(selectors.progress)
      expect(progressBar.exists()).toBeTruthy()
    })
  })
  describe('info', () => {
    it('should show the number of successful items', async () => {
      const { wrapper } = getShallowWrapper()
      wrapper.vm.showInfo = true
      wrapper.vm.successful = ['1', '2']
      await nextTick()

      const info = wrapper.find(selectors.success).text()
      expect(info).toBe('2 items uploaded')
    })
    it('should show the number of failed items', async () => {
      const { wrapper } = getShallowWrapper()
      wrapper.vm.showInfo = true
      wrapper.vm.errors = { '1': new HttpError('', undefined) }
      wrapper.vm.successful = ['1']
      await nextTick()

      const info = wrapper.find(selectors.error).text()
      expect(info).toBe('1 of 2 items failed')
    })
  })
  describe('details', () => {
    it('should hide the info by default', async () => {
      const { wrapper } = getShallowWrapper()
      wrapper.vm.showInfo = true
      await nextTick()

      const info = wrapper.find(selectors.info.items)
      expect(info.exists()).toBeFalsy()
    })
    it('should list all the uploaded files when the info is displayed', async () => {
      const { wrapper } = getShallowWrapper()
      wrapper.vm.showInfo = true
      wrapper.vm.infoExpanded = true
      wrapper.vm.uploads = {
        '1': {
          name: 'file',
          path: '/',
          type: 'file',
          meta: { uploadId: '1' }
        } as unknown as OcUppyFile,
        '2': {
          name: 'file2',
          path: '/',
          type: 'file',
          meta: { uploadId: '2' }
        } as unknown as OcUppyFile
      }

      await nextTick()

      const info = wrapper.find(selectors.info.items)
      expect(info.exists()).toBeTruthy()

      const uploadedItems = wrapper.findAll(selectors.info.item)
      expect(uploadedItems.length).toBe(2)
    })
    it('should show a message on the failed uploaded files', async () => {
      const { wrapper } = getShallowWrapper()
      wrapper.vm.showInfo = true
      wrapper.vm.infoExpanded = true
      ;((wrapper.vm.uploads = {
        '1': {
          name: 'file',
          path: '/',
          type: 'file',
          meta: { uploadId: '1' }
        } as unknown as OcUppyFile,
        '2': {
          name: 'file2',
          path: '/',
          type: 'file',
          meta: { uploadId: '2' }
        } as unknown as OcUppyFile,
        '3': {
          name: 'file3',
          path: '/',
          type: 'file',
          meta: { uploadId: '3' }
        } as unknown as OcUppyFile
      }),
        (wrapper.vm.errors = {
          1: new HttpError('', undefined),
          2: new HttpError('', undefined)
        }))
      await nextTick()

      const info = wrapper.find(selectors.info.items)
      expect(info.exists()).toBeTruthy()

      const infoMessages = wrapper.findAll(selectors.message)
      expect(infoMessages.length).toBe(2)
      expect(infoMessages.at(0).text()).toBe('Unknown error')
      expect(infoMessages.at(1).text()).toBe('Unknown error')
    })
    it('folder is clickable', async () => {
      const { wrapper } = getShallowWrapper()
      wrapper.vm.showInfo = true
      wrapper.vm.infoExpanded = true
      wrapper.vm.uploads = {
        '1': {
          name: 'file',
          type: 'folder',
          targetRoute: { params: { driveAliasAndItem: 'some/drive/alias' } },
          path: '',
          meta: { uploadId: '1', isFolder: true }
        } as unknown as OcUppyFile
      }
      await nextTick()

      const info = wrapper.find(selectors.info.items)
      expect(info.exists()).toBeTruthy()
      const resourceStub = wrapper.findComponent<typeof ResourceListItem>(
        `${selectors.info.item} resource-list-item-stub`
      )
      expect(resourceStub.props().isResourceClickable).toBeTruthy()
    })
  })
  describe('getRemainingTime method', () => {
    it.each([
      { ms: 1000, expected: 'Few seconds left' },
      { ms: 1000 * 60 * 30, expected: '30 minutes left' },
      { ms: 1000 * 60 * 60, expected: '1 hour left' },
      { ms: 1000 * 60 * 60 * 2, expected: '2 hours left' }
    ])('should return the proper string', ({ ms, expected }) => {
      const { wrapper } = getShallowWrapper()
      const estimatedTime = wrapper.vm.getRemainingTime(ms)
      expect(estimatedTime).toBe(expected)
    })
  })
})

function getShallowWrapper() {
  const mocks = {
    ...defaultComponentMocks()
  }

  return {
    mocks,
    wrapper: shallowMount(UploadInfo, {
      global: {
        plugins: [...defaultPlugins()],
        mocks
      }
    })
  }
}
