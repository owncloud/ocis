import {
  ClientService,
  eventBus,
  PreviewService,
  useConfigStore,
  useMessages,
  useResourcesStore,
  useSharesStore,
  useSpacesStore,
  useUserStore
} from '@ownclouders/web-pkg'
import {
  EventSchemaType,
  onSSELinkCreatedEvent,
  onSSELinkRemovedEvent,
  onSSEShareCreatedEvent,
  onSSEShareRemovedEvent,
  onSSEShareUpdatedEvent,
  onSSESpaceMemberAddedEvent,
  onSSESpaceMemberRemovedEvent,
  onSSESpaceShareUpdatedEvent
} from '../../../../src/container/sse'
import { mock, mockDeep } from 'vitest-mock-extended'
import { DriveItem, User } from '@ownclouders/web-client/graph/generated'
import { ShareTypes, Resource, SpaceResource } from '@ownclouders/web-client'
import { createTestingPinia, defaultComponentMocks } from '@ownclouders/web-test-helpers'
import { Language } from 'vue3-gettext'
import PQueue from 'p-queue'
import { RouteLocation } from 'vue-router'

describe('shares events', () => {
  describe('onSSESpaceMemberAddedEvent', () => {
    it('calls "upsertSpace" when space member has been added', async () => {
      const space = mock<SpaceResource>({ id: 'space1' })
      const mocks = getMocks({ currentRouteFilesSpacesGeneric: true })
      const sseData = mock<EventSchemaType>({
        itemid: space.id,
        spaceid: space.id
      })
      mocks.clientService.graphAuthenticated.drives.getDrive.mockResolvedValue(space)
      await onSSESpaceMemberAddedEvent({ sseData, ...mocks })
      expect(mocks.clientService.graphAuthenticated.drives.getDrive).toHaveBeenCalled()
      expect(mocks.spacesStore.upsertSpace).toHaveBeenCalled()
      expect(mocks.resourcesStore.upsertResource).not.toHaveBeenCalled()
    })
    it('calls "upsertResource" when space member has been added and current route equals "files-spaces-projects"', async () => {
      const space = mock<SpaceResource>({ id: 'space1' })
      const mocks = getMocks({ currentRouteFilesSpacesProjects: true })
      const sseData = mock<EventSchemaType>({
        itemid: space.id,
        spaceid: space.id
      })
      mocks.clientService.graphAuthenticated.drives.getDrive.mockResolvedValue(space)
      await onSSESpaceMemberAddedEvent({ sseData, ...mocks })
      expect(mocks.clientService.graphAuthenticated.drives.getDrive).toHaveBeenCalled()
      expect(mocks.spacesStore.upsertSpace).toHaveBeenCalled()
      expect(mocks.resourcesStore.upsertResource).toHaveBeenCalled()
    })
    it('does not trigger any action when initiator ids are identical', async () => {
      const space = mock<SpaceResource>({ id: 'space1' })
      const mocks = getMocks({ currentRouteFilesSpacesProjects: true })
      const sseData = mock<EventSchemaType>({
        itemid: space.id,
        spaceid: space.id,
        initiatorid: 'local1'
      })
      mocks.clientService.graphAuthenticated.drives.getDrive.mockResolvedValue(space)
      await onSSESpaceMemberAddedEvent({ sseData, ...mocks })
      expect(mocks.clientService.graphAuthenticated.drives.getDrive).not.toHaveBeenCalled()
      expect(mocks.spacesStore.upsertSpace).not.toHaveBeenCalled()
      expect(mocks.resourcesStore.upsertResource).not.toHaveBeenCalled()
    })
  })
  describe('onSSESpaceMemberRemovedEvent', () => {
    it('calls "upsertSpace" when space member has been removed and current user is not affected', async () => {
      const space = mock<SpaceResource>({ id: 'space1' })
      const mocks = getMocks()
      const sseData = mock<EventSchemaType>({
        itemid: space.id,
        spaceid: space.id,
        affecteduserids: ['2']
      })
      mocks.clientService.graphAuthenticated.drives.getDrive.mockResolvedValue(space)
      await onSSESpaceMemberRemovedEvent({ sseData, ...mocks })
      expect(mocks.clientService.graphAuthenticated.drives.getDrive).toHaveBeenCalled()
      expect(mocks.spacesStore.upsertSpace).toHaveBeenCalled()
      expect(mocks.spacesStore.removeSpace).not.toHaveBeenCalled()
      expect(mocks.messageStore.showMessage).not.toHaveBeenCalled()
    })
    it('calls "removeSpace" when space member has been removed and current user is affected', async () => {
      const space = mock<SpaceResource>({ id: 'space1' })
      const mocks = getMocks({ currentRouteFilesSpacesProjects: true, spaces: [space] })
      const sseData = mock<EventSchemaType>({
        itemid: space.id,
        spaceid: space.id,
        affecteduserids: ['1']
      })
      mocks.clientService.graphAuthenticated.drives.getDrive.mockResolvedValue(space)
      await onSSESpaceMemberRemovedEvent({ sseData, ...mocks })
      expect(mocks.clientService.graphAuthenticated.drives.getDrive).not.toHaveBeenCalled()
      expect(mocks.spacesStore.upsertSpace).not.toHaveBeenCalled()
      expect(mocks.spacesStore.removeSpace).toHaveBeenCalled()
      expect(mocks.messageStore.showMessage).not.toHaveBeenCalled()
    })
    it('calls "showMessage" when space member has been removed and current user is affected and navigated to space', async () => {
      const space = mock<SpaceResource>({ id: 'space1' })
      const mocks = getMocks({ currentRouteFilesSpacesGeneric: true, spaces: [space] })
      const sseData = mock<EventSchemaType>({
        itemid: space.id,
        spaceid: space.id,
        affecteduserids: ['1']
      })
      mocks.clientService.graphAuthenticated.drives.getDrive.mockResolvedValue(space)
      await onSSESpaceMemberRemovedEvent({ sseData, ...mocks })
      expect(mocks.clientService.graphAuthenticated.drives.getDrive).not.toHaveBeenCalled()
      expect(mocks.spacesStore.upsertSpace).not.toHaveBeenCalled()
      expect(mocks.spacesStore.removeSpace).toHaveBeenCalled()
      expect(mocks.messageStore.showMessage).toHaveBeenCalled()
    })
    it('does not trigger any action when initiator ids are identical', async () => {
      const space = mock<SpaceResource>({ id: 'space1' })
      const mocks = getMocks({ currentRouteFilesSpacesProjects: true, spaces: [space] })
      const sseData = mock<EventSchemaType>({
        itemid: space.id,
        spaceid: space.id,
        affecteduserids: ['1'],
        initiatorid: 'local1'
      })
      mocks.clientService.graphAuthenticated.drives.getDrive.mockResolvedValue(space)
      await onSSESpaceMemberRemovedEvent({ sseData, ...mocks })
      expect(mocks.clientService.graphAuthenticated.drives.getDrive).not.toHaveBeenCalled()
      expect(mocks.spacesStore.upsertSpace).not.toHaveBeenCalled()
      expect(mocks.spacesStore.removeSpace).not.toHaveBeenCalled()
      expect(mocks.messageStore.showMessage).not.toHaveBeenCalled()
    })
  })
  describe('onSSESpaceShareUpdatedEvent', () => {
    it('calls "upsertSpace" when space share has been updated', async () => {
      const space = mock<SpaceResource>({ id: 'space1' })
      const mocks = getMocks()
      const sseData = mock<EventSchemaType>({
        itemid: space.id,
        spaceid: space.id,
        affecteduserids: ['2']
      })
      mocks.clientService.graphAuthenticated.drives.getDrive.mockResolvedValue(space)
      const busStub = vi.spyOn(eventBus, 'publish')
      await onSSESpaceShareUpdatedEvent({ sseData, ...mocks })
      expect(mocks.clientService.graphAuthenticated.drives.getDrive).toHaveBeenCalled()
      expect(mocks.spacesStore.upsertSpace).toHaveBeenCalled()
      expect(busStub).not.toHaveBeenCalled()
    })
    it('calls "eventBus.publish" when space share has been updated and current user is affected and navigated to space', async () => {
      const space = mock<SpaceResource>({ id: 'space1' })
      const mocks = getMocks({ currentRouteFilesSpacesGeneric: true, spaces: [space] })
      const sseData = mock<EventSchemaType>({
        itemid: space.id,
        spaceid: space.id,
        affecteduserids: ['1']
      })
      mocks.clientService.graphAuthenticated.drives.getDrive.mockResolvedValue(space)
      const busStub = vi.spyOn(eventBus, 'publish')
      await onSSESpaceShareUpdatedEvent({ sseData, ...mocks })
      expect(mocks.clientService.graphAuthenticated.drives.getDrive).toHaveBeenCalled()
      expect(mocks.spacesStore.upsertSpace).toHaveBeenCalled()
      expect(busStub).toHaveBeenCalled()
    })
    it('does not trigger any action when initiator ids are identical', async () => {
      const space = mock<SpaceResource>({ id: 'space1' })
      const mocks = getMocks({ currentRouteFilesSpacesGeneric: true, spaces: [space] })
      const sseData = mock<EventSchemaType>({
        itemid: space.id,
        spaceid: space.id,
        affecteduserids: ['1'],
        initiatorid: 'local1'
      })
      mocks.clientService.graphAuthenticated.drives.getDrive.mockResolvedValue(space)
      const busStub = vi.spyOn(eventBus, 'publish')
      await onSSESpaceShareUpdatedEvent({ sseData, ...mocks })
      expect(mocks.clientService.graphAuthenticated.drives.getDrive).not.toHaveBeenCalled()
      expect(mocks.spacesStore.upsertSpace).not.toHaveBeenCalled()
      expect(busStub).not.toHaveBeenCalled()
    })
  })
  describe('onSSEShareCreatedEvent', () => {
    it('calls "upsertResource" when resource has been shared and current user navigated to resource', async () => {
      const sharedResource = mock<Resource>({
        id: 'sharedResource',
        storageId: 'space1',
        parentFolderId: 'currenFolder!currentFolder'
      })
      const space = mock<SpaceResource>({ id: sharedResource.storageId })
      const mocks = getMocks({ currentRouteFilesSpacesGeneric: true, spaces: [space] })
      const sseData = mock<EventSchemaType>({
        itemid: sharedResource.id,
        spaceid: sharedResource.storageId,
        parentitemid: sharedResource.parentFolderId
      })
      mocks.clientService.webdav.getFileInfo.mockResolvedValue(
        mock<Resource>({ ...sharedResource })
      )
      await onSSEShareCreatedEvent({ sseData, ...mocks })
      expect(mocks.clientService.webdav.getFileInfo).toHaveBeenCalled()
      expect(mocks.resourcesStore.upsertResource).toHaveBeenCalled()
      expect(
        mocks.clientService.graphAuthenticated.driveItems.listSharedWithMe
      ).not.toHaveBeenCalled()
      expect(
        mocks.clientService.graphAuthenticated.driveItems.listSharedByMe
      ).not.toHaveBeenCalled()
    })

    it('calls "upsertResource" when resource has been shared and current route equals "files-shares-with-me"', async () => {
      const sharedDrive = mockDeep<DriveItem>({
        id: 'sharedDrive',
        remoteItem: {
          id: 'sharedResource',
          permissions: []
        }
      })
      const mocks = getMocks({
        currentRouteFilesSharesWithMe: true,
        currentFolder: null
      })
      const sseData = mock<EventSchemaType>({
        itemid: sharedDrive.remoteItem.id
      })
      mocks.clientService.graphAuthenticated.driveItems.listSharedWithMe.mockResolvedValue([
        sharedDrive
      ])
      await onSSEShareCreatedEvent({ sseData, ...mocks })
      expect(mocks.clientService.graphAuthenticated.driveItems.listSharedWithMe).toHaveBeenCalled()
      expect(mocks.resourcesStore.upsertResource).toHaveBeenCalled()
      expect(mocks.resourcesStore.updateResourceField).not.toHaveBeenCalled()
      expect(mocks.clientService.webdav.getFileInfo).not.toHaveBeenCalled()
      expect(
        mocks.clientService.graphAuthenticated.driveItems.listSharedByMe
      ).not.toHaveBeenCalled()
    })
    it('calls "upsertResource" when resource has been shared and current route equals "files-shares-with-others"', async () => {
      const sharedDrive = mockDeep<DriveItem>({
        id: 'sharedDrive',
        permissions: []
      })
      const mocks = getMocks({
        currentRouteFilesSharesWithOthers: true,
        currentFolder: null
      })
      const sseData = mock<EventSchemaType>({
        itemid: sharedDrive.id
      })
      mocks.clientService.graphAuthenticated.driveItems.listSharedByMe.mockResolvedValue([
        sharedDrive
      ])
      await onSSEShareCreatedEvent({ sseData, ...mocks })
      expect(mocks.clientService.graphAuthenticated.driveItems.listSharedByMe).toHaveBeenCalled()
      expect(mocks.resourcesStore.upsertResource).toHaveBeenCalled()
      expect(mocks.resourcesStore.updateResourceField).not.toHaveBeenCalled()
      expect(mocks.clientService.webdav.getFileInfo).not.toHaveBeenCalled()
      expect(
        mocks.clientService.graphAuthenticated.driveItems.listSharedWithMe
      ).not.toHaveBeenCalled()
    })
    it('does not trigger any action when initiator ids are identical', async () => {
      const sharedResource = mock<Resource>({
        id: 'sharedResource',
        storageId: 'space1',
        parentFolderId: 'currenFolder!currentFolder'
      })
      const space = mock<SpaceResource>({ id: sharedResource.storageId })
      const mocks = getMocks({ currentRouteFilesSpacesGeneric: true, spaces: [space] })
      const sseData = mock<EventSchemaType>({
        itemid: sharedResource.id,
        spaceid: sharedResource.storageId,
        parentitemid: sharedResource.parentFolderId,
        initiatorid: 'local1'
      })
      await onSSEShareCreatedEvent({ sseData, ...mocks })
      expect(mocks.clientService.webdav.getFileInfo).not.toHaveBeenCalled()
      expect(mocks.resourcesStore.upsertResource).not.toHaveBeenCalled()
      expect(mocks.resourcesStore.updateResourceField).not.toHaveBeenCalled()
      expect(
        mocks.clientService.graphAuthenticated.driveItems.listSharedWithMe
      ).not.toHaveBeenCalled()
      expect(
        mocks.clientService.graphAuthenticated.driveItems.listSharedByMe
      ).not.toHaveBeenCalled()
    })
  })
  describe('onSSEShareUpdatedEvent', () => {
    it('calls "eventBus.publish" when share has been updated and current user is affected and navigated to share', async () => {
      const sharedResource = mock<Resource>({
        id: 'sharedResource',
        storageId: 'space1',
        parentFolderId: 'space1'
      })
      const mocks = getMocks({ currentRouteFilesSpacesGeneric: true })
      const sseData = mock<EventSchemaType>({
        itemid: sharedResource.id,
        spaceid: sharedResource.storageId,
        affecteduserids: ['1']
      })
      const busStub = vi.spyOn(eventBus, 'publish')
      await onSSEShareUpdatedEvent({ sseData, ...mocks })
      expect(busStub).toHaveBeenCalled()
      expect(
        mocks.clientService.graphAuthenticated.driveItems.listSharedWithMe
      ).not.toHaveBeenCalled()
      expect(mocks.resourcesStore.upsertResource).not.toHaveBeenCalled()
    })
    it('calls "upsertResource" when share has been updated and current route equals "files-shares-with-me"', async () => {
      const sharedDrive = mockDeep<DriveItem>({
        id: 'sharedDrive',
        remoteItem: {
          id: 'sharedResource',
          permissions: []
        }
      })
      const mocks = getMocks({
        currentRouteFilesSharesWithMe: true,
        currentFolder: null
      })
      const sseData = mock<EventSchemaType>({
        itemid: sharedDrive.remoteItem.id
      })
      mocks.clientService.graphAuthenticated.driveItems.listSharedWithMe.mockResolvedValue([
        sharedDrive
      ])
      const busStub = vi.spyOn(eventBus, 'publish')
      await onSSEShareUpdatedEvent({ sseData, ...mocks })
      expect(mocks.clientService.graphAuthenticated.driveItems.listSharedWithMe).toHaveBeenCalled()
      expect(mocks.resourcesStore.upsertResource).toHaveBeenCalled()
      expect(busStub).not.toHaveBeenCalled()
    })
    it('does not trigger any action when initiator ids are identical', async () => {
      const sharedResource = mock<Resource>({
        id: 'sharedResource',
        storageId: 'space1',
        parentFolderId: 'space1'
      })
      const mocks = getMocks({ currentRouteFilesSpacesGeneric: true })
      const sseData = mock<EventSchemaType>({
        itemid: sharedResource.id,
        spaceid: sharedResource.storageId,
        affecteduserids: ['1'],
        initiatorid: 'local1'
      })
      const busStub = vi.spyOn(eventBus, 'publish')
      await onSSEShareUpdatedEvent({ sseData, ...mocks })
      expect(busStub).not.toHaveBeenCalled()
      expect(
        mocks.clientService.graphAuthenticated.driveItems.listSharedWithMe
      ).not.toHaveBeenCalled()
      expect(mocks.resourcesStore.upsertResource).not.toHaveBeenCalled()
    })
  })
  describe('onSSEShareRemovedEvent', () => {
    it('calls "showMessage" when share has been removed and current user is affected and navigated to share', async () => {
      const mocks = getMocks({
        currentRouteFilesSpacesGeneric: true
      })
      const sseData = mock<EventSchemaType>({
        itemid: 'sharedResource',
        spaceid: 'space1',
        affecteduserids: ['1']
      })
      await onSSEShareRemovedEvent({ sseData, ...mocks })

      expect(mocks.messageStore.showMessage).toHaveBeenCalled()
      expect(mocks.clientService.webdav.getFileInfo).not.toHaveBeenCalled()
      expect(mocks.resourcesStore.upsertResource).not.toHaveBeenCalled()
      expect(mocks.resourcesStore.updateResourceField).not.toHaveBeenCalled()
      expect(mocks.resourcesStore.removeResources).not.toHaveBeenCalled()
    })

    it('calls "upsertResource" when share has been removed and current user is not affected and share is located in current folder', async () => {
      const sharedResource = mock<Resource>({
        id: 'sharedResource',
        storageId: 'space1',
        parentFolderId: 'currenFolder!currentFolder'
      })
      const space = mock<SpaceResource>({ id: sharedResource.storageId })
      const mocks = getMocks({ currentRouteFilesSpacesGeneric: true, spaces: [space] })
      const sseData = mock<EventSchemaType>({
        itemid: sharedResource.id,
        spaceid: sharedResource.storageId,
        parentitemid: sharedResource.parentFolderId,
        affecteduserids: ['2']
      })
      mocks.clientService.webdav.getFileInfo.mockResolvedValue(
        mock<Resource>({ ...sharedResource })
      )
      await onSSEShareRemovedEvent({ sseData, ...mocks })
      expect(mocks.clientService.webdav.getFileInfo).toHaveBeenCalled()
      expect(mocks.resourcesStore.upsertResource).toHaveBeenCalled()
      expect(mocks.messageStore.showMessage).not.toHaveBeenCalled()
      expect(mocks.resourcesStore.removeResources).not.toHaveBeenCalled()
    })
    it('calls "removeResources" when share has removed and current route equals "files-shared-with-others"', async () => {
      const sharedResource = mock<Resource>({
        id: 'sharedResource',
        storageId: 'space1',
        shareTypes: []
      })
      const space = mock<SpaceResource>({ id: sharedResource.storageId })
      const mocks = getMocks({
        currentRouteFilesSharesWithOthers: true,
        currentFolder: null,
        spaces: [space]
      })
      const sseData = mock<EventSchemaType>({
        itemid: sharedResource.id,
        spaceid: sharedResource.storageId,
        parentitemid: sharedResource.parentFolderId
      })
      mocks.clientService.webdav.getFileInfo.mockResolvedValue(
        mock<Resource>({ ...sharedResource })
      )
      await onSSEShareRemovedEvent({ sseData, ...mocks })
      expect(mocks.resourcesStore.removeResources).toHaveBeenCalled()
      expect(mocks.clientService.webdav.getFileInfo).toHaveBeenCalled()
      expect(mocks.resourcesStore.upsertResource).not.toHaveBeenCalled()
      expect(mocks.resourcesStore.updateResourceField).not.toHaveBeenCalled()
      expect(mocks.messageStore.showMessage).not.toHaveBeenCalled()
    })
    it('does not call "removeResources" when share has been removed and current route equals "files-shares-with-others but link shares are present"', async () => {
      const sharedResource = mock<Resource>({
        id: 'sharedResource',
        storageId: 'space1',
        shareTypes: [ShareTypes.group.value]
      })
      const space = mock<SpaceResource>({ id: sharedResource.storageId })
      const mocks = getMocks({
        currentRouteFilesSharesWithOthers: true,
        currentFolder: null,
        spaces: [space]
      })
      const sseData = mock<EventSchemaType>({
        itemid: sharedResource.id,
        spaceid: sharedResource.storageId,
        parentitemid: sharedResource.parentFolderId
      })
      mocks.clientService.webdav.getFileInfo.mockResolvedValue(
        mock<Resource>({ ...sharedResource })
      )
      await onSSEShareRemovedEvent({ sseData, ...mocks })
      expect(mocks.clientService.webdav.getFileInfo).toHaveBeenCalled()
      expect(mocks.resourcesStore.removeResources).not.toHaveBeenCalled()
      expect(mocks.resourcesStore.upsertResource).not.toHaveBeenCalled()
      expect(mocks.resourcesStore.updateResourceField).not.toHaveBeenCalled()
      expect(mocks.messageStore.showMessage).not.toHaveBeenCalled()
    })
    it('calls "removeResources" when share has removed and current route equals "files-shared-with-me"', async () => {
      const sharedResource = mock<Resource>({
        id: 'sharedResource',
        fileId: 'sharedResource',
        storageId: 'space1'
      })
      const mocks = getMocks({
        currentRouteFilesSharesWithMe: true,
        currentFolder: null,
        resources: [sharedResource]
      })
      const sseData = mock<EventSchemaType>({
        itemid: sharedResource.id,
        spaceid: sharedResource.storageId
      })
      await onSSEShareRemovedEvent({ sseData, ...mocks })
      expect(mocks.resourcesStore.removeResources).toHaveBeenCalled()
      expect(mocks.clientService.webdav.getFileInfo).not.toHaveBeenCalled()
      expect(mocks.resourcesStore.upsertResource).not.toHaveBeenCalled()
      expect(mocks.resourcesStore.updateResourceField).not.toHaveBeenCalled()
      expect(mocks.messageStore.showMessage).not.toHaveBeenCalled()
    })
    it('does not trigger any action when initiator ids are identical', async () => {
      const mocks = getMocks({
        currentRouteFilesSpacesGeneric: true
      })
      const sseData = mock<EventSchemaType>({
        itemid: 'sharedResource',
        spaceid: 'space1',
        affecteduserids: ['1'],
        initiatorid: 'local1'
      })
      await onSSEShareRemovedEvent({ sseData, ...mocks })
      expect(mocks.messageStore.showMessage).not.toHaveBeenCalled()
      expect(mocks.clientService.webdav.getFileInfo).not.toHaveBeenCalled()
      expect(mocks.resourcesStore.upsertResource).not.toHaveBeenCalled()
      expect(mocks.resourcesStore.updateResourceField).not.toHaveBeenCalled()
      expect(mocks.resourcesStore.removeResources).not.toHaveBeenCalled()
    })
  })
  describe('onSSELinkCreatedEvent', () => {
    it('calls "upsertResource" when resource has been shared via link and current user navigated to resource', async () => {
      const sharedResource = mock<Resource>({
        id: 'sharedResource',
        storageId: 'space1',
        parentFolderId: 'currenFolder!currentFolder'
      })
      const space = mock<SpaceResource>({ id: sharedResource.storageId })
      const mocks = getMocks({ currentRouteFilesSpacesGeneric: true, spaces: [space] })
      const sseData = mock<EventSchemaType>({
        itemid: sharedResource.id,
        spaceid: sharedResource.storageId,
        parentitemid: sharedResource.parentFolderId
      })
      mocks.clientService.webdav.getFileInfo.mockResolvedValue(
        mock<Resource>({ ...sharedResource })
      )
      await onSSELinkCreatedEvent({ sseData, ...mocks })
      expect(mocks.clientService.webdav.getFileInfo).toHaveBeenCalled()
      expect(mocks.resourcesStore.upsertResource).toHaveBeenCalled()
      expect(
        mocks.clientService.graphAuthenticated.driveItems.listSharedByMe
      ).not.toHaveBeenCalled()
    })
    it('calls "upsertResource" when resource has been shared via link and current route equals "files-shares-via-link"', async () => {
      const sharedDrive = mockDeep<DriveItem>({
        id: 'sharedDrive',
        permissions: []
      })
      const mocks = getMocks({
        currentRouteFilesSharesViaLink: true,
        currentFolder: null
      })
      const sseData = mock<EventSchemaType>({
        itemid: sharedDrive.id
      })
      mocks.clientService.graphAuthenticated.driveItems.listSharedByMe.mockResolvedValue([
        sharedDrive
      ])
      await onSSELinkCreatedEvent({ sseData, ...mocks })
      expect(mocks.clientService.graphAuthenticated.driveItems.listSharedByMe).toHaveBeenCalled()
      expect(mocks.resourcesStore.upsertResource).toHaveBeenCalled()
      expect(mocks.resourcesStore.updateResourceField).not.toHaveBeenCalled()
      expect(mocks.clientService.webdav.getFileInfo).not.toHaveBeenCalled()
    })
    it('does not trigger any action when initiator ids are identical', async () => {
      const sharedResource = mock<Resource>({
        id: 'sharedResource',
        storageId: 'space1',
        parentFolderId: 'currenFolder!currentFolder'
      })
      const space = mock<SpaceResource>({ id: sharedResource.storageId })
      const mocks = getMocks({ currentRouteFilesSpacesGeneric: true, spaces: [space] })
      const sseData = mock<EventSchemaType>({
        itemid: sharedResource.id,
        spaceid: sharedResource.storageId,
        parentitemid: sharedResource.parentFolderId,
        initiatorid: 'local1'
      })
      mocks.clientService.webdav.getFileInfo.mockResolvedValue(
        mock<Resource>({ ...sharedResource })
      )
      await onSSELinkCreatedEvent({ sseData, ...mocks })
      expect(mocks.clientService.webdav.getFileInfo).not.toHaveBeenCalled()
      expect(mocks.resourcesStore.upsertResource).not.toHaveBeenCalled()
      expect(mocks.resourcesStore.updateResourceField).not.toHaveBeenCalled()
      expect(
        mocks.clientService.graphAuthenticated.driveItems.listSharedByMe
      ).not.toHaveBeenCalled()
    })
  })
  describe('onSSELinkRemovedEvent', () => {
    it('calls "upsertResource" when link share has been removed and resource is in current folder', async () => {
      const sharedResource = mock<Resource>({
        id: 'sharedResource',
        storageId: 'space1',
        parentFolderId: 'currenFolder!currentFolder'
      })
      const space = mock<SpaceResource>({ id: sharedResource.storageId })
      const mocks = getMocks({ currentRouteFilesSpacesGeneric: true, spaces: [space] })
      const sseData = mock<EventSchemaType>({
        itemid: sharedResource.id,
        spaceid: sharedResource.storageId,
        parentitemid: sharedResource.parentFolderId
      })
      mocks.clientService.webdav.getFileInfo.mockResolvedValue(
        mock<Resource>({ ...sharedResource })
      )
      await onSSELinkRemovedEvent({ sseData, ...mocks })
      expect(mocks.clientService.webdav.getFileInfo).toHaveBeenCalled()
      expect(mocks.resourcesStore.upsertResource).toHaveBeenCalled()
      expect(mocks.resourcesStore.removeResources).not.toHaveBeenCalled()
    })
    it('calls "removeResources" when link share has been removed and current route equals "files-shares-via-link"', async () => {
      const sharedResource = mock<Resource>({
        id: 'sharedResource',
        storageId: 'space1',
        parentFolderId: 'currenFolder!currentFolder',
        shareTypes: []
      })
      const space = mock<SpaceResource>({ id: sharedResource.storageId })
      const mocks = getMocks({
        currentRouteFilesSharesViaLink: true,
        currentFolder: null,
        spaces: [space]
      })
      const sseData = mock<EventSchemaType>({
        itemid: sharedResource.id,
        spaceid: sharedResource.storageId,
        parentitemid: sharedResource.parentFolderId
      })
      mocks.clientService.webdav.getFileInfo.mockResolvedValue(
        mock<Resource>({ ...sharedResource })
      )
      await onSSELinkRemovedEvent({ sseData, ...mocks })
      expect(mocks.resourcesStore.removeResources).toHaveBeenCalled()
      expect(mocks.clientService.webdav.getFileInfo).toHaveBeenCalled()
      expect(mocks.resourcesStore.upsertResource).not.toHaveBeenCalled()
      expect(mocks.resourcesStore.updateResourceField).not.toHaveBeenCalled()
    })
    it('does not call "removeResources" when link share has been removed and current route equals "files-shares-via-link but other link shares are present"', async () => {
      const sharedResource = mock<Resource>({
        id: 'sharedResource',
        storageId: 'space1',
        parentFolderId: 'currenFolder!currentFolder',
        shareTypes: [ShareTypes.link.value]
      })
      const space = mock<SpaceResource>({ id: sharedResource.storageId })
      const mocks = getMocks({
        currentRouteFilesSharesViaLink: true,
        currentFolder: null,
        spaces: [space]
      })
      const sseData = mock<EventSchemaType>({
        itemid: sharedResource.id,
        spaceid: sharedResource.storageId,
        parentitemid: sharedResource.parentFolderId
      })
      mocks.clientService.webdav.getFileInfo.mockResolvedValue(
        mock<Resource>({ ...sharedResource })
      )
      await onSSELinkRemovedEvent({ sseData, ...mocks })
      expect(mocks.resourcesStore.removeResources).not.toHaveBeenCalled()
      expect(mocks.clientService.webdav.getFileInfo).toHaveBeenCalled()
      expect(mocks.resourcesStore.upsertResource).not.toHaveBeenCalled()
      expect(mocks.resourcesStore.updateResourceField).not.toHaveBeenCalled()
    })
    it('does not trigger any action when initiator ids are identical', async () => {
      const sharedResource = mock<Resource>({
        id: 'sharedResource',
        storageId: 'space1',
        parentFolderId: 'currenFolder!currentFolder'
      })
      const space = mock<SpaceResource>({ id: sharedResource.storageId })
      const mocks = getMocks({ currentRouteFilesSpacesGeneric: true, spaces: [space] })
      const sseData = mock<EventSchemaType>({
        itemid: sharedResource.id,
        spaceid: sharedResource.storageId,
        parentitemid: sharedResource.parentFolderId,
        initiatorid: 'local1'
      })
      mocks.clientService.webdav.getFileInfo.mockResolvedValue(
        mock<Resource>({ ...sharedResource })
      )
      await onSSELinkRemovedEvent({ sseData, ...mocks })
      expect(mocks.clientService.webdav.getFileInfo).not.toHaveBeenCalled()
      expect(mocks.resourcesStore.upsertResource).not.toHaveBeenCalled()
      expect(mocks.resourcesStore.updateResourceField).not.toHaveBeenCalled()
      expect(mocks.resourcesStore.removeResources).not.not.toHaveBeenCalled()
    })
  })
})

