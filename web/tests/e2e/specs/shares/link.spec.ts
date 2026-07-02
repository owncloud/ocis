import { test } from '../../environment/test'
import * as api from '../../steps/api/api'
import * as ui from '../../steps/ui/index'
import { fileAction, resourcePage, application } from '../../environment/constants'

test.describe('link', () => {
  test.beforeEach(async () => {
    // Given "Admin" creates following users using API
    //   | id    |
    //   | Alice |
    await api.usersHaveBeenCreated({ stepUser: 'Admin', users: ['Alice'] })
  })

  test('public link', { tag: '@predefined-users' }, async () => {
    // Given "Admin" creates following users using API
    //   | id    |
    //   | Brian |
    await api.usersHaveBeenCreated({ stepUser: 'Admin', users: ['Brian'] })

    //  And "Alice" creates the following folders in personal space using API
    //   | name                   |
    //   | folderPublic           |
    //   | folderPublic/SubFolder |
    await api.userHasCreatedFolders({
      stepUser: 'Alice',
      folderNames: ['folderPublic', 'folderPublic/SubFolder']
    })
    // And "Alice" creates the following files into personal space using API
    //   | pathToFile             | content     |
    //   | folderPublic/lorem.txt | lorem ipsum |
    await api.userHasCreatedFiles({
      stepUser: 'Alice',
      files: [{ pathToFile: 'folderPublic/lorem.txt', content: 'lorem ipsum' }]
    })

    // When "Alice" logs in
    await ui.userLogsIn({ stepUser: 'Alice' })
    // And "Alice" creates a public link of following resource using the sidebar panel
    //   | resource     | role             | password |
    //   | folderPublic | Secret File Drop | %public% |
    await ui.userCreatesPublicLink({
      stepUser: 'Alice',
      resource: 'folderPublic',
      role: 'Secret File Drop',
      password: '%public%'
    })
    // And "Alice" renames the most recently created public link of resource "folderPublic" to "myPublicLink"
    await ui.userRenamesMostRecentlyCreatedPublicLinkOfResource({
      stepUser: 'Alice',
      resource: 'folderPublic',
      newName: 'myPublicLink'
    })
    // And "Alice" edits the public link named "myPublicLink" of resource "folderPublic" changing role to "Secret File Drop"
    await ui.userChangesRoleOfPublicLinkOfResource({
      stepUser: 'Alice',
      resource: 'folderPublic',
      linkName: 'myPublicLink',
      newRole: 'Secret File Drop'
    })
    // And "Alice" sets the expiration date of the public link named "myPublicLink" of resource "folderPublic" to "+5 days"
    await ui.userSetsExpirationDateOfThePublicLinkOfResource({
      stepUser: 'Alice',
      resource: 'folderPublic',
      linkName: 'myPublicLink',
      expireDate: '+5 days'
    })
    // When "Anonymous" opens the public link "myPublicLink"
    await ui.userOpensPublicLink({ stepUser: 'Anonymous', name: 'myPublicLink' })
    // And "Anonymous" unlocks the public link with password "%public%"
    await ui.userUnlocksPublicLink({ password: '%public%', stepUser: 'Anonymous' })
    // And "Anonymous" drop uploads following resources
    //   | resource     |
    //   | textfile.txt |
    await ui.userDropUploadsResources({ stepUser: 'Anonymous', resources: ['textfile.txt'] })

    // authenticated user
    // When "Brian" logs in
    await ui.userLogsIn({ stepUser: 'Brian' })
    // And "Brian" opens the public link "myPublicLink"
    await ui.userOpensPublicLink({ stepUser: 'Brian', name: 'myPublicLink' })
    // And "Brian" unlocks the public link with password "%public%"
    await ui.userUnlocksPublicLink({ password: '%public%', stepUser: 'Brian' })
    // And "Brian" drop uploads following resources
    //   | resource   |
    //   | simple.pdf |
    await ui.userDropUploadsResources({ stepUser: 'Brian', resources: ['simple.pdf'] })

    //   When "Alice" opens folder "folderPublic"
    await ui.userOpensResource({ stepUser: 'Alice', resource: 'folderPublic' })
    // Then following resources should be displayed in the files list for user "Alice"
    //   | resource     |
    //   | textfile.txt |
    //   | simple.pdf   |
    await ui.userShouldSeeResources({
      listType: resourcePage.filesList,
      stepUser: 'Alice',
      resources: ['textfile.txt', 'simple.pdf']
    })
    // And "Alice" opens the "files" app
    await ui.userOpensApplication({ stepUser: 'Alice', name: 'files' })
    // And "Alice" edits the public link named "myPublicLink" of resource "folderPublic" changing role to "Can edit"
    await ui.userChangesRoleOfPublicLinkOfResource({
      stepUser: 'Alice',
      resource: 'folderPublic',
      linkName: 'myPublicLink',
      newRole: 'Can edit'
    })
    // And "Brian" refreshes the old link
    await ui.userRefreshesTheOldLink({ stepUser: 'Brian' })

    // Then following resources should be displayed in the files list for user "Brian"
    //   | resource     |
    //   | textfile.txt |
    //   | simple.pdf   |
    //   | SubFolder    |
    //   | lorem.txt    |
    await ui.userShouldSeeResources({
      listType: resourcePage.filesList,
      stepUser: 'Brian',
      resources: ['textfile.txt', 'simple.pdf', 'SubFolder', 'lorem.txt']
    })
    // And "Brian" deletes the following resources from public link using sidebar panel
    //   | resource   |
    //   | simple.pdf |
    await ui.userDeletesResourcesFromPublicLink({
      stepUser: 'Brian',
      actionType: fileAction.sideBarPanel,
      resources: [{ resource: 'simple.pdf' }]
    })
    // And "Brian" logs out
    await ui.userLogsOut({ stepUser: 'Brian' })

    // And "Anonymous" refreshes the old link
    await ui.userRefreshesTheOldLink({ stepUser: 'Anonymous' })
    // And "Anonymous" downloads the following public link resources using the sidebar panel
    //   | resource     | type |
    //   | lorem.txt    | file |
    //   | textfile.txt | file |
    await ui.userDownloadsThePublicLinkResources({
      stepUser: 'Anonymous',
      actionType: fileAction.sideBarPanel,
      resources: [
        { resource: 'lorem.txt', type: 'file' },
        { resource: 'textfile.txt', type: 'file' }
      ]
    })
    // And "Anonymous" uploads the following resources in public link page
    //   | resource      | option  |
    //   | new-lorem.txt |         |
    //   | lorem.txt     | replace |
    await ui.userUploadsResourcesInPublicLink({
      stepUser: 'Anonymous',
      resources: [{ name: 'new-lorem.txt' }, { name: 'lorem.txt', option: 'replace' }]
    })
    // And "Anonymous" creates the following resources
    //   | resource       | type   |
    //   | myfolder/child | folder |
    await ui.userCreatesResources({
      stepUser: 'Anonymous',
      resources: [{ name: 'myfolder/child', type: 'folder' }]
    })

    // And "Anonymous" uploads the following resources in public link page
    //   | resource | type   |
    //   | PARENT   | folder |
    await ui.userUploadsResourcesInPublicLink({
      stepUser: 'Anonymous',
      resources: [{ name: 'PARENT', type: 'folder' }]
    })
    // And "Anonymous" should see the resource "PARENT" in the files list
    await ui.userShouldSeeResources({
      listType: resourcePage.filesList,
      stepUser: 'Anonymous',
      resources: ['PARENT']
    })
    // And "Anonymous" moves the following resource using drag-drop
    //   | resource      | to        |
    //   | new-lorem.txt | SubFolder |
    await ui.userMovesResources({
      stepUser: 'Anonymous',
      actionType: fileAction.dragDrop,
      resources: [{ resource: 'new-lorem.txt', to: 'SubFolder' }]
    })
    // And "Anonymous" copies the following resource using sidebar-panel
    //   | resource  | to       |
    //   | lorem.txt | myfolder |
    await ui.userCopiesResources({
      stepUser: 'Anonymous',
      actionType: fileAction.sideBarPanel,
      resources: [{ resource: 'lorem.txt', to: 'myfolder' }]
    })
    // And "Anonymous" renames the following public link resources
    //   | resource     | as               |
    //   | lorem.txt    | lorem_new.txt    |
    //   | textfile.txt | textfile_new.txt |
    await ui.userRenamesPublicLinkResources({
      stepUser: 'Anonymous',
      resources: [
        { resource: 'lorem.txt', newName: 'lorem_new.txt' },
        { resource: 'textfile.txt', newName: 'textfile_new.txt' }
      ]
    })
    // And "Anonymous" deletes the following resources from public link using batch action
    //   | resource  | from     |
    //   | lorem.txt | myfolder |
    //   | child     | myfolder |
    await ui.userDeletesResourcesFromPublicLink({
      stepUser: 'Anonymous',
      actionType: fileAction.batchAction,
      resources: [
        { resource: 'lorem.txt', parentFolder: 'myfolder' },
        { resource: 'child', parentFolder: 'myfolder' }
      ]
    })
    // And "Alice" removes the public link named "myPublicLink" of resource "folderPublic"
    await ui.userRemovesThePublicLinkOfResource({
      stepUser: 'Alice',
      resource: 'folderPublic',
      linkName: 'myPublicLink'
    })
    // And "Anonymous" should not be able to open the old link "myPublicLink"
    await ui.userShouldNotBeAbleToOpenTheOldLink({
      stepUser: 'Anonymous',
      linkName: 'myPublicLink'
    })
    // And "Alice" logs out
    await ui.userLogsOut({ stepUser: 'Alice' })
  })

  test(
    'public link for folder and file (by authenticated user)',
    { tag: '@predefined-users' },
    async () => {
      // Given "Admin" creates following user using API
      //   | id    |
      //   | Brian |
      //   | Carol |
      await api.usersHaveBeenCreated({ stepUser: 'Admin', users: ['Brian', 'Carol'] })

      // And "Alice" logs in
      await ui.userLogsIn({ stepUser: 'Alice' })
      // And "Alice" creates the following folders in personal space using API
      //   | name         |
      //   | folderPublic |
      await api.userHasCreatedFolders({ stepUser: 'Alice', folderNames: ['folderPublic'] })
      // And "Alice" creates the following files into personal space using API
      //   | pathToFile                    | content   |
      //   | folderPublic/shareToBrian.txt | some text |
      //   | folderPublic/shareToBrian.md  | readme    |
      await api.userHasCreatedFiles({
        stepUser: 'Alice',
        files: [
          { pathToFile: 'folderPublic/shareToBrian.txt', content: 'some text' },
          { pathToFile: 'folderPublic/shareToBrian.md', content: 'readme' }
        ]
      })
      // And "Alice" uploads the following local file into personal space using API
      //   | localFile                     | to             |
      //   | filesForUpload/simple.pdf     | simple.pdf     |
      //   | filesForUpload/testavatar.jpg | testavatar.jpg |
      //   | filesForUpload/test_video.mp4 | test_video.mp4 |
      await api.userHasUploadedFilesInPersonalSpace({
        stepUser: 'Alice',
        filesToUpload: [
          { localFile: 'filesForUpload/simple.pdf', to: 'simple.pdf' },
          { localFile: 'filesForUpload/testavatar.jpg', to: 'testavatar.jpg' },
          { localFile: 'filesForUpload/test_video.mp4', to: 'test_video.mp4' }
        ]
      })

      // And "Alice" shares the following resource using API
      //   | resource       | recipient | type | role     |
      //   | folderPublic   | Brian     | user | Can edit |
      //   | simple.pdf     | Brian     | user | Can edit |
      //   | testavatar.jpg | Brian     | user | Can edit |
      await api.userHasSharedResources({
        stepUser: 'Alice',
        shares: [
          {
            resource: 'folderPublic',
            recipient: 'Brian',
            type: 'user',
            role: 'Can edit',
            resourceType: 'folder'
          },
          {
            resource: 'simple.pdf',
            recipient: 'Brian',
            type: 'user',
            role: 'Can edit',
            resourceType: 'file'
          },
          {
            resource: 'testavatar.jpg',
            recipient: 'Brian',
            type: 'user',
            role: 'Can edit',
            resourceType: 'file'
          }
        ]
      })

      // And "Alice" opens the "files" app
      await ui.userOpensApplication({ stepUser: 'Alice', name: 'files' })

      // And "Alice" creates a public link of following resource using the sidebar panel
      //   | resource     | password |
      //   | folderPublic | %public% |
      await ui.userCreatesPublicLink({
        stepUser: 'Alice',
        resource: 'folderPublic',
        password: '%public%'
      })
      // And "Alice" renames the most recently created public link of resource "folderPublic" to "folderLink"
      await ui.userRenamesMostRecentlyCreatedPublicLinkOfResource({
        stepUser: 'Alice',
        resource: 'folderPublic',
        newName: 'folderLink'
      })
      // And "Alice" creates a public link of following resource using the sidebar panel
      //   | resource                      | password |
      //   | folderPublic/shareToBrian.txt | %public% |
      await ui.userCreatesPublicLink({
        stepUser: 'Alice',
        resource: 'folderPublic/shareToBrian.txt',
        password: '%public%'
      })
      // And "Alice" renames the most recently created public link of resource "folderPublic/shareToBrian.txt" to "textLink"
      await ui.userRenamesMostRecentlyCreatedPublicLinkOfResource({
        stepUser: 'Alice',
        resource: 'folderPublic/shareToBrian.txt',
        newName: 'textLink'
      })
      // And "Alice" creates a public link of following resource using the sidebar panel
      //   | resource                     | password |
      //   | folderPublic/shareToBrian.md | %public% |
      await ui.userCreatesPublicLink({
        stepUser: 'Alice',
        resource: 'folderPublic/shareToBrian.md',
        password: '%public%'
      })
      // And "Alice" renames the most recently created public link of resource "folderPublic/shareToBrian.md" to "markdownLink"
      await ui.userRenamesMostRecentlyCreatedPublicLinkOfResource({
        stepUser: 'Alice',
        resource: 'folderPublic/shareToBrian.md',
        newName: 'markdownLink'
      })
      // And "Alice" creates a public link of following resource using the sidebar panel
      //   | resource   | password |
      //   | simple.pdf | %public% |
      await ui.userCreatesPublicLink({
        stepUser: 'Alice',
        resource: 'simple.pdf',
        password: '%public%'
      })
      // And "Alice" renames the most recently created public link of resource "simple.pdf" to "pdfLink"
      await ui.userRenamesMostRecentlyCreatedPublicLinkOfResource({
        stepUser: 'Alice',
        resource: 'simple.pdf',
        newName: 'pdfLink'
      })
      // And "Alice" creates a public link of following resource using the sidebar panel
      //   | resource       | password |
      //   | testavatar.jpg | %public% |
      await ui.userCreatesPublicLink({
        stepUser: 'Alice',
        resource: 'testavatar.jpg',
        password: '%public%'
      })
      // And "Alice" renames the most recently created public link of resource "testavatar.jpg" to "imageLink"
      await ui.userRenamesMostRecentlyCreatedPublicLinkOfResource({
        stepUser: 'Alice',
        resource: 'testavatar.jpg',
        newName: 'imageLink'
      })
      // And "Alice" creates a public link of following resource using the sidebar panel
      //   | resource       | password |
      //   | test_video.mp4 | %public% |
      await ui.userCreatesPublicLink({
        stepUser: 'Alice',
        resource: 'test_video.mp4',
        password: '%public%'
      })
      // And "Alice" renames the most recently created public link of resource "test_video.mp4" to "videoLink"
      await ui.userRenamesMostRecentlyCreatedPublicLinkOfResource({
        stepUser: 'Alice',
        resource: 'test_video.mp4',
        newName: 'videoLink'
      })
      // And "Alice" logs out
      await ui.userLogsOut({ stepUser: 'Alice' })

      // authenticated user with access to resources. should be redirected to shares with me page
      // And "Brian" logs in
      await ui.userLogsIn({ stepUser: 'Brian' })
      // When "Brian" opens the public link "folderLink"
      await ui.userOpensPublicLink({ stepUser: 'Brian', name: 'folderLink' })
      // And "Brian" unlocks the public link with password "%public%"
      await ui.userUnlocksPublicLink({ password: '%public%', stepUser: 'Brian' })
      // And "Brian" downloads the following public link resources using the sidebar panel
      //   | resource         | type |
      //   | shareToBrian.txt | file |
      await ui.userDownloadsThePublicLinkResources({
        stepUser: 'Brian',
        actionType: 'sidebar panel',
        resources: [{ resource: 'shareToBrian.txt', type: 'file' }]
      })
      // And "Brian" uploads the following resources
      //   | resource  |
      //   | lorem.txt |
      await ui.userUploadsResourcesInPublicLink({
        stepUser: 'Brian',
        resources: [{ name: 'lorem.txt' }]
      })
      // When "Brian" opens the public link "textLink"
      await ui.userOpensPublicLink({ stepUser: 'Brian', name: 'textLink' })
      // And "Brian" unlocks the public link with password "%public%"
      await ui.userUnlocksPublicLink({ stepUser: 'Brian', password: '%public%' })
      // Then "Brian" is in a text-editor
      await ui.userShouldBeInFileViewer({
        stepUser: 'Brian',
        fileViewerType: application.textEditor
      })

      // And "Brian" closes the file viewer
      await ui.userClosesFileViewer({ stepUser: 'Brian' })
      // When "Brian" opens the public link "markdownLink"
      await ui.userOpensPublicLink({ stepUser: 'Brian', name: 'markdownLink' })
      // And "Brian" unlocks the public link with password "%public%"
      await ui.userUnlocksPublicLink({ password: '%public%', stepUser: 'Brian' })
      // Then "Brian" is in a text-editor
      await ui.userShouldBeInFileViewer({
        stepUser: 'Brian',
        fileViewerType: application.textEditor
      })
      // And "Brian" closes the file viewer
      await ui.userClosesFileViewer({ stepUser: 'Brian' })
      // And "Brian" downloads the following public link resources using the single share view
      //   | resource        | type |
      //   | shareToBrian.md | file |
      await ui.userDownloadsThePublicLinkResources({
        stepUser: 'Brian',
        actionType: fileAction.singleShareView,
        resources: [{ resource: 'shareToBrian.md', type: 'file' }]
      })
      // When "Brian" opens the public link "pdfLink"
      await ui.userOpensPublicLink({ stepUser: 'Brian', name: 'pdfLink' })
      // And "Brian" unlocks the public link with password "%public%"
      await ui.userUnlocksPublicLink({ password: '%public%', stepUser: 'Brian' })
      // Then "Brian" is in a pdf-viewer
      await ui.userShouldBeInFileViewer({
        stepUser: 'Brian',
        fileViewerType: application.pdfViewer
      })
      // And "Brian" closes the file viewer
      await ui.userClosesFileViewer({ stepUser: 'Brian' })
      // And "Brian" downloads the following public link resources using the single share view
      //   | resource   | type |
      //   | simple.pdf | file |
      await ui.userDownloadsThePublicLinkResources({
        stepUser: 'Brian',
        actionType: 'SINGLE_SHARE_VIEW',
        resources: [{ resource: 'simple.pdf', type: 'file' }]
      })
      // When "Brian" opens the public link "imageLink"
      await ui.userOpensPublicLink({ stepUser: 'Brian', name: 'imageLink' })
      // And "Brian" unlocks the public link with password "%public%"
      await ui.userUnlocksPublicLink({ password: '%public%', stepUser: 'Brian' })
      // https://github.com/owncloud/ocis/issues/8602
      // Then "Brian" is in a media-viewer
      await ui.userShouldBeInFileViewer({
        stepUser: 'Brian',
        fileViewerType: application.mediaViewer
      })
      // And "Brian" closes the file viewer
      await ui.userClosesFileViewer({ stepUser: 'Brian' })
      // And "Brian" downloads the following public link resources using the single share view
      //   | resource       | type |
      //   | testavatar.jpg | file |
      await ui.userDownloadsThePublicLinkResources({
        stepUser: 'Brian',
        actionType: 'SINGLE_SHARE_VIEW',
        resources: [{ resource: 'testavatar.jpg', type: 'file' }]
      })
      // And "Brian" logs out
      await ui.userLogsOut({ stepUser: 'Brian' })

      // authenticated user without access to resources. should be redirected to the public links page
      // And "Carol" logs in
      await ui.userLogsIn({ stepUser: 'Carol' })
      // When "Carol" opens the public link "folderLink"
      await ui.userOpensPublicLink({ stepUser: 'Carol', name: 'folderLink' })
      // And "Carol" unlocks the public link with password "%public%"
      await ui.userUnlocksPublicLink({ password: '%public%', stepUser: 'Carol' })
      // https://github.com/owncloud/web/issues/10473
      // And "Carol" downloads the following public link resources using the sidebar panel
      //   | resource  | type |
      //   | lorem.txt | file |
      await ui.userDownloadsThePublicLinkResources({
        stepUser: 'Carol',
        actionType: 'sidebar panel',
        resources: [{ resource: 'lorem.txt', type: 'file' }]
      })
      // When "Carol" opens the public link "textLink"
      await ui.userOpensPublicLink({ stepUser: 'Carol', name: 'textLink' })
      // And "Carol" unlocks the public link with password "%public%"
      await ui.userUnlocksPublicLink({ password: '%public%', stepUser: 'Carol' })
      // Then "Carol" is in a text-editor
      await ui.userShouldBeInFileViewer({
        stepUser: 'Carol',
        fileViewerType: application.textEditor
      })
      // And "Carol" closes the file viewer
      await ui.userClosesFileViewer({ stepUser: 'Carol' })
      // When "Carol" opens the public link "markdownLink"
      await ui.userOpensPublicLink({ stepUser: 'Carol', name: 'markdownLink' })
      // And "Carol" unlocks the public link with password "%public%"
      await ui.userUnlocksPublicLink({ password: '%public%', stepUser: 'Carol' })
      // Then "Carol" is in a text-editor
      await ui.userShouldBeInFileViewer({
        stepUser: 'Carol',
        fileViewerType: application.textEditor
      })
      // And "Carol" closes the file viewer
      await ui.userClosesFileViewer({ stepUser: 'Carol' })
      // And "Carol" downloads the following public link resources using the single share view
      //   | resource        | type |
      //   | shareToBrian.md | file |
      await ui.userDownloadsThePublicLinkResources({
        stepUser: 'Carol',
        actionType: 'SINGLE_SHARE_VIEW',
        resources: [{ resource: 'shareToBrian.md', type: 'file' }]
      })
      // When "Carol" opens the public link "pdfLink"
      await ui.userOpensPublicLink({ stepUser: 'Carol', name: 'pdfLink' })
      // And "Carol" unlocks the public link with password "%public%"
      await ui.userUnlocksPublicLink({ password: '%public%', stepUser: 'Carol' })
      // Then "Carol" is in a pdf-viewer
      await ui.userShouldBeInFileViewer({
        stepUser: 'Carol',
        fileViewerType: application.pdfViewer
      })
      // And "Carol" closes the file viewer
      await ui.userClosesFileViewer({ stepUser: 'Carol' })
      // And "Carol" downloads the following public link resources using the single share view
      //   | resource   | type |
      //   | simple.pdf | file |
      await ui.userDownloadsThePublicLinkResources({
        stepUser: 'Carol',
        actionType: 'SINGLE_SHARE_VIEW',
        resources: [{ resource: 'simple.pdf', type: 'file' }]
      })
      // When "Carol" opens the public link "imageLink"
      await ui.userOpensPublicLink({ stepUser: 'Carol', name: 'imageLink' })
      // And "Carol" unlocks the public link with password "%public%"
      await ui.userUnlocksPublicLink({ password: '%public%', stepUser: 'Carol' })
      // https://github.com/owncloud/ocis/issues/8602
      // Then "Carol" is in a media-viewer
      await ui.userShouldBeInFileViewer({
        stepUser: 'Carol',
        fileViewerType: application.mediaViewer
      })
      // And "Carol" closes the file viewer
      await ui.userClosesFileViewer({ stepUser: 'Carol' })
      // And "Carol" downloads the following public link resources using the single share view
      //   | resource       | type |
      //   | testavatar.jpg | file |
      await ui.userDownloadsThePublicLinkResources({
        stepUser: 'Carol',
        actionType: 'SINGLE_SHARE_VIEW',
        resources: [{ resource: 'testavatar.jpg', type: 'file' }]
      })
      // And "Carol" logs out
      await ui.userLogsOut({ stepUser: 'Carol' })

      // Anonymous user
      // When "Anonymous" opens the public link "folderLink"
      await ui.userOpensPublicLink({ stepUser: 'Anonymous', name: 'folderLink' })
      // And "Anonymous" unlocks the public link with password "%public%"
      await ui.userUnlocksPublicLink({ password: '%public%', stepUser: 'Anonymous' })
      // And "Anonymous" downloads the following public link resources using the sidebar panel
      //   | resource  | type |
      //   | lorem.txt | file |
      await ui.userDownloadsThePublicLinkResources({
        stepUser: 'Anonymous',
        actionType: 'sidebar panel',
        resources: [{ resource: 'lorem.txt', type: 'file' }]
      })
      // When "Anonymous" opens the public link "textLink"
      await ui.userOpensPublicLink({ stepUser: 'Anonymous', name: 'textLink' })
      // And "Anonymous" unlocks the public link with password "%public%"
      await ui.userUnlocksPublicLink({ password: '%public%', stepUser: 'Anonymous' })
      // Then "Anonymous" is in a text-editor
      await ui.userShouldBeInFileViewer({
        stepUser: 'Anonymous',
        fileViewerType: application.textEditor
      })
      // And "Anonymous" closes the file viewer
      await ui.userClosesFileViewer({ stepUser: 'Anonymous' })
      // When "Anonymous" opens the public link "markdownLink"
      await ui.userOpensPublicLink({ stepUser: 'Anonymous', name: 'markdownLink' })
      // And "Anonymous" unlocks the public link with password "%public%"
      await ui.userUnlocksPublicLink({ password: '%public%', stepUser: 'Anonymous' })
      // Then "Anonymous" is in a text-editor
      await ui.userShouldBeInFileViewer({
        stepUser: 'Anonymous',
        fileViewerType: application.textEditor
      })
      // And "Anonymous" closes the file viewer
      await ui.userClosesFileViewer({ stepUser: 'Anonymous' })
      // And "Anonymous" downloads the following public link resources using the single share view
      //   | resource        | type |
      //   | shareToBrian.md | file |
      await ui.userDownloadsThePublicLinkResources({
        stepUser: 'Anonymous',
        actionType: 'SINGLE_SHARE_VIEW',
        resources: [{ resource: 'shareToBrian.md', type: 'file' }]
      })
      // When "Anonymous" opens the public link "pdfLink"
      await ui.userOpensPublicLink({ stepUser: 'Anonymous', name: 'pdfLink' })
      // And "Anonymous" unlocks the public link with password "%public%"
      await ui.userUnlocksPublicLink({ password: '%public%', stepUser: 'Anonymous' })
      // Then "Anonymous" is in a pdf-viewer
      await ui.userShouldBeInFileViewer({
        stepUser: 'Anonymous',
        fileViewerType: application.pdfViewer
      })
      // And "Anonymous" closes the file viewer
      await ui.userClosesFileViewer({ stepUser: 'Anonymous' })
      // And "Anonymous" downloads the following public link resources using the single share view
      //   | resource   | type |
      //   | simple.pdf | file |
      await ui.userDownloadsThePublicLinkResources({
        stepUser: 'Anonymous',
        actionType: 'SINGLE_SHARE_VIEW',
        resources: [{ resource: 'simple.pdf', type: 'file' }]
      })
      // When "Anonymous" opens the public link "imageLink"
      await ui.userOpensPublicLink({ stepUser: 'Anonymous', name: 'imageLink' })
      // And "Anonymous" unlocks the public link with password "%public%
      await ui.userUnlocksPublicLink({ password: '%public%', stepUser: 'Anonymous' })
      // https://github.com/owncloud/ocis/issues/8602
      // Then "Anonymous" is in a media-viewer
      await ui.userShouldBeInFileViewer({
        stepUser: 'Anonymous',
        fileViewerType: application.mediaViewer
      })
      // And "Anonymous" closes the file viewer
      await ui.userClosesFileViewer({ stepUser: 'Anonymous' })
      // And "Anonymous" downloads the following public link resources using the single share view
      //   | resource       | type |
      //   | testavatar.jpg | file |
      await ui.userDownloadsThePublicLinkResources({
        stepUser: 'Anonymous',
        actionType: 'SINGLE_SHARE_VIEW',
        resources: [{ resource: 'testavatar.jpg', type: 'file' }]
      })
      // When "Anonymous" opens the public link "videoLink"
      await ui.userOpensPublicLink({ stepUser: 'Anonymous', name: 'videoLink' })
      // And "Anonymous" unlocks the public link with password "%public%"
      await ui.userUnlocksPublicLink({ password: '%public%', stepUser: 'Anonymous' })
      // Then "Anonymous" is in a media-viewer
      await ui.userShouldBeInFileViewer({
        stepUser: 'Anonymous',
        fileViewerType: application.mediaViewer
      })
      // And "Anonymous" closes the file viewer
      await ui.userClosesFileViewer({ stepUser: 'Anonymous' })
    }
  )

  test('add banned password for public link', async () => {
    // When "Alice" logs in
    await ui.userLogsIn({ stepUser: 'Alice' })
    // And "Alice" creates the following files into personal space using API
    //   | pathToFile | content   |
    //   | lorem.txt  | some text |
    await api.userHasCreatedFiles({
      stepUser: 'Alice',
      files: [{ pathToFile: 'lorem.txt', content: 'some text' }]
    })

    // And "Alice" opens the "files" app
    await ui.userOpensApplication({ stepUser: 'Alice', name: 'files' })
    // And "Alice" creates a public link of following resource using the sidebar panel
    //   | resource  | password |
    //   | lorem.txt | %public% |
    await ui.userCreatesPublicLink({
      stepUser: 'Alice',
      resource: 'lorem.txt',
      password: '%public%'
    })
    // When "Alice" tries to sets a new password "ownCloud-1" of the public link named "Unnamed link" of resource "lorem.txt"
    await ui.userChangesPasswordOfThePublicLinkOfResource({
      stepUser: 'Alice',
      resource: 'lorem.txt',
      linkName: 'Unnamed link',
      newPassword: 'ownCloud-1'
    })
    // Then "Alice" should see an error message
    //   """
    //   Unfortunately, your password is commonly used. please pick a harder-to-guess password for your safety
    //   """
    await ui.userShouldSeeAnErrorMessage({
      stepUser: 'Alice',
      errorMessage:
        'Unfortunately, your password is commonly used. please pick a harder-to-guess password for your safety'
    })
    // And "Alice" closes the public link password dialog box
    await ui.userClosesThePublicLinkPasswordDialogBox({ stepUser: 'Alice' })
    // When "Alice" tries to sets a new password "12345678" of the public link named "Unnamed link" of resource "lorem.txt"
    await ui.userChangesPasswordOfThePublicLinkOfResource({
      stepUser: 'Alice',
      resource: 'lorem.txt',
      linkName: 'Unnamed link',
      newPassword: '12345678'
    })
    // Then "Alice" should see an error message
    //   """
    //   Unfortunately, your password is commonly used. please pick a harder-to-guess password for your safety
    //   """
    await ui.userShouldSeeAnErrorMessage({
      stepUser: 'Alice',
      errorMessage:
        'Unfortunately, your password is commonly used. please pick a harder-to-guess password for your safety'
    })
    // And "Alice" reveals the password of the public link
    await ui.userRevealsThePasswordOfThePublicLink({ stepUser: 'Alice' })
    // And "Alice" hides the password of the public link
    await ui.userHidesThePasswordOfThePublicLink({ stepUser: 'Alice' })
    // And "Alice" generates the password for the public link
    await ui.userGeneratesThePasswordForThePublicLink({ stepUser: 'Alice' })
    // And "Alice" copies the password of the public link
    await ui.userCopiesThePasswordOfThePublicLink({ stepUser: 'Alice' })
    // And "Alice" sets the password of the public link
    await ui.userSetsThePasswordOfThePublicLink({ stepUser: 'Alice' })
    // And "Anonymous" opens the public link "Unnamed link"
    await ui.userOpensPublicLink({ stepUser: 'Anonymous', name: 'Unnamed link' })
    // And "Anonymous" unlocks the public link with password "%copied_password%"
    await ui.userUnlocksPublicLink({ password: '%copied_password%', stepUser: 'Anonymous' })
    // And "Alice" logs out
    await ui.userLogsOut({ stepUser: 'Alice' })
  })

  test('edit password of the public link', { tag: '@predefined-users' }, async () => {
    // When "Alice" logs in
    await ui.userLogsIn({ stepUser: 'Alice' })
    // And "Alice" creates the following folders in personal space using API
    //   | name         |
    //   | folderPublic |
    await api.userHasCreatedFolders({ stepUser: 'Alice', folderNames: ['folderPublic'] })
    // And "Alice" creates the following files into personal space using API
    //   | pathToFile             | content     |
    //   | folderPublic/lorem.txt | lorem ipsum |
    await api.userHasCreatedFiles({
      stepUser: 'Alice',
      files: [{ pathToFile: 'folderPublic/lorem.txt', content: 'lorem ipsum' }]
    })
    // And "Alice" opens the "files" app
    await ui.userOpensApplication({ stepUser: 'Alice', name: 'files' })
    // And "Alice" creates a public link of following resource using the sidebar panel
    //   | resource     | role     | password |
    //   | folderPublic | Can edit | %public% |
    await ui.userCreatesPublicLink({
      stepUser: 'Alice',
      resource: 'folderPublic',
      password: '%public%',
      role: 'Can edit'
    })
    // And "Alice" renames the most recently created public link of resource "folderPublic" to "myPublicLink"
    await ui.userRenamesMostRecentlyCreatedPublicLinkOfResource({
      stepUser: 'Alice',
      resource: 'folderPublic',
      newName: 'myPublicLink'
    })
    // When "Anonymous" opens the public link "myPublicLink"
    await ui.userOpensPublicLink({ stepUser: 'Anonymous', name: 'myPublicLink' })
    // And "Anonymous" unlocks the public link with password "%public%"
    await ui.userUnlocksPublicLink({ password: '%public%', stepUser: 'Anonymous' })
    // And "Alice" changes the password of the public link named "myPublicLink" of resource "folderPublic" to "new-strongPass1"
    await ui.userChangesPasswordOfThePublicLinkOfResource({
      stepUser: 'Alice',
      resource: 'folderPublic',
      linkName: 'myPublicLink',
      newPassword: 'new-strongPass1'
    })
    // And "Anonymous" refreshes the old link
    await ui.userRefreshesTheOldLink({ stepUser: 'Anonymous' })
    // And "Anonymous" unlocks the public link with password "new-strongPass1"
    await ui.userUnlocksPublicLink({ password: 'new-strongPass1', stepUser: 'Anonymous' })
    // And "Anonymous" downloads the following public link resources using the sidebar panel
    //   | resource  | type |
    //   | lorem.txt | file |
    await ui.userDownloadsThePublicLinkResources({
      stepUser: 'Anonymous',
      actionType: 'sidebar panel',
      resources: [{ resource: 'lorem.txt', type: 'file' }]
    })
    // And "Alice" logs out
    await ui.userLogsOut({ stepUser: 'Alice' })
  })

  test('link indication', { tag: '@predefined-users' }, async () => {
    // When "Alice" logs in
    await ui.userLogsIn({ stepUser: 'Alice' })
    // And "Alice" creates the following folders in personal space using API
    //   | name         |
    //   | folderPublic |
    await api.userHasCreatedFolders({ stepUser: 'Alice', folderNames: ['folderPublic'] })
    // And "Alice" creates the following files into personal space using API
    //   | pathToFile             | content     |
    //   | folderPublic/lorem.txt | lorem ipsum |
    await api.userHasCreatedFiles({
      stepUser: 'Alice',
      files: [{ pathToFile: 'folderPublic/lorem.txt', content: 'lorem ipsum' }]
    })
    // And "Alice" opens the "files" app
    await ui.userOpensApplication({ stepUser: 'Alice', name: 'files' })
    // And "Alice" creates a public link of following resource using the sidebar panel
    //   | resource     | role     | password |
    //   | folderPublic | Can edit | %public% |
    await ui.userCreatesPublicLink({
      stepUser: 'Alice',
      resource: 'folderPublic',
      password: '%public%',
      role: 'Can edit'
    })
    // When "Alice" opens the "files" app
    await ui.userOpensApplication({ stepUser: 'Alice', name: 'files' })
    // And "Alice" closes the sidebar
    await ui.userClosesSidebar({ stepUser: 'Alice' })
    // Then "Alice" should see link-direct indicator on the folder "folderPublic"
    await ui.userShouldSeeShareIndicatorOnResource({
      stepUser: 'Alice',
      buttonLabel: 'link-direct',
      resource: 'folderPublic'
    })

    // When "Alice" opens folder "folderPublic"
    await ui.userOpensResource({ stepUser: 'Alice', resource: 'folderPublic' })
    // Then "Alice" should see link-indirect indicator on the file "lorem.txt"
    await ui.userShouldSeeShareIndicatorOnResource({
      stepUser: 'Alice',
      buttonLabel: 'link-indirect',
      resource: 'lorem.txt'
    })

    // And "Alice" navigates to the shared via link page
    await ui.userNavigatesToSharedViaLinkPage({ stepUser: 'Alice' })
    // Then following resources should be displayed in the files list for user "Alice"
    //   | resource     |
    //   | folderPublic |
    await ui.userShouldSeeResources({
      listType: 'files list',
      stepUser: 'Alice',
      resources: ['folderPublic']
    })

    // check copy link to clipboard button
    // When "Alice" opens the "files" app
    await ui.userOpensApplication({ stepUser: 'Alice', name: 'files' })
    // And "Alice" copies the link "Unnamed link" of resource "folderPublic"
    await ui.userCopiesLinkOfResource({ stepUser: 'Alice', resource: 'folderPublic' })

    // And "Alice" opens the "%clipboard%" url
    await ui.userOpensClipboardUrl({ stepUser: 'Alice', url: '%clipboard%' })
    // And "Alice" unlocks the public link with password "%public%"
    await ui.userUnlocksPublicLink({ password: '%public%', stepUser: 'Alice' })
    // Then following resources should be displayed in the files list for user "Alice"
    //   | resource  |
    //   | lorem.txt |
    await ui.userShouldSeeResources({
      listType: 'files list',
      stepUser: 'Alice',
      resources: ['lorem.txt']
    })

    // And "Alice" logs out
    await ui.userLogsOut({ stepUser: 'Alice' })
  })
})
