import { Resource, SpaceResource } from '@ownclouders/web-client'
import { createTestingPinia } from '@ownclouders/web-test-helpers'
import {
  ClientService,
  PreviewService,
  useConfigStore,
  useMessages,
  useResourcesStore,
  useSharesStore,
  useSpacesStore,
  useUserStore
} from '@ownclouders/web-pkg'
import { mock, mockDeep } from 'vitest-mock-extended'
import { isItemInCurrentFolder, sseEventWrapper } from '../../../../src/container/sse'
import PQueue from 'p-queue'
import { Language } from 'vue3-gettext'
import { Router } from 'vue-router'

describe('helpers', () => {
  describe('method "sseEventWrapper"', () => {
    it('calls "console.debug" when executed', () => {
      console.debug = vi.fn()
      const topic = 'folder-created'
      const msg = mock<MessageEvent>({ data: JSON.stringify({ itemid: 'newfolder' }) })
      sseEventWrapper({
        msg,
        topic,
        method: () => {},
        ...getMocks()
      })
      expect(console.debug).toHaveBeenCalledWith(`SSE event '${topic}'`, { itemid: 'newfolder' })
    })
    it('calls "console.error" when error was thrown', () => {
      console.error = vi.fn()
      const topic = 'folder-created'
      const msg = mock<MessageEvent>({ data: JSON.stringify({ itemid: 'newfolder' }) })
      const error = new Error('processing failed')
      sseEventWrapper({
        msg,
        topic,
        method: () => {
          throw error
        },
        ...getMocks()
      })
      expect(console.error).toHaveBeenCalledWith(`Unable to process sse event ${topic}`, error)
    })
  })
  describe('method "isItemInCurrentFolder"', () => {
    it('returns "true" when item is in current folder', () => {
      const mocks = getMocks()
      expect(
        isItemInCurrentFolder({
          resourcesStore: mocks.resourcesStore,
          parentFolderId: 'currenFolder!currentFolder'
        })
      ).toBeTruthy()
    })
    it('returns "false" when item is not in current folder', () => {
      const mocks = getMocks()
      expect(
        isItemInCurrentFolder({
          resourcesStore: mocks.resourcesStore,
          parentFolderId: 'differentFolder!differentFolder'
        })
      ).toBeFalsy()
    })
    describe('current folder is space', () => {
      it('returns "true" when item is in current folder', () => {
        const mocks = getMocks({
          currentFolder: mock<SpaceResource>({
            id: 'bbf8b91f-54be-45f0-935e-a50c4922db21$c96eb07d-54a5-47bf-8402-64ad9a007797'
          })
        })
        expect(
          isItemInCurrentFolder({
            resourcesStore: mocks.resourcesStore,
            parentFolderId:
              'bbf8b91f-54be-45f0-935e-a50c4922db21$c96eb07d-54a5-47bf-8402-64ad9a007797!c96eb07d-54a5-47bf-8402-64ad9a007797'
          })
        ).toBeTruthy()
      })
      it('returns "false" when item is not in current folder', () => {
        const mocks = getMocks({
          currentFolder: mock<SpaceResource>({
            id: 'bbf8b91f-54be-45f0-935e-a50c4922db21$c96eb07d-54a5-47bf-8402-64ad9a007797'
          })
        })
        expect(
          isItemInCurrentFolder({
            resourcesStore: mocks.resourcesStore,
            parentFolderId: 'differentFolder!differentFolder'
          })
        ).toBeFalsy()
      })
    })
  })
})

const getMocks = ({
  currentFolder = mock<Resource>({
    id: 'currenFolder!currentFolder',
    isFolder: true,
    storageId: 'space1'
  })
}: { currentFolder?: Resource } = {}) => {
  createTestingPinia()
  const resourcesStore = useResourcesStore()
  resourcesStore.currentFolder = currentFolder
  const spacesStore = useSpacesStore()
  const messageStore = useMessages()
  const userStore = useUserStore()
  const configStore = useConfigStore()
  const sharesStore = useSharesStore()
  const clientService = mockDeep<ClientService>()
  const previewService = mockDeep<PreviewService>()
  const router = mockDeep<Router>()
  const language = mockDeep<Language>()
  const resourceQueue = mockDeep<PQueue>()

  return {
    resourcesStore,
    spacesStore,
    router,
    messageStore,
    userStore,
    sharesStore,
    configStore,
    clientService,
    previewService,
    resourceQueue,
    language
  }
}
