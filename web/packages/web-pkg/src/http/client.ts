import {
  FetchClient,
  type FetchClientOptions,
  type FetchRequestOptions,
  type HttpResponse
} from '@ownclouders/web-client'
import { z } from 'zod'

export type RequestConfig<D = any, S = unknown> = FetchRequestOptions & {
  schema?: S extends z.Schema ? S : never
}

type Resolved<T, S> = HttpResponse<S extends z.Schema ? z.infer<S> : T>

export class HttpClient {
  private readonly client: FetchClient

  constructor(options: FetchClientOptions = {}) {
    this.client = new FetchClient(options)
  }

  public delete<T = any, D = any, S extends z.Schema | T = T>(
    url: string,
    data?: D,
    config?: RequestConfig<D, S>
  ) {
    return this.send<T, S>(url, { ...config, method: 'DELETE', body: data })
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
    return this.send<T, S>(url, { ...config, method: 'PATCH', body: data })
  }

  public post<T = any, D = any, S extends z.Schema | T = T>(
    url: string,
    data?: D,
    config?: RequestConfig<D, S>
  ) {
    return this.send<T, S>(url, { ...config, method: 'POST', body: data })
  }

  public put<T = any, D = any, S extends z.Schema | T = T>(
    url: string,
    data?: D,
    config?: RequestConfig<D, S>
  ) {
    return this.send<T, S>(url, { ...config, method: 'PUT', body: data })
  }

  public request<T = any, D = any, S extends z.Schema | T = T>(
    config: RequestConfig<D, S> & { url?: string; method?: string }
  ) {
    const { url = '', ...rest } = config
    return this.send<T, S>(url, rest)
  }

  private async send<T, S>(
    url: string,
    config: RequestConfig<any, S>
  ): Promise<Resolved<T, S>> {
    const response = await this.client.request<any>(url, config)

    if (config?.schema) {
      return { ...response, data: config.schema.parse(response.data) } as Resolved<T, S>
    }

    return response as Resolved<T, S>
  }
}
