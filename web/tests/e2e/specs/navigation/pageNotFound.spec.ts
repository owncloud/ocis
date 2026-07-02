import { test } from '../../environment/test'
import * as api from '../../steps/api/api.js'
import * as ui from '../../steps/ui/index'

test.describe('Page not found', { tag: '@predefined-users' }, () => {
  test('not found page', async () => {
    // Given "Admin" creates following user using API
    //   | id    |
    //   | Alice |
    await api.usersHaveBeenCreated({ stepUser: 'Admin', users: ['Alice'] })

    // And "Alice" logs in
    await ui.userLogsIn({ stepUser: 'Alice' })

    // When "Alice" navigates to a non-existing page
    await ui.userNavigatesToNonExistingPage({ stepUser: 'Alice' })
    // Then "Alice" should see the not found page
    await ui.userShouldSeeNotFoundPage({ stepUser: 'Alice' })
    // And "Alice" logs out
    await ui.userLogsOut({ stepUser: 'Alice' })
  })
})
