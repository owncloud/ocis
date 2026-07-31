import { Resource } from '@ownclouders/web-client'
import {
  extractParentFolderName,
  isProjectSpaceResource,
  isShareRoot,
  isShareSpaceResource
} from '@ownclouders/web-client'
import { useGettext } from 'vue3-gettext'
import { unref } from 'vue'
import { useGetMatchingSpace } from '../spaces'
import path, { dirname } from 'path'
import { ResourceRouteResolverOptions, useResourceRouteResolver } from '../filesList'
import { createLocationShares, createLocationSpaces } from '../../router'
import { useCapabilityStore } from '../piniaStores'

export const useFolderLink = (options: ResourceRouteResolverOptions = {}) => {
  const capabilityStore = useCapabilityStore()
  const { $gettext } = useGettext()
  const { getInternalSpace, getMatchingSpace, isResourceAccessible } = useGetMatchingSpace()
  const { createFolderLink } = useResourceRouteResolver(options)

  const getPathPrefix = (resource: Resource) => {
    const space = unref(options.space) || getMatchingSpace(resource)

    if (isProjectSpaceResource(space)) {
      if (resource.path === '.') {
        return $gettext('Spaces')
      }
      return path.join($gettext('Spaces'), space.name)
    }

    if (isShareSpaceResource(space)) {
      return path.join($gettext('Shares'), space.name)
    }

    return space.name
  }

  const getFolderLink = (resource: Resource) => {
    return createFolderLink({
      path: resource.path,
      fileId: resource.fileId,
      resource
    })
  }

  const getParentFolderLink = (resource: Resource) => {
    const space = unref(options.space) || getMatchingSpace(resource)
    const parentFolderAccessible = isResourceAccessible({
      space,
      path: dirname(resource.path)
    })
    if ((resource.remoteItemId && resource.path === '/') || !parentFolderAccessible) {
      return createLocationShares('files-shares-with-me')
    }
    if (isProjectSpaceResource(resource)) {
      return createLocationSpaces('files-spaces-projects')
    }
    if (resource.path === '.') {
      return createLocationSpaces('files-spaces-projects')
    }

    return createFolderLink({
      path: dirname(resource.path),
      ...(resource.parentFolderId && { fileId: resource.parentFolderId }),
      resource
    })
  }

  const getParentFolderName = (resource: Resource) => {
    const space = unref(options.space) || getMatchingSpace(resource)
    const parentFolderAccessible = isResourceAccessible({
      space,
      path: dirname(resource.path)
    })

    //FIXME: As soon the backend exposes oc-share-root via webdav, only use isShareRoot fn
    const shareRoot =
      isShareRoot(resource) || (resource.id === space.id && isShareSpaceResource(space))
    if (shareRoot || !parentFolderAccessible) {
      return $gettext('Shared with me')
    }
    const parentFolder = extractParentFolderName(resource)
    if (parentFolder) {
      if (resource.path === '.' && parentFolder === '.') {
        return $gettext('Spaces')
      }
      return parentFolder
    }

    if (isShareSpaceResource(space)) {
      return space.name
    }

    if (capabilityStore.spacesProjects) {
      if (isProjectSpaceResource(resource)) {
        return $gettext('Spaces')
      }
      if (space?.driveType === 'project') {
        return space.name
      }
    }

    return $gettext('Personal')
  }

  const getParentFolderLinkIconAdditionalAttributes = (resource: Resource) => {
    // Identify if resource is project space or is part of a project space and the resource is located in its root
    if (
      isProjectSpaceResource(resource) ||
      (isProjectSpaceResource(getInternalSpace(resource.storageId) || ({} as Resource)) &&
        resource.path.split('/').length === 2)
    ) {
      return {
        name: 'layout-grid',
        'fill-type': 'fill'
      }
    }

    return {}
  }

  return {
    getPathPrefix,
    getFolderLink,
    getParentFolderLink,
    getParentFolderName,
    getParentFolderLinkIconAdditionalAttributes
  }
}
