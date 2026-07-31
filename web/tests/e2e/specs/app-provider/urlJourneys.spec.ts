// To run this feature we need to run the external app-provider service along with wopi, OnlyOffice, Collabora services
// This is a minimal test for the integration of ocis with different online office suites like Collabora and OnlyOffice
// Check that the file can be opened in collabora or onlyoffice using the url. https://github.com/owncloud/web/issues/9897

import { test } from '../../environment/test'
import * as api from '../../steps/api/api'
import * as ui from '../../steps/ui/index'
import { application, client } from '../../environment/constants'

test.describe('url stability for mobile and desktop client', { tag: '@predefined-users' }, () => {
  test.beforeEach(async () => {
    // Given "Admin" creates following users using API
    //   | id    |
    //   | Alice |
    await api.usersHaveBeenCreated({ stepUser: 'Admin', users: ['Alice'] })
    // And "Alice" logs in
    await ui.userLogsIn({ stepUser: 'Alice' })
    // And "Alice" creates the following files into personal space using API
    //   | pathToFile          | content                 |
    //   | OpenDocument.odt    | OpenDocument Content    |
    await api.userHasCreatedFiles({
      stepUser: 'Alice',
      files: [{ pathToFile: 'OpenDocument.odt', content: 'OpenDocument Content' }]
    })
    // And "Alice" creates the following resources
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
    // And "Alice" opens the "files" app
    await ui.userOpensApplication({ stepUser: 'Alice', name: 'files' })
  })

  test('open office suite files with Collabora and onlyOffice', async () => {
    // desktop feature
    // When "Alice" opens the file "OpenDocument.odt" of space "personal" in Collabora through the URL for desktop client
    await ui.userOpensResourceViaUrl({
      stepUser: 'Alice',
      resource: 'OpenDocument.odt',
      space: 'personal',
      editorName: application.collabora,
      clientType: client.desktop
    })
    // Then "Alice" should see the content "OpenDocument Content" in editor "Collabora"
    await ui.userShouldSeeContentInEditor({
      stepUser: 'Alice',
      expectedContent: 'OpenDocument Content',
      editor: 'Collabora'
    })
    // When "Alice" opens the file "MicrosoftWord.docx" of space "personal" in OnlyOffice through the URL for desktop client
    await ui.userOpensResourceViaUrl({
      stepUser: 'Alice',
      resource: 'MicrosoftWord.docx',
      space: 'personal',
      editorName: application.onlyOffice,
      clientType: client.desktop
    })
    // Then "Alice" should see the content "Microsoft Word Content" in editor "OnlyOffice"
    await ui.userShouldSeeContentInEditor({
      stepUser: 'Alice',
      expectedContent: 'Microsoft Word Content',
      editor: 'OnlyOffice'
    })

    // mobile feature
    // When "Alice" opens the file "OpenDocument.odt" of space "personal" in Collabora through the URL for mobile client
    await ui.userOpensResourceViaUrl({
      stepUser: 'Alice',
      resource: 'OpenDocument.odt',
      space: 'personal',
      editorName: application.collabora,
      clientType: client.mobile
    })
    // Then "Alice" should see the content "OpenDocument Content" in editor "Collabora"
    await ui.userShouldSeeContentInEditor({
      stepUser: 'Alice',
      expectedContent: 'OpenDocument Content',
      editor: 'Collabora'
    })
    // When "Alice" opens the file "MicrosoftWord.docx" of space "personal" in OnlyOffice through the URL for mobile client
    await ui.userOpensResourceViaUrl({
      stepUser: 'Alice',
      resource: 'MicrosoftWord.docx',
      space: 'personal',
      editorName: application.onlyOffice,
      clientType: client.mobile
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
})
