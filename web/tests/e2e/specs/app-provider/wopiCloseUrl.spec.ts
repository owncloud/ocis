import { expect } from '@playwright/test'
import { test } from '../../environment/test'
import * as api from '../../steps/api/api'
import * as ui from '../../steps/ui/index'
import { application } from '../../environment/constants'
import { getWorld } from '../../environment/world'

const suites = [
  {
    name: 'Collabora',
    viewer: application.collabora,
    type: 'OpenDocument',
    file: 'closeurl-collabora.odt'
  },
  {
    name: 'OnlyOffice',
    viewer: application.onlyOffice,
    type: 'Microsoft Word',
    file: 'closeurl-onlyoffice.docx'
  }
] as const

test.describe('WOPI CheckFileInfo: CloseUrl', { tag: '@predefined-users' }, () => {
  test.beforeEach(async () => {
    await api.usersHaveBeenCreated({ stepUser: 'Admin', users: ['Alice'] })
    await ui.userLogsIn({ stepUser: 'Alice' })
    await ui.userOpensApplication({ stepUser: 'Alice', name: 'files' })
  })

  for (const suite of suites) {
    test(`closing ${suite.name} navigates back into the files app via CloseUrl`, async () => {
      await ui.userCreatesResources({
        stepUser: 'Alice',
        resources: [{ name: suite.file, type: suite.type, content: 'close url content' }]
      })
      await ui.userOpensResourceInViewer({
        stepUser: 'Alice',
        resource: suite.file,
        viewer: suite.viewer
      })

      const world = getWorld()
      const { page } = world.actorsEnvironment.getActor({ key: 'Alice' })

      // userClosesFileViewer's underlying editor.close() already waits for
      // page.waitForURL(/.*\/files\/(spaces|shares|link|search)\/.*/) internally, so a broken
      // CloseUrl would already cause this call to time out and fail the test
      await ui.userClosesFileViewer({ stepUser: 'Alice' })

      await expect(page).toHaveURL(/.*\/files\/(spaces|shares|link|search)\/.*/)
    })
  }
})
