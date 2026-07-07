import type { Application } from '../generated'
import type { GraphRequestOptions } from '../types'

export interface GraphApplications {
  getApplication(id: string, requestOptions?: GraphRequestOptions): Promise<Application>
  listApplications(requestOptions?: GraphRequestOptions): Promise<Application[]>
}
