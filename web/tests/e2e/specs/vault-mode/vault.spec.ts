import { test } from '../../environment/test'
import * as api from '../../steps/api/api'
import * as ui from '../../steps/ui/index'

test.describe('Vault Mode Access and Authentication', { tag: '@predefined-users' }, () => {
  test.beforeEach(async () => {
    await api.usersHaveBeenCreated({
      stepUser: 'Admin',
      users: ['Alice']
    })
    await ui.userLogsIn({ stepUser: 'Alice' })
  })

  test('User with vault permission can access Vault mode with second factor authentication', async () => {
    await ui.userIsInDriveMode({
      stepUser: 'Alice'
    })

    await ui.userSwitchesToVaultMode({
      stepUser: 'Alice'
    })

    await ui.userIsRedirectedToAuthenticatorPage({
      stepUser: 'Alice'
    })

    await ui.userAuthenticatesToVault({
      stepUser: 'Alice'
    })

    await ui.userIsInVaultMode({
      stepUser: 'Alice'
    })

    await ui.userSwitchesToDriveMode({
      stepUser: 'Alice'
    })

    await ui.userIsInDriveMode({
      stepUser: 'Alice'
    })
  })
})
