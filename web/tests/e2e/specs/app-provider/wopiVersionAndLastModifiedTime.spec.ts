import { test } from '../../environment/test'
import * as api from '../../steps/api/api'
import * as ui from '../../steps/ui/index'
import { application, fileAction } from '../../environment/constants'

const suites = [
  {
    name: 'Collabora',
    viewer: application.collabora,
    type: 'OpenDocument',
    file: 'version-collabora.odt'
  },
  {
    name: 'OnlyOffice',
    viewer: application.onlyOffice,
    type: 'Microsoft Word',
    file: 'version-onlyoffice.docx'
  }
] as const

test.describe(
  'WOPI CheckFileInfo: Version and LastModifiedTime',
  { tag: '@predefined-users' },
  () => {
    test.beforeEach(async () => {
      await api.usersHaveBeenCreated({ stepUser: 'Admin', users: ['Alice', 'Brian'] })
      await ui.userLogsIn({ stepUser: 'Alice' })
      await ui.userOpensApplication({ stepUser: 'Alice', name: 'files' })
    })

    for (const suite of suites) {
      test(`a fresh edit is visible to another session opening ${suite.name} afterward`, async () => {
        await ui.userCreatesResources({
          stepUser: 'Alice',
          resources: [{ name: suite.file, type: suite.type, content: 'initial content' }]
        })
        await ui.userSharesResources({
          stepUser: 'Alice',
          actionType: fileAction.sideBarPanel,
          shares: [
            {
              resource: suite.file,
              recipient: 'Brian',
              type: 'user',
              role: 'Can edit',
              resourceType: 'file'
            }
          ]
        })

        await ui.userOpensResourceInViewer({
          stepUser: 'Alice',
          resource: suite.file,
          viewer: suite.viewer
        })
        await ui.userEditsResources({
          stepUser: 'Alice',
          resources: [{ name: suite.file, type: suite.type, content: 'edited content v2' }]
        })
        await ui.userShouldSeeContentInEditor({
          stepUser: 'Alice',
          expectedContent: 'edited content v2',
          editor: suite.name
        })
        await ui.userClosesFileViewer({ stepUser: 'Alice' })

        await ui.userLogsIn({ stepUser: 'Brian' })
        await ui.userNavigatesToSharedWithMePage({ stepUser: 'Brian' })
        await ui.userOpensResourceInViewer({
          stepUser: 'Brian',
          resource: suite.file,
          viewer: suite.viewer
        })
        // if LastModifiedTime/Version were wrong or stale, the WOPI client could serve Brian a
        // cached pre-edit copy instead of fetching the actual current content
        await ui.userShouldSeeContentInEditor({
          stepUser: 'Brian',
          expectedContent: 'edited content v2',
          editor: suite.name
        })
        await ui.userClosesFileViewer({ stepUser: 'Brian' })
      })
    }
  }
)
