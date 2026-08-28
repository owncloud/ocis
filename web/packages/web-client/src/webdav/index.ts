import { Headers } from 'webdav'
import { FetchClient } from '../http'
import { WebDAV } from './types'
import { CopyFilesFactory } from './copyFiles'
import { CreateFolderFactory } from './createFolder'
import { GetFileContentsFactory } from './getFileContents'
import { GetFileInfoFactory } from './getFileInfo'
import { GetFileUrlFactory } from './getFileUrl'
import { GetPublicFileUrlFactory } from './getPublicFileUrl'
import { ListFilesFactory } from './listFiles'
import { MoveFilesFactory } from './moveFiles'
import { PutFileContentsFactory } from './putFileContents'
import { DeleteFileFactory } from './deleteFile'
import { RestoreFileFactory } from './restoreFile'
import { RestoreFileVersionFactory } from './restoreFileVersion'
import { ClearTrashBinFactory } from './clearTrashBin'
import { SearchFactory } from './search'
import { GetPathForFileIdFactory } from './getPathForFileId'
import { DAV } from './client/dav'
import { ListFileVersionsFactory } from './listFileVersions'
import { SetFavoriteFactory } from './setFavorite'
import { ListFavoriteFilesFactory } from './listFavoriteFiles'
import { shouldResponseTriggerMaintenance } from '../helpers/maintenance'

export * from './constants'
export * from './types'

export type { ListFilesOptions, ListFilesResult } from './listFiles'
export type { GetFileContentsResponse } from './getFileContents'

export const webdav = (
  baseURI: string,
  onSetMaintenance: (value: boolean) => void,
  headers?: () => Headers
): WebDAV => {
  const httpClient = new FetchClient({
    ...(headers && { headers }),
    onResponse: ({ response, status, requestUrl }) => {
      if (response?.ok) {
        onSetMaintenance(false)
        return
      }
      onSetMaintenance(shouldResponseTriggerMaintenance(status, requestUrl))
    }
  })

  const options = { httpClient, baseUrl: baseURI, headers }

  const dav = new DAV({ baseUrl: baseURI, headers, onSetMaintenance })
  const registerExtraProp = (name: string) => {
    dav.extraProps.push(name)
  }

  const pathForFileIdFactory = GetPathForFileIdFactory(dav, options)
  const { getPathForFileId } = pathForFileIdFactory

  const listFilesFactory = ListFilesFactory(dav, pathForFileIdFactory, options)
  const { listFiles } = listFilesFactory

  const getFileInfoFactory = GetFileInfoFactory(listFilesFactory, options)
  const { getFileInfo } = getFileInfoFactory

  const { createFolder } = CreateFolderFactory(dav, getFileInfoFactory, options)
  const getFileContentsFactory = GetFileContentsFactory(dav, options)
  const { getFileContents } = getFileContentsFactory
  const { putFileContents } = PutFileContentsFactory(dav, getFileInfoFactory, options)

  const { getFileUrl, revokeUrl } = GetFileUrlFactory(dav, getFileContentsFactory, options)
  const { getPublicFileUrl } = GetPublicFileUrlFactory(dav, options)

  const { copyFiles } = CopyFilesFactory(dav, options)
  const { moveFiles } = MoveFilesFactory(dav, options)

  const { deleteFile } = DeleteFileFactory(dav, options)
  const { restoreFile } = RestoreFileFactory(dav, options)

  const { listFileVersions } = ListFileVersionsFactory(dav, options)
  const { restoreFileVersion } = RestoreFileVersionFactory(dav, options)

  const { clearTrashBin } = ClearTrashBinFactory(dav, options)

  const { search } = SearchFactory(dav, options)

  const { listFavoriteFiles } = ListFavoriteFilesFactory(dav, options)
  const { setFavorite } = SetFavoriteFactory(dav, options)

  return {
    copyFiles,
    createFolder,
    deleteFile,
    restoreFile,
    restoreFileVersion,
    getFileContents,
    getFileInfo,
    getFileUrl,
    getPublicFileUrl,
    getPathForFileId,
    listFiles,
    listFileVersions,
    moveFiles,
    putFileContents,
    revokeUrl,
    clearTrashBin,
    search,
    listFavoriteFiles,
    setFavorite,

    registerExtraProp
  }
}
