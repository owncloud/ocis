import { isPublicSpaceResource, SpaceResource } from '../helpers'
import { urlJoin } from '../utils'

/**
 * Builds a webdav path based on a given `path` or `fileId`. A `path` takes precedence.
 *
 * Public spaces currently don't support id-based paths, hence `path` needs to be provided.
 * Some id-based requests need a resource `name` appended (mkcol, put, copy, move, restore).
 * In this case, the `fileId` is supposed to be the id of the parent folder.
 **/
export const getWebDavPath = (
  space: SpaceResource,
  { fileId, path, name }: { fileId?: string; path?: string; name?: string }
) => {
  if (path !== undefined) {
    return urlJoin(space.webDavPath, path)
  }

  if (fileId !== undefined) {
    if (isPublicSpaceResource(space)) {
      throw new Error('public spaces need a path provided')
    }

    return urlJoin('spaces', fileId, name || '')
  }

  return space.webDavPath
}
