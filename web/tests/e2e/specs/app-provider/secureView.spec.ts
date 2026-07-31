// To run this feature we need to run the external app-provider service along OnlyOffice, Collabora services

import { test } from '../../environment/test'
import * as api from '../../steps/api/api'
import * as ui from '../../steps/ui/index'
import { searchScope, application, fileAction } from '../../environment/constants'

test.describe('Secure view', { tag: '@predefined-users' }, () => {
  test.beforeEach(async () => {
    // Given "Admin" creates following users using API
    //   | id    |
    //   | Alice |
    //   | Brian |
    await api.usersHaveBeenCreated({ stepUser: 'Admin', users: ['Alice', 'Brian'] })

    // And "Alice" logs in
    await ui.userLogsIn({ stepUser: 'Alice' })

    // And "Brian" logs in
    await ui.userLogsIn({ stepUser: 'Brian' })

    // And "Alice" opens the "files" app
    await ui.userOpensApplication({ stepUser: 'Alice', name: 'files' })
  })

  test('open a secure view file with Collabora', async () => {
    // Given "Alice" creates the following folder in personal space using API
    //   | name          |
    //   | shared folder |
    await api.userHasCreatedFolder({ stepUser: 'Alice', folderName: 'shared folder' })

    // And "Alice" uploads the following local file into personal space using API
    //   | localFile                      | to                            |
    //   | filesForUpload/simple.pdf      | shared folder/simple.pdf      |
    //   | filesForUpload/testavatar.jpeg | shared folder/testavatar.jpeg |
    //   | filesForUpload/lorem.txt       | shared folder/lorem.txt       |
    await api.userHasUploadedFilesInPersonalSpace({
      stepUser: 'Alice',
      filesToUpload: [
        { localFile: 'filesForUpload/simple.pdf', to: 'shared folder/simple.pdf' },
        { localFile: 'filesForUpload/testavatar.jpeg', to: 'shared folder/testavatar.jpeg' },
        { localFile: 'filesForUpload/lorem.txt', to: 'shared folder/lorem.txt ' }
      ]
    })

    // And "Alice" creates the following resources
    //   | resource           | type         | content                 |
    //   | secureDocument.odt | OpenDocument | very important document |
    await ui.userCreatesResources({
      stepUser: 'Alice',
      resources: [
        { name: 'secureDocument.odt', type: 'OpenDocument', content: 'very important document' }
      ]
    })

    // And "Alice" shares the following resources using the sidebar panel
    //   | resource           | recipient | type | role              | resourceType |
    //   | secureDocument.odt | Brian     | user | Can view (secure) | file         |
    //   | shared folder      | Brian     | user | Can view (secure) | folder       |
    await ui.userSharesResources({
      stepUser: 'Alice',
      actionType: fileAction.sideBarPanel,
      shares: [
        {
          resource: 'secureDocument.odt',
          recipient: 'Brian',
          type: 'user',
          role: 'Can view (secure)',
          resourceType: 'file'
        },
        {
          resource: 'shared folder',
          recipient: 'Brian',
          type: 'user',
          role: 'Can view (secure)',
          resourceType: 'folder'
        }
      ]
    })

    // And "Alice" logs out
    await ui.userLogsOut({ stepUser: 'Alice' })

    // When "Brian" navigates to the shared with me page
    await ui.userNavigatesToSharedWithMePage({ stepUser: 'Brian' })

    // And "Brian" opens the following file in Collabora
    //   | resource           |
    //   | secureDocument.odt |
    await ui.userOpensResourceInViewer({
      stepUser: 'Brian',
      resource: 'secureDocument.odt',
      viewer: application.collabora
    })

    // we copy the contents of the file and compare the clipboard with the expected contents.
    // In case the user does not have download permissions and tries to copy file content, the clipboard should be set to “Copying from document disabled”.
    // Then "Brian" should see the content "Copying from the document disabled" in editor "Collabora"
    await ui.userShouldSeeContentInEditor({
      stepUser: 'Brian',
      expectedContent: 'Copying from the document disabled',
      editor: 'Collabora'
    })

    // And "Brian" closes the file viewer
    await ui.userClosesFileViewer({ stepUser: 'Brian' })

    // When "Brian" opens folder "shared folder"
    await ui.userOpensResource({ stepUser: 'Brian', resource: 'shared folder' })

    // And "Brian" opens the following file in Collabora
    //   | resource   |
    //   | simple.pdf |
    await ui.userOpensResourceInViewer({
      stepUser: 'Brian',
      resource: 'simple.pdf',
      viewer: application.collabora
    })

    // Then "Brian" should see the content "Copying from the document disabled" in editor "Collabora"
    await ui.userShouldSeeContentInEditor({
      stepUser: 'Brian',
      expectedContent: 'Copying from the document disabled',
      editor: 'Collabora'
    })

    // And "Brian" closes the file viewer
    await ui.userClosesFileViewer({ stepUser: 'Brian' })

    // And "Brian" opens the following file in Collabora
    //   | resource        |
    //   | testavatar.jpeg |
    await ui.userOpensResourceInViewer({
      stepUser: 'Brian',
      resource: 'testavatar.jpeg',
      viewer: application.collabora
    })

    // Then "Brian" should see the content "Copying from the document disabled" in editor "Collabora"
    await ui.userShouldSeeContentInEditor({
      stepUser: 'Brian',
      expectedContent: 'Copying from the document disabled',
      editor: 'Collabora'
    })

    // And "Brian" closes the file viewer
    await ui.userClosesFileViewer({ stepUser: 'Brian' })

    // And "Brian" opens the following file in Collabora
    //   | resource  |
    //   | lorem.txt |
    await ui.userOpensResourceInViewer({
      stepUser: 'Brian',
      resource: 'lorem.txt',
      viewer: application.collabora
    })

    // Then "Brian" should see the content "Copying from the document disabled" in editor "Collabora"
    await ui.userShouldSeeContentInEditor({
      stepUser: 'Brian',
      expectedContent: 'Copying from the document disabled',
      editor: 'Collabora'
    })

    // And "Brian" logs out
    await ui.userLogsOut({ stepUser: 'Brian' })
  })

  test('check available actions for secure view file', async () => {
    // Given "Alice" creates the following folder in personal space using API
    //   | name          |
    //   | shared folder |
    await api.userHasCreatedFolder({ stepUser: 'Alice', folderName: 'shared folder' })

    // And "Alice" uploads the following local file into personal space using API
    //   | localFile                      | to                             |
    //   | filesForUpload/simple.pdf      | shared folder/secure.pdf       |
    //   | filesForUpload/testavatar.jpeg | shared folder/securePhoto.jpeg |
    //   | filesForUpload/lorem.txt       | shared folder/secureFile.txt   |
    await api.userHasUploadedFilesInPersonalSpace({
      stepUser: 'Alice',
      filesToUpload: [
        { localFile: 'filesForUpload/simple.pdf', to: 'shared folder/secure.pdf' },
        { localFile: 'filesForUpload/testavatar.jpeg', to: 'shared folder/securePhoto.jpeg' },
        { localFile: 'filesForUpload/lorem.txt', to: 'shared folder/secureFile.txt ' }
      ]
    })

    // And "Alice" creates the following resources
    //   | resource           | type         | content                 |
    //   | secureDocument.odt | OpenDocument | very important document |
    await ui.userCreatesResources({
      stepUser: 'Alice',
      resources: [
        { name: 'secureDocument.odt', type: 'OpenDocument', content: 'very important document' }
      ]
    })

    // And "Alice" shares the following resources using the sidebar panel
    //   | resource           | recipient | type | role              | resourceType |
    //   | secureDocument.odt | Brian     | user | Can view (secure) | file         |
    //   | shared folder      | Brian     | user | Can view (secure) | folder       |
    await ui.userSharesResources({
      stepUser: 'Alice',
      actionType: fileAction.sideBarPanel,
      shares: [
        {
          resource: 'secureDocument.odt',
          recipient: 'Brian',
          type: 'user',
          role: 'Can view (secure)',
          resourceType: 'file'
        },
        {
          resource: 'shared folder',
          recipient: 'Brian',
          type: 'user',
          role: 'Can view (secure)',
          resourceType: 'folder'
        }
      ]
    })

    // And "Alice" logs out
    await ui.userLogsOut({ stepUser: 'Alice' })

    // When "Brian" navigates to the shared with me page
    await ui.userNavigatesToSharedWithMePage({ stepUser: 'Brian' })

    // .odt file
    // Then "Brian" should see following actions for file "secureDocument.odt"
    //   | action            |
    //   | Open in Collabora |
    await ui.userShouldSeeActionsForResource({
      stepUser: 'Brian',
      actions: ['Open in Collabora'],
      resource: 'secureDocument.odt'
    })

    // But "Brian" should not see following actions for file "secureDocument.odt"
    //   | action             |
    //   | Download           |
    //   | Copy               |
    //   | Open in OnlyOffice |
    await ui.userShouldNotSeeActionsForResource({
      stepUser: 'Brian',
      actions: ['Download', 'Copy', 'Open in OnlyOffice'],
      resource: 'secureDocument.odt'
    })

    // And "Brian" should not see preview for file "secureDocument.odt"
    await ui.userShouldNotSeePreviewForFile({ stepUser: 'Brian', resource: 'secureDocument.odt' })

    // folder
    // Then "Brian" should not see following actions for folder "shared folder"
    //   | action   |
    //   | Download |
    //   | Copy     |
    await ui.userShouldNotSeeActionsForResource({
      stepUser: 'Brian',
      actions: ['Download', 'Copy'],
      resource: 'shared folder'
    })

    // When "Brian" opens folder "shared folder"
    await ui.userOpensResource({ stepUser: 'Brian', resource: 'shared folder' })

    // .pdf file
    // Then "Brian" should see following actions for file "secure.pdf"
    //   | action            |
    //   | Open in Collabora |
    await ui.userShouldSeeActionsForResource({
      stepUser: 'Brian',
      actions: ['Open in Collabora'],
      resource: 'secure.pdf'
    })

    // But "Brian" should not see following actions for file "secure.pdf"
    //   | action             |
    //   | Download           |
    //   | Copy               |
    //   | Open in PDF Viewer |
    await ui.userShouldNotSeeActionsForResource({
      stepUser: 'Brian',
      actions: ['Download', 'Copy', 'Open in PDF Viewer'],
      resource: 'secure.pdf'
    })

    // And "Brian" should not see thumbnail and preview for file "secure.pdf"
    await ui.userShouldNotSeeThumbnailAndPreviewForFile({
      stepUser: 'Brian',
      resource: 'secure.pdf'
    })

    // .jpeg file
    // Then "Brian" should see following actions for file "securePhoto.jpeg"
    //   | action            |
    //   | Open in Collabora |
    await ui.userShouldSeeActionsForResource({
      stepUser: 'Brian',
      actions: ['Open in Collabora'],
      resource: 'securePhoto.jpeg'
    })

    // But "Brian" should not see following actions for file "securePhoto.jpeg"
    //   | action   |
    //   | Download |
    //   | Copy     |
    //   | Preview  |
    await ui.userShouldNotSeeActionsForResource({
      stepUser: 'Brian',
      actions: ['Download', 'Copy', 'Preview'],
      resource: 'securePhoto.jpeg'
    })

    // And "Brian" should not see thumbnail and preview for file "securePhoto.jpeg"
    await ui.userShouldNotSeeThumbnailAndPreviewForFile({
      stepUser: 'Brian',
      resource: 'securePhoto.jpeg'
    })

    // .txt file
    // Then "Brian" should see following actions for file "secureFile.txt"
    //   | action            |
    //   | Open in Collabora |
    await ui.userShouldSeeActionsForResource({
      stepUser: 'Brian',
      actions: ['Open in Collabora'],
      resource: 'secureFile.txt'
    })

    // But "Brian" should not see following actions for file "secureFile.txt"
    //   | action              |
    //   | Download            |
    //   | Copy                |
    //   | Open in Text Editor |
    //   | Open in OnlyOffice  |
    await ui.userShouldNotSeeActionsForResource({
      stepUser: 'Brian',
      actions: ['Download', 'Copy', 'Open in Text Editor', 'Open in OnlyOffice'],
      resource: 'secureFile.txt'
    })

    // And "Brian" should not see thumbnail and preview for file "secureFile.txt"
    await ui.userShouldNotSeeThumbnailAndPreviewForFile({
      stepUser: 'Brian',
      resource: 'secureFile.txt'
    })

    // When "Brian" searches "secure" using the global search and the "all files" filter and presses enter
    await ui.userSearchesGloballyWithFilter({
      stepUser: 'Brian',
      keyword: 'secure',
      filter: searchScope.allFiles,
      command: 'presses enter'
    })

    // Then following resources should be displayed in the files list for user "Brian"
    //   | resource           |
    //   | secureFile.txt     |
    //   | securePhoto.jpeg   |
    //   | secure.pdf         |
    //   | secureDocument.odt |
    await ui.userShouldSeeResources({
      listType: 'files list',
      stepUser: 'Brian',
      resources: ['secureFile.txt', 'securePhoto.jpeg', 'secure.pdf', 'secureDocument.odt']
    })

    // .txt file
    // Then "Brian" should see following actions for file "secureFile.txt"
    //   | action            |
    //   | Open in Collabora |
    await ui.userShouldSeeActionsForResource({
      stepUser: 'Brian',
      actions: ['Open in Collabora'],
      resource: 'secureFile.txt'
    })

    // But "Brian" should not see following actions for file "secureFile.txt"
    //   | action              |
    //   | Download            |
    //   | Copy                |
    //   | Open in Text Editor |
    //   | Open in OnlyOffice  |
    await ui.userShouldNotSeeActionsForResource({
      stepUser: 'Brian',
      actions: ['Download', 'Copy', 'Open in Text Editor', 'Open in OnlyOffice'],
      resource: 'secureFile.txt'
    })

    // And "Brian" should not see thumbnail and preview for file "secureFile.txt"
    await ui.userShouldNotSeeThumbnailAndPreviewForFile({
      stepUser: 'Brian',
      resource: 'secureFile.txt'
    })

    // .jpeg file
    // Then "Brian" should see following actions for file "securePhoto.jpeg"
    //   | action            |
    //   | Open in Collabora |
    await ui.userShouldSeeActionsForResource({
      stepUser: 'Brian',
      actions: ['Open in Collabora'],
      resource: 'securePhoto.jpeg'
    })

    // But "Brian" should not see following actions for file "securePhoto.jpeg"
    //   | action   |
    //   | Download |
    //   | Copy     |
    //   | Preview  |
    await ui.userShouldNotSeeActionsForResource({
      stepUser: 'Brian',
      actions: ['Download', 'Copy', 'Preview'],
      resource: 'securePhoto.jpeg'
    })

    // And "Brian" should not see thumbnail and preview for file "securePhoto.jpeg"
    await ui.userShouldNotSeeThumbnailAndPreviewForFile({
      stepUser: 'Brian',
      resource: 'securePhoto.jpeg'
    })

    // .pdf file
    // Then "Brian" should see following actions for file "secure.pdf"
    //   | action            |
    //   | Open in Collabora |
    await ui.userShouldSeeActionsForResource({
      stepUser: 'Brian',
      actions: ['Open in Collabora'],
      resource: 'secure.pdf'
    })

    // But "Brian" should not see following actions for file "secure.pdf"
    //   | action             |
    //   | Download           |
    //   | Copy               |
    //   | Open in PDF Viewer |
    await ui.userShouldNotSeeActionsForResource({
      stepUser: 'Brian',
      actions: ['Download', 'Copy', 'Open in PDF Viewer'],
      resource: 'secure.pdf'
    })

    // And "Brian" should not see preview for file "secure.pdf"
    await ui.userShouldNotSeePreviewForFile({ stepUser: 'Brian', resource: 'secure.pdf' })

    // .odt file
    // Then "Brian" should see following actions for file "secureDocument.odt"
    //   | action            |
    //   | Open in Collabora |
    await ui.userShouldSeeActionsForResource({
      stepUser: 'Brian',
      actions: ['Open in Collabora'],
      resource: 'secureDocument.odt'
    })

    // But "Brian" should not see following actions for file "secureDocument.odt"
    //   | action             |
    //   | Download           |
    //   | Copy               |
    //   | Open in OnlyOffice |
    await ui.userShouldNotSeeActionsForResource({
      stepUser: 'Brian',
      actions: ['Download', 'Copy', 'Open in OnlyOffice'],
      resource: 'secureDocument.odt'
    })

    // And "Brian" should not see preview for file "secureDocument.odt"
    await ui.userShouldNotSeeThumbnailAndPreviewForFile({
      stepUser: 'Brian',
      resource: 'secureDocument.odt'
    })

    // And "Brian" logs out
    await ui.userLogsOut({ stepUser: 'Brian' })
  })
})
