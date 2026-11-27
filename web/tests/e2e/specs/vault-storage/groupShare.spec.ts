import { fileAction, resourcePage } from '../../environment/constants'
import { test } from '../../environment/test'
import * as api from '../../steps/api/api'
import * as ui from '../../steps/ui/index'

const entersVault = async (stepUser: string): Promise<void> => {
  await ui.userSwitchesToVaultMode({ stepUser })
  await ui.userIsRedirectedToAuthenticatorPage({ stepUser })
  await ui.userAuthenticatesWithOTP({ stepUser, deviceName: 'test' })
  await ui.userIsInVaultMode({ stepUser })
}

test.describe('Vault group shares', { tag: '@predefined-users' }, () => {
  test.beforeEach(async () => {
    await api.usersHaveBeenCreated({ stepUser: 'Admin', users: ['Alice', 'Brian'] })
    await api.userHasAssignedRolesToUsers({
      stepUser: 'Admin',
      users: [
        { id: 'Alice', role: 'Admin' },
        { id: 'Brian', role: 'Space Admin' }
      ]
    })
    await api.groupsHaveBeenCreated({ stepUser: 'Admin', groupIds: ['sales'] })
    await api.usersHaveBeenAddedToGroup({
      stepUser: 'Admin',
      usersToAdd: [{ user: 'Brian', group: 'sales' }]
    })

    await ui.userLogsIn({ stepUser: 'Alice' })
    await entersVault('Alice')
    await ui.userCreatesResources({
      stepUser: 'Alice',
      resources: [{ name: 'vault-file.txt', type: 'txtFile', content: 'some content' }]
    })
  })

  test('Vault share dialog shows the group notice only for group recipients', async () => {
    // Brian's new role reaches oCIS only on login; without it he is not offered as a vault recipient
    await ui.userLogsIn({ stepUser: 'Brian' })

    await ui.userSelectsShareRecipient({
      stepUser: 'Alice',
      resource: 'vault-file.txt',
      recipient: 'Brian',
      type: 'user'
    })
    await ui.userShouldNotSeeVaultGroupNotice({ stepUser: 'Alice' })

    await ui.userReloadsPage({ stepUser: 'Alice' })
    await ui.userSelectsShareRecipient({
      stepUser: 'Alice',
      resource: 'vault-file.txt',
      recipient: 'sales',
      type: 'group'
    })
    await ui.userShouldSeeVaultGroupNotice({ stepUser: 'Alice' })
  })

  test('Member with vault permission opens a group share from Shared with me', async () => {
    await ui.userLogsIn({ stepUser: 'Brian' })
    await entersVault('Brian')

    await ui.userSharesResources({
      stepUser: 'Alice',
      actionType: fileAction.sideBarPanel,
      shares: [
        {
          resource: 'vault-file.txt',
          recipient: 'sales',
          type: 'group',
          role: 'Can view',
          resourceType: 'file'
        }
      ]
    })

    await ui.userShouldSeeNotifications({
      stepUser: 'Brian',
      expectedMessages: ['%user_alice_displayName% shared vault-file.txt with you']
    })
    await ui.userNavigatesToSharedWithMePage({ stepUser: 'Brian' })
    await ui.userShouldSeeResources({
      listType: resourcePage.shares,
      stepUser: 'Brian',
      resources: ['vault-file.txt']
    })
  })
})
