import { expect } from '@playwright/test'
import { objects } from '../../support'
import { getWorld } from '../../environment/world'

export async function userChangesLanguage({
  stepUser,
  language
}: {
  stepUser: string
  language: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const accountObject = new objects.account.Account({ page })
  const isAnonymousUser = stepUser === 'Anonymous'
  await accountObject.changeLanguage(language, isAnonymousUser)
}

export async function userShouldSeeAccountPageTitle({
  stepUser,
  expectedTitle
}: {
  stepUser: string
  expectedTitle: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const accountObject = new objects.account.Account({ page })
  const actualTitle = await accountObject.getTitle()
  expect(actualTitle).toEqual(expectedTitle)
}

export async function userRequestsGdprExport({ stepUser }: { stepUser: string }): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const accountObject = new objects.account.Account({ page })
  await accountObject.requestGdprExport()
}

export async function userDownloadsGdprExport({ stepUser }: { stepUser: string }): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const accountObject = new objects.account.Account({ page })
  await accountObject.downloadGdprExport()
}

export async function userMarksAllNotificationsAsRead({
  stepUser
}: {
  stepUser: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const application = new objects.runtime.Application({ page })
  await application.markNotificationsAsRead()
}

export async function userOpensAccountPage({ stepUser }: { stepUser: string }): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const accountObject = new objects.account.Account({ page })
  await accountObject.openAccountPage()
}

export async function userDisablesNotificationEvents({
  stepUser,
  events
}: {
  stepUser: string
  events: string[]
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const accountObject = new objects.account.Account({ page })
  for (const event of events) {
    await accountObject.disableNotificationEvent(event)
  }
}

export async function userShouldHaveQuota({
  stepUser,
  quota
}: {
  stepUser: string
  quota: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const accountObject = new objects.account.Account({ page })
  const actualQuota = await accountObject.getQuotaValue()
  expect(actualQuota).toBe(quota)
}
