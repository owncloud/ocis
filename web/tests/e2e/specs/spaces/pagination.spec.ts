import { test } from '../../environment/test'
import * as api from '../../steps/api/api'
import * as ui from '../../steps/ui/index'

test.describe('check files pagination in project space', () => {
  test('pagination', async () => {
    // Given "Admin" creates following user using API
    //   | id    |
    //   | Alice |
    await api.usersHaveBeenCreated({ stepUser: 'Admin', users: ['Alice'] })

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
    //   | name       | id    |
    //   | Developers | dev.1 |
    await api.userHasCreatedProjectSpaces({
      stepUser: 'Alice',
      spaces: [{ name: 'Developers', id: 'dev.1' }]
    })

    // And "Alice" creates 55 folders in space "Developers" using API
    await api.userHasCreatedFoldersInSpace({
      stepUser: 'Alice',
      spaceName: 'Developers',
      folders: Array.from({ length: 55 }, (_, i) => `testFolder${i + 1}`)
    })

    // And "Alice" creates 55 files in space "Developers" using API
    await api.userHasCreatedFilesInsideSpace({
      stepUser: 'Alice',
      files: Array.from({ length: 55 }, (_, i) => ({
        name: `testfile${i + 1}.txt`,
        space: 'Developers',
        content: `This is a test file${i + 1}`
      }))
    })

    // And "Alice" creates the following file in space "Developers" using API
    //   | name                 | content                |
    //   | .hidden-testFile.txt | This is a hidden file. |
    await api.userHasCreatedFilesInsideSpace({
      stepUser: 'Alice',
      files: [
        { name: '.hidden-testFile.txt', space: 'Developers', content: 'This is a hidden file.' }
      ]
    })

    // And "Alice" navigates to the project space "dev.1"
    await ui.userNavigatesToSpace({ stepUser: 'Alice', space: 'dev.1' })

    // When "Alice" navigates to page "2" of the project space files view
    await ui.userNavigatesToPageNumber({ stepUser: 'Alice', pageNumber: '2' })

    // Then "Alice" should see the text "112 items with 1 kB in total (56 files including 1 hidden, 56 folders including 1 hidden)" at the footer of the page
    await ui.userShouldSeeFooterText({
      stepUser: 'Alice',
      expectedText:
        '112 items with 1 kB in total (56 files including 1 hidden, 56 folders including 1 hidden)'
    })

    // And "Alice" should see 10 resources in the project space files view
    await ui.userShouldSeeNumberOfResources({ stepUser: 'Alice', expectedNumberOfResources: 10 })

    // When "Alice" enables the option to display the hidden file
    await ui.userEnablesShowHiddenFilesOption({ stepUser: 'Alice' })

    // And "Alice" should see 12 resources in the project space files view
    await ui.userShouldSeeNumberOfResources({ stepUser: 'Alice', expectedNumberOfResources: 12 })

    // When "Alice" opens file "testfile45.txt"
    await ui.userOpensResource({ stepUser: 'Alice', resource: 'testfile45.txt' })

    // And "Alice" closes the file viewer
    await ui.userClosesFileViewer({ stepUser: 'Alice' })

    // Then "Alice" should be on page "2"
    await ui.userShouldBeOnPage({ stepUser: 'Alice', pageNumber: '2' })

    // When "Alice" changes the items per page to "500"
    await ui.userChangesItemsPerPage({ stepUser: 'Alice', itemsPerPage: '500' })

    // Then "Alice" should not see the pagination in the project space files view
    await ui.userShouldNotSeePagination({ stepUser: 'Alice' })

    // When "Alice" enables flat list
    await ui.userEnablesFlatList({ stepUser: 'Alice' })

    // Then "Alice" should see files being sorted in alphabetic order
    await ui.userShouldSeeFilesSortedAlphabetically({ stepUser: 'Alice' })

    // And "Alice" logs out
    await ui.userLogsOut({ stepUser: 'Alice' })
  })
})
