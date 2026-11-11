import type { Configuration, InitOverrideFunction } from './generated'

export interface GraphFactoryOptions {
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
 *
 * The merge goes through `Headers.set()` rather than an object spread so that an override
 * replaces a generated header whatever case either of them used. An object spread would
 * keep both spellings, and the duplicate would later be joined into a single
 * comma-separated value.
 */
export const toInitOverrides =
  (options?: GraphRequestOptions): InitOverrideFunction =>
  async ({ init }) => {
    const headers = new Headers(init.headers)
    Object.entries(options?.headers ?? {}).forEach(([name, value]) => headers.set(name, value))

    return {
      ...(options?.signal && { signal: options.signal }),
      ...(options?.params && { [undeclaredParams]: options.params }),
      headers
    }
  }
