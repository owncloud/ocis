import {
  defaultComponentMocks,
  getComposableWrapper,
  writable
} from '@ownclouders/web-test-helpers'
import { useOpenFolderAction } from '../../../src/composables/useOpenFolderAction'
import { unref } from 'vue'
import { mock } from 'vitest-mock-extended'
import { Resource, SpaceResource } from '@ownclouders/web-client'
import { useConfigStore, useModals } from '@ownclouders/web-pkg'
import { MockedFunction } from 'vitest'
import FolderViewModal from '../../../src/components/FolderViewModal.vue'

const SERVER_URL = 'https://example.org/'

describe('openFolderAction', () => {
  it('should open a modal with the public link', () => {
    getWrapper({
      async setup(instance) {
        const { dispatchModal } = useModals()

        await unref(instance).handler({
          resources: [mock<Resource>()],
          space: mock<SpaceResource>()
        })

        const modalConfig = (dispatchModal as MockedFunction<typeof dispatchModal>).mock.calls
          .at(0)
          .at(0)
        const attrs = modalConfig.customComponentAttrs()

        expect(dispatchModal).toHaveBeenCalledWith(
          expect.objectContaining({ customComponent: FolderViewModal })
        )
        expect(attrs).toStrictEqual({
          publicLink: 'https://example.org/public-link',
          serverUrl: SERVER_URL
        })
      }
    })
  })

  it('should throw when .psec file URL points to a different server', () => {
    getWrapper({
      body: btoa('https://other.example.com/public-link'),
      async setup(instance) {
        await expect(
          unref(instance).handler({
            resources: [mock<Resource>()],
            space: mock<SpaceResource>()
          })
        ).rejects.toThrow(
          'This folder cannot be opened because the link it contains does not point to this server.'
        )
      }
    })
  })

  it.each(['javascript:alert(1)', 'data:text/html,<script>alert(1)</script>', 'blob:fake'])(
    'should throw when .psec file contains a non-http(s) URL: %s',
    (invalidUrl) => {
      getWrapper({
        body: btoa(invalidUrl),
        async setup(instance) {
          await expect(
            unref(instance).handler({
              resources: [mock<Resource>()],
              space: mock<SpaceResource>()
            })
          ).rejects.toThrow('This folder cannot be opened because the link it contains is invalid.')
        }
      })
    }
  )
})

function getWrapper({
  setup,
  body = btoa('https://example.org/public-link')
}: {
  setup: (
    instance: ReturnType<typeof useOpenFolderAction>,
    mocks: ReturnType<typeof defaultComponentMocks>
  ) => void
  body?: string
}) {
  const mocks = defaultComponentMocks()
  mocks.$clientService.webdav.getFileContents.mockResolvedValue({ body })

  return {
    wrapper: getComposableWrapper(
      () => {
        const configStore = useConfigStore()
        writable(configStore).serverUrl = SERVER_URL
        const instance = useOpenFolderAction()
        setup(instance, mocks)
      },
      {
        mocks,
        provide: mocks
      }
    )
  }
}
