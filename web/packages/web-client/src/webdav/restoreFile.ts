import { isPublicSpaceResource, SpaceResource } from '../helpers'
import { WebDavOptions } from './types'
import { urlJoin } from '../utils'
import { DAV, DAVRequestOptions } from './client'

export const RestoreFileFactory = (dav: DAV, options: WebDavOptions) => {
  return {
    restoreFile(
      space: SpaceResource,
      { id }: { id: string },
      { path: restorePath }: { path: string },
      { overwrite, ...opts }: { overwrite?: boolean } & DAVRequestOptions = {}
    ) {
      if (isPublicSpaceResource(space)) {
        return
      }

      const restoreWebDavPath = urlJoin(space.webDavPath, restorePath)
      return dav.move(urlJoin(space.webDavTrashPath, id), restoreWebDavPath, {
        overwrite,
        ...opts
      })
    }
  }
}
