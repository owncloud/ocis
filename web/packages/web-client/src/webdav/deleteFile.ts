import { urlJoin } from '../utils'
import { SpaceResource } from '../helpers'
import { WebDavOptions } from './types'
import { DAV, DAVRequestOptions } from './client'

export const DeleteFileFactory = (dav: DAV, options: WebDavOptions) => {
  return {
    deleteFile(space: SpaceResource, { path, ...opts }: { path: string } & DAVRequestOptions) {
      // FIXME: id-based not possible because of the server: https://github.com/owncloud/ocis/issues/9619
      return dav.delete(urlJoin(space.webDavPath, path), opts)
    }
  }
}
