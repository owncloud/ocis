import {
  buildCollaboratorShare,
  buildLinkShare,
  CollaboratorShare,
  LinkShare,
  ShareRole
} from '../../helpers'
import {
  CollectionOfPermissionsWithAllowedValues,
  DrivesPermissionsApi,
  DrivesRootApi,
  Permission,
  RoleManagementApi,
  UnifiedRoleDefinition
} from './../generated'
import { toInitOverrides, type GraphFactoryOptions, type GraphRequestOptions } from './../types'
import type { GraphPermissions } from './types'

export const PermissionsFactory = ({ config }: GraphFactoryOptions): GraphPermissions => {
  const drivesRootApi = new DrivesRootApi(config)
  const roleManagementApi = new RoleManagementApi(config)
  const drivesPermissionsApi = new DrivesPermissionsApi(config)

  return {
    async getPermission<T extends CollaboratorShare | LinkShare>(
      driveId: string,
      itemId: string,
      permId: string,
      graphRoles: Record<string, ShareRole>,
      requestOptions: GraphRequestOptions
    ): Promise<T> {
      const permission = await drivesPermissionsApi.getPermission(
        { driveId, itemId, permId },
        toInitOverrides(requestOptions)
      )

      if (permission.link) {
        return buildLinkShare({ graphPermission: permission, resourceId: itemId }) as T
      }

      return buildCollaboratorShare({
        graphPermission: permission,
        resourceId: itemId,
        graphRoles: graphRoles || {}
      }) as T
    },

    async listPermissions(driveId, itemId, graphRoles, options, requestOptions) {
      let responseData: CollectionOfPermissionsWithAllowedValues

      if (driveId === itemId) {
        responseData = await drivesRootApi.listPermissionsSpaceRoot(
          {
            driveId,
            $filter: options?.filter,
            $select: options?.select ? new Set([...options.select]) : null
          },
          toInitOverrides(requestOptions)
        )
      } else {
        responseData = await drivesPermissionsApi.listPermissions(
          {
            driveId,
            itemId,
            $filter: options?.filter,
            $select: options?.select ? new Set([...options.select]) : null
          },
          toInitOverrides(requestOptions)
        )
      }

      const permissions = responseData.value || []
      const allowedActions = responseData.atLibreGraphPermissionsActionsAllowedValues
      const allowedRoles = responseData.atLibreGraphPermissionsRolesAllowedValues

      const shares = permissions.map((permission) => {
        if (permission.link) {
          return buildLinkShare({ graphPermission: permission, resourceId: itemId })
        }

        return buildCollaboratorShare({
          graphPermission: permission,
          resourceId: itemId,
          graphRoles: graphRoles || {}
        })
      })

      return { shares, allowedActions, allowedRoles }
    },

    async updatePermission<T extends CollaboratorShare | LinkShare>(
      driveId: string,
      itemId: string,
      permId: string,
      data: Permission,
      graphRoles: Record<string, ShareRole>,
      requestOptions: GraphRequestOptions
    ): Promise<T> {
      let permission: Permission

      if (driveId === itemId) {
        permission = await drivesRootApi.updatePermissionSpaceRoot(
          { driveId, permId, permission: data },
          toInitOverrides(requestOptions)
        )
      } else {
        permission = await drivesPermissionsApi.updatePermission(
          { driveId, itemId, permId, permission: data },
          toInitOverrides(requestOptions)
        )
      }

      if (permission.link) {
        return buildLinkShare({ graphPermission: permission, resourceId: itemId }) as T
      }

      return buildCollaboratorShare({
        graphPermission: permission,
        resourceId: itemId,
        graphRoles: graphRoles || {}
      }) as T
    },

    async deletePermission(driveId, itemId, permId, requestOptions) {
      if (driveId === itemId) {
        await drivesRootApi.deletePermissionSpaceRoot(
          { driveId, permId },
          toInitOverrides(requestOptions)
        )
        return
      }

      await drivesPermissionsApi.deletePermission(
        { driveId, itemId, permId },
        toInitOverrides(requestOptions)
      )
    },

    async createInvite(driveId, itemId, data, graphRoles, requestOptions) {
      let permission: Permission | undefined

      if (driveId === itemId) {
        const perm = await drivesRootApi.inviteSpaceRoot(
          { driveId, driveItemInvite: data },
          toInitOverrides(requestOptions)
        )

        permission = perm.value?.[0]
      } else {
        const perm = await drivesPermissionsApi.invite(
          { driveId, itemId, driveItemInvite: data },
          toInitOverrides(requestOptions)
        )

        permission = perm.value?.[0]
      }

      if (!permission) {
        throw new Error('no permission returned')
      }

      return buildCollaboratorShare({
        graphPermission: permission,
        resourceId: itemId,
        graphRoles: graphRoles || {}
      })
    },

    async createLink(driveId, itemId, data, requestOptions) {
      let permission: Permission

      if (driveId === itemId) {
        permission = await drivesRootApi.createLinkSpaceRoot(
          { driveId, driveItemCreateLink: data },
          toInitOverrides(requestOptions)
        )
      } else {
        permission = await drivesPermissionsApi.createLink(
          { driveId, itemId, driveItemCreateLink: data },
          toInitOverrides(requestOptions)
        )
      }

      return buildLinkShare({ graphPermission: permission, resourceId: itemId })
    },

    async setPermissionPassword(driveId, itemId, permId, data, requestOptions) {
      let permission: Permission

      if (driveId === itemId) {
        permission = await drivesRootApi.setPermissionPasswordSpaceRoot(
          { driveId, permId, sharingLinkPassword: data },
          toInitOverrides(requestOptions)
        )
      } else {
        permission = await drivesPermissionsApi.setPermissionPassword(
          { driveId, itemId, permId, sharingLinkPassword: data },
          toInitOverrides(requestOptions)
        )
      }

      return buildLinkShare({ graphPermission: permission, resourceId: itemId })
    },

    async listRoleDefinitions(requestOptions) {
      const data = await roleManagementApi.listPermissionRoleDefinitions(
        toInitOverrides(requestOptions)
      )

      // FIXME: graph type is wrong
      return data as UnifiedRoleDefinition[]
    }
  }
}
