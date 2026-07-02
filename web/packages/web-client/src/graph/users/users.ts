import {
  MeChangepasswordApiFactory,
  MeUserApiFactory,
  UserApiFactory,
  UserAppRoleAssignmentApiFactory,
  UsersApiFactory
} from './../generated'
import type { GraphFactoryOptions } from './../types'
import type { GraphUsers } from './types'

export const UsersFactory = ({ axiosClient, config }: GraphFactoryOptions): GraphUsers => {
  const userApiFactory = UserApiFactory(config, config.basePath, axiosClient)
  const usersApiFactory = UsersApiFactory(config, config.basePath, axiosClient)
  const meUserApiFactory = MeUserApiFactory(config, config.basePath, axiosClient)
  const meChangepasswordApiFactory = MeChangepasswordApiFactory(
    config,
    config.basePath,
    axiosClient
  )
  const userAppRoleAssignmentApiFactory = UserAppRoleAssignmentApiFactory(
    config,
    config.basePath,
    axiosClient
  )

  return {
    async getUser(id, options, requestOptions) {
      const { data } = await userApiFactory.getUser(
        id,
        options?.select ? new Set([...options.select]) : null,
        options?.expand
          ? new Set([...options.expand])
          : new Set(['drive', 'memberOf', 'appRoleAssignments']),
        requestOptions
      )
      return data
    },

    async createUser(data, requestOptions) {
      const { data: user } = await usersApiFactory.createUser(data, requestOptions)
      return user
    },

    async editUser(id, data, requestOptions) {
      const { data: user } = await userApiFactory.updateUser(id, data, requestOptions)
      return user
    },

    async deleteUser(id, ifMatch, requestOptions) {
      await userApiFactory.deleteUser(id, ifMatch, requestOptions)
    },

    async listUsers(options, requestOptions) {
      const {
        data: { value }
      } = await usersApiFactory.listUsers(
        options?.search,
        options?.filter,
        options?.orderBy ? new Set([...options.orderBy]) : null,
        options?.select ? new Set([...options.select]) : null,
        options?.expand ? new Set([...options.expand]) : null,
        requestOptions
      )
      return value
    },

    async getMe(options, requestOptions) {
      const { data } = await meUserApiFactory.getOwnUser(
        options?.expand ? new Set([...options.expand]) : new Set(['memberOf']),
        requestOptions
      )
      return data
    },

    async editMe(user, requestOptions) {
      const { data } = await meUserApiFactory.updateOwnUser(user, requestOptions)
      return data
    },

    async changeOwnPassword(change, requestOptions) {
      await meChangepasswordApiFactory.changeOwnPassword(change, requestOptions)
    },

    async exportPersonalData(id, destination, requestOptions) {
      await userApiFactory.exportPersonalData(id, destination, requestOptions)
    },

    async createUserAppRoleAssignment(id, roleAssignment, requestOptions) {
      const { data } = await userAppRoleAssignmentApiFactory.userCreateAppRoleAssignments(
        id,
        roleAssignment,
        requestOptions
      )
      return data
    }
  }
}
