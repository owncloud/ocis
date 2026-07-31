import type { Activity } from '../generated'
import type { GraphRequestOptions } from '../types'

export interface GraphActivities {
  listActivities(kqlTerm?: string, requestOptions?: GraphRequestOptions): Promise<Activity[]>
}
