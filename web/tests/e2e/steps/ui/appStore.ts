import { expect } from '@playwright/test'
import { objects } from '../../support'
import { getWorld } from '../../environment/world'

export async function userOpensAppStore({ stepUser }: { stepUser: string }): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const pageObject = new objects.appStore.AppStore({ page })
  await pageObject.openAppStore()
}

export async function userShouldSeeAppStore({ stepUser }: { stepUser: string }): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const pageObject = new objects.appStore.AppStore({ page })
  await pageObject.waitForAppStoreToBeVisible()
}

export async function userSelectsApp({
  stepUser,
  app
}: {
  stepUser: string
  app: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const pageObject = new objects.appStore.AppStore({ page })
  await pageObject.selectApp(app)
}
export async function userShouldSeeAppDetails({
  stepUser,
  app
}: {
  stepUser: string
  app: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const pageObject = new objects.appStore.AppStore({ page })
  await pageObject.waitForAppDetailsToBeVisible(app)
}

export async function userDownloadsAppVersion({
  stepUser,
  version
}: {
  stepUser: string
  version: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const pageObject = new objects.appStore.AppStore({ page })
  const downloadedVersion = await pageObject.downloadAppVersion(version)
  expect(downloadedVersion).toContain(version)
}

export async function userDownloadsApp({
  stepUser,
  app
}: {
  stepUser: string
  app: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const pageObject = new objects.appStore.AppStore({ page })
  const download = await pageObject.downloadApp(app)
  expect(download).toBeDefined()
}

export async function userNavigatesToAppStoreOverview({
  stepUser
}: {
  stepUser: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const pageObject = new objects.appStore.AppStore({ page })
  await pageObject.navigateToAppStoreOverview()
}

export async function userShouldSeeApps({
  stepUser,
  expectedApps
}: {
  stepUser: string
  expectedApps: string[]
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const pageObject = new objects.appStore.AppStore({ page })
  const apps = await pageObject.getAppsList()
  for (const app of expectedApps) {
    expect(apps).toContain(app)
  }
}

export async function userSetsSearchTerm({
  stepUser,
  searchTerm
}: {
  stepUser: string
  searchTerm: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const pageObject = new objects.appStore.AppStore({ page })
  await pageObject.setSearchTerm(searchTerm)
}

export async function userSelectsAppTag({
  stepUser,
  tag,
  app
}: {
  stepUser: string
  tag: string
  app: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const pageObject = new objects.appStore.AppStore({ page })
  await pageObject.selectAppTag({ tag, app })
}

export async function userSelectsTag({
  stepUser,
  tag
}: {
  stepUser: string
  tag: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const pageObject = new objects.appStore.AppStore({ page })
  await pageObject.selectTag(tag)
}
