import { ApplicationsApiFactory } from './../generated'
import type { GraphFactoryOptions } from './../types'
import type { GraphApplications } from './types'

export const ApplicationsFactory = ({
  axiosClient,
  config
}: GraphFactoryOptions): GraphApplications => {
  const applicationsApiFactory = ApplicationsApiFactory(config, config.basePath, axiosClient)

  return {
    async getApplication(id, requestOptions) {
      const { data } = await applicationsApiFactory.getApplication(id, requestOptions)
      return data
    },

    async listApplications(requestOptions) {
      const {
        data: { value }
      } = await applicationsApiFactory.listApplications(requestOptions)
      return value || []
    }
  }
}
