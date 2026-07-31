import type {
  AppRoleAssignment,
  ExportPersonalDataRequest,
  PasswordChange,
  User
} from '../generated'
import type { GraphRequestOptions } from '../types'

export interface GraphUsers {
  getUser: (
    id: string,
    options?: {
      expand?: Array<'memberOf' | 'drive' | 'drives' | 'appRoleAssignments'>
      select?: Array<
        | 'id'
        | 'displayName'
        | 'drive'
        | 'drives'
        | 'mail'
        | 'memberOf'
        | 'onPremisesSamAccountName'
        | 'surname'
      >
    },
    requestOptions?: GraphRequestOptions
  ) => Promise<User>
  createUser: (data: User, requestOptions?: GraphRequestOptions) => Promise<User>
  editUser: (id: string, data: User, requestOptions?: GraphRequestOptions) => Promise<User>
  deleteUser: (id: string, ifMatch?: string, requestOptions?: GraphRequestOptions) => Promise<void>
  listUsers: (
    options?: {
      expand?: Array<'memberOf' | 'drive' | 'drives' | 'appRoleAssignments'>
      filter?: string
      orderBy?: Array<
        | 'displayName'
        | 'displayName desc'
        | 'mail'
        | 'mail desc'
        | 'onPremisesSamAccountName'
        | 'onPremisesSamAccountName desc'
      >
      search?: string
      select?: Array<
        'id' | 'displayName' | 'mail' | 'memberOf' | 'onPremisesSamAccountName' | 'surname'
      >
    },
    requestOptions?: GraphRequestOptions
  ) => Promise<User[]>
  getMe: (
    options?: {
      expand?: Array<'memberOf'>
    },
    requestOptions?: GraphRequestOptions
  ) => Promise<User>
  editMe: (user: User, requestOptions?: GraphRequestOptions) => Promise<User>
  changeOwnPassword: (change: PasswordChange, requestOptions?: GraphRequestOptions) => Promise<void>
  exportPersonalData: (
    id: string,
    destination?: ExportPersonalDataRequest,
    requestOptions?: GraphRequestOptions
  ) => Promise<void>
  createUserAppRoleAssignment: (
    id: string,
    roleAssignment: AppRoleAssignment,
    requestOptions?: GraphRequestOptions
  ) => Promise<AppRoleAssignment>
}
