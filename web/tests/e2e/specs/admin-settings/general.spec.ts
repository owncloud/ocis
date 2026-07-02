import { test } from '../../environment/test'
import * as ui from '../../steps/ui/index'

test.describe('general management', () => {
  test('logo can be changed in the admin settings', async () => {
    await ui.userLogsIn({ stepUser: 'Admin' })
    await ui.userOpensApplication({ stepUser: 'Admin', name: 'admin-settings' })
    await ui.userNavigatesToGeneralManagementPage({ stepUser: 'Admin' })
    await ui.userUploadsLogoFromLocalPath({
      stepUser: 'Admin',
      localFile: 'filesForUpload/testavatar.png'
    })
    await ui.userNavigatesToGeneralManagementPage({ stepUser: 'Admin' })
    await ui.userResetsLogo({ stepUser: 'Admin' })
    await ui.userLogsOut({ stepUser: 'Admin' })
  })
})
