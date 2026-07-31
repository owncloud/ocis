import { test } from '../../environment/test'
import * as ui from '../../steps/ui/index'

test.describe('general management', () => {
  test('mfa', async () => {
    await ui.userLogsIn({ stepUser: 'Admin' })
    await ui.userOpensApplication({ stepUser: 'Admin', name: 'admin-settings' })
    await ui.userAuthenticatesWithOTP({ stepUser: 'Admin', deviceName: 'test' })
    await ui.userNavigatesToProjectSpaceManagementPage({ stepUser: 'Admin' })
    await ui.userLogsOut({ stepUser: 'Admin' })
    await ui.userLogsIn({ stepUser: 'Admin' })
    await ui.userOpensApplication({ stepUser: 'Admin', name: 'admin-settings' })
    await ui.logInWithOTP({ stepUser: 'Admin' })
    await ui.userNavigatesToProjectSpaceManagementPage({ stepUser: 'Admin' })
  })
})
