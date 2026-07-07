import { test } from '../../environment/test'
import * as api from '../../steps/api/api'
import * as ui from '../../steps/ui/index'

test.describe('GDPR export', { tag: '@predefined-users' }, () => {
  test.beforeEach(async () => {
    await api.usersHaveBeenCreated({ stepUser: 'Admin', users: ['Alice'] })
    await ui.userLogsIn({ stepUser: 'Alice' })
  })

  test('create and download a GDPR export', async () => {
    // And "Alice" opens the user menu
    await ui.userOpensAccountPage({ stepUser: 'Alice' })
    // And "Alice" requests a new GDPR export
    await ui.userRequestsGdprExport({ stepUser: 'Alice' })
    // And "Alice" downloads the GDPR export
    await ui.userDownloadsGdprExport({ stepUser: 'Alice' })
    // And "Alice" logs out
    await ui.userLogsOut({ stepUser: 'Alice' })
  })
})
