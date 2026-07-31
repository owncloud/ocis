import {
  Resource,
  isMountPointSpaceResource,
  OCM_PROVIDER_ID,
  buildSpace
} from '@ownclouders/web-client'
import { computed, unref } from 'vue'
import { useClientService } from '../clientService'
import { urlJoin } from '@ownclouders/web-client'
import { useSpacesStore, useSharesStore } from '../piniaStores'
import { DavProperty } from '@ownclouders/web-client/webdav'

export const useGetResourceContext = () => {
  const clientService = useClientService()
  const spacesStore = useSpacesStore()
  const shareStore = useSharesStore()

  const spaces = computed(() => spacesStore.spaces)

  const getMatchingSpaceByFileId = (id: Resource['id']) => {
    return unref(spaces).find((space) => id.toString().startsWith(space.id.toString()))
  }
  const getMatchingMountPoint = (id: Resource['id']) => {
    return unref(spaces).find(
      (space) => isMountPointSpaceResource(space) && space.root?.remoteItem?.id === id
    )
  }

  const loadFileInfoById = (fileId: string) => {
    const davProperties = [
      DavProperty.FileId,
      DavProperty.FileParent,
      DavProperty.Name,
      DavProperty.ResourceType
    ]

    const tmpSpace = buildSpace({ id: fileId, name: '' }, shareStore.graphRoles)
    return clientService.webdav.getFileInfo(tmpSpace, { fileId }, { davProperties })
  }

  // get context for a resource when only having its id. be careful, this might be very expensive!
  const getResourceContext = async (id: string) => {
    let path: string
    let resource: Resource
    let space = getMatchingSpaceByFileId(id)

    if (space) {
      path = await clientService.webdav.getPathForFileId(id)
      resource = await clientService.webdav.getFileInfo(space, { path })
      return { space, resource, path }
    }

    // no matching space found => the file doesn't lie in own spaces => it's a share.
    // do PROPFINDs on parents until root of accepted share is found in `mountpoint` spaces
    await spacesStore.loadMountPoints({ graphClient: clientService.graphAuthenticated })

    let mountPoint = getMatchingMountPoint(id)
    resource = await loadFileInfoById(id)
    const sharePathSegments = mountPoint ? [] : [unref(resource).name]
    let tmpResource = unref(resource)

    while (!mountPoint) {
      tmpResource = await loadFileInfoById(tmpResource.parentFolderId)
      mountPoint = getMatchingMountPoint(tmpResource.id)
      if (!mountPoint) {
        sharePathSegments.unshift(tmpResource.name)
      }
    }

    space =
      spacesStore.getSpace(mountPoint.root?.remoteItem?.id) ||
      spacesStore.createShareSpace({
        driveAliasPrefix: resource.storageId?.startsWith(OCM_PROVIDER_ID) ? 'ocm-share' : 'share',
        id: mountPoint.root?.remoteItem?.id,
        shareName: mountPoint.name
      })

    path = urlJoin(...sharePathSegments)
    return { space, resource, path }
  }

  return {
    getResourceContext
  }
}
