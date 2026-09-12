import { TagsApi } from './../generated'
import { toInitOverrides, type GraphFactoryOptions } from './../types'
import type { GraphTags } from './types'

export const TagsFactory = ({ config }: GraphFactoryOptions): GraphTags => {
  const tagsApi = new TagsApi(config)

  return {
    async listTags(requestOptions) {
      const { value } = await tagsApi.getTags(toInitOverrides(requestOptions))
      return value || []
    },

    async assignTags(data, requestOptions) {
      await tagsApi.assignTags({ tagAssignment: data }, toInitOverrides(requestOptions))
    },

    async unassignTags(data, requestOptions) {
      await tagsApi.unassignTags({ tagUnassignment: data }, toInitOverrides(requestOptions))
    }
  }
}
