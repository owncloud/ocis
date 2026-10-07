import { test } from '../../environment/test'
import * as api from '../../steps/api/api'
import * as ui from '../../steps/ui/index'
import { application } from '../../environment/constants'

const suites = [
  {
    name: 'Collabora',
    viewer: application.collabora,
    type: 'OpenDocument',
    file: 'supportsupdate-collabora.odt'
  },
  {
    name: 'OnlyOffice',
    viewer: application.onlyOffice,
    type: 'Microsoft Word',
    file: 'supportsupdate-onlyoffice.docx'
  }
] as const

test.describe('WOPI CheckFileInfo: SupportsUpdate', { tag: '@predefined-users' }, () => {
  test.beforeEach(async () => {
    await api.usersHaveBeenCreated({ stepUser: 'Admin', users: ['Alice'] })
    await ui.userLogsIn({ stepUser: 'Alice' })
    await ui.userOpensApplication({ stepUser: 'Alice', name: 'files' })
  })

  for (const suite of suites) {
    test(`SupportsUpdate: a new ${suite.name} document created via the office suite persists its content`, async () => {
      // creating a resource of this type/content via the office suite's "New" menu goes
      // through the WOPI PutRelativeFile flow, the specific operation SupportsUpdate vouches for
      await ui.userCreatesResources({
        stepUser: 'Alice',
        resources: [{ name: suite.file, type: suite.type, content: 'created via office suite' }]
      })
      await ui.userOpensResourceInViewer({
        stepUser: 'Alice',
        resource: suite.file,
        viewer: suite.viewer
      })
      await ui.userShouldSeeContentInEditor({
        stepUser: 'Alice',
        expectedContent: 'created via office suite',
        editor: suite.name
      })
      await ui.userClosesFileViewer({ stepUser: 'Alice' })

      // reopen to confirm the PutRelativeFile write actually landed server-side, not just
      // reflected in the editor's own in-memory state
      await ui.userOpensResourceInViewer({
        stepUser: 'Alice',
        resource: suite.file,
        viewer: suite.viewer
      })
      await ui.userShouldSeeContentInEditor({
        stepUser: 'Alice',
        expectedContent: 'created via office suite',
        editor: suite.name
      })
      await ui.userClosesFileViewer({ stepUser: 'Alice' })
    })
  }
})
