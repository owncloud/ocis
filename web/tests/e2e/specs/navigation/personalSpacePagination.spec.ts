import { test } from '../../environment/test'
import * as api from '../../steps/api/api.js'
import * as ui from '../../steps/ui/index'

test.describe('Personal space pagination', { tag: '@predefined-users' }, () => {
  test('pagination', async () => {
    // Given "Admin" creates following user using API
    //   | id    |
    //   | Alice |
    await api.usersHaveBeenCreated({ stepUser: 'Admin', users: ['Alice'] })

    // And "Alice" logs in
    await ui.userLogsIn({ stepUser: 'Alice' })

    // And "Alice" creates 15 folders in personal space using API
    await api.userHasCreatedFolders({
      stepUser: 'Alice',
      folderNames: Array.from({ length: 15 }, (_, i) => `folder${i + 1}`)
    })
    // And "Alice" creates 10 files in personal space using API
    await api.userHasCreatedFiles({
      stepUser: 'Alice',
      files: Array.from({ length: 10 }, (_, i) => ({
        pathToFile: `file${i + 1}`,
        content: `This is a test file${i + 1}`
      }))
    })
    // And "Alice" creates the following files into personal space using API
    //   | pathToFile           | content                |
    //   | .hidden-testFile.txt | This is a hidden file. |
    await api.userHasCreatedFiles({
      stepUser: 'Alice',
      files: [{ pathToFile: '.hidden-testFile.txt', content: 'This is a hidden file.' }]
    })
    // When "Alice" opens the "files" app
    await ui.userOpensApplication({ stepUser: 'Alice', name: 'files' })
    // And "Alice" changes the items per page to "20"
    await ui.userChangesItemsPerPage({ stepUser: 'Alice', itemsPerPage: '20' })
    // Then "Alice" should see the text "26 items with 223 B in total (11 files including 1 hidden, 15 folders)" at the footer of the page
    await ui.userShouldSeeFooterText({
      stepUser: 'Alice',
      expectedText: '26 items with 223 B in total (11 files including 1 hidden, 15 folders)'
    })
    // When "Alice" navigates to page "2" of the personal space files view
    await ui.userNavigatesToPageNumber({ stepUser: 'Alice', pageNumber: '2' })
    // Then "Alice" should see 5 resources in the personal space files view
    await ui.userShouldSeeNumberOfResources({ stepUser: 'Alice', expectedNumberOfResources: 5 })
    // When "Alice" enables the option to display the hidden file
    await ui.userEnablesShowHiddenFilesOption({ stepUser: 'Alice' })
    // Then "Alice" should see 6 resources in the personal space files view
    await ui.userShouldSeeNumberOfResources({ stepUser: 'Alice', expectedNumberOfResources: 6 })
    // When "Alice" changes the items per page to "500"
    await ui.userChangesItemsPerPage({ stepUser: 'Alice', itemsPerPage: '500' })

    // Then "Alice" should not see the pagination in the personal space files view
    await ui.userShouldNotSeePagination({ stepUser: 'Alice' })

    // And "Alice" logs out
    await ui.userLogsOut({ stepUser: 'Alice' })
  })
})
