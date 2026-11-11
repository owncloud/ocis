import { urlJoin } from '../../utils'
import { GroupApi, GroupsApi } from './../generated'
import { toInitOverrides, type GraphFactoryOptions } from './../types'
import type { GraphGroups } from './types'

export const GroupsFactory = ({ config }: GraphFactoryOptions): GraphGroups => {
  const groupApi = new GroupApi(config)
  const groupsApi = new GroupsApi(config)

  return {
    async getGroup(id, options, requestOptions) {
      return await groupApi.getGroup(
        {
          groupId: id,
          $select: options?.select ? new Set([...options.select]) : null,
          $expand: options?.expand ? new Set([...options.expand]) : new Set(['members'])
        },
        toInitOverrides(requestOptions)
      )
    },

    async createGroup(data, requestOptions) {
      return await groupsApi.createGroup({ group: data }, toInitOverrides(requestOptions))
    },

    async editGroup(id, data, requestOptions) {
      await groupApi.updateGroup({ groupId: id, group: data }, toInitOverrides(requestOptions))
    },

    async deleteGroup(id, ifMatch, requestOptions) {
      await groupApi.deleteGroup({ groupId: id, ifMatch }, toInitOverrides(requestOptions))
    },

    async listGroups(options, requestOptions) {
      const { value } = await groupsApi.listGroups(
        {
          $search: options?.search,
          $orderby: options?.orderBy ? new Set([...options.orderBy]) : null,
          $select: options?.select ? new Set([...options.select]) : null,
          $expand: options?.expand ? new Set([...options.expand]) : null
        },
        toInitOverrides(requestOptions)
      )
      return value
    },

    async addMember(groupId, userId, requestOptions) {
      await groupApi.addMember(
        {
          groupId,
          memberReference: { atOdataId: urlJoin(config.basePath, 'v1.0', 'users', userId) }
        },
        toInitOverrides(requestOptions)
      )
    },

    async deleteMember(groupId, userId, ifMatch, requestOptions) {
      await groupApi.deleteMember(
        { groupId, directoryObjectId: userId, ifMatch },
        toInitOverrides(requestOptions)
      )
    }
  }
}
