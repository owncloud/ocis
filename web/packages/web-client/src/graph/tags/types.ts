import type { TagAssignment, TagUnassignment } from '../generated'
import type { GraphRequestOptions } from '../types'

export interface GraphTags {
  listTags(requestOptions?: GraphRequestOptions): Promise<string[]>
  assignTags(data: TagAssignment, requestOptions?: GraphRequestOptions): Promise<void>
  unassignTags(data: TagUnassignment, requestOptions?: GraphRequestOptions): Promise<void>
}
