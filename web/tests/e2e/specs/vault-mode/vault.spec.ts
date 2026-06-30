import { resourcePage } from '../../environment/constants'
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

  test.describe('Resource Isolation', () => {
    test.beforeEach(async () => {
      await api.userHasCreatedFiles({
        stepUser: 'Alice',
        files: [{ pathToFile: 'drive-file.txt', content: 'drive content' }]
      })
      
      // await api.userHasCreatedVaultFiles({
      //   stepUser: 'Alice',
      //   files: [{ pathToFile: 'vault-file.txt', content: 'vault content' }]
      // })
    })

    test('Drive and Vault resources are isolated', async () => {
      await ui.userIsInDriveMode({
        stepUser: 'Alice'
      })

      await ui.userShouldSeeResources({
        listType: resourcePage.filesList,
        stepUser: 'Alice',
        resources: ['drive-file.txt']
      })

      await ui.userShouldNotSeeTheResources({
        listType: resourcePage.filesList,
        stepUser: 'Alice',
        resources: ['vault-file.txt']
      })

      await ui.userSwitchesToVaultMode({
        stepUser: 'Alice'
      })

      await ui.userAuthenticatesToVault({
        stepUser: 'Alice'
      })

      await ui.userShouldSeeResources({
        listType: resourcePage.filesList,
        stepUser: 'Alice',
        resources: ['vault-file.txt']
      })

      await ui.userShouldNotSeeTheResources({
        listType: resourcePage.filesList,
        stepUser: 'Alice',
        resources: ['drive-file.txt']
      })
    }) 
  })
})