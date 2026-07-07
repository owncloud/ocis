import { createdTokenStore, federatedTokenStore, keycloakTokenStore } from '../store/token'
import { Token, User } from '../types'
import { config } from '../../config'

export type TokenProviderType = 'keycloak' | null | undefined
export type TokenEnvironmentType = KeycloakTokenEnvironment | IdpTokenEnvironment

export function TokenEnvironmentFactory(type?: TokenProviderType) {
  switch (type) {
    case 'keycloak':
      return new KeycloakTokenEnvironment()
    default:
      return new IdpTokenEnvironment()
  }
}

class IdpTokenEnvironment {
  getToken({ user }: { user: User }): Token {
    const store = config.federatedServer ? federatedTokenStore : createdTokenStore
    return store.get(user.id)
  }

  setToken({ user, token }: { user: User; token: Token }): Token {
    const store = config.federatedServer ? federatedTokenStore : createdTokenStore
    store.set(user.id, token)
    return token
  }

  deleteToken({ user }: { user: User }): void {
    createdTokenStore.delete(user.id)
  }
}

class KeycloakTokenEnvironment {
  getToken({ user }: { user: User }): Token {
    return keycloakTokenStore.get(user.id)
  }

  setToken({ user, token }: { user: User; token: Token }): Token {
    keycloakTokenStore.set(user.id, token)
    return token
  }

  deleteToken({ user }: { user: User }): void {
    keycloakTokenStore.delete(user.id)
  }
}
