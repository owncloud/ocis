import { ApplicationsApi } from './../generated'
import { toInitOverrides, type GraphFactoryOptions } from './../types'
import type { GraphApplications } from './types'

export const ApplicationsFactory = ({ config }: GraphFactoryOptions): GraphApplications => {
  const applicationsApi = new ApplicationsApi(config)

  return {
    async getApplication(id, requestOptions) {
      return await applicationsApi.getApplication(
        { applicationId: id },
        toInitOverrides(requestOptions)
      )
    },

    async listApplications(requestOptions) {
      const { value } = await applicationsApi.listApplications(toInitOverrides(requestOptions))
      return value || []
    }
  }
}
