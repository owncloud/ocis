import { test } from '../../environment/test'
import * as api from '../../steps/api/api'
import * as ui from '../../steps/ui/index'
import { fileAction, resourcePage, searchFilter, searchScope } from '../../environment/constants'

test.describe('Search', { tag: '@predefined-users' }, () => {
  test.beforeEach(async () => {
    // Given "Admin" creates following users using API
    //   | id    |
    //   | Alice |
    await api.usersHaveBeenCreated({ stepUser: 'Admin', users: ['Alice'] })
  })

  test('Search in personal spaces', async ({ world }) => {
    // Given "Admin" creates following users using API
    //   | id    |
    //   | Brian |
    await api.usersHaveBeenCreated({ stepUser: 'Admin', users: ['Brian'] })
    // And "Brian" logs in
    await ui.userLogsIn({ stepUser: 'Brian' })

    // And "Alice" logs in
    await ui.userLogsIn({ stepUser: 'Alice' })

    // And "Brian" creates the following folder in personal space using API
    //   | name                 |
    //   | new_share_from_brian |
    await api.userHasCreatedFolder({ stepUser: 'Brian', folderName: 'new_share_from_brian' })

    // And "Brian" uploads the following local file into personal space using API
    //   | localFile                        | to                |
    //   | filesForUpload/new-lorem-big.txt | new-lorem-big.txt |
    await api.userHasUploadedFilesInPersonalSpace({
      stepUser: 'Brian',
      filesToUpload: [{ localFile: 'filesForUpload/new-lorem-big.txt', to: 'new-lorem-big.txt' }]
    })
    // And "Brian" shares the following resource using the sidebar panel
    //   | resource             | recipient | type | role     | resourceType |
    //   | new_share_from_brian | Alice     | user | Can view | folder       |
    //   | new-lorem-big.txt    | Alice     | user | Can view | file         |
    await ui.userSharesResources({
      stepUser: 'Brian',
      actionType: fileAction.sideBarPanel,
      shares: [
        {
          resource: 'new_share_from_brian',
          recipient: 'Alice',
          type: 'user',
          role: 'Can view',
          resourceType: 'folder'
        },
        {
          resource: 'new-lorem-big.txt',
          recipient: 'Alice',
          type: 'user',
          role: 'Can view',
          resourceType: 'file'
        }
      ]
    })
    // And "Brian" logs out
    await ui.userLogsOut({ stepUser: 'Brian' })

    // And "Alice" creates the following resources
    //   | resource                   | type   |
    //   | folder                     | folder |
    //   | FolDer/child-one/child-two | folder |
    //   | strängéनेपालीName          | folder |
    await ui.userCreatesResources({
      stepUser: 'Alice',
      resources: [
        { name: 'folder', type: 'folder' },
        { name: 'FolDer/child-one/child-two', type: 'folder' },
        { name: 'strängéनेपालीName', type: 'folder' }
      ]
    })

    // And "Alice" enables the option to display the hidden file
    await ui.userEnablesShowHiddenFilesOption({ stepUser: 'Alice' })

    // And "Alice" uploads the following resources
    //   | resource         |
    //   | .hidden-file.txt |
    await ui.userUploadsResources({ stepUser: 'Alice', resources: [{ name: '.hidden-file.txt' }] })

    // # search for objects of personal space
    // When "Alice" searches "foldeR" using the global search and the "all files" filter
    await ui.userSearchesGloballyWithFilter({
      stepUser: 'Alice',
      keyword: 'foldeR',
      filter: searchScope.allFiles
    })
    // Then following resources should be displayed in the search list for user "Alice"
    //   | resource |
    //   | folder   |
    //   | FolDer   |
    await ui.userShouldSeeResources({
      listType: resourcePage.searchList,
      stepUser: 'Alice',
      resources: ['folder', 'FolDer']
    })

    // But following resources should not be displayed in the search list for user "Alice"
    //   | resource             |
    //   | new_share_from_brian |
    //   | new-lorem-big.txt    |
    //   | .hidden-file.txt     |
    await ui.userShouldNotSeeTheResources({
      listType: resourcePage.searchList,
      stepUser: 'Alice',
      resources: ['new_share_from_brian', 'new-lorem-big.txt', '.hidden-file.txt']
    })

    // # search for hidden file
    // When "Alice" searches "hidden" using the global search and the "all files" filter
    await ui.userSearchesGloballyWithFilter({
      stepUser: 'Alice',
      keyword: 'hidden',
      filter: searchScope.allFiles
    })
    // Then following resources should be displayed in the search list for user "Alice"
    //   | resource         |
    //   | .hidden-file.txt |
    await ui.userShouldSeeResources({
      listType: resourcePage.searchList,
      stepUser: 'Alice',
      resources: ['.hidden-file.txt']
    })
    // But following resources should not be displayed in the search list for user "Alice"
    //   | resource          |
    //   | folder            |
    //   | FolDer            |
    //   | PARENT            |
    //   | new-lorem-big.txt |
    await ui.userShouldNotSeeTheResources({
      listType: resourcePage.searchList,
      stepUser: 'Alice',
      resources: ['folder', 'FolDer', 'PARENT', 'new-lorem-big.txt']
    })

    // # subfolder search
    // And "Alice" searches "child" using the global search and the "all files" filter
    await ui.userSearchesGloballyWithFilter({
      stepUser: 'Alice',
      keyword: 'child',
      filter: searchScope.allFiles
    })
    // Then following resources should be displayed in the search list for user "Alice"
    //   | resource  |
    //   | child-one |
    //   | child-two |
    await ui.userShouldSeeResources({
      listType: resourcePage.searchList,
      stepUser: 'Alice',
      resources: ['child-one', 'child-two']
    })
    // But following resources should not be displayed in the search list for user "Alice"
    //   | resource          |
    //   | folder            |
    //   | FolDer            |
    //   | folder_from_brian |
    //   | .hidden-file.txt  |
    //   | new-lorem-big.txt |
    await ui.userShouldNotSeeTheResources({
      listType: resourcePage.searchList,
      stepUser: 'Alice',
      resources: ['folder', 'FolDer', 'folder_from_brian', '.hidden-file.txt', 'new-lorem-big.txt']
    })

    // # received shares search
    // And "Alice" searches "NEW" using the global search and the "all files" filter
    await ui.userSearchesGloballyWithFilter({
      stepUser: 'Alice',
      keyword: 'NEW',
      filter: searchScope.allFiles
    })
    // Then following resources should be displayed in the search list for user "Alice"
    //   | resource             |
    //   | new_share_from_brian |
    //   | new-lorem-big.txt    |
    await ui.userShouldSeeResources({
      listType: resourcePage.searchList,
      stepUser: 'Alice',
      resources: ['new_share_from_brian', 'new-lorem-big.txt']
    })
    // But following resources should not be displayed in the search list for user "Alice"
    //   | resource         |
    //   | folder           |
    //   | FolDer           |
    //   | .hidden-file.txt |
    await ui.userShouldNotSeeTheResources({
      listType: resourcePage.searchList,
      stepUser: 'Alice',
      resources: ['folder', 'FolDer', '.hidden-file.txt']
    })
    // And "Alice" opens the "files" app
    await ui.userOpensApplication({ stepUser: 'Alice', name: 'files' })

    // # search renamed resources
    // When "Alice" renames the following resource
    //   | resource | as            |
    //   | folder   | renamedFolder |
    //   | FolDer   | renamedFolDer |
    await ui.userRenamesResource({
      stepUser: 'Alice',
      resource: 'folder',
      newResourceName: 'renamedFolder'
    })
    await ui.userRenamesResource({
      stepUser: 'Alice',
      resource: 'FolDer',
      newResourceName: 'renamedFolDer'
    })
    // And "Alice" searches "rena" using the global search and the "all files" filter
    await ui.userSearchesGloballyWithFilter({
      stepUser: 'Alice',
      keyword: 'rena',
      filter: searchScope.allFiles
    })
    // Then following resources should be displayed in the search list for user "Alice"
    //   | resource      |
    //   | renamedFolder |
    //   | renamedFolDer |
    await ui.userShouldSeeResources({
      listType: resourcePage.searchList,
      stepUser: 'Alice',
      resources: ['renamedFolder', 'renamedFolDer']
    })
    // But following resources should not be displayed in the search list for user "Alice"
    //   | resource |
    //   | folder   |
    //   | FolDer   |
    await ui.userShouldNotSeeTheResources({
      listType: resourcePage.searchList,
      stepUser: 'Alice',
      resources: ['folder', 'FolDer']
    })

    // # search difficult names
    // allow extra time for search indexing of Unicode folder name after rename operations
    const { page } = world.actorsEnvironment.getActor({ key: 'Alice' })
    await page.waitForTimeout(5000)
    // When "Alice" searches "strängéनेपालीName" using the global search and the "all files" filter and presses enter
    await ui.userSearchesGloballyWithFilter({
      stepUser: 'Alice',
      keyword: 'strängéनेपालीName',
      filter: searchScope.allFiles,
      command: 'presses enter'
    })
    // And "Alice" enables the option to search title only
    await ui.userEnablesTitleOnlySearch({ stepUser: 'Alice' })
    // Then following resources should be displayed in the files list for user "Alice"
    //   | strängéनेपालीName |
    await ui.userShouldSeeResources({
      listType: resourcePage.filesList,
      stepUser: 'Alice',
      resources: ['strängéनेपालीName']
    })

    // # deleting folder from search result and search deleted resource
    // When "Alice" deletes the following resource using the sidebar panel
    //   | resource          | from |
    //   | strängéनेपालीName |      |
    await ui.userDeletesResources({
      stepUser: 'Alice',
      actionType: fileAction.sideBarPanel,
      resources: [{ name: 'strängéनेपालीName' }]
    })

    // And "Alice" searches "forDeleting" using the global search and the "all files" filter
    await ui.userSearchesGloballyWithFilter({
      stepUser: 'Alice',
      keyword: 'forDeleting',
      filter: searchScope.allFiles
    })
    // Then following resources should not be displayed in the search list for user "Alice"
    //   | resource          |
    //   | strängéनेपालीName |
    await ui.userShouldNotSeeTheResources({
      listType: 'search list',
      stepUser: 'Alice',
      resources: ['strängéनेपालीName']
    })

    // And "Alice" logs out
    await ui.userLogsOut({ stepUser: 'Alice' })
  })

  test('Search using "current folder" filter', async () => {
    // And "Alice" logs in
    await ui.userLogsIn({ stepUser: 'Alice' })
    // And "Alice" creates the following folders in personal space using API
    //   | name                 |
    //   | mainFolder/subFolder |
    await api.userHasCreatedFolder({ stepUser: 'Alice', folderName: 'mainFolder/subFolder' })
    // And "Alice" creates the following files into personal space using API
    //   | pathToFile                                         | content                   |
    //   | exampleInsideThePersonalSpace.txt                  | I'm in the personal Space |
    //   | mainFolder/exampleInsideTheMainFolder.txt          | I'm in the main folder    |
    //   | mainFolder/subFolder/exampleInsideTheSubFolder.txt | I'm in the sub folder     |
    await api.userHasCreatedFiles({
      stepUser: 'Alice',
      files: [
        { pathToFile: 'exampleInsideThePersonalSpace.txt', content: "I'm in the personal Space" },
        {
          pathToFile: 'mainFolder/exampleInsideTheMainFolder.txt',
          content: "I'm in the main folder"
        },
        {
          pathToFile: 'mainFolder/subFolder/exampleInsideTheSubFolder.txt',
          content: "I'm in the sub folder"
        }
      ]
    })
    // When "Alice" opens folder "mainFolder"
    await ui.userOpensResource({ stepUser: 'Alice', resource: 'mainFolder' })
    // And "Alice" searches "example" using the global search and the "all files" filter
    await ui.userSearchesGloballyWithFilter({
      stepUser: 'Alice',
      keyword: 'example',
      filter: searchScope.allFiles
    })
    // Then following resources should be displayed in the search list for user "Alice"
    //   | resource                          |
    //   | exampleInsideThePersonalSpace.txt |
    //   | exampleInsideTheMainFolder.txt    |
    //   | exampleInsideTheSubFolder.txt     |
    await ui.userShouldSeeResources({
      listType: resourcePage.searchList,
      stepUser: 'Alice',
      resources: [
        'exampleInsideThePersonalSpace.txt',
        'exampleInsideTheMainFolder.txt',
        'exampleInsideTheSubFolder.txt'
      ]
    })

    // When "Alice" searches "example" using the global search and the "current folder" filter
    await ui.userSearchesGloballyWithFilter({
      stepUser: 'Alice',
      keyword: 'example',
      filter: searchScope.currentFolder
    })
    // Then following resources should be displayed in the search list for user "Alice"
    //   | resource                       |
    //   | exampleInsideTheMainFolder.txt |
    //   | exampleInsideTheSubFolder.txt  |
    await ui.userShouldSeeResources({
      listType: resourcePage.searchList,
      stepUser: 'Alice',
      resources: ['exampleInsideTheMainFolder.txt', 'exampleInsideTheSubFolder.txt']
    })
    // But following resources should not be displayed in the search list for user "Alice"
    //   | resource                          |
    //   | exampleInsideThePersonalSpace.txt |
    await ui.userShouldNotSeeTheResources({
      listType: 'search list',
      stepUser: 'Alice',
      resources: ['exampleInsideThePersonalSpace.txt']
    })
    // And "Alice" logs out
    await ui.userLogsOut({ stepUser: 'Alice' })
  })

  test('Search using mediaType filter', async () => {
    // And "Alice" logs in
    await ui.userLogsIn({ stepUser: 'Alice' })
    // And "Alice" creates the following folders in personal space using API
    //   | name      |
    //   | mediaTest |
    await api.userHasCreatedFolder({ stepUser: 'Alice', folderName: 'mediaTest' })
    // And "Alice" uploads the following local file into personal space using API
    //   | localFile                     | to            |
    //   | filesForUpload/testavatar.jpg | mediaTest.jpg |
    await api.userHasUploadedFilesInPersonalSpace({
      stepUser: 'Alice',
      filesToUpload: [{ localFile: 'filesForUpload/testavatar.jpg', to: 'mediaTest.jpg' }]
    })
    // And "Alice" creates the following files into personal space using API
    //   | pathToFile    | content        |
    //   | mediaTest.txt | I'm a Document |
    //   | mediaTest.pdf | I'm a PDF      |
    //   | mediaTest.mp3 | I'm a Audio    |
    //   | mediaTest.zip | I'm a Archive  |
    await api.userHasCreatedFiles({
      stepUser: 'Alice',
      files: [
        { pathToFile: 'mediaTest.txt', content: "I'm a Document" },
        { pathToFile: 'mediaTest.pdf', content: "I'm a PDF" },
        { pathToFile: 'mediaTest.mp3', content: "I'm a Audio" },
        { pathToFile: 'mediaTest.zip', content: "I'm a Archive" }
      ]
    })
    // And "Alice" searches "mediaTest" using the global search and the "all files" filter and presses enter
    await ui.userSearchesGloballyWithFilter({
      stepUser: 'Alice',
      keyword: 'mediaTest',
      filter: searchScope.allFiles,
      command: 'presses enter'
    })
    // And "Alice" selects mediaType "Document" from the search result filter chip
    await ui.userFiltersSearchByMediaType({ stepUser: 'Alice', mediaType: 'Document' })
    // And "Alice" enables the option to search title only
    await ui.userEnablesTitleOnlySearch({ stepUser: 'Alice' })
    // Then following resources should be displayed in the files list for user "Alice"
    //   | resource      |
    //   | mediaTest.txt |
    await ui.userShouldSeeResources({
      listType: resourcePage.filesList,
      stepUser: 'Alice',
      resources: ['mediaTest.txt']
    })
    // And "Alice" clears mediaType filter
    await ui.userClearsFilter({ stepUser: 'Alice', filter: searchFilter.mediaType })
    // When "Alice" selects mediaType "PDF" from the search result filter chip
    await ui.userFiltersSearchByMediaType({ stepUser: 'Alice', mediaType: 'PDF' })
    // Then following resources should be displayed in the files list for user "Alice"
    //   | resource      |
    //   | mediaTest.pdf |
    await ui.userShouldSeeResources({
      listType: resourcePage.filesList,
      stepUser: 'Alice',
      resources: ['mediaTest.pdf']
    })
    // And "Alice" clears mediaType filter
    await ui.userClearsFilter({ stepUser: 'Alice', filter: searchFilter.mediaType })
    // When "Alice" selects mediaType "Audio" from the search result filter chip
    await ui.userFiltersSearchByMediaType({ stepUser: 'Alice', mediaType: 'Audio' })
    // Then following resources should be displayed in the files list for user "Alice"
    //   | resource      |
    //   | mediaTest.mp3 |
    await ui.userShouldSeeResources({
      listType: resourcePage.filesList,
      stepUser: 'Alice',
      resources: ['mediaTest.mp3']
    })
    // And "Alice" clears mediaType filter
    await ui.userClearsFilter({ stepUser: 'Alice', filter: searchFilter.mediaType })
    // When "Alice" selects mediaType "Archive" from the search result filter chip
    await ui.userFiltersSearchByMediaType({ stepUser: 'Alice', mediaType: 'Archive' })
    // Then following resources should be displayed in the files list for user "Alice"
    //   | resource      |
    //   | mediaTest.zip |
    await ui.userShouldSeeResources({
      listType: resourcePage.filesList,
      stepUser: 'Alice',
      resources: ['mediaTest.zip']
    })
    // And "Alice" clears mediaType filter
    await ui.userClearsFilter({ stepUser: 'Alice', filter: searchFilter.mediaType })

    // # multiple choose
    // When "Alice" selects mediaType "Folder" from the search result filter chip
    await ui.userFiltersSearchByMediaType({ stepUser: 'Alice', mediaType: 'Folder' })
    // And "Alice" selects mediaType "Image" from the search result filter chip
    await ui.userFiltersSearchByMediaType({ stepUser: 'Alice', mediaType: 'Image' })
    // Then following resources should be displayed in the files list for user "Alice"
    //   | resource      |
    //   | mediaTest     |
    //   | mediaTest.jpg |
    await ui.userShouldSeeResources({
      listType: resourcePage.filesList,
      stepUser: 'Alice',
      resources: ['mediaTest', 'mediaTest.jpg']
    })
    // But following resources should not be displayed in the files list for user "Alice"
    //   | resource      |
    //   | mediaTest.txt |
    //   | mediaTest.pdf |
    //   | mediaTest.mp3 |
    //   | mediaTest.zip |
    await ui.userShouldNotSeeTheResources({
      listType: 'files list',
      stepUser: 'Alice',
      resources: ['mediaTest.txt', 'mediaTest.pdf', 'mediaTest.mp3', 'mediaTest.zip']
    })
    // And "Alice" logs out
    await ui.userLogsOut({ stepUser: 'Alice' })
  })

  test('Search using lastModified filter', async () => {
    // And "Alice" logs in
    await ui.userLogsIn({ stepUser: 'Alice' })
    // And "Alice" creates the following folders in personal space using API
    //   | name       |
    //   | mainFolder |
    await api.userHasCreatedFolder({ stepUser: 'Alice', folderName: 'mainFolder' })
    // And "Alice" creates the following files with mtime into personal space using API
    //   | pathToFile               | content             | mtimeDeltaDays |
    //   | mainFolder/mediaTest.pdf | created 29 days ago | -29 days       |
    //   | mainFolder/mediaTest.txt | created 5 days ago  | -5 days        |
    //   | mainFolder/mediaTest.md  | created today       |                |
    await api.userHasCreatedFiles({
      stepUser: 'Alice',
      files: [
        {
          pathToFile: 'mainFolder/mediaTest.pdf',
          content: 'created 29 days ago',
          mtimeDeltaDays: '-29 days'
        },
        {
          pathToFile: 'mainFolder/mediaTest.txt',
          content: 'created 5 days ago',
          mtimeDeltaDays: '-5 days'
        },
        { pathToFile: 'mainFolder/mediaTest.md', content: 'created today' }
      ]
    })
    // When "Alice" opens folder "mainFolder"
    await ui.userOpensResource({ stepUser: 'Alice', resource: 'mainFolder' })
    // And "Alice" searches "mediaTest" using the global search and the "current folder" filter and presses enter
    await ui.userSearchesGloballyWithFilter({
      stepUser: 'Alice',
      keyword: 'mediaTest',
      filter: searchScope.currentFolder,
      command: 'presses enter'
    })
    // And "Alice" selects lastModified "last 30 days" from the search result filter chip
    await ui.userFiltersSearchByLastModifiedDate({
      stepUser: 'Alice',
      lastModified: 'last 30 days'
    })
    // And "Alice" enables the option to search title only
    await ui.userEnablesTitleOnlySearch({ stepUser: 'Alice' })
    // Then following resources should be displayed in the files list for user "Alice"
    //   | resource                 |
    //   | mainFolder/mediaTest.pdf |
    //   | mainFolder/mediaTest.txt |
    //   | mainFolder/mediaTest.md  |
    await ui.userShouldSeeResources({
      listType: resourcePage.filesList,
      stepUser: 'Alice',
      resources: ['mainFolder/mediaTest.pdf', 'mainFolder/mediaTest.txt', 'mainFolder/mediaTest.md']
    })
    // When "Alice" selects lastModified "last 7 days" from the search result filter chip
    await ui.userFiltersSearchByLastModifiedDate({ stepUser: 'Alice', lastModified: 'last 7 days' })
    // Then following resources should be displayed in the files list for user "Alice"
    //   | resource                 |
    //   | mainFolder/mediaTest.txt |
    //   | mainFolder/mediaTest.md  |
    await ui.userShouldSeeResources({
      listType: resourcePage.filesList,
      stepUser: 'Alice',
      resources: ['mainFolder/mediaTest.txt', 'mainFolder/mediaTest.md']
    })
    // But following resources should not be displayed in the files list for user "Alice"
    //   | resource                 |
    //   | mainFolder/mediaTest.pdf |
    await ui.userShouldNotSeeTheResources({
      listType: resourcePage.searchList,
      stepUser: 'Alice',
      resources: ['mainFolder/mediaTest.pdf']
    })
    // When "Alice" selects lastModified "today" from the search result filter chip
    await ui.userFiltersSearchByLastModifiedDate({ stepUser: 'Alice', lastModified: 'today' })
    // Then following resources should be displayed in the files list for user "Alice"
    //   | resource                |
    //   | mainFolder/mediaTest.md |
    await ui.userShouldSeeResources({
      listType: resourcePage.filesList,
      stepUser: 'Alice',
      resources: ['mainFolder/mediaTest.md']
    })
    // But following resources should not be displayed in the files list for user "Alice"
    //   | resource                 |
    //   | mainFolder/mediaTest.pdf |
    //   | mainFolder/mediaTest.txt |
    await ui.userShouldNotSeeTheResources({
      listType: resourcePage.searchList,
      stepUser: 'Alice',
      resources: ['mainFolder/mediaTest.pdf', 'mainFolder/mediaTest.txt']
    })
    // And "Alice" clears lastModified filter
    await ui.userClearsFilter({ stepUser: 'Alice', filter: searchFilter.lastModified })
    // And "Alice" logs out
    await ui.userLogsOut({ stepUser: 'Alice' })
  })
})
