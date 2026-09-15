import { ActivitiesApi } from './../generated'
import { toInitOverrides, type GraphFactoryOptions } from './../types'
import type { GraphActivities } from './types'

export const ActivitiesFactory = ({ config }: GraphFactoryOptions): GraphActivities => {
  const activitiesApi = new ActivitiesApi(config)

  return {
    async listActivities(kqlTerm, requestOptions) {
      const { value } = await activitiesApi.getActivities(
        { kql: kqlTerm },
        toInitOverrides(requestOptions)
      )
      return value || []
    }
  }
}
