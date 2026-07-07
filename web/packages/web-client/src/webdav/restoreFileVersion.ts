import { SpaceResource } from '../helpers'
import { WebDavOptions } from './types'
import { urlJoin } from '../utils'
import { DAV, DAVRequestOptions } from './client'
import { getWebDavPath } from './utils'

export const RestoreFileVersionFactory = (dav: DAV, options: WebDavOptions) => {
  return {
    restoreFileVersion(
      space: SpaceResource,
      { parentFolderId, name, path }: { parentFolderId?: string; name?: string; path?: string },
      versionId: string,
      opts: DAVRequestOptions = {}
    ) {
      const webDavPath = getWebDavPath(space, { path, fileId: parentFolderId, name })
      const source = urlJoin('meta', parentFolderId, 'v', versionId, { leadingSlash: true })
      const target = urlJoin('files', webDavPath, { leadingSlash: true })
      return dav.copy(source, target, opts)
    }
  }
}
