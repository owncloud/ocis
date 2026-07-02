import { RouteMeta } from 'vue-router'

type Dictionary<T> = { [key: string]: T }

export type QueryValue = string | (string | null)[]
export type LocationQuery = Dictionary<QueryValue>

export type ParamValue = string

export const authContextValues = ['anonymous', 'user', 'idp', 'publicLink', 'hybrid'] as const
export type AuthContext = (typeof authContextValues)[number]

export interface WebRouteMeta extends RouteMeta {
  title?: string
  authContext?: AuthContext
  patchCleanPath?: boolean
  contextQueryItems?: string[]
}
