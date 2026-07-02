import type { ShareRole, SpaceResource } from '../../helpers'
import type { Drive } from '../generated'
import type { GraphRequestOptions } from '../types'

export interface GraphDrives {
  getDrive: (
    id: string,
    graphRoles: Record<string, ShareRole>,
    requestOptions?: GraphRequestOptions
  ) => Promise<SpaceResource>
  createDrive: (
    data: Drive,
    graphRoles: Record<string, ShareRole>,
    requestOptions?: GraphRequestOptions
  ) => Promise<SpaceResource>
  updateDrive: (
    id: string,
    data: Drive,
    graphRoles: Record<string, ShareRole>,
    requestOptions?: GraphRequestOptions
  ) => Promise<SpaceResource>
  disableDrive: (
    id: string,
    ifMatch?: string,
    requestOptions?: GraphRequestOptions
  ) => Promise<void>
  deleteDrive: (id: string, ifMatch?: string, requestOptions?: GraphRequestOptions) => Promise<void>
  listMyDrives: (
    graphRoles: Record<string, ShareRole>,
    options?: {
      orderBy?: string
      filter?: string
    },
    requestOptions?: GraphRequestOptions
  ) => Promise<SpaceResource[]>
  listAllDrives: (
    graphRoles: Record<string, ShareRole>,
    options?: {
      orderBy?: string
      filter?: string
    },
    requestOptions?: GraphRequestOptions
  ) => Promise<SpaceResource[]>
}
