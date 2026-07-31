import { Permission, User } from '../../graph/generated'
import {
  Ability,
  extractDomSelector,
  extractExtensionFromFile,
  extractNodeId,
  Resource
} from '../resource'
import {
  isPersonalSpaceResource,
  isPublicSpaceResource,
  PublicSpaceResource,
  ShareSpaceResource,
  SpaceMember,
  SpaceResource
} from './types'

import { DavProperty } from '../../webdav/constants'
import { buildWebDavPublicPath, buildWebDavOcmPath } from '../publicLink'
import { urlJoin } from '../../utils'
import { Drive, DriveItem } from '@ownclouders/web-client/graph/generated'
import { GraphSharePermission, ShareRole } from '../share'

export function buildWebDavSpacesPath(storageId: string, path?: string) {
  return urlJoin('spaces', storageId, path, {
    leadingSlash: true
  })
}

export function buildWebDavSpacesTrashPath(storageId: string, path = '') {
  return urlJoin('spaces', 'trash-bin', storageId, path, {
    leadingSlash: true
  })
}

export function getRelativeSpecialFolderSpacePath(space: SpaceResource, type: 'image' | 'readme') {
  const typeMap = { image: 'spaceImageData', readme: 'spaceReadmeData' } as const
  const specialProp = space[typeMap[type]]
  if (!specialProp) {
    return ''
  }
  const webDavPathComponents = decodeURI(specialProp.webDavUrl).split('/')
  const idComponent = webDavPathComponents.find((c) => c.startsWith(space.id))
  if (!idComponent) {
    return ''
  }
  return webDavPathComponents.slice(webDavPathComponents.indexOf(idComponent) + 1).join('/')
}

// although roles are a loose concept, we can assume that there are always space managers
export function getSpaceManagers(space: SpaceResource) {
  return Object.values(space.members).filter(({ permissions }) =>
    // delete permissions implies that the user/group is a manager
    permissions.includes(GraphSharePermission.deletePermissions)
  )
}

export type PublicLinkType = 'ocm' | 'public-link'
export function buildPublicSpaceResource(
  data: any & { publicLinkType: PublicLinkType }
): PublicSpaceResource {
  const publicLinkPassword = data.publicLinkPassword

  const fileId = data.props?.[DavProperty.FileId]
  const publicLinkItemType = data.props?.[DavProperty.PublicLinkItemType]
  const publicLinkPermission = data.props?.[DavProperty.PublicLinkPermission]
  const publicLinkExpiration = data.props?.[DavProperty.PublicLinkExpiration]
  const publicLinkShareDate = data.props?.[DavProperty.PublicLinkShareDate]
  const publicLinkShareOwner = data.props?.[DavProperty.PublicLinkShareOwner]
  const publicLinkShareOwnerDisplayName = data.props?.[DavProperty.OwnerDisplayName]

  let driveAlias
  let webDavPath
  if (data.publicLinkType === 'ocm') {
    driveAlias = `ocm/${data.id}`
    webDavPath = buildWebDavOcmPath(data.id)
  } else {
    driveAlias = `public/${data.id}`
    webDavPath = buildWebDavPublicPath(data.id)
  }

  return Object.assign(
    buildSpace(
      {
        ...data,
        driveType: 'public',
        driveAlias,
        webDavPath
      },
      {}
    ),
    {
      ...(fileId && { fileId }),
      ...(publicLinkPassword && { publicLinkPassword }),
      ...(publicLinkItemType && { publicLinkItemType }),
      ...(publicLinkPermission && { publicLinkPermission: parseInt(publicLinkPermission) }),
      ...(publicLinkExpiration && { publicLinkExpiration }),
      ...(publicLinkShareDate && { publicLinkShareDate }),
      ...(publicLinkShareOwner && { publicLinkShareOwner }),
      ...(publicLinkShareOwnerDisplayName && { publicLinkShareOwnerDisplayName })
    }
  )
}

export function buildShareSpaceResource({
  driveAliasPrefix,
  id,
  shareName,
  serverUrl
}: {
  driveAliasPrefix: 'share' | 'ocm-share'
  id: string
  shareName: string
  serverUrl: string
}): ShareSpaceResource {
  const space = buildSpace(
    {
      id,
      driveAlias: `${driveAliasPrefix}/${shareName}`,
      driveType: 'share',
      name: shareName,
      serverUrl
    },
    {}
  ) as ShareSpaceResource
  space.rename = (newName: string) => {
    space.driveAlias = `${driveAliasPrefix}/${newName}`
    space.name = newName
  }
  return space
}

