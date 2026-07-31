/**
 * To run this feature we need to run the external app-provider service along with wopi, OnlyOffice, Collabora services
 * This is a minimal test for the integration of ocis with different online office suites like Collabora and OnlyOffice
 */

import { test } from '../../environment/test'
import * as api from '../../steps/api/api'
import * as ui from '../../steps/ui/index'
import { fileAction, application, resourcePage } from '../../environment/constants'

test.describe('Integrate with online office suites like Collabora and OnlyOffice', () => {
  test.beforeEach(async () => {
    await api.usersHaveBeenCreated({ stepUser: 'Admin', users: ['Alice', 'Brian'] })
    await ui.userLogsIn({ stepUser: 'Alice' })
  })

  test('create an OpenDocument file with Collabora', { tag: '@predefined-users' }, async () => {
    // Given "Alice" opens the "files" app
    await ui.userOpensApplication({ stepUser: 'Alice', name: 'files' })

    // When "Alice" creates the following resources
    //   | resource         | type         | content              |
    //   | OpenDocument.odt | OpenDocument | OpenDocument Content |
    await ui.userCreatesResources({
      stepUser: 'Alice',
      resources: [
        { name: 'OpenDocument.odt', type: 'OpenDocument', content: 'OpenDocument Content' }
      ]
    })

    // And "Alice" creates a public link of following resource using the sidebar panel
    //   | resource         | role     | password |
    //   | OpenDocument.odt | Can edit | %public% |
    await ui.userCreatesPublicLink({
      stepUser: 'Alice',
      resource: 'OpenDocument.odt',
      role: 'Can edit',
      password: '%public%'
    })

    // And "Alice" opens the following file in Collabora
    //   | resource         |
    //   | OpenDocument.odt |
    await ui.userOpensResourceInViewer({
      stepUser: 'Alice',
      resource: 'OpenDocument.odt',
      viewer: application.collabora
    })

    // And "Anonymous" opens the public link "Unnamed link"
    await ui.userOpensPublicLink({ stepUser: 'Anonymous', name: 'Unnamed link' })

    // And "Anonymous" unlocks the public link with password "%public%"
    await ui.userUnlocksPublicLink({ stepUser: 'Anonymous', password: '%public%' })

    // Then "Anonymous" should see the content "OpenDocument Content" in editor "Collabora"
    await ui.userShouldSeeContentInEditor({
      stepUser: 'Anonymous',
      expectedContent: 'OpenDocument Content',
      editor: 'Collabora'
    })

    // When "Alice" edits the following resource
    //   | resource         | type         | content                           |
    //   | OpenDocument.odt | OpenDocument | Alice Edited OpenDocument Content |
    await ui.userEditsResources({
      stepUser: 'Alice',
      resources: [
        {
          name: 'OpenDocument.odt',
          type: 'OpenDocument',
          content: 'Alice Edited OpenDocument Content'
        }
      ]
    })

    // Then "Anonymous" should see the content "Alice Edited OpenDocument Content" in editor "Collabora"
    await ui.userShouldSeeContentInEditor({
      stepUser: 'Anonymous',
      expectedContent: 'Alice Edited OpenDocument Content',
      editor: 'Collabora'
    })

    // When "Anonymous" edits the following resource
    //   | resource         | type         | content                     |
    //   | OpenDocument.odt | OpenDocument | Edited OpenDocument Content |
    await ui.userEditsResources({
      stepUser: 'Anonymous',
      resources: [
        { name: 'OpenDocument.odt', type: 'OpenDocument', content: 'Edited OpenDocument Content' }
      ]
    })

    // Then "Alice" should see the content "Edited OpenDocument Content" in editor "Collabora"
    await ui.userShouldSeeContentInEditor({
      stepUser: 'Alice',
      expectedContent: 'Edited OpenDocument Content',
      editor: 'Collabora'
    })

    // And "Alice" closes the file viewer
    await ui.userClosesFileViewer({ stepUser: 'Alice' })

    // When "Alice" edits the public link named "Unnamed link" of resource "OpenDocument.odt" changing role to "Can view"
    await ui.userChangesRoleOfPublicLinkOfResource({
      stepUser: 'Alice',
      linkName: 'Unnamed link',
      resource: 'OpenDocument.odt',
      newRole: 'Can view'
    })

    // And "Anonymous" opens the public link "Unnamed link"
    await ui.userOpensPublicLink({ stepUser: 'Anonymous', name: 'Unnamed link' })

    // And "Anonymous" unlocks the public link with password "%public%"
    await ui.userUnlocksPublicLink({ stepUser: 'Anonymous', password: '%public%' })

    // Then "Anonymous" should not be able to edit content of following resource
    //   | resource         | type         | content                     |
    //   | OpenDocument.odt | OpenDocument | Edited OpenDocument Content |
    await ui.userShouldNotBeAbleToEditContentOfResources({
      stepUser: 'Anonymous',
      resources: [
        { name: 'OpenDocument.odt', type: 'OpenDocument', content: 'Edited OpenDocument Content' }
      ]
    })

    // When "Alice" shares the following resource using the sidebar panel
    //   | resource         | recipient | type | role     | resourceType |
    //   | OpenDocument.odt | Brian     | user | Can view | file         |
    await ui.userSharesResources({
      stepUser: 'Alice',
      actionType: fileAction.sideBarPanel,
      shares: [
        {
          resource: 'OpenDocument.odt',
          recipient: 'Brian',
          type: 'user',
          role: 'Can view',
          resourceType: 'file'
        }
      ]
    })

    // And "Brian" logs in
    await ui.userLogsIn({ stepUser: 'Brian' })

    // And "Brian" navigates to the shared with me page
    await ui.userNavigatesToSharedWithMePage({ stepUser: 'Brian' })

    // And "Brian" opens the following file in Collabora
    //   | resource         |
    //   | OpenDocument.odt |
    await ui.userOpensResourceInViewer({
      stepUser: 'Brian',
      resource: 'OpenDocument.odt',
      viewer: application.collabora
    })

    // Then "Brian" should not be able to edit content of following resource
    //   | resource         | type         | content                     |
    //   | OpenDocument.odt | OpenDocument | Edited OpenDocument Content |
    await ui.userShouldNotBeAbleToEditContentOfResources({
      stepUser: 'Brian',
      resources: [
        { name: 'OpenDocument.odt', type: 'OpenDocument', content: 'Edited OpenDocument Content' }
      ]
    })

    // And "Brian" closes the file viewer
    await ui.userClosesFileViewer({ stepUser: 'Brian' })

    // When "Alice" updates following sharee role
    //   | resource         | recipient | type | role                   | resourceType |
    //   | OpenDocument.odt | Brian     | user | Can edit with trashbin | file         |
    await ui.userUpdatesShareeRoles({
      stepUser: 'Alice',
      roleUpdates: [
        {
          resource: 'OpenDocument.odt',
          recipient: 'Brian',
          type: 'user',
          role: 'Can edit with trashbin',
          resourceType: 'file'
        }
      ]
    })

    // And "Alice" opens the following file in Collabora
    //   | resource         |
    //   | OpenDocument.odt |
    await ui.userOpensResourceInViewer({
      stepUser: 'Alice',
      resource: 'OpenDocument.odt',
      viewer: application.collabora
    })

    // And "Brian" opens the following file in Collabora
    //   | resource         |
    //   | OpenDocument.odt |
    await ui.userOpensResourceInViewer({
      stepUser: 'Brian',
      resource: 'OpenDocument.odt',
      viewer: application.collabora
    })

    // And "Alice" edits the following resource
    //   | resource         | type         | content                           |
    //   | OpenDocument.odt | OpenDocument | Alice Edited OpenDocument Content |
    await ui.userEditsResources({
      stepUser: 'Alice',
      resources: [
        {
          name: 'OpenDocument.odt',
          type: 'OpenDocument',
          content: 'Alice Edited OpenDocument Content'
        }
      ]
    })

    // Then "Brian" should see the content "Alice Edited OpenDocument Content" in editor "Collabora"
    await ui.userShouldSeeContentInEditor({
      stepUser: 'Brian',
      expectedContent: 'Alice Edited OpenDocument Content',
      editor: 'Collabora'
    })

    // When "Brian" edits the following resource
    //   | resource         | type         | content                           |
    //   | OpenDocument.odt | OpenDocument | Brian Edited OpenDocument Content |
    await ui.userEditsResources({
      stepUser: 'Brian',
      resources: [
        {
          name: 'OpenDocument.odt',
          type: 'OpenDocument',
          content: 'Brian Edited OpenDocument Content'
        }
      ]
    })

    // Then "Alice" should see the content "Brian Edited OpenDocument Content" in editor "Collabora"
    await ui.userShouldSeeContentInEditor({
      stepUser: 'Alice',
      expectedContent: 'Brian Edited OpenDocument Content',
      editor: 'Collabora'
    })

    // And "Brian" logs out
    await ui.userLogsOut({ stepUser: 'Brian' })

    // And "Alice" logs out
    await ui.userLogsOut({ stepUser: 'Alice' })
  })

  test('create a Microsoft Word file with OnlyOffice', { tag: '@predefined-users' }, async () => {
    // Given "Alice" opens the "files" app
    await ui.userOpensApplication({ stepUser: 'Alice', name: 'files' })

    // When "Alice" creates the following resources
    //   | resource           | type           | content                |
    //   | MicrosoftWord.docx | Microsoft Word | Microsoft Word Content |
    await ui.userCreatesResources({
      stepUser: 'Alice',
      resources: [
        { name: 'MicrosoftWord.docx', type: 'Microsoft Word', content: 'Microsoft Word Content' }
      ]
    })

    // And for "Alice" file "MicrosoftWord.docx" should not be locked
    await ui.resourceShouldNotBeLockedForUser({ stepUser: 'Alice', resource: 'MicrosoftWord.docx' })

    // And "Alice" creates a public link of following resource using the sidebar panel
    //   | resource           | role     | password |
    //   | MicrosoftWord.docx | Can edit | %public% |
    await ui.userCreatesPublicLink({
      stepUser: 'Alice',
      resource: 'MicrosoftWord.docx',
      role: 'Can edit',
      password: '%public%'
    })

    // And "Alice" opens the following file in OnlyOffice
    //   | resource           |
    //   | MicrosoftWord.docx |
    await ui.userOpensResourceInViewer({
      stepUser: 'Alice',
      resource: 'MicrosoftWord.docx',
      viewer: application.onlyOffice
    })

    // And "Anonymous" opens the public link "Unnamed link"
    await ui.userOpensPublicLink({ stepUser: 'Anonymous', name: 'Unnamed link' })

    // And "Anonymous" unlocks the public link with password "%public%"
    await ui.userUnlocksPublicLink({ password: '%public%', stepUser: 'Anonymous' })

    // Then "Anonymous" should see the content "Microsoft Word Content" in editor "OnlyOffice"
    await ui.userShouldSeeContentInEditor({
      stepUser: 'Anonymous',
      expectedContent: 'Microsoft Word Content',
      editor: 'OnlyOffice'
    })

    // When "Alice" edits the following resource
    //   | resource           | type           | content                             |
    //   | MicrosoftWord.docx | Microsoft Word | Alice Edited Microsoft Word Content |
    await ui.userEditsResources({
      stepUser: 'Alice',
      resources: [
        {
          name: 'MicrosoftWord.docx',
          type: 'Microsoft Word',
          content: 'Alice Edited Microsoft Word Content'
        }
      ]
    })

    // Then "Anonymous" should see the content "Alice Edited Microsoft Word Content" in editor "OnlyOffice"
    await ui.userShouldSeeContentInEditor({
      stepUser: 'Anonymous',
      expectedContent: 'Alice Edited Microsoft Word Content',
      editor: 'OnlyOffice'
    })

    // When "Anonymous" edits the following resource
    //   | resource           | type           | content                       |
    //   | MicrosoftWord.docx | Microsoft Word | Edited Microsoft Word Content |
    await ui.userEditsResources({
      stepUser: 'Anonymous',
      resources: [
        {
          name: 'Microsoft Word',
          type: 'Microsoft Word',
          content: 'Edited Microsoft Word Content'
        }
      ]
    })

    // Then "Alice" should see the content "Edited Microsoft Word Content" in editor "OnlyOffice"
    await ui.userShouldSeeContentInEditor({
      stepUser: 'Alice',
      expectedContent: 'Edited Microsoft Word Content',
      editor: 'OnlyOffice'
    })

    // And "Alice" closes the file viewer
    await ui.userClosesFileViewer({ stepUser: 'Alice' })

    // When "Alice" edits the public link named "Unnamed link" of resource "MicrosoftWord.docx" changing role to "Can view"
    await ui.userChangesRoleOfPublicLinkOfResource({
      stepUser: 'Alice',
      linkName: 'Unnamed link',
      resource: 'MicrosoftWord.docx',
      newRole: 'Can view'
    })

    // And "Anonymous" opens the public link "Unnamed link"
    await ui.userOpensPublicLink({ stepUser: 'Anonymous', name: 'Unnamed link' })

    // And "Anonymous" unlocks the public link with password "%public%"
    await ui.userUnlocksPublicLink({ stepUser: 'Anonymous', password: '%public%' })

    // Then "Anonymous" should not be able to edit content of following resource
    //   | resource           | type           | content                       |
    //   | MicrosoftWord.docx | Microsoft Word | Edited Microsoft Word Content |
    await ui.userShouldNotBeAbleToEditContentOfResources({
      stepUser: 'Anonymous',
      resources: [
        {
          name: 'MicrosoftWord.docx',
          type: 'Microsoft Word',
          content: 'Edited OpenDocument Content'
        }
      ]
    })

    // And for "Alice" file "MicrosoftWord.docx" should not be locked
    await ui.resourceShouldNotBeLockedForUser({ stepUser: 'Alice', resource: 'MicrosoftWord.docx' })

    // When "Alice" shares the following resource using the sidebar panel
    //   | resource           | recipient | type | role     | resourceType |
    //   | MicrosoftWord.docx | Brian     | user | Can view | file         |
    await ui.userSharesResources({
      stepUser: 'Alice',
      actionType: fileAction.sideBarPanel,
      shares: [
        {
          resource: 'MicrosoftWord.docx',
          recipient: 'Brian',
          type: 'user',
          role: 'Can view',
          resourceType: 'file'
        }
      ]
    })

    // And "Brian" logs in
    await ui.userLogsIn({ stepUser: 'Brian' })

    // And "Brian" navigates to the shared with me page
    await ui.userNavigatesToSharedWithMePage({ stepUser: 'Brian' })

    // And "Brian" opens the following file in OnlyOffice
    //   | resource           |
    //   | MicrosoftWord.docx |
    await ui.userOpensResourceInViewer({
      stepUser: 'Brian',
      resource: 'MicrosoftWord.docx',
      viewer: application.onlyOffice
    })

    // Then "Brian" should not be able to edit content of following resource
    //   | resource           | type           | content                       |
    //   | MicrosoftWord.docx | Microsoft Word | Edited Microsoft Word Content |
    await ui.userShouldNotBeAbleToEditContentOfResources({
      stepUser: 'Brian',
      resources: [
        {
          name: 'MicrosoftWord.docx',
          type: 'Microsoft Word',
          content: 'Edited Microsoft Word Content'
        }
      ]
    })

    // And "Brian" closes the file viewer
    await ui.userClosesFileViewer({ stepUser: 'Brian' })

    // When "Alice" updates following sharee role
    //   | resource           | recipient | type | role                   | resourceType |
    //   | MicrosoftWord.docx | Brian     | user | Can edit with trashbin | file         |
    await ui.userUpdatesShareeRoles({
      stepUser: 'Alice',
      roleUpdates: [
        {
          resource: 'MicrosoftWord.docx',
          recipient: 'Brian',
          type: 'user',
          role: 'Can edit with trashbin',
          resourceType: 'file'
        }
      ]
    })

    // And "Alice" opens the following file in OnlyOffice
    //   | resource           |
    //   | MicrosoftWord.docx |
    await ui.userOpensResourceInViewer({
      stepUser: 'Alice',
      resource: 'MicrosoftWord.docx',
      viewer: application.onlyOffice
    })

    // And "Brian" opens the following file in OnlyOffice
    //   | resource           |
    //   | MicrosoftWord.docx |
    await ui.userOpensResourceInViewer({
      stepUser: 'Brian',
      resource: 'MicrosoftWord.docx',
      viewer: application.onlyOffice
    })

    // And "Alice" edits the following resource
    //   | resource           | type           | content                             |
    //   | MicrosoftWord.docx | Microsoft Word | Alice Edited Microsoft Word Content |
    await ui.userEditsResources({
      stepUser: 'Alice',
      resources: [
        {
          name: 'MicrosoftWord.docx',
          type: 'Microsoft Word',
          content: 'Alice Edited Microsoft Word Content'
        }
      ]
    })

    // Then "Brian" should see the content "Alice Edited Microsoft Word Content" in editor "OnlyOffice"
    await ui.userShouldSeeContentInEditor({
      stepUser: 'Brian',
      expectedContent: 'Alice Edited Microsoft Word Content',
      editor: 'OnlyOffice'
    })

    // When "Brian" edits the following resource
    //   | resource           | type           | content                             |
    //   | MicrosoftWord.docx | Microsoft Word | Brian Edited Microsoft Word Content |
    await ui.userEditsResources({
      stepUser: 'Brian',
      resources: [
        {
          name: 'MicrosoftWord.docx',
          type: 'Microsoft Word',
          content: 'Brian Edited Microsoft Word Content'
        }
      ]
    })

    // Then "Alice" should see the content "Brian Edited Microsoft Word Content" in editor "OnlyOffice"
    await ui.userShouldSeeContentInEditor({
      stepUser: 'Alice',
      expectedContent: 'Brian Edited Microsoft Word Content',
      editor: 'OnlyOffice'
    })

    // And "Brian" logs out
    await ui.userLogsOut({ stepUser: 'Brian' })

    // And "Alice" logs out
    await ui.userLogsOut({ stepUser: 'Alice' })
  })

  test('public creates a Mircosoft Word file with OnlyOffice', async () => {
    // Given "Admin" assigns following roles to the users using API
    //   | id    | role        |
    //   | Alice | Space Admin |
    await api.userHasAssignedRolesToUsers({
      stepUser: 'Admin',
      users: [{ id: 'Alice', role: 'Space Admin' }]
    })

    // And "Alice" creates the following project space using API
    //   | name      | id          |
    //   | Marketing | marketing.1 |
    await api.userHasCreatedProjectSpaces({
      stepUser: 'Alice',
      spaces: [{ name: 'Marketing', id: 'marketing.1' }]
    })

    // And "Alice" creates the following folder in space "Marketing" using API
    //   | name     |
    //   | myfolder |
    await api.userHasCreatedFoldersInSpace({
      stepUser: 'Alice',
      spaceName: 'Marketing',
      folders: ['myfolder']
    })

    // When "Alice" navigates to the project space "marketing.1"
    await ui.userNavigatesToSpace({ stepUser: 'Alice', space: 'marketing.1' })

    // And "Alice" creates a public link for the space with password "%public%" using the sidebar panel
    await ui.userCreatesPublicLinkOfSpaceWithPassword({ stepUser: 'Alice', password: '%public%' })

    // And "Alice" renames the most recently created public link of space to "spaceLink"
    await ui.userRenamesTheMostRecentlyCreatedPublicLinkOfSpace({
      stepUser: 'Alice',
      newName: 'spaceLink'
    })

    // And "Alice" edits the public link named "spaceLink" of the space changing role to "Can edit"
    await ui.userEditsThePublicLinkOfSpaceChangingRole({
      stepUser: 'Alice',
      linkName: 'spaceLink',
      role: 'Can edit'
    })

    // And "Alice" creates a public link of following resource using the sidebar panel
    //   | resource | role     | password |
    //   | myfolder | Can edit | %public% |
    await ui.userCreatesPublicLink({
      stepUser: 'Alice',
      resource: 'myfolder',
      role: 'Can edit',
      password: '%public%'
    })

    // public create .docx file using spaceLink
    // When "Anonymous" opens the public link "spaceLink"
    await ui.userOpensPublicLink({ stepUser: 'Anonymous', name: 'spaceLink' })

    // And "Anonymous" unlocks the public link with password "%public%"
    await ui.userUnlocksPublicLink({ password: '%public%', stepUser: 'Anonymous' })

    // And "Anonymous" creates the following resources
    //   | resource            | type           | content                                                      |
    //   | usingSpaceLink.docx | Microsoft Word | public can create files in the project space using spaceLink |
    await ui.userCreatesResources({
      stepUser: 'Anonymous',
      resources: [
        {
          name: 'usingSpaceLink.docx',
          type: 'Microsoft Word',
          content: 'public can create files in the project space using spaceLink'
        }
      ]
    })

    // public create .docx file using folderLink
    // When "Anonymous" opens the public link "Unnamed link"
    await ui.userOpensPublicLink({ stepUser: 'Anonymous', name: 'Unnamed link' })

    // And "Anonymous" unlocks the public link with password "%public%"
    await ui.userUnlocksPublicLink({ password: '%public%', stepUser: 'Anonymous' })

    // And "Anonymous" creates the following resources
    //   | resource             | type           | content                |
    //   | usingFolderLink.docx | Microsoft Word | Microsoft Word Content |
    await ui.userCreatesResources({
      stepUser: 'Anonymous',
      resources: [
        { name: 'usingFolderLink.docx', type: 'Microsoft Word', content: 'Microsoft Word Content' }
      ]
    })

    // When "Alice" navigates to the project space "marketing.1"
    await ui.userNavigatesToSpace({ stepUser: 'Alice', space: 'marketing.1' })

    // And for "Alice" file "usingSpaceLink.docx" should not be locked
    await ui.resourceShouldNotBeLockedForUser({
      stepUser: 'Alice',
      resource: 'usingSpaceLink.docx'
    })

    // And "Alice" opens the following file in OnlyOffice
    //   | resource            |
    //   | usingSpaceLink.docx |
    await ui.userOpensResourceInViewer({
      stepUser: 'Alice',
      resource: 'usingSpaceLink.docx',
      viewer: application.onlyOffice
    })

    // Then "Alice" should see the content "public can create files in the project space using spaceLink" in editor "OnlyOffice"
    await ui.userShouldSeeContentInEditor({
      stepUser: 'Alice',
      expectedContent: 'public can create files in the project space using spaceLink',
      editor: 'OnlyOffice'
    })

    // And "Alice" closes the file viewer
    await ui.userClosesFileViewer({ stepUser: 'Alice' })

    // When "Alice" opens folder "myfolder"
    await ui.userOpensResource({ stepUser: 'Alice', resource: 'myfolder' })

    // And for "Alice" file "usingFolderLink.docx" should not be locked
    await ui.resourceShouldNotBeLockedForUser({
      stepUser: 'Alice',
      resource: 'usingFolderLink.docx'
    })

    // And "Alice" opens the following file in OnlyOffice
    //   | resource             |
    //   | usingFolderLink.docx |
    await ui.userOpensResourceInViewer({
      stepUser: 'Alice',
      resource: 'usingFolderLink.docx',
      viewer: application.onlyOffice
    })

    // Then "Alice" should see the content "Microsoft Word Content" in editor "OnlyOffice"
    await ui.userShouldSeeContentInEditor({
      stepUser: 'Alice',
      expectedContent: 'Microsoft Word Content',
      editor: 'OnlyOffice'
    })

    // And "Alice" logs out
    await ui.userLogsOut({ stepUser: 'Alice' })
  })

  test('public creates a Microsoft Word file with Collabora', async () => {
    // Given "Admin" assigns following roles to the users using API
    //   | id    | role        |
    //   | Alice | Space Admin |
    await api.userHasAssignedRolesToUsers({
      stepUser: 'Admin',
      users: [{ id: 'Alice', role: 'Space Admin' }]
    })

    // And "Alice" creates the following project space using API
    //   | name      | id          |
    //   | Marketing | marketing.1 |
    await api.userHasCreatedProjectSpaces({
      stepUser: 'Alice',
      spaces: [{ name: 'Marketing', id: 'marketing.1' }]
    })

    // And "Alice" creates the following folder in space "Marketing" using API
    //   | name     |
    //   | myfolder |
    await api.userHasCreatedFoldersInSpace({
      stepUser: 'Alice',
      spaceName: 'Marketing',
      folders: ['myfolder']
    })

    // When "Alice" navigates to the project space "marketing.1"
    await ui.userNavigatesToSpace({ stepUser: 'Alice', space: 'marketing.1' })

    // And "Alice" creates a public link for the space with password "%public%" using the sidebar panel
    await ui.userCreatesPublicLinkOfSpaceWithPassword({ stepUser: 'Alice', password: '%public%' })

    // And "Alice" renames the most recently created public link of space to "spaceLink"
    await ui.userRenamesTheMostRecentlyCreatedPublicLinkOfSpace({
      stepUser: 'Alice',
      newName: 'spaceLink'
    })

    // And "Alice" edits the public link named "spaceLink" of the space changing role to "Can edit"
    await ui.userEditsThePublicLinkOfSpaceChangingRole({
      stepUser: 'Alice',
      linkName: 'spaceLink',
      role: 'Can edit'
    })

    // And "Alice" creates a public link of following resource using the sidebar panel
    //   | resource | role     | password |
    //   | myfolder | Can edit | %public% |
    await ui.userCreatesPublicLink({
      stepUser: 'Alice',
      resource: 'myfolder',
      role: 'Can edit',
      password: '%public%'
    })

    // public create .odt file using spaceLink
    // When "Anonymous" opens the public link "spaceLink"
    await ui.userOpensPublicLink({ stepUser: 'Anonymous', name: 'spaceLink' })

    // And "Anonymous" unlocks the public link with password "%public%"
    await ui.userUnlocksPublicLink({ password: '%public%', stepUser: 'Anonymous' })

    // And "Anonymous" creates the following resources
    //   | resource           | type         | content                                                      |
    //   | usingSpaceLink.odt | OpenDocument | public can create files in the project space using spaceLink |
    await ui.userCreatesResources({
      stepUser: 'Anonymous',
      resources: [
        {
          name: 'usingSpaceLink.odt',
          type: 'OpenDocument',
          content: 'public can create files in the project space using spaceLink'
        }
      ]
    })

    // public create .odt file using folderLink
    // When "Anonymous" opens the public link "Unnamed link"
    await ui.userOpensPublicLink({ stepUser: 'Anonymous', name: 'Unnamed link' })

    // And "Anonymous" unlocks the public link with password "%public%"
    await ui.userUnlocksPublicLink({ stepUser: 'Anonymous', password: '%public%' })

    // And "Anonymous" creates the following resources
    //   | resource            | type         | content              |
    //   | usingFolderLink.odt | OpenDocument | OpenDocument Content |
    await ui.userCreatesResources({
      stepUser: 'Anonymous',
      resources: [
        { name: 'usingFolderLink.odt', type: 'OpenDocument', content: 'OpenDocument Content' }
      ]
    })

    // When "Alice" navigates to the project space "marketing.1"
    await ui.userNavigatesToSpace({ stepUser: 'Alice', space: 'marketing.1' })

    // And "Alice" opens the following file in Collabora
    //   | resource           |
    //   | usingSpaceLink.odt |
    await ui.userOpensResourceInViewer({
      stepUser: 'Alice',
      resource: 'usingSpaceLink.odt',
      viewer: application.collabora
    })

    // Then "Alice" should see the content "public can create files in the project space using spaceLink" in editor "Collabora"
    await ui.userShouldSeeContentInEditor({
      stepUser: 'Alice',
      expectedContent: 'public can create files in the project space using spaceLink',
      editor: 'Collabora'
    })

    // And "Alice" closes the file viewer
    await ui.userClosesFileViewer({ stepUser: 'Alice' })

    // When "Alice" opens folder "myfolder"
    await ui.userOpensResource({ stepUser: 'Alice', resource: 'myfolder' })

    // And "Alice" opens the following file in Collabora
    //   | resource            |
    //   | usingFolderLink.odt |
    await ui.userOpensResourceInViewer({
      stepUser: 'Alice',
      resource: 'usingFolderLink.odt',
      viewer: application.collabora
    })

    // Then "Alice" should see the content "OpenDocument Content" in editor "Collabora"
    await ui.userShouldSeeContentInEditor({
      stepUser: 'Alice',
      expectedContent: 'OpenDocument Content',
      editor: 'Collabora'
    })

    // And "Alice" logs out
    await ui.userLogsOut({ stepUser: 'Alice' })
  })

  test('create files from office templates', { tag: '@predefined-users' }, async () => {
    // Given "Alice" uploads the following local file into personal space using API
    //   | localFile     | to            |
    //   | Template.dotx | Template.dotx |
    //   | Template.ott  | Template.ott  |
    await api.userHasUploadedFilesInPersonalSpace({
      stepUser: 'Alice',
      filesToUpload: [
        { localFile: 'Template.dotx', to: 'Template.dotx' },
        { localFile: 'Template.ott', to: 'Template.ott' }
      ]
    })

    // And "Alice" opens the "files" app
    await ui.userOpensApplication({ stepUser: 'Alice', name: 'files' })

    // When "Alice" creates a file from template file "Template.dotx" via "OnlyOffice" using the sidebar panel
    await ui.userCreatesFileFromTemplateFile({
      stepUser: 'Alice',
      file: 'Template.dotx',
      webOffice: application.onlyOffice,
      actionType: fileAction.sideBarPanel
    })

    // Then "Alice" should see the content "As a user I want to create a document by clicking on a template file" in editor "OnlyOffice"
    await ui.userShouldSeeContentInEditor({
      stepUser: 'Alice',
      expectedContent: 'As a user I want to create a document by clicking on a template file',
      editor: application.onlyOffice
    })

    // And "Alice" closes the file viewer
    await ui.userClosesFileViewer({ stepUser: 'Alice' })

    // When "Alice" creates a file from template file "Template.ott" via "Collabora" using the context menu
    await ui.userCreatesFileFromTemplateFile({
      stepUser: 'Alice',
      file: 'Template.ott',
      webOffice: application.collabora,
      actionType: fileAction.contextMenu
    })

    // Then "Alice" should see the content "As a user I want to create a document by clicking on a template file" in editor "Collabora"
    await ui.userShouldSeeContentInEditor({
      stepUser: 'Alice',
      expectedContent: 'As a user I want to create a document by clicking on a template file',
      editor: application.collabora
    })

    // When "Alice" closes the file viewer
    await ui.userClosesFileViewer({ stepUser: 'Alice' })

    // And following resources should be displayed in the files list for user "Alice"
    //   | resource      |
    //   | Template.odt  |
    //   | Template.docx |
    await ui.userShouldSeeResources({
      listType: 'files list',
      stepUser: 'Alice',
      resources: ['Template.odt', 'Template.docx']
    })

    // When "Alice" opens file "Template.dotx"
    await ui.userOpensResource({ stepUser: 'Alice', resource: 'Template.dotx' })

    // Then "Alice" should see the content "As a user I want to create a document by clicking on a template file" in editor "OnlyOffice"
    await ui.userShouldSeeContentInEditor({
      stepUser: 'Alice',
      expectedContent: 'As a user I want to create a document by clicking on a template file',
      editor: application.onlyOffice
    })

    // When "Alice" closes the file viewer
    await ui.userClosesFileViewer({ stepUser: 'Alice' })

    // And following resources should be displayed in the files list for user "Alice"
    //   | resource          |
    //   | Template (1).docx |
    await ui.userShouldSeeResources({
      listType: 'files list',
      stepUser: 'Alice',
      resources: ['Template (1).docx']
    })

    // When "Alice" opens template file "Template.ott" via "Collabora" using the context menu
    await ui.userOpensTemplateFileUsingContextMenu({
      stepUser: 'Alice',
      file: 'Template.ott',
      webOffice: 'Collabora'
    })

    // Then "Alice" should see the content "As a user I want to create a document by clicking on a template file" in editor "Collabora"
    await ui.userShouldSeeContentInEditor({
      stepUser: 'Alice',
      expectedContent: 'As a user I want to create a document by clicking on a template file',
      editor: 'Collabora'
    })

    // When "Alice" closes the file viewer
    await ui.userClosesFileViewer({ stepUser: 'Alice' })

    // Then following resources should not be displayed in the files list for user "Alice"
    //   | resource          |
    //   | Template (2).docx |
    await ui.userShouldNotSeeTheResources({
      listType: resourcePage.filesList,
      stepUser: 'Alice',
      resources: ['Template (2).docx']
    })

    // And "Alice" logs out
    await ui.userLogsOut({ stepUser: 'Alice' })
  })
})
