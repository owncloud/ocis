import {
  defaultComponentMocks,
  flushPromises,
  getComposableWrapper,
  writable
} from '@ownclouders/web-test-helpers'
import { useOpenFolderAction } from '../../../src/composables/useOpenFolderAction'
import { unref } from 'vue'
import { mock } from 'vitest-mock-extended'
import {
  PASSWORD_PROTECTED_FOLDER_RENAMED_MESSAGE,
  Resource,
  SpaceResource
} from '@ownclouders/web-client'
import { useConfigStore, useModals, useResourcesStore } from '@ownclouders/web-pkg'
import { MockedFunction } from 'vitest'
import FolderViewModal from '../../../src/components/FolderViewModal.vue'

const SERVER_URL = 'https://example.org/'
const SERVER_ORIGIN = 'https://example.org'

describe('openFolderAction', () => {
  // the composable registers a `message` listener on `window`, which is shared across tests
  // in the same jsdom environment. Track and remove them so a listener from one test cannot
  // react to a message dispatched by another.
  const trackedListeners: Array<[string, EventListenerOrEventListenerObject]> = []
  const originalAddEventListener = window.addEventListener.bind(window)

  beforeEach(() => {
    trackedListeners.length = 0
    vi.spyOn(window, 'addEventListener').mockImplementation((type, listener) => {
      trackedListeners.push([type, listener])
      originalAddEventListener(type, listener)
    })
  })

  afterEach(() => {
    trackedListeners.forEach(([type, listener]) => window.removeEventListener(type, listener))
    vi.restoreAllMocks()
  })

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

  describe('sync of the .psec file when the folder is renamed inside the modal', () => {
    const psecFile = mock<Resource>({ name: 'old.psec', path: '/old.psec' })
    const renamedMessage = (newName: string, origin = SERVER_ORIGIN) =>
      new MessageEvent('message', {
        origin,
        data: { name: PASSWORD_PROTECTED_FOLDER_RENAMED_MESSAGE, data: { newName } }
      })

    const openModal = (instance: ReturnType<typeof useOpenFolderAction>) =>
      unref(instance).handler({
        resources: [psecFile],
        space: mock<SpaceResource>({ webDavPath: '/files/admin' })
      })

    it('renames the .psec file when the folder is renamed in the framed app', async () => {
      let assertions: () => void
      const { setupPromise } = getWrapper({
        async setup(instance, mocks) {
          await openModal(instance)

          window.dispatchEvent(renamedMessage('new'))
          await flushPromises()

          assertions = () => {
            expect(mocks.$clientService.webdav.moveFiles).toHaveBeenCalledWith(
              expect.anything(),
              psecFile,
              expect.anything(),
              { path: '/new.psec' }
            )
            const { upsertResource } = useResourcesStore()
            expect(upsertResource).toHaveBeenCalledTimes(1)
          }
        }
      })
      await setupPromise
      assertions()
    })

    it('ignores messages from a foreign origin', async () => {
      let assertions: () => void
      const { setupPromise } = getWrapper({
        async setup(instance, mocks) {
          await openModal(instance)

          window.dispatchEvent(renamedMessage('new', 'https://evil.example.com'))
          await flushPromises()

          assertions = () => expect(mocks.$clientService.webdav.moveFiles).not.toHaveBeenCalled()
        }
      })
      await setupPromise
      assertions()
    })

    it('ignores unrelated messages', async () => {
      let assertions: () => void
      const { setupPromise } = getWrapper({
        async setup(instance, mocks) {
          await openModal(instance)

          window.dispatchEvent(
            new MessageEvent('message', {
              origin: SERVER_ORIGIN,
              data: { name: 'some-other-message', data: { newName: 'new' } }
            })
          )
          await flushPromises()

          assertions = () => expect(mocks.$clientService.webdav.moveFiles).not.toHaveBeenCalled()
        }
      })
      await setupPromise
      assertions()
    })

    it('does nothing when the name is unchanged', async () => {
      let assertions: () => void
      const { setupPromise } = getWrapper({
        async setup(instance, mocks) {
          await openModal(instance)

          window.dispatchEvent(renamedMessage('old'))
          await flushPromises()

          assertions = () => expect(mocks.$clientService.webdav.moveFiles).not.toHaveBeenCalled()
        }
      })
      await setupPromise
      assertions()
    })

    it('stops listening once the modal is closed', async () => {
      let assertions: () => void
      const { setupPromise } = getWrapper({
        async setup(instance, mocks) {
          await openModal(instance)

          const { dispatchModal } = useModals()
          const modalConfig = (dispatchModal as MockedFunction<typeof dispatchModal>).mock.calls
            .at(0)
            .at(0)
          modalConfig.onCancel()

          window.dispatchEvent(renamedMessage('new'))
          await flushPromises()

          assertions = () => expect(mocks.$clientService.webdav.moveFiles).not.toHaveBeenCalled()
        }
      })
      await setupPromise
      assertions()
    })
  })
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

  let setupPromise: unknown
  const wrapper = getComposableWrapper(
    () => {
      const configStore = useConfigStore()
      writable(configStore).serverUrl = SERVER_URL
      const instance = useOpenFolderAction()
      setupPromise = setup(instance, mocks)
    },
    {
      mocks,
      provide: mocks
    }
  )

  return { wrapper, setupPromise }
}
