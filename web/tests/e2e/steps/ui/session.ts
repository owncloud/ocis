import { config } from '../../config.js'
import { api, objects } from '../../support'
import { User } from '../../support/types'
import { listenSSE } from '../../support/environment/sse.js'
import { test, expect } from '@playwright/test'
import { waitForSSEEvent } from '../../support/utils/locator.js'
import { getWorld } from '../../environment/world'
import { Jimp } from 'jimp'
import { getOtpFromImage } from '../../support/utils/mfa.js'

async function createNewSession(stepUser: string) {
  const world = getWorld()
  const { page } = await world.actorsEnvironment.createActor({
    key: stepUser,
    namespace: world.actorsEnvironment.generateNamespace(stepUser, stepUser)
  })
  return new objects.runtime.Session({ page })
}

async function initUserStates(userKey: string, user: User) {
  const world = getWorld()
  const userInfo = await api.graph.getMeInfo(user)
  world.usersEnvironment.storeCreatedUser(userKey, {
    ...user,
    uuid: userInfo.id,
    email: userInfo.mail
  })
  world.usersEnvironment.saveUserState(userKey, {
    language: userInfo.preferredLanguage,
    autoAcceptShare: await api.settings.getAutoAcceptSharesValue(user)
  })
}

export async function userLogsIn({ stepUser }: { stepUser: string }): Promise<void> {
  const world = getWorld()
  const sessionObject = await createNewSession(stepUser)
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })

  let user = null
  if (stepUser === 'Admin' || config.predefinedUsers) {
    user = world.usersEnvironment.getUser({ key: stepUser })
    // Predefined users exist in OCIS with their real (non-namespaced) username.
    // Override the worker-namespaced id so login uses the real OCIS username.
    user = { ...user, id: user.originalId || user.id }
  } else {
    user = world.usersEnvironment.getCreatedUser({ key: stepUser })
  }

  await page.goto(config.baseUrl)
  await sessionObject.login(user)

  if (test.info().tags.length > 0) {
    // listen to SSE events when running scenarios with '@sse' tag
    if (test.info().tags.includes('@sse')) {
      void listenSSE(config.baseUrl, user)
    }
  }

  await page.locator('#web-content').waitFor()

  // initialize user states: uuid, language, auto-sync
  if (config.predefinedUsers) {
    await initUserStates(stepUser, user)
    // test should run with English language
    await api.settings.changeLanguage({ user, language: 'en' })
    await page.reload({ waitUntil: 'load' })
  }
}

export async function logInWithOTP({ stepUser }: { stepUser: string }): Promise<void> {
  const world = getWorld()
  const sessionObject = await createNewSession(stepUser)
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })

  const image = await Jimp.read('./qr.png')
  const { data, width, height } = image.bitmap
  const errorLocator = page.locator('#input-error-otp')
  for (let attempt = 0; attempt < 2; attempt++) {
    const otp = await getOtpFromImage(data, width, height)
    await sessionObject.keycloakOTPSignIn(String(otp))
    await page.waitForTimeout(1000)
    if (!(await errorLocator.isVisible())) {
      break
    } else {
      await page.waitForTimeout(config.tokenTimeout * 1000)
    }
  }
}

export async function userLogsOut({ stepUser }: { stepUser: string }): Promise<void> {
  const world = getWorld()
  // Check if actor exists (user might not have been logged in)
  let actor
  try {
    actor = world.actorsEnvironment.getActor({ key: stepUser })
  } catch {
    // Actor doesn't exist - user was never logged in, nothing to do
    return
  }

  // When using Keycloak, skip the browser-based logout (OIDC end-session).
  // Keycloak can send a backchannel logout to OCIS which may affect other
  // workers sharing the same user session. Close the browser context directly
  // instead — the server-side session expires naturally via SSO session timeout.
  if (!config.keycloak) {
    const canLogout = !!(await actor.page.locator('#_userMenuButton').count())
    const sessionObject = new objects.runtime.Session({ page: actor.page })
    canLogout && (await sessionObject.logout())
  }
  await actor.close()
}

export async function userShouldGetSSEEvent({
  stepUser,
  event
}: {
  stepUser: string
  event: string
}): Promise<void> {
  const world = getWorld()
  const user = world.usersEnvironment.getCreatedUser({ key: stepUser })
  await waitForSSEEvent(user, event)
}

export async function userClosesTheCurrentTab({ stepUser }: { stepUser: string }): Promise<void> {
  const world = getWorld()
  const actor = world.actorsEnvironment.getActor({ key: stepUser })
  await actor.closeCurrentTab()
}

export async function userNavigatesToNewTab({ stepUser }: { stepUser: string }): Promise<void> {
  const world = getWorld()
  const actor = world.actorsEnvironment.getActor({ key: stepUser })
  await actor.newTab()
}

export async function userWaitsForTokenToExpire({ stepUser }: { stepUser: string }): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  // wait for the token to expire
  await page.waitForTimeout(config.tokenTimeout * 1000)
}

export async function useServer({ server }: { server: 'LOCAL' | 'FEDERATED' }): Promise<void> {
  config.federatedServer = server === 'FEDERATED'
}

export async function userFailsToLogin({ stepUser }: { stepUser: string }): Promise<void> {
  const world = getWorld()
  const sessionObject = await createNewSession(stepUser)
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const user = world.usersEnvironment.getUser({ key: stepUser })

  await page.goto(config.baseUrl)
  await sessionObject.signIn(user.id, user.password)
  expect(page.locator('#oc-login-error-message')).toBeVisible({ timeout: config.timeout })
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['loginErrorMessageLocator'],
    'login error message'
  )
}
