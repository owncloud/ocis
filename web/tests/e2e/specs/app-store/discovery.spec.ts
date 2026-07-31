import { test } from '../../environment/test'
import * as ui from '../../steps/ui/index'

test.describe('details', { tag: '@predefined-users' }, () => {
  test('apps can be searched and downloaded', async () => {
    // When "Admin" logs in
    await ui.userLogsIn({ stepUser: 'Admin' })

    // And "Admin" navigates to the app store
    await ui.userOpensAppStore({ stepUser: 'Admin' })

    // Then "Admin" should see the app store
    await ui.userShouldSeeAppStore({ stepUser: 'Admin' })

    // And "Admin" should see the following apps
    //   | app         |
    //   | Draw.io     |
    //   | JSON Viewer |
    //   | Unzip       |
    await ui.userShouldSeeApps({
      stepUser: 'Admin',
      expectedApps: ['Draw.io', 'JSON Viewer', 'Unzip']
    })

    // When "Admin" enters the search term "draw"
    await ui.userSetsSearchTerm({ stepUser: 'Admin', searchTerm: 'draw' })

    // Then "Admin" should see the following apps
    //   | app     |
    //   | Draw.io |
    await ui.userShouldSeeApps({ stepUser: 'Admin', expectedApps: ['Draw.io'] })

    // When "Admin" clicks on the tag "viewer" of the app "Draw.io"
    await ui.userSelectsAppTag({ stepUser: 'Admin', tag: 'viewer', app: 'Draw.io' })

    // Then "Admin" should see the following apps
    //   | app         |
    //   | JSON Viewer |
    //   | Draw.io     |
    await ui.userShouldSeeApps({ stepUser: 'Admin', expectedApps: ['Draw.io', 'JSON Viewer'] })

    // When "Admin" clicks on the app "JSON Viewer"
    await ui.userSelectsApp({ stepUser: 'Admin', app: 'JSON Viewer' })

    // Then "Admin" should see the app details of "JSON Viewer"
    await ui.userShouldSeeAppDetails({ stepUser: 'Admin', app: 'JSON Viewer' })

    // When "Admin" clicks on the tag "viewer"
    await ui.userSelectsTag({ stepUser: 'Admin', tag: 'viewer' })

    // Then "Admin" should see the app store
    await ui.userShouldSeeAppStore({ stepUser: 'Admin' })

    // Then "Admin" should see the following apps
    //   | app         |
    //   | JSON Viewer |
    //   | Draw.io     |
    await ui.userShouldSeeApps({ stepUser: 'Admin', expectedApps: ['Draw.io', 'JSON Viewer'] })

    // And "Admin" logs out
    await ui.userLogsOut({ stepUser: 'Admin' })
  })
})
