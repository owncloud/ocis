import { expect } from '@playwright/test'
import { objects } from '../../support'
import { UsersEnvironment } from '../../support/environment/userManagement'
import { getWorld } from '../../environment/world'
import { fileAction } from '../../environment/constants'

export async function userNavigatesToGeneralManagementPage({
  stepUser
}: {
  stepUser: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const pageObject = new objects.applicationAdminSettings.page.General({ page })
  await pageObject.navigate()
}

export async function userUploadsLogoFromLocalPath({
  stepUser,
  localFile
}: {
  stepUser: string
  localFile: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const generalObject = new objects.applicationAdminSettings.General({ page })
  const logoPath = world.filesEnvironment.getFile({ name: localFile.split('/').pop() }).path
  await generalObject.uploadLogo({ path: logoPath })
}

export async function userResetsLogo({ stepUser }: { stepUser: string }): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const generalObject = new objects.applicationAdminSettings.General({ page })
  await generalObject.resetLogo()
}

export async function userAllowsLoginForUserUsingContextMenu({
  stepUser,
  key
}: {
  stepUser: string
  key: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const usersObject = new objects.applicationAdminSettings.Users({ page })

  await usersObject.allowLogin({ key, action: fileAction.contextMenu })
}

export async function userForbidsLoginForUserUsingContextMenu({
  stepUser,
  key
}: {
  stepUser: string
  key: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const usersObject = new objects.applicationAdminSettings.Users({ page })

  await usersObject.forbidLogin({ key, action: fileAction.contextMenu })
}

export async function userChangesQuotaOfUserUsingContextMenu({
  stepUser,
  key,
  value
}: {
  stepUser: string
  key: string
  value: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const usersObject = new objects.applicationAdminSettings.Users({ page })

  await usersObject.changeQuota({ key, value, action: fileAction.contextMenu })
}

export async function userChangesQuotaOfUserUsingSidebarPanel({
  stepUser,
  key,
  value
}: {
  stepUser: string
  key: string
  value: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const usersObject = new objects.applicationAdminSettings.Users({ page })

  await usersObject.changeQuota({ key, value, action: fileAction.quickAction })
}

export async function userChangesQuotaForUsersUsingBatchAction({
  stepUser,
  value,
  users
}: {
  stepUser: string
  value: string
  users: string[]
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const usersObject = new objects.applicationAdminSettings.Users({ page })
  for (const user of users) {
    await usersObject.selectUser({ key: user })
  }
  await usersObject.changeQuotaUsingBatchAction({ value, users })
}

async function selectUsersAndGetIds({
  usersObject,
  userKeys
}: {
  usersObject: InstanceType<typeof objects.applicationAdminSettings.Users>
  userKeys: string[]
}): Promise<string[]> {
  const selectedUserIds: string[] = []

  for (const userKey of userKeys) {
    selectedUserIds.push(usersObject.getUUID({ key: userKey }))
    await usersObject.select({ key: userKey })
  }

  return selectedUserIds
}

export async function userAddsUsersToGroupsUsingBatchActions({
  stepUser,
  assignments
}: {
  stepUser: string
  assignments: Array<{ group: string; users: string[] }>
}): Promise<void> {
  const world = getWorld()
  if (assignments.length === 0) {
    return
  }

  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const usersObject = new objects.applicationAdminSettings.Users({ page })

  for (const { group, users } of assignments) {
    const selectedUserIds = await selectUsersAndGetIds({
      usersObject,
      userKeys: users
    })

    await usersObject.addToGroupsBatchAction({
      userIds: selectedUserIds,
      groups: [group]
    })
  }
}

export async function userRemovesUsersFromGroupsUsingBatchActions({
  stepUser,
  assignments
}: {
  stepUser: string
  assignments: Array<{ user: string; groups: string[] }>
}): Promise<void> {
  const world = getWorld()
  if (assignments.length === 0) {
    return
  }

  const users = assignments.map(({ user }) => user)
  const groups = assignments[0].groups
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const usersObject = new objects.applicationAdminSettings.Users({ page })
  const selectedUserIds = await selectUsersAndGetIds({
    usersObject,
    userKeys: users
  })

  await usersObject.removeFromGroupsBatchAtion({
    userIds: selectedUserIds,
    groups
  })
}

export async function userSetsFilters({
  stepUser,
  filters
}: {
  stepUser: string
  filters: { filter: string; values: string[] }[]
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const usersObject = new objects.applicationAdminSettings.Users({ page })

  for (const { filter, values } of filters) {
    await usersObject.filter({ filter, values })
  }
}

export async function usersShouldBeVisible({
  stepUser,
  expectedUsers
}: {
  stepUser: string
  expectedUsers: string[]
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const usersObject = new objects.applicationAdminSettings.Users({ page })
  const displayedUsers = await usersObject.getDisplayedUsers()

  for (const user of expectedUsers) {
    const userId = usersObject.getUUID({ key: user })
    expect(displayedUsers).toContain(userId)
  }
}

export async function usersShouldNotBeVisible({
  stepUser,
  expectedUsers
}: {
  stepUser: string
  expectedUsers: string[]
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const usersObject = new objects.applicationAdminSettings.Users({ page })
  const displayedUsers = await usersObject.getDisplayedUsers()

  for (const user of expectedUsers) {
    const userId = usersObject.getUUID({ key: user })
    expect(displayedUsers).not.toContain(userId)
  }
}

export async function userChangesNameOfUserUsingContextMenu({
  stepUser,
  key,
  value
}: {
  stepUser: string
  key: string
  value: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const usersObject = new objects.applicationAdminSettings.Users({ page })
  await usersObject.changeUser({
    key,
    attribute: 'userName',
    value,
    action: fileAction.contextMenu
  })
}

export async function userUpdatesUserAttributeUsingContextMenu({
  stepUser,
  user,
  attribute,
  value
}: {
  stepUser: string
  user: string
  attribute: string
  value: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const usersObject = new objects.applicationAdminSettings.Users({ page })
  await usersObject.changeUser({ key: user, attribute, value, action: fileAction.contextMenu })
}

export async function userDeletesUsersUsingBatchActions({
  stepUser,
  users
}: {
  stepUser: string
  users: string[]
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const usersObject = new objects.applicationAdminSettings.Users({ page })
  const userIds = []

  for (const user of users) {
    userIds.push(usersObject.getUUID({ key: user }))
    await usersObject.select({ key: user })
  }

  await usersObject.deleteUserUsingBatchAction({ userIds })
}

export async function userDeletesUsersUsingContextMenu({
  stepUser,
  users
}: {
  stepUser: string
  users: string[]
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const usersObject = new objects.applicationAdminSettings.Users({ page })

  for (const user of users) {
    await usersObject.deleteUserUsingContextMenu({ key: user })
  }
}

export async function userShouldHaveSelfInfo({
  stepUser,
  info
}: {
  stepUser: string
  info: { key: string; value: string }[]
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const accountObject = new objects.account.Account({ page })

  for (const { key, value } of info) {
    let expectedValue = value
    // For group info in parallel mode, convert expected displayName to actual displayName
    if (world && key === 'groups') {
      expectedValue = value
        .split(', ')
        .map((name) => world.usersEnvironment.getGroupDisplayName({ displayName: name }))
        .join(', ')
    }
    const actual = await accountObject.getUserInfo(key)
    expect(actual).toBe(expectedValue)
  }
}

export async function userCreatesUser({
  stepUser,
  userData
}: {
  stepUser: string
  userData: { name: string; displayname: string; email: string; password: string }[]
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const usersObject = new objects.applicationAdminSettings.Users({ page })
  for (const info of userData) {
    await usersObject.createUser({
      name: info.name,
      displayname: info.displayname,
      email: info.email,
      password: info.password
    })
  }
}

export async function userOpensEditPanelOfUserUsingQuickAction({
  stepUser,
  actionUser
}: {
  stepUser: string
  actionUser: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const usersObject = new objects.applicationAdminSettings.Users({ page })
  await usersObject.openEditPanel({ key: actionUser, action: fileAction.quickAction })
}

export async function userOpensEditPanelOfUserUsingContextMenu({
  stepUser,
  actionUser
}: {
  stepUser: string
  actionUser: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const usersObject = new objects.applicationAdminSettings.Users({ page })
  await usersObject.openEditPanel({ key: actionUser, action: fileAction.contextMenu })
}

export async function userShouldSeeEditPanel({ stepUser }: { stepUser: string }): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const usersObject = new objects.applicationAdminSettings.Users({ page })
  await usersObject.waitForEditPanelToBeVisible()
}

export async function userAddsUserToGroupsUsingContextMenu({
  stepUser,
  user,
  groups
}: {
  stepUser: string
  user: string
  groups: string[]
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const usersObject = new objects.applicationAdminSettings.Users({ page })
  await usersObject.addToGroups({
    key: user,
    groups,
    action: fileAction.contextMenu
  })
}

export async function userRemovesUserFromGroupsUsingContextMenu({
  stepUser,
  user,
  groups
}: {
  stepUser: string
  user: string
  groups: string[]
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const usersObject = new objects.applicationAdminSettings.Users({ page })
  await usersObject.removeFromGroups({
    key: user,
    groups,
    action: fileAction.contextMenu
  })
}

export async function userNavigatesToGroupsManagementPage({
  stepUser
}: {
  stepUser: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const groupsObject = new objects.applicationAdminSettings.page.Groups({ page })
  await groupsObject.navigate()
}

export async function userNavigatesToUsersManagementPage({
  stepUser
}: {
  stepUser: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const pageObject = new objects.applicationAdminSettings.page.Users({ page })
  await pageObject.navigate()
}

export async function userCreatesGroups({
  stepUser,
  groupIds
}: {
  stepUser: string
  groupIds: string[]
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const groupsObject = new objects.applicationAdminSettings.Groups({ page })
  for (const groupId of groupIds) {
    await groupsObject.createGroup({ key: groupId })
  }
}

export async function userShouldSeeGroupIds({
  stepUser,
  expectedGroupIds
}: {
  stepUser: string
  expectedGroupIds: string[]
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const groupsObject = new objects.applicationAdminSettings.Groups({ page })
  const actualGroupsIds = await groupsObject.getDisplayedGroupsIds()
  for (const group of expectedGroupIds) {
    expect(actualGroupsIds).toContain(groupsObject.getUUID({ key: group }))
  }
}

export async function userShouldNotSeeGroupIds({
  stepUser,
  expectedGroupIds
}: {
  stepUser: string
  expectedGroupIds: string[]
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const groupsObject = new objects.applicationAdminSettings.Groups({ page })
  const actualGroupsIds = await groupsObject.getDisplayedGroupsIds()
  for (const group of expectedGroupIds) {
    expect(actualGroupsIds).not.toContain(groupsObject.getUUID({ key: group }))
  }
}

export async function userShouldSeeGroupDisplayName({
  stepUser,
  groupDisplayName
}: {
  stepUser: string
  groupDisplayName: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const groupsObject = new objects.applicationAdminSettings.Groups({ page })
  const groups = await groupsObject.getGroupsDisplayName()
  // Apply worker suffix to expected display name for parallel test safety,
  // matching the world transformation applied to group renames in Groups.changeGroup
  const expected = world
    ? new UsersEnvironment().getGroupDisplayName({ displayName: groupDisplayName })
    : groupDisplayName
  expect(groups).toContain(expected)
}

export async function userDeletesGroups({
  stepUser,
  actionType,
  groupsToBeDeleted
}: {
  stepUser: string
  actionType: typeof fileAction.batchAction | typeof fileAction.contextMenu
  groupsToBeDeleted: string[]
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const groupsObject = new objects.applicationAdminSettings.Groups({ page })
  const groupIds = []
  switch (actionType) {
    case fileAction.batchAction:
      for (const group of groupsToBeDeleted) {
        groupIds.push(groupsObject.getUUID({ key: group }))
        await groupsObject.selectGroup({ key: group })
      }
      await groupsObject.deleteGroupUsingBatchAction({ groupIds })
      break
    case fileAction.contextMenu:
      for (const group of groupsToBeDeleted) {
        await groupsObject.deleteGroupUsingContextMenu({ key: group })
      }
      break
    default:
      throw new Error(`'${actionType}' not implemented`)
  }
}

export async function userChangesGroup({
  stepUser,
  key,
  attribute,
  value,
  action
}: {
  stepUser: string
  key: string
  attribute: string
  value: string
  action: typeof fileAction.contextMenu | typeof fileAction.quickAction
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const groupsObject = new objects.applicationAdminSettings.Groups({ page })
  await groupsObject.changeGroup({
    key,
    attribute: attribute,
    value,
    action
  })
}

export async function userNavigatesToUserManagementPage({
  stepUser
}: {
  stepUser: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const pageObject = new objects.applicationAdminSettings.page.Users({ page })
  await pageObject.navigate()
}

export async function userAuthenticatesWithOTP({
  stepUser,
  deviceName
}: {
  stepUser: string
  deviceName: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const generalObject = new objects.applicationAdminSettings.General({ page })
  await generalObject.userAuthenticatesWithOTP({ deviceName })
}
