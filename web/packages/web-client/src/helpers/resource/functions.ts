import path, { basename, dirname } from 'path'
import { urlJoin } from '../../utils'
import { DavPermission, DavProperty } from '../../webdav/constants'
import { Resource, SearchResource, TrashResource, WebDavResponseResource } from './types'
import { camelCase } from 'lodash-es'
import { HIDDEN_FILE_EXTENSIONS, PASSWORD_PROTECTED_FOLDER_FILE_EXTENSION } from '../../constants'

const fileExtensions = {
  complex: ['tar.bz2', 'tar.gz', 'tar.xz']
}

export const isTrashResource = (resource: Resource): resource is TrashResource => {
  return Object.hasOwn(resource, 'ddate')
}

export const isSearchResource = (resource: Resource): resource is SearchResource => {
  return Object.hasOwn(resource, 'highlights')
}

export const extractDomSelector = (str: string): string => {
  return str.replace(/[^A-Za-z0-9\-_]/g, '')
}

const extractIdSegment = (id: string, index: number): string => {
  if (!id || typeof id !== 'string') {
    return ''
  }
  return id.indexOf('!') >= 0 ? id.split('!')[index] : ''
}

export const extractStorageId = (id?: string): string => {
  return extractIdSegment(id, 0)
}

export const extractNodeId = (id?: string): string => {
  return extractIdSegment(id, 1)
}

export const extractNameWithoutExtension = (resource?: Resource): string => {
  const extension = resource.extension || ''
  const name = resource.name || ''
  if (!extension.length) {
    return name
  }
  const extensionIndexInName = name.lastIndexOf(`.${extension}`)
  return name.substring(0, extensionIndexInName)
}

export const extractExtensionFromFile = (resource: Resource): string => {
  const name = resource.name
  if (resource.type === 'directory' || resource.isFolder) {
    return ''
  }

  const parts = name.split('.')
  if (parts.length > 2) {
    for (let i = 0; i < parts.length; i++) {
      const possibleExtension = parts.slice(i, parts.length).join('.')
      if (fileExtensions.complex.includes(possibleExtension)) {
        return possibleExtension
      }
    }
  }
  // Fallback if file extension is unknown or no extension
  if (parts.length < 2) {
    return ''
  }
  return parts[parts.length - 1]
}

export const extractParentFolderName = (resource: Resource): string | null => {
  return basename(dirname(resource.path)) || null
}

export const isShareRoot = (resource: Resource) => {
  return typeof resource.isShareRoot === 'function' && resource.isShareRoot()
}

const convertObjectToCamelCaseKeys = (data: Record<string, any>) => {
  if (!data) {
    return data
  }
  const converted: Record<string, any> = {}
  Object.keys(data).forEach((key) => {
    converted[camelCase(key)] = data[key]
  })
  return converted
}