const getMocks = ({
  currentFolder = mockDeep<Resource>({
    id: 'currenFolder!currentFolder',
    isFolder: true,
    storageId: 'space1'
  }),
  resources = [],
  spaces = [mockDeep<SpaceResource>({ id: 'space1' })],
  currentRouteFilesSpacesGeneric = false,
  currentRouteFilesSpacesProjects = false,
  currentRouteFilesSharesViaLink = false,
  currentRouteFilesSharesWithMe = false,
  currentRouteFilesSharesWithOthers = false
}: {
  currentFolder?: Resource
  resources?: Resource[]
  spaces?: SpaceResource[]
  currentRouteFilesSpacesGeneric?: boolean
  currentRouteFilesSpacesProjects?: boolean
  currentRouteFilesSharesViaLink?: boolean
  currentRouteFilesSharesWithMe?: boolean
  currentRouteFilesSharesWithOthers?: boolean
} = {}) => {
  createTestingPinia()
  const resourcesStore = useResourcesStore()
  resourcesStore.currentFolder = currentFolder
  resourcesStore.resources = resources
  const spacesStore = useSpacesStore()
  spacesStore.spaces = spaces
  const messageStore = useMessages()
  const userStore = useUserStore()
  const configStore = useConfigStore()
  userStore.user = mockDeep<User>({ id: '1' })
  const sharesStore = useSharesStore()
  const clientService = mockDeep<ClientService>({ initiatorId: 'local1' })
  const previewService = mockDeep<PreviewService>()

  const language = mockDeep<Language>({
    $gettext: vi.fn((m) => m)
  })
  const resourceQueue = mockDeep<PQueue>()

  let routeName = 'files-spaces-generic'
  if (currentRouteFilesSpacesGeneric) {
    routeName = 'files-spaces-generic'
  }
  if (currentRouteFilesSpacesProjects) {
    routeName = 'files-spaces-projects'
  }
  if (currentRouteFilesSharesViaLink) {
    routeName = 'files-shares-via-link'
  }
  if (currentRouteFilesSharesWithMe) {
    routeName = 'files-shares-with-me'
  }
  if (currentRouteFilesSharesWithOthers) {
    routeName = 'files-shares-with-others'
  }

  const currentRoute = mock<RouteLocation>({ name: routeName })
  const { $router: router } = defaultComponentMocks({ currentRoute })

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
