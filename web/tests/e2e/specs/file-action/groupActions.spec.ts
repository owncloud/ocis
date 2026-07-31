import { test } from '../../environment/test'
import * as ui from '../../steps/ui/index'
import * as api from '../../steps/api/api'
import { fileAction } from '../../environment/constants'

test.describe('Group actions', { tag: '@predefined-users' }, () => {
  test.beforeEach(async () => {
    // Given "Admin" creates following user using API
    //   | id    |
    //   | Alice |
    //   | Brian |
    //   | Carol |
    //   | David |
    //   | Edith |
    await api.usersHaveBeenCreated({
      stepUser: 'Admin',
      users: ['Alice', 'Brian', 'Carol', 'David', 'Edith']
    })
    // And "Admin" creates following group using API
    //   | id       |
    //   | sales    |
    //   | finance  |
    //   | security |
    await api.groupsHaveBeenCreated({
      groupIds: ['sales', 'finance', 'security'],
      stepUser: 'Admin'
    })
    // And "Admin" adds user to the group using API
    //   | user  | group    |
    //   | Brian | sales    |
    //   | Brian | finance  |
    //   | Brian | security |
    await api.usersHaveBeenAddedToGroup({
      stepUser: 'Admin',
      usersToAdd: [
        { user: 'Brian', group: 'sales' },
        { user: 'Brian', group: 'finance' },
        { user: 'Brian', group: 'security' }
      ]
    })
    // And "Brian" logs in
    await ui.userLogsIn({ stepUser: 'Brian' })
  })

  test('batch share a resource to multiple users and groups', async () => {
    // disabling auto accepting to check accepting share
    // And "Brian" disables auto-accepting using API
    await api.userHasDisabledAutoAcceptingShare({ stepUser: 'Brian' })

    // And "Alice" creates the following folders in personal space using API
    //   | name                   |
    //   | sharedFolder           |
    //   | folder1                |
    //   | folder2                |
    //   | folder3                |
    //   | folder4                |
    //   | folder5                |
    //   | parentFolder/SubFolder |
    await api.userHasCreatedFolders({
      stepUser: 'Alice',
      folderNames: [
        'sharedFolder',
        'folder1',
        'folder2',
        'folder3',
        'folder4',
        'folder5',
        'parentFolder/SubFolder'
      ]
    })
    //  And "Alice" shares the following resource using API
    // | resource     | recipient | type | role                   | resourceType |
    // | folder1      | Brian     | user | Can edit with trashbin | folder       |
    // | folder2      | Brian     | user | Can edit with trashbin | folder       |
    // | folder3      | Brian     | user | Can edit with trashbin | folder       |
    // | folder4      | Brian     | user | Can edit with trashbin | folder       |
    // | folder5      | Brian     | user | Can edit with trashbin | folder       |
    // | parentFolder | Brian     | user | Can edit with trashbin | folder       |
    await api.userHasSharedResources({
      stepUser: 'Alice',
      shares: [
        {
          resource: 'folder1',
          recipient: 'Brian',
          type: 'user',
          role: 'Can edit with trashbin',
          resourceType: 'folder'
        },
        {
          resource: 'folder2',
          recipient: 'Brian',
          type: 'user',
          role: 'Can edit with trashbin',
          resourceType: 'folder'
        },
        {
          resource: 'folder3',
          recipient: 'Brian',
          type: 'user',
          role: 'Can edit with trashbin',
          resourceType: 'folder'
        },
        {
          resource: 'folder4',
          recipient: 'Brian',
          type: 'user',
          role: 'Can edit with trashbin',
          resourceType: 'folder'
        },
        {
          resource: 'folder5',
          recipient: 'Brian',
          type: 'user',
          role: 'Can edit with trashbin',
          resourceType: 'folder'
        },
        {
          resource: 'parentFolder',
          recipient: 'Brian',
          type: 'user',
          role: 'Can edit with trashbin',
          resourceType: 'folder'
        }
      ]
    })
    // And "Alice" logs in
    await ui.userLogsIn({ stepUser: 'Alice' })

    // # multiple share
    // And "Alice" shares the following resources using the sidebar panel
    // | resource     | recipient | type  | role                   | resourceType |
    // | sharedFolder | Brian     | user  | Can edit with trashbin | folder       |
    // | sharedFolder | Carol     | user  | Can edit with trashbin | folder       |
    // | sharedFolder | David     | user  | Can edit with trashbin | folder       |
    // | sharedFolder | Edith     | user  | Can edit with trashbin | folder       |
    // | sharedFolder | sales     | group | Can edit with trashbin | folder       |
    // | sharedFolder | finance   | group | Can edit with trashbin | folder       |
    // | sharedFolder | security  | group | Can edit with trashbin | folder       |
    await ui.userSharesResources({
      actionType: fileAction.sideBarPanel,
      stepUser: 'Alice',
      shares: [
        {
          resource: 'sharedFolder',
          recipient: 'Brian',
          type: 'user',
          role: 'Can edit with trashbin',
          resourceType: 'folder'
        },
        {
          resource: 'sharedFolder',
          recipient: 'Carol',
          type: 'user',
          role: 'Can edit with trashbin',
          resourceType: 'folder'
        },
        {
          resource: 'sharedFolder',
          recipient: 'David',
          type: 'user',
          role: 'Can edit with trashbin',
          resourceType: 'folder'
        },
        {
          resource: 'sharedFolder',
          recipient: 'Edith',
          type: 'user',
          role: 'Can edit with trashbin',
          resourceType: 'folder'
        },
        {
          resource: 'sharedFolder',
          recipient: 'sales',
          type: 'group',
          role: 'Can edit with trashbin',
          resourceType: 'folder'
        },
        {
          resource: 'sharedFolder',
          recipient: 'finance',
          type: 'group',
          role: 'Can edit with trashbin',
          resourceType: 'folder'
        },
        {
          resource: 'sharedFolder',
          recipient: 'security',
          type: 'group',
          role: 'Can edit with trashbin',
          resourceType: 'folder'
        }
      ]
    })

    // And "Brian" navigates to the shared with me page
    await ui.userNavigatesToSharedWithMePage({ stepUser: 'Brian' })

    // And "Brian" enables the sync for all shares using the batch action
    await ui.userEnablesSyncForAllShares({ stepUser: 'Brian' })

    // And "Alice" logs out
    await ui.userLogsOut({ stepUser: 'Alice' })
    // And "Brian" logs out
    await ui.userLogsOut({ stepUser: 'Brian' })
  })
})
