import { test } from '../../environment/test'
import * as ui from '../../steps/ui/index'
import * as api from '../../steps/api/api'

test.describe('Access breadcrumb', { tag: '@predefined-users' }, () => {
  test.beforeEach(async () => {
    await api.usersHaveBeenCreated({ stepUser: 'Admin', users: ['Alice'] })
    await ui.userLogsIn({ stepUser: 'Alice' })
  })

  test('Breadcrumb navigation', async () => {
    await ui.userCreatesResources({
      stepUser: 'Alice',
      resources: [{ name: 'parent/folder%2Fwith%2FSlashes', type: 'folder' }]
    })
    await ui.userOpensResource({ stepUser: 'Alice', resource: 'parent/folder%2Fwith%2FSlashes' })
    await ui.userCreatesResources({
      stepUser: 'Alice',
      resources: [{ name: `'single-double quotes"`, type: 'folder' }]
    })
    await ui.userOpensResource({ stepUser: 'Alice', resource: `'single-double quotes"` })
    await ui.userCreatesResources({
      stepUser: 'Alice',
      resources: [{ name: `"inner" double quote`, type: 'folder' }]
    })
    await ui.userOpensResource({ stepUser: 'Alice', resource: `"inner" double quote` })
    await ui.userCreatesResources({
      stepUser: 'Alice',
      resources: [{ name: 'sub-folder', type: 'folder' }]
    })
    await ui.userOpensResource({ stepUser: 'Alice', resource: 'sub-folder' })
    await ui.userNavigatesToFolderViaBreadcrumb({
      stepUser: 'Alice',
      resource: `"inner" double quote`
    })
    await ui.userNavigatesToFolderViaBreadcrumb({
      stepUser: 'Alice',
      resource: `'single-double quotes"`
    })
    await ui.userNavigatesToFolderViaBreadcrumb({
      stepUser: 'Alice',
      resource: 'folder%2Fwith%2FSlashes'
    })
    await ui.userNavigatesToFolderViaBreadcrumb({ stepUser: 'Alice', resource: 'parent' })
    await ui.userLogsOut({ stepUser: 'Alice' })
  })
})
