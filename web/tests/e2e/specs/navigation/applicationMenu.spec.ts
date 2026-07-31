import { test } from '../../environment/test'
import * as ui from '../../steps/ui/index'
import * as api from '../../steps/api/api'
import { resourcePage } from '../../environment/constants'

test.describe('Application menu', { tag: '@predefined-users' }, () => {
  test.beforeEach(async () => {
    await api.usersHaveBeenCreated({ stepUser: 'Admin', users: ['Alice'] })
    await ui.userLogsIn({ stepUser: 'Alice' })
  })

  test('Open text editor via application menu', async () => {
    await ui.userOpensApplication({ stepUser: 'Alice', name: 'text-editor' })
    await ui.userAddsContentInTextEditor({ stepUser: 'Alice', text: 'Hello world' })
    await ui.userSavesTextEditor({ stepUser: 'Alice' })
    await ui.userClosesFileViewer({ stepUser: 'Alice' })
    await ui.userShouldSeeResources({
      listType: resourcePage.filesList,
      stepUser: 'Alice',
      resources: ['New file.txt']
    })
    await ui.userLogsOut({ stepUser: 'Alice' })
  })
})