export function buildSpace(
  data: Drive & {
    path?: string
    serverUrl?: string
    webDavPath?: string
    webDavTrashPath?: string
  },
  graphRoles: Record<string, ShareRole>
): SpaceResource {
  let spaceImageData: DriveItem, spaceReadmeData: DriveItem
  if (data.special) {
    spaceImageData = data.special.find((el) => el.specialFolder.name === 'image')
    spaceReadmeData = data.special.find((el) => el.specialFolder.name === 'readme')

    if (spaceImageData) {
      spaceImageData.webDavUrl = decodeURI(spaceImageData.webDavUrl)
    }

    if (spaceReadmeData) {
      spaceReadmeData.webDavUrl = decodeURI(spaceReadmeData.webDavUrl)
    }
  }

  const disabled = data.root?.deleted?.state === 'trashed'
  const webDavPath = urlJoin(data.webDavPath || buildWebDavSpacesPath(data.id), {
    leadingSlash: true
  })
  const webDavUrl = urlJoin(data.serverUrl, 'dav', webDavPath)
  const webDavTrashPath = urlJoin(data.webDavTrashPath || buildWebDavSpacesTrashPath(data.id), {
    leadingSlash: true
  })
  const webDavTrashUrl = urlJoin(data.serverUrl, 'dav', webDavTrashPath)

  const members = data.root?.permissions?.reduce<Record<string, SpaceMember>>((acc, p) => {
    acc[(p.grantedToV2.user || p.grantedToV2.group).id] = {
      grantedTo: p.grantedToV2,
      permissions: getPermissionsFromGraphPermission(p, graphRoles),
      roleId: p.roles?.[0]
    }
    return acc
  }, {})

  const s = {
    id: data.id,
    fileId: data.id,
    storageId: data.id,
    mimeType: '',
    name: data.name,
    description: data.description,
    extension: '',
    path: '/',
    webDavPath,
    webDavTrashPath,
    driveAlias: data.driveAlias,
    driveType: data.driveType,
    type: 'space',
    isFolder: true,
    mdate: data.lastModifiedDateTime,
    size: data.quota?.used,
    tags: [] as string[],
    permissions: '',
    starred: false,
    etag: '',
    shareTypes: [] as number[],
    privateLink: data.webUrl,
    downloadURL: '',
    owner: data.owner?.user,
    disabled,
    root: data.root,
    spaceQuota: data.quota,
    members: members || {},
    spaceImageData,
    spaceReadmeData,
    spaceId: data.id,
    canUpload: function ({ user }: { user?: User } = {}): boolean {
      if (isPersonalSpaceResource(this) && this.isOwner(user)) {
        return true
      }
      return getPermissionsForSpaceMember(this, user).includes(GraphSharePermission.createUpload)
    },
    canDownload: function () {
      return true
    },
    canBeDeleted: function ({ user, ability }: { user?: User; ability?: Ability } = {}) {
      if (!this.disabled) {
        return false
      }
      if (ability?.can('delete-all', 'Drive')) {
        return true
      }
      // FIXME: server permissions are a mess currently: https://github.com/owncloud/ocis/issues/9862
      return getPermissionsForSpaceMember(this, user).includes(
        GraphSharePermission.deletePermissions
      )
    },
    canRename: function ({ user, ability }: { user?: User; ability?: Ability } = {}) {
      if (ability?.can('update-all', 'Drive')) {
        return true
      }
      // FIXME: server permissions are a mess currently: https://github.com/owncloud/ocis/issues/9862
      return getPermissionsForSpaceMember(this, user).includes(
        GraphSharePermission.deletePermissions
      )
    },
    canEditDescription: function ({ user, ability }: { user?: User; ability?: Ability } = {}) {
      if (ability?.can('update-all', 'Drive')) {
        return true
      }
      // FIXME: server permissions are a mess currently: https://github.com/owncloud/ocis/issues/9862
      return getPermissionsForSpaceMember(this, user).includes(
        GraphSharePermission.deletePermissions
      )
    },
    canRestore: function ({ user, ability }: { user?: User; ability?: Ability } = {}) {
      if (!this.disabled) {
        return false
      }
      if (ability?.can('update-all', 'Drive')) {
        return true
      }
      // FIXME: server permissions are a mess currently: https://github.com/owncloud/ocis/issues/9862
      return getPermissionsForSpaceMember(this, user).includes(
        GraphSharePermission.deletePermissions
      )
    },
    canDisable: function ({ user, ability }: { user?: User; ability?: Ability } = {}) {
      if (this.disabled) {
        return false
      }
      if (ability?.can('delete-all', 'Drive')) {
        return true
      }
      // FIXME: server permissions are a mess currently: https://github.com/owncloud/ocis/issues/9862
      return getPermissionsForSpaceMember(this, user).includes(
        GraphSharePermission.deletePermissions
      )
    },
    canShare: function ({ user }: { user?: User } = {}) {
      return getPermissionsForSpaceMember(this, user).includes(
        GraphSharePermission.createPermissions
      )
    },
    canEditImage: function ({ user }: { user?: User } = {}) {
      if (this.disabled) {
        return false
      }
      // FIXME: server permissions are a mess currently: https://github.com/owncloud/ocis/issues/9862
      return getPermissionsForSpaceMember(this, user).includes(
        GraphSharePermission.deletePermissions
      )
    },
    canEditReadme: function ({ user }: { user?: User } = {}) {
      if (this.disabled) {
        return false
      }
      // FIXME: server permissions are a mess currently: https://github.com/owncloud/ocis/issues/9862
      return getPermissionsForSpaceMember(this, user).includes(
        GraphSharePermission.deletePermissions
      )
    },
    canRestoreFromTrashbin: function ({ user }: { user?: User } = {}) {
      return getPermissionsForSpaceMember(this, user).includes(GraphSharePermission.updateDeleted)
    },
    canDeleteFromTrashBin: function ({ user }: { user?: User } = {}) {
      // FIXME: server permissions are a mess currently: https://github.com/owncloud/ocis/issues/9862
      return getPermissionsForSpaceMember(this, user).includes(
        GraphSharePermission.deletePermissions
      )
    },
    canListVersions: function ({ user }: { user?: User } = {}) {
      if (isPersonalSpaceResource(this) && this.isOwner(user)) {
        return true
      }
      return getPermissionsForSpaceMember(this, user).includes(GraphSharePermission.readVersions)
    },
    canCreate: function () {
      return true
    },
    canEditTags: function () {
      return false
    },
    isMounted: function () {
      return true
    },
    isReceivedShare: function () {
      return false
    },
    isShareRoot: function () {
      return ['share', 'mountpoint', 'public'].includes(data.driveType)
    },
    canDeny: () => false,
    getDomSelector: () => extractDomSelector(data.id),
    getDriveAliasAndItem({ path }: Resource): string {
      return urlJoin(this.driveAlias, path, {
        leadingSlash: false
      })
    },
    getWebDavUrl({ path }: { path: string }): string {
      return urlJoin(webDavUrl, path)
    },
    getWebDavTrashUrl({ path }: { path: string }): string {
      return urlJoin(webDavTrashUrl, path)
    },
    isMember(user: User): boolean {
      if (isPublicSpaceResource(this)) {
        return false
      }
      if (this.isOwner(user) || !!this.members[user.id]) {
        return true
      }
      return user.memberOf?.some((group) => !!this.members[group.id])
    },
    isOwner(user: User): boolean {
      return user?.id === this.owner?.id
    }
  } satisfies SpaceResource
  Object.defineProperty(s, 'nodeId', {
    get() {
      return extractNodeId(this.id)
    }
  })
  return s
}

