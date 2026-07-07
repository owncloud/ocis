import { WebDavOptions } from './types'
import { urlJoin } from '../utils'
import { DAV, DAVRequestOptions } from './client'
import { buildResource } from '../helpers'

export const ListFileVersionsFactory = (dav: DAV, options: WebDavOptions) => {
  return {
    async listFileVersions(id: string, opts: DAVRequestOptions = {}) {
      // eslint-disable-next-line @typescript-eslint/no-unused-vars
      const [currentFolder, ...versions] = await dav.propfind(
        urlJoin('meta', id, 'v', { leadingSlash: true }),
        opts
      )
      return versions.map((v) => buildResource(v, dav.extraProps))
    }
  }
}
