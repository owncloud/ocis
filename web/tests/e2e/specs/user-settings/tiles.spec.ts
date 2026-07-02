import { test } from '../../environment/test'
import * as api from '../../steps/api/api'
import * as ui from '../../steps/ui/index'

test.describe('tiles view', { tag: '@predefined-users' }, () => {
  test.beforeEach(async () => {
    // Given "Admin" creates following user using API
    await api.usersHaveBeenCreated({ stepUser: 'Admin', users: ['Alice'] })
    // And "Alice" logs in
    await ui.userLogsIn({ stepUser: 'Alice' })
    // And "Alice" creates the following resources
    await api.userHasCreatedFolder({ stepUser: 'Alice', folderName: 'tile_folder' })
  })

  test('Users can navigate web via tiles', async () => {
    // When "Alice" switches to the tiles-view
    await ui.userSwitchesToTilesViewMode({ stepUser: 'Alice' })
    // Then "Alice" sees the resources displayed as tiles
    await ui.userShouldSeeResourcesAsTiles({ stepUser: 'Alice' })
    // And "Alice" opens folder "tile_folder"
    await ui.userOpensResource({ stepUser: 'Alice', resource: 'tile_folder' })
    // And "Alice" creates the following resources
    await ui.userCreatesResources({
      stepUser: 'Alice',
      resources: [{ name: 'tile_folder/tile_folder2', type: 'folder' }]
    })
    // And "Alice" sees the resources displayed as tiles
    await ui.userShouldSeeResourcesAsTiles({ stepUser: 'Alice' })
    // And "Alice" logs out
    await ui.userLogsOut({ stepUser: 'Alice' })
  })
})
