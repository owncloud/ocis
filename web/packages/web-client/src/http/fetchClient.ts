import { HttpError } from '../errors'
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
      if (error?.name === 'AbortError') {
        throw error
      }
      // Degrade a transport failure to 500 so maintenance detection still runs.
      this.options.onResponse?.({ response: null, status: 500, requestUrl })
      throw new HttpError(error?.message || 'Network request failed', null, 500)
    }

    this.options.onResponse?.({ response, status: response.status, requestUrl })

    if (!response.ok && throwOnError) {
      throw await this.buildError(response)
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
      headers: response.headers
    }
  }

  private buildHeaders(perRequest?: Record<string, string>): Headers {
    return new Headers({
      ...(this.options.staticHeaders || {}),
      ...(this.options.headers?.() || {}),
      ...(perRequest || {})
    })
  }

  private buildBody(body: unknown, headers: Headers): BodyInit | undefined {
    if (body === undefined || body === null) {
      return undefined
    }
    if (isBodyInit(body)) {
      return body
    }
    if (!headers.has('Content-Type')) {
      headers.set('Content-Type', 'application/json')
    }
    return JSON.stringify(body)
  }

  private async buildError(response: Response): Promise<HttpError> {
    // Clone so that HttpError.response still exposes an unread body to callers.
    const data = await this.readBodySafely(response.clone())
    return new HttpError(
      response.statusText || `Request failed with status ${response.status}`,
      response,
      response.status,
      data
    )
  }

  private async readBodySafely(response: Response): Promise<unknown> {
    try {
      const text = await response.text()
      if (!text) {
        return undefined
      }
      try {
        return JSON.parse(text)
      } catch {
        return text
      }
    } catch {
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
        return text ? JSON.parse(text) : undefined
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
