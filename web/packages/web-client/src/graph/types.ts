import type { AxiosInstance } from 'axios'
import type { Configuration } from './generated'

export interface GraphFactoryOptions {
  axiosClient: AxiosInstance
  config: Configuration
}

export interface GraphRequestOptions {
  headers?: Record<string, string>
  params?: Record<string, string>
  signal?: AbortSignal
}