// build a space image resource based on a given space by its spaceImageData
export function buildSpaceImageResource(space: SpaceResource): Resource {
  return {
    id: space.spaceImageData.id,
    name: space.spaceImageData.name,
    etag: space.spaceImageData.eTag,
    extension: extractExtensionFromFile({ name: space.spaceImageData.name } as Resource),
    mimeType: space.spaceImageData.file.mimeType,
    type: 'file',
    webDavPath: urlJoin(space.webDavPath, '.space', space.spaceImageData.name),
    canDownload: () => true
  } as Resource
}

export function getPermissionsForSpaceMember(space: SpaceResource, user: User) {
  const permissions: string[] = []

  // FIXME: user should always be given, adjust `can...` functions in SpaceResource
  if (!user) {
    return permissions
  }

  // gather permissions from direct user membership
  const member = space.members[user.id]
  if (member) {
    permissions.push(...member.permissions)
  }

  // gather permissions from indirect group membership(s)
  user.memberOf?.forEach((group) => {
    const member = space.members[group.id]
    if (member) {
      permissions.push(...member.permissions)
    }
  })

  return [...new Set(permissions)]
}

/**
 * Get array of permissions from a given graph permission object. If it has '@libre.graph.permissions.actions',
 * then no role exists for this set of permissions. Otherwise, the role is found in the graphRoles array.
 */
function getPermissionsFromGraphPermission(
  permission: Permission,
  graphRoles: Record<string, ShareRole>
): string[] {
  if (permission['@libre.graph.permissions.actions']) {
    return permission['@libre.graph.permissions.actions']
  }
  const role = graphRoles[permission.roles?.[0]]
  if (role) {
    const permissions = role.rolePermissions.find(
      ({ condition }) => condition === 'exists @Resource.Root'
    )
    return permissions?.allowedResourceActions || []
  }
  return []
}
