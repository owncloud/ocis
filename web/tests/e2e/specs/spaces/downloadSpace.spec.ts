import { test } from '../../environment/test'
import * as api from '../../steps/api/api.js'
import * as ui from '../../steps/ui/index'

test.describe('download space', () => {
  test('download space', async () => {
    // Given "Admin" creates following users using API
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

    // Given "Alice" logs in
    await ui.userLogsIn({ stepUser: 'Alice' })

    // And "Alice" creates the following project spaces using API
    //   | name | id     |
    //   | team | team.1 |
    await api.userHasCreatedProjectSpaces({
      stepUser: 'Alice',
      spaces: [{ name: 'team', id: 'team.1' }]
    })

    // And "Alice" creates the following folder in space "team" using API
    //   | name        |
    //   | spaceFolder |
    await api.userHasCreatedFoldersInSpace({
      stepUser: 'Alice',
      spaceName: 'team',
      folders: ['spaceFolder']
    })

    // And "Alice" creates the following file in space "team" using API
    //   | name                  | content    |
    //   | spaceFolder/lorem.txt | space team |
    await api.userHasCreatedFilesInsideSpace({
      stepUser: 'Alice',
      files: [{ name: 'spaceFolder/lorem.txt', space: 'team', content: 'space team' }]
    })

    // And "Alice" navigates to the project space "team.1"
    await ui.userNavigatesToSpace({ stepUser: 'Alice', space: 'team.1' })

    // When "Alice" downloads the space "team.1"
    await ui.userDownloadsSpace({ stepUser: 'Alice' })

    // And "Alice" adds following users to the project space
    //   | user     | role     | kind  |
    //   | Brian    | Can edit | user  |
    await ui.userAddsMembersToSpace({
      stepUser: 'Alice',
      members: [{ user: 'Brian', role: 'Can edit with versions and trash bin', kind: 'user' }]
    })

    // And "Alice" logs out
    await ui.userLogsOut({ stepUser: 'Alice' })
    // And "Brian" logs in
    await ui.userLogsIn({ stepUser: 'Brian' })

    // And "Brian" navigates to the project space "team.1"
    await ui.userNavigatesToSpace({ stepUser: 'Brian', space: 'team.1' })

    // When "Alice" downloads the space "team.1"
    await ui.userDownloadsSpace({ stepUser: 'Brian' })
    // And "Brian" logs out
    await ui.userLogsOut({ stepUser: 'Brian' })
  })
})
