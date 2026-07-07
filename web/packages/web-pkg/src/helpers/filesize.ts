import { filesize, FilesizeOptions } from 'filesize'
import { getLocaleFromLanguage } from './locale'

const mb = 1048576

/**
 * Returns formatted size
 *
 * @param {Number|String} size Unformatted size
 * @param {String} currentLanguage
 * @returns {String} formatted size
 */
export const formatFileSize = (size: number | string, currentLanguage: string) => {
  const parsedSize = typeof size === 'string' ? parseInt(size) : size
  if (parsedSize < 0) {
    return '--'
  }

  if (isNaN(parsedSize)) {
    return '?'
  }

  return filesize<FilesizeOptions>(parsedSize, {
    round: parsedSize < mb ? 0 : 1,
    locale: getLocaleFromLanguage(currentLanguage),
    base: 10,
    output: 'string'
  })
}
