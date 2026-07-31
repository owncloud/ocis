import { test } from '../../environment/test'
import * as api from '../../steps/api/api.js'
import * as ui from '../../steps/ui/index'
import { application, fileAction } from '../../environment/constants'

test.describe('Navigate web directly through urls', () => {
  test('pagination', async () => {
    // Given "Admin" creates following user using API
    //   | id    |
    //   | Alice |
    //   | Brian |
    await api.usersHaveBeenCreated({ stepUser: 'Admin', users: ['Alice', 'Brian'] })
    // And "Admin" assigns following roles to the users using API
    //   | id    | role        |
    //   | Alice | Space Admin |
    await api.userHasAssignedRolesToUsers({
      stepUser: 'Admin',
      users: [{ id: 'Alice', role: 'Space Admin' }]
    })
    // And "Alice" creates the following folders in personal space using API
    //   | name   |
    //   | FOLDER |
    await api.userHasCreatedFolder({ stepUser: 'Alice', folderName: 'FOLDER' })
    // And "Alice" creates the following files into personal space using API
    //   | pathToFile                    | content      |
    //   | FOLDER/file_inside_folder.txt | example text |
    //   | lorem.txt                     | some content |
    //   | test.odt                      | some content |
    //   | lorem.txt                     | new content |
    await api.userHasCreatedFiles({
      stepUser: 'Alice',
      files: [
        {
          pathToFile: 'FOLDER/file_inside_folder.txt',
          content: 'example text'
        },
        {
          pathToFile: 'lorem.txt',
          content: 'some content'
        },
        {
          pathToFile: 'test.odt',
          content: 'some content'
        },
        {
          pathToFile: 'lorem.txt',
          content: 'new content'
        }
      ]
    })
    // And "Alice" creates the following project space using API
    //   | name        | id     |
    //   | Development | team.1 |
    await api.userHasCreatedProjectSpaces({
      stepUser: 'Alice',
      spaces: [{ name: 'Development', id: 'team.1' }]
    })
    // And "Alice" creates the following file in space "Development" using API
    //   | name              | content                   |
    //   | spaceTextfile.txt | This is test file. Cheers |
    await api.userHasCreatedFilesInsideSpace({
      stepUser: 'Alice',
      files: [
        {
          name: 'spaceTextfile.txt',
          space: 'Development',
          content: 'This is test file. Cheers'
        }
      ]
    })
    // And "Alice" logs in
    await ui.userLogsIn({ stepUser: 'Alice' })
    // When "Alice" navigates to "versions" details panel of file "lorem.txt" of space "personal" through the URL
    await ui.userOpensResourceDetailsPanelViaUrl({
      stepUser: 'Alice',
      resource: 'lorem.txt',
      detailsPanel: 'versions',
      space: 'personal'
    })
    // Then "Alice" restores following resources version
    //   | resource  | to | version | openDetailsPanel |
    //   | lorem.txt | /  | 1       | false            |
    await ui.userRestoresResourceVersion({
      stepUser: 'Alice',
      file: 'lorem.txt',
      to: '/',
      version: 1,
      openDetailsPanel: false
    })

    // When "Alice" navigates to "sharing" details panel of file "lorem.txt" of space "personal" through the URL
    await ui.userOpensResourceDetailsPanelViaUrl({
      stepUser: 'Alice',
      resource: 'lorem.txt',
      detailsPanel: 'sharing',
      space: 'personal'
    })
    // Then "Alice" shares the following resource using the direct url navigation
    //   | resource  | recipient | type | role     | resourceType |
    //   | lorem.txt | Brian     | user | Can view | file         |
    await ui.userSharesResources({
      stepUser: 'Alice',
      actionType: fileAction.urlNavigation,
      shares: [
        {
          resource: 'lorem.txt',
          recipient: 'Brian',
          type: 'user',
          role: 'Can view',
          resourceType: 'file'
        }
      ]
    })

    // file that has respective editor will open in the respective editor
    // When "Alice" opens the file "lorem.txt" of space "personal" through the URL
    await ui.userOpensSpaceResourceViaUrl({
      stepUser: 'Alice',
      resource: 'lorem.txt',
      space: 'personal'
    })

    // Then "Alice" is in a text-editor
    await ui.userShouldBeInFileViewer({ stepUser: 'Alice', fileViewerType: application.textEditor })

    // And "Alice" closes the file viewer
    await ui.userClosesFileViewer({ stepUser: 'Alice' })

    // file without the respective editor will show the file in the file list
    // When "Alice" opens the file "test.odt" of space "personal" through the URL
    await ui.userOpensSpaceResourceViaUrl({
      stepUser: 'Alice',
      resource: 'test.odt',
      space: 'personal'
    })
    // Then following resources should be displayed in the files list for user "Alice"
    //   | resource  |
    //   | FOLDER    |
    //   | lorem.txt |
    //   | test.odt  |
    await ui.userShouldSeeResources({
      listType: 'files list',
      stepUser: 'Alice',
      resources: ['FOLDER', 'lorem.txt', 'test.odt']
    })
    // When "Alice" opens the folder "FOLDER" of space "personal" through the URL
    await ui.userOpensSpaceResourceViaUrl({
      stepUser: 'Alice',
      resource: 'FOLDER',
      space: 'personal'
    })
    // And "Alice" opens the following file in texteditor
    //   | resource               |
    //   | file_inside_folder.txt |
    await ui.userOpensResourceInViewer({
      stepUser: 'Alice',
      resource: 'file_inside_folder.txt',
      viewer: application.textEditor
    })
    // Then "Alice" is in a text-editor
    await ui.userShouldBeInFileViewer({ stepUser: 'Alice', fileViewerType: application.textEditor })
    // And "Alice" closes the file viewer
    await ui.userClosesFileViewer({ stepUser: 'Alice' })
    // When "Alice" opens space "Development" through the URL
    await ui.userOpensSpaceViaUrl({ stepUser: 'Alice', space: 'Development' })
    // And "Alice" opens the following file in texteditor
    //   | resource          |
    //   | spaceTextfile.txt |
    await ui.userOpensResourceInViewer({
      stepUser: 'Alice',
      resource: 'spaceTextfile.txt',
      viewer: application.textEditor
    })
    // Then "Alice" is in a text-editor
    await ui.userShouldBeInFileViewer({ stepUser: 'Alice', fileViewerType: application.textEditor })
    // And "Alice" closes the file viewer
    await ui.userClosesFileViewer({ stepUser: 'Alice' })
    // And "Alice" logs out
    await ui.userLogsOut({ stepUser: 'Alice' })
  })
})
