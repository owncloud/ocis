import { createFileRouteOptions, ImageDimension } from '@ownclouders/web-pkg'
import { SSEEventOptions } from './types'
import { isItemInCurrentFolder } from './helpers'

export const onSSEItemRenamedEvent = async ({
  sseData,
  resourcesStore,
  spacesStore,
  clientService,
  router
}: SSEEventOptions) => {
  if (sseData.initiatorid === clientService.initiatorId) {
    // If initiated by current client (browser tab), action unnecessary. Web manages its own logic, return early.
    return
  }

  const currentFolder = resourcesStore.currentFolder
  const resourceIsCurrentFolder = currentFolder?.id === sseData.itemid
  const resource = resourceIsCurrentFolder
    ? currentFolder
    : resourcesStore.resources.find((f) => f.id === sseData.itemid)

  if (!resource) {
    return
  }
  const space = spacesStore.spaces.find((s) => s.id === resource.storageId)

  if (!space) {
    return
  }

  const updatedResource = await clientService.webdav.getFileInfo(space, {
    path: '',
    fileId: sseData.itemid
  })

  if (resourceIsCurrentFolder) {
    resourcesStore.setCurrentFolder(updatedResource)
    return router.push(
      createFileRouteOptions(space, {
        path: updatedResource.path,
        fileId: updatedResource.fileId
      })
    )
  }

  resourcesStore.upsertResource(updatedResource)
}

export const onSSEFileLockingEvent = async ({
  sseData,
  resourcesStore,
  spacesStore,
  clientService
}: SSEEventOptions) => {
  const resource = resourcesStore.resources.find((f) => f.id === sseData.itemid)
  const space = spacesStore.spaces.find((s) => s.id === resource?.storageId)

  if (!resource || !space) {
    return
  }

  const updatedResource = await clientService.webdav.getFileInfo(space, {
    path: '',
    fileId: sseData.itemid
  })

  resourcesStore.upsertResource(updatedResource)
}

export const onSSEProcessingFinishedEvent = async ({
  sseData,
  resourcesStore,
  spacesStore,
  clientService,
  resourceQueue,
  previewService
}: SSEEventOptions) => {
  if (!isItemInCurrentFolder({ resourcesStore, parentFolderId: sseData.parentitemid })) {
    return false
  }
  const resource = resourcesStore.resources.find((f) => f.id === sseData.itemid)
  const space = spacesStore.spaces.find((s) => s.id === sseData.spaceid)
  if (!space) {
    return
  }

  /**
   * If resource is not loaded, it suggests an upload is in progress.
   */
  if (!resource) {
    if (sseData.initiatorid === clientService.initiatorId) {
      /**
       * If the upload is initiated by the current client,
       * there's no necessity to retrieve the resources again.
       */
      return
    }

    return resourceQueue.add(async () => {
      const { resource } = await clientService.webdav.listFiles(space, {
        path: '',
        fileId: sseData.itemid
      })

      // check again for the current folder in case the user has navigated away in the meantime
      if (isItemInCurrentFolder({ resourcesStore, parentFolderId: sseData.parentitemid })) {
        resourcesStore.upsertResource(resource)
      }
    })
  }

  /**
   * Resource not changed, don't fetch more data
   */
  if (resource.etag === sseData.etag) {
    return resourcesStore.updateResourceField({
      id: sseData.itemid,
      field: 'processing',
      value: false
    })
  }

  const updatedResource = await clientService.webdav.getFileInfo(space, {
    path: '',
    fileId: sseData.itemid
  })
  resourcesStore.upsertResource(updatedResource)

  const preview = await previewService.loadPreview({
    resource,
    space,
    dimensions: ImageDimension.Thumbnail
  })

  if (preview) {
    resourcesStore.updateResourceField({
      id: sseData.itemid,
      field: 'thumbnail',
      value: preview
    })
  }
}

export const onSSEItemTrashedEvent = ({
  sseData,
  language,
  messageStore,
  resourcesStore,
  clientService
}: SSEEventOptions) => {
  if (sseData.initiatorid === clientService.initiatorId) {
    // If initiated by current client (browser tab), action unnecessary. Web manages its own logic, return early.
    return
  }
  const { $gettext } = language

  const currentFolder = resourcesStore.currentFolder
  const resourceIsCurrentFolder = currentFolder?.id === sseData.itemid

  if (resourceIsCurrentFolder) {
    return messageStore.showMessage({
      title: $gettext(
        'The folder you were accessing has been removed. Please navigate to another location.'
      )
    })
  }

  const resource = resourcesStore.resources.find((f) => f.id === sseData.itemid)

  if (!resource) {
    return
  }

  resourcesStore.removeResources([resource])
}

export const onSSEItemRestoredEvent = async ({
  sseData,
  resourcesStore,
  spacesStore,
  userStore,
  clientService
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

  if (!resource) {
    return
  }

  if (!isItemInCurrentFolder({ resourcesStore, parentFolderId: resource.parentFolderId })) {
    return false
  }

  resourcesStore.upsertResource(resource)
}

export const onSSEItemMovedEvent = async ({
  sseData,
  resourcesStore,
  spacesStore,
  userStore,
  clientService
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

  if (!resource) {
    return
  }

  if (resource.parentFolderId !== resourcesStore.currentFolder?.id) {
    return resourcesStore.removeResources([resource])
  }

  resourcesStore.upsertResource(resource)
}

/**
 * The FileTouched event is triggered when a new empty file, such as a new text file,
 * is about to be created on the server. This event is necessary because the
 * post-processing event won't be triggered in this case.
 */
export const onSSEFileTouchedEvent = async ({
  sseData,
  resourcesStore,
  spacesStore,
  clientService
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

  if (!resource) {
    return
  }

  if (!isItemInCurrentFolder({ resourcesStore, parentFolderId: resource.parentFolderId })) {
    return false
  }

  resourcesStore.upsertResource(resource)
}

export const onSSEFolderCreatedEvent = async ({
  sseData,
  resourcesStore,
  spacesStore,
  clientService
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

  if (!resource) {
    return
  }

  if (!isItemInCurrentFolder({ resourcesStore, parentFolderId: resource.parentFolderId })) {
    return false
  }

  resourcesStore.upsertResource(resource)
}
