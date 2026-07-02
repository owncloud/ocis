import { test } from '../../environment/test'
import * as api from '../../steps/api/api.js'
import * as ui from '../../steps/ui/index'
import { application, fileAction } from '../../environment/constants'

// For synchronization-related details, see https://owncloud.dev/services/proxy/#claim-updates
test.describe('groups management', () => {
  test('keycloak group sync with oCIS', async () => {
    // Given "Admin" creates following user using API
    //   | id    |
    //   | Alice |
    //   | Brian |
    await api.usersHaveBeenCreated({ stepUser: 'Admin', users: ['Alice', 'Brian'] })

    // And "Alice" creates the following files into personal space using API
    //   | pathToFile          | content              |
    //   | shareToSales.txt    | Keycloak group share |
    //   | shareToSecurity.txt | Keycloak group share |
    await api.userHasCreatedFiles({
      stepUser: 'Alice',
      files: [
        { pathToFile: 'shareToSales.txt', content: 'Keycloak group share' },
        { pathToFile: 'shareToSecurity.txt', content: 'Keycloak group share' }
      ]
    })

    // When "Admin" logs in
    await ui.userLogsIn({ stepUser: 'Admin' })

    // And "Admin" opens the "admin-settings" app
    await ui.userOpensApplication({ stepUser: 'Admin', name: 'admin-settings' })

    // And "Admin" navigates to the groups management page
    await ui.userNavigatesToGroupsManagementPage({ stepUser: 'Admin' })

    // When "Admin" creates the following groups
    //   | id       |
    //   | security |
    //   | sales    |
    await ui.userCreatesGroups({ stepUser: 'Admin', groupIds: ['security', 'sales'] })

    // Then "Admin" should see the following group
    //   | group            |
    //   | security         |
    //   | keycloak sales   |
    //   | keycloak finance |
    await ui.userShouldSeeGroupIds({
      stepUser: 'Admin',
      expectedGroupIds: ['security', 'keycloak sales', 'keycloak finance']
    })

    // When "Admin" navigates to the users management page
    await ui.userNavigatesToUserManagementPage({ stepUser: 'Admin' })
    // And "Admin" adds the user "Brian" to the groups "security,keycloak sales" using the sidebar panel
    await ui.userAddsUserToGroupsUsingContextMenu({
      stepUser: 'Admin',
      groups: ['security', 'keycloak sales'],
      user: 'Brian'
    })

    // And "Admin" logs out
    await ui.userLogsOut({ stepUser: 'Admin' })
    // And "Alice" logs in
    await ui.userLogsIn({ stepUser: 'Alice' })

    // And "Alice" shares the following resource using the sidebar panel
    //   | resource            | recipient      | type  | role                      | resourceType |
    //   | shareToSales.txt    | keycloak sales | group | Can edit without versions | file         |
    //   | shareToSecurity.txt | security       | group | Can edit without versions | file         |
    await ui.userSharesResources({
      stepUser: 'Alice',
      actionType: fileAction.sideBarPanel,
      shares: [
        {
          resource: 'shareToSales.txt',
          recipient: 'keycloak sales',
          type: 'group',
          role: 'Can edit with trashbin',
          resourceType: 'file'
        },
        {
          resource: 'shareToSecurity.txt',
          recipient: 'security',
          type: 'group',
          role: 'Can edit with trashbin',
          resourceType: 'file'
        }
      ]
    })

    // And "Alice" logs out
    await ui.userLogsOut({ stepUser: 'Alice' })

    // And "Brian" logs in
    await ui.userLogsIn({ stepUser: 'Brian' })
    // And "Brian" navigates to the shared with me page
    await ui.userNavigatesToSharedWithMePage({ stepUser: 'Brian' })

    // user should have access to unsynced shares
    // When "Brian" opens the following file in texteditor
    //   | resource         |
    //   | shareToSales.txt |
    await ui.userOpensResourceInViewer({
      stepUser: 'Brian',
      resource: 'shareToSales.txt',
      viewer: application.textEditor
    })
    // And "Brian" closes the file viewer
    await ui.userClosesFileViewer({ stepUser: 'Brian' })
    // And "Brian" edits the following resources
    //   | resource            | content     |
    //   | shareToSecurity.txt | new content |
    await ui.userEditsResources({
      stepUser: 'Brian',
      resources: [
        {
          name: 'shareToSecurity.txt',
          content: 'new content'
        }
      ]
    })
    // And "Brian" logs out
    await ui.userLogsOut({ stepUser: 'Brian' })

    // When "Admin" logs in
    await ui.userLogsIn({ stepUser: 'Admin' })

    // And "Admin" opens the "admin-settings" app
    await ui.userOpensApplication({ stepUser: 'Admin', name: 'admin-settings' })

    // And "Admin" navigates to the groups management page
    await ui.userNavigatesToGroupsManagementPage({ stepUser: 'Admin' })

    // Renaming a Keycloak group results in the creation of a new group on the oCIS server (see https://github.com/owncloud/ocis/issues/10445).
    // After renaming a group, it may take up to 5 minutes for the changes to sync, so avoid using the renamed group in the subsequent steps.
    // And "Admin" changes displayName to "a renamed group" for group "keycloak finance" using the sidebar panel
    await ui.userChangesGroup({
      stepUser: 'Admin',
      attribute: 'displayName',
      key: 'keycloak finance',
      value: 'a renamed group',
      action: 'context menu'
    })

    // When "Admin" deletes the following group using the context menu
    //   | group |
    //   | sales |
    await ui.userDeletesGroups({
      stepUser: 'Admin',
      actionType: fileAction.contextMenu,
      groupsToBeDeleted: ['sales']
    })
    // Then "Admin" should not see the following group
    //   | group |
    //   | sales |
    await ui.userShouldNotSeeGroupIds({ stepUser: 'Admin', expectedGroupIds: ['sales'] })
    // And "Admin" logs out
    await ui.userLogsOut({ stepUser: 'Admin' })
  })
})
