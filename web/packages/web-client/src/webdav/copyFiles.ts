import { SpaceResource } from '../helpers'
import { DAV, DAVRequestOptions } from './client'
import { WebDavOptions } from './types'
import { getWebDavPath } from './utils'

export const CopyFilesFactory = (dav: DAV, options: WebDavOptions) => {
  return {
    copyFiles(
      sourceSpace: SpaceResource,
      { path: sourcePath, fileId: sourceFileId }: { path?: string; fileId?: string },
      targetSpace: SpaceResource,
      {
        path: targetPath,
        parentFolderId,
        name
      }: { path?: string; parentFolderId?: string; name?: string },
      { overwrite, ...opts }: { overwrite?: boolean } & DAVRequestOptions = {}
    ) {
      const sourceWebDavPath = getWebDavPath(sourceSpace, {
        fileId: sourceFileId,
        path: sourcePath
      })

      const targetWebDavPath = getWebDavPath(targetSpace, {
        fileId: parentFolderId,
        path: targetPath,
        name
      })

      return dav.copy(sourceWebDavPath, targetWebDavPath, {
        overwrite: overwrite || false,
        ...opts
      })
    }
  }
}
