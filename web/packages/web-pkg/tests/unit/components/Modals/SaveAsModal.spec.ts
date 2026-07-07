import SaveAsModal from '../../../../src/components/Modals/SaveAsModal.vue'
import {
  defaultComponentMocks,
  defaultPlugins,
  nextTicks,
  shallowMount
} from '@ownclouders/web-test-helpers'
import { mock, mockDeep } from 'vitest-mock-extended'
import { Resource, SpaceResource } from '@ownclouders/web-client'
import { ListFilesResult } from '@ownclouders/web-client/webdav'
import { Modal, useMessages, useModals } from '../../../../src/composables/piniaStores'
import { ClientService } from '../../../../src'

window.open = vi.fn()

describe('SaveAsModal', () => {
  describe('iframe', () => {
    it('sets the iframe src correctly', () => {
      const { wrapper } = getWrapper()
      expect((wrapper.vm as any).iframeSrc).toEqual(
        'http://localhost:3000/files-spaces-generic?hide-logo=true&embed=true&embed-target=location&embed-choose-file-name=true&embed-delegate-authentication=false&embed-choose-file-name-suggestion=test.txt'
      )
    })
    it('sets the iframe title correctly', () => {
      const { wrapper } = getWrapper()
      expect((wrapper.vm as any).iframeTitle).toEqual('ownCloud')
    })
  })
  describe('method "onLocationPick"', () => {
    it('does nothing if the event message does not equal "owncloud-embed:select"', () => {
      const { mocks } = getWrapper()

      expect(mocks.$clientService.webdav.listFiles).not.toHaveBeenCalled()
      expect(mocks.$clientService.webdav.putFileContents).not.toHaveBeenCalled()
      expect(window.open).not.toHaveBeenCalled()
    })
    it('saves the file when message does equal "owncloud-embed:select"', async () => {
      const { wrapper, mocks } = getWrapper()
      const modalStore = useModals()
      const messageStore = useMessages()

      mocks.$clientService.webdav.putFileContents.mockResolvedValue(mock<Resource>())
      ;(wrapper.vm as any).onLocationPick(
        mock<MessageEvent>({
          origin: window.location.origin,
          data: {
            name: 'owncloud-embed:select',
            data: {
              resources: [mock<Resource>({ storageId: '1', spaceId: '1' })],
              fileName: 'test with new name.txt'
            }
          }
        })
      )

      await nextTicks(4)
      expect(messageStore.showMessage).toHaveBeenCalled()
      expect(modalStore.removeModal).toHaveBeenCalled()
      expect(window.open).toHaveBeenCalled()
    })
    it('does nothing when the message originates from an untrusted origin', async () => {
      const { wrapper, mocks } = getWrapper()

      ;(wrapper.vm as any).onLocationPick(
        mock<MessageEvent>({
          origin: 'https://attacker.example.com',
          data: {
            name: 'owncloud-embed:select',
            data: {
              resources: [mock<Resource>({ storageId: '1', spaceId: '1' })],
              fileName: 'test with new name.txt'
            }
          }
        })
      )

      await nextTicks(4)
      expect(mocks.$clientService.webdav.putFileContents).not.toHaveBeenCalled()
      expect(window.open).not.toHaveBeenCalled()
    })
    it('shows an error message when the file when message does equal "owncloud-embed:select and request fails"', async () => {
      console.error = vi.fn()
      const { wrapper, mocks } = getWrapper()
      const modalStore = useModals()
      const messageStore = useMessages()

      mocks.$clientService.webdav.putFileContents.mockRejectedValue(new Error(''))
      ;(wrapper.vm as any).onLocationPick(
        mock<MessageEvent>({
          origin: window.location.origin,
          data: {
            name: 'owncloud-embed:select',
            data: {
              resources: [mock<Resource>({ storageId: '1', spaceId: '1' })],
              fileName: 'test with new name.txt'
            }
          }
        })
      )

      await nextTicks(4)
      expect(messageStore.showErrorMessage).toHaveBeenCalled()
      expect(modalStore.removeModal).toHaveBeenCalled()
      expect(window.open).not.toHaveBeenCalled()
    })
  })
})

function getWrapper() {
  const $clientService = mockDeep<ClientService>()
  const mocks = { ...defaultComponentMocks(), $clientService }
  mocks.$clientService.webdav.listFiles.mockResolvedValue(mock<ListFilesResult>({ children: [] }))

  return {
    mocks,
    wrapper: shallowMount(SaveAsModal, {
      props: {
        modal: mock<Modal>(),
        content: 'some text',
        originalResource: {
          id: '1',
          path: '/test.txt',
          name: 'test.txt',
          extension: 'txt',
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
