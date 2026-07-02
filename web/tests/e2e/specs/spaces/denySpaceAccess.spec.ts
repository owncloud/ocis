import { test } from '../../environment/test'
import * as ui from '../../steps/ui/index'
import * as api from '../../steps/api/api'
import { resourcePage, fileAction } from '../../environment/constants'

test.describe('deny space access', () => {
  test('deny and grant access', async () => {
    // Given "Admin" creates following users using API
    //   | id    |
    //   | Alice |
    //   | Brian |
    await api.usersHaveBeenCreated({ stepUser: 'Admin', users: ['Alice', 'Brian'] })

    // And "Admin" assigns following roles to the users using API
    //   | id    | role        |
    //   | Alice | Space Admin |
    await api.userHasAssignedRolesToUsers({
      stepUser: 'Admin',
      users: [{ id: 'Alice', role: 'Space Admin' }]
    })

    // And "Alice" logs in
    await ui.userLogsIn({ stepUser: 'Alice' })

    // And "Alice" creates the following project space using API
    //   | name  | id    |
    //   | sales | sales |
    await api.userHasCreatedProjectSpaces({
      stepUser: 'Alice',
      spaces: [{ name: 'sales', id: 'sales' }]
    })

    // And "Alice" creates the following folder in space "sales" using API
    //   | name |
    //   | f1   |
    //   | f2   |
    await api.userHasCreatedFoldersInSpace({
      stepUser: 'Alice',
      spaceName: 'sales',
      folders: ['f1', 'f2']
    })

    // And "Alice" adds the following members to the space "sales" using API
    //   | user  | role     | shareType |
    //   | Brian | Can edit | user      |
    await api.userHasAddedMembersToSpace({
      stepUser: 'Alice',
      space: 'sales',
      sharee: [{ user: 'Brian', role: 'Can edit with versions and trash bin', shareType: 'user' }]
    })

    // When "Alice" navigates to the project space "sales"
    await ui.userNavigatesToSpace({ stepUser: 'Alice', space: 'sales' })

    // When "Alice" shares the following resource using the sidebar panel
    //   | resource | recipient | type | role          | resourceType |
    //   | f1       | Brian     | user | Cannot access | folder       |
    await ui.userSharesResources({
      stepUser: 'Alice',
      actionType: fileAction.sideBarPanel,
      shares: [
        {
          resource: 'f1',
          recipient: 'Brian',
          type: 'user',
          role: 'Cannot access',
          resourceType: 'folder'
        }
      ]
    })

    // And "Brian" logs in
    await ui.userLogsIn({ stepUser: 'Brian' })

    // And "Brian" navigates to the project space "sales"
    await ui.userNavigatesToSpace({ stepUser: 'Brian', space: 'sales' })

    // Then following resources should not be displayed in the files list for user "Brian"
    //   | resource |
    //   | f1       |
    await ui.userShouldNotSeeTheResources({
      listType: resourcePage.filesList,
      stepUser: 'Brian',
      resources: ['f1']
    })

    // But following resources should be displayed in the files list for user "Brian"
    //   | resource |
    //   | f2       |
    await ui.userShouldSeeResources({
      listType: resourcePage.filesList,
      stepUser: 'Brian',
      resources: ['f2']
    })
    // allow access - deleting "Cannot access" share
    // When "Alice" removes following sharee
    //   | resource | recipient |
    //   | f1       | Brian     |
    await ui.userRemovesSharees({
      stepUser: 'Alice',
      sharees: [
        {
          resource: 'f1',
          recipient: 'Brian'
        }
      ]
    })

    // And "Brian" navigates to the project space "sales"
    await ui.userNavigatesToSpace({ stepUser: 'Brian', space: 'sales' })

    // Then following resources should be displayed in the files list for user "Brian"
    //   | resource |
    //   | f1       |
    //   | f2       |
    await ui.userShouldSeeResources({
      listType: resourcePage.filesList,
      stepUser: 'Brian',
      resources: ['f1', 'f2']
    })

    // And "Alice" logs out
    await ui.userLogsOut({ stepUser: 'Alice' })
    // And "Brian" logs out
    await ui.userLogsOut({ stepUser: 'Brian' })
  })
})
