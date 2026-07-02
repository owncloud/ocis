import { SpaceResource } from '../helpers'
import { GetFileInfoFactory } from './getFileInfo'
import { WebDavOptions } from './types'
import { DAV, DAVRequestOptions } from './client'
import { ProgressEventCallback } from 'webdav'
import { getWebDavPath } from './utils'

type PutFileContentsOptions = {
  fileName?: string
  parentFolderId?: string
  path?: string
  content?: string | ArrayBuffer
  previousEntityTag?: string
  overwrite?: boolean
  onUploadProgress?: ProgressEventCallback
} & DAVRequestOptions

export const PutFileContentsFactory = (
  dav: DAV,
  getFileInfoFactory: ReturnType<typeof GetFileInfoFactory>,
  options: WebDavOptions
) => {
  return {
    async putFileContents(
      space: SpaceResource,
      {
        fileName,
        path,
        parentFolderId,
        content = '',
        previousEntityTag = '',
        overwrite,
        onUploadProgress = null,
        ...opts
      }: PutFileContentsOptions
    ) {
      const webDavPath = getWebDavPath(space, { fileId: parentFolderId, name: fileName, path })
      const { result } = await dav.put(webDavPath, content, {
        previousEntityTag,
        overwrite,
        onUploadProgress,
        ...opts
      })

      return getFileInfoFactory.getFileInfo(space, {
        fileId: result.headers.get('Oc-Fileid'),
        path
      })
    }
  }
}
