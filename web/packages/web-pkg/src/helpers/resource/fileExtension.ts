import { Resource } from '@ownclouders/web-client'
import { ApplicationFileExtension } from '../../apps/types'

/**
 * Checks whether a resource matches a registered application file extension,
 * either by its file extension or by its mime type.
 */
export const isResourceMatchingFileExtension = (
  resource: Resource,
  fileExtension: ApplicationFileExtension
): boolean => {
  if (!resource || !fileExtension) {
    return false
  }

  if (resource.extension && fileExtension.extension) {
    return resource.extension.toLowerCase() === fileExtension.extension.toLowerCase()
  }

  if (resource.mimeType && fileExtension.mimeType) {
    return (
      resource.mimeType.toLowerCase() === fileExtension.mimeType.toLowerCase() ||
      resource.mimeType.split('/')[0].toLowerCase() === fileExtension.mimeType.toLowerCase()
    )
  }

  return false
}

const normalizeMimeType = (mimeType: string) => mimeType.split(';')[0].trim().toLowerCase()

/**
 * Strict check for apps that render the file content directly: the mime type reported by the
 * server must match the registered one. The file extension is not taken into account because
 * the name of a resource (e.g. a mounted share) can differ from the name of the underlying file,
 * while the mime type always reflects the latter.
 */
export const isResourceMatchingFileType = (
  resource: Resource,
  fileExtension: ApplicationFileExtension
): boolean => {
  if (!resource?.mimeType || !fileExtension?.mimeType) {
    return false
  }

  return normalizeMimeType(resource.mimeType) === normalizeMimeType(fileExtension.mimeType)
}
