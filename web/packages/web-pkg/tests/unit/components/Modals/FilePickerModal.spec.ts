import FilePickerModal from '../../../../src/components/Modals/FilePickerModal.vue'
import { defaultComponentMocks, defaultPlugins, shallowMount } from '@ownclouders/web-test-helpers'
import { mock } from 'vitest-mock-extended'
import { Resource, SpaceResource } from '@ownclouders/web-client'
import { Modal, useModals } from '../../../../src/composables/piniaStores'
import { RouteLocation } from 'vue-router'

window.open = vi.fn()

describe('FilePickerModal', () => {
  describe('iframe', () => {
    it('sets the iframe src correctly', () => {
      const { wrapper } = getWrapper()
      expect((wrapper.vm as any).iframeSrc).toEqual(
        'http://localhost:3000/files-spaces-generic?hide-logo=true&embed=true&embed-target=file&embed-delegate-authentication=false&embed-file-types=text%2Cmd%2Ctext%2Frtf'
      )
    })
    it('sets the iframe title correctly', () => {
      const { wrapper } = getWrapper()
      expect((wrapper.vm as any).iframeTitle).toEqual('ownCloud')
    })
  })
  describe('method "onFilePick"', () => {
    it('does nothing if the event message does not equal "owncloud-embed:file-pick"', () => {
      const { wrapper } = getWrapper()
      ;(wrapper.vm as any).onFilePick(mock<MessageEvent>({ data: { name: 'some-other-event' } }))
      expect(window.open).not.toHaveBeenCalled()
    })
    it('opens resource in new window when message does equal "owncloud-embed:file-pick"', () => {
      const { wrapper } = getWrapper()
      const modalStore = useModals()
      ;(wrapper.vm as any).onFilePick(
        mock<MessageEvent>({
          origin: window.location.origin,
          data: {
            name: 'owncloud-embed:file-pick',
            data: {
              resource: mock<Resource>({ spaceId: '1' }),
              originRoute: mock<RouteLocation>()
            }
          }
        })
      )
      expect(modalStore.removeModal).toHaveBeenCalled()
      expect(window.open).toHaveBeenCalled()
    })
    it('does nothing when the message originates from an untrusted origin', () => {
      const { wrapper } = getWrapper()
      ;(wrapper.vm as any).onFilePick(
        mock<MessageEvent>({
          origin: 'https://attacker.example.com',
          data: {
            name: 'owncloud-embed:file-pick',
            data: {
              resource: mock<Resource>({ spaceId: '1' }),
              originRoute: mock<RouteLocation>()
            }
          }
        })
      )
      expect(window.open).not.toHaveBeenCalled()
    })
  })
})

function getWrapper() {
  const mocks = defaultComponentMocks()

  return {
    mocks,
    wrapper: shallowMount(FilePickerModal, {
      props: {
        modal: mock<Modal>(),
        resource: mock<Resource>(),
        app: {
          name: 'text-editor',
          extensions: [{ extension: 'text' }, { extension: 'md' }, { mimeType: 'text/rtf' }]
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
