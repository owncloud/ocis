import {
  FetchClient,
  type FetchClientOptions,
  type FetchRequestOptions,
  type HttpResponse
} from '@ownclouders/web-client'
import { z } from 'zod'

export type RequestConfig<D = any, S = unknown> = Omit<FetchRequestOptions, 'body'> & {
  /**
   * The request body. Named `data` rather than the core's `body` so that the
   * config this client has always taken keeps working unchanged.
   */
  data?: D
  schema?: S extends z.Schema ? S : never
}

type Resolved<T, S> = HttpResponse<S extends z.Schema ? z.infer<S> : T>

/**
 * Ties a per-request signal to the client-wide one without leaking a listener per request:
 * the caller disposes once the request has settled.
 */
const combineSignals = (clientSignal: AbortSignal, requestSignal?: AbortSignal) => {
  if (!requestSignal) {
    return { signal: clientSignal, dispose: () => undefined }
  }

  const controller = new AbortController()
  const signals = [clientSignal, requestSignal]
  const abort = (source: AbortSignal) => controller.abort(source.reason)
  const listeners = signals.map((source) => {
    const listener = () => abort(source)
    source.addEventListener('abort', listener)
    return () => source.removeEventListener('abort', listener)
  })

  const aborted = signals.find((source) => source.aborted)
  if (aborted) {
    abort(aborted)
  }

  return { signal: controller.signal, dispose: () => listeners.forEach((remove) => remove()) }
}

export class HttpClient {
  private readonly client: FetchClient
  private readonly controller = new AbortController()

  constructor(options: FetchClientOptions = {}) {
    this.client = new FetchClient(options)
  }

  /**
   * Aborts every request this client has in flight. As with the `CancelToken` this
   * replaces, the client is spent afterwards and rejects immediately.
   */
  public cancel(msg?: string): void {
    this.controller.abort(msg ? new DOMException(msg, 'AbortError') : undefined)
  }

  public delete<T = any, D = any, S extends z.Schema | T = T>(
    url: string,
    data?: D,
    config?: RequestConfig<D, S>
  ) {
    return this.send<T, S>(url, { ...config, method: 'DELETE', data })
  }

  public get<T = unknown, D = any, S extends z.Schema | T = T>(
    url: string,
    config?: RequestConfig<D, S>
  ) {
    return this.send<T, S>(url, { ...config, method: 'GET' })
  }

  public head<T = any, D = any, S extends z.Schema | T = T>(
    url: string,
    config?: RequestConfig<D, S>
  ) {
    return this.send<T, S>(url, { ...config, method: 'HEAD' })
  }

  public options<T = any, D = any, S extends z.Schema | T = T>(
    url: string,
    config?: RequestConfig<D, S>
  ) {
    return this.send<T, S>(url, { ...config, method: 'OPTIONS' })
  }

  public patch<T = any, D = any, S extends z.Schema | T = T>(
    url: string,
    data?: D,
    config?: RequestConfig<D, S>
  ) {
    return this.send<T, S>(url, { ...config, method: 'PATCH', data })
  }

  public post<T = any, D = any, S extends z.Schema | T = T>(
    url: string,
    data?: D,
    config?: RequestConfig<D, S>
  ) {
    return this.send<T, S>(url, { ...config, method: 'POST', data })
  }

  public put<T = any, D = any, S extends z.Schema | T = T>(
    url: string,
    data?: D,
    config?: RequestConfig<D, S>
  ) {
    return this.send<T, S>(url, { ...config, method: 'PUT', data })
  }

  public request<T = any, D = any, S extends z.Schema | T = T>(
    config: RequestConfig<D, S> & { url?: string; method?: string }
  ) {
    const { url = '', ...rest } = config
    return this.send<T, S>(url, rest)
  }

  private async send<T, S>(url: string, config: RequestConfig<any, S>): Promise<Resolved<T, S>> {
    const { data, schema, signal, ...rest } = config
    const { signal: combined, dispose } = combineSignals(this.controller.signal, signal)

    try {
      const response = await this.client.request<any>(url, {
        ...rest,
        body: data,
        signal: combined
      })

      if (schema) {
        return { ...response, data: schema.parse(response.data) } as Resolved<T, S>
      }

      return response as Resolved<T, S>
    } finally {
      dispose()
    }
  }
}
