import { application, fileAction } from '../../environment/constants'
import { test } from '../../environment/test'
import * as api from '../../steps/api/api'
import * as ui from '../../steps/ui/index'

test.describe('Vault access enforcement', { tag: '@predefined-users' }, () => {
  test('User without vault permission is denied when manipulating the URL', async () => {
    await api.usersHaveBeenCreated({ stepUser: 'Admin', users: ['Alice'] })
    await api.userHasAssignedRolesToUsers({
      stepUser: 'Admin',
      users: [{ id: 'Alice', role: 'User' }]
    })
    await ui.userLogsIn({ stepUser: 'Alice' })

    await ui.userNavigatesToVaultViaUrl({ stepUser: 'Alice' })
    await ui.userIsOnAccessDeniedPage({ stepUser: 'Alice' })
  })

  test('User Light without vault permission is denied when manipulating the URL', async () => {
    await api.usersHaveBeenCreated({ stepUser: 'Admin', users: ['Alice'] })
    await api.userHasAssignedRolesToUsers({
      stepUser: 'Admin',
      users: [{ id: 'Alice', role: 'User Light' }]
    })
    await ui.userLogsIn({ stepUser: 'Alice' })

    await ui.userNavigatesToVaultViaUrl({ stepUser: 'Alice' })
    await ui.userIsOnAccessDeniedPage({ stepUser: 'Alice' })
  })

  test('Space Admin manipulating the URL is routed to MFA, not denied', async () => {
    await api.usersHaveBeenCreated({ stepUser: 'Admin', users: ['Alice'] })
    await api.userHasAssignedRolesToUsers({
      stepUser: 'Admin',
      users: [{ id: 'Alice', role: 'Space Admin' }]
    })
    await ui.userLogsIn({ stepUser: 'Alice' })

    await ui.userNavigatesToVaultViaUrl({ stepUser: 'Alice' })
    await ui.userIsRedirectedToAuthenticatorPage({ stepUser: 'Alice' })
  })

  test('Space Admin keeps vault context after a reload', async () => {
    await api.usersHaveBeenCreated({ stepUser: 'Admin', users: ['Alice'] })
    await api.userHasAssignedRolesToUsers({
      stepUser: 'Admin',
      users: [{ id: 'Alice', role: 'Space Admin' }]
    })
    await ui.userLogsIn({ stepUser: 'Alice' })

    await ui.userSwitchesToVaultMode({ stepUser: 'Alice' })
    await ui.userIsRedirectedToAuthenticatorPage({ stepUser: 'Alice' })
    await ui.userAuthenticatesWithOTP({ stepUser: 'Alice', deviceName: 'test' })
    await ui.userIsInVaultMode({ stepUser: 'Alice' })

    await ui.userReloadsPage({ stepUser: 'Alice' })
    await ui.userIsInVaultMode({ stepUser: 'Alice' })
  })

  test('Opening a permanent link copied in Vault mode opens the file successfully', async () => {
    await api.usersHaveBeenCreated({ stepUser: 'Admin', users: ['Alice', 'Brian'] })
    await api.userHasAssignedRolesToUsers({
      stepUser: 'Admin',
      users: [
        { id: 'Alice', role: 'Admin' },
        { id: 'Brian', role: 'Space Admin' }
      ]
    })
    await ui.userLogsIn({ stepUser: 'Alice' })
    await ui.userSwitchesToVaultMode({ stepUser: 'Alice' })
    await ui.userIsRedirectedToAuthenticatorPage({ stepUser: 'Alice' })
    await ui.userAuthenticatesWithOTP({ stepUser: 'Alice', deviceName: 'test' })
    await ui.userIsInVaultMode({ stepUser: 'Alice' })

    await ui.userCreatesResources({
      stepUser: 'Alice',
      resources: [{ name: 'test.txt', type: 'txtFile', content: 'some content' }]
    })

    await ui.userLogsIn({ stepUser: 'Brian' })
    await ui.userSwitchesToVaultMode({ stepUser: 'Brian' })
    await ui.userIsRedirectedToAuthenticatorPage({ stepUser: 'Brian' })
    await ui.userAuthenticatesWithOTP({ stepUser: 'Brian', deviceName: 'test' })
    await ui.userIsInVaultMode({ stepUser: 'Brian' })

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
    await ui.userNavigatesToSharedWithMePage({ stepUser: 'Brian' })

    const permanentLink = await ui.userCopiesResourcePermanentLink({
      stepUser: 'Brian',
      resource: 'test.txt'
    })
    await ui.userOpensPermanentLink({ stepUser: 'Brian', permanentLink })

    await ui.userShouldSeeContentInEditor({
      stepUser: 'Brian',
      expectedContent: 'some content',
      editor: application.textEditor
    })
  })
})
