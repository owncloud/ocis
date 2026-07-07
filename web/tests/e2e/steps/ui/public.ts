import { objects } from '../../support'
import { editor } from '../../support/objects/app-files/utils'
import { substitute } from '../../support/utils'
import { expect } from '@playwright/test'
import { getWorld } from '../../environment/world'
import { fileAction, application } from '../../environment/constants'

export async function userOpensPublicLink({
  stepUser,
  name
}: {
  stepUser: string
  name: string
}): Promise<void> {
  const world = getWorld()
  const { page } = await world.actorsEnvironment.createActor({
    key: stepUser,
    namespace: world.actorsEnvironment.generateNamespace(stepUser, stepUser)
  })

  const { url } = world.linksEnvironment.getLink({ name })
  const pageObject = new objects.applicationFiles.page.Public({ page })
  await pageObject.open({ url })
}

export async function userCreatesPublicLink({
  stepUser,
  resource,
  password,
  role,
  name = 'Unnamed link'
}: {
  stepUser: string
  resource: string
  password: string
  role?: string
  name?: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const publicObject = new objects.applicationFiles.Link({ page })
  await publicObject.create({ resource, password: substitute(password), role, name })
}

export async function anonymousUserOpensPublicLink({
  stepUser,
  name
}: {
  stepUser: string
  name: string
}): Promise<void> {
  const world = getWorld()
  const { page } = await world.actorsEnvironment.createActor({
    key: stepUser,
    namespace: world.actorsEnvironment.generateNamespace(
      `${stepUser} user language change`,
      'Anonymous'
    )
  })

  const { url } = world.linksEnvironment.getLink({ name })
  const pageObject = new objects.applicationFiles.page.Public({ page })
  await pageObject.open({ url })
}

export async function userUnlocksPublicLink({
  password,
  stepUser
}: {
  password: string
  stepUser: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const pageObject = new objects.applicationFiles.page.Public({ page })
  if (password === '%copied_password%') {
    // Use world-specific stored password instead of clipboard (parallel safety)
    password = world.linksEnvironment.copiedPassword
  } else {
    password = substitute(password)
  }
  await pageObject.authenticate({ password })
}

export async function userUploadsResourcesInPublicLink({
  stepUser,
  resources
}: {
  stepUser: string
  resources: { name: string; to?: string; option?: string; type?: string }[]
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const pageObject = new objects.applicationFiles.page.Public({ page })
  for (const resource of resources) {
    await pageObject.upload({
      to: resource.to,
      resources: [world.filesEnvironment.getFile({ name: resource.name })],
      option: resource.option,
      type: resource.type
    })
  }
}

export async function userDeletesResourcesFromPublicLink({
  stepUser,
  actionType = fileAction.sideBarPanel,
  resources
}: {
  stepUser: string
  actionType: typeof fileAction.sideBarPanel | typeof fileAction.batchAction
  resources: { resource: string; parentFolder?: string }[]
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const pageObject = new objects.applicationFiles.page.Public({ page })
  for (const resource of resources) {
    await pageObject.delete({
      folder: resource.parentFolder ? resource.parentFolder : null,
      resourcesWithInfo: [{ name: resource.resource }],
      via: actionType
    })
  }
}

export async function userShouldBeInFileViewer({
  stepUser,
  fileViewerType
}: {
  stepUser: string
  fileViewerType:
    | typeof application.mediaViewer
    | typeof application.pdfViewer
    | typeof application.textEditor
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const fileViewerLocator = editor.fileViewerLocator({ page, fileViewerType })
  await expect(fileViewerLocator).toBeVisible()
}

export async function userTriesToUnlockPasswordProtectedFolderWithPassword({
  stepUser,
  password
}: {
  stepUser: string
  password: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const pageObject = new objects.applicationFiles.page.Public({ page })
  const linkObject = new objects.applicationFiles.Link({ page })
  password = substitute(password)
  await pageObject.authenticate({
    password,
    passwordProtectedFolder: true,
    expectToSucceed: false
  })
  const actualErrorMessage = await linkObject.checkErrorMessage({ passwordProtectedFolder: true })
  expect(actualErrorMessage).toBe('Incorrect password')
}

export async function userUnlocksPasswordProtectedFolderWithPassword({
  stepUser,
  password
}: {
  stepUser: string
  password: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const pageObject = new objects.applicationFiles.page.Public({ page })
  password = substitute(password)
  await pageObject.authenticate({
    password,
    passwordProtectedFolder: true
  })
}

export async function userDropUploadsResources({
  stepUser,
  resources
}: {
  stepUser: string
  resources: string[]
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const pageObject = new objects.applicationFiles.page.Public({ page })
  for (const resource of resources) {
    await pageObject.dropUpload({
      resources: [world.filesEnvironment.getFile({ name: resource })]
    })
  }
}

export async function userRefreshesTheOldLink({ stepUser }: { stepUser: string }): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const pageObject = new objects.applicationFiles.page.Public({ page })
  await pageObject.reload()
}

export async function userShouldNotBeAbleToOpenTheOldLink({
  stepUser,
  linkName
}: {
  stepUser: string
  linkName: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const pageObject = new objects.applicationFiles.page.Public({ page })
  const { url } = world.linksEnvironment.getLink({ name: linkName })
  await pageObject.expectThatLinkIsDeleted({ url })
}

export async function userRenamesPublicLinkResources({
  stepUser,
  resources
}: {
  stepUser: string
  resources: { resource: string; newName: string }[]
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const pageObject = new objects.applicationFiles.page.Public({ page })
  for (const resource of resources) {
    await pageObject.rename({ resource: resource.resource, newName: resource.newName })
  }
}
