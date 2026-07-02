import { BrowserContext, Page } from '@playwright/test'

export interface Link {
  name: string
  url: string
  password?: string
}

export interface Space {
  name: string
  id: string
  driveType?: string
}

export interface Actor {
  context: BrowserContext
  page: Page
  close(): Promise<void>
  savePage(page: Page): void
  newTab(): Promise<Page>
  closeCurrentTab(): Promise<void>
}

export interface User {
  /**
   * actual id, that will be exposed by the graph api
   */
  uuid?: string
  id: string
  displayName: string
  password: string
  email: string
  mail?: string
  role?: string
  preferredLanguage?: string
  /**
   * original id preserved for token lookups (used in parallel test scenarios)
   */
  originalId?: string
}

export interface File {
  name: string
  path: string
}

export interface Me {
  id: string
}

export interface Group {
  uuid?: string
  id: string
  displayName: string
  groupTypes?: string[]
  originalId?: string
}

export interface Token {
  userId: string
  accessToken: string
  refreshToken: string
}

// keycloak realm role
export interface KeycloakRealmRole {
  id: string
  name: string
}

export interface ApplicationEntity {
  appRoles: AppRole[]
  displayName: string
  id: string
}

export interface AppRole {
  displayName: string
  id: string
}

export interface FederatedShareInvitation {
  code: string
  description?: string
  email?: string
}

// previous user state
export interface UserState {
  autoAcceptShare?: boolean
  language?: string
}
