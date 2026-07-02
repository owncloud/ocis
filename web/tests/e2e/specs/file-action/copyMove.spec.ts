import { test } from '../../environment/test'
import * as api from '../../steps/api/api'
import * as ui from '../../steps/ui/index'
import { fileAction, resourcePage } from '../../environment/constants'

test.describe('file action - copy/move', { tag: '@predefined-users' }, () => {
  test.beforeEach(async () => {
    await api.usersHaveBeenCreated({ stepUser: 'Admin', users: ['Alice'] })
    await ui.userLogsIn({ stepUser: 'Alice' })
  })

  test('Users can copy a file from one folder to another', async () => {
    // Given "Alice" creates the following folders in personal space using API
    //   | name                        |
    //   | PARENTCopy1                 |
    //   | PARENTCopy2                 |
    //   | PARENTMove                  |
    //   | PARENTCopy3                 |
    //   | PARENTCopy4/Sub1/Sub2       |
    //   | PARENT                      |
    //   | PARENT/Sub1/Sub             |
    //   | PARENT/Sub2                 |
    //   | PARENT/Sub3                 |
    //   | PARENT/Sub4                 |
    //   | PARENT/Sub5                 |
    //   | Duplicate                   |
    //   | Duplicate/folderToDuplicate |
    await api.userHasCreatedFolders({
      stepUser: 'Alice',
      folderNames: [
        'PARENTCopy1',
        'PARENTCopy2',
        'PARENTMove',
        'PARENTCopy3',
        'PARENTCopy4/Sub1/Sub2',
        'PARENT',
        'PARENT/Sub1/Sub',
        'PARENT/Sub2',
        'PARENT/Sub3',
        'PARENT/Sub4',
        'PARENT/Sub5',
        'Duplicate',
        'Duplicate/folderToDuplicate'
      ]
    })

    // And "Alice" creates the following files into personal space using API
    //   | pathToFile               | content                             |
    //   | PARENTCopy3/example1.txt | example text                        |
    //   | PARENTCopy3/example2.txt | example text                        |
    //   | KeyboardExample.txt      | copy with the help of keyboard      |
    //   | dragDrop.txt             | copy with the help of drag-drop     |
    //   | sidebar.txt              | copy with the help of sidebar panel |
    //   | duplicate.txt            | duplicate file                      |
    //   | Duplicate/duplicate.txt  | duplicate file                      |
    //   | PARENT/fileToCopy1.txt   | some content                        |
    //   | PARENT/fileToCopy2.txt   | some content                        |
    //   | PARENT/fileToCopy3.txt   | some content                        |
    //   | PARENT/fileToCopy4.txt   | some content                        |
    //   | PARENT/fileToCopy5.txt   | some content
    await api.userHasCreatedFiles({
      stepUser: 'Alice',
      files: [
        { pathToFile: 'PARENTCopy3/example1.txt', content: 'example text' },
        { pathToFile: 'PARENTCopy3/example2.txt', content: 'example text' },
        { pathToFile: 'KeyboardExample.txt', content: 'copy with the help of keyboard' },
        { pathToFile: 'dragDrop.txt', content: 'copy with the help of drag-drop' },
        { pathToFile: 'sidebar.txt', content: 'copy with the help of sidebar panel' },
        { pathToFile: 'duplicate.txt', content: 'duplicate file' },
        { pathToFile: 'Duplicate/duplicate.txt', content: 'duplicate file' },
        { pathToFile: 'PARENT/fileToCopy1.txt', content: 'some content' },
        { pathToFile: 'PARENT/fileToCopy2.txt', content: 'some content' },
        { pathToFile: 'PARENT/fileToCopy3.txt', content: 'some content' },
        { pathToFile: 'PARENT/fileToCopy4.txt', content: 'some content' },
        { pathToFile: 'PARENT/fileToCopy5.txt', content: 'some content' }
      ]
    })

    // When "Alice" duplicates the following resource using sidebar-panel
    //   | resource      |
    //   | duplicate.txt |
    await ui.userDuplicatesResources({
      stepUser: 'Alice',
      method: fileAction.sideBarPanel,
      resources: ['duplicate.txt']
    })

    // And "Alice" duplicates the following resource using dropdown-menu
    //   | resource  |
    //   | Duplicate |
    await ui.userDuplicatesResources({
      stepUser: 'Alice',
      method: fileAction.dropDownMenu,
      resources: ['Duplicate']
    })
    // Then following resources should be displayed in the files list for user "Alice"
    //   | resource          |
    //   | duplicate (1).txt |
    //   | Duplicate (1)     |
    await ui.userShouldSeeResources({
      listType: resourcePage.filesList,
      stepUser: 'Alice',
      resources: ['duplicate (1).txt', 'Duplicate (1)']
    })
    // When "Alice" opens folder "Duplicate"
    await ui.userOpensResource({ stepUser: 'Alice', resource: 'Duplicate' })
    // When "Alice" duplicates the following resource using batch-action
    //   | resource      |
    //   | duplicate.txt |
    await ui.userDuplicatesResources({
      stepUser: 'Alice',
      method: fileAction.batchAction,
      resources: ['duplicate.txt']
    })
    // And "Alice" duplicates the following resource at once using batch-action
    //   | resource          |
    //   | folderToDuplicate |
    //   | duplicate.txt     |
    await ui.userDuplicatesResources({
      stepUser: 'Alice',
      method: fileAction.batchAction,
      resources: ['folderToDuplicate', 'duplicate.txt']
    })
    // And "Alice" duplicates the following resource at once using dropdown-menu
    //   | resource          |
    //   | folderToDuplicate |
    //   | duplicate.txt     |
    await ui.userDuplicatesResources({
      stepUser: 'Alice',
      method: fileAction.dropDownMenu,
      resources: ['folderToDuplicate', 'duplicate.txt']
    })
    // Then following resources should be displayed in the files list for user "Alice"
    //   | resource              |
    //   | duplicate (1).txt     |
    //   | duplicate (2).txt     |
    //   | duplicate (3).txt     |
    //   | folderToDuplicate (1) |
    //   | folderToDuplicate (2) |
    await ui.userShouldSeeResources({
      listType: 'files list',
      stepUser: 'Alice',
      resources: [
        'duplicate (1).txt',
        'duplicate (2).txt',
        'duplicate (3).txt',
        'folderToDuplicate (1)',
        'folderToDuplicate (2)'
      ]
    })
    // And "Alice" opens the "files" app
    await ui.userOpensApplication({ stepUser: 'Alice', name: 'files' })
    // When "Alice" copies the following resource using sidebar-panel
    //   | resource    | to          |
    //   | sidebar.txt | PARENTCopy2 |
    await ui.userCopiesResources({
      stepUser: 'Alice',
      actionType: fileAction.sideBarPanel,
      resources: [{ resource: 'sidebar.txt', to: 'PARENTCopy2' }]
    })
    // And "Alice" copies the following resource using dropdown-menu
    //   | resource                 | to          |
    //   | PARENTCopy3/example1.txt | PARENTCopy1 |
    await ui.userCopiesResources({
      stepUser: 'Alice',
      actionType: fileAction.dropDownMenu,
      resources: [{ resource: 'PARENTCopy3/example1.txt', to: 'PARENTCopy1' }]
    })
    // And "Alice" copies the following resource using batch-action
    //   | resource                 | to          |
    //   | PARENTCopy3/example2.txt | PARENTCopy1 |
    await ui.userCopiesResources({
      stepUser: 'Alice',
      actionType: fileAction.batchAction,
      resources: [{ resource: 'PARENTCopy3/example2.txt', to: 'PARENTCopy1' }]
    })
    // And "Alice" copies the following resource using keyboard
    //   | resource            | to          |
    //   | KeyboardExample.txt | PARENTCopy3 |
    await ui.userCopiesResources({
      stepUser: 'Alice',
      actionType: fileAction.keyboard,
      resources: [{ resource: 'KeyboardExample.txt', to: 'PARENTCopy3' }]
    })
    // And "Alice" moves the following resource using drag-drop
    //   | resource     | to          |
    //   | dragDrop.txt | PARENTCopy2 |
    await ui.userMovesResources({
      stepUser: 'Alice',
      actionType: fileAction.dragDrop,
      resources: [{ resource: 'dragDrop.txt', to: 'PARENTCopy2' }]
    })
    // And "Alice" moves the following resource using dropdown-menu
    //   | resource                 | to         |
    //   | PARENTCopy1/example1.txt | PARENTMove |
    await ui.userMovesResources({
      stepUser: 'Alice',
      actionType: fileAction.dropDownMenu,
      resources: [{ resource: 'PARENTCopy1/example1.txt', to: 'PARENTMove' }]
    })
    // And "Alice" moves the following resource using batch-action
    //   | resource                 | to         |
    //   | PARENTCopy1/example2.txt | PARENTMove |
    await ui.userMovesResources({
      stepUser: 'Alice',
      actionType: fileAction.batchAction,
      resources: [{ resource: 'PARENTCopy1/example2.txt', to: 'PARENTMove' }]
    })
    // And "Alice" moves the following resource using keyboard
    //   | resource    | to         |
    //   | PARENTCopy2 | PARENTMove |
    await ui.userMovesResources({
      stepUser: 'Alice',
      actionType: fileAction.keyboard,
      resources: [{ resource: 'PARENTCopy2', to: 'PARENTMove' }]
    })
    // And "Alice" moves the following resource using sidebar-panel
    //   | resource    | to         |
    //   | PARENTCopy3 | PARENTMove |
    await ui.userMovesResources({
      stepUser: 'Alice',
      actionType: fileAction.sideBarPanel,
      resources: [{ resource: 'PARENTCopy3', to: 'PARENTMove' }]
    })
    // And "Alice" opens folder "PARENTCopy4"
    await ui.userOpensResource({ stepUser: 'Alice', resource: 'PARENTCopy4' })
    // And "Alice" opens folder "Sub1"
    await ui.userOpensResource({ stepUser: 'Alice', resource: 'Sub1' })
    // And "Alice" moves the following resource using drag-drop-breadcrumb
    //   | resource | to          |
    //   | Sub2     | PARENTCopy4 |
    await ui.userMovesResources({
      stepUser: 'Alice',
      actionType: fileAction.dragDropBreadcrumb,
      resources: [{ resource: 'Sub2', to: 'PARENTCopy4' }]
    })

    // And "Alice" opens the "files" app
    await ui.userOpensApplication({ stepUser: 'Alice', name: 'files' })
    // And "Alice" opens folder "PARENT"
    await ui.userOpensResource({ stepUser: 'Alice', resource: 'PARENT' })
    // And "Alice" copies the following resources to "PARENT/Sub1" at once using dropdown-menu
    //   | resource        |
    //   | fileToCopy1.txt |
    //   | fileToCopy2.txt |
    //   | fileToCopy3.txt |
    //   | fileToCopy4.txt |
    //   | fileToCopy5.txt |
    //   | Sub4            |
    //   | Sub5            |
    await ui.userCopiesResourcesAtOnce({
      stepUser: 'Alice',
      newLocation: 'PARENT/Sub1',
      method: fileAction.dropDownMenu,
      resources: [
        'fileToCopy1.txt',
        'fileToCopy2.txt',
        'fileToCopy3.txt',
        'fileToCopy4.txt',
        'fileToCopy5.txt',
        'Sub4',
        'Sub5'
      ]
    })
    // And "Alice" copies the following resources to "PARENT/Sub2" at once using batch-action
    //   | resource        |
    //   | fileToCopy1.txt |
    //   | fileToCopy2.txt |
    //   | fileToCopy3.txt |
    //   | fileToCopy4.txt |
    //   | fileToCopy5.txt |
    //   | Sub4            |
    //   | Sub5            |
    await ui.userCopiesResourcesAtOnce({
      stepUser: 'Alice',
      newLocation: 'PARENT/Sub2',
      method: fileAction.batchAction,
      resources: [
        'fileToCopy1.txt',
        'fileToCopy2.txt',
        'fileToCopy3.txt',
        'fileToCopy4.txt',
        'fileToCopy5.txt',
        'Sub4',
        'Sub5'
      ]
    })
    // And "Alice" copies the following resources to "PARENT/Sub3" at once using keyboard
    //   | resource        |
    //   | fileToCopy1.txt |
    //   | fileToCopy2.txt |
    //   | fileToCopy3.txt |
    //   | fileToCopy4.txt |
    //   | fileToCopy5.txt |
    //   | Sub4            |
    //   | Sub5            |
    await ui.userCopiesResourcesAtOnce({
      stepUser: 'Alice',
      newLocation: 'PARENT/Sub3',
      method: fileAction.keyboard,
      resources: [
        'fileToCopy1.txt',
        'fileToCopy2.txt',
        'fileToCopy3.txt',
        'fileToCopy4.txt',
        'fileToCopy5.txt',
        'Sub4',
        'Sub5'
      ]
    })
    // And "Alice" opens folder "Sub1"
    await ui.userOpensResource({ stepUser: 'Alice', resource: 'Sub1' })
    // And "Alice" moves the following resources to "PARENT/Sub1/Sub" at once using dropdown-menu
    //   | resource        |
    //   | fileToCopy1.txt |
    //   | fileToCopy2.txt |
    //   | fileToCopy3.txt |
    //   | fileToCopy4.txt |
    //   | fileToCopy5.txt |
    //   | Sub4            |
    //   | Sub5            |
    await ui.userMovesResourcesAtOnce({
      stepUser: 'Alice',
      newLocation: 'PARENT/Sub1/Sub',
      method: fileAction.dropDownMenu,
      resources: [
        'fileToCopy1.txt',
        'fileToCopy2.txt',
        'fileToCopy3.txt',
        'fileToCopy4.txt',
        'fileToCopy5.txt',
        'Sub4',
        'Sub5'
      ]
    })
    // And "Alice" opens folder "Sub"
    await ui.userOpensResource({ stepUser: 'Alice', resource: 'Sub' })
    // And "Alice" moves the following resources to "PARENT/Sub1" at once using batch-action
    //   | resource        |
    //   | fileToCopy1.txt |
    //   | fileToCopy2.txt |
    //   | fileToCopy3.txt |
    //   | fileToCopy4.txt |
    //   | fileToCopy5.txt |
    //   | Sub4            |
    //   | Sub5            |
    await ui.userMovesResourcesAtOnce({
      stepUser: 'Alice',
      newLocation: 'PARENT/Sub1',
      method: fileAction.batchAction,
      resources: [
        'fileToCopy1.txt',
        'fileToCopy2.txt',
        'fileToCopy3.txt',
        'fileToCopy4.txt',
        'fileToCopy5.txt',
        'Sub4',
        'Sub5'
      ]
    })
    // And "Alice" opens the "files" app
    await ui.userOpensApplication({ stepUser: 'Alice', name: 'files' })
    // And "Alice" opens folder "PARENT"
    await ui.userOpensResource({ stepUser: 'Alice', resource: 'PARENT' })
    // And "Alice" moves the following resources to "Sub4" at once using drag-drop
    //   | resource        |
    //   | fileToCopy1.txt |
    //   | fileToCopy2.txt |
    //   | fileToCopy3.txt |
    //   | fileToCopy4.txt |
    //   | fileToCopy5.txt |
    //   | Sub1            |
    //   | Sub2            |
    await ui.userMovesResourcesAtOnce({
      stepUser: 'Alice',
      newLocation: 'Sub4',
      method: fileAction.dragDrop,
      resources: [
        'fileToCopy1.txt',
        'fileToCopy2.txt',
        'fileToCopy3.txt',
        'fileToCopy4.txt',
        'fileToCopy5.txt',
        'Sub1',
        'Sub2'
      ]
    })
    // And "Alice" opens folder "Sub4"
    await ui.userOpensResource({ stepUser: 'Alice', resource: 'Sub4' })
    // And "Alice" moves the following resources to "PARENT" at once using drag-drop-breadcrumb
    //   | resource        |
    //   | fileToCopy1.txt |
    //   | fileToCopy2.txt |
    //   | fileToCopy3.txt |
    //   | fileToCopy4.txt |
    //   | fileToCopy5.txt |
    //   | Sub1            |
    //   | Sub2            |
    await ui.userMovesResourcesAtOnce({
      stepUser: 'Alice',
      newLocation: 'PARENT',
      method: fileAction.dragDropBreadcrumb,
      resources: [
        'fileToCopy1.txt',
        'fileToCopy2.txt',
        'fileToCopy3.txt',
        'fileToCopy4.txt',
        'fileToCopy5.txt',
        'Sub1',
        'Sub2'
      ]
    })
    // And "Alice" opens the "files" app
    await ui.userOpensApplication({ stepUser: 'Alice', name: 'files' })
    // And "Alice" opens folder "PARENT"
    await ui.userOpensResource({ stepUser: 'Alice', resource: 'PARENT' })
    // And "Alice" moves the following resources to "PARENT/Sub4" at once using keyboard
    //   | resource        |
    //   | fileToCopy1.txt |
    //   | fileToCopy2.txt |
    //   | fileToCopy3.txt |
    //   | fileToCopy4.txt |
    //   | fileToCopy5.txt |
    //   | Sub1            |
    //   | Sub2            |
    await ui.userMovesResourcesAtOnce({
      stepUser: 'Alice',
      newLocation: 'PARENT/Sub4',
      method: 'keyboard',
      resources: [
        'fileToCopy1.txt',
        'fileToCopy2.txt',
        'fileToCopy3.txt',
        'fileToCopy4.txt',
        'fileToCopy5.txt',
        'Sub1',
        'Sub2'
      ]
    })
    // And "Alice" logs out
    await ui.userLogsOut({ stepUser: 'Alice' })
  })

  test('Copy and move resources with same name in personal space', async () => {
    // And "Alice" creates the following folders in personal space using API
    //   | name         |
    //   | sub          |
    //   | folder1      |
    //   | sub/folder1  |
    //   | sub1/folder1 |
    await api.userHasCreatedFolders({
      stepUser: 'Alice',
      folderNames: ['sub', 'folder1', 'sub/folder1', 'sub1/folder1']
    })
    // And "Alice" creates the following files into personal space using API
    //   | pathToFile                | content                 |
    //   | example1.txt              | personal space location |
    //   | folder1/example1.txt      | folder1 location        |
    //   | sub/folder1/example1.txt  | sub/folder1 location    |
    //   | sub1/folder1/example1.txt | sub1/folder1 location   |
    await api.userHasCreatedFiles({
      stepUser: 'Alice',
      files: [
        { pathToFile: 'example1.txt', content: 'personal space location' },
        { pathToFile: 'folder1/example1.txt', content: 'folder1 location' },
        { pathToFile: 'sub/folder1/example1.txt', content: 'sub/folder1 location' },
        { pathToFile: 'sub1/folder1/example1.txt', content: 'sub1/folder1 location' }
      ]
    })
    // copy and move file
    // When "Alice" copies the following resource using sidebar-panel
    //   | resource     | to      | option    |
    //   | example1.txt | folder1 | keep both |
    //   | example1.txt | folder1 | replace   |
    await ui.userCopiesResources({
      stepUser: 'Alice',
      actionType: fileAction.sideBarPanel,
      resources: [
        { resource: 'example1.txt', to: 'folder1', option: 'keep both' },
        { resource: 'example1.txt', to: 'folder1', option: 'replace' }
      ]
    })
    // And "Alice" moves the following resource using sidebar-panel
    //   | resource             | to          | option    |
    //   | example1.txt         | sub/folder1 | keep both |
    //   | folder1/example1.txt | sub/folder1 | replace   |
    await ui.userMovesResources({
      stepUser: 'Alice',
      actionType: fileAction.sideBarPanel,
      resources: [
        { resource: 'example1.txt', to: 'sub/folder1', option: 'keep both' },
        { resource: 'folder1/example1.txt', to: 'sub/folder1', option: 'replace' }
      ]
    })
    // copy and move folder
    // And "Alice" copies the following resource using sidebar-panel
    //   | resource | to  | option    |
    //   | folder1  | sub | keep both |
    //   issue https://github.com/owncloud/web/issues/10515
    //   | folder1  | sub | replace   |
    await ui.userCopiesResources({
      stepUser: 'Alice',
      actionType: fileAction.sideBarPanel,
      resources: [
        { resource: 'folder1', to: 'sub', option: 'keep both' },
        { resource: 'folder1', to: 'sub', option: 'replace' }
      ]
    })
    // And "Alice" moves the following resource using sidebar-panel
    //   | resource     | to  | option    |
    //   | folder1      | sub | keep both |
    //   issue https://github.com/owncloud/web/issues/10515
    //   | sub1/folder1 | sub | replace   |
    await ui.userMovesResources({
      stepUser: 'Alice',
      actionType: fileAction.sideBarPanel,
      resources: [
        { resource: 'folder1', to: 'sub', option: 'keep both' },
        { resource: 'sub1/folder1', to: 'sub', option: 'replace' }
      ]
    })
    // And "Alice" logs out
    await ui.userLogsOut({ stepUser: 'Alice' })
  })
})
