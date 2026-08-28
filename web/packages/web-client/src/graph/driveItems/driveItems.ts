import { DriveItemApi, DrivesRootApi, MeDriveApi } from './../generated'
import { toInitOverrides, type GraphFactoryOptions } from './../types'
import type { GraphDriveItems } from './types'

export const DriveItemsFactory = ({ config }: GraphFactoryOptions): GraphDriveItems => {
  const driveItemApi = new DriveItemApi(config)
  const drivesRootApi = new DrivesRootApi(config)
  const meDriveApi = new MeDriveApi(config)

  return {
    async getDriveItem(driveId, itemId, requestOptions) {
      return await driveItemApi.getDriveItem({ driveId, itemId }, toInitOverrides(requestOptions))
    },

    async createDriveItem(driveId, data, requestOptions) {
      return await drivesRootApi.createDriveItem(
        { driveId, driveItem: data },
        toInitOverrides(requestOptions)
      )
    },

    async updateDriveItem(driveId, itemId, data, requestOptions) {
      return await driveItemApi.updateDriveItem(
        { driveId, itemId, driveItem: data },
        toInitOverrides(requestOptions)
      )
    },

    async deleteDriveItem(driveId, itemId, requestOptions) {
      await driveItemApi.deleteDriveItem({ driveId, itemId }, toInitOverrides(requestOptions))
    },

    async listSharedByMe(requestOptions) {
      const data = await meDriveApi.listSharedByMe(toInitOverrides(requestOptions))
      return data?.value || []
    },

    async listSharedWithMe(requestOptions) {
      const data = await meDriveApi.listSharedWithMe(toInitOverrides(requestOptions))
      return data?.value || []
    }
  }
}
