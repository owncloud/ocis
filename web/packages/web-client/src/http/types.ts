export type ResponseType = 'json' | 'text' | 'blob' | 'arraybuffer' | 'none'

export interface OnResponseArgs {
  /** null on a transport-level failure */
  response: Response | null
  /** 500 on a transport-level failure */
  status: number
  /** the URL exactly as the caller passed it, before any baseUrl join */
  requestUrl: string
}

export interface FetchClientOptions {
  baseUrl?: string
  /** headers that never change for the lifetime of the client */
  staticHeaders?: Record<string, string>
  /** evaluated per request */
  headers?: () => Record<string, string>
  /** invoked for every outcome, before any throw */
  onResponse?: (args: OnResponseArgs) => void
}

export interface FetchRequestOptions {
  method?: string
  headers?: Record<string, string>
  params?: Record<string, string | number | boolean>
  /** JSON-encoded unless it is already a BodyInit */
  body?: unknown
  responseType?: ResponseType
  signal?: AbortSignal
  /** default true; false returns the envelope for non-2xx instead of throwing */
  throwOnError?: boolean
}

export interface HttpResponse<T = unknown> {
  data: T
  status: number
  statusText: string
  headers: Headers
}
