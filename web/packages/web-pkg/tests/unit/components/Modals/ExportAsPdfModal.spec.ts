import ExportAsPdfModal from '../../../../src/components/Modals/ExportAsPdfModal.vue'
import { defaultComponentMocks, defaultPlugins, shallowMount } from '@ownclouders/web-test-helpers'
import { mock } from 'vitest-mock-extended'
import { Resource, SpaceResource } from '@ownclouders/web-client'
import { Modal, useMessages, useModals } from '../../../../src/composables/piniaStores'
import {
  ExportAsPdfWorkerReturnData,
  useExportAsPdfWorker
} from '../../../../src/composables/webWorkers/exportAsPdfWorker'

vi.mock('../../../../src/composables/webWorkers/exportAsPdfWorker')

// The worker runs PDF generation off-thread (and its import can't be resolved in tests),
// so we stub it: capture the callback the modal passes in to drive the result handling.
let workerCallback: (result: ExportAsPdfWorkerReturnData) => void
const startWorker = vi.fn((_folder, _space, _fileName, _content, callback) => {
  workerCallback = callback
  return Promise.resolve()
})

beforeEach(() => {
  startWorker.mockClear()
  vi.mocked(useExportAsPdfWorker).mockReturnValue({ startWorker })
})

const messageEvent = (name: string, origin = window.location.origin) =>
  mock<MessageEvent>({
    origin,
    data: {
      name,
      data: {
        resources: [mock<Resource>({ storageId: '1', spaceId: '1' })],
        fileName: 'test.pdf'
      }
    }
  })

describe('ExportAsPdfModal', () => {
  describe('iframe', () => {
    it('sets the iframe src correctly', () => {
      const { wrapper } = getWrapper()
      expect((wrapper.vm as any).iframeUrl.href).toEqual(
        'http://localhost:3000/files-spaces-generic?hide-logo=true&embed=true&embed-target=location&embed-choose-file-name=true&embed-delegate-authentication=false&embed-choose-file-name-suggestion=test.pdf'
      )
    })
    it('sets the iframe title correctly', () => {
      const { wrapper } = getWrapper()
      expect((wrapper.vm as any).iframeTitle).toEqual('ownCloud')
    })
  })

  describe('method "handleMessage"', () => {
    it('does nothing if the event message is neither "owncloud-embed:select" nor "owncloud-embed:cancel"', () => {
      const { wrapper } = getWrapper()
      const modalStore = useModals()
      ;(wrapper.vm as any).handleMessage(messageEvent('some-other-event'))
      expect(startWorker).not.toHaveBeenCalled()
      expect(modalStore.removeModal).not.toHaveBeenCalled()
    })

    it('starts the export worker when message does equal "owncloud-embed:select"', () => {
      const { wrapper } = getWrapper()
      const modalStore = useModals()
      ;(wrapper.vm as any).handleMessage(messageEvent('owncloud-embed:select'))
      expect(startWorker).toHaveBeenCalled()
      expect(modalStore.removeModal).toHaveBeenCalled()
    })

    it('closes the modal when message does equal "owncloud-embed:cancel"', () => {
      const { wrapper } = getWrapper()
      const modalStore = useModals()
      ;(wrapper.vm as any).handleMessage(messageEvent('owncloud-embed:cancel'))
      expect(startWorker).not.toHaveBeenCalled()
      expect(modalStore.removeModal).toHaveBeenCalled()
    })

    it('does nothing when the message originates from an untrusted origin', () => {
      const { wrapper } = getWrapper()
      const modalStore = useModals()
      ;(wrapper.vm as any).handleMessage(
        messageEvent('owncloud-embed:select', 'https://attacker.example.com')
      )
      expect(startWorker).not.toHaveBeenCalled()
      expect(modalStore.removeModal).not.toHaveBeenCalled()
    })
  })

  describe('export result handling', () => {
    it('shows a success message when the export succeeds', () => {
      const { wrapper } = getWrapper()
      const messageStore = useMessages()
      ;(wrapper.vm as any).handleMessage(messageEvent('owncloud-embed:select'))

      workerCallback({ successful: [mock<Resource>()], failed: [] })

      expect(messageStore.showMessage).toHaveBeenCalled()
      expect(messageStore.showErrorMessage).not.toHaveBeenCalled()
    })

    it('shows an error message when the export fails', () => {
      console.error = vi.fn()
      const { wrapper } = getWrapper()
      const messageStore = useMessages()
      ;(wrapper.vm as any).handleMessage(messageEvent('owncloud-embed:select'))

      workerCallback({ successful: [], failed: [{ resourceName: 'test.pdf', error: mock() }] })

      expect(messageStore.showErrorMessage).toHaveBeenCalled()
      expect(messageStore.showMessage).not.toHaveBeenCalled()
    })
  })
})

function getWrapper() {
  const mocks = defaultComponentMocks()

  return {
    mocks,
    wrapper: shallowMount(ExportAsPdfModal, {
      props: {
        modal: mock<Modal>(),
        content: 'some text',
        originalResource: {
          id: '1',
          path: '/test.md',
          name: 'test.md',
          extension: 'md',
          spaceId: '1'
        },
        parentFolderLink: {
          name: 'files-spaces-generic',
          params: {
            driveAliasAndItem: 'personal/admin'
          },
          query: {
            fileId:
              '61dcd768-0bc4-4dd5-975a-2fe2bc9bc664$f1e4f3ec-1f24-460d-9f9a-4416ab6ddb6b!36cce768-8c9d-45e4-9c7d-4c9611962a75'
          }
        }
      },
      global: {
        plugins: [
          ...defaultPlugins({
            piniaOptions: { spacesState: { spaces: [mock<SpaceResource>({ id: '1' })] } }
          })
        ],
        mocks,
        provide: mocks
      }
    })
  }
}
