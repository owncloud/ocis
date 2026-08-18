import { test } from '../../environment/test'
import * as api from '../../steps/api/api'
import * as ui from '../../steps/ui/index'
import { application, fileAction } from '../../environment/constants'

const suites = [
  {
    name: 'Collabora',
    viewer: application.collabora,
    type: 'OpenDocument',
    file: 'readonly-collabora.odt'
  },
  {
    name: 'OnlyOffice',
    viewer: application.onlyOffice,
    type: 'Microsoft Word',
    file: 'readonly-onlyoffice.docx'
  }
] as const

test.describe('WOPI CheckFileInfo: ReadOnly', { tag: '@predefined-users' }, () => {
  test.beforeEach(async () => {
    await api.usersHaveBeenCreated({ stepUser: 'Admin', users: ['Alice', 'Brian'] })
    await ui.userLogsIn({ stepUser: 'Alice' })
    await ui.userOpensApplication({ stepUser: 'Alice', name: 'files' })
  })

  for (const suite of suites) {
    test(`ReadOnly is false for the file owner in ${suite.name}`, async () => {
      await ui.userCreatesResources({
        stepUser: 'Alice',
        resources: [{ name: suite.file, type: suite.type, content: 'owner content' }]
      })
      await ui.userOpensResourceInViewer({
        stepUser: 'Alice',
        resource: suite.file,
        viewer: suite.viewer
      })
      await ui.userEditsResources({
        stepUser: 'Alice',
        resources: [{ name: suite.file, type: suite.type, content: 'owner edited content' }]
      })
      await ui.userShouldSeeContentInEditor({
        stepUser: 'Alice',
        expectedContent: 'owner edited content',
        editor: suite.name
      })
      await ui.userClosesFileViewer({ stepUser: 'Alice' })
    })

    test(`ReadOnly is true for a view-only share recipient in ${suite.name}`, async () => {
      await ui.userCreatesResources({
        stepUser: 'Alice',
        resources: [{ name: suite.file, type: suite.type, content: 'owner content' }]
      })
      await ui.userSharesResources({
        stepUser: 'Alice',
        actionType: fileAction.sideBarPanel,
        shares: [
          {
            resource: suite.file,
            recipient: 'Brian',
            type: 'user',
            role: 'Can view',
            resourceType: 'file'
          }
        ]
      })
      await ui.userLogsIn({ stepUser: 'Brian' })
      await ui.userNavigatesToSharedWithMePage({ stepUser: 'Brian' })
      await ui.userOpensResourceInViewer({
        stepUser: 'Brian',
        resource: suite.file,
        viewer: suite.viewer
      })
      await ui.userShouldNotBeAbleToEditContentOfResources({
        stepUser: 'Brian',
        resources: [{ name: suite.file, type: suite.type, content: 'owner content' }]
      })
      await ui.userClosesFileViewer({ stepUser: 'Brian' })
    })
  }
})
