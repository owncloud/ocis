import { urlJoin } from '../utils'
import { WebDavOptions } from './types'
import { DAV, DAVRequestOptions } from './client'
import { DavProperties, DavPropertyValue } from './constants'

export const ListFavoriteFilesFactory = (dav: DAV, options: WebDavOptions) => {
  return {
    listFavoriteFiles({
      davProperties = DavProperties.Default,
      username = '',
      ...opts
    }: { davProperties?: DavPropertyValue[]; username?: string } & DAVRequestOptions = {}) {
      return dav.report(urlJoin('files', username), {
        properties: davProperties,
        filterRules: { favorite: 1 },
        ...opts
      })
    }
  }
}
