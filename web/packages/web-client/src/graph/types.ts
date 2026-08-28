import type { Configuration, InitOverrideFunction } from './generated'
import type { FetchClient } from '../http'

export interface GraphFactoryOptions {
  httpClient: FetchClient
  config: Configuration
}

export interface GraphRequestOptions {
  headers?: Record<string, string>
  signal?: AbortSignal
  /**
   * Query parameters the libre-graph spec does not declare, such as the `template` of a
   * newly created drive. Declared parameters belong in the generated request object.
   */
  params?: Record<string, string>
}

/**
 * Channel for {@link GraphRequestOptions.params}. The generated client assembles the URL
 * before it applies `initOverrides`, so undeclared query parameters cannot reach it there.
 * They ride along on the `RequestInit` instead and are unpacked by our `fetchApi` bridge.
 */
export const undeclaredParams = Symbol('graph.undeclaredParams')

/**
 * Adapts our request options to the generated client's `initOverrides`.
 *
 * A plain object would be spread shallowly over the generated `RequestInit`, replacing
 * its headers wholesale and dropping `Content-Type`. A function receives the built init
 * and can merge instead.
 */
export const toInitOverrides =
  (options?: GraphRequestOptions): InitOverrideFunction =>
  async ({ init }) => ({
    ...(options?.signal && { signal: options.signal }),
    ...(options?.params && { [undeclaredParams]: options.params }),
    headers: { ...(init.headers as Record<string, string>), ...(options?.headers ?? {}) }
  })
