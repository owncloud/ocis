import {
  Headers,
  ProgressEventCallback,
  RequestOptions,
  RequestOptionsCustom,
  WebDAVClient,
  createClient
} from 'webdav'
import { v4 as uuidV4 } from 'uuid'
import { encodePath, urlJoin } from '../../utils'
import { DavMethod, DavPropertyValue } from '../constants'
import { buildPropFindBody, buildPropPatchBody } from './builders'
import { parseError, parseMultiStatus, parseTusHeaders } from './parsers'
import { WebDavResponseResource } from '../../helpers'
import { DavHttpError } from '../../errors'
import { normalize } from 'path'

export interface DAVOptions {
  baseUrl: string
  headers?: () => Headers
  onSetMaintenance: (value: boolean) => void
}

export interface DavResult {
  body: WebDavResponseResource[] | undefined
  status: number
  result: Response
}

export type DAVRequestOptions = {
  headers?: Headers
  signal?: AbortSignal
}

export class DAV {
  private client: WebDAVClient
  private davPath: string
  private headers: () => Headers
  public extraProps: string[]
  public remoteBasePath: string

  #onSetMaintenance: (value: boolean) => void

  constructor({ baseUrl, headers, onSetMaintenance }: DAVOptions) {
    const remoteBasePath = extractUrlPath(baseUrl)

    this.davPath = urlJoin(baseUrl, 'dav')
    this.client = createClient(this.davPath, { remoteBasePath: remoteBasePath })
    this.headers = headers
    this.extraProps = []
    this.#onSetMaintenance = onSetMaintenance
    this.remoteBasePath = remoteBasePath
  }

  public mkcol(path: string, opts: DAVRequestOptions = {}) {
    return this.request(path, { method: DavMethod.mkcol, ...opts })
  }

  public async propfind(
    path: string,
    {
      depth = 1,
      properties = [],
      headers = {},
      ...opts
    }: { depth?: number; properties?: DavPropertyValue[] } & DAVRequestOptions = {}
  ) {
    const requestHeaders = { ...headers, Depth: depth.toString() }
    const { body, result } = await this.request(path, {
      method: DavMethod.propfind,
      data: buildPropFindBody(properties, { extraProps: this.extraProps }),
      headers: requestHeaders,
      ...opts
    })

    if (body?.length) {
      // add tus info to current folder only
      body[0].tusSupport = parseTusHeaders(result.headers)
    }

    return body
  }

  public async report(
    path: string,
    {
      pattern = '',
      filterRules = null,
      limit = 30,
      properties,
      ...opts
    }: {
      pattern?: string
      filterRules?: Partial<Record<DavPropertyValue, unknown>>
      limit?: number
      properties?: DavPropertyValue[]
    } & DAVRequestOptions = {}
  ) {
    const { body, result } = await this.request(path, {
      method: DavMethod.report,
      data: buildPropFindBody(properties, {
        pattern,
        filterRules,
        limit,
        extraProps: this.extraProps
      }),
      ...opts
    })

    return {
      results: body,
      range: result.headers.get('content-range')
    }
  }

  public copy(
    source: string,
    target: string,
    { overwrite = false, headers = {}, ...opts }: { overwrite?: boolean } & DAVRequestOptions = {}
  ) {
    const targetUrl = urlJoin(this.davPath, encodePath(target))
    return this.request(source, {
      method: DavMethod.copy,
      headers: { ...headers, Destination: targetUrl, overwrite: overwrite ? 'T' : 'F' },
      ...opts
    })
  }

  public move(
    source: string,
    target: string,
    { overwrite = false, headers = {}, ...opts }: { overwrite?: boolean } & DAVRequestOptions = {}
  ) {
    const targetUrl = urlJoin(this.davPath, encodePath(target))
    return this.request(source, {
      method: DavMethod.move,
      headers: { ...headers, Destination: targetUrl, overwrite: overwrite ? 'T' : 'F' },
      ...opts
    })
  }

  public put(
    path: string,
    content: string | ArrayBuffer,
    {
      headers = {},
      onUploadProgress,
      previousEntityTag,
      overwrite,
      ...opts
    }: {
      onUploadProgress?: ProgressEventCallback
      previousEntityTag?: string
      overwrite?: boolean
    } & DAVRequestOptions = {}
  ) {
    const requestHeaders = { ...headers }
    if (previousEntityTag) {
      // will ensure that no other client uploaded a different version meanwhile
      requestHeaders['If-Match'] = previousEntityTag
    } else if (!overwrite) {
      // will trigger 412 precondition failed if a file already exists
      requestHeaders['If-None-Match'] = '*'
    }

    return this.request(path, {
      method: DavMethod.put,
      data: content,
      headers: requestHeaders,
      onUploadProgress,
      ...opts
    })
  }

  public delete(path: string, opts: DAVRequestOptions = {}) {
    return this.request(path, { method: DavMethod.delete, ...opts })
  }

  public propPatch(
    path: string,
    properties: Partial<Record<DavPropertyValue, unknown>>,
    opts: DAVRequestOptions = {}
  ) {
    const body = buildPropPatchBody(properties)
    return this.request(path, { method: DavMethod.proppatch, data: body, ...opts })
  }

  public getFileUrl(path: string) {
    return urlJoin(this.davPath, encodePath(path))
  }

  private buildHeaders(headers: Headers = {}): Headers {
    return {
      'Content-Type': 'application/xml; charset=utf-8',
      'X-Requested-With': 'XMLHttpRequest',
      'X-Request-ID': uuidV4(),
      ...(this.headers && { ...this.headers() }),
      ...headers
    }
  }

  private async request(path: string, options: RequestOptionsCustom): Promise<DavResult> {
    const url = urlJoin(this.davPath, encodePath(path), { leadingSlash: true })

    const requestOptions = {
      ...options,
      url,
      headers: this.buildHeaders(options.headers || {})
    } as RequestOptions

    try {
      const result = (await this.client.customRequest('', requestOptions)) as unknown as Response

      let resultBody: WebDavResponseResource[]
      if (result.status === 207) {
        const parsedBody = await result.text()
        resultBody = await parseMultiStatus(parsedBody, this.remoteBasePath)
      }

      this.#onSetMaintenance(false)

      return {
        body: resultBody,
        status: result.status,
        result
      }
    } catch (error) {
      const { response } = error

      if (response.status === 503) {
        this.#onSetMaintenance(true)
      }

      const body = await response.text()
      const errorMessage = parseError(body)
      throw new DavHttpError(
        errorMessage.message,
        errorMessage.errorCode,
        response,
        response.status
      )
    }
  }
}

function extractUrlPath(urlString: string) {
  const url = new URL(urlString)
  const path = url.pathname
  return normalize(path)
}
