import { Resource } from '../resource'
import {
  ShareResource,
  OutgoingShareResource,
  IncomingShareResource,
  CollaboratorShare,
  GraphSharePermission,
  LinkShare,
  ShareRole
} from './types'
import { extractDomSelector, extractExtensionFromFile, extractStorageId } from '../resource'
import { ShareTypes } from './type'
import { buildWebDavSpacesPath } from '../space'
import { DriveItem, Identity, Permission, UnifiedRoleDefinition, User } from '../../graph/generated'
import { urlJoin } from '../../utils'
import { uniq } from 'lodash-es'

export const isShareResource = (resource: Resource): resource is ShareResource => {
  return Object.hasOwn(resource, 'sharedWith')
}

export const isOutgoingShareResource = (resource: Resource): resource is OutgoingShareResource => {
  return isShareResource(resource) && resource.outgoing
}

export const isIncomingShareResource = (resource: Resource): resource is IncomingShareResource => {
  return isShareResource(resource) && !resource.outgoing
}

export const isCollaboratorShare = (
  share: CollaboratorShare | LinkShare
): share is CollaboratorShare => {
  return Object.hasOwn(share, 'role')
}

export const isLinkShare = (share: CollaboratorShare | LinkShare): share is LinkShare => {
  return !Object.hasOwn(share, 'role')
}

export const getShareResourceRoles = ({
  driveItem,
  graphRoles
}: {
  driveItem: DriveItem
  graphRoles: Record<string, ShareRole>
}) => {
  return driveItem.remoteItem?.permissions.reduce<UnifiedRoleDefinition[]>((acc, permission) => {
    permission.roles?.forEach((roleId) => {
      const role = graphRoles[roleId]
      if (role && !acc.some(({ id }) => id === role.id)) {
        acc.push(role)
      }
    })

    return acc
  }, [])
}

export const getShareResourcePermissions = ({
  driveItem,
  shareRoles
}: {
  driveItem: DriveItem
  shareRoles: UnifiedRoleDefinition[]
}): GraphSharePermission[] => {
  if (!shareRoles.length) {
    // the server lists plain permissions if it doesn't find a corresponding role
    const permissions = driveItem.remoteItem?.permissions.reduce<GraphSharePermission[]>(
      (acc, permission) => {
        const permissions = permission['@libre.graph.permissions.actions'] as GraphSharePermission[]
        if (permissions) {
          acc.push(...permissions)
        }

        return acc
      },
      []
    )
    return [...new Set(permissions)]
  }

  const permissions = shareRoles.reduce((acc, role) => {
    role.rolePermissions.forEach((permission) => {
      acc.push(...permission.allowedResourceActions)
    })
    return acc
  }, [])

  return [...new Set(permissions)]
}

export function buildIncomingShareResource({
  driveItem,
  graphRoles,
  serverUrl
}: {
  driveItem: DriveItem
  graphRoles: Record<string, ShareRole>
  serverUrl: string
}): IncomingShareResource {
  const resourceName = driveItem.name || driveItem.remoteItem.name
  const storageId = extractStorageId(driveItem.remoteItem.id)

  const sharedWith = driveItem.remoteItem.permissions.map((permission) => {
    const { grantedToV2 } = permission
    const identity = grantedToV2.group || grantedToV2.user
    return {
      ...identity,
      shareType: getShareTypeFromPermission(permission)
    }
  })

  const sharedBy = driveItem.remoteItem.permissions.reduce<Identity[]>((acc, permission) => {
    const sharedBy = permission.invitation.invitedBy.user
    if (!acc.some(({ id }) => id === sharedBy.id)) {
      acc.push(sharedBy)
    }
    return acc
  }, [])

  let shareTypes = uniq(driveItem.remoteItem.permissions.map(getShareTypeFromPermission))
  const isExternal = sharedBy.some((s) => s['@libre.graph.userType'] === 'Federated')
  if (isExternal) {
    shareTypes = [ShareTypes.remote.value]
  }

  const shareRoles = getShareResourceRoles({ driveItem, graphRoles })
  const sharePermissions = getShareResourcePermissions({ driveItem, shareRoles })

  const resource: IncomingShareResource = {
    id: driveItem.id,
    remoteItemId: driveItem.remoteItem.id,
    driveId: driveItem.parentReference?.driveId,
    path: '/',
    name: resourceName,
    fileId: driveItem.remoteItem.id,
    size: driveItem.size,
    storageId,
    parentFolderId: driveItem.parentReference?.id,
    sdate: driveItem.remoteItem.permissions[0].createdDateTime,
    tags: [],
    webDavPath: buildWebDavSpacesPath(driveItem.remoteItem.id, '/'),
    sharedBy,
    owner: driveItem.remoteItem.createdBy?.user,
    sharedWith,
    shareTypes,
    isFolder: !!driveItem.folder,
    type: !!driveItem.folder ? 'folder' : 'file',
    mimeType: driveItem.file?.mimeType || 'httpd/unix-directory',
    mdate: driveItem.lastModifiedDateTime
      ? new Date(driveItem.lastModifiedDateTime).toUTCString()
      : undefined,
    syncEnabled: driveItem['@client.synchronize'],
    hidden: driveItem['@UI.Hidden'],
    shareRoles,
    sharePermissions,
    outgoing: false,
    privateLink: urlJoin(serverUrl, 'f', driveItem.remoteItem.id),
    spaceId: driveItem.remoteItem.spaceId,
    canRename: () => driveItem['@client.synchronize'],
    canDownload: () => sharePermissions.includes(GraphSharePermission.readContent),
    canUpload: () => sharePermissions.includes(GraphSharePermission.createUpload),
    canCreate: () => sharePermissions.includes(GraphSharePermission.createChildren),
    canBeDeleted: () => sharePermissions.includes(GraphSharePermission.deleteStandard),
    canEditTags: () => sharePermissions.includes(GraphSharePermission.createChildren),
    canListVersions: () => sharePermissions.includes(GraphSharePermission.readVersions),
    isMounted: () => false,
    isReceivedShare: () => true,
    canShare: () => false,
    canDeny: () => false,
    getDomSelector: () => extractDomSelector(driveItem.id)
  }

  resource.extension = extractExtensionFromFile(resource)

  return resource
}

