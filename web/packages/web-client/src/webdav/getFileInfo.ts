import { Resource, SpaceResource } from '../helpers'
import { ListFilesFactory, ListFilesOptions } from './listFiles'
import { WebDavOptions } from './types'

export const GetFileInfoFactory = (
  listFilesFactory: ReturnType<typeof ListFilesFactory>,
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  options?: WebDavOptions
) => {
  return {
    async getFileInfo(
      space: SpaceResource,
      resource?: { path?: string; fileId?: string },
      options?: ListFilesOptions
    ): Promise<Resource> {
      return (
        await listFilesFactory.listFiles(space, resource, {
          depth: 0,
          ...options
        })
      ).resource
    }
  }
}
