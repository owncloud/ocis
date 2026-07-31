import {
  buildIncomingShareResource,
  buildOutgoingShareResource,
  ShareTypes
} from '@ownclouders/web-client'
import { eventBus, isLocationSharesActive, isLocationSpacesActive } from '@ownclouders/web-pkg'
import { SSEEventOptions } from './types'
import { isItemInCurrentFolder } from './helpers'

export const onSSESpaceMemberAddedEvent = async ({
  sseData,
  resourcesStore,
  spacesStore,
  sharesStore,
  clientService,
  router
}: SSEEventOptions) => {
  if (sseData.initiatorid === clientService.initiatorId) {
    // If initiated by current client (browser tab), action unnecessary. Web manages its own logic, return early.
    return
  }

  const space = await clientService.graphAuthenticated.drives.getDrive(
    sseData.itemid,
    sharesStore.graphRoles
  )
  spacesStore.upsertSpace(space)

  if (!isLocationSpacesActive(router, 'files-spaces-projects')) {
    return
  }

  resourcesStore.upsertResource(space)
}

export const onSSESpaceMemberRemovedEvent = async ({
  sseData,
  resourcesStore,
  spacesStore,
  sharesStore,
  messageStore,
  clientService,
  language,
  userStore,
  router
}: SSEEventOptions) => {
  if (sseData.initiatorid === clientService.initiatorId) {
    // If initiated by current client (browser tab), action unnecessary. Web manages its own logic, return early.
    return
  }

  if (!sseData.affecteduserids?.includes(userStore.user.id)) {
    const space = await clientService.graphAuthenticated.drives.getDrive(
      sseData.itemid,
      sharesStore.graphRoles
    )
    return spacesStore.upsertSpace(space)
  }

  const removedSpace = spacesStore.spaces.find((space) => space.id === sseData.spaceid)
  if (!removedSpace) {
    return
  }

  const { $gettext } = language

  spacesStore.removeSpace(removedSpace)

  if (
    isLocationSpacesActive(router, 'files-spaces-generic') &&
    removedSpace.id === resourcesStore.currentFolder?.storageId
  ) {
    // Fixme: Message triggers when user membership was revoked, but still is in a group membership, which is wrong
    return messageStore.showMessage({
      title: $gettext(
        'Your access to this space has been revoked. Please navigate to another location.'
      )
    })
  }

  if (isLocationSpacesActive(router, 'files-spaces-projects')) {
    return resourcesStore.removeResources([removedSpace])
  }
}

export const onSSESpaceShareUpdatedEvent = async ({
  sseData,
  resourcesStore,
  spacesStore,
  sharesStore,
  clientService,
  userStore,
  router
}: SSEEventOptions) => {
  if (sseData.initiatorid === clientService.initiatorId) {
    // If initiated by current client (browser tab), action unnecessary. Web manages its own logic, return early.
    return
  }

  const space = await clientService.graphAuthenticated.drives.getDrive(
    sseData.itemid,
    sharesStore.graphRoles
  )
  spacesStore.upsertSpace(space)

  if (
    isLocationSpacesActive(router, 'files-spaces-generic') &&
    sseData.affecteduserids?.includes(userStore.user.id) &&
    resourcesStore.currentFolder?.storageId === sseData.spaceid
  ) {
    return eventBus.publish('app.files.list.load')
  }
}

export const onSSEShareCreatedEvent = async ({
  sseData,
  resourcesStore,
  spacesStore,
  sharesStore,
  userStore,
  configStore,
  clientService,
  router
}: SSEEventOptions) => {
  if (sseData.initiatorid === clientService.initiatorId) {
    // If initiated by current client (browser tab), action unnecessary. Web manages its own logic, return early.
    return
  }

  if (
    isLocationSpacesActive(router, 'files-spaces-generic') &&
    isItemInCurrentFolder({
      resourcesStore,
      parentFolderId: sseData.parentitemid
    })
  ) {
    const space = spacesStore.spaces.find((space) => space.id === sseData.spaceid)
    if (!space) {
      return
    }

    const resource = await clientService.webdav.getFileInfo(space, {
      path: '',
      fileId: sseData.itemid
    })

    resourcesStore.upsertResource(resource)
  }

  if (isLocationSharesActive(router, 'files-shares-with-me')) {
    // FIXME: get drive item by id as soon as server supports it
    const driveItems = await clientService.graphAuthenticated.driveItems.listSharedWithMe()
    const driveItem = driveItems.find(({ remoteItem }) => remoteItem.id === sseData.itemid)
    if (!driveItem) {
      return
    }
    const resource = buildIncomingShareResource({
      driveItem,
      graphRoles: sharesStore.graphRoles,
      serverUrl: configStore.serverUrl
    })
    return resourcesStore.upsertResource(resource)
  }

  if (isLocationSharesActive(router, 'files-shares-with-others')) {
    // FIXME: get drive item by id as soon as server supports it
    const driveItems = await clientService.graphAuthenticated.driveItems.listSharedByMe()
    const driveItem = driveItems.find(({ id }) => id === sseData.itemid)
    if (!driveItem) {
      return
    }
    const resource = buildOutgoingShareResource({
      driveItem,
      user: userStore.user,
      serverUrl: configStore.serverUrl
    })
    return resourcesStore.upsertResource(resource)
  }
}
export const onSSEShareUpdatedEvent = async ({
  sseData,
  resourcesStore,
  sharesStore,
  clientService,
  userStore,
  configStore,
  router
}: SSEEventOptions) => {
  if (sseData.initiatorid === clientService.initiatorId) {
    // If initiated by current client (browser tab), action unnecessary. Web manages its own logic, return early.
    return
  }

  if (
    isLocationSpacesActive(router, 'files-spaces-generic') &&
    sseData.affecteduserids?.includes(userStore.user.id) &&
    resourcesStore.currentFolder?.storageId === sseData.spaceid
  ) {
    return eventBus.publish('app.files.list.load')
  }

  if (isLocationSharesActive(router, 'files-shares-with-me')) {
    // FIXME: get drive item by id as soon as server supports it
    const driveItems = await clientService.graphAuthenticated.driveItems.listSharedWithMe()
    const driveItem = driveItems.find(({ remoteItem }) => remoteItem.id === sseData.itemid)
    if (!driveItem) {
      return
    }
    const resource = buildIncomingShareResource({
      driveItem,
      graphRoles: sharesStore.graphRoles,
      serverUrl: configStore.serverUrl
    })
    return resourcesStore.upsertResource(resource)
  }
}

