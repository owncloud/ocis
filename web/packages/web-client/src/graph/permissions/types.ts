import { CollaboratorShare, LinkShare, ShareRole } from '../../helpers'
import type {
  DriveItemCreateLink,
  DriveItemInvite,
  ListPermissionsSpaceRootSelectEnum,
  Permission,
  SharingLinkPassword,
  UnifiedRoleDefinition
} from '../generated'
import type { GraphRequestOptions } from '../types'

type Share = CollaboratorShare | LinkShare

type ListPermissionsResponse = {
  shares: Share[]
  allowedActions: string[]
  allowedRoles: UnifiedRoleDefinition[]
}

export interface GraphPermissions {
  getPermission<T extends Share>(
    driveId: string,
    itemId: string,
    permId: string,
    graphRoles?: Record<string, ShareRole>,
    requestOptions?: GraphRequestOptions
  ): Promise<T>
  listPermissions(
    driveId: string,
    itemId: string,
    graphRoles?: Record<string, ShareRole>,
    options?: {
      filter?: string
      select?: Array<ListPermissionsSpaceRootSelectEnum>
    },
    requestOptions?: GraphRequestOptions
  ): Promise<ListPermissionsResponse>
  updatePermission<T extends Share>(
    driveId: string,
    itemId: string,
    permId: string,
    data: Permission,
    graphRoles?: Record<string, ShareRole>,
    requestOptions?: GraphRequestOptions
  ): Promise<T>
  deletePermission(
    driveId: string,
    itemId: string,
    permId: string,
    requestOptions?: GraphRequestOptions
  ): Promise<void>
  createInvite(
    driveId: string,
    itemId: string,
    data: DriveItemInvite,
    graphRoles?: Record<string, ShareRole>,
    requestOptions?: GraphRequestOptions
  ): Promise<CollaboratorShare>
  createLink(
    driveId: string,
    itemId: string,
    data: DriveItemCreateLink,
    requestOptions?: GraphRequestOptions
  ): Promise<LinkShare>
  setPermissionPassword(
    driveId: string,
    itemId: string,
    permId: string,
    data: SharingLinkPassword,
    requestOptions?: GraphRequestOptions
  ): Promise<LinkShare>
  listRoleDefinitions(requestOptions?: GraphRequestOptions): Promise<UnifiedRoleDefinition[]>
}
