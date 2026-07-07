import { test } from '../../environment/test'
import * as ui from '../../steps/ui/index'
import * as api from '../../steps/api/api'
import { fileAction } from '../../environment/constants'

test.describe('groups management', () => {
  test.beforeEach(async () => {
    await ui.userLogsIn({ stepUser: 'Admin' })
  })

  test.afterEach(async () => {
    await ui.userLogsOut({ stepUser: 'Admin' })
  })

  test('admin creates group', async () => {
    await ui.userOpensApplication({ stepUser: 'Admin', name: 'admin-settings' })
    await ui.userNavigatesToGroupsManagementPage({ stepUser: 'Admin' })
    await ui.userCreatesGroups({ stepUser: 'Admin', groupIds: ['sales', 'security'] })
    await ui.userShouldSeeGroupIds({ stepUser: 'Admin', expectedGroupIds: ['sales', 'security'] })
  })

  test('admin deletes group', async () => {
    await api.groupsHaveBeenCreated({
      groupIds: ['sales', 'security', 'finance'],
      stepUser: 'Admin'
    })
    await ui.userOpensApplication({ stepUser: 'Admin', name: 'admin-settings' })
    await ui.userNavigatesToGroupsManagementPage({ stepUser: 'Admin' })
    await ui.userDeletesGroups({
      stepUser: 'Admin',
      actionType: fileAction.contextMenu,
      groupsToBeDeleted: ['sales']
    })

    await ui.userShouldNotSeeGroupIds({ stepUser: 'Admin', expectedGroupIds: ['sales'] })

    await ui.userDeletesGroups({
      stepUser: 'Admin',
      actionType: fileAction.batchAction,
      groupsToBeDeleted: ['security', 'finance']
    })

    await ui.userShouldNotSeeGroupIds({
      stepUser: 'Admin',
      expectedGroupIds: ['security', 'finance']
    })
  })

  test('edit groups', async () => {
    await api.groupsHaveBeenCreated({ groupIds: ['sales'], stepUser: 'Admin' })
    await ui.userOpensApplication({ stepUser: 'Admin', name: 'admin-settings' })
    await ui.userNavigatesToGroupsManagementPage({ stepUser: 'Admin' })
    await ui.userChangesGroup({
      stepUser: 'Admin',
      key: 'sales',
      attribute: 'displayName',
      value: 'a renamed group',
      action: fileAction.contextMenu
    })
    await ui.userShouldSeeGroupDisplayName({
      stepUser: 'Admin',
      groupDisplayName: 'a renamed group'
    })
  })
})
