import { test } from '../../environment/test'
import * as ui from '../../steps/ui/index'
import * as api from '../../steps/api/api'

test.describe('users management', () => {
  test.beforeEach(async () => {
    await ui.userLogsIn({ stepUser: 'Admin' })
  })

  test.afterEach(async () => {
    await ui.userLogsOut({ stepUser: 'Admin' })
  })

  test('user login can be managed in the admin settings', async () => {
    await api.usersHaveBeenCreated({ stepUser: 'Admin', users: ['Alice'] })
    await ui.userOpensApplication({ stepUser: 'Admin', name: 'admin-settings' })
    await ui.userNavigatesToUsersManagementPage({ stepUser: 'Admin' })
    await ui.userForbidsLoginForUserUsingContextMenu({ stepUser: 'Admin', key: 'Alice' })
    await ui.userFailsToLogin({ stepUser: 'Alice' })
    await ui.userOpensApplication({ stepUser: 'Admin', name: 'admin-settings' })
    await ui.userNavigatesToUsersManagementPage({ stepUser: 'Admin' })
    await ui.userAllowsLoginForUserUsingContextMenu({ stepUser: 'Admin', key: 'Alice' })
    await ui.userLogsIn({ stepUser: 'Alice' })
    await ui.userLogsOut({ stepUser: 'Alice' })
  })

  test('admin user can change personal quotas for users', async () => {
    await api.usersHaveBeenCreated({ stepUser: 'Admin', users: ['Alice', 'Brian'] })
    await ui.userLogsIn({ stepUser: 'Alice' })
    await ui.userLogsIn({ stepUser: 'Brian' })
    await ui.userOpensApplication({ stepUser: 'Admin', name: 'admin-settings' })
    await ui.userNavigatesToUsersManagementPage({ stepUser: 'Admin' })
    await ui.userChangesQuotaOfUserUsingContextMenu({
      stepUser: 'Admin',
      key: 'Alice',
      value: '500'
    })
    await ui.userShouldHaveQuota({ stepUser: 'Alice', quota: '500' })
    await ui.userChangesQuotaForUsersUsingBatchAction({
      stepUser: 'Admin',
      value: '20',
      users: ['Alice', 'Brian']
    })
    await ui.userShouldHaveQuota({ stepUser: 'Alice', quota: '20' })
    await ui.userShouldHaveQuota({ stepUser: 'Brian', quota: '20' })
    await ui.userLogsOut({ stepUser: 'Alice' })
    await ui.userLogsOut({ stepUser: 'Brian' })
  })

  test('user group assignments can be handled via batch actions', async () => {
    await api.usersHaveBeenCreated({ stepUser: 'Admin', users: ['Alice', 'Brian', 'Carol'] })
    await api.groupsHaveBeenCreated({ stepUser: 'Admin', groupIds: ['sales', 'finance'] })

    await ui.userOpensApplication({ stepUser: 'Admin', name: 'admin-settings' })
    await ui.userNavigatesToUsersManagementPage({ stepUser: 'Admin' })
    await ui.userAddsUsersToGroupsUsingBatchActions({
      stepUser: 'Admin',
      assignments: [
        { group: 'sales', users: ['Alice', 'Brian', 'Carol'] },
        { group: 'finance', users: ['Alice', 'Brian', 'Carol'] }
      ]
    })
    await ui.userSetsFilters({
      stepUser: 'Admin',
      filters: [{ filter: 'groups', values: ['sales department', 'finance department'] }]
    })
    await ui.usersShouldBeVisible({ stepUser: 'Admin', expectedUsers: ['Alice', 'Brian', 'Carol'] })
    await ui.userRemovesUsersFromGroupsUsingBatchActions({
      stepUser: 'Admin',
      assignments: [
        { user: 'Alice', groups: ['sales', 'finance'] },
        { user: 'Brian', groups: ['sales', 'finance'] }
      ]
    })
    await ui.userReloadsPage({ stepUser: 'Admin' })
    await ui.usersShouldBeVisible({ stepUser: 'Admin', expectedUsers: ['Carol'] })
    await ui.usersShouldNotBeVisible({ stepUser: 'Admin', expectedUsers: ['Alice', 'Brian'] })
  })

  test('edit user', async () => {
    await api.usersHaveBeenCreated({ stepUser: 'Admin', users: ['Alice'] })
    await ui.userOpensApplication({ stepUser: 'Admin', name: 'admin-settings' })
    await ui.userNavigatesToUsersManagementPage({ stepUser: 'Admin' })
    await ui.userUpdatesUserAttributeUsingContextMenu({
      stepUser: 'Admin',
      user: 'Alice',
      attribute: 'userName',
      value: 'anna'
    })
    await ui.userUpdatesUserAttributeUsingContextMenu({
      stepUser: 'Admin',
      user: 'anna',
      attribute: 'displayName',
      value: 'Anna Murphy'
    })
    await ui.userUpdatesUserAttributeUsingContextMenu({
      stepUser: 'Admin',
      user: 'anna',
      attribute: 'email',
      value: 'anna@example.org'
    })
    await ui.userUpdatesUserAttributeUsingContextMenu({
      stepUser: 'Admin',
      user: 'anna',
      attribute: 'password',
      value: 'password'
    })
    await ui.userUpdatesUserAttributeUsingContextMenu({
      stepUser: 'Admin',
      user: 'anna',
      attribute: 'role',
      value: 'Space Admin'
    })
    await ui.userLogsIn({ stepUser: 'anna' })
    await ui.userShouldHaveSelfInfo({
      stepUser: 'anna',
      info: [
        { key: 'username', value: 'anna' },
        { key: 'displayname', value: 'Anna Murphy' },
        { key: 'email', value: 'anna@example.org' }
      ]
    })
    await ui.userLogsOut({ stepUser: 'anna' })
  })

  test('assign user to groups', async () => {
    await api.usersHaveBeenCreated({ stepUser: 'Admin', users: ['Alice'] })
    await api.groupsHaveBeenCreated({
      stepUser: 'Admin',
      groupIds: ['sales', 'finance', 'security']
    })
    await api.usersHaveBeenAddedToGroup({
      stepUser: 'Admin',
      usersToAdd: [{ user: 'Alice', group: 'sales' }]
    })
    await ui.userOpensApplication({ stepUser: 'Admin', name: 'admin-settings' })
    await ui.userNavigatesToUsersManagementPage({ stepUser: 'Admin' })
    await ui.userAddsUserToGroupsUsingContextMenu({
      stepUser: 'Admin',
      user: 'Alice',
      groups: ['finance', 'security']
    })
    await ui.userRemovesUserFromGroupsUsingContextMenu({
      stepUser: 'Admin',
      user: 'Alice',
      groups: ['sales']
    })
    await ui.userLogsIn({ stepUser: 'Alice' })
    await ui.userShouldHaveSelfInfo({
      stepUser: 'Alice',
      info: [{ key: 'groups', value: 'finance department, security department' }]
    })
    await ui.userLogsOut({ stepUser: 'Alice' })
  })

  test('delete user', async () => {
    await api.usersHaveBeenCreated({
      stepUser: 'Admin',
      users: ['Alice', 'Brian', 'Carol', 'David']
    })
    await ui.userOpensApplication({ stepUser: 'Admin', name: 'admin-settings' })
    await ui.userNavigatesToUsersManagementPage({ stepUser: 'Admin' })
    await ui.userUpdatesUserAttributeUsingContextMenu({
      stepUser: 'Admin',
      user: 'David',
      attribute: 'role',
      value: 'Space Admin'
    })
    await ui.userSetsFilters({
      stepUser: 'Admin',
      filters: [{ filter: 'roles', values: ['User', 'Admin'] }]
    })
    await ui.usersShouldBeVisible({ stepUser: 'Admin', expectedUsers: ['Alice', 'Brian', 'Carol'] })
    await ui.usersShouldNotBeVisible({ stepUser: 'Admin', expectedUsers: ['David'] })
    await ui.userDeletesUsersUsingBatchActions({ stepUser: 'Admin', users: ['Alice', 'Brian'] })
    await ui.userDeletesUsersUsingContextMenu({ stepUser: 'Admin', users: ['Carol'] })
    await ui.usersShouldNotBeVisible({
      stepUser: 'Admin',
      expectedUsers: ['Alice', 'Brian', 'Carol']
    })
  })

  test('admin creates user', async () => {
    await ui.userOpensApplication({ stepUser: 'Admin', name: 'admin-settings' })
    await ui.userNavigatesToUsersManagementPage({ stepUser: 'Admin' })
    await ui.userCreatesUser({
      stepUser: 'Admin',
      userData: [
        {
          name: 'max',
          displayname: 'Max Testing',
          email: 'maxtesting@owncloud.com',
          password: '12345678'
        }
      ]
    })
    await ui.userLogsIn({ stepUser: 'Max' })
    await ui.userShouldHaveSelfInfo({
      stepUser: 'Max',
      info: [
        { key: 'username', value: 'max' },
        { key: 'displayname', value: 'Max Testing' },
        { key: 'email', value: 'maxtesting@owncloud.com' }
      ]
    })
    await ui.userLogsOut({ stepUser: 'Max' })
  })

  test('edit panel can be opened via quick action and context menu', async () => {
    await api.usersHaveBeenCreated({ stepUser: 'Admin', users: ['Alice', 'Brian', 'Carol'] })
    await ui.userOpensApplication({ stepUser: 'Admin', name: 'admin-settings' })
    await ui.userNavigatesToUsersManagementPage({ stepUser: 'Admin' })
    await ui.userOpensEditPanelOfUserUsingQuickAction({ stepUser: 'Admin', actionUser: 'Brian' })
    await ui.userShouldSeeEditPanel({ stepUser: 'Admin' })
    await ui.userOpensEditPanelOfUserUsingContextMenu({ stepUser: 'Admin', actionUser: 'Brian' })
    await ui.userShouldSeeEditPanel({ stepUser: 'Admin' })
  })
})
