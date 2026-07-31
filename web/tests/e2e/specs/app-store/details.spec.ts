import { test } from '../../environment/test'
import * as ui from '../../steps/ui/index'

test.describe('details', { tag: '@predefined-users' }, () => {
  test('Apps can be viewed and downloaded', async () => {
    // When "Admin" logs in
    await ui.userLogsIn({ stepUser: 'Admin' })

    // And "Admin" navigates to the app store
    await ui.userOpensAppStore({ stepUser: 'Admin' })

    // Then "Admin" should see the app store
    await ui.userShouldSeeAppStore({ stepUser: 'Admin' })

    // When "Admin" clicks on the app "Development boilerplate"
    await ui.userSelectsApp({ stepUser: 'Admin', app: 'Development boilerplate' })

    // Then "Admin" should see the app details of "Development boilerplate"
    await ui.userShouldSeeAppDetails({ stepUser: 'Admin', app: 'Development boilerplate' })

    // And "Admin" downloads app version "0.1.0"
    await ui.userDownloadsAppVersion({ stepUser: 'Admin', version: '0.1.0' })

    // When "Admin" navigates back to the app store overview
    await ui.userNavigatesToAppStoreOverview({ stepUser: 'Admin' })

    // Then "Admin" should see the app store
    await ui.userShouldSeeAppStore({ stepUser: 'Admin' })

    // And "Admin" downloads the latest version of the app "Development boilerplate"
    await ui.userDownloadsApp({ stepUser: 'Admin', app: 'Development boilerplate' })

    // And "Admin" logs out
    await ui.userLogsOut({ stepUser: 'Admin' })
  })
})
