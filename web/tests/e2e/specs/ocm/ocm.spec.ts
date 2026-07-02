import { test } from '../../environment/test'
import * as api from '../../steps/api/api'
import * as ui from '../../steps/ui/index'
import { application, fileAction } from '../../environment/constants'

test.describe('federation management', { tag: '@ocm' }, async () => {
  test('user creates a federated share', async () => {
    // Given using "LOCAL" server
    await ui.useServer({
      server: 'LOCAL'
    })

    // And "Admin" creates following user using API
    //   | id    |
    //   | Alice |
    await api.usersHaveBeenCreated({ stepUser: 'Admin', users: ['Alice'] })

    // And "Alice" creates the following folders in personal space using API
    //   | name         |
    //   | folderPublic |
    await api.userHasCreatedFolders({ stepUser: 'Alice', folderNames: ['folderPublic'] })

    // And "Alice" uploads the following local file into personal space using API
    //   | localFile                     | to                      |
    //   | filesForUpload/simple.pdf     | folderPublic/simple.pdf |
    //   | filesForUpload/testavatar.jpg | testavatar.jpg          |
    //   | filesForUpload/sampleGif.gif  | sampleGif.gif           |
    await api.userHasUploadedFilesInPersonalSpace({
      stepUser: 'Alice',
      filesToUpload: [
        { localFile: 'filesForUpload/simple.pdf', to: 'folderPublic/simple.gif' },
        { localFile: 'filesForUpload/testavatar.jpg', to: 'testavatar.jpg' },
        { localFile: 'filesForUpload/sampleGif.gif', to: 'sampleGif.gif' }
      ]
    })

    // And "Alice" logs in
    await ui.userLogsIn({ stepUser: 'Alice' })

    // And "Alice" opens the "open-cloud-mesh" app
    await ui.userOpensApplication({ stepUser: 'Alice', name: 'open-cloud-mesh' })

    // And "Alice" generates invitation token for the federation share
    await ui.userGeneratesInvitationTokenForTheFederationShare({ stepUser: 'Alice' })

    // Given using "FEDERATED" server
    await ui.useServer({
      server: 'FEDERATED'
    })

    // And "Admin" creates following user using API
    //   | id    |
    //   | Brian |
    await api.usersHaveBeenCreated({ stepUser: 'Admin', users: ['Brian'] })

    // And "Brian" logs in
    await ui.userLogsIn({ stepUser: 'Brian' })

    // And "Brian" opens the "open-cloud-mesh" app
    await ui.userOpensApplication({ stepUser: 'Brian', name: 'open-cloud-mesh' })

    // When "Brian" accepts federated share invitation by local user "Alice"
    await ui.userAcceptsFederatedShareInvitationByLocalUser({ stepUser: 'Brian', sharer: 'Alice' })

    // Then "Brian" should see the following federated connections:
    //   | user                     | email              |
    //   | %user_alice_displayName% | %user_alice_email% |
    await ui.userShouldSeeTheFederatedConnections({
      stepUser: 'Brian',
      federation: [{ user: '%user_alice_displayName%', email: '%user_alice_email%' }]
    })

    // And using "LOCAL" server
    await ui.useServer({
      server: 'LOCAL'
    })

    // When "Alice" reloads the page
    await ui.userReloadsPage({ stepUser: 'Alice' })

    // Then "Alice" should see the following federated connections:
    //   | user                     | email              |
    //   | %user_brian_displayName% | %user_brian_email% |
    await ui.userShouldSeeTheFederatedConnections({
      stepUser: 'Alice',
      federation: [{ user: '%user_brian_displayName%', email: '%user_brian_email%' }]
    })

    // And "Alice" opens the "files" app
    await ui.userOpensApplication({ stepUser: 'Alice', name: 'files' })

    // And "Alice" shares the following resource using the sidebar panel
    //   | resource       | recipient | type | role                      | resourceType | shareType |
    //   | folderPublic   | Brian     | user | Can edit with trashbin    | folder       | external  |
    //   | sampleGif.gif  | Brian     | user | Can edit with trashbin    | file         | external  |
    //   | testavatar.jpg | Brian     | user | Can view                  | file         | external  |
    await ui.userSharesResources({
      stepUser: 'Alice',
      actionType: fileAction.sideBarPanel,
      shares: [
        {
          resource: 'folderPublic',
          recipient: 'Brian',
          type: 'user',
          role: 'Can edit with trashbin',
          resourceType: 'folder',
          shareType: 'external'
        },
        {
          resource: 'sampleGif.gif',
          recipient: 'Brian',
          type: 'user',
          role: 'Can edit with trashbin',
          resourceType: 'file',
          shareType: 'external'
        },
        {
          resource: 'testavatar.jpg',
          recipient: 'Brian',
          type: 'user',
          role: 'Can view',
          resourceType: 'file',
          shareType: 'external'
        }
      ]
    })

    // And "Alice" checks the following access details of share "folderPublic" for user "Brian"
    //   | Name   | Brian Murphy |
    //   | Type   | External     |
    await ui.userChecksAccessDetailsOfShare({
      stepUser: 'Alice',
      resource: 'folderPublic',
      sharee: { name: 'brian', type: 'user' },
      accessDetails: { Name: 'Brian Murphy', Type: 'External' }
    })

    // And "Alice" should see the following access details of share "folderPublic" for federated user "Brian"
    //   | detail |
    //   | Domain |
    await ui.userShouldSeeAccessDetailsOfShareForFederatedUser({
      stepUser: 'Alice',
      resource: 'folderPublic',
      collaboratorName: 'Brian',
      detail: 'Domain'
    })

    // And "Alice" logs out
    await ui.userLogsOut({ stepUser: 'Alice' })

    // And using "FEDERATED" server
    await ui.useServer({
      server: 'FEDERATED'
    })

    // And "Brian" opens the "files" app
    await ui.userOpensApplication({ stepUser: 'Brian', name: 'files' })

    // When "Brian" navigates to the shared with me page
    await ui.userNavigatesToSharedWithMePage({ stepUser: 'Brian' })

    // Then "Brian" should see thumbnail and preview for file "testavatar.jpg"
    await ui.userShouldSeeThumbnailAndPreviewForFile({
      stepUser: 'Brian',
      resource: 'testavatar.jpg'
    })

    // And "Brian" should see thumbnail and preview for file "sampleGif.gif"
    await ui.userShouldSeeThumbnailAndPreviewForFile({
      stepUser: 'Brian',
      resource: 'sampleGif.gif'
    })

    // When "Brian" opens the following file in mediaviewer
    //   | resource       |
    //   | testavatar.jpg |
    await ui.userOpensResourceInViewer({
      stepUser: 'Brian',
      resource: 'testavatar.jpg',
      viewer: application.mediaViewer
    })

    // Then "Brian" is in a media-viewer
    await ui.userShouldBeInFileViewer({
      stepUser: 'Brian',
      fileViewerType: application.mediaViewer
    })

    // And "Brian" navigates to the next media resource
    await ui.userNavigatesToMediaResource({ stepUser: 'Brian', navigationType: 'next' })

    // And "Brian" closes the file viewer
    await ui.userClosesFileViewer({ stepUser: 'Brian' })

    // And "Brian" downloads the following resources using the sidebar panel
    //   | resource       | type   |
    //   | folderPublic   | folder |
    //   | sampleGif.gif  | file   |
    //   | testavatar.jpg | file   |
    await ui.userDownloadsResource({
      stepUser: 'Brian',
      resourceToDownload: [
        { resource: 'folderPublic', type: 'folder' },
        { resource: 'sampleGif.gif', type: 'file' },
        { resource: 'testavatar.jpg', type: 'file' }
      ],
      actionType: fileAction.sideBarPanel
    })

    // When "Brian" uploads the following resources
    //   | resource       | to           |
    //   | testavatar.png | folderPublic |
    await ui.userUploadsResources({
      stepUser: 'Brian',
      resources: [{ name: 'testavatar.png', to: 'folderPublic' }]
    })

    // And "Brian" opens folder "folderPublic"
    await ui.userOpensResource({ stepUser: 'Brian', resource: 'folderPublic' })

    // Then "Brian" should see thumbnail and preview for file "testavatar.png"
    await ui.userShouldSeeThumbnailAndPreviewForFile({
      stepUser: 'Brian',
      resource: 'testavatar.png'
    })

    // When "Brian" opens the following file in mediaviewer
    //   | resource       |
    //   | testavatar.png |
    await ui.userOpensResourceInViewer({
      stepUser: 'Brian',
      resource: 'testavatar.png',
      viewer: application.mediaViewer
    })

    // Then "Brian" is in a media-viewer
    await ui.userShouldBeInFileViewer({
      stepUser: 'Brian',
      fileViewerType: application.mediaViewer
    })

    // And "Brian" logs out
    await ui.userLogsOut({ stepUser: 'Brian' })
  })
})
