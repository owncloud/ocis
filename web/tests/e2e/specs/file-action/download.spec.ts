import { test } from '../../environment/test'
import * as api from '../../steps/api/api'
import * as ui from '../../steps/ui/index'
import { application, fileAction } from '../../environment/constants'

test.describe('Download', { tag: '@predefined-users' }, () => {
  test.beforeEach(async () => {
    await api.usersHaveBeenCreated({ stepUser: 'Admin', users: ['Alice', 'Brian'] })
  })

  test('download resources', async () => {
    // Given "Alice" logs in
    // Given "Brian" logs in
    await ui.userLogsIn({ stepUser: 'Alice' })
    await ui.userLogsIn({ stepUser: 'Brian' })

    // And "Alice" creates the following folders in personal space using API
    //   | name         |
    //   | folderPublic |
    //   | emptyFolder  |
    await api.userHasCreatedFolders({
      stepUser: 'Alice',
      folderNames: ['folderPublic', 'emptyFolder']
    })

    // And "Alice" creates the following files into personal space using API
    //   | pathToFile                | content     |
    //   | folderPublic/new file.txt | lorem ipsum |
    await api.userHasCreatedFiles({
      stepUser: 'Alice',
      files: [{ pathToFile: 'folderPublic/new file.txt', content: 'lorem ipsum' }]
    })

    // And "Alice" uploads the following local file into personal space using API
    //   | localFile                     | to             |
    //   | filesForUpload/testavatar.jpg | testavatar.jpg |
    await api.userHasUploadedFilesInPersonalSpace({
      stepUser: 'Alice',
      filesToUpload: [{ localFile: 'filesForUpload/testavatar.jpg', to: 'testavatar.jpg' }]
    })

    // And "Alice" shares the following resource using API
    //   | resource       | recipient | type | role                                | resourceType |
    //   | folderPublic   | Brian     | user | Can edit with versions and trashbin | folder       |
    //   | emptyFolder    | Brian     | user | Can edit with versions and trashbin | folder       |
    //   | testavatar.jpg | Brian     | user | Can edit with versions and trashbin | file         |
    await api.userHasSharedResources({
      stepUser: 'Alice',
      shares: [
        {
          resource: 'folderPublic',
          recipient: 'Brian',
          type: 'user',
          role: 'Can edit with versions and trash bin',
          resourceType: 'folder'
        },
        {
          resource: 'emptyFolder',
          recipient: 'Brian',
          type: 'user',
          role: 'Can edit with versions and trash bin',
          resourceType: 'folder'
        },
        {
          resource: 'testavatar.jpg',
          recipient: 'Brian',
          type: 'user',
          role: 'Can edit with versions and trash bin',
          resourceType: 'file'
        }
      ]
    })

    // When "Alice" downloads the following resources using the batch action
    //   | resource       | type   |
    //   | folderPublic   | folder |
    //   | emptyFolder    | folder |
    //   | testavatar.jpg | file   |
    const resourceToDownloadInBatch = [
      { resource: 'folderPublic', type: 'folder' },
      { resource: 'emptyFolder', type: 'folder' },
      { resource: 'testavatar.jpg', type: 'file' }
    ]
    await ui.userDownloadsResource({
      stepUser: 'Alice',
      resourceToDownload: resourceToDownloadInBatch,
      actionType: fileAction.batchAction
    })

    // And "Alice" opens the following file in mediaviewer
    //   | resource       |
    //   | testavatar.jpg |
    await ui.userOpensResourceInViewer({
      stepUser: 'Alice',
      resource: 'testavatar.jpg',
      viewer: application.mediaViewer
    })

    // And "Alice" downloads the following resources using the preview topbar
    //   | resource       | type |
    //   | testavatar.jpg | file |
    const downloadImage = [{ resource: 'testavatar.jpg', type: 'file' }]
    await ui.userDownloadsResource({
      stepUser: 'Alice',
      resourceToDownload: downloadImage,
      actionType: fileAction.previewTopBar
    })

    // And "Alice" closes the file viewer
    await ui.userClosesFileViewer({ stepUser: 'Alice' })

    // And "Alice" logs out
    await ui.userLogsOut({ stepUser: 'Alice' })

    // When "Brian" navigates to the shared with me page
    await ui.userNavigatesToSharedWithMePage({ stepUser: 'Brian' })
    // And "Brian" downloads the following resources using the batch action
    //   | resource       | type   |
    //   | folderPublic   | folder |
    //   | emptyFolder    | folder |
    //   | testavatar.jpg | file   |
    await ui.userDownloadsResource({
      stepUser: 'Brian',
      resourceToDownload: resourceToDownloadInBatch,
      actionType: fileAction.batchAction
    })

    // And "Brian" downloads the following resources using the sidebar panel
    //   | resource       | from         | type   |
    //   | new file.txt   | folderPublic | file   |
    //   | testavatar.jpg |              | file   |
    //   | folderPublic   |              | folder |
    //   | emptyFolder    |              | folder |
    const resourceToDownloadSidebar = [
      { resource: 'new file.txt', from: 'folderPublic', type: 'file' },
      { resource: 'testavatar.jpg', type: 'file' },
      { resource: 'folderPublic', type: 'folder' },
      { resource: 'emptyFolder', type: 'folder' }
    ]
    await ui.userDownloadsResource({
      stepUser: 'Brian',
      resourceToDownload: resourceToDownloadSidebar,
      actionType: fileAction.sideBarPanel
    })

    // And "Brian" opens the following file in mediaviewer
    //   | resource       |
    //   | testavatar.jpg |
    await ui.userOpensResourceInViewer({
      stepUser: 'Brian',
      resource: 'testavatar.jpg',
      viewer: application.mediaViewer
    })

    // And "Brian" downloads the following resources using the preview topbar
    //   | resource       | type |
    //   | testavatar.jpg | file |
    await ui.userDownloadsResource({
      stepUser: 'Brian',
      resourceToDownload: downloadImage,
      actionType: fileAction.previewTopBar
    })

    // And "Brian" logs out
    await ui.userLogsOut({ stepUser: 'Brian' })
  })
})
