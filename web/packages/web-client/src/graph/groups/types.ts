import type { Group } from '../generated'
import type { GraphRequestOptions } from '../types'

export interface GraphGroups {
  getGroup: (
    id: string,
    options?: {
      expand?: Array<'members'>
      select?: Array<'id' | 'description' | 'displayName' | 'members'>
    },
    requestOptions?: GraphRequestOptions
  ) => Promise<Group>
  createGroup: (data: Group, requestOptions?: GraphRequestOptions) => Promise<Group>
  editGroup: (id: string, data: Group, requestOptions?: GraphRequestOptions) => Promise<void>
  deleteGroup: (id: string, ifMatch?: string, requestOptions?: GraphRequestOptions) => Promise<void>
  listGroups: (
    options?: {
      expand?: Array<'members'>
      orderBy?: Array<'displayName' | 'displayName desc'>
      search?: string
      select?: Array<'id' | 'description' | 'displayName' | 'mail' | 'members'>
    },
    requestOptions?: GraphRequestOptions
  ) => Promise<Group[]>
  addMember: (
    groupId: string,
    userId: string,
    requestOptions?: GraphRequestOptions
  ) => Promise<void>
  deleteMember: (
    groupId: string,
    userId: string,
    ifMatch?: string,
    requestOptions?: GraphRequestOptions
  ) => Promise<void>
}
