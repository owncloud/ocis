import { isPasswordProtectedFolderFileResource, Resource } from '@ownclouders/web-client'

/**
 * Asserts whether given resource can be moved
 * @param {String} resource Resource which is to be moved
 * @param {Object} parentPath Path of the parent folder of the resource
 * @return {Boolean} can be moved
 */
export function canBeMoved(resource: Resource, parentPath: string) {
  // TODO: Find a better way to prevent moving shared resources than by checking if the current folder is root
  // TODO: Find a way to disable move action when shares are mounted in different folder then root
  const isExternal = resource.isReceivedShare() || resource.isMounted()
  const isMountedInRoot = parentPath === '' && isExternal

  if (isPasswordProtectedFolderFileResource(resource.name)) {
    return false
  }

  return resource.canBeDeleted() && !isMountedInRoot
}
