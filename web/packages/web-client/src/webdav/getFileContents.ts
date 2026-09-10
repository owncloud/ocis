import { SpaceResource } from '../helpers'
import { WebDavOptions } from './types'
import { DAV, DAVRequestOptions } from './client'
import { HttpError } from '../errors'
import type { ResponseType } from '../http'
import { getWebDavPath } from './utils'

export type GetFileContentsResponse = {
  body: any
  [key: string]: any
}

export const GetFileContentsFactory = (dav: DAV, { httpClient }: WebDavOptions) => {
  return {
    async getFileContents(
      space: SpaceResource,
      { fileId, path }: { fileId?: string; path?: string },
      {
        responseType = 'text',
        noCache = true,
        headers,
        ...opts
      }: {
        responseType?: ResponseType
        noCache?: boolean
      } & DAVRequestOptions = {}
    ): Promise<GetFileContentsResponse> {
      try {
        const webDavPath = getWebDavPath(space, { fileId, path })
        const response = await httpClient.request(dav.getFileUrl(webDavPath), {
          responseType,
          headers: {
            ...(noCache && { 'Cache-Control': 'no-cache' }),
            ...(headers || {})
          },
          ...opts
        })
        return {
          response,
          body: response.data,
          headers: {
            // bracket access, not get(): an absent header stays undefined here rather than null
            ETag: response.headers['etag'],
            'OC-ETag': response.headers['oc-etag'],
            'OC-FileId': response.headers['oc-fileid']
          }
        }
      } catch (error) {
        // The core already throws an HttpError carrying the response and status, and lets an
        // abort through verbatim. Both have to reach the caller unchanged — wrapping an abort
        // would hide the `AbortError` name that callers distinguish a cancelled load by, and
        // turn a restarted load into a visible error.
        if (error instanceof HttpError || error?.name === 'AbortError' || opts.signal?.aborted) {
          throw error
        }
        throw new HttpError(error?.message, error?.response, error?.statusCode)
      }
    }
  }
}
