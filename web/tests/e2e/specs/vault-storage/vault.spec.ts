import { test } from '../../environment/test'
import * as api from '../../steps/api/api'
import * as ui from '../../steps/ui/index'

test.describe('Vault Mode', { tag: '@predefined-users' }, () => {
  test.beforeEach(async () => {
    await api.usersHaveBeenCreated({
      stepUser: 'Admin',
      users: ['Alice']
    })
    await api.userHasAssignedRolesToUsers({
      stepUser: 'Admin',
      users: [{ id: 'Alice', role: 'Space Admin' }]
    })
    await ui.userLogsIn({ stepUser: 'Alice' })
  })

  test('Switching between Drive and Vault modes', async () => {
    await ui.userIsInDriveMode({ stepUser: 'Alice' })

    await ui.userCreatesResources({
      stepUser: 'Alice',
      resources: [{ name: 'test.txt', type: 'txtFile', content: 'some content' }]
    })
    await ui.userShouldSeeResources({
      listType: 'files list',
      stepUser: 'Alice',
      resources: ['test.txt']
    })

    await ui.userSwitchesToVaultMode({ stepUser: 'Alice' })
    await ui.userIsRedirectedToAuthenticatorPage({ stepUser: 'Alice' })
    await ui.userAuthenticatesWithOTP({ stepUser: 'Alice', deviceName: 'test' })
    await ui.userIsInVaultMode({ stepUser: 'Alice' })

    await ui.userCreatesResources({
      stepUser: 'Alice',
      resources: [{ name: 'lorem.txt', type: 'txtFile', content: 'some content' }]
    })
    await ui.userShouldNotSeeTheResources({
      listType: 'files list',
      stepUser: 'Alice',
      resources: ['test.txt']
    })
    await ui.userShouldSeeResources({
      listType: 'files list',
      stepUser: 'Alice',
      resources: ['lorem.txt']
    })

    await ui.userSwitchesToDriveMode({ stepUser: 'Alice' })
    await ui.userIsInDriveMode({ stepUser: 'Alice' })

    await ui.userShouldNotSeeTheResources({
      listType: 'files list',
      stepUser: 'Alice',
      resources: ['lorem.txt']
    })

    await ui.userLogsOut({ stepUser: 'Alice' })
  })
})
