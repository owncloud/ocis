/**
 * A native `Headers` that also answers to bracket access.
 *
 * `.get('etag')` is the shape a fetch `Response` produces; `['etag']` is the shape axios
 * produced and what consumers outside this repo were written against. Supporting both
 * keeps the axios removal invisible to them.
 */
export interface HttpHeaders extends Headers {
  [key: string]: any
}

/**
 * Views a native `Headers` as {@link HttpHeaders}.
 *
 * Enumeration (`Object.keys`, spread, `JSON.stringify`) yields the lowercase header names,
 * matching what axios put on its response headers object. Bracket access is
 * case-insensitive, which axios's was not, and resolves to `undefined` rather than `null`
 * for an absent header — again what a plain object did.
 */
export const httpHeaders = (headers: Headers): HttpHeaders =>
  new Proxy(headers, {
    get(target, property) {
      const value = Reflect.get(target, property, target)
      if (typeof value === 'function') {
        // methods need the real Headers as their receiver, not the proxy
        return value.bind(target)
      }
      if (value !== undefined || typeof property !== 'string') {
        return value
      }
      return target.get(property) ?? undefined
    },
    has(target, property) {
      return Reflect.has(target, property) || (typeof property === 'string' && target.has(property))
    },
    ownKeys(target) {
      // lowercased explicitly: browsers do it in the iterator, jsdom does not
      return [...target.keys()].map((name) => name.toLowerCase())
    },
    getOwnPropertyDescriptor(target, property) {
      if (typeof property === 'string' && target.has(property)) {
        return { value: target.get(property), enumerable: true, configurable: true }
      }
      return Reflect.getOwnPropertyDescriptor(target, property)
    }
  }) as HttpHeaders
