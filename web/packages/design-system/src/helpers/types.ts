import { RouteLocationRaw } from 'vue-router'
import type * as components from '../components'

export type OcComponents = {
  [K in keyof typeof components]: InstanceType<(typeof components)[K]>
}

export interface ContextualHelperDataListItem {
  text: string
  headline?: boolean
}

export interface ContextualHelperData {
  title: string
  text?: string
  list?: ContextualHelperDataListItem[]
  readMoreLink?: string
}

export interface ContextualHelper {
  isEnabled: boolean
  data: ContextualHelperData
}

export interface PasswordPolicyRule {
  code: string
  message: string
  helperMessage?: string
  format: (number | string)[]
  verified: boolean
}

export interface PasswordPolicy {
  rules: unknown[]

  check(password: string): boolean

  missing(password: string): {
    rules: PasswordPolicyRule[]
  }
}

// FIXME: ideally the id should not be optional, but some generated types (e.g. User and Group) need this
export type Item = {
  id?: string
}

export type FieldType = {
  name: string
  title?: string
  headerType?: string
  type?: string
  callback?: any
  alignH?: string
  alignV?: string
  width?: string
  wrap?: string
  thClass?: string
  tdClass?: string
  sortable?: boolean
  sortDir?: string
  prop?: string
  accessibleLabelCallback?: (item: Item) => string
}

export type Recipient = {
  name: string
  icon?: {
    name?: string
    label?: string
  }
  isLoadingAvatar?: boolean
  hasAvatar?: boolean
  avatar?: string
}

export interface BreadcrumbItem {
  id?: string
  text: string
  to?: RouteLocationRaw
  allowContextActions?: boolean
  onClick?: () => void
  isTruncationPlaceholder?: boolean
  isStaticNav?: boolean
}

export type AvailableSizeType =
  | 'xsmall'
  | 'small'
  | 'medium'
  | 'large'
  | 'xlarge'
  | 'xxlarge'
  | 'xxxlarge'
