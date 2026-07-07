import { objects } from '../../support'
import { getWorld } from '../../environment/world'
import { application, client } from '../../environment/constants'

export async function userNavigatesToNonExistingPage({
  stepUser
}: {
  stepUser: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const urlNavObject = new objects.urlNavigation.URLNavigation({ page })
  await urlNavObject.navigateToNonExistingPage()
}

export async function userShouldSeeNotFoundPage({ stepUser }: { stepUser: string }): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const urlNavObject = new objects.urlNavigation.URLNavigation({ page })
  await urlNavObject.waitForNotFoundPageToBeVisible()
}

export async function userOpensResourceViaUrl({
  stepUser,
  resource,
  space,
  editorName,
  clientType
}: {
  stepUser: string
  resource: string
  space: string
  editorName: typeof application.collabora | typeof application.onlyOffice
  clientType: keyof typeof client
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const user = world.usersEnvironment.getUser({ key: stepUser })
  const urlNavObject = new objects.urlNavigation.URLNavigation({ page })
  await urlNavObject.openResourceViaUrl({ resource, user, space, editorName, client: clientType })
}

export async function userOpensSpaceResourceViaUrl({
  stepUser,
  resource,
  space
}: {
  stepUser: string
  resource: string
  space: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const user = world.usersEnvironment.getUser({ key: stepUser })
  const urlNavObject = new objects.urlNavigation.URLNavigation({ page })
  await urlNavObject.openResourceViaUrl({ resource, user, space })
}

export async function userOpensResourceDetailsPanelViaUrl({
  stepUser,
  resource,
  detailsPanel,
  space
}: {
  stepUser: string
  resource: string
  detailsPanel: string
  space: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const user = world.usersEnvironment.getUser({ key: stepUser })
  const urlNavObject = new objects.urlNavigation.URLNavigation({ page })
  await urlNavObject.navigateToDetailsPanelOfResource({ resource, detailsPanel, user, space })
}

export async function userOpensSpaceViaUrl({
  stepUser,
  space
}: {
  stepUser: string
  space: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const user = world.usersEnvironment.getUser({ key: stepUser })
  const urlNavObject = new objects.urlNavigation.URLNavigation({ page })
  await urlNavObject.openSpaceViaUrl({ user, space })
}