export const onSSEShareRemovedEvent = async ({
  sseData,
  resourcesStore,
  spacesStore,
  userStore,
  clientService,
  messageStore,
  language,
  router
}: SSEEventOptions) => {
  if (sseData.initiatorid === clientService.initiatorId) {
    // If initiated by current client (browser tab), action unnecessary. Web manages its own logic, return early.
    return
  }
  const { $gettext } = language

  if (
    isLocationSpacesActive(router, 'files-spaces-generic') &&
    sseData.affecteduserids?.includes(userStore.user.id) &&
    resourcesStore.currentFolder?.storageId === sseData.spaceid
  ) {
    // Fixme: Message triggers when user share was revoked, but still is in a group share, which is wrong
    return messageStore.showMessage({
      title: $gettext(
        'Your access to this share has been revoked. Please navigate to another location.'
      )
    })
  }

  if (
    isLocationSpacesActive(router, 'files-spaces-generic') &&
    isItemInCurrentFolder({
      resourcesStore,
      parentFolderId: sseData.parentitemid
    })
  ) {
    const space = spacesStore.spaces.find((space) => space.id === sseData.spaceid)
    if (!space) {
      return
    }

    const resource = await clientService.webdav.getFileInfo(space, {
      path: '',
      fileId: sseData.itemid
    })

    resourcesStore.upsertResource(resource)
  }

  if (isLocationSharesActive(router, 'files-shares-with-others')) {
    const space = spacesStore.spaces.find((space) => space.id === sseData.spaceid)
    if (!space) {
      return
    }

    const resource = await clientService.webdav.getFileInfo(space, {
      path: '',
      fileId: sseData.itemid
    })

    if (
      !resource.shareTypes.includes(ShareTypes.user.value) &&
      !resource.shareTypes.includes(ShareTypes.group.value)
    ) {
      return resourcesStore.removeResources([resource])
    }
  }

  if (isLocationSharesActive(router, 'files-shares-with-me')) {
    const removedShareResource = resourcesStore.resources.find((r) => r.fileId === sseData.itemid)
    if (!removedShareResource) {
      return
    }
    return resourcesStore.removeResources([removedShareResource])
  }
}

export const onSSELinkCreatedEvent = async ({
  sseData,
  resourcesStore,
  spacesStore,
  userStore,
  configStore,
  clientService,
  router
}: SSEEventOptions) => {
  if (sseData.initiatorid === clientService.initiatorId) {
    // If initiated by current client (browser tab), action unnecessary. Web manages its own logic, return early.
    return
  }

  if (
    isLocationSpacesActive(router, 'files-spaces-generic') &&
    isItemInCurrentFolder({
      resourcesStore,
      parentFolderId: sseData.parentitemid
    })
  ) {
    const space = spacesStore.spaces.find((space) => space.id === sseData.spaceid)
    if (!space) {
      return
    }

    const resource = await clientService.webdav.getFileInfo(space, {
      path: '',
      fileId: sseData.itemid
    })

    resourcesStore.upsertResource(resource)
  }

  if (isLocationSharesActive(router, 'files-shares-via-link')) {
    // FIXME: get drive item by id as soon as server supports it
    const driveItems = await clientService.graphAuthenticated.driveItems.listSharedByMe()
    const driveItem = driveItems.find(({ id }) => id === sseData.itemid)
    if (!driveItem) {
      return
    }
    const resource = buildOutgoingShareResource({
      driveItem,
      user: userStore.user,
      serverUrl: configStore.serverUrl
    })
    return resourcesStore.upsertResource(resource)
  }
}

export const onSSELinkUpdatedEvent = ({ sseData, clientService }: SSEEventOptions) => {
  if (sseData.initiatorid === clientService.initiatorId) {
    // If initiated by current client (browser tab), action unnecessary. Web manages its own logic, return early.
    return
  }
}

export const onSSELinkRemovedEvent = async ({
  sseData,
  resourcesStore,
  spacesStore,
  userStore,
  clientService,
  router
}: SSEEventOptions) => {
  if (sseData.initiatorid === clientService.initiatorId) {
    // If initiated by current client (browser tab), action unnecessary. Web manages its own logic, return early.
    return
  }

  const space = spacesStore.spaces.find((space) => space.id === sseData.spaceid)
  if (!space) {
    return
  }

  const resource = await clientService.webdav.getFileInfo(space, {
    path: '',
    fileId: sseData.itemid
  })

  if (
    isLocationSpacesActive(router, 'files-spaces-generic') &&
    isItemInCurrentFolder({
      resourcesStore,
      parentFolderId: sseData.parentitemid
    })
  ) {
    resourcesStore.upsertResource(resource)
  }

  if (isLocationSharesActive(router, 'files-shares-via-link')) {
    if (!resource.shareTypes.includes(ShareTypes.link.value)) {
      return resourcesStore.removeResources([resource])
    }
  }
}
