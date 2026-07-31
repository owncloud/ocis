import { test } from '../../environment/test'
import * as api from '../../steps/api/api'
import * as ui from '../../steps/ui/index'
import { fileAction, resourcePage } from '../../environment/constants'

test.describe('deny share access', () => {
  test.beforeEach(async () => {
    await api.usersHaveBeenCreated({ stepUser: 'Admin', users: ['Alice', 'Brian'] })

    await ui.userLogsIn({ stepUser: 'Alice' })
    await ui.userLogsIn({ stepUser: 'Brian' })

    await api.userHasCreatedFolders({
      stepUser: 'Alice',
      folderNames: [
        'folder_to_shared',
        'folder_to_shared/folder',
        'folder_to_shared/folder_to_deny'
      ]
    })

    await ui.userOpensApplication({ stepUser: 'Alice', name: 'files' })

    await ui.userSharesResources({
      stepUser: 'Alice',
      actionType: fileAction.sideBarPanel,
      shares: [
        {
          resource: 'folder_to_shared',
          recipient: 'Brian',
          type: 'user',
          role: 'Can view',
          resourceType: 'folder'
        }
      ]
    })

    await ui.userOpensResource({ stepUser: 'Alice', resource: 'folder_to_shared' })

    await ui.userSharesResources({
      stepUser: 'Alice',
      actionType: fileAction.sideBarPanel,
      shares: [
        {
          resource: 'folder_to_deny',
          recipient: 'Brian',
          type: 'user',
          role: 'Cannot access',
          resourceType: 'folder'
        }
      ]
    })
  })

  test('deny and grant access', async () => {
    // deny access
    await ui.userOpensApplication({ stepUser: 'Brian', name: 'files' })

    await ui.userNavigatesToSharedWithMePage({ stepUser: 'Brian' })
    await ui.userOpensResource({ stepUser: 'Brian', resource: 'folder_to_shared' })

    await ui.userShouldNotSeeTheResources({
      listType: resourcePage.filesList,
      stepUser: 'Brian',
      resources: ['folder_to_deny']
    })

    await ui.userOpensApplication({ stepUser: 'Alice', name: 'files' })

    await ui.userOpensResource({ stepUser: 'Alice', resource: 'folder_to_shared' })

    // allow access - deleting "Cannot access" share

    await ui.userRemovesSharees({
      stepUser: 'Alice',
      sharees: [
        {
          resource: 'folder_to_deny',
          recipient: 'Brian'
        }
      ]
    })
    await ui.userOpensApplication({ stepUser: 'Brian', name: 'files' })
    await ui.userNavigatesToSharedWithMePage({ stepUser: 'Brian' })
    await ui.userOpensResource({ stepUser: 'Brian', resource: 'folder_to_shared' })
    await ui.userShouldSeeResources({
      listType: resourcePage.filesList,
      stepUser: 'Brian',
      resources: ['folder_to_deny']
    })

    await ui.userLogsOut({ stepUser: 'Brian' })
    await ui.userLogsOut({ stepUser: 'Alice' })
  })
})
