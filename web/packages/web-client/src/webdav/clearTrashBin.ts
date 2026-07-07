import { Resource, SpaceResource, buildWebDavSpacesTrashPath } from '../helpers'
import { WebDavOptions } from './types'
import { DAV, DAVRequestOptions } from './client'
import { urlJoin } from '../utils'

type ClearTrashBinOptions = {
  id?: Resource['id']
} & DAVRequestOptions

export const ClearTrashBinFactory = (dav: DAV, options: WebDavOptions) => {
  return {
    clearTrashBin(space: SpaceResource, { id, ...opts }: ClearTrashBinOptions = {}) {
      let path = buildWebDavSpacesTrashPath(space.id)

      if (id) {
        path = urlJoin(path, id)
      }

      return dav.delete(path, opts)
    }
  }
}
