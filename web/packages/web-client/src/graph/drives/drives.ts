import { buildSpace } from '../../helpers'
import { Drive, DrivesApi, DrivesGetDrivesApi, MeDrivesApi } from './../generated'
import { toInitOverrides, type GraphFactoryOptions } from './../types'
import type { GraphDrives } from './types'

const getServerUrlFromDrive = (drive: Drive) => new URL(drive.webUrl).origin

export const DrivesFactory = ({ config }: GraphFactoryOptions): GraphDrives => {
  const drivesApi = new DrivesApi(config)
  const meDrivesApi = new MeDrivesApi(config)
  const allDrivesApi = new DrivesGetDrivesApi(config)

  return {
    async getDrive(id, graphRoles, requestOptions) {
      const drive = await drivesApi.getDriveBeta({ driveId: id }, toInitOverrides(requestOptions))
      return buildSpace({ ...drive, serverUrl: getServerUrlFromDrive(drive) }, graphRoles)
    },

    async createDrive(data, graphRoles, requestOptions) {
      const drive = await drivesApi.createDriveBeta(
        { drive: data },
        toInitOverrides(requestOptions)
      )
      return buildSpace({ ...drive, serverUrl: getServerUrlFromDrive(drive) }, graphRoles)
    },

    async updateDrive(id, data, graphRoles, requestOptions) {
      const drive = await drivesApi.updateDriveBeta(
        { driveId: id, driveUpdate: data },
        toInitOverrides(requestOptions)
      )
      return buildSpace({ ...drive, serverUrl: getServerUrlFromDrive(drive) }, graphRoles)
    },

    async disableDrive(id, ifMatch, requestOptions) {
      await drivesApi.deleteDriveBeta({ driveId: id, ifMatch }, toInitOverrides(requestOptions))
    },

    async deleteDrive(id, ifMatch, requestOptions) {
      await drivesApi.deleteDriveBeta(
        { driveId: id, ifMatch },
        toInitOverrides({
          ...requestOptions,
          headers: { ...(requestOptions?.headers || {}), Purge: 'T' }
        })
      )
    },

    async listMyDrives(graphRoles, options, requestOptions) {
      const { value } = await meDrivesApi.listMyDrivesBeta(
        { $orderby: options?.orderBy, $filter: options?.filter },
        toInitOverrides(requestOptions)
      )
      return value.map((d) => buildSpace({ ...d, serverUrl: getServerUrlFromDrive(d) }, graphRoles))
    },

    async listAllDrives(graphRoles, options, requestOptions) {
      const { value } = await allDrivesApi.listAllDrivesBeta(
        { $orderby: options?.orderBy, $filter: options?.filter },
        toInitOverrides(requestOptions)
      )
      return value.map((d) => buildSpace({ ...d, serverUrl: getServerUrlFromDrive(d) }, graphRoles))
    }
  }
}
