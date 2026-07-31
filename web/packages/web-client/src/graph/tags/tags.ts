import { TagsApiFactory } from './../generated'
import type { GraphFactoryOptions } from './../types'
import type { GraphTags } from './types'

export const TagsFactory = ({ axiosClient, config }: GraphFactoryOptions): GraphTags => {
  const tagsApiFactory = TagsApiFactory(config, config.basePath, axiosClient)

  return {
    async listTags(requestOptions) {
      const {
        data: { value }
      } = await tagsApiFactory.getTags(requestOptions)
      return value || []
    },

    async assignTags(data, requestOptions) {
      await tagsApiFactory.assignTags(data, requestOptions)
    },

    async unassignTags(data, requestOptions) {
      await tagsApiFactory.unassignTags(data, requestOptions)
    }
  }
}
