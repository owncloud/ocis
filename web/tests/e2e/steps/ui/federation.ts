import { expect } from '@playwright/test'
import { objects } from '../../support'
import { substitute } from '../../support/utils'
import { getWorld } from '../../environment/world'

export async function userGeneratesInvitationTokenForTheFederationShare({
  stepUser
}: {
  stepUser: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const pageObject = new objects.scienceMesh.Federation({ page })
  const user = world.usersEnvironment.getUser({ key: stepUser })
  await pageObject.generateInvitation(user.id)
}

export async function userAcceptsFederatedShareInvitationByLocalUser({
  stepUser,
  sharer
}: {
  stepUser: string
  sharer: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const sharerUser = world.usersEnvironment.getUser({ key: sharer })
  const pageObject = new objects.scienceMesh.Federation({ page })
  await pageObject.acceptInvitation(sharerUser.id)
}

export async function userShouldSeeTheFederatedConnections({
  stepUser,
  federation
}: {
  stepUser: string
  federation: { user: string; email: string }[]
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const pageObject = new objects.scienceMesh.Federation({ page })
  for (const fed of federation) {
    fed.user = substitute(fed.user)
    fed.email = substitute(fed.email)
    const isConnectionExist = await pageObject.connectionExists(fed)
    expect(isConnectionExist).toBe(true)
  }
}