export function buildOutgoingShareResource({
  driveItem,
  user,
  serverUrl
}: {
  driveItem: DriveItem
  user: User
  serverUrl: string
}): OutgoingShareResource {
  const storageId = extractStorageId(driveItem.id)
  const path =
    driveItem.parentReference.path === '.'
      ? driveItem.parentReference.path
      : urlJoin(driveItem.parentReference.path, driveItem.name)

  const resource: OutgoingShareResource = {
    id: driveItem.id,
    driveId: driveItem.parentReference?.driveId,
    path,
    name: driveItem.name,
    fileId: driveItem.id,
    size: driveItem.size,
    storageId,
    parentFolderId: driveItem.parentReference?.id,
    sdate: driveItem.permissions[0].createdDateTime,
    tags: [],
    webDavPath: buildWebDavSpacesPath(storageId, path),
    sharedBy: [{ id: user.id, displayName: user.displayName }],
    owner: { id: user.id, displayName: user.displayName },
    sharedWith: driveItem.permissions.map((p) => {
      if (p.link) {
        return {
          id: p.id,
          displayName: p.link['@libre.graph.displayName'],
          shareType: ShareTypes.link.value
        }
      }
      const shareType = getShareTypeFromPermission(p)
      return { ...(p.grantedToV2.user || p.grantedToV2.group), shareType }
    }),
    shareTypes: driveItem.permissions.map(getShareTypeFromPermission),
    isFolder: !!driveItem.folder,
    type: !!driveItem.folder ? 'folder' : 'file',
    mimeType: driveItem.file?.mimeType || 'httpd/unix-directory',
    outgoing: true,
    privateLink: urlJoin(serverUrl, 'f', driveItem.id),
    spaceId: driveItem.parentReference?.driveId || '',
    canRename: () => true,
    canDownload: () => true,
    canUpload: () => true,
    canCreate: () => true,
    canBeDeleted: () => true,
    canEditTags: () => true,
    isMounted: () => false,
    isReceivedShare: () => true,
    canShare: () => true,
    canDeny: () => true,
    getDomSelector: () => extractDomSelector(driveItem.id)
  }

  resource.extension = extractExtensionFromFile(resource)

  return resource
}

export function buildCollaboratorShare({
  graphPermission,
  graphRoles,
  resourceId,
  indirect = false
}: {
  graphPermission: Permission
  graphRoles: Record<string, ShareRole>
  resourceId: string
  indirect?: boolean
}): CollaboratorShare {
  const role = graphRoles[graphPermission.roles?.[0]]
  const invitedBy = graphPermission.invitation?.invitedBy?.user

  return {
    id: graphPermission.id,
    resourceId,
    indirect,
    shareType: getShareTypeFromPermission(graphPermission),
    role,
    sharedBy: { id: invitedBy?.id, displayName: invitedBy?.displayName },
    sharedWith: graphPermission.grantedToV2.user || graphPermission.grantedToV2.group,
    permissions: (graphPermission['@libre.graph.permissions.actions']
      ? graphPermission['@libre.graph.permissions.actions']
      : role.rolePermissions.flatMap((p) => p.allowedResourceActions)) as GraphSharePermission[],
    createdDateTime: graphPermission.createdDateTime,
    expirationDateTime: graphPermission.expirationDateTime
  }
}

export function buildLinkShare({
  graphPermission,
  resourceId,
  indirect = false
}: {
  graphPermission: Permission
  resourceId: string
  indirect?: boolean
}): LinkShare {
  const invitedBy = graphPermission.invitation?.invitedBy?.user

  return {
    id: graphPermission.id,
    resourceId,
    indirect,
    shareType: ShareTypes.link.value,
    sharedBy: { id: invitedBy?.id, displayName: invitedBy?.displayName },
    hasPassword: graphPermission.hasPassword,
    createdDateTime: graphPermission.createdDateTime,
    expirationDateTime: graphPermission.expirationDateTime,
    displayName: graphPermission.link['@libre.graph.displayName'],
    isQuickLink: graphPermission.link['@libre.graph.quickLink'],
    type: graphPermission.link.type,
    webUrl: graphPermission.link.webUrl,
    preventsDownload: graphPermission.link.preventsDownload
  }
}

function getShareTypeFromPermission({ link, grantedToV2 }: Permission) {
  if (link) {
    return ShareTypes.link.value
  }
  if (grantedToV2?.group) {
    return ShareTypes.group.value
  }
  if (grantedToV2?.user?.['@libre.graph.userType'] === 'Federated') {
    return ShareTypes.remote.value
  }
  return ShareTypes.user.value
}