export function buildResource(
  resource: WebDavResponseResource,
  extraPropNames: string[] = []
): Resource {
  const name = resource.props[DavProperty.Name]?.toString() || basename(resource.filename)
  const id = resource.props[DavProperty.FileId]

  const isFolder = resource.type === 'directory'
  let resourcePath: string

  if (resource.filename.startsWith('/files') || resource.filename.startsWith('/space')) {
    resourcePath = resource.filename.split('/').slice(3).join('/')
  } else {
    resourcePath = resource.filename
  }

  if (!resourcePath.startsWith('/')) {
    resourcePath = `/${resourcePath}`
  }

  const extension = extractExtensionFromFile({
    ...resource,
    id,
    name,
    path: resourcePath,
    spaceId: resource.props[DavProperty.SpaceId]
  })

  const lock = resource.props[DavProperty.LockDiscovery]
  let activeLock: { [DavProperty.LockOwner]?: string; [DavProperty.LockTime]?: string }
  let lockOwner: string, lockTime: string
  if (lock) {
    activeLock = lock[DavProperty.ActiveLock]
    lockOwner = activeLock[DavProperty.LockOwner]
    lockTime = activeLock[DavProperty.LockTime]
  }

  let shareTypes: number[] = []
  if (resource.props[DavProperty.ShareTypes]) {
    shareTypes = resource.props[DavProperty.ShareTypes]['share-type']
    if (!Array.isArray(shareTypes)) {
      shareTypes = [shareTypes]
    }
  }

  const extraProps: Record<string, unknown> = {}
  for (const name of extraPropNames || []) {
    // only use the tag name, the namespace is ignored by our WebDav library
    // https://github.com/perry-mitchell/webdav-client/issues/210
    const extraPropName = name.split(':').pop()
    if (resource.props[extraPropName]) {
      extraProps[name] = resource.props[extraPropName]
    }
  }

  const r: Resource = {
    id,
    fileId: id,
    storageId: extractStorageId(id),
    parentFolderId: resource.props[DavProperty.FileParent],
    mimeType: resource.props[DavProperty.MimeType],
    name,
    extension: isFolder ? '' : extension,
    path: resourcePath,
    webDavPath: resource.filename,
    type: isFolder ? 'folder' : resource.type,
    isFolder,
    locked: !!activeLock,
    lockOwner,
    lockTime,
    processing: resource.processing || false,
    mdate: resource.props[DavProperty.LastModifiedDate],
    size: isFolder
      ? resource.props[DavProperty.ContentSize]?.toString() || '0'
      : resource.props[DavProperty.ContentLength]?.toString() || '0',
    permissions: resource.props[DavProperty.Permissions] || '',
    starred: resource.props[DavProperty.IsFavorite] !== 0,
    etag: resource.props[DavProperty.ETag],
    shareTypes,
    privateLink: resource.props[DavProperty.PrivateLink],
    downloadURL: resource.props[DavProperty.DownloadURL],
    remoteItemId: resource.props[DavProperty.RemoteItemId],
    remoteItemPath: resource.props[DavProperty.ShareRoot],
    owner: {
      id: resource.props[DavProperty.OwnerId],
      displayName: resource.props[DavProperty.OwnerDisplayName]
    },
    tags: (resource.props[DavProperty.Tags] || '').toString().split(',').filter(Boolean),
    audio: convertObjectToCamelCaseKeys(resource.props[DavProperty.Audio]),
    location: convertObjectToCamelCaseKeys(resource.props[DavProperty.Location]),
    image: convertObjectToCamelCaseKeys(resource.props[DavProperty.Image]),
    photo: convertObjectToCamelCaseKeys(resource.props[DavProperty.Photo]),
    extraProps,
    spaceId: resource.props[DavProperty.SpaceId],
    canUpload: function () {
      return this.permissions.indexOf(DavPermission.FolderCreateable) >= 0
    },
    canDownload: function () {
      return (
        // TODO: we should later on separate this from the hidden extensions and use some better logic
        !HIDDEN_FILE_EXTENSIONS.includes(this.extension) &&
        this.permissions.indexOf(DavPermission.SecureView) === -1
      )
    },
    canBeDeleted: function () {
      return this.permissions.indexOf(DavPermission.Deletable) >= 0
    },
    canRename: function () {
      return (
        !HIDDEN_FILE_EXTENSIONS.includes(this.extension) &&
        this.permissions.indexOf(DavPermission.Renameable) >= 0
      )
    },
    canShare: function ({ ability }) {
      return (
        !HIDDEN_FILE_EXTENSIONS.includes(this.extension) &&
        ability.can('create-all', 'Share') &&
        this.permissions.indexOf(DavPermission.Shareable) >= 0
      )
    },
    canCreate: function () {
      return this.permissions.indexOf(DavPermission.FolderCreateable) >= 0
    },
    canEditTags: function () {
      return (
        !HIDDEN_FILE_EXTENSIONS.includes(this.extension) &&
        (this.permissions.indexOf(DavPermission.Updateable) >= 0 ||
          this.permissions.indexOf(DavPermission.FileUpdateable) >= 0)
      )
    },
    isMounted: function () {
      return this.permissions.indexOf(DavPermission.Mounted) >= 0
    },
    isReceivedShare: function () {
      return this.permissions.indexOf(DavPermission.Shared) >= 0
    },
    isShareRoot(): boolean {
      return resource.props[DavProperty.ShareRoot]
        ? resource.filename.split('/').length === 3
        : false
    },
    canDeny: function () {
      return this.permissions.indexOf(DavPermission.Deny) >= 0
    },
    getDomSelector: () => extractDomSelector(id)
  }
  Object.defineProperty(r, 'nodeId', {
    get() {
      return extractNodeId(this.id)
    }
  })

  if (resource.props[DavProperty.SignatureAuth]) {
    r.signatureAuth = {
      signature: resource.props[DavProperty.SignatureAuth].signature,
      expiration: new Date(resource.props[DavProperty.SignatureAuth].expiration)
    }
  }

  return r
}

export function buildDeletedResource(resource: WebDavResponseResource): TrashResource {
  const isFolder = resource.type === 'directory'
  const fullName = resource.props[DavProperty.TrashbinOriginalFilename]?.toString()
  const extension = extractExtensionFromFile({ name: fullName, type: resource.type } as Resource)
  const id = path.basename(resource.filename)
  return {
    type: isFolder ? 'folder' : resource.type,
    isFolder,
    ddate: resource.props[DavProperty.TrashbinDeletedDate],
    name: path.basename(fullName),
    extension,
    path: urlJoin(resource.props[DavProperty.TrashbinOriginalLocation], { leadingSlash: true }),
    id,
    parentFolderId: resource.props[DavProperty.FileParent],
    webDavPath: '',
    spaceId: '',
    canUpload: () => false,
    canDownload: () => false,
    canBeDeleted: () => {
      /** FIXME: once https://github.com/owncloud/ocis/issues/3339 gets implemented,
       * we want to add a check if the permission is set.
       * We might to be careful and do an early return true if DavProperty.Permissions is not set
       * as oc10 does not support it.
       **/
      return true
    },
    canBeRestored: function () {
      /** FIXME: once https://github.com/owncloud/ocis/issues/3339 gets implemented,
       * we want to add a check if the permission is set.
       * We might to be careful and do an early return true if DavProperty.Permissions is not set
       * as oc10 does not support it.
       **/
      return true
    },
    canRename: () => false,
    canShare: () => false,
    canCreate: () => false,
    isMounted: () => false,
    isReceivedShare: () => false,
    getDomSelector: () => extractDomSelector(id)
  }
}

export function isPasswordProtectedFolderFileResource(resourceName: string): boolean {
  return resourceName.endsWith(PASSWORD_PROTECTED_FOLDER_FILE_EXTENSION)
}
