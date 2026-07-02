import { test } from '../../environment/test'
import * as ui from '../../steps/ui/index'
import * as api from '../../steps/api/api'

test.describe('spaces management', () => {
  test.beforeEach(async () => {
    await api.usersHaveBeenCreated({ stepUser: 'Admin', users: ['Alice'] })
  })

  test.afterEach(async () => {
    await ui.userLogsOut({ stepUser: 'Alice' })
  })

  test('spaces can be created', async () => {
    await api.userHasAssignedRolesToUsers({
      stepUser: 'Admin',
      users: [{ id: 'Alice', role: 'Space Admin' }]
    })
    await api.userHasCreatedProjectSpaces({
      stepUser: 'Alice',
      spaces: [{ name: 'team A', id: 'team.a' }]
    })
    await ui.userLogsIn({ stepUser: 'Alice' })
    await ui.userOpensApplication({ stepUser: 'Alice', name: 'admin-settings' })
    await ui.userNavigatesToProjectSpaceManagementPage({ stepUser: 'Alice' })
    await ui.userCreatesProjectSpaces({
      stepUser: 'Alice',
      spaces: [{ name: 'team B', id: 'team.b' }]
    })
    await ui.userShouldSeeSpaces({ stepUser: 'Alice', expectedSpaceIds: ['team.a', 'team.b'] })
  })

  test('spaces can be managed in the admin settings via the context menu', async () => {
    await api.usersHaveBeenCreated({ stepUser: 'Admin', users: ['Brian'] })
    await api.userHasAssignedRolesToUsers({
      stepUser: 'Admin',
      users: [
        { id: 'Alice', role: 'Space Admin' },
        { id: 'Brian', role: 'Space Admin' }
      ]
    })
    await api.userHasCreatedProjectSpaces({
      stepUser: 'Alice',
      spaces: [
        { name: 'team A', id: 'team.a' },
        { name: 'team B', id: 'team.b' }
      ]
    })
    await ui.userLogsIn({ stepUser: 'Alice' })
    await ui.userOpensApplication({ stepUser: 'Alice', name: 'admin-settings' })
    await ui.userNavigatesToProjectSpaceManagementPage({ stepUser: 'Alice' })
    await ui.userUpdatesSpaceUsingContextMenu({
      stepUser: 'Alice',
      spaceId: 'team.a',
      updates: [
        { attribute: 'name', value: 'developer team' },
        { attribute: 'subtitle', value: 'developer team-subtitle' },
        { attribute: 'quota', value: '50' }
      ]
    })
    await ui.userDisablesSpaceUsingContextMenu({ stepUser: 'Alice', spaceId: 'team.a' })
    await ui.userEnablesSpaceUsingContextMenu({ stepUser: 'Alice', spaceId: 'team.a' })
    await ui.userShouldSeeSpaces({ stepUser: 'Alice', expectedSpaceIds: ['team.a'] })
    await ui.userLogsIn({ stepUser: 'Brian' })
    await ui.userOpensApplication({ stepUser: 'Brian', name: 'admin-settings' })
    await ui.userNavigatesToProjectSpaceManagementPage({ stepUser: 'Brian' })
    await ui.userDisablesSpaceUsingContextMenu({ stepUser: 'Brian', spaceId: 'team.b' })
    await ui.userDeletesSpaceUsingContextMenu({ stepUser: 'Brian', spaceId: 'team.b' })
    await ui.userShouldNotSeeSpaces({ stepUser: 'Brian', expectedSpaceIds: ['team.b'] })
    await ui.userLogsOut({ stepUser: 'Brian' })
  })

  test('multiple spaces can be managed at once in the admin settings via the batch actions', async () => {
    await api.userHasAssignedRolesToUsers({
      stepUser: 'Admin',
      users: [{ id: 'Alice', role: 'Space Admin' }]
    })
    await api.userHasCreatedProjectSpaces({
      stepUser: 'Alice',
      spaces: [
        { name: 'team A', id: 'team.a' },
        { name: 'team B', id: 'team.b' },
        { name: 'team C', id: 'team.c' },
        { name: 'team D', id: 'team.d' }
      ]
    })
    await ui.userLogsIn({ stepUser: 'Alice' })
    await ui.userOpensApplication({ stepUser: 'Alice', name: 'admin-settings' })
    await ui.userNavigatesToProjectSpaceManagementPage({ stepUser: 'Alice' })
    await ui.userDisablesSpacesUsingBatchActions({
      stepUser: 'Alice',
      spaceIds: ['team.a', 'team.b', 'team.c', 'team.d']
    })
    await ui.userEnablesSpacesUsingBatchActions({
      stepUser: 'Alice',
      spaceIds: ['team.a', 'team.b', 'team.c', 'team.d']
    })
    await ui.userChangesSpaceQuotaUsingBatchActions({
      stepUser: 'Alice',
      spaceIds: ['team.a', 'team.b', 'team.c', 'team.d'],
      value: '50'
    })
    await ui.userDisablesSpacesUsingBatchActions({
      stepUser: 'Alice',
      spaceIds: ['team.a', 'team.b', 'team.c', 'team.d']
    })
    await ui.userDeletesSpacesUsingBatchActions({
      stepUser: 'Alice',
      spaceIds: ['team.a', 'team.b', 'team.c', 'team.d']
    })
    await ui.userShouldNotSeeSpaces({
      stepUser: 'Alice',
      expectedSpaceIds: ['team.a', 'team.b', 'team.c', 'team.d']
    })
  })

  test('list members via sidebar', async () => {
    await api.usersHaveBeenCreated({
      stepUser: 'Admin',
      users: ['Brian', 'Carol', 'David', 'Edith']
    })
    await api.userHasAssignedRolesToUsers({
      stepUser: 'Admin',
      users: [{ id: 'Alice', role: 'Space Admin' }]
    })
    await api.userHasCreatedProjectSpaces({
      stepUser: 'Admin',
      spaces: [{ name: 'team A', id: 'team.a' }]
    })
    await api.userHasAddedMembersToSpace({
      stepUser: 'Admin',
      space: 'team A',
      sharee: [
        { user: 'Brian', shareType: 'user', role: 'Can edit with versions and trash bin' },
        { user: 'Carol', shareType: 'user', role: 'Can view' },
        { user: 'David', shareType: 'user', role: 'Can view' },
        { user: 'Edith', shareType: 'user', role: 'Can view' }
      ]
    })
    await ui.userLogsIn({ stepUser: 'Alice' })
    await ui.userOpensApplication({ stepUser: 'Alice', name: 'admin-settings' })
    await ui.userNavigatesToProjectSpaceManagementPage({ stepUser: 'Alice' })
    await ui.userListsMembersOfProjectSpaceUsingSidebarPanel({ stepUser: 'Alice', space: 'team.a' })
    await ui.userShouldSeeUsersInSidebarPanelOfSpacesAdminSettings({
      stepUser: 'Alice',
      expectedMembers: [
        { user: 'Admin', role: 'Can manage' },
        { user: 'Brian', role: 'Can edit with versions and trash bin' },
        { user: 'Carol', role: 'Can view' },
        { user: 'David', role: 'Can view' },
        { user: 'Edith', role: 'Can view' }
      ]
    })
  })

  test('admin user can manage the spaces created by other space admin user', async () => {
    await api.usersHaveBeenCreated({ stepUser: 'Admin', users: ['Brian', 'Carol'] })
    await api.userHasAssignedRolesToUsers({
      stepUser: 'Admin',
      users: [
        { id: 'Alice', role: 'Admin' },
        { id: 'Brian', role: 'Space Admin' },
        { id: 'Carol', role: 'Space Admin' }
      ]
    })
    await api.userHasCreatedProjectSpaces({
      stepUser: 'Brian',
      spaces: [{ name: 'team A', id: 'team.a' }]
    })
    await api.userHasCreatedProjectSpaces({
      stepUser: 'Carol',
      spaces: [{ name: 'team B', id: 'team.b' }]
    })
    await ui.userLogsIn({ stepUser: 'Alice' })
    await ui.userOpensApplication({ stepUser: 'Alice', name: 'admin-settings' })
    await ui.userNavigatesToProjectSpaceManagementPage({ stepUser: 'Alice' })
    await ui.userChangesSpaceQuotaUsingBatchActions({
      stepUser: 'Alice',
      spaceIds: ['team.a', 'team.b'],
      value: '50'
    })
    await ui.userDisablesSpacesUsingBatchActions({
      stepUser: 'Alice',
      spaceIds: ['team.a', 'team.b']
    })
    await ui.userEnablesSpacesUsingBatchActions({
      stepUser: 'Alice',
      spaceIds: ['team.a', 'team.b']
    })
    await ui.userDisablesSpacesUsingBatchActions({
      stepUser: 'Alice',
      spaceIds: ['team.a', 'team.b']
    })
    await ui.userDeletesSpacesUsingBatchActions({
      stepUser: 'Alice',
      spaceIds: ['team.a', 'team.b']
    })
    await ui.userShouldNotSeeSpaces({ stepUser: 'Alice', expectedSpaceIds: ['team.a', 'team.b'] })
  })
})
