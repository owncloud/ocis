import { test } from '../../environment/test'
import * as api from '../../steps/api/api.js'
import * as ui from '../../steps/ui/index'

// Regression test for https://github.com/owncloud/ocis/issues/11467.
//
// A Keycloak user who carries no ocis realm role produces a token the proxy
// role_mapping cannot match. Before the fix the proxy answered that login with a
// bare 500 and the web UI showed an access denied page whose "log in again" button
// returned to the same page, because the login itself had succeeded -- clearing the
// browser cookies was the only way out. That is the normal state for users
// federated into the realm from an external user directory.
//
// PROXY_ROLE_ASSIGNMENT_OIDC_DEFAULT_ROLE names the role to fall back to. It is
// empty by default; tests/acceptance/run-e2e.py sets it to "User" for the keycloak
// suite. A matching role_mapping entry still wins over it, which is why Alice --
// created the normal way, with the ocisUser realm role -- is exercised first in the
// same test.
test.describe('oidc default role', () => {
  test('user with no ocis role in the token can log in', async () => {
    // Given "Admin" creates following user using API
    //   | id    |
    //   | Alice |
    await api.usersHaveBeenCreated({ stepUser: 'Admin', users: ['Alice'] })

    // When "Alice" logs in
    // Alice has the ocisUser realm role, so she is assigned through the role
    // mapping and never touches the default role.
    await ui.userLogsIn({ stepUser: 'Alice' })
    // And "Alice" navigates to the personal space page
    await ui.userNavigatesToPersonalSpacePage({ stepUser: 'Alice' })
    // Then "Alice" should see an empty personal space
    await ui.userShouldSeeEmptyPersonalSpace({ stepUser: 'Alice' })
    // And "Alice" logs out
    await ui.userLogsOut({ stepUser: 'Alice' })

    // When "Admin" creates the following user in Keycloak without any ocis realm role
    //   | id    |
    //   | Brian |
    await api.keycloakUsersWithoutRoleHaveBeenCreated({ stepUser: 'Admin', users: ['Brian'] })

    // Then "Brian" can log in
    // This is the assertion the issue is about. Without the fallback this step
    // fails: the login succeeds at the IDP and the proxy then refuses to resolve
    // the account, so #web-content never appears.
    await ui.userLogsIn({ stepUser: 'Brian' })

    // And "Brian" navigates to the personal space page
    await ui.userNavigatesToPersonalSpacePage({ stepUser: 'Brian' })
    // And "Brian" should see an empty personal space
    await ui.userShouldSeeEmptyPersonalSpace({ stepUser: 'Brian' })

    // And "Brian" navigates to the project spaces page
    await ui.userNavigatesToSpacesPage({ stepUser: 'Brian' })
    // And "Brian" should see no project space
    await ui.userShouldSeeNoSpaces({ stepUser: 'Brian' })

    // And "Brian" logs out
    await ui.userLogsOut({ stepUser: 'Brian' })
  })
})
