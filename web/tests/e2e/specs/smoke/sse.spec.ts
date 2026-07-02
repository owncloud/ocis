import { test } from '../../environment/test'
import * as api from '../../steps/api/api'
import * as ui from '../../steps/ui/index'
import { fileAction, shareIndicator, resourcePage } from '../../environment/constants'

test.describe('server sent events', { tag: '@sse' }, () => {
  test.beforeEach(async () => {
    // Given "Admin" creates following users using API
    //   | id    |
    //   | Alice |
    //   | Brian |
    //   | Carol |
    await api.usersHaveBeenCreated({ stepUser: 'Admin', users: ['Alice', 'Brian', 'Carol'] })
  })

  test('space sse events', async () => {
    // Given "Admin" assigns following roles to the users using API
    //   | id    | role        |
    //   | Alice | Space Admin |
    await api.userHasAssignedRolesToUsers({
      stepUser: 'Admin',
      users: [{ id: 'Alice', role: 'Space Admin' }]
    })
    // And "Brian" logs in
    await ui.userLogsIn({ stepUser: 'Brian' })
    // And "Brian" navigates to the projects space page
    await ui.userNavigatesToSpacesPage({ stepUser: 'Brian' })

    // And "Alice" logs in
    await ui.userLogsIn({ stepUser: 'Alice' })
    // And "Alice" creates the following project space using API
    //   | name      | id        |
    //   | Marketing | marketing |
    await api.userHasCreatedProjectSpaces({
      stepUser: 'Alice',
      spaces: [{ name: 'Marketing', id: 'marketing' }]
    })
    //  space-member-added
    // When "Alice" adds the following members to the space "Marketing" using API
    //   | user  | role     | shareType |
    //   | Brian | Can view | user      |
    await api.userHasAddedMembersToSpace({
      stepUser: 'Alice',
      space: 'Marketing',
      sharee: [{ user: 'Brian', role: 'Can view', shareType: 'user' }]
    })
    // Then "Alice" should get "space-member-added" SSE event
    await ui.userShouldGetSSEEvent({ stepUser: 'Alice', event: 'space-member-added' })
    // And "Brian" should get "userlog-notification" SSE event
    await ui.userShouldGetSSEEvent({ stepUser: 'Brian', event: 'userlog-notification' })
    // And "Brian" should get "space-member-added" SSE event
    await ui.userShouldGetSSEEvent({ stepUser: 'Brian', event: 'space-member-added' })
    // And "Brian" should see space "marketing"
    await ui.userShouldSeeSpace({ stepUser: 'Brian', space: 'marketing' })

    // folder-created
    // When "Brian" navigates to the project space "marketing"
    await ui.userNavigatesToSpace({ stepUser: 'Brian', space: 'marketing' })
    // And "Alice" creates the following folder in space "Marketing" using API
    //   | name         |
    //   | space-folder |
    await api.userHasCreatedFoldersInSpace({
      stepUser: 'Alice',
      spaceName: 'Marketing',
      folders: ['space-folder']
    })
    // Then "Alice" should get "folder-created" SSE event
    await ui.userShouldGetSSEEvent({ stepUser: 'Alice', event: 'folder-created' })
    // And "Brian" should get "folder-created" SSE event
    await ui.userShouldGetSSEEvent({ stepUser: 'Brian', event: 'folder-created' })
    // And following resources should be displayed in the files list for user "Brian"
    //   | resource     |
    //   | space-folder |
    await ui.userShouldSeeResources({
      listType: resourcePage.filesList,
      stepUser: 'Brian',
      resources: ['space-folder']
    })
    // And "Brian" should not be able to edit folder "space-folder"
    await ui.userShouldNotBeAbleToEditResource({ stepUser: 'Brian', resource: 'space-folder' })

    // space-share-updated
    // When "Alice" navigates to the project space "marketing"
    await ui.userNavigatesToSpace({ stepUser: 'Alice', space: 'marketing' })
    // And "Alice" changes the roles of the following users in the project space
    //   | user  | role     |
    //   | Brian | Can edit |
    await ui.userChangesMemberRole({ stepUser: 'Alice', role: 'Can edit', sharee: 'Brian' })
    // Then "Alice" should get "space-share-updated" SSE event
    await ui.userShouldGetSSEEvent({ stepUser: 'Alice', event: 'space-share-updated' })
    // And "Brian" should get "space-share-updated" SSE event
    await ui.userShouldGetSSEEvent({ stepUser: 'Brian', event: 'space-share-updated' })
    // And "Brian" should be able to edit folder "space-folder"
    await ui.userShouldBeAbleToEditResource({ stepUser: 'Brian', resource: 'space-folder' })

    // share-created
    // When "Alice" shares the following resource using the sidebar panel
    //   | resource     | recipient | type | role     | resourceType |
    //   | space-folder | Carol     | user | Can view | folder       |
    await ui.userSharesResources({
      stepUser: 'Alice',
      actionType: fileAction.sideBarPanel,
      shares: [
        {
          resource: 'space-folder',
          recipient: 'Carol',
          type: 'user',
          role: 'Can view',
          resourceType: 'folder'
        }
      ]
    })

    // Then "Alice" should get "share-created" SSE event
    await ui.userShouldGetSSEEvent({ stepUser: 'Alice', event: 'share-created' })
    // And "Brian" should get "share-created" SSE event
    await ui.userShouldGetSSEEvent({ stepUser: 'Brian', event: 'share-created' })

    // And "Brian" closes the sidebar
    await ui.userClosesSidebar({ stepUser: 'Brian' })
    // And "Brian" should see user-direct indicator on the folder "space-folder"
    await ui.userShouldSeeShareIndicatorOnResource({
      stepUser: 'Brian',
      buttonLabel: shareIndicator.userDirect,
      resource: 'space-folder'
    })

    // share-updated
    // When "Alice" updates following sharee role
    //   | resource     | recipient | type | role     | resourceType |
    //   | space-folder | Carol     | user | Can view | folder       |
    await ui.userUpdatesShareeRoles({
      stepUser: 'Alice',
      roleUpdates: [
        {
          resource: 'space-folder',
          recipient: 'Carol',
          type: 'user',
          role: 'Can view',
          resourceType: 'folder'
        }
      ]
    })
    // Then "Alice" should get "share-updated" SSE event
    await ui.userShouldGetSSEEvent({ stepUser: 'Alice', event: 'share-updated' })
    // And "Brian" should get "share-updated" SSE event
    await ui.userShouldGetSSEEvent({ stepUser: 'Brian', event: 'share-updated' })

    // # link-created
    // When "Alice" creates a public link of following resource using the sidebar panel
    //   | resource     | password |
    //   | space-folder | %public% |
    await ui.userCreatesPublicLink({
      stepUser: 'Alice',
      resource: 'space-folder',
      password: '%public%'
    })
    // Then "Alice" should get "link-created" SSE event
    await ui.userShouldGetSSEEvent({ stepUser: 'Alice', event: 'link-created' })
    // Then "Brian" should get "link-created" SSE event
    await ui.userShouldGetSSEEvent({ stepUser: 'Brian', event: 'link-created' })
    // And "Brian" should see link-direct indicator on the folder "space-folder"
    await ui.userShouldSeeShareIndicatorOnResource({
      stepUser: 'Brian',
      buttonLabel: shareIndicator.linkDirect,
      resource: 'space-folder'
    })

    // link-updated
    // When "Alice" renames the most recently created public link of resource "space-folder" to "myLink"
    await ui.userRenamesMostRecentlyCreatedPublicLinkOfResource({
      stepUser: 'Alice',
      resource: 'space-folder',
      newName: 'myLink'
    })
    // Then "Alice" should get "link-updated" SSE event
    await ui.userShouldGetSSEEvent({ stepUser: 'Alice', event: 'link-updated' })
    // And "Brian" should get "link-updated" SSE event
    await ui.userShouldGetSSEEvent({ stepUser: 'Brian', event: 'link-updated' })

    // share-removed
    // When "Alice" removes following sharee
    //   | resource     | recipient |
    //   | space-folder | Carol     |
    await ui.userRemovesSharees({
      stepUser: 'Alice',
      sharees: [{ resource: 'space-folder', recipient: 'Carol' }]
    })
    // Then "Alice" should get "share-removed" SSE event
    await ui.userShouldGetSSEEvent({ stepUser: 'Alice', event: 'share-removed' })
    // And "Brian" should get "share-removed" SSE event
    await ui.userShouldGetSSEEvent({ stepUser: 'Brian', event: 'share-removed' })
    // And "Brian" should not see user-direct indicator on the folder "space-folder"
    await ui.userShouldNotSeeShareIndicatorOnResource({
      stepUser: 'Brian',
      buttonLabel: shareIndicator.userDirect,
      resource: 'space-folder'
    })

    // link-removed
    // When "Alice" removes the public link named "myLink" of resource "space-folder"
    await ui.userRemovesThePublicLinkOfResource({
      stepUser: 'Alice',
      linkName: 'myLink',
      resource: 'space-folder'
    })
    // Then "Alice" should get "link-removed" SSE event
    await ui.userShouldGetSSEEvent({ stepUser: 'Alice', event: 'link-removed' })
    // And "Brian" should get "link-removed" SSE event
    await ui.userShouldGetSSEEvent({ stepUser: 'Brian', event: 'link-removed' })
    // And "Brian" should not see link-direct indicator on the folder "space-folder"
    await ui.userShouldNotSeeShareIndicatorOnResource({
      stepUser: 'Brian',
      buttonLabel: shareIndicator.linkDirect,
      resource: 'space-folder'
    })

    // # space-member-removed
    // When "Brian" navigates to the projects space page
    await ui.userNavigatesToSpacesPage({ stepUser: 'Brian' })
    // And "Alice" navigates to the project space "marketing"
    await ui.userNavigatesToSpace({ stepUser: 'Alice', space: 'marketing' })
    // And "Alice" removes access to following users from the project space
    //   | user  |
    //   | Brian |
    await ui.userRemovesAccessToMember({ stepUser: 'Alice', reciver: 'Brian' })
    // Then "Alice" should get "space-member-removed" SSE event
    await ui.userShouldGetSSEEvent({ stepUser: 'Alice', event: 'space-member-removed' })
    // And "Brian" should get "space-member-removed" SSE event
    await ui.userShouldGetSSEEvent({ stepUser: 'Brian', event: 'space-member-removed' })
    // And "Brian" should not see space "marketing"
    await ui.userShouldNotSeeSpace({ stepUser: 'Brian', space: 'marketing' })

    // And "Brian" logs out
    await ui.userLogsOut({ stepUser: 'Brian' })
    // And "Alice" logs out
    await ui.userLogsOut({ stepUser: 'Alice' })
  })

  test('share sse events', async () => {
    // When "Brian" logs in
    await ui.userLogsIn({ stepUser: 'Brian' })
    // And "Brian" navigates to the shared with me page
    await ui.userNavigatesToSharedWithMePage({ stepUser: 'Brian' })
    // And "Alice" logs in
    await ui.userLogsIn({ stepUser: 'Alice' })
    // And "Alice" creates the following folder in personal space using API
    //   | name                   |
    //   | sharedFolder/subFolder |
    await api.userHasCreatedFolders({ stepUser: 'Alice', folderNames: ['sharedFolder/subFolder'] })

    // share-created
    // When "Alice" shares the following resource using the sidebar panel
    //   | resource     | recipient | type | role     | resourceType |
    //   | sharedFolder | Brian     | user | Can view | folder       |
    await ui.userSharesResources({
      stepUser: 'Alice',
      actionType: fileAction.sideBarPanel,
      shares: [
        {
          resource: 'sharedFolder',
          recipient: 'Brian',
          type: 'user',
          role: 'Can view',
          resourceType: 'folder'
        }
      ]
    })
    // Then "Alice" should get "share-created" SSE event
    await ui.userShouldGetSSEEvent({ stepUser: 'Alice', event: 'share-created' })
    // And "Brian" should get "share-created" SSE event
    await ui.userShouldGetSSEEvent({ stepUser: 'Brian', event: 'share-created' })
    // And "Brian" should not be able to edit folder "sharedFolder"
    await ui.userShouldNotBeAbleToEditResource({ stepUser: 'Brian', resource: 'sharedFolder' })

    // share-updated
    // When "Alice" updates following sharee role
    //   | resource     | recipient | type | role                   | resourceType |
    //   | sharedFolder | Brian     | user | Can edit with trashbin | folder       |
    await ui.userUpdatesShareeRoles({
      stepUser: 'Alice',
      roleUpdates: [
        {
          resource: 'sharedFolder',
          recipient: 'Brian',
          type: 'user',
          role: 'Can edit with trashbin',
          resourceType: 'folder'
        }
      ]
    })
    // Then "Alice" should get "share-updated" SSE event
    await ui.userShouldGetSSEEvent({ stepUser: 'Alice', event: 'share-updated' })
    // And "Brian" should get "share-updated" SSE event
    await ui.userShouldGetSSEEvent({ stepUser: 'Brian', event: 'share-updated' })
    // And "Brian" opens folder "sharedFolder"
    await ui.userOpensResource({ stepUser: 'Brian', resource: 'sharedFolder' })
    // And "Brian" should be able to edit folder "subFolder"
    await ui.userShouldBeAbleToEditResource({ stepUser: 'Brian', resource: 'subFolder' })

    // share-removed
    // When "Alice" removes following sharee
    //   | resource     | recipient |
    //   | sharedFolder | Brian     |
    await ui.userRemovesSharees({
      stepUser: 'Alice',
      sharees: [{ resource: 'sharedFolder', recipient: 'Brian' }]
    })

    // Then "Alice" should get "share-removed" SSE event
    await ui.userShouldGetSSEEvent({ stepUser: 'Alice', event: 'share-removed' })
    // And "Brian" should get "share-removed" SSE event
    await ui.userShouldGetSSEEvent({ stepUser: 'Brian', event: 'share-removed' })
    // And "Brian" should see the message "Your access to this share has been revoked. Please navigate to another location." on the webUI
    await ui.userShouldSeeMessageOnWebUI({
      stepUser: 'Brian',
      message: 'Your access to this share has been revoked. Please navigate to another location.'
    })
    // And "Brian" logs out
    await ui.userLogsOut({ stepUser: 'Brian' })
    // And "Alice" logs out
    await ui.userLogsOut({ stepUser: 'Alice' })
  })

  test('sse events on file operations', async () => {
    // Given "Admin" assigns following roles to the users using API
    //   | id    | role        |
    //   | Alice | Space Admin |
    await api.userHasAssignedRolesToUsers({
      stepUser: 'Admin',
      users: [{ id: 'Alice', role: 'Space Admin' }]
    })
    // And "Alice" creates the following project space using API
    //   | name      | id        |
    //   | Marketing | marketing |
    await api.userHasCreatedProjectSpaces({
      stepUser: 'Alice',
      spaces: [{ name: 'Marketing', id: 'marketing' }]
    })
    // And "Alice" adds the following members to the space "Marketing" using API
    //   | user  | role                   | shareType |
    //   | Brian | Can edit with trashbin | user      |
    await api.userHasAddedMembersToSpace({
      stepUser: 'Alice',
      space: 'Marketing',
      sharee: [{ user: 'Brian', role: 'Can edit with trashbin', shareType: 'user' }]
    })
    // And "Alice" creates the following folder in space "Marketing" using API
    //   | name         |
    //   | space-folder |
    await api.userHasCreatedFoldersInSpace({
      stepUser: 'Alice',
      spaceName: 'Marketing',
      folders: ['space-folder']
    })
    // And "Alice" logs in
    await ui.userLogsIn({ stepUser: 'Alice' })
    // When "Alice" navigates to the project space "marketing"
    await ui.userNavigatesToSpace({ stepUser: 'Alice', space: 'marketing' })

    // And "Brian" logs in
    await ui.userLogsIn({ stepUser: 'Brian' })
    // And "Brian" navigates to the project space "marketing"
    await ui.userNavigatesToSpace({ stepUser: 'Brian', space: 'marketing' })

    // postprocessing-finished - upload file
    // When "Brian" uploads the following resources
    //   | resource   |
    //   | simple.pdf |
    await ui.userUploadsResources({ stepUser: 'Brian', resources: [{ name: 'simple.pdf' }] })
    // Then "Brian" should get "postprocessing-finished" SSE event
    await ui.userShouldGetSSEEvent({ stepUser: 'Brian', event: 'postprocessing-finished' })
    // And "Alice" should get "postprocessing-finished" SSE event
    await ui.userShouldGetSSEEvent({ stepUser: 'Alice', event: 'postprocessing-finished' })
    // And following resources should be displayed in the files list for user "Alice"
    //   | resource   |
    //   | simple.pdf |
    await ui.userShouldSeeResources({
      listType: resourcePage.filesList,
      stepUser: 'Alice',
      resources: ['simple.pdf']
    })

    // postprocessing-finished - create file
    // file-touched -create file
    // When "Alice" creates the following resources
    //   | resource    | type    | content   |
    //   | example.txt | txtFile | some text |
    await ui.userCreatesResources({
      stepUser: 'Alice',
      resources: [{ name: 'example.txt', type: 'txtFile', content: 'some text' }]
    })
    // Then "Alice" should get "postprocessing-finished" SSE event
    await ui.userShouldGetSSEEvent({ stepUser: 'Alice', event: 'postprocessing-finished' })
    // And "Alice" should get "file-touched" SSE event
    await ui.userShouldGetSSEEvent({ stepUser: 'Alice', event: 'file-touched' })
    // And "Brian" should get "postprocessing-finished" SSE event
    await ui.userShouldGetSSEEvent({ stepUser: 'Brian', event: 'postprocessing-finished' })
    // And "Brian" should get "file-touched" SSE event
    await ui.userShouldGetSSEEvent({ stepUser: 'Brian', event: 'file-touched' })
    // And following resources should be displayed in the files list for user "Brian"
    //   | resource    |
    //   | example.txt |
    await ui.userShouldSeeResources({
      listType: resourcePage.filesList,
      stepUser: 'Brian',
      resources: ['example.txt']
    })

    // item-renamed
    // When "Brian" renames the following resource
    //   | resource   | as                 |
    //   | simple.pdf | simple-renamed.pdf |
    await ui.userRenamesResource({
      stepUser: 'Brian',
      resource: 'simple.pdf',
      newResourceName: 'simple-renamed.pdf'
    })
    // Then "Brian" should get "item-renamed" SSE event
    await ui.userShouldGetSSEEvent({ stepUser: 'Brian', event: 'item-renamed' })
    // And "Alice" should get "item-renamed" SSE event
    await ui.userShouldGetSSEEvent({ stepUser: 'Alice', event: 'item-renamed' })
    // And following resources should be displayed in the files list for user "Alice"
    //   | resource           |
    //   | simple-renamed.pdf |
    await ui.userShouldSeeResources({
      listType: resourcePage.filesList,
      stepUser: 'Alice',
      resources: ['simple-renamed.pdf']
    })

    // item-trashed
    // When "Alice" deletes the following resource using the sidebar panel
    //   | resource    |
    //   | example.txt |
    await ui.userDeletesResources({
      stepUser: 'Alice',
      actionType: fileAction.sideBarPanel,
      resources: [{ name: 'example.txt' }]
    })
    // Then "Alice" should get "item-trashed" SSE event
    await ui.userShouldGetSSEEvent({ stepUser: 'Alice', event: 'item-trashed' })
    // And "Brian" should get "item-trashed" SSE event
    await ui.userShouldGetSSEEvent({ stepUser: 'Brian', event: 'item-trashed' })
    // And following resources should not be displayed in the files list for user "Brian"
    //   | resource    |
    //   | example.txt |
    await ui.userShouldNotSeeTheResources({
      listType: resourcePage.filesList,
      stepUser: 'Brian',
      resources: ['example.txt']
    })

    // item-restored
    // When "Brian" navigates to the trashbin of the project space "marketing"
    await ui.userNavigatesToTrashbinOfSpace({ stepUser: 'Brian', space: 'marketing' })
    // And "Brian" restores the following resources from trashbin
    //   | resource    |
    //   | example.txt |
    await ui.userRestoresResourcesFromTrashbin({ stepUser: 'Brian', resources: ['example.txt'] })

    // Then "Brian" should get "item-restored" SSE event
    await ui.userShouldGetSSEEvent({ stepUser: 'Brian', event: 'item-restored' })
    // And "Alice" should get "item-restored" SSE event
    await ui.userShouldGetSSEEvent({ stepUser: 'Alice', event: 'item-restored' })
    // And following resources should be displayed in the files list for user "Alice"
    //   | resource    |
    //   | example.txt |
    await ui.userShouldSeeResources({
      listType: resourcePage.filesList,
      stepUser: 'Alice',
      resources: ['example.txt']
    })

    // # item-moved
    // When "Brian" navigates to the project space "marketing"
    await ui.userNavigatesToSpace({ stepUser: 'Brian', space: 'marketing' })
    // And "Brian" opens folder "space-folder"
    await ui.userOpensResource({ stepUser: 'Brian', resource: 'space-folder' })
    // And "Alice" moves the following resource using drag-drop
    //   | resource           | to           |
    //   | simple-renamed.pdf | space-folder |
    await ui.userMovesResources({
      stepUser: 'Alice',
      actionType: fileAction.dragDrop,
      resources: [{ resource: 'simple-renamed.pdf', to: 'space-folder' }]
    })
    // Then "Alice" should get "item-moved" SSE event
    await ui.userShouldGetSSEEvent({ stepUser: 'Alice', event: 'item-moved' })
    // And "Brian" should get "item-moved" SSE event
    await ui.userShouldGetSSEEvent({ stepUser: 'Brian', event: 'item-moved' })
    // And following resources should be displayed in the files list for user "Brian"
    //   | resource           |
    //   | simple-renamed.pdf |
    await ui.userShouldSeeResources({
      listType: resourcePage.filesList,
      stepUser: 'Brian',
      resources: ['simple-renamed.pdf']
    })

    // And "Brian" logs out
    await ui.userLogsOut({ stepUser: 'Brian' })
    // And "Alice" logs out
    await ui.userLogsOut({ stepUser: 'Alice' })
  })
})
