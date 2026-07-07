import { objects } from '../../support'
import { expect } from '@playwright/test'
import { getWorld } from '../../environment/world'
import { substitute } from '../../support/utils'

export async function userRenamesMostRecentlyCreatedPublicLinkOfResource({
  stepUser,
  resource,
  newName
}: {
  stepUser: string
  resource: string
  newName: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const linkObject = new objects.applicationFiles.Link({ page })
  const linkName = await linkObject.changeName({ resource, newName })
  expect(linkName).toBe(newName)
}

export async function userCopiesThePasswordOfThePublicLink({
  stepUser
}: {
  stepUser: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const linkObject = new objects.applicationFiles.Link({ page })
  // Copy password and store in linksEnvironment for parallel test safety
  world.linksEnvironment.copiedPassword = await linkObject.copyEnteredPassword()
}

export async function userCopiesTheLinkOfPasswordProtectedFolder({
  stepUser,
  resource
}: {
  stepUser: string
  resource: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const linkObject = new objects.applicationFiles.Link({ page })
  await linkObject.copyLinkToClipboard({
    resource: resource,
    resourceType: 'passwordProtectedFolder'
  })
}

export async function userClosesThePasswordProtectedFolderModal({
  stepUser
}: {
  stepUser: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const linkObject = new objects.applicationFiles.Link({ page })
  await linkObject.closeFolderModal()
}

export async function userChangesRoleOfPublicLinkOfResource({
  stepUser,
  resource,
  linkName,
  newRole
}: {
  stepUser: string
  resource: string
  linkName: string
  newRole: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const linkObject = new objects.applicationFiles.Link({ page })
  const roleText = await linkObject.changeRole({ linkName, resource, role: newRole })
  expect(roleText.toLowerCase()).toBe(newRole.toLowerCase())
}

export async function userSetsExpirationDateOfThePublicLinkOfResource({
  stepUser,
  resource,
  linkName,
  expireDate
}: {
  stepUser: string
  resource: string
  linkName: string
  expireDate: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const linkObject = new objects.applicationFiles.Link({ page })
  await linkObject.addExpiration({ resource, linkName, expireDate })
}

export async function userRemovesThePublicLinkOfResource({
  stepUser,
  resource,
  linkName
}: {
  stepUser: string
  resource: string
  linkName: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const linkObject = new objects.applicationFiles.Link({ page })
  await linkObject.delete({ resourceName: resource, name: linkName })
}

export async function userCreatesPublicLinkOfSpaceWithPassword({
  stepUser,
  password
}: {
  stepUser: string
  password: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const spaceObject = new objects.applicationFiles.Spaces({ page })
  password = substitute(password)
  await spaceObject.createPublicLink({ password })
}

export async function userRenamesTheMostRecentlyCreatedPublicLinkOfSpace({
  stepUser,
  newName
}: {
  stepUser: string
  newName: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const linkObject = new objects.applicationFiles.Link({ page })
  const linkName = await linkObject.changeName({ newName, space: true })
  expect(linkName).toBe(newName)
}

export async function userEditsThePublicLinkOfSpaceChangingRole({
  stepUser,
  linkName,
  role
}: {
  stepUser: string
  linkName: string
  role: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const linkObject = new objects.applicationFiles.Link({ page })
  const newPermission = await linkObject.changeRole({ linkName, role, space: true })
  expect(newPermission.toLowerCase()).toBe(role.toLowerCase())
}

export async function userShouldNotBeAbleToEditThePublicLink({
  stepUser,
  linkName
}: {
  stepUser: string
  linkName: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const linkObject = new objects.applicationFiles.Link({ page })
  const isVisible = await linkObject.islinkEditButtonVisibile(linkName)
  expect(isVisible).toBe(false)
}

export async function userChangesPasswordOfThePublicLinkOfResource({
  stepUser,
  resource,
  linkName,
  newPassword
}: {
  stepUser: string
  resource: string
  linkName: string
  newPassword: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const linkObject = new objects.applicationFiles.Link({ page })
  await linkObject.fillPassword({ resource, linkName, newPassword })
}

export async function userShouldSeeAnErrorMessage({
  stepUser,
  errorMessage
}: {
  stepUser: string
  errorMessage: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const linkObject = new objects.applicationFiles.Link({ page })
  const actualErrorMessage = await linkObject.checkErrorMessage()
  expect(actualErrorMessage).toBe(errorMessage)
}

export async function userClosesThePublicLinkPasswordDialogBox({
  stepUser
}: {
  stepUser: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const linkObject = new objects.applicationFiles.Link({ page })
  await linkObject.clickOnCancelButton()
}

export async function userRevealsThePasswordOfThePublicLink({
  stepUser
}: {
  stepUser: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const linkObject = new objects.applicationFiles.Link({ page })
  await linkObject.showOrHidePassword({ showOrHide: 'reveals' })
}

export async function userHidesThePasswordOfThePublicLink({
  stepUser
}: {
  stepUser: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const linkObject = new objects.applicationFiles.Link({ page })
  await linkObject.showOrHidePassword({ showOrHide: 'hides' })
}

export async function userGeneratesThePasswordForThePublicLink({
  stepUser
}: {
  stepUser: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const linkObject = new objects.applicationFiles.Link({ page })
  await linkObject.generatePassword()
}

export async function userSetsThePasswordOfThePublicLink({
  stepUser
}: {
  stepUser: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const linkObject = new objects.applicationFiles.Link({ page })
  await linkObject.setPassword()
}

export async function userCopiesLinkOfResource({
  stepUser,
  resource
}: {
  stepUser: string
  resource: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const linkObject = new objects.applicationFiles.Link({ page })
  await linkObject.copyLinkToClipboard({ resource })
}

export async function userChangesThePasswordOfPublicLink({
  stepUser,
  linkName,
  resource,
  newPassword
}: {
  stepUser: string
  linkName: string
  resource: string
  newPassword: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const linkObject = new objects.applicationFiles.Link({ page })
  await linkObject.addPassword({ resource, linkName, newPassword })
}
