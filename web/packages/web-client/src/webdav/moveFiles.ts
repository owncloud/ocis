import { SpaceResource } from '../helpers'
import { WebDavOptions } from './types'
import { DAV, DAVRequestOptions } from './client'
import { getWebDavPath } from './utils'

export const MoveFilesFactory = (dav: DAV, options: WebDavOptions) => {
  return {
    moveFiles(
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

      return dav.move(sourceWebDavPath, targetWebDavPath, {
        overwrite: overwrite || false,
        ...opts
      })
    }
  }
}
