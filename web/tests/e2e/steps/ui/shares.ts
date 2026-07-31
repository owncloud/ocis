import { expect } from '@playwright/test'
import { objects } from '../../support'
import { getDynamicRoleIdByName, ResourceType } from '../../support/api/share/share'
import { CollaboratorType, ICollaborator } from '../../support/objects/app-files/share/collaborator'
import { substitute } from '../../support/utils/substitute'
import { getWorld } from '../../environment/world'
import { fileAction } from '../../environment/constants'

const parseShareTable = function (
  resource: string,
  recipient: string,
  type: CollaboratorType,
  role: string,
  resourceType: string,
  expirationDate?: string,
  shareType?: string
) {
  const world = getWorld()
  const stepTable = [
    {
      resource,
      recipient,
      type,
      role,
      resourceType,
      expirationDate,
      shareType
    }
  ]
  return stepTable.reduce<Record<string, ICollaborator[]>>((acc, stepRow) => {
    const { resource, recipient, type, role, resourceType, expirationDate, shareType } = stepRow

    if (!acc[resource]) {
      acc[resource] = []
    }

    acc[resource].push({
      collaborator:
        type === 'group'
          ? world.usersEnvironment.getGroup({ key: recipient })
          : world.usersEnvironment.getUser({ key: recipient }),
      role,
      type: type,
      resourceType,
      expirationDate,
      shareType
    })

    return acc
  }, {})
}

