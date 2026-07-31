import { urlJoin } from '@ownclouders/web-client'

/**
 * Return all absolute parent paths.
 *
 * For example if passing in "a/b/c" it will return
 * ["/a/b", "/a", ""]

 * If an empty string or "/" is passed in, an empty array is returned.
 *
 * @param {String} path path to process
 * @param {Boolean} includeCurrent whether to include the current path (with leading slash)
 * @return {Array.<String>} parent paths
 */
export function getParentPaths(path = '', includeCurrent = false) {
  // remove potential leading and trailing slash from current path (so that the resulting array doesn't start with an empty string).
  // then reintroduce the leading slash, because we know that we need it.
  const s = urlJoin(path, {
    leadingSlash: true
  })
  if (s === '/') {
    return []
  }

  const paths: string[] = []
  const sections = s.split('/')

  if (includeCurrent) {
    paths.push(s)
  }

  sections.pop()
  while (sections.length > 0) {
    if (!sections.join('/')) {
      sections.pop()
      continue
    }
    paths.push(sections.join('/'))
    sections.pop()
  }

  return paths
}
