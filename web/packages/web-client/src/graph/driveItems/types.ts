import { DriveItem } from '../generated'
import type { GraphRequestOptions } from '../types'

export interface GraphDriveItems {
  getDriveItem: (
    driveId: string,
    itemId: string,
    requestOptions?: GraphRequestOptions
  ) => Promise<DriveItem>
  createDriveItem: (
    driveId: string,
    data: DriveItem,
    requestOptions?: GraphRequestOptions
  ) => Promise<DriveItem>
  updateDriveItem: (
    driveId: string,
    itemId: string,
    data: DriveItem,
    requestOptions?: GraphRequestOptions
  ) => Promise<DriveItem>
  deleteDriveItem: (
    driveId: string,
    itemId: string,
    requestOptions?: GraphRequestOptions
  ) => Promise<void>
  listSharedByMe: (requestOptions?: GraphRequestOptions) => Promise<DriveItem[]>
  listSharedWithMe: (requestOptions?: GraphRequestOptions) => Promise<DriveItem[]>
}
