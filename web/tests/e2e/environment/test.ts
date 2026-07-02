import { test as base, TestInfo } from '@playwright/test'
import { User, UserState, Group } from '../support/types'
import { config } from '../config.js'
import { api, store, environment, utils } from '../support'
import { World, setWorld } from './world'
import { state } from './shared'
import { getBrowserLaunchOptions } from '../support/environment/actor/shared'
import { Browser, chromium, firefox, webkit } from '@playwright/test'

export const test = base.extend<{
  world: World
  globalCleanup: void
  globalBeforeHook: void
}>({
  world: async ({}, use, testInfo) => {
    const world = new World(testInfo.workerIndex, testInfo.testId)
    setWorld(world)
    await use(world)
  },
  globalCleanup: [
    async ({ world }: { world: World }, use, testInfo) => {
      await use()

      config.federatedServer = false
      await world.actorsEnvironment.close()

      const adminUser = config.keycloak
        ? world.usersEnvironment.getUser({ key: config.keycloakAdminUser })
        : world.usersEnvironment.getUser({ key: config.adminUsername })

      if (!config.predefinedUsers && !config.mfa && adminUser) {
        if (config.keycloak) {
          await api.keycloak.refreshAccessTokenForKeycloakUser(adminUser)
          await api.keycloak.refreshAccessTokenForKeycloakOcisUser(adminUser)
        } else {
          await api.token.refreshAccessToken(adminUser)
        }

        if (isOcm(testInfo)) {
          config.federatedServer = true
          await api.token.refreshAccessToken(adminUser)
          await cleanUpUser(store.federatedUserStore, adminUser)
          config.federatedServer = false
        }
      }

      await cleanUpUser(store.createdUserStore, adminUser)
      await cleanUpGroup(adminUser)
      await cleanUpSpaces(adminUser)

      store.createdLinkStore.clear()
      store.createdTokenStore.clear()
      store.federatedTokenStore.clear()
      store.keycloakTokenStore.clear()

      utils.removeTempUploadDirectory()
      environment.closeSSEConnections()
    },
    { auto: true }
  ],
  globalBeforeHook: [
    async ({ world }: { world: World }, use, testInfo) => {
      if (!config.basicAuth && !config.predefinedUsers && !config.mfa) {
        if (config.keycloak) {
          const user = world.usersEnvironment.getUser({ key: config.keycloakAdminUser })
          await api.keycloak.setAccessTokenForKeycloakOcisUser(user)
          await api.keycloak.setAccessTokenForKeycloakUser(user)
          await storeKeycloakGroups(user)
        } else {
          const user = world.usersEnvironment.getUser({ key: config.adminUsername })
          await api.token.setAccessAndRefreshToken(user)
          if (isOcm(testInfo)) {
            config.federatedServer = true
            await api.token.setAccessAndRefreshToken(user)
            config.federatedServer = false
          }
        }
      }
      await use()
    },
    { auto: true }
  ]
})

test.beforeAll(async () => {
  const browserConfiguration = getBrowserLaunchOptions()
  const browsers: Record<string, () => Promise<Browser>> = {
    firefox: async () => await firefox.launch(browserConfiguration),
    webkit: async () => await webkit.launch(browserConfiguration),
    chrome: async () => await chromium.launch({ ...browserConfiguration, channel: 'chrome' }),
    chromium: async () => await chromium.launch(browserConfiguration)
  }
  state.browser = await browsers[config.browser]()

  if (config.keycloak) {
    // Get base user directly from store — no world available in beforeAll
    const adminUser = store.userStore.get(config.keycloakAdminUser.toLowerCase())
    if (adminUser) {
      api.keycloak.setupKeycloakAdminUser(adminUser)
    }
  }
})

const cleanUpUser = async (createdUserStore: Map<string, User>, adminUser: User) => {
  if (!adminUser) {
    return
  }

  const requests: Promise<User>[] = []
  for (const [key, user] of createdUserStore.entries()) {
    if (!config.predefinedUsers) {
      requests.push(api.provision.deleteUser({ user, admin: adminUser }))
    } else {
      await cleanupPredefinedUser(key, user)
    }
  }

  await awaitAllOrThrow('Failed to clean up users', requests)

  createdUserStore.clear()
  store.keycloakCreatedUser.clear()
}

const cleanupPredefinedUser = async (userKey: string, user: User) => {
  // delete the personal space resources
  const resources = await api.dav.listSpaceResources({ user, spaceType: 'personal' })
  for (const fileId in resources) {
    await api.dav.deleteSpaceResource({ user, fileId })
  }

  // cleanup trashbin if resources have been deleted
  if (Object.keys(resources).length) {
    await api.dav.emptyTrashbin({ user, spaceType: 'personal' })
  }

  // revert user state
  const usersEnvironment = new environment.UsersEnvironment()
  const userState: UserState = usersEnvironment.getUserState(userKey)
  if (userState.hasOwnProperty('autoAcceptShare')) {
    await api.settings.configureAutoAcceptShare({ user, state: userState.autoAcceptShare })
  }
  if (userState.hasOwnProperty('language')) {
    await api.settings.changeLanguage({ user, language: userState.language })
  }
}

const cleanUpGroup = async (adminUser: User) => {
  if (config.predefinedUsers || !adminUser) {
    return
  }
  const requests: Promise<Group>[] = []
  store.createdGroupStore.forEach((group) => {
    if (!group.id.startsWith('keycloak')) {
      requests.push(api.graph.deleteGroup({ group, admin: adminUser }))
    }
  })

  await awaitAllOrThrow('Failed to clean up groups', requests)
  store.createdGroupStore.clear()
}

const cleanUpSpaces = async (adminUser: User) => {
  if (config.predefinedUsers) {
    return
  }
  const requests: Promise<void>[] = []
  store.createdSpaceStore.forEach((space) => {
    requests.push(
      api.graph
        .disableSpace({
          user: adminUser,
          space
        })
        .then(async (res) => {
          if (res.status === 204) {
            await api.graph.deleteSpace({
              user: adminUser,
              space
            })
          }
        })
    )
  })
  await awaitAllOrThrow('Space clean up failed', requests)
  store.createdSpaceStore.clear()
}

const isOcm = (testInfo: TestInfo): boolean => {
  return testInfo.tags.includes('@ocm')
}

const awaitAllOrThrow = async <T>(label: string, requests: Promise<T>[]) => {
  const results = await Promise.allSettled(requests)
  const failures = results.filter((r) => r.status === 'rejected')
  if (failures.length > 0) {
    throw new Error(`${label}: ${failures.map((f) => f.reason?.message || f.reason).join(', ')}`)
  }
}

const storeKeycloakGroups = async (adminUser: User) => {
  const groups = await api.graph.getGroups(adminUser)

  store.dummyKeycloakGroupStore.forEach((dummyGroup) => {
    const matchingGroup = groups.find((group) => group.displayName === dummyGroup.displayName)
    if (matchingGroup) {
      const groupKey = (dummyGroup.originalId || dummyGroup.id).toLowerCase()
      store.createdGroupStore.set(groupKey, { ...dummyGroup, uuid: matchingGroup.id })
    }
  })
}
