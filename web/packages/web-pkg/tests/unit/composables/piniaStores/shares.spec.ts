import { createTestingPinia, getComposableWrapper } from '@ownclouders/web-test-helpers'
import {
  AddLinkOptions,
  AddShareOptions,
  DeleteLinkOptions,
  DeleteShareOptions,
  UpdateLinkOptions,
  UpdateShareOptions,
  useSharesStore,
  useUserStore
} from '../../../../src/composables/piniaStores'
import { mock, mockDeep } from 'vitest-mock-extended'
import { ClientService } from '../../../../src/services'
import { CollaboratorShare, LinkShare, Resource } from '@ownclouders/web-client'
import { User } from '@ownclouders/web-client/graph/generated'

describe('useSharesStore', () => {
  beforeEach(() => {
    createTestingPinia({
      stubActions: false,
      initialState: { resources: { currentFolder: mock<Resource>() } }
    })
  })

  describe('addShare', () => {
    it('adds a collaborator share', () => {
      getWrapper({
        setup: async (instance) => {
          const resource = { id: '1' } as Resource
          const share = mock<CollaboratorShare>({ id: '1' })
          const user = { id: '1' } as User

          const clientService = mockDeep<ClientService>()
          clientService.graphAuthenticated.permissions.createInvite.mockResolvedValue(share)

          const userStore = useUserStore()
          userStore.user = user

          await instance.addShare(mock<AddShareOptions>({ clientService, resource }))

          expect(clientService.graphAuthenticated.permissions.createInvite).toHaveBeenCalledTimes(1)
          expect(instance.collaboratorShares.length).toBe(1)
        }
      })
    })
  })
  describe('updateShare', () => {
    it('updates a collaborator share', () => {
      getWrapper({
        setup: async (instance) => {
          const resource = { id: '1' } as Resource
          const share = mock<CollaboratorShare>({ id: '1' })
          const user = { id: '1' } as User

          const clientService = mockDeep<ClientService>()
          clientService.graphAuthenticated.permissions.updatePermission.mockResolvedValue(share)

          const userStore = useUserStore()
          userStore.user = user

          await instance.updateShare(mock<UpdateShareOptions>({ clientService, resource }))

          expect(
            clientService.graphAuthenticated.permissions.updatePermission
          ).toHaveBeenCalledTimes(1)
        }
      })
    })
  })
  describe('deleteShare', () => {
    it('deletes a collaborator share', () => {
      getWrapper({
        setup: async (instance) => {
          const resource = { id: '1' } as Resource
          const clientService = mockDeep<ClientService>()
          clientService.graphAuthenticated.permissions.deletePermission.mockResolvedValue(undefined)

          await instance.deleteShare(mock<DeleteShareOptions>({ clientService, resource }))

          expect(
            clientService.graphAuthenticated.permissions.deletePermission
          ).toHaveBeenCalledTimes(1)
        }
      })
    })
  })

  describe('addLink', () => {
    it('adds a link share', () => {
      getWrapper({
        setup: async (instance) => {
          const resource = { id: '1' } as Resource
          const link = mock<LinkShare>({ id: '1' })
          const user = { id: '1' } as User

          const clientService = mockDeep<ClientService>()
          clientService.graphAuthenticated.permissions.createLink.mockResolvedValue(link)

          const userStore = useUserStore()
          userStore.user = user

          await instance.addLink(mock<AddLinkOptions>({ clientService, resource }))

          expect(clientService.graphAuthenticated.permissions.createLink).toHaveBeenCalledTimes(1)
          expect(instance.linkShares.length).toBe(1)
        }
      })
    })
  })
  describe('updateLink', () => {
    it('updates a link share', () => {
      getWrapper({
        setup: async (instance) => {
          const resource = { id: '1' } as Resource
          const link = mock<LinkShare>({ id: '1' })
          const user = { id: '1' } as User

          const clientService = mockDeep<ClientService>()
          clientService.graphAuthenticated.permissions.updatePermission.mockResolvedValue(link)

          const userStore = useUserStore()
          userStore.user = user

          await instance.updateLink(mock<UpdateLinkOptions>({ clientService, resource }))

          expect(
            clientService.graphAuthenticated.permissions.updatePermission
          ).toHaveBeenCalledTimes(1)
        }
      })
    })
  })
  describe('deleteLink', () => {
    it('deletes a link share', () => {
      getWrapper({
        setup: async (instance) => {
          const resource = { id: '1' } as Resource
          const clientService = mockDeep<ClientService>()
          clientService.graphAuthenticated.permissions.deletePermission.mockResolvedValue(undefined)

          await instance.deleteLink(mock<DeleteLinkOptions>({ clientService, resource }))

          expect(
            clientService.graphAuthenticated.permissions.deletePermission
          ).toHaveBeenCalledTimes(1)
        }
      })
    })
  })
})

function getWrapper({ setup }: { setup: (instance: ReturnType<typeof useSharesStore>) => void }) {
  return {
    wrapper: getComposableWrapper(
      () => {
        const instance = useSharesStore()
        setup(instance)
      },
      { pluginOptions: { pinia: false } }
    )
  }
}
