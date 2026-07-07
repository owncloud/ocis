import { test } from '../../environment/test'
import * as ui from '../../steps/ui/index'
import * as api from '../../steps/api/api'

test.describe('details', () => {
  test('access token renewal via iframe', async () => {
    // Given "Admin" creates following users using API
    //   | id    |
    //   | Alice |
    await api.usersHaveBeenCreated({ stepUser: 'Admin', users: ['Alice'] })

    // And "Admin" assigns following roles to the users using API
    //   | id    | role        |
    //   | Alice | Space Admin |
    await api.userHasAssignedRolesToUsers({
      stepUser: 'Admin',
      users: [{ id: 'Alice', role: 'Space Admin' }]
    })

    // And "Alice" logs in
    await ui.userLogsIn({ stepUser: 'Alice' })

    // And "Alice" opens the "files" app
    await ui.userOpensApplication({ stepUser: 'Alice', name: 'files' })

    // And "Alice" navigates to the projects space page
    await ui.userNavigatesToSpacesPage({ stepUser: 'Alice' })

    // And "Alice" creates the following project spaces
    //   | name | id     |
    //   | team | team.1 |
    await ui.userCreatesProjectSpaces({
      stepUser: 'Alice',
      spaces: [{ name: 'team', id: 'team.1' }]
    })

    // When "Alice" waits for token renewal via refresh token
    await ui.userWaitsForTokenRenewal({ stepUser: 'Alice', renewalType: 'refresh token' })

    // And "Alice" navigates to the project space "team.1"
    await ui.userNavigatesToSpace({ stepUser: 'Alice', space: 'team.1' })

    // And "Alice" creates the following resources
    //   | resource     | type   |
    //   | space-folder | folder |
    await ui.userCreatesResources({
      stepUser: 'Alice',
      resources: [{ name: 'space-folder', type: 'folder' }]
    })

    // Then following resources should be displayed in the files list for user "Alice"
    //   | resource     |
    //   | space-folder |
    await ui.userShouldSeeResources({
      listType: 'files list',
      stepUser: 'Alice',
      resources: ['space-folder']
    })

    // When "Alice" navigates to new tab
    await ui.userNavigatesToNewTab({ stepUser: 'Alice' })

    // And "Alice" waits for token to expire
    await ui.userWaitsForTokenToExpire({ stepUser: 'Alice' })

    // And "Alice" closes the current tab
    await ui.userClosesTheCurrentTab({ stepUser: 'Alice' })

    // And "Alice" opens the "files" app
    await ui.userOpensApplication({ stepUser: 'Alice', name: 'files' })

    // And "Alice" creates the following resources
    //  | resource          | type    | content   |
    //  | PARENT/parent.txt | txtFile | some text |
    await ui.userCreatesResources({
      stepUser: 'Alice',
      resources: [{ name: 'PARENT/parent.txt', type: 'txtFile', content: 'some text' }]
    })

    // And "Alice" logs out
    await ui.userLogsOut({ stepUser: 'Alice' })
  })
})