export async function userNavigatesToSharedWithMePage({
  stepUser
}: {
  stepUser: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const pageObject = new objects.applicationFiles.page.shares.WithMe({ page })
  await pageObject.navigate()
}

export async function userNavigatesToSharedWithOthersPage({
  stepUser
}: {
  stepUser: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const pageObject = new objects.applicationFiles.page.shares.WithOthers({ page })
  await pageObject.navigate()
}

export async function userUpdatesShareeRoles({
  stepUser,
  roleUpdates
}: {
  stepUser: string
  roleUpdates: {
    resource: string
    recipient: string
    type: CollaboratorType
    role: string
    resourceType: string
    expirationDate?: string
    shareType?: string
  }[]
}) {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const shareObject = new objects.applicationFiles.Share({ page })
  const sharer = world.usersEnvironment.getUser({ key: stepUser })

  for (const update of roleUpdates) {
    const shareInfo = parseShareTable(
      update.resource,
      update.recipient,
      update.type,
      update.role,
      update.resourceType,
      update.expirationDate,
      update.shareType
    )

    for (const [resource, shareObj] of Object.entries(shareInfo)) {
      const roleId = await getDynamicRoleIdByName(
        sharer,
        shareObj[0].role,
        shareObj[0].resourceType as ResourceType
      )
      shareObj.forEach((item) => (item.role = roleId))
      await shareObject.changeShareeRole({
        resource,
        recipients: shareObj
      })
    }
  }
}

export async function userSharesResources({
  stepUser,
  actionType,
  shares
}: {
  stepUser: string
  actionType:
    typeof fileAction.sideBarPanel | typeof fileAction.quickAction | typeof fileAction.urlNavigation
  shares: {
    resource: string
    recipient: string
    type: CollaboratorType
    role: string
    resourceType: string
    expirationDate?: string
    shareType?: string
  }[]
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const shareObject = new objects.applicationFiles.Share({ page })
  const sharer = world.usersEnvironment.getUser({ key: stepUser })

  for (const resource of shares) {
    const shareRecipient = {
      collaborator:
        resource.type === 'group'
          ? world.usersEnvironment.getGroup({ key: resource.recipient })
          : world.usersEnvironment.getUser({ key: resource.recipient }),
      role: resource.role,
      type: resource.type as CollaboratorType,
      resourceType: resource.resourceType,
      expirationDate: resource.expirationDate,
      shareType: resource.shareType
    }

    shareRecipient.role = await getDynamicRoleIdByName(
      sharer,
      resource.role,
      resource.resourceType as ResourceType
    )
    await shareObject.create({
      resource: resource.resource,
      recipients: [shareRecipient],
      via: actionType
    })
  }
}

export async function userRemovesSharees({
  stepUser,
  sharees
}: {
  stepUser: string
  sharees: { resource: string; recipient: string; type?: 'group' | 'user' }[]
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const shareObject = new objects.applicationFiles.Share({ page })

  for (const sharee of sharees) {
    await shareObject.removeSharee({
      resource: sharee.resource,
      recipients: [
        {
          collaborator:
            sharee.type === 'group'
              ? world.usersEnvironment.getGroup({ key: sharee.recipient })
              : world.usersEnvironment.getUser({ key: sharee.recipient }),
          type: sharee.type as CollaboratorType
        }
      ]
    })
  }
}

export async function userAddsUsersToProjectSpace({
  stepUser,
  space,
  members
}: {
  stepUser: string
  space: string
  members: { reciver: string; role: string; kind: 'user' | 'group' }[]
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const spacesObject = new objects.applicationFiles.Spaces({ page })
  const sharer = world.usersEnvironment.getUser({ key: stepUser })

  for (const member of members) {
    const collaborator =
      member.kind === 'user'
        ? world.usersEnvironment.getUser({ key: member.reciver })
        : world.usersEnvironment.getGroup({ key: member.reciver })
    const roleId = await getDynamicRoleIdByName(sharer, member.role, 'space' as ResourceType)
    const collaboratorWithRole = {
      collaborator,
      role: roleId
    }
    await spacesObject.addMembers({ users: [collaboratorWithRole] })
  }
}

export async function userShouldBeAbleToManageShareOfFile({
  stepUser,
  resource,
  recipient
}: {
  stepUser: string
  resource: string
  recipient: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const shareObject = new objects.applicationFiles.Share({ page })
  const changeRole = shareObject.changeRoleLocator(
    world.usersEnvironment.getUser({ key: recipient })
  )
  const changeShare = shareObject.changeShareLocator(
    world.usersEnvironment.getUser({ key: recipient })
  )

  await shareObject.openSharingPanel(resource)

  const canChangeRole = !(await changeRole.isDisabled())
  const canChangeShare = !(await changeShare.isDisabled())

  expect(canChangeRole).toBe(true)
  expect(canChangeShare).toBe(true)
}

export async function userShouldNotBeAbleToManageShareOfFile({
  stepUser,
  resource,
  recipient
}: {
  stepUser: string
  resource: string
  recipient: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const shareObject = new objects.applicationFiles.Share({ page })
  const changeRole = shareObject.changeRoleLocator(
    world.usersEnvironment.getUser({ key: recipient })
  )
  const changeShare = shareObject.changeShareLocator(
    world.usersEnvironment.getUser({ key: recipient })
  )

  await shareObject.openSharingPanel(resource)

  const canChangeRole = !(await changeRole.isDisabled())
  const canChangeShare = !(await changeShare.isDisabled())

  expect(canChangeRole).toBe(false)
  expect(canChangeShare).toBe(false)
}

export async function userDisablesSyncForShares({
  stepUser,
  shares
}: {
  stepUser: string
  shares: string[]
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const shareObject = new objects.applicationFiles.Share({ page })
  for (const share of shares) {
    await shareObject.disableSync({ resource: share })
  }
}

export async function userEnablesSyncForShares({
  stepUser,
  shares
}: {
  stepUser: string
  shares: string[]
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const shareObject = new objects.applicationFiles.Share({ page })
  for (const share of shares) {
    await shareObject.enableSync({ resource: share, via: fileAction.contextMenu })
  }
}

export async function sharesShouldHaveSyncStatus({
  stepUser,
  shares
}: {
  stepUser: string
  shares: string[]
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const shareObject = new objects.applicationFiles.Share({ page })
  for (const share of shares) {
    expect(await shareObject.resourceIsSynced(share)).toBe(true)
  }
}

export async function sharesShouldNotHaveSyncStatus({
  stepUser,
  shares
}: {
  stepUser: string
  shares: string[]
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const shareObject = new objects.applicationFiles.Share({ page })
  for (const share of shares) {
    expect(await shareObject.resourceIsSynced(share)).toBe(false)
  }
}

export async function userShouldNotSeeShare({
  stepUser,
  resource,
  owner
}: {
  stepUser: string
  resource: string
  owner: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const shareObject = new objects.applicationFiles.Share({ page })
  const isAcceptedSharePresent = await shareObject.isAcceptedSharePresent(
    resource,
    substitute(owner)
  )
  expect(isAcceptedSharePresent).toBe(false)
}

export async function userEnablesSyncForAllShares({
  stepUser
}: {
  stepUser: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const shareObject = new objects.applicationFiles.Share({ page })
  await shareObject.syncAll()
}

export async function userChecksAccessDetailsOfShare({
  stepUser,
  resource,
  sharee,
  accessDetails
}: {
  stepUser: string
  resource: string
  sharee: { name: string; type: 'user' | 'group' }
  accessDetails: { Name: string; Type: string }
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const shareObject = new objects.applicationFiles.Share({ page })

  let selectorType = sharee.type
  // NOTE: external users have group type element selector
  if (accessDetails.hasOwnProperty('Type') && accessDetails.Type === 'External') {
    selectorType = 'group'
  }
  const expectedAccessDetails = {
    ...accessDetails,
    Name: substitute(accessDetails.Name).replace(/ \(\d+\)$/, '')
  }

  const actualDetails = await shareObject.getAccessDetails({
    resource,
    collaborator: {
      collaborator:
        sharee.type === 'group'
          ? world.usersEnvironment.getGroup({ key: sharee.name })
          : world.usersEnvironment.getUser({ key: sharee.name }),
      type: selectorType
    } as ICollaborator
  })

  const normalizedActualDetails = {
    ...actualDetails,
    Name: (actualDetails.Name || '').replace(/ \(\d+\)$/, '')
  }

  expect(normalizedActualDetails).toMatchObject(expectedAccessDetails)
}

export async function userShouldSeeAccessDetailsOfShareForFederatedUser({
  stepUser,
  resource,
  collaboratorName,
  detail
}: {
  stepUser: string
  resource: string
  collaboratorName: string
  detail: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const shareObject = new objects.applicationFiles.Share({ page })

  const actualDetails = await shareObject.getAccessDetails({
    resource,
    collaborator: {
      collaborator: world.usersEnvironment.getUser({ key: collaboratorName }),
      type: 'group'
    } as ICollaborator
  })

  expect(actualDetails).toHaveProperty(detail)
}

export async function userSetsExpirationDateOfShare({
  stepUser,
  resource,
  collaboratorType,
  collaboratorName,
  expirationDate
}: {
  stepUser: string
  resource: string
  collaboratorType: 'user' | 'group'
  collaboratorName: string
  expirationDate: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const shareObject = new objects.applicationFiles.Share({ page })
  const collaborator =
    collaboratorType === 'group'
      ? world.usersEnvironment.getGroup({ key: collaboratorName })
      : world.usersEnvironment.getUser({ key: collaboratorName })
  await shareObject.addExpirationDate({
    resource,
    collaborator: { collaborator, type: collaboratorType } as ICollaborator,
    expirationDate
  })
}

export async function userShouldSeeMessageOnWebUI({
  stepUser,
  message
}: {
  stepUser: string
  message: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const shareObject = new objects.applicationFiles.Share({ page })
  const actualMessage = await shareObject.getMessage()
  expect(actualMessage).toBe(message)
}

export async function userNavigatesToSharedViaLinkPage({
  stepUser
}: {
  stepUser: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const pageObject = new objects.applicationFiles.page.shares.ViaLink({ page })
  await pageObject.navigate()
}
