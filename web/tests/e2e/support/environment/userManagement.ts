import { Group, User, UserState } from '../types'
import { getWorld } from '../../environment/world'
import {
  userStore,
  dummyGroupStore,
  createdUserStore,
  createdGroupStore,
  keycloakCreatedUser,
  federatedUserStore,
  dummyKeycloakGroupStore,
  userStateStore
} from '../store'
import { config } from '../../config'

export class UsersEnvironment {
  getUser({ key }: { key: string }): User {
    const userKey = key.toLowerCase()

    const base = userStore.get(userKey)

    if (!base) {
      throw new Error(`user with key '${userKey}' not found`)
    }

    const world = getWorld()
    if (world) {
      const id = world.getUserId(userKey)
      // unique email per worker to avoid Keycloak duplicate email (409) on parallel runs
      const [localPart, domain] = base.email.split('@')
      const email = `${localPart}+w${world.workerIndex}-${world.testId}@${domain}`

      return {
        ...base,
        id,
        email,
        // Keep original id for token lookup and original displayName for UI readability
        originalId: base.id
      }
    }

    return base
  }

  createUser({ key, user }: { key: string; user: User }): User {
    const userKey = key.toLowerCase()

    if (userStore.has(userKey)) {
      throw new Error(`user with key '${userKey}' already exists`)
    }

    userStore.set(userKey, user)

    return user
  }

  storeCreatedUser(key: string, user: User): User {
    const userKey = key.toLowerCase()
    const store = config.federatedServer ? federatedUserStore : createdUserStore
    if (store.has(userKey)) {
      throw new Error(`user '${userKey}' already exists`)
    }
    store.set(userKey, user)
    return user
  }

  private resolveCreatedUserKey(key: string): string | null {
    const store = config.federatedServer ? federatedUserStore : createdUserStore
    const world = getWorld()
    if (world) {
      const worldKey = world.getUserId(key).toLowerCase()
      if (store.has(worldKey)) {
        return worldKey
      }
    }
    const userKey = key.toLowerCase()
    if (store.has(userKey)) {
      return userKey
    }
    return null
  }

  getCreatedUser({ key }: { key: string }): User {
    const store = config.federatedServer ? federatedUserStore : createdUserStore
    const storeKey = this.resolveCreatedUserKey(key)
    if (storeKey) {
      return store.get(storeKey)
    }
    const userKey = key.toLowerCase()
    throw new Error(`user with key '${userKey}' not found`)
  }

  updateCreatedUser({ key, user }: { key: string; user: User }): User {
    const store = config.federatedServer ? federatedUserStore : createdUserStore
    const storeKey = this.resolveCreatedUserKey(key)

    if (storeKey) {
      store.delete(storeKey)
      // add to new key if the username is changed
      if (storeKey === user.id.toLowerCase()) {
        store.set(storeKey, user)
      } else {
        store.set(user.id, user)
      }
      return user
    }

    // Fall back to original key
    const userKey = key.toLowerCase()
    if (!store.has(userKey)) {
      throw new Error(`user '${userKey}' not found`)
    }
    store.delete(userKey)
    // add to new key if the username is changed
    if (userKey === user.id.toLowerCase()) {
      store.set(userKey, user)
    } else {
      store.set(user.id, user)
    }
    return user
  }

  removeCreatedUser({ key }: { key: string }): boolean {
    const store = config.federatedServer ? federatedUserStore : createdUserStore

    const storeKey = this.resolveCreatedUserKey(key)
    if (storeKey) {
      return store.delete(storeKey)
    }

    // If not found by key, try to find by user.id
    const userKey = key.toLowerCase()
    for (const [storedKey, storedUser] of store.entries()) {
      if (storedUser.id.toLowerCase() === userKey) {
        return store.delete(storedKey)
      }
    }

    throw new Error(`user '${userKey}' not found`)
  }

  getGroup({ key }: { key: string }): Group {
    const groupKey = key.toLowerCase()
    const store = groupKey.startsWith('keycloak') ? dummyKeycloakGroupStore : dummyGroupStore

    if (!store.has(groupKey)) {
      throw new Error(`group with key '${groupKey}' not found`)
    }

    const base = store.get(groupKey)

    // keycloak groups pre-exist and are not created by tests, so skip worker suffix
    const world = getWorld()
    if (world && !groupKey.startsWith('keycloak')) {
      const id = world.getGroupId(groupKey)
      const displayName = `${base.displayName} (${world.workerIndex})`

      return {
        ...base,
        id,
        displayName
      }
    }

    return base
  }

  getGroupDisplayName({ displayName }: { displayName: string }): string {
    const world = getWorld()
    if (!world) return displayName
    return `${displayName} (${world.workerIndex})`
  }

  getCreatedGroup({ key }: { key: string }): Group {
    const world = getWorld()
    // If world is available, try world-aware key first for parallel test safety
    if (world) {
      const worldKey = world.getGroupId(key).toLowerCase()
      if (createdGroupStore.has(worldKey)) {
        return createdGroupStore.get(worldKey)
      }
    }
    // Fall back to original key (for backward compatibility)
    const groupKey = key.toLowerCase()
    if (!createdGroupStore.has(groupKey)) {
      throw new Error(`group with key '${groupKey}' not found`)
    }
    return createdGroupStore.get(groupKey)
  }

  storeCreatedGroup({ group }: { group: Group }): Group {
    const groupKey = (group.originalId || group.id).toLowerCase()
    if (createdGroupStore.has(groupKey)) {
      throw new Error(`group with key '${groupKey}' already exists`)
    }
    createdGroupStore.set(groupKey, group)
    return group
  }

  storeCreatedKeycloakUser({ user }: { user: User }): User {
    if (keycloakCreatedUser.has(user.id)) {
      throw new Error(`Keycloak user '${user.id}' already exists`)
    }
    keycloakCreatedUser.set(user.id, user)
    return user
  }

  getCreatedKeycloakUser({ key }: { key: string }): User {
    const userKey = key.toLowerCase()
    if (!keycloakCreatedUser.has(userKey)) {
      throw new Error(`Keycloak user with key '${userKey}' not found`)
    }

    return keycloakCreatedUser.get(userKey)
  }

  removeCreatedKeycloakUser({ key }: { key: string }): boolean {
    const userKey = key.toLowerCase()

    if (!keycloakCreatedUser.has(userKey)) {
      throw new Error(`Keycloak user with key '${userKey}' not found`)
    }

    return keycloakCreatedUser.delete(userKey)
  }

  saveUserState(key: string, states: UserState): void {
    key = key.toLowerCase()
    let userStates = {}
    if (userStateStore.has(key)) {
      userStates = userStateStore.get(key)
    }
    userStateStore.set(key, { ...userStates, ...states })
  }

  getUserState(key: string): UserState {
    const userKey = key.toLowerCase()
    if (!userStateStore.has(userKey)) {
      throw new Error(`User key '${userKey}' not found`)
    }

    return userStateStore.get(userKey)
  }
}
