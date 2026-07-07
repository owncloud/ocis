import { test } from '../../environment/test'
import * as api from '../../steps/api/api'
import * as ui from '../../steps/ui/index'
import { fileAction } from '../../environment/constants'

test.describe('Users can use web to organize tags', () => {
  test.beforeEach(async () => {
    // Given "Admin" creates following users using API
    //   | id    |
    //   | Alice |
    //   | Brian |
    await api.usersHaveBeenCreated({ stepUser: 'Admin', users: ['Alice', 'Brian'] })

    // Given "Alice" logs in
    await ui.userLogsIn({ stepUser: 'Alice' })
  })

  test('Tag management', { tag: '@predefined-users' }, async () => {
    // When "Alice" creates the following files into personal space using API
    //   | pathToFile | content     |
    //   | lorem.txt  | lorem ipsum |
    await api.userHasCreatedFiles({
      stepUser: 'Alice',
      files: [{ pathToFile: 'lorem.txt', content: 'lorem ipsum' }]
    })
    // And "Alice" adds the following tags for the following resources using the sidebar panel
    //   | resource  | tags         |
    //   | lorem.txt | tag 1, tag 2 |
    await ui.userAddsFollowingTagsForResourcesUsingSidebarPanel({
      stepUser: 'Alice',
      resources: [{ name: 'lorem.txt', tags: ['tag 1', 'tag 2'] }]
    })
    // Then the following resources should contain the following tags in the files list for user "Alice"
    //   | resource  | tags         |
    //   | lorem.txt | tag 1, tag 2 |
    await ui.resourceShouldContainTagsInFileList({
      stepUser: 'Alice',
      resources: [{ name: 'lorem.txt', tags: ['tag 1', 'tag 2'] }]
    })
    // Then the following resources should contain the following tags in the details panel for user "Alice"
    //   | resource  | tags         |
    //   | lorem.txt | tag 1, tag 2 |
    await ui.resourceShouldContainTagsInDetailPanel({
      stepUser: 'Alice',
      resources: [{ name: 'lorem.txt', tags: ['tag 1', 'tag 2'] }]
    })
    // When "Alice" removes the following tags for the following resources using the sidebar panel
    //   | resource  | tags  |
    //   | lorem.txt | tag 1 |
    await ui.userRemovesTagsFromResourcesUsingSideBar({
      stepUser: 'Alice',
      resources: [{ name: 'lorem.txt', tags: ['tag 1'] }]
    })
    // Then the following resources should contain the following tags in the files list for user "Alice"
    //   | resource  | tags  |
    //   | lorem.txt | tag 2 |
    await ui.resourceShouldContainTagsInFileList({
      stepUser: 'Alice',
      resources: [{ name: 'lorem.txt', tags: ['tag 2'] }]
    })
    // Then the following resources should contain the following tags in the details panel for user "Alice"
    //   | resource  | tags  |
    //   | lorem.txt | tag 2 |
    await ui.resourceShouldContainTagsInDetailPanel({
      stepUser: 'Alice',
      resources: [{ name: 'lorem.txt', tags: ['tag 2'] }]
    })
    // And "Alice" logs out
    await ui.userLogsOut({ stepUser: 'Alice' })
  })

  test('Long tag name', async () => {
    // And "Alice" creates the following files into personal space using API
    //   | pathToFile | content     |
    //   | lorem.txt  | lorem ipsum |
    await api.userHasCreatedFiles({
      stepUser: 'Alice',
      files: [{ pathToFile: 'lorem.txt', content: 'lorem ipsum' }]
    })
    // When "Alice" tries to add the following tag for the following resources using the sidebar panel
    //   | resource  | tags                                                                                                       |
    //   | lorem.txt | Loremipsum dolor sit amet consectetur adipiscing elit sed do eiusmod tempor incididunt ut labore et dolore |
    await ui.userTriesToAddTagForResourceUsingSidebarPanel({
      stepUser: 'Alice',
      resources: [
        {
          name: 'lorem.txt',
          tags: [
            'Loremipsum dolor sit amet consectetur adipiscing elit sed do eiusmod tempor incididunt ut labore et dolore'
          ]
        }
      ]
    })
    // Then "Alice" should see the following tag validation message:
    //   """
    //   Tags must not be longer than 100 characters
    //   """
    await ui.userShouldSeeFollowingTagValidationMessages({
      stepUser: 'Alice',
      message: 'Tags must not be longer than 100 characters'
    })
    // And "Alice" logs out
    await ui.userLogsOut({ stepUser: 'Alice' })
  })

  test('Tag search', { tag: '@predefined-users' }, async () => {
    // When "Alice" creates the following files into personal space using API
    //   | pathToFile   | content     |
    //   | lorem.txt    | lorem ipsum |
    //   | textfile.txt | test file   |
    await api.userHasCreatedFiles({
      stepUser: 'Alice',
      files: [
        { pathToFile: 'lorem.txt', content: 'lorem ipsum' },
        { pathToFile: 'textfile.txt', content: 'test file' }
      ]
    })
    // And "Alice" adds the following tags for the following resources using the sidebar panel
    //   | resource  | tags       |
    //   | lorem.txt | tag1, tag2 |
    await ui.userAddsFollowingTagsForResourcesUsingSidebarPanel({
      stepUser: 'Alice',
      resources: [{ name: 'lorem.txt', tags: ['tag1', 'tag2'] }]
    })
    // And "Alice" clicks the tag "tag1" on the resource "lorem.txt"
    await ui.userClicksTheTagOnResource({
      stepUser: 'Alice',
      resourceName: 'lorem.txt',
      tagName: 'tag1'
    })
    // Then the following resources should contain the following tags in the files list for user "Alice"
    //   | resource  | tags |
    //   | lorem.txt | tag1 |
    await ui.resourceShouldContainTagsInFileList({
      stepUser: 'Alice',
      resources: [{ name: 'lorem.txt', tags: ['tag1'] }]
    })
    // Then following resources should not be displayed in the files list for user "Alice"
    //   | resource     |
    //   | textfile.txt |
    await ui.userShouldNotSeeTheResources({
      listType: 'files list',
      stepUser: 'Alice',
      resources: ['textfile.txt']
    })
    // And "Alice" logs out
    await ui.userLogsOut({ stepUser: 'Alice' })
  })

  test('Tags in shared resources', { tag: '@predefined-users' }, async () => {
    // And "Brian" logs in
    await ui.userLogsIn({ stepUser: 'Brian' })
    // And "Alice" creates the following folders in personal space using API
    //   | name             |
    //   | folder_to_shared |
    await api.userHasCreatedFolders({ stepUser: 'Alice', folderNames: ['folder_to_shared'] })
    // And "Alice" creates the following files into personal space using API
    //   | pathToFile                 | content     |
    //   | folder_to_shared/lorem.txt | lorem ipsum |
    await api.userHasCreatedFiles({
      stepUser: 'Alice',
      files: [{ pathToFile: 'folder_to_shared/lorem.txt', content: 'lorem ipsum' }]
    })
    // And "Alice" adds the following tags for the following resources using the sidebar panel
    //   | resource                   | tags         |
    //   | folder_to_shared/lorem.txt | tag 1, tag 2 |
    await ui.userAddsFollowingTagsForResourcesUsingSidebarPanel({
      stepUser: 'Alice',
      resources: [{ name: 'folder_to_shared/lorem.txt', tags: ['tag 1', 'tag 2'] }]
    })
    // When "Alice" shares the following resource using the sidebar panel
    //   | resource         | recipient | type | role                   | resourceType |
    //   | folder_to_shared | Brian     | user | Can edit with trashbin | folder       |
    await ui.userSharesResources({
      stepUser: 'Alice',
      actionType: fileAction.sideBarPanel,
      shares: [
        {
          resource: 'folder_to_shared',
          recipient: 'Brian',
          type: 'user',
          role: 'Can edit with trashbin',
          resourceType: 'folder'
        }
      ]
    })
    // And "Alice" logs out
    await ui.userLogsOut({ stepUser: 'Alice' })
    // And "Brian" navigates to the shared with me page
    await ui.userNavigatesToSharedWithMePage({ stepUser: 'Brian' })
    // Then the following resources should contain the following tags in the files list for user "Brian"
    //   | resource                   | tags         |
    //   | folder_to_shared/lorem.txt | tag 1, tag 2 |
    await ui.resourceShouldContainTagsInFileList({
      stepUser: 'Brian',
      resources: [{ name: 'folder_to_shared/lorem.txt', tags: ['tag 1', 'tag 2'] }]
    })
    // And "Brian" logs out
    await ui.userLogsOut({ stepUser: 'Brian' })
  })
})
