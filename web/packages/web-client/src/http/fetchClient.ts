import { HttpError } from '../errors'
import { httpHeaders } from './headers'
import type { FetchClientOptions, FetchRequestOptions, HttpResponse, ResponseType } from './types'

const isBodyInit = (value: unknown): value is BodyInit =>
  typeof value === 'string' ||
  value instanceof Blob ||
  value instanceof FormData ||
  value instanceof URLSearchParams ||
  value instanceof ArrayBuffer ||
  value instanceof ReadableStream ||
  ArrayBuffer.isView(value)

const hasScheme = (url: string) => /^[a-z][a-z0-9+.-]*:/i.test(url)

export class FetchClient {
  constructor(private readonly options: FetchClientOptions = {}) {}

  public async fetch(url: string, options: FetchRequestOptions = {}): Promise<Response> {
    const { method = 'GET', params, body, signal, throwOnError = true } = options

    // Reported to onResponse as-is, never response.url: the maintenance mode allow-list
    // is matched against relative paths.
    const requestUrl = this.appendParams(url, params)
    const headers = this.buildHeaders(options.headers)
    const payload = this.buildBody(body, headers)

    let response: Response
    try {
      response = await fetch(this.resolveUrl(requestUrl), {
        method,
        headers,
        ...(payload !== undefined && { body: payload }),
        ...(signal && { signal })
      })
    } catch (error) {
      // An abort is a caller decision, not a transport failure: propagate it verbatim.
      // `signal.aborted` and not just the error name, because `AbortController.abort(reason)`
      // rejects with that reason as-is — it is only a DOMException named `AbortError` when no
      // reason was given. A caller-supplied reason must not be mistaken for a network failure
      // and reported as a 500, which would also trip maintenance detection.
      if (error?.name === 'AbortError' || signal?.aborted) {
        throw error
      }
      // Degrade a transport failure to 500 so maintenance detection still runs.
      this.options.onResponse?.({ response: null, status: 500, requestUrl })
      throw new HttpError(error?.message || 'Network request failed', null, 500)
    }

    this.options.onResponse?.({ response, status: response.status, requestUrl })

    if (!response.ok && throwOnError) {
      throw await this.buildError(response, options.responseType)
    }

    return response
  }

  public async request<T = unknown>(
    url: string,
    options: FetchRequestOptions = {}
  ): Promise<HttpResponse<T>> {
    const response = await this.fetch(url, options)

    return {
      data: (await this.readBody(response, options.responseType)) as T,
      status: response.status,
      statusText: response.statusText,
      headers: httpHeaders(response.headers)
    }
  }

  /**
   * Merges the three header layers, later ones replacing earlier ones.
   *
   * Deliberately not an object spread into `new Headers()`: header names are
   * case-insensitive, but object keys are not, so `Authorization` and `authorization` would
   * both survive the spread — and a record init is applied with `append`, which joins the
   * two into `Bearer stale, Bearer fresh` instead of overriding. `set()` per entry is
   * case-insensitive and replaces, which is what a layered merge means.
   */
  private buildHeaders(perRequest?: HeadersInit): Headers {
    const headers = new Headers()
    const apply = (layer?: HeadersInit) => {
      if (layer) {
        new Headers(layer).forEach((value, name) => headers.set(name, value))
      }
    }

    apply(this.options.staticHeaders)
    apply(this.options.headers?.())
    apply(perRequest)

    return headers
  }

  private buildBody(body: unknown, headers: Headers): BodyInit | undefined {
    if (body === undefined || body === null) {
      return undefined
    }
    if (body instanceof FormData) {
      // Only fetch knows the multipart boundary it is about to generate, so a caller-supplied
      // `multipart/form-data` header would reach the server without one and make the body
      // unparseable. Dropping it lets fetch fill in the complete header, as axios did.
      if (!headers.get('Content-Type')?.includes('boundary=')) {
        headers.delete('Content-Type')
      }
      return body
    }
    if (isBodyInit(body)) {
      return body
    }
    if (!headers.has('Content-Type')) {
      headers.set('Content-Type', 'application/json')
    }
    return JSON.stringify(body)
  }

  private async buildError(response: Response, responseType?: ResponseType): Promise<HttpError> {
    // Read as the caller asked for the success body: axios applied `responseType` to error
    // bodies too, so a `blob` caller keeps getting a Blob on `error.data`.
    const data = await this.readBodySafely(response, responseType)

    // `error.data` and `error.statusCode` are this repo's convention, but `error.response`
    // is reachable from outside it, where `.response.data` and `.response.headers['x']`
    // were the only spellings. Keep those working too; the `headers` override shadows the
    // prototype accessor with a superset of it. The body itself is already consumed —
    // `error.response.data` is the way to it, not `error.response.json()`.
    Object.defineProperty(response, 'data', { value: data })
    Object.defineProperty(response, 'headers', { value: httpHeaders(response.headers) })

    return new HttpError(
      `Request failed with status code ${response.status}`,
      response,
      response.status,
      data
    )
  }

  private async readBodySafely(response: Response, responseType?: ResponseType): Promise<unknown> {
    try {
      return await this.readBody(response, responseType)
    } catch {
      // an unreadable body must not replace the HTTP error with a read error
      return undefined
    }
  }

  private async readBody(response: Response, responseType: ResponseType = 'json') {
    if (responseType === 'none' || response.status === 204) {
      return undefined
    }

    switch (responseType) {
      case 'text':
        return await response.text()
      case 'blob':
        return await response.blob()
      case 'arraybuffer':
        return await response.arrayBuffer()
      default: {
        const text = await response.text()
        if (!text) {
          return undefined
        }
        try {
          return JSON.parse(text)
        } catch {
          // A caller that does not set a responseType still expects the raw body for a
          // non-JSON response — a WebDAV multistatus, say. Axios did the same, and
          // throwing here would turn a working request into a SyntaxError.
          return text
        }
      }
    }
  }

  private appendParams(url: string, params?: FetchRequestOptions['params']): string {
    if (!params) {
      return url
    }

    const entries = Object.entries(params)
      .filter(([, value]) => value !== undefined && value !== null)
      .map(([key, value]) => [key, String(value)])

    if (!entries.length) {
      return url
    }

    const search = new URLSearchParams(entries).toString()
    return url.includes('?') ? `${url}&${search}` : `${url}?${search}`
  }

  private resolveUrl(url: string): string {
    const { baseUrl } = this.options
    if (!baseUrl || hasScheme(url)) {
      return url
    }
    return `${baseUrl.replace(/\/+$/, '')}/${url.replace(/^\/+/, '')}`
  }
}
