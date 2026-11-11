import {
  MeChangepasswordApi,
  MeUserApi,
  UserApi,
  UserAppRoleAssignmentApi,
  UsersApi
} from './../generated'
import { toInitOverrides, type GraphFactoryOptions } from './../types'
import type { GraphUsers } from './types'

export const UsersFactory = ({ config }: GraphFactoryOptions): GraphUsers => {
  const userApi = new UserApi(config)
  const usersApi = new UsersApi(config)
  const meUserApi = new MeUserApi(config)
  const meChangepasswordApi = new MeChangepasswordApi(config)
  const userAppRoleAssignmentApi = new UserAppRoleAssignmentApi(config)

  return {
    async getUser(id, options, requestOptions) {
      return await userApi.getUser(
        {
          userId: id,
          $select: options?.select ? new Set([...options.select]) : null,
          $expand: options?.expand
            ? new Set([...options.expand])
            : new Set(['drive', 'memberOf', 'appRoleAssignments'])
        },
        toInitOverrides(requestOptions)
      )
    },

    async createUser(data, requestOptions) {
      return await usersApi.createUser({ user: data }, toInitOverrides(requestOptions))
    },

    async editUser(id, data, requestOptions) {
      return await userApi.updateUser(
        { userId: id, userUpdate: data },
        toInitOverrides(requestOptions)
      )
    },

    async deleteUser(id, ifMatch, requestOptions) {
      await userApi.deleteUser({ userId: id, ifMatch }, toInitOverrides(requestOptions))
    },

    async listUsers(options, requestOptions) {
      const { value } = await usersApi.listUsers(
        {
          $search: options?.search,
          $filter: options?.filter,
          $orderby: options?.orderBy ? new Set([...options.orderBy]) : null,
          $select: options?.select ? new Set([...options.select]) : null,
          $expand: options?.expand ? new Set([...options.expand]) : null
        },
        toInitOverrides(requestOptions)
      )
      return value
    },

    async getMe(options, requestOptions) {
      return await meUserApi.getOwnUser(
        { $expand: options?.expand ? new Set([...options.expand]) : new Set(['memberOf']) },
        toInitOverrides(requestOptions)
      )
    },

    async editMe(user, requestOptions) {
      return await meUserApi.updateOwnUser({ userUpdate: user }, toInitOverrides(requestOptions))
    },

    async changeOwnPassword(change, requestOptions) {
      await meChangepasswordApi.changeOwnPassword(
        { passwordChange: change },
        toInitOverrides(requestOptions)
      )
    },

    async exportPersonalData(id, destination, requestOptions) {
      await userApi.exportPersonalData(
        { userId: id, exportPersonalDataRequest: destination },
        toInitOverrides(requestOptions)
      )
    },

    async createUserAppRoleAssignment(id, roleAssignment, requestOptions) {
      return await userAppRoleAssignmentApi.userCreateAppRoleAssignments(
        { userId: id, appRoleAssignment: roleAssignment },
        toInitOverrides(requestOptions)
      )
    }
  }
}
