/**
 * A copy of https://github.com/moxystudio/js-proper-url-join/blob/master/src/index.js
 * but without the query handling.
 */

const urlRegExp = /^(\w+:\/\/[^/?]+)?(.*?)$/

export interface UrlJoinOptions {
  /**
   * Add a leading slash.
   *
   * **Default**: `true`
   */
  leadingSlash?: boolean | 'keep' | undefined
  /**
   * Add a trailing slash.
   *
   * **Default**: `false`
   */
  trailingSlash?: boolean | 'keep' | undefined
}

const normalizeParts = (parts: string[]) =>
  parts
    // Filter non-string or non-numeric values
    .filter((part) => typeof part === 'string' || typeof part === 'number')
    // Convert to strings
    .map((part) => `${part}`)
    // Remove empty parts
    .filter((part) => part)

const parseParts = (parts: string[]) => {
  const partsStr = parts.join('/')
  const [, prefix = '', pathname = ''] = partsStr.match(urlRegExp) || []

  return {
    prefix,
    pathname: {
      parts: pathname.split('/').filter((part) => part !== ''),
      hasLeading: /^\/+/.test(pathname),
      hasTrailing: /\/+$/.test(pathname)
    }
  }
}

const buildUrl = (parsedParts: ReturnType<typeof parseParts>, options: UrlJoinOptions) => {
  const { prefix, pathname } = parsedParts
  const { parts: pathnameParts, hasLeading, hasTrailing } = pathname
  const { leadingSlash, trailingSlash } = options

  const addLeading = leadingSlash === true || (leadingSlash === 'keep' && hasLeading)
  const addTrailing = trailingSlash === true || (trailingSlash === 'keep' && hasTrailing)

  // Start with prefix if not empty (http://google.com)
  let url = prefix

  // Add the parts
  if (pathnameParts.length > 0) {
    if (url || addLeading) {
      url += '/'
    }

    url += pathnameParts.join('/')
  }

  // Add trailing to the end
  if (addTrailing) {
    url += '/'
  }

  // Add leading if URL is still empty
  if (!url && addLeading) {
    url += '/'
  }

  return url
}

export const urlJoin = (...parts: Array<string | UrlJoinOptions>) => {
  const lastArg = parts[parts.length - 1]
  let options: UrlJoinOptions

  // If last argument is an object, then it's the options
  // Note that null is an object, so we verify if is truthy
  if (lastArg && typeof lastArg === 'object') {
    options = lastArg
    parts = parts.slice(0, -1)
  }

  // Parse options
  options = {
    leadingSlash: true,
    trailingSlash: false,
    ...options
  } as UrlJoinOptions

  // Normalize parts before parsing them
  parts = normalizeParts(parts as string[])

  // Split the parts into prefix, pathname
  // (scheme://host)(/pathnameParts.join('/'))
  const parsedParts = parseParts(parts as string[])

  // Finally build the url based on the parsedParts
  return buildUrl(parsedParts, options)
}
