import { test } from '../../environment/test'
import * as api from '../../steps/api/api.js'
import * as ui from '../../steps/ui/index'
import { application, fileAction } from '../../environment/constants'

test.describe(
  'Users can create shortcuts for resources and sites',
  { tag: '@predefined-users' },
  () => {
    test.beforeEach(async () => {
      //   Given "Admin" creates following users using API
      //    | id    |
      //    | Alice |
      //    | Brian |
      await api.usersHaveBeenCreated({ stepUser: 'Admin', users: ['Alice', 'Brian'] })
    })

    test('shortcut', async () => {
      // Given "Alice" logs in
      await ui.userLogsIn({ stepUser: 'Alice' })
      // And "Brian" logs in
      await ui.userLogsIn({ stepUser: 'Brian' })

      // And "Alice" creates the following folders in personal space using API
      //   | name |
      //   | docs |
      await api.userHasCreatedFolder({ stepUser: 'Alice', folderName: 'docs' })

      // And "Alice" creates the following files into personal space using API
      // | pathToFile      | content           |
      // | docs/notice.txt | important content |
      await api.userHasCreatedFiles({
        stepUser: 'Alice',
        files: [{ pathToFile: 'docs/notice.txt', content: 'important content' }]
      })

      // And "Alice" uploads the following local file into personal space using API
      // | localFile                     | to             |
      // | filesForUpload/testavatar.jpg | testavatar.jpg |
      await api.userHasUploadedFilesInPersonalSpace({
        stepUser: 'Alice',
        filesToUpload: [{ localFile: 'filesForUpload/testavatar.jpg', to: 'testavatar.jpg' }]
      })
      // And "Alice" shares the following resource using API
      // | resource       | recipient | type | role     | resourceType |
      // | testavatar.jpg | Brian     | user | Can view | file         |
      await api.userHasSharedResources({
        stepUser: 'Alice',
        shares: [
          {
            resource: 'testavatar.jpg',
            recipient: 'Brian',
            type: 'user',
            role: 'Can view',
            resourceType: 'file'
          }
        ]
      })
      // And "Alice" creates a public link of following resource using API
      // | resource        | password |
      // | docs/notice.txt | %public% |
      await api.userHasCreatedPublicLinkOfResource({
        stepUser: 'Alice',
        resource: 'docs/notice.txt',
        password: '%public%'
      })
      // And "Alice" renames the most recently created public link of resource "docs/notice.txt" to "myPublicLink"
      await ui.userRenamesMostRecentlyCreatedPublicLinkOfResource({
        stepUser: 'Alice',
        resource: 'docs/notice.txt',
        newName: 'myPublicLink'
      })
      // When "Alice" opens the "files" app
      await ui.userOpensApplication({ stepUser: 'Alice', name: 'files' })

      // # create a shortcut to file folder website
      // And "Alice" creates a shortcut for the following resources
      // | resource                   | name           | type    |
      // | notice.txt                 | important file | file    |
      // | docs                       |                | folder  |
      // | https://owncloud.com/news/ | companyNews    | website |
      await ui.userCreatesShortcutForResources({
        stepUser: 'Alice',
        resources: [
          { resource: 'notice.txt', name: 'important file', type: 'file' },
          { resource: 'docs', name: '', type: 'folder' },
          { resource: 'https://owncloud.com/blogs/', name: 'companyNews', type: 'website' }
        ]
      })

      // And "Alice" downloads the following resources using the sidebar panel
      // | resource           | type |
      // | important file.url | file |
      const resourceToDownload = [{ resource: 'important file.url', type: 'file' }]
      await ui.userDownloadsResource({
        stepUser: 'Alice',
        resourceToDownload: resourceToDownload,
        actionType: fileAction.sideBarPanel
      })

      // When "Alice" opens a shortcut "important file.url"
      await ui.userOpensShortcut({ stepUser: 'Alice', name: 'important file.url' })
      // Then "Alice" is in a text-editor
      await ui.userShouldBeInFileViewer({
        stepUser: 'Alice',
        fileViewerType: application.textEditor
      })
      // And "Alice" closes the file viewer
      await ui.userClosesFileViewer({ stepUser: 'Alice' })
      // And "Alice" opens the "files" app
      await ui.userOpensApplication({ stepUser: 'Alice', name: 'files' })
      // Then "Alice" can open a shortcut "companyNews.url" with external url "https://owncloud.com/news/"
      await ui.userCanOpenShortcutWithExternalUrl({
        stepUser: 'Alice',
        name: 'companyNews.url',
        url: 'https://owncloud.com/blogs/'
      })
      // And "Alice" logs out
      await ui.userLogsOut({ stepUser: 'Alice' })

      // # create a shortcut to the shared file
      // When "Brian" creates a shortcut for the following resources
      // | resource       | name | type |
      // | testavatar.jpg | logo | file |
      await ui.userCreatesShortcutForResources({
        stepUser: 'Brian',
        resources: [{ resource: 'testavatar.jpg', name: 'logo', type: 'file' }]
      })
      // And "Brian" opens a shortcut "logo.url"
      await ui.userOpensShortcut({ stepUser: 'Brian', name: 'logo.url' })
      // Then "Brian" is in a media-viewer
      await ui.userShouldBeInFileViewer({
        stepUser: 'Brian',
        fileViewerType: application.mediaViewer
      })
      // And "Brian" closes the file viewer
      await ui.userClosesFileViewer({ stepUser: 'Brian' })

      // # create a shortcut to the public link
      // When "Brian" opens the "files" app
      await ui.userOpensApplication({ stepUser: 'Brian', name: 'files' })
      // And "Brian" creates a shortcut for the following resources
      // | resource     | name             | type        |
      // | myPublicLink | linkToNoticeFile | public link |
      await ui.userCreatesShortcutForResources({
        stepUser: 'Brian',
        resources: [{ resource: 'myPublicLink', name: 'linkToNoticeFile', type: 'public link' }]
      })
      // And "Brian" opens a shortcut "linkToNoticeFile.url"
      await ui.userOpensShortcut({ stepUser: 'Brian', name: 'linkToNoticeFile.url' })
      // And "Brian" unlocks the public link with password "%public%"
      await ui.userUnlocksPublicLink({ password: '%public%', stepUser: 'Brian' })
      // Then "Brian" is in a text-editor
      await ui.userShouldBeInFileViewer({
        stepUser: 'Brian',
        fileViewerType: application.textEditor
      })
      // And "Brian" logs out
      await ui.userLogsOut({ stepUser: 'Brian' })
    })
  }
)
