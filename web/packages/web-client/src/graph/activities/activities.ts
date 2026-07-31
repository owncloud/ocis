import { ActivitiesApiFactory } from './../generated'
import type { GraphFactoryOptions } from './../types'
import type { GraphActivities } from './types'

export const ActivitiesFactory = ({
  axiosClient,
  config
}: GraphFactoryOptions): GraphActivities => {
  const activitiesApiFactory = ActivitiesApiFactory(config, config.basePath, axiosClient)

  return {
    async listActivities(kqlTerm, requestOptions) {
      const {
        data: { value }
      } = await activitiesApiFactory.getActivities(kqlTerm, requestOptions)
      return value || []
    }
  }
}
