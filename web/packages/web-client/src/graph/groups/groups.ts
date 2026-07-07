import { urlJoin } from '../../utils'
import { GroupApiFactory, GroupsApiFactory } from './../generated'
import type { GraphFactoryOptions } from './../types'
import type { GraphGroups } from './types'

export const GroupsFactory = ({ axiosClient, config }: GraphFactoryOptions): GraphGroups => {
  const groupApiFactory = GroupApiFactory(config, config.basePath, axiosClient)
  const groupsApiFactory = GroupsApiFactory(config, config.basePath, axiosClient)

  return {
    async getGroup(id, options, requestOptions) {
      const { data } = await groupApiFactory.getGroup(
        id,
        options?.select ? new Set([...options.select]) : null,
        options?.expand ? new Set([...options.expand]) : new Set(['members']),
        requestOptions
      )
      return data
    },

    async createGroup(data, requestOptions) {
      const { data: group } = await groupsApiFactory.createGroup(data, requestOptions)
      return group
    },

    async editGroup(id, data, requestOptions) {
      const { data: group } = await groupApiFactory.updateGroup(id, data, requestOptions)
      return group
    },

    async deleteGroup(id, ifMatch, requestOptions) {
      await groupApiFactory.deleteGroup(id, ifMatch, requestOptions)
    },

    async listGroups(options, requestOptions) {
      const {
        data: { value }
      } = await groupsApiFactory.listGroups(
        options?.search,
        options?.orderBy ? new Set([...options.orderBy]) : null,
        options?.select ? new Set([...options.select]) : null,
        options?.expand ? new Set([...options.expand]) : null,
        requestOptions
      )
      return value
    },

    async addMember(groupId, userId, requestOptions) {
      await groupApiFactory.addMember(
        groupId,
        { '@odata.id': urlJoin(config.basePath, 'v1.0', 'users', userId) },
        requestOptions
      )
    },

    async deleteMember(groupId, userId, ifMatch, requestOptions) {
      await groupApiFactory.deleteMember(groupId, userId, ifMatch, requestOptions)
    }
  }
}
