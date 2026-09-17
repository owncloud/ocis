import { test } from '../../environment/test'
import * as api from '../../steps/api/api'
import * as ui from '../../steps/ui/index'

test.describe('range selection with shift-click', { tag: '@predefined-users' }, () => {
  test.beforeEach(async () => {
    // Given "Admin" creates following user using API
    //   | id    |
    //   | Alice |
    await api.usersHaveBeenCreated({ stepUser: 'Admin', users: ['Alice'] })
    // And "Alice" creates the following folders in personal space using API
    //   | name    |
    //   | folder1 |
    //   | folder2 |
    //   | folder3 |
    //   | folder4 |
    await api.userHasCreatedFolders({
      stepUser: 'Alice',
      folderNames: ['folder1', 'folder2', 'folder3', 'folder4']
    })
    // And "Alice" logs in
    await ui.userLogsIn({ stepUser: 'Alice' })
    // And "Alice" navigates to the personal space page
    await ui.userNavigatesToPersonalSpacePage({ stepUser: 'Alice' })
  })

  test.afterEach(async () => {
    // And "Alice" logs out
    await ui.userLogsOut({ stepUser: 'Alice' })
  })

  test('shift-click selects every resource in between in the list view', async () => {
    // When "Alice" selects resource "folder1"
    await ui.userSelectsResource({ stepUser: 'Alice', resource: 'folder1' })
    // And "Alice" shift-clicks resource "folder3"
    await ui.userSelectsResource({
      stepUser: 'Alice',
      resource: 'folder3',
      modifiers: ['Shift']
    })
    // Then following resources should be selected for user "Alice"
    //   | resource |
    //   | folder1  |
    //   | folder2  |
    //   | folder3  |
    await ui.userShouldSeeSelectedResources({
      stepUser: 'Alice',
      resources: ['folder1', 'folder2', 'folder3']
    })
    // And following resources should not be selected for user "Alice"
    //   | resource |
    //   | folder4  |
    await ui.userShouldNotSeeSelectedResources({ stepUser: 'Alice', resources: ['folder4'] })
    // And "Alice" should not see any highlighted text
    await ui.userShouldNotSeeHighlightedText({ stepUser: 'Alice' })

    // When "Alice" shift-clicks resource "folder4" with ctrl/cmd pressed
    await ui.userSelectsResource({
      stepUser: 'Alice',
      resource: 'folder4',
      modifiers: ['Shift', 'ControlOrMeta']
    })
    // Then following resources should be selected for user "Alice"
    //   | resource |
    //   | folder1  |
    //   | folder2  |
    //   | folder3  |
    //   | folder4  |
    await ui.userShouldSeeSelectedResources({
      stepUser: 'Alice',
      resources: ['folder1', 'folder2', 'folder3', 'folder4']
    })
  })

  test('shift-click selects backwards and replaces the previous selection', async () => {
    // When "Alice" selects resource "folder4"
    await ui.userSelectsResource({ stepUser: 'Alice', resource: 'folder4' })
    // And "Alice" shift-clicks resource "folder2"
    await ui.userSelectsResource({
      stepUser: 'Alice',
      resource: 'folder2',
      modifiers: ['Shift']
    })
    // Then following resources should be selected for user "Alice"
    //   | resource |
    //   | folder2  |
    //   | folder3  |
    //   | folder4  |
    await ui.userShouldSeeSelectedResources({
      stepUser: 'Alice',
      resources: ['folder2', 'folder3', 'folder4']
    })
    // And following resources should not be selected for user "Alice"
    //   | resource |
    //   | folder1  |
    await ui.userShouldNotSeeSelectedResources({ stepUser: 'Alice', resources: ['folder1'] })

    // When "Alice" selects resource "folder1"
    await ui.userSelectsResource({ stepUser: 'Alice', resource: 'folder1' })
    // And "Alice" shift-clicks resource "folder2"
    await ui.userSelectsResource({
      stepUser: 'Alice',
      resource: 'folder2',
      modifiers: ['Shift']
    })
    // Then following resources should be selected for user "Alice"
    //   | resource |
    //   | folder1  |
    //   | folder2  |
    await ui.userShouldSeeSelectedResources({
      stepUser: 'Alice',
      resources: ['folder1', 'folder2']
    })
    // And following resources should not be selected for user "Alice"
    //   | resource |
    //   | folder3  |
    //   | folder4  |
    await ui.userShouldNotSeeSelectedResources({
      stepUser: 'Alice',
      resources: ['folder3', 'folder4']
    })
  })

  test('shift-click selects every resource in between in the tiles view', async () => {
    // When "Alice" switches to the tiles-view
    await ui.userSwitchesToTilesViewMode({ stepUser: 'Alice' })
    // And "Alice" selects resource "folder1"
    await ui.userSelectsResource({ stepUser: 'Alice', resource: 'folder1' })
    // And "Alice" shift-clicks resource "folder3"
    await ui.userSelectsResource({
      stepUser: 'Alice',
      resource: 'folder3',
      modifiers: ['Shift']
    })
    // Then following resources should be selected for user "Alice"
    //   | resource |
    //   | folder1  |
    //   | folder2  |
    //   | folder3  |
    await ui.userShouldSeeSelectedResources({
      stepUser: 'Alice',
      resources: ['folder1', 'folder2', 'folder3']
    })
    // And following resources should not be selected for user "Alice"
    //   | resource |
    //   | folder4  |
    await ui.userShouldNotSeeSelectedResources({ stepUser: 'Alice', resources: ['folder4'] })
    // And "Alice" should not see any highlighted text
    await ui.userShouldNotSeeHighlightedText({ stepUser: 'Alice' })
  })
})
