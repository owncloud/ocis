import { KeycloakRealmRole, User, Group } from '../types'

export const keycloakRealmRoles = new Map<string, KeycloakRealmRole>()
export const keycloakCreatedUser = new Map<string, User>()

export const dummyKeycloakGroupStore = new Map<string, Group>([
  [
    'keycloak sales',
    {
      id: 'keycloak sales',
      displayName: 'keycloak sales department'
    }
  ],
  [
    'keycloak finance',
    {
      id: 'keycloak finance',
      displayName: 'keycloak finance department'
    }
  ]
])
