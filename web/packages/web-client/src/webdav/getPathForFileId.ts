import { urlJoin } from '../utils'
import { DAV, DAVRequestOptions } from './client'
import { DavProperty } from './constants'
import { WebDavOptions } from './types'

export const GetPathForFileIdFactory = (dav: DAV, options: WebDavOptions) => {
  return {
    async getPathForFileId(id: string, opts: DAVRequestOptions = {}) {
      const result = await dav.propfind(urlJoin('meta', id, { leadingSlash: true }), {
        properties: [DavProperty.MetaPathForUser],
        ...opts
      })

      return result[0].props[DavProperty.MetaPathForUser]
    }
  }
}
