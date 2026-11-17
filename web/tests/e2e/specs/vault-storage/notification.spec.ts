import { fileAction } from '../../environment/constants'
import { test } from '../../environment/test'
import * as api from '../../steps/api/api'
import * as ui from '../../steps/ui/index'

test.describe('Vault Notifications', { tag: '@predefined-users' }, () => {
  test('Vault notifications are isolated from Drive notifications', async () => {
    await api.usersHaveBeenCreated({ stepUser: 'Admin', users: ['Alice', 'Brian'] })
    await api.userHasAssignedRolesToUsers({
      stepUser: 'Admin',
      users: [
        { id: 'Alice', role: 'Admin' },
        { id: 'Brian', role: 'Space Admin' }
      ]
    })

    await ui.userLogsIn({ stepUser: 'Alice' })
    await ui.userCreatesResources({
      stepUser: 'Alice',
      resources: [{ name: 'test.txt', type: 'txtFile', content: 'some content' }]
    })
    await ui.userSharesResources({
      stepUser: 'Alice',
      actionType: fileAction.sideBarPanel,
      shares: [
        {
          resource: 'test.txt',
          recipient: 'Brian',
          type: 'user',
          role: 'Can view',
          resourceType: 'file'
        }
      ]
    })
    await ui.userSwitchesToVaultMode({ stepUser: 'Alice' })
    await ui.userIsRedirectedToAuthenticatorPage({ stepUser: 'Alice' })
    await ui.userAuthenticatesWithOTP({ stepUser: 'Alice', deviceName: 'test' })
    await ui.userIsInVaultMode({ stepUser: 'Alice' })

    await ui.userCreatesResources({
      stepUser: 'Alice',
      resources: [{ name: 'lorem.txt', type: 'txtFile', content: 'some content' }]
    })

    await ui.userLogsIn({ stepUser: 'Brian' })

    await ui.userSharesResources({
      stepUser: 'Alice',
      actionType: fileAction.sideBarPanel,
      shares: [
        {
          resource: 'lorem.txt',
          recipient: 'Brian',
          type: 'user',
          role: 'Can view',
          resourceType: 'file'
        }
      ]
    })
    await ui.userLogsOut({ stepUser: 'Alice' })

    await ui.userShouldSeeNotifications({
      stepUser: 'Brian',
      expectedMessages: ['%user_alice_displayName% shared test.txt with you']
    })
    await ui.userShouldNotSeeNotifications({
      stepUser: 'Brian',
      unexpectedMessages: ['%user_alice_displayName% shared lorem.txt with you']
    })

    await ui.userSwitchesToVaultMode({ stepUser: 'Brian' })
    await ui.userIsRedirectedToAuthenticatorPage({ stepUser: 'Brian' })
    await ui.userAuthenticatesWithOTP({ stepUser: 'Brian', deviceName: 'test' })
    await ui.userIsInVaultMode({ stepUser: 'Brian' })

    await ui.userShouldSeeNotifications({
      stepUser: 'Brian',
      expectedMessages: ['%user_alice_displayName% shared lorem.txt with you']
    })
    await ui.userShouldNotSeeNotifications({
      stepUser: 'Brian',
      unexpectedMessages: ['%user_alice_displayName% shared test.txt with you']
    })
    await ui.userLogsOut({ stepUser: 'Brian' })
  })
})
