import { expect } from '@playwright/test'
import { objects } from '../../support'
import { config } from '../../config'
import { getWorld } from '../../environment/world'
import { substitute } from '../../support/utils'

export async function userOpensApplication({
  stepUser,
  name
}: {
  stepUser: string
  name: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const applicationObject = new objects.runtime.Application({ page })
  await applicationObject.open({ name })
}

export async function userShouldSeeNotifications({
  stepUser,
  expectedMessages
}: {
  stepUser: string
  expectedMessages: string[]
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const application = new objects.runtime.Application({ page })
  const messages = await application.getNotificationMessages()
  for (const message of expectedMessages) {
    expect(messages).toContain(substitute(message))
  }
}

export async function userShouldSeeNoNotifications({
  stepUser
}: {
  stepUser: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const application = new objects.runtime.Application({ page })
  const messages = await application.getNotificationMessages()
  expect(messages).toHaveLength(0)
}

export async function userWaitsForTokenRenewal({
  stepUser,
  renewalType
}: {
  stepUser: string
  renewalType: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const application = new objects.runtime.Application({ page })

  if (renewalType === 'iframe') {
    return await application.waitForTokenRenewalViaIframe()
  }
  return await application.waitForTokenRenewalViaRefreshToken()
}

export async function userOpensClipboardUrl({
  stepUser,
  url
}: {
  stepUser: string
  url: string
}): Promise<void> {
  const world = getWorld()
  const { page } = await world.actorsEnvironment.createActor({
    key: stepUser,
    namespace: world.actorsEnvironment.generateNamespace(stepUser, stepUser)
  })

  const applicationObject = new objects.runtime.Application({ page })
  // This is required as reading from clipboard is only possible when the browser is opened.
  await applicationObject.openUrl(config.baseUrl)
  url = url === '%clipboard%' ? await page.evaluate('navigator.clipboard.readText()') : url
  await applicationObject.openUrl(url)
}

export async function userReloadsPage({ stepUser }: { stepUser: string }): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const applicationObject = new objects.runtime.Application({ page })
  await applicationObject.reloadPage()
}

export async function userClosesSidebar({ stepUser }: { stepUser: string }): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const applicationObject = new objects.runtime.Application({ page })
  await applicationObject.closeSidebar()
}
