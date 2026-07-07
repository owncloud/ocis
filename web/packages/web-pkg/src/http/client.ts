import axios, {
  AxiosError,
  AxiosInstance,
  AxiosRequestConfig,
  AxiosResponse,
  CancelTokenSource,
  InternalAxiosRequestConfig
} from 'axios'
import merge from 'lodash-es/merge'
import { z } from 'zod'

export type RequestConfig<D, S> = AxiosRequestConfig<D> & {
  schema?: S extends z.Schema ? S : never
}
export class HttpClient {
  private readonly instance: AxiosInstance
  private readonly cancelToken: CancelTokenSource

  constructor({
    config,
    requestInterceptor,
    responseInterceptor
  }: {
    config?: AxiosRequestConfig
    requestInterceptor?: (
      value: InternalAxiosRequestConfig<any>
    ) => InternalAxiosRequestConfig<any> | Promise<InternalAxiosRequestConfig<any>>
    responseInterceptor?: [
      (response: AxiosResponse<any>) => AxiosResponse<any> | Promise<AxiosResponse<any>>,
      (error: AxiosError<any>) => AxiosResponse<any> | Promise<AxiosError<any>>
    ]
  } = {}) {
    this.cancelToken = axios.CancelToken.source()
    this.instance = axios.create(config)

    if (requestInterceptor) {
      this.instance.interceptors.request.use(requestInterceptor)
    }

    if (responseInterceptor) {
      this.instance.interceptors.response.use(responseInterceptor[0], responseInterceptor[1])
    }
  }

  public cancel(msg?: string): void {
    this.cancelToken.cancel(msg)
  }

  public async delete<T = any, D = any, S extends z.Schema | T = T>(
    url: string,
    data?: D,
    config?: RequestConfig<D, S>
  ) {
    return await this.internalRequestWithData('delete', url, data, config)
  }

  public get<T = unknown, D = any, S extends z.Schema | T = T>(
    url: string,
    config?: RequestConfig<D, S>
  ) {
    return this.internalRequest('get', url, config)
  }

  public head<T = any, D = any, S extends z.Schema | T = T>(
    url: string,
    config?: RequestConfig<D, S>
  ) {
    return this.internalRequest('head', url, config)
  }

  public options<T = any, D = any, S extends z.Schema | T = T>(
    url: string,
    config?: RequestConfig<D, S>
  ) {
    return this.internalRequest('options', url, config)
  }

  public patch<T = any, D = any, S extends z.Schema | T = T>(
    url: string,
    data?: D,
    config?: RequestConfig<D, S>
  ) {
    return this.internalRequestWithData('patch', url, data, config)
  }

  public post<T = any, D = any, S extends z.Schema | T = T>(
    url: string,
    data?: D,
    config?: RequestConfig<D, S>
  ) {
    return this.internalRequestWithData('post', url, data, config)
  }

  public put<T = any, D = any, S extends z.Schema | T = T>(
    url: string,
    data?: D,
    config?: RequestConfig<D, S>
  ) {
    return this.internalRequestWithData('put', url, data, config)
  }

  public async request<T = any, D = any, S extends z.Schema | T = T>(config: RequestConfig<D, S>) {
    const response = await this.instance.request<S, AxiosResponse<S>, D>(
      this.obtainConfig<D>(config)
    )
    return this.processResponse(response, config)
  }

  private obtainConfig<D = any>(config?: AxiosRequestConfig): AxiosRequestConfig<D> {
    return merge({ cancelToken: this.cancelToken.token }, config)
  }

  private processResponse<T, S extends z.Schema | T = T>(
    response: AxiosResponse<T>,
    config?: RequestConfig<any, S>
  ): AxiosResponse<S extends z.Schema ? z.infer<S> : T> {
    if (config?.schema) {
      const data = config.schema.parse(response.data)
      return { ...response, data } as AxiosResponse<S extends z.Schema ? z.infer<S> : T>
    }

    return response as AxiosResponse<S extends z.Schema ? z.infer<S> : T>
  }

  private async internalRequest<T = any, D = any, S extends z.Schema | T = T>(
    method: 'delete' | 'get' | 'head' | 'options',
    url: string,
    config: RequestConfig<D, S>
  ) {
    const response = await this.instance[method]<S, AxiosResponse<S>, D>(
      url,
      this.obtainConfig<D>(config)
    )

    return this.processResponse(response, config)
  }

  private async internalRequestWithData<T = any, D = any, S extends z.Schema | T = T>(
    method: 'post' | 'put' | 'patch' | 'delete',
    url: string,
    data: D,
    config: RequestConfig<D, S>
  ) {
    const response = await this.instance[method]<S, AxiosResponse<S>, D>(
      url,
      data,
      this.obtainConfig<D>(config)
    )

    return this.processResponse(response, config)
  }
}
